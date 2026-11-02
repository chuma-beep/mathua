package engine

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/statistics"
	"github.com/chuma-beep/mathua/internal/storage"
)

func TestEngine_MasteryPathStatisticsChain(t *testing.T) {
	chain := []struct {
		id      string
		prereqs []string
		answer  string
	}{
		{id: "stat.data.mean", prereqs: nil, answer: "10"},
		{id: "stat.adv.variance", prereqs: nil, answer: "20"},
		{id: "stat.spread.std_dev", prereqs: []string{"stat.data.mean", "stat.adv.variance"}, answer: "30"},
		{id: "stat.reg.correlation", prereqs: []string{"stat.spread.std_dev"}, answer: "40"},
		{id: "stat.reg.least_squares", prereqs: []string{"stat.reg.correlation"}, answer: "50"},
	}

	defs := make([]concepts.Concept, len(chain))
	answers := make(map[string]string, len(chain))
	for i, c := range chain {
		defs[i] = concepts.Concept{
			ID: c.id, Label: c.id, Domain: "statistics",
			GradingType: "numeric", Prerequisites: c.prereqs,
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		}
		answers[c.id] = c.answer
	}

	dag, err := concepts.Build(defs)
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	reg := generator.NewRegistry()
	for id, ans := range answers {
		if err := reg.Register(id, &testGen{answer: ans}); err != nil {
			t.Fatalf("register %s: %v", id, err)
		}
	}
	e := New(store, dag, reg, nil, nil)
	st, err := e.CreateStudent("synthetic_stats_chain")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	for _, id := range []string{"stat.spread.std_dev", "stat.reg.correlation", "stat.reg.least_squares"} {
		elig, err := e.EligibilityFor(st.ID, id)
		if err != nil {
			t.Fatalf("eligibility for %s: %v", id, err)
		}
		if elig.Eligible {
			t.Fatalf("%s should be locked before its prerequisites are mastered", id)
		}
	}

	// A locked concept must not be recommended while only one of its
	// prerequisites is mastered.
	if err := e.masterThroughStudy(st.ID, chain[0].id, answers[chain[0].id], 3); err != nil {
		t.Fatalf("master %s: %v", chain[0].id, err)
	}
	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	e.mu.Lock()
	e.sessions[sess.ID] = &activeSession{}
	e.mu.Unlock()
	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if q.ConceptID == "stat.spread.std_dev" {
		t.Fatal("recommended std_dev while variance is still unmastered")
	}

	for i, c := range chain {
		if i == 0 {
			continue // already mastered above
		}
		if err := e.masterThroughStudy(st.ID, c.id, answers[c.id], 3); err != nil {
			t.Fatalf("master %s: %v", c.id, err)
		}
	}

	progress, err := e.GetProgress(st.ID)
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	for _, c := range chain {
		p, ok := progress[c.id]
		if !ok {
			t.Fatalf("no progress row for %s", c.id)
		}
		if p.Status != "MASTERED" {
			t.Errorf("%s status = %s, want MASTERED", c.id, p.Status)
		}
		if p.MasteredAt == nil {
			t.Errorf("%s has no mastered timestamp", c.id)
		}
		if p.NextReviewDue == nil {
			t.Errorf("%s has no next review date", c.id)
		}
		if p.SM2Repetitions == 0 || p.SM2Interval == 0 || p.SM2EFactor == 0 {
			t.Errorf("%s has no SM-2 schedule: reps=%d interval=%d ef=%.2f", c.id, p.SM2Repetitions, p.SM2Interval, p.SM2EFactor)
		}
		if p.FSRSStability == nil || p.FSRSDifficulty == nil {
			t.Errorf("%s has no FSRS schedule", c.id)
		}
		t.Logf("%s: streak=%d sm2_reps=%d sm2_interval=%d sm2_ef=%.2f fsrs_stability=%v fsrs_difficulty=%v next_review=%v",
			c.id, p.Streak, p.SM2Repetitions, p.SM2Interval, p.SM2EFactor, p.FSRSStability, p.FSRSDifficulty, p.NextReviewDue)
	}

	// The real generators for this chain must serve all three sections.
	realReg := generator.NewRegistry()
	statistics.Register(realReg)
	for _, c := range chain {
		seen := map[string]bool{}
		for ordinal := 1; ordinal <= 3; ordinal++ {
			p, err := realReg.GenerateContext(c.id, generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: ordinal})
			if err != nil {
				t.Fatalf("generate %s kp%d: %v", c.id, ordinal, err)
			}
			if p.Question == "" || p.Answer == "" {
				t.Errorf("%s kp%d produced an empty question or answer", c.id, ordinal)
			}
			seen[p.Question] = true
		}
		if len(seen) < 2 {
			t.Errorf("%s produced fewer than two distinct questions across its three sections", c.id)
		}
	}
}

func (e *Engine) masterThroughStudy(studentID, conceptID, answer string, needed int) error {
	e.SetStudyAnchor(studentID, conceptID, answer)
	for i := 0; i < needed+5; i++ {
		if _, err := e.SubmitStudyAnswer(studentID, conceptID, answer, 1.0, "Synthetic mastery-path question"); err != nil {
			return err
		}
		progress, err := e.GetProgress(studentID)
		if err != nil {
			return err
		}
		if p := progress[conceptID]; p != nil && p.Status == "MASTERED" {
			return nil
		}
	}
	return nil
}
