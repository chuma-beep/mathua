package diagnostic

import (
	"math"
	"sort"
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
	// minTotalQuestions is the floor before the diagnostic considers convergence.
	minTotalQuestions = 15
	// maxTotalQuestions hard caps an adaptive session (supplemental included).
	maxTotalQuestions = 45
	// coverSize is the target size of the compressed covering set.
	coverSize = 30
	// beliefThreshold is the frontier cutoff: belief >= threshold is "known".
	beliefThreshold = 0.6
)

type Session struct {
	sync.Mutex

	ID        string
	StudentID string
	State     State
	Attempts  []Attempt
	order     []*concepts.Concept

	// LastProblem stores the current problem the student saw.
	LastProblem     *generator.Problem
	LastConceptID   string
	LastConceptName string

	// Per-concept student model (belief 0-1 + confidence 0-1).
	beliefs     map[string]float64
	confidence  map[string]float64
	probeCounts map[string]int
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

// Progress is the backend truth for the progress bar (MA parity: no exact
// total is promised — the Diagnostic is adaptive). CoverDone/CoverSize is
// the monotonic coverage fraction (never decreases); Answered counts
// questions answered; EstimatedTotal is a legacy hint kept for compat and
// must NOT be displayed as a promise.
type Progress struct {
	Answered       int  `json:"answered"`
	EstimatedTotal int  `json:"estimated_total"`
	MinTotal       int  `json:"min_total"`
	MaxTotal       int  `json:"max_total"`
	CoverDone      int  `json:"cover_done"`
	CoverSize      int  `json:"cover_size"`
	Done           bool `json:"done"`
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

// StartWithPath begins an adaptive session over the given path (topo order).
// It compresses the path into a minimal covering set and picks the first
// question by info gain.
func (e *Engine) StartWithPath(path []*concepts.Concept) *Session {
	if len(path) == 0 {
		return &Session{State: StateDone, order: path}
	}
	s := &Session{
		State:        StateProbing,
		order:        path,
		beliefs:      make(map[string]float64, len(path)),
		confidence:   make(map[string]float64, len(path)),
		probeCounts:  make(map[string]int, len(path)),
		correctCount: make(map[string]int, len(path)),
		totalCount:   make(map[string]int, len(path)),
		doneSet:      make(map[string]bool, len(path)),
	}
	for _, c := range path {
		s.beliefs[c.ID] = 0.5
	}
	return s
}

// compressedCover returns the smallest covering set: roots + evenly spaced
// concepts across the topo order at coverSize granularity.
func (e *Engine) compressedCover(s *Session) []*concepts.Concept {
	if len(s.order) <= coverSize {
		return s.order
	}
	step := int(math.Ceil(float64(len(s.order)) / float64(coverSize)))
	// Root concepts (no prerequisites within the path) are always covered.
	isRoot := func(c *concepts.Concept) bool {
		return len(c.Prerequisites) == 0
	}
	var out []*concepts.Concept
	seen := make(map[string]bool)
	add := func(c *concepts.Concept) {
		if !seen[c.ID] {
			seen[c.ID] = true
			out = append(out, c)
		}
	}
	for i := 0; i < len(s.order); i += step {
		add(s.order[i])
	}
	for _, c := range s.order {
		if isRoot(c) {
			add(c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return s.index(out[i].ID) < s.index(out[j].ID) })
	return out
}

func (s *Session) index(id string) int {
	for i, c := range s.order {
		if c.ID == id {
			return i
		}
	}
	return len(s.order)
}

// entropy is the binary entropy of a belief p (uncertainty).
func entropy(p float64) float64 {
	if p <= 0 || p >= 1 {
		return 0
	}
	return -(p*math.Log2(p) + (1-p)*math.Log2(1-p))
}

// infoGain scores a candidate concept: uncertainty times coverage influence.
// Coverage = how many other cover concepts it reaches through prereqs/dependents.
func (e *Engine) infoGain(s *Session, c *concepts.Concept, coverSet map[string]bool) float64 {
	unc := entropy(s.beliefs[c.ID])
	influence := 1.0
	for _, dep := range e.dag.DependentsOf(c.ID) {
		if coverSet[dep.ID] {
			influence += 0.25
		}
	}
	for _, pr := range e.dag.PrereqsOf(c.ID) {
		if coverSet[pr.ID] {
			influence += 0.25
		}
	}
	return unc * influence
}

// pickByInfoGain selects the highest-gain cover concept not yet settled.
func (e *Engine) pickByInfoGain(s *Session) string {
	cover := e.compressedCover(s)
	coverSet := make(map[string]bool, len(cover))
	for _, c := range cover {
		coverSet[c.ID] = true
	}
	best := ""
	bestScore := -1.0
	for _, c := range cover {
		if s.doneSet[c.ID] {
			continue
		}
		score := e.infoGain(s, c, coverSet)
		if score > bestScore {
			bestScore = score
			best = c.ID
		}
	}
	return best
}

func (e *Engine) NextQuestion(s *Session) (*generator.Problem, string, error) {
	s.Lock()
	defer s.Unlock()
	if s.State == StateDone {
		s.LastProblem = nil
		return nil, "", nil
	}
	cid := e.pickByInfoGain(s)
	if cid == "" {
		s.State = StateDone
		s.LastProblem = nil
		return nil, "", nil
	}
	concept := s.order[s.index(cid)]
	s.probeCounts[cid]++
	// Vary difficulty: start low, increase on later probes.
	difficulty := 0.3
	if s.probeCounts[cid] >= 2 {
		difficulty = 0.6
	}
	if s.probeCounts[cid] >= 3 {
		difficulty = 0.8
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

// RecordAnswer processes an answer: belief update + evidence propagation +
// confidence update + supplemental + stop conditions.
func (e *Engine) RecordAnswer(s *Session, conceptID string, correct, fast bool) {
	elapsed := 0.0
	if fast {
		elapsed = -1 // marker: fast without measured elapsed (legacy callers/tests)
	} else {
		elapsed = -2 // marker: not-fast without measured elapsed
	}
	e.RecordAnswerTimed(s, conceptID, correct, elapsed, 0)
}

// RecordAnswerTimed is the MA-parity path: elapsed seconds and the
// per-concept time threshold drive automaticity weighting. Correct but
// excessively slow answers get diminished weight (MA: "higher likelihood
// the student has not yet learned the topic well enough"), while fast
// correct answers keep the +0.1 automaticity bonus. Incorrect answers are
// unaffected by timing. elapsed < 0 preserves legacy fast/not-fast-only
// behavior (tests, /api/diagnostic/answer).
func (e *Engine) RecordAnswerTimed(s *Session, conceptID string, correct bool, elapsed, timeThresh float64) {
	s.Lock()
	defer s.Unlock()
	fast := false
	if elapsed < 0 {
		fast = elapsed == -1
	} else {
		thresh := timeThresh
		if thresh <= 0 {
			thresh = 10.0
		}
		fast = elapsed < thresh
	}
	s.Attempts = append(s.Attempts, Attempt{
		ConceptID:      conceptID,
		Correct:        correct,
		Fast:           fast,
		ElapsedSeconds: elapsed,
		Timestamp:      time.Now().UTC(),
	})
	s.totalCount[conceptID]++
	s.totalAsked++
	if correct {
		s.correctCount[conceptID]++
	}

	// Belief update for the concept itself (MA parity: slow correct answers
	// get diminished weight — "higher likelihood the student has not yet
	// learned the topic well enough to build on it").
	slowFactor := 1.0
	if correct && elapsed >= 0 {
		thresh := timeThresh
		if thresh <= 0 {
			thresh = 10.0
		}
		switch {
		case elapsed > 4*thresh:
			slowFactor = 0.25
		case elapsed > 2*thresh:
			slowFactor = 0.5
		}
	}
	b := s.beliefs[conceptID]
	if correct {
		b += 0.35 * slowFactor * (1 - b)
		if fast {
			b += 0.1 * (1 - b)
		}
	} else {
		b *= 0.55
	}
	if b < 0.05 {
		b = 0.05
	}
	if b > 0.95 {
		b = 0.95
	}
	s.beliefs[conceptID] = b

	// Evidence propagation: correct -> prerequisites more likely known
	// (diminished when slow); incorrect -> dependents more likely unknown.
	if correct {
		for _, pr := range e.dag.PrereqsOf(conceptID) {
			if pr.ID == conceptID {
				continue
			}
			s.beliefs[pr.ID] += 0.3 * slowFactor * (1 - s.beliefs[pr.ID]) * 0.5
		}
	} else {
		for _, dep := range e.dag.DependentsOf(conceptID) {
			if dep.ID == conceptID {
				continue
			}
			s.beliefs[dep.ID] *= (1 - 0.3*0.5)
		}
		// MA fall-back: a miss also strips "conditionally completed" credit
		// from barely-passed prerequisites — the frontier falls backwards
		// along the learning path immediately on struggle.
		for _, pr := range e.dag.PrereqsOf(conceptID) {
			if pr.ID == conceptID {
				continue
			}
			if b := s.beliefs[pr.ID]; b >= beliefThreshold && b < beliefThreshold+0.15 {
				s.beliefs[pr.ID] = b * 0.8
			}
		}
	}

	// Confidence: 2 probes settle a concept.
	s.confidence[conceptID] = math.Min(1.0, float64(s.probeCounts[conceptID])*0.5)

	// Supplemental: re-probe low-confidence cover concepts until 45 Q cap.
	s.doneSet[conceptID] = s.totalCount[conceptID] >= 2

	if s.shouldStop() {
		s.State = StateDone
	}
}

func (s *Session) shouldStop() bool {
	if s.totalAsked >= maxTotalQuestions {
		return true
	}
	cover := s.compressedCoverSafe()
	coverDone := true
	for _, c := range cover {
		if !s.doneSet[c.ID] {
			coverDone = false
			break
		}
	}
	if coverDone && s.totalAsked >= minTotalQuestions {
		return true
	}
	return false
}

// compressedCoverSafe mirrors compressedCover without an Engine receiver
// (session-only helper for stop conditions).
func (s *Session) compressedCoverSafe() []*concepts.Concept {
	if len(s.order) <= coverSize {
		return s.order
	}
	step := int(math.Ceil(float64(len(s.order)) / float64(coverSize)))
	var out []*concepts.Concept
	seen := make(map[string]bool)
	add := func(c *concepts.Concept) {
		if !seen[c.ID] {
			seen[c.ID] = true
			out = append(out, c)
		}
	}
	for i := 0; i < len(s.order); i += step {
		add(s.order[i])
	}
	for _, c := range s.order {
		if len(c.Prerequisites) == 0 {
			add(c)
		}
	}
	return out
}

func (e *Engine) IsComplete(s *Session) bool {
	s.Lock()
	defer s.Unlock()
	return s.State == StateDone
}

// Progress returns the backend truth for the progress bar (MA parity).
// CoverDone/CoverSize is monotonic (doneSet only grows); display the bar
// from cover fraction, never from the moving estimated_total denominator.
func (e *Engine) Progress(s *Session) Progress {
	s.Lock()
	defer s.Unlock()
	answered := s.totalAsked
	cover := s.compressedCoverSafeLocked()
	coverDone := 0
	for _, c := range cover {
		if s.doneSet[c.ID] {
			coverDone++
		}
	}
	remaining := len(cover) - coverDone
	est := answered + remaining
	if est < minTotalQuestions {
		est = minTotalQuestions
	}
	if est > maxTotalQuestions {
		est = maxTotalQuestions
	}
	if est < answered {
		est = answered
	}
	return Progress{
		Answered:       answered,
		EstimatedTotal: est,
		MinTotal:       minTotalQuestions,
		MaxTotal:       maxTotalQuestions,
		CoverDone:      coverDone,
		CoverSize:      len(cover),
		Done:           s.State == StateDone,
	}
}

// compressedCoverSafeLocked is compressedCoverSafe without locking
// (caller must hold s.Lock).
func (s *Session) compressedCoverSafeLocked() []*concepts.Concept {
	return s.compressedCoverSafe()
}

// FrontierEstimate returns the highest order index whose belief is known.
func (e *Engine) FrontierEstimate(s *Session) int {
	s.Lock()
	defer s.Unlock()
	idx := -1
	for i, c := range s.order {
		if s.beliefs[c.ID] >= beliefThreshold {
			idx = i
		}
	}
	return idx
}

// KnowledgeConfidence returns per-concept confidence 0-1 (copy, lock-safe).
func (e *Engine) KnowledgeConfidence(s *Session) map[string]float64 {
	s.Lock()
	defer s.Unlock()
	out := make(map[string]float64, len(s.confidence))
	for k, v := range s.confidence {
		out[k] = v
	}
	return out
}
