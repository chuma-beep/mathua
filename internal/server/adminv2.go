package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chuma-beep/mathua/internal/admin"
	"github.com/chuma-beep/mathua/internal/leaderboard"
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

// GET /api/admin/me, who the caller is to this surface, and what they may do. The frontend uses
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

// GET /api/admin/reports?status=&limit=&offset=, the moderation queue.
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
	if search, ok := s.repo.(storage.ReportSearchStore); ok {
		query := strings.TrimSpace(q.Get("q"))
		limit := adminLimit(q.Get("limit"), 20, 100)
		offset := adminOffset(q.Get("offset"))
		found, total, err := search.SearchReports(status, query, limit, offset)
		if err != nil {
			writeError(w, "failed to list reports", 500)
			return
		}
		if found == nil {
			found = []storage.QuestionReport{}
		}
		writeJSON(w, map[string]interface{}{"reports": found, "total": total, "limit": limit, "offset": offset})
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
	// /api/admin/reports/{id}/image is the optional screenshot. It is fetched
	// as a blob by the moderation UI because an <img> tag cannot carry the
	// bearer token, and it stays behind the same read permission as the report.
	if strings.HasSuffix(idStr, "/image") {
		id, err := strconv.ParseInt(strings.TrimSuffix(idStr, "/image"), 10, 64)
		if err != nil || id <= 0 {
			writeError(w, "invalid report id", 400)
			return
		}
		if r.Method != http.MethodGet {
			writeError(w, "method not allowed", 405)
			return
		}
		store, ok := s.repo.(storage.ReportImageStore)
		if !ok {
			writeError(w, "screenshots are not supported by this store", 500)
			return
		}
		img, found, err := store.GetReportImage(id)
		if err != nil {
			writeError(w, "failed to load screenshot", 500)
			return
		}
		if !found {
			writeError(w, "this report has no screenshot", 404)
			return
		}
		w.Header().Set("Content-Type", img.ContentType)
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Set("Content-Disposition", "inline")
		w.Write(img.Bytes)
		return
	}

	// /api/admin/reports/{id}/grade is the one sub-resource, and it is handled
	// here rather than as a second pattern so the permission boundary and the
	// id parsing cannot drift apart.
	if strings.HasSuffix(idStr, "/grade") {
		id, err := strconv.ParseInt(strings.TrimSuffix(idStr, "/grade"), 10, 64)
		if err != nil || id <= 0 {
			writeError(w, "invalid report id", 400)
			return
		}
		if r.Method != http.MethodPost {
			writeError(w, "method not allowed", 405)
			return
		}
		s.handleAdminReportGrade(w, r, id)
		return
	}
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

// POST /api/admin/reports/{id}/grade, replay a report's answer through the
// same grader the learner's submission used, and write nothing.
//
// The queue could show the question and the expected answer, but not whether an
// answer was actually accepted, so "is this fixed?" was an opinion. This calls
// engine.GradeAnswer, the single canonical entry point the learner's submission
// calls, so a disagreement with the verdict stored on the report is evidence of
// a bug rather than a re-grade by a different code path.
//
// It deliberately touches no session, attempt, XP or streak: GradeAnswer is
// pure, and this handler never enters SubmitAnswer. Replaying a report is not
// the learner answering it, and a moderation tool that moved a learner's
// mastery would make the act of investigating a report change the evidence.
func (s *Server) handleAdminReportGrade(w http.ResponseWriter, r *http.Request, id int64) {
	var req struct {
		Answer string `json:"answer"`
	}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, "invalid request", 400)
		return
	}
	report, err := s.repo.GetReport(id)
	if err != nil {
		writeError(w, "failed to read report", 500)
		return
	}
	if report == nil {
		writeError(w, "report not found", 404)
		return
	}
	answer := strings.TrimSpace(req.Answer)
	if answer == "" {
		answer = strings.TrimSpace(report.AnswerGiven)
	}
	if report.ConceptID == "" || report.Expected == "" || answer == "" {
		writeError(w, "this report has no concept, expected answer and given answer to replay", 400)
		return
	}
	s.writeGradeResult(w, report.ConceptID, report.Expected, answer)
}

