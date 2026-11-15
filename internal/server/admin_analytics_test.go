package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		Activity          []storage.AdminActivityPoint `json:"activity"`
		Previous          *storage.AdminMetricTotals   `json:"previous"`
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
	if out.Previous == nil {
		t.Fatal("previous equal window should be present")
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
