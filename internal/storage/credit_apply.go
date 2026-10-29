package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/fsrs"
)

// Active FIRe application. Unlike the shadow projection in credit_*.go this
// path writes concept_progress, so it verifies the mode-independent target
// shape, the ordinal's due baseline and the progress row in one serialised
// transaction. It never changes mastery status, evidence, streaks, placement or
// XP: an implicit review is a scheduling event.
var _ CreditApplyStore = (*SQLiteStore)(nil)
var _ CreditOrdinalStore = (*SQLiteStore)(nil)
var _ CreditPreviewStore = (*SQLiteStore)(nil)

type creditTargetRow struct {
	Status          string
	Due             sql.NullString
	LastReviewed    sql.NullString
	SM2Interval     int
	PlacementSeeded int
	IsProvisional   int
	Stability       sql.NullFloat64
	Difficulty      sql.NullFloat64
	SRSLastReview   sql.NullString
}

// creditActivationSQLite reads the durable activation inside the application
// transaction. A missing row, a non-active row or a cohort mismatch fails
// closed before any progress or ordinal write.
func creditActivationSQLite(ctx context.Context, conn *sql.Conn, cohort string) string {
	cohort = strings.TrimSpace(cohort)
	if cohort == "" {
		return "not_active"
	}
	var mode, activeCohort string
	err := conn.QueryRowContext(ctx, `
		SELECT mode, cohort FROM fire_activations WHERE id = 'global'
	`).Scan(&mode, &activeCohort)
	if errors.Is(err, sql.ErrNoRows) {
		return "not_active"
	}
	if err != nil {
		return "activation_unavailable"
	}
	if strings.ToLower(strings.TrimSpace(mode)) != "active" {
		return "not_active"
	}
	if strings.TrimSpace(activeCohort) != cohort {
		return "cohort_mismatch"
	}
	return ""
}

func parseCreditTime(null sql.NullString) (time.Time, bool) {
	if !null.Valid {
		return time.Time{}, false
	}
	return parseCreditDueText(null.String)
}

