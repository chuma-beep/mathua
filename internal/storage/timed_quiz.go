package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// timedQuizAccess is the shared timed-quiz/review-event implementation. The
// two stores differ only in placeholder spelling, so every statement is
// written with ? and rebound for PostgreSQL.
type timedQuizAccess struct {
	db     *sql.DB
	rebind func(string) string
}

func (s *SQLiteStore) timedQuiz() *timedQuizAccess {
	return &timedQuizAccess{db: s.db, rebind: rebindQuestion}
}

func (s *PostgresStore) timedQuiz() *timedQuizAccess {
	return &timedQuizAccess{db: s.db, rebind: rebindDollars}
}

func (s *SQLiteStore) EnsureTimedQuizBaseline(studentID string) (TimedQuizBaseline, error) {
	return s.timedQuiz().ensureBaseline(studentID)
}
func (s *PostgresStore) EnsureTimedQuizBaseline(studentID string) (TimedQuizBaseline, error) {
	return s.timedQuiz().ensureBaseline(studentID)
}
func (s *SQLiteStore) GetTimedQuizBaseline(studentID string) (*TimedQuizBaseline, error) {
	return s.timedQuiz().getBaseline(studentID)
}
func (s *PostgresStore) GetTimedQuizBaseline(studentID string) (*TimedQuizBaseline, error) {
	return s.timedQuiz().getBaseline(studentID)
}
func (s *SQLiteStore) TimedQuizEligibleXP(studentID string, since time.Time) (int, error) {
	return s.timedQuiz().eligibleXP(studentID, since)
}
func (s *PostgresStore) TimedQuizEligibleXP(studentID string, since time.Time) (int, error) {
	return s.timedQuiz().eligibleXP(studentID, since)
}
func (s *SQLiteStore) CreateTimedQuizRun(run TimedQuizRun, items []TimedQuizItem) error {
	return s.timedQuiz().createRun(run, items)
}
func (s *PostgresStore) CreateTimedQuizRun(run TimedQuizRun, items []TimedQuizItem) error {
	return s.timedQuiz().createRun(run, items)
}
func (s *SQLiteStore) GetActiveTimedQuizRun(studentID string) (*TimedQuizRun, error) {
	return s.timedQuiz().activeRun(studentID)
}
func (s *PostgresStore) GetActiveTimedQuizRun(studentID string) (*TimedQuizRun, error) {
	return s.timedQuiz().activeRun(studentID)
}
func (s *SQLiteStore) GetTimedQuizRun(id, studentID string) (*TimedQuizRun, error) {
	return s.timedQuiz().getRun(id, studentID)
}
func (s *PostgresStore) GetTimedQuizRun(id, studentID string) (*TimedQuizRun, error) {
	return s.timedQuiz().getRun(id, studentID)
}
func (s *SQLiteStore) LatestTimedQuizRun(studentID string) (*TimedQuizRun, error) {
	return s.timedQuiz().latestRun(studentID)
}
func (s *PostgresStore) LatestTimedQuizRun(studentID string) (*TimedQuizRun, error) {
	return s.timedQuiz().latestRun(studentID)
}
func (s *SQLiteStore) GetTimedQuizItems(runID string) ([]TimedQuizItem, error) {
	return s.timedQuiz().items(runID)
}
func (s *PostgresStore) GetTimedQuizItems(runID string) ([]TimedQuizItem, error) {
	return s.timedQuiz().items(runID)
}
func (s *SQLiteStore) ClaimTimedQuizItem(itemID, runID, submissionID string) (bool, error) {
	return s.timedQuiz().claimItem(itemID, runID, submissionID)
}
func (s *PostgresStore) ClaimTimedQuizItem(itemID, runID, submissionID string) (bool, error) {
	return s.timedQuiz().claimItem(itemID, runID, submissionID)
}
func (s *SQLiteStore) UpdateTimedQuizItem(itemID, runID, answer, submissionID string, correct bool, elapsed float64, xp int, status string) (bool, error) {
	return s.timedQuiz().updateItem(itemID, runID, answer, submissionID, correct, elapsed, xp, status, 0, "", "", "")
}
func (s *PostgresStore) UpdateTimedQuizItem(itemID, runID, answer, submissionID string, correct bool, elapsed float64, xp int, status string) (bool, error) {
	return s.timedQuiz().updateItem(itemID, runID, answer, submissionID, correct, elapsed, xp, status, 0, "", "", "")
}
func (s *SQLiteStore) UpdateTimedQuizItemWithAdjustment(itemID, runID, answer, submissionID string, correct bool, elapsed float64, xp int, status string, xpAdjustment int, xpNote, xpReasonCode, xpSource string) (bool, error) {
	return s.timedQuiz().updateItem(itemID, runID, answer, submissionID, correct, elapsed, xp, status, xpAdjustment, xpNote, xpReasonCode, xpSource)
}
func (s *PostgresStore) UpdateTimedQuizItemWithAdjustment(itemID, runID, answer, submissionID string, correct bool, elapsed float64, xp int, status string, xpAdjustment int, xpNote, xpReasonCode, xpSource string) (bool, error) {
	return s.timedQuiz().updateItem(itemID, runID, answer, submissionID, correct, elapsed, xp, status, xpAdjustment, xpNote, xpReasonCode, xpSource)
}
func (s *SQLiteStore) FinalizeTimedQuizRun(runID, studentID string, now time.Time) (*TimedQuizRun, error) {
	return s.timedQuiz().finalize(runID, studentID, now)
}
func (s *PostgresStore) FinalizeTimedQuizRun(runID, studentID string, now time.Time) (*TimedQuizRun, error) {
	return s.timedQuiz().finalize(runID, studentID, now)
}
func (s *SQLiteStore) RecordReviewEvent(event ReviewEvent) error {
	return s.timedQuiz().recordReviewEvent(event)
}
func (s *PostgresStore) RecordReviewEvent(event ReviewEvent) error {
	return s.timedQuiz().recordReviewEvent(event)
}

