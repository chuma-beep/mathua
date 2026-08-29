package quiz

import (
	"sort"
	"sync"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

// Session is an actionable quiz every ~150 XP, own grading path (attemptID).
// Guest unlimited retake: each POST /api/quiz/session creates a new session.
type Session struct {
	sync.Mutex
	ID        string
	StudentID string
	Order     []*concepts.Concept
	Index     int
	Attempts  []Attempt
	LastProblem   *generator.Problem
	LastConceptID string
}

type Attempt struct {
	ConceptID string
	Correct   bool
}

type Engine struct {
	dag      *concepts.DAG
	registry *generator.Registry
}

func NewEngine(dag *concepts.DAG, reg *generator.Registry) *Engine {
	return &Engine{dag: dag, registry: reg}
}

// Start picks diverse recent weak concepts via weakness hint, fallback to DAG order.
// Caller provides weakness-filtered concepts; we cap to 5 for stub.
func (e *Engine) Start(concepts []*concepts.Concept) *Session {
	if len(concepts) == 0 {
		concepts = e.dag.Order()
		if len(concepts) > 5 {
			concepts = concepts[:5]
		}
	}
	if len(concepts) > 5 {
		concepts = concepts[:5]
	}
	return &Session{
		Order: concepts,
	}
}

func (e *Engine) NextQuestion(s *Session) (*generator.Problem, string, error) {
	s.Lock()
	defer s.Unlock()
	if s.Index >= len(s.Order) {
		s.LastProblem = nil
		return nil, "", nil
	}
	c := s.Order[s.Index]
	// 80% difficulty target via weakness — stub uses 0.7 constant (close to 0.8)
	p, err := e.registry.Generate(c.ID, 0.7)
	if err != nil {
		return nil, "", err
	}
	s.LastProblem = &p
	s.LastConceptID = c.ID
	return s.LastProblem, c.ID, nil
}

func (e *Engine) RecordAnswer(s *Session, conceptID string, correct bool) {
	s.Lock()
	defer s.Unlock()
	s.Attempts = append(s.Attempts, Attempt{ConceptID: conceptID, Correct: correct})
	// Advance index after answer
	if s.Index < len(s.Order) && s.Order[s.Index].ID == conceptID {
		s.Index++
	}
}

func (e *Engine) IsComplete(s *Session) bool {
	s.Lock()
	defer s.Unlock()
	return s.Index >= len(s.Order)
}

// PickQuizConcepts selects diverse recent weak concepts for quiz.
// Stub: sorts by weakness desc, picks top 5 distinct subdomains if possible.
func PickQuizConcepts(dag *concepts.DAG, weakness map[string]float64) []*concepts.Concept {
	type scored struct {
		c *concepts.Concept
		w float64
	}
	var list []scored
	for _, c := range dag.Order() {
		w := weakness[c.ID]
		if w > 0.3 {
			list = append(list, scored{c, w})
		}
	}
	if len(list) == 0 {
		// fallback to 5 available concepts
		order := dag.Order()
		if len(order) > 5 {
			order = order[:5]
		}
		return order
	}
	sort.Slice(list, func(i, j int) bool { return list[i].w > list[j].w })
	// distinct subdomain pick
	seen := map[string]bool{}
	var out []*concepts.Concept
	for _, s := range list {
		if !seen[s.c.Subdomain] {
			out = append(out, s.c)
			seen[s.c.Subdomain] = true
		}
		if len(out) >= 5 {
			break
		}
	}
	// fill up to 5 if still short
	for _, s := range list {
		if len(out) >= 5 {
			break
		}
		found := false
		for _, c := range out {
			if c.ID == s.c.ID {
				found = true
				break
			}
		}
		if !found {
			out = append(out, s.c)
		}
	}
	return out
}
