package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/storage"
)

func seedNotification(t *testing.T, store storage.Repository, studentID, id, dedupe string, createdAt time.Time, read bool) {
	t.Helper()
	notif := store.(storage.NotificationStore)
	n := &storage.Notification{
		ID:        id,
		StudentID: studentID,
		Kind:      "review",
		DedupeKey: dedupe,
		Title:     "Review reminder",
		Body:      "You have reviews waiting.",
		Href:      "/review",
		CreatedAt: createdAt,
	}
	if read {
		n.ReadAt = createdAt.Add(time.Minute)
	}
	if _, err := notif.CreateNotification(n); err != nil {
		t.Fatalf("seed notification %s: %v", id, err)
	}
}

func notificationRequest(t *testing.T, mux *http.ServeMux, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestNotificationsRequireAuthentication(t *testing.T) {
	_, mux, _, _, _ := setupTestServerWithPlanner(t)
	for _, tc := range []struct{ method, path, body string }{
		{http.MethodGet, "/api/notifications", ""},
		{http.MethodGet, "/api/notifications/unread-count", ""},
		{http.MethodPost, "/api/notifications/read", `{"id":"n1"}`},
		{http.MethodPost, "/api/notifications/read-all", ""},
		{http.MethodGet, "/api/notifications/preferences", ""},
		{http.MethodPut, "/api/notifications/preferences", `{}`},
	} {
		rec := notificationRequest(t, mux, tc.method, tc.path, "", tc.body)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s %s without auth = %d, want 401", tc.method, tc.path, rec.Code)
		}
	}
}

