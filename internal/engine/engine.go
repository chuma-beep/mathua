package engine

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/leaderboard"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/scoring"
	"github.com/chuma-beep/mathua/internal/storage"
	xppolicy "github.com/chuma-beep/mathua/internal/xp"
)

type activeSession struct {
	conceptID      string
	conceptName    string
	expectedAnswer string
	explanation    string
	requiredStreak int
	timeThreshold  float64
	isReview       bool
	answered       bool
	attemptID      string
	questionText   string
	sessionReview  int
	sessionNew     int
	lastConceptID  string
	recentConcepts []string
	// PR 1.5 halt/re-attempt + negative XP
	consecutiveMisses  int
	halted             bool
	remedialQueue      []string
	remedialDifficulty float64
	// Difficulty the currently-served question was generated at. Zero means the question
	// was not generated at a known difficulty, and the attempt is written with a NULL
	// rather than a fabricated one.
	questionDifficulty float64
	rushCount          int
}

// ErrNoActiveQuestion is returned when a session has no unanswered question
// ready for grading — e.g. a duplicate/stale submission racing the current one.
var ErrNoActiveQuestion = fmt.Errorf("no active question")

// ErrUnknownConcept is returned when an answer references a concept_id absent
// from the DAG. Rejecting (instead of persisting a fallback progress row)
// keeps garbage IDs from farming XP and polluting progress/attempt tables.
var ErrUnknownConcept = fmt.Errorf("unknown concept")

// ErrNoStudyAnchor means the server has no record of the question being
// answered, so it cannot grade it. Distinct from a wrong answer: nothing is
// recorded — no attempt, no streak, no weakness, no XP — and the caller should
// re-serve the question. The anchor TTL is 24h, so the realistic trigger is a
// tab left open overnight.
var ErrNoStudyAnchor = fmt.Errorf("no server-side anchor for this question")

const (
	// serverSessionStudyExpected is the server_sessions kind for H1b
	// anti-cheat anchors. The name is historical and deliberately unchanged:
	// it is a persisted value referenced by the reset statements in both
	// stores, and renaming it would buy nothing but churn there.
	// studyAnchorTTL bounds the durable row: a practice question answered
	// within a day is normal; older rows are stale and treated as missing.
	serverSessionStudyExpected = "study_expected"
	studyAnchorTTL             = 24 * time.Hour
	// Ceiling on how many served-but-unanswered questions one student/concept
	// keeps anchored. Generators draw from a finite parameter space (nine
	// values per operand for arith.add.single, for one), so the set of live
	// question texts is small — but a long session with difficulty stepping and
	// variant changes could otherwise grow the blob without limit. At the cap the
	// batch just served is kept in full and older entries are dropped, which
	// degrades those questions to the 409 re-serve path rather than to a wrong
	// grade.
	studyAnchorMaxEntries = 64
)

// isFutureRFC3339 reports whether exp (RFC3339 UTC) is still in the future.
// Unparseable timestamps are treated as expired (fail closed).
func isFutureRFC3339(exp string) bool {
	t, err := time.Parse(time.RFC3339, exp)
	if err != nil {
		return false
	}
	return t.After(time.Now().UTC())
}

type Question struct {
	ConceptID   string          `json:"concept_id"`
	ConceptName string          `json:"concept_name"`
	Question    string          `json:"question"`
	IsReview    bool            `json:"is_review"`
	AttemptID   string          `json:"attempt_id,omitempty"`
	Lesson      *lessons.Lesson `json:"lesson,omitempty"`
	Diagram     string          `json:"diagram,omitempty"`
	// GradingType lets every answer host (study, review, quiz, diagnostic)
	// show the right input hint/keyboard. Empty means "unknown".
	GradingType string `json:"grading_type,omitempty"`
}

type AnswerResult struct {
	Correct        bool           `json:"correct"`
	Feedback       string         `json:"feedback"`
	NewStatus      mastery.Status `json:"new_status"`
	Explanation    string         `json:"explanation"`
	Streak         int            `json:"streak"`
	RequiredStreak int            `json:"required_streak"`
	XP             int            `json:"xp"`
	ExpectedAnswer string         `json:"expected_answer"`
	Halted         bool           `json:"halted,omitempty"`
	Remedial       []string       `json:"remedial,omitempty"`
	// Ungraded marks a grader infrastructure fault: the answer was NOT
	// evaluated and nothing was recorded. Clients should retry, not penalize.
	Ungraded bool `json:"ungraded,omitempty"`
	// Diagnosis names the mistake in one sentence, when it can be determined
	// with certainty. Empty most of the time, and never before grading.
	Diagnosis string `json:"diagnosis,omitempty"`
	// EvidenceScore is the authoritative mastery evidence score for this concept after
	// this answer, 0..1, from internal/mastery.BuildEvidence.
	//
	// It is sent so a client can render the learner's progress from the same number the
	// ladder decides on. `web/next-app/lib/progression.ts` used to compute its own score
	// from its own attempt history, which made a second definition of "how well are they
	// doing" that could not be reconciled with this one — it weighted difficulty per
	// attempt against this file's mean-of-correct-attempts, counted a wider population for
	// the variety term, and applied a flat 60s time cliff instead of the concept's own
	// threshold. Removing that copy is only possible because the number is here.
	//
	// Zero when the answer was ungraded, because nothing was recorded and there is no new
	// evidence. It is a *progress* reading, not a mastery verdict: the status is
	// `NewStatus`, and this score alone never means anything without the evidence floor
	// (mastery.MasteryEvidenceFloor) that the ladder also requires.
	EvidenceScore float64 `json:"evidence_score"`
	// EvidenceBand is EvidenceScore in plain language, from mastery.EvidenceBand.
	// Presentational only, and the words deliberately claim nothing about retention.
	EvidenceBand string `json:"evidence_band"`
}

type Engine struct {
	dag      *concepts.DAG
	repo     storage.Repository
	sched    *scheduler.Scheduler
	registry *generator.Registry
	gr       *grader.Router
	machine  *mastery.Machine
	scorer   *scoring.Updater
	lboard   *leaderboard.Computer
	diag     *diagnostic.Engine
	ll       *lessons.Loader
	planner  *planning.Planner

	mu         sync.Mutex
	sessions   map[string]*activeSession
	activePath map[string]map[string]bool
	// PR 1.5: study-path consecutive-miss tracker (studentID|conceptID → count)
	studyMisses map[string]int
	// studyRush counts sub-RushSeconds wrong answers per concept on the study/quiz path.
	// The practice path keeps its own count on the active session; both now apply the same
	// xp.RushSeconds/xp.RushAfter rule instead of two different ones.
	studyRush map[string]int
	// PR 1.5: stable per-student session id for study/quiz attempt FK.
	studySessions map[string]string
	// H1b: server-side anchors for the study seam — the answer each served
	// question must be graded against plus the explanation owed for it.
	studyAnchor map[string]string
	// Batch 1: immediate remedial queue from quiz misses (studentID → conceptIDs).
	// In-memory by design, same as diag/quiz sessions: a restart just means
	// a retake; weakness propagation (+0.2/miss) is the durable signal.
	quizRemedial map[string][]string
}

// QuizGateXP is the mastery-check interval. It delegates so there is one definition: the
// value was written out in three packages (engine, planning, scoring) with a comment in one
// of them admitting it was a copy.
const QuizGateXP = xppolicy.QuizGateXP

// maxQuizRemedial caps the per-student quiz remedial queue.
const maxQuizRemedial = 10

func New(repo storage.Repository, dag *concepts.DAG, reg *generator.Registry, ll *lessons.Loader, planner *planning.Planner) *Engine {
	return &Engine{
		dag:           dag,
		repo:          repo,
		sched:         scheduler.New(dag),
		registry:      reg,
		gr:            grader.NewRouter(),
		machine:       &mastery.Machine{},
		scorer:        scoring.NewUpdater(dag, repo),
		lboard:        leaderboard.NewComputer(repo),
		diag:          diagnostic.NewEngine(dag, reg),
		ll:            ll,
		planner:       planner,
		sessions:      make(map[string]*activeSession),
		activePath:    make(map[string]map[string]bool),
		studyMisses:   make(map[string]int),
		studyRush:     make(map[string]int),
		studySessions: make(map[string]string),
		studyAnchor:   make(map[string]string),
		quizRemedial:  make(map[string][]string),
	}
}

// Anchor is the server-side record of one served question: the answer it must
// be graded against, and the explanation the learner is owed for it. Both come
// from the same generator call over the same values, so storing them together
// is what lets a later submit address the exact instance the learner saw — the
// only place in the system that knows which instance a question text refers to.
type Anchor struct {
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
	// Difficulty the served question was produced at, or nil when it is genuinely unknown
	// (a curated question with no recorded difficulty).
	//
	// The anchor is the only durable record of what the learner was actually served, keyed
	// by question text, so it is the right place for this: `submitAnswerWithTask` already
	// reads the anchor to grade, and reading the difficulty from the same place means the
	// attempt is credited with the difficulty of the question that was served rather than
	// whatever the grading path happens to know.
	Difficulty *float64 `json:"difficulty,omitempty"`
}

// SetStudyAnchor stores a single anchor for a concept (question key "").
// Prefer SetStudyAnchorBatch for multi-question sets.
func (e *Engine) SetStudyAnchor(studentID, conceptID, expected string) {
	e.SetStudyAnchorBatch(studentID, conceptID, map[string]Anchor{"": {Answer: expected}})
}

// SetStudyAnchorBatch stores the anchor for every question in a served set,
// keyed by question text. Anchoring *every* question (not just the first) is
// what stops a later question from being graded against the first one's answer;
// carrying the explanation alongside is what lets the submit return the
// solution for the question actually answered. Best-effort durable
// write-through for restarts.
//
// The blob MERGES rather than replaces, and that is load-bearing. The Learn
// client keeps a buffer of served-but-unanswered questions across several
// fetches, and serves a new batch whenever the buffer runs low — which happens
// while the learner is still reading the question in front of them. Replacing
// meant each new batch erased the previous one, so the next submit looked up a
// question the server had genuinely forgotten, got ErrNoStudyAnchor, and came
// back 409. The client treats 409 as "that question expired" and swaps in a
// fresh question, which is what looked like Learn skipping questions the learner
// had just answered. Merging keeps every question still in flight gradeable.
//
// Incoming entries win on identical text. That is safe rather than lossy
// because the generators are deterministic per (concept, seed, difficulty), so
// the same text always implies the same answer and explanation.
//
// The server_sessions kind deliberately stays "study_expected": the stored
// value is an opaque JSON blob, and renaming the kind would touch every reset
// statement in both stores plus their tests for no behavioural gain.
func (e *Engine) SetStudyAnchorBatch(studentID, conceptID string, byQuestion map[string]Anchor) {
	if len(byQuestion) == 0 {
		return
	}
	key := studentID + "|" + conceptID

	// Read-merge-write under one lock. Two serves racing could otherwise
	// interleave and drop a batch even though each merge was individually
	// correct.
	e.mu.Lock()
	if e.studyAnchor == nil {
		e.studyAnchor = make(map[string]string)
	}
	merged := make(map[string]Anchor, len(byQuestion))
	// Only decode a blob that actually exists. decodeStudyAnchor's final
	// fallback for unparseable input is a single ""-keyed entry, so handing it
	// an empty string would plant a phantom empty anchor -- and studyAnchorFor
	// falls back to the "" key, which would then report every unknown question
	// as "found" (with an empty answer).
	if existing := e.studyAnchor[key]; existing != "" {
		merged = decodeStudyAnchor(existing)
	}
	for question, anchor := range byQuestion {
		merged[question] = anchor
	}
	if len(merged) > studyAnchorMaxEntries {
		// Keep the batch just served in full: every question the learner is
		// looking at right now is in it. Older entries are evicted in map order
		// (deliberately not sorted — any order is fine and this stays O(1) per
		// eviction). An evicted question degrades to the existing 409 recovery,
		// which is the correct response for an anchor that is genuinely gone.
		for question := range merged {
			if _, keep := byQuestion[question]; keep {
				continue
			}
			delete(merged, question)
			if len(merged) <= studyAnchorMaxEntries {
				break
			}
		}
	}
	blob, err := json.Marshal(merged)
	if err != nil {
		e.mu.Unlock()
		log.Printf("warning: marshal study anchor %s: %v", key, err)
		return
	}
	e.studyAnchor[key] = string(blob)
	e.mu.Unlock()

	if e.repo != nil {
		exp := time.Now().UTC().Add(studyAnchorTTL).Format(time.RFC3339)
		if err := e.repo.UpsertServerSession(serverSessionStudyExpected, key, string(blob), exp); err != nil {
			log.Printf("warning: persist study anchor %s: %v", key, err)
		}
	}
}

