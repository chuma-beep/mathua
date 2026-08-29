package server

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/quiz"
	"github.com/chuma-beep/mathua/internal/scoring"
	"github.com/chuma-beep/mathua/internal/storage"
)

// MinAnswerSeconds is the minimum time in seconds a human should realistically
// take to read a question and type an answer. Answers faster than this are likely
// automated or copy-pasted submissions.
const MinAnswerSeconds = 0.3

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next(w, r)
	}
}

type Server struct {
	eng          *engine.Engine
	repo         storage.Repository
	auth         *auth.AuthService
	diagSessions map[string]*diagnostic.Session
	diagCreated  map[string]time.Time
	quizSessions map[string]*quiz.Session
	quizCreated  map[string]time.Time
	mu           sync.Mutex
	authLimiter  *rateLimiter
}

func New(eng *engine.Engine, repo storage.Repository, auth *auth.AuthService) *Server {
	s := &Server{
		eng:          eng,
		repo:         repo,
		auth:         auth,
		diagSessions: make(map[string]*diagnostic.Session),
		diagCreated:  make(map[string]time.Time),
		quizSessions: make(map[string]*quiz.Session),
		quizCreated:  make(map[string]time.Time),
		authLimiter:  newRateLimiter(5, 10, time.Minute),
	}
	// Clean up abandoned diagnostic/quiz sessions older than 1 hour
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			s.mu.Lock()
			cutoff := time.Now().Add(-1 * time.Hour)
			for id, created := range s.diagCreated {
				if created.Before(cutoff) {
					delete(s.diagSessions, id)
					delete(s.diagCreated, id)
				}
			}
			for id, created := range s.quizCreated {
				if created.Before(cutoff) {
					delete(s.quizSessions, id)
					delete(s.quizCreated, id)
				}
			}
			s.mu.Unlock()
		}
	}()
	return s
}

func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/signup", logRequest(cors(s.authLimiter.middleware(s.handleSignup))))
	mux.HandleFunc("/api/auth/login", logRequest(cors(s.authLimiter.middleware(s.handleLogin))))
	mux.HandleFunc("/api/auth/me", logRequest(cors(s.handleMe)))

	mux.HandleFunc("/api/session", logRequest(cors(s.optionalAuthMiddleware(s.handleSession))))
	mux.HandleFunc("/api/session/current", logRequest(cors(s.optionalAuthMiddleware(s.handleSessionCurrent))))
	mux.HandleFunc("/api/answer", logRequest(cors(s.optionalAuthMiddleware(s.handleAnswer))))
	mux.HandleFunc("/api/progress/", logRequest(cors(s.optionalAuthMiddleware(s.handleProgress))))
	mux.HandleFunc("/api/scores/", logRequest(cors(s.optionalAuthMiddleware(s.handleScores))))
	mux.HandleFunc("/api/config", logRequest(cors(s.handleConfig)))
	mux.HandleFunc("/api/graph", logRequest(cors(s.handleGraph)))
	mux.HandleFunc("/api/leaderboard", logRequest(cors(s.handleLeaderboard)))
	mux.HandleFunc("/api/courses", logRequest(cors(s.authMiddleware(s.handleCourses))))
	mux.HandleFunc("/api/courses/", logRequest(cors(s.authMiddleware(s.handleCourseDiagnostic))))
	mux.HandleFunc("/api/diagnostic", logRequest(cors(s.handleDiagnosticStart)))
	mux.HandleFunc("/api/diagnostic/answer", logRequest(cors(s.handleDiagnosticAnswer)))
	mux.HandleFunc("/api/goal", logRequest(cors(s.authMiddleware(s.handleGoal))))
	mux.HandleFunc("/api/goal/diagnostic", logRequest(cors(s.optionalAuthMiddleware(s.handleGoalDiagnosticStart))))
	mux.HandleFunc("/api/goal/diagnostic/answer", logRequest(cors(s.optionalAuthMiddleware(s.handleGoalDiagnosticAnswer))))
	mux.HandleFunc("/api/goal/plan", logRequest(cors(s.optionalAuthMiddleware(s.handleGoalPlan))))
	mux.HandleFunc("/api/weaknesses", logRequest(cors(s.authMiddleware(s.handleWeaknesses))))
	mux.HandleFunc("/api/goals/xp", logRequest(cors(s.authMiddleware(s.handleSetDailyXPGoal))))
	mux.HandleFunc("/api/settings", logRequest(cors(s.authMiddleware(s.handleSettings))))
	mux.HandleFunc("/api/reviews/due", logRequest(cors(s.authMiddleware(s.handleDueReviews))))
	mux.HandleFunc("/api/reviews/session", logRequest(cors(s.authMiddleware(s.handleReviewsSession))))
	mux.HandleFunc("/api/reviews/answer", logRequest(cors(s.authMiddleware(s.handleReviewsAnswer))))
	mux.HandleFunc("/api/lessons", logRequest(cors(s.handleLessons)))
	mux.HandleFunc("/api/lessons/body", logRequest(cors(s.handleLessonBody)))
	mux.HandleFunc("/api/lessons/", logRequest(cors(s.handleLessonConcept)))
	mux.HandleFunc("/api/concepts/", logRequest(cors(s.handleConceptDetail)))
	mux.HandleFunc("/api/study/answer", logRequest(cors(s.optionalAuthMiddleware(s.handleStudyAnswer))))
	mux.HandleFunc("/api/quiz/session", logRequest(cors(s.optionalAuthMiddleware(s.handleQuizSession))))
	mux.HandleFunc("/api/quiz/answer", logRequest(cors(s.optionalAuthMiddleware(s.handleQuizAnswer))))
	mux.HandleFunc("/api/health", logRequest(cors(s.handleHealth)))
	mux.HandleFunc("/api/activity", logRequest(cors(s.authMiddleware(s.handleActivity))))
}

