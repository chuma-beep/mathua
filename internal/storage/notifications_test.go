package storage

import (
	"errors"
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

func TestNotificationFeedStorage(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	learner := seedStudent(t, store)
	other := seedStudent(t, store)

	base := time.Now().UTC().Truncate(time.Second)
	items := []*Notification{
		{ID: "n-old", StudentID: learner.ID, Kind: "review", DedupeKey: "old", Title: "Old", CreatedAt: base.Add(-3 * time.Hour)},
		{ID: "n-a", StudentID: learner.ID, Kind: "review", DedupeKey: "a", Title: "A", CreatedAt: base.Add(-time.Hour)},
		{ID: "n-b", StudentID: learner.ID, Kind: "review", DedupeKey: "b", Title: "B", CreatedAt: base.Add(-time.Hour)},
		{ID: "n-canceled", StudentID: learner.ID, Kind: "review", DedupeKey: "canceled", Title: "Canceled", CreatedAt: base, CanceledAt: base},
		{ID: "n-other", StudentID: other.ID, Kind: "review", DedupeKey: "other", Title: "Other", CreatedAt: base},
	}
	for _, item := range items {
		created, err := ns.CreateNotification(item)
		if err != nil {
			t.Fatalf("create %s: %v", item.ID, err)
		}
		if !created {
			t.Errorf("create %s reported duplicate", item.ID)
		}
	}
	// The dedupe key is per learner, and a retry must not create a second row.
	dup, err := ns.CreateNotification(&Notification{
		ID: "n-duplicate", StudentID: learner.ID, Kind: "review", DedupeKey: "a", Title: "Duplicate",
		CreatedAt: base,
	})
	if err != nil {
		t.Fatalf("create duplicate: %v", err)
	}
	if dup {
		t.Error("duplicate dedupe key created a row")
	}

	page, err := ns.ListNotifications(learner.ID, NotificationListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("list first page: %v", err)
	}
	if len(page.Notifications) != 2 || !page.HasMore {
		t.Fatalf("first page = %+v, want two items and more", page)
	}
	// Same-second rows order by id descending, and canceled/other rows never
	// appear.
	if page.Notifications[0].ID != "n-b" || page.Notifications[1].ID != "n-a" {
		t.Fatalf("first page order = %+v, want n-b, n-a", page.Notifications)
	}
	last := page.Notifications[len(page.Notifications)-1]
	next, err := ns.ListNotifications(learner.ID, NotificationListOptions{
		Limit:           2,
		BeforeCreatedAt: last.CreatedAt.Unix(),
		BeforeID:        last.ID,
	})
	if err != nil {
		t.Fatalf("list next page: %v", err)
	}
	if len(next.Notifications) != 1 || next.Notifications[0].ID != "n-old" || next.HasMore {
		t.Fatalf("next page = %+v, want n-old only", next)
	}
	for _, n := range append(page.Notifications, next.Notifications...) {
		if n.ID == "n-canceled" || n.ID == "n-other" {
			t.Fatalf("list returned excluded notification %s", n.ID)
		}
	}

	unread, err := ns.NotificationUnreadCount(learner.ID)
	if err != nil || unread != 3 {
		t.Fatalf("unread = %d (%v), want 3", unread, err)
	}
	if err := ns.MarkNotificationRead(learner.ID, "n-b", time.Now().UTC()); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	if unread, _ = ns.NotificationUnreadCount(learner.ID); unread != 2 {
		t.Fatalf("unread after one read = %d, want 2", unread)
	}
	// Already read is idempotent; another learner's id is not found.
	if err := ns.MarkNotificationRead(learner.ID, "n-b", time.Now().UTC()); err != nil {
		t.Fatalf("repeat mark read: %v", err)
	}
	if err := ns.MarkNotificationRead(other.ID, "n-a", time.Now().UTC()); !errors.Is(err, ErrNotificationNotFound) {
		t.Fatalf("foreign mark read = %v, want ErrNotificationNotFound", err)
	}
	unreadOnly, err := ns.ListNotifications(learner.ID, NotificationListOptions{Limit: 10, UnreadOnly: true})
	if err != nil {
		t.Fatalf("list unread: %v", err)
	}
	if len(unreadOnly.Notifications) != 2 {
		t.Fatalf("unread list = %+v, want 2", unreadOnly.Notifications)
	}

	if err := ns.MarkAllNotificationsRead(learner.ID, time.Now().UTC()); err != nil {
		t.Fatalf("mark all read: %v", err)
	}
	if unread, _ = ns.NotificationUnreadCount(learner.ID); unread != 0 {
		t.Fatalf("learner unread after mark all = %d, want 0", unread)
	}
	if unread, _ = ns.NotificationUnreadCount(other.ID); unread != 1 {
		t.Fatalf("other unread after mark all = %d, want 1", unread)
	}
}

func TestNotificationRecipientsIncludeStudentDeliveryFields(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	now := time.Now().UTC()
	first := seedStudent(t, store)
	second := seedStudent(t, store)
	if err := store.SetEmail(first.ID, "first@example.com"); err != nil {
		t.Fatalf("set first email: %v", err)
	}
	if err := store.SetEmailVerified(first.ID, true); err != nil {
		t.Fatalf("verify first email: %v", err)
	}
	if err := store.SetEmail(second.ID, "second@example.com"); err != nil {
		t.Fatalf("set second email: %v", err)
	}
	for _, id := range []string{first.ID, second.ID} {
		if _, err := ns.EnsureNotificationPreferences(id, now); err != nil {
			t.Fatalf("ensure %s: %v", id, err)
		}
	}

	seen := map[string]NotificationRecipient{}
	cursor := ""
	for {
		page, err := ns.ListNotificationRecipients(cursor, 1)
		if err != nil {
			t.Fatalf("walk recipients: %v", err)
		}
		if len(page) == 0 {
			break
		}
		if len(page) > 1 {
			t.Fatalf("recipient page = %d rows, want at most 1", len(page))
		}
		for _, r := range page {
			seen[r.StudentID] = r
		}
		cursor = page[len(page)-1].StudentID
	}
	r, ok := seen[first.ID]
	if !ok {
		t.Fatalf("first learner missing from recipient walk: %+v", seen)
	}
	if r.Name == "" || r.Email != "first@example.com" || !r.EmailVerified {
		t.Fatalf("first recipient = %+v, want name, email and verified flag", r)
	}
	if !r.Preferences.InAppEnabled || r.Preferences.Timezone != NotificationDefaultTimezone {
		t.Fatalf("first preferences = %+v, want defaults", r.Preferences)
	}
	if r2, ok := seen[second.ID]; !ok || r2.EmailVerified {
		t.Fatalf("second recipient = %+v (found %v), want unverified", r2, ok)
	}
}

func TestNotificationDeliveryClaimLifecycle(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	learner := seedStudent(t, store)
	now := time.Now().UTC().Truncate(time.Second)
	n := &Notification{
		ID: "n-delivery", StudentID: learner.ID, Kind: "review", DedupeKey: "delivery",
		Title: "Reviews are waiting", Body: "You have reviews due.", Href: "/review", CreatedAt: now,
	}
	if created, err := ns.CreateNotification(n); err != nil || !created {
		t.Fatalf("create notification: %v %v", created, err)
	}

	inserted, err := ns.CreateNotificationDelivery(&NotificationDelivery{
		NotificationID: n.ID,
		Channel:        NotificationChannelEmail,
		NextAttemptAt:  now,
	})
	if err != nil || !inserted {
		t.Fatalf("create delivery: %v %v", inserted, err)
	}
	if dup, err := ns.CreateNotificationDelivery(&NotificationDelivery{
		NotificationID: n.ID, Channel: NotificationChannelEmail, NextAttemptAt: now,
	}); err != nil || dup {
		t.Fatalf("duplicate delivery = %v (%v), want false", dup, err)
	}

	jobs, err := ns.ClaimNotificationDeliveries(now, time.Minute, 10)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if len(jobs) != 1 || jobs[0].NotificationID != n.ID || jobs[0].Attempts != 1 {
		t.Fatalf("claimed = %+v, want one job at attempt 1", jobs)
	}
	if jobs[0].StudentID != learner.ID || jobs[0].Title != n.Title || jobs[0].Body != n.Body {
		t.Fatalf("job projection = %+v", jobs[0])
	}
	// A claimed row is invisible until its lease expires.
	if again, err := ns.ClaimNotificationDeliveries(now.Add(30*time.Second), time.Minute, 10); err != nil || len(again) != 0 {
		t.Fatalf("claim during lease = %+v (%v), want none", again, err)
	}
	if err := ns.CompleteNotificationDelivery(n.ID, NotificationChannelEmail, "mailer", now); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if again, err := ns.ClaimNotificationDeliveries(now.Add(2*time.Minute), time.Minute, 10); err != nil || len(again) != 0 {
		t.Fatalf("claim after complete = %+v (%v), want none", again, err)
	}

	// A retry is claimable only at its next attempt time, and increments the
	// attempt count.
	if err := ns.RetryNotificationDelivery(n.ID, NotificationChannelEmail, "temporary", now.Add(5*time.Minute)); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if again, err := ns.ClaimNotificationDeliveries(now.Add(4*time.Minute), time.Minute, 10); err != nil || len(again) != 0 {
		t.Fatalf("claim before retry = %+v (%v), want none", again, err)
	}
	again, err := ns.ClaimNotificationDeliveries(now.Add(5*time.Minute), time.Minute, 10)
	if err != nil {
		t.Fatalf("claim retry: %v", err)
	}
	if len(again) != 1 || again[0].Attempts != 2 {
		t.Fatalf("retry claim = %+v, want attempt 2", again)
	}
	if err := ns.FailNotificationDelivery(n.ID, NotificationChannelEmail, "permanent", now); err != nil {
		t.Fatalf("fail: %v", err)
	}
	if again, err := ns.ClaimNotificationDeliveries(now.Add(24*time.Hour), time.Minute, 10); err != nil || len(again) != 0 {
		t.Fatalf("claim after fail = %+v (%v), want none", again, err)
	}
}

func TestNotificationDeliverySuppressionAndScanState(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	learner := seedStudent(t, store)
	now := time.Now().UTC().Truncate(time.Second)
	n := &Notification{
		ID: "n-suppress", StudentID: learner.ID, Kind: "review", DedupeKey: "suppress",
		Title: "Reviews are waiting", Body: "Body", CreatedAt: now,
	}
	if _, err := ns.CreateNotification(n); err != nil {
		t.Fatalf("create notification: %v", err)
	}
	if _, err := ns.CreateNotificationDelivery(&NotificationDelivery{
		NotificationID: n.ID, Channel: NotificationChannelEmail, NextAttemptAt: now,
	}); err != nil {
		t.Fatalf("create delivery: %v", err)
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now, time.Minute, 10); err != nil || len(jobs) != 1 {
		t.Fatalf("claim = %+v (%v)", jobs, err)
	}
	if err := ns.SuppressNotificationDelivery(n.ID, NotificationChannelEmail, "email_unverified", now); err != nil {
		t.Fatalf("suppress: %v", err)
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now.Add(time.Hour), time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("claim after suppress = %+v (%v), want none", jobs, err)
	}

	if got, err := ns.GetNotificationScanState(learner.ID, "review"); err != nil || got != 0 {
		t.Fatalf("initial scan state = %d (%v), want 0", got, err)
	}
	if err := ns.SetNotificationScanState(learner.ID, "review", 20261111, now); err != nil {
		t.Fatalf("set scan state: %v", err)
	}
	if got, err := ns.GetNotificationScanState(learner.ID, "review"); err != nil || got != 20261111 {
		t.Fatalf("scan state = %d (%v), want 20261111", got, err)
	}
	if err := ns.SetNotificationScanState(learner.ID, "review", 20261112, now.Add(time.Minute)); err != nil {
		t.Fatalf("update scan state: %v", err)
	}
	if got, _ := ns.GetNotificationScanState(learner.ID, "review"); got != 20261112 {
		t.Fatalf("updated scan state = %d, want 20261112", got)
	}
}

func ExampleNotificationPreferences() {
	p := DefaultNotificationPreferences("s1", time.Unix(0, 0).UTC())
	fmt.Println(p.Timezone, p.PreferredTimeMin, p.EmailEnabled, p.InAppEnabled)
	// Output: UTC 1080 false true
}

func TestNotificationDigestRoundTrip(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	learner := seedStudent(t, store)

	prefs, err := ns.EnsureNotificationPreferences(learner.ID, time.Now().UTC())
	if err != nil {
		t.Fatalf("ensure preferences: %v", err)
	}
	if prefs.Digest != NotificationDigestImmediate {
		t.Fatalf("digest default = %q, want %q", prefs.Digest, NotificationDigestImmediate)
	}
	prefs.Digest = NotificationDigestWeekly
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("upsert digest: %v", err)
	}
	after, err := ns.GetNotificationPreferences(learner.ID)
	if err != nil {
		t.Fatalf("get digest: %v", err)
	}
	if after.Digest != NotificationDigestWeekly {
		t.Fatalf("digest = %q, want weekly", after.Digest)
	}
	recipients, err := ns.ListNotificationRecipients("", 10)
	if err != nil {
		t.Fatalf("list recipients: %v", err)
	}
	found := false
	for _, recipient := range recipients {
		if recipient.StudentID != learner.ID {
			continue
		}
		found = true
		if recipient.Preferences.Digest != NotificationDigestWeekly {
			t.Fatalf("recipient digest = %q, want weekly", recipient.Preferences.Digest)
		}
	}
	if !found {
		t.Fatal("learner missing from the recipient projection")
	}
}

