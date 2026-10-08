package server

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/admin"
	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/diagnostic"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/lessons"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/quiz"
	"github.com/chuma-beep/mathua/internal/scoring"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/xp"
)

// MinAnswerSeconds is the minimum time in seconds a human should realistically
// take to read a question and type an answer. Answers faster than this are likely
// automated or copy-pasted submissions.
const MinAnswerSeconds = 0.3

var (
	allowedOriginsOnce     sync.Once
	allowedOriginsExact    map[string]bool
	allowedOriginPatterns  []string
	allowedOriginsAllowAll bool
)

func parseAllowedOrigins() {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")
	raw = strings.TrimSpace(raw)
	if raw == "" {
		allowedOriginsAllowAll = true
		allowedOriginsExact = nil
		allowedOriginPatterns = nil
		return
	}
	allowedOriginsAllowAll = false
	allowedOriginsExact = make(map[string]bool)
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			allowedOriginsAllowAll = true
			continue
		}
		if strings.Contains(o, "*") {
			allowedOriginPatterns = append(allowedOriginPatterns, o)
		} else {
			allowedOriginsExact[o] = true
		}
	}
	if len(allowedOriginsExact) == 0 && len(allowedOriginPatterns) == 0 && !allowedOriginsAllowAll {
		allowedOriginsAllowAll = true
	}
}

func isAllowedOrigin(origin string) bool {
	allowedOriginsOnce.Do(parseAllowedOrigins)
	if allowedOriginsAllowAll {
		return true
	}
	if allowedOriginsExact != nil && allowedOriginsExact[origin] {
		return true
	}
	for _, pat := range allowedOriginPatterns {
		if matchOriginPattern(origin, pat) {
			return true
		}
	}
	return false
}

// matchOriginPattern supports "*" wildcards via prefix/suffix split.
// e.g. "*.vercel.app" matches "https://foo.vercel.app"
//
//	"https://*.vercel.app" matches "https://mathua.vercel.app" but not "http://..."
func matchOriginPattern(origin, pattern string) bool {
	if !strings.Contains(pattern, "*") {
		return origin == pattern
	}
	if pattern == "*" {
		return true
	}
	prefix, suffix, _ := strings.Cut(pattern, "*")
	// Only single wildcard is supported; treat multiple "*" as requiring both affixes.
	if strings.Contains(suffix, "*") {
		// fallback: require prefix and suffix of first wildcard plus that suffix contains "*"
		// simple check: origin must have prefix and suffix around the first "*"
		return strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix)
	}
	return strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix)
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Vary", "Origin")
			if isAllowedOrigin(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == http.MethodOptions {
			// Preflight: 204 with CORS headers already set
			w.WriteHeader(http.StatusNoContent)
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
	// Admin triage sessions live in server_sessions (kind "admin"), not in
	// memory: logins survive restarts and work behind a second replica.
	// Diag/quiz sessions stay in memory by design (Fix 6): short-lived
	// capability UUIDs swept hourly; a restart just means a retake.
	mu           sync.Mutex
	authLimiter  *rateLimiter
	writeLimiter *rateLimiter
	shareLimiter *rateLimiter
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
		writeLimiter: newRateLimiter(20, 20, 3*time.Second),
		shareLimiter: newRateLimiter(10, 10, 6*time.Second),
	}
	// MATHUA_LOADTEST=1 relaxes rate limits so load generators measure the
	// app and database instead of the limiter. Local/staging use only —
	// never set in production.
	if os.Getenv("MATHUA_LOADTEST") == "1" {
		log.Print("loadtest mode: rate limits relaxed")
		s.authLimiter = newRateLimiter(100000, 100000, time.Second)
		s.writeLimiter = newRateLimiter(100000, 100000, time.Second)
		s.shareLimiter = newRateLimiter(100000, 100000, time.Second)
	}
	// Clean up abandoned diagnostic/quiz sessions older than 1 hour, plus
	// expired durable server_sessions rows (admin logins, study anchors).
	//nolint:goroutinelint // process-lifetime janitor: runs until the process exits
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
			if err := s.repo.SweepServerSessions(); err != nil {
				log.Printf("warning: sweep server sessions: %v", err)
			}
		}
	}()
	return s
}

// dropStudentLiveSessions evicts every in-memory diagnostic and quiz session
// belonging to studentID.
//
// The age-based janitor above only sweeps after an hour, which is the right
// default for abandoned tabs and the wrong one for a reset: a learner who
// wipes their record and then reopens the tab they left mid-diagnostic gets it
// all back, because handleGoalDiagnosticAnswer re-inserts the attempt via
// EnsureSession + RecordAttempt, and handleQuizAnswer writes concept_progress
// and XP through SubmitQuizAnswer. Reset and delete must therefore evict these
// eagerly — a wipe that a stale in-memory session can undo is not a wipe.
func (s *Server) dropStudentLiveSessions(studentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, sess := range s.diagSessions {
		if sess != nil && sess.StudentID == studentID {
			delete(s.diagSessions, id)
			delete(s.diagCreated, id)
		}
	}
	for id, sess := range s.quizSessions {
		if sess != nil && sess.StudentID == studentID {
			delete(s.quizSessions, id)
			delete(s.quizCreated, id)
		}
	}
}

func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/signup", logRequest(cors(s.authLimiter.middleware(s.handleSignup))))
	mux.HandleFunc("/api/auth/login", logRequest(cors(s.authLimiter.middleware(s.handleLogin))))
	mux.HandleFunc("/api/auth/guest", logRequest(cors(s.authLimiter.middleware(s.handleGuestToken))))
	mux.HandleFunc("/api/auth/password", logRequest(cors(s.authLimiter.middleware(s.authMiddleware(s.handleChangePassword)))))
	mux.HandleFunc("/api/auth/reset/request", logRequest(cors(s.authLimiter.middleware(s.handleResetRequest))))
	mux.HandleFunc("/api/auth/reset/complete", logRequest(cors(s.authLimiter.middleware(s.handleResetComplete))))
	mux.HandleFunc("/api/auth/google", logRequest(cors(s.authLimiter.middleware(s.handleGoogleOneTap))))
	mux.HandleFunc("/api/auth/google/login", logRequest(cors(s.authLimiter.middleware(s.handleGoogleLogin))))
	mux.HandleFunc("/api/auth/google/callback", logRequest(cors(s.authLimiter.middleware(s.handleGoogleCallback))))
	mux.HandleFunc("/api/auth/google/connect", logRequest(cors(s.authLimiter.middleware(s.authMiddleware(s.handleGoogleConnect)))))
	// Generic OAuth providers (github/facebook/microsoft/apple).
	mux.HandleFunc("/api/auth/github/login", logRequest(cors(s.authLimiter.middleware(s.handleOAuthLogin))))
	mux.HandleFunc("/api/auth/github/callback", logRequest(cors(s.authLimiter.middleware(s.handleOAuthCallback))))
	mux.HandleFunc("/api/auth/facebook/login", logRequest(cors(s.authLimiter.middleware(s.handleOAuthLogin))))
	mux.HandleFunc("/api/auth/facebook/callback", logRequest(cors(s.authLimiter.middleware(s.handleOAuthCallback))))
	mux.HandleFunc("/api/auth/microsoft/login", logRequest(cors(s.authLimiter.middleware(s.handleOAuthLogin))))
	mux.HandleFunc("/api/auth/microsoft/callback", logRequest(cors(s.authLimiter.middleware(s.handleOAuthCallback))))
	mux.HandleFunc("/api/auth/apple/login", logRequest(cors(s.authLimiter.middleware(s.handleOAuthLogin))))
	mux.HandleFunc("/api/auth/apple/callback", logRequest(cors(s.authLimiter.middleware(s.handleOAuthCallback))))
	mux.HandleFunc("/api/auth/link-token", logRequest(cors(s.authMiddleware(s.handleLinkToken))))
	mux.HandleFunc("/api/auth/identities", logRequest(cors(s.authMiddleware(s.handleIdentities))))
	mux.HandleFunc("/api/auth/identities/", logRequest(cors(s.authMiddleware(s.handleIdentityDelete))))
	mux.HandleFunc("/api/auth/email/request", logRequest(cors(s.authMiddleware(s.handleEmailRequest))))
	mux.HandleFunc("/api/auth/email/verify", logRequest(cors(s.handleEmailVerify)))
	mux.HandleFunc("/api/auth/me", logRequest(cors(s.handleMe)))
	// Alias: older frontend bundles validate against /api/me (same handler).
	mux.HandleFunc("/api/me", logRequest(cors(s.handleMe)))
	mux.HandleFunc("/api/profile", logRequest(cors(s.authMiddleware(s.handleProfileUpdate))))

	mux.HandleFunc("/api/progress/", logRequest(cors(s.optionalAuthMiddleware(s.handleProgress))))
	mux.HandleFunc("/api/scores/", logRequest(cors(s.optionalAuthMiddleware(s.handleScores))))
	mux.HandleFunc("/api/config", logRequest(cors(getOnly(s.handleConfig))))
	mux.HandleFunc("/api/graph", logRequest(cors(s.handleGraph)))
	mux.HandleFunc("/api/curriculum/domains", logRequest(cors(s.optionalAuthMiddleware(s.handleCurriculumDomains))))
	mux.HandleFunc("/api/curriculum/domains/", logRequest(cors(s.optionalAuthMiddleware(s.handleCurriculumDomain))))
	mux.HandleFunc("/api/leaderboard", logRequest(cors(s.handleLeaderboard)))
	mux.HandleFunc("/api/leagues", logRequest(cors(s.authMiddleware(s.handleLeagues))))
	mux.HandleFunc("/api/share", logRequest(cors(s.authMiddleware(s.handleShareToggle))))
	mux.HandleFunc("/api/share/", logRequest(cors(s.shareLimiter.middleware(s.handleShareReport))))
	mux.HandleFunc("/api/courses", logRequest(cors(s.authMiddleware(s.handleCourses))))
	mux.HandleFunc("/api/courses/", logRequest(cors(s.authMiddleware(s.handleCourseRoute))))
	mux.HandleFunc("/api/destinations", logRequest(cors(s.authMiddleware(s.handleDestinations))))
	mux.HandleFunc("/api/destinations/", logRequest(cors(s.authMiddleware(s.handleDestinationEstimate))))
	mux.HandleFunc("/api/plans", logRequest(cors(s.authMiddleware(s.handlePlanSave))))
	mux.HandleFunc("/api/plans/current", logRequest(cors(s.authMiddleware(s.handlePlanCurrent))))
	mux.HandleFunc("/api/transcript", logRequest(cors(s.authMiddleware(s.handleTranscript))))
	mux.HandleFunc("/api/attempts", logRequest(cors(s.authMiddleware(s.handleAttempts))))
	mux.HandleFunc("/api/goal", logRequest(cors(s.authMiddleware(s.handleGoal))))
	mux.HandleFunc("/api/goal/diagnostic", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleGoalDiagnosticStart)))))
	mux.HandleFunc("/api/goal/diagnostic/answer", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleGoalDiagnosticAnswer)))))
	mux.HandleFunc("/api/goal/diagnostic/skip", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleGoalDiagnosticSkip)))))
	mux.HandleFunc("/api/goal/diagnostic/retry", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleGoalDiagnosticRetry)))))
	mux.HandleFunc("/api/goal/diagnostic/resume", logRequest(cors(s.optionalAuthMiddleware(s.handleGoalDiagnosticResume))))
	mux.HandleFunc("/api/goal/plan", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleGoalPlan)))))
	mux.HandleFunc("/api/weaknesses", logRequest(cors(s.authMiddleware(s.handleWeaknesses))))
	mux.HandleFunc("/api/goals/xp", logRequest(cors(s.authMiddleware(s.handleSetDailyXPGoal))))
	mux.HandleFunc("/api/account/reset", logRequest(cors(s.writeLimiter.middleware(s.authMiddleware(s.handleAccountReset)))))
	mux.HandleFunc("/api/account", logRequest(cors(s.writeLimiter.middleware(s.authMiddleware(s.handleAccountDelete)))))
	mux.HandleFunc("/api/settings", logRequest(cors(s.authMiddleware(s.handleSettings))))
	mux.HandleFunc("/api/avatar", logRequest(cors(s.writeLimiter.middleware(s.authMiddleware(s.handleAvatar)))))
	mux.HandleFunc("/api/avatar/me", logRequest(cors(s.authMiddleware(s.handleAvatarMe))))
	// Public per-student photo for leaderboard avatars (custom uploads are
	// board-visible by design; 404 when the student has none).
	mux.HandleFunc("/api/avatar/", logRequest(cors(s.handleAvatarPublic)))
	mux.HandleFunc("/api/reviews/due", logRequest(cors(s.authMiddleware(s.handleDueReviews))))
	// GET /api/next — the one answer to "what should I do now?". Optional auth so a guest
	// landing on the app gets the diagnostic recommendation rather than an error.
	mux.HandleFunc("/api/next", logRequest(cors(s.optionalAuthMiddleware(s.handleNext))))
	mux.HandleFunc("/api/reviews/session", logRequest(cors(s.writeLimiter.middleware(s.authMiddleware(s.handleReviewsSession)))))
	mux.HandleFunc("/api/reviews/answer", logRequest(cors(s.writeLimiter.middleware(s.authMiddleware(s.handleReviewsAnswer)))))
	// The lesson routes stay publicly readable — Study is reference and is
	// indexed — so auth is optional, but it must be *present*: the practice
	// branch writes the server-side answer anchor, and it can only key that
	// anchor by student once the bearer token has been resolved into the
	// request context. Without the middleware here the anchor was never
	// written, and every study answer fell back to the client's own expected.
	mux.HandleFunc("/api/lessons", logRequest(cors(s.optionalAuthMiddleware(s.handleLessons))))
	mux.HandleFunc("/api/lessons/body", logRequest(cors(s.optionalAuthMiddleware(s.handleLessonBody))))
	mux.HandleFunc("/api/lessons/", logRequest(cors(s.optionalAuthMiddleware(s.handleLessonConcept))))
	mux.HandleFunc("/api/study/answer", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleStudyAnswer)))))
	// Report creation only. Listing and resolving reports is administration now and lives under
	// /api/admin/reports behind the role boundary; the old shared-password triage endpoints were
	// removed in V2 so there is exactly one moderation path.
	mux.HandleFunc("/api/reports", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleCreateReport)))))
	// Admin V2. Every route names the permission it needs, and requirePermission enforces it
	// server-side on every request from the role the database holds — never from anything the
	// client sends. The password login above is the old content-triage gate; it is not this
	// boundary and does not open any of these routes.
	mux.HandleFunc("/api/admin/me", logRequest(cors(s.requirePermission(admin.Access, s.handleAdminMe))))
	mux.HandleFunc("/api/admin/overview", logRequest(cors(s.requirePermission(admin.Access, s.handleAdminOverview))))
	mux.HandleFunc("/api/admin/reports", logRequest(cors(s.requirePermission(admin.ReportsRead, s.handleAdminReports))))
	mux.HandleFunc("/api/admin/reports/", logRequest(cors(s.requirePermission(admin.ReportsRead, s.handleAdminReport))))
	mux.HandleFunc("/api/admin/users", logRequest(cors(s.requirePermission(admin.UsersRead, s.handleAdminUsers))))
	mux.HandleFunc("/api/admin/users/", logRequest(cors(s.requirePermission(admin.UsersRead, s.handleAdminUser))))
	mux.HandleFunc("/api/admin/admins", logRequest(cors(s.requirePermission(admin.AdminsRead, s.handleAdminAdmins))))
	// Accepting an invitation is the one route here reachable by a non-staff account: it is how
	// a learner becomes staff. It still requires authentication, and it is bound to the invited
	// address.
	mux.HandleFunc("/api/admin/invitations/accept", logRequest(cors(s.writeLimiter.middleware(s.authMiddleware(s.handleAdminInvitationAccept)))))
	mux.HandleFunc("/api/admin/invitations", logRequest(cors(s.requirePermission(admin.AdminsRead, s.handleAdminInvitations))))
	mux.HandleFunc("/api/admin/invitations/", logRequest(cors(s.requirePermission(admin.AdminsRead, s.handleAdminInvitationRevoke))))
	mux.HandleFunc("/api/admin/audit", logRequest(cors(s.requirePermission(admin.AuditRead, s.handleAdminAudit))))
	mux.HandleFunc("/api/admin/content", logRequest(cors(s.requirePermission(admin.ContentRead, s.handleAdminContent))))
	mux.HandleFunc("/api/quiz/session", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleQuizSession)))))
	mux.HandleFunc("/api/quiz/answer", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleQuizAnswer)))))
	mux.HandleFunc("/api/quiz/skip", logRequest(cors(s.writeLimiter.middleware(s.optionalAuthMiddleware(s.handleQuizSkip)))))
	mux.HandleFunc("/api/health", logRequest(cors(getOnly(s.handleHealth))))
	mux.HandleFunc("/api/activity", logRequest(cors(s.authMiddleware(s.handleActivity))))
	mux.HandleFunc("/api/efficacy", logRequest(cors(s.authMiddleware(s.handleEfficacy))))
	mux.HandleFunc("/api/efficacy/all", logRequest(cors(s.shareLimiter.middleware(s.handleEfficacyAll))))
	mux.HandleFunc("/api/efficacy/trend", logRequest(cors(s.shareLimiter.middleware(s.handleEfficacyTrend))))
}