// POST /api/session
type startSessionReq struct {
	Name string `json:"name"` // Only used when no auth token
}
type startSessionRes struct {
	StudentID string           `json:"student_id"`
	SessionID string           `json:"session_id"`
	Question  *engine.Question `json:"question"`
}

type authStudentKey struct{}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.auth == nil {
			next(w, r)
			return
		}
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, `{"error":"missing authorization"}`, 401)
			return
		}
		studentID, err := s.auth.ValidateToken(token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, 401)
			return
		}
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, authStudentKey{}, studentID))
		next(w, r)
	}
}

func (s *Server) optionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.auth == nil {
			next(w, r)
			return
		}
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			next(w, r)
			return
		}
		studentID, err := s.auth.ValidateToken(token)
		if err != nil {
			// Invalid token treated as unauthenticated — allow guest flow to proceed
			log.Printf("optionalAuth: invalid token: %v", err)
			next(w, r)
			return
		}
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, authStudentKey{}, studentID))
		next(w, r)
	}
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	// If authenticated, prefer the authenticated identity; otherwise fall back to body-provided
	// name/student_id so visitors can take diagnostics and practice without an account.
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID != "" {
		sess, err := s.repo.CreateSession(studentID)
		if err != nil {
			writeError(w, "failed to create session", 500)
			return
		}
		q, err := s.eng.NextQuestion(sess.ID, studentID)
		if err != nil {
			writeError(w, "failed to get question", 500)
			return
		}
		writeJSON(w, startSessionRes{StudentID: studentID, SessionID: sess.ID, Question: q})
		return
	}
	// Unauthenticated / guest path
	var req struct {
		Name      string `json:"name"`
		StudentID string `json:"student_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, 400)
		return
	}
	sid := req.StudentID
	if sid == "" {
		if req.Name == "" {
			http.Error(w, `{"error":"name or student_id is required"}`, 400)
			return
		}
		st, err := s.eng.CreateStudent(req.Name)
		if err != nil {
			writeError(w, "failed to create student", 500)
			return
		}
		sid = st.ID
	}
	sess, err := s.repo.CreateSession(sid)
	if err != nil {
		writeError(w, "failed to create session", 500)
		return
	}
	q, err := s.eng.NextQuestion(sess.ID, sid)
	if err != nil {
		writeError(w, "failed to get question", 500)
		return
	}
	writeJSON(w, startSessionRes{StudentID: sid, SessionID: sess.ID, Question: q})
}

// GET /api/session/current?session_id=xxx — verbatim peek of the active question
func (s *Server) handleSessionCurrent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		writeError(w, "session_id required", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		ses, _ := s.repo.GetSession(sessionID)
		if ses != nil {
			studentID = ses.StudentID
		}
	}
	if studentID == "" {
		writeError(w, "missing student id", 400)
		return
	}
	ses, _ := s.repo.GetSession(sessionID)
	if ses == nil {
		writeError(w, "session not found", 404)
		return
	}
	if ses.StudentID != studentID {
		writeError(w, "session does not belong to student", 403)
		return
	}
	q, err := s.eng.GetCurrentQuestion(sessionID, studentID)
	if err != nil {
		writeError(w, "failed to get current question", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"question": q})
}

// POST /api/answer
type answerReq struct {
	SessionID string  `json:"session_id"`
	AttemptID string  `json:"attempt_id"`
	Answer    string  `json:"answer"`
	Elapsed   float64 `json:"elapsed"`
}
type answerRes struct {
	Result       *engine.AnswerResult `json:"result"`
	NextQuestion *engine.Question     `json:"next_question"`
	Done         bool                 `json:"done"`
}

func (s *Server) handleAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req answerReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		ses, _ := s.repo.GetSession(req.SessionID)
		if ses != nil {
			studentID = ses.StudentID
		}
	}
	if studentID == "" {
		writeError(w, "missing student id", 400)
		return
	}
	sesCheck, _ := s.repo.GetSession(req.SessionID)
	if sesCheck != nil && sesCheck.StudentID != studentID {
		writeError(w, "session does not belong to student", 403)
		return
	}
	if req.Elapsed < MinAnswerSeconds {
		writeError(w, "answer submitted too quickly", 400)
		return
	}
	result, err := s.eng.SubmitAnswer(req.SessionID, studentID, req.AttemptID, req.Answer, req.Elapsed)
	if err != nil {
		if errors.Is(err, engine.ErrNoActiveQuestion) {
			writeError(w, "no active question", 409)
			return
		}
		writeError(w, "failed to submit answer", 500)
		return
	}
	next, err := s.eng.NextQuestion(req.SessionID, studentID)
	if err != nil {
		log.Printf("handleAnswer: NextQuestion failed for session %s: %v (returning result without next question)", req.SessionID, err)
		writeJSON(w, answerRes{Result: result, NextQuestion: nil, Done: true})
		return
	}
	writeJSON(w, answerRes{Result: result, NextQuestion: next, Done: next == nil})
}

// GET /api/progress/{student_id}
func (s *Server) handleProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		studentID = strings.TrimPrefix(r.URL.Path, "/api/progress/")
		if studentID == "" {
			writeError(w, "student_id required", 400)
			return
		}
	}
	progress, err := s.eng.GetProgress(studentID)
	if err != nil {
		writeError(w, "failed to get progress", 500)
		return
	}
	writeJSON(w, progress)
}

// GET /api/scores/{student_id}
func (s *Server) handleScores(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		studentID = strings.TrimPrefix(r.URL.Path, "/api/scores/")
		if studentID == "" {
			writeError(w, "student_id required", 400)
			return
		}
	}
	scores, err := s.eng.GetScores(studentID)
	if err != nil {
		writeError(w, "failed to get scores", 500)
		return
	}
	writeJSON(w, scores)
}

// GET /api/config
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]interface{}{
		"auth_enabled": s.auth != nil,
	})
}

// GET /api/graph
func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	dag := s.eng.GetDAG()
	type node struct {
		ID            string   `json:"id"`
		Label         string   `json:"label"`
		Domain        string   `json:"domain"`
		Prerequisites []string `json:"prerequisites"`
		GradingType   string   `json:"grading_type"`
	}
	nodes := make([]node, 0, dag.Count())
	for _, c := range dag.Order() {
		nodes = append(nodes, node{
			ID:            c.ID,
			Label:         c.Label,
			Domain:        c.Domain,
			Prerequisites: c.Prerequisites,
			GradingType:   c.GradingType,
		})
	}
	writeJSON(w, map[string]interface{}{"nodes": nodes, "count": len(nodes)})
}

// GET /api/leaderboard
func (s *Server) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	entries, err := s.eng.GetLeaderboard()
	if err != nil {
		writeError(w, "failed to get leaderboard", 500)
		return
	}
	writeJSON(w, entries)
}

// POST /api/diagnostic
func (s *Server) handleDiagnosticStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	sess := s.eng.StartDiagnostic()
	sess.ID = newUUID()
	_, cid, err := s.eng.NextDiagnosticQuestion(sess)
	if err != nil {
		writeError(w, "failed to get diagnostic question", 500)
		return
	}
	s.mu.Lock()
	s.diagSessions[sess.ID] = sess
	s.diagCreated[sess.ID] = time.Now()
	s.mu.Unlock()
	writeJSON(w, map[string]interface{}{
		"session_id": sess.ID,
		"concept_id": cid,
	})
}

// POST /api/diagnostic/answer
type diagAnswerReq struct {
	SessionID string `json:"session_id"`
	ConceptID string `json:"concept_id"`
	Correct   bool   `json:"correct"`
	Fast      bool   `json:"fast"`
}

func (s *Server) handleDiagnosticAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req diagAnswerReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	s.mu.Lock()
	sess := s.diagSessions[req.SessionID]
	s.mu.Unlock()
	if sess == nil {
		writeError(w, "diagnostic session not found", 404)
		return
	}
	s.eng.SubmitDiagnosticAnswer(sess, req.ConceptID, req.Correct, req.Fast)
	if s.eng.IsDiagnosticComplete(sess) {
		frontier := s.eng.DiagnosticFrontier(sess)
		report := s.eng.DiagnosticReport(sess)
		s.mu.Lock()
		delete(s.diagSessions, req.SessionID)
		s.mu.Unlock()
		writeJSON(w, map[string]interface{}{"done": true, "frontier": frontier, "report": report})
		return
	}
	_, cid, err := s.eng.NextDiagnosticQuestion(sess)
	if err != nil {
		writeError(w, "failed to get next question", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"done": false, "concept_id": cid})
}

// GET /api/health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

// GET /api/activity?days=365
func (s *Server) handleActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	days := 365
	if d := r.URL.Query().Get("days"); d != "" {
		if n, err := fmt.Sscanf(d, "%d", &days); err != nil || n != 1 || days < 1 || days > 3660 {
			writeError(w, "days must be between 1 and 3660", 400)
			return
		}
	}
	activity, err := s.repo.GetDailyActivity(studentID, days)
	if err != nil {
		writeError(w, "failed to get activity", 500)
		return
	}
	if activity == nil {
		activity = []storage.DailyActivity{}
	}
	writeJSON(w, activity)
}

// POST /api/goal
// Body: { "concept_ids": [...], "domain": "..." }
func (s *Server) handleGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		ConceptIDs []string `json:"concept_ids"`
		Domain     string   `json:"domain"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	ids := req.ConceptIDs
	if req.Domain != "" && len(ids) == 0 {
		for _, c := range s.eng.GetDAG().Order() {
			if c.Domain == req.Domain {
				ids = append(ids, c.ID)
			}
		}
	}
	if len(ids) == 0 {
		writeError(w, "no concepts specified", 400)
		return
	}
	path, err := s.eng.GetPlanner().PrerequisitesOf(ids)
	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	type conceptInfo struct {
		ID      string   `json:"id"`
		Label   string   `json:"label"`
		Domain  string   `json:"domain"`
		Prereqs []string `json:"prerequisites"`
	}
	chain := make([]conceptInfo, 0, len(path.Concepts))
	for _, c := range path.Concepts {
		chain = append(chain, conceptInfo{
			ID:      c.ID,
			Label:   c.Label,
			Domain:  c.Domain,
			Prereqs: c.Prerequisites,
		})
	}
	writeJSON(w, map[string]interface{}{
		"concepts": chain,
		"count":    len(chain),
	})
}

