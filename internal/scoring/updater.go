package scoring

import (
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
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

type Scores struct {
	LifetimePoints   int     `json:"lifetime_points"`
	WeeklyScore      int     `json:"weekly_score"`
	SpeedBonus       float64 `json:"speed_bonus"`
	ConceptsMastered int     `json:"concepts_mastered"`
	CurrentStreak    int     `json:"current_streak"`
	Level            string  `json:"level"`
	XPTotal          int     `json:"xp_total"`
	XPToday          int     `json:"xp_today"`
}

type Updater struct {
	dag  *concepts.DAG
	repo storage.Repository
}

func NewUpdater(dag *concepts.DAG, repo storage.Repository) *Updater {
	return &Updater{dag: dag, repo: repo}
}

func (u *Updater) Compute(studentID string) (*Scores, error) {
	progress, err := u.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, err
	}
	monday := weekStart(time.Now().UTC())

	mastered := 0
	weeklyCount := 0
	speedBonus := 0.0

	for _, p := range progress {
		if p.Status == "MASTERED" {
			mastered++
			if p.MasteredAt != nil && !p.MasteredAt.Before(monday) {
				weeklyCount++
			}
			speedBonus += computeSpeedBonus(p)
		}
	}

	weeklyScore := weeklyCount*100 + int(speedBonus)
	lifetimePoints := mastered * 100
	level := computeLevel(mastered)
	streak := computeCurrentStreak(progress)
	xpTotal, xpToday, _ := u.repo.GetXP(studentID)

	return &Scores{
		LifetimePoints:   lifetimePoints,
		WeeklyScore:      weeklyScore,
		SpeedBonus:       speedBonus,
		ConceptsMastered: mastered,
		CurrentStreak:    streak,
		Level:            level,
		XPTotal:          xpTotal,
		XPToday:          xpToday,
	}, nil
}

func computeSpeedBonus(p *storage.ConceptProgress) float64 {
	if p.AvgResponseTime <= 0 {
		return 0
	}
	bonus := 50.0 / p.AvgResponseTime
	if bonus > 25 {
		bonus = 25
	}
	return bonus
}

func computeCurrentStreak(progress map[string]*storage.ConceptProgress) int {
	streak := 0
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	check := today
	for {
		day := check
		found := false
		for _, p := range progress {
			if p.LastAttempted != nil {
				attemptDay := time.Date(p.LastAttempted.Year(), p.LastAttempted.Month(), p.LastAttempted.Day(), 0, 0, 0, 0, time.UTC)
				if attemptDay.Equal(day) {
					found = true
					break
				}
			}
		}
		if !found {
			break
		}
		streak++
		check = check.AddDate(0, 0, -1)
		if check.Before(today.AddDate(0, 0, -365)) {
			break
		}
	}
	return streak
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

func weekStart(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	offset := int(weekday) - int(time.Monday)
	start := t.AddDate(0, 0, -offset)
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
}
