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
	Dependents    []PrereqInfo     `json:"dependents,omitempty"`
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

	for _, dep := range e.dag.DependentsOf(conceptID) {
		status := string(mastery.StatusUnseen)
		masteryPct := 0.0
		if p, ok := progress[dep.ID]; ok {
			status = p.Status
		}
		detail.Dependents = append(detail.Dependents, PrereqInfo{
			ID:         dep.ID,
			Label:      dep.Label,
			Status:     status,
			MasteryPct: masteryPct,
		})
	}

	return detail, nil
}

// LessonPrerequisites computes the aggregate prerequisites for a lesson
// (the set of concepts a lesson covers). It returns unique prerequisites
// that are NOT covered by the lesson itself, with optional progress info.
func (e *Engine) LessonPrerequisites(conceptIDs []string, progress map[string]*storage.ConceptProgress) []PrereqInfo {
	inLesson := make(map[string]bool, len(conceptIDs))
	for _, cid := range conceptIDs {
		inLesson[cid] = true
	}

	prereqSet := make(map[string]bool)
	for _, cid := range conceptIDs {
		c := e.dag.Concept(cid)
		if c == nil {
			continue
		}
		for _, pid := range c.Prerequisites {
			if !inLesson[pid] {
				prereqSet[pid] = true
			}
		}
	}

	prereqs := make([]PrereqInfo, 0, len(prereqSet))
	for pid := range prereqSet {
		c := e.dag.Concept(pid)
		label := pid
		if c != nil {
			label = c.Label
		}

		status := string(mastery.StatusUnseen)
		masteryPct := 0.0
		if p, ok := progress[pid]; ok {
			status = p.Status
			if c != nil && c.MasteryThreshold.Streak > 0 {
				masteryPct = float64(p.Streak) / float64(c.MasteryThreshold.Streak)
				if masteryPct > 1 {
					masteryPct = 1
				}
			}
		}

		prereqs = append(prereqs, PrereqInfo{
			ID:         pid,
			Label:      label,
			Status:     status,
			MasteryPct: masteryPct,
		})
	}
	return prereqs
}

// LessonDependents computes the concepts that depend on any of the given
// concept IDs. Useful for showing "what to study next" after a lesson.
func (e *Engine) LessonDependents(conceptIDs []string, progress map[string]*storage.ConceptProgress) []PrereqInfo {
	inSet := make(map[string]bool, len(conceptIDs))
	for _, cid := range conceptIDs {
		inSet[cid] = true
	}

	depSet := make(map[string]bool)
	for _, cid := range conceptIDs {
		for _, dep := range e.dag.DependentsOf(cid) {
			if !inSet[dep.ID] {
				depSet[dep.ID] = true
			}
		}
	}

	deps := make([]PrereqInfo, 0, len(depSet))
	for did := range depSet {
		c := e.dag.Concept(did)
		label := did
		if c != nil {
			label = c.Label
		}

		status := string(mastery.StatusUnseen)
		masteryPct := 0.0
		if p, ok := progress[did]; ok {
			status = p.Status
		}

		deps = append(deps, PrereqInfo{
			ID:         did,
			Label:      label,
			Status:     status,
			MasteryPct: masteryPct,
		})
	}
	return deps
}

// LessonDependentsByDomain groups dependents by their concept domain.
func (e *Engine) LessonDependentsByDomain(conceptIDs []string, progress map[string]*storage.ConceptProgress) map[string][]PrereqInfo {
	deps := e.LessonDependents(conceptIDs, progress)
	byDomain := make(map[string][]PrereqInfo)
	for _, d := range deps {
		c := e.dag.Concept(d.ID)
		domain := "general"
		if c != nil {
			domain = c.Domain
		}
		byDomain[domain] = append(byDomain[domain], d)
	}
	return byDomain
}
