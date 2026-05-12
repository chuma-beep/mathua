package leaderboard

import (
	"sort"

	"github.com/chuma-beep/mathua/internal/storage"
)

var levels = []struct {
	Min   int
	Title string
}{
	{0, "Novice"},
	{32, "Apprentice"},
	{64, "Student"},
	{96, "Scholar"},
	{128, "Adept"},
	{160, "Expert"},
	{192, "Master"},
	{224, "Grandmaster"},
	{256, "Math Architect"},
}

type Entry struct {
	Rank     int    `json:"rank"`
	Name     string `json:"name"`
	Mastered int    `json:"mastered"`
	Streak   int    `json:"streak"`
	Level    string `json:"level"`
	Score    int    `json:"score"`
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
		entries = append(entries, Entry{
			Name:     r.Name,
			Mastered: r.TotalMastered,
			Level:    computeLevel(r.TotalMastered),
			Streak:   0,
			Score:    score,
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

func computeLevel(mastered int) string {
	title := levels[0].Title
	for _, l := range levels {
		if mastered >= l.Min {
			title = l.Title
		}
	}
	return title
}
