package server

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestPushKeyReportsConfiguration(t *testing.T) {
	s, mux, _, token, _ := setupTestServerWithPlanner(t)

	rec := notificationRequest(t, mux, http.MethodGet, "/api/notifications/push/key", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("key = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Configured bool   `json:"configured"`
		PublicKey  string `json:"public_key"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Configured || body.PublicKey != "" {
		t.Fatalf("key without a VAPID pair = %+v, want unconfigured", body)
	}

	s.pushPublicKey = "BPublicKeyExample"
	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications/push/key", token, "")
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !body.Configured || body.PublicKey != "BPublicKeyExample" {
		t.Fatalf("key = %+v, want the configured public key", body)
	}
}

func TestPushSubscribeStoresDeviceAndEnablesChannel(t *testing.T) {
	s, mux, store, token, studentID := setupTestServerWithPlanner(t)
	s.pushPublicKey = "BPublicKeyExample"

	payload := `{"endpoint":"https://push.example/device-1","keys":{"p256dh":"p256dh-value","auth":"auth-value"}}`
	rec := notificationRequest(t, mux, http.MethodPost, "/api/notifications/push/subscribe", token, payload)
	if rec.Code != http.StatusOK {
		t.Fatalf("subscribe = %d: %s", rec.Code, rec.Body.String())
	}

	subs, err := store.(storage.NotificationStore).ListNotificationPushSubscriptions(studentID)
	if err != nil || len(subs) != 1 {
		t.Fatalf("subscriptions = %+v (%v), want one", subs, err)
	}
	if subs[0].Endpoint != "https://push.example/device-1" || subs[0].P256DH != "p256dh-value" {
		t.Fatalf("subscription = %+v", subs[0])
	}
	prefs, err := store.(storage.NotificationStore).GetNotificationPreferences(studentID)
	if err != nil || prefs == nil || !prefs.PushEnabled {
		t.Fatalf("prefs = %+v (%v), want push enabled", prefs, err)
	}

	// Re-subscribing with the same endpoint updates the row instead of adding a
	// second device.
	again := `{"endpoint":"https://push.example/device-1","keys":{"p256dh":"rotated","auth":"rotated"}}`
	rec = notificationRequest(t, mux, http.MethodPost, "/api/notifications/push/subscribe", token, again)
	if rec.Code != http.StatusOK {
		t.Fatalf("resubscribe = %d: %s", rec.Code, rec.Body.String())
	}
	subs, _ = store.(storage.NotificationStore).ListNotificationPushSubscriptions(studentID)
	if len(subs) != 1 || subs[0].P256DH != "rotated" {
		t.Fatalf("subscriptions after resubscribe = %+v", subs)
	}
}

func TestPushSubscribeRejectsBadEndpoints(t *testing.T) {
	s, mux, _, token, _ := setupTestServerWithPlanner(t)
	s.pushPublicKey = "BPublicKeyExample"

	for name, body := range map[string]string{
		"not a url":     `{"endpoint":"not-a-url","keys":{"p256dh":"p","auth":"a"}}`,
		"insecure":      `{"endpoint":"http://push.example/device","keys":{"p256dh":"p","auth":"a"}}`,
		"no keys":       `{"endpoint":"https://push.example/device","keys":{}}`,
		"empty payload": `{}`,
	} {
		t.Run(name, func(t *testing.T) {
			rec := notificationRequest(t, mux, http.MethodPost, "/api/notifications/push/subscribe", token, body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("subscribe = %d, want 400", rec.Code)
			}
		})
	}
}

func TestPushSubscribeUnavailableWithoutKeys(t *testing.T) {
	_, mux, _, token, _ := setupTestServerWithPlanner(t)
	body := `{"endpoint":"https://push.example/device","keys":{"p256dh":"p","auth":"a"}}`
	rec := notificationRequest(t, mux, http.MethodPost, "/api/notifications/push/subscribe", token, body)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("subscribe = %d, want 503", rec.Code)
	}
}

func TestPushUnsubscribeRemovesOnlyOwnDevice(t *testing.T) {
	s, mux, store, token, studentID := setupTestServerWithPlanner(t)
	s.pushPublicKey = "BPublicKeyExample"
	notif := store.(storage.NotificationStore)

	for _, endpoint := range []string{"https://push.example/one", "https://push.example/two"} {
		if err := notif.UpsertNotificationPushSubscription(&storage.NotificationPushSubscription{
			StudentID: studentID, Endpoint: endpoint, P256DH: "p", Auth: "a",
		}); err != nil {
			t.Fatalf("seed device: %v", err)
		}
	}

	rec := notificationRequest(t, mux, http.MethodPost, "/api/notifications/push/unsubscribe", token,
		`{"endpoint":"https://push.example/one"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("unsubscribe = %d: %s", rec.Code, rec.Body.String())
	}
	subs, _ := notif.ListNotificationPushSubscriptions(studentID)
	if len(subs) != 1 || subs[0].Endpoint != "https://push.example/two" {
		t.Fatalf("subscriptions = %+v, want only the second device", subs)
	}
	prefs, _ := notif.GetNotificationPreferences(studentID)
	if prefs == nil || prefs.PushEnabled {
		t.Fatalf("prefs = %+v, want push disabled", prefs)
	}
}

func TestAdminNotificationsReportsDeliveryHealth(t *testing.T) {
	e := newAdminEnv(t)
	notif := e.store.(storage.NotificationStore)
	now := time.Now().UTC().Truncate(time.Second)
	n := &storage.Notification{
		ID: "n-health", StudentID: e.studentID, Kind: "review", DedupeKey: "review:2026-11-11",
		Title: "Reviews are waiting", Body: "Two reviews are due.", Href: "/review", CreatedAt: now,
	}
	if created, err := notif.CreateNotification(n); err != nil || !created {
		t.Fatalf("seed notification: %v %v", created, err)
	}
	if inserted, err := notif.CreateNotificationDelivery(&storage.NotificationDelivery{
		NotificationID: n.ID,
		Channel:        storage.NotificationChannelEmail,
		Status:         storage.NotificationDeliveryFailed,
		Attempts:       3,
		LastError:      "provider unavailable",
		NextAttemptAt:  now,
	}); err != nil || !inserted {
		t.Fatalf("seed delivery: %v %v", inserted, err)
	}
	sent := &storage.Notification{
		ID: "n-health-sent", StudentID: e.studentID, Kind: "study", DedupeKey: "study:2026-11-11",
		Title: "Study today", Body: "A short session is waiting.", Href: "/study", CreatedAt: now,
	}
	if created, err := notif.CreateNotification(sent); err != nil || !created {
		t.Fatalf("seed sent notification: %v %v", created, err)
	}
	if inserted, err := notif.CreateNotificationDelivery(&storage.NotificationDelivery{
		NotificationID: sent.ID,
		Channel:        storage.NotificationChannelInApp,
		Status:         storage.NotificationDeliverySent,
		SentAt:         now,
	}); err != nil || !inserted {
		t.Fatalf("seed sent delivery: %v %v", inserted, err)
	}
	if err := notif.MarkNotificationOpened(e.studentID, sent.ID, now.Add(time.Minute)); err != nil {
		t.Fatalf("mark opened: %v", err)
	}

	rec := e.do(t, http.MethodGet, "/api/admin/notifications", e.adminTok, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("admin notifications = %d: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		Stats []struct {
			Channel string `json:"channel"`
			Status  string `json:"status"`
			Count   int    `json:"count"`
		} `json:"stats"`
		Engagement []struct {
			Kind    string `json:"kind"`
			Channel string `json:"channel"`
			Sent    int    `json:"sent"`
			Opened  int    `json:"opened"`
			Clicked int    `json:"clicked"`
		} `json:"engagement"`
		Problems []struct {
			StudentID string `json:"student_id"`
			Channel   string `json:"channel"`
			Status    string `json:"status"`
			Attempts  int    `json:"attempts"`
			LastError string `json:"last_error"`
		} `json:"problems"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	foundFailed := false
	for _, stat := range body.Stats {
		if stat.Channel == "email" && stat.Status == "failed" && stat.Count == 1 {
			foundFailed = true
		}
	}
	if !foundFailed {
		t.Fatalf("stats = %+v, want one failed email row", body.Stats)
	}
	if len(body.Problems) != 1 || body.Problems[0].StudentID != e.studentID ||
		body.Problems[0].Attempts != 3 || body.Problems[0].LastError != "provider unavailable" {
		t.Fatalf("problems = %+v", body.Problems)
	}
	// Only sent rows count, so the failed email delivery is absent and the
	// opened in-app one is the single engagement row.
	if len(body.Engagement) != 1 || body.Engagement[0].Kind != "study" ||
		body.Engagement[0].Channel != "in_app" || body.Engagement[0].Sent != 1 ||
		body.Engagement[0].Opened != 1 || body.Engagement[0].Clicked != 0 {
		t.Fatalf("engagement = %+v", body.Engagement)
	}
}
