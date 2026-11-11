package notifications

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/mailer"
	"github.com/chuma-beep/mathua/internal/storage"
)

type fakeSender struct {
	messages []mailer.Message
	err      error
}

func (f *fakeSender) Send(msg mailer.Message) error {
	f.messages = append(f.messages, msg)
	return f.err
}

func deliveryFixture(t *testing.T, emailEnabled bool) (storage.NotificationStore, string, time.Time) {
	t.Helper()
	_, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.EmailEnabled = emailEnabled
		p.PreferredTimeMin = 0
		p.QuietStartMin = 0
		p.QuietEndMin = 0
	})
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	n := &storage.Notification{
		ID: "n-email", StudentID: learner.ID, Kind: ReviewKind, DedupeKey: "review:2026-11-11",
		Title: "Reviews are waiting", Body: "You have 3 due reviews.", Href: "/review", CreatedAt: now,
	}
	if created, err := ns.CreateNotification(n); err != nil || !created {
		t.Fatalf("create notification: %v %v", created, err)
	}
	if inserted, err := ns.CreateNotificationDelivery(&storage.NotificationDelivery{
		NotificationID: n.ID,
		Channel:        storage.NotificationChannelEmail,
		NextAttemptAt:  now,
	}); err != nil || !inserted {
		t.Fatalf("create delivery: %v %v", inserted, err)
	}
	return ns, n.ID, now
}

func newTestDelivery(ns storage.NotificationStore, sender mailer.Sender, now time.Time, max int) *Delivery {
	return NewDelivery(DeliveryOptions{
		Store:       ns,
		Sender:      sender,
		FrontendURL: "https://mathua.test/",
		UnsubscribeToken: func(studentID string) (string, error) {
			return "token-" + studentID, nil
		},
		Now:         func() time.Time { return now },
		MaxAttempts: max,
	})
}

func TestDeliverySendsReviewReminder(t *testing.T) {
	ns, notificationID, now := deliveryFixture(t, true)
	sender := &fakeSender{}
	result, err := newTestDelivery(ns, sender, now, 5).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Sent != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v, want one sent", result)
	}
	if len(sender.messages) != 1 {
		t.Fatalf("messages = %d, want 1", len(sender.messages))
	}
	msg := sender.messages[0]
	if msg.To != "ada@example.com" || msg.Subject != "Reviews are waiting" {
		t.Fatalf("message = %+v", msg)
	}
	for _, want := range []string{
		"Hello Ada Example,",
		"You have 3 due reviews.",
		"https://mathua.test/review",
		"https://mathua.test/settings#notifications",
		"https://mathua.test/notifications/unsubscribe?token=token-",
	} {
		if !strings.Contains(msg.Text, want) {
			t.Errorf("email missing %q:\n%s", want, msg.Text)
		}
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now.Add(time.Minute), time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("claim after send = %+v (%v), want none for %s", jobs, err, notificationID)
	}
}

func TestDeliveryRetriesThenFails(t *testing.T) {
	ns, _, now := deliveryFixture(t, true)
	sender := &fakeSender{err: errors.New("provider unavailable")}
	worker := newTestDelivery(ns, sender, now, 2)

	result, err := worker.RunOnce()
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	if result.Retried != 1 || len(sender.messages) != 1 {
		t.Fatalf("first run = %+v with %d messages", result, len(sender.messages))
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now.Add(time.Minute), time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("early claim = %+v (%v), want none before backoff", jobs, err)
	}
	// The second claim happens inside the worker: attempt 2 reaches the cap,
	// so the row becomes failed instead of retried again.
	worker = newTestDelivery(ns, sender, now.Add(5*time.Minute), 2)
	result, err = worker.RunOnce()
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if result.Claimed != 1 || result.Failed != 1 {
		t.Fatalf("second run = %+v, want one job failed at the attempt cap", result)
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now.Add(24*time.Hour), time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("claim after failure = %+v (%v), want none", jobs, err)
	}
}

func TestDeliverySuppressesWhenEmailWasTurnedOff(t *testing.T) {
	ns, _, now := deliveryFixture(t, false)
	sender := &fakeSender{}
	result, err := newTestDelivery(ns, sender, now, 5).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Suppressed != 1 || len(sender.messages) != 0 {
		t.Fatalf("result = %+v with %d messages, want suppressed", result, len(sender.messages))
	}
}

func TestDeliverySuppressesUnverifiedAddress(t *testing.T) {
	ns, _, now := deliveryFixture(t, true)
	store := ns.(*storage.SQLiteStore)
	studentID := mustStudentID(t, ns)
	if err := store.SetEmail(studentID, "ada@example.com"); err != nil {
		t.Fatalf("set email: %v", err)
	}
	if err := store.SetEmailVerified(studentID, false); err != nil {
		t.Fatalf("unverify email: %v", err)
	}
	sender := &fakeSender{}
	result, err := newTestDelivery(ns, sender, now, 5).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Suppressed != 1 || len(sender.messages) != 0 {
		t.Fatalf("result = %+v with %d messages, want suppressed", result, len(sender.messages))
	}
}

func mustStudentID(t *testing.T, ns storage.NotificationStore) string {
	t.Helper()
	recipients, err := ns.ListNotificationRecipients("", 10)
	if err != nil || len(recipients) == 0 {
		t.Fatalf("recipients = %+v (%v)", recipients, err)
	}
	return recipients[0].StudentID
}

func TestDeliveryDefersOutsideWindow(t *testing.T) {
	ns, _, now := deliveryFixture(t, true)
	prefs, err := ns.GetNotificationPreferences(mustStudentID(t, ns))
	if err != nil || prefs == nil {
		t.Fatalf("prefs: %v %+v", err, prefs)
	}
	prefs.QuietStartMin = 12 * 60
	prefs.QuietEndMin = 13 * 60
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("save quiet hours: %v", err)
	}
	sender := &fakeSender{}
	worker := newTestDelivery(ns, sender, now, 5)
	result, err := worker.RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Retried != 1 || len(sender.messages) != 0 {
		t.Fatalf("quiet run = %+v with %d messages, want deferred", result, len(sender.messages))
	}
	if jobs, err := ns.ClaimNotificationDeliveries(now.Add(30*time.Minute), time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("claim during quiet = %+v (%v), want none", jobs, err)
	}
	afterQuiet := time.Date(2026, 11, 11, 13, 0, 0, 0, time.UTC)
	jobs, err := ns.ClaimNotificationDeliveries(afterQuiet, time.Minute, 10)
	if err != nil || len(jobs) != 1 {
		t.Fatalf("claim after quiet = %+v (%v), want one", jobs, err)
	}
}

func TestDeliveryFailsWithoutSender(t *testing.T) {
	ns, _, now := deliveryFixture(t, true)
	worker := NewDelivery(DeliveryOptions{Store: ns, Now: func() time.Time { return now }})
	if _, err := worker.RunOnce(); err == nil {
		t.Fatal("delivery ran without a sender")
	}
}