// POST /api/goal/diagnostic
// Body (auth): { "concept_ids": [...] }
// Body (no-auth): { "name": "...", "concept_ids": [...] }
func (s *Server) handleGoalDiagnosticStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	var req struct {
		ConceptIDs []string `json:"concept_ids"`
		Name       string   `json:"name"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if len(req.ConceptIDs) == 0 {
		writeError(w, "concept_ids required", 400)
		return
	}
	if studentID == "" && req.Name != "" {
		trimmed := strings.TrimSpace(req.Name)
		if trimmed == "" {
			writeError(w, "name is required for guest diagnostic", 400)
			return
		}
		st, err := s.eng.CreateStudent(trimmed)
		if err != nil {
			writeError(w, "failed to create student", 500)
			return
		}
		studentID = st.ID
	}
	if studentID == "" {
		writeError(w, "not authenticated: provide Authorization Bearer token or name", 401)
		return
	}
	session, question, err := s.eng.StartGoalDiagnostic(studentID, req.ConceptIDs)
	if err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	session.ID = newUUID()
	session.StudentID = studentID
	s.mu.Lock()
	s.diagSessions[session.ID] = session
	s.diagCreated[session.ID] = time.Now()
	s.mu.Unlock()
	if question == nil {
		writeJSON(w, map[string]interface{}{"session_id": session.ID, "done": true})
		return
	}
	writeJSON(w, map[string]interface{}{
		"session_id":   session.ID,
		"student_id":   studentID,
		"concept_id":   question.ConceptID,
		"concept_name": question.ConceptName,
		"question":     question.Question,
	})
}

// POST /api/goal/diagnostic/answer
// Body: { "session_id": "...", "concept_id": "...", "answer": "...", "elapsed": 0.0 }
func (s *Server) handleGoalDiagnosticAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		SessionID string  `json:"session_id"`
		ConceptID string  `json:"concept_id"`
		Answer    string  `json:"answer"`
		Elapsed   float64 `json:"elapsed"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.Elapsed < MinAnswerSeconds {
		writeError(w, "answer submitted too quickly", 400)
		return
	}
	s.mu.Lock()
	session := s.diagSessions[req.SessionID]
	s.mu.Unlock()
	if session == nil {
		writeError(w, "diagnostic session not found", 404)
		return
	}
	// If caller is authenticated, verify session ownership
	if authStudentID, _ := r.Context().Value(authStudentKey{}).(string); authStudentID != "" {
		session.Lock()
		owner := session.StudentID
		session.Unlock()
		if owner != "" && owner != authStudentID {
			writeError(w, "diagnostic session does not belong to authenticated user", 403)
			return
		}
	}

	// Grade against the stored problem (the one the user actually saw)
	if session.LastProblem == nil {
		// Session is already done — reject the answer
		writeJSON(w, map[string]interface{}{
			"done":     true,
			"correct":  false,
			"feedback": "diagnostic session already complete",
		})
		return
	}
	if session.LastConceptID != "" && session.LastConceptID != req.ConceptID {
		writeError(w, "concept_id does not match the current question", 400)
		return
	}
	concept := s.eng.GetDAG().Concept(req.ConceptID)
	timeThresh := 10.0
	if concept != nil {
		timeThresh = concept.MasteryThreshold.AvgTimeSeconds
	}
	correct := false
	explanation := ""
	session.Lock()
	gt := "numeric"
	if concept != nil {
		gt = concept.GradingType
	}
	expected := session.LastProblem.Answer
	expExplanation := session.LastProblem.Explanation
	session.Unlock()
	graderRouter := s.eng.GetGrader()
	if graderRouter != nil {
		grResult := graderRouter.Grade(grader.GradingType(gt), expected, req.Answer)
		correct = grResult.Correct
	} else {
		correct = expected == req.Answer
	}
	explanation = expExplanation
	fast := req.Elapsed < timeThresh

	s.eng.SubmitDiagnosticAnswer(session, req.ConceptID, correct, fast)
	if s.eng.IsDiagnosticComplete(session) {
		report := s.eng.DiagnosticReport(session)
		writeJSON(w, map[string]interface{}{
			"done":     true,
			"correct":  correct,
			"feedback": explanation,
			"report":   report,
		})
		return
	}
	nextProb, cid, err := s.eng.NextDiagnosticQuestion(session)
	if err != nil {
		writeError(w, "failed to get next question", 500)
		return
	}
	c := s.eng.GetDAG().Concept(cid)
	name := cid
	if c != nil {
		name = c.Label
	}
	writeJSON(w, map[string]interface{}{
		"done":         false,
		"correct":      correct,
		"feedback":     explanation,
		"concept_id":   cid,
		"concept_name": name,
		"question":     nextProb.Question,
	})
}

