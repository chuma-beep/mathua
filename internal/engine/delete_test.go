package engine

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestDeleteAccount_PhraseGate(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("del_gate")
	if err := e.DeleteAccount(st.ID, "delete"); err == nil {
		t.Error("wrong phrase must fail")
	}
	if err := e.DeleteAccount(st.ID, " delete my account "); err != nil {
		t.Errorf("padded phrase should pass: %v", err)
	}
}

func TestDeleteAccount_ClearsRowsAndMemory(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("del_state")
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{StudentID: st.ID, ConceptID: "a", Status: "PRACTICING", Streak: 3}); err != nil {
		t.Fatal(err)
	}
	if err := e.repo.SetEmail(st.ID, "gone@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := e.repo.CreateIdentity("google", "gid-del", st.ID, "gone@example.com", true); err != nil {
		t.Fatal(err)
	}
	e.SetStudyExpected(st.ID, "a", "42")
	e.studyMisses[st.ID+"|a"] = 2
	if err := e.savePlanPrefs(st.ID, &StudyPlanPrefs{Destination: "math-for-ml", DeadlineDays: 90, RestDays: 1}); err != nil {
		t.Fatal(err)
	}
	if err := e.repo.UpsertServerSession("study_plan", st.ID, `{"destination":"math-for-ml","xp_remaining":100}`, futureTS()); err != nil {
		t.Fatal(err)
	}
	if err := e.DeleteAccount(st.ID, DeletePhrase); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got, _ := e.repo.GetStudent(st.ID); got != nil {
		t.Error("student row survives")
	}
	if p, _ := e.repo.GetAllProgress(st.ID); len(p) != 0 {
		t.Errorf("progress survives: %d", len(p))
	}
	e.mu.Lock()
	for k := range e.studyExpected {
		if strings.HasPrefix(k, st.ID+"|") {
			t.Errorf("memory anchor survives: %q", k)
		}
	}
	for k := range e.studyMisses {
		if strings.HasPrefix(k, st.ID+"|") {
			t.Errorf("memory miss survives: %q", k)
		}
	}
	e.mu.Unlock()
	// Retry-safe: second run is a no-op success.
	if err := e.DeleteAccount(st.ID, DeletePhrase); err != nil {
		t.Errorf("retry: %v", err)
	}
}

func TestDeleteAccount_PhraseDistinctFromReset(t *testing.T) {
	if DeletePhrase == ResetPhrase {
		t.Error("delete and reset phrases must differ so confirmations cannot cross-authorize")
	}
}