// GET /api/progress/{student_id}
func (s *Server) handleProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		// No ID enumeration: unauthenticated callers cannot read
		// arbitrary students. Guest IDs are the exception — the
		// unguessable `guest_` token in localStorage IS the credential
		// (same capability model as quiz/diag session UUIDs).
		// (Auth-disabled dev/test deployments keep the legacy fallback.)
		pathID := strings.TrimPrefix(r.URL.Path, "/api/progress/")
		if s.auth != nil && !strings.HasPrefix(pathID, "guest_") {
			writeError(w, "missing authorization", 401)
			return
		}
		studentID = pathID
		if studentID == "" {
			writeError(w, "student_id required", 400)
			return
		}
	}
	// Authed: always serve the caller's own progress — the path ID is
	// ignored so one student can never pull another's via enumeration
	// (e.g. IDs harvested from the public leaderboard).
	// mastery_pct is derived server-side, beside the state machine that enforces the
	// threshold. The client used to recompute it from `streak` and a bundled copy of the
	// corpus, which is a second source of truth that drifts from the server's whenever
	// the corpus is rebuilt.
	progress, err := s.eng.ProgressWithMasteryPct(studentID)
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
		// No ID enumeration: unauthenticated callers cannot read
		// arbitrary students. Guest IDs are the exception — the
		// unguessable `guest_` token in localStorage IS the credential
		// (same capability model as quiz/diag session UUIDs).
		// (Auth-disabled dev/test deployments keep the legacy fallback.)
		pathID := strings.TrimPrefix(r.URL.Path, "/api/scores/")
		if s.auth != nil && !strings.HasPrefix(pathID, "guest_") {
			writeError(w, "missing authorization", 401)
			return
		}
		studentID = pathID
		if studentID == "" {
			writeError(w, "student_id required", 400)
			return
		}
	}
	// Authed: always serve the caller's own scores — the path ID is ignored.
	scores, err := s.eng.GetScores(studentID)
	if err != nil {
		writeError(w, "failed to get scores", 500)
		return
	}
	writeJSON(w, scores)
}

// GET /api/config
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	cfg := map[string]interface{}{
		"auth_enabled": s.auth != nil,
		"providers":    auth.ConfiguredProviders(),
	}
	if v := os.Getenv("GOOGLE_CLIENT_ID"); v != "" {
		cfg["google_client_id"] = v
	} else if v := os.Getenv("GOOGLE_OAUTH_CLIENT_ID"); v != "" {
		cfg["google_client_id"] = v
	}
	writeJSON(w, cfg)
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

// GET /api/curriculum/domains — read-only curriculum picker: each domain with
// per-state counts. The server decides the states; the client renders them.
func (s *Server) handleCurriculumDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	domains, err := s.eng.CurriculumDomains(s.lessonStudentID(r))
	if err != nil {
		writeError(w, "failed to load curriculum", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"domains": domains})
}

// GET /api/curriculum/domains/{domainId} — one domain's concepts with state.
func (s *Server) handleCurriculumDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	domainID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/curriculum/domains/"), "/")
	if domainID == "" {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	curriculum, err := s.eng.CurriculumDomain(s.lessonStudentID(r), domainID)
	if err != nil {
		writeError(w, "domain not found", 404)
		return
	}
	writeJSON(w, map[string]interface{}{
		"domain": map[string]string{"id": curriculum.Domain},
		"topics": curriculum.Topics,
	})
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

// GET /api/leagues — weekly league standings with promotion/demotion.
func (s *Server) handleLeagues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	board, err := s.eng.GetLeagues()
	if err != nil {
		writeError(w, "failed to get leagues", 500)
		return
	}
	writeJSON(w, board)
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
// gradingTypeOf returns the DAG grading type for a concept ("" when unknown).
// Served alongside every question so clients can state the expected answer
// form up front instead of letting learners guess it.
func (s *Server) gradingTypeOf(cid string) string {
	if c := s.eng.GetDAG().Concept(cid); c != nil {
		return c.GradingType
	}
	return ""
}

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
	prog := s.eng.DiagnosticProgress(session)
	if question == nil {
		writeJSON(w, map[string]interface{}{"session_id": session.ID, "done": true, "progress": prog})
		return
	}
	writeJSON(w, map[string]interface{}{
		"session_id":   session.ID,
		"student_id":   studentID,
		"concept_id":   question.ConceptID,
		"concept_name": question.ConceptName,
		"question":     question.Question,
		"grading_type": s.gradingTypeOf(question.ConceptID),
		"progress":     prog,
	})
}