// studyAnchorFor returns the anchor for a specific served question. Falls back
// to the single-anchor ("") form for callers that store one, and to the durable
// row after a restart.
func (e *Engine) studyAnchorFor(studentID, conceptID, question string) (Anchor, bool) {
	key := studentID + "|" + conceptID
	e.mu.Lock()
	blob, ok := e.studyAnchor[key]
	e.mu.Unlock()
	if !ok && e.repo != nil {
		if v, exp, found, err := e.repo.GetServerSession(serverSessionStudyExpected, key); err == nil && found {
			if exp == "" || isFutureRFC3339(exp) {
				blob, ok = v, true
			}
		} else if err != nil {
			log.Printf("warning: read study anchor %s: %v", key, err)
		}
	}
	if !ok {
		return Anchor{}, false
	}
	byQuestion := decodeStudyAnchor(blob)
	if question != "" {
		if a, found := byQuestion[question]; found {
			return a, true
		}
	}
	if a, found := byQuestion[""]; found {
		return a, true
	}
	return Anchor{}, false
}

// decodeStudyAnchor reads an anchor blob, tolerating the shape written before
// the explanation was carried: {"q": "answer"}. A session opened before a
// deploy therefore still grades; it simply has no explanation to return, and
// the caller falls back. The final bare-string return is a defensive
// last resort — every writer marshals a map, so no stored blob has that shape.
func decodeStudyAnchor(blob string) map[string]Anchor {
	var byQuestion map[string]Anchor
	if err := json.Unmarshal([]byte(blob), &byQuestion); err == nil {
		return byQuestion
	}
	var legacy map[string]string
	if err := json.Unmarshal([]byte(blob), &legacy); err == nil {
		byQuestion = make(map[string]Anchor, len(legacy))
		for question, answer := range legacy {
			byQuestion[question] = Anchor{Answer: answer}
		}
		return byQuestion
	}
	return map[string]Anchor{"": {Answer: blob}}
}

func (e *Engine) CreateStudent(name string) (*storage.Student, error) {
	return e.repo.CreateStudent(name)
}

func (e *Engine) GetStudent(id string) (*storage.Student, error) {
	return e.repo.GetStudent(id)
}

func (e *Engine) SetActiveCourse(studentID, courseID string) (*planning.Path, error) {
	if e.planner == nil {
		return nil, fmt.Errorf("no planner configured")
	}
	path, err := e.planner.PathForCourse(courseID)
	if err != nil {
		return nil, err
	}
	if err := e.repo.SetCourseID(studentID, courseID); err != nil {
		return nil, fmt.Errorf("save course: %w", err)
	}
	set := make(map[string]bool, len(path.Concepts))
	for _, c := range path.Concepts {
		set[c.ID] = true
	}
	e.mu.Lock()
	e.activePath[studentID] = set
	e.mu.Unlock()
	return path, nil
}

func (e *Engine) GetPlanner() *planning.Planner { return e.planner }
func (e *Engine) GetGrader() *grader.Router     { return e.gr }
func (e *Engine) PlannerCourses() []*planning.Course {
	if e.planner == nil {
		return nil
	}
	return e.planner.Courses()
}

// CourseStatus is one course with the student's accreditation-track progress.
type CourseStatus struct {
	*planning.Course
	Progress *planning.CourseProgress `json:"progress"`
}

// CourseCatalog returns all courses with per-student progress + estimates.
func (e *Engine) CourseCatalog(studentID string) ([]CourseStatus, error) {
	if e.planner == nil {
		return nil, nil
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, err
	}
	goal := 10
	if st, err := e.repo.GetStudent(studentID); err == nil && st != nil && st.DailyXPGoal > 0 {
		goal = st.DailyXPGoal
	}
	out := make([]CourseStatus, 0, len(e.planner.Courses()))
	for _, c := range e.planner.Courses() {
		cp, err := e.planner.ProgressForCourse(c, progress, goal)
		if err != nil {
			return nil, err
		}
		out = append(out, CourseStatus{Course: c, Progress: cp})
	}
	return out, nil
}

func (e *Engine) ActivePath(studentID string) map[string]bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	path := e.activePath[studentID]
	if path == nil {
		return nil
	}
	cp := make(map[string]bool, len(path))
	for k, v := range path {
		cp[k] = v
	}
	return cp
}

// computeDifficulty returns a difficulty score (0.3–1.0) based on the student's
// performance for the given concept. Weak/struggling students get easier questions.
func (e *Engine) computeDifficulty(studentID, conceptID string) float64 {
	return difficultyFromWeakness(e.WeaknessMap(studentID), conceptID)
}

// difficultyFromWeakness is the scan-free core of computeDifficulty: callers
// that already hold a weakness map (built from a single progress fetch) use
// this instead of triggering another DB round trip.
func difficultyFromWeakness(m map[string]float64, conceptID string) float64 {
	weakness := 0.5
	if m != nil {
		if w, ok := m[conceptID]; ok {
			weakness = w
		}
	}
	// Higher weakness → lower difficulty (easier questions)
	d := 0.3 + (1-weakness)*0.5
	if d < 0.3 {
		d = 0.3
	}
	if d > 1.0 {
		d = 1.0
	}
	return d
}

// accommodatedThreshold applies the student's accommodations
// (settings.accommodations.extra_time 0-2.0 → 1.0-2.5x) to a time threshold.
func (e *Engine) accommodatedThreshold(studentID string, base float64) float64 {
	if base <= 0 {
		return base
	}
	if e.repo == nil {
		return base
	}
	raw, err := e.repo.GetSettings(studentID)
	if err != nil || raw == "" || raw == "{}" {
		return base
	}
	var cfg struct {
		Accommodations struct {
			ExtraTime float64 `json:"extra_time"`
		} `json:"accommodations"`
	}
	if json.Unmarshal([]byte(raw), &cfg) != nil || cfg.Accommodations.ExtraTime <= 0 {
		return base
	}
	m := 1.0 + cfg.Accommodations.ExtraTime
	if m > 2.5 {
		m = 2.5
	}
	return base * m
}

func (e *Engine) NextQuestion(sessionID, studentID string) (*Question, error) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, fmt.Errorf("load progress: %w", err)
	}
	snapshots := e.conceptSnapshotsFrom(studentID, progress)

	e.mu.Lock()
	defer e.mu.Unlock()

	// Filter to active course path if set
	if path := e.activePath[studentID]; len(path) > 0 {
		for cid := range snapshots {
			if !path[cid] {
				delete(snapshots, cid)
			}
		}
	}

	as := e.sessions[sessionID]
	if as == nil {
		as = &activeSession{}
	}

	// PR 1.5 halt: after 2 consecutive misses, serve the remedial queue first
	// (KeyPrerequisites then the concept at an easier variant) before resuming.
	if as.halted && len(as.remedialQueue) > 0 {
		cid := as.remedialQueue[0]
		as.remedialQueue = as.remedialQueue[1:]
		if len(as.remedialQueue) == 0 {
			as.halted = false
		}
		as.answered = false
		as.attemptID = newAttemptID()
		as.conceptID = cid
		c := e.dag.Concept(cid)
		if c == nil {
			return nil, nil
		}
		as.conceptName = c.Label
		as.requiredStreak = c.MasteryThreshold.Streak
		as.timeThreshold = e.accommodatedThreshold(studentID, c.MasteryThreshold.AvgTimeSeconds)
		diff := 0.3
		if as.remedialDifficulty > 0 {
			diff = as.remedialDifficulty
		}
		ctx := generator.GeneratorContext{Difficulty: diff, Seed: hashSeed(studentID + "|" + cid + "|" + as.attemptID)}
		prob, err := e.registry.GenerateContext(cid, ctx)
		if err != nil {
			return nil, fmt.Errorf("generate remedial problem: %w", err)
		}
		as.expectedAnswer = prob.Answer
		as.explanation = prob.Explanation
		as.isReview = false
		as.questionText = prob.Question
		as.questionDifficulty = diff
		e.sessions[sessionID] = as
		e.persistActiveSession(sessionID, studentID, as, as.questionText)
		var lesson *lessons.Lesson
		if e.ll != nil {
			lesson = e.ll.Lesson(cid)
		}
		return &Question{
			ConceptID:   cid,
			ConceptName: c.Label,
			Question:    prob.Question,
			GradingType: e.gradingTypeFor(cid),
			IsReview:    false,
			AttemptID:   as.attemptID,
			Lesson:      lesson,
		}, nil
	}

	// PR 1.4: recent-window (last 2) for interleaving + non-interference.
	recent := as.recentConcepts
	if as.lastConceptID != "" {
		recent = append([]string{as.lastConceptID}, recent...)
		if len(recent) > 2 {
			recent = recent[:2]
		}
	}
	// Single progress fetch threads through scheduler + difficulty: no
	// extra DB scans on the hot path.
	weakness := e.weaknessMapFromProgress(progress)
	cands := e.sched.NextSmart(snapshots, recent, as.sessionReview, as.sessionNew, weakness)
	if len(cands) == 0 {
		return nil, nil
	}
	next := cands[0]
	difficulty := difficultyFromWeakness(weakness, next.Concept.ID)
	prevQuestion := as.questionText
	attemptID := newAttemptID()
	seedBase := hashSeed(studentID + "|" + next.Concept.ID + "|" + attemptID)
	ctx := generator.GeneratorContext{Difficulty: difficulty, Seed: seedBase}
	prob, err := e.registry.GenerateContext(next.Concept.ID, ctx)
	if err != nil {
		return nil, fmt.Errorf("generate problem: %w", err)
	}
	// Dedup verbatim: if same concept and question text collides with previous, re-roll with incremented seed.
	for i := 0; i < 3 && prevQuestion != "" && prob.Question == prevQuestion; i++ {
		ctx.Seed = seedBase + int64(i+1)
		if p2, err2 := e.registry.GenerateContext(next.Concept.ID, ctx); err2 == nil {
			prob = p2
		}
	}
	as.conceptID = next.Concept.ID
	as.conceptName = next.Concept.Label
	as.expectedAnswer = prob.Answer
	as.explanation = prob.Explanation
	as.requiredStreak = next.Concept.MasteryThreshold.Streak
	as.timeThreshold = e.accommodatedThreshold(studentID, next.Concept.MasteryThreshold.AvgTimeSeconds)
	as.isReview = next.IsReview
	as.answered = false
	as.attemptID = attemptID
	as.questionText = prob.Question
	as.questionDifficulty = difficulty
	as.recentConcepts = append(as.recentConcepts, next.Concept.ID)
	if len(as.recentConcepts) > 3 {
		as.recentConcepts = as.recentConcepts[len(as.recentConcepts)-3:]
	}
	e.sessions[sessionID] = as
	e.persistActiveSession(sessionID, studentID, as, as.questionText)

	var lesson *lessons.Lesson
	if e.ll != nil {
		lesson = e.ll.Lesson(next.Concept.ID)
	}

	diagram := diagramForConcept(next.Concept.ID)

	return &Question{
		ConceptID:   next.Concept.ID,
		ConceptName: next.Concept.Label,
		Question:    prob.Question,
		GradingType: e.gradingTypeFor(next.Concept.ID),
		IsReview:    next.IsReview,
		AttemptID:   as.attemptID,
		Lesson:      lesson,
		Diagram:     diagram,
	}, nil
}

// NextReviewQuestion only serves concepts due for spaced-repetition review.
func (e *Engine) NextReviewQuestion(sessionID, studentID string) (*Question, error) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, fmt.Errorf("load progress: %w", err)
	}
	snapshots := e.conceptSnapshotsFrom(studentID, progress)

	e.mu.Lock()
	defer e.mu.Unlock()

	as := e.sessions[sessionID]
	if as == nil {
		as = &activeSession{}
	}
	next := e.sched.NextReview(snapshots, as.lastConceptID)
	if next == nil {
		return nil, nil
	}
	difficulty := difficultyFromWeakness(e.weaknessMapFromProgress(progress), next.Concept.ID)
	prevQuestion := as.questionText
	attemptID := newAttemptID()
	seedBase := hashSeed(studentID + "|" + next.Concept.ID + "|" + attemptID)
	ctx := generator.GeneratorContext{Difficulty: difficulty, Seed: seedBase}
	prob, err := e.registry.GenerateContext(next.Concept.ID, ctx)
	if err != nil {
		return nil, fmt.Errorf("generate problem: %w", err)
	}
	for i := 0; i < 3 && prevQuestion != "" && prob.Question == prevQuestion; i++ {
		ctx.Seed = seedBase + int64(i+1)
		if p2, err2 := e.registry.GenerateContext(next.Concept.ID, ctx); err2 == nil {
			prob = p2
		}
	}
	as.conceptID = next.Concept.ID
	as.conceptName = next.Concept.Label
	as.expectedAnswer = prob.Answer
	as.explanation = prob.Explanation
	as.requiredStreak = next.Concept.MasteryThreshold.Streak
	as.timeThreshold = e.accommodatedThreshold(studentID, next.Concept.MasteryThreshold.AvgTimeSeconds)
	as.isReview = true
	as.answered = false
	as.attemptID = attemptID
	as.questionText = prob.Question
	as.questionDifficulty = difficulty
	e.sessions[sessionID] = as
	e.persistActiveSession(sessionID, studentID, as, as.questionText)

	var lesson *lessons.Lesson
	if e.ll != nil {
		lesson = e.ll.Lesson(next.Concept.ID)
	}
	diagram := diagramForConcept(next.Concept.ID)

	return &Question{
		ConceptID:   next.Concept.ID,
		ConceptName: next.Concept.Label,
		Question:    prob.Question,
		GradingType: e.gradingTypeFor(next.Concept.ID),
		IsReview:    true,
		AttemptID:   as.attemptID,
		Lesson:      lesson,
		Diagram:     diagram,
	}, nil
}

