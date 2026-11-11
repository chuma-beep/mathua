package notifications

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

type fakePush struct {
	sent      []storage.NotificationPushSubscription
	payloads  [][]byte
	errFor    map[string]error
	lastError error
}

func (f *fakePush) Send(sub storage.NotificationPushSubscription, payload []byte) error {
	f.sent = append(f.sent, sub)
	f.payloads = append(f.payloads, payload)
	if f.errFor != nil {
		if err, ok := f.errFor[sub.Endpoint]; ok {
			f.lastError = err
			return err
		}
	}
	return nil
}

func pushFixture(t *testing.T) (storage.NotificationStore, *storage.SQLiteStore, string, time.Time) {
	t.Helper()
	store, ns, learner := generatorStore(t)
	savePrefs(t, ns, learner.ID, func(p *storage.NotificationPreferences) {
		p.InAppEnabled = false
		p.EmailEnabled = false
		p.PushEnabled = true
		p.PreferredTimeMin = 0
		p.QuietStartMin = 0
		p.QuietEndMin = 0
	})
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	n := &storage.Notification{
		ID: "n-push", StudentID: learner.ID, Kind: StudyKind, DedupeKey: "study:2026-11-11",
		Title: "Today is a study day", Body: "You chose today for study.", Href: "/study", CreatedAt: now,
	}
	if created, err := ns.CreateNotification(n); err != nil || !created {
		t.Fatalf("create notification: %v %v", created, err)
	}
	if inserted, err := ns.CreateNotificationDelivery(&storage.NotificationDelivery{
		NotificationID: n.ID,
		Channel:        storage.NotificationChannelPush,
		NextAttemptAt:  now,
	}); err != nil || !inserted {
		t.Fatalf("create delivery: %v %v", inserted, err)
	}
	return ns, store, learner.ID, now
}

func addDevice(t *testing.T, ns storage.NotificationStore, studentID, endpoint string) {
	t.Helper()
	if err := ns.UpsertNotificationPushSubscription(&storage.NotificationPushSubscription{
		StudentID: studentID,
		Endpoint:  endpoint,
		P256DH:    "p256dh-" + endpoint,
		Auth:      "auth-" + endpoint,
		UserAgent: "test",
	}); err != nil {
		t.Fatalf("upsert device %s: %v", endpoint, err)
	}
}

func newPushDelivery(ns storage.NotificationStore, push PushSender, now time.Time) *Delivery {
	return NewDelivery(DeliveryOptions{
		Store: ns,
		Push:  push,
		Now:   func() time.Time { return now },
	})
}

func TestDeliverySendsPushToEveryDeviceOnce(t *testing.T) {
	ns, _, studentID, now := pushFixture(t)
	addDevice(t, ns, studentID, "https://push.example/one")
	addDevice(t, ns, studentID, "https://push.example/two")
	push := &fakePush{}

	result, err := newPushDelivery(ns, push, now).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Sent != 1 || result.Failed != 0 || result.Retried != 0 {
		t.Fatalf("result = %+v, want one delivery row sent", result)
	}
	if len(push.sent) != 2 {
		t.Fatalf("devices reached = %d, want 2", len(push.sent))
	}
	var msg PushMessage
	if err := json.Unmarshal(push.payloads[0], &msg); err != nil {
		t.Fatalf("payload: %v", err)
	}
	if msg.Title != "Today is a study day" || msg.Href != "/study" || msg.NotificationID != "n-push" {
		t.Fatalf("payload = %+v", msg)
	}
	subs, err := ns.ListNotificationPushSubscriptions(studentID)
	if err != nil || len(subs) != 2 {
		t.Fatalf("subscriptions = %+v (%v)", subs, err)
	}
	for _, sub := range subs {
		if sub.Failures != 0 || sub.LastSuccess.IsZero() {
			t.Fatalf("subscription health = %+v, want a recorded success", sub)
		}
	}
	// The single delivery row is complete, so nothing is claimable afterwards.
	if jobs, err := ns.ClaimNotificationDeliveries(now.Add(time.Minute), time.Minute, 10); err != nil || len(jobs) != 0 {
		t.Fatalf("claim after send = %+v (%v), want none", jobs, err)
	}
}

func TestDeliveryPrunesGoneEndpoints(t *testing.T) {
	ns, _, studentID, now := pushFixture(t)
	addDevice(t, ns, studentID, "https://push.example/gone")
	addDevice(t, ns, studentID, "https://push.example/live")
	push := &fakePush{errFor: map[string]error{
		"https://push.example/gone": &PushGoneError{Status: 410, Endpoint: "https://push.example/gone"},
	}}

	result, err := newPushDelivery(ns, push, now).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Sent != 1 {
		t.Fatalf("result = %+v, want sent because one device accepted", result)
	}
	subs, err := ns.ListNotificationPushSubscriptions(studentID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(subs) != 1 || subs[0].Endpoint != "https://push.example/live" {
		t.Fatalf("subscriptions = %+v, want only the live endpoint", subs)
	}
}

func TestDeliveryRetriesWhenNoDeviceAccepts(t *testing.T) {
	ns, _, studentID, now := pushFixture(t)
	addDevice(t, ns, studentID, "https://push.example/one")
	push := &fakePush{errFor: map[string]error{
		"https://push.example/one": errors.New("push service is down"),
	}}

	result, err := newPushDelivery(ns, push, now).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Retried != 1 || result.Sent != 0 || result.Failed != 0 {
		t.Fatalf("result = %+v, want a retry", result)
	}
	subs, _ := ns.ListNotificationPushSubscriptions(studentID)
	if len(subs) != 1 || subs[0].Failures != 1 {
		t.Fatalf("subscriptions = %+v, want one recorded failure", subs)
	}
}

func TestDeliverySuppressesPushWithoutADevice(t *testing.T) {
	ns, _, _, now := pushFixture(t)
	push := &fakePush{}

	result, err := newPushDelivery(ns, push, now).RunOnce()
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Suppressed != 1 || len(push.sent) != 0 {
		t.Fatalf("result = %+v with %d sends, want suppressed", result, len(push.sent))
	}
}