func parseCreditDueText(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

// creditBaselineMatches compares a stored ordinal baseline with the current
// due date by time, not by string shape. Older rows and imported databases may
// store the same instant as "2006-01-02 15:04:05" while the current value is
// RFC3339. Treating that formatting difference as a moved due date would stall
// every application for the target.
func creditBaselineMatches(baseline, actual string) bool {
	if strings.TrimSpace(baseline) == "" {
		return true
	}
	from, okFrom := parseCreditDueText(baseline)
	to, okTo := parseCreditDueText(actual)
	if okFrom && okTo {
		return from.Equal(to)
	}
	return baseline == actual
}

func creditReasonEligible(row creditTargetRow) string {
	switch strings.ToUpper(strings.TrimSpace(row.Status)) {
	case "MASTERED", "DECAYING":
	default:
		return "not_mastered"
	}
	if row.PlacementSeeded != 0 {
		return "placement_seeded"
	}
	if row.IsProvisional != 0 {
		return "provisional"
	}
	if _, ok := parseCreditTime(row.Due); !ok {
		return "no_due_date"
	}
	return ""
}

func creditNextState(row creditTargetRow, due, now time.Time) (stability, difficulty float64, next time.Time, err error) {
	p := fsrs.DefaultParameters()
	stability = 0
	if row.Stability.Valid && !math.IsNaN(row.Stability.Float64) && row.Stability.Float64 > 0 {
		stability = row.Stability.Float64
	}
	if stability <= 0 {
		if seeded, seedErr := fsrs.SeedFromInterval(row.SM2Interval, due, now); seedErr == nil {
			stability = seeded.Stability
		} else {
			stability = p.InitialStability(fsrs.Good)
		}
	}
	if math.IsNaN(stability) || math.IsInf(stability, 0) || stability <= 0 {
		return 0, 0, time.Time{}, fmt.Errorf("credit: target has no usable stability")
	}
	difficulty = p.InitialDifficulty(fsrs.Good)
	if row.Difficulty.Valid && !math.IsNaN(row.Difficulty.Float64) && row.Difficulty.Float64 > 0 {
		difficulty = row.Difficulty.Float64
	}
	last := due
	if t, ok := parseCreditTime(row.SRSLastReview); ok {
		last = t
	} else if t, ok := parseCreditTime(row.LastReviewed); ok {
		last = t
	}
	elapsed := now.Sub(last).Hours() / 24
	if elapsed < 0 || math.IsNaN(elapsed) || math.IsInf(elapsed, 0) {
		elapsed = 0
	}
	nextStability := p.NextStability(stability, difficulty, p.Retrievability(elapsed, stability), fsrs.Good)
	nextDifficulty := p.NextDifficulty(difficulty, fsrs.Good)
	step := p.Interval(nextStability)
	if step < 1 {
		step = 1
	}
	base := due
	if now.After(base) {
		base = now
	}
	if t, ok := parseCreditTime(row.LastReviewed); ok && t.After(base) {
		base = t
	}
	next = base.AddDate(0, 0, step)
	if next.After(now.AddDate(0, 0, fsrs.MaxIntervalDays)) {
		next = now.AddDate(0, 0, fsrs.MaxIntervalDays)
	}
	return nextStability, nextDifficulty, next, nil
}

func (s *SQLiteStore) ApplyCredit(req CreditApplyRequest) (CreditApplication, error) {
	out := CreditApplication{
		StudentID:       strings.TrimSpace(req.StudentID),
		TargetConceptID: strings.TrimSpace(req.TargetConceptID),
		SourceConceptID: strings.TrimSpace(req.SourceConceptID),
		SourceAttemptID: strings.TrimSpace(req.SourceAttemptID),
		Epoch:           req.Epoch,
	}
	if out.StudentID == "" || out.TargetConceptID == "" || out.SourceAttemptID == "" {
		return out, fmt.Errorf("credit: apply needs student, target and source attempt")
	}
	now := req.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	err := s.runImmediate(context.Background(), func(conn *sql.Conn) error {
		ctx := context.Background()
		var epoch int64
		if err := conn.QueryRowContext(ctx, `
			SELECT credit_epoch FROM implicit_credit_targets
			WHERE student_id = ? AND target_concept_id = ?
		`, out.StudentID, out.TargetConceptID).Scan(&epoch); errors.Is(err, sql.ErrNoRows) {
			out.Reason = "no_credit"
			return nil
		} else if err != nil {
			return fmt.Errorf("credit: read epoch: %w", err)
		}
		if req.Epoch > 0 && epoch != req.Epoch {
			out.Epoch = epoch
			out.Reason = "epoch_closed"
			return nil
		}
		out.Epoch = epoch
		if existing, ok, err := readCreditApplicationSQLite(ctx, conn, out.StudentID, out.TargetConceptID, out.SourceAttemptID, epoch); err != nil {
			return err
		} else if ok {
			existing.Replayed = true
			out = existing
			return nil
		}
		if reason := creditActivationSQLite(ctx, conn, req.Cohort); reason != "" {
			out.Reason = reason
			return nil
		}
		row, ok, err := readCreditTargetSQLite(ctx, conn, out.StudentID, out.TargetConceptID)
		if err != nil {
			return err
		}
		if !ok {
			out.Reason = "no_progress"
			return nil
		}
		if out.SourceConceptID != "" && out.SourceConceptID == out.TargetConceptID {
			out.Reason = "source_is_target"
			return nil
		}
		if reason := creditReasonEligible(row); reason != "" {
			out.Reason = reason
			return nil
		}
		actualDue, _ := parseCreditTime(row.Due)
		if actualDue.After(now) {
			out.Reason = "not_due"
			return nil
		}
		ordinal, baseline, ok, err := pickCreditOrdinalSQLite(ctx, conn, out.StudentID, out.TargetConceptID, epoch, req.Ordinal)
		if err != nil {
			return err
		}
		if !ok {
			out.Reason = "no_banked_ordinal"
			return nil
		}
		out.Ordinal = ordinal
		out.ActualDue = actualDue.Format(time.RFC3339)
		if !creditBaselineMatches(baseline, out.ActualDue) {
			out.Reason = "stale_baseline"
			return nil
		}
		stability, difficulty, next, err := creditNextState(row, actualDue, now)
		if err != nil {
			out.Reason = "invalid_fsrs_state"
			return nil
		}
		out.ProposedDue = next.Format(time.RFC3339)
		out.StabilityBefore = creditFloat(row.Stability)
		out.DifficultyBefore = creditFloat(row.Difficulty)
		out.StabilityAfter = stability
		out.DifficultyAfter = difficulty
		out.ReviewsDeferred = 1
		stabilityArg := stability
		difficultyArg := difficulty
		var dueArg any
		if row.Due.Valid {
			dueArg = row.Due.String
		}
		res, err := conn.ExecContext(ctx, `
			UPDATE concept_progress
			SET next_review_due = ?, fsrs_stability = ?, fsrs_difficulty = ?,
			    fsrs_last_review = ?, last_reviewed = ?
			WHERE student_id = ? AND concept_id = ? AND next_review_due IS ?
		`, next.Format(time.RFC3339), stabilityArg, difficultyArg, now.Format(time.RFC3339), now.Format(time.RFC3339), out.StudentID, out.TargetConceptID, dueArg)
		if err != nil {
			return fmt.Errorf("credit: update progress: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			out.Reason = "due_moved"
			return nil
		}
		res, err = conn.ExecContext(ctx, `
			UPDATE implicit_credit_ordinals
			SET state = 'projected', projected_at = ?, baseline_due = ?, counterfactual_due = ?, attempts = attempts + 1
			WHERE student_id = ? AND target_concept_id = ? AND credit_epoch = ? AND ordinal = ? AND state = 'banked'
		`, now.UnixNano(), out.ActualDue, out.ProposedDue, out.StudentID, out.TargetConceptID, epoch, ordinal)
		if err != nil {
			return fmt.Errorf("credit: consume ordinal: %w", err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("%w: ordinal %d changed", ErrCreditOrdinalStateChanged, ordinal)
		}
		_, err = conn.ExecContext(ctx, `
			INSERT INTO implicit_credit_applications
				(student_id, target_concept_id, credit_epoch, ordinal, source_concept_id, source_attempt_id,
				 applied_at, actual_due, proposed_due, stability_before, stability_after,
				 difficulty_before, difficulty_after, reviews_deferred)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT DO NOTHING
		`, out.StudentID, out.TargetConceptID, epoch, ordinal, out.SourceConceptID, out.SourceAttemptID,
			now.UnixNano(), out.ActualDue, out.ProposedDue, out.StabilityBefore, out.StabilityAfter,
			out.DifficultyBefore, out.DifficultyAfter, out.ReviewsDeferred)
		if err != nil {
			return fmt.Errorf("credit: record application: %w", err)
		}
		out.Applied = true
		return nil
	})
	return out, err
}

// PreviewCredit mirrors ApplyCredit's decision but never writes. It is the only
// method the dry-run command is allowed to call.
func (s *SQLiteStore) PreviewCredit(req CreditApplyRequest) (CreditApplication, error) {
	out := CreditApplication{
		StudentID:       strings.TrimSpace(req.StudentID),
		TargetConceptID: strings.TrimSpace(req.TargetConceptID),
		SourceConceptID: strings.TrimSpace(req.SourceConceptID),
		SourceAttemptID: strings.TrimSpace(req.SourceAttemptID),
		Epoch:           req.Epoch,
	}
	if out.StudentID == "" || out.TargetConceptID == "" {
		return out, fmt.Errorf("credit: preview needs student and target")
	}
	now := req.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	ctx := context.Background()
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return out, fmt.Errorf("credit: preview connection: %w", err)
	}
	defer conn.Close()
	var epoch int64
	if err := conn.QueryRowContext(ctx, `
		SELECT credit_epoch FROM implicit_credit_targets
		WHERE student_id = ? AND target_concept_id = ?
	`, out.StudentID, out.TargetConceptID).Scan(&epoch); errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_credit"
		return out, nil
	} else if err != nil {
		return out, fmt.Errorf("credit: preview epoch: %w", err)
	}
	if req.Epoch > 0 && epoch != req.Epoch {
		out.Epoch = epoch
		out.Reason = "epoch_closed"
		return out, nil
	}
	out.Epoch = epoch
	row, ok, err := readCreditTargetSQLite(ctx, conn, out.StudentID, out.TargetConceptID)
	if err != nil {
		return out, err
	}
	if !ok {
		out.Reason = "no_progress"
		return out, nil
	}
	if out.SourceConceptID != "" && out.SourceConceptID == out.TargetConceptID {
		out.Reason = "source_is_target"
		return out, nil
	}
	if reason := creditReasonEligible(row); reason != "" {
		out.Reason = reason
		return out, nil
	}
	actualDue, _ := parseCreditTime(row.Due)
	if actualDue.After(now) {
		out.Reason = "not_due"
		return out, nil
	}
	ordinal, baseline, ok, err := pickCreditOrdinalSQLite(ctx, conn, out.StudentID, out.TargetConceptID, epoch, req.Ordinal)
	if err != nil {
		return out, err
	}
	if !ok {
		out.Reason = "no_banked_ordinal"
		return out, nil
	}
	out.Ordinal = ordinal
	out.ActualDue = actualDue.Format(time.RFC3339)
	if !creditBaselineMatches(baseline, out.ActualDue) {
		out.Reason = "stale_baseline"
		return out, nil
	}
	stability, difficulty, next, err := creditNextState(row, actualDue, now)
	if err != nil {
		out.Reason = "invalid_fsrs_state"
		return out, nil
	}
	out.ProposedDue = next.Format(time.RFC3339)
	out.StabilityBefore = creditFloat(row.Stability)
	out.DifficultyBefore = creditFloat(row.Difficulty)
	out.StabilityAfter = stability
	out.DifficultyAfter = difficulty
	out.ReviewsDeferred = 1
	out.Applied = true
	return out, nil
}

func readCreditTargetSQLite(ctx context.Context, conn *sql.Conn, studentID, conceptID string) (creditTargetRow, bool, error) {
	var row creditTargetRow
	err := conn.QueryRowContext(ctx, `
		SELECT status, next_review_due, last_reviewed, sm2_interval,
		       placement_seeded, is_provisional, fsrs_stability, fsrs_difficulty, fsrs_last_review
		FROM concept_progress WHERE student_id = ? AND concept_id = ?
	`, studentID, conceptID).Scan(&row.Status, &row.Due, &row.LastReviewed, &row.SM2Interval,
		&row.PlacementSeeded, &row.IsProvisional, &row.Stability, &row.Difficulty, &row.SRSLastReview)
	if errors.Is(err, sql.ErrNoRows) {
		return row, false, nil
	}
	if err != nil {
		return row, false, fmt.Errorf("credit: read progress: %w", err)
	}
	return row, true, nil
}

func pickCreditOrdinalSQLite(ctx context.Context, conn *sql.Conn, studentID, targetID string, epoch int64, requested int) (int, string, bool, error) {
	query := `
		SELECT ordinal, COALESCE(baseline_due, '') FROM implicit_credit_ordinals
		WHERE student_id = ? AND target_concept_id = ? AND credit_epoch = ? AND state = 'banked'`
	args := []any{studentID, targetID, epoch}
	if requested > 0 {
		query += " AND ordinal = ?"
		args = append(args, requested)
	}
	query += " ORDER BY ordinal LIMIT 1"
	var ordinal int
	var baseline string
	err := conn.QueryRowContext(ctx, query, args...).Scan(&ordinal, &baseline)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", false, nil
	}
	if err != nil {
		return 0, "", false, fmt.Errorf("credit: pick ordinal: %w", err)
	}
	return ordinal, baseline, true, nil
}

