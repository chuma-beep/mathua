package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/notifications"
	"github.com/chuma-beep/mathua/internal/storage"
)

type adminDeliveryStatJSON struct {
	Channel string `json:"channel"`
	Status  string `json:"status"`
	Count   int    `json:"count"`
}

type adminDeliveryProblemJSON struct {
	NotificationID   string `json:"notification_id"`
	StudentID        string `json:"student_id"`
	StudentName      string `json:"student_name"`
	Title            string `json:"title"`
	Channel          string `json:"channel"`
	Status           string `json:"status"`
	Attempts         int    `json:"attempts"`
	LastError        string `json:"last_error,omitempty"`
	SuppressedReason string `json:"suppressed_reason,omitempty"`
	CreatedAt        string `json:"created_at"`
}

// GET /api/admin/notifications reports notification delivery health: one count
// per channel and status, plus the deliveries that will never move again.
// Retries are deliberately absent from the problem list because they are still
// in flight; a queue that is retrying is not a queue that is stuck.
func (s *Server) handleAdminNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	stats, err := store.NotificationDeliveryStats()
	if err != nil {
		writeError(w, "failed to load delivery statistics", http.StatusInternalServerError)
		return
	}
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	problems, err := store.ListNotificationDeliveryProblems(limit)
	if err != nil {
		writeError(w, "failed to load delivery problems", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{
		"stats":    adminDeliveryStats(stats),
		"problems": adminDeliveryProblems(problems),
		"retention": map[string]interface{}{
			"days":               notifications.RetentionDays,
			"push_failure_days":  notifications.PushFailureRetentionDays,
			"push_failure_limit": notifications.PushFailureLimit,
		},
	})
}

func adminDeliveryStats(stats []storage.NotificationDeliveryStat) []adminDeliveryStatJSON {
	out := make([]adminDeliveryStatJSON, 0, len(stats))
	for _, stat := range stats {
		out = append(out, adminDeliveryStatJSON{
			Channel: stat.Channel,
			Status:  stat.Status,
			Count:   stat.Count,
		})
	}
	return out
}

func adminDeliveryProblems(records []storage.NotificationDeliveryRecord) []adminDeliveryProblemJSON {
	out := make([]adminDeliveryProblemJSON, 0, len(records))
	for _, rec := range records {
		created := ""
		if !rec.CreatedAt.IsZero() {
			created = rec.CreatedAt.UTC().Format(time.RFC3339)
		}
		out = append(out, adminDeliveryProblemJSON{
			NotificationID:   rec.NotificationID,
			StudentID:        rec.StudentID,
			StudentName:      rec.StudentName,
			Title:            rec.Title,
			Channel:          rec.Channel,
			Status:           rec.Status,
			Attempts:         rec.Attempts,
			LastError:        rec.LastError,
			SuppressedReason: rec.SuppressedReason,
			CreatedAt:        created,
		})
	}
	return out
}

type adminDeliveryDetailJSON struct {
	NotificationID   string `json:"notification_id"`
	Channel          string `json:"channel"`
	Status           string `json:"status"`
	Attempts         int    `json:"attempts"`
	ProviderID       string `json:"provider_id,omitempty"`
	LastError        string `json:"last_error,omitempty"`
	SuppressedReason string `json:"suppressed_reason,omitempty"`
	NextAttemptAt    string `json:"next_attempt_at,omitempty"`
	SentAt           string `json:"sent_at,omitempty"`
	OpenedAt         string `json:"opened_at,omitempty"`
	ClickedAt        string `json:"clicked_at,omitempty"`
}

// POST /api/admin/notifications/replay puts one failed or suppressed delivery
// back on the queue. Only those two terminal states are replayable, so a
// successful send cannot be duplicated by an operator click.
func (s *Server) handleAdminNotificationReplay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		NotificationID string `json:"notification_id"`
		Channel        string `json:"channel"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	notificationID := strings.TrimSpace(req.NotificationID)
	channel := strings.TrimSpace(req.Channel)
	if notificationID == "" || channel == "" {
		writeError(w, "notification_id and channel are required", http.StatusBadRequest)
		return
	}
	if channel != storage.NotificationChannelEmail && channel != storage.NotificationChannelPush {
		writeError(w, "only email and push deliveries can be replayed", http.StatusBadRequest)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	if err := store.ReplayNotificationDelivery(notificationID, channel, time.Now().UTC()); err != nil {
		if errors.Is(err, storage.ErrNotificationNotFound) {
			writeError(w, "delivery is not replayable", http.StatusConflict)
			return
		}
		writeError(w, "failed to replay delivery", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// GET /api/admin/notifications/history returns the delivery rows for one
// notification, or the recent rows for one learner. The two shapes answer the
// two questions an operator asks: "why did this message not arrive" and "what
// has this learner been sent".
func (s *Server) handleAdminNotificationHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	store, ok := notificationStore(w, s.repo)
	if !ok {
		return
	}
	notificationID := strings.TrimSpace(r.URL.Query().Get("notification_id"))
	studentID := strings.TrimSpace(r.URL.Query().Get("student_id"))
	switch {
	case notificationID != "":
		deliveries, err := store.ListNotificationDeliveries(notificationID)
		if err != nil {
			writeError(w, "failed to load delivery history", http.StatusInternalServerError)
			return
		}
		out := make([]adminDeliveryDetailJSON, 0, len(deliveries))
		for _, d := range deliveries {
			out = append(out, adminDeliveryDetail(d))
		}
		writeJSON(w, map[string]interface{}{"deliveries": out})
	case studentID != "":
		limit := 50
		if raw := r.URL.Query().Get("limit"); raw != "" {
			if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		records, err := store.ListStudentNotificationDeliveries(studentID, limit)
		if err != nil {
			writeError(w, "failed to load learner history", http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{"deliveries": adminDeliveryProblems(records)})
	default:
		writeError(w, "notification_id or student_id is required", http.StatusBadRequest)
	}
}

func adminDeliveryDetail(d storage.NotificationDelivery) adminDeliveryDetailJSON {
	detail := adminDeliveryDetailJSON{
		NotificationID:   d.NotificationID,
		Channel:          d.Channel,
		Status:           d.Status,
		Attempts:         d.Attempts,
		ProviderID:       d.ProviderID,
		LastError:        d.LastError,
		SuppressedReason: d.SuppressedReason,
	}
	detail.NextAttemptAt = optionalRFC3339(d.NextAttemptAt)
	detail.SentAt = optionalRFC3339(d.SentAt)
	detail.OpenedAt = optionalRFC3339(d.OpenedAt)
	detail.ClickedAt = optionalRFC3339(d.ClickedAt)
	return detail
}

func optionalRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
