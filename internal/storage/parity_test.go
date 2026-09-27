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

// TestParity_SQLiteVsPostgres runs one identical scenario against both stores
// and requires the same observable result. Skipped unless TEST_POSTGRES_DSN is set.
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
