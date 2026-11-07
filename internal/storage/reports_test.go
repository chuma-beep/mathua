package storage

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func seedAttempt(t *testing.T, repo Repository, sid, concept, source string, correct bool, at time.Time) {
	t.Helper()
	sess, err := repo.CreateSession(sid)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	_, err = repo.RecordAttempt(AttemptEntry{
		SessionID:      sess.ID,
		StudentID:      sid,
		ConceptID:      concept,
		Answer:         "answer",
		Expected:       "expected",
		Correct:        correct,
		ElapsedSeconds: 5,
		Timestamp:      at.UTC(),
		Question:       "question",
		Source:         source,
	})
	if err != nil {
		t.Fatalf("record attempt: %v", err)
	}
}

func TestAssessmentAttemptsFiltersAndOrders(t *testing.T) {
	repo := newTestStore(t)
	st, err := repo.CreateStudent("Learner")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	seedAttempt(t, repo, st.ID, "a", "practice", true, now.Add(-30*time.Minute))
	seedAttempt(t, repo, st.ID, "a", "quiz", false, now.Add(-20*time.Minute))
	seedAttempt(t, repo, st.ID, "a", "diagnostic", true, now.Add(-time.Hour))
	seedAttempt(t, repo, st.ID, "a", "quiz", true, now.Add(-3*time.Hour))

	got, err := repo.AssessmentAttempts(st.ID, now.Add(-90*time.Minute), now, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 assessment attempts in range, got %d", len(got))
	}
	if got[0].Source != "quiz" || got[1].Source != "diagnostic" {
		t.Fatalf("want newest first quiz then diagnostic, got %q then %q", got[0].Source, got[1].Source)
	}
	for _, a := range got {
		if a.Source == "practice" {
			t.Fatal("practice attempt leaked into assessment history")
		}
	}
}