// POST /api/admin/reports/bulk, one decision across a concept's unsettled
// reports.
//
// A fix is usually one bug, not one report: the answer-format defect produced
// eleven reports across eight concepts, and closing them one at a time made the
// operator retype the same reason and left eleven unrelated-looking audit rows.
// This applies the same decision to a set while keeping each report's own
// transaction, stamp and audit row, so the trail never becomes one row claiming
// to speak for many. Settled reports are never revisited: closing a decided
// report again would overwrite who decided it and when.
func (s *Server) handleAdminReportsBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", 405)
		return
	}
	var req struct {
		IDs        []int64 `json:"ids"`
		ConceptID  string  `json:"concept_id"`
		Status     string  `json:"status"`
		Resolution string  `json:"resolution"`
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
	ids := req.IDs
	conceptID := strings.TrimSpace(req.ConceptID)
	if len(ids) == 0 && conceptID != "" {
		unsettled, err := s.repo.ListUnsettledReportIDsForConcept(conceptID)
		if err != nil {
			writeError(w, "failed to list reports", 500)
			return
		}
		ids = unsettled
	}
	if len(ids) == 0 {
		writeError(w, "select at least one report, or a concept with unsettled reports", 400)
		return
	}
	if len(ids) > 200 {
		writeError(w, "too many reports in one decision", 400)
		return
	}
	for _, id := range ids {
		if id <= 0 {
			writeError(w, "invalid report id", 400)
			return
		}
	}
	actor := actorFrom(r).ID
	updated, skipped := 0, 0
	failed := []map[string]interface{}{}
	for _, id := range ids {
		_, err := s.repo.UpdateReportStatusAudited(actor, id, status, resolution)
		switch {
		case err == nil:
			updated++
		case errors.Is(err, storage.ErrReportUnchanged):
			skipped++
		case errors.Is(err, sql.ErrNoRows):
			failed = append(failed, map[string]interface{}{"id": id, "error": "not found"})
		default:
			log.Printf("admin: bulk moderation of report %d failed: %v", id, err)
			failed = append(failed, map[string]interface{}{"id": id, "error": "update failed"})
		}
	}
	writeJSON(w, map[string]interface{}{"updated": updated, "skipped": skipped, "failed": failed})
}

// ── Admins and contributors ───────────────────────────────────────────────

// GET /api/admin/admins, who currently holds a staff role, most privileged first.
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

// GET /api/admin/invitations, list; POST /api/admin/invitations, create.
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

// POST /api/admin/invitations/accept {token}, any authenticated account may accept the
// invitation addressed to it.
//
// This is the one route under /api/admin that is not staff-gated: it is how a learner becomes
// staff. It requires an authenticated account whose email matches the invitation's, so knowing
// the token is not by itself enough, the invitation is bound to an address, and only the
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
	// Every failure of the invitation itself, unknown, expired, revoked, already used, reads
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

// GET /api/admin/content?q=&limit=&offset=, read-only inspection of the corpus.
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
	// overridden=1 lists only the concepts whose served answer format differs
	// from the default their grading type would give. The lab renders that as a
	// table so a new override, or a removed one, is visible without reading the
	// generator source.
	overriddenOnly := r.URL.Query().Get("overridden") == "1"
	type row struct {
		ID               string   `json:"id"`
		Label            string   `json:"label"`
		Domain           string   `json:"domain"`
		Subdomain        string   `json:"subdomain"`
		GradingType      string   `json:"grading_type"`
		Prerequisites    []string `json:"prerequisites"`
		AnswerHint       string   `json:"answer_hint"`
		AnswerInputMode  string   `json:"answer_input_mode"`
		AnswerOverridden bool     `json:"answer_overridden"`
	}
	all := dag.Order()
	out := []row{}
	matched := 0
	for _, c := range all {
		if q != "" && !strings.Contains(strings.ToLower(c.ID), q) && !strings.Contains(strings.ToLower(c.Label), q) {
			continue
		}
		format, overridden := s.answerFormatIsOverride(c.ID, c.GradingType)
		if overriddenOnly && !overridden {
			continue
		}
		matched++
		if matched <= offset || len(out) >= limit {
			continue
		}
		out = append(out, row{
			ID: c.ID, Label: c.Label, Domain: c.Domain, Subdomain: c.Subdomain,
			GradingType: c.GradingType, Prerequisites: c.Prerequisites,
			AnswerHint: format.Hint, AnswerInputMode: format.InputMode,
			AnswerOverridden: overridden,
		})
	}
	writeJSON(w, map[string]interface{}{"concepts": out, "total": matched})
}

