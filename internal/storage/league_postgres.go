package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Weekly leagues (Postgres). Mirrors league_sqlite.go statement for statement so
// the two stores cannot disagree about who moves.

func (s *PostgresStore) ResolveLeagueWeek(week string, prevFrom, prevTo time.Time) error {
	start, err := WeekFromString(week)
	if err != nil {
		return err
	}
	prev := start.AddDate(0, 0, -7).Format("2006-01-02")

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("resolve league week: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		`INSERT INTO league_weeks (week, closed_at) VALUES ($1, $2)
		 ON CONFLICT (week) DO NOTHING`,
		week, time.Now().UTC().UnixNano())
	if err != nil {
		return fmt.Errorf("claim league week: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return nil // another request closed it first
	}

	rows, err := tx.Query(
		`SELECT m.tier, m.league_no, m.student_id, s.name,
		        COALESCE((SELECT SUM(e.amount) FROM xp_events e
		                  WHERE e.student_id = m.student_id
		                    AND e.created_at >= $1 AND e.created_at < $2), 0) AS xp
		 FROM league_members m
		 JOIN students s ON s.id = m.student_id
		 WHERE m.week = $3
		 ORDER BY m.tier, m.league_no, xp DESC, s.name ASC, m.student_id ASC`,
		prevFrom.UnixNano(), prevTo.UnixNano(), prev)
	if err != nil {
		return fmt.Errorf("read previous leagues: %w", err)
	}
	var flat []leagueCohortRow
	for rows.Next() {
		var r leagueCohortRow
		if err := rows.Scan(&r.Tier, &r.LeagueNo, &r.StudentID, &r.Name, &r.XP); err != nil {
			rows.Close()
			return fmt.Errorf("scan previous league: %w", err)
		}
		flat = append(flat, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read previous leagues: %w", err)
	}

	if len(flat) == 0 {
		if err := s.bootstrapLeagues(tx, week); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit league seed: %w", err)
		}
		return nil
	}

	cohorts := groupCohorts(flat)
	tierOf, moveOf := planWeek(cohorts)

	for _, c := range cohorts {
		for i, m := range c.Members {
			if _, err := tx.Exec(
				`UPDATE league_members SET xp = $1, rank = $2 WHERE week = $3 AND student_id = $4`,
				m.XP, i+1, prev, m.StudentID); err != nil {
				return fmt.Errorf("finalise league member: %w", err)
			}
		}
	}

	byTier := make(map[string][]leagueMember)
	for _, c := range cohorts {
		for _, m := range c.Members {
			tier := tierOf[m.StudentID]
			m.Movement = moveOf[m.StudentID]
			byTier[tier] = append(byTier[tier], m)
		}
	}
	var next []leagueCohort
	for _, tier := range LeagueTiers {
		next = append(next, cutCohorts(tier, byTier[tier])...)
	}
	if err := insertCohorts(tx, week, next, rebindDollars); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit league week: %w", err)
	}
	return nil
}

func (s *PostgresStore) bootstrapLeagues(tx *sql.Tx, week string) error {
	rows, err := tx.Query(
		`SELECT id, name, xp_total
		 FROM students
		 WHERE xp_total > 0 AND id NOT LIKE 'guest_%'
		 ORDER BY xp_total DESC, name ASC, id ASC`)
	if err != nil {
		return fmt.Errorf("seed leagues: %w", err)
	}
	var members []leagueMember
	for rows.Next() {
		var id, name string
		var xp int
		if err := rows.Scan(&id, &name, &xp); err != nil {
			rows.Close()
			return fmt.Errorf("scan seed learner: %w", err)
		}
		members = append(members, leagueMember{StudentID: id, Name: name, XP: xp})
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("seed leagues: %w", err)
	}
	return insertCohorts(tx, week, cutCohorts(LeagueTiers[0], members), rebindDollars)
}

func (s *PostgresStore) EnsureLeagueMembership(studentID, week string) error {
	if studentID == "" {
		return nil
	}
	var exists int
	err := s.db.QueryRow(
		`SELECT 1 FROM league_members WHERE week = $1 AND student_id = $2`,
		week, studentID).Scan(&exists)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check league membership: %w", err)
	}

	var xpTotal, ranked int
	if err := s.db.QueryRow(`SELECT s.xp_total, CASE WHEN `+rankedStudentWhere+` THEN 1 ELSE 0 END
		FROM students s WHERE s.id = $1`, studentID).
		Scan(&xpTotal, &ranked); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("check league eligibility: %w", err)
	}
	// An adopted guest keeps its guest_ id but has credentials, so the ranked
	// predicate, not the id prefix, decides eligibility.
	if xpTotal <= 0 || ranked == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("join league: %w", err)
	}
	defer tx.Rollback()

	if err := tx.QueryRow(
		`SELECT 1 FROM league_members WHERE week = $1 AND student_id = $2`,
		week, studentID).Scan(&exists); err == nil {
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("recheck league membership: %w", err)
	}

	no := 0
	err = tx.QueryRow(
		`SELECT league_no FROM league_members
		 WHERE week = $1 AND tier = $2
		 GROUP BY league_no HAVING COUNT(*) < $3
		 ORDER BY league_no LIMIT 1`,
		week, LeagueTiers[0], LeagueSize).Scan(&no)
	if errors.Is(err, sql.ErrNoRows) {
		if err := tx.QueryRow(
			`SELECT COALESCE(MAX(league_no), 0) + 1 FROM league_members
			 WHERE week = $1 AND tier = $2`,
			week, LeagueTiers[0]).Scan(&no); err != nil {
			return fmt.Errorf("open a league: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("find a league with room: %w", err)
	}

	if _, err := tx.Exec(
		`INSERT INTO league_members (week, student_id, tier, league_no, xp, rank, movement, from_tier)
		 VALUES ($1, $2, $3, $4, 0, 0, 'new', '')
		 ON CONFLICT (week, student_id) DO NOTHING`,
		week, studentID, LeagueTiers[0], no); err != nil {
		return fmt.Errorf("insert league member: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit league membership: %w", err)
	}
	return nil
}

func (s *PostgresStore) LeagueMembershipOf(studentID, week string) (string, int, bool, error) {
	var tier string
	var no int
	err := s.db.QueryRow(
		`SELECT tier, league_no FROM league_members WHERE week = $1 AND student_id = $2`,
		week, studentID).Scan(&tier, &no)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, false, nil
	}
	if err != nil {
		return "", 0, false, fmt.Errorf("read league membership: %w", err)
	}
	return validTier(tier), no, true, nil
}

func (s *PostgresStore) LeagueRoster(week, tier string, leagueNo int, from, to time.Time) ([]LeagueRow, error) {
	rows, err := s.db.Query(
		`SELECT m.student_id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings,
		        COALESCE((SELECT SUM(e.amount) FROM xp_events e
		                  WHERE e.student_id = m.student_id
		                    AND e.created_at >= $1 AND e.created_at < $2), 0) AS xp
		 FROM league_members m
		 JOIN students s ON s.id = m.student_id
		 WHERE m.week = $3 AND m.tier = $4 AND m.league_no = $5
		 ORDER BY xp DESC, s.name ASC, m.student_id ASC`,
		from.UnixNano(), to.UnixNano(), week, validTier(tier), leagueNo)
	if err != nil {
		return nil, fmt.Errorf("read league roster: %w", err)
	}
	defer rows.Close()
	var out []LeagueRow
	for rows.Next() {
		var r LeagueRow
		if err := rows.Scan(&r.StudentID, &r.Name, &r.Username, &r.AvatarURL, &r.Settings, &r.XP); err != nil {
			return nil, fmt.Errorf("scan league roster: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read league roster: %w", err)
	}
	rankInCohort(out)
	return out, nil
}

func (s *PostgresStore) LeagueRosterFinal(week, tier string, leagueNo int) ([]LeagueRow, error) {
	rows, err := s.db.Query(
		`SELECT m.student_id, s.name, COALESCE(s.username, ''), s.avatar_url, s.settings, m.xp
		 FROM league_members m
		 JOIN students s ON s.id = m.student_id
		 WHERE m.week = $1 AND m.tier = $2 AND m.league_no = $3
		 ORDER BY m.xp DESC, s.name ASC, m.student_id ASC`,
		week, validTier(tier), leagueNo)
	if err != nil {
		return nil, fmt.Errorf("read final league roster: %w", err)
	}
	defer rows.Close()
	var out []LeagueRow
	for rows.Next() {
		var r LeagueRow
		if err := rows.Scan(&r.StudentID, &r.Name, &r.Username, &r.AvatarURL, &r.Settings, &r.XP); err != nil {
			return nil, fmt.Errorf("scan final league roster: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read final league roster: %w", err)
	}
	rankInCohort(out)
	return out, nil
}

func (s *PostgresStore) LeagueOutcome(studentID, week string) (movement, fromTier, toTier string, found bool, err error) {
	err = s.db.QueryRow(
		`SELECT movement, from_tier, tier FROM league_members
		 WHERE week = $1 AND student_id = $2`, week, studentID).
		Scan(&movement, &fromTier, &toTier)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", "", false, nil
	}
	if err != nil {
		return "", "", "", false, fmt.Errorf("read league outcome: %w", err)
	}
	if fromTier != "" {
		fromTier = validTier(fromTier)
	}
	if toTier != "" {
		toTier = validTier(toTier)
	}
	return movement, fromTier, toTier, true, nil
}

func (s *PostgresStore) EarliestLeagueWeek() (string, bool, error) {
	var week sql.NullString
	if err := s.db.QueryRow(`SELECT MIN(week) FROM league_members`).Scan(&week); err != nil {
		return "", false, fmt.Errorf("earliest league week: %w", err)
	}
	if !week.Valid || week.String == "" {
		return "", false, nil
	}
	return week.String, true, nil
}

// LeagueCohorts lists the week's leagues in pyramid order, exactly as the
// SQLite twin does.
func (s *PostgresStore) LeagueCohorts(week string) ([]LeagueCohort, error) {
	rows, err := s.db.Query(
		`SELECT tier, league_no, COUNT(*) FROM league_members
		 WHERE week = $1
		 GROUP BY tier, league_no`, week)
	if err != nil {
		return nil, fmt.Errorf("read league cohorts: %w", err)
	}
	defer rows.Close()
	var out []LeagueCohort
	for rows.Next() {
		var c LeagueCohort
		if err := rows.Scan(&c.Tier, &c.LeagueNo, &c.Size); err != nil {
			return nil, fmt.Errorf("scan league cohort: %w", err)
		}
		c.Tier = validTier(c.Tier)
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read league cohorts: %w", err)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Tier != out[j].Tier {
			return tierIndex(out[i].Tier) < tierIndex(out[j].Tier)
		}
		return out[i].LeagueNo < out[j].LeagueNo
	})
	return out, nil
}
