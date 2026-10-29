package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/repetition"
)

// PostgreSQL mirror of credit_sqlite.go. The accounting is identical; what
// differs is the locking. Every write takes the shared row in
// implicit_credit_targets with SELECT ... FOR UPDATE inside the transaction, so
// allocation, projection, reset and void for one learner's target queue instead
// of interleaving. A duplicate allocation is detected by the primary-key
// conflict before any write, and then read back from committed rows: a replay
// never banks, stales or projects.
var _ CreditStore = (*PostgresStore)(nil)

func (s *PostgresStore) pgTargetEpoch(ctx context.Context, tx *sql.Tx, studentID, conceptID string) (int64, error) {
	if err := lockStudentCreditTx(ctx, tx, studentID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO implicit_credit_targets (student_id, target_concept_id, credit_epoch, created_at)
		VALUES ($1, $2, 1, $3)
		ON CONFLICT (student_id, target_concept_id) DO NOTHING
	`, studentID, conceptID, time.Now().UTC().UnixNano()); err != nil {
		return 0, fmt.Errorf("credit: ensure target: %w", err)
	}
	var epoch int64
	if err := tx.QueryRowContext(ctx, `
		SELECT credit_epoch FROM implicit_credit_targets
		WHERE student_id = $1 AND target_concept_id = $2
		FOR UPDATE
	`, studentID, conceptID).Scan(&epoch); err != nil {
		return 0, fmt.Errorf("credit: lock target: %w", err)
	}
	return epoch, nil
}

// lockStudentCreditTx serialises every credit write for one learner with
// ResetProgress and DeleteAccount. Row locks on credit tables alone are not
// enough: reset deletes those rows before the progress row, so an application
// that already held its ordinal could insert its audit row after reset had
// swept applications. The advisory lock is taken first by every writer, so
// reset either runs before the application (which then finds no target) or
// after it (which then removes the application with the rest of the credit).
func lockStudentCreditTx(ctx context.Context, tx *sql.Tx, studentID string) error {
	if strings.TrimSpace(studentID) == "" {
		return fmt.Errorf("credit: empty student")
	}
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, studentID); err != nil {
		return fmt.Errorf("credit: lock student: %w", err)
	}
	return nil
}

func (s *PostgresStore) pgAccount(ctx context.Context, q interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}, studentID, conceptID string, epoch int64) (CreditAccount, error) {
	acc := CreditAccount{ConceptID: conceptID, Epoch: epoch}
	if err := q.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(retained_units), 0) FROM implicit_credit_events
		WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3 AND voided_at IS NULL
	`, studentID, conceptID, epoch).Scan(&acc.RetainedUnits); err != nil {
		return acc, fmt.Errorf("credit: sum retained: %w", err)
	}
	if err := q.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN state = 'projected' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN state = 'stale' THEN 1 ELSE 0 END), 0)
		FROM implicit_credit_ordinals
		WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3
	`, studentID, conceptID, epoch).Scan(&acc.Banked, &acc.Projected, &acc.Stale); err != nil {
		return acc, fmt.Errorf("credit: count ordinals: %w", err)
	}
	return acc, nil
}

func (s *PostgresStore) AllocateCredit(studentID, attemptID string, epoch int64, candidates []CreditCandidate) (CreditAllocation, bool, error) {
	out := CreditAllocation{StudentID: studentID, AttemptID: attemptID, Epoch: epoch}
	if studentID == "" || attemptID == "" {
		return out, false, fmt.Errorf("credit: empty student or attempt")
	}
	raw := make([]repetition.Candidate, 0, len(candidates))
	for _, c := range candidates {
		raw = append(raw, repetition.Candidate{ConceptID: c.ConceptID, RawUnits: c.RawUnits})
	}
	scaled, err := repetition.Allocate(raw, repetition.Micro)
	if err != nil {
		return out, false, fmt.Errorf("credit: allocate: %w", err)
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return out, false, err
	}
	defer tx.Rollback()
	if err := lockStudentCreditTx(ctx, tx, studentID); err != nil {
		return out, false, err
	}

	var marker int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO implicit_credit_allocations
			(student_id, source_attempt_id, raw_units, allocated_units, target_count, dropped_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (student_id, source_attempt_id) DO NOTHING
		RETURNING 1
	`, studentID, attemptID, sumRaw(raw), sumScaled(scaled), len(scaled), len(raw)-len(scaled), time.Now().UTC().UnixNano()).Scan(&marker)
	if errors.Is(err, sql.ErrNoRows) {
		committed, err := s.pgReplayAllocation(ctx, tx, studentID, attemptID)
		if err != nil {
			return CreditAllocation{}, false, err
		}
		if err := tx.Commit(); err != nil {
			return CreditAllocation{}, false, err
		}
		return committed, false, nil
	}
	if err != nil {
		return CreditAllocation{}, false, fmt.Errorf("credit: insert allocation: %w", err)
	}

	out.RawUnits = sumRaw(raw)
	out.AllocatedUnits = sumScaled(scaled)
	out.TargetCount = len(scaled)
	out.DroppedCount = len(raw) - len(scaled)
	out.Events = make([]CreditEvent, 0, len(scaled))

	for _, a := range scaled {
		targetEpoch, err := s.pgTargetEpoch(ctx, tx, studentID, a.ConceptID)
		if err != nil {
			return CreditAllocation{}, false, err
		}
		if epoch > 0 && targetEpoch != epoch {
			return CreditAllocation{}, false, fmt.Errorf("%w: target %s is at epoch %d, not %d", ErrCreditClosedEpoch, a.ConceptID, targetEpoch, epoch)
		}
		useEpoch := targetEpoch
		if epoch > 0 {
			useEpoch = epoch
		}
		acc, err := s.pgAccount(ctx, tx, studentID, a.ConceptID, useEpoch)
		if err != nil {
			return CreditAllocation{}, false, err
		}
		retained := a.Units
		capacity := int64(repetition.MaxBank)*repetition.Micro - acc.RetainedUnits
		if capacity < 0 {
			capacity = 0
		}
		discarded := int64(0)
		reason := ""
		if retained > capacity {
			discarded = retained - capacity
			retained = capacity
			reason = "epoch_bank_limit"
		}
		banked := int64(0)
		if want := (acc.RetainedUnits + retained) / repetition.Micro; want > int64(acc.Banked) {
			banked = want - int64(acc.Banked)
		}

		// One implicit advance per review cycle. Ten answers in one sitting earn
		// ten units against this target, and each used to become its own ordinal,
		// which ProjectCredit chained into ten spaced reviews: a three-day
		// interval reached 30,595 days. The cycle closes when the review this
		// credit stands for is consumed, so at most one ordinal may be banked per
		// epoch. Refused units are recorded, not dropped.
		var advanced int
		if err := tx.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM implicit_credit_ordinals
			WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3
			  AND state IN ('banked', 'projected')
		`, studentID, a.ConceptID, useEpoch).Scan(&advanced); err != nil {
			return CreditAllocation{}, false, fmt.Errorf("credit: cycle check: %w", err)
		}
		if advanced > 0 && banked > 0 {
			discarded += retained
			retained = 0
			banked = 0
			reason = "cycle_already_advanced"
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO implicit_credit_events
				(student_id, source_attempt_id, target_concept_id, credit_epoch,
				 scaled_units, retained_units, discarded_units, discard_reason, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (student_id, source_attempt_id, target_concept_id) DO NOTHING
		`, studentID, attemptID, a.ConceptID, useEpoch, a.Units, retained, discarded, reason, time.Now().UTC().UnixNano()); err != nil {
			return CreditAllocation{}, false, fmt.Errorf("credit: insert event: %w", err)
		}
		for i := int64(0); i < banked; i++ {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO implicit_credit_ordinals
					(student_id, target_concept_id, credit_epoch, ordinal, state, banked_at, attempts)
				VALUES ($1, $2, $3, $4, 'banked', $5, 0)
			`, studentID, a.ConceptID, useEpoch, int64(acc.Banked)+i+1, time.Now().UTC().UnixNano()); err != nil {
				return CreditAllocation{}, false, fmt.Errorf("credit: insert ordinal: %w", err)
			}
		}
		out.Events = append(out.Events, CreditEvent{
			AttemptID: attemptID, ConceptID: a.ConceptID, Epoch: useEpoch,
			ScaledUnits: a.Units, RetainedUnits: retained, DiscardedUnits: discarded,
			DiscardReason: reason, Banked: int(banked),
		})
	}
	if err := tx.Commit(); err != nil {
		return CreditAllocation{}, false, fmt.Errorf("credit: commit: %w", err)
	}
	return out, true, nil
}

func (s *PostgresStore) pgReplayAllocation(ctx context.Context, tx *sql.Tx, studentID, attemptID string) (CreditAllocation, error) {
	out := CreditAllocation{StudentID: studentID, AttemptID: attemptID}
	if err := tx.QueryRowContext(ctx, `
		SELECT raw_units, allocated_units, target_count, dropped_count
		FROM implicit_credit_allocations WHERE student_id = $1 AND source_attempt_id = $2
	`, studentID, attemptID).Scan(&out.RawUnits, &out.AllocatedUnits, &out.TargetCount, &out.DroppedCount); err != nil {
		return out, fmt.Errorf("credit: replay allocation: %w", err)
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT target_concept_id, credit_epoch, scaled_units, retained_units, discarded_units, discard_reason
		FROM implicit_credit_events
		WHERE student_id = $1 AND source_attempt_id = $2
		ORDER BY target_concept_id ASC
	`, studentID, attemptID)
	if err != nil {
		return out, fmt.Errorf("credit: replay events: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ev CreditEvent
		ev.AttemptID = attemptID
		if err := rows.Scan(&ev.ConceptID, &ev.Epoch, &ev.ScaledUnits, &ev.RetainedUnits, &ev.DiscardedUnits, &ev.DiscardReason); err != nil {
			return out, fmt.Errorf("credit: scan replay: %w", err)
		}
		out.Events = append(out.Events, ev)
	}
	if len(out.Events) != out.TargetCount {
		return out, fmt.Errorf("credit: allocation %s records %d targets but has %d events", attemptID, out.TargetCount, len(out.Events))
	}
	return out, rows.Err()
}

