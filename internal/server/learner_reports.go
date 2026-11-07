package server

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/scheduler"
	"github.com/chuma-beep/mathua/internal/storage"
)

// learnerReportDefaultHistory is the first page returned with the report. The
// history endpoint serves the rest, so the page never carries an unbounded
// stream of attempts.
const learnerReportDefaultHistory = 25

// learnerReportMaxHistory caps a single history page.
const learnerReportMaxHistory = 100

type learnerReportPeriod struct {
	Days  int       `json:"days"`
	Since time.Time `json:"since,omitempty"`
	Until time.Time `json:"until"`
	Label string    `json:"label"`
}

type learnerReportOverview struct {
	Total        int `json:"total"`
	Mastered     int `json:"mastered"`
	InProgress   int `json:"in_progress"`
	Unlocked     int `json:"unlocked"`
	Locked       int `json:"locked"`
	XPLifetime   int `json:"xp_lifetime"`
	XPPeriod     int `json:"xp_period"`
	ActivityDays int `json:"activity_days"`
	Questions    int `json:"questions"`
	Correct      int `json:"correct"`
}

type learnerReportGap struct {
	ConceptID string `json:"concept_id"`
	Label     string `json:"label"`
	State     string `json:"state,omitempty"`
}

type learnerMasteryRow struct {
	ConceptID       string             `json:"concept_id"`
	Label           string             `json:"label"`
	Domain          string             `json:"domain"`
	State           string             `json:"state"`
	Status          string             `json:"status"`
	Decaying        bool               `json:"decaying"`
	MasteryPct      float64            `json:"mastery_pct"`
	MasteredAt      *time.Time         `json:"mastered_at,omitempty"`
	Attempts        int                `json:"attempts"`
	Streak          int                `json:"streak"`
	LastAttempted   *time.Time         `json:"last_attempted,omitempty"`
	NextReview      *time.Time         `json:"next_review,omitempty"`
	PlacementSeeded bool               `json:"placement_seeded"`
	Provisional     bool               `json:"provisional"`
	Gaps            []learnerReportGap `json:"prerequisite_gaps"`
	Href            string             `json:"href"`
}

type learnerAssessment struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"`
	Title       string             `json:"title"`
	Status      string             `json:"status"`
	StartedAt   *time.Time         `json:"started_at,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	Questions   int                `json:"questions"`
	Correct     int                `json:"correct"`
	Accuracy    *float64           `json:"accuracy,omitempty"`
	Concepts    []string           `json:"concepts"`
	Findings    []learnerReportGap `json:"findings"`
	Href        string             `json:"href,omitempty"`
	Note        string             `json:"note,omitempty"`
}

type learnerHistoryEvent struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	At        time.Time `json:"at"`
	Title     string    `json:"title"`
	ConceptID string    `json:"concept_id,omitempty"`
	Source    string    `json:"source,omitempty"`
	Detail    string    `json:"detail,omitempty"`
	Href      string    `json:"href,omitempty"`
	Correct   *bool     `json:"correct,omitempty"`
	XP        *int      `json:"xp,omitempty"`
}

type learnerReportResponse struct {
	GeneratedAt  time.Time               `json:"generated_at"`
	Learner      map[string]string       `json:"learner"`
	Period       learnerReportPeriod     `json:"period"`
	Overview     learnerReportOverview   `json:"overview"`
	Domains      []engine.DomainSummary  `json:"domains"`
	Mastery      []learnerMasteryRow     `json:"mastery"`
	Assessments  []learnerAssessment     `json:"assessments"`
	Activity     []storage.DailyActivity `json:"activity"`
	History      []learnerHistoryEvent   `json:"history"`
	HistoryTotal int                     `json:"history_total"`
	Limitations  []string                `json:"limitations"`
}

type learnerHistoryResponse struct {
	Events []learnerHistoryEvent `json:"events"`
	Total  int                   `json:"total"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}

