package storage

import (
	"fmt"
	"testing"
	"time"
)

func notificationStore(t *testing.T, repo Repository) NotificationStore {
	t.Helper()
	store, ok := repo.(NotificationStore)
	if !ok {
		t.Fatalf("repository %T does not implement NotificationStore", repo)
	}
	return store
}

func assertNotificationDefaults(t *testing.T, p *NotificationPreferences) {
	t.Helper()
	if p == nil {
		t.Fatal("notification preferences are nil")
	}
	if p.EmailEnabled {
		t.Error("email should default off")
	}
	if !p.InAppEnabled {
		t.Error("in-app should default on")
	}
	if !p.ReviewsEnabled {
		t.Error("reviews should default on")
	}
	if p.StudyEnabled {
		t.Error("study reminders should default off")
	}
	if !p.MilestonesEnabled {
		t.Error("milestones should default on")
	}
	if p.InactivityEnabled {
		t.Error("inactivity should default off")
	}
	if p.PreferredTimeMin != NotificationDefaultPreferredTimeMin {
		t.Errorf("preferred time = %d, want %d", p.PreferredTimeMin, NotificationDefaultPreferredTimeMin)
	}
	if p.StudyTimeMin != NotificationDefaultStudyTimeMin {
		t.Errorf("study time = %d, want %d", p.StudyTimeMin, NotificationDefaultStudyTimeMin)
	}
	if p.Timezone != NotificationDefaultTimezone {
		t.Errorf("timezone = %q, want %q", p.Timezone, NotificationDefaultTimezone)
	}
	if p.QuietStartMin != NotificationDefaultQuietStartMin || p.QuietEndMin != NotificationDefaultQuietEndMin {
		t.Errorf("quiet hours = %d-%d, want %d-%d", p.QuietStartMin, p.QuietEndMin, NotificationDefaultQuietStartMin, NotificationDefaultQuietEndMin)
	}
	if p.PausedUntil != nil {
		t.Errorf("paused_until = %v, want nil", p.PausedUntil)
	}
}