var conceptDiagrams = map[string]string{
	// Integrals (dual-coding worked examples, improve.md:39)
	"calc.integral.definite":          "/diagrams/algebrica/definite-integrals-1.svg",
	"calc.integral.ftc":               "/diagrams/algebrica/fundamental-theorem-of-calculus-1.svg",
	"calc.integral.area_between":      "/diagrams/algebrica/finding-areas-by-integration-1.svg",
	"calc.integral.volume":            "/diagrams/algebrica/finding-areas-by-integration-2.svg",
	"calc.integral.improper":          "/diagrams/algebrica/improper-integrals-1.svg",
	"calc.integral.numerical":         "/diagrams/algebrica/improper-integrals-2.svg",
	"calc.integral.riemann_criteria":  "/diagrams/algebrica/riemann-integrability-criteria-2.svg",
	"calc.integral.arc_length":        "/diagrams/algebrica/arc-length-of-a-curve-1.svg",
	"calc.integral.trig_substitution": "/diagrams/algebrica/trigonometric-substitution-for-integrals-1.svg",

	// Limits
	"calc.limit.continuity": "/diagrams/algebrica/riemann-integrability-criteria-1.svg",
	"calc.limit.supremum":   "/diagrams/algebrica/supremum-and-infimum-1.svg",

	// Equations / quadratics
	"alg.quad.solve_factor":    "/diagrams/algebrica/quadratic-equations.svg",
	"alg.quad.formula":         "/diagrams/algebrica/quadratic-equations.svg",
	"alg.quad.quadratic":       "/diagrams/algebrica/quadratic-equations.svg",
	"alg.quad.incomplete":      "/diagrams/algebrica/incomplete-quadratic-equations.svg",
	"alg.quad.complete_square": "/diagrams/algebrica/completing-square.svg",

	// Linear / polynomials
	"alg.linear.graph":           "/diagrams/algebrica/linear-equation-graph.svg",
	"alg.linear.slope":           "/diagrams/algebrica/linear-equation-graph.svg",
	"alg.linear.slope_intercept": "/diagrams/algebrica/linear-equation-graph.svg",
	"alg.poly.roots":             "/diagrams/algebrica/polynomial-roots-graph.svg",

	// Number lines / sets
	"arith.neg.abs_value":     "/diagrams/algebrica/number-line-absolute-value.svg",
	"alg.ineq.absolute_value": "/diagrams/algebrica/number-line-absolute-value.svg",
	"alg.ineq.interval":       "/diagrams/algebrica/number-line-intervals.svg",
	"prealg.real.concept":     "/diagrams/algebrica/number-line-real.svg",
	"prealg.types":            "/diagrams/algebrica/number-types-venn.svg",

	// Complex
	"complex.basics.concept": "/diagrams/algebrica/complex-plane.svg",
	"complex.adv.polar":      "/diagrams/algebrica/complex-plane.svg",

	// Trigonometry
	"trig.hyperbolic.sinh_cosh":   "/diagrams/algebrica/hyperbolic-functions.svg",
	"trig.hyperbolic.tanh_coth":   "/diagrams/algebrica/hyperbolic-functions.svg",
	"trig.adv.inverse":            "/diagrams/algebrica/inverse-trig-graphs.svg",
	"trig.adv.arctan":             "/diagrams/algebrica/inverse-trig-graphs.svg",
	"trig.adv.law_cosines":        "/diagrams/algebrica/law-of-cosines.svg",
	"trig.adv.law_sines":          "/diagrams/algebrica/law-of-sines.svg",
	"geo.triangle.pythagorean":    "/diagrams/algebrica/pythagorean-theorem.svg",
	"trig.ident.pythagorean":      "/diagrams/algebrica/pythagorean-theorem.svg",
	"trig.basics.reference_angle": "/diagrams/algebrica/reference-angles.svg",
	"trig.basics.right_triangle":  "/diagrams/algebrica/right-triangle-trig.svg",
	"trig.basics.sin_cos_def":     "/diagrams/algebrica/right-triangle-unit-circle.svg",
	"trig.basics.radians":         "/diagrams/algebrica/unit-circle-sine-cosine.svg",
	"trig.basics.tan_def":         "/diagrams/algebrica/unit-circle-tangent.svg",
	"trig.basics.reciprocal":      "/diagrams/algebrica/sec-csc-cot-graphs.svg",
	"trig.basics.unit_circle":     "/diagrams/algebrica/unit-circle-labeled.svg",
	"trig.basics.special_angles":  "/diagrams/algebrica/unit-circle-labeled.svg",
	"trig.graph.sin":              "/diagrams/algebrica/sine-cosine-graph.svg",
	"trig.graph.cos":              "/diagrams/algebrica/sine-cosine-graph.svg",
	"trig.graph.period":           "/diagrams/algebrica/sine-cosine-graph.svg",
	"trig.func.sine":              "/diagrams/algebrica/sine-and-cosine-3.svg",
	"trig.func.cosine":            "/diagrams/algebrica/sine-and-cosine-4.svg",
	"trig.func.tangent":           "/diagrams/algebrica/tangent-and-cotangent-4.svg",
	"trig.func.cotangent":         "/diagrams/algebrica/tangent-and-cotangent-5.svg",
	"trig.func.secant":            "/diagrams/algebrica/secant-and-cosecant-3.svg",
	"trig.func.cosecant":          "/diagrams/algebrica/secant-and-cosecant-4.svg",
	"trig.func.arcsine":           "/diagrams/algebrica/arcsine-function-1.svg",
	"trig.func.arccosine":         "/diagrams/algebrica/arccosine-function-1.svg",
	"trig.func.arctangent":        "/diagrams/algebrica/arctangent-function-1.svg",
	"trig.func.arccotangent":      "/diagrams/algebrica/arccotangent-function-1.svg",

	// Combinatorics
	"discrete.combinatorics.pascal":           "/diagrams/algebrica/pascals-triangle.svg",
	"discrete.combinatorics.binomial_theorem": "/diagrams/algebrica/pascals-triangle.svg",
	"precalc.binomial_theorem":                "/diagrams/algebrica/pascals-triangle.svg",

	// Vectors
	"linalg.vector.add":     "/diagrams/algebrica/vector-addition.svg",
	"linalg.vector.dot":     "/diagrams/algebrica/vector-addition.svg",
	"linalg.vector.concept": "/diagrams/algebrica/vector-arrow.svg",

	// Conics / sequences / stats (PNG assets)
	"alg.conic.circle":    "/diagrams/algebrica/conic-circle-1.png",
	"alg.conic.ellipse":   "/diagrams/algebrica/ellipse-1.svg",
	"alg.conic.hyperbola": "/diagrams/algebrica/hyperbola-1.svg",
	"alg.conic.parabola":  "/diagrams/algebrica/parabola-1.svg",
	"alg.seq.arithmetic":  "/diagrams/algebrica/arithmetic-sequence.png",
	"alg.seq.geometric":   "/diagrams/algebrica/geometri-sequence-1.png",
	"stat.dist.normal":    "/diagrams/algebrica/normal-distribution-1.svg",
	"stat.func.gaussian":  "/diagrams/algebrica/gaussian-function-1.svg",

	// Functions (G3)
	"alg.func.absolute_value":    "/diagrams/algebrica/absolute-value-1.svg",
	"alg.exp.concept":            "/diagrams/algebrica/exponential-function-1.svg",
	"alg.func.even_odd":          "/diagrams/algebrica/even-odd-functions-1.png",
	"alg.func.monotonicity":      "/diagrams/algebrica/increasing-decreasing-function.png",
	"alg.func.sigmoid":           "/diagrams/algebrica/sigmoid-function.png",
	"alg.func.sign":              "/diagrams/algebrica/sign-function-1.svg",
	"alg.func.composite":         "/diagrams/algebrica/composite-functions-1-1.png",
	"alg.func.injectivity":       "/diagrams/algebrica/injective-surjective-and-bijective-functions-1.svg",
	"alg.func.power":             "/diagrams/algebrica/power-function-1.svg",
	"alg.func.irrational":        "/diagrams/algebrica/irrational-functions-1.svg",
	"precalc.func.floor_ceiling": "/diagrams/algebrica/floor-and-ceiling-functions-1.svg",
	"precalc.func.heaviside":     "/diagrams/algebrica/heaviside-function-1.svg",
	"alg.ineq.quadratic":         "/diagrams/algebrica/quadratic-inequalities-1-1.png",
	"alg.ineq.sign_analysis":     "/diagrams/algebrica/sign-analysis-1.png",
	"alg.quad.complex":           "/diagrams/algebrica/quad-eq-complex-roots.png",
	"alg.quad.parametric":        "/diagrams/algebrica/quadratic-equations-params-6.png",
	"trig.eq.basic":              "/diagrams/algebrica/trigonometric-equations-1.png",
	"trig.ineq.basic":            "/diagrams/algebrica/trigonometric-inequalities-1.svg",

	// Calculus (G3)
	"calc.limit.asymptotes":            "/diagrams/algebrica/asymptotes-1.svg",
	"calc.limit.squeeze":               "/diagrams/algebrica/squeeze-theorem.png",
	"calc.limit.uniform_continuity":    "/diagrams/algebrica/uniform-continuity-1.svg",
	"calc.limit.big_o":                 "/diagrams/algebrica/big-o-notation-1.svg",
	"calc.limit.weierstrass":           "/diagrams/algebrica/weierstrass-theorem-1.svg",
	"calc.deriv.difference_quotient":   "/diagrams/algebrica/difference-quotient-2.png",
	"calc.deriv.partial":               "/diagrams/algebrica/partial-derivatives-1.png",
	"calc.deriv.non_differentiability": "/diagrams/algebrica/non-differentiable-points-1.png",
	"calc.deriv.convexity":             "/diagrams/algebrica/convexity-1-1.png",
	"calc.deriv.rolle":                 "/diagrams/algebrica/rolle-theorem-1-1.png",
	"calc.deriv.applications":          "/diagrams/algebrica/velocity-1-1.png",
	"calc.seq.convergence":             "/diagrams/algebrica/sequences-conv-1.png",
	"calc.seq.cauchy":                  "/diagrams/algebrica/cauchy-sequence-1.svg",
	"calc.seq.euler":                   "/diagrams/algebrica/euler-number-limit-sequence-1.svg",
	"calc.series.harmonic":             "/diagrams/algebrica/harmonic-series-1-2.png",
	"calc.series.function_series":      "/diagrams/algebrica/sequence-functions-1.png",
	"calc.series.cauchy_criterion":     "/diagrams/algebrica/series-cauchy-1.png",

	// Stats / vectors / misc (G3)
	"stat.dist.student_t":             "/diagrams/algebrica/student-t-distribution.png",
	"stat.dist.uniform":               "/diagrams/algebrica/uniform-distribution.png",
	"stat.dist.beta":                  "/diagrams/algebrica/beta-distribution-1.svg",
	"stat.dist.gamma":                 "/diagrams/algebrica/gamma-distribution.png",
	"stat.dist.exponential":           "/diagrams/algebrica/exponential-distribution-1.png",
	"stat.dist.chi_square":            "/diagrams/algebrica/chi-squared-distribution.png",
	"stat.infer.confidence":           "/diagrams/algebrica/confidence-intervals.png",
	"stat.prob.continuous_rv":         "/diagrams/algebrica/continuous-random-vars-1.png",
	"linalg.vector.cosine_similarity": "/diagrams/algebrica/cosine-similarity.png",
	"linalg.vector.parametric":        "/diagrams/algebrica/vector-equation-line-1.png",
	"ml.backpropagation":              "/diagrams/algebrica/neural-network.png",
	"geo.coord.polar":                 "/diagrams/algebrica/polar-coordinates-1-1.png",
	"precalc.polar.coordinates":       "/diagrams/algebrica/polar-coordinates-1-1.png",
	"precalc.polar.graph":             "/diagrams/algebrica/polar-coordinates-2-1.png",
	"ode.basics.concept":              "/diagrams/algebrica/differential-equations-1-1.png",
	// Batch C1: pool remap — verified against rendered image content,
	// one asset per concept; wrong mappings dropped, not replaced.
	"calc.limit.epsilon_delta":   "/diagrams/algebrica/limits-1.svg",
	"calc.limit.infinite":        "/diagrams/algebrica/limits-2.svg",
	"calc.limit.infinity":        "/diagrams/algebrica/limits-3.svg",
	"precalc.conic.hyperbola":    "/diagrams/algebrica/hyperbola-2.svg",
	"ml.forward_pass":            "/diagrams/algebrica/neural-network-2.png",
	"calc.deriv.global_extrema":  "/diagrams/algebrica/max-min-1.png",
	"calc.deriv.fermat":          "/diagrams/algebrica/fermat-1.png",
	"precalc.parametric.graph":   "/diagrams/algebrica/velocity-2.svg",
	"calc.limit.discontinuity":   "/diagrams/algebrica/Heaviside-function-1.png",
	"calc.limit.concept":         "/diagrams/algebrica/continuous-functions.png",
	"calc.deriv.concept":         "/diagrams/algebrica/derivatives-1.svg",
	"ode.first_order.separable":  "/diagrams/algebrica/differential-equations-2.png",
	"precalc.conic.ellipse":      "/diagrams/algebrica/ellipse-2.svg",
	"geo.coord.distance":         "/diagrams/algebrica/euclidean-distance.png",
	"alg.systems.concept":        "/diagrams/algebrica/lines-5.png",
	"alg.func.concept":           "/diagrams/algebrica/functions-5.png",
	"precalc.seq.geometric":      "/diagrams/algebrica/geometri-sequence-2-1.png",
	"alg.ineq.compound":          "/diagrams/algebrica/inequalities-1.png",
	"calc.series.integral_test":  "/diagrams/algebrica/integral-test-series-1.png",
	"alg.func.inverse":           "/diagrams/algebrica/inverse-function-1.svg",
	"alg.ineq.irrational":        "/diagrams/algebrica/irr-ineq-1.png",
	"calc.deriv.mvt":             "/diagrams/algebrica/lagrange-theoreme-4.png",
	"alg.linear.standard_form":   "/diagrams/algebrica/line-1.png",
	"calc.limit.little_o":        "/diagrams/algebrica/little-o-1.svg",
	"alg.log.concept":            "/diagrams/algebrica/logharithm-6.png",
	"calc.deriv.critical_points": "/diagrams/algebrica/max-min-3-1.png",
	"calc.deriv.inflection":      "/diagrams/algebrica/maximum-minimum-7.png",
	"stat.dist.z_table":          "/diagrams/algebrica/normal-distribution-standard-1.png",
	"alg.func.quad":              "/diagrams/algebrica/parabola-2.svg",
	"precalc.conic.parabola":     "/diagrams/algebrica/parabola-5-1.png",
	"calc.deriv.tangent_line":    "/diagrams/algebrica/parabola-7.svg",
	"alg.quad.discriminant":      "/diagrams/algebrica/quadratic-inequalities-2.svg",
	"geo.circle.parts":           "/diagrams/algebrica/circumference-5.svg",
	// Batch C5a: hand-authored arithmetic diagrams (add/sub/place/round).
	"arith.add.single":      "/diagrams/arithmetic/add-single.svg",
	"arith.add.double":      "/diagrams/arithmetic/add-double.svg",
	"arith.add.carry":       "/diagrams/arithmetic/add-carry.svg",
	"arith.add.triple":      "/diagrams/arithmetic/add-triple.svg",
	"arith.add.word":        "/diagrams/arithmetic/add-word.svg",
	"arith.sub.single":      "/diagrams/arithmetic/sub-single.svg",
	"arith.sub.double":      "/diagrams/arithmetic/sub-double.svg",
	"arith.sub.borrow":      "/diagrams/arithmetic/sub-borrow.svg",
	"arith.sub.word":        "/diagrams/arithmetic/sub-word.svg",
	"arith.place.tens":      "/diagrams/arithmetic/place-tens.svg",
	"arith.place.hundreds":  "/diagrams/arithmetic/place-hundreds.svg",
	"arith.place.thousands": "/diagrams/arithmetic/place-thousands.svg",
	"arith.round.tens":      "/diagrams/arithmetic/round-tens.svg",
	"arith.round.hundreds":  "/diagrams/arithmetic/round-hundreds.svg",
	"arith.round.thousands": "/diagrams/arithmetic/round-thousands.svg",
	// Batch C5b: hand-authored arithmetic diagrams (mult/div/factors).
	"arith.mult.concept":      "/diagrams/arithmetic/mult-concept.svg",
	"arith.mult.2_5_10":       "/diagrams/arithmetic/mult-2-5-10.svg",
	"arith.mult.tables":       "/diagrams/arithmetic/mult-tables.svg",
	"arith.mult.double":       "/diagrams/arithmetic/mult-double.svg",
	"arith.mult.triple":       "/diagrams/arithmetic/mult-triple.svg",
	"arith.mult.word":         "/diagrams/arithmetic/mult-word.svg",
	"arith.div.concept":       "/diagrams/arithmetic/div-concept.svg",
	"arith.div.basic":         "/diagrams/arithmetic/div-basic.svg",
	"arith.div.remainder":     "/diagrams/arithmetic/div-remainder.svg",
	"arith.div.long":          "/diagrams/arithmetic/div-long.svg",
	"arith.div.word":          "/diagrams/arithmetic/div-word.svg",
	"arith.factor.find":       "/diagrams/arithmetic/factor-find.svg",
	"arith.factor.prime":      "/diagrams/arithmetic/factor-prime.svg",
	"arith.factor.prime_fact": "/diagrams/arithmetic/factor-prime-fact.svg",
	"arith.factor.gcf":        "/diagrams/arithmetic/factor-gcf.svg",
	"arith.factor.lcm":        "/diagrams/arithmetic/factor-lcm.svg",
	"arith.factor.composite":  "/diagrams/arithmetic/factor-composite.svg",
	// Batch C5c: hand-authored arithmetic diagrams (exp/sqrt/neg/order/dec).
	"arith.exp.concept":       "/diagrams/arithmetic/exp-concept.svg",
	"arith.exp.evaluate":      "/diagrams/arithmetic/exp-evaluate.svg",
	"arith.exp.product_rule":  "/diagrams/arithmetic/exp-product-rule.svg",
	"arith.exp.quotient_rule": "/diagrams/arithmetic/exp-quotient-rule.svg",
	"arith.exp.power_rule":    "/diagrams/arithmetic/exp-power-rule.svg",
	"arith.sqrt.perfect":      "/diagrams/arithmetic/sqrt-perfect.svg",
	"arith.sqrt.simplify":     "/diagrams/arithmetic/sqrt-simplify.svg",
	"arith.neg.number_line":   "/diagrams/arithmetic/neg-number-line.svg",
	"arith.neg.add_sub":       "/diagrams/arithmetic/neg-add-sub.svg",
	"arith.neg.mult_div":      "/diagrams/arithmetic/neg-mult-div.svg",
	"arith.order_ops.basic":   "/diagrams/arithmetic/order-basic.svg",
	"arith.order_ops.full":    "/diagrams/arithmetic/order-full.svg",
	"arith.order_ops.nested":  "/diagrams/arithmetic/order-nested.svg",
	"arith.dec.intro":         "/diagrams/arithmetic/dec-intro.svg",
}