// POST /api/goal/plan
// Body: { "session_id": "..." }
// Returns readiness, weak areas by domain, strong areas
func (s *Server) handleGoalPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	s.mu.Lock()
	session := s.diagSessions[req.SessionID]
	s.mu.Unlock()
	if session == nil {
		writeError(w, "diagnostic session not found", 404)
		return
	}
	// If caller is authenticated, verify session ownership
	if authStudentID, _ := r.Context().Value(authStudentKey{}).(string); authStudentID != "" {
		session.Lock()
		owner := session.StudentID
		session.Unlock()
		if owner != "" && owner != authStudentID {
			writeError(w, "diagnostic session does not belong to authenticated user", 403)
			return
		}
	}
	session.Lock()
	studentID := session.StudentID
	attempts := make([]diagnostic.Attempt, len(session.Attempts))
	copy(attempts, session.Attempts)
	session.Unlock()

	// Persist diagnostic results
	if err := s.eng.ApplyGoalResults(studentID, session); err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	if err := s.repo.SetDiagnosticCompleted(studentID); err != nil {
		log.Printf("warning: failed to set diagnostic completed: %v", err)
	}

	// Clean up session
	s.mu.Lock()
	delete(s.diagSessions, req.SessionID)
	delete(s.diagCreated, req.SessionID)
	s.mu.Unlock()

	// Compute readiness and weak areas (use copied attempts for thread safety)
	weakByDomain := make(map[string][]map[string]interface{})
	strongByDomain := make(map[string][]string)

	for _, att := range attempts {
		c := s.eng.GetDAG().Concept(att.ConceptID)
		domain := "unknown"
		label := att.ConceptID
		if c != nil {
			domain = c.Domain
			label = c.Label
		}
		isWeak := (!att.Correct || !att.Fast)
		if isWeak {
			weakByDomain[domain] = append(weakByDomain[domain], map[string]interface{}{
				"id":    att.ConceptID,
				"label": label,
			})
		} else {
			strongByDomain[domain] = append(strongByDomain[domain], label)
		}
	}

	total := len(attempts)
	correct := 0
	for _, att := range attempts {
		if att.Correct {
			correct++
		}
	}
	readiness := 0.0
	if total > 0 {
		readiness = float64(correct) / float64(total)
	}

	writeJSON(w, map[string]interface{}{
		"readiness":     readiness,
		"total_tested":  total,
		"correct_count": correct,
		"weak_areas":    weakByDomain,
		"strong_areas":  strongByDomain,
	})
}