func TestNotificationSchemaTablesExist(t *testing.T) {
	store := newTestStore(t)
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Close() error
		Err() error
	}
	var err error
	switch s := store.(type) {
	case *SQLiteStore:
		rows, err = s.db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name IN
			('notification_preferences','notifications','notification_deliveries','notification_scan_state')`)
	case *PostgresStore:
		rows, err = s.db.Query(`SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename IN
			('notification_preferences','notifications','notification_deliveries','notification_scan_state')`)
	default:
		t.Fatalf("unsupported store %T", store)
	}
	if err != nil {
		t.Fatalf("list notification tables: %v", err)
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table: %v", err)
		}
		seen[name] = true
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("table rows: %v", err)
	}
	for _, name := range []string{"notification_preferences", "notifications", "notification_deliveries", "notification_scan_state"} {
		if !seen[name] {
			t.Errorf("missing table %s", name)
		}
	}
}

func TestNotificationPreferencesDefaultsAndRoundTrip(t *testing.T) {
	store := newTestStore(t)
	prefs := notificationStore(t, store)
	learner := seedStudent(t, store)

	got, err := prefs.GetNotificationPreferences(learner.ID)
	if err != nil {
		t.Fatalf("get defaults: %v", err)
	}
	assertNotificationDefaults(t, got)

	pause := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	got.EmailEnabled = true
	got.PreferredTimeMin = 7 * 60
	got.StudyEnabled = true
	got.StudyDays = "mon,wed,fri"
	got.StudyTimeMin = 6 * 60
	got.Timezone = "America/New_York"
	got.QuietStartMin = 21 * 60
	got.QuietEndMin = 6 * 60
	got.PausedUntil = &pause
	if err := prefs.UpsertNotificationPreferences(got); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	after, err := prefs.GetNotificationPreferences(learner.ID)
	if err != nil {
		t.Fatalf("get after upsert: %v", err)
	}
	if !after.EmailEnabled || !after.StudyEnabled {
		t.Errorf("booleans did not round-trip: %+v", after)
	}
	if after.PreferredTimeMin != 7*60 || after.StudyTimeMin != 6*60 {
		t.Errorf("times did not round-trip: %+v", after)
	}
	if after.StudyDays != "mon,wed,fri" || after.Timezone != "America/New_York" {
		t.Errorf("schedule did not round-trip: %+v", after)
	}
	if after.QuietStartMin != 21*60 || after.QuietEndMin != 6*60 {
		t.Errorf("quiet hours did not round-trip: %+v", after)
	}
	if after.PausedUntil == nil || !after.PausedUntil.Equal(pause) {
		t.Errorf("pause did not round-trip: got %v, want %v", after.PausedUntil, pause)
	}
}

func TestNotificationPreferencesBackfillExistingStudent(t *testing.T) {
	store := newTestStore(t)
	prefs := notificationStore(t, store)
	const id = "legacy-student"
	var removeMarker, rerunMigration, countQuery func() error
	var insertLegacy func() error
	switch s := store.(type) {
	case *SQLiteStore:
		insertLegacy = func() error {
			_, err := s.db.Exec("INSERT INTO students (id, name, settings) VALUES (?, ?, '{}')", id, "Legacy")
			return err
		}
		removeMarker = func() error {
			_, err := s.db.Exec("DELETE FROM schema_migrations WHERE version = 7")
			return err
		}
		rerunMigration = func() error { return authMigrate(s.db) }
		countQuery = func() error {
			var count int
			if err := s.db.QueryRow("SELECT COUNT(*) FROM notification_preferences WHERE student_id = ?", id).Scan(&count); err != nil {
				return err
			}
			if count != 1 {
				return fmt.Errorf("notification_preferences rows = %d, want 1", count)
			}
			return nil
		}
	case *PostgresStore:
		insertLegacy = func() error {
			_, err := s.db.Exec("INSERT INTO students (id, name, settings) VALUES ($1, $2, '{}')", id, "Legacy")
			return err
		}
		removeMarker = func() error {
			_, err := s.db.Exec("DELETE FROM schema_migrations WHERE version = 7")
			return err
		}
		rerunMigration = func() error { return pgAuthMigrate(s.db) }
		countQuery = func() error {
			var count int
			if err := s.db.QueryRow("SELECT COUNT(*) FROM notification_preferences WHERE student_id = $1", id).Scan(&count); err != nil {
				return err
			}
			if count != 1 {
				return fmt.Errorf("notification_preferences rows = %d, want 1", count)
			}
			return nil
		}
	default:
		t.Fatalf("unsupported store %T", store)
	}

	if err := insertLegacy(); err != nil {
		t.Fatalf("insert legacy student: %v", err)
	}
	if err := removeMarker(); err != nil {
		t.Fatalf("remove migration marker: %v", err)
	}
	if err := rerunMigration(); err != nil {
		t.Fatalf("run v7 migration: %v", err)
	}
	got, err := prefs.GetNotificationPreferences(id)
	if err != nil {
		t.Fatalf("get backfilled preferences: %v", err)
	}
	assertNotificationDefaults(t, got)

	// The backfill must not overwrite a preference a learner already set.
	got.EmailEnabled = true
	if err := prefs.UpsertNotificationPreferences(got); err != nil {
		t.Fatalf("set preference before rerun: %v", err)
	}
	if err := removeMarker(); err != nil {
		t.Fatalf("remove migration marker again: %v", err)
	}
	if err := rerunMigration(); err != nil {
		t.Fatalf("rerun v7 migration: %v", err)
	}
	again, err := prefs.GetNotificationPreferences(id)
	if err != nil {
		t.Fatalf("get after rerun: %v", err)
	}
	if !again.EmailEnabled {
		t.Error("v7 rerun overwrote an existing preference")
	}
	if err := countQuery(); err != nil {
		t.Fatal(err)
	}
}

func TestNotificationPreferencesInitializedOnAccountCreation(t *testing.T) {
	cases := []struct {
		name   string
		create func(*testing.T, Repository) (*Student, error)
	}{
		{"CreateStudent", func(_ *testing.T, s Repository) (*Student, error) { return s.CreateStudent("Learner") }},
		{"CreateUser", func(_ *testing.T, s Repository) (*Student, error) { return s.CreateUser("Learner", "learner", "hash") }},
		{"CreateGoogleUser", func(_ *testing.T, s Repository) (*Student, error) {
			return s.CreateGoogleUser("Google", "google@example.com", "google-id", "")
		}},
		{"CreateOAuthUser", func(_ *testing.T, s Repository) (*Student, error) {
			return s.CreateOAuthUser("github", "github-id", "OAuth", "oauth@example.com", true, "")
		}},
		{"ClaimGuestStudent", func(_ *testing.T, s Repository) (*Student, error) {
			return s.ClaimGuestStudent("guest-notify", "Guest")
		}},
		{"AdoptGuestAsUser", func(t *testing.T, s Repository) (*Student, error) {
			guest, err := s.ClaimGuestStudent("guest-adopt", "Guest")
			if err != nil {
				t.Fatalf("claim guest: %v", err)
			}
			return s.AdoptGuestAsUser(guest.ID, "Adopted", "adopted", "hash")
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := newTestStore(t)
			prefs := notificationStore(t, store)
			learner, err := tc.create(t, store)
			if err != nil {
				t.Fatalf("create account: %v", err)
			}
			got, err := prefs.GetNotificationPreferences(learner.ID)
			if err != nil {
				t.Fatalf("get preferences: %v", err)
			}
			assertNotificationDefaults(t, got)
		})
	}
}

func TestResetPreservesNotificationPreferencesDeleteRemovesThem(t *testing.T) {
	store := newTestStore(t)
	prefs := notificationStore(t, store)
	learner := seedStudent(t, store)

	if err := store.ResetProgress(learner.ID); err != nil {
		t.Fatalf("reset progress: %v", err)
	}
	got, err := prefs.GetNotificationPreferences(learner.ID)
	if err != nil {
		t.Fatalf("get after reset: %v", err)
	}
	if got == nil {
		t.Fatal("ResetProgress removed notification preferences")
	}

	if err := store.DeleteAccount(learner.ID); err != nil {
		t.Fatalf("delete account: %v", err)
	}
	got, err = prefs.GetNotificationPreferences(learner.ID)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if got != nil {
		t.Errorf("DeleteAccount left notification preferences: %+v", got)
	}
}

func TestNotificationPreferencesEnsureIsIdempotent(t *testing.T) {
	store := newTestStore(t)
	prefs := notificationStore(t, store)
	learner := seedStudent(t, store)
	now := time.Now().UTC()
	for i := 0; i < 3; i++ {
		if _, err := prefs.EnsureNotificationPreferences(learner.ID, now); err != nil {
			t.Fatalf("ensure %d: %v", i, err)
		}
	}
	var count int
	switch s := store.(type) {
	case *SQLiteStore:
		if err := s.db.QueryRow("SELECT COUNT(*) FROM notification_preferences WHERE student_id = ?", learner.ID).Scan(&count); err != nil {
			t.Fatalf("count: %v", err)
		}
	case *PostgresStore:
		if err := s.db.QueryRow("SELECT COUNT(*) FROM notification_preferences WHERE student_id = $1", learner.ID).Scan(&count); err != nil {
			t.Fatalf("count: %v", err)
		}
	default:
		t.Fatalf("unsupported store %T", store)
	}
	if count != 1 {
		t.Errorf("notification_preferences rows = %d, want 1", count)
	}
}

func ExampleNotificationPreferences() {
	p := DefaultNotificationPreferences("s1", time.Unix(0, 0).UTC())
	fmt.Println(p.Timezone, p.PreferredTimeMin, p.EmailEnabled, p.InAppEnabled)
	// Output: UTC 1080 false true
}
