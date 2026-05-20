package engine

import (
	"fmt"

	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

type ConceptDetail struct {
	Concept       ConceptInfo      `json:"concept"`
	Lesson        *LessonInfo      `json:"lesson,omitempty"`
	Prerequisites []PrereqInfo     `json:"prerequisites"`
	Unlocked      bool             `json:"unlocked"`
	Progress      *ProgressInfo    `json:"progress,omitempty"`
}

type ConceptInfo struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
}

type LessonInfo struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Concepts []string `json:"concepts"`
}

type PrereqInfo struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	Status     string  `json:"status"`
	MasteryPct float64 `json:"mastery_pct"`
}

type ProgressInfo struct {
	Status         string  `json:"status"`
	Streak         int     `json:"streak"`
	RequiredStreak int     `json:"required_streak"`
	MasteryPct     float64 `json:"mastery_pct"`
}

func (e *Engine) ConceptDetail(studentID, conceptID string) (*ConceptDetail, error) {
	c := e.dag.Concept(conceptID)
	if c == nil {
		return nil, fmt.Errorf("concept %q not found", conceptID)
	}

	var progress map[string]*storage.ConceptProgress
	if studentID != "" {
		progress, _ = e.repo.GetAllProgress(studentID)
	}

	detail := &ConceptDetail{
		Concept: ConceptInfo{
			ID:        c.ID,
			Label:     c.Label,
			Domain:    c.Domain,
			Subdomain: c.Subdomain,
		},
	}

	if e.ll != nil {
		if lesson := e.ll.Lesson(conceptID); lesson != nil {
			detail.Lesson = &LessonInfo{
				Title:    lesson.Title,
				Body:     lesson.Body,
				Concepts: lesson.Concepts,
			}
		}
	}

	detail.Unlocked = true
	for _, pid := range c.Prerequisites {
		prereq := e.dag.Concept(pid)
		p, hasProgress := progress[pid]
		status := string(mastery.StatusUnseen)
		masteryPct := 0.0

		if hasProgress {
			status = p.Status
			if prereq != nil && prereq.MasteryThreshold.Streak > 0 {
				masteryPct = float64(p.Streak) / float64(prereq.MasteryThreshold.Streak)
				if masteryPct > 1 {
					masteryPct = 1
				}
			}
		}

		if status != string(mastery.StatusMastered) {
			detail.Unlocked = false
		}

		label := pid
		if prereq != nil {
			label = prereq.Label
		}

		detail.Prerequisites = append(detail.Prerequisites, PrereqInfo{
			ID:         pid,
			Label:      label,
			Status:     status,
			MasteryPct: masteryPct,
		})
	}

	if p, ok := progress[conceptID]; ok {
		streak := p.Streak
		reqStreak := c.MasteryThreshold.Streak
		masteryPct := 0.0
		if reqStreak > 0 {
			masteryPct = float64(streak) / float64(reqStreak)
			if masteryPct > 1 {
				masteryPct = 1
			}
		}
		detail.Progress = &ProgressInfo{
			Status:         p.Status,
			Streak:         streak,
			RequiredStreak: reqStreak,
			MasteryPct:     masteryPct,
		}
	}

	return detail, nil
}