func TestXPInRangeMatchesLedger(t *testing.T) {
	repo := newTestStore(t)
	st, err := repo.CreateStudent("Learner")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.AwardXP(st.ID, 10, "test", "k1"); err != nil {
		t.Fatal(err)
	}
	if err := repo.AwardXP(st.ID, 5, "test", "k2"); err != nil {
		t.Fatal(err)
	}
	// The same key must not add XP twice.
	if err := repo.AwardXP(st.ID, 99, "test", "k1"); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	got, err := repo.XPInRange(st.ID, now.Add(-time.Hour), now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if got != 15 {
		t.Fatalf("want 15 XP from the ledger, got %d", got)
	}
	future, err := repo.XPInRange(st.ID, now.Add(time.Hour), now.Add(2*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if future != 0 {
		t.Fatalf("want 0 XP in an empty range, got %d", future)
	}
}

func TestLearningEventsPageSubstantiatesAndPaginates(t *testing.T) {
	repo := newTestStore(t)
	st, err := repo.CreateStudent("Learner")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	masteredAt := now.Add(-4 * time.Hour)
	if err := repo.UpsertProgress(&ConceptProgress{
		StudentID:  st.ID,
		ConceptID:  "a",
		Status:     "MASTERED",
		MasteredAt: &masteredAt,
	}); err != nil {
		t.Fatal(err)
	}
	seedAttempt(t, repo, st.ID, "a", "quiz", true, now.Add(-3*time.Hour))
	if err := repo.RecordQuizCompletion(st.ID, 120); err != nil {
		t.Fatal(err)
	}
	if err := repo.SetDiagnosticCompleted(st.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.AwardXP(st.ID, -5, "accountability:test", "adjustment:1"); err != nil {
		t.Fatal(err)
	}

	first, total, err := repo.LearningEventsPage(st.ID, 3, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total != 5 {
		t.Fatalf("want 5 substantiated events including the XP adjustment, got %d", total)
	}
	if len(first) != 3 {
		t.Fatalf("want a page of 3, got %d", len(first))
	}
	second, total2, err := repo.LearningEventsPage(st.ID, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if total2 != total {
		t.Fatalf("total changed between pages: %d then %d", total, total2)
	}
	seen := map[string]bool{}
	byKind := map[string]int{}
	for _, e := range append(first, second...) {
		key := e.Kind + ":" + e.Ref
		if seen[key] {
			t.Fatalf("duplicate event %s", key)
		}
		seen[key] = true
		byKind[e.Kind]++
	}
	for _, kind := range []string{"answer", "mastery", "quiz", "diagnostic", "xp"} {
		if byKind[kind] != 1 {
			t.Fatalf("want exactly one %s event, got %d", kind, byKind[kind])
		}
	}
	for _, e := range append(first, second...) {
		if e.Kind == "xp" {
			if e.XP != -5 || e.Reason != "accountability:test" {
				t.Fatalf("xp event = %+v, want -5 accountability:test", e)
			}
		}
	}
	if len(seen) != 5 {
		t.Fatalf("want 5 distinct events, got %d", len(seen))
	}
}

func TestCompletionCutoffParsesEveryWrittenFormat(t *testing.T) {
	base := time.Date(2026, 11, 7, 3, 5, 44, 0, time.UTC)
	cases := []struct {
		raw  string
		want time.Time
	}{
		{base.Format(time.RFC3339Nano), base},
		{base.Format(time.RFC3339), base},
		{"2026-11-07 03:05:44.123456+00", time.Date(2026, 11, 7, 3, 5, 44, 123456000, time.UTC)},
		{"2026-11-07 03:05:44+00", base},
		{"2026-11-07 03:05:44", base.Add(time.Second)},
	}
	for _, tc := range cases {
		got, ok := completionCutoff(tc.raw)
		if !ok || !got.Equal(tc.want) {
			t.Errorf("completionCutoff(%q) = %v, %v; want %v, true", tc.raw, got, ok, tc.want)
		}
	}
	if _, ok := completionCutoff("not a timestamp"); ok {
		t.Fatal("invalid timestamp was accepted")
	}
}

func TestQuizCompletionsRespectSince(t *testing.T) {
	repo := newTestStore(t)
	st, err := repo.CreateStudent("Learner")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.RecordQuizCompletion(st.ID, 80); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	recent, err := repo.QuizCompletions(st.ID, now.Add(-time.Hour), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(recent) != 1 || recent[0].XPTotal != 80 {
		t.Fatalf("want one completion with XP 80, got %+v", recent)
	}
	older, err := repo.QuizCompletions(st.ID, now.Add(time.Hour), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(older) != 0 {
		t.Fatalf("want no completions after now, got %d", len(older))
	}
}

// TestReportsParity exercises the four report queries through both engines, so
// the SQLite and PostgreSQL implementations cannot drift while both compile.
// The Postgres half runs only when the storage job supplies TEST_POSTGRES_DSN.
func TestReportsParity(t *testing.T) {
	stores := map[string]Repository{"sqlite": newTestStore(t)}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		pg, err := NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { pg.Close() })
		stores["postgres"] = pg
	}

	for name, repo := range stores {
		t.Run(name, func(t *testing.T) {
			st, err := repo.CreateStudent(fmt.Sprintf("Parity %s %d", name, time.Now().UnixNano()))
			if err != nil {
				t.Fatalf("create student: %v", err)
			}
			sid := st.ID
			now := time.Now().UTC()
			seedAttempt(t, repo, sid, "c1", "quiz", true, now.Add(-2*time.Hour))
			seedAttempt(t, repo, sid, "c1", "practice", true, now.Add(-90*time.Minute))
			seedAttempt(t, repo, sid, "c2", "diagnostic", false, now.Add(-30*time.Minute))

			masteredAt := now.Add(-3 * time.Hour)
			if err := repo.UpsertProgress(&ConceptProgress{
				StudentID: sid, ConceptID: "c1", Status: "MASTERED", MasteredAt: &masteredAt,
			}); err != nil {
				t.Fatalf("upsert progress: %v", err)
			}
			if err := repo.AwardXP(sid, 12, "parity", "parity-xp"); err != nil {
				t.Fatalf("award xp: %v", err)
			}
			if err := repo.RecordQuizCompletion(sid, 100); err != nil {
				t.Fatalf("record quiz completion: %v", err)
			}
			if err := repo.SetDiagnosticCompleted(sid); err != nil {
				t.Fatalf("set diagnostic completed: %v", err)
			}

			assessments, err := repo.AssessmentAttempts(sid, now.Add(-3*time.Hour), now.Add(time.Hour), 10)
			if err != nil {
				t.Fatalf("assessment attempts: %v", err)
			}
			if len(assessments) != 2 || assessments[0].ConceptID != "c2" || assessments[1].ConceptID != "c1" {
				t.Fatalf("assessment attempts = %+v, want diagnostic c2 then quiz c1", assessments)
			}

			xp, err := repo.XPInRange(sid, now.Add(-time.Hour), now.Add(time.Hour))
			if err != nil {
				t.Fatalf("xp in range: %v", err)
			}
			if xp != 12 {
				t.Fatalf("xp in range = %d, want 12", xp)
			}

			events, total, err := repo.LearningEventsPage(sid, 10, 0)
			if err != nil {
				t.Fatalf("learning events: %v", err)
			}
			if total != 6 {
				t.Fatalf("learning events total = %d, want 6", total)
			}
			if len(events) != 6 {
				t.Fatalf("learning events = %d, want 6", len(events))
			}

			completions, err := repo.QuizCompletions(sid, now.Add(-time.Hour), 5)
			if err != nil {
				t.Fatalf("quiz completions: %v", err)
			}
			if len(completions) != 1 || completions[0].XPTotal != 100 {
				t.Fatalf("quiz completions = %+v, want one with xp 100", completions)
			}
		})
	}
}
