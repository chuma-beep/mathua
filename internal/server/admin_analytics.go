package server

import (
	"net/http"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

const adminAnalyticsCacheSeconds = 30

// adminAnalyticsRange is the resolved window for one historical request.
// Rolling ranges are measured back from now; daily ranges start at UTC
// midnight and include the partial current day.
type adminAnalyticsRange struct {
	Name   string
	Since  time.Time
	Until  time.Time
	Bucket time.Duration
}

// adminAnalyticsResponse embeds the storage projection so the JSON shape and
// the storage types cannot drift into two definitions.
type adminAnalyticsResponse struct {
	Range             string                     `json:"range"`
	From              time.Time                  `json:"from"`
	To                time.Time                  `json:"to"`
	BucketSeconds     int64                      `json:"bucket_seconds"`
	GeneratedAt       time.Time                  `json:"generated_at"`
	PartialLastBucket bool                       `json:"partial_last_bucket"`
	Previous          *storage.AdminMetricTotals `json:"previous,omitempty"`
	*storage.AdminMetrics
}

type adminAnalyticsLiveResponse struct {
	OnlineNow           int       `json:"online_now"`
	OnlineWindowSeconds int       `json:"online_window_seconds"`
	PracticingNow       int       `json:"practicing_now"`
	GuestPracticeNow    int       `json:"guest_practice_now"`
	Answers15m          int       `json:"answers_15m"`
	ActiveSessions      int       `json:"active_sessions"`
	GeneratedAt         time.Time `json:"generated_at"`
}

func (s *Server) handleAdminAnalytics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	store, ok := s.repo.(storage.AdminMetricsStore)
	if !ok {
		writeError(w, "analytics are not available", http.StatusNotImplemented)
		return
	}
	window, ok := resolveAdminAnalyticsRange(r.URL.Query().Get("range"), time.Now().UTC())
	if !ok {
		writeError(w, "invalid range", http.StatusBadRequest)
		return
	}
	metrics, err := store.AdminMetrics(window.Since, window.Until, window.Bucket)
	if err != nil {
		writeError(w, "failed to load analytics", http.StatusInternalServerError)
		return
	}
	resp := adminAnalyticsResponse{
		Range:             window.Name,
		From:              window.Since,
		To:                window.Until,
		BucketSeconds:     int64(window.Bucket / time.Second),
		GeneratedAt:       time.Now().UTC(),
		PartialLastBucket: lastBucketPartial(window),
		AdminMetrics:      metrics,
	}
	// A previous-period comparison is only meaningful for event counts.
	// Current-state fields such as mastery and report status deliberately do
	// not appear in Previous; they are not reconstructed history.
	span := window.Until.Sub(window.Since)
	if span > 0 {
		if previous, prevErr := store.AdminMetrics(window.Since.Add(-span), window.Since, window.Bucket); prevErr == nil {
			resp.Previous = &previous.Totals
		}
	}
	w.Header().Set("Cache-Control", "private, max-age=30")
	writeJSON(w, resp)
}

func (s *Server) handleAdminAnalyticsLive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	store, ok := s.repo.(storage.AdminMetricsStore)
	if !ok {
		writeError(w, "analytics are not available", http.StatusNotImplemented)
		return
	}
	now := time.Now().UTC()
	since := now.Add(-15 * time.Minute)
	activeNow, activeGuests, _, _, err := s.repo.ActiveLearnerCounts(
		since, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))
	if err != nil {
		writeError(w, "failed to count active learners", http.StatusInternalServerError)
		return
	}
	live, err := store.AdminLiveMetrics(since)
	if err != nil {
		writeError(w, "failed to load live analytics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, adminAnalyticsLiveResponse{
		OnlineNow:           s.presence.Count(now),
		OnlineWindowSeconds: 300,
		PracticingNow:       activeNow,
		GuestPracticeNow:    activeGuests,
		Answers15m:          live.Answers15m,
		ActiveSessions:      live.ActiveSessions,
		GeneratedAt:         now,
	})
}

// resolveAdminAnalyticsRange maps the public range names to a window and a
// bucket. The daily ranges begin at UTC midnight so the first and last labels
// match the calendar a learner sees elsewhere in the product.
func resolveAdminAnalyticsRange(raw string, now time.Time) (adminAnalyticsRange, bool) {
	now = now.UTC()
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	switch raw {
	case "", "30d":
		return adminAnalyticsRange{Name: "30d", Since: dayStart.AddDate(0, 0, -29), Until: now, Bucket: 24 * time.Hour}, true
	case "15m":
		return adminAnalyticsRange{Name: "15m", Since: now.Add(-15 * time.Minute), Until: now, Bucket: time.Minute}, true
	case "1h":
		return adminAnalyticsRange{Name: "1h", Since: now.Add(-time.Hour), Until: now, Bucket: 5 * time.Minute}, true
	case "24h":
		return adminAnalyticsRange{Name: "24h", Since: now.Add(-24 * time.Hour), Until: now, Bucket: time.Hour}, true
	case "7d":
		return adminAnalyticsRange{Name: "7d", Since: dayStart.AddDate(0, 0, -6), Until: now, Bucket: 24 * time.Hour}, true
	case "90d":
		return adminAnalyticsRange{Name: "90d", Since: dayStart.AddDate(0, 0, -89), Until: now, Bucket: 24 * time.Hour}, true
	default:
		return adminAnalyticsRange{}, false
	}
}

func lastBucketPartial(window adminAnalyticsRange) bool {
	last := window.Until.Truncate(window.Bucket)
	return last.Add(window.Bucket).After(window.Until)
}
