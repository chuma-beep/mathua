package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// EffortSignal is the durable, minimal record of one graded answer as seen by
// XP accountability. It deliberately stores no answer text, only a one-way
// signature, and no mastery verdict, so reading it can never change evidence.
//
// The table is an audit and detection source, not an XP system. The ledger stays
// the only place a balance moves.
type EffortSignal struct {
	StudentID       string
	TaskKind        string
	TaskID          string
	ReviewRunID     string
	ConceptID       string
	AttemptKey      string
	Correct         bool
	DontKnow        bool
	Mastered        bool
	LegacyPenalty   bool
	ElapsedSeconds  float64
	Difficulty      float64
	XPEarned        int
	AnswerSignature string
	CreatedAt       time.Time
}

// EffortStore is the optional persistence surface for XP accountability.
// Stores that do not implement it keep their old behaviour: the engine treats
// the layer as unavailable and awards XP exactly as before.
type EffortStore interface {
	// RecordEffortSignal appends one observation. It is idempotent on
	// (student, task kind, task id, attempt key), so a retried request cannot
	// double count.
	RecordEffortSignal(sig EffortSignal) error
	// EffortSignalsSince returns up to limit signals at or after since, newest
	// first. It is bounded because detection only needs a recent window.
	EffortSignalsSince(studentID string, since time.Time, limit int) ([]EffortSignal, error)
	// CountXPDeductions counts confirmed effort deductions at or after since.
	// It drives the escalating ladder and never counts ordinary awards or a
	// performance deduction, so the two dimensions stay separate.
	CountXPDeductions(studentID string, since time.Time) (int, error)
	// SumXPDeductionsSince returns the positive magnitude of every negative
	// ledger event at or after since, including legacy rush penalties and
	// accountability adjustments. It is an audit helper; progression gates
	// use SumXPEarnedSince instead.
	SumXPDeductionsSince(studentID string, since time.Time) (int, error)
	// SumXPEarnedSince returns the positive XP ledger total at or after since.
	// It is the earned-XP basis for progression gates, independent of the
	// current signed balance.
	SumXPEarnedSince(studentID string, since time.Time) (int, error)
	// ApplyXPAdjustment writes a signed ledger row and moves the counters in one
	// transaction. Negative amounts are applied in full, so lifetime XP may
	// become negative. It returns the amount actually applied (zero when the
	// event already exists), so a caller can explain the real outcome.
	ApplyXPAdjustment(studentID string, amount int, reason, eventKey string) (int, error)
}

// effortSignalColumns is shared by both stores so a new field cannot be added to
// one SELECT and forgotten in the other.
const effortSignalColumns = `student_id, task_kind, task_id, review_run_id, concept_id, attempt_key,
	correct, dont_know, mastered, legacy_penalty, elapsed_seconds, difficulty,
	xp_earned, answer_signature, created_at`

type effortExecer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// recordEffortSignalTx is the single INSERT both stores use. The unique key is
// the duplicate guard, so the caller does not need a read-before-write.
func recordEffortSignalTx(db *sql.DB, rebind func(string) string, sig EffortSignal) error {
	return recordEffortSignalExec(db, rebind, sig)
}

// recordEffortSignalExec accepts a DB or an existing transaction. The study
// path uses the transaction form so the signal and the answer commit together.
func recordEffortSignalExec(db effortExecer, rebind func(string) string, sig EffortSignal) error {
	if sig.StudentID == "" || sig.TaskKind == "" || sig.TaskID == "" || sig.AttemptKey == "" {
		return fmt.Errorf("effort signal: missing scope")
	}
	_, err := db.Exec(rebind(`
		INSERT INTO effort_signals
			(student_id, task_kind, task_id, review_run_id, concept_id, attempt_key, correct,
			 dont_know, mastered, legacy_penalty, elapsed_seconds, difficulty,
			 xp_earned, answer_signature, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(student_id, task_kind, task_id, attempt_key) DO NOTHING
	`), sig.StudentID, sig.TaskKind, sig.TaskID, sig.ReviewRunID, sig.ConceptID, sig.AttemptKey,
		boolToInt(sig.Correct), boolToInt(sig.DontKnow), boolToInt(sig.Mastered),
		boolToInt(sig.LegacyPenalty), sig.ElapsedSeconds, sig.Difficulty,
		sig.XPEarned, sig.AnswerSignature, signalTime(sig.CreatedAt))
	if err != nil {
		return fmt.Errorf("effort signal: %w", err)
	}
	return nil
}

func effortSignalsSince(db *sql.DB, rebind func(string) string, studentID string, since time.Time, limit int) ([]EffortSignal, error) {
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}
	rows, err := db.Query(rebind(`
		SELECT `+effortSignalColumns+`
		FROM effort_signals
		WHERE student_id = ? AND created_at >= ?
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`), studentID, signalTime(since), limit)
	if err != nil {
		return nil, fmt.Errorf("effort signals: %w", err)
	}
	defer rows.Close()

	var out []EffortSignal
	for rows.Next() {
		var sig EffortSignal
		var correct, dontKnow, mastered, legacy, created int
		if err := rows.Scan(&sig.StudentID, &sig.TaskKind, &sig.TaskID, &sig.ReviewRunID, &sig.ConceptID, &sig.AttemptKey,
			&correct, &dontKnow, &mastered, &legacy, &sig.ElapsedSeconds, &sig.Difficulty,
			&sig.XPEarned, &sig.AnswerSignature, &created); err != nil {
			return nil, fmt.Errorf("scan effort signal: %w", err)
		}
		sig.Correct = correct != 0
		sig.DontKnow = dontKnow != 0
		sig.Mastered = mastered != 0
		sig.LegacyPenalty = legacy != 0
		sig.CreatedAt = time.Unix(0, int64(created)).UTC()
		out = append(out, sig)
	}
	return out, rows.Err()
}

