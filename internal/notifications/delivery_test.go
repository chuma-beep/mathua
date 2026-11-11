package notifications

import (
	"encoding/json"
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

func TestDeliverySendsOneDigestForABatch(t *testing.T) {
	ns, firstID, now := deliveryFixture(t, true)
	recipients, err := ns.ListNotificationRecipients("", 5)
	if err != nil || len(recipients) == 0 {
		t.Fatalf("list recipients: %v %v", recipients, err)
	}
	studentID := recipients[0].StudentID
	prefs, err := ns.GetNotificationPreferences(studentID)
	if err != nil || prefs == nil {
		t.Fatalf("load preferences: %v", err)
	}
	prefs.Digest = storage.NotificationDigestDaily
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("save digest preference: %v", err)
	}
	second := &storage.Notification{
		ID: "n-email-2", StudentID: studentID, Kind: StudyKind, DedupeKey: "study:2026-11-11",
		Title: "Study session", Body: "Today is a study day.", Href: "/study", CreatedAt: now,
	}
	if created, err := ns.CreateNotification(second); err != nil || !created {
		t.Fatalf("create second notification: %v %v", created, err)
	}
	if inserted, err := ns.CreateNotificationDelivery(&storage.NotificationDelivery{
		NotificationID: second.ID, Channel: storage.NotificationChannelEmail, NextAttemptAt: now,
	}); err != nil || !inserted {
		t.Fatalf("create second delivery: %v %v", inserted, err)
	}

	sender := &fakeSender{}
	result, err := newTestDelivery(ns, sender, now, 5).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(sender.messages) != 1 {
		t.Fatalf("messages = %d, want one digest", len(sender.messages))
	}
	if !strings.Contains(sender.messages[0].Text, "You have 2 notifications waiting.") {
		t.Fatalf("digest did not list both items:\n%s", sender.messages[0].Text)
	}
	if result.Sent != 2 {
		t.Fatalf("sent = %d, want both items counted", result.Sent)
	}
	for _, id := range []string{firstID, second.ID} {
		deliveries, err := ns.ListNotificationDeliveries(id)
		if err != nil || len(deliveries) != 1 {
			t.Fatalf("deliveries for %s: %v %v", id, deliveries, err)
		}
		if deliveries[0].Status != storage.NotificationDeliverySent {
			t.Fatalf("delivery for %s = %s, want sent", id, deliveries[0].Status)
		}
	}
}

func TestPushDigestCarriesANotificationID(t *testing.T) {
	_, ns, learner := generatorStore(t)
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.PushEnabled = true
		p.Digest = storage.NotificationDigestDaily
		p.PreferredTimeMin = 0
		p.QuietStartMin = 0
		p.QuietEndMin = 0
	})
	if err := ns.UpsertNotificationPushSubscription(&storage.NotificationPushSubscription{
		StudentID: learner.ID, Endpoint: "https://push.example/one", P256DH: "key", Auth: "secret",
	}); err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	for _, id := range []string{"n-push-1", "n-push-2"} {
		if _, err := ns.CreateNotification(&storage.Notification{
			ID: id, StudentID: learner.ID, Kind: ReviewKind, DedupeKey: id,
			Title: "Waiting", Body: "A reminder is waiting.", Href: "/review", CreatedAt: now,
		}); err != nil {
			t.Fatalf("create %s: %v", id, err)
		}
		if _, err := ns.CreateNotificationDelivery(&storage.NotificationDelivery{
			NotificationID: id, Channel: storage.NotificationChannelPush,
			NextAttemptAt: now,
		}); err != nil {
			t.Fatalf("create %s delivery: %v", id, err)
		}
	}

	push := &fakePush{}
	delivery := NewDelivery(DeliveryOptions{
		Store: ns, Push: push, Now: func() time.Time { return now }, MaxAttempts: 5,
	})
	result, err := delivery.RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(push.payloads) != 1 {
		t.Fatalf("push payloads = %d, want one digest", len(push.payloads))
	}
	var message PushMessage
	if err := json.Unmarshal(push.payloads[0], &message); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if message.NotificationID == "" {
		t.Fatal("digest push carries no notification id, so a click can never be reported")
	}
	if message.Href != "/notifications" {
		t.Fatalf("digest href = %q, want the archive", message.Href)
	}
	if result.Sent != 2 {
		t.Fatalf("sent = %d, want both items settled", result.Sent)
	}
	if _, err := json.Marshal(message); err != nil {
		t.Fatalf("marshal: %v", err)
	}
}

// Two workers over one database must not both send the same digest. The first
// pass claims and settles the batch; the second finds nothing to claim.
func TestSecondWorkerDoesNotRepeatASettledDigest(t *testing.T) {
	ns, firstID, now := deliveryFixture(t, true)
	recipients, err := ns.ListNotificationRecipients("", 5)
	if err != nil || len(recipients) == 0 {
		t.Fatalf("list recipients: %v %v", recipients, err)
	}
	studentID := recipients[0].StudentID
	prefs, err := ns.GetNotificationPreferences(studentID)
	if err != nil || prefs == nil {
		t.Fatalf("load preferences: %v", err)
	}
	prefs.Digest = storage.NotificationDigestDaily
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("save digest preference: %v", err)
	}
	second := &storage.Notification{
		ID: "n-two-workers", StudentID: studentID, Kind: StudyKind, DedupeKey: "study:two-workers",
		Title: "Study session", Body: "Today is a study day.", Href: "/study", CreatedAt: now,
	}
	if created, err := ns.CreateNotification(second); err != nil || !created {
		t.Fatalf("create second notification: %v %v", created, err)
	}
	if inserted, err := ns.CreateNotificationDelivery(&storage.NotificationDelivery{
		NotificationID: second.ID, Channel: storage.NotificationChannelEmail, NextAttemptAt: now,
	}); err != nil || !inserted {
		t.Fatalf("create second delivery: %v %v", inserted, err)
	}

	firstSender := &fakeSender{}
	first, err := newTestDelivery(ns, firstSender, now, 5).RunOnce()
	if err != nil {
		t.Fatalf("first worker: %v", err)
	}
	if len(firstSender.messages) != 1 || first.Sent != 2 {
		t.Fatalf("first worker sent %d messages and %d rows", len(firstSender.messages), first.Sent)
	}

	secondSender := &fakeSender{}
	secondRun, err := newTestDelivery(ns, secondSender, now, 5).RunOnce()
	if err != nil {
		t.Fatalf("second worker: %v", err)
	}
	if len(secondSender.messages) != 0 || secondRun.Claimed != 0 || secondRun.Sent != 0 {
		t.Fatalf("second worker repeated the digest: messages %d, result %+v", len(secondSender.messages), secondRun)
	}
	for _, id := range []string{firstID, second.ID} {
		rows, err := ns.ListNotificationDeliveries(id)
		if err != nil || len(rows) != 1 || rows[0].Status != storage.NotificationDeliverySent {
			t.Fatalf("delivery for %s = %+v, err %v", id, rows, err)
		}
	}
}