// POST /api/goal/diagnostic/answer
// Body: { "session_id": "...", "concept_id": "...", "answer": "...", "elapsed": 0.0, "dont_know": false }
// dont_know records an admitted unknown as negative evidence (grader
// bypassed, too-quick check exempt — an instant admit is honesty, not spam).
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
		DontKnow  bool    `json:"dont_know"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if !req.DontKnow && req.Elapsed < MinAnswerSeconds {
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
	// Sliding expiry: answering keeps a long diagnostic alive.
	s.mu.Lock()
	s.diagCreated[req.SessionID] = time.Now()
	s.mu.Unlock()

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
	expected := session.LastProblem.Answer
	expExplanation := session.LastProblem.Explanation
	questionText := session.LastProblem.Question
	// The difficulty this diagnostic question was generated at. The attempt row is the
	// mistakes transcript, but it is also read back by the mastery evidence window
	// (GetRecentAttemptsForConcept), so an attempt recorded here without difficulty is an
	// attempt that counts as "correct at unknown difficulty" against any concept the
	// learner later practises.
	servedDifficulty := session.LastProblem.Difficulty
	session.Unlock()
	if req.DontKnow {
		s.eng.SubmitDiagnosticDontKnow(session, req.ConceptID, req.Elapsed, timeThresh)
	} else {
		correct = s.eng.GradeAnswer(req.ConceptID, expected, req.Answer).Correct
		s.eng.SubmitDiagnosticAnswerTimed(session, req.ConceptID, correct, req.Elapsed, timeThresh)
	}
	explanation = expExplanation
	// Persist the attempt for the mistakes transcript (best-effort: a
	// recording failure must not fail the graded answer).
	session.Lock()
	diagStudentID := session.StudentID
	session.Unlock()
	if diagStudentID != "" {
		// Ephemeral diagnostic UUIDs have no sessions row; ensure one so
		// the attempts FK holds. Best-effort throughout: recording must
		// never fail the graded answer.
		if err := s.repo.EnsureSession(req.SessionID, diagStudentID); err != nil {
			log.Printf("handleGoalDiagnosticAnswer: ensure session failed for %s: %v", diagStudentID, err)
		} else if err := s.repo.RecordAttempt(storage.AttemptEntry{
			SessionID:      req.SessionID,
			StudentID:      diagStudentID,
			ConceptID:      req.ConceptID,
			Answer:         req.Answer,
			Expected:       expected,
			Correct:        correct,
			ElapsedSeconds: req.Elapsed,
			Timestamp:      time.Now(),
			Question:       questionText,
			Source:         "diagnostic",
			Explanation:    expExplanation,
			Difficulty:     servedDifficulty,
		}); err != nil {
			log.Printf("handleGoalDiagnosticAnswer: record attempt failed for %s: %v", diagStudentID, err)
		}
	}
	if s.eng.IsDiagnosticComplete(session) {
		report := s.eng.DiagnosticReport(session)
		prog := s.eng.DiagnosticProgress(session)
		writeAssessmentVerdict(w, map[string]interface{}{
			"done":     true,
			"correct":  correct,
			"feedback": explanation,
			"report":   report,
			"progress": prog,
		})
		return
	}
	nextProb, cid, err := s.eng.NextDiagnosticQuestion(session)
	if err != nil {
		writeError(w, "failed to get next question", 500)
		return
	}
	if cid == "" || nextProb == nil {
		// Cover exhausted (incl. supplemental): close with the report,
		// never serve an empty question with done:false.
		report := s.eng.DiagnosticReport(session)
		prog := s.eng.DiagnosticProgress(session)
		writeAssessmentVerdict(w, map[string]interface{}{
			"done":     true,
			"correct":  correct,
			"feedback": explanation,
			"report":   report,
			"progress": prog,
		})
		return
	}
	// Methodology gate, evaluated against the staged next question: the
	// offer covers the just-answered miss given what CAT serves next.
	retryAvailable := s.eng.DiagnosticRetryAvailableFor(session, req.ConceptID, cid)
	c := s.eng.GetDAG().Concept(cid)
	name := cid
	if c != nil {
		name = c.Label
	}
	prog := s.eng.DiagnosticProgress(session)
	writeAssessmentVerdict(w, map[string]interface{}{
		"done":            false,
		"correct":         correct,
		"feedback":        explanation,
		"concept_id":      cid,
		"concept_name":    name,
		"question":        nextProb.Question,
		"grading_type":    s.gradingTypeOf(cid),
		"progress":        prog,
		"retry_available": retryAvailable,
	})
}

// POST /api/goal/diagnostic/skip
// Body: { "session_id": "..." }
// Escape hatch for questions that can't be answered (bad content, repeated
// submit failures). Settles the pending concept with no evidence recorded —
// settling (not re-asking) is what moves CAT past a poisoned transition,
// since NextQuestion is deterministic on session state.
func (s *Server) handleGoalDiagnosticSkip(w http.ResponseWriter, r *http.Request) {
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
	// Sliding expiry: skipping keeps a long diagnostic alive.
	s.mu.Lock()
	s.diagCreated[req.SessionID] = time.Now()
	s.mu.Unlock()

	if s.eng.SettleDiagnosticCurrent(session) == "" {
		// Nothing pending — session is already done.
		writeJSON(w, map[string]interface{}{
			"done":     true,
			"correct":  false,
			"feedback": "diagnostic session already complete",
		})
		return
	}
	if s.eng.IsDiagnosticComplete(session) {
		report := s.eng.DiagnosticReport(session)
		prog := s.eng.DiagnosticProgress(session)
		writeJSON(w, map[string]interface{}{
			"done":     true,
			"correct":  false,
			"feedback": "Skipped — no evidence recorded.",
			"report":   report,
			"progress": prog,
		})
		return
	}
	nextProb, cid, err := s.eng.NextDiagnosticQuestion(session)
	if err != nil {
		writeError(w, "failed to get next question", 500)
		return
	}
	if cid == "" || nextProb == nil {
		report := s.eng.DiagnosticReport(session)
		prog := s.eng.DiagnosticProgress(session)
		writeJSON(w, map[string]interface{}{
			"done":     true,
			"correct":  false,
			"feedback": "Skipped — no evidence recorded.",
			"report":   report,
			"progress": prog,
		})
		return
	}
	c := s.eng.GetDAG().Concept(cid)
	name := cid
	if c != nil {
		name = c.Label
	}
	prog := s.eng.DiagnosticProgress(session)
	writeJSON(w, map[string]interface{}{
		"done":         false,
		"correct":      false,
		"feedback":     "Skipped — no evidence recorded.",
		"concept_id":   cid,
		"concept_name": name,
		"question":     nextProb.Question,
		"grading_type": s.gradingTypeOf(cid),
		"progress":     prog,
	})
}

// POST /api/goal/diagnostic/retry
// Body: { "session_id": "...", "concept_id": "..." }
// "Silly mistake" retry: voids the superseded question's miss and restores
// it as pending. The concept_id names the question being retried and must
// match the voidable attempt — a stale client that already moved on gets a
// 400 instead of voiding the wrong answer. Methodology-gated (engine):
// only surprising misses qualify, once per question. Never for quiz.
func (s *Server) handleGoalDiagnosticRetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
		ConceptID string `json:"concept_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.ConceptID == "" {
		writeError(w, "concept_id required", 400)
		return
	}
	s.mu.Lock()
	session := s.diagSessions[req.SessionID]
	s.mu.Unlock()
	if session == nil {
		writeError(w, "diagnostic session not found", 404)
		return
	}
	if authStudentID, _ := r.Context().Value(authStudentKey{}).(string); authStudentID != "" {
		session.Lock()
		owner := session.StudentID
		session.Unlock()
		if owner != "" && owner != authStudentID {
			writeError(w, "diagnostic session does not belong to authenticated user", 403)
			return
		}
	}
	// Sliding expiry: retrying keeps a long diagnostic alive.
	s.mu.Lock()
	s.diagCreated[req.SessionID] = time.Now()
	s.mu.Unlock()

	prob, cid, name, err := s.eng.RetryDiagnosticQuestion(session, req.ConceptID)
	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	prog := s.eng.DiagnosticProgress(session)
	writeJSON(w, map[string]interface{}{
		"done":         false,
		"concept_id":   cid,
		"concept_name": name,
		"question":     prob.Question,
		"grading_type": s.gradingTypeOf(cid),
		"progress":     prog,
	})
}