func TestNotificationEngagementReplayAndRetention(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	learner := seedStudent(t, store)
	now := time.Now().UTC().Truncate(time.Second)

	if _, err := ns.CreateNotification(&Notification{
		ID: "n-engage", StudentID: learner.ID, Kind: "review", DedupeKey: "engage",
		Title: "Reviews are due", CreatedAt: now,
	}); err != nil {
		t.Fatalf("create notification: %v", err)
	}
	if _, err := ns.CreateNotificationDelivery(&NotificationDelivery{
		NotificationID: "n-engage", Channel: NotificationChannelEmail,
		Status: NotificationDeliveryFailed, Attempts: 3, LastError: "smtp down",
	}); err != nil {
		t.Fatalf("create delivery: %v", err)
	}

	if err := ns.MarkNotificationOpened(learner.ID, "n-engage", now); err != nil {
		t.Fatalf("mark opened: %v", err)
	}
	page, err := ns.ListNotifications(learner.ID, NotificationListOptions{})
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(page.Notifications) != 1 || page.Notifications[0].ReadAt.IsZero() {
		t.Fatalf("open did not mark the notification read: %+v", page.Notifications)
	}
	deliveries, err := ns.ListNotificationDeliveries("n-engage")
	if err != nil {
		t.Fatalf("list deliveries: %v", err)
	}
	if len(deliveries) != 1 || deliveries[0].OpenedAt.IsZero() {
		t.Fatalf("open was not stamped on the delivery: %+v", deliveries)
	}

	if err := ns.MarkNotificationClicked(learner.ID, "n-engage", now.Add(time.Minute)); err != nil {
		t.Fatalf("mark clicked: %v", err)
	}
	deliveries, _ = ns.ListNotificationDeliveries("n-engage")
	if deliveries[0].ClickedAt.IsZero() {
		t.Fatal("click was not stamped on the delivery")
	}

	if err := ns.ReplayNotificationDelivery("n-engage", NotificationChannelEmail, now); err != nil {
		t.Fatalf("replay delivery: %v", err)
	}
	deliveries, _ = ns.ListNotificationDeliveries("n-engage")
	if deliveries[0].Status != NotificationDeliveryPending || deliveries[0].Attempts != 0 || deliveries[0].LastError != "" {
		t.Fatalf("replay did not reset the delivery: %+v", deliveries[0])
	}
	if err := ns.ReplayNotificationDelivery("n-engage", NotificationChannelEmail, now); err == nil {
		t.Fatal("a pending delivery was replayed twice")
	}

	old := now.AddDate(0, 0, -400)
	if _, err := ns.CreateNotification(&Notification{
		ID: "n-old", StudentID: learner.ID, Kind: "review", DedupeKey: "old",
		Title: "Old", CreatedAt: old,
	}); err != nil {
		t.Fatalf("create old notification: %v", err)
	}
	if _, err := ns.CreateNotificationDelivery(&NotificationDelivery{
		NotificationID: "n-old", Channel: NotificationChannelInApp,
		Status: NotificationDeliverySent, SentAt: old,
	}); err != nil {
		t.Fatalf("create old delivery: %v", err)
	}
	pruned, err := ns.PruneNotificationHistory(now.AddDate(0, 0, -180))
	if err != nil {
		t.Fatalf("prune history: %v", err)
	}
	if pruned.Notifications != 1 || pruned.Deliveries != 1 {
		t.Fatalf("pruned = %+v, want one notification and one delivery", pruned)
	}
	remaining, _ := ns.ListNotificationDeliveries("n-old")
	if len(remaining) != 0 {
		t.Fatalf("old deliveries survived the prune: %+v", remaining)
	}
	if _, err := ns.ListNotificationDeliveries("n-engage"); err != nil {
		t.Fatalf("prune removed the live notification: %v", err)
	}
}

