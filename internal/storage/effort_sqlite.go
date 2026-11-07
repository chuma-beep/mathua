package storage

import (
	"time"
)

// RecordEffortSignal appends one accountability observation; see EffortStore.
func (s *SQLiteStore) RecordEffortSignal(sig EffortSignal) error {
	return recordEffortSignalTx(s.db, func(q string) string { return q }, sig)
}

// EffortSignalsSince returns recent observations, newest first.
func (s *SQLiteStore) EffortSignalsSince(studentID string, since time.Time, limit int) ([]EffortSignal, error) {
	return effortSignalsSince(s.db, func(q string) string { return q }, studentID, since, limit)
}

// CountXPDeductions counts prior accountability deductions in the window.
func (s *SQLiteStore) CountXPDeductions(studentID string, since time.Time) (int, error) {
	return countXPDeductions(s.db, func(q string) string { return q }, studentID, since)
}

// SumXPDeductionsSince returns the positive deduction magnitude.
func (s *SQLiteStore) SumXPDeductionsSince(studentID string, since time.Time) (int, error) {
	return sumAccountabilityDeductionsSince(s.db, func(q string) string { return q }, studentID, since)
}

// SumXPEarnedSince returns the positive XP earned since a point in time.
func (s *SQLiteStore) SumXPEarnedSince(studentID string, since time.Time) (int, error) {
	return sumXPEarnedSince(s.db, func(q string) string { return q }, studentID, since)
}

// ApplyXPAdjustment writes a capped signed ledger row atomically.
func (s *SQLiteStore) ApplyXPAdjustment(studentID string, amount int, reason, eventKey string) (int, error) {
	return applyXPAdjustmentTx(s.db, func(q string) string { return q }, "", studentID, amount, reason, eventKey)
}