// adminUserActivity is the counts-only projection shown beside an account. It answers "how
// much has this person used the product and what did they send us" without exposing a single
// question, answer or attempt row: a moderator investigating a report needs the shape of the
// account, not a window onto the learner's work.
type adminUserActivity struct {
	XPTotal               int    `json:"xp_total"`
	XPToday               int    `json:"xp_today"`
	ConceptsTouched       int    `json:"concepts_touched"`
	ConceptsMastered      int    `json:"concepts_mastered"`
	Attempts              int    `json:"attempts"`
	LastAttemptAt         string `json:"last_attempt_at"`
	LastActiveDay         string `json:"last_active_day"`
	DiagnosticCompletedAt string `json:"diagnostic_completed_at"`
	ReportsFiled          int    `json:"reports_filed"`
}

// adminUserActivity reads the counts from the accessors that already own them. It returns an
// error rather than zeros on a failed read, because a profile that silently says "0 attempts"
// for a database error is worse than one that says it could not load.
func (s *Server) adminUserActivity(id string) (adminUserActivity, error) {
	var out adminUserActivity
	total, today, err := s.repo.GetXP(id)
	if err != nil {
		return out, err
	}
	out.XPTotal, out.XPToday = total, today

	progress, err := s.repo.GetAllProgress(id)
	if err != nil {
		return out, err
	}
	for _, p := range progress {
		if p == nil {
			continue
		}
		out.ConceptsTouched++
		if p.Status == "MASTERED" {
			out.ConceptsMastered++
		}
	}

	attempts, err := s.repo.GetAttemptsForStudent(id)
	if err != nil {
		return out, err
	}
	out.Attempts = len(attempts)
	for _, a := range attempts {
		if ts := a.Timestamp.UTC().Format(time.RFC3339); ts > out.LastAttemptAt {
			out.LastAttemptAt = ts
		}
	}

	// Five years is not a retention policy, it is wide enough that a dormant account still
	// shows its last active day instead of looking like an account that never started.
	days, err := s.repo.GetDailyActivity(id, 365*5)
	if err != nil {
		return out, err
	}
	if len(days) > 0 {
		out.LastActiveDay = days[len(days)-1].Date
	}

	if at, ok, err := s.repo.DiagnosticCompletedAt(id); err != nil {
		return out, err
	} else if ok {
		out.DiagnosticCompletedAt = at.UTC().Format(time.RFC3339)
	}

	filed, err := s.repo.CountReportsByReporter(id)
	if err != nil {
		return out, err
	}
	out.ReportsFiled = filed
	return out, nil
}

// adminAccountAvatar is the account-detail view of the same avatar fields the
// league roster carries. It is separate from storage.AdminUser because that
// projection is also written into audit snapshots, and a role change has no
// reason to record a picture.
type adminAccountAvatar struct {
	AvatarURL      string             `json:"avatar_url,omitempty"`
	AvatarDicebear *storage.AvatarRef `json:"avatar_dicebear,omitempty"`
	AvatarPreset   *int               `json:"avatar_preset,omitempty"`
}