// handleLearnerReport builds the learner-facing report from stored state. It
// deliberately does not grade anything: mastery comes from the stored progress
// rows and the eligibility rule, XP from the ledger, assessments from the
// recorded answers and completion rows. Missing history is reported as missing.
func (s *Server) handleLearnerReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sid, _ := r.Context().Value(authStudentKey{}).(string)
	if sid == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	period, err := learnerReportPeriodFrom(r)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	domain := strings.TrimSpace(r.URL.Query().Get("domain"))

	student, err := s.repo.GetStudent(sid)
	if err != nil || student == nil {
		writeError(w, "learner not found", http.StatusNotFound)
		return
	}
	scores, err := s.eng.GetScores(sid)
	if err != nil {
		writeError(w, "could not read scores", http.StatusInternalServerError)
		return
	}
	progress, err := s.repo.GetAllProgress(sid)
	if err != nil {
		writeError(w, "could not read progress", http.StatusInternalServerError)
		return
	}
	snaps, err := s.eng.LearnerSnapshots(sid)
	if err != nil {
		writeError(w, "could not read learning state", http.StatusInternalServerError)
		return
	}

	dag := s.eng.GetDAG()
	now := time.Now()

	domains, err := s.eng.CurriculumDomains(sid)
	if err != nil {
		writeError(w, "could not read domains", http.StatusInternalServerError)
		return
	}
	if domain != "" {
		filtered := domains[:0]
		for _, d := range domains {
			if d.ID == domain {
				filtered = append(filtered, d)
			}
		}
		domains = filtered
	}

	overview := learnerReportOverview{XPLifetime: scores.XPTotal}
	for _, d := range domains {
		overview.Total += d.ConceptCount
		overview.Mastered += d.MasteredCount
		overview.InProgress += d.InProgressCount
		overview.Unlocked += d.UnlockedCount
	}
	overview.Locked = overview.Total - overview.Mastered - overview.InProgress - overview.Unlocked
	if overview.Locked < 0 {
		overview.Locked = 0
	}

	xpPeriod, err := s.repo.XPInRange(sid, period.Since, period.Until)
	if err != nil {
		writeError(w, "could not read XP for the period", http.StatusInternalServerError)
		return
	}
	overview.XPPeriod = xpPeriod

	days := period.Days
	if days <= 0 {
		// GetDailyActivity takes a day count, not a range. Ten years is a bounded
		// stand-in for "all available" and is labelled as such below.
		days = 3650
	}
	activity, err := s.repo.GetDailyActivity(sid, days)
	if err != nil {
		writeError(w, "could not read activity", http.StatusInternalServerError)
		return
	}
	overview.ActivityDays = len(activity)
	for _, day := range activity {
		overview.Questions += day.Questions
		overview.Correct += day.Correct
	}

	masteryRows := learnerMasteryRows(dag, progress, snaps, now, domain)

	assessments, err := s.learnerAssessments(sid, period)
	if err != nil {
		writeError(w, "could not read assessments", http.StatusInternalServerError)
		return
	}

	history, historyTotal, err := s.learnerHistory(sid, dag, learnerReportDefaultHistory, 0)
	if err != nil {
		writeError(w, "could not read history", http.StatusInternalServerError)
		return
	}

	learner := map[string]string{"id": student.ID, "name": student.Name}
	if learner["name"] == "" {
		learner["name"] = student.Username
	}

	writeJSON(w, learnerReportResponse{
		GeneratedAt:  now.UTC(),
		Learner:      learner,
		Period:       period,
		Overview:     overview,
		Domains:      domains,
		Mastery:      masteryRows,
		Assessments:  assessments,
		Activity:     activity,
		History:      history,
		HistoryTotal: historyTotal,
		Limitations:  learnerReportLimitations(domain),
	})
}

// handleLearnerReportHistory serves further pages of the same event stream.
func (s *Server) handleLearnerReportHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sid, _ := r.Context().Value(authStudentKey{}).(string)
	if sid == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	limit := learnerReportDefaultHistory
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > learnerReportMaxHistory {
			writeError(w, fmt.Sprintf("limit must be between 1 and %d", learnerReportMaxHistory), http.StatusBadRequest)
			return
		}
		limit = n
	}
	offset := 0
	if raw := r.URL.Query().Get("offset"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			writeError(w, "offset must be zero or greater", http.StatusBadRequest)
			return
		}
		offset = n
	}
	events, total, err := s.learnerHistory(sid, s.eng.GetDAG(), limit, offset)
	if err != nil {
		writeError(w, "could not read history", http.StatusInternalServerError)
		return
	}
	writeJSON(w, learnerHistoryResponse{Events: events, Total: total, Limit: limit, Offset: offset})
}

