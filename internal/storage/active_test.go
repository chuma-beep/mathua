package storage

import (
	"database/sql"
	"testing"
	"time"
)

// TestActiveLearnerCounts_SplitsGuestsAndWindows pins the two numbers the admin
// overview shows. activeNow is a live session, activeToday is an answered
// question, and neither may count the same learner twice across sessions or
// count a session older than the window.
func TestActiveLearnerCounts_SplitsGuestsAndWindows(t *testing.T) {
	store := newTestStore(t)
	// A fixed instant at midday, not time.Now().
	//
	// The test records an attempt an hour before "now" and expects it to count as
	// active today, which is only true while now is at least an hour into the UTC day.
	// Running the suite between 00:00 and 01:00 UTC put that attempt on the previous
	// day and the assertion failed with "active today = 0 accounts", on a code path
	// nobody had touched. Every timestamp the test writes is already explicit, so
	// pinning the instant costs nothing and makes the case independent of the clock.
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	ada, err := store.CreateStudent("ada")
	if err != nil {
		t.Fatalf("create ada: %v", err)
	}
	bob, err := store.CreateStudent("bob")
	if err != nil {
		t.Fatalf("create bob: %v", err)
	}
	guest, err := store.ClaimGuestStudent("guest_active_test", "Guest")
	if err != nil {
		t.Fatalf("guest: %v", err)
	}

	live := []*Student{ada, bob, guest}
	for i, st := range live {
		sid := "sess-live-" + st.ID
		if err := store.EnsureSession(sid, st.ID); err != nil {
			t.Fatalf("session %d: %v", i, err)
		}
		seedActiveSession(t, store, sid, st.ID, "att-live-"+st.ID, now)
	}
	// Ada also has a stale seat from earlier; it must not add to the count.
	oldSid := "sess-old-" + ada.ID
	if err := store.EnsureSession(oldSid, ada.ID); err != nil {
		t.Fatalf("old session: %v", err)
	}
	seedActiveSession(t, store, oldSid, ada.ID, "att-old-"+ada.ID, now.Add(-2*time.Hour))

	// Ada and the guest answered today; Bob's only attempt is from yesterday.
	record := func(st *Student, when time.Time, key string) {
		t.Helper()
		if _, err := store.RecordAttempt(AttemptEntry{
			SessionID:      "sess-live-" + st.ID,
			StudentID:      st.ID,
			ConceptID:      "c1",
			Answer:         "1",
			Expected:       "1",
			Correct:        true,
			ElapsedSeconds: 4,
			Timestamp:      when,
			Question:       "q",
			Source:         "practice",
		}); err != nil {
			t.Fatalf("attempt %s: %v", key, err)
		}
	}
	record(ada, now.Add(-time.Hour), "ada")
	record(guest, now.Add(-time.Hour), "guest")
	record(bob, now.AddDate(0, 0, -1), "bob-yesterday")

	activeNow, activeNowGuests, activeToday, activeTodayGuests, err :=
		store.ActiveLearnerCounts(now.Add(-15*time.Minute), dayStart)
	if err != nil {
		t.Fatalf("counts: %v", err)
	}
	if activeNow != 2 || activeNowGuests != 1 {
		t.Errorf("active now = %d accounts / %d guests, want 2/1", activeNow, activeNowGuests)
	}
	if activeToday != 1 || activeTodayGuests != 1 {
		t.Errorf("active today = %d accounts / %d guests, want 1/1", activeToday, activeTodayGuests)
	}
}

// TestActiveLearnerCounts_DoesNotTreatGuestXAsGuest pins the wildcard fix.
// `LIKE 'guest_%'` treats `_` as one character, so it classified guestX as a
// guest. The prefix test must be exact.
func TestActiveLearnerCounts_DoesNotTreatGuestXAsGuest(t *testing.T) {
	store := newTestStore(t)
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)

	uuid, err := store.CreateStudent("uuid guest")
	if err != nil {
		t.Fatalf("create uuid guest: %v", err)
	}
	guestX, err := store.ClaimGuestStudent("guestX", "guestX")
	if err != nil {
		t.Fatalf("create guestX: %v", err)
	}
	guest, err := store.ClaimGuestStudent("guest_abc", "real guest")
	if err != nil {
		t.Fatalf("create real guest: %v", err)
	}

	for _, st := range []*Student{uuid, guestX, guest} {
		sid := "sess-prefix-" + st.ID
		if err := store.EnsureSession(sid, st.ID); err != nil {
			t.Fatalf("session %s: %v", st.ID, err)
		}
		seedActiveSession(t, store, sid, st.ID, "att-prefix-"+st.ID, now)
	}

	activeNow, guests, _, _, err := store.ActiveLearnerCounts(now.Add(-15*time.Minute), now)
	if err != nil {
		t.Fatalf("counts: %v", err)
	}
	if activeNow != 2 || guests != 1 {
		t.Fatalf("active now = %d accounts / %d guests, want 2/1", activeNow, guests)
	}
}

// seedActiveSession writes the row EnsureSession leaves behind, with the fields
// it fills in for itself. The sweep is tested at a fixed instant, so the stale
// seat needs an updated_at two hours old and no exported method can express
// that. The two engines store the column differently: SQLite keeps a datetime
// string and compares it with datetime(), while Postgres keeps text and casts
// it to timestamptz, so the value is formatted for whichever engine is running.
func seedActiveSession(t *testing.T, repo Repository, sessionID, studentID, attemptID string, updatedAt time.Time) {
	t.Helper()
	const columns = "(session_id, student_id, concept_id, expected_answer, attempt_id, question, updated_at)"
	var (
		db    *sql.DB
		query string
		stamp string
	)
	switch s := repo.(type) {
	case *SQLiteStore:
		db = s.db
		query = "INSERT INTO active_sessions " + columns + " VALUES (?, ?, 'c1', '1', ?, 'q', ?)"
		stamp = updatedAt.UTC().Format("2006-01-02 15:04:05")
	case *PostgresStore:
		db = s.db
		query = "INSERT INTO active_sessions " + columns + " VALUES ($1, $2, 'c1', '1', $3, 'q', $4)"
		stamp = updatedAt.UTC().Format("2006-01-02 15:04:05.000000+00")
	default:
		t.Fatalf("seedActiveSession: unsupported store %T", repo)
	}
	if _, err := db.Exec(query, sessionID, studentID, attemptID, stamp); err != nil {
		t.Fatalf("insert active session %s: %v", sessionID, err)
	}
}