// DiagramFor exposes the dual-coding diagram for a concept (exported for API).
func (e *Engine) DiagramFor(conceptID string) string {
	return diagramForConcept(conceptID)
}

func diagramForConcept(id string) string {
	if d, ok := conceptDiagrams[id]; ok {
		return d
	}
	return ""
}

// gradeAnswer uses the generator's own grader when available, otherwise falls back to the type-based router.
func (e *Engine) gradeAnswer(conceptID string, expectedAnswer, userAnswer string) grader.Result {
	gradingType := grader.GradingNumeric
	if c := e.dag.Concept(conceptID); c != nil {
		gradingType = grader.GradingType(c.GradingType)
	}
	var res grader.Result
	if gen, err := e.registry.Get(conceptID); err == nil {
		if gg, ok := gen.(generator.GradedGenerator); ok {
			res = gg.Grade(expectedAnswer, userAnswer)
		}
	}
	if res == (grader.Result{}) {
		res = e.gr.Grade(gradingType, expectedAnswer, userAnswer)
	}
	// Diagnosis is applied here rather than inside the graders so it covers a
	// custom GradedGenerator as well as the type router. It never touches
	// Correct: it describes a miss after the fact, and is empty whenever the
	// mistake cannot be named with certainty.
	if !res.Correct && !res.Unavailable {
		res.Diagnosis = grader.Diagnose(gradingType, expectedAnswer, userAnswer)
	}
	return res
}

// GradeAnswer is the single canonical grading entry point: it honours a
// generator's custom GradedGenerator, then falls back to the type router
// (Go graders + SymPy). Every server path must use this so custom-grader
// concepts (e.g. "5 R 3", "3 sqrt(2)", matrices, ratios) grade identically
// everywhere instead of being mis-graded by the bare router.
func (e *Engine) GradeAnswer(conceptID, expectedAnswer, userAnswer string) grader.Result {
	return e.gradeAnswer(conceptID, expectedAnswer, userAnswer)
}

// gradingTypeFor returns a concept's grading_type ("" if unknown), so served
// questions can tell every answer host which input to show.
func (e *Engine) gradingTypeFor(conceptID string) string {
	if c := e.dag.Concept(conceptID); c != nil {
		return c.GradingType
	}
	return ""
}