func (a *timedQuizAccess) ensureBaselineAt(studentID string, at time.Time) error {
	if _, err := a.db.Exec(a.rebind(`INSERT INTO timed_quiz_baselines (student_id, baseline_xp, baseline_at)
		VALUES (?, 0, ?) ON CONFLICT (student_id) DO NOTHING`), studentID, at.UnixNano()); err != nil {
		return fmt.Errorf("ensure timed quiz baseline: %w", err)
	}
	return nil
}

func (a *timedQuizAccess) ensureBaseline(studentID string) (TimedQuizBaseline, error) {
	now := time.Now().UTC()
	if err := a.ensureBaselineAt(studentID, now); err != nil {
		return TimedQuizBaseline{}, err
	}
	b, err := a.getBaseline(studentID)
	if err != nil {
		return TimedQuizBaseline{}, err
	}
	if b == nil {
		return TimedQuizBaseline{}, fmt.Errorf("timed quiz baseline missing after insert")
	}
	return *b, nil
}

func (a *timedQuizAccess) getBaseline(studentID string) (*TimedQuizBaseline, error) {
	var b TimedQuizBaseline
	var at int64
	err := a.db.QueryRow(a.rebind(`SELECT student_id, baseline_xp, baseline_at
		FROM timed_quiz_baselines WHERE student_id = ?`), studentID).Scan(&b.StudentID, &b.BaselineXP, &at)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get timed quiz baseline: %w", err)
	}
	b.BaselineAt = time.Unix(0, at).UTC()
	return &b, nil
}

func (a *timedQuizAccess) eligibleXP(studentID string, since time.Time) (int, error) {
	var total int
	err := a.db.QueryRow(a.rebind(`SELECT COALESCE(SUM(amount), 0) FROM xp_events
		WHERE student_id = ? AND created_at >= ? AND amount > 0`), studentID, since.UnixNano()).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("sum timed quiz eligible XP: %w", err)
	}
	return total, nil
}