func TestNotificationsOwnershipPaginationAndReadState(t *testing.T) {
	_, mux, store, token, sid := setupTestServerWithPlanner(t)
	otherToken, otherSID := quizGuest(t, mux)

	now := time.Now().UTC().Truncate(time.Second)
	seedNotification(t, store, sid, "n-old", "old", now.Add(-3*time.Hour), false)
	seedNotification(t, store, sid, "n-mid", "mid", now.Add(-2*time.Hour), true)
	seedNotification(t, store, sid, "n-new", "new", now.Add(-time.Hour), false)
	seedNotification(t, store, otherSID, "other", "other", now.Add(-30*time.Minute), false)

	rec := notificationRequest(t, mux, http.MethodGet, "/api/notifications?limit=2", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list = %d: %s", rec.Code, rec.Body.String())
	}
	var page struct {
		Notifications []notificationJSON `json:"notifications"`
		Unread        int                `json:"unread"`
		HasMore       bool               `json:"has_more"`
		NextCursor    string             `json:"next_cursor"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(page.Notifications) != 2 || page.Notifications[0].ID != "n-new" || page.Notifications[1].ID != "n-mid" {
		t.Fatalf("first page = %+v, want n-new then n-mid", page.Notifications)
	}
	if !page.HasMore || page.NextCursor == "" {
		t.Fatalf("first page should continue: %+v", page)
	}
	if page.Unread != 2 {
		t.Errorf("unread = %d, want 2", page.Unread)
	}
	for _, n := range page.Notifications {
		if n.ID == "other" {
			t.Fatal("another learner's notification leaked into the feed")
		}
	}

	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications?limit=2&cursor="+page.NextCursor, token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("next page = %d: %s", rec.Code, rec.Body.String())
	}
	var next struct {
		Notifications []notificationJSON `json:"notifications"`
		HasMore       bool               `json:"has_more"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &next); err != nil {
		t.Fatalf("decode next page: %v", err)
	}
	if len(next.Notifications) != 1 || next.Notifications[0].ID != "n-old" || next.HasMore {
		t.Fatalf("next page = %+v, want n-old and no continuation", next)
	}

	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications?unread=true", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("unread feed = %d: %s", rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode unread feed: %v", err)
	}
	if len(page.Notifications) != 2 {
		t.Fatalf("unread feed = %+v, want two unread", page.Notifications)
	}
	for _, n := range page.Notifications {
		if n.Read {
			t.Errorf("unread feed returned read notification %s", n.ID)
		}
	}

	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications/unread-count", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("unread count = %d: %s", rec.Code, rec.Body.String())
	}
	var count struct {
		Unread int `json:"unread"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &count); err != nil || count.Unread != 2 {
		t.Fatalf("unread count = %+v (%v), want 2", count, err)
	}

	// Another learner's id is not found, and it must not change that learner's
	// unread state.
	rec = notificationRequest(t, mux, http.MethodPost, "/api/notifications/read", token, `{"id":"other"}`)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("reading another learner's notification = %d, want 404", rec.Code)
	}
	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications/unread-count", otherToken, "")
	if err := json.Unmarshal(rec.Body.Bytes(), &count); err != nil || count.Unread != 1 {
		t.Fatalf("other unread after foreign read = %+v (%v), want 1", count, err)
	}

	rec = notificationRequest(t, mux, http.MethodPost, "/api/notifications/read", token, `{"id":"n-new"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("mark read = %d: %s", rec.Code, rec.Body.String())
	}
	var mark struct {
		OK     bool `json:"ok"`
		Unread int  `json:"unread"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &mark); err != nil || !mark.OK || mark.Unread != 1 {
		t.Fatalf("mark read response = %+v (%v), want ok and unread 1", mark, err)
	}

	// Repeating the same read is idempotent.
	rec = notificationRequest(t, mux, http.MethodPost, "/api/notifications/read", token, `{"id":"n-new"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("repeat mark read = %d: %s", rec.Code, rec.Body.String())
	}

	rec = notificationRequest(t, mux, http.MethodPost, "/api/notifications/read-all", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("mark all = %d: %s", rec.Code, rec.Body.String())
	}
	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications/unread-count", token, "")
	if err := json.Unmarshal(rec.Body.Bytes(), &count); err != nil || count.Unread != 0 {
		t.Fatalf("unread after mark all = %+v (%v), want 0", count, err)
	}
	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications/unread-count", otherToken, "")
	if err := json.Unmarshal(rec.Body.Bytes(), &count); err != nil || count.Unread != 1 {
		t.Fatalf("other unread after mark all = %+v (%v), want 1", count, err)
	}
}

func TestNotificationFeedHonoursInAppPreference(t *testing.T) {
	_, mux, store, token, sid := setupTestServerWithPlanner(t)
	now := time.Now().UTC().Truncate(time.Second)
	seedNotification(t, store, sid, "n1", "one", now, false)

	ns := store.(storage.NotificationStore)
	prefs, err := ns.EnsureNotificationPreferences(sid, now)
	if err != nil || prefs == nil {
		t.Fatalf("ensure prefs: %v %+v", err, prefs)
	}
	prefs.InAppEnabled = false
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("save prefs: %v", err)
	}

	rec := notificationRequest(t, mux, http.MethodGet, "/api/notifications", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list = %d: %s", rec.Code, rec.Body.String())
	}
	var page struct {
		Notifications []notificationJSON `json:"notifications"`
		Unread        int                `json:"unread"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(page.Notifications) != 0 || page.Unread != 0 {
		t.Fatalf("feed with in-app off = %+v, want hidden and zero", page)
	}
	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications/unread-count", token, "")
	var count struct {
		Unread int `json:"unread"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &count); err != nil || count.Unread != 0 {
		t.Fatalf("count with in-app off = %+v (%v), want 0", count, err)
	}

	prefs.InAppEnabled = true
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("restore prefs: %v", err)
	}
	rec = notificationRequest(t, mux, http.MethodGet, "/api/notifications", token, "")
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode restored: %v", err)
	}
	if len(page.Notifications) != 1 || page.Unread != 1 {
		t.Fatalf("feed with in-app on = %+v, want the hidden item back", page)
	}
}