// adminUserDetail backs GET /api/admin/users/{id}. It returns the non-sensitive projection, the
// account's role history from the audit trail, the counts-only activity summary and the reports
// that account filed. The reporter's own reports are included because the question an
// administrator arrives with is usually "who is this and what did they send", and the reports
// list is the only place the two meet.
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
	student, err := s.repo.GetStudent(id)
	if err != nil {
		log.Printf("admin user avatar %s: %v", id, err)
		writeError(w, "failed to read account avatar", 500)
		return
	}
	avatar := adminAccountAvatar{}
	if student != nil {
		avatar.AvatarURL, avatar.AvatarDicebear, avatar.AvatarPreset =
			adminAvatar(student.ID, student.AvatarURL, student.Settings)
	}
	history, err := s.repo.ListAdminAuditFor("user", id, 50)
	if err != nil {
		writeError(w, "failed to read role history", 500)
		return
	}
	activity, err := s.adminUserActivity(id)
	if err != nil {
		log.Printf("admin user activity %s: %v", id, err)
		writeError(w, "failed to read account activity", 500)
		return
	}
	reports, err := s.repo.ListReportsByReporter(id, 20)
	if err != nil {
		writeError(w, "failed to read reports", 500)
		return
	}
	writeJSON(w, map[string]interface{}{
		"user": user, "avatar": avatar, "history": history, "activity": activity, "reports": reports,
	})
}

