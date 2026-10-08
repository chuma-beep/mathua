package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/admin"
	"github.com/chuma-beep/mathua/internal/storage"
)

// Role-based administration. Every route under /api/admin is gated by requirePermission,
// which reads the caller's role from the database on each request and checks it against the
// permission the route names. The shared-password content triage V1 had is gone: Admin V2 folds
// moderation into the same role model, so there is exactly one administrative mechanism.

// requireAdmin is the floor of the administrative security boundary: it lets through any staff
// member (moderator or above) and injects the caller's identity and role into the request
// context. It is `requirePermission(admin.Access, …)` under the old name, kept because most
// routes mean exactly "any staff member" and should not have to name a permission to say so.
//
// The role is read from the database on **every** request rather than trusted from the token.
// A role baked into a token keeps granting access until the token expires, so demoting someone
// would leave them able to act for the life of their session — the failure this boundary exists
// to prevent. The cost is one indexed lookup, which is cheaper than the mistake.
func (s *Server) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return s.requirePermission(admin.Access, next)
}

// requirePermission is the real boundary. The sequence is authenticate → identify → load role →
// check the permission → allow, and every staff route names the permission it needs.
//
// Status codes are distinguishable on purpose: 401 means "we do not know who you are", 403 means
// "we know, and the answer is no". Collapsing them tells a prober that a token of any kind is
// getting closer.
func (s *Server) requirePermission(perm admin.Permission, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// With auth disabled there is no identity to authorize, so the surface is closed rather
		// than open. `-no-auth` is a development convenience; it must not become an
		// unauthenticated admin.
		if s.auth == nil {
			writeError(w, "authentication is disabled", 401)
			return
		}
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			writeError(w, "missing authorization", 401)
			return
		}
		studentID, err := s.auth.ValidateToken(token)
		if err != nil || studentID == "" {
			writeError(w, "invalid token", 401)
			return
		}
		role, err := s.repo.GetStudentRole(studentID)
		if err != nil {
			writeError(w, "could not read role", 500)
			return
		}
		if !admin.Permitted(role, perm) {
			// A learner who is authenticated but not staff gets 403, the same as a staff member
			// who lacks one permission. The two are not distinguished because the reason is the
			// same to the caller: this account may not do this.
			writeError(w, "insufficient permission", 403)
			return
		}
		ctx := context.WithValue(r.Context(), authStudentKey{}, studentID)
		ctx = context.WithValue(ctx, authRoleKey{}, role)
		next(w, r.WithContext(ctx))
	}
}

// actorFrom returns the authenticated actor id and role the middleware injected. Both are
// present on every request that reached a handler through requirePermission.
func actorFrom(r *http.Request) (string, storage.Role) {
	id, _ := r.Context().Value(authStudentKey{}).(string)
	role, _ := r.Context().Value(authRoleKey{}).(storage.Role)
	return id, role
}

// adminLimit clamps a paging parameter. A client asking for 100000 rows gets the cap, not an
// unbounded read: these endpoints are read by a human on a page.
func adminLimit(q string, def, max int) int {
	n, err := strconv.Atoi(strings.TrimSpace(q))
	if err != nil || n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

func adminOffset(q string) int {
	n, err := strconv.Atoi(strings.TrimSpace(q))
	if err != nil || n < 0 {
		return 0
	}
	return n
}

// GET /api/admin/overview — counts, and nothing else.
//
// Every figure here is read live from the system on request. Nothing is precomputed into a
// dashboard table, because a number that has to be refreshed is a number that will be stale in
// the screenshot someone later uses to argue about something. Concepts and domains come from
// the loaded corpus rather than the database: the corpus is the definition of what Mathua
// teaches, so a count taken from a query would report what has been visited.
func (s *Server) handleAdminOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	users, err := s.repo.CountAdminUsers()
	if err != nil {
		writeError(w, "failed to count users", 500)
		return
	}
	staff, err := s.repo.CountStaff()
	if err != nil {
		writeError(w, "failed to count staff", 500)
		return
	}
	owners, err := s.repo.CountOwners()
	if err != nil {
		writeError(w, "failed to count owners", 500)
		return
	}
	questions, err := s.repo.CountStoredQuestions()
	if err != nil {
		writeError(w, "failed to count questions", 500)
		return
	}
	reportCounts, err := s.repo.CountReportsByStatus()
	if err != nil {
		writeError(w, "failed to count reports", 500)
		return
	}
	// Outstanding moderation work, split so the dashboard can answer "what needs me now".
	recent, err := s.repo.ListAdminAudit(5, 0)
	if err != nil {
		writeError(w, "failed to read audit log", 500)
		return
	}
	conceptCount, domainCount := s.eng.CorpusSize()
	writeJSON(w, map[string]interface{}{
		"users":     users,
		"staff":     staff,
		"owners":    owners,
		"concepts":  conceptCount,
		"domains":   domainCount,
		"questions": questions,
		"reports": map[string]int{
			"open":      reportCounts["open"],
			"reviewing": reportCounts["reviewing"],
			"resolved":  reportCounts["resolved"],
			"dismissed": reportCounts["dismissed"],
			// outstanding is what the operator must still act on.
			"outstanding": reportCounts["open"] + reportCounts["reviewing"],
		},
		"recent_activity": recent,
	})
}

