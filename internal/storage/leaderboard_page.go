package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Leaderboard page queries.
//
// Three properties are load-bearing and easy to lose in an edit:
//
//  1. Eligibility is part of the query, not a filter applied afterwards. A guest
//     row can never occupy a LIMIT slot, which is what kept the old board full
//     of zero-XP Guests.
//  2. Ordering ends with a stable key. (week XP, lifetime XP, created_at, id) is
//     a total order over distinct rows, so no learner moves between pages
//     because the database returned ties in another order.
//  3. Rank and total come from the same predicate as the rows. Computing them
//     in Go from a page would rank the page, not the board.
//
// A period is a half-open window [start, end). The upper bound is not optional:
// without it a past week also counted every later week, so an archive could
// never agree with the board it claimed to show.

// rankedStudentWhere is the one definition of a ranked account. A guest row is
// created with only id, name, settings and created_at. It may have a generated
// display username, so username is not a credential and must not make a guest
// eligible. A real account has a password, an email or an OAuth identity. An
// adopted guest keeps its id but gains a password, and must stay ranked.
const rankedStudentWhere = `(COALESCE(s.password_hash, '') <> '' OR COALESCE(s.email, '') <> '' OR EXISTS (SELECT 1 FROM identities i WHERE i.student_id = s.id))`

// LeaderboardPage returns one page of a board and the number of eligible rows.
// view is "weekly", "monthly" or "all_time". start and end bound the XP window;
// both are ignored for all_time. leagueWeek is the Monday used to pick the tier
// frame, so a month never borrows next month's tier.
func (s *SQLiteStore) LeaderboardPage(view string, start, end, leagueWeek time.Time, limit, offset int) ([]LeaderboardRow, int, error) {
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
			`SELECT COUNT(*) FROM (`+sqliteFilteredRows+`)`,
			masteredFrom, masteredTo, leagueKey, start.UnixNano(), end.UnixNano()).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("%s leaderboard total: %w", view, err)
		}
		rows, err = s.db.Query(sqliteFilteredRows+`
			ORDER BY xp DESC, s.xp_total DESC, s.created_at ASC, s.id ASC
			LIMIT ? OFFSET ?`,
			masteredFrom, masteredTo, leagueKey, start.UnixNano(), end.UnixNano(), limit, offset)
	case "all_time":
		if err := s.db.QueryRow(
			`SELECT COUNT(*) FROM students s WHERE s.xp_total > 0 AND ` + rankedStudentWhere).Scan(&total); err != nil {
			return nil, 0, fmt.Errorf("all-time leaderboard total: %w", err)
		}
		rows, err = s.db.Query(sqliteAllTimeRows+`
			ORDER BY xp DESC, s.created_at ASC, s.id ASC
			LIMIT ? OFFSET ?`, leagueKey, limit, offset)
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

// LeaderboardPosition returns the learner's rank and score, or found=false when
// they are not eligible for that view. It is a separate query on purpose: the
// page must not grow to include the learner, and the learner must be findable
// while sitting on page forty.
func (s *SQLiteStore) LeaderboardPosition(studentID, view string, start, end time.Time) (rank, xp int, found bool, err error) {
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
			JOIN xp_events e ON e.student_id = s.id AND e.created_at >= ? AND e.created_at < ?
			WHERE ` + rankedStudentWhere + `
			GROUP BY s.id
			HAVING SUM(e.amount) > 0
		) WHERE sid = ?`
		// created_at is Unix nanoseconds; mastered_at is RFC3339 text. The two
		// placeholders take different types, which is why one is converted here.
		err = s.db.QueryRow(q, start.UnixNano(), end.UnixNano(), studentID).Scan(&rank, &xp)
	case "all_time":
		q = `SELECT pos, xp FROM (
			SELECT s.id AS sid, s.xp_total AS xp,
			       ROW_NUMBER() OVER (ORDER BY s.xp_total DESC, s.created_at ASC, s.id ASC) AS pos
			FROM students s
			WHERE s.xp_total > 0 AND ` + rankedStudentWhere + `
		) WHERE sid = ?`
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

// EarliestLeaderboardEvent is the first recorded XP event, the lower bound an
// archive picker may offer. It is one indexed aggregate, not a table scan.
func (s *SQLiteStore) EarliestLeaderboardEvent() (time.Time, bool, error) {
	var ns sql.NullInt64
	if err := s.db.QueryRow(`SELECT MIN(created_at) FROM xp_events`).Scan(&ns); err != nil {
		return time.Time{}, false, fmt.Errorf("earliest xp event: %w", err)
	}
	if !ns.Valid || ns.Int64 <= 0 {
		return time.Time{}, false, nil
	}
	return time.Unix(0, ns.Int64).UTC(), true, nil
}

const sqliteFilteredRows = `
	SELECT s.id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
	       s.xp_total,
	       (SELECT COUNT(*) FROM concept_progress cp
	         WHERE cp.student_id = s.id AND cp.status = 'MASTERED'),
	       (SELECT COUNT(*) FROM concept_progress cp
	         WHERE cp.student_id = s.id AND cp.status = 'MASTERED'
	           AND cp.mastered_at >= ? AND cp.mastered_at < ?),
	       SUM(e.amount) AS xp,
	       COALESCE((SELECT lm.tier FROM league_members lm
	                 WHERE lm.student_id = s.id AND lm.week <= ?
	                 ORDER BY lm.week DESC LIMIT 1), ''),
	       s.xp_total,
	       s.created_at
	FROM students s
	JOIN xp_events e ON e.student_id = s.id AND e.created_at >= ? AND e.created_at < ?
	WHERE ` + rankedStudentWhere + `
	GROUP BY s.id
	HAVING SUM(e.amount) > 0`

const sqliteAllTimeRows = `
	SELECT s.id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
	       s.xp_total,
	       (SELECT COUNT(*) FROM concept_progress cp
	         WHERE cp.student_id = s.id AND cp.status = 'MASTERED'),
	       0,
	       s.xp_total AS xp,
	       COALESCE((SELECT lm.tier FROM league_members lm
	                 WHERE lm.student_id = s.id AND lm.week <= ?
	                 ORDER BY lm.week DESC LIMIT 1), ''),
	       s.xp_total,
	       s.created_at
	FROM students s
	WHERE s.xp_total > 0 AND ` + rankedStudentWhere

func scanLeaderboardRows(rows *sql.Rows) ([]LeaderboardRow, error) {
	var out []LeaderboardRow
	for rows.Next() {
		var r LeaderboardRow
		if err := rows.Scan(
			&r.StudentID, &r.Name, &r.Username, &r.AvatarURL, &r.AvatarSettings,
			&r.XPTotal, &r.TotalMastered, &r.WeeklyMastered, &r.XP, &r.League,
			&r.LeagueXP, &r.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan leaderboard row: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// IsGuestID reports the legacy guest_ credential prefix. Ranking eligibility
// does not use this prefix: guests can have UUID ids, so rankedStudentWhere
// checks for real credentials instead.
func IsGuestID(id string) bool {
	return strings.HasPrefix(id, "guest_")
}