// GET /api/goal/diagnostic/resume?session_id=...
// MA parity: the Diagnostic doesn't have to be completed at once — a paused
// session resumes on its current question with coverage progress intact.
func (s *Server) handleGoalDiagnosticResume(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	sid := r.URL.Query().Get("session_id")
	if sid == "" {
		writeError(w, "session_id required", 400)
		return
	}
	s.mu.Lock()
	session := s.diagSessions[sid]
	s.mu.Unlock()
	if session == nil {
		writeError(w, "diagnostic session not found", 404)
		return
	}
	if authStudentID, _ := r.Context().Value(authStudentKey{}).(string); authStudentID != "" {
		session.Lock()
		owner := session.StudentID
		session.Unlock()
		if owner != "" && owner != authStudentID {
			writeError(w, "diagnostic session does not belong to authenticated user", 403)
			return
		}
	}
	if s.eng.IsDiagnosticComplete(session) {
		writeJSON(w, map[string]interface{}{
			"session_id": sid,
			"done":       true,
			"progress":   s.eng.DiagnosticProgress(session),
		})
		return
	}
	session.Lock()
	prob := session.LastProblem
	cid := session.LastConceptID
	name := session.LastConceptName
	session.Unlock()
	if prob == nil {
		nextProb, nc, err := s.eng.NextDiagnosticQuestion(session)
		if err != nil {
			writeError(w, "failed to get next question", 500)
			return
		}
		if nextProb == nil {
			writeJSON(w, map[string]interface{}{
				"session_id": sid,
				"done":       true,
				"progress":   s.eng.DiagnosticProgress(session),
			})
			return
		}
		prob = nextProb
		cid = nc
		if c := s.eng.GetDAG().Concept(nc); c != nil {
			name = c.Label
		} else {
			name = nc
		}
	}
	// Sliding expiry: an actively resumed session doesn't age out mid-test.
	s.mu.Lock()
	s.diagCreated[sid] = time.Now()
	s.mu.Unlock()
	writeJSON(w, map[string]interface{}{
		"session_id":   sid,
		"done":         false,
		"concept_id":   cid,
		"concept_name": name,
		"question":     prob.Question,
		"grading_type": s.gradingTypeOf(cid),
		"progress":     s.eng.DiagnosticProgress(session),
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

	// MA-parity report fields: placement + completion estimates from the
	// adaptive session (frontier, course recommendation, dates).
	report := s.eng.DiagnosticReport(session)

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
		"readiness":               readiness,
		"total_tested":            total,
		"correct_count":           correct,
		"weak_areas":              weakByDomain,
		"strong_areas":            strongByDomain,
		"frontier_label":          report.FrontierLabel,
		"frontier_idx":            report.FrontierIdx,
		"frontier_conditional":    report.FrontierConditional,
		"conditionally_completed": report.ConditionallyCompleted,
		"placement_course_id":     report.PlacementCourseID,
		"completion_estimates":    report.CompletionEstimates,
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
	// Only report concepts the student has actually attempted (has a progress row).
	// The 0.5 default for unseen concepts is for scheduling/difficulty only and
	// must not surface as a "struggle" for fresh accounts.
	progress, _ := s.repo.GetAllProgress(studentID)
	byDomain := make(map[string][]map[string]interface{})
	for _, c := range s.eng.GetDAG().Order() {
		if _, ok := progress[c.ID]; !ok {
			continue
		}
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
		// An empty object, not an empty array: the client's LessonsResSchema is
		// a record of domain → lessons, so returning [] rejects the whole
		// response and turns a missing corpus into a load error rather than an
		// empty catalog.
		writeJSON(w, map[string]interface{}{"lessons": map[string]interface{}{}})
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

	// Identity comes from the validated bearer token first, and only then from
	// an explicit student_id. The readiness and practice branches below already
	// resolve it this way; reading only the query param here meant a caller who
	// simply sent its token — which is every guest, since authedFetch attaches
	// the guest token — got the catalog back with no per-concept status at all.
	// Preferring the validated identity also means a mismatched param can never
	// select whose progress is returned.
	//
	// A foreign id is still ignored: the catalog stays public, but no one's
	// progress leaks and no one's study anchor can be poisoned through this
	// param.
	sid, _ := r.Context().Value(authStudentKey{}).(string)
	if sid == "" {
		sid = r.URL.Query().Get("student_id")
	}
	var progressMap map[string]map[string]interface{}
	if s.ownsStudentID(r, sid) {
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
	// Assets travel with the body so a client can address figures as data. `/study` renders
	// the body verbatim and needs nothing from them, but 133 of the corpus's figures have no
	// alt text and no caption, and there is no way to fix or even count those without a
	// list. An empty list is [] rather than null so clients need no null check.
	assets := l.Assets
	if assets == nil {
		assets = []lessons.Asset{}
	}
	writeJSON(w, map[string]interface{}{
		"title":  l.Title,
		"body":   l.Body,
		"source": l.Source,
		"assets": assets,
	})
}

func (s *Server) handleLessonConcept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 || (parts[1] != "practice" && parts[1] != "kp" && parts[1] != "readiness") {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	conceptID := parts[0]

	// GET /api/lessons/{id}/readiness — prerequisite detail and the server's
	// eligibility verdict. `ready/weak/missing` remains the soft banner the
	// /learn UI has always rendered; `eligible/state/reason/prerequisites` is the
	// canonical answer a learner's curriculum selection is gated on. A client
	// must ask, never decide (invariant 5).
	if parts[1] == "readiness" {
		sid := s.lessonStudentID(r)
		var rawProgress map[string]*storage.ConceptProgress
		if sid != "" {
			if p, err := s.eng.GetProgress(sid); err == nil {
				rawProgress = make(map[string]*storage.ConceptProgress)
				for cid, cp := range p {
					rawProgress[cid] = cp
				}
			}
		}
		prereqs := s.eng.LessonPrerequisites([]string{conceptID}, rawProgress)
		weak := []engine.PrereqInfo{}
		missing := []engine.PrereqInfo{}
		for _, p := range prereqs {
			if p.Status == string(mastery.StatusMastered) {
				continue
			}
			if p.Status == string(mastery.StatusUnseen) || p.Status == "" {
				missing = append(missing, p)
			} else {
				weak = append(weak, p)
			}
		}
		resp := map[string]interface{}{
			"concept_id": conceptID,
			"ready":      len(weak) == 0 && len(missing) == 0,
			"weak":       weak,
			"missing":    missing,
		}
		if el, err := s.eng.EligibilityFor(sid, conceptID); err == nil {
			resp["eligible"] = el.Eligible
			resp["state"] = el.State
			resp["reason"] = el.Reason
			resp["prerequisites"] = el.Prerequisites
		}
		writeJSON(w, resp)
		return
	}

	// GET /api/lessons/{id}/kp — knowledge-point shards with worked examples.
	if parts[1] == "kp" {
		ll := s.eng.GetLessonLoader()
		if ll == nil {
			http.Error(w, `{"error":"lessons unavailable"}`, 503)
			return
		}
		kps := ll.KPs(conceptID)
		type kpInfo struct {
			Label         string          `json:"label"`
			Section       string          `json:"section"`
			Subgoals      []string        `json:"subgoals"`
			WorkedExample string          `json:"worked_example"`
			Assets        []lessons.Asset `json:"assets"`
		}
		out := make([]kpInfo, 0, len(kps))
		for _, kp := range kps {
			// TeachingSlice never falls back to the lesson body. It used to: when a shard's
			// section was empty the server substituted l.Body, so 69 of 1,971 steps across 43
			// concepts served the whole reference article to a learner who asked for one
			// question — 33,784 words, and discrete.logic.propositions served its 2,724-word
			// article three times over. The reference panel keeps the article reachable (ADR-047);
			// this is the learning surface.
			we, _ := ll.TeachingSlice(conceptID, kp)
			assets := ll.KPAssets(conceptID, kp)
			if assets == nil {
				assets = []lessons.Asset{}
			}
			out = append(out, kpInfo{
				Label:         kp.Label,
				Section:       kp.Section,
				Subgoals:      kp.Subgoals,
				WorkedExample: we,
				Assets:        assets,
			})
		}
		writeJSON(w, map[string]interface{}{
			"concept_id": conceptID,
			"kps":        out,
			"diagram":    s.eng.DiagramFor(conceptID),
		})
		return
	}

	count := 5
	if c := r.URL.Query().Get("count"); c != "" {
		var n int
		if _, err := fmt.Sscanf(c, "%d", &n); err == nil && n > 0 && n <= 20 {
			count = n
		}
	}

	// Hard gate: a topic whose prerequisites are unmet cannot be started
	// through /learn, whether the learner arrived from a recommendation, the
	// curriculum browser, or a hand-typed URL. The readiness endpoint lets the
	// UI explain this in advance; this is the enforcement behind that promise
	// (invariant 3). Serving questions and anchoring answers for a locked
	// concept would let the frontend's choice stand in for knowledge.
	if el, err := s.eng.EligibilityFor(s.lessonStudentID(r), conceptID); err != nil {
		writeError(w, "unknown concept", 404)
		return
	} else if !el.Eligible {
		writeTopicLocked(w, el)
		return
	}
	// P1 fresh variants: optional seed for a different parameter set, and
	// exclude[] question texts the client already saw (miss => new variant,
	// never re-serve identical text). difficulty (0.3-1.0) steps the
	// generator's staged ladder; defaults to 0.5.
	var seed int64
	if s := r.URL.Query().Get("seed"); s != "" {
		fmt.Sscanf(s, "%d", &seed)
	}
	diff := 0.5
	if d := r.URL.Query().Get("difficulty"); d != "" {
		var v float64
		if _, err := fmt.Sscanf(d, "%f", &v); err == nil && v >= 0 && v <= 1 {
			diff = v
		}
	}
	excluded := make(map[string]bool)
	for _, ex := range r.URL.Query()["exclude"] {
		for _, part := range strings.Split(ex, "\n") {
			if part != "" {
				excluded[part] = true
			}
		}
	}

	type qInfo struct {
		Question    string `json:"question"`
		Answer      string `json:"answer"`
		Explanation string `json:"explanation"`
		Source      string `json:"source,omitempty"`
	}

	// Try DB first (filtered by exclude so retries vary)
	dbQs, err := s.eng.GetQuestions(conceptID, count+len(excluded)+5)
	if err == nil && len(dbQs) > 0 {
		filtered := dbQs[:0]
		for _, q := range dbQs {
			if !excluded[q.Question] {
				filtered = append(filtered, q)
			}
			if len(filtered) >= count {
				break
			}
		}
		if len(filtered) > 0 {
			dbQs = filtered
			questions := make([]qInfo, len(dbQs))
			for i, q := range dbQs {
				src := q.Source
				if src == "" {
					src = "curated"
				}
				questions[i] = qInfo{Question: q.Question, Answer: q.Answer, Explanation: q.Explanation, Source: src}
			}
			qa := make(map[string]engine.Anchor, len(dbQs))
			for _, q := range dbQs {
				// questions.difficulty is NOT NULL DEFAULT 0.5, so a curated row always has
				// a real value. It is the same field the generator path fills in, so both
				// branches of this serve hand the grader the same kind of evidence.
				d := q.Difficulty
				qa[q.Question] = engine.Anchor{Answer: q.Answer, Explanation: q.Explanation, Difficulty: &d}
			}
			if studentID, _ := r.Context().Value(authStudentKey{}).(string); s.ownsStudentID(r, studentID) {
				s.eng.SetStudyAnchorBatch(studentID, conceptID, qa)
			} else if sid := r.URL.Query().Get("student_id"); s.ownsStudentID(r, sid) {
				s.eng.SetStudyAnchorBatch(sid, conceptID, qa)
			}
			writeJSON(w, map[string]interface{}{"questions": questions, "concept_id": conceptID})
			return
		}
	}

	// Fallback to generator (seeded + exclude-filtered for fresh variants)
	reg := s.eng.GetGeneratorRegistry()
	if reg == nil {
		writeJSON(w, map[string]interface{}{"questions": []interface{}{}})
		return
	}
	problems, err := reg.BatchGenerateContext(conceptID, count+len(excluded)+5, generator.GeneratorContext{Difficulty: diff, Seed: seed})
	if err != nil {
		writeError(w, err.Error(), 404)
		return
	}
	fresh := problems[:0]
	for _, p := range problems {
		if !excluded[p.Question] {
			fresh = append(fresh, p)
		}
		if len(fresh) >= count {
			break
		}
	}
	// If exclusion emptied the set, serve the batch anyway rather than returning
	// `questions: []` with a 200. A learner whose session has already seen every
	// variant a concept has must still be able to answer something.
	//
	// This is not hypothetical: `geo.basic.points_lines` has exactly four distinct
	// question texts, so after four answers the filtered set was empty and the Learn
	// feed read that as success, stopped appending, and left the last verdict as the end
	// of the page with no next step and no message. Repeats are a far better outcome than
	// a stranded learner, and the variety shortfall is measured separately in
	// `internal/generator/all/bank_test.go` — 378 of 657 concepts sit under a 12-question
	// reference, which is a content backlog, not something an empty array should be
	// answering.
	if len(fresh) == 0 && len(problems) > 0 {
		fresh = problems
	}
	problems = fresh
	if len(problems) == 0 {
		// The generator produced nothing at all, which is a different failure from "the
		// learner has seen them all" and must not look like it. An empty 200 reads as
		// success to every caller.
		http.Error(w, `{"error":"no questions available for this concept"}`, 503)
		return
	}
	// H1b: anchor the whole served set (not just the first question) so each
	// question is graded against its own answer and answered with its own
	// explanation. The two travel together: both come from one generator call.
	qa := make(map[string]engine.Anchor, len(problems))
	for _, p := range problems {
		qa[p.Question] = engine.Anchor{Answer: p.Answer, Explanation: p.Explanation, Difficulty: p.Difficulty}
	}
	if studentID, _ := r.Context().Value(authStudentKey{}).(string); s.ownsStudentID(r, studentID) {
		s.eng.SetStudyAnchorBatch(studentID, conceptID, qa)
	} else if sid := r.URL.Query().Get("student_id"); s.ownsStudentID(r, sid) {
		s.eng.SetStudyAnchorBatch(sid, conceptID, qa)
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
	if req.Goal < xp.MinDailyGoal || req.Goal > xp.MaxDailyGoal {
		writeError(w, fmt.Sprintf("goal must be between %d and %d", xp.MinDailyGoal, xp.MaxDailyGoal), 400)
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

// POST /api/avatar (multipart "avatar") uploads a custom profile photo;
// DELETE /api/avatar removes it. Photos are magic-byte sniffed (PNG/JPEG/GIF
// /WebP only) and capped at 512KB. A settings flag (avatar_custom) records
// the explicit choice so "remove" falls back to the DiceBear pick.
func (s *Server) handleAvatar(w http.ResponseWriter, r *http.Request) {
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	switch r.Method {
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 600*1024)
		if err := r.ParseMultipartForm(600 * 1024); err != nil {
			writeError(w, "avatar too large (512KB max)", 400)
			return
		}
		f, _, err := r.FormFile("avatar")
		if err != nil {
			writeError(w, "missing avatar file", 400)
			return
		}
		defer f.Close()
		data, err := io.ReadAll(io.LimitReader(f, 512*1024+1))
		if err != nil {
			writeError(w, "failed to read avatar", 400)
			return
		}
		if len(data) == 0 || len(data) > 512*1024 {
			writeError(w, "avatar too large (512KB max)", 400)
			return
		}
		ct := http.DetectContentType(data)
		switch ct {
		case "image/png", "image/jpeg", "image/gif", "image/webp":
		default:
			writeError(w, "avatar must be PNG, JPEG, GIF, or WebP", 400)
			return
		}
		if err := s.repo.SetAvatarImage(studentID, ct, data); err != nil {
			writeError(w, "failed to save avatar", 500)
			return
		}
		if err := s.setSettingsFlag(studentID, "avatar_custom", true); err != nil {
			writeError(w, "failed to save avatar", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"ok": true, "content_type": ct, "bytes": len(data)})
	case http.MethodDelete:
		if err := s.repo.ClearAvatarImage(studentID); err != nil {
			writeError(w, "failed to remove avatar", 500)
			return
		}
		if err := s.setSettingsFlag(studentID, "avatar_custom", false); err != nil {
			writeError(w, "failed to remove avatar", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"ok": true})
	default:
		http.Error(w, `{"error":"method not allowed"}`, 405)
	}
}

// GET /api/avatar/me serves the student's custom photo (404 when none).
func (s *Server) handleAvatarMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	ct, data, found, err := s.repo.GetAvatarImage(studentID)
	if err != nil {
		writeError(w, "failed to load avatar", 500)
		return
	}
	if !found {
		writeError(w, "no custom avatar", 404)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "private, max-age=86400")
	w.Write(data)
}

// GET /api/avatar/{student_id} serves another student's custom photo for
// leaderboard avatars. Public (the board is visible logged-out); 404 unless
// the student uploaded a photo. Version-busting is client-side (?v=).
func (s *Server) handleAvatarPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID := strings.TrimPrefix(r.URL.Path, "/api/avatar/")
	if studentID == "" || strings.Contains(studentID, "/") {
		writeError(w, "student id is required", 400)
		return
	}
	ct, data, found, err := s.repo.GetAvatarImage(studentID)
	if err != nil {
		writeError(w, "failed to load avatar", 500)
		return
	}
	if !found {
		writeError(w, "no custom avatar", 404)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(data)
}

// setSettingsFlag merges one key into the student's opaque settings JSON.
func (s *Server) setSettingsFlag(studentID, key string, value bool) error {
	raw, err := s.repo.GetSettings(studentID)
	if err != nil {
		return err
	}
	m := map[string]interface{}{}
	if strings.TrimSpace(raw) != "" {
		_ = json.Unmarshal([]byte(raw), &m)
		if m == nil {
			m = map[string]interface{}{}
		}
	}
	m[key] = value
	out, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return s.repo.UpdateSettings(studentID, string(out))
}

// GET /api/next — the learner's recommended task and their alternatives.
//
// The whole point of this endpoint is that the decision happens once. `/learn` and `/profile`
// used to recompute it independently in the browser from five API calls and a bundled copy of
// the corpus, while internal/scheduler held the eligibility rules, the decay rule and the
// priority weights where neither surface could reach them. The two had already drifted: the
// client kept its own copy of the XP award math and still halved reviews, a rule ADR-020
// removed from the server.
//
// `exclude` lets a caller drop a concept it has just finished, so it is never offered back as
// the next thing to do.
func (s *Server) handleNext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID != "" && paused(studentID, s.repo) {
		// A paused learner's reviews are hidden, so a recommendation that sends them to
		// /review would contradict the pause.
		writeJSON(w, s.eng.RecommendPaused(studentID))
		return
	}
	var exclude []string
	for _, e := range r.URL.Query()["exclude"] {
		for _, id := range strings.Split(e, ",") {
			if id = strings.TrimSpace(id); id != "" {
				exclude = append(exclude, id)
			}
		}
	}
	resp, err := s.eng.RecommendNext(studentID, exclude)
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"primary":      nil,
			"alternatives": []interface{}{},
			"generatedAt":  time.Now().UTC(),
		})
		return
	}
	writeJSON(w, resp)
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
	// S3-3: paused students see no due reviews.
	if paused(studentID, s.repo) {
		writeJSON(w, map[string]interface{}{"count": 0})
		return
	}
	// The count lives on the engine now: the recommender needs the same number, and having it
	// written twice is how the two could disagree about what "due" means.
	count, err := s.eng.ReviewsDue(studentID)
	if err != nil {
		writeJSON(w, map[string]interface{}{"count": 0})
		return
	}
	writeJSON(w, map[string]interface{}{"count": count})
}

// POST /api/share — enable (body {enabled:true|false}) the read-only share link.
func (s *Server) handleShareToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if !req.Enabled {
		if err := s.eng.DisableShare(studentID); err != nil {
			writeError(w, "failed to disable share", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"enabled": false, "token": ""})
		return
	}
	token, err := s.eng.EnableShare(studentID)
	if err != nil {
		writeError(w, "failed to enable share", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"enabled": true, "token": token})
}