func learnerReportPeriodFrom(r *http.Request) (learnerReportPeriod, error) {
	days := 30
	if raw := r.URL.Query().Get("days"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return learnerReportPeriod{}, fmt.Errorf("days must be one of 7, 30, 90 or 0")
		}
		days = n
	}
	switch days {
	case 7, 30, 90, 0:
	default:
		return learnerReportPeriod{}, fmt.Errorf("days must be one of 7, 30, 90 or 0")
	}
	now := time.Now().UTC()
	p := learnerReportPeriod{Days: days, Until: now}
	if days == 0 {
		p.Since = time.Unix(0, 0).UTC()
		p.Label = "All recorded history"
	} else {
		p.Since = now.AddDate(0, 0, -days)
		p.Label = fmt.Sprintf("Last %d days", days)
	}
	return p, nil
}

// learnerMasteryRows joins the stored progress to the eligibility rule. The
// rule is asked, not re-derived: lock, in-progress and mastered stay distinct,
// and a placement seed is visible as placement rather than as knowledge.
func learnerMasteryRows(dag *concepts.DAG, progress map[string]*storage.ConceptProgress, snaps map[string]*scheduler.ConceptSnapshot, now time.Time, domain string) []learnerMasteryRow {
	rows := make([]learnerMasteryRow, 0, len(progress))
	for _, c := range dag.Order() {
		if domain != "" && c.Domain != domain {
			continue
		}
		p := progress[c.ID]
		snap := snaps[c.ID]
		elig := scheduler.EligibilityOf(snaps, now, c)

		status := string(mastery.StatusUnseen)
		row := learnerMasteryRow{
			ConceptID: c.ID,
			Label:     c.Label,
			Domain:    c.Domain,
			State:     elig.State,
			Status:    status,
			Gaps:      learnerReportGaps(elig.Missing, snaps, dag),
			Href:      "/study?concept=" + c.ID,
		}
		if p != nil {
			row.Status = p.Status
			row.Attempts = p.Attempts
			row.Streak = p.Streak
			row.MasteredAt = p.MasteredAt
			row.PlacementSeeded = p.PlacementSeeded
			row.Provisional = p.IsProvisional
			if p.LastAttempted != nil && !p.LastAttempted.IsZero() {
				t := *p.LastAttempted
				row.LastAttempted = &t
			}
			if p.NextReviewDue != nil && !p.NextReviewDue.IsZero() {
				t := *p.NextReviewDue
				row.NextReview = &t
			}
			if snap != nil {
				row.MasteryPct = mastery.MasteryPct(mastery.Status(p.Status), snap.Streak, snap.RequiredStreak)
				if p.Status == string(mastery.StatusMastered) {
					last := masteryReferenceTime(snap, p)
					days := 0.0
					if !last.IsZero() {
						days = now.Sub(last).Hours() / 24
					}
					if mastery.EffectiveStatusFSRS(mastery.Status(p.Status), days, snap.FSRSStability, mastery.DecayDays) == "DECAYING" {
						row.Decaying = true
						row.State = "decaying"
					}
				}
			}
		}
		if row.State == "" {
			row.State = elig.State
		}
		rows = append(rows, row)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Domain != rows[j].Domain {
			return rows[i].Domain < rows[j].Domain
		}
		return rows[i].Label < rows[j].Label
	})
	return rows
}

func masteryReferenceTime(snap *scheduler.ConceptSnapshot, p *storage.ConceptProgress) time.Time {
	if snap != nil && snap.FSRSLastReview != nil {
		return *snap.FSRSLastReview
	}
	if snap != nil && !snap.LastReviewed.IsZero() {
		return snap.LastReviewed
	}
	if p != nil && p.MasteredAt != nil {
		return *p.MasteredAt
	}
	return time.Time{}
}

func learnerReportGaps(missing []string, snaps map[string]*scheduler.ConceptSnapshot, dag *concepts.DAG) []learnerReportGap {
	gaps := make([]learnerReportGap, 0, len(missing))
	for _, id := range missing {
		label := id
		if c := dag.Concept(id); c != nil {
			label = c.Label
		}
		state := ""
		if snap := snaps[id]; snap != nil {
			state = string(snap.Status)
		}
		gaps = append(gaps, learnerReportGap{ConceptID: id, Label: label, State: state})
	}
	return gaps
}

type quizSession struct {
	id        string
	first     time.Time
	last      time.Time
	questions int
	correct   int
	concepts  map[string]bool
	missed    map[string]bool
}

