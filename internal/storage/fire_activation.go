package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// fireActivationAccess is the shared SQL surface for the durable activation
// record. The record is global: it belongs to the deployment, not a learner.
type fireActivationAccess struct {
	db     *sql.DB
	rebind func(string) string
}

func (s *SQLiteStore) fireActivation() fireActivationAccess {
	return fireActivationAccess{db: s.db, rebind: func(q string) string { return q }}
}

func (s *PostgresStore) fireActivation() fireActivationAccess {
	return fireActivationAccess{db: s.db, rebind: rebindDollars}
}

const fireActivationID = "global"

var _ FireActivationStore = (*SQLiteStore)(nil)
var _ FireActivationStore = (*PostgresStore)(nil)

func (a fireActivationAccess) get() (*FireActivation, error) {
	var out FireActivation
	var ackAt, updatedAt int64
	err := a.db.QueryRowContext(context.Background(), a.rebind(`
		SELECT mode, cohort, acknowledged_by, acknowledged_at, updated_at, revision
		FROM fire_activations WHERE id = ?
	`), fireActivationID).Scan(&out.Mode, &out.Cohort, &out.AcknowledgedBy, &ackAt, &updatedAt, &out.Revision)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fire activation: read: %w", err)
	}
	out.AcknowledgedAt = time.Unix(0, ackAt).UTC()
	out.UpdatedAt = time.Unix(0, updatedAt).UTC()
	return &out, nil
}

func (a fireActivationAccess) set(in FireActivation) (FireActivation, error) {
	mode := strings.ToLower(strings.TrimSpace(in.Mode))
	cohort := strings.TrimSpace(in.Cohort)
	by := strings.TrimSpace(in.AcknowledgedBy)
	if mode != "active" && mode != "shadow" && mode != "off" {
		return FireActivation{}, fmt.Errorf("fire activation: invalid mode %q", in.Mode)
	}
	if by == "" {
		return FireActivation{}, fmt.Errorf("fire activation: acknowledged_by is required")
	}
	if mode == "active" && cohort == "" {
		return FireActivation{}, fmt.Errorf("fire activation: cohort is required for active")
	}
	now := in.UpdatedAt.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	ackAt := in.AcknowledgedAt.UTC()
	if ackAt.IsZero() {
		ackAt = now
	}
	tx, err := a.db.BeginTx(context.Background(), nil)
	if err != nil {
		return FireActivation{}, err
	}
	defer tx.Rollback()
	// The revision is part of the audit trail, so it must advance once per
	// accepted change. Reading it first and writing revision+1 lets two
	// operators read the same value and commit the same revision. RETURNING
	// makes the increment happen in the row update itself on both engines.
	var revision int64
	err = tx.QueryRowContext(context.Background(), a.rebind(`
		INSERT INTO fire_activations (id, mode, cohort, acknowledged_by, acknowledged_at, updated_at, revision)
		VALUES (?, ?, ?, ?, ?, ?, 1)
		ON CONFLICT(id) DO UPDATE SET
			mode = excluded.mode,
			cohort = excluded.cohort,
			acknowledged_by = excluded.acknowledged_by,
			acknowledged_at = excluded.acknowledged_at,
			updated_at = excluded.updated_at,
			revision = fire_activations.revision + 1
		RETURNING revision
	`), fireActivationID, mode, cohort, by, ackAt.UnixNano(), now.UnixNano()).Scan(&revision)
	if err != nil {
		return FireActivation{}, fmt.Errorf("fire activation: write: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return FireActivation{}, fmt.Errorf("fire activation: commit: %w", err)
	}
	return FireActivation{
		Mode:           mode,
		Cohort:         cohort,
		AcknowledgedBy: by,
		AcknowledgedAt: ackAt,
		UpdatedAt:      now,
		Revision:       revision,
	}, nil
}

func (s *SQLiteStore) GetFireActivation() (*FireActivation, error) {
	return s.fireActivation().get()
}

func (s *SQLiteStore) SetFireActivation(a FireActivation) (FireActivation, error) {
	return s.fireActivation().set(a)
}

func (s *PostgresStore) GetFireActivation() (*FireActivation, error) {
	return s.fireActivation().get()
}

func (s *PostgresStore) SetFireActivation(a FireActivation) (FireActivation, error) {
	return s.fireActivation().set(a)
}
