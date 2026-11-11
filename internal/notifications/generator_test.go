package notifications

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

type fakeDue struct {
	count int
	err   error
	calls int
}

func (f *fakeDue) ReviewsDueReadOnly(string) (int, error) {
	f.calls++
	return f.count, f.err
}

func generatorStore(t *testing.T) (*storage.SQLiteStore, storage.NotificationStore, *storage.Student) {
	t.Helper()
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	student, err := store.CreateStudent("Ada Example")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := store.SetEmail(student.ID, "ada@example.com"); err != nil {
		t.Fatalf("set email: %v", err)
	}
	if err := store.SetEmailVerified(student.ID, true); err != nil {
		t.Fatalf("verify email: %v", err)
	}
	return store, store, student
}

func savePrefs(t *testing.T, store storage.NotificationStore, studentID string, mutate func(*storage.NotificationPreferences)) {
	t.Helper()
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	prefs := storage.DefaultNotificationPreferences(studentID, now)
	prefs.PreferredTimeMin = 0
	prefs.QuietStartMin = 0
	prefs.QuietEndMin = 0
	if mutate != nil {
		mutate(&prefs)
	}
	if err := store.UpsertNotificationPreferences(&prefs); err != nil {
		t.Fatalf("save prefs: %v", err)
	}
}

func generatorAt(store storage.NotificationStore, due DueCounter, now time.Time, emailAvailable bool) *Generator {
	return NewGenerator(GeneratorOptions{
		Store:          store,
		Due:            due,
		EmailAvailable: emailAvailable,
		Now:            func() time.Time { return now },
	})
}

func TestGeneratorCreatesOneReminderPerLocalDay(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.EmailEnabled = true
	})
	due := &fakeDue{count: 3}
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	gen := generatorAt(ns, due, now, true)

	result, err := gen.RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 1 || result.EmailPending != 1 || result.NoReviews != 0 {
		t.Fatalf("first run = %+v, want one notification and one email", result)
	}
	if due.calls != 1 {
		t.Fatalf("due calls = %d, want 1", due.calls)
	}

	page, err := ns.ListNotifications(learner.ID, storage.NotificationListOptions{Limit: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(page.Notifications) != 1 {
		t.Fatalf("notifications = %d, want 1", len(page.Notifications))
	}
	n := page.Notifications[0]
	if n.Kind != ReviewKind || n.Href != "/review" || n.DedupeKey != "review:2026-11-11" {
		t.Fatalf("notification = %+v", n)
	}
	if n.Title != "Reviews are waiting" || n.Body == "" {
		t.Fatalf("copy = %q / %q", n.Title, n.Body)
	}

	jobs, err := ns.ClaimNotificationDeliveries(now, time.Minute, 10)
	if err != nil || len(jobs) != 1 {
		t.Fatalf("email claim = %+v (%v), want one", jobs, err)
	}
	if jobs[0].NotificationID != n.ID || jobs[0].Email != "ada@example.com" || !jobs[0].EmailVerified {
		t.Fatalf("email job = %+v", jobs[0])
	}

	// The same local day is not evaluated twice, and the review counter is not
	// called again.
	result, err = gen.RunOnce()
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if result.AlreadyGenerated != 1 || result.Notifications != 0 {
		t.Fatalf("second run = %+v, want already generated", result)
	}
	if due.calls != 1 {
		t.Fatalf("due calls after second run = %d, want 1", due.calls)
	}

	// The next local day is a new dedupe key and a new reminder.
	tomorrow := now.Add(24 * time.Hour)
	gen = generatorAt(ns, due, tomorrow, true)
	result, err = gen.RunOnce()
	if err != nil {
		t.Fatalf("tomorrow run: %v", err)
	}
	if result.Notifications != 1 || result.AlreadyGenerated != 0 {
		t.Fatalf("tomorrow run = %+v, want a new notification", result)
	}
	if due.calls != 2 {
		t.Fatalf("due calls after tomorrow = %d, want 2", due.calls)
	}
}

func TestGeneratorRespectsWindowQuietHoursAndPause(t *testing.T) {
	cases := []struct {
		name   string
		now    time.Time
		mutate func(*storage.NotificationPreferences)
	}{
		{
			name: "before preferred time",
			now:  time.Date(2026, 11, 11, 9, 0, 0, 0, time.UTC),
			mutate: func(p *storage.NotificationPreferences) {
				p.PreferredTimeMin = 18 * 60
			},
		},
		{
			name: "inside quiet hours",
			now:  time.Date(2026, 11, 11, 12, 30, 0, 0, time.UTC),
			mutate: func(p *storage.NotificationPreferences) {
				p.QuietStartMin = 12 * 60
				p.QuietEndMin = 13 * 60
			},
		},
		{
			name: "paused",
			now:  time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC),
			mutate: func(p *storage.NotificationPreferences) {
				pause := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
				p.PausedUntil = &pause
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, ns, learner := generatorStore(t)
			savePrefs(t, ns, learner.ID, tc.mutate)
			due := &fakeDue{count: 4}
			result, err := generatorAt(ns, due, tc.now, true).RunOnce()
			if err != nil {
				t.Fatalf("run: %v", err)
			}
			if result.OutsideWindow != 1 || result.Notifications != 0 {
				t.Fatalf("result = %+v, want outside window and no notification", result)
			}
			if due.calls != 0 {
				t.Fatalf("due calls = %d, want 0 outside the window", due.calls)
			}
		})
	}
}