// GET /api/weaknesses
// Returns weakness scores grouped by domain for the authenticated student
func (s *Server) handleWeaknesses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	weakMap := s.eng.WeaknessMap(studentID)
	if weakMap == nil {
		writeJSON(w, map[string]interface{}{"by_domain": map[string]interface{}{}})
		return
	}
	byDomain := make(map[string][]map[string]interface{})
	for _, c := range s.eng.GetDAG().Order() {
		w := weakMap[c.ID]
		if w > 0.2 {
			byDomain[c.Domain] = append(byDomain[c.Domain], map[string]interface{}{
				"id":       c.ID,
				"label":    c.Label,
				"weakness": w,
			})
		}
	}
	writeJSON(w, map[string]interface{}{"by_domain": byDomain})
}

// GET /api/lessons  Optional: ?student_id=... for per-concept progress
func (s *Server) handleLessons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	ll := s.eng.GetLessonLoader()
	if ll == nil {
		writeJSON(w, map[string]interface{}{"lessons": []interface{}{}})
		return
	}
	byDomain := ll.LessonsByDomain()

	// Collect all concept IDs to fetch progress in one batch
	allConcepts := make(map[string]bool)
	for _, lessons := range byDomain {
		for _, l := range lessons {
			for _, cid := range l.Concepts {
				allConcepts[cid] = true
			}
		}
	}

	// Try to load progress if student_id is provided
	var progressMap map[string]map[string]interface{}
	if sid := r.URL.Query().Get("student_id"); sid != "" {
		if p, err := s.eng.GetProgress(sid); err == nil && p != nil {
			progressMap = make(map[string]map[string]interface{})
			for cid, cp := range p {
				if allConcepts[cid] {
					mastery := float64(0)
					if cp.Status == "MASTERED" {
						mastery = 1.0
					} else if cp.Status == "PRACTICING" {
						mastery = 0.6
					} else if cp.Status == "LEARNING" {
						mastery = 0.3
					}
					progressMap[cid] = map[string]interface{}{
						"status":  cp.Status,
						"mastery": mastery,
						"streak":  cp.Streak,
					}
				}
			}
		}
	}

	type prereqInfo struct {
		ID         string  `json:"id"`
		Label      string  `json:"label"`
		Status     string  `json:"status"`
		MasteryPct float64 `json:"mastery_pct"`
	}

	type lessonInfo struct {
		Title         string                            `json:"title"`
		Body          string                            `json:"body,omitempty"`
		Concepts      []string                          `json:"concepts"`
		Progress      map[string]map[string]interface{} `json:"progress,omitempty"`
		Prerequisites []prereqInfo                      `json:"prerequisites,omitempty"`
		Dependents    []prereqInfo                      `json:"dependents,omitempty"`
	}
	result := make(map[string][]lessonInfo)
	for domain, lessons := range byDomain {
		for _, l := range lessons {
			info := lessonInfo{
				Title:    l.Title,
				Concepts: l.Concepts,
			}
			if progressMap != nil {
				info.Progress = make(map[string]map[string]interface{})
				for _, cid := range l.Concepts {
					if p, ok := progressMap[cid]; ok {
						info.Progress[cid] = p
					}
				}
				if len(info.Progress) == 0 {
					info.Progress = nil
				}
			}

			// Convert progress map to the format LessonPrerequisites expects
			var rawProgress map[string]*storage.ConceptProgress
			if progressMap != nil {
				rawProgress = make(map[string]*storage.ConceptProgress)
				for cid, p := range progressMap {
					status, _ := p["status"].(string)
					streakFloat, _ := p["streak"].(float64)
					rawProgress[cid] = &storage.ConceptProgress{
						Status: status,
						Streak: int(streakFloat),
					}
				}
			}

			prereqs := s.eng.LessonPrerequisites(l.Concepts, rawProgress)
			if len(prereqs) > 0 {
				info.Prerequisites = make([]prereqInfo, len(prereqs))
				for i, p := range prereqs {
					info.Prerequisites[i] = prereqInfo{
						ID:         p.ID,
						Label:      p.Label,
						Status:     p.Status,
						MasteryPct: p.MasteryPct,
					}
				}
			}

			deps := s.eng.LessonDependents(l.Concepts, rawProgress)
			if len(deps) > 0 {
				info.Dependents = make([]prereqInfo, len(deps))
				for i, d := range deps {
					info.Dependents[i] = prereqInfo{
						ID:         d.ID,
						Label:      d.Label,
						Status:     d.Status,
						MasteryPct: d.MasteryPct,
					}
				}
			}

			result[domain] = append(result[domain], info)
		}
	}
	writeJSON(w, map[string]interface{}{"lessons": result})
}

// GET /api/lessons/{conceptId}/practice

// GET /api/lessons/body?title=... returns a single lesson's markdown body.
// The list endpoint intentionally omits bodies (multi-MB payload); clients
// fetch the selected lesson's content here on demand.
func (s *Server) handleLessonBody(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	ll := s.eng.GetLessonLoader()
	if ll == nil {
		http.Error(w, `{"error":"lessons unavailable"}`, 503)
		return
	}
	title := strings.TrimSpace(r.URL.Query().Get("title"))
	if title == "" {
		http.Error(w, `{"error":"title required"}`, 400)
		return
	}
	l := ll.LessonByTitle(title)
	if l == nil {
		http.Error(w, `{"error":"lesson not found"}`, 404)
		return
	}
	writeJSON(w, map[string]string{"title": l.Title, "body": l.Body})
}

