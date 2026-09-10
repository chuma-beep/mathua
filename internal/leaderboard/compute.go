package leaderboard

import (
	"sort"

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
}

func NewComputer(repo storage.Repository) *Computer {
	return &Computer{repo: repo}
}

func (c *Computer) Weekly() ([]Entry, error) {
	rows, err := c.repo.GetWeeklyLeaderboard()
	if err != nil {
		return nil, err
	}
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
	return entries, nil
}
