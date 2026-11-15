package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestAdminAnalyticsRequiresPermissionAndReturnsWindow(t *testing.T) {
	env := newAdminEnv(t)

	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/admin/analytics?range=7d", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous analytics = %d, want 401", rec.Code)
	}

	req := httptest.NewRequest("GET", "/api/admin/analytics?range=7d", nil)
	req.Header.Set("Authorization", "Bearer "+env.adminTok)
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("analytics = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Range             string                       `json:"range"`
		BucketSeconds     int64                        `json:"bucket_seconds"`
		PartialLastBucket bool                         `json:"partial_last_bucket"`
		Accounts          storage.AdminAccountCounts   `json:"accounts"`
		Totals            *storage.AdminMetricTotals   `json:"totals"`
		Activity          []storage.AdminActivityPoint `json:"activity"`
		Previous          *storage.AdminMetricTotals   `json:"previous"`
		PreviousFrom      *time.Time                   `json:"previous_from"`
		PreviousTo        *time.Time                   `json:"previous_to"`
		From              time.Time                    `json:"from"`
		To                time.Time                    `json:"to"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode analytics: %v", err)
	}
	if out.Range != "7d" || out.BucketSeconds != 86400 {
		t.Fatalf("range = %q bucket = %d, want 7d/86400", out.Range, out.BucketSeconds)
	}
	if len(out.Activity) != 7 {
		t.Fatalf("activity buckets = %d, want 7", len(out.Activity))
	}
	if !out.PartialLastBucket {
		t.Fatal("7d window should mark the current day partial")
	}
	if out.Accounts.Total < 3 {
		t.Fatalf("accounts total = %d, want at least the three test accounts", out.Accounts.Total)
	}
	if out.Totals == nil || out.Totals.DiagnosticCohortRegistered < 3 {
		t.Fatalf("diagnostic cohort denominator missing from response: %+v", out.Totals)
	}
	if out.Previous == nil || out.PreviousFrom == nil || out.PreviousTo == nil {
		t.Fatal("previous equal window should be present")
	}
	if !out.PreviousTo.Equal(out.From) {
		t.Fatalf("previous window ends at %v, current starts at %v", out.PreviousTo, out.From)
	}
	if out.PreviousTo.Sub(*out.PreviousFrom) != out.To.Sub(out.From) {
		t.Fatalf("previous span = %v, current span = %v", out.PreviousTo.Sub(*out.PreviousFrom), out.To.Sub(out.From))
	}

	req = httptest.NewRequest("GET", "/api/admin/analytics?range=nonsense", nil)
	req.Header.Set("Authorization", "Bearer "+env.adminTok)
	rec = httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid range = %d, want 400", rec.Code)
	}
}

func TestAdminAnalyticsLiveReturnsCurrentCounts(t *testing.T) {
	env := newAdminEnv(t)
	req := httptest.NewRequest("GET", "/api/admin/analytics/live", nil)
	req.Header.Set("Authorization", "Bearer "+env.adminTok)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("live analytics = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		OnlineNow           int `json:"online_now"`
		OnlineWindowSeconds int `json:"online_window_seconds"`
		PracticingNow       int `json:"practicing_now"`
		Answers15m          int `json:"answers_15m"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode live analytics: %v", err)
	}
	if out.OnlineWindowSeconds != 300 {
		t.Fatalf("online window = %d, want 300", out.OnlineWindowSeconds)
	}
	if out.OnlineNow < 0 || out.PracticingNow < 0 || out.Answers15m < 0 {
		t.Fatalf("live counts must not be negative: %+v", out)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("live cache-control = %q, want no-store", got)
	}
}

func TestAdminAnalyticsSectionsOmitOtherSections(t *testing.T) {
	env := newAdminEnv(t)
	req := httptest.NewRequest("GET", "/api/admin/analytics?range=7d&sections=activity", nil)
	req.Header.Set("Authorization", "Bearer "+env.adminTok)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("section analytics = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode sections: %v", err)
	}
	if _, ok := body["activity"]; !ok {
		t.Fatal("activity section is missing")
	}
	if _, ok := body["accounts"]; ok {
		t.Fatal("accounts should be omitted when only activity is requested")
	}
	if _, ok := body["totals"]; ok {
		t.Fatal("totals should be omitted when only activity is requested")
	}
}

func TestAdminAnalyticsViewMapsToSections(t *testing.T) {
	env := newAdminEnv(t)
	req := httptest.NewRequest("GET", "/api/admin/analytics?range=7d&view=usage", nil)
	req.Header.Set("Authorization", "Bearer "+env.adminTok)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("view analytics = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var body struct {
		View     string   `json:"view"`
		Sections []string `json:"sections"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode view: %v", err)
	}
	if body.View != "usage" || len(body.Sections) != 4 {
		t.Fatalf("view = %q sections = %v, want usage with four sections", body.View, body.Sections)
	}
}

func TestAdminAnalyticsRejectsUnknownSection(t *testing.T) {
	env := newAdminEnv(t)
	req := httptest.NewRequest("GET", "/api/admin/analytics?range=7d&sections=secrets", nil)
	req.Header.Set("Authorization", "Bearer "+env.adminTok)
	rec := httptest.NewRecorder()
	env.mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unknown section = %d, want 400", rec.Code)
	}
}

func TestAdminAnalyticsCacheTTL(t *testing.T) {
	cache := newAdminAnalyticsCache()
	now := time.Now()
	if _, flight, leader, cached := cache.begin("key", now); !leader || cached || flight == nil {
		t.Fatal("first request should lead a flight")
	}
	cache.finish("key", []byte("payload"), nil, now)
	if got, flight, leader, cached := cache.begin("key", now.Add(29*time.Second)); !cached || leader || flight != nil || string(got) != "payload" {
		t.Fatalf("cache at 29s = %q, flight=%v leader=%v cached=%v; want payload hit", got, flight, leader, cached)
	}
	if _, flight, leader, cached := cache.begin("key", now.Add(31*time.Second)); !leader || cached || flight == nil {
		t.Fatal("expired entry should start a new flight")
	}
}

func TestAdminAnalyticsCacheDeduplicatesConcurrentMisses(t *testing.T) {
	cache := newAdminAnalyticsCache()
	now := time.Now()
	_, first, leader, cached := cache.begin("key", now)
	if !leader || cached {
		t.Fatal("first request should lead")
	}
	_, second, leader, cached := cache.begin("key", now)
	if leader || cached {
		t.Fatal("second request should wait on the in-flight computation")
	}
	if first != second {
		t.Fatal("both requests should share one flight")
	}
	cache.finish("key", []byte("shared"), nil, now)
	if got, err := cache.wait(context.Background(), second); err != nil || string(got) != "shared" {
		t.Fatalf("waiter got %q, %v; want shared, nil", got, err)
	}
	if _, flight, leader, cached := cache.begin("key", now.Add(time.Second)); leader || !cached || flight != nil {
		t.Fatal("completed flight should be a cache hit")
	}
}
