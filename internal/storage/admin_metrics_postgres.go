package storage

import (
	"fmt"
	"sort"
	"time"
)

// adminPostgresClass mirrors adminSQLiteClass for a joined students row.
const adminPostgresClass = `CASE
	WHEN substr(a.student_id, 1, 6) = 'guest_' THEN 'guest'
	WHEN COALESCE(s.username, '') <> '' OR COALESCE(s.password_hash, '') <> ''
	  OR COALESCE(s.email, '') <> '' OR COALESCE(s.google_id, '') <> ''
	THEN 'registered'
	ELSE 'unclassified'
END`

func (s *PostgresStore) AdminMetrics(since, until time.Time, bucket time.Duration, sections ...string) (*AdminMetrics, error) {
	if !until.After(since) {
		return nil, fmt.Errorf("admin metrics: until must be after since")
	}
	bucketSec := int64(bucket / time.Second)
	if bucketSec <= 0 {
		bucketSec = int64((24 * time.Hour) / time.Second)
	}
	out := &AdminMetrics{
		MasteredByConcept: map[string]int{},
		ReportStatus:      map[string]int{},
	}
	wantsTotals := adminSectionRequested(sections, "totals")
	if adminSectionRequested(sections, "accounts") {
		if err := s.adminPostgresAccounts(out); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "activity") || wantsTotals {
		if err := s.adminPostgresActivity(out, since, until, bucketSec); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "signups") || wantsTotals {
		if err := s.adminPostgresSignups(out, since, until, bucketSec); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "mastery") {
		if err := s.adminPostgresMastery(out, since, until, bucketSec); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "xp") || wantsTotals {
		if err := s.adminPostgresXP(out, since, until, bucketSec); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "quizzes") || adminSectionRequested(sections, "diagnostics") || wantsTotals {
		if err := s.adminPostgresQuizzes(out, since, until, bucketSec); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "reports") || adminSectionRequested(sections, "status") || wantsTotals {
		if err := s.adminPostgresReports(out, since, until, bucketSec); err != nil {
			return nil, err
		}
	}
	if adminSectionRequested(sections, "concepts") {
		if err := s.adminPostgresConcepts(out, since, until); err != nil {
			return nil, err
		}
	}
	starts := adminBucketStarts(since, until, bucket)
	if adminSectionRequested(sections, "activity") {
		out.Activity = adminFillActivity(out.Activity, starts)
	}
	if adminSectionRequested(sections, "signups") {
		out.Signups = adminFillSeries(out.Signups, starts)
	}
	if adminSectionRequested(sections, "mastery") {
		out.Mastery = adminFillSeries(out.Mastery, starts)
	}
	if adminSectionRequested(sections, "xp") {
		out.XP = adminFillXP(out.XP, starts)
	}
	if adminSectionRequested(sections, "reports") {
		out.Reports = adminFillReports(out.Reports, starts)
	}
	if adminSectionRequested(sections, "quizzes") {
		out.Quizzes = adminFillSeries(out.Quizzes, starts)
		out.TimedQuizzes = adminFillSeries(out.TimedQuizzes, starts)
	}
	if adminSectionRequested(sections, "diagnostics") {
		out.DiagnosticCompletions = adminFillSeries(out.DiagnosticCompletions, starts)
	}
	return out, nil
}

func (s *PostgresStore) AdminLiveMetrics(since time.Time) (AdminLiveMetrics, error) {
	var live AdminLiveMetrics
	dateLo, dateHi := adminDatePrefilter(since, time.Now().UTC())
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM attempts
		 WHERE timestamp::timestamptz >= $1 AND timestamp >= $2 AND timestamp < $3`,
		since.UTC(), dateLo, dateHi).Scan(&live.Answers15m); err != nil {
		return live, fmt.Errorf("admin live answers: %w", err)
	}
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM active_sessions WHERE updated_at::timestamptz >= $1`,
		since.UTC()).Scan(&live.ActiveSessions); err != nil {
		return live, fmt.Errorf("admin live sessions: %w", err)
	}
	return live, nil
}

