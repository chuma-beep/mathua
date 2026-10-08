package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/admin"
	"github.com/chuma-beep/mathua/internal/storage"
)

// Admin V2 handlers: moderation, contributors/invitations, content inspection and the caller's
// own capability readout. Account role changes live in adminrole.go with the V1 handler they
// extend.
//
// Every route here is wrapped in requirePermission, so the authorization is a single visible
// line at the registration site and the check itself is a lookup in internal/admin. A handler
// never decides "is this caller an admin" from the body, the query or the frontend; it reads the
// role the middleware loaded from the database for this request.

// hasPermission re-checks a permission inside a handler for the few routes that carry two
// permissions on one path (a GET that reads and a PATCH that writes). The middleware already
// proved the floor permission; this is the stricter one.
func hasPermission(r *http.Request, perm admin.Permission) bool {
	return admin.Permitted(actorFrom(r).Role, perm)
}

// GET /api/admin/me — who the caller is to this surface, and what they may do. The frontend uses
// it to show only reachable sections; the sections still enforce themselves server-side.
func (s *Server) handleAdminMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	a := actorFrom(r)
	perms := admin.Permissions(a.Role)
	names := make([]string, 0, len(perms))
	for _, p := range perms {
		names = append(names, string(p))
	}
	assignable := make([]string, 0)
	for _, rr := range admin.AssignableRoles(a.Role) {
		assignable = append(assignable, string(rr))
	}
	writeJSON(w, map[string]interface{}{
		"id":               a.ID,
		"role":             string(a.Role),
		"permissions":      names,
		"assignable_roles": assignable,
	})
}

// ── Reported content ──────────────────────────────────────────────────────

var moderationStatuses = map[string]bool{
	"open": true, "reviewing": true, "resolved": true, "dismissed": true,
}

// GET /api/admin/reports?status=&limit=&offset= — the moderation queue.
func (s *Server) handleAdminReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := r.URL.Query()
	status := strings.TrimSpace(q.Get("status"))
	if status != "" && status != "all" && !moderationStatuses[status] {
		writeError(w, "invalid status", 400)
		return
	}
	reports, err := s.repo.ListReports(status, adminLimit(q.Get("limit"), 50, 200), adminOffset(q.Get("offset")))
	if err != nil {
		writeError(w, "failed to list reports", 500)
		return
	}
	if reports == nil {
		reports = []storage.QuestionReport{}
	}
	writeJSON(w, map[string]interface{}{"reports": reports})
}

// GET /api/admin/reports/{id} and PATCH /api/admin/reports/{id}.
//
// The detail read carries the report's own audit trail, so a moderator sees what was decided
// before them, not only the reporter's complaint. The patch records the decision with the actor
// attached; a resolved or dismissed report without a reason is refused, because a decision
// nobody explained is not one the next moderator can learn from.
func (s *Server) handleAdminReport(w http.ResponseWriter, r *http.Request) {
	idStr := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/reports/"), "/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, "invalid report id", 400)
		return
	}
	switch r.Method {
	case http.MethodGet:
		report, err := s.repo.GetReport(id)
		if err != nil {
			writeError(w, "failed to read report", 500)
			return
		}
		if report == nil {
			writeError(w, "report not found", 404)
			return
		}
		history, err := s.repo.ListAdminAuditFor("report", idStr, 50)
		if err != nil {
			writeError(w, "failed to read moderation history", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"report": report, "history": history})
	case http.MethodPatch, http.MethodPut:
		if !hasPermission(r, admin.ReportsModerate) {
			writeError(w, "insufficient permission", 403)
			return
		}
		var req struct {
			Status     string `json:"status"`
			Resolution string `json:"resolution"`
		}
		if err := decodeJSON(w, r, &req); err != nil {
			writeError(w, "invalid request", 400)
			return
		}
		status := strings.TrimSpace(req.Status)
		if !moderationStatuses[status] {
			writeError(w, "status must be open, reviewing, resolved or dismissed", 400)
			return
		}
		resolution := strings.TrimSpace(req.Resolution)
		if (status == "resolved" || status == "dismissed") && resolution == "" {
			writeError(w, "a resolution is required when resolving or dismissing", 400)
			return
		}
		report, err := s.repo.UpdateReportStatusAudited(actorFrom(r).ID, id, status, resolution)
		switch {
		case errors.Is(err, storage.ErrReportUnchanged):
			writeJSON(w, map[string]interface{}{"report": report, "changed": false})
			return
		case errors.Is(err, sql.ErrNoRows):
			writeError(w, "report not found", 404)
			return
		case err != nil:
			log.Printf("admin: moderation of report %d failed: %v", id, err)
			writeError(w, "failed to update report", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"report": report, "changed": true})
	default:
		writeError(w, "method not allowed", 405)
	}
}

