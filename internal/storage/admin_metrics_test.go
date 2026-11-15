package storage

import (
	"os"
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
	if m.Accounts.Total != 3 || m.Accounts.Registered != 1 || m.Accounts.Guest != 1 || m.Accounts.Unclassified != 1 {
		t.Fatalf("accounts = %+v, want 1 registered / 1 guest / 1 unclassified", m.Accounts)
	}
	if m.Totals.UniqueLearners != 3 || m.Totals.UniqueRegistered != 1 ||
		m.Totals.UniqueGuests != 1 || m.Totals.UniqueUnclassified != 1 {
		t.Fatalf("unique totals = %+v, want 1 of each class", m.Totals)
	}
	if m.Totals.Answers != 3 || m.Totals.Correct != 2 {
		t.Fatalf("answers = %d correct = %d, want 3/2", m.Totals.Answers, m.Totals.Correct)
	}
	if m.Totals.NewRegistrations != 1 {
		t.Fatalf("new registrations = %d, want 1", m.Totals.NewRegistrations)
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
	if m.Totals.RegisteredDiagnosticCohort != 1 {
		t.Fatalf("registered diagnostic cohort = %d, want 1", m.Totals.RegisteredDiagnosticCohort)
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
		testAdminMetricsTimestampBoundaries(t, newPostgresStore(t))
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
	timestamps := []string{
		"2099-01-01T00:00:00Z",
		"2099-01-01 23:59:59+00:00",
		"2099-01-01T12:00:00.123456Z",
		"2099-01-02T00:00:00Z", // outside the half-open window
	}
	for i, ts := range timestamps {
		query := `INSERT INTO attempts
			(session_id, student_id, concept_id, answer, expected, correct,
			 elapsed_seconds, timestamp, question, source, explanation)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		var err error
		switch s := store.(type) {
		case *SQLiteStore:
			_, err = s.db.Exec(query, "sess-boundary-2099", st.ID, "boundary.concept", "1", "1", 1, 1.0, ts, "q", "practice", "")
		case *PostgresStore:
			_, err = s.db.Exec(rebindDollars(query), "sess-boundary-2099", st.ID, "boundary.concept", "1", "1", 1, 1.0, ts, "q", "practice", "")
		default:
			t.Fatalf("unexpected store %T", store)
		}
		if err != nil {
			t.Fatalf("insert attempt %d: %v", i, err)
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
	if m.Totals.Answers != 3 {
		t.Fatalf("answers = %d, want 3 inside the half-open window", m.Totals.Answers)
	}
	if len(m.Activity) != 1 {
		t.Fatalf("activity buckets = %d, want 1", len(m.Activity))
	}
	if m.Activity[0].Answers != 3 {
		t.Fatalf("activity answers = %d, want 3", m.Activity[0].Answers)
	}
}