func (s *Server) handleLessonConcept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 || (parts[1] != "practice" && parts[1] != "kp") {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	conceptID := parts[0]

	// GET /api/lessons/{id}/kp — knowledge-point shards with worked examples.
	if parts[1] == "kp" {
		ll := s.eng.GetLessonLoader()
		if ll == nil {
			http.Error(w, `{"error":"lessons unavailable"}`, 503)
			return
		}
		kps := ll.KPs(conceptID)
		type kpInfo struct {
			Label         string   `json:"label"`
			Section       string   `json:"section"`
			Subgoals      []string `json:"subgoals"`
			WorkedExample string   `json:"worked_example"`
		}
		out := make([]kpInfo, 0, len(kps))
		for _, kp := range kps {
			we, ok := ll.KPSectionBody(conceptID, kp.Section)
			if !ok {
				if l := ll.Lesson(conceptID); l != nil {
					we = l.Body
				}
			}
			out = append(out, kpInfo{
				Label:         kp.Label,
				Section:       kp.Section,
				Subgoals:      kp.Subgoals,
				WorkedExample: we,
			})
		}
		writeJSON(w, map[string]interface{}{"concept_id": conceptID, "kps": out})
		return
	}

	count := 5
	if c := r.URL.Query().Get("count"); c != "" {
		var n int
		if _, err := fmt.Sscanf(c, "%d", &n); err == nil && n > 0 && n <= 20 {
			count = n
		}
	}

	type qInfo struct {
		Question    string `json:"question"`
		Answer      string `json:"answer"`
		Explanation string `json:"explanation"`
		Source      string `json:"source,omitempty"`
	}

	// Try DB first
	dbQs, err := s.eng.GetQuestions(conceptID, count)
	if err == nil && len(dbQs) > 0 {
		questions := make([]qInfo, len(dbQs))
		for i, q := range dbQs {
			src := q.Source
			if src == "" {
				src = "curated"
			}
			questions[i] = qInfo{Question: q.Question, Answer: q.Answer, Explanation: q.Explanation, Source: src}
		}
		writeJSON(w, map[string]interface{}{"questions": questions, "concept_id": conceptID})
		return
	}

	// Fallback to generator
	reg := s.eng.GetGeneratorRegistry()
	if reg == nil {
		writeJSON(w, map[string]interface{}{"questions": []interface{}{}})
		return
	}
	problems, err := reg.BatchGenerate(conceptID, count, 0.5)
	if err != nil {
		writeError(w, err.Error(), 404)
		return
	}
	questions := make([]qInfo, len(problems))
	for i, p := range problems {
		questions[i] = qInfo{Question: p.Question, Answer: p.Answer, Explanation: p.Explanation, Source: "generator"}
	}
	writeJSON(w, map[string]interface{}{"questions": questions, "concept_id": conceptID})
}

// POST /api/goals/xp  Body: { "goal": 200 }
func (s *Server) handleSetDailyXPGoal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	var req struct {
		Goal int `json:"goal"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.Goal < 1 || req.Goal > 10000 {
		writeError(w, "goal must be between 1 and 10000", 400)
		return
	}
	if err := s.repo.SetDailyXPGoal(studentID, req.Goal); err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]interface{}{"goal": req.Goal})
}

// GET|PUT /api/settings
func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	switch r.Method {
	case http.MethodGet:
		settings, err := s.repo.GetSettings(studentID)
		if err != nil {
			writeError(w, "failed to get settings", 500)
			return
		}
		var parsed interface{}
		if err := json.Unmarshal([]byte(settings), &parsed); err != nil {
			log.Printf("warning: failed to parse settings JSON: %v; returning empty", err)
			parsed = map[string]interface{}{}
		}
		writeJSON(w, parsed)
	case http.MethodPut:
		var req map[string]interface{}
		if err := decodeJSON(w, r, &req); err != nil {
			writeError(w, "invalid request", 400)
			return
		}
		bytes, err := json.Marshal(req)
		if err != nil {
			writeError(w, "failed to encode settings", 500)
			return
		}
		if err := s.repo.UpdateSettings(studentID, string(bytes)); err != nil {
			writeError(w, "failed to update settings", 500)
			return
		}
		writeJSON(w, req)
	default:
		http.Error(w, `{"error":"method not allowed"}`, 405)
	}
}

// GET /api/reviews/due — returns count of concepts due for review
func (s *Server) handleDueReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeJSON(w, map[string]interface{}{"count": 0})
		return
	}
	progress, err := s.eng.GetProgress(studentID)
	if err != nil {
		writeJSON(w, map[string]interface{}{"count": 0})
		return
	}
	now := time.Now()
	count := 0
	for _, p := range progress {
		if p.NextReviewDue != nil && p.NextReviewDue.Before(now) && p.Status != string(mastery.StatusMastered) {
			count++
		}
	}
	writeJSON(w, map[string]interface{}{"count": count})
}

// POST /api/reviews/session — creates a review-only session
func (s *Server) handleReviewsSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, ok := r.Context().Value(authStudentKey{}).(string)
	if !ok || studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	sess, err := s.repo.CreateSession(studentID)
	if err != nil {
		writeError(w, "failed to create session", 500)
		return
	}
	q, err := s.eng.NextReviewQuestion(sess.ID, studentID)
	if err != nil {
		writeError(w, "failed to get review question", 500)
		return
	}
	writeJSON(w, startSessionRes{StudentID: studentID, SessionID: sess.ID, Question: q})
}

// POST /api/reviews/answer — submits a review answer and returns the next review question
func (s *Server) handleReviewsAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req answerReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "missing student id", 400)
		return
	}
	sesCheck, _ := s.repo.GetSession(req.SessionID)
	if sesCheck != nil && sesCheck.StudentID != studentID {
		writeError(w, "session does not belong to student", 403)
		return
	}
	if req.Elapsed < MinAnswerSeconds {
		writeError(w, "answer submitted too quickly", 400)
		return
	}
	result, err := s.eng.SubmitAnswer(req.SessionID, studentID, req.AttemptID, req.Answer, req.Elapsed)
	if err != nil {
		if errors.Is(err, engine.ErrNoActiveQuestion) {
			writeError(w, "no active question", 409)
			return
		}
		writeError(w, "failed to submit answer", 500)
		return
	}
	next, err := s.eng.NextReviewQuestion(req.SessionID, studentID)
	if err != nil {
		log.Printf("handleReviewsAnswer: NextReviewQuestion failed for session %s: %v (returning result without next question)", req.SessionID, err)
		writeJSON(w, answerRes{Result: result, NextQuestion: nil, Done: true})
		return
	}
	writeJSON(w, answerRes{Result: result, NextQuestion: next, Done: next == nil})
}

// GET /api/courses
func (s *Server) handleCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	courses := s.eng.PlannerCourses()
	writeJSON(w, map[string]interface{}{"courses": courses})
}

// POST /api/courses/{id}/diagnostic
func (s *Server) handleCourseDiagnostic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	courseID := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	courseID = strings.TrimSuffix(courseID, "/diagnostic")
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	path, err := s.eng.SetActiveCourse(studentID, courseID)
	if err != nil {
		writeError(w, "failed to set course: "+err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{
		"course_id":   courseID,
		"path_length": len(path.Concepts),
	})
}

// GET /api/concepts/{id}
func (s *Server) handleConceptDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	conceptID := strings.TrimPrefix(r.URL.Path, "/api/concepts/")
	if conceptID == "" {
		writeError(w, "missing concept id", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	detail, err := s.eng.ConceptDetail(studentID, conceptID)
	if err != nil {
		writeError(w, err.Error(), 404)
		return
	}
	writeJSON(w, detail)
}

// POST /api/study/answer — Study seam: LessonQuiz → SubmitAnswer (CONTEXT.md Seam)
// Body: { concept_id, answer, expected, elapsed, student_id? } student_id used for guest (mathua_guest_id)
func (s *Server) handleStudyAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		ConceptID string  `json:"concept_id"`
		Answer    string  `json:"answer"`
		Expected  string  `json:"expected"`
		Elapsed   float64 `json:"elapsed"`
		StudentID string  `json:"student_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.ConceptID == "" {
		writeError(w, "concept_id required", 400)
		return
	}
	if req.Expected == "" {
		writeError(w, "expected required", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		studentID = req.StudentID
	}
	// Allow unauthenticated without student_id: grade only, no persistence.
	if studentID == "" {
		concept := s.eng.GetDAG().Concept(req.ConceptID)
		gt := "numeric"
		if concept != nil {
			gt = concept.GradingType
		}
		gr := s.eng.GetGrader().Grade(grader.GradingType(gt), req.Expected, req.Answer)
		writeJSON(w, map[string]interface{}{
			"correct":  gr.Correct,
			"feedback": gr.Feedback,
			"xp":       0,
		})
		return
	}
	if req.Elapsed < 0 {
		req.Elapsed = 0
	}
	if req.Elapsed > 600 {
		req.Elapsed = 600
	}
	res, err := s.eng.SubmitStudyAnswer(studentID, req.ConceptID, req.Answer, req.Expected, req.Elapsed)
	if err != nil {
		writeError(w, "failed to submit study answer", 500)
		return
	}
	writeJSON(w, map[string]interface{}{
		"correct":         res.Correct,
		"feedback":        res.Feedback,
		"explanation":     res.Explanation,
		"new_status":      res.NewStatus,
		"streak":          res.Streak,
		"required_streak": res.RequiredStreak,
		"xp":              res.XP,
		"expected_answer": res.ExpectedAnswer,
		"halted":          res.Halted,
	})
}