// ── Admins and contributors ───────────────────────────────────────────────

// GET /api/admin/admins — who currently holds a staff role, most privileged first.
func (s *Server) handleAdminAdmins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := r.URL.Query()
	staff, err := s.repo.ListStaff(adminLimit(q.Get("limit"), 100, 200), adminOffset(q.Get("offset")))
	if err != nil {
		writeError(w, "failed to list staff", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"admins": staff})
}

// GET /api/admin/invitations — list; POST /api/admin/invitations — create.
//
// Creation mints the raw token once and returns it in the response; the database holds only the
// hash, so the invitation is shown to the inviter and cannot be recovered afterwards. That is
// the point of hashing it: a leaked backup yields nothing replayable.
func (s *Server) handleAdminInvitations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query()
		invs, err := s.repo.ListAdminInvitations(adminLimit(q.Get("limit"), 50, 200), adminOffset(q.Get("offset")))
		if err != nil {
			writeError(w, "failed to list invitations", 500)
			return
		}
		if invs == nil {
			invs = []storage.AdminInvitation{}
		}
		writeJSON(w, map[string]interface{}{"invitations": invs})
	case http.MethodPost:
		if !hasPermission(r, admin.AdminsManage) {
			writeError(w, "insufficient permission", 403)
			return
		}
		var req struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}
		if err := decodeJSON(w, r, &req); err != nil {
			writeError(w, "invalid request", 400)
			return
		}
		email := strings.ToLower(strings.TrimSpace(req.Email))
		if email == "" || !strings.Contains(email, "@") {
			writeError(w, "a valid email is required", 400)
			return
		}
		role, ok := storage.ParseRole(strings.TrimSpace(req.Role))
		if !ok || role == storage.RoleStudent {
			writeError(w, "role must be moderator, admin or owner", 400)
			return
		}
		a := actorFrom(r)
		if !admin.CanAssign(a.Role, storage.RoleStudent, role) {
			writeError(w, "you cannot invite someone to that role", 403)
			return
		}
		raw, hash, err := admin.NewInvitationToken()
		if err != nil {
			writeError(w, "failed to create invitation", 500)
			return
		}
		now := time.Now().UTC()
		id, err := s.repo.CreateAdminInvitation(storage.AdminInvitation{
			Email:     email,
			Role:      role,
			InvitedBy: a.ID,
			CreatedAt: now.Format(time.RFC3339),
			ExpiresAt: now.Add(admin.InvitationTTL).Format(time.RFC3339),
		}, hash)
		if err != nil {
			writeError(w, "failed to create invitation", 500)
			return
		}
		// The audit records that the invitation was created, and to whom, without the token.
		afterJSON, _ := json.Marshal(map[string]string{"email": email, "role": string(role)})
		_, _ = s.repo.RecordAdminAudit(storage.AdminAuditRecord{
			ActorID: a.ID, Action: storage.ActionAdminInviteCreate,
			EntityType: "invitation", EntityID: strconv.FormatInt(id, 10),
			AfterJSON: string(afterJSON),
			CreatedAt: now.Format(time.RFC3339),
		})
		// The raw token is shown once, here, and never stored.
		writeJSON(w, map[string]interface{}{
			"id":         id,
			"token":      raw,
			"email":      email,
			"role":       string(role),
			"expires_at": now.Add(admin.InvitationTTL).Format(time.RFC3339),
		})
	default:
		writeError(w, "method not allowed", 405)
	}
}

// POST /api/admin/invitations/{id}/revoke.
func (s *Server) handleAdminInvitationRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	if !hasPermission(r, admin.AdminsManage) {
		writeError(w, "insufficient permission", 403)
		return
	}
	idStr := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/invitations/"), "/")
	idStr = strings.TrimSuffix(idStr, "/revoke")
	id, err := strconv.ParseInt(strings.Trim(idStr, "/"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, "invalid invitation id", 400)
		return
	}
	inv, err := s.repo.RevokeAdminInvitation(actorFrom(r).ID, id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, "invitation not found or already settled", 404)
		return
	}
	if err != nil {
		writeError(w, "failed to revoke invitation", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"invitation": inv})
}

