package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/notifications"
	"github.com/chuma-beep/mathua/internal/storage"
)

const (
	notificationPageDefault = 20
	notificationPageMax     = 50
)

type notificationJSON struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Href      string `json:"href"`
	CreatedAt string `json:"created_at"`
	Read      bool   `json:"read"`
	ReadAt    string `json:"read_at,omitempty"`
}

func notificationDTO(n storage.Notification) notificationJSON {
	out := notificationJSON{
		ID:        n.ID,
		Kind:      n.Kind,
		Title:     n.Title,
		Body:      n.Body,
		Href:      n.Href,
		CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
		Read:      !n.ReadAt.IsZero(),
	}
	if !n.ReadAt.IsZero() {
		out.ReadAt = n.ReadAt.UTC().Format(time.RFC3339)
	}
	return out
}

func notificationStore(w http.ResponseWriter, repo storage.Repository) (storage.NotificationStore, bool) {
	store, ok := repo.(storage.NotificationStore)
	if !ok {
		writeError(w, "notifications are not supported by this store", http.StatusInternalServerError)
		return nil, false
	}
	return store, true
}

// notificationInAppEnabled is the Phase 3 policy read. A learner with no
// preference row gets the Phase 1 default, which is in-app on. Turning the
// preference off hides the feed and zeroes the badge without deleting history.
func notificationInAppEnabled(store storage.NotificationStore, studentID string) (bool, error) {
	prefs, err := store.GetNotificationPreferences(studentID)
	if err != nil {
		return false, err
	}
	if prefs == nil {
		return true, nil
	}
	return prefs.InAppEnabled, nil
}

func notificationCursor(createdAt int64, id string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d|%s", createdAt, id)))
}

func parseNotificationCursor(raw string) (int64, string, error) {
	if raw == "" {
		return 0, "", nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return 0, "", fmt.Errorf("invalid cursor")
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 || parts[1] == "" {
		return 0, "", fmt.Errorf("invalid cursor")
	}
	createdAt, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || createdAt <= 0 {
		return 0, "", fmt.Errorf("invalid cursor")
	}
	return createdAt, parts[1], nil
}

// GET /api/notifications returns one stable page of the authenticated
// learner's in-app feed. Ordering is newest first, with the id as tie-breaker,
// so a cursor cannot skip or repeat an item created in the same second.
func (s *Server) handleNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}

	limit := notificationPageDefault
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > notificationPageMax {
			writeError(w, "limit must be between 1 and 50", http.StatusBadRequest)
			return
		}
		limit = n
	}
	unreadOnly := false
	if raw := r.URL.Query().Get("unread"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			writeError(w, "unread must be true or false", http.StatusBadRequest)
			return
		}
		unreadOnly = v
	}
	beforeCreatedAt, beforeID, err := parseNotificationCursor(r.URL.Query().Get("cursor"))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	inApp, err := notificationInAppEnabled(store, studentID)
	if err != nil {
		writeError(w, "failed to load notifications", http.StatusInternalServerError)
		return
	}
	if !inApp {
		writeJSON(w, map[string]interface{}{
			"notifications": []notificationJSON{},
			"unread":        0,
			"has_more":      false,
			"next_cursor":   "",
		})
		return
	}

	page, err := store.ListNotifications(studentID, storage.NotificationListOptions{
		Limit:           limit,
		BeforeCreatedAt: beforeCreatedAt,
		BeforeID:        beforeID,
		UnreadOnly:      unreadOnly,
	})
	if err != nil {
		writeError(w, "failed to load notifications", http.StatusInternalServerError)
		return
	}
	unread, err := store.NotificationUnreadCount(studentID)
	if err != nil {
		writeError(w, "failed to count notifications", http.StatusInternalServerError)
		return
	}
	out := make([]notificationJSON, 0, len(page.Notifications))
	for _, n := range page.Notifications {
		out = append(out, notificationDTO(n))
	}
	next := ""
	if page.HasMore && len(page.Notifications) > 0 {
		last := page.Notifications[len(page.Notifications)-1]
		next = notificationCursor(last.CreatedAt.UTC().Unix(), last.ID)
	}
	writeJSON(w, map[string]interface{}{
		"notifications": out,
		"unread":        unread,
		"has_more":      page.HasMore,
		"next_cursor":   next,
	})
}

