package storage

import (
	"database/sql"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClassifyStudent(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		username string
		password string
		email    string
		google   string
		want     StudentClass
	}{
		{name: "guest prefix", id: "guest_abc", want: StudentClassGuest},
		{name: "guestX is not a guest", id: "guestX", want: StudentClassUnclassified},
		{name: "uuid without credentials", id: "3f0f6a52-0000-4000-8000-000000000000", want: StudentClassUnclassified},
		{name: "username", id: "3f0f6a52-0000-4000-8000-000000000001", username: "ada", want: StudentClassRegistered},
		{name: "password", id: "3f0f6a52-0000-4000-8000-000000000002", password: "hash", want: StudentClassRegistered},
		{name: "email", id: "3f0f6a52-0000-4000-8000-000000000003", email: "ada@example.com", want: StudentClassRegistered},
		{name: "google identity", id: "3f0f6a52-0000-4000-8000-000000000004", google: "google-sub", want: StudentClassRegistered},
		{name: "whitespace only", id: "3f0f6a52-0000-4000-8000-000000000005", username: "  ", password: "\t", email: " ", want: StudentClassUnclassified},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ClassifyStudent(tc.id, tc.username, tc.password, tc.email, tc.google); got != tc.want {
				t.Fatalf("ClassifyStudent = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAdminMetricsCore(t *testing.T) {
	store := newTestStore(t)
	now := time.Now().UTC()
	since := now.Add(-2 * time.Hour)

	registered, err := store.CreateUser("Ada", "ada_metrics", "hash")
	if err != nil {
		t.Fatalf("create registered: %v", err)
	}
	guest, err := store.ClaimGuestStudent("guest_metrics", "Guest")
	if err != nil {
		t.Fatalf("create guest: %v", err)
	}
	unclassified, err := store.CreateStudent("uuid guest")
	if err != nil {
		t.Fatalf("create unclassified: %v", err)
	}
	if _, err := store.CreateUser("Grace", "grace_metrics", "hash"); err != nil {
		t.Fatalf("create second registered: %v", err)
	}

	record := func(st *Student, correct bool, key string) {
		t.Helper()
		sid := "sess-metrics-" + key
		if err := store.EnsureSession(sid, st.ID); err != nil {
			t.Fatalf("session %s: %v", key, err)
		}
		if _, err := store.RecordAttempt(AttemptEntry{
			SessionID:      sid,
			StudentID:      st.ID,
			ConceptID:      "admin.metric.concept",
			Answer:         "1",
			Expected:       "1",
			Correct:        correct,
			ElapsedSeconds: 4,
			Timestamp:      now.Add(-time.Hour),
			Question:       "q",
			Source:         "practice",
		}); err != nil {
			t.Fatalf("attempt %s: %v", key, err)
		}
	}
	record(registered, true, "registered")
	record(guest, true, "guest")
	record(unclassified, false, "unclassified")

	if _, err := store.CreateReport(QuestionReport{
		ReporterID: registered.ID,
		ConceptID:  "admin.metric.concept",
		Kind:       "question",
		Question:   "q",
		Expected:   "1",
		Source:     "practice",
	}); err != nil {
		t.Fatalf("create report: %v", err)
	}
	if err := store.SetDiagnosticCompleted(registered.ID); err != nil {
		t.Fatalf("diagnostic completed: %v", err)
	}
	if err := store.RecordQuizCompletion(registered.ID, 60); err != nil {
		t.Fatalf("quiz completion: %v", err)
	}
	if err := store.AwardXP(registered.ID, 100, "test", "metrics:xp:award"); err != nil {
		t.Fatalf("award xp: %v", err)
	}
	if err := store.AwardXP(registered.ID, -25, "test", "metrics:xp:deduct"); err != nil {
		t.Fatalf("deduct xp: %v", err)
	}
	tqStore, ok := store.(TimedQuizStore)
	if !ok {
		t.Fatalf("store %T does not implement TimedQuizStore", store)
	}
	run := TimedQuizRun{ID: "tq-metrics", StudentID: registered.ID, StartedAt: now,
		DeadlineAt: now.Add(10 * time.Minute), Status: "active", CreatedAt: now}
	if err := tqStore.CreateTimedQuizRun(run, nil); err != nil {
		t.Fatalf("timed quiz run: %v", err)
	}
	if _, err := tqStore.FinalizeTimedQuizRun(run.ID, registered.ID, now); err != nil {
		t.Fatalf("finalize timed quiz: %v", err)
	}

	metricsStore, ok := store.(AdminMetricsStore)
	if !ok {
		t.Fatalf("store %T does not implement AdminMetricsStore", store)
	}
	m, err := metricsStore.AdminMetrics(since, now.Add(time.Minute), time.Hour)
	if err != nil {
		t.Fatalf("admin metrics: %v", err)
	}
	if m.Accounts.Total != 4 || m.Accounts.Registered != 2 || m.Accounts.Guest != 1 || m.Accounts.Unclassified != 1 {
		t.Fatalf("accounts = %+v, want 2 registered / 1 guest / 1 unclassified", m.Accounts)
	}
	if m.Totals.UniqueLearners != 3 || m.Totals.UniqueRegistered != 1 ||
		m.Totals.UniqueGuests != 1 || m.Totals.UniqueUnclassified != 1 {
		t.Fatalf("unique totals = %+v, want 1 of each class", m.Totals)
	}
	if m.Totals.Answers != 3 || m.Totals.Correct != 2 {
		t.Fatalf("answers = %d correct = %d, want 3/2", m.Totals.Answers, m.Totals.Correct)
	}
	if m.Totals.NewRegistrations != 2 {
		t.Fatalf("new registrations = %d, want 2", m.Totals.NewRegistrations)
	}
	if m.Totals.QuizCompletions != 1 {
		t.Fatalf("quiz completions = %d, want 1", m.Totals.QuizCompletions)
	}
	if m.Totals.TimedQuizCompletions != 1 {
		t.Fatalf("timed quiz completions = %d, want 1", m.Totals.TimedQuizCompletions)
	}
	if m.Totals.XPAwarded != 100 || m.Totals.XPDeducted != 25 {
		t.Fatalf("xp = %d/%d, want 100/25", m.Totals.XPAwarded, m.Totals.XPDeducted)
	}
	if m.Totals.DiagnosticCompletions != 1 {
		t.Fatalf("diagnostic completions = %d, want 1", m.Totals.DiagnosticCompletions)
	}
	if m.Totals.DiagnosticCohortRegistered != 2 || m.Totals.DiagnosticCohortCompleted != 1 {
		t.Fatalf("diagnostic cohort = %d/%d, want 1 completed of 2 registered",
			m.Totals.DiagnosticCohortCompleted, m.Totals.DiagnosticCohortRegistered)
	}
	if m.Totals.ReportsOpened != 1 || m.ReportStatus["open"] != 1 {
		t.Fatalf("reports = %+v status = %+v, want one open", m.Totals, m.ReportStatus)
	}
	if len(m.Activity) < 3 {
		t.Fatalf("activity buckets = %d, want at least 3", len(m.Activity))
	}
	var answers, activityUnclassified int
	for _, p := range m.Activity {
		answers += p.Answers
		activityUnclassified += p.Unclassified
	}
	if answers != 3 {
		t.Fatalf("activity answers = %d, want 3", answers)
	}
	if activityUnclassified != 1 {
		t.Fatalf("activity unclassified = %d, want 1 separate segment", activityUnclassified)
	}
	if len(m.TopConcepts) != 1 || m.TopConcepts[0].ConceptID != "admin.metric.concept" ||
		m.TopConcepts[0].Attempts != 3 || m.TopConcepts[0].Correct != 2 || m.TopConcepts[0].Reports != 1 {
		t.Fatalf("top concepts = %+v", m.TopConcepts)
	}
	if len(m.MostReported) != 1 || m.MostReported[0].Reports != 1 {
		t.Fatalf("most reported = %+v", m.MostReported)
	}
}

// TestAdminMetricsTimestampBoundaries covers the two storage shapes the
// timestamp filter must survive: RFC3339 in SQLite and now()::text style in
// Postgres, plus a fractional second and a space separator. The window is in
// 2099 so a shared Postgres database cannot contribute rows.
func TestAdminMetricsTimestampBoundaries(t *testing.T) {
	t.Run("sqlite", func(t *testing.T) {
		testAdminMetricsTimestampBoundaries(t, newSQLiteStore(t))
	})
	t.Run("postgres", func(t *testing.T) {
		if os.Getenv("TEST_POSTGRES_DSN") == "" {
			t.Skip("TEST_POSTGRES_DSN not set; skipping Postgres boundary check")
		}
		pg := newPostgresStore(t)
		truncatePostgres(t, pg)
		testAdminMetricsTimestampBoundaries(t, pg)
	})
}

func testAdminMetricsTimestampBoundaries(t *testing.T, store Repository) {
	t.Helper()
	st, err := store.CreateStudent("boundary learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := store.EnsureSession("sess-boundary-2099", st.ID); err != nil {
		t.Fatalf("ensure session: %v", err)
	}
	// Every value is a 2099-01-01 UTC instant expressed in a different shape.
	// Positive and negative offsets move the local calendar day, fractional
	// seconds sit on both sides of a boundary, and the two midnight rows test
	// the half-open interval exactly.
	cases := []struct {
		timestamp string
		inside    bool
	}{
		{"2099-01-01T00:00:00Z", true},
		{"2099-01-01T02:00:00+02:00", true},
		{"2098-12-31T22:00:00-02:00", true},
		{"2099-01-01T12:00:00.123456Z", true},
		{"2099-01-01T12:00:00.987654+05:30", true},
		{"2099-01-01T23:59:59Z", true},
		{"2099-01-01T23:59:59.999999-00:00", true},
		{"2099-01-02T01:00:00+02:00", true},
		{"2099-01-02T00:00:00Z", false},
		{"2099-01-01T23:30:00-02:00", false},
		{"2098-12-31T23:59:59+01:00", false},
		{"2099-01-02T02:00:00+02:00", false},
	}
	want := 0
	for i, tc := range cases {
		if tc.inside {
			want++
		}
		query := `INSERT INTO attempts
			(session_id, student_id, concept_id, answer, expected, correct,
			 elapsed_seconds, timestamp, question, source, explanation)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		var err error
		switch s := store.(type) {
		case *SQLiteStore:
			_, err = s.db.Exec(query, "sess-boundary-2099", st.ID, "boundary.concept", "1", "1", 1, 1.0, tc.timestamp, "q", "practice", "")
		case *PostgresStore:
			_, err = s.db.Exec(rebindDollars(query), "sess-boundary-2099", st.ID, "boundary.concept", "1", "1", 1, 1.0, tc.timestamp, "q", "practice", "")
		default:
			t.Fatalf("unexpected store %T", store)
		}
		if err != nil {
			t.Fatalf("insert attempt %d (%s): %v", i, tc.timestamp, err)
		}
	}

	metricsStore, ok := store.(AdminMetricsStore)
	if !ok {
		t.Fatalf("store %T does not implement AdminMetricsStore", store)
	}
	since := time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 1, 2, 0, 0, 0, 0, time.UTC)
	m, err := metricsStore.AdminMetrics(since, until, 24*time.Hour)
	if err != nil {
		t.Fatalf("admin metrics: %v", err)
	}
	if m.Totals.Answers != want {
		t.Fatalf("answers = %d, want %d inside the half-open window", m.Totals.Answers, want)
	}
	if len(m.Activity) != 1 {
		t.Fatalf("activity buckets = %d, want 1", len(m.Activity))
	}
	if m.Activity[0].Answers != want {
		t.Fatalf("activity answers = %d, want %d", m.Activity[0].Answers, want)
	}
}

// TestAdminMetricsParity runs one identical fixture through both stores and
// compares the complete projection. It is the check that catches a SQLite-only
// date expression or a Postgres cast that rounds a fractional second.
func TestAdminMetricsParity(t *testing.T) {
	if os.Getenv("TEST_POSTGRES_DSN") == "" {
		t.Skip("TEST_POSTGRES_DSN not set; skipping Postgres parity")
	}
	sqlite := newSQLiteStore(t)
	pg := newPostgresStore(t)
	truncatePostgres(t, pg)

	now := time.Now().UTC().Truncate(time.Second)
	since := now.Add(-2 * time.Hour)
	until := now.Add(time.Minute)

	seed := func(store Repository) {
		t.Helper()
		ada, err := store.CreateUser("Ada", "ada_parity", "hash")
		if err != nil {
			t.Fatalf("create Ada: %v", err)
		}
		if _, err := store.CreateUser("Grace", "grace_parity", "hash"); err != nil {
			t.Fatalf("create Grace: %v", err)
		}
		guest, err := store.ClaimGuestStudent("guest_parity", "Guest")
		if err != nil {
			t.Fatalf("create guest: %v", err)
		}
		unclassified, err := store.CreateStudent("uuid parity")
		if err != nil {
			t.Fatalf("create unclassified: %v", err)
		}
		if err := store.EnsureSession("sess-parity", ada.ID); err != nil {
			t.Fatalf("ensure session: %v", err)
		}
		for i, st := range []*Student{ada, guest, unclassified} {
			if _, err := store.RecordAttempt(AttemptEntry{
				SessionID: "sess-parity", StudentID: st.ID, ConceptID: "parity.concept",
				Answer: "1", Expected: "1", Correct: i != 2, ElapsedSeconds: 4,
				Timestamp: now.Add(-time.Hour), Question: "q", Source: "practice",
			}); err != nil {
				t.Fatalf("attempt %d: %v", i, err)
			}
		}
		if _, err := store.CreateReport(QuestionReport{
			ReporterID: ada.ID, ConceptID: "parity.concept", Kind: "question",
			Question: "q", Expected: "1", Source: "practice",
		}); err != nil {
			t.Fatalf("create report: %v", err)
		}
		if err := store.SetDiagnosticCompleted(ada.ID); err != nil {
			t.Fatalf("diagnostic: %v", err)
		}
		if err := store.RecordQuizCompletion(ada.ID, 60); err != nil {
			t.Fatalf("quiz: %v", err)
		}
		if err := store.AwardXP(ada.ID, 100, "test", "parity:xp:award"); err != nil {
			t.Fatalf("award xp: %v", err)
		}
		if err := store.AwardXP(ada.ID, -25, "test", "parity:xp:deduct"); err != nil {
			t.Fatalf("deduct xp: %v", err)
		}
		tq, ok := store.(TimedQuizStore)
		if !ok {
			t.Fatalf("store %T lacks TimedQuizStore", store)
		}
		run := TimedQuizRun{ID: "tq-parity", StudentID: ada.ID, StartedAt: now,
			DeadlineAt: now.Add(10 * time.Minute), Status: "active", CreatedAt: now}
		if err := tq.CreateTimedQuizRun(run, nil); err != nil {
			t.Fatalf("timed quiz: %v", err)
		}
		if _, err := tq.FinalizeTimedQuizRun(run.ID, ada.ID, now); err != nil {
			t.Fatalf("finalize timed quiz: %v", err)
		}
		insertMastery := func(db interface {
			Exec(string, ...interface{}) (sql.Result, error)
		}, query string) {
			if _, err := db.Exec(query, ada.ID, "parity.mastered", now.Add(-time.Hour).Format(time.RFC3339)); err != nil {
				t.Fatalf("insert mastery: %v", err)
			}
		}
		switch s := store.(type) {
		case *SQLiteStore:
			insertMastery(s.db, `INSERT INTO concept_progress (student_id, concept_id, status, mastered_at) VALUES (?, ?, 'MASTERED', ?)`)
		case *PostgresStore:
			insertMastery(s.db, `INSERT INTO concept_progress (student_id, concept_id, status, mastered_at) VALUES ($1, $2, 'MASTERED', $3)`)
		}
	}

	metrics := func(store Repository) *AdminMetrics {
		t.Helper()
		ms, ok := store.(AdminMetricsStore)
		if !ok {
			t.Fatalf("store %T lacks AdminMetricsStore", store)
		}
		m, err := ms.AdminMetrics(since, until, time.Hour)
		if err != nil {
			t.Fatalf("admin metrics: %v", err)
		}
		return m
	}

	seed(sqlite)
	seed(pg)
	if got, want := metrics(sqlite), metrics(pg); !reflect.DeepEqual(got, want) {
		t.Fatalf("SQLite and Postgres disagree:\nSQLite: %+v\nPostgres: %+v", got, want)
	}
}

// TestSQLiteAttemptsQueryUsesTimestampIndex keeps the planner constraint
// honest. If a future schema or query change removes the index or makes the
// prefilter unusable, the 90-day path silently becomes a full scan again.
func TestSQLiteAttemptsQueryUsesTimestampIndex(t *testing.T) {
	store := newSQLiteStore(t)
	rows, err := store.db.Query(`EXPLAIN QUERY PLAN
		SELECT COUNT(DISTINCT a.student_id)
		FROM attempts a INDEXED BY idx_attempts_timestamp
		LEFT JOIN students s ON s.id = a.student_id
		WHERE datetime(a.timestamp) >= datetime(?)
		  AND datetime(a.timestamp) < datetime(?)
		  AND a.timestamp >= ? AND a.timestamp < ?`,
		"2026-08-18T00:00:00Z", "2026-11-16T00:00:00Z", "2026-08-18", "2026-11-16")
	if err != nil {
		t.Fatalf("explain attempts query: %v", err)
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var id, parent, notUsed int
		var detail string
		if err := rows.Scan(&id, &parent, &notUsed, &detail); err != nil {
			t.Fatalf("scan plan: %v", err)
		}
		lines = append(lines, detail)
	}
	plan := strings.Join(lines, "\n")
	if err := rows.Err(); err != nil {
		t.Fatalf("plan rows: %v", err)
	}
	if !strings.Contains(plan, "idx_attempts_timestamp") || !strings.Contains(plan, "SEARCH") {
		t.Fatalf("attempts query does not use the timestamp index:\n%s", plan)
	}
}
