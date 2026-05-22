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
}

type ConceptProgress struct {
	StudentID       string
	ConceptID       string
	Status          string // UNSEEN | LEARNING | PRACTICING | MASTERED
	Streak          int
	BestStreak      int
	AvgResponseTime float64
	Attempts        int
	LastAttempted   *time.Time
	LastReviewed    *time.Time
	NextReviewDue   *time.Time
	SM2Repetitions  int
	SM2Interval     int
	SM2EFactor      float64
	MasteredAt      *time.Time
	WeaknessScore   float64
}

type Session struct {
	ID        string
	StudentID string
	StartedAt time.Time
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

// Repository interface

type Repository interface {
	CreateStudent(name string) (*Student, error)
	GetStudent(id string) (*Student, error)
	FindByUsername(username string) (*Student, error)
	CreateUser(name, username, passwordHash string) (*Student, error)
	SetCourseID(studentID, courseID string) error

	GetProgress(studentID, conceptID string) (*ConceptProgress, error)
	GetAllProgress(studentID string) (map[string]*ConceptProgress, error)
	UpsertProgress(p *ConceptProgress) error

	CreateSession(studentID string) (*Session, error)
	GetSession(id string) (*Session, error)
	RecordAttempt(entry AttemptEntry) error
	GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error)

	GetQuestions(conceptID string, count int) ([]Question, error)
	GetQuestionCount(conceptID string) (int, error)
	ImportQuestions(qs []Question) error

	GetWeeklyLeaderboard() ([]LeaderboardRow, error)

	AddXP(studentID string, amount int) error
	GetXP(studentID string) (total int, today int, err error)
	SetDiagnosticCompleted(studentID string) error
	SetDailyXPGoal(studentID string, goal int) error
	GetSettings(studentID string) (string, error)
	UpdateSettings(studentID string, settings string) error

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
