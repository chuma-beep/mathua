package engine

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chuma-beep/mathua/internal/planning"
)

// serverSessionStudyPlan is the server_sessions kind for saved study plans
// (destination + rate + deadline snapshot for ahead/behind tracking).
const serverSessionStudyPlan = "study_plan"

// planSnapshotTTL bounds a saved plan: plans older than a term are stale.
const planSnapshotTTL = 90 * 24 * time.Hour

// StudyPlanSnapshot is the persisted plan baseline. The frontier is never
// stored — it derives live from knowledge state on every estimate.
type StudyPlanSnapshot struct {
	Destination  string  `json:"destination"`
	DailyGoal    int     `json:"daily_goal"`
	DeadlineDays int     `json:"deadline_days"`
	RestDays     int     `json:"rest_days"`
	XPRemaining  float64 `json:"xp_remaining"`
	Created      string  `json:"created"`
}

// DestinationStatus is one destination with the student's progress.
type DestinationStatus struct {
	*planning.Destination
	Total    int     `json:"total"`
	Mastered int     `json:"mastered"`
	Pct      float64 `json:"pct"`
}

// Destinations returns all destination bundles with per-student progress.
func (e *Engine) Destinations(studentID string) ([]DestinationStatus, error) {
	if e.planner == nil {
		return nil, nil
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, err
	}
	var out []DestinationStatus
	for _, d := range e.planner.Destinations() {
		path, err := e.planner.PathForDestination(d.ID)
		if err != nil {
			continue
		}
		total := len(path.Concepts)
		mastered := len(path.Concepts) - len(e.planner.Unmastered(path, progress))
		pct := 0.0
		if total > 0 {
			pct = float64(mastered) / float64(total)
		}
		out = append(out, DestinationStatus{Destination: d, Total: total, Mastered: mastered, Pct: pct})
	}
	return out, nil
}

// AttemptStatsFor builds per-concept decisive history from stored attempts.
// Diagnostic probes are excluded (assessment, not learning). An incorrect
// attempt with an empty answer is an admitted don't-know; anything else
// incorrect is a failed attempt.
func (e *Engine) AttemptStatsFor(studentID string) (map[string]planning.AttemptStats, error) {
	entries, err := e.repo.GetAttemptsForStudent(studentID)
	if err != nil {
		return nil, err
	}
	out := map[string]planning.AttemptStats{}
	for _, a := range entries {
		if a.Source == "diagnostic" {
			continue
		}
		s := out[a.ConceptID]
		switch {
		case a.Correct:
			s.Correct++
		case a.Answer == "":
			s.DontKnow++
		default:
			s.Failed++
		}
		out[a.ConceptID] = s
	}
	return out, nil
}

// PaceEstimate is the measured daily XP rate: the stored goal scaled by
// trailing adherence (active days / 14). No history means goal-based math,
// labeled so the UI can say "no history yet".
type PaceEstimate struct {
	Rate        float64 `json:"rate"`
	Source      string  `json:"source"` // "measured" | "goal"
	Adherence   float64 `json:"adherence"`
	ActiveDays  int     `json:"active_days"`
	TrailingDay int     `json:"trailing_days"`
}

// PaceFor measures the learner's expected earned XP per active day.
func (e *Engine) PaceFor(studentID string, dailyGoal int) (*PaceEstimate, error) {
	const window = 14
	activity, err := e.repo.GetDailyActivity(studentID, window)
	if err != nil || len(activity) == 0 {
		return &PaceEstimate{Rate: float64(dailyGoal), Source: "goal", TrailingDay: window}, nil
	}
	active := 0
	for _, d := range activity {
		if d.Questions > 0 {
			active++
		}
	}
	if active == 0 {
		return &PaceEstimate{Rate: float64(dailyGoal), Source: "goal", TrailingDay: window}, nil
	}
	adherence := float64(active) / float64(window)
	return &PaceEstimate{
		Rate:        float64(dailyGoal) * adherence,
		Source:      "measured",
		Adherence:   adherence,
		ActiveDays:  active,
		TrailingDay: window,
	}, nil
}

