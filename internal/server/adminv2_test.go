package server

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/admin"
	"github.com/chuma-beep/mathua/internal/storage"
)

// Admin V2 authorization, at the API layer. Every case here is a request against a real mux with
// a real database, because the property being tested is "the server refuses", not "the UI hides".

// The routes that a non-staff account must not reach, and must not reach with any verb.
var staffOnlyRoutes = []struct{ method, path string }{
	{"GET", "/api/admin/me"},
	{"GET", "/api/admin/overview"},
	{"GET", "/api/admin/reports"},
	{"PATCH", "/api/admin/reports/1"},
	{"GET", "/api/admin/users"},
	{"GET", "/api/admin/admins"},
	{"GET", "/api/admin/invitations"},
	{"POST", "/api/admin/invitations"},
	{"GET", "/api/admin/audit"},
	{"GET", "/api/admin/content"},
}

func TestAdminV2RejectsUnauthenticated(t *testing.T) {
	e := newAdminEnv(t)
	for _, rt := range staffOnlyRoutes {
		rec := e.do(t, rt.method, rt.path, "", "")
		if rec.Code != 401 {
			t.Errorf("%s %s: got %d, want 401", rt.method, rt.path, rec.Code)
		}
	}
}

func TestAdminV2RejectsLearner(t *testing.T) {
	e := newAdminEnv(t)
	for _, rt := range staffOnlyRoutes {
		rec := e.do(t, rt.method, rt.path, e.studentTok, "")
		if rec.Code != 403 {
			t.Errorf("%s %s as a learner: got %d, want 403 (%s)", rt.method, rt.path, rec.Code, rec.Body.String())
		}
	}
}

// A moderator works the moderation queue and inspects content; it does not see accounts, staff,
// invitations or the audit trail.
func TestModeratorScopeIsReportsAndContentOnly(t *testing.T) {
	e := newAdminEnv(t)
	if err := e.store.SetStudentRole(e.studentID, storage.RoleModerator); err != nil {
		t.Fatalf("make moderator: %v", err)
	}
	allowed := []struct{ method, path string }{
		{"GET", "/api/admin/me"},
		{"GET", "/api/admin/overview"},
		{"GET", "/api/admin/reports"},
		{"GET", "/api/admin/content"},
	}
	for _, rt := range allowed {
		if rec := e.do(t, rt.method, rt.path, e.studentTok, ""); rec.Code != 200 {
			t.Errorf("moderator %s %s: got %d, want 200 (%s)", rt.method, rt.path, rec.Code, rec.Body.String())
		}
	}
	denied := []struct{ method, path string }{
		{"GET", "/api/admin/users"},
		{"GET", "/api/admin/admins"},
		{"GET", "/api/admin/invitations"},
		{"GET", "/api/admin/audit"},
	}
	for _, rt := range denied {
		if rec := e.do(t, rt.method, rt.path, e.studentTok, ""); rec.Code != 403 {
			t.Errorf("moderator %s %s: got %d, want 403", rt.method, rt.path, rec.Code)
		}
	}
}

// An administrator cannot promote themselves and cannot bestow the owner role: both are the
// escalation the role model exists to stop.
func TestAdminCannotEscalate(t *testing.T) {
	e := newAdminEnv(t)
	// "third" is made an admin — a real manager, but not an owner.
	if err := e.store.SetStudentRole(e.thirdID, storage.RoleAdmin); err != nil {
		t.Fatalf("promote third: %v", err)
	}
	// Admin attempts to grant itself owner.
	if rec := e.do(t, "PATCH", "/api/admin/users/"+e.thirdID, e.thirdTok, `{"role":"owner"}`); rec.Code != 403 {
		t.Errorf("admin self-escalation to owner: got %d, want 403 (%s)", rec.Code, rec.Body.String())
	}
	if role, _ := e.store.GetStudentRole(e.thirdID); role != storage.RoleAdmin {
		t.Errorf("self-escalation changed the role to %q", role)
	}
	// Admin attempts to promote a learner straight to owner.
	if rec := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.thirdTok, `{"role":"owner"}`); rec.Code != 403 {
		t.Errorf("admin granting owner to another: got %d, want 403", rec.Code)
	}
}

// ── Invitations ───────────────────────────────────────────────────────────

