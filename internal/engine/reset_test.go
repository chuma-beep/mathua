package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func futureTS() string {
	return time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
}

func TestResetAccountProgress_PhraseGate(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("reset_gate")
	if err := e.ResetAccountProgress(st.ID, "reset"); err == nil {
		t.Error("wrong phrase must fail")
	}
	if err := e.ResetAccountProgress(st.ID, " reset my progress "); err != nil {
		t.Errorf("padded phrase should pass (server trims whitespace): %v", err)
	}
	if err := e.ResetAccountProgress(st.ID, ResetPhrase); err != nil {
		t.Errorf("exact phrase: %v", err)
	}
	// Retry-safe: second run is a no-op success.
	if err := e.ResetAccountProgress(st.ID, ResetPhrase); err != nil {
		t.Errorf("retry: %v", err)
	}
}

func TestResetAccountProgress_ClearsState(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("reset_state")
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{StudentID: st.ID, ConceptID: "a", Status: "PRACTICING", Streak: 3}); err != nil {
		t.Fatal(err)
	}
	e.SetStudyAnchor(st.ID, "a", "42")
	e.studyMisses[st.ID+"|a"] = 2
	if err := e.savePlanPrefs(st.ID, &StudyPlanPrefs{Destination: "math-for-ml", DeadlineDays: 90, RestDays: 1}); err != nil {
		t.Fatal(err)
	}
	if err := e.repo.UpsertServerSession("study_plan", st.ID, `{"destination":"math-for-ml","xp_remaining":100}`, futureTS()); err != nil {
		t.Fatal(err)
	}
	if err := e.ResetAccountProgress(st.ID, ResetPhrase); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if p, _ := e.repo.GetAllProgress(st.ID); len(p) != 0 {
		t.Errorf("progress survives: %d", len(p))
	}
	if _, ok := e.studyAnchorFor(st.ID, "a", ""); ok {
		t.Error("in-memory anchor survives")
	}
	if _, ok := e.studyMisses[st.ID+"|a"]; ok {
		t.Error("in-memory misses survive")
	}
	// Evidence gone, prefs kept with a reset marker.
	if _, ok := e.GetPlan(st.ID); ok {
		t.Error("evidence snapshot survives")
	}
	prefs, ok := e.GetPlanPrefs(st.ID)
	if !ok || prefs.Destination != "math-for-ml" || prefs.DeadlineDays != 90 || prefs.ResetAt == "" {
		t.Errorf("prefs must survive with reset marker: %+v %v", prefs, ok)
	}
}

func TestPlanPrefs_LegacyMigration(t *testing.T) {
	e := testEngine(t)
	st, _ := e.CreateStudent("legacy_plan")
	// Pre-split single record (old SavePlan shape with daily goal inline).
	legacy := `{"destination":"math-for-ml","daily_goal":10,"deadline_days":90,"rest_days":1,"xp_remaining":50}`
	if err := e.repo.UpsertServerSession("study_plan", st.ID, legacy, futureTS()); err != nil {
		t.Fatal(err)
	}
	prefs, ok := e.GetPlanPrefs(st.ID)
	if !ok || prefs.Destination != "math-for-ml" || prefs.DeadlineDays != 90 || prefs.RestDays != 1 {
		t.Errorf("migration failed: %+v %v", prefs, ok)
	}
	// Evidence row itself is untouched by the migration.
	if _, ok := e.GetPlan(st.ID); !ok {
		t.Error("evidence must survive migration")
	}
}

// A confirmation typed on a phone arrives capitalised: Safari upper-cases the first
// letter of a `type="text"` input and Android keyboards autocorrect. An exact
// comparison rejected people who had typed exactly what the label asked for, and told
// them they had mistyped it.
func TestPhraseMatches_CasingAndWhitespace(t *testing.T) {
	for _, typed := range []string{
		"reset my progress",
		"Reset my progress",
		"RESET MY PROGRESS",
		"  reset my progress  ",
		"Reset My Progress",
	} {
		if !PhraseMatches(typed, ResetPhrase) {
			t.Errorf("PhraseMatches(%q, ResetPhrase) = false, want true", typed)
		}
	}
	for _, typed := range []string{
		"reset",
		"reset my progress!",
		"reset my progres",
		"delete my account", // the other phrase must not unlock this one
		"",
	} {
		if PhraseMatches(typed, ResetPhrase) {
			t.Errorf("PhraseMatches(%q, ResetPhrase) = true, want false", typed)
		}
	}
	if !PhraseMatches("Delete My Account", DeletePhrase) {
		t.Error("the delete phrase must tolerate casing too")
	}
}
