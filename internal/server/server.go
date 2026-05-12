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
	mux.HandleFunc("/api/graph", logRequest(cors(s.handleGraph)))
	mux.HandleFunc("/api/leaderboard", logRequest(cors(s.handleLeaderboard)))
	mux.HandleFunc("/api/diagnostic", logRequest(cors(s.handleDiagnosticStart)))
	mux.HandleFunc("/api/diagnostic/answer", logRequest(cors(s.handleDiagnosticAnswer)))
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
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			http.Error(w, `{"error":"name is required"}`, 400)
			return
		}
		st, err := s.eng.CreateStudent(req.Name)
		if err != nil {
			writeError(w, "failed to create student", 500)
			return
		}
		sess, err := s.repo.CreateSession(st.ID)
		if err != nil {
			writeError(w, "failed to create session", 500)
			return
		}
		q, err := s.eng.NextQuestion(sess.ID, st.ID)
		if err != nil {
			writeError(w, "failed to get question", 500)
			return
		}
		writeJSON(w, startSessionRes{StudentID: st.ID, SessionID: sess.ID, Question: q})
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
	Token     string `json:"token"`
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
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
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name})
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
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
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
		"student_id":        st.ID,
		"name":              st.Name,
		"username":          st.Username,
		"concepts_mastered": scores.ConceptsMastered,
		"current_streak":    scores.CurrentStreak,
		"level":             scores.Level,
	})
}

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