// GET /api/share/{token} — public read-only oversight view (parent/teacher).
func (s *Server) handleShareReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	token := strings.TrimPrefix(r.URL.Path, "/api/share/")
	token = strings.TrimSuffix(token, "/")
	if token == "" {
		writeError(w, "token required", 400)
		return
	}
	report, err := s.eng.GetShareReport(token)
	if err != nil {
		writeError(w, "invalid or disabled share link", 404)
		return
	}
	writeJSON(w, report)
}

// paused reports whether settings.pause_until is set to a future date.
func paused(studentID string, repo storage.Repository) bool {
	raw, err := repo.GetSettings(studentID)
	if err != nil || raw == "" || raw == "{}" {
		return false
	}
	var cfg struct {
		PauseUntil string `json:"pause_until"`
	}
	if json.Unmarshal([]byte(raw), &cfg) != nil || cfg.PauseUntil == "" {
		return false
	}
	t, err := time.Parse("2006-01-02", cfg.PauseUntil)
	if err != nil {
		return false
	}
	return t.After(time.Now().UTC())
}

// startSessionRes and the answer shapes below were named for the session flow
// that used to live here. Those routes are gone; the review handlers still use
// them, so the shapes stay and the comments now say so.
type startSessionRes struct {
	StudentID string           `json:"student_id"`
	SessionID string           `json:"session_id"`
	Question  *engine.Question `json:"question"`
}

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

// GET /api/courses/{id} | POST /api/courses/{id}/diagnostic
func (s *Server) handleCourseRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	if strings.HasSuffix(path, "/diagnostic") {
		s.handleCourseDiagnostic(w, r)
		return
	}
	s.handleCourseDetail(w, r)
}

// GET /api/destinations — destination bundles with per-student progress.
func (s *Server) handleDestinations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	list, err := s.eng.Destinations(studentID)
	if err != nil {
		writeError(w, "failed to load destinations", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"destinations": list})
}

// GET /api/destinations/{id}/estimate — two-track workload estimate.
// Query: daily_goal (default stored goal), deadline_days (-1 none),
// rest_days (0-6), diagnostic_min. Pure estimator: never mutates mastery,
// quiz eligibility, or XP.
func (s *Server) handleDestinationEstimate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	rest := strings.TrimPrefix(r.URL.Path, "/api/destinations/")
	rest = strings.TrimSuffix(rest, "/")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[1] != "estimate" || parts[0] == "" {
		http.Error(w, `{"error":"not found"}`, 404)
		return
	}
	q := r.URL.Query()
	dailyGoal := 0
	if v, err := strconv.Atoi(q.Get("daily_goal")); err == nil && v > 0 {
		dailyGoal = v
	}
	if dailyGoal == 0 {
		if st, err := s.repo.GetStudent(studentID); err == nil && st != nil && st.DailyXPGoal > 0 {
			dailyGoal = st.DailyXPGoal
		} else {
			// Must match xp.DefaultDailyGoal, engine/plan.go's fallback and the client's
			// lib/plan.ts DEFAULT_DAILY_GOAL. This one said 30 while its sibling route
			// handlePlanCurrent said 10, so the same learner got a 3x different pace from
			// two endpoints.
			dailyGoal = xp.DefaultDailyGoal
		}
	}
	deadlineDays := -1
	if v, err := strconv.Atoi(q.Get("deadline_days")); err == nil && v > 0 {
		deadlineDays = v
	}
	restDays := 0
	if v, err := strconv.Atoi(q.Get("rest_days")); err == nil && v >= 0 {
		restDays = v
	}
	diagMin := 0.0
	if v, err := strconv.ParseFloat(q.Get("diagnostic_min"), 64); err == nil && v > 0 {
		diagMin = v
	}
	resp, err := s.eng.DestinationEstimate(studentID, parts[0], dailyGoal, deadlineDays, restDays, diagMin)
	if err != nil {
		writeError(w, err.Error(), 404)
		return
	}
	writeJSON(w, resp)
}

// POST /api/plans — snapshot the current estimate as the ahead/behind baseline.
func (s *Server) handlePlanSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	var req struct {
		Destination  string `json:"destination"`
		DailyGoal    int    `json:"daily_goal"`
		DeadlineDays int    `json:"deadline_days"`
		RestDays     int    `json:"rest_days"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.Destination == "" {
		writeError(w, "destination required", 400)
		return
	}
	snap, err := s.eng.SavePlan(studentID, req.Destination, req.DailyGoal, req.DeadlineDays, req.RestDays)
	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, snap)
}

// GET /api/plans/current — saved baseline plus fresh ahead/behind delta.
// Prefs (destination/deadline/pace) come from the kept prefs record; the
// daily rate reads live so goal changes apply without re-saving.
func (s *Server) handlePlanCurrent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	prefs, prefsOK := s.eng.GetPlanPrefs(studentID)
	snap, snapOK := s.eng.GetPlan(studentID)
	if !prefsOK && !snapOK {
		writeJSON(w, map[string]interface{}{"plan": nil})
		return
	}
	dest := ""
	deadline, rest := -1, 0
	if prefsOK {
		dest, deadline, rest = prefs.Destination, prefs.DeadlineDays, prefs.RestDays
	}
	if dest == "" && snapOK {
		dest = snap.Destination
	}
	goal := xp.DefaultDailyGoal
	if st, err := s.repo.GetStudent(studentID); err == nil && st != nil && st.DailyXPGoal > 0 {
		goal = st.DailyXPGoal
	}
	resp, err := s.eng.DestinationEstimate(studentID, dest, goal, deadline, rest, 0)
	if err != nil {
		writeJSON(w, map[string]interface{}{"plan": snap, "prefs": prefs})
		return
	}
	writeJSON(w, map[string]interface{}{"plan": snap, "prefs": prefs, "delta_days": resp.PlanDelta, "fresh_start": resp.FreshStart, "xp_remaining": resp.Estimate.XPRemaining})
}

// POST /api/account/reset — wipe learning record after typed confirmation.
// Body: { "phrase": "reset my progress" }. Retry-safe; see engine docs.
func (s *Server) handleAccountReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	var req struct {
		Phrase string `json:"phrase"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if err := s.eng.ResetAccountProgress(studentID, req.Phrase); err != nil {
		if errors.Is(err, engine.ErrResetPhraseMismatch) {
			writeError(w, "confirmation phrase does not match", 400)
			return
		}
		// An infrastructure fault must not read as a typo: the learner would
		// retype the same correct phrase forever.
		log.Printf("account reset failed for %s: %v", studentID, err)
		writeError(w, "reset failed and nothing was changed — please try again", 500)
		return
	}
	// Evict live diagnostic/quiz sessions so a tab left open cannot rewrite
	// the rows we just deleted.
	s.dropStudentLiveSessions(studentID)
	writeJSON(w, map[string]interface{}{"reset": true})
}

// DELETE /api/account — delete the student and every owned row after typed
// confirmation, plus the current password for password accounts (OAuth-only
// and guest rows need the phrase alone). Irreversible.
func (s *Server) handleAccountDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	var req struct {
		Phrase   string `json:"phrase"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	st, err := s.repo.GetStudent(studentID)
	if err != nil || st == nil {
		writeError(w, "student not found", 404)
		return
	}
	if st.PasswordHash != "" {
		if strings.TrimSpace(req.Password) == "" {
			writeError(w, "current password is required", 401)
			return
		}
		if err := s.auth.CheckPassword(studentID, req.Password); err != nil {
			writeError(w, "current password is incorrect", 401)
			return
		}
	}
	if err := s.eng.DeleteAccount(studentID, req.Phrase); err != nil {
		if errors.Is(err, engine.ErrDeletePhraseMismatch) {
			writeError(w, "confirmation phrase does not match", 400)
			return
		}
		log.Printf("account delete failed for %s: %v", studentID, err)
		writeError(w, "deletion failed and nothing was changed — please try again", 500)
		return
	}
	// Evict live sessions too. On Postgres these FKs are absent, so without
	// this an in-flight answer re-creates attempts rows for an identity that
	// no longer exists.
	s.dropStudentLiveSessions(studentID)
	log.Printf("auth: account deleted for student %s", studentID)
	writeJSON(w, map[string]interface{}{"deleted": true})
}

// GET /api/courses/{id} — one course with the student's progress.
func (s *Server) handleCourseDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	courseID := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	courseID = strings.TrimSuffix(courseID, "/")
	if courseID == "" {
		writeError(w, "missing course id", 400)
		return
	}
	catalog, err := s.eng.CourseCatalog(studentID)
	if err != nil {
		writeError(w, "failed to load courses", 500)
		return
	}
	for _, cs := range catalog {
		if cs.ID == courseID {
			writeJSON(w, cs)
			return
		}
	}
	writeError(w, "course not found", 404)
}

// GET /api/attempts — mistakes transcript source, always self-scoped.
// Query: source=diagnostic|quiz|practice|review, concept_id=...,
// incorrect_only=1, limit (default 100, max 500), offset (default 0).
// Newest first. Pre-migration rows carry empty question/source.
func (s *Server) handleAttempts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	q := r.URL.Query()
	sourceFilter := q.Get("source")
	conceptFilter := q.Get("concept_id")
	incorrectOnly := q.Get("incorrect_only") == "1" || strings.EqualFold(q.Get("incorrect_only"), "true")
	limit := 100
	if v, err := strconv.Atoi(q.Get("limit")); err == nil {
		limit = v
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 500 {
		limit = 500
	}
	offset := 0
	if v, err := strconv.Atoi(q.Get("offset")); err == nil && v > 0 {
		offset = v
	}
	all, err := s.repo.GetAttemptsForStudent(studentID)
	if err != nil {
		writeError(w, "failed to load attempts", 500)
		return
	}
	type attemptRow struct {
		storage.AttemptEntry
		ConceptName string `json:"concept_name"`
	}
	filtered := make([]attemptRow, 0, len(all))
	for _, a := range all {
		if sourceFilter != "" && a.Source != sourceFilter {
			continue
		}
		if conceptFilter != "" && a.ConceptID != conceptFilter {
			continue
		}
		if incorrectOnly && a.Correct {
			continue
		}
		name := a.ConceptID
		if c := s.eng.GetDAG().Concept(a.ConceptID); c != nil {
			name = c.Label
		}
		filtered = append(filtered, attemptRow{AttemptEntry: a, ConceptName: name})
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})
	total := len(filtered)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	writeJSON(w, map[string]interface{}{"attempts": filtered[offset:end], "total": total})
}

// GET /api/transcript — accreditation-track completion overview
// (?format=csv returns a downloadable transcript).
func (s *Server) handleTranscript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	catalog, err := s.eng.CourseCatalog(studentID)
	if err != nil {
		writeError(w, "failed to load transcript", 500)
		return
	}
	if r.URL.Query().Get("format") == "csv" {
		writeTranscriptCSV(w, catalog)
		return
	}
	writeJSON(w, map[string]interface{}{"courses": catalog})
}

func writeTranscriptCSV(w http.ResponseWriter, catalog []engine.CourseStatus) {
	var b strings.Builder
	b.WriteString("course,grade,total,mastered,pct,days_remaining,status\n")
	for _, cs := range catalog {
		p := cs.Progress
		if p == nil {
			b.WriteString(fmt.Sprintf("%s,%s,0,0,0,0,not_started\n", csvCell(cs.Name), csvCell(cs.Grade)))
			continue
		}
		status := "not_started"
		if p.Mastered > 0 && p.Mastered < p.Total {
			status = "in_progress"
		} else if p.Total > 0 && p.Mastered == p.Total {
			status = "complete"
		}
		b.WriteString(fmt.Sprintf("%s,%s,%d,%d,%.0f,%d,%s\n",
			csvCell(cs.Name), csvCell(cs.Grade), p.Total, p.Mastered, p.Pct*100, p.DaysRemaining, status))
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="mathua-transcript.csv"`)
	w.Write([]byte(b.String()))
}

