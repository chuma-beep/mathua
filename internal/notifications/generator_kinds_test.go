package notifications

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// weekdayClock is a Wednesday at noon UTC. 2026-11-11 is chosen so the study
// day test can name a day rather than compute an index.
func weekdayClock() time.Time {
	return time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
}

func TestGeneratorCreatesOneStudyReminderOnAChosenDay(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.ReviewsEnabled = false
		p.StudyEnabled = true
		p.StudyDays = "wed"
		p.StudyTimeMin = 0
	})
	now := weekdayClock()
	result, err := generatorAt(ns, &fakeDue{}, now, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 1 || result.StudyReminders != 1 {
		t.Fatalf("result = %+v, want one study reminder", result)
	}
	page, err := ns.ListNotifications(learner.ID, storage.NotificationListOptions{Limit: 5})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(page.Notifications) != 1 {
		t.Fatalf("notifications = %d, want 1", len(page.Notifications))
	}
	n := page.Notifications[0]
	if n.Kind != StudyKind || n.DedupeKey != "study:2026-11-11" || n.Href != "/study" {
		t.Fatalf("notification = %+v", n)
	}

	again, err := generatorAt(ns, &fakeDue{}, now, true).RunOnce()
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if again.Notifications != 0 || again.AlreadyGenerated != 1 {
		t.Fatalf("second run = %+v, want already generated", again)
	}
}

func TestGeneratorSkipsStudyReminderOnAnotherDay(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.ReviewsEnabled = false
		p.StudyEnabled = true
		p.StudyDays = "mon,tue"
	})
	result, err := generatorAt(ns, &fakeDue{}, weekdayClock(), true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 0 || result.StudyReminders != 0 {
		t.Fatalf("result = %+v, want no study reminder", result)
	}
	if got, _ := ns.GetNotificationScanState(learner.ID, "study"); got != 0 {
		t.Fatalf("study watermark = %d, want none", got)
	}
}

func TestGeneratorAnnouncesOnlyTheHighestMilestone(t *testing.T) {
	store, ns, learner := generatorStore(t)
	if err := store.AwardXP(learner.ID, 1200, "test", "milestone-fixture"); err != nil {
		t.Fatalf("award xp: %v", err)
	}
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.ReviewsEnabled = false
	})
	result, err := generatorAt(ns, &fakeDue{}, weekdayClock(), true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 1 || result.Milestones != 1 {
		t.Fatalf("result = %+v, want one milestone", result)
	}
	page, _ := ns.ListNotifications(learner.ID, storage.NotificationListOptions{Limit: 10})
	if len(page.Notifications) != 1 {
		t.Fatalf("notifications = %d, want 1", len(page.Notifications))
	}
	n := page.Notifications[0]
	if n.Kind != MilestoneKind || n.DedupeKey != "milestone:xp:1000" {
		t.Fatalf("notification = %+v, want the highest crossed threshold", n)
	}
	if n.Title != "1000 XP reached" {
		t.Fatalf("title = %q", n.Title)
	}
}

func TestGeneratorNudgesInactivityAtMostOncePerWeek(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.ReviewsEnabled = false
		p.MilestonesEnabled = false
		p.InactivityEnabled = true
	})
	// The account was created on 11 November; a month later the learner has
	// not studied, so one nudge is due.
	first := time.Date(2026, 12, 11, 12, 0, 0, 0, time.UTC)
	result, err := generatorAt(ns, &fakeDue{}, first, true).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Inactivity != 1 || result.Notifications != 1 {
		t.Fatalf("result = %+v, want one nudge", result)
	}

	sameDay, err := generatorAt(ns, &fakeDue{}, first, true).RunOnce()
	if err != nil {
		t.Fatalf("same day: %v", err)
	}
	if sameDay.Inactivity != 0 {
		t.Fatalf("same day = %+v, want no second nudge", sameDay)
	}

	soon, err := generatorAt(ns, &fakeDue{}, first.AddDate(0, 0, 3), true).RunOnce()
	if err != nil {
		t.Fatalf("three days later: %v", err)
	}
	if soon.Inactivity != 0 {
		t.Fatalf("three days later = %+v, want the gap enforced", soon)
	}

	week, err := generatorAt(ns, &fakeDue{}, first.AddDate(0, 0, 7), true).RunOnce()
	if err != nil {
		t.Fatalf("a week later: %v", err)
	}
	if week.Inactivity != 1 {
		t.Fatalf("a week later = %+v, want another nudge", week)
	}
}

func TestGeneratorCreatesPushDeliveryWhenConfigured(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.InAppEnabled = false
		p.ReviewsEnabled = false
		p.MilestonesEnabled = false
		p.StudyEnabled = true
		p.StudyDays = "wed"
		p.StudyTimeMin = 0
		p.PushEnabled = true
	})
	now := weekdayClock()
	gen := NewGenerator(GeneratorOptions{
		Store:         ns,
		Due:           &fakeDue{},
		PushAvailable: true,
		Now:           func() time.Time { return now },
	})
	result, err := gen.RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Notifications != 1 || result.PushPending != 1 {
		t.Fatalf("result = %+v, want one pending push", result)
	}
	page, _ := ns.ListNotifications(learner.ID, storage.NotificationListOptions{Limit: 5})
	jobs, err := ns.ClaimNotificationDeliveries(now, time.Minute, 10)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("push jobs = %d, want 1", len(jobs))
	}
	if jobs[0].Channel != storage.NotificationChannelPush || jobs[0].StudentID != learner.ID {
		t.Fatalf("push job = %+v", jobs[0])
	}
	if len(page.Notifications) != 1 || page.Notifications[0].Kind != StudyKind {
		t.Fatalf("notification = %+v", page.Notifications)
	}
	if jobs[0].Kind != StudyKind {
		t.Fatalf("job kind = %q, want %q", jobs[0].Kind, StudyKind)
	}
}
