package server

import (
	"net/http"
	"strconv"
	"time"

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