// GET /api/admin/users?q=&limit=&offset= — search and list accounts.
func (s *Server) handleAdminUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := r.URL.Query()
	users, err := s.repo.SearchAdminUsers(q.Get("q"), adminLimit(q.Get("limit"), 50, 200), adminOffset(q.Get("offset")))
	if err != nil {
		writeError(w, "failed to search users", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"users": users})
}

// PATCH /api/admin/users/{id} {"role":"admin"|"student"} — the only mutation V1 has.
//
// The body contributes exactly one field. `role` is read and validated; every other key in the
// payload is ignored rather than applied, so a request carrying `{"role":"admin","xp_total":0,
// "password_hash":"..."}` changes the role and nothing else. That is asserted by a test rather
// than left to review, because mass assignment is the failure this endpoint is most able to
// have: it already holds the ability to write to an arbitrary row.
func (s *Server) handleAdminUser(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/users/"), "/")
	if id == "" {
		writeError(w, "invalid user id", 400)
		return
	}
	if r.Method == http.MethodGet {
		s.adminUserDetail(w, r, id)
		return
	}
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		writeError(w, "method not allowed", 405)
		return
	}
	// The route's floor permission is UsersRead (so a moderator-area caller can be told apart
	// from a stranger), but changing a role is stricter. Re-checked here because the same path
	// serves the read.
	if !hasPermission(r, admin.AdminsManage) {
		writeError(w, "insufficient permission", 403)
		return
	}
	var req struct {
		Role string `json:"role"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	role, ok := storage.ParseRole(req.Role)
	if !ok {
		writeError(w, "role must be student, moderator, admin or owner", 400)
		return
	}

	actorID, actorRole := actorFrom(r)
	// Read the target first: the authorization decision is a comparison between the actor, the
	// target's *current* role and the requested one, so it needs all three before it can allow
	// anything. A missing target is 404 here rather than a student-shaped default.
	target, err := s.repo.GetAdminUser(id)
	if err != nil {
		writeError(w, "failed to read user", 500)
		return
	}
	if target == nil {
		writeError(w, "user not found", 404)
		return
	}
	// Self-promotion is the specific self-change worth naming. Everything else about acting on
	// yourself is either harmless or already blocked by the rank rule (you cannot outrank
	// yourself), but "I made myself an owner" is the escalation to refuse in so many words.
	if id == actorID && admin.CouldEscalate(target.Role, role) {
		writeError(w, "you cannot change your own role", 403)
		return
	}
	if !admin.CanAssign(actorRole, target.Role, role) {
		writeError(w, "you cannot assign that role", 403)
		return
	}

	// One call. The repository applies the mutation and writes the audit event in a single
	// transaction, so there is no window in which the role has changed and the record has not
	// — and no code path in which this handler could forget to record.
	res, err := s.repo.SetRoleAudited(actorID, id, role)
	switch {
	case errors.Is(err, storage.ErrLastOwner):
		// 409 rather than 400: the request was well-formed and permitted in general, and it
		// conflicted with current state. A client can retry it after creating another owner.
		writeError(w, "cannot remove the last owner", 409)
		return
	case errors.Is(err, storage.ErrRoleUnchanged):
		writeJSON(w, map[string]interface{}{"user": res.After, "changed": false})
		return
	case errors.Is(err, sql.ErrNoRows):
		writeError(w, "user not found", 404)
		return
	case err != nil:
		log.Printf("admin: role change for %s failed: %v", id, err)
		writeError(w, "failed to change role", 500)
		return
	}

	writeJSON(w, map[string]interface{}{"user": res.After, "changed": true})
}