func readCreditApplicationSQLite(ctx context.Context, conn *sql.Conn, studentID, targetID, attemptID string, epoch int64) (CreditApplication, bool, error) {
	var out CreditApplication
	err := conn.QueryRowContext(ctx, `
		SELECT student_id, target_concept_id, credit_epoch, ordinal, source_concept_id, source_attempt_id,
		       actual_due, proposed_due, COALESCE(stability_before, 0), COALESCE(stability_after, 0),
		       COALESCE(difficulty_before, 0), COALESCE(difficulty_after, 0), reviews_deferred
		FROM implicit_credit_applications
		WHERE student_id = ? AND source_attempt_id = ? AND target_concept_id = ? AND credit_epoch = ?
	`, studentID, attemptID, targetID, epoch).Scan(&out.StudentID, &out.TargetConceptID, &out.Epoch,
		&out.Ordinal, &out.SourceConceptID, &out.SourceAttemptID, &out.ActualDue, &out.ProposedDue,
		&out.StabilityBefore, &out.StabilityAfter, &out.DifficultyBefore, &out.DifficultyAfter, &out.ReviewsDeferred)
	if errors.Is(err, sql.ErrNoRows) {
		return out, false, nil
	}
	if err != nil {
		return out, false, fmt.Errorf("credit: read application: %w", err)
	}
	out.Applied = true
	return out, true, nil
}