func (e *Engine) SubmitAnswer(sessionID, studentID, attemptID, answer string, elapsedSeconds float64) (*AnswerResult, error) {
	// Narrow critical section: snapshot the active question under lock, mark
	// answered to reject concurrent duplicates, then release before blocking
	// I/O (grading + DB). This prevents a 12s SymPy grading call from
	// serializing NextQuestion for all users.
	e.mu.Lock()
	as := e.sessions[sessionID]
	if as == nil {
		if persisted := e.rehydrateActiveSession(sessionID, studentID); persisted != nil {
			as = persisted
			e.sessions[sessionID] = as
		}
	}
	if as == nil || as.conceptID == "" || as.answered || as.attemptID != attemptID {
		e.mu.Unlock()
		return nil, fmt.Errorf("%w for session %q", ErrNoActiveQuestion, sessionID)
	}
	sessionFields := struct {
		conceptID          string
		expectedAnswer     string
		explanation        string
		questionText       string
		questionDifficulty float64
		requiredStreak     int
		timeThreshold      float64
		isReview           bool
	}{
		conceptID:          as.conceptID,
		expectedAnswer:     as.expectedAnswer,
		explanation:        as.explanation,
		questionText:       as.questionText,
		questionDifficulty: as.questionDifficulty,
		requiredStreak:     as.requiredStreak,
		timeThreshold:      as.timeThreshold,
		isReview:           as.isReview,
	}
	// Mark answered while still holding the lock so a concurrent
	// SubmitAnswer for the same attemptID is rejected as ErrNoActiveQuestion.
	as.answered = true
	e.mu.Unlock()

	gr := e.gradeAnswer(sessionFields.conceptID, sessionFields.expectedAnswer, answer)
	if gr.Unavailable {
		// Grader infrastructure fault — never a student miss. Re-arm the
		// pending question so the client can retry, and record nothing.
		e.mu.Lock()
		if cur := e.sessions[sessionID]; cur == as {
			as.answered = false
		}
		e.mu.Unlock()
		return &AnswerResult{Ungraded: true, Feedback: gr.Feedback}, nil
	}

	progress, err := e.repo.GetProgress(studentID, sessionFields.conceptID)
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}
	if progress == nil {
		progress = &storage.ConceptProgress{
			StudentID:  studentID,
			ConceptID:  sessionFields.conceptID,
			Status:     string(mastery.StatusUnseen),
			SM2EFactor: 2.5,
		}
	}

	newStatus, evidence, err := e.applyGradedAttempt(progress, gradedAttempt{
		studentID:      studentID,
		conceptID:      sessionFields.conceptID,
		requiredStreak: sessionFields.requiredStreak,
		timeThreshold:  sessionFields.timeThreshold,
		elapsed:        elapsedSeconds,
		correct:        gr.Correct,
		questionText:   sessionFields.questionText,
		difficulty:     difficultyOrNil(sessionFields.questionDifficulty),
	})
	if err != nil {
		return nil, fmt.Errorf("save progress: %w", err)
	}

	// Record attempt
	if err := e.repo.RecordAttempt(storage.AttemptEntry{
		SessionID:      sessionID,
		StudentID:      studentID,
		ConceptID:      sessionFields.conceptID,
		Answer:         answer,
		Expected:       sessionFields.expectedAnswer,
		Correct:        gr.Correct,
		ElapsedSeconds: elapsedSeconds,
		Timestamp:      nowUTC(),
		Question:       sessionFields.questionText,
		Source:         sourceForTask(sessionFields.isReview),
		Explanation:    sessionFields.explanation,
		Difficulty:     difficultyOrNil(sessionFields.questionDifficulty),
	}); err != nil {
		return nil, fmt.Errorf("record attempt: %w", err)
	}

	// The active session already holds the explanation the generator produced
	// for the question that was served (sessionFields.explanation, from
	// Problem.Explanation). Return it on both branches: a correct answer still
	// needs the reasoning that made it right. Grader feedback is empty on a
	// correct grade, so gating on the verdict returned nothing at all there.
	explanation := sessionFields.explanation
	if explanation == "" {
		explanation = gr.Feedback
	}

	taskType := TaskLesson
	if sessionFields.isReview {
		taskType = TaskReview
	} else if strings.HasSuffix(sessionFields.conceptID, ".word") {
		// ADR-002: word problems are TaskMultistep 15 everywhere (practice and
		// study), not the flat TaskLesson 10.
		taskType = TaskMultistep
	}
	xp := computeXPForTask(gr.Correct, elapsedSeconds, sessionFields.timeThreshold, progress.Streak, taskType)

	// PR 1.5: consecutive-miss tracking — halt after 2 (XP=0) + remedial queue.
	halted := false
	var remedial []string
	if !gr.Correct {
		as.consecutiveMisses++
		if as.consecutiveMisses >= 2 && !as.halted {
			halted = true
			xp = 0
			as.halted = true
			as.remedialDifficulty = 0.3
			if c := e.dag.Concept(sessionFields.conceptID); c != nil {
				if len(c.Variants) > 0 {
					as.remedialDifficulty = c.Variants[0].Difficulty
				}
				for _, kp := range c.KeyPrerequisites {
					if e.dag.Concept(kp) != nil {
						as.remedialQueue = append(as.remedialQueue, kp)
						remedial = append(remedial, kp)
					}
				}
			}
			as.remedialQueue = append(as.remedialQueue, sessionFields.conceptID)
			remedial = append(remedial, sessionFields.conceptID)
		}
	} else {
		as.consecutiveMisses = 0
		// Remedial recovery: clear the halt once the queue is drained.
		if as.halted && len(as.remedialQueue) == 0 {
			as.halted = false
		}
	}

	// PR 1.5: negative XP for rushing/guessing — elapsed <2s incorrect twice.
	if !gr.Correct && elapsedSeconds < xppolicy.RushSeconds {
		as.rushCount++
		if as.rushCount >= xppolicy.RushAfter {
			xp = xppolicy.RushPenalty
		}
	}
	if xp != 0 {
		if err := e.repo.AddXP(studentID, xp); err != nil {
			log.Printf("warning: failed to add XP for student %s: %v", studentID, err)
		}
	}

	e.mu.Lock()
	// Re-acquire to update session counters and clear the answered question.
	// as is still the same pointer we marked answered=true above.
	as.lastConceptID = as.conceptID
	if as.isReview {
		as.sessionReview++
	} else {
		as.sessionNew++
	}
	// answered already true from the early mark; keep it true.
	as.conceptID = ""
	as.questionText = ""
	e.persistActiveSession(sessionID, studentID, as, "")
	e.mu.Unlock()

	return &AnswerResult{
		Correct:        gr.Correct,
		Feedback:       gr.Feedback,
		Diagnosis:      gr.Diagnosis,
		NewStatus:      newStatus,
		EvidenceScore:  evidence.Score,
		EvidenceBand:   mastery.EvidenceBand(evidence),
		Explanation:    explanation,
		Streak:         progress.Streak,
		RequiredStreak: sessionFields.requiredStreak,
		XP:             xp,
		ExpectedAnswer: sessionFields.expectedAnswer,
		Halted:         halted,
		Remedial:       remedial,
	}, nil
}

// SubmitStudyAnswer records a Learn-loop answer (LessonQuiz seam per CONTEXT.md:
// Seam). It grades via the concept's grading_type, updates mastery/SM-2/
// weakness/XP, and awards TaskMultistep 15 for *.word else TaskLesson 10.
//
// It takes no expected answer, and that is the point: the answer and the
// explanation both come from the server-side anchor written when the question
// was served, which is keyed by question text. There is no parameter a caller
// could use to grade its own homework, and no code path that grades a study
// answer against client input.
//
// A question with no anchor — the anchor's TTL expired, or the tab sat open
// longer than it — returns ErrNoStudyAnchor so the caller can re-serve rather
// than accept an ungradeable attempt.
// It takes no difficulty argument, and does not need one: the difficulty comes from the
// study anchor, which is the durable record of the question that was actually served. That
// is a stronger source than a caller could supply, for the same reason the expected answer
// does.
func (e *Engine) SubmitStudyAnswer(studentID, conceptID, answer string, elapsedSeconds float64, questionText string) (*AnswerResult, error) {
	taskType := TaskLesson
	if strings.HasSuffix(conceptID, ".word") {
		taskType = TaskMultistep
	}
	return e.submitAnswerWithTask(studentID, conceptID, answer, "", elapsedSeconds, taskType, true, false, questionText, nil)
}

// SubmitQuizAnswer is the single-path quiz grader: TaskQuiz base XP (20),
// progress update, and exactly one AddXP — DB and response agree by
// construction. Unlike SubmitStudyAnswer it never consults the study anchor:
// the quiz expected answer comes from the quiz session.
//
// difficulty is the difficulty the quiz question was generated at, from
// `quiz.Session.LastProblem.Difficulty`. It is a parameter rather than something re-derived
// here because re-deriving it from the learner's *current* weakness would credit the
// attempt with a difficulty nobody actually served — and because the whole point of this
// argument is that the quiz path produces the same evidence as the study path.
func (e *Engine) SubmitQuizAnswer(studentID, conceptID, answer, expected string, elapsedSeconds float64, questionText string, difficulty *float64) (*AnswerResult, error) {
	return e.submitAnswerWithTask(studentID, conceptID, answer, expected, elapsedSeconds, TaskQuiz, false, false, questionText, difficulty)
}

// SubmitQuizDontKnow records an admitted unknown ("I don't know" button) as
// a quiz miss: weakness up, streak reset, remedial queued — but no XP and no
// rushing penalty (an instant admit is honesty, not rushing).
func (e *Engine) SubmitQuizDontKnow(studentID, conceptID, expected string, elapsedSeconds float64, questionText string, difficulty *float64) (*AnswerResult, error) {
	return e.submitAnswerWithTask(studentID, conceptID, "", expected, elapsedSeconds, TaskQuiz, false, true, questionText, difficulty)
}

// difficulty is the difficulty the question that produced this attempt was served at, or
// nil when it is genuinely unknown. For study it comes from the anchor; for quiz it is
// passed by the caller from the served problem. It is never derived from the learner's
// current weakness, because that is a different number: it is the difficulty a *future*
// question would get, not the one the learner just answered.
func (e *Engine) submitAnswerWithTask(studentID, conceptID, answer, expected string, elapsedSeconds float64, taskType string, useAnchor bool, dontKnow bool, questionText string, difficulty *float64) (*AnswerResult, error) {
	if e.dag.Concept(conceptID) == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnknownConcept, conceptID)
	}
	if useAnchor {
		anchor, found := e.studyAnchorFor(studentID, conceptID, questionText)
		if !found || anchor.Answer == "" {
			// No server-side record of this question. Refuse rather than fall
			// back to anything the caller supplied: the one thing a study
			// answer must never be graded against is a guess from the client.
			return nil, ErrNoStudyAnchor
		}
		expected = anchor.Answer
		// The anchor is the record of what was served, so it wins over the argument. This
		// is the study path's equivalent of the quiz passing LastProblem.Difficulty: both
		// read the number from the question the learner actually saw.
		difficulty = anchor.Difficulty
	}
	// Narrow lock: only protect studySessions lookup/creation and studyMisses update.
	e.mu.Lock()
	sessionID := e.studySessions[studentID]
	e.mu.Unlock()
	if sessionID == "" {
		s, err := e.repo.CreateSession(studentID)
		if err != nil {
			return nil, fmt.Errorf("create study session: %w", err)
		}
		sessionID = s.ID
		e.mu.Lock()
		// Double-check after re-acquiring to avoid race creating duplicate sessions
		if existing := e.studySessions[studentID]; existing != "" {
			sessionID = existing
		} else {
			e.studySessions[studentID] = sessionID
		}
		e.mu.Unlock()
	}

	c := e.dag.Concept(conceptID)
	requiredStreak := 3
	timeThreshold := 10.0
	if c != nil {
		requiredStreak = c.MasteryThreshold.Streak
		timeThreshold = c.MasteryThreshold.AvgTimeSeconds
	}
	timeThreshold = e.accommodatedThreshold(studentID, timeThreshold)
	if timeThreshold == 0 {
		timeThreshold = 10.0
	}
	// Grade using concept's grading type via registry/router — skipped for
	// admitted unknowns, which are forced misses (no XP, no rushing penalty).
	gr := e.gradeAnswer(conceptID, expected, answer)
	if dontKnow {
		gr = grader.Result{Correct: false}
	}
	if gr.Unavailable {
		// Grader infrastructure fault — never a student miss; record nothing.
		return &AnswerResult{Ungraded: true, Feedback: gr.Feedback}, nil
	}

	progress, err := e.repo.GetProgress(studentID, conceptID)
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}
	if progress == nil {
		progress = &storage.ConceptProgress{
			StudentID:  studentID,
			ConceptID:  conceptID,
			Status:     string(mastery.StatusUnseen),
			SM2EFactor: 2.5,
		}
	}
	newStatus, evidence, err := e.applyGradedAttempt(progress, gradedAttempt{
		studentID:      studentID,
		conceptID:      conceptID,
		requiredStreak: requiredStreak,
		timeThreshold:  timeThreshold,
		elapsed:        elapsedSeconds,
		correct:        gr.Correct,
		questionText:   questionText,
		difficulty:     difficulty,
	})
	if err != nil {
		return nil, fmt.Errorf("save progress: %w", err)
	}
	// The explanation is the one the generator produced for the question this
	// learner actually saw, held in the anchor written when it was served. It
	// is returned whether the answer was right or wrong: a correct answer still
	// needs the reasoning that made it right, and a wrong one needs the same
	// reasoning to reconstruct the solution. Grader feedback is a status token
	// ("Incorrect"), not teaching, so it is only a last resort for an anchor
	// that predates the explanation being carried.
	explanation := ""
	if anchor, found := e.studyAnchorFor(studentID, conceptID, questionText); found {
		explanation = anchor.Explanation
	}
	if explanation == "" {
		explanation = gr.Feedback
	}
	if err := e.repo.RecordAttempt(storage.AttemptEntry{
		SessionID:      sessionID,
		StudentID:      studentID,
		ConceptID:      conceptID,
		Answer:         answer,
		Expected:       expected,
		Correct:        gr.Correct,
		ElapsedSeconds: elapsedSeconds,
		Timestamp:      nowUTC(),
		Question:       questionText,
		Source:         sourceForTaskType(taskType),
		Explanation:    explanation,
		// The same value the evidence window was just given, so the persisted row and the
		// decision it fed cannot disagree. This column used to be left NULL on this path,
		// which is what made quiz and study evidence incomparable.
		Difficulty: difficulty,
	}); err != nil {
		return nil, fmt.Errorf("record attempt: %w", err)
	}
	xp := computeXPForTask(gr.Correct, elapsedSeconds, timeThreshold, progress.Streak, taskType)

	// PR 1.5: negative XP for rushing/guessing + halt flag at 2 consecutive misses.
	missKey := studentID + "|" + conceptID
	halted := false
	var remedial []string
	e.mu.Lock()
	if !gr.Correct {
		e.studyMisses[missKey]++
		if e.studyMisses[missKey] >= 2 {
			halted = true
		}
		// Same rule as the practice path: a wrong answer answered implausibly fast, twice.
		// It used to require two misses on the concept as well, so a learner who rushed once
		// and was wrong twice for ordinary reasons was never penalised here — the same
		// behaviour was penalised on one route and not the other.
		if !dontKnow && elapsedSeconds < xppolicy.RushSeconds {
			e.studyRush[missKey]++
			if e.studyRush[missKey] >= xppolicy.RushAfter {
				xp = xppolicy.RushPenalty
			}
		}
		// Batch 1: immediate remedial enqueue on quiz miss — missed concept
		// plus its key prerequisites surface first in Study after the quiz.
		if taskType == TaskQuiz {
			e.enqueueQuizRemedialLocked(studentID, conceptID)
			remedial = append([]string(nil), e.quizRemedial[studentID]...)
		}
	} else {
		delete(e.studyMisses, missKey)
		delete(e.studyRush, missKey)
	}
	e.mu.Unlock()
	if xp != 0 {
		if err := e.repo.AddXP(studentID, xp); err != nil {
			log.Printf("warning: failed to add XP for student %s: %v", studentID, err)
		}
	}
	return &AnswerResult{
		Correct:        gr.Correct,
		Feedback:       gr.Feedback,
		Diagnosis:      gr.Diagnosis,
		NewStatus:      newStatus,
		EvidenceScore:  evidence.Score,
		EvidenceBand:   mastery.EvidenceBand(evidence),
		Explanation:    explanation,
		Streak:         progress.Streak,
		RequiredStreak: requiredStreak,
		XP:             xp,
		ExpectedAnswer: expected,
		Halted:         halted,
		Remedial:       remedial,
	}, nil
}

