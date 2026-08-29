package diagnostic

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

type fakeGen struct{}

func (f *fakeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "test?",
		Answer:      "test",
		Explanation: "test",
	}
}

func testDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	raw := make([]concepts.Concept, 10)
	for i := range raw {
		raw[i] = concepts.Concept{
			ID:               idFor(i),
			Domain:           "d",
			Prerequisites:    prereqsFor(i),
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		}
	}
	d, err := concepts.Build(raw)
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	return d
}

func idFor(i int) string {
	letters := "abcdefghij"
	return string(letters[i])
}

func prereqsFor(i int) []string {
	if i == 0 {
		return nil
	}
	return []string{idFor(i - 1)}
}

func mustRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	r := generator.NewRegistry()
	for i := 0; i < 10; i++ {
		if err := r.Register(idFor(i), &fakeGen{}); err != nil {
			t.Fatalf("register: %v", err)
		}
	}
	return r
}

func TestEngine_Start(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	if s.State != StateProbing {
		t.Errorf("expected probing, got %q", s.State)
	}
}

func TestEngine_NextQuestion(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	p, cid, err := e.NextQuestion(s)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if p.Question == "" || cid == "" {
		t.Error("expected non-empty question and concept ID")
	}
	if s.LastConceptID != cid {
		t.Errorf("expected LastConceptID=%q, got %q", cid, s.LastConceptID)
	}
}

func TestEngine_CompressedCover(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	cover := e.compressedCover(s)
	if len(cover) == 0 {
		t.Fatal("expected non-empty covering set")
	}
	if cover[0].ID != "a" {
		t.Errorf("expected root a first in cover, got %s", cover[0].ID)
	}
}

func TestEngine_InfoGain_UncertaintyFirst(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// Settle concept "f" with a high belief (2 correct probes).
	s.beliefs["f"] = 0.9
	s.probeCounts["f"] = 2
	s.confidence["f"] = 1.0
	s.doneSet["f"] = true
	cid := e.pickByInfoGain(s)
	if cid == "" {
		t.Fatal("expected a candidate after settling f")
	}
	if cid == "f" {
		t.Error("settled concept f should not be re-picked")
	}
	if s.beliefs[cid] == 0.9 {
		t.Errorf("expected an unsettled (0.5-belief) concept to win info gain, got %s", cid)
	}
}

func TestEngine_EvidencePropagation_CorrectBoostsPrereqs(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// "e" depends on "d" (prereq chain a<-b<-c<-d<-e). Correct on e should
	// raise belief on its prerequisite d.
	e.RecordAnswer(s, "e", true, false)
	if s.beliefs["e"] < 0.6 {
		t.Errorf("expected e belief up after correct, got %f", s.beliefs["e"])
	}
	if s.beliefs["d"] <= 0.5 {
		t.Errorf("expected prereq d belief boosted above 0.5, got %f", s.beliefs["d"])
	}
}

func TestEngine_EvidencePropagation_WrongLowersDependents(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// Wrong on "d" should lower belief on dependent "e".
	e.RecordAnswer(s, "d", false, false)
	if s.beliefs["d"] >= 0.5 {
		t.Errorf("expected d belief below 0.5 after wrong, got %f", s.beliefs["d"])
	}
	if s.beliefs["e"] >= 0.5 {
		t.Errorf("expected dependent e belief lowered, got %f", s.beliefs["e"])
	}
}

// answer answers the session's current question via the real Next→Record flow.
func answer(t *testing.T, e *Engine, s *Session, correct, fast bool) string {
	t.Helper()
	_, cid, err := e.NextQuestion(s)
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	if cid == "" {
		return ""
	}
	e.RecordAnswer(s, cid, correct, fast)
	return cid
}

func TestEngine_Frontier_StrongStudent(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// Strong student: everything correct+fast through the real flow.
	for i := 0; i < 60 && !e.IsComplete(s); i++ {
		if answer(t, e, s, true, true) == "" {
			break
		}
	}
	rep := e.Report(s)
	if rep.FrontierIdx < 8 {
		t.Errorf("expected strong student frontier near top (>=8), got %d", rep.FrontierIdx)
	}
}

func TestEngine_Frontier_WeakStudent(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	for i := 0; i < 60 && !e.IsComplete(s); i++ {
		if answer(t, e, s, false, false) == "" {
			break
		}
	}
	rep := e.Report(s)
	if rep.FrontierIdx > 3 {
		t.Errorf("expected weak student frontier low (<=3), got %d", rep.FrontierIdx)
	}
	if len(rep.GapsByDomain["d"]) == 0 {
		t.Error("expected gaps recorded for weak student")
	}
}

func TestEngine_Report_Fields(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	answer(t, e, s, true, true)
	answer(t, e, s, true, true)
	rep := e.Report(s)
	if rep.TotalQuestions != 2 {
		t.Errorf("expected total_questions 2, got %d", rep.TotalQuestions)
	}
	if rep.MasteryLevels[s.LastConceptID] < 0.6 {
		t.Errorf("expected mastery of answered concept high, got %f", rep.MasteryLevels[s.LastConceptID])
	}
	if rep.Confidence[s.LastConceptID] < 0.5 {
		t.Errorf("expected confidence >= 0.5 after probing, got %f", rep.Confidence[s.LastConceptID])
	}
	if rep.CompletionEstimates[150] == "" {
		t.Error("expected completion estimate for 150 XP")
	}
}

func TestEngine_Supplemental_LowConfidenceExtends(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	// Answer only 1 probe per concept -> confidence 0.5 (< 0.7) -> not done
	// until min questions met and all settled.
	for i := 0; i < 10; i++ {
		e.RecordAnswer(s, idFor(i), true, true)
	}
	if e.IsComplete(s) {
		t.Error("expected not complete with only 1 probe per concept (low confidence)")
	}
	// Second probe per concept settles them.
	for i := 0; i < 10; i++ {
		e.RecordAnswer(s, idFor(i), true, true)
	}
	if !e.IsComplete(s) {
		t.Error("expected complete after 2 probes per concept")
	}
}

func TestEngine_AllProbedCompletes(t *testing.T) {
	e := NewEngine(testDAG(t), mustRegistry(t))
	s := e.Start()
	for i := 0; i < len(s.order) && !e.IsComplete(s); i++ {
		cid := idFor(i)
		if s.totalCount[cid] < 2 {
			e.RecordAnswer(s, cid, true, true)
			e.RecordAnswer(s, cid, true, true)
		}
	}
	if !e.IsComplete(s) {
		t.Errorf("expected complete after all probed, state=%s asked=%d", s.State, s.totalAsked)
	}
}
