package engine

import (
	"fmt"
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
	sessionReview  int
	sessionNew     int
	lastConceptID  string
}

type Question struct {
	ConceptID   string          `json:"concept_id"`
	ConceptName string          `json:"concept_name"`
	Question    string          `json:"question"`
	IsReview    bool            `json:"is_review"`
	Lesson      *lessons.Lesson `json:"lesson,omitempty"`
	Diagram     string          `json:"diagram,omitempty"`
}

type AnswerResult struct {
	Correct        bool           `json:"correct"`
	Feedback       string         `json:"feedback"`
	NewStatus      mastery.Status `json:"new_status"`
	Explanation    string         `json:"explanation,omitempty"`
	Streak         int            `json:"streak"`
	RequiredStreak int            `json:"required_streak"`
	XP             int            `json:"xp"`
	ExpectedAnswer string         `json:"expected_answer,omitempty"`
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
func (e *Engine) PlannerCourses() []*planning.Course {
	if e.planner == nil {
		return nil
	}
	return e.planner.Courses()
}

func (e *Engine) ActivePath(studentID string) map[string]bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.activePath[studentID]
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
	prob, err := e.registry.Generate(next.Concept.ID, 0.5)
	if err != nil {
		return nil, fmt.Errorf("generate problem: %w", err)
	}
	as.conceptID = next.Concept.ID
	as.conceptName = next.Concept.Label
	as.expectedAnswer = prob.Answer
	as.explanation = prob.Explanation
	as.requiredStreak = next.Concept.MasteryThreshold.Streak
	as.timeThreshold = next.Concept.MasteryThreshold.AvgTimeSeconds
	as.isReview = next.IsReview
	e.sessions[sessionID] = as

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
	"calc.integral.substitution":      "/diagrams/algebrica/integration-by-substitution.svg",
	"calc.integral.partial_fractions": "/diagrams/algebrica/integral-of-rational-functions.svg",
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

// gradeAnswer delegates to the generator's own grader if it implements
// GradedGenerator, otherwise falls back to the type-based router.
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

func (e *Engine) SubmitAnswer(sessionID, studentID string, answer string, elapsedSeconds float64) (*AnswerResult, error) {
	e.mu.Lock()
	as := e.sessions[sessionID]
	e.mu.Unlock()
	if as == nil {
		return nil, fmt.Errorf("no active question for session %q", sessionID)
	}

	gr := e.gradeAnswer(as.conceptID, as.expectedAnswer, answer)

	progress, err := e.repo.GetProgress(studentID, as.conceptID)
	if err != nil {
		return nil, fmt.Errorf("get progress: %w", err)
	}
	if progress == nil {
		progress = &storage.ConceptProgress{
			StudentID:  studentID,
			ConceptID:  as.conceptID,
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
	conceptThreshold := as.timeThreshold
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
		RequiredStreak:    as.requiredStreak,
		AvgResponseTime:   progress.AvgResponseTime,
		ResponseThreshold: as.timeThreshold,
	}
	newStatus := e.machine.Next(mastery.Status(progress.Status), ctx)
	oldStatus := progress.Status
	progress.Status = string(newStatus)

	// Reset streak when advancing to a new mastery level.
	if newStatus != mastery.Status(oldStatus) && gr.Correct {
		progress.Streak = 1
	}

	// SM-2 update
	quality := mastery.SM2Quality(
		gr.Correct && progress.Streak >= as.requiredStreak,
		progress.AvgResponseTime/as.timeThreshold,
	)
	prevSM2 := scheduler.SM2{
		Repetitions: progress.SM2Repetitions,
		EFactor:     progress.SM2EFactor,
		Interval:    progress.SM2Interval,
	}
	nextSM2 := scheduler.ComputeSM2(prevSM2, quality)
	progress.SM2Repetitions = nextSM2.Repetitions
	progress.SM2EFactor = nextSM2.EFactor
	progress.SM2Interval = nextSM2.Interval

	if newStatus == mastery.StatusMastered && mastery.Status(progress.Status) != mastery.StatusMastered {
		now := nowUTC()
		progress.MasteredAt = &now
	}
	if newStatus == mastery.StatusMastered {
		progress.WeaknessScore = 0.0
	}
	if progress.Streak >= as.requiredStreak && progress.AvgResponseTime <= as.timeThreshold {
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
		ConceptID:      as.conceptID,
		Answer:         answer,
		Expected:       as.expectedAnswer,
		Correct:        gr.Correct,
		ElapsedSeconds: elapsedSeconds,
		Timestamp:      nowUTC(),
	}); err != nil {
		return nil, fmt.Errorf("record attempt: %w", err)
	}

	explanation := ""
	if !gr.Correct {
		explanation = as.explanation
	}

	xp := computeXP(gr.Correct, elapsedSeconds, as.timeThreshold, progress.Streak, as.isReview)
	if xp > 0 {
		_ = e.repo.AddXP(studentID, xp)
	}

	e.mu.Lock()
	as.lastConceptID = as.conceptID
	if as.isReview {
		as.sessionReview++
	} else {
		as.sessionNew++
	}
	as.conceptID = ""
	e.mu.Unlock()

	return &AnswerResult{
		Correct:        gr.Correct,
		Feedback:       gr.Feedback,
		NewStatus:      newStatus,
		Explanation:    explanation,
		Streak:         progress.Streak,
		RequiredStreak: as.requiredStreak,
		XP:             xp,
		ExpectedAnswer: as.expectedAnswer,
	}, nil
}

func computeXP(correct bool, elapsed, timeThreshold float64, streak int, isReview bool) int {
	if !correct {
		return 0
	}
	base := 10
	if isReview {
		base = 5
	}
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

func (e *Engine) GetProgress(studentID string) (map[string]*storage.ConceptProgress, error) {
	return e.repo.GetAllProgress(studentID)
}

func (e *Engine) GetScores(studentID string) (*scoring.Scores, error) {
	return e.scorer.Compute(studentID)
}

func (e *Engine) GetDAG() *concepts.DAG {
	return e.dag
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
	for _, att := range session.Attempts {
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
		_ = e.repo.UpsertProgress(prog)
	}
}

func (e *Engine) AdjustPlan(studentID string) {
	e.mu.Lock()
	path := e.activePath[studentID]
	e.mu.Unlock()
	if len(path) == 0 {
		return
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return
	}
	e.mu.Lock()
	for cid := range path {
		if p, ok := progress[cid]; ok && p.Status == string(mastery.StatusMastered) && p.WeaknessScore < 0.2 {
			delete(path, cid)
		}
	}
	e.activePath[studentID] = path
	e.mu.Unlock()
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

func nowUTC() time.Time {
	return time.Now().UTC()
}
