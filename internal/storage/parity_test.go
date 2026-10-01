package storage

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func newPostgresStore(t *testing.T) *PostgresStore {
	t.Helper()
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN not set; skipping Postgres parity")
	}
	s, err := NewPostgresStore(dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// snapshot is the observable state a scenario produces, so the two stores can
// be compared for parity without touching shared rows.
type snapshot struct {
	Status          string
	Streak          int
	BestStreak      int
	Attempts        int
	AvgResponseTime float64
	WeaknessScore   float64
	SM2EFactor      float64
	SM2Repetitions  int
	SM2Interval     int
	AttemptPoints   []bool
	DailyQuestions  int
	DailyCorrect    int
}

// exercise runs the same read/write sequence against either store and returns
// the resulting snapshot. A unique student keeps runs isolated.
func exercise(t *testing.T, repo Repository, tag string) snapshot {
	t.Helper()
	name := fmt.Sprintf("parity-%s-%d", tag, time.Now().UnixNano())
	stu, err := repo.CreateStudent(name)
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	sess, err := repo.CreateSession(stu.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	p := &ConceptProgress{
		StudentID: stu.ID, ConceptID: "c1", Status: "LEARNING",
		Streak: 2, BestStreak: 3, Attempts: 4, AvgResponseTime: 6.5,
		WeaknessScore: 0.3, SM2EFactor: 2.4, SM2Repetitions: 2, SM2Interval: 6,
	}
	if err := repo.UpsertProgress(p); err != nil {
		t.Fatalf("upsert progress: %v", err)
	}
	got, err := repo.GetProgress(stu.ID, "c1")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}

	base := time.Now().UTC().Add(-time.Hour)
	_ = repo.RecordAttempt(AttemptEntry{SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c1",
		Answer: "4 and 1/10", Expected: "4 1/10", Correct: false, ElapsedSeconds: 5, Timestamp: base})
	_ = repo.RecordAttempt(AttemptEntry{SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c1",
		Answer: "5", Expected: "4 1/10", Correct: false, ElapsedSeconds: 9, Timestamp: base.Add(time.Minute)})

	atts, err := repo.GetAttemptsForStudent(stu.ID)
	if err != nil {
		t.Fatalf("get attempts: %v", err)
	}
	// Flip the first attempt (postgres/sqlite must both expose a stable id).
	if len(atts) > 0 {
		if err := repo.UpdateAttemptCorrect(atts[0].ID, true); err != nil {
			t.Fatalf("update attempt: %v", err)
		}
	}
	atts, err = repo.GetAttemptsForStudent(stu.ID)
	if err != nil {
		t.Fatalf("reload attempts: %v", err)
	}
	points := make([]bool, len(atts))
	for i, a := range atts {
		points[i] = a.Correct
	}

	da, err := repo.GetDailyActivity(stu.ID, 5)
	if err != nil {
		t.Fatalf("daily activity: %v", err)
	}
	var questions, correct int
	for _, d := range da {
		questions += d.Questions
		correct += d.Correct
	}

	return snapshot{
		Status: got.Status, Streak: got.Streak, BestStreak: got.BestStreak,
		Attempts: got.Attempts, AvgResponseTime: got.AvgResponseTime,
		WeaknessScore: got.WeaknessScore, SM2EFactor: got.SM2EFactor,
		SM2Repetitions: got.SM2Repetitions, SM2Interval: got.SM2Interval,
		AttemptPoints: points, DailyQuestions: questions, DailyCorrect: correct,
	}
}

// TestDeleteAccount_RemovesEveryOwnedRow seeds one student across tables,
// deletes, and requires every owned row gone plus the address reusable.
// Runs against SQLite always and Postgres when TEST_POSTGRES_DSN is set.
func TestDeleteAccount_RemovesEveryOwnedRow(t *testing.T) {
	stores := map[string]Repository{"sq": newTestStore(t)}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		pg, err := NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { pg.Close() })
		stores["pg"] = pg
	}
	for tag, repo := range stores {
		stu, err := repo.CreateUser(fmt.Sprintf("del-%s", tag), fmt.Sprintf("deluser-%s-%d", tag, time.Now().UnixNano()), "hash")
		if err != nil {
			t.Fatalf("[%s] create user: %v", tag, err)
		}
		if err := repo.SetEmail(stu.ID, "gone@example.com"); err != nil {
			t.Fatalf("[%s] set email: %v", tag, err)
		}
		if err := repo.CreateIdentity("google", "gid-"+stu.ID, stu.ID, "gone@example.com", true); err != nil {
			t.Fatalf("[%s] create identity: %v", tag, err)
		}
		sess, err := repo.CreateSession(stu.ID)
		if err != nil {
			t.Fatalf("[%s] create session: %v", tag, err)
		}
		if err := repo.UpsertProgress(&ConceptProgress{StudentID: stu.ID, ConceptID: "c1", Status: "LEARNING", Streak: 1}); err != nil {
			t.Fatalf("[%s] upsert progress: %v", tag, err)
		}
		if err := repo.RecordAttempt(AttemptEntry{SessionID: sess.ID, StudentID: stu.ID, ConceptID: "c1", Answer: "1", Expected: "1", Correct: true, Timestamp: time.Now().UTC()}); err != nil {
			t.Fatalf("[%s] record attempt: %v", tag, err)
		}
		if err := repo.RecordQuizCompletion(stu.ID, 60); err != nil {
			t.Fatalf("[%s] quiz completion: %v", tag, err)
		}
		if _, err := repo.CreateReport(QuestionReport{ReporterID: stu.ID, ConceptID: "c1", Kind: "question", Question: "q", Status: "open", CreatedAt: time.Now().UTC()}); err != nil {
			t.Fatalf("[%s] create report: %v", tag, err)
		}
		if err := repo.DeleteAccount(stu.ID); err != nil {
			t.Fatalf("[%s] delete: %v", tag, err)
		}
		if got, err := repo.GetStudent(stu.ID); err != nil || got != nil {
			t.Errorf("[%s] student row survives: %+v %v", tag, got, err)
		}
		if got, err := repo.GetProgress(stu.ID, "c1"); err != nil || got != nil {
			t.Errorf("[%s] progress survives: %+v %v", tag, got, err)
		}
		if atts, err := repo.GetAttemptsForStudent(stu.ID); err != nil || len(atts) != 0 {
			t.Errorf("[%s] attempts survive: %d %v", tag, len(atts), err)
		}
		if reps, err := repo.ListReports("all", 10, 0); err != nil {
			t.Errorf("[%s] list reports: %v", tag, err)
		} else {
			for _, r := range reps {
				if r.ReporterID == stu.ID {
					t.Errorf("[%s] report survives", tag)
				}
			}
		}
		// Address reusable by a fresh account (one-email-one-account lives on).
		other, err := repo.CreateUser("Next", fmt.Sprintf("next-%s-%d", tag, time.Now().UnixNano()), "hash")
		if err != nil {
			t.Fatalf("[%s] create next user: %v", tag, err)
		}
		if err := repo.SetEmail(other.ID, "gone@example.com"); err != nil {
			t.Errorf("[%s] address not reusable: %v", tag, err)
		}
		// Retry-safe: second delete is a no-op success.
		if err := repo.DeleteAccount(stu.ID); err != nil {
			t.Errorf("[%s] repeat delete: %v", tag, err)
		}
	}
}
func TestParity_SQLiteVsPostgres(t *testing.T) {
	pg := newPostgresStore(t)
	if got := exercise(t, pg, "pg"); reflect.DeepEqual(got, snapshot{}) {
		t.Fatal("postgres produced an empty snapshot")
	}
	pgSnap := exercise(t, pg, "pg2")
	sqSnap := exercise(t, newTestStore(t), "sq")

	if !reflect.DeepEqual(pgSnap, sqSnap) {
		t.Errorf("storage parity mismatch\n postgres: %+v\n sqlite:   %+v", pgSnap, sqSnap)
	}
}