// POST /api/admin/invitations/accept {token} — any authenticated account may accept the
// invitation addressed to it.
//
// This is the one route under /api/admin that is not staff-gated: it is how a learner becomes
// staff. It requires an authenticated account whose email matches the invitation's, so knowing
// the token is not by itself enough — the invitation is bound to an address, and only the
// account holding that address can turn it into a role.
func (s *Server) handleAdminInvitationAccept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	studentID := actorFrom(r).ID
	if studentID == "" {
		writeError(w, "authentication required", 401)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	token := strings.TrimSpace(req.Token)
	if token == "" {
		writeError(w, "token required", 400)
		return
	}
	hash := admin.HashInvitationToken(token)
	inv, err := s.repo.FindAdminInvitationByTokenHash(hash)
	if err != nil {
		writeError(w, "failed to read invitation", 500)
		return
	}
	// Every failure of the invitation itself — unknown, expired, revoked, already used — reads
	// the same to the caller. A prober learns "that token is not usable", never which near-miss
	// they are closest to.
	if inv == nil || inv.AcceptedAt != "" || inv.RevokedAt != "" || expired(inv.ExpiresAt) {
		writeError(w, "invitation is not valid", 400)
		return
	}
	// Bind the invitation to the address it was sent to.
	if inv.Email != "" {
		st, err := s.repo.GetStudent(studentID)
		if err != nil || st == nil {
			writeError(w, "account not found", 403)
			return
		}
		if !strings.EqualFold(strings.TrimSpace(st.Email), inv.Email) {
			writeError(w, "this invitation was issued to a different address", 403)
			return
		}
	}
	accepted, ok, err := s.repo.ConsumeAdminInvitation(hash, studentID)
	if err != nil {
		writeError(w, "failed to accept invitation", 500)
		return
	}
	if !ok {
		// Lost a race with another acceptance of the same token.
		writeError(w, "invitation is not valid", 400)
		return
	}
	writeJSON(w, map[string]interface{}{"role": string(accepted.Role)})
}

func expired(expiresAt string) bool {
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return true
	}
	return time.Now().UTC().After(t)
}

// ── Content inspection ────────────────────────────────────────────────────

// GET /api/admin/content?q=&limit=&offset= — read-only inspection of the corpus.
//
// Inspection, not editing. The content is code-adjacent corpus data (concepts, generated
// questions) with no safe write path, and the task's guidance is explicit: provide inspection
// rather than invent a fragile editor. What an operator actually needs here is to answer "what
// is this concept, and has anyone complained about it".
func (s *Server) handleAdminContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	limit := adminLimit(r.URL.Query().Get("limit"), 100, 500)
	offset := adminOffset(r.URL.Query().Get("offset"))
	dag := s.eng.GetDAG()
	if dag == nil {
		writeJSON(w, map[string]interface{}{"concepts": []interface{}{}, "total": 0})
		return
	}
	type row struct {
		ID            string   `json:"id"`
		Label         string   `json:"label"`
		Domain        string   `json:"domain"`
		Subdomain     string   `json:"subdomain"`
		GradingType   string   `json:"grading_type"`
		Prerequisites []string `json:"prerequisites"`
	}
	all := dag.Order()
	out := []row{}
	matched := 0
	for _, c := range all {
		if q != "" && !strings.Contains(strings.ToLower(c.ID), q) && !strings.Contains(strings.ToLower(c.Label), q) {
			continue
		}
		matched++
		if matched <= offset || len(out) >= limit {
			continue
		}
		out = append(out, row{
			ID: c.ID, Label: c.Label, Domain: c.Domain, Subdomain: c.Subdomain,
			GradingType: c.GradingType, Prerequisites: c.Prerequisites,
		})
	}
	writeJSON(w, map[string]interface{}{"concepts": out, "total": matched})
}

// adminUserDetail backs GET /api/admin/users/{id}. It returns the non-sensitive projection and
// the account's role history from the audit trail, so "when were they made an admin and by whom"
// is answered from the trail rather than a maintained column.
func (s *Server) adminUserDetail(w http.ResponseWriter, r *http.Request, id string) {
	user, err := s.repo.GetAdminUser(id)
	if err != nil {
		writeError(w, "failed to read user", 500)
		return
	}
	if user == nil {
		writeError(w, "user not found", 404)
		return
	}
	history, err := s.repo.ListAdminAuditFor("user", id, 50)
	if err != nil {
		writeError(w, "failed to read role history", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"user": user, "history": history})
}

// GET /api/admin/audit?entity_type=&entity_id=&limit=&offset= — the trail, newest first.
//
// Read-only. There is no PUT, PATCH or DELETE for this table on any store, and no route to
// create one, which is what makes "the audit log must not be editable" a property of the system
// rather than a promise about the UI.
func (s *Server) handleAdminAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	q := r.URL.Query()
	entityType := strings.TrimSpace(q.Get("entity_type"))
	entityID := strings.TrimSpace(q.Get("entity_id"))
	if entityType != "" && entityID != "" {
		records, err := s.repo.ListAdminAuditFor(entityType, entityID, adminLimit(q.Get("limit"), 100, 500))
		if err != nil {
			writeError(w, "failed to read audit log", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"events": records})
		return
	}
	records, err := s.repo.ListAdminAudit(adminLimit(q.Get("limit"), 50, 200), adminOffset(q.Get("offset")))
	if err != nil {
		writeError(w, "failed to read audit log", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"events": records})
}
