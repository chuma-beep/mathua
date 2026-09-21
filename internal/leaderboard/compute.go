package leaderboard

import (
	"sort"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/levels"
	"github.com/chuma-beep/mathua/internal/storage"
)

type Entry struct {
	Rank           int                `json:"rank"`
	StudentID      string             `json:"student_id"`
	Name           string             `json:"name"`
	Username       string             `json:"username,omitempty"`
	AvatarURL      string             `json:"avatar_url,omitempty"`
	AvatarDicebear *storage.AvatarRef `json:"avatar_dicebear,omitempty"`
	AvatarCustom   bool               `json:"avatar_custom,omitempty"`
	AvatarVersion  int                `json:"avatar_version,omitempty"`
	Mastered       int                `json:"mastered"`
	Streak         int                `json:"streak"`
	Level          string             `json:"level"`
	Score          int                `json:"score"`
}

type Computer struct {
	repo storage.Repository

	// Weekly board is a full students×progress scan per call; cache it
	// briefly so leaderboard polling under load doesn't queue behind
	// the write path on the SQLite pool.
	mu      sync.Mutex
	cached  []Entry
	filled  time.Time
	ttl     time.Duration
}

// DefaultWeeklyTTL is the cache lifetime for the weekly board. Weekly
// granularity data changes slowly; 30s keeps it fresh while absorbing
// poll storms from concurrent readers.
const DefaultWeeklyTTL = 30 * time.Second

func NewComputer(repo storage.Repository) *Computer {
	return NewComputerWithTTL(repo, DefaultWeeklyTTL)
}

func NewComputerWithTTL(repo storage.Repository, ttl time.Duration) *Computer {
	return &Computer{repo: repo, ttl: ttl}
}

func (c *Computer) Weekly() ([]Entry, error) {
	c.mu.Lock()
	if c.cached != nil && time.Since(c.filled) < c.ttl {
		out := c.cached
		c.mu.Unlock()
		return out, nil
	}
	c.mu.Unlock()

	rows, err := c.repo.GetWeeklyLeaderboard()
	if err != nil {
		return nil, err
	}
	entries := buildEntries(rows)

	c.mu.Lock()
	c.cached = entries
	c.filled = time.Now()
	c.mu.Unlock()
	return entries, nil
}

func buildEntries(rows []storage.LeaderboardRow) []Entry {
	entries := make([]Entry, 0, len(rows))
	for _, r := range rows {
		score := r.WeeklyMastered * 100
		custom, pick, version := storage.AvatarBits(r.AvatarSettings)
		entries = append(entries, Entry{
			StudentID:      r.StudentID,
			Name:           r.Name,
			Username:       r.Username,
			AvatarURL:      r.AvatarURL,
			AvatarDicebear: pick,
			AvatarCustom:   custom,
			AvatarVersion:  version,
			Mastered:       r.TotalMastered,
			Level:          levels.Compute(r.TotalMastered),
			Streak:         0,
			Score:          score,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Name < entries[j].Name
	})
	for i := range entries {
		entries[i].Rank = i + 1
	}
	return entries
}