func csvCell(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
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

// POST /api/study/answer — Learn seam: LessonQuiz → SubmitAnswer (CONTEXT.md Seam).
// The endpoint keeps its name; it is Learn's answer route and has never been Study's page.
// Body: { concept_id, answer, elapsed, question, student_id? }
// student_id is used only for a guest (mathua_guest_id).
//
// The request carries no expected answer and needs none: the answer and the
// explanation both come from the server-side anchor written when the question
// was served. An `expected` a client does send is ignored — there is nothing
// left for it to grade against.
func (s *Server) handleStudyAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		ConceptID string  `json:"concept_id"`
		Answer    string  `json:"answer"`
		Elapsed   float64 `json:"elapsed"`
		StudentID string  `json:"student_id"`
		Question  string  `json:"question"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.ConceptID == "" {
		writeError(w, "concept_id required", 400)
		return
	}
	if len(req.Answer) > 4096 {
		writeError(w, "answer too long", 400)
		return
	}
	if len(req.Question) > 4096 {
		writeError(w, "question too long", 400)
		return
	}
	authStudentID, _ := r.Context().Value(authStudentKey{}).(string)
	if authStudentID != "" && req.StudentID != "" && req.StudentID != authStudentID {
		writeError(w, "student_id does not match authenticated user", 403)
		return
	}
	studentID := authStudentID
	if studentID == "" {
		studentID = req.StudentID
	}
	if studentID != "" && !s.ownsStudentID(r, studentID) {
		writeError(w, "student_id not permitted", 403)
		return
	}
	// An unidentifiable caller gets nothing: there is no answer to grade
	// against and nowhere to record a result, and this used to be the one path
	// that graded a study answer against the client's own expected. Every
	// client sends a guest id or a token, so this is a malformed request.
	if studentID == "" {
		writeError(w, "student identity required", 400)
		return
	}
	if req.Elapsed < 0 {
		req.Elapsed = 0
	}
	if req.Elapsed > 600 {
		req.Elapsed = 600
	}
	if req.Elapsed < MinAnswerSeconds {
		writeError(w, "answer submitted too quickly", 400)
		return
	}
	// Defense in depth behind the practice gate: an answer for a concept whose
	// prerequisites are unmet is refused rather than recorded, so no path can
	// turn a learner's choice into evidence for a topic they were not ready for.
	if el, err := s.eng.EligibilityFor(studentID, req.ConceptID); err == nil && !el.Eligible {
		writeTopicLocked(w, el)
		return
	}
	res, err := s.eng.SubmitStudyAnswer(studentID, req.ConceptID, req.Answer, req.Elapsed, req.Question)
	if err != nil {
		if errors.Is(err, engine.ErrUnknownConcept) {
			writeError(w, "unknown concept", 404)
			return
		}
		if errors.Is(err, engine.ErrNoStudyAnchor) {
			// Not a wrong answer: nothing was recorded, no streak or weakness
			// moved. 409 tells the client to re-serve, which is the only way
			// this becomes gradeable again.
			writeError(w, "no server-side record of this question; re-serve it", 409)
			return
		}
		writeError(w, "failed to submit study answer", 500)
		return
	}
	writeJSON(w, map[string]interface{}{
		"correct":         res.Correct,
		"feedback":        res.Feedback,
		"diagnosis":       res.Diagnosis,
		"explanation":     res.Explanation,
		"new_status":      res.NewStatus,
		"streak":          res.Streak,
		"required_streak": res.RequiredStreak,
		"xp":              res.XP,
		"expected_answer": res.ExpectedAnswer,
		"halted":          res.Halted,
	})
}

type authStudentKey struct{}

// authRoleKey carries the caller's role from requirePermission into the handler, so a handler
// that needs to compare itself against a target (may it manage this account? is this a
// self-change?) does not have to re-read the role it was already authorized on.
type authRoleKey struct{}

// ownsStudentID reports whether the caller may act as sid. The validated
// Bearer identity always wins; unauthenticated callers may only act as
// client-generated guest IDs (unguessable `guest_` tokens kept in
// localStorage — knowledge of the token IS the credential). Auth-disabled
// deployments (dev/test) allow all, preserving legacy behavior.
func (s *Server) ownsStudentID(r *http.Request, sid string) bool {
	if sid == "" {
		return false
	}
	if authID, _ := r.Context().Value(authStudentKey{}).(string); authID != "" {
		return authID == sid
	}
	if s.auth != nil {
		return strings.HasPrefix(sid, "guest_")
	}
	return true
}

// lessonStudentID resolves whose learning state a lesson request describes.
// The validated bearer identity wins; otherwise an owned guest id from the query
// is used. Empty means anonymous, whose state is all-unseen.
func (s *Server) lessonStudentID(r *http.Request) string {
	if sid, _ := r.Context().Value(authStudentKey{}).(string); sid != "" {
		return sid
	}
	if sid := r.URL.Query().Get("student_id"); s.ownsStudentID(r, sid) {
		return sid
	}
	return ""
}

// writeTopicLocked answers a request to start a concept whose prerequisites are
// unmet, in the project's error envelope. The code is stable for the client to
// branch on; the prerequisite ids let the UI name what to learn first.
func writeTopicLocked(w http.ResponseWriter, el *engine.ConceptEligibility) {
	ids := make([]string, 0, len(el.Prerequisites))
	for _, p := range el.Prerequisites {
		if !p.Met {
			ids = append(ids, p.ConceptID)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "TOPIC_LOCKED",
			"message": "This topic has unmet prerequisites.",
			"details": map[string]interface{}{"prerequisiteConceptIds": ids},
		},
	})
}

// getOnly rejects non-GET requests for read-only endpoints (health/config).
func getOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

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
	authStudentID, _ := r.Context().Value(authStudentKey{}).(string)
	if authStudentID != "" && req.StudentID != "" && req.StudentID != authStudentID {
		writeError(w, "student_id does not match authenticated user", 403)
		return
	}
	studentID := authStudentID
	if studentID == "" {
		studentID = req.StudentID
	}
	if studentID != "" && !s.ownsStudentID(r, studentID) {
		writeError(w, "student_id not permitted", 403)
		return
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
		"grading_type": s.gradingTypeOf(cid),
		"done":         false,
		// Batch 1: timed closed-book contract — per-question limit
		// (accommodated), total count, no-lesson closed book.
		"closed_book":        true,
		"time_limit_seconds": s.eng.TimeLimitFor(studentID, cid),
		"questions_total":    len(sess.Order),
	})
}

// POST /api/quiz/answer — own grading path via SubmitStudyAnswer with TaskQuiz 20
// Body: { session_id, concept_id, answer, elapsed, student_id?, dont_know? }
// dont_know records an admitted unknown as a miss (remedial queued, no XP,
// no rushing penalty); the grader is bypassed.
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
		DontKnow  bool    `json:"dont_know"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.SessionID == "" || req.ConceptID == "" {
		writeError(w, "session_id and concept_id required", 400)
		return
	}
	if len(req.Answer) > 4096 {
		writeError(w, "answer too long", 400)
		return
	}
	if !req.DontKnow && req.Elapsed < MinAnswerSeconds {
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
	authStudentID, _ := r.Context().Value(authStudentKey{}).(string)
	if authStudentID != "" && req.StudentID != "" && req.StudentID != authStudentID {
		writeError(w, "student_id does not match authenticated user", 403)
		return
	}
	studentID := authStudentID
	if studentID == "" {
		studentID = req.StudentID
		if studentID == "" {
			studentID = sess.StudentID
		}
	}
	if studentID != "" && !s.ownsStudentID(r, studentID) {
		writeError(w, "student_id not permitted", 403)
		return
	}
	if authStudentID != "" && sess.StudentID != "" && sess.StudentID != authStudentID {
		writeError(w, "quiz session does not belong to authenticated user", 403)
		return
	}
	// Grade via the canonical engine grader (custom GradedGenerator → router),
	// the same dispatch SubmitQuizAnswer uses, so the response, the recorded
	// answer, and the XP/mastery update can never disagree.
	expected := sess.LastProblem.Answer
	gr := s.eng.GradeAnswer(req.ConceptID, expected, req.Answer)
	// The session already holds the explanation the generator produced for the
	// question that was served. Send it on every verdict — a correct answer
	// still needs the reasoning that made it right — and fall back to the
	// grader's status token only if it is somehow empty. Same contract as
	// handleGoalDiagnosticAnswer, which was the one path already doing this.
	explanation := sess.LastProblem.Explanation
	feedback := gr.Feedback
	if explanation == "" {
		explanation = feedback
	}
	if req.DontKnow {
		// Admitted unknown: forced miss, teaching content as feedback.
		gr = grader.Result{Correct: false}
		feedback = explanation
		if feedback == "" {
			feedback = "Not quite."
		}
	}
	// Record via quiz engine (advance index)
	qEng := quiz.NewEngine(s.eng.GetDAG(), s.eng.GetGeneratorRegistry())
	qEng.RecordAnswer(sess, req.ConceptID, gr.Correct)
	// Single-path quiz grading: the engine owns TaskQuiz XP atomically —
	// one progress update, one AddXP, DB and response agree by construction.
	xp := 0
	var newStatus string
	var remedial []string
	if studentID != "" {
		var res *engine.AnswerResult
		var err error
		// The difficulty comes off the problem that was served, which quiz.NextQuestion
		// stamped from the GeneratorContext it generated at. Reading it here rather than
		// recomputing it is what makes a quiz attempt count for mastery exactly as much as
		// the equivalent study attempt.
		servedDifficulty := sess.LastProblem.Difficulty
		if req.DontKnow {
			res, err = s.eng.SubmitQuizDontKnow(studentID, req.ConceptID, expected, req.Elapsed, sess.LastProblem.Question, servedDifficulty)
		} else {
			res, err = s.eng.SubmitQuizAnswer(studentID, req.ConceptID, req.Answer, expected, req.Elapsed, sess.LastProblem.Question, servedDifficulty)
		}
		if err != nil {
			if errors.Is(err, engine.ErrUnknownConcept) {
				writeError(w, "unknown concept", 404)
				return
			}
			// Persistence failure must not kill quiz progression: the
			// grade above stands, XP just isn't awarded.
			log.Printf("handleQuizAnswer: SubmitQuizAnswer failed for %s/%s: %v", studentID, req.ConceptID, err)
			xp = 0
		} else if res != nil {
			xp = res.XP
			newStatus = string(res.NewStatus)
			// Batch 1: immediate remedial on quiz miss.
			remedial = res.Remedial
		}
	}
	if qEng.IsComplete(sess) {
		if studentID != "" {
			// Batch 1: completing a quiz resets the 150 XP gate baseline.
			if err := s.eng.RecordQuizCompletion(studentID); err != nil {
				log.Printf("handleQuizAnswer: RecordQuizCompletion failed for %s: %v", studentID, err)
			}
		}
		s.mu.Lock()
		delete(s.quizSessions, req.SessionID)
		delete(s.quizCreated, req.SessionID)
		s.mu.Unlock()
		// The `diagnosis` and `explanation` keys are still passed so the helper is the single
		// place that decides they do not cross, rather than the call site remembering to omit
		// them. A verdict built anywhere else leaks by default; through the helper it does not.
		writeAssessmentVerdict(w, map[string]interface{}{"done": true, "correct": gr.Correct, "feedback": feedback, "diagnosis": gr.Diagnosis, "explanation": explanation, "xp": xp, "new_status": newStatus, "remedial": remedial, "retake_available": true})
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
	writeAssessmentVerdict(w, map[string]interface{}{
		"done":               false,
		"correct":            gr.Correct,
		"feedback":           feedback,
		"diagnosis":          gr.Diagnosis,
		"explanation":        explanation,
		"xp":                 xp,
		"new_status":         newStatus,
		"remedial":           remedial,
		"concept_id":         cid,
		"concept_name":       name2,
		"question":           prob.Question,
		"grading_type":       s.gradingTypeOf(cid),
		"closed_book":        true,
		"time_limit_seconds": s.eng.TimeLimitFor(studentID, cid),
		"questions_total":    len(sess.Order),
	})
}

// POST /api/quiz/skip — escape hatch for unanswerable quiz questions.
// Body: { session_id, concept_id?, student_id? }
// Advances past the pending question with no grade, no XP, no remedial.
// The 150 XP gate baseline resets only if at least one answer was recorded,
// so a fully-skipped quiz can't dodge the gate.
func (s *Server) handleQuizSkip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
		ConceptID string `json:"concept_id"`
		StudentID string `json:"student_id"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if req.SessionID == "" {
		writeError(w, "session_id required", 400)
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
	authStudentID, _ := r.Context().Value(authStudentKey{}).(string)
	if authStudentID != "" && req.StudentID != "" && req.StudentID != authStudentID {
		writeError(w, "student_id does not match authenticated user", 403)
		return
	}
	studentID := authStudentID
	if studentID == "" {
		studentID = req.StudentID
		if studentID == "" {
			studentID = sess.StudentID
		}
	}
	if studentID != "" && !s.ownsStudentID(r, studentID) {
		writeError(w, "student_id not permitted", 403)
		return
	}
	if authStudentID != "" && sess.StudentID != "" && sess.StudentID != authStudentID {
		writeError(w, "quiz session does not belong to authenticated user", 403)
		return
	}
	qEng := quiz.NewEngine(s.eng.GetDAG(), s.eng.GetGeneratorRegistry())
	qEng.SkipQuestion(sess)
	if qEng.IsComplete(sess) {
		if studentID != "" && qEng.HasAttempts(sess) {
			if err := s.eng.RecordQuizCompletion(studentID); err != nil {
				log.Printf("handleQuizSkip: RecordQuizCompletion failed for %s: %v", studentID, err)
			}
		}
		s.mu.Lock()
		delete(s.quizSessions, req.SessionID)
		delete(s.quizCreated, req.SessionID)
		s.mu.Unlock()
		writeJSON(w, map[string]interface{}{"done": true, "correct": false, "feedback": "Skipped — no XP awarded.", "xp": 0, "retake_available": true})
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
		"done":               false,
		"correct":            false,
		"feedback":           "Skipped — no XP awarded.",
		"xp":                 0,
		"new_status":         "",
		"remedial":           []string{},
		"concept_id":         cid,
		"concept_name":       name2,
		"question":           prob.Question,
		"grading_type":       s.gradingTypeOf(cid),
		"closed_book":        true,
		"time_limit_seconds": s.eng.TimeLimitFor(studentID, cid),
		"questions_total":    len(sess.Order),
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

// ---------------------------------------------------------------------------
// The assessment boundary.
//
// Diagnostics and tests measure knowledge. They do not teach. So a verdict from
// an assessment may carry the assessment's own state — was it right, how far
// through are we, what should be looked at next — and must not carry anything
// that tells the learner the answer.
//
// `done: true` is not the line. It means "the assessment is over", not "now
// start teaching". A final verdict that returns the worked solution teaches on
// exactly the questions that were just missed, which is the outcome the whole
// CAT exists to avoid: the learner can read the solution instead of having
// demonstrated the knowledge, so the next measurement is contaminated.
//
// The stripped fields, and why each one goes:
//
//   explanation — the generator's worked solution. The whole solution.
//   feedback    — for these two handlers it *is* the explanation; the graders'
//                 own tokens ("Not equivalent", "Incorrect comparison") are
//                 kept because they say how the answer was judged and nothing
//                 about what the answer was.
//   diagnosis   — "You are exactly one away, so a count or a boundary is off by
//                 one." That is not a hint, it is a derivation: the learner
//                 already holds their own wrong answer, so "exactly one away"
//                 hands them the right one.
//
// Stripping happens on the way out, not on the way in. `LastProblem.Explanation`
// stays on the session and `RecordAttempt` still stores the expected answer, so
// the mistakes transcript and the grading system are unchanged; the solution
// simply never travels with a verdict.
//
// `/api/study/answer` deliberately does not use this. That is Learn's seam:
// Learn teaches, and teaching is what it is for. One boundary, two rules, and
// the difference between them is which handler the request reached.
// ---------------------------------------------------------------------------

// assessmentVerdictFields are the keys an assessment response is built from. A
// verdict assembled as map[string]interface{} is passed through
// writeAssessmentVerdict rather than writeJSON, so that adding a field to an
// assessment response cannot accidentally start leaking it: the allowlist
// below is what a response may contain, and a field nobody thought about is
// dropped by default.
const (
	// Keys that report assessment state. These are the whole of what may cross.
	assessmentKeyCorrect         = "correct"
	assessmentKeyDone            = "done"
	assessmentKeyFeedback        = "feedback"
	assessmentKeyProgress        = "progress"
	assessmentKeyReport          = "report"
	assessmentKeyXP              = "xp"
	assessmentKeyNewStatus       = "new_status"
	assessmentKeyRemedial        = "remedial"
	assessmentKeyRetryAvailable  = "retry_available"
	assessmentKeyRetakeAvailable = "retake_available"
	assessmentKeyConceptID       = "concept_id"
	assessmentKeyConceptName     = "concept_name"
	assessmentKeyQuestion        = "question"
	assessmentKeyGradingType     = "grading_type"
	assessmentKeyTimeLimit       = "time_limit_seconds"
	assessmentKeyQuestionsTotal  = "questions_total"
	assessmentKeyClosedBook      = "closed_book"
)

// assessmentVerdictKeys is the allowlist, as a set for lookup. Declared
// separately so the list above reads as documentation and this reads as data.
var assessmentVerdictKeys = map[string]bool{
	assessmentKeyCorrect:         true,
	assessmentKeyDone:            true,
	assessmentKeyFeedback:        true,
	assessmentKeyProgress:        true,
	assessmentKeyReport:          true,
	assessmentKeyXP:              true,
	assessmentKeyNewStatus:       true,
	assessmentKeyRemedial:        true,
	assessmentKeyRetryAvailable:  true,
	assessmentKeyRetakeAvailable: true,
	assessmentKeyConceptID:       true,
	assessmentKeyConceptName:     true,
	assessmentKeyQuestion:        true,
	assessmentKeyGradingType:     true,
	assessmentKeyTimeLimit:       true,
	assessmentKeyQuestionsTotal:  true,
	assessmentKeyClosedBook:      true,
}

// assessmentGraderTokens are the only feedback strings allowed through, and
// only because the graders emit them as bare status and they describe the
// comparison rather than its outcome's value. Anything else in `feedback` is
// treated as instructional and dropped: the quiz handler sets `feedback` to the
// served explanation when the learner answers "don't know", and the diagnostic
// handler sets it to the explanation unconditionally, so keying the check on
// the value rather than the key is what stops a future handler reintroducing
// the leak by reusing the field.
//
// The list is exactly the strings internal/grader can produce. A grader that
// starts emitting something new loses its token until this list names it, which
// is the safe direction: the frontend falls back to the verdict it computed
// from `correct` and shows no invented text.
var assessmentGraderTokens = map[string]bool{
	"Could not parse numeric values": true,
	"Grading service busy":           true,
	"Grading service error":          true,
	"Grading service not found":      true,
	"Grading service response error": true,
	"Grading service timed out":      true,
	"Grading service write error":    true,
	"Incorrect":                      true,
	"Incorrect choice":               true,
	"Incorrect comparison":           true,
	"Incorrect number of elements":   true,
	"Incorrect order":                true,
	"Internal error":                 true,
	"Not equivalent":                 true,
}

// writeAssessmentVerdict writes an assessment response, holding back anything
// that would teach. `verdict` is a map assembled by the handler; every key not
// on the allowlist is dropped, and `feedback` survives only if it is one of the
// graders' bare status tokens.
//
// Dropping rather than blanking matters for a second reason: an absent
// explanation cannot be inspected in devtools, cannot be found by a test that
// greps the response body, and cannot be read off a `undefined` field by a
// client that guesses the key. There is nothing there to leak.
func writeAssessmentVerdict(w http.ResponseWriter, verdict map[string]interface{}) {
	safe := make(map[string]interface{}, len(verdict))
	for key, value := range verdict {
		if !assessmentVerdictKeys[key] {
			continue
		}
		if key == "feedback" {
			text, _ := value.(string)
			if !assessmentGraderTokens[text] {
				continue
			}
		}
		safe[key] = value
	}
	writeJSON(w, safe)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// GET /api/efficacy — first-pass / second-pass instrumentation.
func (s *Server) handleEfficacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	report, err := s.eng.Efficacy(studentID)
	if err != nil {
		writeError(w, "failed to compute efficacy", 500)
		return
	}
	writeJSON(w, report)
}

// GET /api/efficacy/all — public product-wide first-pass / second-pass rates.
func (s *Server) handleEfficacyAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	report, err := s.eng.AggregateEfficacy()
	if err != nil {
		writeError(w, "failed to compute efficacy", 500)
		return
	}
	writeJSON(w, report)
}

// GET /api/efficacy/trend — public longitudinal weekly efficacy + retention
// (docs/efficacy.md). Rate-limited like /api/efficacy/all.
func (s *Server) handleEfficacyTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	trend, err := s.eng.EfficacyTrend()
	if err != nil {
		writeError(w, "failed to compute efficacy trend", 500)
		return
	}
	writeJSON(w, trend)
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
	Email    string `json:"email"`
}

type authRes struct {
	Token               string `json:"token"`
	StudentID           string `json:"student_id"`
	Name                string `json:"name"`
	DiagnosticCompleted bool   `json:"diagnostic_completed"`
}

// POST /api/auth/guest — mint (or reclaim) a bearer token for a guest.
// Body: { student_id? (existing guest ID to claim), name? }.
// Claimable rows are credential-less by construction: no password hash, no
// email, no linked OAuth identity. Anything with a credential → 403 (use
// signup/login — this is what stops guest-claim account takeover).
// Unknown IDs are adopted only under the `guest_` prefix (progress rows
// keyed by that ID predate any students row); anything else → 404.
// Rate-limited with the other auth endpoints.
func (s *Server) handleGuestToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	var req struct {
		StudentID string `json:"student_id"`
		Name      string `json:"name"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	sid := strings.TrimSpace(req.StudentID)
	name := strings.TrimSpace(req.Name)
	if len([]rune(name)) > 50 {
		writeError(w, "name must be 1-50 characters", 400)
		return
	}
	if sid == "" {
		if name == "" {
			name = "Guest"
		}
		st, err := s.repo.CreateStudent(name)
		if err != nil {
			writeError(w, "failed to create guest", 500)
			return
		}
		sid = st.ID
	} else {
		st, err := s.repo.GetStudent(sid)
		if err != nil {
			writeError(w, "failed to look up student", 500)
			return
		}
		if st == nil {
			// Adopt pre-existing progress keyed by a client guest ID.
			if !strings.HasPrefix(sid, "guest_") {
				writeError(w, "unknown student — sign up for a new account", 404)
				return
			}
			if _, err := s.repo.ClaimGuestStudent(sid, name); err != nil {
				writeError(w, "failed to claim guest", 500)
				return
			}
		} else {
			// A row with any credential is a real account — not claimable.
			if st.PasswordHash != "" || strings.TrimSpace(st.Email) != "" {
				writeError(w, "account already registered — log in instead", 403)
				return
			}
			if ids, err := s.repo.ListIdentities(sid); err != nil {
				writeError(w, "failed to look up student", 500)
				return
			} else if len(ids) > 0 {
				writeError(w, "account already registered — log in instead", 403)
				return
			}
		}
	}
	token, err := auth.IssueGuestToken(sid)
	if err != nil {
		writeError(w, "failed to issue token", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"token": token, "student_id": sid})
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
	email := auth.NormalizeEmail(req.Email)
	if email == "" {
		writeError(w, "email is required", 400)
		return
	}
	if err := auth.ValidateEmail(email); err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	// Strict one-email-one-account: any existing row with this address —
	// verified or not — is a 409. Point the caller at login/reset instead
	// of forking a duplicate (same enumeration class as "username is taken").
	if existing, err := s.repo.FindByEmail(email); err != nil {
		writeError(w, "signup failed: "+err.Error(), 500)
		return
	} else if existing != nil {
		writeError(w, "email already in use", 409)
		return
	}
	token, st, err := s.auth.Signup(req.Name, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUsernameTaken) || storage.IsUniqueViolation(err) {
			writeError(w, "username is taken", 409)
			return
		}
		writeError(w, "signup failed: "+err.Error(), 400)
		return
	}
	if err := s.repo.SetEmail(st.ID, email); err != nil {
		writeError(w, "signup failed: "+err.Error(), 500)
		return
	}
	// Best-effort verification dispatch — never fails signup when mail is
	// unconfigured or the send fails (Settings → Verify email stays available).
	if err := s.auth.RequestEmailVerification(st.ID); err != nil {
		log.Printf("auth: signup verification dispatch skipped: %v", err)
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
		if errors.Is(err, auth.ErrGoogleOnly) {
			writeError(w, "this account uses Google sign-in — continue with Google", 401)
			return
		}
		writeError(w, "invalid username or password", 401)
		return
	}
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name, DiagnosticCompleted: st.DiagnosticCompleted})
}

