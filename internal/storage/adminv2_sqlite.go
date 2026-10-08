package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Admin V2 storage (SQLite). The Postgres twin lives in adminv2_postgres.go and mirrors this
// statement for statement; ADR-012 requires the two stores agree on an observable result.
//
// Everything here is additive to V1. The V1 methods (SetRoleAudited, ListAdminAudit, the report
// CRUD) keep their names and contracts; these are the roles above `admin`, the moderation state
// on reports, staff invitations, and the per-entity audit query.

func (s *SQLiteStore) CountStaff() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE role IN ('owner','admin','moderator')").Scan(&n)
	return n, err
}

func (s *SQLiteStore) CountOwners() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE role = 'owner'").Scan(&n)
	return n, err
}

// CountReportsByStatus returns one count per status that has rows. Statuses with none are
// absent, so a caller that wants zeros supplies them.
func (s *SQLiteStore) CountReportsByStatus() (map[string]int, error) {
	rows, err := s.db.Query("SELECT status, COUNT(*) FROM question_reports GROUP BY status")
	if err != nil {
		return nil, fmt.Errorf("count reports: %w", err)
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			return nil, err
		}
		out[status] = n
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetReport(id int64) (*QuestionReport, error) {
	row := s.db.QueryRow(reportSelectColumns+` WHERE id = ?`, id)
	r, err := scanReportRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateReportStatusAudited applies a moderation decision and records it together.
//
// The report row is edited, never deleted. Resolving or dismissing stamps who decided and when
// and keeps the moderator's reason; reopening clears those stamps so the row never claims a
// resolution that a later decision superseded. The audit event and the state change share one
// transaction, so a decision cannot exist without a record of who made it.
func (s *SQLiteStore) UpdateReportStatusAudited(actorID string, id int64, status, resolution string) (QuestionReport, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return QuestionReport{}, fmt.Errorf("begin moderation: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	before, err := scanReportRow(tx.QueryRow(reportSelectColumns+` WHERE id = ?`, id))
	if err == sql.ErrNoRows {
		return QuestionReport{}, sql.ErrNoRows
	}
	if err != nil {
		return QuestionReport{}, err
	}
	if before.Status == status {
		return before, ErrReportUnchanged
	}

	settled := status == "resolved" || status == "dismissed"
	resolvedBy, resolvedAt := "", ""
	if settled {
		resolvedBy = actorID
		resolvedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if _, err := tx.Exec(
		`UPDATE question_reports SET status = ?, resolution = ?, resolved_by = ?, resolved_at = ? WHERE id = ?`,
		status, resolution, resolvedBy, resolvedAt, id,
	); err != nil {
		return QuestionReport{}, fmt.Errorf("update report: %w", err)
	}

	after, err := scanReportRow(tx.QueryRow(reportSelectColumns+` WHERE id = ?`, id))
	if err != nil {
		return QuestionReport{}, err
	}
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES (?, ?, 'report', ?, ?, ?, ?)`,
		actorID, ActionReportStatusChange, fmt.Sprintf("%d", id), string(beforeJSON), string(afterJSON),
		time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return QuestionReport{}, fmt.Errorf("record moderation: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return QuestionReport{}, fmt.Errorf("commit moderation: %w", err)
	}
	return after, nil
}

func (s *SQLiteStore) ListAdminAuditFor(entityType, entityID string, limit int) ([]AdminAuditRecord, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.Query(
		`SELECT id, actor_id, action, entity_type, entity_id, before_json, after_json, created_at
		 FROM admin_audit WHERE entity_type = ? AND entity_id = ?
		 ORDER BY id DESC LIMIT ?`, entityType, entityID, limit)
	if err != nil {
		return nil, fmt.Errorf("list audit for entity: %w", err)
	}
	defer rows.Close()
	return scanAuditRows(rows)
}

// ── Invitations ───────────────────────────────────────────────────────────

func (s *SQLiteStore) CreateAdminInvitation(inv AdminInvitation, tokenHash string) (int64, error) {
	res, err := s.db.Exec(
		`INSERT INTO admin_invitations (token_hash, email, role, invited_by, created_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		tokenHash, strings.ToLower(strings.TrimSpace(inv.Email)), string(inv.Role), inv.InvitedBy,
		inv.CreatedAt, inv.ExpiresAt,
	)
	if err != nil {
		return 0, fmt.Errorf("create invitation: %w", err)
	}
	return res.LastInsertId()
}

func (s *SQLiteStore) ListAdminInvitations(limit, offset int) ([]AdminInvitation, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.db.Query(invitationSelectColumns+` ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list invitations: %w", err)
	}
	defer rows.Close()
	return scanInvitationRows(rows)
}

func (s *SQLiteStore) RevokeAdminInvitation(actorID string, id int64) (AdminInvitation, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return AdminInvitation{}, fmt.Errorf("begin revoke: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	before, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE id = ?`, id))
	if err == sql.ErrNoRows {
		return AdminInvitation{}, sql.ErrNoRows
	}
	if err != nil {
		return AdminInvitation{}, err
	}
	// A settled invitation cannot be revoked again: it is either accepted (the person has the
	// role; revoke the role instead) or already revoked. Both answer ErrNoRows so the caller
	// cannot act on a token that is no longer live.
	if before.AcceptedAt != "" || before.RevokedAt != "" {
		return AdminInvitation{}, sql.ErrNoRows
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := tx.Exec(`UPDATE admin_invitations SET revoked_at = ?, revoked_by = ? WHERE id = ?`, now, actorID, id); err != nil {
		return AdminInvitation{}, fmt.Errorf("revoke invitation: %w", err)
	}
	after, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE id = ?`, id))
	if err != nil {
		return AdminInvitation{}, err
	}
	afterJSON, _ := json.Marshal(after)
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES (?, ?, 'invitation', ?, '', ?, ?)`,
		actorID, ActionAdminInviteRevoke, fmt.Sprintf("%d", id), string(afterJSON), now,
	); err != nil {
		return AdminInvitation{}, fmt.Errorf("record revoke: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return AdminInvitation{}, fmt.Errorf("commit revoke: %w", err)
	}
	return after, nil
}

// ConsumeAdminInvitation burns a token and grants the role it carried, in one transaction. The
// role change and the audit event are part of the same unit as the burn, so there is no state in
// which someone holds a staff role with no record of how they got it — and no state in which a
// token is spent but the role was not granted.
func (s *SQLiteStore) ConsumeAdminInvitation(tokenHash, studentID string) (*AdminInvitation, bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, false, fmt.Errorf("begin accept: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	inv, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE token_hash = ?`, tokenHash))
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	exp, perr := time.Parse(time.RFC3339, inv.ExpiresAt)
	if perr != nil {
		return nil, false, fmt.Errorf("parse invitation expiry: %w", perr)
	}
	if inv.AcceptedAt != "" || inv.RevokedAt != "" || time.Now().UTC().After(exp) {
		return nil, false, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := tx.Exec(
		`UPDATE admin_invitations SET accepted_at = ?, accepted_by = ? WHERE id = ?`,
		now, studentID, inv.ID,
	); err != nil {
		return nil, false, fmt.Errorf("accept invitation: %w", err)
	}
	if _, err := tx.Exec(`UPDATE students SET role = ? WHERE id = ?`, string(inv.Role), studentID); err != nil {
		return nil, false, fmt.Errorf("grant invited role: %w", err)
	}
	after, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE id = ?`, inv.ID))
	if err != nil {
		return nil, false, err
	}
	afterJSON, _ := json.Marshal(after)
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES (?, ?, 'invitation', ?, '', ?, ?)`,
		studentID, ActionAdminInviteAccept, fmt.Sprintf("%d", inv.ID), string(afterJSON), now,
	); err != nil {
		return nil, false, fmt.Errorf("record accept: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("commit accept: %w", err)
	}
	return &after, true, nil
}

// ── shared scan helpers ───────────────────────────────────────────────────

const reportSelectColumns = `SELECT id, reporter_id, concept_id, kind, question, expected, explanation,
	lesson_id, source, session_id, attempt_id, reason, detail, status,
	resolution, resolved_by, resolved_at, created_at
	FROM question_reports`

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanReportRow(row rowScanner) (QuestionReport, error) {
	var r QuestionReport
	var created string
	err := row.Scan(
		&r.ID, &r.ReporterID, &r.ConceptID, &r.Kind, &r.Question, &r.Expected,
		&r.Explanation, &r.LessonID, &r.Source, &r.SessionID, &r.AttemptID,
		&r.Reason, &r.Detail, &r.Status, &r.Resolution, &r.ResolvedBy, &r.ResolvedAt, &created,
	)
	if err != nil {
		return QuestionReport{}, err
	}
	r.CreatedAt = parseReportTime(created)
	return r, nil
}

const invitationSelectColumns = `SELECT id, email, role, invited_by, created_at, expires_at,
	accepted_at, accepted_by, revoked_at, revoked_by
	FROM admin_invitations`

func scanInvitationRow(row rowScanner) (AdminInvitation, error) {
	var inv AdminInvitation
	var role string
	err := row.Scan(
		&inv.ID, &inv.Email, &role, &inv.InvitedBy, &inv.CreatedAt, &inv.ExpiresAt,
		&inv.AcceptedAt, &inv.AcceptedBy, &inv.RevokedAt, &inv.RevokedBy,
	)
	if err != nil {
		return AdminInvitation{}, err
	}
	inv.Role = NormalizeRole(role)
	return inv, nil
}

func scanInvitationRows(rows *sql.Rows) ([]AdminInvitation, error) {
	out := []AdminInvitation{}
	for rows.Next() {
		inv, err := scanInvitationRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, rows.Err()
}

func scanAuditRows(rows *sql.Rows) ([]AdminAuditRecord, error) {
	out := []AdminAuditRecord{}
	for rows.Next() {
		var r AdminAuditRecord
		var createdAt sql.NullString
		if err := rows.Scan(&r.ID, &r.ActorID, &r.Action, &r.EntityType, &r.EntityID,
			&r.BeforeJSON, &r.AfterJSON, &createdAt); err != nil {
			return nil, err
		}
		r.CreatedAt = createdAt.String
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) FindAdminInvitationByTokenHash(tokenHash string) (*AdminInvitation, error) {
	inv, err := scanInvitationRow(s.db.QueryRow(invitationSelectColumns+` WHERE token_hash = ?`, tokenHash))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (s *SQLiteStore) ListStaff(limit, offset int) ([]AdminUser, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return queryAdminUsers(s.db,
		`SELECT `+adminUserColumns+` FROM students
		 WHERE role IN ('owner','admin','moderator')
		 ORDER BY CASE role WHEN 'owner' THEN 0 WHEN 'admin' THEN 1 ELSE 2 END, created_at DESC
		 LIMIT ? OFFSET ?`, limit, offset)
}
