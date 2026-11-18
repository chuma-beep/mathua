package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// Leaderboard page queries (Postgres). Mirrors leaderboard_page.go statement for
// statement; the eligibility predicate, the tie-breakers and the window rank
// must be identical or the two stores would publish different podiums.

func (s *PostgresStore) LeaderboardPage(view string, start, end, leagueWeek time.Time, limit, offset int) ([]LeaderboardRow, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	if offset < 0 {
		offset = 0
	}
	leagueKey := leagueWeek.UTC().Format("2006-01-02")
	var rows *sql.Rows
	var err error
	var total int
	switch view {
	case "weekly", "monthly":
		if end.IsZero() || !end.After(start) {
			return nil, 0, fmt.Errorf("leaderboard view %q needs a bounded window", view)
		}
		masteredFrom := start.UTC().Format(time.RFC3339)
		masteredTo := end.UTC().Format(time.RFC3339)
		if err := s.db.QueryRow(
			`SELECT COUNT(*) FROM (`+pgFilteredRows+`) AS ranked`,
			masteredFrom, masteredTo, leagueKey, start.UnixNano(), end.UnixNano()).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("%s leaderboard total: %w", view, err)
		}
		rows, err = s.db.Query(pgFilteredRows+`
			ORDER BY xp DESC, s.xp_total DESC, s.created_at ASC, s.id ASC
			LIMIT $6 OFFSET $7`,
			masteredFrom, masteredTo, leagueKey, start.UnixNano(), end.UnixNano(), limit, offset)
	case "all_time":
		if err := s.db.QueryRow(
			`SELECT COUNT(*) FROM students s WHERE s.xp_total > 0 AND ` + rankedStudentWhere).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("all-time leaderboard total: %w", err)
		}
		rows, err = s.db.Query(pgAllTimeRows+`
			ORDER BY xp DESC, s.created_at ASC, s.id ASC
			LIMIT $2 OFFSET $3`, leagueKey, limit, offset)
	default:
		return nil, 0, fmt.Errorf("leaderboard view %q", view)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("leaderboard page: %w", err)
	}
	defer rows.Close()
	out, err := scanLeaderboardRows(rows)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (s *PostgresStore) LeaderboardPosition(studentID, view string, start, end time.Time) (rank, xp int, found bool, err error) {
	var q string
	switch view {
	case "weekly", "monthly":
		if end.IsZero() || !end.After(start) {
			return 0, 0, false, fmt.Errorf("leaderboard view %q needs a bounded window", view)
		}
		q = `SELECT pos, xp FROM (
			SELECT s.id AS sid, SUM(e.amount) AS xp,
			       ROW_NUMBER() OVER (ORDER BY SUM(e.amount) DESC, s.xp_total DESC, s.created_at ASC, s.id ASC) AS pos
			FROM students s
			JOIN xp_events e ON e.student_id = s.id AND e.created_at >= $1 AND e.created_at < $2
			WHERE ` + rankedStudentWhere + `
			GROUP BY s.id
			HAVING SUM(e.amount) > 0
		) AS ranked WHERE sid = $3`
		err = s.db.QueryRow(q, start.UnixNano(), end.UnixNano(), studentID).Scan(&rank, &xp)
	case "all_time":
		q = `SELECT pos, xp FROM (
			SELECT s.id AS sid, s.xp_total AS xp,
			       ROW_NUMBER() OVER (ORDER BY s.xp_total DESC, s.created_at ASC, s.id ASC) AS pos
			FROM students s
			WHERE s.xp_total > 0 AND ` + rankedStudentWhere + `
		) AS ranked WHERE sid = $1`
		err = s.db.QueryRow(q, studentID).Scan(&rank, &xp)
	default:
		return 0, 0, false, fmt.Errorf("leaderboard view %q", view)
	}
	if err == sql.ErrNoRows {
		return 0, 0, false, nil
	}
	if err != nil {
		return 0, 0, false, fmt.Errorf("leaderboard position: %w", err)
	}
	return rank, xp, true, nil
}

func (s *PostgresStore) EarliestLeaderboardEvent() (time.Time, bool, error) {
	var ns sql.NullInt64
	if err := s.db.QueryRow(`SELECT MIN(created_at) FROM xp_events`).Scan(&ns); err != nil {
		return time.Time{}, false, fmt.Errorf("earliest xp event: %w", err)
	}
	if !ns.Valid || ns.Int64 <= 0 {
		return time.Time{}, false, nil
	}
	return time.Unix(0, ns.Int64).UTC(), true, nil
}

const pgFilteredRows = `
	SELECT s.id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
	       s.xp_total,
	       (SELECT COUNT(*) FROM concept_progress cp
	         WHERE cp.student_id = s.id AND cp.status = 'MASTERED'),
	       (SELECT COUNT(*) FROM concept_progress cp
	         WHERE cp.student_id = s.id AND cp.status = 'MASTERED'
	           AND cp.mastered_at >= $1 AND cp.mastered_at < $2),
	       SUM(e.amount) AS xp,
	       COALESCE((SELECT lm.tier FROM league_members lm
	                 WHERE lm.student_id = s.id AND lm.week <= $3
	                 ORDER BY lm.week DESC LIMIT 1), ''),
	       s.xp_total,
	       s.created_at
	FROM students s
	JOIN xp_events e ON e.student_id = s.id AND e.created_at >= $4 AND e.created_at < $5
	WHERE ` + rankedStudentWhere + `
	GROUP BY s.id
	HAVING SUM(e.amount) > 0`

const pgAllTimeRows = `
	SELECT s.id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
	       s.xp_total,
	       (SELECT COUNT(*) FROM concept_progress cp
	         WHERE cp.student_id = s.id AND cp.status = 'MASTERED'),
	       0,
	       s.xp_total AS xp,
	       COALESCE((SELECT lm.tier FROM league_members lm
	                 WHERE lm.student_id = s.id AND lm.week <= $1
	                 ORDER BY lm.week DESC LIMIT 1), ''),
	       s.xp_total,
	       s.created_at
	FROM students s
	WHERE s.xp_total > 0 AND ` + rankedStudentWhere
