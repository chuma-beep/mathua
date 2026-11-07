package storage

import (
	"time"
)

// RecordEffortSignal appends one accountability observation; see EffortStore.
func (s *PostgresStore) RecordEffortSignal(sig EffortSignal) error {
	return recordEffortSignalTx(s.db, rebindDollars, sig)
}

// EffortSignalsSince returns recent observations, newest first.
func (s *PostgresStore) EffortSignalsSince(studentID string, since time.Time, limit int) ([]EffortSignal, error) {
	return effortSignalsSince(s.db, rebindDollars, studentID, since, limit)
}

// CountXPDeductions counts prior accountability deductions in the window.
func (s *PostgresStore) CountXPDeductions(studentID string, since time.Time) (int, error) {
	return countXPDeductions(s.db, rebindDollars, studentID, since)
}

// SumXPDeductionsSince returns the positive deduction magnitude.
func (s *PostgresStore) SumXPDeductionsSince(studentID string, since time.Time) (int, error) {
	return sumAccountabilityDeductionsSince(s.db, rebindDollars, studentID, since)
}

// SumXPEarnedSince returns the positive XP earned since a point in time.
func (s *PostgresStore) SumXPEarnedSince(studentID string, since time.Time) (int, error) {
	return sumXPEarnedSince(s.db, rebindDollars, studentID, since)
}

// ApplyXPAdjustment writes a capped signed ledger row atomically.
func (s *PostgresStore) ApplyXPAdjustment(studentID string, amount int, reason, eventKey string) (int, error) {
	return applyXPAdjustmentTx(s.db, rebindDollars, " FOR UPDATE", studentID, amount, reason, eventKey)
}
