package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	unread, err := store.NotificationUnreadCount(studentID)
	if err != nil {
		writeError(w, "failed to count notifications", http.StatusInternalServerError)
		return
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
