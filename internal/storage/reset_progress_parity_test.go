package storage

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestResetProgress_RoundTrip_Parity runs the reset contract against SQLite
// always and Postgres when TEST_POSTGRES_DSN is set. ADR-012 makes Postgres a
// first-class store, but ResetProgress previously had SQLite-only coverage, so
// its Postgres statement list — including the server_sessions LIKE predicate —
// was never executed by any test.
func TestResetProgress_RoundTrip_Parity(t *testing.T) {
	stores := map[string]Repository{"sq": newTestStore(t)}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		pg, err := NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { pg.Close() })
		stores["pg"] = pg
	}

	for tag, store := range stores {
		run := func(step string, err error) {
			t.Helper()
			if err != nil {
				t.Fatalf("[%s] %s: %v", tag, step, err)
			}
		}
		sid := fmt.Sprintf("resetter-%s-%d", tag, time.Now().UnixNano())

		st, err := store.CreateStudent(sid)
		run("create student", err)
		sess, err := store.CreateSession(st.ID)
		run("create session", err)
		run("seed progress", store.UpsertProgress(&ConceptProgress{StudentID: st.ID, ConceptID: "a", Status: "PRACTICING", Streak: 3}))
		run("seed attempt", store.RecordAttempt(AttemptEntry{SessionID: sess.ID, StudentID: st.ID, ConceptID: "a", Answer: "1", Expected: "2", Correct: false, ElapsedSeconds: 4, Timestamp: time.Now().UTC()}))
		run("seed speed", store.UpsertTopicSpeed(&TopicSpeed{StudentID: st.ID, ConceptID: "a", EFactor: 2.5, Interval: 1, Repetitions: 1, LearningSpeed: 1.0}))
		run("seed quiz", store.RecordQuizCompletion(st.ID, 60))
		run("seed xp", store.AddXP(st.ID, 42))
		run("seed anchor", store.UpsertServerSession("study_expected", st.ID+"|a", `{"q":"1"}`, time.Now().UTC().Add(time.Hour).Format(time.RFC3339)))
		run("seed plan", store.UpsertServerSession("study_plan", st.ID, `{"destination":"d"}`, time.Now().UTC().Add(time.Hour).Format(time.RFC3339)))
		run("seed prefs", store.UpsertServerSession("study_plan_prefs", st.ID, `{"pace":"steady"}`, time.Now().UTC().Add(time.Hour).Format(time.RFC3339)))
		run("seed admin", store.UpsertServerSession("admin_login", st.ID, "tok", time.Now().UTC().Add(time.Hour).Format(time.RFC3339)))

		run("reset", store.ResetProgress(st.ID))

		if p, _ := store.GetAllProgress(st.ID); len(p) != 0 {
			t.Errorf("[%s] progress survives: %d rows", tag, len(p))
		}
		if a, _ := store.GetAttemptsForStudent(st.ID); len(a) != 0 {
			t.Errorf("[%s] attempts survive: %d rows", tag, len(a))
		}
		if s, _ := store.GetAllTopicSpeeds(st.ID); len(s) != 0 {
			t.Errorf("[%s] speeds survive: %d rows", tag, len(s))
		}
		if q, _ := store.LastQuizCompletion(st.ID); q != nil {
			t.Errorf("[%s] quiz completion survives: %+v", tag, q)
		}
		if _, _, found, _ := store.GetServerSession("study_expected", st.ID+"|a"); found {
			t.Errorf("[%s] study anchor survives", tag)
		}
		if _, _, found, _ := store.GetServerSession("study_plan", st.ID); found {
			t.Errorf("[%s] study plan evidence survives", tag)
		}
		// ADR-022: prefs survive a reset (destination/deadline/pace are kept).
		if _, _, found, _ := store.GetServerSession("study_plan_prefs", st.ID); !found {
			t.Errorf("[%s] study_plan_prefs wrongly wiped", tag)
		}
		if _, _, found, _ := store.GetServerSession("admin_login", st.ID); !found {
			t.Errorf("[%s] admin session wrongly wiped", tag)
		}
		got, err := store.GetStudent(st.ID)
		if err != nil || got == nil {
			t.Fatalf("[%s] account must survive: %v", tag, err)
		}
		if got.XPTotal != 0 || got.XPToday != 0 || got.DiagnosticCompleted || got.Name != sid {
			t.Errorf("[%s] student row wrong after reset: %+v", tag, got)
		}
		// Idempotent: a second reset is a no-op success.
		if err := store.ResetProgress(st.ID); err != nil {
			t.Errorf("[%s] second reset: %v", tag, err)
		}
	}
}