func createInvite(t *testing.T, e *adminEnv, email, role, token string) (int, string, int64) {
	t.Helper()
	body := `{"email":"` + email + `","role":"` + role + `"}`
	rec := e.do(t, "POST", "/api/admin/invitations", token, body)
	if rec.Code != 200 {
		return rec.Code, "", 0
	}
	var out struct {
		ID    int64  `json:"id"`
		Token string `json:"token"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec.Code, out.Token, out.ID
}

func TestInvitationCreateAndAcceptGrantsRole(t *testing.T) {
	e := newAdminEnv(t)
	code, raw, _ := createInvite(t, e, "ada@example.com", "moderator", e.adminTok)
	if code != 200 || raw == "" {
		t.Fatalf("create invitation: %d token=%q", code, raw)
	}
	// The invited learner accepts with their own token.
	rec := e.do(t, "POST", "/api/admin/invitations/accept", e.studentTok, `{"token":"`+raw+`"}`)
	if rec.Code != 200 {
		t.Fatalf("accept: got %d (%s)", rec.Code, rec.Body.String())
	}
	if role, _ := e.store.GetStudentRole(e.studentID); role != storage.RoleModerator {
		t.Errorf("role after accept = %q, want moderator", role)
	}
	// The learner's existing token now carries the role — read per request.
	if rec := e.do(t, "GET", "/api/admin/reports", e.studentTok, ""); rec.Code != 200 {
		t.Errorf("promoted account still denied reports: %d", rec.Code)
	}
}

func TestInvitationCannotBeReused(t *testing.T) {
	e := newAdminEnv(t)
	_, raw, _ := createInvite(t, e, "ada@example.com", "moderator", e.adminTok)
	if rec := e.do(t, "POST", "/api/admin/invitations/accept", e.studentTok, `{"token":"`+raw+`"}`); rec.Code != 200 {
		t.Fatalf("first accept: %d", rec.Code)
	}
	// Second accept, same token, even by the same account.
	if rec := e.do(t, "POST", "/api/admin/invitations/accept", e.studentTok, `{"token":"`+raw+`"}`); rec.Code == 200 {
		t.Errorf("a spent invitation was accepted again")
	}
}

func TestInvitationRevokedCannotBeAccepted(t *testing.T) {
	e := newAdminEnv(t)
	_, raw, id := createInvite(t, e, "ada@example.com", "moderator", e.adminTok)
	rec := e.do(t, "POST", "/api/admin/invitations/"+itoa(id)+"/revoke", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("revoke: got %d (%s)", rec.Code, rec.Body.String())
	}
	if rec := e.do(t, "POST", "/api/admin/invitations/accept", e.studentTok, `{"token":"`+raw+`"}`); rec.Code == 200 {
		t.Errorf("a revoked invitation was accepted")
	}
}

func TestInvitationExpiredCannotBeAccepted(t *testing.T) {
	e := newAdminEnv(t)
	raw := "expired-token-value"
	hash := admin.HashInvitationToken(raw)
	now := time.Now().UTC()
	if _, err := e.store.CreateAdminInvitation(storage.AdminInvitation{
		Email: "ada@example.com", Role: storage.RoleModerator, InvitedBy: e.adminID,
		CreatedAt: now.Add(-48 * time.Hour).Format(time.RFC3339),
		ExpiresAt: now.Add(-time.Hour).Format(time.RFC3339),
	}, hash); err != nil {
		t.Fatalf("seed expired invitation: %v", err)
	}
	if rec := e.do(t, "POST", "/api/admin/invitations/accept", e.studentTok, `{"token":"`+raw+`"}`); rec.Code == 200 {
		t.Errorf("an expired invitation was accepted")
	}
	if role, _ := e.store.GetStudentRole(e.studentID); role != storage.RoleStudent {
		t.Errorf("expired invitation granted a role anyway: %q", role)
	}
}

// The invitation is bound to the invited address. Knowing the token is not enough if it was
// issued to someone else.
func TestInvitationBoundToInvitedAddress(t *testing.T) {
	e := newAdminEnv(t)
	// Invitation for Ada, but "third" (third@example.com) tries to accept it.
	_, raw, _ := createInvite(t, e, "ada@example.com", "moderator", e.adminTok)
	if rec := e.do(t, "POST", "/api/admin/invitations/accept", e.thirdTok, `{"token":"`+raw+`"}`); rec.Code != 403 {
		t.Errorf("wrong-address accept: got %d, want 403 (%s)", rec.Code, rec.Body.String())
	}
	if role, _ := e.store.GetStudentRole(e.thirdID); role != storage.RoleStudent {
		t.Errorf("the wrong account gained a role: %q", role)
	}
}

func TestInvitationCreateRequiresManagePermission(t *testing.T) {
	e := newAdminEnv(t)
	// Moderator cannot invite.
	if err := e.store.SetStudentRole(e.studentID, storage.RoleModerator); err != nil {
		t.Fatalf("moderator: %v", err)
	}
	if rec := e.do(t, "POST", "/api/admin/invitations", e.studentTok, `{"email":"x@example.com","role":"admin"}`); rec.Code != 403 {
		t.Errorf("moderator created an invitation: %d", rec.Code)
	}
}

func TestInvitationCreateCannotExceedActorAuthority(t *testing.T) {
	e := newAdminEnv(t)
	if err := e.store.SetStudentRole(e.thirdID, storage.RoleAdmin); err != nil {
		t.Fatalf("admin: %v", err)
	}
	if code, _, _ := createInvite(t, e, "x@example.com", "owner", e.thirdTok); code != 403 {
		t.Errorf("an admin invited someone as owner: got %d, want 403", code)
	}
}

// ── Reported content ──────────────────────────────────────────────────────

func mustCreateReport(t *testing.T, e *adminEnv) int64 {
	t.Helper()
	rec := e.do(t, "POST", "/api/reports", "", `{"concept_id":"a","kind":"question","question":"2+2=?","reason":"wrong_answer","reporter_id":"guest_1"}`)
	if rec.Code != 200 {
		t.Fatalf("create report: %d (%s)", rec.Code, rec.Body.String())
	}
	var out struct {
		ID int64 `json:"id"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return out.ID
}

func TestModerationDecisionIsAudited(t *testing.T) {
	e := newAdminEnv(t)
	id := mustCreateReport(t, e)

	rec := e.do(t, "PATCH", "/api/admin/reports/"+itoa(id), e.adminTok, `{"status":"resolved","resolution":"fixed the answer"}`)
	if rec.Code != 200 {
		t.Fatalf("moderate: got %d (%s)", rec.Code, rec.Body.String())
	}
	var out struct {
		Report storage.QuestionReport `json:"report"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if out.Report.Status != "resolved" || out.Report.ResolvedBy != e.adminID || out.Report.Resolution != "fixed the answer" {
		t.Errorf("report not stamped correctly: %+v", out.Report)
	}

	// The decision is in the trail, attributed to the actor.
	detail := e.do(t, "GET", "/api/admin/reports/"+itoa(id), e.adminTok, "")
	if detail.Code != 200 {
		t.Fatalf("report detail: %d", detail.Code)
	}
	var d struct {
		History []storage.AdminAuditRecord `json:"history"`
	}
	_ = json.Unmarshal(detail.Body.Bytes(), &d)
	if len(d.History) != 1 || d.History[0].ActorID != e.adminID || d.History[0].Action != storage.ActionReportStatusChange {
		t.Errorf("moderation history = %+v", d.History)
	}
}

func TestResolvedReportRequiresResolution(t *testing.T) {
	e := newAdminEnv(t)
	id := mustCreateReport(t, e)
	if rec := e.do(t, "PATCH", "/api/admin/reports/"+itoa(id), e.adminTok, `{"status":"resolved"}`); rec.Code != 400 {
		t.Errorf("resolve without a reason: got %d, want 400", rec.Code)
	}
}

func TestModeratorCanModerateButNotResolvePrivilege(t *testing.T) {
	e := newAdminEnv(t)
	if err := e.store.SetStudentRole(e.studentID, storage.RoleModerator); err != nil {
		t.Fatalf("moderator: %v", err)
	}
	id := mustCreateReport(t, e)
	if rec := e.do(t, "PATCH", "/api/admin/reports/"+itoa(id), e.studentTok, `{"status":"dismissed","resolution":"not a bug"}`); rec.Code != 200 {
		t.Errorf("moderator could not moderate: %d (%s)", rec.Code, rec.Body.String())
	}
	// But still cannot reach the audit trail or account management.
	if rec := e.do(t, "GET", "/api/admin/audit", e.studentTok, ""); rec.Code != 403 {
		t.Errorf("moderator read the audit trail: %d", rec.Code)
	}
}

// ── Audit ─────────────────────────────────────────────────────────────────

func TestAuditTrailIsReadOnlyAndForbiddenToNonStaff(t *testing.T) {
	e := newAdminEnv(t)
	if rec := e.do(t, "GET", "/api/admin/audit", e.studentTok, ""); rec.Code != 403 {
		t.Errorf("learner read the audit trail: %d", rec.Code)
	}
	for _, m := range []string{"POST", "PATCH", "PUT", "DELETE"} {
		if rec := e.do(t, m, "/api/admin/audit", e.adminTok, `{"action":"FORGED"}`); rec.Code != 405 {
			t.Errorf("%s /api/admin/audit: got %d, want 405", m, rec.Code)
		}
	}
}

// ── /api/admin/me ─────────────────────────────────────────────────────────

func TestAdminMeReportsRoleAndPermissions(t *testing.T) {
	e := newAdminEnv(t)
	if err := e.store.SetStudentRole(e.studentID, storage.RoleModerator); err != nil {
		t.Fatalf("moderator: %v", err)
	}
	rec := e.do(t, "GET", "/api/admin/me", e.studentTok, "")
	if rec.Code != 200 {
		t.Fatalf("me: %d", rec.Code)
	}
	var out struct {
		Role        string   `json:"role"`
		Permissions []string `json:"permissions"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	if out.Role != "moderator" {
		t.Errorf("role = %q, want moderator", out.Role)
	}
	joined := strings.Join(out.Permissions, ",")
	if !strings.Contains(joined, "reports.moderate") {
		t.Errorf("moderator missing reports.moderate: %v", out.Permissions)
	}
	if strings.Contains(joined, "admins.manage") {
		t.Errorf("moderator was granted admins.manage: %v", out.Permissions)
	}
}

func itoa(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
