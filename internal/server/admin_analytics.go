package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

const adminAnalyticsCacheTTL = 30 * time.Second

// adminAnalyticsCache is an in-process TTL cache for historical aggregates.
// The live strip never touches it: live values must be current, and caching
// them would make the stale badge lie.
type adminAnalyticsCache struct {
	mu       sync.Mutex
	entries  map[string]adminAnalyticsCacheEntry
	inflight map[string]*adminAnalyticsFlight
}

type adminAnalyticsCacheEntry struct {
	expires time.Time
	payload []byte
}

// adminAnalyticsFlight is one in-process computation. A concurrent request for
// the same key waits on done instead of issuing the same expensive queries.
type adminAnalyticsFlight struct {
	done    chan struct{}
	payload []byte
	err     error
}

func newAdminAnalyticsCache() *adminAnalyticsCache {
	return &adminAnalyticsCache{
		entries:  make(map[string]adminAnalyticsCacheEntry),
		inflight: make(map[string]*adminAnalyticsFlight),
	}
}

// begin returns a cached payload when one is fresh. Otherwise exactly one
// caller becomes the leader and the rest wait on the returned flight.
func (c *adminAnalyticsCache) begin(key string, now time.Time) (payload []byte, flight *adminAnalyticsFlight, leader, cached bool) {
	if c == nil {
		return nil, nil, true, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = make(map[string]adminAnalyticsCacheEntry)
	}
	if c.inflight == nil {
		c.inflight = make(map[string]*adminAnalyticsFlight)
	}
	if entry, ok := c.entries[key]; ok && now.Before(entry.expires) {
		return entry.payload, nil, false, true
	}
	if existing := c.inflight[key]; existing != nil {
		return nil, existing, false, false
	}
	flight = &adminAnalyticsFlight{done: make(chan struct{})}
	c.inflight[key] = flight
	return nil, flight, true, false
}

func (c *adminAnalyticsCache) wait(ctx context.Context, flight *adminAnalyticsFlight) ([]byte, error) {
	if flight == nil {
		return nil, errors.New("analytics flight is missing")
	}
	select {
	case <-flight.done:
		return flight.payload, flight.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *adminAnalyticsCache) finish(key string, payload []byte, err error, now time.Time) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	flight := c.inflight[key]
	if flight == nil {
		return
	}
	flight.payload = payload
	flight.err = err
	close(flight.done)
	delete(c.inflight, key)
	if err != nil {
		return
	}
	// Bound the map: only the handful of range and section combinations the
	// admin page requests should ever be present.
	if len(c.entries) > 64 {
		for k := range c.entries {
			delete(c.entries, k)
			if len(c.entries) <= 32 {
				break
			}
		}
	}
	c.entries[key] = adminAnalyticsCacheEntry{expires: now.Add(adminAnalyticsCacheTTL), payload: payload}
}

// adminAnalyticsRange is the resolved window for one historical request.
// Rolling ranges are measured back from now; daily ranges start at UTC
// midnight and include the partial current day.
type adminAnalyticsRange struct {
	Name string
	// Key identifies the window in the cache. Presets key on the name, but a
	// custom window must include its dates or every custom request would
	// share the first response for the next 30 seconds.
	Key    string
	Since  time.Time
	Until  time.Time
	Bucket time.Duration
}

