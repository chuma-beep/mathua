package diagnostic

import (
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

type State string

const (
	StateProbing State = "probing"
	StateDone    State = "done"
)

const (
	// probesPerConcept is how many questions to ask for each concept.
	probesPerConcept = 2
	// minTotalQuestions is the minimum number of questions before
	// the diagnostic will consider itself converged.
	minTotalQuestions = 12
)

type Session struct {
	sync.Mutex

	ID        string
	StudentID string
	Position  int
	Low       int
	High      int
	State     State
	Attempts  []Attempt
	order     []*concepts.Concept

	// LastProblem stores the current problem the student saw.
	LastProblem     *generator.Problem
	LastConceptID   string
	LastConceptName string

	// Per-concept tracking
	probeCounts  map[string]int
	correctCount map[string]int
	totalCount   map[string]int
	totalAsked   int
	doneSet      map[string]bool
}

type Attempt struct {
	ConceptID      string
	Correct        bool
	Fast           bool
	ElapsedSeconds float64
	Timestamp      time.Time
}

type Engine struct {
	dag      *concepts.DAG
	registry *generator.Registry
}

func NewEngine(dag *concepts.DAG, reg *generator.Registry) *Engine {
	return &Engine{dag: dag, registry: reg}
}

func (e *Engine) Start() *Session {
	order := e.dag.Order()
	return e.StartWithPath(order)
}

func (e *Engine) StartWithPath(path []*concepts.Concept) *Session {
	if len(path) == 0 {
		return &Session{State: StateDone, order: path}
	}
	mid := len(path) / 2
	return &Session{
		Position:     mid,
		Low:          0,
		High:         len(path) - 1,
		State:        StateProbing,
		order:        path,
		probeCounts:  make(map[string]int),
		correctCount: make(map[string]int),
		totalCount:   make(map[string]int),
		doneSet:      make(map[string]bool),
	}
}

func (e *Engine) NextQuestion(s *Session) (*generator.Problem, string, error) {
	s.Lock()
	defer s.Unlock()
	if s.State == StateDone {
		s.LastProblem = nil
		return nil, "", nil
	}

	// Find the next concept that needs probing, working outward
	// from Position through the binary search range.
	cid := s.nextConceptToProbe()
	if cid == "" {
		s.State = StateDone
		s.LastProblem = nil
		return nil, "", nil
	}
	concept := s.order[s.Position]
	s.probeCounts[cid]++

	// Vary difficulty: start low, increase on later probes
	difficulty := 0.3
	if s.probeCounts[cid] >= 2 {
		difficulty = 0.6
	}

	p, err := e.registry.Generate(cid, difficulty)
	if err != nil {
		return nil, "", err
	}
	s.LastProblem = &p
	s.LastConceptID = cid
	s.LastConceptName = concept.Label
	return s.LastProblem, cid, nil
}

// nextConceptToProbe finds a concept within the binary search range [Low, High]
// that hasn't been probed enough yet, starting from Position.
func (s *Session) nextConceptToProbe() string {
	if s.Position < 0 || s.Position >= len(s.order) {
		return ""
	}
	cid := s.order[s.Position].ID
	if s.probeCounts[cid] < probesPerConcept && s.Position >= s.Low && s.Position <= s.High && !s.doneSet[cid] {
		return cid
	}
	// Try concepts within the range [Low, High]
	for i := s.Position; i <= s.High; i++ {
		cid2 := s.order[i].ID
		if s.probeCounts[cid2] < probesPerConcept && i >= s.Low && !s.doneSet[cid2] {
			s.Position = i
			return cid2
		}
	}
	for i := s.Position; i >= s.Low; i-- {
		cid2 := s.order[i].ID
		if s.probeCounts[cid2] < probesPerConcept && i <= s.High && !s.doneSet[cid2] {
			s.Position = i
			return cid2
		}
	}
	return ""
}

// RecordAnswer processes a diagnostic answer and updates the session state.
func (e *Engine) RecordAnswer(s *Session, conceptID string, correct, fast bool) {
	s.Lock()
	defer s.Unlock()
	s.Attempts = append(s.Attempts, Attempt{
		ConceptID: conceptID,
		Correct:   correct,
		Fast:      fast,
		Timestamp: time.Now().UTC(),
	})
	s.totalCount[conceptID]++
	s.totalAsked++
	if correct {
		s.correctCount[conceptID]++
	}

	needed := probesPerConcept
	// If a concept is clearly right or clearly wrong after fewer probes,
	// decide early. But respect minimum probes.
	if s.totalCount[conceptID] >= needed {
		s.evaluateConcept(conceptID)
	}

	// Check stopping conditions
	if s.shouldStop() {
		s.State = StateDone
	}
}

// evaluateConcept decides whether a concept is "known" or "not known"
// and adjusts the binary search bounds accordingly.
func (s *Session) evaluateConcept(conceptID string) {
	total := s.totalCount[conceptID]
	correct := s.correctCount[conceptID]
	if total == 0 || s.doneSet[conceptID] {
		return
	}
	s.doneSet[conceptID] = true
	ratio := float64(correct) / float64(total)

	if ratio >= 0.75 {
		// Student knows this and everything below it.
		s.Low = s.Position + 1
	} else if ratio <= 0.25 {
		// Student doesn't know this — everything above it is
		// also likely unknown.
		s.High = s.Position - 1
	}
	// For 0.5 (1/2), don't move bounds aggressively — let more
	// probing inform the decision.

	// Recompute binary search position
	if s.Low <= s.High {
		s.Position = (s.Low + s.High) / 2
	}
}

// shouldStop returns true when the diagnostic has converged.
func (s *Session) shouldStop() bool {
	if s.Low > s.High {
		return s.totalAsked >= minTotalQuestions || s.allProbed()
	}
	// Also stop if all concepts in range are done and minimum met
	if s.totalAsked >= minTotalQuestions && s.allInRangeDone() {
		return true
	}
	return false
}

// allProbed returns true when every concept in the order has been probed
// at least once.
func (s *Session) allProbed() bool {
	for _, c := range s.order {
		if s.totalCount[c.ID] == 0 {
			return false
		}
	}
	return true
}

// allInRangeDone returns true when every concept in [Low, High] has been
// fully evaluated.
func (s *Session) allInRangeDone() bool {
	for i := s.Low; i <= s.High && i < len(s.order); i++ {
		cid := s.order[i].ID
		if !s.doneSet[cid] && s.totalCount[cid] < probesPerConcept {
			return false
		}
	}
	return true
}

func (e *Engine) IsComplete(s *Session) bool {
	s.Lock()
	defer s.Unlock()
	return s.State == StateDone
}

// FrontierEstimate returns the estimated concept index where the student's
// knowledge frontier lies (the highest concept index they can answer).
func (e *Engine) FrontierEstimate(s *Session) int {
	s.Lock()
	defer s.Unlock()
	return s.High
}
