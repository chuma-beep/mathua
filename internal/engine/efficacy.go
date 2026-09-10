package engine

import (
	"sort"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// EfficacyReport is the instrumentation summary (improve.md:82).
type EfficacyReport struct {
	ConceptsTouched       int     `json:"concepts_touched"`
	FirstPassRate         float64 `json:"first_pass_rate"`  // correct on attempt 1
	SecondPassRate        float64 `json:"second_pass_rate"` // correct within first 2 attempts
	AvgAttemptsPerConcept float64 `json:"avg_attempts_per_concept"`
	TotalAttempts         int     `json:"total_attempts"`
	StudentsTracked       int     `json:"students_tracked,omitempty"`
}

// Efficacy computes first-pass / second-pass rates from the attempt log.
func (e *Engine) Efficacy(studentID string) (*EfficacyReport, error) {
	attempts, err := e.repo.GetAttemptsForStudent(studentID)
	if err != nil {
		return nil, err
	}
	return computeEfficacy(attempts), nil
}

// AggregateEfficacy computes product-wide first-pass / second-pass rates
// across all students (the MA parity metric: 93% / 98%).
func (e *Engine) AggregateEfficacy() (*EfficacyReport, error) {
	attempts, err := e.repo.GetAllAttempts()
	if err != nil {
		return nil, err
	}
	rep := computeEfficacy(attempts)
	students := make(map[string]bool)
	for _, a := range attempts {
		students[a.StudentID] = true
	}
	rep.StudentsTracked = len(students)
	return rep, nil
}

// computeEfficacy groups attempts by (student, concept) in timestamp order
// and derives the pass-rate metrics.
func computeEfficacy(attempts []storage.AttemptEntry) *EfficacyReport {
	rep := &EfficacyReport{TotalAttempts: len(attempts)}
	if len(attempts) == 0 {
		return rep
	}
	type seq struct{ correct []bool }
	byKey := make(map[string]*seq)
	var order []string
	for _, a := range attempts {
		key := a.StudentID + "|" + a.ConceptID
		s, ok := byKey[key]
		if !ok {
			s = &seq{}
			byKey[key] = s
			order = append(order, key)
		}
		s.correct = append(s.correct, a.Correct)
	}
	firstPass, secondPass := 0, 0
	totalAttempts := 0
	for _, key := range order {
		s := byKey[key]
		rep.ConceptsTouched++
		if len(s.correct) > 0 && s.correct[0] {
			firstPass++
		}
		if (len(s.correct) >= 1 && s.correct[0]) || (len(s.correct) >= 2 && s.correct[1]) {
			secondPass++
		}
		totalAttempts += len(s.correct)
	}
	rep.AvgAttemptsPerConcept = float64(totalAttempts) / float64(rep.ConceptsTouched)
	if rep.ConceptsTouched > 0 {
		rep.FirstPassRate = float64(firstPass) / float64(rep.ConceptsTouched)
		rep.SecondPassRate = float64(secondPass) / float64(rep.ConceptsTouched)
	}
	return rep
}

// EfficacyWeek is one Monday-to-Sunday bucket of the longitudinal trend.
type EfficacyWeek struct {
	WeekStart       string  `json:"week_start"` // Monday, YYYY-MM-DD (UTC)
	Attempts        int     `json:"attempts"`
	ConceptsTouched int     `json:"concepts_touched"`
	FirstPassRate   float64 `json:"first_pass_rate"`
	SecondPassRate  float64 `json:"second_pass_rate"`
	ActiveStudents  int     `json:"active_students"`
}

// EfficacyTrend is the longitudinal view: per-week efficacy plus cohort
// retention. FirstPassTrend is the last week minus the first week (positive
// means first-pass accuracy improved over the observed window).
type EfficacyTrend struct {
	Weeks             []EfficacyWeek `json:"weeks"`
	TotalStudents     int            `json:"total_students"`
	ReturningStudents int            `json:"returning_students"`
	RetentionRate     float64        `json:"retention_rate"`
	FirstPassTrend    float64        `json:"first_pass_trend"`
}

// EfficacyTrend buckets the full attempt log by week (UTC, Monday start) and
// computes cohort retention. It is the dataset behind docs/efficacy.md.
func (e *Engine) EfficacyTrend() (*EfficacyTrend, error) {
	attempts, err := e.repo.GetAllAttempts()
	if err != nil {
		return nil, err
	}
	return computeEfficacyTrend(attempts), nil
}

func computeEfficacyTrend(attempts []storage.AttemptEntry) *EfficacyTrend {
	trend := &EfficacyTrend{Weeks: []EfficacyWeek{}}
	if len(attempts) == 0 {
		return trend
	}

	byWeek := make(map[string][]storage.AttemptEntry)
	weeksSeen := make(map[string]map[string]bool) // week -> set of students
	studentWeeks := make(map[string]map[string]bool)
	for _, a := range attempts {
		wk := weekStartUTC(a.Timestamp).Format("2006-01-02")
		byWeek[wk] = append(byWeek[wk], a)
		if weeksSeen[wk] == nil {
			weeksSeen[wk] = make(map[string]bool)
		}
		weeksSeen[wk][a.StudentID] = true
		if studentWeeks[a.StudentID] == nil {
			studentWeeks[a.StudentID] = make(map[string]bool)
		}
		studentWeeks[a.StudentID][wk] = true
	}

	weekKeys := make([]string, 0, len(byWeek))
	for wk := range byWeek {
		weekKeys = append(weekKeys, wk)
	}
	sort.Strings(weekKeys)

	for _, wk := range weekKeys {
		rep := computeEfficacy(byWeek[wk])
		trend.Weeks = append(trend.Weeks, EfficacyWeek{
			WeekStart:       wk,
			Attempts:        rep.TotalAttempts,
			ConceptsTouched: rep.ConceptsTouched,
			FirstPassRate:   rep.FirstPassRate,
			SecondPassRate:  rep.SecondPassRate,
			ActiveStudents:  len(weeksSeen[wk]),
		})
	}

	trend.TotalStudents = len(studentWeeks)
	for _, wks := range studentWeeks {
		if len(wks) >= 2 {
			trend.ReturningStudents++
		}
	}
	if trend.TotalStudents > 0 {
		trend.RetentionRate = float64(trend.ReturningStudents) / float64(trend.TotalStudents)
	}
	if len(trend.Weeks) >= 2 {
		trend.FirstPassTrend = trend.Weeks[len(trend.Weeks)-1].FirstPassRate - trend.Weeks[0].FirstPassRate
	}
	return trend
}

// weekStartUTC returns the Monday 00:00 UTC of t's week.
func weekStartUTC(t time.Time) time.Time {
	t = t.UTC()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday -> 7 so Monday is day 1
	}
	offset := weekday - int(time.Monday)
	start := t.AddDate(0, 0, -offset)
	return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
}
