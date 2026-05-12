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
}

type AnswerResult struct {
	Correct        bool           `json:"correct"`
	Feedback       string         `json:"feedback"`
	NewStatus      mastery.Status `json:"new_status"`
	Explanation    string         `json:"explanation,omitempty"`
	Streak         int            `json:"streak"`
	RequiredStreak int            `json:"required_streak"`
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

	mu       sync.Mutex
	sessions map[string]*activeSession
}

func New(repo storage.Repository, dag *concepts.DAG, reg *generator.Registry, ll *lessons.Loader) *Engine {
	return &Engine{
		dag:      dag,
		repo:     repo,
		sched:    scheduler.New(dag),
		registry: reg,
		gr:       grader.NewRouter(),
		machine:  &mastery.Machine{},
		scorer:   scoring.NewUpdater(dag, repo),
		lboard:   leaderboard.NewComputer(repo),
		diag:     diagnostic.NewEngine(dag, reg),
		ll:       ll,
		sessions: make(map[string]*activeSession),
	}
}

func (e *Engine) CreateStudent(name string) (*storage.Student, error) {
	return e.repo.CreateStudent(name)
}

func (e *Engine) GetStudent(id string) (*storage.Student, error) {
	return e.repo.GetStudent(id)
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

	return &Question{
		ConceptID:   next.Concept.ID,
		ConceptName: next.Concept.Label,
		Question:    prob.Question,
		IsReview:    next.IsReview,
		Lesson:      lesson,
	}, nil
}

func (e *Engine) SubmitAnswer(sessionID, studentID string, answer string, elapsedSeconds float64) (*AnswerResult, error) {
	e.mu.Lock()
	as := e.sessions[sessionID]
	e.mu.Unlock()
	if as == nil {
		return nil, fmt.Errorf("no active question for session %q", sessionID)
	}

	concept := e.dag.Concept(as.conceptID)
	gradingType := grader.GradingNumeric
	if concept != nil {
		gradingType = grader.GradingType(concept.GradingType)
	}
	gr := e.gr.Grade(gradingType, as.expectedAnswer, answer)

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
	if progress.Streak >= as.requiredStreak && progress.AvgResponseTime <= as.timeThreshold {
		now := nowUTC()
		progress.LastReviewed = &now
		nextReview := now.AddDate(0, 0, nextSM2.Interval)
		progress.NextReviewDue = &nextReview
	}

	if err := e.repo.UpsertProgress(progress); err != nil {
		return nil, fmt.Errorf("save progress: %w", err)
	}

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
	}, nil
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