func (s *PostgresStore) adminPostgresAccounts(out *AdminMetrics) error {
	row := s.db.QueryRow(`SELECT COUNT(*),
		COALESCE(SUM(CASE WHEN class = 'registered' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN class = 'guest' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN class = 'unclassified' THEN 1 ELSE 0 END), 0)
		FROM (
			SELECT CASE
				WHEN substr(id, 1, 6) = 'guest_' THEN 'guest'
				WHEN COALESCE(username, '') <> '' OR COALESCE(password_hash, '') <> ''
				  OR COALESCE(email, '') <> '' OR COALESCE(google_id, '') <> ''
				THEN 'registered'
				ELSE 'unclassified'
			END AS class FROM students
		) classified`)
	if err := row.Scan(&out.Accounts.Total, &out.Accounts.Registered,
		&out.Accounts.Guest, &out.Accounts.Unclassified); err != nil {
		return fmt.Errorf("admin accounts: %w", err)
	}
	return nil
}

func (s *PostgresStore) adminPostgresActivity(out *AdminMetrics, since, until time.Time, bucketSec int64) error {
	dateLo, dateHi := adminDatePrefilter(since, until)
	// At 324k rows and a 90-day window that contains every row, the index
	// plan measured 293ms against 218ms for a sequential scan. The planner is
	// deliberately left free to choose; the index is not forced here.
	activity, err := s.db.Query(`SELECT bucket,
		COUNT(DISTINCT CASE WHEN class = 'registered' THEN student_id END),
		COUNT(DISTINCT CASE WHEN class = 'guest' THEN student_id END),
		COUNT(DISTINCT CASE WHEN class = 'unclassified' THEN student_id END),
		COUNT(*), COALESCE(SUM(correct), 0)
		FROM (
			SELECT a.student_id, a.correct,
				(FLOOR(EXTRACT(EPOCH FROM a.timestamp::timestamptz))::bigint / $1) * $1 AS bucket,
				`+adminPostgresClass+` AS class
			FROM attempts a LEFT JOIN students s ON s.id = a.student_id
			WHERE a.timestamp::timestamptz >= $2 AND a.timestamp::timestamptz < $3
			  AND a.timestamp >= $4 AND a.timestamp < $5
		) classified
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UTC(), until.UTC(), dateLo, dateHi)
	if err != nil {
		return fmt.Errorf("admin activity: %w", err)
	}
	defer activity.Close()
	for activity.Next() {
		var p AdminActivityPoint
		if err := activity.Scan(&p.Start, &p.Registered, &p.Guests,
			&p.Unclassified, &p.Answers, &p.Correct); err != nil {
			return fmt.Errorf("admin activity scan: %w", err)
		}
		out.Activity = append(out.Activity, p)
	}
	if err := activity.Err(); err != nil {
		return fmt.Errorf("admin activity rows: %w", err)
	}

	var totals struct {
		all, registered, guests, unclassified, answers, correct int
	}
	if err := s.db.QueryRow(`SELECT
		COUNT(DISTINCT student_id),
		COUNT(DISTINCT CASE WHEN class = 'registered' THEN student_id END),
		COUNT(DISTINCT CASE WHEN class = 'guest' THEN student_id END),
		COUNT(DISTINCT CASE WHEN class = 'unclassified' THEN student_id END),
		COUNT(*), COALESCE(SUM(correct), 0)
		FROM (
			SELECT a.student_id, a.correct, `+adminPostgresClass+` AS class
			FROM attempts a LEFT JOIN students s ON s.id = a.student_id
			WHERE a.timestamp::timestamptz >= $1 AND a.timestamp::timestamptz < $2
			  AND a.timestamp >= $3 AND a.timestamp < $4
		) classified`, since.UTC(), until.UTC(), dateLo, dateHi).Scan(&totals.all, &totals.registered,
		&totals.guests, &totals.unclassified, &totals.answers, &totals.correct); err != nil {
		return fmt.Errorf("admin activity totals: %w", err)
	}
	out.Totals.UniqueLearners = totals.all
	out.Totals.UniqueRegistered = totals.registered
	out.Totals.UniqueGuests = totals.guests
	out.Totals.UniqueUnclassified = totals.unclassified
	out.Totals.Answers = totals.answers
	out.Totals.Correct = totals.correct
	return nil
}

func (s *PostgresStore) adminPostgresSignups(out *AdminMetrics, since, until time.Time, bucketSec int64) error {
	rows, err := s.db.Query(`SELECT bucket, COUNT(*) FROM (
		SELECT (FLOOR(EXTRACT(EPOCH FROM created_at::timestamptz))::bigint / $1) * $1 AS bucket,
			CASE
				WHEN substr(id, 1, 6) = 'guest_' THEN 'guest'
				WHEN COALESCE(username, '') <> '' OR COALESCE(password_hash, '') <> ''
				  OR COALESCE(email, '') <> '' OR COALESCE(google_id, '') <> ''
				THEN 'registered'
				ELSE 'unclassified'
			END AS class
		FROM students
		WHERE created_at::timestamptz >= $2 AND created_at::timestamptz < $3
	) classified WHERE class = 'registered' GROUP BY bucket ORDER BY bucket`,
		bucketSec, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin signups: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p AdminSeriesPoint
		if err := rows.Scan(&p.Start, &p.Value); err != nil {
			return fmt.Errorf("admin signups scan: %w", err)
		}
		out.Signups = append(out.Signups, p)
		out.Totals.NewRegistrations += p.Value
	}
	if err := rows.Err(); err != nil {
		return err
	}
	out.Totals.DiagnosticCohortRegistered = out.Totals.NewRegistrations
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM students
		WHERE diagnostic_completed_at <> ''
		  AND created_at::timestamptz >= $1 AND created_at::timestamptz < $2
		  AND substr(id, 1, 6) <> 'guest_'
		  AND (COALESCE(username, '') <> '' OR COALESCE(password_hash, '') <> ''
		    OR COALESCE(email, '') <> '' OR COALESCE(google_id, '') <> '')`,
		since.UTC(), until.UTC()).Scan(&out.Totals.DiagnosticCohortCompleted); err != nil {
		return fmt.Errorf("admin diagnostic cohort: %w", err)
	}
	return nil
}

