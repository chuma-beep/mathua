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
	rushCount          int
}

// ErrNoActiveQuestion is returned when a session has no unanswered question
// ready for grading — e.g. a duplicate/stale submission racing the current one.
var ErrNoActiveQuestion = fmt.Errorf("no active question")

// ErrUnknownConcept is returned when an answer references a concept_id absent
// from the DAG. Rejecting (instead of persisting a fallback progress row)
// keeps garbage IDs from farming XP and polluting progress/attempt tables.
var ErrUnknownConcept = fmt.Errorf("unknown concept")

const (
	// serverSessionStudyExpected is the server_sessions kind for H1b
	// anti-cheat anchors. studyExpectedTTL bounds the durable row: a
	// practice question answered within a day is normal; older rows are
	// stale and treated as missing.
	serverSessionStudyExpected = "study_expected"
	studyExpectedTTL           = 24 * time.Hour
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
	// PR 1.5: stable per-student session id for study/quiz attempt FK.
	studySessions map[string]string
	// H1b: server-side expected answers for study seam to prevent client cheat
	studyExpected map[string]string
	// Batch 1: immediate remedial queue from quiz misses (studentID → conceptIDs).
	// In-memory by design, same as diag/quiz sessions: a restart just means
	// a retake; weakness propagation (+0.2/miss) is the durable signal.
	quizRemedial map[string][]string
}

// QuizGateXP is the MA-verbatim mastery-check interval (CONTEXT.md Q3 lock).
const QuizGateXP = 150

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
		studySessions: make(map[string]string),
		studyExpected: make(map[string]string),
		quizRemedial:  make(map[string][]string),
	}
}

func (e *Engine) SetStudyExpected(studentID, conceptID, expected string) {
	key := studentID + "|" + conceptID
	e.mu.Lock()
	if e.studyExpected == nil {
		e.studyExpected = make(map[string]string)
	}
	e.studyExpected[key] = expected
	e.mu.Unlock()
	// Write-through to durable storage so a restart doesn't silently drop
	// the H1b anti-cheat anchor (Fix 6). Best-effort: memory is the fast
	// path, the row is the fallback.
	if e.repo != nil {
		exp := time.Now().UTC().Add(studyExpectedTTL).Format(time.RFC3339)
		if err := e.repo.UpsertServerSession(serverSessionStudyExpected, key, expected, exp); err != nil {
			log.Printf("warning: persist study expected %s: %v", key, err)
		}
	}
}