// POST /api/auth/reset/request {identifier} — always 200 (no enumeration),
// 503 when outbound mail is unconfigured.
func (s *Server) handleResetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	var req struct {
		Identifier string `json:"identifier"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if strings.TrimSpace(req.Identifier) == "" {
		writeError(w, "identifier is required", 400)
		return
	}
	if !auth.SMTPConfigured() {
		writeError(w, "password reset not configured", 503)
		return
	}
	if _, err := s.auth.RequestPasswordReset(req.Identifier); err != nil {
		writeError(w, "reset request failed", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// POST /api/auth/reset/complete {token, password} — burns the token,
// sets the password, and logs the user in.
func (s *Server) handleResetComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	token, st, err := s.auth.CompletePasswordReset(req.Token, req.Password)
	if err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name, DiagnosticCompleted: st.DiagnosticCompleted})
}

// PUT /api/auth/password {current_password, new_password} — logged-in
// password change (no email involved).
func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if err := s.auth.ChangePassword(studentID, req.CurrentPassword, req.NewPassword); err != nil {
		writeError(w, err.Error(), 400)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

// PUT /api/profile — update display name and/or recovery email.
// Username is immutable. Email must be valid when given; empty clears it.
func (s *Server) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	studentID, _ := r.Context().Value(authStudentKey{}).(string)
	if studentID == "" {
		writeError(w, "not authenticated", 401)
		return
	}
	var req struct {
		Name  string  `json:"name"`
		Email *string `json:"email,omitempty"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || len([]rune(name)) > 50 {
		writeError(w, "name must be 1-50 characters", 400)
		return
	}
	if err := s.repo.UpdateStudentName(studentID, name); err != nil {
		writeError(w, "failed to update profile", 500)
		return
	}
	email := ""
	if req.Email != nil {
		email = auth.NormalizeEmail(*req.Email)
		if err := auth.ValidateEmail(email); err != nil {
			writeError(w, err.Error(), 400)
			return
		}
		// Strict one-email-one-account (mirrors signup): another account
		// holding this address — verified or not — is a 409. Re-saving
		// your own address stays a no-op success.
		if email != "" {
			if existing, err := s.repo.FindByEmail(email); err != nil {
				writeError(w, "failed to update profile", 500)
				return
			} else if existing != nil && existing.ID != studentID {
				writeError(w, "email already in use", 409)
				return
			}
		}
		if err := s.repo.SetEmail(studentID, email); err != nil {
			writeError(w, "failed to update profile", 500)
			return
		}
	} else if st, err := s.repo.GetStudent(studentID); err == nil && st != nil {
		email = st.Email
	}
	writeJSON(w, map[string]interface{}{"student_id": studentID, "name": name, "email": email})
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
	// The role, so the client can decide whether to *offer* the Admin entry.
	//
	// This is a convenience and nothing more. Hiding a link is not authorization, and a client
	// that lies about the role gets nothing: every /api/admin route re-reads it server-side.
	// It is read here rather than added to storage.Student for the reason given there — a
	// privilege field that silently defaults to "student" when unloaded is the wrong shape.
	role, err := s.repo.GetStudentRole(studentID)
	if err != nil {
		log.Printf("warning: failed to read role for %s: %v", studentID, err)
		role = storage.RoleStudent
	}
	// Placement, so the client reads completion from the server rather than
	// inferring it from local storage or from having displayed the last question.
	// The timestamp is the diagnostic's historical record; the boolean is the
	// gate. Absent (pre-column) means no date is shown, never an invented one.
	var diagnosticCompletedAt string
	if t, ok, err := s.repo.DiagnosticCompletedAt(studentID); err == nil && ok {
		diagnosticCompletedAt = t.UTC().Format(time.RFC3339)
	}
	writeJSON(w, map[string]interface{}{
		"student_id":              st.ID,
		"name":                    st.Name,
		"username":                st.Username,
		"concepts_mastered":       scores.ConceptsMastered,
		"current_streak":          scores.CurrentStreak,
		"level":                   scores.Level,
		"diagnostic_completed":    st.DiagnosticCompleted,
		"diagnostic_completed_at": diagnosticCompletedAt,
		"avatar_url":              st.AvatarURL,
		"email":                   st.Email,
		"email_verified":          st.EmailVerified,
		"has_password":            st.PasswordHash != "",
		"role":                    string(role),
	})
}