func (s *Server) learnerAssessments(sid string, period learnerReportPeriod) ([]learnerAssessment, error) {
	attempts, err := s.repo.AssessmentAttempts(sid, period.Since, period.Until, 2000)
	if err != nil {
		return nil, err
	}
	dag := s.eng.GetDAG()
	label := func(id string) string {
		if c := dag.Concept(id); c != nil {
			return c.Label
		}
		return id
	}

	quizzes := map[string]*quizSession{}
	var quizOrder []string
	var diagnostics []storage.AttemptEntry
	for _, a := range attempts {
		if a.Source == "diagnostic" {
			diagnostics = append(diagnostics, a)
			continue
		}
		if a.Source != "quiz" {
			continue
		}
		key := a.SessionID
		if key == "" {
			key = fmt.Sprintf("attempt-%d", a.ID)
		}
		q := quizzes[key]
		if q == nil {
			q = &quizSession{id: key, first: a.Timestamp, last: a.Timestamp, concepts: map[string]bool{}, missed: map[string]bool{}}
			quizzes[key] = q
			quizOrder = append(quizOrder, key)
		}
		if a.Timestamp.Before(q.first) {
			q.first = a.Timestamp
		}
		if a.Timestamp.After(q.last) {
			q.last = a.Timestamp
		}
		q.questions++
		if a.Correct {
			q.correct++
		} else {
			q.missed[a.ConceptID] = true
		}
		q.concepts[a.ConceptID] = true
	}
	sort.Slice(quizOrder, func(i, j int) bool {
		return quizzes[quizOrder[i]].last.Before(quizzes[quizOrder[j]].last)
	})

	completions, err := s.repo.QuizCompletions(sid, period.Since, 500)
	if err != nil {
		return nil, err
	}
	type completion struct {
		at  time.Time
		raw storage.QuizCompletion
	}
	parsed := make([]completion, 0, len(completions))
	for _, c := range completions {
		if t, ok := parseReportTime(c.CompletedAt); ok && !t.Before(period.Since) && t.Before(period.Until) {
			parsed = append(parsed, completion{at: t, raw: c})
		}
	}
	sort.Slice(parsed, func(i, j int) bool { return parsed[i].at.Before(parsed[j].at) })

	// Completion rows do not carry a session id. A completion is matched to the
	// latest answered session that ended before it, which is the only ordering
	// the records support. Sessions left over are shown as incomplete.
	matched := map[string]time.Time{}
	used := map[string]bool{}
	for _, c := range parsed {
		for i := len(quizOrder) - 1; i >= 0; i-- {
			key := quizOrder[i]
			if used[key] || quizzes[key].last.After(c.at) {
				continue
			}
			used[key] = true
			matched[key] = c.at
			break
		}
	}

	out := make([]learnerAssessment, 0, len(quizOrder)+len(parsed)+1)
	for _, key := range quizOrder {
		q := quizzes[key]
		row := learnerAssessment{
			ID:        "quiz:" + q.id,
			Type:      "quiz",
			Title:     "Mastery check",
			Questions: q.questions,
			Correct:   q.correct,
			Accuracy:  reportAccuracy(q.correct, q.questions),
			Concepts:  sortedKeys(q.concepts),
			Findings:  reportFindings(q.missed, label),
			Note:      "The record stores answers but not a pass or fail threshold.",
			Href:      "/quiz",
		}
		first := q.first
		row.StartedAt = &first
		if at, ok := matched[key]; ok {
			row.Status = "completed"
			row.CompletedAt = &at
		} else {
			row.Status = "incomplete"
			row.Note = "Answers were recorded, but no completion record follows this session."
		}
		out = append(out, row)
	}
	// A completion with no matching session still deserves a row: it is a real
	// record even when its answers fall outside the selected period.
	for _, c := range parsed {
		found := false
		for _, at := range matched {
			if at.Equal(c.at) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		at := c.at
		out = append(out, learnerAssessment{
			ID:          fmt.Sprintf("quiz-completion:%d", len(out)),
			Type:        "quiz",
			Title:       "Mastery check",
			Status:      "completed",
			CompletedAt: &at,
			Concepts:    []string{},
			Findings:    []learnerReportGap{},
			Note:        "The completion record stores the XP gate snapshot, not a score, so no accuracy is shown.",
			Href:        "/quiz",
		})
	}

	completedAt, hasCompletion, err := s.repo.DiagnosticCompletedAt(sid)
	if err != nil {
		return nil, err
	}
	var diagStart *time.Time
	diagCorrect := 0
	diagConcepts := map[string]bool{}
	diagMissed := map[string]bool{}
	for _, a := range diagnostics {
		if diagStart == nil || a.Timestamp.Before(*diagStart) {
			t := a.Timestamp
			diagStart = &t
		}
		if a.Correct {
			diagCorrect++
		} else {
			diagMissed[a.ConceptID] = true
		}
		diagConcepts[a.ConceptID] = true
	}
	if len(diagnostics) > 0 || (hasCompletion && !completedAt.Before(period.Since) && completedAt.Before(period.Until)) {
		row := learnerAssessment{
			ID:        "diagnostic",
			Type:      "diagnostic",
			Title:     "Placement diagnostic",
			StartedAt: diagStart,
			Questions: len(diagnostics),
			Correct:   diagCorrect,
			Accuracy:  reportAccuracy(diagCorrect, len(diagnostics)),
			Concepts:  sortedKeys(diagConcepts),
			Findings:  reportFindings(diagMissed, label),
			Href:      "/diagnose",
			Note:      "A diagnostic opens concepts. It does not establish mastery.",
		}
		if hasCompletion {
			row.Status = "completed"
			t := completedAt
			row.CompletedAt = &t
		} else {
			row.Status = "incomplete"
		}
		out = append(out, row)
	}

	sort.SliceStable(out, func(i, j int) bool {
		a, b := assessmentTime(out[i]), assessmentTime(out[j])
		return a.After(b)
	})
	return out, nil
}

func assessmentTime(a learnerAssessment) time.Time {
	if a.CompletedAt != nil {
		return *a.CompletedAt
	}
	if a.StartedAt != nil {
		return *a.StartedAt
	}
	return time.Time{}
}

func (s *Server) learnerHistory(sid string, dag *concepts.DAG, limit, offset int) ([]learnerHistoryEvent, int, error) {
	events, total, err := s.repo.LearningEventsPage(sid, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	label := func(id string) string {
		if id == "" {
			return ""
		}
		if c := dag.Concept(id); c != nil {
			return c.Label
		}
		return id
	}
	out := make([]learnerHistoryEvent, 0, len(events))
	for _, e := range events {
		ev := learnerHistoryEvent{
			ID:        e.Kind + ":" + e.Ref,
			Kind:      e.Kind,
			At:        e.At.UTC(),
			ConceptID: e.ConceptID,
			Source:    e.Source,
		}
		switch e.Kind {
		case "answer":
			ev.Title = label(e.ConceptID)
			ev.Detail = "Answer recorded"
			if e.Correct {
				ev.Detail = "Answer correct"
			} else {
				ev.Detail = "Answer incorrect"
			}
			correct := e.Correct
			ev.Correct = &correct
			ev.Href = "/study?concept=" + e.ConceptID
		case "mastery":
			ev.Title = label(e.ConceptID)
			ev.Detail = "Reached mastery"
			ev.Href = "/study?concept=" + e.ConceptID
		case "quiz":
			ev.Title = "Mastery check completed"
			ev.Detail = "Completion recorded"
			ev.Href = "/quiz"
		case "diagnostic":
			ev.Title = "Placement diagnostic completed"
			ev.Detail = "Completion recorded"
			ev.Href = "/diagnose"
		case "xp":
			ev.Title = "XP adjustment"
			ev.Detail = "XP ledger entry"
			if e.Reason != "" {
				ev.Source = e.Reason
			}
			v := e.XP
			ev.XP = &v
		default:
			ev.Title = e.Kind
		}
		out = append(out, ev)
	}
	return out, total, nil
}

func learnerReportLimitations(domain string) []string {
	limitations := []string{
		"Mastery counts describe the current state. The system records when a concept reached mastery, but not every earlier transition, so there is no full audit trail.",
		"Quiz completion records store the XP gate snapshot, not a score. A quiz percentage is therefore not shown.",
		"Activity comes from recorded answers. A day with no record is absent from the series rather than shown as zero.",
		"Assessment sessions are matched to completion records by time because completion rows do not store a session id.",
		"Session duration is not stored. Individual answer times are recorded, but their sum is not the wall-clock length of the assessment.",
	}
	if domain != "" {
		limitations = append(limitations, "The domain filter applies to the overview, domain and mastery views. XP, activity and history are not domain scoped because answers and XP events do not store a domain.")
	}
	return limitations
}

func reportAccuracy(correct, questions int) *float64 {
	if questions <= 0 {
		return nil
	}
	v := float64(correct) / float64(questions)
	return &v
}

func reportFindings(missed map[string]bool, label func(string) string) []learnerReportGap {
	out := make([]learnerReportGap, 0, len(missed))
	for id := range missed {
		out = append(out, learnerReportGap{ConceptID: id, Label: label(id)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func parseReportTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}
