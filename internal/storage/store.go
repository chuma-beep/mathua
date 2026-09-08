package storage

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
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
	EmailVerified       bool
	GoogleID            string
	AvatarURL           string
}

// Identity links an external OAuth provider account to a student.
// (provider, provider_id) is globally unique across all providers.
type Identity struct {
	Provider      string
	ProviderID    string
	StudentID     string
	Email         string
	EmailVerified bool
	LinkedAt      time.Time
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
	Username       string
	AvatarURL      string
	AvatarSettings string
	TotalMastered  int
	WeeklyMastered int
}

// LeagueMember is a student's row inside a weekly league tier.
// The snake_case tags matter: this struct is serialized straight into
// GET /api/leagues, and without them the keys came out PascalCase, which
// the frontend could not read — league names never rendered.
type LeagueMember struct {
	StudentID      string     `json:"student_id"`
	Name           string     `json:"name"`
	Username       string     `json:"username,omitempty"`
	AvatarURL      string     `json:"avatar_url,omitempty"`
	AvatarDicebear *AvatarRef `json:"avatar_dicebear,omitempty"`
	AvatarCustom   bool       `json:"avatar_custom,omitempty"`
	AvatarVersion  int        `json:"avatar_version,omitempty"`
	AvatarSettings string     `json:"-"`
	Tier           string     `json:"tier"`
	TotalMastered  int        `json:"total_mastered"`
	WeeklyMastered int        `json:"weekly_mastered"`
	Moved          int        `json:"moved"` // +1 promoted, -1 demoted, 0 stayed (last reset)
}

// AvatarRef is the public dicebear pick surfaced on leaderboard entries.
type AvatarRef struct {
	Style string `json:"style"`
	Seed  string `json:"seed"`
}

// avatarSettingsJSON mirrors the opaque settings keys the frontend owns
// (web/next-app UserSettings); only leaderboard-relevant bits are decoded.
type avatarSettingsJSON struct {
	AvatarCustom   bool       `json:"avatar_custom"`
	AvatarDicebear *AvatarRef `json:"avatar_dicebear"`
	AvatarVersion  int        `json:"avatar_version"`
}

// AvatarBits extracts leaderboard-relevant avatar data from a student's
// opaque settings blob: custom-upload flag, dicebear pick, photo version.
// Corrupt/empty blobs yield zero values (initial-avatar fallback downstream).
func AvatarBits(raw string) (custom bool, pick *AvatarRef, version int) {
	if strings.TrimSpace(raw) == "" {
		return false, nil, 0
	}
	var s avatarSettingsJSON
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return false, nil, 0
	}
	return s.AvatarCustom, s.AvatarDicebear, s.AvatarVersion
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

// QuestionReport is a user complaint about a question, explanation,
// lesson body, worked example, or diagram.
type QuestionReport struct {
	ID          int64     `json:"id"`
	ReporterID  string    `json:"reporter_id"`
	ConceptID   string    `json:"concept_id"`
	Kind        string    `json:"kind"` // question | explanation | lesson_body | worked_example | diagram
	Question    string    `json:"question"`
	Expected    string    `json:"expected"`
	Explanation string    `json:"explanation"`
	LessonID    string    `json:"lesson_id"`
	Source      string    `json:"source"`
	SessionID   string    `json:"session_id"`
	AttemptID   string    `json:"attempt_id"`
	Reason      string    `json:"reason"` // wrong_answer | bad_explanation | unclear | formatting | other
	Detail      string    `json:"detail"`
	Status      string    `json:"status"` // open | confirmed | fixed | dismissed
	CreatedAt   time.Time `json:"created_at"`
}

// Repository interface

type Repository interface {
	CreateStudent(name string) (*Student, error)
	GetStudent(id string) (*Student, error)
	FindByUsername(username string) (*Student, error)
	FindByEmail(email string) (*Student, error)
	CreateUser(name, username, passwordHash string) (*Student, error)
	SetUsername(studentID, username string) error
	ListStudentsMissingUsernames() ([]string, error)
	CreateGoogleUser(name, email, googleID, avatarURL string) (*Student, error)
	CreateOAuthUser(provider, providerID, name, email string, emailVerified bool, avatarURL string) (*Student, error)
	SetAvatarURL(studentID, avatarURL string) error
	SetCourseID(studentID, courseID string) error
	UpdateStudentName(studentID, name string) error
	SetEmail(studentID, email string) error
	SetEmailVerified(studentID string, verified bool) error
	SetPasswordHash(studentID, hash string) error
	CreatePasswordReset(tokenHash, studentID string, expiresAt time.Time) error
	ConsumePasswordReset(tokenHash string) (studentID string, ok bool, err error)

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

	CreateReport(r QuestionReport) (int64, error)
	ListReports(status string, limit, offset int) ([]QuestionReport, error)
	UpdateReportStatus(id int64, status string) error

	CreateIdentity(provider, providerID, studentID, email string, emailVerified bool) error
	FindStudentByIdentity(provider, providerID string) (*Student, error)
	ListIdentities(studentID string) ([]Identity, error)
	DeleteIdentity(provider, studentID string) error
	CreateEmailVerification(tokenHash, studentID string, expiresAt time.Time) error
	ConsumeEmailVerification(tokenHash string) (studentID string, ok bool, err error)
	CreateLinkToken(tokenHash, studentID string, expiresAt time.Time) error
	ConsumeLinkToken(tokenHash string) (studentID string, ok bool, err error)

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

// IsUniqueViolation reports UNIQUE-constraint failures across stores
// (SQLite "UNIQUE constraint failed", Postgres "duplicate key").
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "duplicate key")
}
