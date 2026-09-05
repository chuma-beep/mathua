package storage

import (
	"crypto/rand"
	"fmt"
	"time"
)

// Data types

type Student struct {
	ID                  string
	Name                string
	Username            string
	PasswordHash        string
	CourseID            string
	XPTotal             int
	XPToday             int
	XPTodayDate         string
	DiagnosticCompleted bool
	DailyXPGoal         int
	Settings            string
	CreatedAt           time.Time
	League              string
	LeagueWeek          string
	LeagueMoved         int
	ShareToken          string
	Email               string
	GoogleID            string
	AvatarURL           string
}

type ConceptProgress struct {
	StudentID       string     `json:"student_id,omitempty"`
	ConceptID       string     `json:"concept_id,omitempty"`
	Status          string     `json:"status"` // UNSEEN | LEARNING | PRACTICING | MASTERED
	Streak          int        `json:"streak"`
	BestStreak      int        `json:"best_streak"`
	AvgResponseTime float64    `json:"avg_response_time"`
	Attempts        int        `json:"attempts"`
	LastAttempted   *time.Time `json:"last_attempted,omitempty"`
	LastReviewed    *time.Time `json:"last_reviewed,omitempty"`
	NextReviewDue   *time.Time `json:"next_review_due,omitempty"`
	SM2Repetitions  int        `json:"sm2_repetitions,omitempty"`
	SM2Interval     int        `json:"sm2_interval,omitempty"`
	SM2EFactor      float64    `json:"sm2_efactor,omitempty"`
	MasteredAt      *time.Time `json:"mastered_at,omitempty"`
	WeaknessScore   float64    `json:"weakness_score,omitempty"`
}

type Session struct {
	ID        string
	StudentID string
	StartedAt time.Time
}

type ActiveSession struct {
	SessionID      string
	StudentID      string
	ConceptID      string
	ConceptName    string
	ExpectedAnswer string
	AttemptID      string
	Question       string
	Explanation    string
	Diagram        string
	IsReview       bool
	Answered       bool
	LastConceptID  string
	SessionReview  int
	SessionNew     int
	UpdatedAt      time.Time
}

type AttemptEntry struct {
	SessionID      string
	StudentID      string
	ConceptID      string
	Answer         string
	Expected       string
	Correct        bool
	ElapsedSeconds float64
	Timestamp      time.Time
}

type Question struct {
	ID          int
	ConceptID   string
	Question    string
	Answer      string
	Explanation string
	Source      string
	Difficulty  float64
}

type LeaderboardRow struct {
	StudentID      string
	Name           string
	TotalMastered  int
	WeeklyMastered int
}

// LeagueMember is a student's row inside a weekly league tier.
type LeagueMember struct {
	StudentID      string
	Name           string
	Tier           string
	TotalMastered  int
	WeeklyMastered int
	Moved          int // +1 promoted, -1 demoted, 0 stayed (last reset)
}

type DailyActivity struct {
	Date      string   `json:"date"`
	Questions int      `json:"questions"`
	Correct   int      `json:"correct"`
	Concepts  []string `json:"concepts"`
}

type TopicSpeed struct {
	StudentID     string  `json:"student_id"`
	ConceptID     string  `json:"concept_id"`
	EFactor       float64 `json:"efactor"`
	Interval      int     `json:"interval"`
	Repetitions   int     `json:"repetitions"`
	LearningSpeed float64 `json:"learning_speed"`
}

// Repository interface

type Repository interface {
	CreateStudent(name string) (*Student, error)
	GetStudent(id string) (*Student, error)
	FindByUsername(username string) (*Student, error)
	FindByGoogleID(googleID string) (*Student, error)
	FindByEmail(email string) (*Student, error)
	CreateUser(name, username, passwordHash string) (*Student, error)
	CreateGoogleUser(name, email, googleID, avatarURL string) (*Student, error)
	LinkGoogleID(studentID, googleID, avatarURL string) error
	SetCourseID(studentID, courseID string) error
	UpdateStudentName(studentID, name string) error

	GetProgress(studentID, conceptID string) (*ConceptProgress, error)
	GetAllProgress(studentID string) (map[string]*ConceptProgress, error)
	UpsertProgress(p *ConceptProgress) error

	CreateSession(studentID string) (*Session, error)
	GetSession(id string) (*Session, error)
	GetActiveSession(sessionID string) (*ActiveSession, error)
	UpsertActiveSession(a *ActiveSession) error
	DeleteActiveSession(sessionID string) error
	RecordAttempt(entry AttemptEntry) error
	GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error)
	GetAttemptsForStudent(studentID string) ([]AttemptEntry, error)
	GetAllAttempts() ([]AttemptEntry, error)

	GetQuestions(conceptID string, count int) ([]Question, error)
	GetQuestionCount(conceptID string) (int, error)
	ImportQuestions(qs []Question) error
	PurgeGeneratedQuestions(conceptIDs map[string]bool) (int64, error)

	GetWeeklyLeaderboard() ([]LeaderboardRow, error)
	GetLeagueStandings() ([]LeagueMember, error)
	SetLeague(studentID, tier, week string, moved int) error
	SetShareToken(studentID, token string) error
	GetStudentByShareToken(token string) (*Student, error)
	GetDailyActivity(studentID string, days int) ([]DailyActivity, error)

	AddXP(studentID string, amount int) error
	GetXP(studentID string) (total int, today int, err error)
	SetDiagnosticCompleted(studentID string) error
	SetDailyXPGoal(studentID string, goal int) error
	GetSettings(studentID string) (string, error)
	UpdateSettings(studentID string, settings string) error
	SetAvatarImage(studentID, contentType string, data []byte) error
	GetAvatarImage(studentID string) (contentType string, data []byte, found bool, err error)
	ClearAvatarImage(studentID string) error

	GetTopicSpeed(studentID, conceptID string) (*TopicSpeed, error)
	GetAllTopicSpeeds(studentID string) (map[string]*TopicSpeed, error)
	UpsertTopicSpeed(ts *TopicSpeed) error

	Migrate() error
	Close() error
}

// UUID

func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("uuid-%x", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
