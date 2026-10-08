package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/storage"
)

// Role-based administration.
//
// This is a *different* mechanism from the one in admin.go, and the two are deliberately not
// merged. That file gates content triage behind a shared password from the operator's
// environment — the trust model of someone who can already read the deployment's secrets.
// This one gates account administration behind a role on an individual account, because
// promoting and demoting people is a decision about *who*, and a password shared by everyone
// who holds it cannot record which one of them decided.
//
// They therefore answer different questions, and the distinction is worth stating because
// both live under /admin: `/admin/reports` is triage and is gated by ADMIN_PASSWORD;
// `/admin/users` and `/admin/audit` are account administration and are gated by the role below.
// Neither authorizes the other.

// requireAdmin is the administrative security boundary. Every route under /api/admin that is
// not the triage login goes through it, and there is no other way in.
//
// The sequence is authenticate → identify → load role → allow, and the role step is a database
// read on **every request**. That is the part worth arguing for: a role baked into the token
// would keep granting admin until the token expired, so demoting an administrator would leave
// them able to promote themselves back for as long as their session lived. Reading it per
// request means a demotion takes effect on the next call. The cost is one indexed lookup, which
// is cheaper than the mistake.
//
// Status codes are the three the boundary is specified to produce, and they are distinguishable
// on purpose. 401 means "we do not know who you are"; 403 means "we know, and the answer is
// no". Collapsing them to one would tell a prober that a token of any kind is getting closer.
func (s *Server) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// With auth disabled there is no identity to authorize, so the surface is closed rather
		// than open. `-no-auth` is a development convenience; it must not silently become an
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
		if role != storage.RoleAdmin {
			writeError(w, "admin role required", 403)
			return
		}
		ctx := context.WithValue(r.Context(), authStudentKey{}, studentID)
		next(w, r.WithContext(ctx))
	}
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
	admins, err := s.repo.CountAdmins()
	if err != nil {
		writeError(w, "failed to count admins", 500)
		return
	}
	questions, err := s.repo.CountStoredQuestions()
	if err != nil {
		writeError(w, "failed to count questions", 500)
		return
	}
	recent, err := s.repo.ListAdminAudit(5, 0)
	if err != nil {
		writeError(w, "failed to read audit log", 500)
		return
	}
	conceptCount, domainCount := s.eng.CorpusSize()
	writeJSON(w, map[string]interface{}{
		"users":           users,
		"admins":          admins,
		"concepts":        conceptCount,
		"domains":         domainCount,
		"questions":       questions,
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
func (s *Server) handleAdminUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		writeError(w, "method not allowed", 405)
		return
	}
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/users/"), "/")
	if id == "" {
		writeError(w, "invalid user id", 400)
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
		writeError(w, "role must be student or admin", 400)
		return
	}

	// One call. The repository applies the mutation and writes the audit event in a single
	// transaction, so there is no window in which the role has changed and the record has not
	// — and no code path in which this handler could forget to record.
	actorID, _ := r.Context().Value(authStudentKey{}).(string)
	res, err := s.repo.SetRoleAudited(actorID, id, role)
	switch {
	case errors.Is(err, storage.ErrLastAdmin):
		// 409 rather than 400: the request was well-formed and permitted in general, and it
		// conflicted with current state. A client can retry it after creating another
		// administrator.
		writeError(w, "cannot demote the last administrator", 409)
		return
	case errors.Is(err, storage.ErrRoleUnchanged):
		// Not an error and not a new event: an audit row for a change that did not happen is
		// noise, and the log's value is that every row is a fact. The transaction rolled back,
		// so nothing reached either table.
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

// GET /api/admin/audit?limit=&offset= — the trail, newest first.
//
// Read-only. There is no PUT, PATCH or DELETE for this table on any store, and no route to
// create one, which is what makes "the audit log must not be editable" a property of the
// system rather than a promise about the UI.
func (s *Server) handleAdminAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := r.URL.Query()
	records, err := s.repo.ListAdminAudit(adminLimit(q.Get("limit"), 50, 200), adminOffset(q.Get("offset")))
	if err != nil {
		writeError(w, "failed to read audit log", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"events": records})
}