func (e *Engine) popStudyExpected(studentID, conceptID string) (string, bool) {
	key := studentID + "|" + conceptID
	e.mu.Lock()
	val, ok := e.studyExpected[key]
	if ok {
		delete(e.studyExpected, key)
	}
	e.mu.Unlock()
	if ok {
		if e.repo != nil {
			_ = e.repo.DeleteServerSession(serverSessionStudyExpected, key)
		}
		return val, true
	}
	// Restart fallback: answer against the durable row, then consume it.
	if e.repo != nil {
		if v, exp, found, err := e.repo.GetServerSession(serverSessionStudyExpected, key); err == nil && found {
			if exp == "" || isFutureRFC3339(exp) {
				_ = e.repo.DeleteServerSession(serverSessionStudyExpected, key)
				return v, true
			}
			_ = e.repo.DeleteServerSession(serverSessionStudyExpected, key)
		} else if err != nil {
			log.Printf("warning: read study expected %s: %v", key, err)
		}
	}
	return "", false
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
	goal := 30
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
	snapshots := make(map[string]*scheduler.ConceptSnapshot, len(progress))
	for cid, p := range progress {
		c := e.dag.Concept(cid)
		reqStreak := 0
		timeThresh := 0.0
		if c != nil {
			reqStreak = c.MasteryThreshold.Streak
			timeThresh = c.MasteryThreshold.AvgTimeSeconds
		}
		status := mastery.Status(p.Status)
		snapshots[cid] = &scheduler.ConceptSnapshot{
			Status:         status,
			Streak:         p.Streak,
			LastAttempted:  timeOrZero(p.LastAttempted),
			LastReviewed:   timeOrZero(p.LastReviewed),
			NextReviewDue:  p.NextReviewDue,
			RequiredStreak: reqStreak,
			TimeThreshold:  timeThresh,
			WeaknessScore:  p.WeaknessScore,
		}
	}
	for _, c := range e.dag.Order() {
		if _, ok := snapshots[c.ID]; !ok {
			snapshots[c.ID] = &scheduler.ConceptSnapshot{
				Status:         mastery.StatusUnseen,
				RequiredStreak: c.MasteryThreshold.Streak,
				TimeThreshold:  c.MasteryThreshold.AvgTimeSeconds,
			}
		}
	}

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
	snapshots := make(map[string]*scheduler.ConceptSnapshot, len(progress))
	for cid, p := range progress {
		c := e.dag.Concept(cid)
		reqStreak := 0
		timeThresh := 0.0
		if c != nil {
			reqStreak = c.MasteryThreshold.Streak
			timeThresh = c.MasteryThreshold.AvgTimeSeconds
		}
		snapshots[cid] = &scheduler.ConceptSnapshot{
			Status:         mastery.Status(p.Status),
			Streak:         p.Streak,
			LastAttempted:  timeOrZero(p.LastAttempted),
			LastReviewed:   timeOrZero(p.LastReviewed),
			NextReviewDue:  p.NextReviewDue,
			RequiredStreak: reqStreak,
			TimeThreshold:  timeThresh,
			WeaknessScore:  p.WeaknessScore,
		}
	}

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
	"alg.conic.ellipse":   "/diagrams/algebrica/ellipse-1.png",
	"alg.conic.hyperbola": "/diagrams/algebrica/hyperbola-1.png",
	"alg.conic.parabola":  "/diagrams/algebrica/parabola-1.png",
	"alg.seq.arithmetic":  "/diagrams/algebrica/arithmetic-sequence.png",
	"alg.seq.geometric":   "/diagrams/algebrica/geometri-sequence-1.png",
	"stat.dist.normal":    "/diagrams/algebrica/normal-distribution-1.png",

	// Functions (G3)
	"alg.func.absolute_value": "/diagrams/algebrica/absolute-value-1.png",
	"alg.exp.concept":         "/diagrams/algebrica/exponential-function-1.png",
	"alg.func.even_odd":       "/diagrams/algebrica/even-odd-functions-1.png",
	"alg.func.monotonicity":   "/diagrams/algebrica/increasing-decreasing-function.png",
	"alg.func.sigmoid":        "/diagrams/algebrica/sigmoid-function.png",
	"alg.func.sign":           "/diagrams/algebrica/sign-function-1.png",
	"alg.func.composite":      "/diagrams/algebrica/composite-functions-1-1.png",
	"alg.ineq.quadratic":      "/diagrams/algebrica/quadratic-inequalities-1-1.png",
	"alg.ineq.sign_analysis":  "/diagrams/algebrica/sign-analysis-1.png",
	"alg.quad.complex":        "/diagrams/algebrica/quad-eq-complex-roots.png",
	"alg.quad.parametric":     "/diagrams/algebrica/quadratic-equations-params-6.png",
	"trig.eq.basic":           "/diagrams/algebrica/trigonometric-equations-1.png",

	// Calculus (G3)
	"calc.limit.asymptotes":            "/diagrams/algebrica/asymptotes-1.png",
	"calc.limit.squeeze":               "/diagrams/algebrica/squeeze-theorem.png",
	"calc.limit.uniform_continuity":    "/diagrams/algebrica/uniform-continuity-1.png",
	"calc.limit.big_o":                 "/diagrams/algebrica/big-o-notation-1.png",
	"calc.limit.weierstrass":           "/diagrams/algebrica/weierstrass-theorem-1.png",
	"calc.deriv.difference_quotient":   "/diagrams/algebrica/difference-quotient-2.png",
	"calc.deriv.partial":               "/diagrams/algebrica/partial-derivatives-1.png",
	"calc.deriv.non_differentiability": "/diagrams/algebrica/non-differentiable-points-1.png",
	"calc.deriv.convexity":             "/diagrams/algebrica/convexity-1-1.png",
	"calc.deriv.rolle":                 "/diagrams/algebrica/rolle-theorem-1-1.png",
	"calc.deriv.applications":          "/diagrams/algebrica/velocity-1-1.png",
	"calc.seq.convergence":             "/diagrams/algebrica/sequences-conv-1.png",
	"calc.seq.cauchy":                  "/diagrams/algebrica/cauchy-sequence-1.png",
	"calc.series.harmonic":             "/diagrams/algebrica/harmonic-series-1-2.png",
	"calc.series.function_series":      "/diagrams/algebrica/sequence-functions-1.png",
	"calc.series.cauchy_criterion":     "/diagrams/algebrica/series-cauchy-1.png",

	// Stats / vectors / misc (G3)
	"stat.dist.student_t":             "/diagrams/algebrica/student-t-distribution.png",
	"stat.dist.uniform":               "/diagrams/algebrica/uniform-distribution.png",
	"stat.dist.beta":                  "/diagrams/algebrica/beta-distribution-1.png",
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
	"calc.limit.epsilon_delta":   "/diagrams/algebrica/limits-1.png",
	"calc.limit.infinite":        "/diagrams/algebrica/limits-2.png",
	"calc.limit.infinity":        "/diagrams/algebrica/limits-3.png",
	"precalc.conic.hyperbola":    "/diagrams/algebrica/hyperbola-2.png",
	"ml.forward_pass":            "/diagrams/algebrica/neural-network-2.png",
	"calc.deriv.global_extrema":  "/diagrams/algebrica/max-min-1.png",
	"calc.deriv.fermat":          "/diagrams/algebrica/fermat-1.png",
	"precalc.parametric.graph":   "/diagrams/algebrica/velocity-2.png",
	"calc.limit.discontinuity":   "/diagrams/algebrica/Heaviside-function-1.png",
	"calc.limit.concept":         "/diagrams/algebrica/continuous-functions.png",
	"calc.deriv.concept":         "/diagrams/algebrica/derivatives-1.png",
	"ode.first_order.separable":  "/diagrams/algebrica/differential-equations-2.png",
	"precalc.conic.ellipse":      "/diagrams/algebrica/ellipse-2.png",
	"geo.coord.distance":         "/diagrams/algebrica/euclidean-distance.png",
	"alg.systems.concept":        "/diagrams/algebrica/lines-5.png",
	"alg.func.concept":           "/diagrams/algebrica/functions-5.png",
	"precalc.seq.geometric":      "/diagrams/algebrica/geometri-sequence-2-1.png",
	"alg.ineq.compound":          "/diagrams/algebrica/inequalities-1.png",
	"calc.series.integral_test":  "/diagrams/algebrica/integral-test-series-1.png",
	"alg.func.inverse":           "/diagrams/algebrica/inverse-function-1.png",
	"alg.ineq.irrational":        "/diagrams/algebrica/irr-ineq-1.png",
	"calc.deriv.mvt":             "/diagrams/algebrica/lagrange-theoreme-4.png",
	"alg.linear.standard_form":   "/diagrams/algebrica/line-1.png",
	"calc.limit.little_o":        "/diagrams/algebrica/little-o-1.png",
	"alg.log.concept":            "/diagrams/algebrica/logharithm-6.png",
	"calc.deriv.critical_points": "/diagrams/algebrica/max-min-3-1.png",
	"calc.deriv.inflection":      "/diagrams/algebrica/maximum-minimum-7.png",
	"stat.dist.z_table":          "/diagrams/algebrica/normal-distribution-standard-1.png",
	"alg.func.quad":              "/diagrams/algebrica/parabola-2.png",
	"precalc.conic.parabola":     "/diagrams/algebrica/parabola-5-1.png",
	"calc.deriv.tangent_line":    "/diagrams/algebrica/parabola-7.png",
	"alg.quad.discriminant":      "/diagrams/algebrica/quadratic-inequalities-2.png",
	"geo.circle.parts":           "/diagrams/algebrica/circumference-5.png",
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
	if gen, err := e.registry.Get(conceptID); err == nil {
		if gg, ok := gen.(generator.GradedGenerator); ok {
			return gg.Grade(expectedAnswer, userAnswer)
		}
	}
	c := e.dag.Concept(conceptID)
	gradingType := grader.GradingNumeric
	if c != nil {
		gradingType = grader.GradingType(c.GradingType)
	}
	return e.gr.Grade(gradingType, expectedAnswer, userAnswer)
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
		conceptID      string
		expectedAnswer string
		explanation    string
		requiredStreak int
		timeThreshold  float64
		isReview       bool
	}{
		conceptID:      as.conceptID,
		expectedAnswer: as.expectedAnswer,
		explanation:    as.explanation,
		requiredStreak: as.requiredStreak,
		timeThreshold:  as.timeThreshold,
		isReview:       as.isReview,
	}
	// Mark answered while still holding the lock so a concurrent
	// SubmitAnswer for the same attemptID is rejected as ErrNoActiveQuestion.
	as.answered = true
	e.mu.Unlock()

	gr := e.gradeAnswer(sessionFields.conceptID, sessionFields.expectedAnswer, answer)

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

	progress.Attempts++
	progress.LastAttempted = ptrTime(nowUTC())
	if gr.Correct {
		progress.Streak++
		if progress.Streak > progress.BestStreak {
			progress.BestStreak = progress.Streak
		}
		totalTime := progress.AvgResponseTime*float64(progress.Attempts-1) + elapsedSeconds
		progress.AvgResponseTime = totalTime / float64(progress.Attempts)
	} else {
		progress.Streak = 0
	}

	// Weakness score update
	conceptThreshold := sessionFields.timeThreshold
	if conceptThreshold == 0 {
		conceptThreshold = 10.0
	}
	if gr.Correct {
		progress.WeaknessScore *= 0.5
		if elapsedSeconds <= conceptThreshold {
			progress.WeaknessScore *= 0.3
		}
	} else {
		progress.WeaknessScore += 0.2
		if progress.WeaknessScore > 1.0 {
			progress.WeaknessScore = 1.0
		}
	}

	// Mastery transition
	ctx := mastery.TransitionCtx{
		Streak:            progress.Streak,
		RequiredStreak:    sessionFields.requiredStreak,
		AvgResponseTime:   progress.AvgResponseTime,
		ResponseThreshold: sessionFields.timeThreshold,
	}
	newStatus := e.machine.Next(mastery.Status(progress.Status), ctx)
	oldStatus := progress.Status
	progress.Status = string(newStatus)

	// Reset streak when advancing to a new mastery level.
	if newStatus != mastery.Status(oldStatus) && gr.Correct {
		progress.Streak = 1
	}

	// SM-2 update — PR 1.2 student model: scale interval by per-topic learningSpeed
	quality := mastery.SM2Quality(
		gr.Correct && progress.Streak >= sessionFields.requiredStreak,
		progress.AvgResponseTime/sessionFields.timeThreshold,
	)
	prevSM2 := scheduler.SM2{
		Repetitions: progress.SM2Repetitions,
		EFactor:     progress.SM2EFactor,
		Interval:    progress.SM2Interval,
	}
	learningSpeed := 1.0
	if ts, err := e.repo.GetTopicSpeed(studentID, sessionFields.conceptID); err == nil && ts != nil {
		learningSpeed = ts.LearningSpeed
	}
	nextSM2 := scheduler.ComputeSM2WithSpeed(prevSM2, quality, learningSpeed)
	progress.SM2Repetitions = nextSM2.Repetitions
	progress.SM2EFactor = nextSM2.EFactor
	progress.SM2Interval = nextSM2.Interval
	// Persist updated per-topic speed from timeRatio + streak
	newSpeed := scheduler.UpdateLearningSpeed(learningSpeed, gr.Correct, progress.AvgResponseTime/sessionFields.timeThreshold, progress.Streak)
	if err := e.repo.UpsertTopicSpeed(&storage.TopicSpeed{
		StudentID:     studentID,
		ConceptID:     sessionFields.conceptID,
		EFactor:       nextSM2.EFactor,
		Interval:      nextSM2.Interval,
		Repetitions:   nextSM2.Repetitions,
		LearningSpeed: newSpeed,
	}); err != nil {
		log.Printf("warning: upsert topic speed %s/%s: %v", studentID, sessionFields.conceptID, err)
	}

	if newStatus == mastery.StatusMastered && mastery.Status(oldStatus) != mastery.StatusMastered {
		now := nowUTC()
		progress.MasteredAt = &now
	}
	if newStatus == mastery.StatusMastered {
		progress.WeaknessScore = 0.0
	}
	if progress.Streak >= sessionFields.requiredStreak && progress.AvgResponseTime <= sessionFields.timeThreshold {
		now := nowUTC()
		progress.LastReviewed = &now
		nextReview := now.AddDate(0, 0, nextSM2.Interval)
		progress.NextReviewDue = &nextReview
	}

	if err := e.repo.UpsertProgress(progress); err != nil {
		return nil, fmt.Errorf("save progress: %w", err)
	}

	e.PropagateWeakness(studentID)

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
	}); err != nil {
		return nil, fmt.Errorf("record attempt: %w", err)
	}

	explanation := ""
	if !gr.Correct {
		explanation = sessionFields.explanation
	}

	taskType := TaskLesson
	if sessionFields.isReview {
		taskType = TaskReview
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
	if !gr.Correct && elapsedSeconds < 2.0 {
		as.rushCount++
		if as.rushCount >= 2 {
			xp = -5
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
		NewStatus:      newStatus,
		Explanation:    explanation,
		Streak:         progress.Streak,
		RequiredStreak: sessionFields.requiredStreak,
		XP:             xp,
		ExpectedAnswer: sessionFields.expectedAnswer,
		Halted:         halted,
		Remedial:       remedial,
	}, nil
}

// SubmitStudyAnswer records a Study-library answer (LessonQuiz seam per CONTEXT.md: Seam).
// It grades via the concept's grading_type, updates mastery/SM-2/weakness/XP, and
// awards TaskMultistep 15 for *.word else TaskLesson 10 (Q2 lock).
// H1b: if a server-side expected was stored via SetStudyExpected (practice
// generation), it is used instead of the client-supplied expected to prevent
// trivial cheat (client sending expected==answer).
func (e *Engine) SubmitStudyAnswer(studentID, conceptID, answer, expected string, elapsedSeconds float64) (*AnswerResult, error) {
	taskType := TaskLesson
	if strings.HasSuffix(conceptID, ".word") {
		taskType = TaskMultistep
	}
	return e.submitAnswerWithTask(studentID, conceptID, answer, expected, elapsedSeconds, taskType, true)
}

// SubmitQuizAnswer is the single-path quiz grader: TaskQuiz base XP (20),
// progress update, and exactly one AddXP — DB and response agree by
// construction. Unlike SubmitStudyAnswer it never consults studyExpected:
// the quiz expected answer comes from the quiz session.
func (e *Engine) SubmitQuizAnswer(studentID, conceptID, answer, expected string, elapsedSeconds float64) (*AnswerResult, error) {
	return e.submitAnswerWithTask(studentID, conceptID, answer, expected, elapsedSeconds, TaskQuiz, false)
}

func (e *Engine) submitAnswerWithTask(studentID, conceptID, answer, expected string, elapsedSeconds float64, taskType string, useStudyExpected bool) (*AnswerResult, error) {
	if e.dag.Concept(conceptID) == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnknownConcept, conceptID)
	}
	if useStudyExpected {
		if stored, ok := e.popStudyExpected(studentID, conceptID); ok && stored != "" {
			expected = stored
		}
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
	// Grade using concept's grading type via registry/router.
	gr := e.gradeAnswer(conceptID, expected, answer)

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
	progress.Attempts++
	progress.LastAttempted = ptrTime(nowUTC())
	if gr.Correct {
		progress.Streak++
		if progress.Streak > progress.BestStreak {
			progress.BestStreak = progress.Streak
		}
		totalTime := progress.AvgResponseTime*float64(progress.Attempts-1) + elapsedSeconds
		progress.AvgResponseTime = totalTime / float64(progress.Attempts)
	} else {
		progress.Streak = 0
	}
	conceptThreshold := timeThreshold
	if gr.Correct {
		progress.WeaknessScore *= 0.5
		if elapsedSeconds <= conceptThreshold {
			progress.WeaknessScore *= 0.3
		}
	} else {
		progress.WeaknessScore += 0.2
		if progress.WeaknessScore > 1.0 {
			progress.WeaknessScore = 1.0
		}
	}
	ctx := mastery.TransitionCtx{
		Streak:            progress.Streak,
		RequiredStreak:    requiredStreak,
		AvgResponseTime:   progress.AvgResponseTime,
		ResponseThreshold: timeThreshold,
	}
	newStatus := e.machine.Next(mastery.Status(progress.Status), ctx)
	oldStatus := progress.Status
	progress.Status = string(newStatus)
	if newStatus != mastery.Status(oldStatus) && gr.Correct {
		progress.Streak = 1
	}
	quality := mastery.SM2Quality(
		gr.Correct && progress.Streak >= requiredStreak,
		progress.AvgResponseTime/timeThreshold,
	)
	prevSM2 := scheduler.SM2{
		Repetitions: progress.SM2Repetitions,
		EFactor:     progress.SM2EFactor,
		Interval:    progress.SM2Interval,
	}
	learningSpeed := 1.0
	if ts, err := e.repo.GetTopicSpeed(studentID, conceptID); err == nil && ts != nil {
		learningSpeed = ts.LearningSpeed
	}
	nextSM2 := scheduler.ComputeSM2WithSpeed(prevSM2, quality, learningSpeed)
	progress.SM2Repetitions = nextSM2.Repetitions
	progress.SM2EFactor = nextSM2.EFactor
	progress.SM2Interval = nextSM2.Interval
	newSpeed := scheduler.UpdateLearningSpeed(learningSpeed, gr.Correct, progress.AvgResponseTime/timeThreshold, progress.Streak)
	if err := e.repo.UpsertTopicSpeed(&storage.TopicSpeed{
		StudentID:     studentID,
		ConceptID:     conceptID,
		EFactor:       nextSM2.EFactor,
		Interval:      nextSM2.Interval,
		Repetitions:   nextSM2.Repetitions,
		LearningSpeed: newSpeed,
	}); err != nil {
		log.Printf("warning: upsert topic speed %s/%s: %v", studentID, conceptID, err)
	}
	if newStatus == mastery.StatusMastered && mastery.Status(oldStatus) != mastery.StatusMastered {
		now := nowUTC()
		progress.MasteredAt = &now
	}
	if newStatus == mastery.StatusMastered {
		progress.WeaknessScore = 0.0
	}
	if progress.Streak >= requiredStreak && progress.AvgResponseTime <= timeThreshold {
		now := nowUTC()
		progress.LastReviewed = &now
		nextReview := now.AddDate(0, 0, nextSM2.Interval)
		progress.NextReviewDue = &nextReview
	}
	if err := e.repo.UpsertProgress(progress); err != nil {
		return nil, fmt.Errorf("save progress: %w", err)
	}
	e.PropagateWeakness(studentID)
	if err := e.repo.RecordAttempt(storage.AttemptEntry{
		SessionID:      sessionID,
		StudentID:      studentID,
		ConceptID:      conceptID,
		Answer:         answer,
		Expected:       expected,
		Correct:        gr.Correct,
		ElapsedSeconds: elapsedSeconds,
		Timestamp:      nowUTC(),
	}); err != nil {
		return nil, fmt.Errorf("record attempt: %w", err)
	}
	explanation := ""
	if !gr.Correct {
		// Prefer generator explanation; fallback to grader feedback.
		if c != nil {
			// Try to fetch explanation via expected (already passed) — keep grader feedback as explanation.
			explanation = gr.Feedback
		}
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
		if elapsedSeconds < 2.0 && e.studyMisses[missKey] >= 2 {
			xp = -5
		}
		// Batch 1: immediate remedial enqueue on quiz miss — missed concept
		// plus its key prerequisites surface first in Study after the quiz.
		if taskType == TaskQuiz {
			e.enqueueQuizRemedialLocked(studentID, conceptID)
			remedial = append([]string(nil), e.quizRemedial[studentID]...)
		}
	} else {
		delete(e.studyMisses, missKey)
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
		NewStatus:      newStatus,
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

// taskBaseXP maps MA task types to base XP (10/5/15/20).
func taskBaseXP(taskType string) int {
	switch taskType {
	case TaskReview:
		return 5
	case TaskMultistep:
		return 15
	case TaskQuiz:
		return 20
	case TaskLesson:
		fallthrough
	default:
		return 10
	}
}

// computeXPForTask is the MA-differed XP calculator (10/5/15/20).
func computeXPForTask(correct bool, elapsed, timeThreshold float64, streak int, taskType string) int {
	if !correct {
		return 0
	}
	base := taskBaseXP(taskType)
	ratio := elapsed / timeThreshold
	if ratio <= 0 {
		ratio = 0.01
	}
	timeMultiplier := 2.0 - ratio
	if timeMultiplier < 0.5 {
		timeMultiplier = 0.5
	}
	if timeMultiplier > 1.5 {
		timeMultiplier = 1.5
	}
	streakMultiplier := 1.0 + float64(min(streak, 10))*0.1
	return int(float64(base) * timeMultiplier * streakMultiplier)
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

// QuizDue reports whether the 150 XP mastery-check gate is reached.
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

func (e *Engine) GetProgress(studentID string) (map[string]*storage.ConceptProgress, error) {
	return e.repo.GetAllProgress(studentID)
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

// EfficacyReport is the instrumentation summary (improve.md:82).
type EfficacyReport struct {
	ConceptsTouched       int     `json:"concepts_touched"`
	FirstPassRate         float64 `json:"first_pass_rate"`  // correct on attempt 1
	SecondPassRate        float64 `json:"second_pass_rate"` // correct within first 2 attempts
	AvgAttemptsPerConcept float64 `json:"avg_attempts_per_concept"`
	TotalAttempts         int     `json:"total_attempts"`
	StudentsTracked       int     `json:"students_tracked,omitempty"`
}

// Efficacy computes first-pass / second-pass rates from the attempt log.
func (e *Engine) Efficacy(studentID string) (*EfficacyReport, error) {
	attempts, err := e.repo.GetAttemptsForStudent(studentID)
	if err != nil {
		return nil, err
	}
	return computeEfficacy(attempts), nil
}

// AggregateEfficacy computes product-wide first-pass / second-pass rates
// across all students (the MA parity metric: 93% / 98%).
func (e *Engine) AggregateEfficacy() (*EfficacyReport, error) {
	attempts, err := e.repo.GetAllAttempts()
	if err != nil {
		return nil, err
	}
	rep := computeEfficacy(attempts)
	students := make(map[string]bool)
	for _, a := range attempts {
		students[a.StudentID] = true
	}
	rep.StudentsTracked = len(students)
	return rep, nil
}

// computeEfficacy groups attempts by (student, concept) in timestamp order
// and derives the pass-rate metrics.
func computeEfficacy(attempts []storage.AttemptEntry) *EfficacyReport {
	rep := &EfficacyReport{TotalAttempts: len(attempts)}
	if len(attempts) == 0 {
		return rep
	}
	type seq struct{ correct []bool }
	byKey := make(map[string]*seq)
	var order []string
	for _, a := range attempts {
		key := a.StudentID + "|" + a.ConceptID
		s, ok := byKey[key]
		if !ok {
			s = &seq{}
			byKey[key] = s
			order = append(order, key)
		}
		s.correct = append(s.correct, a.Correct)
	}
	firstPass, secondPass := 0, 0
	totalAttempts := 0
	for _, key := range order {
		s := byKey[key]
		rep.ConceptsTouched++
		if len(s.correct) > 0 && s.correct[0] {
			firstPass++
		}
		if (len(s.correct) >= 1 && s.correct[0]) || (len(s.correct) >= 2 && s.correct[1]) {
			secondPass++
		}
		totalAttempts += len(s.correct)
	}
	rep.AvgAttemptsPerConcept = float64(totalAttempts) / float64(rep.ConceptsTouched)
	if rep.ConceptsTouched > 0 {
		rep.FirstPassRate = float64(firstPass) / float64(rep.ConceptsTouched)
		rep.SecondPassRate = float64(secondPass) / float64(rep.ConceptsTouched)
	}
	return rep
}

// ShareReport is the read-only parent/teacher view of a student.
type ShareReport struct {
	StudentID string                              `json:"student_id"`
	Name      string                              `json:"name"`
	Scores    *scoring.Scores                     `json:"scores"`
	Activity  []storage.DailyActivity             `json:"activity"`
	Progress  map[string]*storage.ConceptProgress `json:"progress"`
	Weakness  map[string]float64                  `json:"weakness"`
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
	}, nil
}

func (e *Engine) StartDiagnostic() *diagnostic.Session {
	return e.diag.Start()
}

func (e *Engine) NextDiagnosticQuestion(s *diagnostic.Session) (*generator.Problem, string, error) {
	return e.diag.NextQuestion(s)
}

func (e *Engine) SubmitDiagnosticAnswer(s *diagnostic.Session, conceptID string, correct, fast bool) {
	e.diag.RecordAnswer(s, conceptID, correct, fast)
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
	}, nil
}

func (e *Engine) ApplyGoalResults(studentID string, session *diagnostic.Session) error {
	session.Lock()
	attempts := make([]diagnostic.Attempt, len(session.Attempts))
	copy(attempts, session.Attempts)
	session.Unlock()
	batch := make([]*storage.ConceptProgress, 0, len(attempts))
	for _, att := range attempts {
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
		batch = append(batch, &storage.ConceptProgress{
			StudentID:     studentID,
			ConceptID:     att.ConceptID,
			Status:        status,
			WeaknessScore: weakness,
		})
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
