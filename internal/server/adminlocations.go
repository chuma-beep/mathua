package server

import (
	"net/http"
	"sort"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// adminLocationCountry is one row of the map and the table under it. Growth is
// expressed twice because the two questions are different: active counts say
// where people are practising this week, new counts say where the newest
// accounts came from. A country can grow in activity without a single signup.
type adminLocationCountry struct {
	Code        string `json:"code"`
	Known       int    `json:"known"`
	Active7     int    `json:"active_7d"`
	ActivePrev7 int    `json:"active_prev_7d"`
	New7        int    `json:"new_7d"`
	NewPrev7    int    `json:"new_prev_7d"`
}

type adminLocationDay struct {
	Day    string `json:"day"`
	Active int    `json:"active"`
}

// adminLocationDevice is one class in the audience split. active_7d counts
// distinct learners, not rows, so one learner using two devices counts once in
// each class they used but never twice in the same class.
type adminLocationDevice struct {
	Device      string `json:"device"`
	Active7     int    `json:"active_7d"`
	ActivePrev7 int    `json:"active_prev_7d"`
}

type adminLocationsResponse struct {
	Countries   []adminLocationCountry `json:"countries"`
	Located     int                    `json:"located"`
	Total       int                    `json:"total"`
	Active7     int                    `json:"active_7d"`
	ActivePrev7 int                    `json:"active_prev_7d"`
	New7        int                    `json:"new_7d"`
	NewPrev7    int                    `json:"new_prev_7d"`
	Visits7     int                    `json:"visits_7d"`
	VisitsPrev7 int                    `json:"visits_prev_7d"`
	Visits30    int                    `json:"visits_30d"`
	Devices     []adminLocationDevice  `json:"devices"`
	Daily       []adminLocationDay     `json:"daily"`
	GeneratedAt string                 `json:"generated_at"`
}

// handleAdminLocations answers "where are my learners" with country granularity.
//
// The edge resolves the country from the connection and the origin stores only
// that code. There is no city, no coordinate and no IP address in the aggregate,
// so the endpoint cannot be turned into a tracker even by an administrator. The
// page says the country is approximate because a VPN or a carrier's egress can
// put a learner somewhere they are not.
func (s *Server) handleAdminLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	window := storage.LocationWindow{
		RecentFrom: now.AddDate(0, 0, -6).Format("2006-01-02"),
		PriorFrom:  now.AddDate(0, 0, -13).Format("2006-01-02"),
		DailyFrom:  now.AddDate(0, 0, -29).Format("2006-01-02"),
		Today:      today,
	}
	locations, err := s.repo.StudentLocations(window)
	if err != nil {
		writeError(w, "failed to load locations", http.StatusInternalServerError)
		return
	}

	codes := make(map[string]struct{}, len(locations.Known))
	for code := range locations.Known {
		codes[code] = struct{}{}
	}
	for code := range locations.Active7 {
		codes[code] = struct{}{}
	}
	for code := range locations.ActivePrev7 {
		codes[code] = struct{}{}
	}

	countries := make([]adminLocationCountry, 0, len(codes))
	for code := range codes {
		countries = append(countries, adminLocationCountry{
			Code:        code,
			Known:       locations.Known[code],
			Active7:     locations.Active7[code],
			ActivePrev7: locations.ActivePrev7[code],
			New7:        locations.New7[code],
			NewPrev7:    locations.NewPrev7[code],
		})
	}
	// Known accounts first, then the busiest this week, then the code, so two
	// reads of unchanged data always render in the same order.
	sort.Slice(countries, func(i, j int) bool {
		if countries[i].Known != countries[j].Known {
			return countries[i].Known > countries[j].Known
		}
		if countries[i].Active7 != countries[j].Active7 {
			return countries[i].Active7 > countries[j].Active7
		}
		return countries[i].Code < countries[j].Code
	})

	byDay := make(map[string]int, len(locations.Daily))
	for _, d := range locations.Daily {
		byDay[d.Day] = d.Active
	}
	daily := make([]adminLocationDay, 0, 30)
	start, err := time.Parse("2006-01-02", window.DailyFrom)
	if err != nil {
		writeError(w, "failed to load locations", http.StatusInternalServerError)
		return
	}
	end, err := time.Parse("2006-01-02", window.Today)
	if err != nil {
		writeError(w, "failed to load locations", http.StatusInternalServerError)
		return
	}
	// Zero-fill the series. A missing day is a day with no activity, and a
	// sparkline that joins two active days across a gap draws a lie.
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		daily = append(daily, adminLocationDay{Day: key, Active: byDay[key]})
	}

	new7, newPrev7 := 0, 0
	for _, n := range locations.New7 {
		new7 += n
	}
	for _, n := range locations.NewPrev7 {
		newPrev7 += n
	}
	visits30 := 0
	for _, d := range daily {
		visits30 += d.Active
	}
	// A fixed class order keeps the split from reshuffling between reads.
	devices := make([]adminLocationDevice, 0, 3)
	for _, class := range []string{"mobile", "tablet", "desktop"} {
		if locations.Device7[class] == 0 && locations.DevicePrev7[class] == 0 {
			continue
		}
		devices = append(devices, adminLocationDevice{
			Device:      class,
			Active7:     locations.Device7[class],
			ActivePrev7: locations.DevicePrev7[class],
		})
	}

	w.Header().Set("Cache-Control", "private, max-age=60")
	writeJSON(w, adminLocationsResponse{
		Countries:   countries,
		Located:     locations.Located,
		Total:       locations.Total,
		Active7:     locations.Active7Total,
		ActivePrev7: locations.ActivePrev7Total,
		New7:        new7,
		NewPrev7:    newPrev7,
		Visits7:     locations.Visits7,
		VisitsPrev7: locations.VisitsPrev7,
		Visits30:    visits30,
		Devices:     devices,
		Daily:       daily,
		GeneratedAt: now.Format(time.RFC3339),
	})
}
