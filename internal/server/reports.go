package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Report endpoints: users (authed or guest) complain about questions,
// explanations, lesson bodies, worked examples, or diagrams. Listing and resolving are
// administration now and live under /api/admin/reports behind the role boundary, see adminv2.go.
// This file only creates reports.

var (
	validReportKinds = map[string]bool{
		"question": true, "explanation": true, "lesson_body": true,
		"worked_example": true, "diagram": true,
		// A message is not a complaint about one block. It is the contact box's
		// only destination, and it rides this table so that it lands in the same
		// audited queue instead of an inbox nobody remembers to open.
		"message": true,
	}
	validReportReasons = map[string]bool{
		"wrong_answer": true, "bad_explanation": true, "unclear": true,
		"formatting": true, "other": true,
	}
	validReportSources = map[string]bool{
		"study": true, "practice": true, "diagnostic": true, "quiz": true,
		"lesson": true, "concept": true, "review": true, "contact": true,
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

const maxReportImageBytes = 4 << 20

// reportCreateReq is both the JSON body and the multipart form fields. A
// screenshot uses the same request so a report and its image cannot arrive as
// two independently retryable operations.
type reportCreateReq struct {
	ConceptID   string `json:"concept_id"`
	Kind        string `json:"kind"`
	Question    string `json:"question"`
	Expected    string `json:"expected"`
	Explanation string `json:"explanation"`
	LessonID    string `json:"lesson_id"`
	Source      string `json:"source"`
	SessionID   string `json:"session_id"`
	AttemptID   string `json:"attempt_id"`
	// What the learner typed and the format the field offered. The client
	// already holds all four on the served question; sending them makes the
	// report reproducible instead of a screenshot the moderator cannot see.
	AnswerGiven     string `json:"answer_given"`
	GradingType     string `json:"grading_type"`
	AnswerHint      string `json:"answer_hint"`
	AnswerInputMode string `json:"answer_input_mode"`
	Reason          string `json:"reason"`
	Detail          string `json:"detail"`
	ReporterID      string `json:"reporter_id"`
}

// decodeReportRequest accepts either the original JSON body or a multipart form
// carrying an optional "image". The image is validated before the report row is
// created, so a bad file never leaves a report that claims a screenshot it does
// not have.
func decodeReportRequest(w http.ResponseWriter, r *http.Request, req *reportCreateReq) (*storage.ReportImage, bool) {
	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		if err := decodeJSON(w, r, req); err != nil {
			writeError(w, "invalid request", 400)
			return nil, false
		}
		return nil, true
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxReportImageBytes+(1<<20))
	if err := r.ParseMultipartForm(maxReportImageBytes + (1 << 20)); err != nil {
		writeError(w, "screenshot is too large (4MB max)", 400)
		return nil, false
	}
	req.ConceptID = r.FormValue("concept_id")
	req.Kind = r.FormValue("kind")
	req.Question = r.FormValue("question")
	req.Expected = r.FormValue("expected")
	req.Explanation = r.FormValue("explanation")
	req.LessonID = r.FormValue("lesson_id")
	req.Source = r.FormValue("source")
	req.SessionID = r.FormValue("session_id")
	req.AttemptID = r.FormValue("attempt_id")
	req.AnswerGiven = r.FormValue("answer_given")
	req.GradingType = r.FormValue("grading_type")
	req.AnswerHint = r.FormValue("answer_hint")
	req.AnswerInputMode = r.FormValue("answer_input_mode")
	req.Reason = r.FormValue("reason")
	req.Detail = r.FormValue("detail")
	req.ReporterID = r.FormValue("reporter_id")

	f, _, err := r.FormFile("image")
	if errors.Is(err, http.ErrMissingFile) {
		return nil, true
	}
	if err != nil {
		writeError(w, "could not read screenshot", 400)
		return nil, false
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxReportImageBytes+1))
	if err != nil {
		writeError(w, "could not read screenshot", 400)
		return nil, false
	}
	if len(data) == 0 {
		return nil, true
	}
	if len(data) > maxReportImageBytes {
		writeError(w, "screenshot is too large (4MB max)", 400)
		return nil, false
	}
	ct := http.DetectContentType(data)
	switch ct {
	case "image/png", "image/jpeg", "image/webp", "image/gif":
	default:
		writeError(w, "screenshot must be PNG, JPEG, GIF, or WebP", 400)
		return nil, false
	}
	return &storage.ReportImage{ContentType: ct, Bytes: data}, true
}

// POST /api/reports, open to guests. Body mirrors submitStudyAnswer:
// authed student_id wins, else body reporter_id (guest id) is stored as-is.
func (s *Server) handleCreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req reportCreateReq
	image, ok := decodeReportRequest(w, r, &req)
	if !ok {
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
	// A message has no failure mode to classify; "other" is what the contact box
	// sends and what a hand-written client gets if it sends nothing.
	if kind == "message" && reason == "" {
		reason = "other"
	}
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
	// A message names no concept, lesson or question, and putting a subject line
	// into `question` would make the moderation card lie about what was asked.
	// The text lives in detail, so a message is held to a length instead.
	if kind == "message" {
		if len(strings.TrimSpace(req.Detail)) < 10 {
			writeError(w, "message is too short", 400)
			return
		}
	} else if conceptID == "" && lessonID == "" && question == "" {
		writeError(w, "concept_id, lesson_id, or question required", 400)
		return
	}
	// The input mode is an enum, so a bad value is a client bug worth reporting
	// rather than storing: the moderation page reads it to rebuild the field.
	inputMode := strings.TrimSpace(req.AnswerInputMode)
	if inputMode != "" && inputMode != "numeric" && inputMode != "text" {
		writeError(w, "answer_input_mode must be numeric or text", 400)
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
		ReporterID:      truncate(reporter, 128),
		ConceptID:       truncate(conceptID, 128),
		Kind:            kind,
		Question:        question,
		Expected:        truncate(req.Expected, 2000),
		Explanation:     truncate(req.Explanation, 4000),
		LessonID:        lessonID,
		Source:          truncate(source, 64),
		SessionID:       truncate(req.SessionID, 128),
		AttemptID:       truncate(req.AttemptID, 128),
		AnswerGiven:     truncate(req.AnswerGiven, 2000),
		GradingType:     truncate(req.GradingType, 64),
		AnswerHint:      truncate(req.AnswerHint, 256),
		AnswerInputMode: inputMode,
		Reason:          reason,
		Detail:          truncate(req.Detail, 2000),
		Status:          "open",
		CreatedAt:       time.Now().UTC(),
	})
	if err != nil {
		writeError(w, "failed to save report", 500)
		return
	}
	if image != nil {
		store, ok := s.repo.(storage.ReportImageStore)
		if !ok {
			log.Printf("reports: store does not implement ReportImageStore")
			writeError(w, "failed to save screenshot", 500)
			return
		}
		if err := store.SetReportImage(id, image.ContentType, image.Bytes); err != nil {
			log.Printf("reports: save image for %d: %v", id, err)
			writeError(w, "failed to save screenshot", 500)
			return
		}
	}
	writeJSON(w, map[string]interface{}{"id": id, "status": "open", "has_image": image != nil})
}
