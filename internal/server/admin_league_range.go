package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func (s *Server) handleAdminLeagueRange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := r.URL.Query()
	from := strings.TrimSpace(q.Get("from"))
	to := strings.TrimSpace(q.Get("to"))
	group := strings.TrimSpace(q.Get("group"))
	if group != "week" && group != "month" {
		group = "week"
	}
	start, err := time.Parse("2006-01-02", from)
	if err != nil || start.Weekday() != time.Monday {
		writeError(w, "from must be a Monday", 400)
		return
	}
	end, err := time.Parse("2006-01-02", to)
	if err != nil || end.Weekday() != time.Monday || end.Before(start) {
		writeError(w, "to must be a Monday on or after from", 400)
		return
	}
	if end.Sub(start) > 5*365*24*time.Hour {
		writeError(w, "range is too wide", 400)
		return
	}
	store, ok := s.repo.(interface {
		LeagueRange(string, string, string) (*storage.LeagueRange, error)
	})
	if !ok {
		writeError(w, "league range is unavailable", 501)
		return
	}
	out, err := store.LeagueRange(from, to, group)
	if err != nil {
		writeError(w, "failed to read league range", 500)
		return
	}
	writeJSON(w, adminLeagueRangeResponse{
		LeagueRange: out,
		Top:         s.adminLeagueRangeLearners(out.Top),
		Latest:      s.adminLeagueRangeLearners(out.Latest),
	})
}

// A ranked learner row carries a picture for the same reason the account list
// does: an administrator reading a name wants to recognise the account, and a
// second request per row from the browser would be worse than one lookup here.
// The range holds at most twenty rows, so the per-row read stays bounded.
type adminLeagueRangeLearner struct {
	storage.LeagueRangeLearner
	AvatarURL      string             `json:"avatar_url,omitempty"`
	AvatarDicebear *storage.AvatarRef `json:"avatar_dicebear,omitempty"`
	AvatarPreset   *int               `json:"avatar_preset,omitempty"`
}

// adminLeagueRangeResponse shadows the embedded range's two learner lists so
// each row can carry its picture. The other fields marshal unchanged from the
// storage read model.
type adminLeagueRangeResponse struct {
	*storage.LeagueRange
	Top    []adminLeagueRangeLearner `json:"top_range"`
	Latest []adminLeagueRangeLearner `json:"top_latest"`
}

func (s *Server) adminLeagueRangeLearners(in []storage.LeagueRangeLearner) []adminLeagueRangeLearner {
	out := make([]adminLeagueRangeLearner, 0, len(in))
	for _, learner := range in {
		row := adminLeagueRangeLearner{LeagueRangeLearner: learner}
		student, err := s.repo.GetStudent(learner.StudentID)
		if err != nil || student == nil {
			// A missing picture must not hide the ranking, so the row is kept
			// without one and the next account read can explain the failure.
			out = append(out, row)
			continue
		}
		row.AvatarURL, row.AvatarDicebear, row.AvatarPreset = adminAvatar(student.ID, student.AvatarURL, student.Settings)
		out = append(out, row)
	}
	return out
}