func TestNotificationUnsubscribeIsSignedAndOnlyDisablesEmail(t *testing.T) {
	_, mux, store, _, sid := setupTestServerWithPlanner(t)
	now := time.Now().UTC()
	ns := store.(storage.NotificationStore)
	prefs, err := ns.EnsureNotificationPreferences(sid, now)
	if err != nil || prefs == nil {
		t.Fatalf("ensure prefs: %v %+v", err, prefs)
	}
	prefs.EmailEnabled = true
	if err := ns.UpsertNotificationPreferences(prefs); err != nil {
		t.Fatalf("save prefs: %v", err)
	}

	token, err := auth.NotificationUnsubscribeToken(sid)
	if err != nil {
		t.Fatalf("mint token: %v", err)
	}
	rec := notificationRequest(t, mux, http.MethodPost, "/api/notifications/unsubscribe", "", fmt.Sprintf(`{"token":%q}`, token))
	if rec.Code != http.StatusOK {
		t.Fatalf("unsubscribe = %d: %s", rec.Code, rec.Body.String())
	}
	updated, err := ns.GetNotificationPreferences(sid)
	if err != nil || updated == nil {
		t.Fatalf("load prefs: %v %+v", err, updated)
	}
	if updated.EmailEnabled {
		t.Error("unsubscribe left email enabled")
	}
	if !updated.InAppEnabled {
		t.Error("unsubscribe changed the in-app preference")
	}

	rec = notificationRequest(t, mux, http.MethodPost, "/api/notifications/unsubscribe", "", `{"token":"not-a-token"}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid token = %d, want 400", rec.Code)
	}
}

func TestNotificationPreferencesAPI(t *testing.T) {
	_, mux, store, token, sid := setupTestServerWithPlanner(t)

	rec := notificationRequest(t, mux, http.MethodGet, "/api/notifications/preferences", token, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get prefs = %d: %s", rec.Code, rec.Body.String())
	}
	var prefs storage.NotificationPreferences
	if err := json.Unmarshal(rec.Body.Bytes(), &prefs); err != nil {
		t.Fatalf("decode prefs: %v", err)
	}
	if prefs.StudentID != sid || !prefs.InAppEnabled || prefs.Timezone != "UTC" {
		t.Fatalf("default prefs = %+v", prefs)
	}

	body := fmt.Sprintf(`{
		"student_id":"someone-else",
		"email_enabled":true,
		"in_app_enabled":false,
		"reviews_enabled":false,
		"study_enabled":true,
		"milestones_enabled":false,
		"inactivity_enabled":true,
		"preferred_time_min":420,
		"study_days":"fri,mon",
		"study_time_min":480,
		"timezone":"America/New_York",
		"quiet_start_min":1260,
		"quiet_end_min":360
	}`)
	rec = notificationRequest(t, mux, http.MethodPut, "/api/notifications/preferences", token, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("put prefs = %d: %s", rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &prefs); err != nil {
		t.Fatalf("decode saved prefs: %v", err)
	}
	if prefs.StudentID != sid {
		t.Errorf("saved student id = %q, want authenticated %q", prefs.StudentID, sid)
	}
	if prefs.StudyDays != "mon,fri" {
		t.Errorf("study days = %q, want canonical mon,fri", prefs.StudyDays)
	}
	if !prefs.EmailEnabled || prefs.InAppEnabled || prefs.ReviewsEnabled || !prefs.StudyEnabled || prefs.MilestonesEnabled || !prefs.InactivityEnabled {
		t.Errorf("saved booleans = %+v", prefs)
	}
	if prefs.PreferredTimeMin != 420 || prefs.StudyTimeMin != 480 || prefs.Timezone != "America/New_York" {
		t.Errorf("saved schedule = %+v", prefs)
	}

	// The body's student_id was ignored, and the preferences persist.
	stored, err := store.(storage.NotificationStore).GetNotificationPreferences(sid)
	if err != nil || stored == nil {
		t.Fatalf("get stored prefs: %v %+v", err, stored)
	}
	if stored.StudentID != sid || stored.Timezone != "America/New_York" {
		t.Errorf("stored prefs = %+v", stored)
	}

	for _, invalid := range []string{
		`{"timezone":"Mars/Olympus","preferred_time_min":420,"study_time_min":480,"quiet_start_min":1260,"quiet_end_min":360}`,
		`{"preferred_time_min":1440,"study_time_min":480,"quiet_start_min":1260,"quiet_end_min":360}`,
		`{"study_enabled":true,"study_days":"","preferred_time_min":420,"study_time_min":480,"quiet_start_min":1260,"quiet_end_min":360}`,
		`{"study_days":"funday","preferred_time_min":420,"study_time_min":480,"quiet_start_min":1260,"quiet_end_min":360}`,
		`{"paused_until":"2000-01-01T00:00:00Z","preferred_time_min":420,"study_time_min":480,"quiet_start_min":1260,"quiet_end_min":360}`,
	} {
		rec = notificationRequest(t, mux, http.MethodPut, "/api/notifications/preferences", token, invalid)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("invalid prefs %s = %d, want 400: %s", invalid, rec.Code, rec.Body.String())
		}
	}
}