// GET /api/notifications/unread-count is the small badge read. It has no
// pagination and never marks anything read.
func (s *Server) handleNotificationUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	inApp, err := notificationInAppEnabled(store, studentID)
	if err != nil {
		writeError(w, "failed to load notifications", http.StatusInternalServerError)
		return
	}
	if !inApp {
		writeJSON(w, map[string]interface{}{"unread": 0})
		return
	}
	unread, err := store.NotificationUnreadCount(studentID)
	if err != nil {
		writeError(w, "failed to count notifications", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"unread": unread})
}

// POST /api/notifications/read marks one owned notification read. Another
// learner's id is indistinguishable from a missing id.
func (s *Server) handleNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.ID) == "" {
		writeError(w, "notification id is required", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	if err := store.MarkNotificationRead(studentID, req.ID, time.Now().UTC()); err != nil {
		if errors.Is(err, storage.ErrNotificationNotFound) {
			writeError(w, "notification not found", http.StatusNotFound)
			return
		}
		writeError(w, "failed to mark notification read", http.StatusInternalServerError)
		return
	}
	unread := 0
	if inApp, err := notificationInAppEnabled(store, studentID); err == nil && inApp {
		unread, err = store.NotificationUnreadCount(studentID)
		if err != nil {
			writeError(w, "failed to count notifications", http.StatusInternalServerError)
			return
		}
	}
	writeJSON(w, map[string]interface{}{"ok": true, "unread": unread})
}

// POST /api/notifications/read-all marks every unread item owned by the caller.
func (s *Server) handleNotificationsReadAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	if err := store.MarkAllNotificationsRead(studentID, time.Now().UTC()); err != nil {
		writeError(w, "failed to mark notifications read", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "unread": 0})
}

// GET/PUT /api/notifications/preferences reads and persists the Phase 1
// preference policy. The student id always comes from the authenticated
// context, never from the request body.
func (s *Server) handleNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPut {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	if r.Method == http.MethodGet {
		prefs, err := store.EnsureNotificationPreferences(studentID, now)
		if err != nil {
			writeError(w, "failed to load notification preferences", http.StatusInternalServerError)
			return
		}
		writeJSON(w, prefs)
		return
	}

	var body storage.NotificationPreferences
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	body.StudentID = studentID
	body.CreatedAt = now
	if existing, err := store.GetNotificationPreferences(studentID); err == nil && existing != nil {
		body.CreatedAt = existing.CreatedAt
	}
	normalized, err := notifications.NormalizePreferencesAt(body, now)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	normalized.StudentID = studentID
	normalized.UpdatedAt = now
	if err := store.UpsertNotificationPreferences(&normalized); err != nil {
		writeError(w, "failed to save notification preferences", http.StatusInternalServerError)
		return
	}
	writeJSON(w, normalized)
}