// POST /api/quiz/session — actionable quiz every 150 XP, own grading path, guest unlimited retake
// Body: { student_id? } optional for guest
func (s *Server) handleQuizSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		StudentID string `json:"student_id"`
	}
	_ = decodeJSON(w, r, &req)
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		studentID = req.StudentID
	}
	// Build quiz picking diverse recent weak concepts; fallback to DAG order
	var quizConcepts []*concepts.Concept
	if studentID != "" {
		weak := s.eng.WeaknessMap(studentID)
		if weak != nil {
			quizConcepts = quiz.PickQuizConcepts(s.eng.GetDAG(), weak)
		}
	}
	qEng := quiz.NewEngine(s.eng.GetDAG(), s.eng.GetGeneratorRegistry())
	sess := qEng.Start(quizConcepts)
	// PR 1.4: per-concept 80% difficulty target from weakness map.
	if studentID != "" {
		for _, c := range quizConcepts {
			sess.Difficulties[c.ID] = s.eng.DifficultyFor(studentID, c.ID)
		}
	}
	sess.ID = newUUID()
	sess.StudentID = studentID
	s.mu.Lock()
	s.quizSessions[sess.ID] = sess
	s.quizCreated[sess.ID] = time.Now()
	s.mu.Unlock()
	prob, cid, err := qEng.NextQuestion(sess)
	if err != nil {
		writeError(w, "failed to get quiz question", 500)
		return
	}
	if prob == nil || cid == "" {
		writeJSON(w, map[string]interface{}{"session_id": sess.ID, "done": true})
		return
	}
	c := s.eng.GetDAG().Concept(cid)
	name := cid
	if c != nil {
		name = c.Label
	}
	writeJSON(w, map[string]interface{}{
		"session_id":   sess.ID,
		"student_id":   studentID,
		"concept_id":   cid,
		"concept_name": name,
		"question":     prob.Question,
		"done":         false,
	})
}

