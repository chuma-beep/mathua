package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/storage"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
	mu           sync.Mutex
}

func New(eng *engine.Engine, repo storage.Repository, auth *auth.AuthService) *Server {
	return &Server{
		eng:          eng,
		repo:         repo,
		auth:         auth,
		diagSessions: make(map[string]*diagnostic.Session),
	}
}

func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/signup", logRequest(cors(s.handleSignup)))
	mux.HandleFunc("/api/auth/login", logRequest(cors(s.handleLogin)))
	mux.HandleFunc("/api/auth/me", logRequest(cors(s.handleMe)))

	mux.HandleFunc("/api/session", logRequest(cors(s.authMiddleware(s.handleSession))))
	mux.HandleFunc("/api/answer", logRequest(cors(s.authMiddleware(s.handleAnswer))))
	mux.HandleFunc("/api/progress/", logRequest(cors(s.authMiddleware(s.handleProgress))))
	mux.HandleFunc("/api/scores/", logRequest(cors(s.authMiddleware(s.handleScores))))
	mux.HandleFunc("/api/config", logRequest(cors(s.handleConfig)))
	mux.HandleFunc("/api/graph", logRequest(cors(s.handleGraph)))
	mux.HandleFunc("/api/leaderboard", logRequest(cors(s.handleLeaderboard)))
	mux.HandleFunc("/api/courses", logRequest(cors(s.authMiddleware(s.handleCourses))))
	mux.HandleFunc("/api/courses/", logRequest(cors(s.authMiddleware(s.handleCourseDiagnostic))))
	mux.HandleFunc("/api/diagnostic", logRequest(cors(s.handleDiagnosticStart)))
	mux.HandleFunc("/api/diagnostic/answer", logRequest(cors(s.handleDiagnosticAnswer)))
	mux.HandleFunc("/api/goal", logRequest(cors(s.authMiddleware(s.handleGoal))))
	mux.HandleFunc("/api/goal/diagnostic", logRequest(cors(s.authMiddleware(s.handleGoalDiagnosticStart))))
	mux.HandleFunc("/api/goal/diagnostic/answer", logRequest(cors(s.authMiddleware(s.handleGoalDiagnosticAnswer))))
	mux.HandleFunc("/api/goal/plan", logRequest(cors(s.authMiddleware(s.handleGoalPlan))))
	mux.HandleFunc("/api/weaknesses", logRequest(cors(s.authMiddleware(s.handleWeaknesses))))
	mux.HandleFunc("/api/goals/xp", logRequest(cors(s.authMiddleware(s.handleSetDailyXPGoal))))
	mux.HandleFunc("/api/settings", logRequest(cors(s.authMiddleware(s.handleSettings))))
	mux.HandleFunc("/api/lessons", logRequest(cors(s.handleLessons)))
	mux.HandleFunc("/api/lessons/", logRequest(cors(s.handleLessonConcept)))
	mux.HandleFunc("/api/concepts/", logRequest(cors(s.handleConceptDetail)))
	mux.HandleFunc("/api/health", logRequest(cors(s.handleHealth)))
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

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	// If no auth, accept name from body (testing / legacy)
	if s.auth == nil {
		var req struct {
			Name      string `json:"name"`
			StudentID string `json:"student_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, 400)
			return
		}
		studentID := req.StudentID
		if studentID == "" {
			if req.Name == "" {
				http.Error(w, `{"error":"name or student_id is required"}`, 400)
				return
			}
			st, err := s.eng.CreateStudent(req.Name)
			if err != nil {
				writeError(w, "failed to create student", 500)
				return
			}
			studentID = st.ID
		}
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
	studentID := r.Context().Value(authStudentKey{}).(string)
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
}

// POST /api/answer
type answerReq struct {
	SessionID string  `json:"session_id"`
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" && s.auth == nil {
		ses, _ := s.repo.GetSession(req.SessionID)
		if ses != nil {
			studentID = ses.StudentID
		}
	}
	if studentID == "" {
		writeError(w, "missing student id", 400)
		return
	}
	result, err := s.eng.SubmitAnswer(req.SessionID, studentID, req.Answer, req.Elapsed)
	if err != nil {
		writeError(w, "failed to submit answer", 500)
		return
	}
	next, err := s.eng.NextQuestion(req.SessionID, studentID)
	if err != nil {
		writeError(w, "failed to get next question", 500)
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
	if studentID == "" && s.auth == nil {
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
	if studentID == "" && s.auth == nil {
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
		s.mu.Lock()
		delete(s.diagSessions, req.SessionID)
		s.mu.Unlock()
		writeJSON(w, map[string]interface{}{"done": true, "frontier": frontier})
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if len(req.ConceptIDs) == 0 {
		writeError(w, "concept_ids required", 400)
		return
	}
	if studentID == "" && req.Name != "" && s.auth == nil {
		st, err := s.eng.CreateStudent(req.Name)
		if err != nil {
			writeError(w, "failed to create student", 500)
			return
		}
		studentID = st.ID
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	// Grade against the stored problem (the one the user actually saw)
	concept := s.eng.GetDAG().Concept(req.ConceptID)
	timeThresh := 10.0
	if concept != nil {
		timeThresh = concept.MasteryThreshold.AvgTimeSeconds
	}
	correct := false
	explanation := ""
	if session.LastProblem != nil {
		gt := "numeric"
		if concept != nil {
			gt = concept.GradingType
		}
		graderRouter := s.eng.GetGrader()
		if graderRouter != nil {
			grResult := graderRouter.Grade(grader.GradingType(gt), session.LastProblem.Answer, req.Answer)
			correct = grResult.Correct
		} else {
			correct = session.LastProblem.Answer == req.Answer
		}
		explanation = session.LastProblem.Explanation
	}
	fast := req.Elapsed < timeThresh

	s.eng.SubmitDiagnosticAnswer(session, req.ConceptID, correct, fast)
	if s.eng.IsDiagnosticComplete(session) {
		writeJSON(w, map[string]interface{}{
			"done":     true,
			"correct":  correct,
			"feedback": explanation,
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	studentID := session.StudentID

	// Persist diagnostic results
	if err := s.eng.ApplyGoalResults(studentID, session); err != nil {
		writeError(w, err.Error(), 500)
		return
	}
	_ = s.repo.SetDiagnosticCompleted(studentID)

	// Clean up session
	s.mu.Lock()
	delete(s.diagSessions, req.SessionID)
	s.mu.Unlock()

	// Compute readiness and weak areas
	weakByDomain := make(map[string][]map[string]interface{})
	strongByDomain := make(map[string][]string)

	for _, att := range session.Attempts {
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

	total := len(session.Attempts)
	correct := 0
	for _, att := range session.Attempts {
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

	type lessonInfo struct {
		Title    string                          `json:"title"`
		Body     string                          `json:"body"`
		Concepts []string                        `json:"concepts"`
		Progress map[string]map[string]interface{} `json:"progress,omitempty"`
	}
	result := make(map[string][]lessonInfo)
	for domain, lessons := range byDomain {
		for _, l := range lessons {
			info := lessonInfo{
				Title:    l.Title,
				Body:     l.Body,
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
			result[domain] = append(result[domain], info)
		}
	}
	writeJSON(w, map[string]interface{}{"lessons": result})
}

// GET /api/lessons/{conceptId}/practice
func (s *Server) handleLessonConcept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 || parts[1] != "practice" {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	conceptID := parts[0]
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
		json.Unmarshal([]byte(settings), &parsed)
		writeJSON(w, parsed)
	case http.MethodPut:
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	scores, _ := s.eng.GetScores(studentID)
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
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
