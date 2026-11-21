package server

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAdminLocationsAggregatesCountryActivity(t *testing.T) {
	e := newAdminEnv(t)
	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	lastWeek := now.AddDate(0, 0, -8).Format("2006-01-02")

	// Two accounts in one country, one active this week and one last week.
	// The point of the daily table is that recording this week must not erase
	// last week's row, so both windows have to survive the same student row.
	if err := e.store.RecordStudentActivity(e.studentID, "NG", "mobile", today, now); err != nil {
		t.Fatalf("record student activity: %v", err)
	}
	if err := e.store.RecordStudentActivity(e.thirdID, "NG", "desktop", lastWeek, now); err != nil {
		t.Fatalf("record third activity: %v", err)
	}

	// A learner may not read the aggregate; it is account data.
	rec := e.do(t, "GET", "/api/admin/locations", e.studentTok, "")
	if rec.Code != 403 {
		t.Fatalf("learner read: %d %s, want 403", rec.Code, rec.Body.String())
	}

	rec = e.do(t, "GET", "/api/admin/locations", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("admin read: %d %s", rec.Code, rec.Body.String())
	}
	var body adminLocationsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Located != 2 {
		t.Fatalf("located = %d, want 2", body.Located)
	}
	if body.Total < 3 {
		t.Fatalf("total = %d, want at least 3 signups", body.Total)
	}
	var ng *adminLocationCountry
	for i := range body.Countries {
		if body.Countries[i].Code == "NG" {
			ng = &body.Countries[i]
		}
	}
	if ng == nil {
		t.Fatalf("NG missing from %+v", body.Countries)
	}
	if ng.Known != 2 || ng.Active7 != 1 || ng.ActivePrev7 != 1 {
		t.Fatalf("NG = %+v, want known 2, active 1, prior 1", *ng)
	}
	if body.Active7 != 1 || body.ActivePrev7 != 1 {
		t.Fatalf("window totals = %d/%d, want 1/1", body.Active7, body.ActivePrev7)
	}
	if body.Visits7 != 1 || body.VisitsPrev7 != 1 {
		t.Fatalf("visits = %d/%d, want 1/1", body.Visits7, body.VisitsPrev7)
	}
	if deviceCount(body.Devices, "mobile") != 1 || deviceCount(body.Devices, "desktop") != 0 {
		t.Fatalf("device split = %+v, want one mobile learner", body.Devices)
	}
	if len(body.Daily) != 30 {
		t.Fatalf("daily has %d days, want 30", len(body.Daily))
	}
	if body.Daily[len(body.Daily)-1].Day != today {
		t.Fatalf("daily ends %s, want %s", body.Daily[len(body.Daily)-1].Day, today)
	}
	if body.GeneratedAt == "" {
		t.Fatal("generated_at is empty")
	}
}

func deviceCount(devices []adminLocationDevice, class string) int {
	for _, d := range devices {
		if d.Device == class {
			return d.Active7
		}
	}
	return 0
}

func TestValidCountryCodeRejectsNonPlaces(t *testing.T) {
	cases := map[string]bool{
		"NG":  true,
		"US":  true,
		"us":  false, // the caller upper-cases first; the validator only accepts stored form
		"XX":  false,
		"T1":  false,
		"":    false,
		"N":   false,
		"NGA": false,
		"1G":  false,
	}
	for code, want := range cases {
		if got := validCountryCode(code); got != want {
			t.Errorf("validCountryCode(%q) = %v, want %v", code, got, want)
		}
	}
}