func (s *PostgresStore) adminPostgresMastery(out *AdminMetrics, since, until time.Time, bucketSec int64) error {
	rows, err := s.db.Query(`SELECT
		(FLOOR(EXTRACT(EPOCH FROM mastered_at::timestamptz))::bigint / $1) * $1 AS bucket,
		COUNT(*)
		FROM concept_progress
		WHERE status = 'MASTERED' AND mastered_at IS NOT NULL AND mastered_at <> ''
		  AND mastered_at::timestamptz >= $2 AND mastered_at::timestamptz < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin mastery: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p AdminSeriesPoint
		if err := rows.Scan(&p.Start, &p.Value); err != nil {
			return fmt.Errorf("admin mastery scan: %w", err)
		}
		out.Mastery = append(out.Mastery, p)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM concept_progress WHERE status = 'MASTERED'`).
		Scan(&out.MasteredCurrent); err != nil {
		return fmt.Errorf("admin mastered current: %w", err)
	}
	byConcept, err := s.db.Query(`SELECT concept_id, COUNT(*) FROM concept_progress
		WHERE status = 'MASTERED' GROUP BY concept_id`)
	if err != nil {
		return fmt.Errorf("admin mastery by concept: %w", err)
	}
	defer byConcept.Close()
	for byConcept.Next() {
		var id string
		var count int
		if err := byConcept.Scan(&id, &count); err != nil {
			return fmt.Errorf("admin mastery by concept scan: %w", err)
		}
		out.MasteredByConcept[id] = count
	}
	return byConcept.Err()
}

