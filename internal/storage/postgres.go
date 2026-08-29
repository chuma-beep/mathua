package storage

import "fmt"

type PostgresStore struct{}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	return nil, fmt.Errorf("PostgreSQL backend not yet implemented")
}

func (s *PostgresStore) CreateStudent(name string) (*Student, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetStudent(id string) (*Student, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) FindByUsername(username string) (*Student, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) SetCourseID(studentID, courseID string) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) CreateUser(name, username, passwordHash string) (*Student, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetProgress(studentID, conceptID string) (*ConceptProgress, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetAllProgress(studentID string) (map[string]*ConceptProgress, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) UpsertProgress(p *ConceptProgress) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) CreateSession(studentID string) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetSession(id string) (*Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetActiveSession(sessionID string) (*ActiveSession, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) UpsertActiveSession(a *ActiveSession) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) DeleteActiveSession(sessionID string) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) RecordAttempt(entry AttemptEntry) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetSessionAttempts(studentID, sessionID string) ([]AttemptEntry, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetWeeklyLeaderboard() ([]LeaderboardRow, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) AddXP(studentID string, amount int) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetXP(studentID string) (int, int, error) {
	return 0, 0, fmt.Errorf("not implemented")
}

func (s *PostgresStore) SetDiagnosticCompleted(studentID string) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) SetDailyXPGoal(studentID string, goal int) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetSettings(studentID string) (string, error) {
	return "{}", fmt.Errorf("not implemented")
}

func (s *PostgresStore) UpdateSettings(studentID string, settings string) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetDailyActivity(studentID string, days int) ([]DailyActivity, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) Migrate() error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) Close() error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetQuestions(conceptID string, count int) ([]Question, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *PostgresStore) GetQuestionCount(conceptID string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (s *PostgresStore) ImportQuestions(qs []Question) error {
	return fmt.Errorf("not implemented")
}

func (s *PostgresStore) PurgeGeneratedQuestions(conceptIDs map[string]bool) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}