// POST /api/quiz/answer — own grading path via SubmitStudyAnswer with TaskQuiz 20
// Body: { session_id, concept_id, answer, elapsed, student_id? }
func (s *Server) handleQuizAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		SessionID string  `json:"session_id"`
		ConceptID string  `json:"concept_id"`
		Answer    string  `json:"answer"`
		Elapsed   float64 `json:"elapsed"`
		StudentID string  `json:"student_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.SessionID == "" || req.ConceptID == "" {
		writeError(w, "session_id and concept_id required", 400)
		return
	}
	if req.Elapsed < MinAnswerSeconds {
		writeError(w, "answer submitted too quickly", 400)
		return
	}
	s.mu.Lock()
	sess := s.quizSessions[req.SessionID]
	s.mu.Unlock()
	if sess == nil {
		writeError(w, "quiz session not found", 404)
		return
	}
	if sess.LastProblem == nil {
		writeJSON(w, map[string]interface{}{"done": true, "correct": false, "feedback": "quiz session already complete"})
		return
	}
	if sess.LastConceptID != "" && sess.LastConceptID != req.ConceptID {
		writeError(w, "concept_id does not match current question", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		studentID = req.StudentID
		if studentID == "" {
			studentID = sess.StudentID
		}
	}
	// Grade via DAG grading_type
	concept := s.eng.GetDAG().Concept(req.ConceptID)
	gt := "numeric"
	if concept != nil {
		gt = concept.GradingType
	}
	expected := sess.LastProblem.Answer
	gr := s.eng.GetGrader().Grade(grader.GradingType(gt), expected, req.Answer)
	// Record via quiz engine (advance index)
	qEng := quiz.NewEngine(s.eng.GetDAG(), s.eng.GetGeneratorRegistry())
	qEng.RecordAnswer(sess, req.ConceptID, gr.Correct)
	// Also persist XP/progress via SubmitStudyAnswer with TaskQuiz override if student known
	xp := 0
	var newStatus string
	if studentID != "" {
		// Use dedicated quiz XP: TaskQuiz 20
		// Temporarily call SubmitStudyAnswer then override base via re-grade? Instead directly compute TaskQuiz XP here
		// We call SubmitStudyAnswer for progress but it would award TaskLesson/Multistep; we want TaskQuiz.
		// So we call engine helper for quiz XP: use SubmitStudyAnswer then patch XP to TaskQuiz 20 equivalent
		// Simpler: call SubmitStudyAnswer then recompute XP as TaskQuiz
		res, err := s.eng.SubmitStudyAnswer(studentID, req.ConceptID, req.Answer, expected, req.Elapsed)
		if err == nil && res != nil {
			// Override XP to TaskQuiz 20 (recompute)
			conceptThresh := 10.0
			if concept != nil {
				conceptThresh = concept.MasteryThreshold.AvgTimeSeconds
			}
			// recompute with TaskQuiz base
			// We already have res.XP as lesson/multistep; recompute correctly
			// Use same streak from res
			xp = res.XP
			// If TaskQuiz differs, scale: TaskQuiz 20 vs TaskLesson 10 => double
			if res.Correct {
				// recompute TaskQuiz XP properly
				// Inline compute to avoid re-calling private: approximate via engine's exported helper? Use simple ratio
				// We know base 20 vs base from concept suffix; but just call with TaskQuiz via reflection: we can ask engine to compute
				// For MVP, double XP if not word else 20/15
				// Instead call engine's public TaskQuiz compute via new helper
				xp = s.eng.QuizXP(res.Correct, req.Elapsed, conceptThresh, res.Streak)
			}
			newStatus = string(res.NewStatus)
		} else {
			xp = 0
		}
	}
	if qEng.IsComplete(sess) {
		s.mu.Lock()
		delete(s.quizSessions, req.SessionID)
		delete(s.quizCreated, req.SessionID)
		s.mu.Unlock()
		writeJSON(w, map[string]interface{}{"done": true, "correct": gr.Correct, "feedback": gr.Feedback, "xp": xp, "new_status": newStatus})
		return
	}
	prob, cid, err := qEng.NextQuestion(sess)
	if err != nil {
		writeError(w, "failed to get next quiz question", 500)
		return
	}
	c2 := s.eng.GetDAG().Concept(cid)
	name2 := cid
	if c2 != nil {
		name2 = c2.Label
	}
	writeJSON(w, map[string]interface{}{
		"done":         false,
		"correct":      gr.Correct,
		"feedback":     gr.Feedback,
		"xp":           xp,
		"new_status":   newStatus,
		"concept_id":   cid,
		"concept_name": name2,
		"question":     prob.Question,
	})
}

const maxBodySize int64 = 1 << 20 // 1 MB

func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("warning: failed to encode JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

type signupReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type authRes struct {
	Token               string `json:"token"`
	StudentID           string `json:"student_id"`
	Name                string `json:"name"`
	DiagnosticCompleted bool   `json:"diagnostic_completed"`
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	var req signupReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.Username == "" || req.Password == "" || req.Name == "" {
		writeError(w, "name, username, and password are required", 400)
		return
	}
	token, st, err := s.auth.Signup(req.Name, req.Username, req.Password)
	if err != nil {
		writeError(w, "signup failed: "+err.Error(), 400)
		return
	}
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name, DiagnosticCompleted: st.DiagnosticCompleted})
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	var req loginReq
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	token, st, err := s.auth.Login(req.Username, req.Password)
	if err != nil || token == "" {
		writeError(w, "invalid username or password", 401)
		return
	}
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name, DiagnosticCompleted: st.DiagnosticCompleted})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		writeError(w, "missing authorization", 401)
		return
	}
	studentID, err := s.auth.ValidateToken(token)
	if err != nil {
		writeError(w, "invalid token", 401)
		return
	}
	st, err := s.repo.GetStudent(studentID)
	if err != nil || st == nil {
		writeError(w, "student not found", 404)
		return
	}
	scores, err := s.eng.GetScores(studentID)
	if err != nil {
		log.Printf("warning: failed to get scores for %s: %v", studentID, err)
		scores = &scoring.Scores{}
	}
	writeJSON(w, map[string]interface{}{
		"student_id":           st.ID,
		"name":                 st.Name,
		"username":             st.Username,
		"concepts_mastered":    scores.ConceptsMastered,
		"current_streak":       scores.CurrentStreak,
		"level":                scores.Level,
		"diagnostic_completed": st.DiagnosticCompleted,
	})
}

func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Printf("warning: crypto/rand.Read failed: %v; using time-based UUID", err)
		return fmt.Sprintf("uuid-%x", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// gzipResponseWriter transparently gzips responses for clients that send
// Accept-Encoding: gzip. JSON lesson/graph payloads compress ~8-10x.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if g.Header().Get("Content-Type") == "" {
		g.Header().Set("Content-Type", "application/json")
	}
	return g.gz.Write(b)
}

// GzipMiddleware compresses responses for gzip-capable clients.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")
		gz := gzip.NewWriter(w)
		defer func() {
			if err := gz.Close(); err != nil {
				log.Printf("warning: gzip close: %v", err)
			}
		}()
		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, gz: gz}, r)
	})
}