const (
	TaskLesson    = "lesson"
	TaskReview    = "review"
	TaskMultistep = "multistep"
	TaskQuiz      = "quiz"
)

// sourceForTask maps the study flow to the mistakes-report source bucket.
func sourceForTask(isReview bool) string {
	if isReview {
		return "review"
	}
	return "practice"
}

// sourceForTaskType maps MA task types to the mistakes-report source bucket.
func sourceForTaskType(taskType string) string {
	switch taskType {
	case TaskQuiz:
		return "quiz"
	case TaskReview:
		return "review"
	default:
		return "practice"
	}
}

// taskBaseXP maps MA task types to base XP (10/5/15/20).
// Single source of truth lives in internal/xp; this delegates to it.
func taskBaseXP(taskType string) int {
	return xppolicy.BaseXP(taskType)
}

// computeXPForTask is the MA-differed XP calculator (10/5/15/20).
// Delegates to internal/xp so estimators share the exact award math.
func computeXPForTask(correct bool, elapsed, timeThreshold float64, streak int, taskType string) int {
	return xppolicy.Award(correct, elapsed, timeThreshold, streak, taskType)
}

// QuizXP awards TaskQuiz 20 for actionable quiz path (own grading path per Q3).
func (e *Engine) QuizXP(correct bool, elapsed, timeThreshold float64, streak int) int {
	return computeXPForTask(correct, elapsed, timeThreshold, streak, TaskQuiz)
}

// DifficultyFor exposes weakness→difficulty (0.3-1.0) for quiz 80% targeting.
func (e *Engine) DifficultyFor(studentID, conceptID string) float64 {
	return e.computeDifficulty(studentID, conceptID)
}

// TimeLimitFor exposes the accommodated per-question time limit for quiz
// timed closed-book metadata.
func (e *Engine) TimeLimitFor(studentID, conceptID string) float64 {
	base := 10.0
	if c := e.dag.Concept(conceptID); c != nil {
		base = c.MasteryThreshold.AvgTimeSeconds
	}
	return e.accommodatedThreshold(studentID, base)
}

// QuizXPSince returns lifetime XP earned since the last completed quiz
// (or all XP when no quiz completed yet). Floor 0: XP never decreases
// except rushing penalties, which must not manufacture quiz eligibility.
func (e *Engine) QuizXPSince(studentID string) (int, error) {
	total, _, err := e.repo.GetXP(studentID)
	if err != nil {
		return 0, err
	}
	last, err := e.repo.LastQuizCompletion(studentID)
	if err != nil {
		return 0, err
	}
	if last == nil {
		return total, nil
	}
	since := total - last.XPTotal
	if since < 0 {
		since = 0
	}
	return since, nil
}

// QuizDue reports whether the mastery-check gate is reached.
func (e *Engine) QuizDue(studentID string) (bool, error) {
	since, err := e.QuizXPSince(studentID)
	if err != nil {
		return false, err
	}
	return since >= QuizGateXP, nil
}

// RecordQuizCompletion snapshots current lifetime XP as the gate baseline.
func (e *Engine) RecordQuizCompletion(studentID string) error {
	total, _, err := e.repo.GetXP(studentID)
	if err != nil {
		return err
	}
	return e.repo.RecordQuizCompletion(studentID, total)
}

// enqueueQuizRemedialLocked queues a missed quiz concept plus its key
// prerequisites for immediate Study review. Caller holds e.mu.
func (e *Engine) enqueueQuizRemedialLocked(studentID, conceptID string) {
	queue := e.quizRemedial[studentID]
	seen := make(map[string]bool, len(queue)+4)
	for _, id := range queue {
		seen[id] = true
	}
	push := func(id string) {
		if seen[id] || len(queue) >= maxQuizRemedial {
			return
		}
		if e.dag.Concept(id) == nil {
			return
		}
		seen[id] = true
		queue = append(queue, id)
	}
	if c := e.dag.Concept(conceptID); c != nil {
		for _, kp := range c.KeyPrerequisites {
			push(kp)
		}
	}
	push(conceptID)
	e.quizRemedial[studentID] = queue
}

// QuizRemedial returns a snapshot of pending quiz-miss remedial concepts.
func (e *Engine) QuizRemedial(studentID string) []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]string(nil), e.quizRemedial[studentID]...)
}

// ClearQuizRemedial drains the queue (called when Study serves it).
func (e *Engine) ClearQuizRemedial(studentID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.quizRemedial, studentID)
}

// computeXP retains bool-based API for backward compatibility (isReview=true→review 5, false→lesson 10).
func computeXP(correct bool, elapsed, timeThreshold float64, streak int, isReview bool) int {
	taskType := TaskLesson
	if isReview {
		taskType = TaskReview
	}
	return computeXPForTask(correct, elapsed, timeThreshold, streak, taskType)
}

// GetProgress returns the learner's progress with decay applied at read time.
//
// The database stores MASTERED forever — decay is deliberately not persisted
// (mastery.EffectiveStatus) — so the raw row claims a concept is mastered long
// after its review went stale. That made Profile count it as mastered while the
// scheduler simultaneously treated it as DECAYING and scheduled it for review:
// two parts of the app disagreeing about the same fact, with nothing to show
// for it. ConceptProgress carries no last_reviewed for the client to correct
// it with, so it has to happen here.
//
// Every reader goes through this: /api/progress, the Study catalog, readiness,
// and the due-review count.
func (e *Engine) GetProgress(studentID string) (map[string]*storage.ConceptProgress, error) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	for _, p := range progress {
		if mastery.Status(p.Status) != mastery.StatusMastered {
			continue
		}
		// No review timestamp means we cannot claim it is fresh, so treat it as
		// fully decayed rather than pretending it was reviewed just now.
		days := 999.0
		if p.LastReviewed != nil {
			days = now.Sub(*p.LastReviewed).Hours() / 24
		}
		p.Status = string(mastery.EffectiveStatus(mastery.StatusMastered, days, mastery.DecayDays))
	}
	return progress, nil
}

// ProgressWithMasteryPct is GetProgress plus the one derived number the client used to
// recompute for itself.
//
// Additive rather than a change to GetProgress: readiness, the due-review count and the
// study catalogue all read GetProgress and none of them need a percentage, and widening
// its return type would touch every one of them for no gain.
//
// It lives here rather than in the server because the threshold lives here, next to the
// state machine that enforces it. The graph used to derive the same ratio client-side
// from `streak` and a *bundled* copy of the corpus, which is a second source of truth
// that silently disagrees with the server's whenever the corpus is rebuilt.
func (e *Engine) ProgressWithMasteryPct(studentID string) (map[string]ProgressView, error) {
	progress, err := e.GetProgress(studentID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]ProgressView, len(progress))
	for id, p := range progress {
		reqStreak := 0
		if c := e.dag.Concept(id); c != nil {
			reqStreak = c.MasteryThreshold.Streak
		}
		out[id] = ProgressView{
			ConceptProgress: p,
			MasteryPct:      mastery.MasteryPct(mastery.Status(p.Status), p.Streak, reqStreak),
		}
	}
	return out, nil
}

func (e *Engine) GetScores(studentID string) (*scoring.Scores, error) {
	return e.scorer.Compute(studentID)
}

func (e *Engine) GetDAG() *concepts.DAG {
	return e.dag
}

func (e *Engine) GetGeneratorRegistry() *generator.Registry { return e.registry }

func (e *Engine) GetQuestions(conceptID string, count int) ([]storage.Question, error) {
	return e.repo.GetQuestions(conceptID, count)
}

func (e *Engine) GetQuestionCount(conceptID string) (int, error) {
	return e.repo.GetQuestionCount(conceptID)
}

func (e *Engine) GetLessonLoader() *lessons.Loader {
	return e.ll
}

func (e *Engine) GetLeaderboard() ([]leaderboard.Entry, error) {
	return e.lboard.Weekly()
}

func (e *Engine) GetLeagues() (*leaderboard.LeagueBoard, error) {
	return leaderboard.Standings(e.repo, nowUTC())
}

// ShareReport is the read-only parent/teacher view of a student.
type ShareReport struct {
	StudentID string                              `json:"student_id"`
	Name      string                              `json:"name"`
	Scores    *scoring.Scores                     `json:"scores"`
	Activity  []storage.DailyActivity             `json:"activity"`
	Progress  map[string]*storage.ConceptProgress `json:"progress"`
	Weakness  map[string]float64                  `json:"weakness"`
	Attempts  []storage.AttemptEntry              `json:"attempts"`
}

// EnableShare mints a read-only share token for the student.
func (e *Engine) EnableShare(studentID string) (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate share token: %w", err)
	}
	token := "s_" + hex.EncodeToString(b)
	if err := e.repo.SetShareToken(studentID, token); err != nil {
		return "", err
	}
	return token, nil
}

// DisableShare clears the student's share token.
func (e *Engine) DisableShare(studentID string) error {
	return e.repo.SetShareToken(studentID, "")
}

// GetShareReport builds the read-only oversight view for a share token.
func (e *Engine) GetShareReport(token string) (*ShareReport, error) {
	st, err := e.repo.GetStudentByShareToken(token)
	if err != nil || st == nil {
		return nil, fmt.Errorf("invalid share token")
	}
	scores, err := e.scorer.Compute(st.ID)
	if err != nil {
		return nil, err
	}
	activity, err := e.repo.GetDailyActivity(st.ID, 365)
	if err != nil {
		activity = []storage.DailyActivity{}
	}
	progress, err := e.repo.GetAllProgress(st.ID)
	if err != nil {
		progress = map[string]*storage.ConceptProgress{}
	}
	attempts, err := e.repo.GetAttemptsForStudent(st.ID)
	if err != nil {
		attempts = []storage.AttemptEntry{}
	}
	// Newest first, capped: the share payload stays small while the
	// mistakes history itself is uncapped server-side.
	sort.SliceStable(attempts, func(i, j int) bool {
		return attempts[i].Timestamp.After(attempts[j].Timestamp)
	})
	if len(attempts) > 500 {
		attempts = attempts[:500]
	}
	weakAll := e.weaknessMapFromProgress(progress)
	filteredWeak := make(map[string]float64)
	for id, w := range weakAll {
		if _, ok := progress[id]; !ok {
			continue
		}
		filteredWeak[id] = w
	}
	return &ShareReport{
		StudentID: st.ID,
		Name:      st.Name,
		Scores:    scores,
		Activity:  activity,
		Progress:  progress,
		Weakness:  filteredWeak,
		Attempts:  attempts,
	}, nil
}

func (e *Engine) StartDiagnostic() *diagnostic.Session {
	return e.diag.Start()
}

func (e *Engine) NextDiagnosticQuestion(s *diagnostic.Session) (*generator.Problem, string, error) {
	return e.diag.NextQuestion(s)
}

// SubmitDiagnosticAnswerTimed is the MA-parity path: elapsed seconds + the
// per-concept threshold drive automaticity weighting, with the student's
// accommodations (extra_time) applied to the threshold.
func (e *Engine) SubmitDiagnosticAnswerTimed(s *diagnostic.Session, conceptID string, correct bool, elapsed, timeThresh float64) {
	studentID := ""
	s.Lock()
	studentID = s.StudentID
	s.Unlock()
	e.diag.RecordAnswerTimed(s, conceptID, correct, elapsed, e.accommodatedThreshold(studentID, timeThresh))
}

// SubmitDiagnosticDontKnow records an admitted unknown ("I don't know"
// button) as negative evidence through the standard incorrect path.
func (e *Engine) SubmitDiagnosticDontKnow(s *diagnostic.Session, conceptID string, elapsed, timeThresh float64) {
	studentID := ""
	s.Lock()
	studentID = s.StudentID
	s.Unlock()
	e.diag.RecordDontKnow(s, conceptID, elapsed, e.accommodatedThreshold(studentID, timeThresh))
}

// SettleDiagnosticCurrent settles the session's pending question without
// recording an answer (skip escape hatch). Returns the settled concept ID.
func (e *Engine) SettleDiagnosticCurrent(s *diagnostic.Session) string {
	return e.diag.SettleCurrent(s)
}

