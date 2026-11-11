package notifications

import (
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestNextDigestTime(t *testing.T) {
	loc := time.UTC
	monday := time.Date(2026, 11, 9, 10, 0, 0, 0, loc)
	cases := []struct {
		name string
		now  time.Time
		mode string
		want time.Time
	}{
		{
			name: "daily before the preferred time keeps today",
			now:  monday,
			mode: storage.NotificationDigestDaily,
			want: time.Date(2026, 11, 9, 18, 0, 0, 0, loc),
		},
		{
			name: "daily after the preferred time moves to tomorrow",
			now:  time.Date(2026, 11, 9, 19, 0, 0, 0, loc),
			mode: storage.NotificationDigestDaily,
			want: time.Date(2026, 11, 10, 18, 0, 0, 0, loc),
		},
		{
			name: "weekly on Monday before the time keeps today",
			now:  monday,
			mode: storage.NotificationDigestWeekly,
			want: time.Date(2026, 11, 9, 18, 0, 0, 0, loc),
		},
		{
			name: "weekly from Tuesday moves to next Monday",
			now:  time.Date(2026, 11, 10, 10, 0, 0, 0, loc),
			mode: storage.NotificationDigestWeekly,
			want: time.Date(2026, 11, 16, 18, 0, 0, 0, loc),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := nextDigestTime(tc.now, tc.mode, 18*60); !got.Equal(tc.want) {
				t.Fatalf("nextDigestTime = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDigestEmailListsEveryItem(t *testing.T) {
	jobs := []storage.NotificationDeliveryJob{
		{NotificationID: "n1", Kind: ReviewKind, Title: "Reviews are due", Body: "Four reviews are waiting.", Href: "/review"},
		{NotificationID: "n2", Kind: StudyKind, Title: "Study session", Body: "Today is one of your study days.", Href: "/study"},
	}
	body := DigestEmail("Ada", jobs, "https://mathua.com", "https://mathua.com/settings#notifications", "https://mathua.com/unsubscribe")
	for _, want := range []string{
		"Hello Ada,",
		"You have 2 notifications waiting.",
		"Reviews are due",
		"https://mathua.com/review",
		"Study session",
		"https://mathua.com/study",
		"Notification settings:",
		"Unsubscribe from email notifications:",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("digest body is missing %q:\n%s", want, body)
		}
	}
}

func TestGeneratorDefersDigestDeliveries(t *testing.T) {
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.EmailEnabled = true
		p.Digest = storage.NotificationDigestDaily
		p.PreferredTimeMin = 18 * 60
	})
	now := time.Date(2026, 11, 11, 19, 0, 0, 0, time.UTC)
	generator := generatorAt(ns, &fakeDue{count: 3}, now, true)
	if _, err := generator.RunOnce(); err != nil {
		t.Fatalf("generate: %v", err)
	}

	page, err := ns.ListNotifications(learner.ID, storage.NotificationListOptions{})
	if err != nil {
		t.Fatalf("list notifications: %v", err)
	}
	if len(page.Notifications) != 1 {
		t.Fatalf("notifications = %d, want one", len(page.Notifications))
	}
	deliveries, err := ns.ListNotificationDeliveries(page.Notifications[0].ID)
	if err != nil {
		t.Fatalf("list deliveries: %v", err)
	}
	var email *storage.NotificationDelivery
	for i := range deliveries {
		if deliveries[i].Channel == storage.NotificationChannelEmail {
			email = &deliveries[i]
		}
	}
	if email == nil {
		t.Fatal("email delivery was not created")
	}
	want := time.Date(2026, 11, 12, 18, 0, 0, 0, time.UTC)
	if !email.NextAttemptAt.Equal(want) {
		t.Fatalf("digest delivery is due at %v, want %v", email.NextAttemptAt, want)
	}
}

func TestNormalizePreferencesDigest(t *testing.T) {
	base := storage.DefaultNotificationPreferences("st-1", time.Now())
	base.Digest = ""
	normalized, err := NormalizePreferencesAt(base, time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("normalize empty digest: %v", err)
	}
	if normalized.Digest != storage.NotificationDigestImmediate {
		t.Fatalf("empty digest = %q, want immediate", normalized.Digest)
	}
	base.Digest = "hourly"
	if _, err := NormalizePreferencesAt(base, time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)); err == nil {
		t.Fatal("unknown digest was accepted")
	}
	base.Digest = "weekly"
	normalized, err = NormalizePreferencesAt(base, time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC))
	if err != nil || normalized.Digest != storage.NotificationDigestWeekly {
		t.Fatalf("weekly digest = %q, %v", normalized.Digest, err)
	}
}