// adminAnalyticsResponse is explicit rather than an embedding so a request
// for one section can omit the others and the payload is genuinely lazy.
type adminAnalyticsResponse struct {
	Range             string                     `json:"range"`
	View              string                     `json:"view,omitempty"`
	Sections          []string                   `json:"sections,omitempty"`
	From              time.Time                  `json:"from"`
	To                time.Time                  `json:"to"`
	BucketSeconds     int64                      `json:"bucket_seconds"`
	GeneratedAt       time.Time                  `json:"generated_at"`
	PartialLastBucket bool                       `json:"partial_last_bucket"`
	Previous          *storage.AdminMetricTotals `json:"previous,omitempty"`
	PreviousFrom      *time.Time                 `json:"previous_from,omitempty"`
	PreviousTo        *time.Time                 `json:"previous_to,omitempty"`

	Accounts              *storage.AdminAccountCounts  `json:"accounts,omitempty"`
	Totals                *storage.AdminMetricTotals   `json:"totals,omitempty"`
	Activity              []storage.AdminActivityPoint `json:"activity,omitempty"`
	Signups               []storage.AdminSeriesPoint   `json:"signups,omitempty"`
	Mastery               []storage.AdminSeriesPoint   `json:"mastery,omitempty"`
	XP                    []storage.AdminXPPoint       `json:"xp,omitempty"`
	Reports               []storage.AdminReportPoint   `json:"reports,omitempty"`
	Quizzes               []storage.AdminSeriesPoint   `json:"quizzes,omitempty"`
	TimedQuizzes          []storage.AdminSeriesPoint   `json:"timed_quizzes,omitempty"`
	DiagnosticCompletions []storage.AdminSeriesPoint   `json:"diagnostic_completions,omitempty"`
	TopConcepts           []storage.AdminConceptMetric `json:"top_concepts,omitempty"`
	MostReported          []storage.AdminConceptMetric `json:"most_reported,omitempty"`
	MasteredByConcept     map[string]int               `json:"mastered_by_concept,omitempty"`
	MasteredCurrent       *int                         `json:"mastered_current,omitempty"`
	ReportStatus          map[string]int               `json:"report_status,omitempty"`
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

// adminAnalyticsViews maps the UI view names to the sections they need. A
// caller can also pass sections directly; view is a convenience, not a second
// data path.
var adminAnalyticsViews = map[string][]string{
	"overview":       {"accounts", "totals", "activity"},
	"usage":          {"activity", "signups", "quizzes", "diagnostics"},
	"learning":       {"mastery", "xp", "quizzes", "concepts"},
	"operations":     {"reports"},
	"administration": {"accounts"},
}

var adminAnalyticsSectionNames = map[string]bool{
	"all": true, "accounts": true, "totals": true, "activity": true,
	"signups": true, "mastery": true, "xp": true, "reports": true,
	"quizzes": true, "timed_quizzes": true, "diagnostics": true,
	"concepts": true, "status": true,
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
	now := time.Now().UTC()
	window, ok := resolveAdminAnalyticsRange(
		r.URL.Query().Get("range"),
		r.URL.Query().Get("from"),
		r.URL.Query().Get("to"),
		now,
	)
	if !ok {
		writeError(w, "invalid range", http.StatusBadRequest)
		return
	}
	view := strings.TrimSpace(r.URL.Query().Get("view"))
	sections, err := resolveAdminAnalyticsSections(r.URL.Query().Get("sections"), view)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	cacheKey := adminAnalyticsCacheKey(window.Key, view, sections)
	payload, flight, leader, cached := s.analyticsCache.begin(cacheKey, now)
	if cached {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private, max-age=30")
		_, _ = w.Write(payload)
		return
	}
	if !leader {
		payload, waitErr := s.analyticsCache.wait(r.Context(), flight)
		if waitErr != nil {
			writeError(w, "analytics request was canceled", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "private, max-age=30")
		_, _ = w.Write(payload)
		return
	}

	metrics, err := store.AdminMetrics(window.Since, window.Until, window.Bucket)
	if err != nil {
		s.analyticsCache.finish(cacheKey, nil, err, now)
		writeError(w, "failed to load analytics", http.StatusInternalServerError)
		return
	}
	resp := buildAdminAnalyticsResponse(window, metrics, sections, view, now)

	// A previous-period comparison is only meaningful for event counts and
	// only when totals are requested. It is an equal elapsed window, so a
	// partial current day is never compared with a full prior day.
	if adminSectionWanted(sections, "totals") {
		span := window.Until.Sub(window.Since)
		if span > 0 {
			prevFrom := window.Since.Add(-span)
			prevTo := window.Since
			if previous, prevErr := store.AdminMetrics(prevFrom, prevTo, window.Bucket, "totals"); prevErr == nil {
				resp.Previous = &previous.Totals
				resp.PreviousFrom = &prevFrom
				resp.PreviousTo = &prevTo
			}
		}
	}

	payload, err = json.Marshal(resp)
	if err != nil {
		s.analyticsCache.finish(cacheKey, nil, err, now)
		writeError(w, "failed to encode analytics", http.StatusInternalServerError)
		return
	}
	s.analyticsCache.finish(cacheKey, payload, nil, now)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "private, max-age=30")
	_, _ = w.Write(payload)
}

func buildAdminAnalyticsResponse(window adminAnalyticsRange, metrics *storage.AdminMetrics, sections []string, view string, now time.Time) adminAnalyticsResponse {
	resp := adminAnalyticsResponse{
		Range:             window.Name,
		View:              view,
		Sections:          sections,
		From:              window.Since,
		To:                window.Until,
		BucketSeconds:     int64(window.Bucket / time.Second),
		GeneratedAt:       now,
		PartialLastBucket: lastBucketPartial(window),
	}
	if adminSectionWanted(sections, "accounts") {
		resp.Accounts = &metrics.Accounts
	}
	if adminSectionWanted(sections, "totals") {
		resp.Totals = &metrics.Totals
	}
	if adminSectionWanted(sections, "activity") {
		resp.Activity = metrics.Activity
	}
	if adminSectionWanted(sections, "signups") {
		resp.Signups = metrics.Signups
	}
	if adminSectionWanted(sections, "mastery") {
		resp.Mastery = metrics.Mastery
		resp.MasteredByConcept = metrics.MasteredByConcept
		resp.MasteredCurrent = &metrics.MasteredCurrent
	}
	if adminSectionWanted(sections, "xp") {
		resp.XP = metrics.XP
	}
	if adminSectionWanted(sections, "reports") || adminSectionWanted(sections, "status") {
		resp.Reports = metrics.Reports
		resp.ReportStatus = metrics.ReportStatus
	}
	if adminSectionWanted(sections, "quizzes") {
		resp.Quizzes = metrics.Quizzes
		resp.TimedQuizzes = metrics.TimedQuizzes
	}
	if adminSectionWanted(sections, "diagnostics") {
		resp.DiagnosticCompletions = metrics.DiagnosticCompletions
	}
	if adminSectionWanted(sections, "concepts") {
		resp.TopConcepts = metrics.TopConcepts
		resp.MostReported = metrics.MostReported
	}
	return resp
}

func adminSectionWanted(sections []string, name string) bool {
	if len(sections) == 0 {
		return true
	}
	for _, section := range sections {
		if section == "all" || section == name {
			return true
		}
	}
	return false
}

func resolveAdminAnalyticsSections(raw, view string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		if view == "" {
			return nil, nil
		}
		sections, ok := adminAnalyticsViews[view]
		if !ok {
			return nil, errInvalidAdminView
		}
		return sections, nil
	}
	seen := map[string]bool{}
	out := make([]string, 0, 4)
	for _, part := range strings.Split(raw, ",") {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		if !adminAnalyticsSectionNames[name] {
			return nil, errInvalidAdminSection
		}
		if !seen[name] {
			seen[name] = true
			out = append(out, name)
		}
	}
	sort.Strings(out)
	if len(out) == 0 {
		return nil, errInvalidAdminSection
	}
	return out, nil
}