// POST /api/notifications/unsubscribe is public because the signed link in an
// email is the credential. It only turns email delivery off; in-app items and
// the account itself are untouched.
func (s *Server) handleNotificationUnsubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	studentID, err := auth.ParseNotificationUnsubscribeToken(req.Token)
	if err != nil {
		writeError(w, "unsubscribe link is invalid or expired", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	prefs, err := store.GetNotificationPreferences(studentID)
	if err != nil {
		writeError(w, "failed to load notification preferences", http.StatusInternalServerError)
		return
	}
	if prefs == nil {
		prefs, err = store.EnsureNotificationPreferences(studentID, now)
		if err != nil {
			writeError(w, "failed to load notification preferences", http.StatusInternalServerError)
			return
		}
	}
	prefs.EmailEnabled = false
	prefs.UpdatedAt = now
	if err := store.UpsertNotificationPreferences(prefs); err != nil {
		writeError(w, "failed to save notification preferences", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// GET /api/notifications/push/key returns the VAPID application-server key a
// browser needs before it can subscribe. configured is false when the
// deployment has no key pair, and the client hides the toggle instead of
// offering a subscription that could never be delivered.
func (s *Server) handleNotificationPushKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, map[string]interface{}{
		"configured": s.pushPublicKey != "",
		"public_key": s.pushPublicKey,
	})
}

type notificationPushKeys struct {
	P256DH string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type notificationPushRequest struct {
	Endpoint string               `json:"endpoint"`
	Keys     notificationPushKeys `json:"keys"`
}

// POST /api/notifications/push/subscribe stores one device endpoint and turns
// the push channel on for that learner. The endpoint is the identity, so a
// browser that re-subscribes updates its keys rather than adding a device.
func (s *Server) handleNotificationPushSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.pushPublicKey == "" {
		writeError(w, "browser push is not configured", http.StatusServiceUnavailable)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	var req notificationPushRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	endpoint := strings.TrimSpace(req.Endpoint)
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		writeError(w, "a push endpoint must be an https URL", http.StatusBadRequest)
		return
	}
	p256dh := strings.TrimSpace(req.Keys.P256DH)
	authSecret := strings.TrimSpace(req.Keys.Auth)
	if p256dh == "" || authSecret == "" {
		writeError(w, "push keys are required", http.StatusBadRequest)
		return
	}
	if len(endpoint) > 2048 || len(p256dh) > 256 || len(authSecret) > 128 {
		writeError(w, "push subscription is too large", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	if err := store.UpsertNotificationPushSubscription(&storage.NotificationPushSubscription{
		StudentID: studentID,
		Endpoint:  endpoint,
		P256DH:    p256dh,
		Auth:      authSecret,
		UserAgent: truncateUserAgent(r.UserAgent()),
	}); err != nil {
		writeError(w, "failed to save push subscription", http.StatusInternalServerError)
		return
	}
	prefs, err := store.EnsureNotificationPreferences(studentID, now)
	if err != nil {
		writeError(w, "failed to load notification preferences", http.StatusInternalServerError)
		return
	}
	prefs.PushEnabled = true
	prefs.UpdatedAt = now
	if err := store.UpsertNotificationPreferences(prefs); err != nil {
		writeError(w, "failed to save notification preferences", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "devices": 1})
}

// POST /api/notifications/push/unsubscribe removes one device endpoint and
// turns the push channel off. Deleting is scoped to the authenticated learner,
// so an endpoint guessed from elsewhere cannot be removed by someone else.
func (s *Server) handleNotificationPushUnsubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	var req struct {
		Endpoint string `json:"endpoint"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	endpoint := strings.TrimSpace(req.Endpoint)
	if endpoint == "" {
		writeError(w, "endpoint is required", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	if err := store.DeleteNotificationPushSubscription(studentID, endpoint); err != nil {
		writeError(w, "failed to remove push subscription", http.StatusInternalServerError)
		return
	}
	prefs, err := store.EnsureNotificationPreferences(studentID, now)
	if err != nil {
		writeError(w, "failed to load notification preferences", http.StatusInternalServerError)
		return
	}
	prefs.PushEnabled = false
	prefs.UpdatedAt = now
	if err := store.UpsertNotificationPreferences(prefs); err != nil {
		writeError(w, "failed to save notification preferences", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// notificationDeviceID is the public handle for a push endpoint. The raw
// endpoint is a capability URL, so the API never returns it; the truncated
// hash is stable enough for a learner to recognise and revoke one device.
func notificationDeviceID(endpoint string) string {
	sum := sha256.Sum256([]byte(endpoint))
	return hex.EncodeToString(sum[:8])
}

type notificationDeviceJSON struct {
	ID          string `json:"id"`
	UserAgent   string `json:"user_agent,omitempty"`
	CreatedAt   string `json:"created_at"`
	LastSuccess string `json:"last_success_at,omitempty"`
	Failures    int    `json:"failures"`
}

// GET /api/notifications/push/devices lists the learner's registered browsers
// with a hash instead of the endpoint, so a stolen response cannot be used to
// send a push to that device.
func (s *Server) handleNotificationPushDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	subs, err := store.ListNotificationPushSubscriptions(studentID)
	if err != nil {
		writeError(w, "failed to list push devices", http.StatusInternalServerError)
		return
	}
	devices := make([]notificationDeviceJSON, 0, len(subs))
	for _, sub := range subs {
		device := notificationDeviceJSON{
			ID:        notificationDeviceID(sub.Endpoint),
			UserAgent: sub.UserAgent,
			CreatedAt: sub.CreatedAt.UTC().Format(time.RFC3339),
			Failures:  sub.Failures,
		}
		if !sub.LastSuccess.IsZero() {
			device.LastSuccess = sub.LastSuccess.UTC().Format(time.RFC3339)
		}
		devices = append(devices, device)
	}
	writeJSON(w, map[string]interface{}{"devices": devices})
}

// POST /api/notifications/push/revoke removes one device by its public hash.
// The lookup is scoped to the authenticated learner, so a hash from another
// account cannot revoke someone else's browser.
func (s *Server) handleNotificationPushRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := strings.TrimSpace(req.ID)
	if id == "" {
		writeError(w, "device id is required", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	subs, err := store.ListNotificationPushSubscriptions(studentID)
	if err != nil {
		writeError(w, "failed to list push devices", http.StatusInternalServerError)
		return
	}
	removed := 0
	for _, sub := range subs {
		if notificationDeviceID(sub.Endpoint) != id {
			continue
		}
		if err := store.DeleteNotificationPushSubscription(studentID, sub.Endpoint); err != nil {
			writeError(w, "failed to revoke push device", http.StatusInternalServerError)
			return
		}
		removed++
	}
	if removed == 0 {
		writeError(w, "push device not found", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "revoked": removed})
}

type notificationEngagementRequest struct {
	ID string `json:"id"`
}

// POST /api/notifications/open records that a learner opened a notification in
// the feed. It marks the row read as well, because an opened item the learner
// has seen should not keep the unread badge.
func (s *Server) handleNotificationOpen(w http.ResponseWriter, r *http.Request) {
	s.markNotificationEngagement(w, r, true)
}

// POST /api/notifications/click records that a learner followed a push or feed
// link. A click implies an open, so both stamps are written.
func (s *Server) handleNotificationClick(w http.ResponseWriter, r *http.Request) {
	s.markNotificationEngagement(w, r, false)
}

func (s *Server) markNotificationEngagement(w http.ResponseWriter, r *http.Request, opened bool) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	var req notificationEngagementRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := strings.TrimSpace(req.ID)
	if id == "" {
		writeError(w, "notification id is required", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	now := time.Now().UTC()
	var err error
	if opened {
		err = store.MarkNotificationOpened(studentID, id, now)
	} else {
		err = store.MarkNotificationClicked(studentID, id, now)
	}
	if errors.Is(err, storage.ErrNotificationNotFound) {
		writeError(w, "notification not found", http.StatusNotFound)
		return
	}
	if err != nil {
		writeError(w, "failed to record notification engagement", http.StatusInternalServerError)
		return
	}
	unread, _ := store.NotificationUnreadCount(studentID)
	writeJSON(w, map[string]interface{}{"ok": true, "unread": unread})
}

// truncateUserAgent keeps a stored user agent useful for telling two devices
// apart without letting a client write an unbounded string into the row.
func truncateUserAgent(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) > 200 {
		return raw[:200]
	}
	return raw
}