func (s *PostgresStore) adminPostgresXP(out *AdminMetrics, since, until time.Time, bucketSec int64) error {
	rows, err := s.db.Query(`SELECT
		(created_at / 1000000000 / $1) * $1 AS bucket,
		COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN amount < 0 THEN -amount ELSE 0 END), 0)
		FROM xp_events WHERE created_at >= $2 AND created_at < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UnixNano(), until.UnixNano())
	if err != nil {
		return fmt.Errorf("admin xp: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p AdminXPPoint
		if err := rows.Scan(&p.Start, &p.Awarded, &p.Deducted); err != nil {
			return fmt.Errorf("admin xp scan: %w", err)
		}
		out.XP = append(out.XP, p)
		out.Totals.XPAwarded += p.Awarded
		out.Totals.XPDeducted += p.Deducted
	}
	return rows.Err()
}

func (s *PostgresStore) adminPostgresQuizzes(out *AdminMetrics, since, until time.Time, bucketSec int64) error {
	rows, err := s.db.Query(`SELECT
		(FLOOR(EXTRACT(EPOCH FROM completed_at::timestamptz))::bigint / $1) * $1 AS bucket, COUNT(*)
		FROM quiz_completions
		WHERE completed_at::timestamptz >= $2 AND completed_at::timestamptz < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin quizzes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var p AdminSeriesPoint
		if err := rows.Scan(&p.Start, &p.Value); err != nil {
			return fmt.Errorf("admin quizzes scan: %w", err)
		}
		out.Quizzes = append(out.Quizzes, p)
		out.Totals.QuizCompletions += p.Value
	}
	if err := rows.Err(); err != nil {
		return err
	}

	timed, err := s.db.Query(`SELECT
		(completed_at / 1000000000 / $1) * $1 AS bucket, COUNT(*)
		FROM timed_quiz_runs WHERE status = 'completed' AND completed_at > 0
		  AND completed_at >= $2 AND completed_at < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UnixNano(), until.UnixNano())
	if err != nil {
		return fmt.Errorf("admin timed quizzes: %w", err)
	}
	defer timed.Close()
	for timed.Next() {
		var p AdminSeriesPoint
		if err := timed.Scan(&p.Start, &p.Value); err != nil {
			return fmt.Errorf("admin timed quizzes scan: %w", err)
		}
		out.TimedQuizzes = append(out.TimedQuizzes, p)
		out.Totals.TimedQuizCompletions += p.Value
	}
	if err := timed.Err(); err != nil {
		return err
	}

	diag, err := s.db.Query(`SELECT
		(FLOOR(EXTRACT(EPOCH FROM diagnostic_completed_at::timestamptz))::bigint / $1) * $1 AS bucket,
		COUNT(*)
		FROM students
		WHERE diagnostic_completed_at <> ''
		  AND diagnostic_completed_at::timestamptz >= $2
		  AND diagnostic_completed_at::timestamptz < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin diagnostics: %w", err)
	}
	defer diag.Close()
	for diag.Next() {
		var p AdminSeriesPoint
		if err := diag.Scan(&p.Start, &p.Value); err != nil {
			return fmt.Errorf("admin diagnostics scan: %w", err)
		}
		out.DiagnosticCompletions = append(out.DiagnosticCompletions, p)
		out.Totals.DiagnosticCompletions += p.Value
	}
	return diag.Err()
}