func creditFloat(n sql.NullFloat64) float64 {
	if n.Valid && !math.IsNaN(n.Float64) {
		return n.Float64
	}
	return 0
}

func (s *SQLiteStore) BankedOrdinalDetails(studentID string) ([]CreditOrdinalDetail, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT o.student_id, o.target_concept_id, o.credit_epoch, o.ordinal, o.state,
		       COALESCE(o.baseline_due, ''), COALESCE(o.counterfactual_due, ''), o.attempts,
		       COALESCE(p.status, ''), COALESCE(p.next_review_due, ''),
		       COALESCE(p.placement_seeded, 0), COALESCE(p.is_provisional, 0),
		       CASE WHEN p.student_id IS NULL THEN 0 ELSE 1 END
		FROM implicit_credit_ordinals o
		LEFT JOIN concept_progress p
		  ON p.student_id = o.student_id AND p.concept_id = o.target_concept_id
		WHERE o.student_id = ? AND o.state = 'banked'
		ORDER BY o.target_concept_id, o.credit_epoch, o.ordinal
	`, studentID)
	if err != nil {
		return nil, fmt.Errorf("credit: banked details: %w", err)
	}
	defer rows.Close()
	return scanCreditOrdinalDetails(rows)
}

// AllBankedOrdinalDetails returns every banked ordinal. The dry-run command is
// the only caller; it is read-only and intentionally has no student filter.
func (s *SQLiteStore) AllBankedOrdinalDetails() ([]CreditOrdinalDetail, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT o.student_id, o.target_concept_id, o.credit_epoch, o.ordinal, o.state,
		       COALESCE(o.baseline_due, ''), COALESCE(o.counterfactual_due, ''), o.attempts,
		       COALESCE(p.status, ''), COALESCE(p.next_review_due, ''),
		       COALESCE(p.placement_seeded, 0), COALESCE(p.is_provisional, 0),
		       CASE WHEN p.student_id IS NULL THEN 0 ELSE 1 END
		FROM implicit_credit_ordinals o
		LEFT JOIN concept_progress p
		  ON p.student_id = o.student_id AND p.concept_id = o.target_concept_id
		WHERE o.state = 'banked'
		ORDER BY o.student_id, o.target_concept_id, o.credit_epoch, o.ordinal
	`)
	if err != nil {
		return nil, fmt.Errorf("credit: all banked details: %w", err)
	}
	defer rows.Close()
	return scanCreditOrdinalDetails(rows)
}