// POST /api/auth/google — One-Tap id_token
func (s *Server) handleGoogleOneTap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, 405)
		return
	}
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	var req struct {
		IDToken string `json:"id_token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	if strings.TrimSpace(req.IDToken) == "" {
		writeError(w, "id_token required", 400)
		return
	}
	profile, err := s.auth.VerifyGoogleIDToken(r.Context(), req.IDToken)
	if err != nil {
		writeError(w, "invalid google token: "+err.Error(), 401)
		return
	}
	token, st, err := s.auth.LoginOrCreateGoogle(profile.Name, profile.Email, profile.GoogleID, profile.Picture)
	if err != nil {
		writeError(w, "google login failed: "+err.Error(), 500)
		return
	}
	writeJSON(w, authRes{Token: token, StudentID: st.ID, Name: st.Name, DiagnosticCompleted: st.DiagnosticCompleted})
}

// GET /api/auth/google/login — redirect to Google OAuth consent
func (s *Server) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		clientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	}
	if strings.TrimSpace(clientID) == "" {
		writeError(w, "google auth not configured", 500)
		return
	}
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if strings.TrimSpace(redirectURL) == "" {
		scheme := "https"
		if r.Header.Get("X-Forwarded-Proto") != "" {
			scheme = r.Header.Get("X-Forwarded-Proto")
		} else if strings.HasPrefix(r.Host, "localhost") {
			scheme = "http"
		}
		redirectURL = fmt.Sprintf("%s://%s/api/auth/google/callback", scheme, r.Host)
	}
	// CSRF state
	state := newUUID()
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state, Path: "/", HttpOnly: true, Secure: r.Header.Get("X-Forwarded-Proto") == "https" || r.TLS != nil, SameSite: http.SameSiteLaxMode, MaxAge: 600})
	// Preserve frontend return path
	ret := r.URL.Query().Get("return")
	if ret == "" {
		ret = "/profile"
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_return", Value: ret, Path: "/", HttpOnly: false, SameSite: http.SameSiteLaxMode, MaxAge: 600})
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURL)
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile")
	q.Set("state", state)
	q.Set("access_type", "offline")
	q.Set("prompt", "select_account")
	http.Redirect(w, r, "https://accounts.google.com/o/oauth2/v2/auth?"+q.Encode(), http.StatusFound)
}

// GET /api/auth/google/callback?code=...&state=...
func (s *Server) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if s.auth == nil {
		writeError(w, "authentication is disabled", 400)
		return
	}
	if r.URL.Query().Get("error") != "" {
		http.Redirect(w, r, "/login?error=google_denied", http.StatusFound)
		return
	}
	state := r.URL.Query().Get("state")
	cookie, _ := r.Cookie("oauth_state")
	if state == "" || cookie == nil || cookie.Value != state {
		writeError(w, "invalid oauth state", 400)
		return
	}
	code := r.URL.Query().Get("code")
	if strings.TrimSpace(code) == "" {
		writeError(w, "missing code", 400)
		return
	}
	profile, err := s.auth.ExchangeGoogleCode(r.Context(), code, r)
	if err != nil {
		log.Printf("google code exchange: %v", err)
		http.Redirect(w, r, "/login?error=google_failed", http.StatusFound)
		return
	}
	token, st, err := s.auth.LoginOrCreateGoogle(profile.Name, profile.Email, profile.GoogleID, profile.Picture)
	if err != nil {
		log.Printf("google login create: %v", err)
		http.Redirect(w, r, "/login?error=google_failed", http.StatusFound)
		return
	}
	// Clear state cookies
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: "", Path: "/", MaxAge: -1})
	ret := "/profile"
	if c, err := r.Cookie("oauth_return"); err == nil && c.Value != "" {
		ret = c.Value
		http.SetCookie(w, &http.Cookie{Name: "oauth_return", Value: "", Path: "/", MaxAge: -1})
	}
	// Redirect to frontend with token fragment (frontend will capture and store)
	// Use query ?token=... so static export can read; token is short-lived HS256.
	// FRONTEND_URL (e.g. http://localhost:3000 in split dev) makes the redirect
	// absolute so the token lands on the frontend origin; unset keeps the
	// relative redirect for single-binary deployments.
	u := frontendRedirect(os.Getenv("FRONTEND_URL"), ret, token, st.Name, st.ID)
	http.Redirect(w, r, u, http.StatusFound)
}

// frontendRedirect builds the post-OAuth landing URL. With an http(s) base it
// returns an absolute URL on the frontend origin; otherwise (empty or invalid
// base) it falls back to the relative path so single-binary mode is unchanged.
func frontendRedirect(base, ret, token, name, id string) string {
	q := url.Values{}
	q.Set("token", token)
	q.Set("name", name)
	q.Set("id", id)
	if ret == "" {
		ret = "/profile"
	}
	if !strings.HasPrefix(ret, "/") {
		ret = "/" + ret
	}
	suffix := ret + "?" + q.Encode()
	if strings.TrimSpace(base) == "" {
		return suffix
	}
	parsed, err := url.Parse(strings.TrimSpace(base))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		log.Printf("warning: ignoring invalid FRONTEND_URL %q (want http(s)://host)", base)
		return suffix
	}
	parsed.Path = strings.TrimSuffix(parsed.Path, "/") + ret
	parsed.RawQuery = q.Encode()
	parsed.Fragment = ""
	return parsed.String()
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
