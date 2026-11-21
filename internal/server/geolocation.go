package server

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// The Cloudflare Worker resolves the caller's country from request.cf and
// overwrites this header on every proxied /api request. The origin never reads
// a forwarding header from the client, and the Worker deletes any inbound copy
// before setting its own, so a caller cannot choose its country.
const countryHeader = "X-Mathua-Country"

// noteActivity records the edge-reported country, a coarse device class and the
// fact that a signed-in learner was active today, at most once per UTC day.
//
// The throttle is in memory on purpose. A restart costs one extra write per
// active learner, which is cheaper than a database read on every authenticated
// request, and the write itself is an idempotent upsert of today's row. The
// map holds one short string per account that has been seen; it is not a
// cache of personal data and it never holds an IP address or a user agent.
func (s *Server) noteActivity(r *http.Request, studentID string) {
	if s.repo == nil || studentID == "" {
		return
	}
	code := strings.ToUpper(strings.TrimSpace(r.Header.Get(countryHeader)))
	if !validCountryCode(code) {
		code = ""
	}
	device := deviceClass(r.UserAgent())
	// A request with neither fact adds nothing: the daily row would carry two
	// empty strings and the admin map would gain a visit it cannot place.
	if code == "" && device == "" {
		return
	}

	now := time.Now().UTC()
	day := now.Format("2006-01-02")
	if last, ok := s.countrySeen.Load(studentID); ok && last == day {
		return
	}
	// Store before the write so two concurrent first requests of the day do
	// not both reach the database in the common case. A failed write clears
	// the entry so the next request retries instead of losing the day.
	s.countrySeen.Store(studentID, day)
	if err := s.repo.RecordStudentActivity(studentID, code, device, day, now); err != nil {
		log.Printf("record activity for %s: %v", studentID, err)
		s.countrySeen.Delete(studentID)
	}
}

// deviceClass turns a user agent into one of three coarse classes. The agent
// itself is never stored: a class answers "are learners on phones" without
// keeping a fingerprint, and a wrong guess costs a bar in a chart, not a
// learner's data.
func deviceClass(ua string) string {
	s := strings.ToLower(strings.TrimSpace(ua))
	if s == "" {
		return ""
	}
	switch {
	case strings.Contains(s, "ipad"), strings.Contains(s, "tablet"):
		return "tablet"
	case strings.Contains(s, "android") && !strings.Contains(s, "mobile"):
		// Android tablets omit "Mobile" from their user agent.
		return "tablet"
	case strings.Contains(s, "mobi"), strings.Contains(s, "iphone"), strings.Contains(s, "ipod"):
		return "mobile"
	default:
		return "desktop"
	}
}

// validCountryCode accepts the ISO 3166-1 alpha-2 shape and rejects the codes
// Cloudflare uses for "unknown" and Tor exits. Those are not places, and
// counting them as countries would put fictional markers on the map.
func validCountryCode(code string) bool {
	if len(code) != 2 {
		return false
	}
	for i := 0; i < len(code); i++ {
		if code[i] < 'A' || code[i] > 'Z' {
			return false
		}
	}
	return code != "XX" && code != "T1"
}