// DestinationEstimateResponse is the full planner answer for one destination.
type DestinationEstimateResponse struct {
	Destination string                     `json:"destination"`
	Estimate    *planning.WorkloadEstimate `json:"estimate"`
	Pace        *PaceEstimate              `json:"pace"`
	Probes      []string                   `json:"probes"`
	PlanDelta   *float64                   `json:"plan_delta_days,omitempty"`
}

// DestinationEstimate prices a destination's remaining workload on both
// tracks (XP + time), both directions (effort→date, deadline→rate).
func (e *Engine) DestinationEstimate(studentID, destID string, dailyGoal, deadlineDays, restDays int, diagnosticMin float64) (*DestinationEstimateResponse, error) {
	if e.planner == nil {
		return nil, fmt.Errorf("planner unavailable")
	}
	path, err := e.planner.PathForDestination(destID)
	if err != nil {
		return nil, err
	}
	if dailyGoal <= 0 {
		dailyGoal = 30
	}
	if restDays < 0 {
		restDays = 0
	}
	if restDays > 6 {
		restDays = 6
	}
	progress, err := e.repo.GetAllProgress(studentID)
	if err != nil {
		return nil, err
	}
	stats, err := e.AttemptStatsFor(studentID)
	if err != nil {
		return nil, err
	}
	pace, err := e.PaceFor(studentID, dailyGoal)
	if err != nil {
		return nil, err
	}
	est := planning.EstimateWorkload(path, progress, stats, planning.EstimateOpts{
		DailyXPRate:     pace.Rate,
		RestDaysPerWeek: restDays,
		DiagnosticMin:   diagnosticMin,
		DeadlineDays:    deadlineDays,
		FeasibleCap:     3 * pace.Rate,
	})
	resp := &DestinationEstimateResponse{
		Destination: destID,
		Estimate:    est,
		Pace:        pace,
		Probes:      planning.SampleClosure(path, 18),
	}
	if snap, ok := e.GetPlan(studentID); ok && snap.Destination == destID && pace.Rate > 0 {
		delta := (snap.XPRemaining - est.XPRemaining) / pace.Rate
		resp.PlanDelta = &delta
	}
	return resp, nil
}

// SavePlan persists the current estimate as the ahead/behind baseline.
func (e *Engine) SavePlan(studentID, destID string, dailyGoal, deadlineDays, restDays int) (*StudyPlanSnapshot, error) {
	resp, err := e.DestinationEstimate(studentID, destID, dailyGoal, deadlineDays, restDays, 0)
	if err != nil {
		return nil, err
	}
	snap := &StudyPlanSnapshot{
		Destination:  destID,
		DailyGoal:    dailyGoal,
		DeadlineDays: deadlineDays,
		RestDays:     restDays,
		XPRemaining:  resp.Estimate.XPRemaining,
		Created:      time.Now().UTC().Format(time.RFC3339),
	}
	blob, err := json.Marshal(snap)
	if err != nil {
		return nil, err
	}
	exp := time.Now().UTC().Add(planSnapshotTTL).Format(time.RFC3339)
	if err := e.repo.UpsertServerSession(serverSessionStudyPlan, studentID, string(blob), exp); err != nil {
		return nil, err
	}
	return snap, nil
}

// GetPlan returns the saved plan baseline, if any and unexpired.
func (e *Engine) GetPlan(studentID string) (*StudyPlanSnapshot, bool) {
	if e.repo == nil {
		return nil, false
	}
	v, exp, found, err := e.repo.GetServerSession(serverSessionStudyPlan, studentID)
	if err != nil || !found || v == "" {
		return nil, false
	}
	if exp != "" && !isFutureRFC3339(exp) {
		return nil, false
	}
	var snap StudyPlanSnapshot
	if err := json.Unmarshal([]byte(v), &snap); err != nil {
		return nil, false
	}
	return &snap, true
}