func countXPDeductions(db *sql.DB, rebind func(string) string, studentID string, since time.Time) (int, error) {
	var n int
	if err := db.QueryRow(rebind(`
		SELECT COUNT(*) FROM xp_events
		WHERE student_id = ? AND amount < 0
		  AND reason IN ('accountability:insufficient_effort',
		                 'accountability:repeated_guessing',
		                 'accountability:farming')
		  AND created_at >= ?
	`), studentID, signalTime(since)).Scan(&n); err != nil {
		return 0, fmt.Errorf("count xp deductions: %w", err)
	}
	return n, nil
}

func sumAccountabilityDeductionsSince(db *sql.DB, rebind func(string) string, studentID string, since time.Time) (int, error) {
	var total int
	sinceN := int64(0)
	if !since.IsZero() {
		sinceN = signalTime(since)
	}
	if err := db.QueryRow(rebind(`
		SELECT COALESCE(SUM(-amount), 0) FROM xp_events
		WHERE student_id = ? AND amount < 0
		  AND created_at >= ?
	`), studentID, sinceN).Scan(&total); err != nil {
		return 0, fmt.Errorf("sum xp deductions: %w", err)
	}
	return total, nil
}

func sumXPEarnedSince(db *sql.DB, rebind func(string) string, studentID string, since time.Time) (int, error) {
	var total int
	sinceN := int64(0)
	if !since.IsZero() {
		sinceN = signalTime(since)
	}
	if err := db.QueryRow(rebind(`
		SELECT COALESCE(SUM(amount), 0) FROM xp_events
		WHERE student_id = ? AND amount > 0
		  AND created_at >= ?
	`), studentID, sinceN).Scan(&total); err != nil {
		return 0, fmt.Errorf("sum xp earned: %w", err)
	}
	return total, nil
}

// applyXPAdjustmentTx moves a signed amount through the ledger and the counters
// in one transaction. Negative amounts are applied in full; lifetime XP may
// become negative. Daily XP remains capped at zero.
func applyXPAdjustmentTx(db *sql.DB, rebind func(string) string, forUpdate string, studentID string, amount int, reason, eventKey string) (int, error) {
	if studentID == "" || eventKey == "" {
		return 0, fmt.Errorf("xp adjustment: missing scope")
	}
	if amount == 0 {
		return 0, nil
	}
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("xp adjustment: %w", err)
	}
	defer tx.Rollback()

	applied, err := applyXPAdjustmentInTx(tx, rebind, forUpdate, studentID, amount, reason, eventKey)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("xp adjustment: %w", err)
	}
	return applied, nil
}

// applyXPAdjustmentInTx is the shared body used by a standalone adjustment and
// by the atomic study submission transaction. It never commits.
func applyXPAdjustmentInTx(tx *sql.Tx, rebind func(string) string, forUpdate string, studentID string, amount int, reason, eventKey string) (int, error) {
	if studentID == "" || eventKey == "" {
		return 0, fmt.Errorf("xp adjustment: missing scope")
	}
	if amount == 0 {
		return 0, nil
	}
	var exists int
	if err := tx.QueryRow(rebind(`SELECT 1 FROM students WHERE id = ?`+forUpdate), studentID).
		Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("xp adjustment: lock balance row: %w", err)
	}
	applied := amount
	now := time.Now().UTC().UnixNano()
	res, err := tx.Exec(rebind(`
		INSERT INTO xp_events (student_id, amount, reason, event_key, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(student_id, event_key) DO NOTHING
	`), studentID, applied, reason, eventKey, now)
	if err != nil {
		return 0, fmt.Errorf("xp adjustment: insert event: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, nil
	}
	todayStr := time.Now().UTC().Format("2006-01-02")
	if _, err := tx.Exec(rebind(`
		UPDATE students
		SET xp_total = xp_total + ?,
		    xp_today = CASE
		        WHEN xp_date = ? THEN (CASE WHEN xp_today + ? < 0 THEN 0 ELSE xp_today + ? END)
		        ELSE (CASE WHEN ? < 0 THEN 0 ELSE ? END)
		    END,
		    xp_date = ?
		WHERE id = ?
	`), applied, todayStr, applied, applied, applied, applied, todayStr, studentID); err != nil {
		return 0, fmt.Errorf("xp adjustment: update counters: %w", err)
	}
	return applied, nil
}

func signalTime(t time.Time) int64 {
	if t.IsZero() {
		return time.Now().UTC().UnixNano()
	}
	return t.UTC().UnixNano()
}
