package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/abstract"
	"github.com/chuma-beep/mathua/internal/generator/algebra"
	"github.com/chuma-beep/mathua/internal/generator/arithmetic"
	"github.com/chuma-beep/mathua/internal/generator/calculus"
	"github.com/chuma-beep/mathua/internal/generator/complex"
	"github.com/chuma-beep/mathua/internal/generator/discrete"
	"github.com/chuma-beep/mathua/internal/generator/fractions"
	"github.com/chuma-beep/mathua/internal/generator/geometry"
	"github.com/chuma-beep/mathua/internal/generator/linalg"
	"github.com/chuma-beep/mathua/internal/generator/machinelearning"
	"github.com/chuma-beep/mathua/internal/generator/numtheory"
	"github.com/chuma-beep/mathua/internal/generator/odes"
	"github.com/chuma-beep/mathua/internal/generator/prealgebra"
	"github.com/chuma-beep/mathua/internal/generator/precalculus"
	"github.com/chuma-beep/mathua/internal/generator/statistics"
	"github.com/chuma-beep/mathua/internal/generator/topology"
	"github.com/chuma-beep/mathua/internal/generator/trigonometry"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

// productionRegistry mirrors cmd/mathua wiring: every domain registered.
func productionRegistry() *generator.Registry {
	reg := generator.NewRegistry()
	arithmetic.Register(reg)
	fractions.Register(reg)
	geometry.Register(reg)
	prealgebra.Register(reg)
	algebra.Register(reg)
	trigonometry.Register(reg)
	precalculus.Register(reg)
	statistics.Register(reg)
	numtheory.Register(reg)
	complex.Register(reg)
	linalg.Register(reg)
	machinelearning.Register(reg)
	discrete.Register(reg)
	calculus.Register(reg)
	odes.Register(reg)
	abstract.Register(reg)
	topology.Register(reg)
	return reg
}

func productionDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	return dag
}

// Every DAG concept must have a working generator at every difficulty the
// diagnostic can request (0.3 new probe, 0.6/0.8 re-probes). A missing or
// failing generator turns into HTTP 500 "failed to get next question" and
// bricks the diagnostic mid-flow with no way forward.
func TestDiagnostic_AllConceptsGenerate(t *testing.T) {
	dag := productionDAG(t)
	reg := productionRegistry()
	for _, c := range dag.Order() {
		if !reg.Has(c.ID) {
			t.Errorf("no generator registered for concept %q (%s)", c.ID, c.Label)
			continue
		}
		for _, d := range []float64{0.3, 0.6, 0.8} {
			p, err := reg.Generate(c.ID, d)
			if err != nil {
				t.Errorf("generate %q at difficulty %v: %v", c.ID, d, err)
				continue
			}
			if p.Question == "" || p.Answer == "" {
				t.Errorf("generate %q at difficulty %v: empty question/answer", c.ID, d)
			}
		}
	}
}

// Exact numeric grading has a single machine-readable answer, so
// meaning-phrasing ("What does X mean?", "describe…", "explain…") is
// ungradeable under typed input: valid descriptions are unbounded
// paraphrases. Such items must be rephrased to the graded granularity or
// retired from test rotation — never shipped as-is.
//
// NOTE: symbolic-graded meaning items (statistics/calculus "what does…"
// probes) are a separate, larger backlog: production has no python3/sympy
// (runtime is debian-slim), so gradeSymPy always returns incorrect there.
// They need a grading strategy (bundle sympy, custom graders, or regrade),
// not just rephrasing — tracked outside this lint.
var meaningPhrasing = []string{"what does", "describe", "explain"}

func TestDiagnostic_NoMeaningPhrasingWithExactGrading(t *testing.T) {
	dag := productionDAG(t)
	reg := productionRegistry()
	for _, c := range dag.Order() {
		if c.GradingType != "numeric" {
			continue
		}
		if !reg.Has(c.ID) {
			continue // reported by TestDiagnostic_AllConceptsGenerate
		}
		// Sample several generations: single-shape generators (the usual
		// case) are caught deterministically; multi-shape ones on luck.
		seen := map[string]bool{}
		for i := 0; i < 10; i++ {
			p, err := reg.Generate(c.ID, 0.3)
			if err != nil {
				break
			}
			lower := strings.ToLower(p.Question)
			for _, phrase := range meaningPhrasing {
				if strings.Contains(lower, phrase) && !seen[p.Question] {
					seen[p.Question] = true
					t.Errorf("concept %q graded %q asks meaning-phrasing: %q",
						c.ID, c.GradingType, p.Question)
				}
			}
		}
	}
}