func scanCreditOrdinalDetails(rows *sql.Rows) ([]CreditOrdinalDetail, error) {
	var out []CreditOrdinalDetail
	for rows.Next() {
		var d CreditOrdinalDetail
		var hasProgress int
		if err := rows.Scan(&d.StudentID, &d.ConceptID, &d.Epoch, &d.Ordinal, &d.State,
			&d.BaselineDue, &d.CounterfactualDue, &d.Attempts, &d.TargetStatus, &d.TargetDue,
			&d.PlacementSeeded, &d.IsProvisional, &hasProgress); err != nil {
			return nil, err
		}
		d.HasProgress = hasProgress != 0
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CreditApplications(studentID, targetConceptID string, limit int) ([]CreditApplication, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	ctx := context.Background()
	query := `
		SELECT student_id, target_concept_id, credit_epoch, ordinal, source_concept_id, source_attempt_id,
		       actual_due, proposed_due, COALESCE(stability_before, 0), COALESCE(stability_after, 0),
		       COALESCE(difficulty_before, 0), COALESCE(difficulty_after, 0), reviews_deferred
		FROM implicit_credit_applications WHERE student_id = ?`
	args := []any{studentID}
	if targetConceptID != "" {
		query += " AND target_concept_id = ?"
		args = append(args, targetConceptID)
	}
	query += " ORDER BY applied_at DESC LIMIT ?"
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("credit: list applications: %w", err)
	}
	defer rows.Close()
	var out []CreditApplication
	for rows.Next() {
		var a CreditApplication
		if err := rows.Scan(&a.StudentID, &a.TargetConceptID, &a.Epoch, &a.Ordinal, &a.SourceConceptID,
			&a.SourceAttemptID, &a.ActualDue, &a.ProposedDue, &a.StabilityBefore, &a.StabilityAfter,
			&a.DifficultyBefore, &a.DifficultyAfter, &a.ReviewsDeferred); err != nil {
			return nil, err
		}
		a.Applied = true
		out = append(out, a)
	}
	return out, rows.Err()
}

// PostgreSQL

var _ CreditApplyStore = (*PostgresStore)(nil)
var _ CreditOrdinalStore = (*PostgresStore)(nil)
var _ CreditPreviewStore = (*PostgresStore)(nil)

func (s *PostgresStore) ApplyCredit(req CreditApplyRequest) (CreditApplication, error) {
	out := CreditApplication{
		StudentID:       strings.TrimSpace(req.StudentID),
		TargetConceptID: strings.TrimSpace(req.TargetConceptID),
		SourceConceptID: strings.TrimSpace(req.SourceConceptID),
		SourceAttemptID: strings.TrimSpace(req.SourceAttemptID),
		Epoch:           req.Epoch,
	}
	if out.StudentID == "" || out.TargetConceptID == "" || out.SourceAttemptID == "" {
		return out, fmt.Errorf("credit: apply needs student, target and source attempt")
	}
	now := req.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return out, err
	}
	defer tx.Rollback()
	var epoch int64
	err = tx.QueryRowContext(ctx, `
		SELECT credit_epoch FROM implicit_credit_targets
		WHERE student_id = $1 AND target_concept_id = $2 FOR UPDATE
	`, out.StudentID, out.TargetConceptID).Scan(&epoch)
	if errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_credit"
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("credit: lock target: %w", err)
	}
	if req.Epoch > 0 && epoch != req.Epoch {
		out.Epoch = epoch
		out.Reason = "epoch_closed"
		return out, nil
	}
	out.Epoch = epoch
	var row creditTargetRow
	err = tx.QueryRowContext(ctx, `
		SELECT status, next_review_due, last_reviewed, sm2_interval,
		       placement_seeded, is_provisional, fsrs_stability, fsrs_difficulty, fsrs_last_review
		FROM concept_progress WHERE student_id = $1 AND concept_id = $2 FOR UPDATE
	`, out.StudentID, out.TargetConceptID).Scan(&row.Status, &row.Due, &row.LastReviewed, &row.SM2Interval,
		&row.PlacementSeeded, &row.IsProvisional, &row.Stability, &row.Difficulty, &row.SRSLastReview)
	if errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_progress"
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("credit: lock progress: %w", err)
	}
	var existing CreditApplication
	err = tx.QueryRowContext(ctx, `
		SELECT student_id, target_concept_id, credit_epoch, ordinal, source_concept_id, source_attempt_id,
		       actual_due, proposed_due, COALESCE(stability_before, 0), COALESCE(stability_after, 0),
		       COALESCE(difficulty_before, 0), COALESCE(difficulty_after, 0), reviews_deferred
		FROM implicit_credit_applications
		WHERE student_id = $1 AND source_attempt_id = $2 AND target_concept_id = $3 AND credit_epoch = $4
	`, out.StudentID, out.SourceAttemptID, out.TargetConceptID, epoch).Scan(&existing.StudentID, &existing.TargetConceptID,
		&existing.Epoch, &existing.Ordinal, &existing.SourceConceptID, &existing.SourceAttemptID,
		&existing.ActualDue, &existing.ProposedDue, &existing.StabilityBefore, &existing.StabilityAfter,
		&existing.DifficultyBefore, &existing.DifficultyAfter, &existing.ReviewsDeferred)
	if err == nil {
		existing.Applied = true
		existing.Replayed = true
		return existing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return out, fmt.Errorf("credit: read application: %w", err)
	}
	if reason := creditActivationPostgres(ctx, tx, req.Cohort); reason != "" {
		out.Reason = reason
		return out, nil
	}
	if out.SourceConceptID != "" && out.SourceConceptID == out.TargetConceptID {
		out.Reason = "source_is_target"
		return out, nil
	}
	if reason := creditReasonEligible(row); reason != "" {
		out.Reason = reason
		return out, nil
	}
	actualDue, _ := parseCreditTime(row.Due)
	if actualDue.After(now) {
		out.Reason = "not_due"
		return out, nil
	}
	var ordinal int
	var baseline string
	query := `
		SELECT ordinal, COALESCE(baseline_due, '') FROM implicit_credit_ordinals
		WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3 AND state = 'banked'`
	args := []any{out.StudentID, out.TargetConceptID, epoch}
	if req.Ordinal > 0 {
		query += " AND ordinal = $4"
		args = append(args, req.Ordinal)
	}
	query += " ORDER BY ordinal LIMIT 1"
	err = tx.QueryRowContext(ctx, query, args...).Scan(&ordinal, &baseline)
	if errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_banked_ordinal"
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("credit: pick ordinal: %w", err)
	}
	out.Ordinal = ordinal
	out.ActualDue = actualDue.Format(time.RFC3339)
	if !creditBaselineMatches(baseline, out.ActualDue) {
		out.Reason = "stale_baseline"
		return out, nil
	}
	stability, difficulty, next, err := creditNextState(row, actualDue, now)
	if err != nil {
		out.Reason = "invalid_fsrs_state"
		return out, nil
	}
	out.ProposedDue = next.Format(time.RFC3339)
	out.StabilityBefore = creditFloat(row.Stability)
	out.DifficultyBefore = creditFloat(row.Difficulty)
	out.StabilityAfter = stability
	out.DifficultyAfter = difficulty
	out.ReviewsDeferred = 1
	var dueArg any
	if row.Due.Valid {
		dueArg = row.Due.String
	}
	res, err := tx.ExecContext(ctx, `
		UPDATE concept_progress
		SET next_review_due = $1, fsrs_stability = $2, fsrs_difficulty = $3,
		    fsrs_last_review = $4, last_reviewed = $5
		WHERE student_id = $6 AND concept_id = $7 AND next_review_due IS NOT DISTINCT FROM $8
	`, next.Format(time.RFC3339), stability, difficulty, now.Format(time.RFC3339), now.Format(time.RFC3339),
		out.StudentID, out.TargetConceptID, dueArg)
	if err != nil {
		return out, fmt.Errorf("credit: update progress: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		out.Reason = "due_moved"
		return out, nil
	}
	res, err = tx.ExecContext(ctx, `
		UPDATE implicit_credit_ordinals
		SET state = 'projected', projected_at = $1, baseline_due = $2, counterfactual_due = $3, attempts = attempts + 1
		WHERE student_id = $4 AND target_concept_id = $5 AND credit_epoch = $6 AND ordinal = $7 AND state = 'banked'
	`, now.UnixNano(), out.ActualDue, out.ProposedDue, out.StudentID, out.TargetConceptID, epoch, ordinal)
	if err != nil {
		return out, fmt.Errorf("credit: consume ordinal: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return out, fmt.Errorf("%w: ordinal %d changed", ErrCreditOrdinalStateChanged, ordinal)
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO implicit_credit_applications
			(student_id, target_concept_id, credit_epoch, ordinal, source_concept_id, source_attempt_id,
			 applied_at, actual_due, proposed_due, stability_before, stability_after,
			 difficulty_before, difficulty_after, reviews_deferred)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT DO NOTHING
	`, out.StudentID, out.TargetConceptID, epoch, ordinal, out.SourceConceptID, out.SourceAttemptID,
		now.UnixNano(), out.ActualDue, out.ProposedDue, out.StabilityBefore, out.StabilityAfter,
		out.DifficultyBefore, out.DifficultyAfter, out.ReviewsDeferred)
	if err != nil {
		return out, fmt.Errorf("credit: record application: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return out, fmt.Errorf("credit: commit application: %w", err)
	}
	out.Applied = true
	return out, nil
}

// PreviewCredit mirrors ApplyCredit's decision but never writes. It runs in a
// read-only transaction so a copied database or read replica is safe.
func (s *PostgresStore) PreviewCredit(req CreditApplyRequest) (CreditApplication, error) {
	out := CreditApplication{
		StudentID:       strings.TrimSpace(req.StudentID),
		TargetConceptID: strings.TrimSpace(req.TargetConceptID),
		SourceConceptID: strings.TrimSpace(req.SourceConceptID),
		SourceAttemptID: strings.TrimSpace(req.SourceAttemptID),
		Epoch:           req.Epoch,
	}
	if out.StudentID == "" || out.TargetConceptID == "" {
		return out, fmt.Errorf("credit: preview needs student and target")
	}
	now := req.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return out, err
	}
	defer tx.Rollback()
	var epoch int64
	if err := tx.QueryRowContext(ctx, `
		SELECT credit_epoch FROM implicit_credit_targets
		WHERE student_id = $1 AND target_concept_id = $2
	`, out.StudentID, out.TargetConceptID).Scan(&epoch); errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_credit"
		return out, nil
	} else if err != nil {
		return out, fmt.Errorf("credit: preview epoch: %w", err)
	}
	if req.Epoch > 0 && epoch != req.Epoch {
		out.Epoch = epoch
		out.Reason = "epoch_closed"
		return out, nil
	}
	out.Epoch = epoch
	var row creditTargetRow
	err = tx.QueryRowContext(ctx, `
		SELECT status, next_review_due, last_reviewed, sm2_interval,
		       placement_seeded, is_provisional, fsrs_stability, fsrs_difficulty, fsrs_last_review
		FROM concept_progress WHERE student_id = $1 AND concept_id = $2
	`, out.StudentID, out.TargetConceptID).Scan(&row.Status, &row.Due, &row.LastReviewed, &row.SM2Interval,
		&row.PlacementSeeded, &row.IsProvisional, &row.Stability, &row.Difficulty, &row.SRSLastReview)
	if errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_progress"
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("credit: preview progress: %w", err)
	}
	if out.SourceConceptID != "" && out.SourceConceptID == out.TargetConceptID {
		out.Reason = "source_is_target"
		return out, nil
	}
	if reason := creditReasonEligible(row); reason != "" {
		out.Reason = reason
		return out, nil
	}
	actualDue, _ := parseCreditTime(row.Due)
	if actualDue.After(now) {
		out.Reason = "not_due"
		return out, nil
	}
	query := `
		SELECT ordinal, COALESCE(baseline_due, '') FROM implicit_credit_ordinals
		WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3 AND state = 'banked'`
	args := []any{out.StudentID, out.TargetConceptID, epoch}
	if req.Ordinal > 0 {
		query += " AND ordinal = $4"
		args = append(args, req.Ordinal)
	}
	query += " ORDER BY ordinal LIMIT 1"
	var ordinal int
	var baseline string
	err = tx.QueryRowContext(ctx, query, args...).Scan(&ordinal, &baseline)
	if errors.Is(err, sql.ErrNoRows) {
		out.Reason = "no_banked_ordinal"
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("credit: preview ordinal: %w", err)
	}
	out.Ordinal = ordinal
	out.ActualDue = actualDue.Format(time.RFC3339)
	if !creditBaselineMatches(baseline, out.ActualDue) {
		out.Reason = "stale_baseline"
		return out, nil
	}
	stability, difficulty, next, err := creditNextState(row, actualDue, now)
	if err != nil {
		out.Reason = "invalid_fsrs_state"
		return out, nil
	}
	out.ProposedDue = next.Format(time.RFC3339)
	out.StabilityBefore = creditFloat(row.Stability)
	out.DifficultyBefore = creditFloat(row.Difficulty)
	out.StabilityAfter = stability
	out.DifficultyAfter = difficulty
	out.ReviewsDeferred = 1
	out.Applied = true
	return out, nil
}

// creditActivationPostgres mirrors creditActivationSQLite inside the write
// transaction, so a concurrent deactivation cannot race an application.
func creditActivationPostgres(ctx context.Context, tx *sql.Tx, cohort string) string {
	cohort = strings.TrimSpace(cohort)
	if cohort == "" {
		return "not_active"
	}
	var mode, activeCohort string
	err := tx.QueryRowContext(ctx, `
		SELECT mode, cohort FROM fire_activations WHERE id = 'global'
	`).Scan(&mode, &activeCohort)
	if errors.Is(err, sql.ErrNoRows) {
		return "not_active"
	}
	if err != nil {
		return "activation_unavailable"
	}
	if strings.ToLower(strings.TrimSpace(mode)) != "active" {
		return "not_active"
	}
	if strings.TrimSpace(activeCohort) != cohort {
		return "cohort_mismatch"
	}
	return ""
}

func (s *PostgresStore) BankedOrdinalDetails(studentID string) ([]CreditOrdinalDetail, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT o.student_id, o.target_concept_id, o.credit_epoch, o.ordinal, o.state,
		       COALESCE(o.baseline_due, ''), COALESCE(o.counterfactual_due, ''), o.attempts,
		       COALESCE(p.status, ''), COALESCE(p.next_review_due, ''),
		       COALESCE(p.placement_seeded, 0), COALESCE(p.is_provisional, 0),
		       CASE WHEN p.student_id IS NULL THEN 0 ELSE 1 END
		FROM implicit_credit_ordinals o
		LEFT JOIN concept_progress p
		  ON p.student_id = o.student_id AND p.concept_id = o.target_concept_id
		WHERE o.student_id = $1 AND o.state = 'banked'
		ORDER BY o.target_concept_id, o.credit_epoch, o.ordinal
	`, studentID)
	if err != nil {
		return nil, fmt.Errorf("credit: banked details: %w", err)
	}
	defer rows.Close()
	return scanCreditOrdinalDetails(rows)
}

// AllBankedOrdinalDetails is the PostgreSQL half of the read-only dry-run
// query. It takes no locks and writes nothing.
func (s *PostgresStore) AllBankedOrdinalDetails() ([]CreditOrdinalDetail, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT o.student_id, o.target_concept_id, o.credit_epoch, o.ordinal, o.state,
		       COALESCE(o.baseline_due, ''), COALESCE(o.counterfactual_due, ''), o.attempts,
		       COALESCE(p.status, ''), COALESCE(p.next_review_due, ''),
		       COALESCE(p.placement_seeded, 0), COALESCE(p.is_provisional, 0),
		       CASE WHEN p.student_id IS NULL THEN 0 ELSE 1 END
		FROM implicit_credit_ordinals o
		LEFT JOIN concept_progress p
		  ON p.student_id = o.student_id AND p.concept_id = o.target_concept_id
		WHERE o.state = 'banked'
		ORDER BY o.student_id, o.target_concept_id, o.credit_epoch, o.ordinal
	`)
	if err != nil {
		return nil, fmt.Errorf("credit: all banked details: %w", err)
	}
	defer rows.Close()
	return scanCreditOrdinalDetails(rows)
}

func (s *PostgresStore) CreditApplications(studentID, targetConceptID string, limit int) ([]CreditApplication, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	ctx := context.Background()
	query := `
		SELECT student_id, target_concept_id, credit_epoch, ordinal, source_concept_id, source_attempt_id,
		       actual_due, proposed_due, COALESCE(stability_before, 0), COALESCE(stability_after, 0),
		       COALESCE(difficulty_before, 0), COALESCE(difficulty_after, 0), reviews_deferred
		FROM implicit_credit_applications WHERE student_id = $1`
	args := []any{studentID}
	if targetConceptID != "" {
		query += " AND target_concept_id = $2"
		args = append(args, targetConceptID)
	}
	query += fmt.Sprintf(" ORDER BY applied_at DESC LIMIT $%d", len(args)+1)
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("credit: list applications: %w", err)
	}
	defer rows.Close()
	var out []CreditApplication
	for rows.Next() {
		var a CreditApplication
		if err := rows.Scan(&a.StudentID, &a.TargetConceptID, &a.Epoch, &a.Ordinal, &a.SourceConceptID,
			&a.SourceAttemptID, &a.ActualDue, &a.ProposedDue, &a.StabilityBefore, &a.StabilityAfter,
			&a.DifficultyBefore, &a.DifficultyAfter, &a.ReviewsDeferred); err != nil {
			return nil, err
		}
		a.Applied = true
		out = append(out, a)
	}
	return out, rows.Err()
}