var (
	errInvalidAdminView    = &adminAnalyticsError{"invalid view"}
	errInvalidAdminSection = &adminAnalyticsError{"invalid sections"}
)

type adminAnalyticsError struct{ msg string }

func (e *adminAnalyticsError) Error() string { return e.msg }

func adminAnalyticsCacheKey(windowKey, view string, sections []string) string {
	copied := append([]string(nil), sections...)
	sort.Strings(copied)
	return windowKey + "|" + view + "|" + strings.Join(copied, ",")
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
//
// A custom range carries a start and end calendar day, both inclusive and both
// UTC. The end day is extended to the end of that day and clamped to now, so an
// operator asking for today does not get a bucket for the rest of the day.
func resolveAdminAnalyticsRange(raw, fromRaw, toRaw string, now time.Time) (adminAnalyticsRange, bool) {
	now = now.UTC()
	if raw == "custom" {
		return resolveCustomAdminAnalyticsRange(fromRaw, toRaw, now)
	}
	// Dates on any other range are a caller mistake, not something to ignore:
	// answering with a preset would show a different window than the one asked
	// for.
	if fromRaw != "" || toRaw != "" {
		return adminAnalyticsRange{}, false
	}
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	switch raw {
	case "", "30d":
		return adminAnalyticsRange{Name: "30d", Key: "30d", Since: dayStart.AddDate(0, 0, -29), Until: now, Bucket: 24 * time.Hour}, true
	case "15m":
		return adminAnalyticsRange{Name: "15m", Key: "15m", Since: now.Add(-15 * time.Minute), Until: now, Bucket: time.Minute}, true
	case "1h":
		return adminAnalyticsRange{Name: "1h", Key: "1h", Since: now.Add(-time.Hour), Until: now, Bucket: 5 * time.Minute}, true
	case "24h":
		return adminAnalyticsRange{Name: "24h", Key: "24h", Since: now.Add(-24 * time.Hour), Until: now, Bucket: time.Hour}, true
	case "7d":
		return adminAnalyticsRange{Name: "7d", Key: "7d", Since: dayStart.AddDate(0, 0, -6), Until: now, Bucket: 24 * time.Hour}, true
	case "90d":
		return adminAnalyticsRange{Name: "90d", Key: "90d", Since: dayStart.AddDate(0, 0, -89), Until: now, Bucket: 24 * time.Hour}, true
	default:
		return adminAnalyticsRange{}, false
	}
}

// adminAnalyticsMaxCustomDays caps a custom window. The storage aggregate scans
// the whole window, so an unbounded range would let one request walk the entire
// event history.
const adminAnalyticsMaxCustomDays = 366

func resolveCustomAdminAnalyticsRange(fromRaw, toRaw string, now time.Time) (adminAnalyticsRange, bool) {
	if fromRaw == "" || toRaw == "" {
		return adminAnalyticsRange{}, false
	}
	from, err := time.Parse("2006-01-02", fromRaw)
	if err != nil {
		return adminAnalyticsRange{}, false
	}
	to, err := time.Parse("2006-01-02", toRaw)
	if err != nil || to.Before(from) {
		return adminAnalyticsRange{}, false
	}
	requestedUntil := to.AddDate(0, 0, 1)
	span := requestedUntil.Sub(from)
	if span > adminAnalyticsMaxCustomDays*24*time.Hour {
		return adminAnalyticsRange{}, false
	}
	until := requestedUntil
	if until.After(now) {
		until = now
	}
	if !from.Before(until) {
		return adminAnalyticsRange{}, false
	}
	// The bucket follows the requested days, not the clamped span, so a range
	// that ends today keeps hourly buckets early in the morning instead of
	// switching to minutes because only a few hours have elapsed. A custom
	// range is day-granular, so the smallest useful bucket is one hour.
	bucket := 7 * 24 * time.Hour
	switch {
	case span <= 3*24*time.Hour:
		bucket = time.Hour
	case span <= 120*24*time.Hour:
		bucket = 24 * time.Hour
	}
	return adminAnalyticsRange{
		Name:   "custom",
		Key:    "custom:" + fromRaw + ":" + toRaw,
		Since:  from,
		Until:  until,
		Bucket: bucket,
	}, true
}

func lastBucketPartial(window adminAnalyticsRange) bool {
	last := window.Until.Truncate(window.Bucket)
	return last.Add(window.Bucket).After(window.Until)
}