func TestNotificationPushSubscriptionHealth(t *testing.T) {
	store := newTestStore(t)
	ns := notificationStore(t, store)
	learner := seedStudent(t, store)
	now := time.Now().UTC().Truncate(time.Second)

	sub := &NotificationPushSubscription{
		StudentID: learner.ID, Endpoint: "https://push.example/abc",
		P256DH: "key", Auth: "secret", UserAgent: "Firefox",
	}
	if err := ns.UpsertNotificationPushSubscription(sub); err != nil {
		t.Fatalf("upsert subscription: %v", err)
	}
	subs, err := ns.ListNotificationPushSubscriptions(learner.ID)
	if err != nil {
		t.Fatalf("list subscriptions: %v", err)
	}
	if len(subs) != 1 || subs[0].P256DH != "key" {
		t.Fatalf("subscription did not round-trip: %+v", subs)
	}

	if err := ns.MarkPushSubscriptionSuccess(sub.Endpoint, now); err != nil {
		t.Fatalf("mark success: %v", err)
	}
	if err := ns.MarkPushSubscriptionFailure(sub.Endpoint, now); err != nil {
		t.Fatalf("mark failure: %v", err)
	}
	subs, _ = ns.ListNotificationPushSubscriptions(learner.ID)
	if subs[0].Failures != 1 || subs[0].LastSuccess.IsZero() {
		t.Fatalf("health did not update: %+v", subs[0])
	}

	if removed, err := ns.PrunePushSubscriptions(now.AddDate(0, 0, 1), 1); err != nil || removed != 1 {
		t.Fatalf("prune dead endpoints = %d, %v, want one", removed, err)
	}
	subs, _ = ns.ListNotificationPushSubscriptions(learner.ID)
	if len(subs) != 0 {
		t.Fatalf("dead endpoint survived the prune: %+v", subs)
	}
}
