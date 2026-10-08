package server

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Report endpoints: users (authed or guest) complain about questions,
// explanations, lesson bodies, worked examples, or diagrams. Listing and resolving are
// administration now and live under /api/admin/reports behind the role boundary — see adminv2.go.
// This file only creates reports.

var (
	validReportKinds = map[string]bool{
		"question": true, "explanation": true, "lesson_body": true,
		"worked_example": true, "diagram": true,
	}
	validReportReasons = map[string]bool{
		"wrong_answer": true, "bad_explanation": true, "unclear": true,
		"formatting": true, "other": true,
	}
	validReportSources = map[string]bool{
		"study": true, "diagnostic": true, "quiz": true, "lesson": true,
		"concept": true, "review": true,
	}
)

// Per-reporter spam guard: 10 reports/hour (writeLimiter covers per-IP burst).
var reportLimiter = struct {
	sync.Mutex
	hits map[string][]time.Time
}{hits: make(map[string][]time.Time)}

func reportAllowed(reporter string) bool {
	if reporter == "" {
		reporter = "anon"
	}
	now := time.Now()
	cutoff := now.Add(-time.Hour)
	reportLimiter.Lock()
	defer reportLimiter.Unlock()
	kept := reportLimiter.hits[reporter][:0]
	for _, t := range reportLimiter.hits[reporter] {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= 10 {
		reportLimiter.hits[reporter] = kept
		return false
	}
	reportLimiter.hits[reporter] = append(kept, now)
	// Opportunistic sweep: reporter keys never seen again would otherwise
	// accumulate forever (no janitor for this map). Bound the scan.
	if len(reportLimiter.hits) > 5000 {
		for k, times := range reportLimiter.hits {
			still := times[:0]
			for _, t := range times {
				if t.After(cutoff) {
					still = append(still, t)
				}
			}
			if len(still) == 0 {
				delete(reportLimiter.hits, k)
			} else {
				reportLimiter.hits[k] = still
			}
		}
	}
	return true
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// POST /api/reports — open to guests. Body mirrors submitStudyAnswer:
// authed student_id wins, else body reporter_id (guest id) is stored as-is.
func (s *Server) handleCreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		ConceptID   string `json:"concept_id"`
		Kind        string `json:"kind"`
		Question    string `json:"question"`
		Expected    string `json:"expected"`
		Explanation string `json:"explanation"`
		LessonID    string `json:"lesson_id"`
		Source      string `json:"source"`
		SessionID   string `json:"session_id"`
		AttemptID   string `json:"attempt_id"`
		Reason      string `json:"reason"`
		Detail      string `json:"detail"`
		ReporterID  string `json:"reporter_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	kind := strings.TrimSpace(req.Kind)
	if kind == "" {
		kind = "question"
	}
	if !validReportKinds[kind] {
		writeError(w, "invalid kind", 400)
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if !validReportReasons[reason] {
		writeError(w, "invalid reason", 400)
		return
	}
	conceptID := strings.TrimSpace(req.ConceptID)
	lessonID := truncate(req.LessonID, 256)
	question := truncate(req.Question, 4000)
	source := strings.TrimSpace(req.Source)
	if source == "" {
		source = "study"
	}
	if !validReportSources[source] {
		writeError(w, "invalid source", 400)
		return
	}
	if conceptID == "" && lessonID == "" && question == "" {
		writeError(w, "concept_id, lesson_id, or question required", 400)
		return
	}
	reporter := ""
	if authID, _ := r.Context().Value(authStudentKey{}).(string); authID != "" {
		reporter = authID
	} else {
		reporter = strings.TrimSpace(req.ReporterID)
	}
	if !reportAllowed(reporter) {
		writeError(w, "too many reports, try again later", 429)
		return
	}
	id, err := s.repo.CreateReport(storage.QuestionReport{
		ReporterID:  truncate(reporter, 128),
		ConceptID:   truncate(conceptID, 128),
		Kind:        kind,
		Question:    question,
		Expected:    truncate(req.Expected, 2000),
		Explanation: truncate(req.Explanation, 4000),
		LessonID:    lessonID,
		Source:      truncate(source, 64),
		SessionID:   truncate(req.SessionID, 128),
		AttemptID:   truncate(req.AttemptID, 128),
		Reason:      reason,
		Detail:      truncate(req.Detail, 2000),
		Status:      "open",
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		writeError(w, "failed to save report", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"id": id, "status": "open"})
}
