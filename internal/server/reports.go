package server

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Report endpoints: users (authed or guest) complain about questions,
// explanations, lesson bodies, worked examples, or diagrams.
// Triage (list/resolve) is gated by ADMIN_TOKEN.

var (
	validReportKinds = map[string]bool{
		"question": true, "explanation": true, "lesson_body": true,
		"worked_example": true, "diagram": true,
	}
	validReportReasons = map[string]bool{
		"wrong_answer": true, "bad_explanation": true, "unclear": true,
		"formatting": true, "other": true,
	}
	validReportStatuses = map[string]bool{
		"open": true, "confirmed": true, "fixed": true, "dismissed": true,
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

func (s *Server) adminAuthorized(r *http.Request) bool {
	bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if bearer != "" {
		// 24h password-login session …
		if s.validAdminSession(bearer) {
			return true
		}
		// … or the legacy static ADMIN_TOKEN.
		if token := strings.TrimSpace(os.Getenv("ADMIN_TOKEN")); token != "" && bearer == token {
			return true
		}
	}
	// Legacy token via query param (local-dev convenience; session tokens
	// are header-only so they never leak into logs via URLs).
	if token := strings.TrimSpace(os.Getenv("ADMIN_TOKEN")); token != "" && r.URL.Query().Get("token") == token {
		return true
	}
	return false
}

// GET /api/reports?status=&limit=&offset= — admin only.
func (s *Server) handleListReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if !s.adminAuthorized(r) {
		writeError(w, "not found", 404)
		return
	}
	q := r.URL.Query()
	status := strings.TrimSpace(q.Get("status"))
	if status != "" && status != "all" && !validReportStatuses[status] {
		writeError(w, "invalid status", 400)
		return
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	reports, err := s.repo.ListReports(status, limit, offset)
	if err != nil {
		writeError(w, "failed to list reports", 500)
		return
	}
	if reports == nil {
		reports = []storage.QuestionReport{}
	}
	writeJSON(w, map[string]interface{}{"reports": reports})
}

// PATCH /api/reports/{id} {status} — admin only.
func (s *Server) handleUpdateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if !s.adminAuthorized(r) {
		writeError(w, "not found", 404)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/reports/")
	idStr = strings.Trim(idStr, "/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, "invalid report id", 400)
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	status := strings.TrimSpace(req.Status)
	if !validReportStatuses[status] {
		writeError(w, "invalid status", 400)
		return
	}
	if err := s.repo.UpdateReportStatus(id, status); err != nil {
		writeError(w, "report not found", 404)
		return
	}
	writeJSON(w, map[string]interface{}{"id": id, "status": status})
}
