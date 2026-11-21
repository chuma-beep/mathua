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

	rec = e.do(t, "GET", "/api/admin/locations?days=7", e.adminTok, "")
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
	if ng.Known != 2 || ng.Active != 1 || ng.ActivePrev != 1 {
		t.Fatalf("NG = %+v, want known 2, active 1, prior 1", *ng)
	}
	if body.Active != 1 || body.ActivePrev != 1 {
		t.Fatalf("window totals = %d/%d, want 1/1", body.Active, body.ActivePrev)
	}
	if body.Visits != 1 || body.VisitsPrev != 1 {
		t.Fatalf("visits = %d/%d, want 1/1", body.Visits, body.VisitsPrev)
	}
	if body.Days != 7 {
		t.Fatalf("days = %d, want 7", body.Days)
	}
	if deviceCount(body.Devices, "mobile") != 1 || deviceCount(body.Devices, "desktop") != 0 {
		t.Fatalf("device split = %+v, want one mobile learner", body.Devices)
	}
	if len(body.Daily) != 7 {
		t.Fatalf("daily has %d days, want 7", len(body.Daily))
	}
	if len(body.DeviceDaily) != 7 {
		t.Fatalf("device daily has %d days, want 7", len(body.DeviceDaily))
	}
	deviceDays := make(map[string]adminLocationDeviceDay, len(body.DeviceDaily))
	for _, d := range body.DeviceDaily {
		deviceDays[d.Day] = d
	}
	if deviceDays[today].Mobile != 1 {
		t.Fatalf("mobile on %s = %d, want 1", today, deviceDays[today].Mobile)
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
			return d.Active
		}
	}
	return 0
}

// The whole page moves with the window: the day series, the active and visit
// totals, the prior comparison and the device trend. A range that only shortened
// the chart would leave the headline cards claiming a week while the chart
// showed a quarter.
func TestAdminLocationsHonoursTheWindow(t *testing.T) {
	e := newAdminEnv(t)
	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	eightDaysAgo := now.AddDate(0, 0, -8).Format("2006-01-02")

	if err := e.store.RecordStudentActivity(e.studentID, "NG", "mobile", today, now); err != nil {
		t.Fatalf("record student activity: %v", err)
	}
	if err := e.store.RecordStudentActivity(e.thirdID, "NG", "desktop", eightDaysAgo, now); err != nil {
		t.Fatalf("record third activity: %v", err)
	}

	cases := []struct {
		query string
		days  int
	}{
		{"?days=7", 7},
		{"?days=90", 90},
	}
	for _, tc := range cases {
		rec := e.do(t, "GET", "/api/admin/locations"+tc.query, e.adminTok, "")
		if rec.Code != 200 {
			t.Fatalf("%s: %d %s", tc.query, rec.Code, rec.Body.String())
		}
		var body adminLocationsResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s decode: %v", tc.query, err)
		}
		if body.Days != tc.days {
			t.Fatalf("%s: days = %d, want %d", tc.query, body.Days, tc.days)
		}
		if len(body.Daily) != tc.days || len(body.DeviceDaily) != tc.days {
			t.Fatalf("%s: daily/device daily = %d/%d, want %d each",
				tc.query, len(body.Daily), len(body.DeviceDaily), tc.days)
		}
		if body.Daily[len(body.Daily)-1].Day != today {
			t.Fatalf("%s: series ends %s, want %s", tc.query, body.Daily[len(body.Daily)-1].Day, today)
		}
		// Eight days ago is outside a seven-day window and inside its prior
		// window; inside a ninety-day window it counts as activity.
		if tc.days == 7 {
			if body.Active != 1 || body.ActivePrev != 1 {
				t.Fatalf("7 day totals = %d/%d, want 1/1", body.Active, body.ActivePrev)
			}
			if deviceCount(body.Devices, "desktop") != 0 {
				t.Fatalf("7 day devices = %+v, want desktop only in the prior window", body.Devices)
			}
		} else if body.Active != 2 || body.ActivePrev != 0 {
			t.Fatalf("90 day totals = %d/%d, want 2/0", body.Active, body.ActivePrev)
		}
	}

	for _, query := range []string{"?days=0", "?days=366", "?days=week"} {
		rec := e.do(t, "GET", "/api/admin/locations"+query, e.adminTok, "")
		if rec.Code != 400 {
			t.Fatalf("%s: %d %s, want 400", query, rec.Code, rec.Body.String())
		}
	}
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