func (s *PostgresStore) adminPostgresReports(out *AdminMetrics, since, until time.Time, bucketSec int64) error {
	opened, err := s.db.Query(`SELECT
		(FLOOR(EXTRACT(EPOCH FROM created_at::timestamptz))::bigint / $1) * $1 AS bucket, COUNT(*)
		FROM question_reports
		WHERE created_at::timestamptz >= $2 AND created_at::timestamptz < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin reports opened: %w", err)
	}
	defer opened.Close()
	for opened.Next() {
		var p AdminReportPoint
		if err := opened.Scan(&p.Start, &p.Opened); err != nil {
			return fmt.Errorf("admin reports opened scan: %w", err)
		}
		out.Reports = adminMergeReport(out.Reports, p)
		out.Totals.ReportsOpened += p.Opened
	}
	if err := opened.Err(); err != nil {
		return err
	}
	resolved, err := s.db.Query(`SELECT
		(FLOOR(EXTRACT(EPOCH FROM resolved_at::timestamptz))::bigint / $1) * $1 AS bucket, COUNT(*)
		FROM question_reports
		WHERE resolved_at <> ''
		  AND resolved_at::timestamptz >= $2 AND resolved_at::timestamptz < $3
		GROUP BY bucket ORDER BY bucket`, bucketSec, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin reports resolved: %w", err)
	}
	defer resolved.Close()
	for resolved.Next() {
		var p AdminReportPoint
		if err := resolved.Scan(&p.Start, &p.Resolved); err != nil {
			return fmt.Errorf("admin reports resolved scan: %w", err)
		}
		out.Reports = adminMergeReport(out.Reports, p)
		out.Totals.ReportsResolved += p.Resolved
	}
	if err := resolved.Err(); err != nil {
		return err
	}
	status, err := s.db.Query(`SELECT status, COUNT(*) FROM question_reports GROUP BY status`)
	if err != nil {
		return fmt.Errorf("admin report status: %w", err)
	}
	defer status.Close()
	for status.Next() {
		var name string
		var count int
		if err := status.Scan(&name, &count); err != nil {
			return fmt.Errorf("admin report status scan: %w", err)
		}
		out.ReportStatus[name] = count
	}
	return status.Err()
}

func (s *PostgresStore) adminPostgresConcepts(out *AdminMetrics, since, until time.Time) error {
	dateLo, dateHi := adminDatePrefilter(since, until)
	byConcept := map[string]*AdminConceptMetric{}
	rows, err := s.db.Query(`SELECT concept_id, COUNT(*), COALESCE(SUM(correct), 0)
		FROM attempts
		WHERE timestamp::timestamptz >= $1 AND timestamp::timestamptz < $2
		  AND timestamp >= $3 AND timestamp < $4
		GROUP BY concept_id`, since.UTC(), until.UTC(), dateLo, dateHi)
	if err != nil {
		return fmt.Errorf("admin concept attempts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var m AdminConceptMetric
		if err := rows.Scan(&m.ConceptID, &m.Attempts, &m.Correct); err != nil {
			return fmt.Errorf("admin concept attempts scan: %w", err)
		}
		byConcept[m.ConceptID] = &m
	}
	if err := rows.Err(); err != nil {
		return err
	}
	reports, err := s.db.Query(`SELECT concept_id, COUNT(*) FROM question_reports
		WHERE concept_id <> '' AND created_at::timestamptz >= $1
		  AND created_at::timestamptz < $2
		GROUP BY concept_id`, since.UTC(), until.UTC())
	if err != nil {
		return fmt.Errorf("admin concept reports: %w", err)
	}
	defer reports.Close()
	for reports.Next() {
		var id string
		var count int
		if err := reports.Scan(&id, &count); err != nil {
			return fmt.Errorf("admin concept reports scan: %w", err)
		}
		m := byConcept[id]
		if m == nil {
			m = &AdminConceptMetric{ConceptID: id}
			byConcept[id] = m
		}
		m.Reports = count
	}
	if err := reports.Err(); err != nil {
		return err
	}
	for _, m := range byConcept {
		out.TopConcepts = append(out.TopConcepts, *m)
	}
	sort.Slice(out.TopConcepts, func(i, j int) bool {
		if out.TopConcepts[i].Attempts != out.TopConcepts[j].Attempts {
			return out.TopConcepts[i].Attempts > out.TopConcepts[j].Attempts
		}
		return out.TopConcepts[i].ConceptID < out.TopConcepts[j].ConceptID
	})
	if len(out.TopConcepts) > 20 {
		out.TopConcepts = out.TopConcepts[:20]
	}
	for _, m := range byConcept {
		if m.Reports > 0 {
			out.MostReported = append(out.MostReported, *m)
		}
	}
	sort.Slice(out.MostReported, func(i, j int) bool {
		if out.MostReported[i].Reports != out.MostReported[j].Reports {
			return out.MostReported[i].Reports > out.MostReported[j].Reports
		}
		return out.MostReported[i].ConceptID < out.MostReported[j].ConceptID
	})
	if len(out.MostReported) > 10 {
		out.MostReported = out.MostReported[:10]
	}
	return nil
}