// adminLeagueRow is one seat in an administrator's view of a league. It carries the account id
// because the page links each row to account detail; the public board deliberately does not,
// and this is the surface where that difference is intended.
type adminLeagueRow struct {
	Rank      int    `json:"rank"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
	XP        int    `json:"xp"`
	// AvatarURL is already resolved for the admin view: a custom upload uses
	// the public by-id route so the admin page never renders its own avatar,
	// and a Google photo keeps its remote URL. The remaining fields let the
	// page build a DiceBear image or fall back to a legacy preset.
	AvatarURL      string             `json:"avatar_url,omitempty"`
	AvatarDicebear *storage.AvatarRef `json:"avatar_dicebear,omitempty"`
	AvatarPreset   *int               `json:"avatar_preset,omitempty"`
}

// adminLeague is one league in the week's pyramid. Members is omitted unless the caller asked
// for that league, so browsing every league costs one grouped query rather than one roster read
// per league.
type adminLeague struct {
	Tier     string           `json:"tier"`
	LeagueNo int              `json:"league_no"`
	Size     int              `json:"size"`
	Movers   int              `json:"movers"`
	Members  []adminLeagueRow `json:"members,omitempty"`
}

// adminAvatar resolves the pieces an administrator's page needs to draw one
// picture. A custom upload is served by the public by-id route, because the
// admin page must not render its own avatar; a Google photo is already a
// complete URL. The same precedence applies to every admin surface that shows
// an account.
func adminAvatar(studentID, avatarURL, settings string) (string, *storage.AvatarRef, *int) {
	custom, dicebear, version := storage.AvatarBits(settings)
	if custom && studentID != "" {
		avatarURL = "/api/avatar/" + url.PathEscape(studentID)
		if version > 0 {
			avatarURL += "?v=" + strconv.Itoa(version)
		}
	}
	return avatarURL, dicebear, storage.AvatarPreset(settings)
}

// GET /api/admin/leagues?week=YYYY-MM-DD&tier=bronze&league=3
//
// The whole week's leagues, or one league's standings when tier and league are given. The public
// league endpoint answers only "where am I", which leaves an owner unable to see the pyramid at
// all: a cohort reporting a problem cannot be located, and a tier that is empty or has one
// learner in it is invisible. This is the read that closes that gap.
//
// A week key must be the Monday the league week starts, and is validated here rather than
// trusted, because a bare date is a join key into a table of seats. Without a week the current
// one is used, and the previous week is settled first exactly as a learner's own read does, so
// the numbers an administrator sees are the numbers the transition would act on.
func (s *Server) handleAdminLeagues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", 405)
		return
	}
	now := time.Now().UTC()
	start := leaderboard.WeekStart(now)
	if raw := strings.TrimSpace(r.URL.Query().Get("week")); raw != "" {
		parsed, err := storage.WeekFromString(raw)
		if err != nil {
			writeError(w, "week must be a Monday in YYYY-MM-DD form", 400)
			return
		}
		start = parsed
	}
	week := start.Format("2006-01-02")
	end := start.AddDate(0, 0, 7)

	// Only the current week needs settling, and only because its seats come from the previous
	// week's final standings. A past week is already written and must not be rewritten by a read.
	if start.Equal(leaderboard.WeekStart(now)) {
		if err := s.repo.ResolveLeagueWeek(week, start.AddDate(0, 0, -7), start); err != nil {
			log.Printf("admin leagues resolve %s: %v", week, err)
			writeError(w, "failed to settle the league week", 500)
			return
		}
	}

	cohorts, err := s.repo.LeagueCohorts(week)
	if err != nil {
		log.Printf("admin leagues cohorts %s: %v", week, err)
		writeError(w, "failed to read leagues", 500)
		return
	}

	tier := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tier")))
	leagueNo := 0
	if raw := strings.TrimSpace(r.URL.Query().Get("league")); raw != "" {
		n, convErr := strconv.Atoi(raw)
		if convErr != nil || n < 1 {
			writeError(w, "league must be a positive number", 400)
			return
		}
		leagueNo = n
	}
	if (tier == "") != (leagueNo == 0) {
		writeError(w, "tier and league are read together", 400)
		return
	}

	out := make([]adminLeague, 0, len(cohorts))
	seats := 0
	found := leagueNo == 0
	for _, c := range cohorts {
		seats += c.Size
		entry := adminLeague{
			Tier:     c.Tier,
			LeagueNo: c.LeagueNo,
			Size:     c.Size,
			Movers:   storage.LeagueMovers(c.Size),
		}
		if leagueNo != 0 {
			if c.Tier != tier || c.LeagueNo != leagueNo {
				continue
			}
			found = true
			rows, err := s.repo.LeagueRoster(week, c.Tier, c.LeagueNo, start, end)
			if err != nil {
				log.Printf("admin leagues roster %s %s/%d: %v", week, c.Tier, c.LeagueNo, err)
				writeError(w, "failed to read that league", 500)
				return
			}
			for _, row := range rows {
				avatarURL, pick, preset := adminAvatar(row.StudentID, row.AvatarURL, row.Settings)
				entry.Members = append(entry.Members, adminLeagueRow{
					Rank:           row.Rank,
					Name:           row.Name,
					Username:       row.Username,
					StudentID:      row.StudentID,
					XP:             row.XP,
					AvatarURL:      avatarURL,
					AvatarDicebear: pick,
					AvatarPreset:   preset,
				})
			}
		}
		out = append(out, entry)
	}
	if !found {
		// An administrator who typed a league that does not exist is told so, rather than shown
		// an empty roster that reads like a league nobody joined.
		writeError(w, "no such league in that week", 404)
		return
	}

	writeJSON(w, map[string]interface{}{
		"week":       week,
		"week_start": start,
		"week_end":   end,
		"promote":    storage.PromoteCount,
		"demote":     storage.DemoteCount,
		"seats":      seats,
		"leagues":    out,
	})
}

// GET /api/admin/audit?entity_type=&entity_id=&limit=&offset=, the trail, newest first.
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
	if search, ok := s.repo.(storage.AuditSearchStore); ok {
		query := strings.TrimSpace(q.Get("q"))
		limit := adminLimit(q.Get("limit"), 50, 200)
		offset := adminOffset(q.Get("offset"))
		records, total, err := search.SearchAdminAudit(query, limit, offset)
		if err != nil {
			writeError(w, "failed to read audit log", 500)
			return
		}
		writeJSON(w, map[string]interface{}{"events": records, "total": total, "limit": limit, "offset": offset})
		return
	}
	records, err := s.repo.ListAdminAudit(adminLimit(q.Get("limit"), 50, 200), adminOffset(q.Get("offset")))
	if err != nil {
		writeError(w, "failed to read audit log", 500)
		return
	}
	writeJSON(w, map[string]interface{}{"events": records})
}