func twoConceptServer(t *testing.T) (*Server, *http.ServeMux) {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "b", Label: "B", Domain: "d", GradingType: "numeric", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	reg.Register("b", &testGen{})
	s := New(engine.New(store, d, reg, nil, planning.New(d)), store, nil)
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux
}

func postJSON(t *testing.T, mux *http.ServeMux, path, body string) (int, map[string]interface{}) {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", path, bytes.NewReader([]byte(body))))
	var res map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	return rec.Code, res
}

// Skipping settles the pending concept with no evidence and moves CAT to a
// different question — the escape hatch for poisoned transitions.
func TestGoalDiagnosticSkip(t *testing.T) {
	_, mux := twoConceptServer(t)

	code, start := postJSON(t, mux, "/api/goal/diagnostic",
		`{"name":"tester","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["session_id"].(string)
	first, _ := start["concept_id"].(string)
	if sid == "" || first == "" {
		t.Fatalf("expected session + first concept, got %v", start)
	}

	code, skipped := postJSON(t, mux, "/api/goal/diagnostic/skip",
		`{"session_id":"`+sid+`"}`)
	if code != 200 {
		t.Fatalf("skip: %d %v", code, skipped)
	}
	if skipped["done"] == true {
		t.Fatalf("expected another question after first skip, got %v", skipped)
	}
	nextCID, _ := skipped["concept_id"].(string)
	if nextCID == "" || nextCID == first {
		t.Errorf("expected skip to advance past %q, got %q", first, nextCID)
	}
	if q, _ := skipped["question"].(string); q == "" {
		t.Error("expected a staged question after skip")
	}
	if gt, _ := skipped["grading_type"].(string); gt != "numeric" {
		t.Errorf("expected grading_type passthrough, got %v", skipped)
	}

	// Answering after a skip still works — progression intact.
	code, ans := postJSON(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"`+nextCID+`","answer":"4","elapsed":5.0}`)
	if code != 200 {
		t.Fatalf("answer after skip: %d %v", code, ans)
	}

	// Unknown session → 404, not 500.
	code, _ = postJSON(t, mux, "/api/goal/diagnostic/skip", `{"session_id":"nope"}`)
	if code != 404 {
		t.Errorf("expected 404 for unknown session, got %d", code)
	}
}

// Quiz skip advances with no grade/XP; an all-skipped quiz still completes
// (retake offered) but must not reset the 150 XP gate baseline.
func TestQuizSkip(t *testing.T) {
	_, mux := twoConceptServer(t)

	code, sess := postJSON(t, mux, "/api/quiz/session", `{}`)
	if code != 200 {
		t.Fatalf("session: %d %v", code, sess)
	}
	sid, _ := sess["session_id"].(string)
	first, _ := sess["concept_id"].(string)
	if sid == "" || first == "" {
		t.Fatalf("expected session + first concept, got %v", sess)
	}

	code, skipped := postJSON(t, mux, "/api/quiz/skip",
		`{"session_id":"`+sid+`"}`)
	if code != 200 {
		t.Fatalf("skip: %d %v", code, skipped)
	}
	if skipped["done"] == true {
		t.Fatalf("expected another question after first skip, got %v", skipped)
	}
	if skipped["correct"] == true {
		t.Errorf("skip must not grade correct, got %v", skipped)
	}
	if xp, _ := skipped["xp"].(float64); xp != 0 {
		t.Errorf("skip must award no XP, got %v", skipped)
	}
	nextCID, _ := skipped["concept_id"].(string)
	if nextCID == "" || nextCID == first {
		t.Errorf("expected skip to advance past %q, got %q", first, nextCID)
	}
	if tl, _ := skipped["time_limit_seconds"].(float64); tl <= 0 {
		t.Errorf("expected a time limit for the revealed question, got %v", skipped)
	}
	if gt, _ := skipped["grading_type"].(string); gt != "numeric" {
		t.Errorf("expected grading_type passthrough, got %v", skipped)
	}

	// Skip to the end: completes without gate reset side effects.
	code, done := postJSON(t, mux, "/api/quiz/skip", `{"session_id":"`+sid+`"}`)
	if code != 200 {
		t.Fatalf("final skip: %d %v", code, done)
	}
	if done["done"] != true {
		t.Fatalf("expected done after skipping all, got %v", done)
	}
	if done["retake_available"] != true {
		t.Errorf("expected retake_available, got %v", done)
	}

	// Unknown session → 404.
	code, _ = postJSON(t, mux, "/api/quiz/skip", `{"session_id":"nope"}`)
	if code != 404 {
		t.Errorf("expected 404 for unknown session, got %d", code)
	}
}