func (s *PostgresStore) AccountOf(studentID, conceptID string, epoch int64) (CreditAccount, error) {
	ctx := context.Background()
	if epoch <= 0 {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return CreditAccount{}, err
		}
		defer tx.Rollback()
		epoch, err = s.pgTargetEpoch(ctx, tx, studentID, conceptID)
		if err != nil {
			return CreditAccount{}, err
		}
		if err := tx.Commit(); err != nil {
			return CreditAccount{}, err
		}
	}
	return s.pgAccount(ctx, s.db, studentID, conceptID, epoch)
}

func (s *PostgresStore) ProjectCredit(studentID, conceptID string, epoch int64, stepDays int, now time.Time) ([]CreditOrdinal, error) {
	if stepDays < 1 || stepDays > repetition.MaxStepDays {
		return nil, fmt.Errorf("credit: step %d outside [1,%d]", stepDays, repetition.MaxStepDays)
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := lockStudentCreditTx(ctx, tx, studentID); err != nil {
		return nil, err
	}

	if epoch <= 0 {
		var err error
		epoch, err = s.pgTargetEpoch(ctx, tx, studentID, conceptID)
		if err != nil {
			return nil, err
		}
	}

	var due sql.NullString
	if err := tx.QueryRowContext(ctx, `
		SELECT next_review_due FROM concept_progress
		WHERE student_id = $1 AND concept_id = $2 FOR UPDATE
	`, studentID, conceptID).Scan(&due); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("credit: read due: %w", err)
	}
	if !due.Valid || due.String == "" {
		return nil, nil
	}
	baseline := due.String
	anchor, err := time.Parse(time.RFC3339, baseline)
	if err != nil {
		anchor, err = time.Parse("2006-01-02 15:04:05", baseline)
		if err != nil {
			return nil, fmt.Errorf("credit: parse due %q: %w", baseline, err)
		}
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT ordinal, COALESCE(counterfactual_due, ''), state, attempts
		FROM implicit_credit_ordinals
		WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3
		ORDER BY ordinal ASC
	`, studentID, conceptID, epoch)
	if err != nil {
		return nil, fmt.Errorf("credit: read ordinals: %w", err)
	}
	type ordinalState struct {
		ordinal int
		due     string
		state   string
		tries   int
	}
	var all []ordinalState
	for rows.Next() {
		var o ordinalState
		if err := rows.Scan(&o.ordinal, &o.due, &o.state, &o.tries); err != nil {
			rows.Close()
			return nil, err
		}
		all = append(all, o)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var out []CreditOrdinal
	prev := anchor
	for i := range all {
		o := &all[i]
		switch o.state {
		case "stale":
			return nil, fmt.Errorf("%w: ordinal %d", ErrCreditOrdinalStale, o.ordinal)
		case "projected":
			parsed, err := time.Parse(time.RFC3339, o.due)
			if err != nil {
				return nil, fmt.Errorf("credit: parse stored counterfactual: %w", err)
			}
			prev = parsed
			out = append(out, CreditOrdinal{ConceptID: conceptID, Epoch: epoch, Ordinal: o.ordinal, State: o.state, BaselineDue: baseline, CounterfactualDue: o.due, Attempts: o.tries})
		case "banked":
			next := prev.AddDate(0, 0, stepDays)
			if floor := now.UTC().AddDate(0, 0, 1); next.Before(floor) {
				next = floor
			}
			res, err := tx.ExecContext(ctx, `
				UPDATE implicit_credit_ordinals
				SET state = 'projected', projected_at = $1, baseline_due = $2,
				    counterfactual_due = $3, attempts = attempts + 1
				WHERE student_id = $4 AND target_concept_id = $5 AND credit_epoch = $6
				  AND ordinal = $7 AND state = 'banked'
			`, time.Now().UTC().UnixNano(), baseline, next.UTC().Format(time.RFC3339), studentID, conceptID, epoch, o.ordinal)
			if err != nil {
				return nil, fmt.Errorf("credit: project ordinal: %w", err)
			}
			if n, _ := res.RowsAffected(); n != 1 {
				return nil, fmt.Errorf("%w: ordinal %d", ErrCreditOrdinalStateChanged, o.ordinal)
			}
			prev = next
			out = append(out, CreditOrdinal{ConceptID: conceptID, Epoch: epoch, Ordinal: o.ordinal, State: "projected", BaselineDue: baseline, CounterfactualDue: next.UTC().Format(time.RFC3339)})
		default:
			return nil, fmt.Errorf("%w: ordinal %d is %q", ErrCreditOrdinalStateChanged, o.ordinal, o.state)
		}
	}

	if s.creditProbe != nil {
		s.creditProbe(func(query string, args ...any) error {
			_, err := tx.ExecContext(ctx, query, args...)
			return err
		})
	}

	// The chain was built from the due date read above. If it moved, the whole
	// projection is stale; rolling back is the only honest outcome.
	var nowDue sql.NullString
	if err := tx.QueryRowContext(ctx, `
		SELECT next_review_due FROM concept_progress WHERE student_id = $1 AND concept_id = $2
	`, studentID, conceptID).Scan(&nowDue); err != nil {
		return nil, fmt.Errorf("credit: reread due: %w", err)
	}
	if !nowDue.Valid || nowDue.String != baseline {
		return nil, fmt.Errorf("%w: was %q, now %q", ErrCreditDueConflict, baseline, nowDue.String)
	}
	// The review this chain was built for is now this cycle's evidence, so the
	// cycle closes: the next implicit credit may bank once the review is next
	// due. Without this the epoch would stay at 1 and the guard in allocate would
	// allow one implicit advance in a lifetime.
	// Only a call that actually consumed an ordinal closes the cycle. A retry
	// with nothing banked must not advance the epoch again, or a repeated
	// projection would walk the cycle forward for free.
	if _, err := tx.ExecContext(ctx, `
		UPDATE implicit_credit_targets SET credit_epoch = $1
		WHERE student_id = $2 AND target_concept_id = $3 AND credit_epoch = $4
		  AND EXISTS (
			SELECT 1 FROM implicit_credit_ordinals
			WHERE student_id = $2 AND target_concept_id = $3 AND credit_epoch = $4
			  AND state = 'projected'
		  )
	`, epoch+1, studentID, conceptID, epoch); err != nil {
		return nil, fmt.Errorf("credit: advance epoch: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("credit: commit projection: %w", err)
	}
	return out, nil
}

func (s *PostgresStore) ResetCreditEpoch(studentID, conceptID string) (int64, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	if err := lockStudentCreditTx(ctx, tx, studentID); err != nil {
		return 0, err
	}

	epoch, err := s.pgTargetEpoch(ctx, tx, studentID, conceptID)
	if err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE implicit_credit_ordinals SET state = 'stale'
		WHERE student_id = $1 AND target_concept_id = $2 AND credit_epoch = $3 AND state = 'banked'
	`, studentID, conceptID, epoch); err != nil {
		return 0, fmt.Errorf("credit: stale ordinals: %w", err)
	}
	next := epoch + 1
	if _, err := tx.ExecContext(ctx, `
		UPDATE implicit_credit_targets SET credit_epoch = $1
		WHERE student_id = $2 AND target_concept_id = $3 AND credit_epoch = $4
	`, next, studentID, conceptID, epoch); err != nil {
		return 0, fmt.Errorf("credit: bump epoch: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return next, nil
}

func (s *PostgresStore) VoidCreditEvent(studentID, attemptID, conceptID, reason string) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := lockStudentCreditTx(ctx, tx, studentID); err != nil {
		return err
	}

	epoch, err := s.pgTargetEpoch(ctx, tx, studentID, conceptID)
	if err != nil {
		return err
	}
	var evEpoch, retained int64
	var voidedAt sql.NullInt64
	if err := tx.QueryRowContext(ctx, `
		SELECT credit_epoch, retained_units, voided_at FROM implicit_credit_events
		WHERE student_id = $1 AND source_attempt_id = $2 AND target_concept_id = $3
	`, studentID, attemptID, conceptID).Scan(&evEpoch, &retained, &voidedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCreditOrdinalMissing
		}
		return err
	}
	if evEpoch != epoch {
		return fmt.Errorf("%w: event epoch %d, target epoch %d", ErrCreditClosedEpoch, evEpoch, epoch)
	}
	if voidedAt.Valid {
		return nil
	}
	acc, err := s.pgAccount(ctx, tx, studentID, conceptID, epoch)
	if err != nil {
		return err
	}
	if acc.RetainedUnits-retained < int64(acc.Banked)*repetition.Micro {
		return fmt.Errorf("%w: retained %d minus %d cannot fund %d ordinals", ErrCreditFundsProjection, acc.RetainedUnits, retained, acc.Banked)
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE implicit_credit_events
		SET voided_at = $1, voided_by = 'credit', void_reason = $2
		WHERE student_id = $3 AND source_attempt_id = $4 AND target_concept_id = $5
	`, time.Now().UTC().UnixNano(), reason, studentID, attemptID, conceptID); err != nil {
		return fmt.Errorf("credit: void event: %w", err)
	}
	return tx.Commit()
}

func (s *PostgresStore) GetBankedOrdinals(studentID string) (map[string]int, error) {
	ctx := context.Background()
	rows, err := s.db.QueryContext(ctx, `
		SELECT target_concept_id, COUNT(*)
		FROM implicit_credit_ordinals
		WHERE student_id = $1 AND state = 'banked'
		GROUP BY target_concept_id
	`, studentID)
	if err != nil {
		return nil, fmt.Errorf("credit: get banked ordinals: %w", err)
	}
	defer rows.Close()
	out := make(map[string]int)
	for rows.Next() {
		var cid string
		var count int
		if err := rows.Scan(&cid, &count); err != nil {
			return nil, err
		}
		out[cid] = count
	}
	return out, rows.Err()
}
