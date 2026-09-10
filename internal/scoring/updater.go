package scoring

import (
	"encoding/json"
	"math"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/levels"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

type Scores struct {
	LifetimePoints    int                 `json:"lifetime_points"`
	WeeklyScore       int                 `json:"weekly_score"`
	SpeedBonus        float64             `json:"speed_bonus"`
	ConceptsMastered  int                 `json:"concepts_mastered"`
	CurrentStreak     int                 `json:"current_streak"`
	Level             string              `json:"level"`
	XPTotal           int                 `json:"xp_total"`
	XPToday           int                 `json:"xp_today"`
	DailyXPGoal       int                 `json:"daily_xp_goal"`
	SpacedReps        map[string]float64  `json:"spaced_reps,omitempty"`
	AvgLearningSpeed  float64             `json:"avg_learning_speed"`
	PausedUntil       string              `json:"paused_until,omitempty"`
	// Batch 1: Quiz 150 XP gate signal (MA verbatim).
	XPSinceQuiz       int                 `json:"xp_since_quiz"`
	QuizDue           bool                `json:"quiz_due"`
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
		if p.Status == string(mastery.StatusMastered) {
			mastered++
			if p.MasteredAt != nil && !p.MasteredAt.Before(monday) {
				weeklyCount++
			}
			speedBonus += computeSpeedBonus(p)
		}
	}

	weeklyScore := weeklyCount*100 + int(speedBonus)
	lifetimePoints := mastered * 100
	level := levels.Compute(mastered)
	streak := computeCurrentStreak(progress)
	xpTotal, xpToday, _ := u.repo.GetXP(studentID)
	dailyGoal := 30
	if st, err := u.repo.GetStudent(studentID); err == nil && st != nil && st.DailyXPGoal > 0 {
		dailyGoal = st.DailyXPGoal
	}

	// PR 1.2: per-topic spaced-reps profile (blue-oval darkness) + avg learning speed
	spacedReps := make(map[string]float64)
	var speedSum float64
	var speedCount int
	if speeds, err := u.repo.GetAllTopicSpeeds(studentID); err == nil {
		for cid, ts := range speeds {
			// darkness = stability = repetitions * learning_speed clamped to [0,1] feel
			dark := ts.LearningSpeed * (0.3 + float64(ts.Repetitions)*0.15)
			if dark > 1 {
				dark = 1
			}
			spacedReps[cid] = math.Round(dark*100) / 100
			speedSum += ts.LearningSpeed
			speedCount++
		}
	}
	avgSpeed := 1.0
	if speedCount > 0 {
		avgSpeed = speedSum / float64(speedCount)
	}

	// Batch 1: XP earned since last completed quiz drives the 150 XP gate.
	xpSinceQuiz := xpTotal
	if last, err := u.repo.LastQuizCompletion(studentID); err == nil && last != nil {
		xpSinceQuiz = xpTotal - last.XPTotal
		if xpSinceQuiz < 0 {
			xpSinceQuiz = 0
		}
	}

	return &Scores{
		LifetimePoints:    lifetimePoints,
		WeeklyScore:       weeklyScore,
		SpeedBonus:        speedBonus,
		ConceptsMastered:  mastered,
		CurrentStreak:     streak,
		Level:             level,
		XPTotal:           xpTotal,
		XPToday:           xpToday,
		DailyXPGoal:       dailyGoal,
		SpacedReps:        spacedReps,
		AvgLearningSpeed:  math.Round(avgSpeed*100) / 100,
		PausedUntil:       pausedUntil(u, studentID),
		XPSinceQuiz:       xpSinceQuiz,
		// 150 mirrors engine.QuizGateXP (import cycle forbids sharing).
		QuizDue:           xpSinceQuiz >= 150,
	}, nil
}

// pausedUntil reads settings.pause_until (ISO date) if still in the future.
func pausedUntil(u *Updater, studentID string) string {
	raw, err := u.repo.GetSettings(studentID)
	if err != nil || raw == "" || raw == "{}" {
		return ""
	}
	var cfg struct {
		PauseUntil string `json:"pause_until"`
	}
	if json.Unmarshal([]byte(raw), &cfg) != nil || cfg.PauseUntil == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", cfg.PauseUntil)
	if err != nil || t.Before(time.Now().UTC()) {
		return ""
	}
	return cfg.PauseUntil
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

func weekStart(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	offset := int(weekday) - int(time.Monday)
	start := t.AddDate(0, 0, -offset)
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
}
