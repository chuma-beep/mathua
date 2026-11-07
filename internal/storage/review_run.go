package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrReviewRunNotFound means the session has no durable review-run row.
var ErrReviewRunNotFound = errors.New("review run not found")

// ReviewRun is the durable completion record for one review session. The
// session ID is the run ID, so the server, the effort signals and the XP
// ledger can all name the same task.
type ReviewRun struct {
	ID           string
	StudentID    string
	StartedAt    time.Time
	CompletedAt  time.Time
	FinalizedAt  time.Time
	Status       string
	XPAdjustment int
	XPNote       string
	XPReasonCode string
	XPSource     string
}

// ReviewFinalization is the one signed adjustment a completed review run may
// receive. NoteForApplied runs inside the transaction with the amount that
// actually reached the ledger.
type ReviewFinalization struct {
	At             time.Time
	Amount         int
	Reason         string
	EventKey       string
	NoteForApplied func(applied int) (note, reasonCode, source string)
}

// ReviewRunStore is the optional persistence surface for review-run
// accountability. Stores that do not implement it keep review effort-only.
type ReviewRunStore interface {
	EnsureReviewRun(id, studentID string, at time.Time) error
	ReviewRunSignals(studentID, id string) ([]EffortSignal, error)
	CountReviewRunEffortDeductions(studentID, id string) (int, error)
	// FinalizeReviewRun atomically marks the run complete and applies one
	// signed adjustment. A concurrent or repeated finalizer returns the stored
	// outcome with finalizedNow false.
	FinalizeReviewRun(studentID, id string, fin ReviewFinalization) (ReviewRun, bool, error)
}

func ensureReviewRun(db *sql.DB, rebind func(string) string, id, studentID string, at time.Time) error {
	if id == "" || studentID == "" {
		return fmt.Errorf("review run: missing scope")
	}
	_, err := db.Exec(rebind(`
		INSERT INTO review_runs (id, student_id, started_at, status)
		VALUES (?, ?, ?, 'active')
		ON CONFLICT(id) DO NOTHING
	`), id, studentID, signalTime(at))
	if err != nil {
		return fmt.Errorf("ensure review run: %w", err)
	}
	return nil
}