func (a *timedQuizAccess) createRun(run TimedQuizRun, items []TimedQuizItem) error {
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("begin timed quiz run: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.Exec(a.rebind(`INSERT INTO timed_quiz_runs
		(id, student_id, started_at, deadline_at, completed_at, status, question_count,
		 correct_count, unanswered_count, xp_earned, created_at)
		VALUES (?, ?, ?, ?, 0, ?, ?, 0, 0, 0, ?)`),
		run.ID, run.StudentID, run.StartedAt.UnixNano(), run.DeadlineAt.UnixNano(),
		run.Status, len(items), run.CreatedAt.UnixNano()); err != nil {
		return fmt.Errorf("insert timed quiz run: %w", err)
	}
	for _, item := range items {
		if _, err := tx.Exec(a.rebind(`INSERT INTO timed_quiz_items
			(id, run_id, concept_id, position, question, expected_answer, explanation,
			 grading_type, difficulty, served_at, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`),
			item.ID, run.ID, item.ConceptID, item.Position, item.Question, item.ExpectedAnswer,
			item.Explanation, item.GradingType, nullableFloat(item.Difficulty), item.ServedAt.UnixNano(),
			item.Status); err != nil {
			return fmt.Errorf("insert timed quiz item: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit timed quiz run: %w", err)
	}
	return nil
}

func (a *timedQuizAccess) activeRun(studentID string) (*TimedQuizRun, error) {
	return a.queryRun(`SELECT id, student_id, started_at, deadline_at, completed_at, status,
		question_count, correct_count, unanswered_count, xp_earned, created_at
		FROM timed_quiz_runs WHERE student_id = ? AND status = 'active'
		ORDER BY started_at DESC LIMIT 1`, studentID)
}

func (a *timedQuizAccess) getRun(id, studentID string) (*TimedQuizRun, error) {
	return a.queryRun(`SELECT id, student_id, started_at, deadline_at, completed_at, status,
		question_count, correct_count, unanswered_count, xp_earned, created_at
		FROM timed_quiz_runs WHERE id = ? AND student_id = ?`, id, studentID)
}

func (a *timedQuizAccess) latestRun(studentID string) (*TimedQuizRun, error) {
	return a.queryRun(`SELECT id, student_id, started_at, deadline_at, completed_at, status,
		question_count, correct_count, unanswered_count, xp_earned, created_at
		FROM timed_quiz_runs WHERE student_id = ? AND status = 'completed'
		ORDER BY completed_at DESC LIMIT 1`, studentID)
}

func (a *timedQuizAccess) queryRun(query string, args ...any) (*TimedQuizRun, error) {
	var run TimedQuizRun
	var started, deadline, completed, created int64
	err := a.db.QueryRow(a.rebind(query), args...).Scan(
		&run.ID, &run.StudentID, &started, &deadline, &completed, &run.Status,
		&run.QuestionCount, &run.CorrectCount, &run.UnansweredCount, &run.XPEarned, &created,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get timed quiz run: %w", err)
	}
	run.StartedAt = time.Unix(0, started).UTC()
	run.DeadlineAt = time.Unix(0, deadline).UTC()
	run.CreatedAt = time.Unix(0, created).UTC()
	if completed > 0 {
		t := time.Unix(0, completed).UTC()
		run.CompletedAt = &t
	}
	return &run, nil
}

func (a *timedQuizAccess) items(runID string) ([]TimedQuizItem, error) {
	rows, err := a.db.Query(a.rebind(`SELECT id, run_id, concept_id, position, question,
		expected_answer, explanation, grading_type, difficulty, served_at, answered_at,
		elapsed_seconds, answer, correct, status, submission_id, xp_earned,
		xp_adjustment, xp_note, xp_reason_code, xp_source
		FROM timed_quiz_items WHERE run_id = ? ORDER BY position`), runID)
	if err != nil {
		return nil, fmt.Errorf("list timed quiz items: %w", err)
	}
	defer rows.Close()
	var out []TimedQuizItem
	for rows.Next() {
		var it TimedQuizItem
		var served, answered int64
		var diff sql.NullFloat64
		var correct int
		if err := rows.Scan(&it.ID, &it.RunID, &it.ConceptID, &it.Position, &it.Question,
			&it.ExpectedAnswer, &it.Explanation, &it.GradingType, &diff, &served, &answered,
			&it.ElapsedSeconds, &it.Answer, &correct, &it.Status, &it.SubmissionID, &it.XPEarned,
			&it.XPAdjustment, &it.XPNote, &it.XPReasonCode, &it.XPSource); err != nil {
			return nil, fmt.Errorf("scan timed quiz item: %w", err)
		}
		it.ServedAt = time.Unix(0, served).UTC()
		if answered > 0 {
			t := time.Unix(0, answered).UTC()
			it.AnsweredAt = &t
		}
		if diff.Valid {
			v := diff.Float64
			it.Difficulty = &v
		}
		it.Correct = correct != 0
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate timed quiz items: %w", err)
	}
	return out, nil
}

func (a *timedQuizAccess) claimItem(itemID, runID, submissionID string) (bool, error) {
	res, err := a.db.Exec(a.rebind(`UPDATE timed_quiz_items
		SET status = 'submitting', submission_id = ?
		WHERE id = ? AND run_id = ? AND status = 'served'`), submissionID, itemID, runID)
	if err != nil {
		return false, fmt.Errorf("claim timed quiz item: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("timed quiz claim rows affected: %w", err)
	}
	return n == 1, nil
}

func (a *timedQuizAccess) updateItem(itemID, runID, answer, submissionID string, correct bool, elapsed float64, xp int, status string, xpAdjustment int, xpNote, xpReasonCode, xpSource string) (bool, error) {
	correctInt := 0
	if correct {
		correctInt = 1
	}
	res, err := a.db.Exec(a.rebind(`UPDATE timed_quiz_items
		SET answer = ?, submission_id = ?, correct = ?, elapsed_seconds = ?, xp_earned = ?,
			xp_adjustment = ?, xp_note = ?, xp_reason_code = ?, xp_source = ?,
			status = ?, answered_at = ?
		WHERE id = ? AND run_id = ? AND submission_id = ? AND status IN ('served', 'submitting')`),
		answer, submissionID, correctInt, elapsed, xp, xpAdjustment, xpNote, xpReasonCode, xpSource,
		status, time.Now().UTC().UnixNano(), itemID, runID, submissionID)
	if err != nil {
		return false, fmt.Errorf("update timed quiz item: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("timed quiz item rows affected: %w", err)
	}
	return n == 1, nil
}

func (a *timedQuizAccess) finalize(runID, studentID string, now time.Time) (*TimedQuizRun, error) {
	_, err := a.db.Exec(a.rebind(`UPDATE timed_quiz_runs SET
		status = 'completed', completed_at = ?,
		question_count = (SELECT COUNT(*) FROM timed_quiz_items WHERE run_id = ?),
		correct_count = (SELECT COUNT(*) FROM timed_quiz_items WHERE run_id = ? AND correct = 1),
		unanswered_count = (SELECT COUNT(*) FROM timed_quiz_items WHERE run_id = ? AND status <> 'answered'),
		xp_earned = (SELECT COALESCE(SUM(xp_earned), 0) FROM timed_quiz_items WHERE run_id = ?)
		WHERE id = ? AND student_id = ? AND status = 'active'`),
		now.UnixNano(), runID, runID, runID, runID, runID, studentID)
	if err != nil {
		return nil, fmt.Errorf("finalize timed quiz run: %w", err)
	}
	run, err := a.getRun(runID, studentID)
	if err != nil {
		return nil, err
	}
	if run == nil {
		return nil, fmt.Errorf("timed quiz run not found")
	}
	return run, nil
}

func (a *timedQuizAccess) recordReviewEvent(e ReviewEvent) error {
	_, err := a.db.Exec(a.rebind(`INSERT INTO review_events
		(student_id, concept_id, attempt_id, rating, correct, elapsed_seconds, reviewed_at,
		 stability_before, difficulty_before, due_before, stability_after, difficulty_after,
		 due_after, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (attempt_id) DO NOTHING`),
		e.StudentID, e.ConceptID, e.AttemptID, e.Rating, boolInt(e.Correct), e.ElapsedSeconds,
		e.ReviewedAt.UnixNano(), nullableFloat(e.StabilityBefore), nullableFloat(e.DifficultyBefore),
		nullableTime(e.DueBefore), nullableFloat(e.StabilityAfter), nullableFloat(e.DifficultyAfter),
		nullableTime(e.DueAfter), time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("record review event: %w", err)
	}
	return nil
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nullableFloat(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullableTime(t *time.Time) any {
	if t == nil || t.IsZero() {
		return nil
	}
	return t.UTC().Format(time.RFC3339Nano)
}
