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

	GetWeeklyLeaderboard() ([]LeaderboardRow, error)

	AddXP(studentID string, amount int) error
	GetXP(studentID string) (total int, today int, err error)
	SetDiagnosticCompleted(studentID string) error

	Migrate() error
	Close() error
}

// UUID

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