func reviewRunSignals(db *sql.DB, rebind func(string) string, studentID, id string) ([]EffortSignal, error) {
	rows, err := db.Query(rebind(`
		SELECT `+effortSignalColumns+`
		FROM effort_signals
		WHERE student_id = ? AND review_run_id = ?
		ORDER BY created_at ASC, id ASC
	`), studentID, id)
	if err != nil {
		return nil, fmt.Errorf("review run signals: %w", err)
	}
	defer rows.Close()

	var out []EffortSignal
	for rows.Next() {
		var sig EffortSignal
		var correct, dontKnow, mastered, legacy, created int
		if err := rows.Scan(&sig.StudentID, &sig.TaskKind, &sig.TaskID, &sig.ReviewRunID, &sig.ConceptID, &sig.AttemptKey,
			&correct, &dontKnow, &mastered, &legacy, &sig.ElapsedSeconds, &sig.Difficulty,
			&sig.XPEarned, &sig.AnswerSignature, &created); err != nil {
			return nil, fmt.Errorf("scan review signal: %w", err)
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

func countReviewRunEffortDeductions(db *sql.DB, rebind func(string) string, studentID, id string) (int, error) {
	var n int
	if err := db.QueryRow(rebind(`
		SELECT COUNT(*) FROM xp_events
		WHERE student_id = ? AND amount < 0
		  AND event_key LIKE ?
	`), studentID, "accountability:review:"+id+":%").Scan(&n); err != nil {
		return 0, fmt.Errorf("count review effort deductions: %w", err)
	}
	return n, nil
}

func finalizeReviewRun(db *sql.DB, rebind func(string) string, forUpdate string, studentID, id string, fin ReviewFinalization) (ReviewRun, bool, error) {
	if studentID == "" || id == "" {
		return ReviewRun{}, false, fmt.Errorf("review run: missing scope")
	}
	tx, err := db.Begin()
	if err != nil {
		return ReviewRun{}, false, fmt.Errorf("finalize review run: %w", err)
	}
	defer tx.Rollback()

	var run ReviewRun
	var started, completed, finalized int64
	if err := tx.QueryRow(rebind(`
		SELECT id, student_id, started_at, completed_at, finalized_at, status,
		       xp_adjustment, xp_note, xp_reason_code, xp_source
		FROM review_runs WHERE id = ? AND student_id = ?`+forUpdate), id, studentID).Scan(
		&run.ID, &run.StudentID, &started, &completed, &finalized, &run.Status,
		&run.XPAdjustment, &run.XPNote, &run.XPReasonCode, &run.XPSource,
	); err != nil {
		if err == sql.ErrNoRows {
			return ReviewRun{}, false, ErrReviewRunNotFound
		}
		return ReviewRun{}, false, fmt.Errorf("read review run: %w", err)
	}
	run.StartedAt = unixOrZero(started)
	run.CompletedAt = unixOrZero(completed)
	run.FinalizedAt = unixOrZero(finalized)
	if finalized > 0 {
		if err := tx.Commit(); err != nil {
			return ReviewRun{}, false, fmt.Errorf("finalize review run: commit: %w", err)
		}
		return run, false, nil
	}

	applied := 0
	if fin.Amount != 0 {
		applied, err = applyXPAdjustmentInTx(tx, rebind, forUpdate, studentID, fin.Amount, fin.Reason, fin.EventKey)
		if err != nil {
			return ReviewRun{}, false, err
		}
	}
	note, reasonCode, source := "", fin.Reason, ""
	if fin.NoteForApplied != nil {
		note, reasonCode, source = fin.NoteForApplied(applied)
	}
	completedAt := fin.At
	if completedAt.IsZero() {
		completedAt = time.Now().UTC()
	}
	finalizedAt := time.Now().UTC()
	if _, err := tx.Exec(rebind(`
		UPDATE review_runs
		SET completed_at = ?, finalized_at = ?, status = 'completed',
		    xp_adjustment = ?, xp_note = ?, xp_reason_code = ?, xp_source = ?
		WHERE id = ? AND student_id = ?
	`), signalTime(completedAt), signalTime(finalizedAt), applied, note, reasonCode, source, id, studentID); err != nil {
		return ReviewRun{}, false, fmt.Errorf("update review run: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return ReviewRun{}, false, fmt.Errorf("finalize review run: commit: %w", err)
	}
	run.CompletedAt = completedAt
	run.FinalizedAt = finalizedAt
	run.Status = "completed"
	run.XPAdjustment = applied
	run.XPNote = note
	run.XPReasonCode = reasonCode
	run.XPSource = source
	return run, true, nil
}

func unixOrZero(n int64) time.Time {
	if n == 0 {
		return time.Time{}
	}
	return time.Unix(0, n).UTC()
}

func (s *SQLiteStore) EnsureReviewRun(id, studentID string, at time.Time) error {
	return ensureReviewRun(s.db, func(q string) string { return q }, id, studentID, at)
}

func (s *PostgresStore) EnsureReviewRun(id, studentID string, at time.Time) error {
	return ensureReviewRun(s.db, rebindDollars, id, studentID, at)
}

func (s *SQLiteStore) ReviewRunSignals(studentID, id string) ([]EffortSignal, error) {
	return reviewRunSignals(s.db, func(q string) string { return q }, studentID, id)
}

func (s *PostgresStore) ReviewRunSignals(studentID, id string) ([]EffortSignal, error) {
	return reviewRunSignals(s.db, rebindDollars, studentID, id)
}

func (s *SQLiteStore) CountReviewRunEffortDeductions(studentID, id string) (int, error) {
	return countReviewRunEffortDeductions(s.db, func(q string) string { return q }, studentID, id)
}

func (s *PostgresStore) CountReviewRunEffortDeductions(studentID, id string) (int, error) {
	return countReviewRunEffortDeductions(s.db, rebindDollars, studentID, id)
}

func (s *SQLiteStore) FinalizeReviewRun(studentID, id string, fin ReviewFinalization) (ReviewRun, bool, error) {
	return finalizeReviewRun(s.db, func(q string) string { return q }, "", studentID, id, fin)
}

func (s *PostgresStore) FinalizeReviewRun(studentID, id string, fin ReviewFinalization) (ReviewRun, bool, error) {
	return finalizeReviewRun(s.db, rebindDollars, " FOR UPDATE", studentID, id, fin)
}
