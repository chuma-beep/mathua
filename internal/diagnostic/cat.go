package diagnostic

import (
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

type State string

const (
	StateProbing State = "probing"
	StateDone    State = "done"
)

type Session struct {
	ID            string
	StudentID     string
	Position      int
	Low           int
	High          int
	ConsecutiveOK int
	State         State
	Attempts      []Attempt
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
	mid := len(e.dag.Order()) / 2
	return &Session{
		Position: mid,
		Low:      0,
		High:     len(e.dag.Order()) - 1,
		State:    StateProbing,
	}
}

func (e *Engine) NextQuestion(s *Session) (*generator.Problem, string, error) {
	if s.State == StateDone {
		return nil, "", nil
	}
	concept := e.dag.Order()[s.Position]
	threshold := concept.MasteryThreshold.AvgTimeSeconds
	p, err := e.registry.Generate(concept.ID, 0.5)
	if err != nil {
		return nil, "", err
	}
	_ = threshold
	return &p, concept.ID, nil
}

// RecordAnswer processes a diagnostic answer and updates the session state.
// threshold is the concept's mastery time threshold.
func (e *Engine) RecordAnswer(s *Session, conceptID string, correct, fast bool) {
	s.Attempts = append(s.Attempts, Attempt{
		ConceptID: conceptID,
		Correct:   correct,
		Fast:      fast,
		Timestamp: time.Now().UTC(),
	})
	if correct && fast {
		s.Low = s.Position + 1
		s.ConsecutiveOK++
	} else {
		s.High = s.Position - 1
		s.ConsecutiveOK = 0
	}
	if s.ConsecutiveOK >= 3 {
		s.State = StateDone
		return
	}
	if s.Low > s.High {
		s.State = StateDone
		return
	}
	// Binary search midpoint
	s.Position = (s.Low + s.High) / 2
}

func (e *Engine) IsComplete(s *Session) bool {
	return s.State == StateDone
}

// FrontierEstimate returns the estimated concept index where the student's
// knowledge frontier lies (the highest concept index they can answer).
func (e *Engine) FrontierEstimate(s *Session) int {
	return s.High
}
