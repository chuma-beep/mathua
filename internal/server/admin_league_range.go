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
	writeJSON(w, out)
}
