package engine

import (
	"crypto/rand"
	"encoding/hex"
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
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/scoring"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/planning"
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
}

// ErrNoActiveQuestion is returned when a session has no unanswered question
// ready for grading — e.g. a duplicate/stale submission racing the current one.
var ErrNoActiveQuestion = fmt.Errorf("no active question")

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
}

func New(repo storage.Repository, dag *concepts.DAG, reg *generator.Registry, ll *lessons.Loader, planner *planning.Planner) *Engine {
	return &Engine{
		dag:        dag,
		repo:       repo,
		sched:      scheduler.New(dag),
		registry:   reg,
		gr:         grader.NewRouter(),
		machine:    &mastery.Machine{},
		scorer:     scoring.NewUpdater(dag, repo),
		lboard:     leaderboard.NewComputer(repo),
		diag:       diagnostic.NewEngine(dag, reg),
		ll:         ll,
		planner:    planner,
		sessions:   make(map[string]*activeSession),
		activePath: make(map[string]map[string]bool),
	}
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
	weakness := 0.5
	if m := e.WeaknessMap(studentID); m != nil {
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
	next := e.sched.Next(snapshots, as.lastConceptID, as.sessionReview, as.sessionNew)
	if next == nil {
		return nil, nil
	}
	difficulty := e.computeDifficulty(studentID, next.Concept.ID)
	prevQuestion := as.questionText
	attemptID := newAttemptID()
	seedBase := hashSeed(studentID + "|" + next.Concept.ID + "|" + attemptID + fmt.Sprintf("|%d", time.Now().UnixNano()))
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
	as.timeThreshold = next.Concept.MasteryThreshold.AvgTimeSeconds
	as.isReview = next.IsReview
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
	difficulty := e.computeDifficulty(studentID, next.Concept.ID)
	prevQuestion := as.questionText
	attemptID := newAttemptID()
	seedBase := hashSeed(studentID + "|" + next.Concept.ID + "|" + attemptID + fmt.Sprintf("|%d", time.Now().UnixNano()))
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
	as.timeThreshold = next.Concept.MasteryThreshold.AvgTimeSeconds
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
	// Integrals
	"calc.integral.definite":          "/diagrams/algebrica/definite-integrals-1.svg",
	"calc.integral.ftc":               "/diagrams/algebrica/fundamental-theorem-of-calculus-1.svg",
	"calc.integral.area_between":      "/diagrams/algebrica/finding-areas-by-integration-1.svg",
	"calc.integral.volume":            "/diagrams/algebrica/finding-areas-by-integration-2.svg",

	// Equations
	"alg.quad.solve_factor": "/diagrams/algebrica/quadratic-equations.svg",
	"alg.quad.formula":      "/diagrams/algebrica/quadratic-equations.svg",
	// Calculus limits
	"calc.limit.continuity":  "/diagrams/algebrica/riemann-integrability-criteria-1.svg",
	"calc.integral.improper": "/diagrams/algebrica/improper-integrals-1.svg",
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
	// Hold e.mu for the whole grade + advance cycle so a concurrent
	// SubmitAnswer/NextQuestion pair can never grade against a question the
	// client was not shown. The `answered` flag rejects duplicate/stale
	// submissions for the same question, and the per-question attemptID
	// rejects submissions for questions the session has already moved past.
	e.mu.Lock()
	defer e.mu.Unlock()

	as := e.sessions[sessionID]
	if as == nil {
		if persisted := e.rehydrateActiveSession(sessionID, studentID); persisted != nil {
			as = persisted
			e.sessions[sessionID] = as
		}
	}
	if as == nil || as.conceptID == "" || as.answered || as.attemptID != attemptID {
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
	if xp > 0 {
		if err := e.repo.AddXP(studentID, xp); err != nil {
			log.Printf("warning: failed to add XP for student %s: %v", studentID, err)
		}
	}

	as.lastConceptID = as.conceptID
	if as.isReview {
		as.sessionReview++
	} else {
		as.sessionNew++
	}
	as.answered = true
	// Keep attemptID unique to avoid UNIQUE constraint collision on tombstone.
	// as.attemptID remains the last question's ID (unique per session).
	as.conceptID = ""
	as.questionText = ""
	e.persistActiveSession(sessionID, studentID, as, "")

	return &AnswerResult{
		Correct:        gr.Correct,
		Feedback:       gr.Feedback,
		NewStatus:      newStatus,
		Explanation:    explanation,
		Streak:         progress.Streak,
		RequiredStreak: sessionFields.requiredStreak,
		XP:             xp,
		ExpectedAnswer: sessionFields.expectedAnswer,
	}, nil
}

// SubmitStudyAnswer records a Study-library answer (LessonQuiz seam per CONTEXT.md: Seam).
// It grades via the concept's grading_type, updates mastery/SM-2/weakness/XP, and
// awards TaskMultistep 15 for *.word else TaskLesson 10 (Q2 lock).
func (e *Engine) SubmitStudyAnswer(studentID, conceptID, answer, expected string, elapsedSeconds float64) (*AnswerResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	c := e.dag.Concept(conceptID)
	requiredStreak := 3
	timeThreshold := 10.0
	if c != nil {
		requiredStreak = c.MasteryThreshold.Streak
		timeThreshold = c.MasteryThreshold.AvgTimeSeconds
	}
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
		SessionID:      "study",
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
	taskType := TaskLesson
	if strings.HasSuffix(conceptID, ".word") {
		taskType = TaskMultistep
	}
	xp := computeXPForTask(gr.Correct, elapsedSeconds, timeThreshold, progress.Streak, taskType)
	if xp > 0 {
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

func (e *Engine) StartDiagnostic() *diagnostic.Session {
	return e.diag.Start()
}

func (e *Engine) NextDiagnosticQuestion(s *diagnostic.Session) (*generator.Problem, string, error) {
	return e.diag.NextQuestion(s)
}

func (e *Engine) SubmitDiagnosticAnswer(s *diagnostic.Session, conceptID string, correct, fast bool) {
	e.diag.RecordAnswer(s, conceptID, correct, fast)
}

func (e *Engine) IsDiagnosticComplete(s *diagnostic.Session) bool {
	return e.diag.IsComplete(s)
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
		prog := &storage.ConceptProgress{
			StudentID:     studentID,
			ConceptID:     att.ConceptID,
			Status:        status,
			WeaknessScore: weakness,
		}
		if err := e.repo.UpsertProgress(prog); err != nil {
			return fmt.Errorf("save diagnostic progress: %w", err)
		}
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

	difficulty := e.computeDifficulty(studentID, conceptID)
	attemptID := newAttemptID()
	seedBase := hashSeed(studentID + "|" + conceptID + "|" + attemptID + fmt.Sprintf("|%d", time.Now().UnixNano()))
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
	as.timeThreshold = c.MasteryThreshold.AvgTimeSeconds
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
	ID           string  `json:"id"`
	Label        string  `json:"label"`
	Unlocked     bool    `json:"unlocked"`
	MasteryPct   float64 `json:"mastery_pct"`
	Streak       int     `json:"streak"`
	Status       string  `json:"status"`
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
	weakness := e.WeaknessMap(studentID)
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
	for _, u := range updates {
		prog, err := e.repo.GetProgress(studentID, u.cid)
		if err != nil || prog == nil {
			prog = &storage.ConceptProgress{
				StudentID: studentID,
				ConceptID: u.cid,
				Status:    string(mastery.StatusUnseen),
			}
		}
		prog.WeaknessScore = u.w
		if err := e.repo.UpsertProgress(prog); err != nil {
			log.Printf("warning: failed to propagate weakness for %s/%s: %v", studentID, u.cid, err)
		}
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