// Admitted unknowns record clean negative evidence: incorrect, flagged,
// progressed — and exempt from the too-quick floor (an instant admit is
// honesty, not spam).
func TestGoalDiagnosticDontKnow(t *testing.T) {
	s, mux := twoConceptServer(t)

	code, start := postJSON(t, mux, "/api/goal/diagnostic",
		`{"name":"tester","concept_ids":["a","b"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["session_id"].(string)
	first, _ := start["concept_id"].(string)

	code, res := postJSON(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"`+first+`","answer":"","elapsed":0.1,"dont_know":true}`)
	if code != 200 {
		t.Fatalf("dont_know: %d %v", code, res)
	}
	if res["correct"] == true {
		t.Errorf("dont_know must grade incorrect, got %v", res)
	}
	if res["done"] == true {
		t.Fatalf("expected another question after dont_know, got %v", res)
	}
	if nextCID, _ := res["concept_id"].(string); nextCID == "" || nextCID == first {
		t.Errorf("expected advance past %q, got %q", first, nextCID)
	}

	s.mu.Lock()
	sess := s.diagSessions[sid]
	s.mu.Unlock()
	if sess == nil {
		t.Fatal("expected session to survive dont_know")
	}
	sess.Lock()
	defer sess.Unlock()
	if len(sess.Attempts) != 1 {
		t.Fatalf("expected 1 recorded attempt, got %d", len(sess.Attempts))
	}
	att := sess.Attempts[0]
	if att.Correct || !att.DontKnow {
		t.Errorf("expected flagged incorrect attempt, got %+v", att)
	}
}

// Quiz admits: miss recorded, no XP, progression intact.
func TestQuizDontKnow(t *testing.T) {
	s, mux := twoConceptServer(t)

	code, sess := postJSON(t, mux, "/api/quiz/session", `{}`)
	if code != 200 {
		t.Fatalf("session: %d %v", code, sess)
	}
	sid, _ := sess["session_id"].(string)
	first, _ := sess["concept_id"].(string)

	code, res := postJSON(t, mux, "/api/quiz/answer",
		`{"session_id":"`+sid+`","concept_id":"`+first+`","answer":"","elapsed":0.1,"dont_know":true}`)
	if code != 200 {
		t.Fatalf("dont_know: %d %v", code, res)
	}
	if res["correct"] == true {
		t.Errorf("dont_know must grade incorrect, got %v", res)
	}
	if xp, _ := res["xp"].(float64); xp != 0 {
		t.Errorf("dont_know must award no XP, got %v", res)
	}
	if res["done"] == true {
		t.Fatalf("expected another question after dont_know, got %v", res)
	}

	s.mu.Lock()
	qsess := s.quizSessions[sid]
	s.mu.Unlock()
	if qsess == nil {
		t.Fatal("expected quiz session to survive dont_know")
	}
	qsess.Lock()
	defer qsess.Unlock()
	if len(qsess.Attempts) != 1 || qsess.Attempts[0].Correct {
		t.Errorf("expected 1 recorded miss, got %+v", qsess.Attempts)
	}
}