func TestGeneratorUsesLearnerTimezone(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.Timezone = "America/New_York"
		p.PreferredTimeMin = 18 * 60
	})
	// 23:30 UTC is 18:30 in New York on 11 November (EST).
	at := time.Date(2026, 11, 11, 23, 30, 0, 0, time.UTC)
	result, err := generatorAt(ns, &fakeDue{count: 2}, at, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 1 {
		t.Fatalf("result = %+v, want one reminder in local evening", result)
	}
	page, _ := ns.ListNotifications(learner.ID, storage.NotificationListOptions{Limit: 5})
	if len(page.Notifications) != 1 || page.Notifications[0].DedupeKey != "review:2026-11-11" {
		t.Fatalf("local dedupe key = %+v", page.Notifications)
	}
}

func TestGeneratorSuppressesUnverifiedEmailButKeepsInApp(t *testing.T) {
	store, ns, learner := generatorStore(t)
	if err := store.SetEmailVerified(learner.ID, false); err != nil {
		t.Fatalf("unverify: %v", err)
	}
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.EmailEnabled = true
	})
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	result, err := generatorAt(ns, &fakeDue{count: 2}, now, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 1 || result.EmailSuppressed != 1 || result.EmailPending != 0 {
		t.Fatalf("result = %+v, want one in-app notification and one suppressed email", result)
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now, time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("email claim = %+v (%v), want none", jobs, err)
	}
}

func TestGeneratorSkipsWhenNoDeliverableChannel(t *testing.T) {
	store, ns, learner := generatorStore(t)
	if err := store.SetEmailVerified(learner.ID, false); err != nil {
		t.Fatalf("unverify: %v", err)
	}
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.InAppEnabled = false
		p.EmailEnabled = true
	})
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	result, err := generatorAt(ns, &fakeDue{count: 2}, now, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.NoChannel != 1 || result.Notifications != 0 {
		t.Fatalf("result = %+v, want no channel", result)
	}
	if page, _ := ns.ListNotifications(learner.ID, storage.NotificationListOptions{Limit: 5}); len(page.Notifications) != 0 {
		t.Fatalf("notifications = %+v, want none", page.Notifications)
	}
}

func TestGeneratorSkipsWhenEmailProviderIsNotConfigured(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.InAppEnabled = false
		p.EmailEnabled = true
	})
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	result, err := generatorAt(ns, &fakeDue{count: 2}, now, false).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.NoChannel != 1 || result.Notifications != 0 {
		t.Fatalf("result = %+v, want no channel when mail is unconfigured", result)
	}
}

func TestGeneratorRecordsCleanDayWithoutCreating(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, nil)
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	due := &fakeDue{count: 0}
	gen := generatorAt(ns, due, now, true)
	result, err := gen.RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.NoReviews != 1 || result.Notifications != 0 {
		t.Fatalf("result = %+v, want a clean day", result)
	}
	if got, _ := ns.GetNotificationScanState(learner.ID, "review"); got != 20261111 {
		t.Fatalf("watermark = %d, want 20261111", got)
	}
	result, err = gen.RunOnce()
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if result.AlreadyGenerated != 1 || due.calls != 1 {
		t.Fatalf("second run = %+v after %d calls, want already generated", result, due.calls)
	}
}

func TestGeneratorReportsDueCounterFailure(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, nil)
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	result, err := generatorAt(ns, &fakeDue{err: errors.New("read failed")}, now, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Failed != 1 || result.Notifications != 0 {
		t.Fatalf("result = %+v, want one failure", result)
	}
	if got, _ := ns.GetNotificationScanState(learner.ID, "review"); got != 0 {
		t.Fatalf("watermark = %d, want 0 after a failed count", got)
	}
}

func TestGeneratorCapsNotificationsPerLocalDay(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, nil)
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	for i := 0; i < maxNotificationsPerLearnerPerDay; i++ {
		if _, err := ns.CreateNotification(&storage.Notification{
			ID:        fmt.Sprintf("n-cap-%d", i),
			StudentID: learner.ID,
			Kind:      "study",
			DedupeKey: fmt.Sprintf("cap:%d", i),
			Title:     "Already sent today",
			CreatedAt: now.Add(-time.Duration(i+1) * time.Minute),
		}); err != nil {
			t.Fatalf("seed notification %d: %v", i, err)
		}
	}

	due := &fakeDue{count: 3}
	result, err := generatorAt(ns, due, now, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 0 || result.Capped != 1 {
		t.Fatalf("result = %+v, want one capped reminder and no notification", result)
	}
	count, err := ns.NotificationCountSince(learner.ID, now.Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("count since: %v", err)
	}
	if count != maxNotificationsPerLearnerPerDay {
		t.Fatalf("notifications today = %d, want %d", count, maxNotificationsPerLearnerPerDay)
	}
}