// DiagnosticRetryAvailableFor reports whether the missed concept's latest
// miss qualifies for a "silly mistake" retry (methodology gate), given the
// staged next concept ("" when none is staged yet).
func (e *Engine) DiagnosticRetryAvailableFor(s *diagnostic.Session, missedCID, nextCID string) bool {
	return e.diag.RetryAvailableFor(s, missedCID, nextCID)
}

// RetryDiagnosticQuestion voids the superseded question's miss and restores
// it as pending for a "silly mistake" retry. Returns restored problem,
// concept ID and name.
func (e *Engine) RetryDiagnosticQuestion(s *diagnostic.Session, conceptID string) (*generator.Problem, string, string, error) {
	return e.diag.RetryServe(s, conceptID)
}

func (e *Engine) IsDiagnosticComplete(s *diagnostic.Session) bool {
	return e.diag.IsComplete(s)
}
func (e *Engine) DiagnosticProgress(s *diagnostic.Session) diagnostic.Progress {
	return e.diag.Progress(s)
}

func (e *Engine) DiagnosticReport(s *diagnostic.Session) *diagnostic.DiagnosticReport {
	return e.diag.Report(s)
}

func (e *Engine) DiagnosticFrontier(s *diagnostic.Session) int {
	return e.diag.FrontierEstimate(s)
}

func (e *Engine) StartGoalDiagnostic(studentID string, conceptIDs []string) (*diagnostic.Session, *Question, error) {
	if e.planner == nil {
		return nil, nil, fmt.Errorf("no planner configured")
	}
	path, err := e.planner.PrerequisitesOf(conceptIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("goal path: %w", err)
	}
	session := e.diag.StartWithPath(path.Concepts)
	if session.State == diagnostic.StateDone || len(path.Concepts) == 0 {
		return session, nil, nil
	}
	prob, cid, err := e.diag.NextQuestion(session)
	if err != nil {
		return nil, nil, err
	}
	c := e.dag.Concept(cid)
	label := cid
	if c != nil {
		label = c.Label
	}
	return session, &Question{
		ConceptID:   cid,
		ConceptName: label,
		Question:    prob.Question,
		GradingType: e.gradingTypeFor(cid),
	}, nil
}

func (e *Engine) ApplyGoalResults(studentID string, session *diagnostic.Session) error {
	session.Lock()
	attempts := make([]diagnostic.Attempt, len(session.Attempts))
	copy(attempts, session.Attempts)
	session.Unlock()
	batch := make([]*storage.ConceptProgress, 0, len(attempts))
	for _, att := range attempts {
		// The diagnostic seeds placement; it must not destroy existing mastery.
		// UpsertProgressBatch overwrites every column, so load the current row
		// and only seed concepts that are still unseen. A concept the student
		// has already started/advanced is left completely untouched.
		existing, err := e.repo.GetProgress(studentID, att.ConceptID)
		if err != nil {
			return fmt.Errorf("load diagnostic progress %s: %w", att.ConceptID, err)
		}
		if existing == nil {
			existing = &storage.ConceptProgress{
				StudentID:  studentID,
				ConceptID:  att.ConceptID,
				SM2EFactor: 2.5,
			}
		} else if existing.Status != string(mastery.StatusUnseen) && existing.Status != "" {
			continue
		}
		status := string(mastery.StatusUnseen)
		weakness := 1.0
		if att.Correct && att.Fast {
			status = string(mastery.StatusPracticing)
			weakness = 0.1
		} else if att.Correct {
			status = string(mastery.StatusLearning)
			weakness = 0.3
		} else {
			weakness = 0.8
		}
		existing.Status = status
		existing.WeaknessScore = weakness
		batch = append(batch, existing)
	}
	if err := e.repo.UpsertProgressBatch(batch); err != nil {
		return fmt.Errorf("save diagnostic progress: %w", err)
	}
	// Set active path to the diagnostic's concept set for focused practice
	if len(session.Attempts) > 0 {
		e.mu.Lock()
		pathSet := make(map[string]bool)
		for _, att := range session.Attempts {
			pathSet[att.ConceptID] = true
		}
		e.activePath[studentID] = pathSet
		e.mu.Unlock()
	}
	return nil
}

func (e *Engine) WeaknessMap(studentID string) map[string]float64 {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil
	}
	return e.weaknessMapFromProgress(progress)
}

// weaknessMapFromProgress builds the weakness map from an already-fetched
// progress snapshot. Hot paths (NextQuestion, reviews, practice, share
// reports) fetch progress once and thread it through here instead of
// re-scanning the table per derivation.
func (e *Engine) weaknessMapFromProgress(progress map[string]*storage.ConceptProgress) map[string]float64 {
	result := make(map[string]float64)
	for _, c := range e.dag.Order() {
		p, ok := progress[c.ID]
		if !ok {
			result[c.ID] = 0.5
			continue
		}
		if p.Status == string(mastery.StatusMastered) {
			result[c.ID] = 0.0
			continue
		}
		result[c.ID] = p.WeaknessScore
	}
	return result
}

func (e *Engine) PracticeConcept(sessionID, studentID, conceptID string) (*Question, error) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, fmt.Errorf("load progress: %w", err)
	}

	c := e.dag.Concept(conceptID)
	if c == nil {
		return nil, fmt.Errorf("concept %q not found", conceptID)
	}

	// Check prerequisites are met
	for _, pid := range c.Prerequisites {
		p, ok := progress[pid]
		if !ok || p.Status != string(mastery.StatusMastered) {
			return nil, fmt.Errorf("prerequisite %q not mastered", pid)
		}
	}

	difficulty := difficultyFromWeakness(e.weaknessMapFromProgress(progress), conceptID)
	attemptID := newAttemptID()
	seedBase := hashSeed(studentID + "|" + conceptID + "|" + attemptID)
	ctx := generator.GeneratorContext{Difficulty: difficulty, Seed: seedBase}
	prob, err := e.registry.GenerateContext(conceptID, ctx)
	if err != nil {
		return nil, fmt.Errorf("generate problem: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	as := e.sessions[sessionID]
	if as == nil {
		as = &activeSession{}
	}
	prevQuestion := as.questionText
	for i := 0; i < 3 && prevQuestion != "" && prob.Question == prevQuestion; i++ {
		ctx.Seed = seedBase + int64(i+1)
		if p2, err2 := e.registry.GenerateContext(conceptID, ctx); err2 == nil {
			prob = p2
		}
	}
	as.conceptID = conceptID
	as.conceptName = c.Label
	as.expectedAnswer = prob.Answer
	as.explanation = prob.Explanation
	as.requiredStreak = c.MasteryThreshold.Streak
	as.timeThreshold = e.accommodatedThreshold(studentID, c.MasteryThreshold.AvgTimeSeconds)
	as.answered = false
	as.attemptID = attemptID
	as.questionText = prob.Question
	as.questionDifficulty = difficulty
	e.sessions[sessionID] = as
	e.persistActiveSession(sessionID, studentID, as, as.questionText)

	var lesson *lessons.Lesson
	if e.ll != nil {
		lesson = e.ll.Lesson(conceptID)
	}

	return &Question{
		ConceptID:   conceptID,
		ConceptName: c.Label,
		Question:    prob.Question,
		GradingType: e.gradingTypeFor(conceptID),
		AttemptID:   as.attemptID,
		Lesson:      lesson,
	}, nil
}

type ConceptTreeNode struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	Unlocked   bool    `json:"unlocked"`
	MasteryPct float64 `json:"mastery_pct"`
	Streak     int     `json:"streak"`
	Status     string  `json:"status"`
}

type SubdomainNode struct {
	Name     string            `json:"name"`
	Concepts []ConceptTreeNode `json:"concepts"`
}

type DomainNode struct {
	Name       string          `json:"name"`
	Subdomains []SubdomainNode `json:"subdomains"`
}

// CorpusSize reports how many concepts and domains this engine was loaded with.
//
// The corpus is the definition of what Mathua teaches, so any surface counting it should read
// it here rather than counting rows: a database query would report what has been visited.
// Zero means the corpus failed to load, which is the honest answer and not a bug to paper over
// with a fallback.
func (e *Engine) CorpusSize() (concepts int, domains int) {
	if e == nil || e.dag == nil {
		return 0, 0
	}
	return e.dag.Count(), len(e.dag.Domains())
}

func (e *Engine) ConceptTree(studentID string) []DomainNode {
	progress, _ := e.repo.GetAllProgress(studentID)
	domains := make(map[string]map[string][]ConceptTreeNode)
	domainOrder := make([]string, 0)
	domainSeen := make(map[string]bool)

	for _, c := range e.dag.Order() {
		if !domainSeen[c.Domain] {
			domainSeen[c.Domain] = true
			domainOrder = append(domainOrder, c.Domain)
		}
		if domains[c.Domain] == nil {
			domains[c.Domain] = make(map[string][]ConceptTreeNode)
		}

		unlocked := true
		for _, pid := range c.Prerequisites {
			p, ok := progress[pid]
			if !ok || p.Status != string(mastery.StatusMastered) {
				unlocked = false
				break
			}
		}

		masteryPct := 0.0
		streak := 0
		status := "UNSEEN"
		if p, ok := progress[c.ID]; ok {
			streak = p.Streak
			status = p.Status
			if c.MasteryThreshold.Streak > 0 {
				masteryPct = float64(streak) / float64(c.MasteryThreshold.Streak)
				if masteryPct > 1 {
					masteryPct = 1
				}
			}
		}

		domains[c.Domain][c.Subdomain] = append(domains[c.Domain][c.Subdomain], ConceptTreeNode{
			ID:         c.ID,
			Label:      c.Label,
			Unlocked:   unlocked,
			MasteryPct: masteryPct,
			Streak:     streak,
			Status:     status,
		})
	}

	var result []DomainNode
	for _, name := range domainOrder {
		d := domains[name]
		var subs []SubdomainNode
		subKeys := make([]string, 0, len(d))
		for k := range d {
			subKeys = append(subKeys, k)
		}
		sort.Strings(subKeys)
		for _, sk := range subKeys {
			subs = append(subs, SubdomainNode{Name: sk, Concepts: d[sk]})
		}
		result = append(result, DomainNode{Name: name, Subdomains: subs})
	}
	return result
}

func (e *Engine) PropagateWeakness(studentID string) {
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return
	}
	weakness := e.weaknessMapFromProgress(progress)
	type update struct {
		cid string
		w   float64
	}
	var updates []update
	for cid, w := range weakness {
		if w > 0.3 {
			for _, dep := range e.dag.DependentsOf(cid) {
				if weakness[dep.ID] < w*0.3 {
					updates = append(updates, update{dep.ID, w * 0.3})
				}
			}
		}
	}
	if len(updates) == 0 {
		return
	}
	// One batched write from the single fetched snapshot — no per-dependent
	// Get/Upsert round trips.
	batch := make([]*storage.ConceptProgress, 0, len(updates))
	for _, u := range updates {
		prog := progress[u.cid]
		if prog == nil {
			prog = &storage.ConceptProgress{
				StudentID: studentID,
				ConceptID: u.cid,
				Status:    string(mastery.StatusUnseen),
			}
		} else {
			cp := *prog
			prog = &cp
		}
		prog.WeaknessScore = u.w
		batch = append(batch, prog)
	}
	if err := e.repo.UpsertProgressBatch(batch); err != nil {
		log.Printf("warning: failed to propagate weakness batch for %s: %v", studentID, err)
	}
}

func (e *Engine) AdjustPlan(studentID string) {
	e.mu.Lock()
	path := e.activePath[studentID]
	if len(path) == 0 {
		e.mu.Unlock()
		return
	}
	// Copy the path under lock to avoid holding the mutex during the DB call
	pathCopy := make(map[string]bool, len(path))
	for k, v := range path {
		pathCopy[k] = v
	}
	e.mu.Unlock()

	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	// Reconcile: start with the current active path, not the stale copy
	for cid := range e.activePath[studentID] {
		if p, ok := progress[cid]; ok && p.Status == string(mastery.StatusMastered) && p.WeaknessScore < 0.2 {
			delete(e.activePath[studentID], cid)
		}
	}
}

