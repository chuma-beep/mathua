package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Admin V2 storage (Postgres). Mirrors adminv2_sqlite.go statement for statement so the two
// stores cannot disagree about a rule — a guard only one store enforces is a guard a Postgres
// deployment does not have. The scan helpers are shared from the SQLite file on purpose.

func (s *PostgresStore) CountStaff() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE role IN ('owner','admin','moderator')").Scan(&n)
	return n, err
}

func (s *PostgresStore) CountOwners() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM students WHERE role = 'owner'").Scan(&n)
	return n, err
}

func (s *PostgresStore) CountReportsByStatus() (map[string]int, error) {
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

func (s *PostgresStore) GetReport(id int64) (*QuestionReport, error) {
	r, err := scanReportRow(s.db.QueryRow(reportSelectColumns+` WHERE id = $1`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *PostgresStore) UpdateReportStatusAudited(actorID string, id int64, status, resolution string) (QuestionReport, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return QuestionReport{}, fmt.Errorf("begin moderation: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	before, err := scanReportRow(tx.QueryRow(reportSelectColumns+` WHERE id = $1`, id))
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
		`UPDATE question_reports SET status = $1, resolution = $2, resolved_by = $3, resolved_at = $4 WHERE id = $5`,
		status, resolution, resolvedBy, resolvedAt, id,
	); err != nil {
		return QuestionReport{}, fmt.Errorf("update report: %w", err)
	}

	after, err := scanReportRow(tx.QueryRow(reportSelectColumns+` WHERE id = $1`, id))
	if err != nil {
		return QuestionReport{}, err
	}
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES ($1, $2, 'report', $3, $4, $5, $6)`,
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

func (s *PostgresStore) ListAdminAuditFor(entityType, entityID string, limit int) ([]AdminAuditRecord, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.Query(
		`SELECT id, actor_id, action, entity_type, entity_id, before_json, after_json, created_at
		 FROM admin_audit WHERE entity_type = $1 AND entity_id = $2
		 ORDER BY id DESC LIMIT $3`, entityType, entityID, limit)
	if err != nil {
		return nil, fmt.Errorf("list audit for entity: %w", err)
	}
	defer rows.Close()
	return scanAuditRows(rows)
}

func (s *PostgresStore) CreateAdminInvitation(inv AdminInvitation, tokenHash string) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		`INSERT INTO admin_invitations (token_hash, email, role, invited_by, created_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		tokenHash, strings.ToLower(strings.TrimSpace(inv.Email)), string(inv.Role), inv.InvitedBy,
		inv.CreatedAt, inv.ExpiresAt,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create invitation: %w", err)
	}
	return id, nil
}

func (s *PostgresStore) ListAdminInvitations(limit, offset int) ([]AdminInvitation, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.db.Query(invitationSelectColumns+` ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list invitations: %w", err)
	}
	defer rows.Close()
	return scanInvitationRows(rows)
}

func (s *PostgresStore) RevokeAdminInvitation(actorID string, id int64) (AdminInvitation, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return AdminInvitation{}, fmt.Errorf("begin revoke: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	before, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE id = $1`, id))
	if err == sql.ErrNoRows {
		return AdminInvitation{}, sql.ErrNoRows
	}
	if err != nil {
		return AdminInvitation{}, err
	}
	if before.AcceptedAt != "" || before.RevokedAt != "" {
		return AdminInvitation{}, sql.ErrNoRows
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := tx.Exec(`UPDATE admin_invitations SET revoked_at = $1, revoked_by = $2 WHERE id = $3`, now, actorID, id); err != nil {
		return AdminInvitation{}, fmt.Errorf("revoke invitation: %w", err)
	}
	after, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE id = $1`, id))
	if err != nil {
		return AdminInvitation{}, err
	}
	afterJSON, _ := json.Marshal(after)
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES ($1, $2, 'invitation', $3, '', $4, $5)`,
		actorID, ActionAdminInviteRevoke, fmt.Sprintf("%d", id), string(afterJSON), now,
	); err != nil {
		return AdminInvitation{}, fmt.Errorf("record revoke: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return AdminInvitation{}, fmt.Errorf("commit revoke: %w", err)
	}
	return after, nil
}

func (s *PostgresStore) ConsumeAdminInvitation(tokenHash, studentID string) (*AdminInvitation, bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, false, fmt.Errorf("begin accept: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	inv, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE token_hash = $1`, tokenHash))
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
		`UPDATE admin_invitations SET accepted_at = $1, accepted_by = $2 WHERE id = $3`,
		now, studentID, inv.ID,
	); err != nil {
		return nil, false, fmt.Errorf("accept invitation: %w", err)
	}
	if _, err := tx.Exec(`UPDATE students SET role = $1 WHERE id = $2`, string(inv.Role), studentID); err != nil {
		return nil, false, fmt.Errorf("grant invited role: %w", err)
	}
	after, err := scanInvitationRow(tx.QueryRow(invitationSelectColumns+` WHERE id = $1`, inv.ID))
	if err != nil {
		return nil, false, err
	}
	afterJSON, _ := json.Marshal(after)
	if _, err := tx.Exec(
		`INSERT INTO admin_audit (actor_id, action, entity_type, entity_id, before_json, after_json, created_at)
		 VALUES ($1, $2, 'invitation', $3, '', $4, $5)`,
		studentID, ActionAdminInviteAccept, fmt.Sprintf("%d", inv.ID), string(afterJSON), now,
	); err != nil {
		return nil, false, fmt.Errorf("record accept: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("commit accept: %w", err)
	}
	return &after, true, nil
}

func (s *PostgresStore) FindAdminInvitationByTokenHash(tokenHash string) (*AdminInvitation, error) {
	inv, err := scanInvitationRow(s.db.QueryRow(invitationSelectColumns+` WHERE token_hash = $1`, tokenHash))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (s *PostgresStore) ListStaff(limit, offset int) ([]AdminUser, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return queryAdminUsersPG(s.db,
		`SELECT `+pgAdminUserColumns+` FROM students
		 WHERE role IN ('owner','admin','moderator')
		 ORDER BY CASE role WHEN 'owner' THEN 0 WHEN 'admin' THEN 1 ELSE 2 END, created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset)
}

func queryAdminUsersPG(q adminQueryerPG, query string, args ...interface{}) ([]AdminUser, error) {
	rows, err := q.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query admin users: %w", err)
	}
	defer rows.Close()
	out := []AdminUser{}
	for rows.Next() {
		var u AdminUser
		var name, username, email, role, createdAt sql.NullString
		var hasPw bool
		if err := rows.Scan(&u.ID, &name, &username, &email, &role, &createdAt, &hasPw); err != nil {
			return nil, err
		}
		u.Name = name.String
		u.Username = username.String
		u.Email = email.String
		u.Role = NormalizeRole(role.String)
		u.CreatedAt = createdAt.String
		u.HasPassword = hasPw
		out = append(out, u)
	}
	return out, rows.Err()
}
