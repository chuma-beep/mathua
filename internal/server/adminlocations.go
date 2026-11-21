package server

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// adminLocationCountry is one row of the map and the table under it. Growth is
// expressed twice because the two questions are different: active counts say
// where people are practising in the selected window, new counts say where the
// newest accounts came from. A country can grow in activity without a single
// signup. The field names carry no length because the window is a query
// parameter, not a constant.
type adminLocationCountry struct {
	Code       string `json:"code"`
	Known      int    `json:"known"`
	Active     int    `json:"active"`
	ActivePrev int    `json:"active_prev"`
	New        int    `json:"new"`
	NewPrev    int    `json:"new_prev"`
}

type adminLocationDay struct {
	Day    string `json:"day"`
	Active int    `json:"active"`
}

// adminLocationDevice is one class in the audience split. Active counts
// distinct learners, not rows, so one learner using two devices counts once in
// each class they used but never twice in the same class.
type adminLocationDevice struct {
	Device     string `json:"device"`
	Active     int    `json:"active"`
	ActivePrev int    `json:"active_prev"`
}

// adminLocationDeviceDay is one day of the per-class trend. The three numbers
// are independent distinct-learner counts, so they can sum above the day's
// active total when someone uses more than one class that day.
type adminLocationDeviceDay struct {
	Day     string `json:"day"`
	Mobile  int    `json:"mobile"`
	Tablet  int    `json:"tablet"`
	Desktop int    `json:"desktop"`
}