func timeOrZero(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// newAttemptID returns a random token binding a submitted answer to the exact
// question the client was shown. A stale submission carrying an older token is
// rejected instead of being graded against the session's current question.
func newAttemptID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func hashSeed(s string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return int64(h.Sum64())
}

func (e *Engine) persistActiveSession(sessionID, studentID string, as *activeSession, questionText string) {
	if e.repo == nil {
		return
	}
	rec := &storage.ActiveSession{
		SessionID:      sessionID,
		StudentID:      studentID,
		ConceptID:      as.conceptID,
		ConceptName:    as.conceptName,
		ExpectedAnswer: as.expectedAnswer,
		AttemptID:      as.attemptID,
		Question:       questionText,
		Explanation:    as.explanation,
		Diagram:        diagramForConcept(as.conceptID),
		IsReview:       as.isReview,
		Answered:       as.answered,
		LastConceptID:  as.lastConceptID,
		SessionReview:  as.sessionReview,
		SessionNew:     as.sessionNew,
		Difficulty:     attemptDifficulty(as),
	}
	if err := e.repo.UpsertActiveSession(rec); err != nil {
		log.Printf("warning: persist active session %s: %v", sessionID, err)
	}
}

func (e *Engine) rehydrateActiveSession(sessionID, studentID string) *activeSession {
	if e.repo == nil {
		return nil
	}
	rec, err := e.repo.GetActiveSession(sessionID)
	if err != nil {
		log.Printf("warning: rehydrate session %s: %v", sessionID, err)
		return nil
	}
	if rec == nil {
		return nil
	}
	if rec.StudentID != "" && rec.StudentID != studentID {
		return nil
	}
	// If no active question (tombstone after SubmitAnswer), return session with
	// carry-over lastConceptID/session counters so NextQuestion can honor skip.
	if rec.ConceptID == "" || rec.AttemptID == "" || rec.Answered {
		if rec.LastConceptID == "" && rec.SessionReview == 0 && rec.SessionNew == 0 {
			return nil
		}
		return &activeSession{
			lastConceptID: rec.LastConceptID,
			sessionReview: rec.SessionReview,
			sessionNew:    rec.SessionNew,
			answered:      true,
		}
	}
	c := e.dag.Concept(rec.ConceptID)
	timeThresh := 10.0
	reqStreak := 0
	if c != nil {
		timeThresh = c.MasteryThreshold.AvgTimeSeconds
		reqStreak = c.MasteryThreshold.Streak
	}
	return &activeSession{
		conceptID:      rec.ConceptID,
		conceptName:    rec.ConceptName,
		expectedAnswer: rec.ExpectedAnswer,
		explanation:    rec.Explanation,
		requiredStreak: reqStreak,
		timeThreshold:  timeThresh,
		isReview:       rec.IsReview,
		answered:       rec.Answered,
		attemptID:      rec.AttemptID,
		questionText:   rec.Question,
		lastConceptID:  rec.LastConceptID,
		sessionReview:  rec.SessionReview,
		sessionNew:     rec.SessionNew,
	}
}

// GetCurrentQuestion returns the verbatim active question for the session
// without advancing. Used for same-session auto-recovery on 409.
func (e *Engine) GetCurrentQuestion(sessionID, studentID string) (*Question, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	as := e.sessions[sessionID]
	if as == nil {
		if persisted := e.rehydrateActiveSession(sessionID, studentID); persisted != nil {
			as = persisted
			e.sessions[sessionID] = as
		}
	}
	if as == nil || as.conceptID == "" || as.answered || as.attemptID == "" || as.questionText == "" {
		return nil, nil
	}
	var lesson *lessons.Lesson
	if e.ll != nil {
		lesson = e.ll.Lesson(as.conceptID)
	}
	return &Question{
		ConceptID:   as.conceptID,
		ConceptName: as.conceptName,
		Question:    as.questionText,
		GradingType: e.gradingTypeFor(as.conceptID),
		IsReview:    as.isReview,
		AttemptID:   as.attemptID,
		Lesson:      lesson,
		Diagram:     diagramForConcept(as.conceptID),
	}, nil
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func (e *Engine) Close() {
	if e.gr != nil {
		e.gr.Close()
	}
}

// difficultyOrNil turns "we generated this at a known difficulty" into a pointer and
// "we did not" into nil.
//
// The distinction is the whole point. An evidence model that weights answers by
// difficulty has to be able to say "not recorded" as well as "trivial", and a stored 0
// would claim the second. Pre-migration attempts and hand-written quiz or diagnostic
// questions land in the nil case honestly rather than being backfilled with a number
// nobody measured.
func difficultyOrNil(d float64) *float64 {
	if d <= 0 {
		return nil
	}
	return &d
}

// attemptDifficulty is the persisted form for the active session row.
func attemptDifficulty(as *activeSession) *float64 {
	if as == nil {
		return nil
	}
	return difficultyOrNil(as.questionDifficulty)
}

// buildEvidenceWindow assembles the evidence the mastery ladder reads.
//
// It fetches the concept's most recent attempts and appends the one being graded, which
// has not been written yet. Reading a bounded window rather than the learner's history is
// the point: GetAttemptsForStudent returns everything — 15,774 rows for someone who has
// worked through the corpus — and this runs on every graded answer.
//
// A failed read is not fatal. With no window the ladder simply does not advance, and the
// learner's next answer will see a full one. The alternative — advancing on no evidence —
// would be mastery without a record of what it was based on.
func (e *Engine) buildEvidenceWindow(studentID, conceptID string, current mastery.Attempt) mastery.Evidence {
	opts := mastery.EvidenceOptions{}
	if c := e.dag.Concept(conceptID); c != nil {
		opts.TimeThreshold = c.MasteryThreshold.AvgTimeSeconds
	}

	prior := make([]mastery.Attempt, 0, mastery.EvidenceWindow)
	if e.repo != nil {
		rows, err := e.repo.GetRecentAttemptsForConcept(studentID, conceptID, mastery.EvidenceWindow)
		if err != nil {
			log.Printf("warning: recent attempts for %s/%s: %v", studentID, conceptID, err)
		} else {
			for _, r := range rows {
				prior = append(prior, mastery.Attempt{
					Correct:    r.Correct,
					Elapsed:    r.ElapsedSeconds,
					Difficulty: r.Difficulty,
					Instance:   r.Question,
				})
			}
		}
	}

	window := make([]mastery.Attempt, 0, len(prior)+1)
	window = append(window, prior...)
	window = append(window, current)
	return mastery.BuildEvidence(window, opts)
}

// gradedAttempt is what one graded answer contributes to a concept's progress.
type gradedAttempt struct {
	studentID      string
	conceptID      string
	requiredStreak int
	timeThreshold  float64
	elapsed        float64
	correct        bool
	questionText   string
	// difficulty the question that produced this attempt was generated at, or nil when it
	// is genuinely unknown. A pointer, not a float64 with 0 meaning unknown: 0 is not a
	// difficulty the engine ever serves, so overloading it would let a forgotten value and
	// a deliberate "unknown" be indistinguishable at the call site — which is exactly the
	// bug that let the quiz path record unknown difficulty for every attempt.
	difficulty *float64
}

// applyGradedAttempt is the single mastery authority.
//
// It used to exist twice — once in SubmitAnswer for study and review, once in
// submitAnswerWithTask for quizzes and the diagnostic — with identical logic and different
// variable names. Two copies of the state machine is how a quiz answer ends up graded
// differently from a study answer, which is the class of divergence ADR-037 spent a week
// removing. There is one now, and both hosts call it.
//
// It owns: the attempt counters, the weakness update, the evidence window, the ladder
// transition, the streak reset, SM-2 and the per-topic speed, the mastery timestamp, and
// the decay clock. Recording the attempt row is left to the caller because the two hosts
// classify it differently.
// It returns the status it settled on and the evidence it decided from, because both are
// answers a caller legitimately needs: the status is reported to the client as
// `new_status`, and the evidence is what lets a client render progress from the same
// number rather than keeping a second copy of the formula.
func (e *Engine) applyGradedAttempt(progress *storage.ConceptProgress, o gradedAttempt) (mastery.Status, mastery.Evidence, error) {
	progress.Attempts++
	progress.LastAttempted = ptrTime(nowUTC())
	// AvgResponseTime reflects every attempt, not just correct ones; excluding
	// misses biases the average low and inflates mastery readiness.
	totalTime := progress.AvgResponseTime*float64(progress.Attempts-1) + o.elapsed
	progress.AvgResponseTime = totalTime / float64(progress.Attempts)
	if o.correct {
		progress.Streak++
		if progress.Streak > progress.BestStreak {
			progress.BestStreak = progress.Streak
		}
	} else {
		progress.Streak = 0
	}

	// Weakness score update
	conceptThreshold := o.timeThreshold
	if conceptThreshold == 0 {
		conceptThreshold = 10.0
	}
	if o.correct {
		progress.WeaknessScore *= 0.5
		if o.elapsed <= conceptThreshold {
			progress.WeaknessScore *= 0.3
		}
	} else {
		progress.WeaknessScore += 0.2
		if progress.WeaknessScore > 1.0 {
			progress.WeaknessScore = 1.0
		}
	}

	// Mastery transition, on evidence rather than on the streak.
	//
	// The window includes the attempt being graded. It has not been written yet —
	// RecordAttempt runs at the end of this function — so without appending it here the
	// very first answer on a concept would be judged on an empty window and a learner
	// could not leave UNSEEN until they had answered twice before being credited with
	// anything.
	evidence := e.buildEvidenceWindow(o.studentID, o.conceptID, mastery.Attempt{
		Correct:    o.correct,
		Elapsed:    o.elapsed,
		Difficulty: o.difficulty,
		Instance:   o.questionText,
	})

	ctx := mastery.TransitionCtx{
		Evidence:          evidence,
		Streak:            progress.Streak,
		RequiredStreak:    o.requiredStreak,
		AvgResponseTime:   progress.AvgResponseTime,
		ResponseThreshold: o.timeThreshold,
	}
	newStatus := e.machine.Next(mastery.Status(progress.Status), ctx)
	oldStatus := progress.Status
	progress.Status = string(newStatus)
	advanced := newStatus != mastery.Status(oldStatus)

	// Reset the streak on a tier advance. It no longer gates anything — the ladder above
	// does — but it is still the honest running count of consecutive correct answers, it
	// still feeds SM-2 quality, and the shelf still reads it to tell "resume" from "new".
	if advanced && o.correct {
		progress.Streak = 1
	}

	// SM-2 update — PR 1.2 student model: scale interval by per-topic learningSpeed
	quality := mastery.SM2Quality(
		o.correct && progress.Streak >= o.requiredStreak,
		progress.AvgResponseTime/o.timeThreshold,
	)
	prevSM2 := scheduler.SM2{
		Repetitions: progress.SM2Repetitions,
		EFactor:     progress.SM2EFactor,
		Interval:    progress.SM2Interval,
	}
	learningSpeed := 1.0
	if ts, err := e.repo.GetTopicSpeed(o.studentID, o.conceptID); err == nil && ts != nil {
		learningSpeed = ts.LearningSpeed
	}
	nextSM2 := scheduler.ComputeSM2WithSpeed(prevSM2, quality, learningSpeed)
	progress.SM2Repetitions = nextSM2.Repetitions
	progress.SM2EFactor = nextSM2.EFactor
	progress.SM2Interval = nextSM2.Interval
	// Persist updated per-topic speed from timeRatio + streak
	newSpeed := scheduler.UpdateLearningSpeed(learningSpeed, o.correct, progress.AvgResponseTime/o.timeThreshold, progress.Streak)
	if err := e.repo.UpsertTopicSpeed(&storage.TopicSpeed{
		StudentID:     o.studentID,
		ConceptID:     o.conceptID,
		EFactor:       nextSM2.EFactor,
		Interval:      nextSM2.Interval,
		Repetitions:   nextSM2.Repetitions,
		LearningSpeed: newSpeed,
	}); err != nil {
		log.Printf("warning: upsert topic speed %s/%s: %v", o.studentID, o.conceptID, err)
	}

	if newStatus == mastery.StatusMastered && mastery.Status(oldStatus) != mastery.StatusMastered {
		now := nowUTC()
		progress.MasteredAt = &now
	}
	if newStatus == mastery.StatusMastered {
		progress.WeaknessScore = 0.0
	}
	// The decay clock starts when a tier is attained, not when a streak happens to reach
	// the old threshold.
	//
	// It was gated on `streak >= required && avg time <= threshold`, which was a proxy for
	// "a tier was just advanced" back when that was what advancement meant. It is not any
	// more: the ladder advances on evidence, so a concept can reach MASTERED with a streak
	// of 1. That left `LastReviewed` nil, and `GetProgress` treats a nil `LastReviewed` as
	// fully decayed — so a concept could be MASTERED and simultaneously read as DECAYING,
	// depending on which rule fired first. Seen in the engine tests as the ladder passing
	// through DECAYING on its way up.
	//
	// A tier attainment is exactly what the decay clock should measure from: the last time
	// the learner demonstrated this concept.
	if advanced || (progress.Streak >= o.requiredStreak && progress.AvgResponseTime <= o.timeThreshold) {
		now := nowUTC()
		progress.LastReviewed = &now
		nextReview := now.AddDate(0, 0, nextSM2.Interval)
		progress.NextReviewDue = &nextReview
	}

	if err := e.repo.UpsertProgress(progress); err != nil {
		return mastery.Status(progress.Status), evidence, err
	}
	e.PropagateWeakness(o.studentID)
	return newStatus, evidence, nil
}
