package storage

import (
	"database/sql"
	"fmt"
	"strings"
)

// Search is an optional extension of the admin reads. Keeping it separate from
// Repository avoids widening an interface every test fake has to satisfy, and
// the server falls back to the unfiltered list when a store does not implement
// it.
type ReportSearchStore interface {
	SearchReports(status, query string, limit, offset int) ([]QuestionReport, int, error)
}

type AuditSearchStore interface {
	SearchAdminAudit(query string, limit, offset int) ([]AdminAuditRecord, int, error)
}

type adminSearchAccess struct {
	db     *sql.DB
	rebind func(string) string
}

func (s *SQLiteStore) adminSearch() *adminSearchAccess {
	return &adminSearchAccess{db: s.db, rebind: rebindQuestion}
}

func (s *PostgresStore) adminSearch() *adminSearchAccess {
	return &adminSearchAccess{db: s.db, rebind: rebindDollars}
}

func (s *SQLiteStore) SearchReports(status, query string, limit, offset int) ([]QuestionReport, int, error) {
	return s.adminSearch().searchReports(status, query, limit, offset)
}

func (s *PostgresStore) SearchReports(status, query string, limit, offset int) ([]QuestionReport, int, error) {
	return s.adminSearch().searchReports(status, query, limit, offset)
}

func (s *SQLiteStore) SearchAdminAudit(query string, limit, offset int) ([]AdminAuditRecord, int, error) {
	return s.adminSearch().searchAudit(query, limit, offset)
}

func (s *PostgresStore) SearchAdminAudit(query string, limit, offset int) ([]AdminAuditRecord, int, error) {
	return s.adminSearch().searchAudit(query, limit, offset)
}

// reportSearchWhere builds the shared predicate. The reporter's display name
// and username are searched through the same subqueries the list renders, so
// typing a name finds the rows the operator can see.
func reportSearchWhere(status, query string) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	if status != "" && status != "all" {
		clauses = append(clauses, "status = ?")
		args = append(args, status)
	}
	if q := strings.ToLower(strings.TrimSpace(query)); q != "" {
		like := "%" + q + "%"
		clauses = append(clauses, `(lower(concept_id) LIKE ?
			OR lower(question) LIKE ?
			OR lower(detail) LIKE ?
			OR lower(reporter_id) LIKE ?
			OR lower(resolution) LIKE ?
			OR lower(COALESCE((SELECT username FROM students WHERE students.id = question_reports.reporter_id), '')) LIKE ?
			OR lower(COALESCE((SELECT name FROM students WHERE students.id = question_reports.reporter_id), '')) LIKE ?)`)
		for i := 0; i < 7; i++ {
			args = append(args, like)
		}
	}
	if len(clauses) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

func (a *adminSearchAccess) searchReports(status, query string, limit, offset int) ([]QuestionReport, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	where, whereArgs := reportSearchWhere(status, query)

	var total int
	if err := a.db.QueryRow(a.rebind("SELECT COUNT(*) FROM question_reports"+where), whereArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reports: %w", err)
	}

	args := append(append([]interface{}{}, whereArgs...), limit, offset)
	rows, err := a.db.Query(a.rebind(reportSelectColumns+where+" ORDER BY created_at DESC LIMIT ? OFFSET ?"), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()
	out := []QuestionReport{}
	for rows.Next() {
		r, err := scanReportRow(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan report: %w", err)
		}
		out = append(out, r)
	}
	return out, total, rows.Err()
}

func (a *adminSearchAccess) searchAudit(query string, limit, offset int) ([]AdminAuditRecord, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var clauses []string
	var args []interface{}
	if q := strings.ToLower(strings.TrimSpace(query)); q != "" {
		like := "%" + q + "%"
		clauses = append(clauses, `(lower(action) LIKE ?
			OR lower(actor_id) LIKE ?
			OR lower(entity_type) LIKE ?
			OR lower(entity_id) LIKE ?
			OR lower(before_json) LIKE ?
			OR lower(after_json) LIKE ?)`)
		for i := 0; i < 6; i++ {
			args = append(args, like)
		}
	}
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	var total int
	if err := a.db.QueryRow(a.rebind("SELECT COUNT(*) FROM admin_audit"+where), args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count audit: %w", err)
	}

	pageArgs := append(append([]interface{}{}, args...), limit, offset)
	rows, err := a.db.Query(a.rebind("SELECT id, actor_id, action, entity_type, entity_id, before_json, after_json, created_at FROM admin_audit"+where+" ORDER BY id DESC LIMIT ? OFFSET ?"), pageArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin audit: %w", err)
	}
	defer rows.Close()
	out := []AdminAuditRecord{}
	for rows.Next() {
		var r AdminAuditRecord
		var createdAt sql.NullString
		if err := rows.Scan(&r.ID, &r.ActorID, &r.Action, &r.EntityType, &r.EntityID,
			&r.BeforeJSON, &r.AfterJSON, &createdAt); err != nil {
			return nil, 0, fmt.Errorf("scan audit: %w", err)
		}
		r.CreatedAt = createdAt.String
		out = append(out, r)
	}
	return out, total, rows.Err()
}