type adminLocationsResponse struct {
	// Days is the length of the selected window. Every active, new, visits and
	// device count on the page refers to it, and daily holds exactly this many
	// points, so a client never has to infer the range from a field name.
	Days int `json:"days"`
	// From and To are the first and last day of the window. The page prints
	// them when the operator picked a range, so a label can never describe a
	// different window from the numbers under it.
	From        string                   `json:"from"`
	To          string                   `json:"to"`
	Countries   []adminLocationCountry   `json:"countries"`
	Located     int                      `json:"located"`
	Total       int                      `json:"total"`
	Active      int                      `json:"active"`
	ActivePrev  int                      `json:"active_prev"`
	New         int                      `json:"new"`
	NewPrev     int                      `json:"new_prev"`
	Visits      int                      `json:"visits"`
	VisitsPrev  int                      `json:"visits_prev"`
	Devices     []adminLocationDevice    `json:"devices"`
	DeviceDaily []adminLocationDeviceDay `json:"device_daily"`
	Daily       []adminLocationDay       `json:"daily"`
	GeneratedAt string                   `json:"generated_at"`
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
	// The page offers 7, 30 and 90 day presets and an arbitrary range. The
	// endpoint accepts `days`, which ends today, or `from` and `to`, never
	// both: a result has one window, and accepting two descriptions of it
	// invites the page and the series to disagree.
	now := time.Now().UTC()
	midnight := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}
	start := midnight(now).AddDate(0, 0, -29)
	end := midnight(now)
	days := 30
	if raw := r.URL.Query().Get("days"); raw != "" {
		n, parseErr := strconv.Atoi(raw)
		if parseErr != nil || n < 1 || n > 365 {
			writeError(w, "days must be between 1 and 365", http.StatusBadRequest)
			return
		}
		days = n
		start = midnight(now).AddDate(0, 0, -(days - 1))
	}
	fromRaw, toRaw := r.URL.Query().Get("from"), r.URL.Query().Get("to")
	if fromRaw != "" || toRaw != "" {
		if r.URL.Query().Get("days") != "" {
			writeError(w, "give either days or from and to, not both", http.StatusBadRequest)
			return
		}
		if fromRaw == "" || toRaw == "" {
			writeError(w, "from and to must be given together", http.StatusBadRequest)
			return
		}
		var parseErr error
		start, parseErr = time.Parse("2006-01-02", fromRaw)
		if parseErr != nil {
			writeError(w, "from must be a date as YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		end, parseErr = time.Parse("2006-01-02", toRaw)
		if parseErr != nil {
			writeError(w, "to must be a date as YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		if start.After(end) {
			writeError(w, "from must not be after to", http.StatusBadRequest)
			return
		}
		// Tomorrow's date is not a window, and a future end would make the
		// zero-filled series claim days that have not happened.
		if end.After(now) {
			writeError(w, "to must not be in the future", http.StatusBadRequest)
			return
		}
		days = int(end.Sub(start).Hours()/24) + 1
		if days > 365 {
			writeError(w, "window must be 365 days or fewer", http.StatusBadRequest)
			return
		}
	}
	window := storage.LocationWindow{
		RecentFrom: start.Format("2006-01-02"),
		PriorFrom:  start.AddDate(0, 0, -days).Format("2006-01-02"),
		DailyFrom:  start.Format("2006-01-02"),
		Today:      end.Format("2006-01-02"),
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
	for code := range locations.Active {
		codes[code] = struct{}{}
	}
	for code := range locations.ActivePrev {
		codes[code] = struct{}{}
	}

	countries := make([]adminLocationCountry, 0, len(codes))
	for code := range codes {
		countries = append(countries, adminLocationCountry{
			Code:       code,
			Known:      locations.Known[code],
			Active:     locations.Active[code],
			ActivePrev: locations.ActivePrev[code],
			New:        locations.New[code],
			NewPrev:    locations.NewPrev[code],
		})
	}
	// Known accounts first, then the busiest in this window, then the code, so
	// two reads of unchanged data always render in the same order.
	sort.Slice(countries, func(i, j int) bool {
		if countries[i].Known != countries[j].Known {
			return countries[i].Known > countries[j].Known
		}
		if countries[i].Active != countries[j].Active {
			return countries[i].Active > countries[j].Active
		}
		return countries[i].Code < countries[j].Code
	})

	byDay := make(map[string]int, len(locations.Daily))
	for _, d := range locations.Daily {
		byDay[d.Day] = d.Active
	}
	daily := make([]adminLocationDay, 0, days)
	// Zero-fill the series. A missing day is a day with no activity, and a
	// sparkline that joins two active days across a gap draws a lie.
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		daily = append(daily, adminLocationDay{Day: key, Active: byDay[key]})
	}

	newAccounts, newPrevAccounts := 0, 0
	for _, n := range locations.New {
		newAccounts += n
	}
	for _, n := range locations.NewPrev {
		newPrevAccounts += n
	}
	// A fixed class order keeps the split from reshuffling between reads.
	devices := make([]adminLocationDevice, 0, 3)
	for _, class := range []string{"mobile", "tablet", "desktop"} {
		if locations.Device[class] == 0 && locations.DevicePrev[class] == 0 {
			continue
		}
		devices = append(devices, adminLocationDevice{
			Device:     class,
			Active:     locations.Device[class],
			ActivePrev: locations.DevicePrev[class],
		})
	}

	deviceDaily := make([]adminLocationDeviceDay, 0, len(locations.DeviceDaily))
	for _, d := range locations.DeviceDaily {
		deviceDaily = append(deviceDaily, adminLocationDeviceDay{
			Day:     d.Day,
			Mobile:  d.Mobile,
			Tablet:  d.Tablet,
			Desktop: d.Desktop,
		})
	}

	w.Header().Set("Cache-Control", "private, max-age=60")
	writeJSON(w, adminLocationsResponse{
		Days:        days,
		From:        window.RecentFrom,
		To:          window.Today,
		Countries:   countries,
		Located:     locations.Located,
		Total:       locations.Total,
		Active:      locations.ActiveTotal,
		ActivePrev:  locations.ActivePrevTotal,
		New:         newAccounts,
		NewPrev:     newPrevAccounts,
		Visits:      locations.Visits,
		VisitsPrev:  locations.VisitsPrev,
		Devices:     devices,
		DeviceDaily: deviceDaily,
		Daily:       daily,
		GeneratedAt: now.Format(time.RFC3339),
	})
}
