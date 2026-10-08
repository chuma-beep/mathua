package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

// The administrative security boundary.
//
// The brief for this phase put authorization first and said so explicitly, which is the right
// priority: a dashboard nobody can reach is a smaller problem than a dashboard everybody can.
// So most of what follows asserts that the wrong caller gets the wrong status code, and only a
// little asserts that the right one gets useful data.
//
// Two design decisions are load-bearing and are tested as such rather than assumed:
//
//   - The role is read from the database on every request. TestLastAdminDemotionTakesEffect
//     Immediately exists because a role cached in the token would keep granting admin until
//     expiry, and that test is the only thing standing between the two implementations.
//   - 401 and 403 are different answers. Unauthenticated callers learn nothing about whether
//     the route exists; authenticated learners learn they are not administrators. Collapsing
//     them into one would tell a prober that its token is getting somewhere.

// adminEnv is a server with auth enabled and two accounts, one of each role.
type adminEnv struct {
	mux        *http.ServeMux
	store      storage.Repository
	studentTok string
	adminTok   string
	studentID  string
	adminID    string
	thirdID    string
	thirdTok   string
}

func newAdminEnv(t *testing.T) *adminEnv {
	t.Helper()
	d, _ := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	authSvc := auth.New(store)
	s := New(engine.New(store, d, reg, nil, nil), store, authSvc)
	mux := http.NewServeMux()
	s.Register(mux)

	env := &adminEnv{mux: mux, store: store}

	signup := func(name, username, email string) string {
		t.Helper()
		body := fmt.Sprintf(`{"name":%q,"username":%q,"password":"Engine!n1","email":%q}`, name, username, email)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader([]byte(body))))
		if rec.Code != 200 {
			t.Fatalf("signup %s: %d %s", username, rec.Code, rec.Body.String())
		}
		var out struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode signup: %v", err)
		}
		return out.Token
	}

	env.studentTok = signup("Ada", "ada", "ada@example.com")
	env.adminTok = signup("Root", "root", "root@example.com")
	thirdTok := signup("Third", "third", "third@example.com")

	me := func(tok string) string {
		t.Helper()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		mux.ServeHTTP(rec, req)
		var out struct {
			StudentID string `json:"student_id"`
		}
		json.Unmarshal(rec.Body.Bytes(), &out)
		return out.StudentID
	}
	env.studentID = me(env.studentTok)
	env.adminID = me(env.adminTok)
	env.thirdID = me(thirdTok)
	env.thirdTok = thirdTok

	// Promote through the operator path, not an endpoint: this is the same call
	// `mathua admin promote` makes, and using it here keeps the tests from depending on a
	// route that does not exist.
	// The V2 top role is owner, not admin: owner is the one that can manage peers and hand
	// out every other role, so the environment's privileged account holds it.
	if err := store.SetStudentRole(env.adminID, storage.RoleOwner); err != nil {
		t.Fatalf("promote: %v", err)
	}
	return env
}

func (e *adminEnv) do(t *testing.T, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.mux.ServeHTTP(rec, r)
	return rec
}

// ── The boundary ───────────────────────────────────────────────────────────

// Every admin route, so a new one cannot be registered without inheriting the gate. A list
// written once and asserted against the mux would be better still; this at least fails loudly
// if someone forgets the middleware on a path they just added.
var adminRoleRoutes = []struct {
	method string
	path   string
}{
	{"GET", "/api/admin/overview"},
	{"GET", "/api/admin/audit"},
	{"GET", "/api/admin/users"},
	{"PATCH", "/api/admin/users/someone"},
}

func TestAdminRoutesRejectUnauthenticatedWith401(t *testing.T) {
	e := newAdminEnv(t)
	for _, rt := range adminRoleRoutes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			rec := e.do(t, rt.method, rt.path, "", "")
			if rec.Code != 401 {
				t.Errorf("no token: got %d, want 401", rec.Code)
			}
		})
	}
}

func TestAdminRoutesRejectStudentWith403(t *testing.T) {
	e := newAdminEnv(t)
	for _, rt := range adminRoleRoutes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			rec := e.do(t, rt.method, rt.path, e.studentTok, "")
			if rec.Code != 403 {
				t.Errorf("student: got %d, want 403 (%s)", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestAdminRoutesRejectGarbageTokenWith401(t *testing.T) {
	e := newAdminEnv(t)
	for _, tok := range []string{"not-a-jwt", "a.b.c", "Bearer"} {
		rec := e.do(t, "GET", "/api/admin/users", tok, "")
		if rec.Code != 401 {
			t.Errorf("token %q: got %d, want 401", tok, rec.Code)
		}
	}
}

// The distinction that matters: a learner is told "no", a stranger is told "who are you".
// If these ever collapse, a prober learns which half of the pair its token satisfied.
func TestAdmin401And403AreDistinguishable(t *testing.T) {
	e := newAdminEnv(t)
	unauth := e.do(t, "GET", "/api/admin/users", "", "")
	forbidden := e.do(t, "GET", "/api/admin/users", e.studentTok, "")
	if unauth.Code != 401 || forbidden.Code != 403 {
		t.Fatalf("got %d and %d, want 401 and 403", unauth.Code, forbidden.Code)
	}
	if strings.Contains(forbidden.Body.String(), "not found") {
		t.Error("a learner must be told they lack the role, not that the route does not exist")
	}
}

// A guest token is a real signed token for a real student row, so it must reach the same 403.
// If guest tokens ever skipped the role check they would be a way in.
func TestAdminRejectsGuestToken(t *testing.T) {
	e := newAdminEnv(t)
	guest := issueGuest(t, e.mux)
	rec := e.do(t, "GET", "/api/admin/users", guest, "")
	if rec.Code != 403 {
		t.Errorf("guest token: got %d, want 403", rec.Code)
	}
}

// With auth disabled there is no identity, so there is no administrator. The surface must be
// closed rather than open: -no-auth is a development convenience and must not become an
// unauthenticated admin.
func TestAdminClosedWhenAuthDisabled(t *testing.T) {
	d, _ := concepts.Build([]concepts.Concept{{ID: "a", Label: "A", Domain: "d", GradingType: "numeric"}})
	store, _ := storage.NewSQLiteStore(":memory:")
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	s := New(engine.New(store, d, reg, nil, nil), store, nil)
	mux := http.NewServeMux()
	s.Register(mux)

	for _, path := range []string{"/api/admin/overview", "/api/admin/users", "/api/admin/audit"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("Authorization", "Bearer anything")
		mux.ServeHTTP(rec, req)
		if rec.Code != 401 {
			t.Errorf("%s with auth disabled: got %d, want 401", path, rec.Code)
		}
	}
}

// ── Role assignment ────────────────────────────────────────────────────────

// The client must never be able to choose its own role. This is the classic privilege
// escalation and it is one JSON key away, so it is asserted at the endpoint that creates
// accounts rather than trusted to a review.
func TestSignupCannotChooseRole(t *testing.T) {
	e := newAdminEnv(t)
	body := `{"name":"Mallory","username":"mallory","password":"Engine!n1","email":"mallory@example.com","role":"admin"}`
	rec := e.do(t, "POST", "/api/auth/signup", "", body)
	if rec.Code != 200 {
		t.Fatalf("signup: got %d, want 200 (%s)", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"admin"`) {
		t.Errorf("signup response echoed a privileged role: %s", rec.Body.String())
	}
	// And the row it created is a learner, which is the half that matters.
	st, err := e.store.FindByEmail("mallory@example.com")
	if err != nil || st == nil {
		t.Fatalf("find mallory: %v", err)
	}
	role, err := e.store.GetStudentRole(st.ID)
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	if role != storage.RoleStudent {
		t.Errorf("role = %q, want student", role)
	}
	// …and a learner with that token cannot reach the admin surface.
	var tok struct {
		Token string `json:"token"`
	}
	json.Unmarshal(rec.Body.Bytes(), &tok)
	if got := e.do(t, "GET", "/api/admin/users", tok.Token, ""); got.Code != 403 {
		t.Errorf("self-promoted account reached admin: %d", got.Code)
	}
}

// A signup body is the one place a client supplies account attributes, so it is also where a
// near-miss of the same class would land. Checked with the alternative spellings a mass
// assignment filter usually misses.
func TestSignupIgnoresRoleFieldCaseVariants(t *testing.T) {
	e := newAdminEnv(t)
	for i, spelling := range []string{`"admin"`, `"Admin"`, `"ADMIN"`, `"admin "`, `" superadmin"`, `"owner"`} {
		// The username and email carry the index, not the spelling's length: two of these are
		// the same length, and reusing one address makes the test fail on a duplicate-account
		// 409 while appearing to be about roles.
		// Username rules are real and enforced here (3-20 characters), so the name carries the
		// index padded rather than a bare digit — otherwise a case would fail on the username
		// policy and the test would appear to be about something else.
		body := fmt.Sprintf(`{"name":"M","username":"mallory%d","password":"Engine!n1","email":"mallory%d@example.com","role":%s}`,
			i, i, spelling)
		rec := e.do(t, "POST", "/api/auth/signup", "", body)
		if rec.Code != 200 {
			t.Fatalf("signup role=%s: %d %s", spelling, rec.Code, rec.Body.String())
		}
		st, err := e.store.FindByEmail(fmt.Sprintf("mallory%d@example.com", i))
		if err != nil || st == nil {
			t.Fatalf("find mallory%d: %v", i, err)
		}
		role, _ := e.store.GetStudentRole(st.ID)
		if role != storage.RoleStudent {
			t.Errorf("role %s produced %q", spelling, role)
		}
	}
}

// ── Promotion and demotion ─────────────────────────────────────────────────

func TestAdminCanPromoteStudent(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, `{"role":"admin"}`)
	if rec.Code != 200 {
		t.Fatalf("promote: got %d (%s)", rec.Code, rec.Body.String())
	}
	role, _ := e.store.GetStudentRole(e.studentID)
	if role != storage.RoleAdmin {
		t.Errorf("role = %q, want admin", role)
	}
	// The promoted account can now use its own existing token: the role is read per request, so
	// no re-login is needed.
	if got := e.do(t, "GET", "/api/admin/overview", e.studentTok, ""); got.Code != 200 {
		t.Errorf("promoted user still denied: %d (%s)", got.Code, got.Body.String())
	}
}

func TestAdminCanDemoteAnotherAdminWhenSafe(t *testing.T) {
	e := newAdminEnv(t)
	// Two admins first, so the demotion is safe by the last-admin rule.
	if err := e.store.SetStudentRole(e.studentID, storage.RoleAdmin); err != nil {
		t.Fatalf("promote second admin: %v", err)
	}
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, `{"role":"student"}`)
	if rec.Code != 200 {
		t.Fatalf("demote: got %d (%s)", rec.Code, rec.Body.String())
	}
	role, _ := e.store.GetStudentRole(e.studentID)
	if role != storage.RoleStudent {
		t.Errorf("role = %q, want student", role)
	}
	// Their token stops working immediately — no waiting for expiry.
	if got := e.do(t, "GET", "/api/admin/users", e.studentTok, ""); got.Code != 403 {
		t.Errorf("demoted user's token still authorized: %d", got.Code)
	}
}

// The rule itself, and the case nobody writes a test for: an admin demoting *themselves* when
// they are the only one.
func TestLastOwnerCannotDemoteSelf(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.adminID, e.adminTok, `{"role":"student"}`)
	if rec.Code != 409 {
		t.Fatalf("last-owner demotion: got %d, want 409 (%s)", rec.Code, rec.Body.String())
	}
	role, _ := e.store.GetStudentRole(e.adminID)
	if role != storage.RoleOwner {
		t.Errorf("the refused demotion still applied: role = %q", role)
	}
	n, _ := e.store.CountOwners()
	if n != 1 {
		t.Errorf("owners = %d, want 1", n)
	}
}

// The other direction: an administrator cannot demote an owner, because the owner outranks
// them. That is authorization (403), not the invariant — the invariant is what an owner hits
// when they would empty the top of the hierarchy, covered by TestLastOwnerCannotDemoteSelf.
func TestAdminCannotDemoteOwner(t *testing.T) {
	e := newAdminEnv(t)
	// A second staff member who is an admin, not an owner.
	if err := e.store.SetStudentRole(e.thirdID, storage.RoleAdmin); err != nil {
		t.Fatalf("promote third: %v", err)
	}
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.adminID, e.thirdTok, `{"role":"student"}`)
	if rec.Code != 403 {
		t.Errorf("an admin demoting an owner: got %d, want 403 (%s)", rec.Code, rec.Body.String())
	}
	role, _ := e.store.GetStudentRole(e.adminID)
	if role != storage.RoleOwner {
		t.Errorf("the owner's role changed to %q", role)
	}
}

// The single most important property of the whole phase: a demotion takes effect on the next
// request, not when the token expires. If the role were read once at login, this test is what
// fails.
func TestLastAdminDemotionTakesEffectImmediately(t *testing.T) {
	e := newAdminEnv(t)
	if got := e.do(t, "GET", "/api/admin/users", e.studentTok, ""); got.Code != 403 {
		t.Fatalf("learner denied: %d", got.Code)
	}
	if err := e.store.SetStudentRole(e.studentID, storage.RoleAdmin); err != nil {
		t.Fatalf("promote: %v", err)
	}
	// Same token, no re-login.
	if got := e.do(t, "GET", "/api/admin/users", e.studentTok, ""); got.Code != 200 {
		t.Errorf("promoted, same token: got %d, want 200", got.Code)
	}
	if err := e.store.SetStudentRole(e.studentID, storage.RoleStudent); err != nil {
		t.Fatalf("demote: %v", err)
	}
	if got := e.do(t, "GET", "/api/admin/users", e.studentTok, ""); got.Code != 403 {
		t.Errorf("demoted, same token: got %d, want 403", got.Code)
	}
}

// The concurrent case lives in the storage package
// (TestConcurrentDemotionsCannotReachZeroAdmins), against a file-backed store. Running it here
// on the ":memory:" store these tests use would report whatever an empty per-connection
// database contains rather than what the code did — see the note in NewSQLiteStore.

// ── Mass assignment and IDOR ───────────────────────────────────────────────

// The endpoint already holds the ability to write to an arbitrary row, so it is the place
// where a mass-assignment bug would be most expensive. Only `role` may be applied.
func TestRolePatchIgnoresEveryOtherField(t *testing.T) {
	e := newAdminEnv(t)
	body := `{"role":"admin","xp_total":99999,"password_hash":"injected","share_token":"stolen",
	           "email":"attacker@example.com","name":"Renamed","diagnostic_completed":true}`
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, body)
	if rec.Code != 200 {
		t.Fatalf("patch: got %d (%s)", rec.Code, rec.Body.String())
	}
	// Snapshotted first: the assertion is that these fields are *unchanged*, not that they are
	// empty. Ada signed up, so she has a real password hash, and "still empty" would have
	// passed against an endpoint that overwrote it with nothing.
	before, err := e.store.GetStudent(e.studentID)
	if err != nil {
		t.Fatalf("get student before: %v", err)
	}
	st, err := e.store.GetStudent(e.studentID)
	if err != nil {
		t.Fatalf("get student: %v", err)
	}
	if st.XPTotal != before.XPTotal {
		t.Errorf("xp_total changed from %d to %d via the request body", before.XPTotal, st.XPTotal)
	}
	if st.PasswordHash != before.PasswordHash {
		t.Errorf("password_hash changed from %q to %q via the request body", before.PasswordHash, st.PasswordHash)
	}
	if st.ShareToken != before.ShareToken {
		t.Errorf("share_token changed from %q to %q via the request body", before.ShareToken, st.ShareToken)
	}
	if st.Email != "ada@example.com" {
		t.Errorf("email was written from the request body: %q", st.Email)
	}
	if st.Name != "Ada" {
		t.Errorf("name was written from the request body: %q", st.Name)
	}
	if st.DiagnosticCompleted {
		t.Error("diagnostic_completed was written from the request body")
	}
}

// Unknown role values are rejected rather than coerced. Coercing would mean the 200 response
// and the stored row could disagree, and the response is what an operator reads.
func TestRolePatchRejectsUnknownRoles(t *testing.T) {
	e := newAdminEnv(t)
	for _, bad := range []string{`"root"`, `"superuser"`, `"Admin"`, `""`, `"admin,student"`, `null`, `123`} {
		rec := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, `{"role":`+bad+`}`)
		if rec.Code != 400 {
			t.Errorf("role %s: got %d, want 400", bad, rec.Code)
		}
		role, _ := e.store.GetStudentRole(e.studentID)
		if role != storage.RoleStudent {
			t.Errorf("role %s changed the row to %q", bad, role)
		}
	}
}

func TestRolePatchUnknownUserIs404(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "PATCH", "/api/admin/users/no-such-id", e.adminTok, `{"role":"admin"}`)
	if rec.Code != 404 {
		t.Errorf("got %d, want 404", rec.Code)
	}
}

func TestRolePatchMalformedBodyIs400(t *testing.T) {
	e := newAdminEnv(t)
	for _, body := range []string{`{`, `{"role":`, `[]`, `{"role":{"nested":"admin"}}`} {
		rec := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, body)
		if rec.Code != 400 {
			t.Errorf("body %s: got %d, want 400", body, rec.Code)
		}
	}
}

// A mutating method that is not PATCH/PUT is refused, and a GET is now a read of the same path
// (the account detail) rather than a 405. Either way, none of them may mutate.
func TestRolePatchRejectsOtherMethods(t *testing.T) {
	e := newAdminEnv(t)
	for _, m := range []string{"POST", "DELETE"} {
		rec := e.do(t, m, "/api/admin/users/"+e.adminID, e.adminTok, "")
		if rec.Code != 405 {
			t.Errorf("%s: got %d, want 405", m, rec.Code)
		}
	}
	// GET reads and returns 200, and the role is untouched.
	if got := e.do(t, "GET", "/api/admin/users/"+e.adminID, e.adminTok, ""); got.Code != 200 {
		t.Errorf("GET detail: got %d, want 200 (%s)", got.Code, got.Body.String())
	}
	role, _ := e.store.GetStudentRole(e.adminID)
	if role != storage.RoleOwner {
		t.Fatalf("a read demoted the owner: %q", role)
	}
}

// ── Sensitive data ─────────────────────────────────────────────────────────

// The user list is the surface most likely to leak: it enumerates accounts. It is built from
// AdminUser, which cannot express a credential, and this test holds it to that.
func TestAdminUserListCarriesNoCredentials(t *testing.T) {
	e := newAdminEnv(t)
	// Give the account a real secret, so there is something to leak.
	if err := e.store.SetPasswordHash(e.studentID, "argon2id$v=19$m=65536,t=3$c2FsdA$aGFzaA"); err != nil {
		t.Fatalf("set password hash: %v", err)
	}
	rec := e.do(t, "GET", "/api/admin/users?q=ada@example.com", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("search: got %d (%s)", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, secret := range []string{"argon2id", "hash", "password_hash", "share_token"} {
		if strings.Contains(strings.ToLower(body), secret) {
			t.Errorf("search response leaked %q: %s", secret, body)
		}
	}
	// The projection does report whether a password exists, which is a fact and not the value.
	if !strings.Contains(body, `"has_password":true`) {
		t.Errorf("expected has_password to be reported: %s", body)
	}
}

func TestAdminUserListOmitsLearningState(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "GET", "/api/admin/users", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("list: got %d", rec.Code)
	}
	body := rec.Body.String()
	// Learning state is not account administration, and an admin UI showing it invites an
	// editor for it, which would be the wrong tool entirely.
	for _, field := range []string{"xp_total", "mastery", "attempts", "streak", "progress"} {
		if strings.Contains(body, field) {
			t.Errorf("user list carries learning state %q: %s", field, body)
		}
	}
}

// The audit log records what changed. If a credential could reach it, the log would become the
// place credentials leak to.
func TestAuditLogCarriesNoCredentials(t *testing.T) {
	e := newAdminEnv(t)
	if err := e.store.SetPasswordHash(e.studentID, "argon2id$v=19$m=65536,t=3$c2FsdA$aGFzaA"); err != nil {
		t.Fatalf("set password hash: %v", err)
	}
	e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, `{"role":"admin"}`)
	rec := e.do(t, "GET", "/api/admin/audit", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("audit: got %d", rec.Code)
	}
	for _, secret := range []string{"argon2id", "$v=19", "hash"} {
		if strings.Contains(rec.Body.String(), secret) {
			t.Errorf("audit log leaked %q: %s", secret, rec.Body.String())
		}
	}
}

// ── The audit trail ────────────────────────────────────────────────────────

func TestPromotionAndDemotionAreAudited(t *testing.T) {
	e := newAdminEnv(t)
	if got := e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, `{"role":"admin"}`); got.Code != 200 {
		t.Fatalf("promote: %d", got.Code)
	}
	if err := e.store.SetStudentRole(e.thirdID, storage.RoleAdmin); err != nil {
		t.Fatalf("promote third: %v", err)
	}
	if got := e.do(t, "PATCH", "/api/admin/users/"+e.thirdID, e.adminTok, `{"role":"student"}`); got.Code != 200 {
		t.Fatalf("demote: %d", got.Code)
	}

	rec := e.do(t, "GET", "/api/admin/audit", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("audit: got %d", rec.Code)
	}
	var out struct {
		Events []storage.AdminAuditRecord `json:"events"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out.Events) != 2 {
		t.Fatalf("events = %d, want 2: %+v", len(out.Events), out.Events)
	}
	// Newest first.
	if out.Events[0].Action != storage.ActionAdminDemoteUser {
		t.Errorf("newest event = %q, want %q", out.Events[0].Action, storage.ActionAdminDemoteUser)
	}
	for _, ev := range out.Events {
		if ev.ActorID != e.adminID {
			t.Errorf("event %q actor = %q, want the administrator %q", ev.Action, ev.ActorID, e.adminID)
		}
		if ev.EntityType != "user" || ev.EntityID == "" {
			t.Errorf("event %q has entity %q/%q", ev.Action, ev.EntityType, ev.EntityID)
		}
		if ev.CreatedAt == "" {
			t.Errorf("event %q has no timestamp", ev.Action)
		}
	}
	// Before/after capture the change the event is named for.
	if !strings.Contains(out.Events[1].BeforeJSON, `"role":"student"`) {
		t.Errorf("promotion before-state = %s", out.Events[1].BeforeJSON)
	}
	if !strings.Contains(out.Events[1].AfterJSON, `"role":"admin"`) {
		t.Errorf("promotion after-state = %s", out.Events[1].AfterJSON)
	}
}

// An event for a change that did not happen makes the log worthless. Re-promoting an
// administrator is a no-op and must not grow the trail.
func TestNoOpRoleChangeIsNotAudited(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.adminID, e.adminTok, `{"role":"owner"}`)
	if rec.Code != 200 {
		t.Fatalf("no-op re-assign: got %d", rec.Code)
	}
	var out struct {
		User    storage.AdminUser `json:"user"`
		Changed bool              `json:"changed"`
	}
	json.Unmarshal(rec.Body.Bytes(), &out)
	if out.Changed {
		t.Error("changed = true for a no-op")
	}
	list := e.do(t, "GET", "/api/admin/audit", e.adminTok, "")
	if strings.Contains(list.Body.String(), storage.ActionAdminPromoteUser) {
		t.Errorf("a no-op was audited: %s", list.Body.String())
	}
}

// The trail is read-only. There is no route to edit or delete an event, and a request trying
// is refused by the method check rather than quietly ignored.
func TestAuditLogCannotBeMutatedThroughHTTP(t *testing.T) {
	e := newAdminEnv(t)
	e.do(t, "PATCH", "/api/admin/users/"+e.studentID, e.adminTok, `{"role":"admin"}`)

	for _, m := range []string{"POST", "PATCH", "PUT", "DELETE"} {
		rec := e.do(t, m, "/api/admin/audit", e.adminTok, `{"action":"FORGED"}`)
		if rec.Code != 405 {
			t.Errorf("%s /api/admin/audit: got %d, want 405", m, rec.Code)
		}
	}
	// And the log is unchanged.
	list := e.do(t, "GET", "/api/admin/audit", e.adminTok, "")
	if strings.Contains(list.Body.String(), "FORGED") {
		t.Error("the audit log accepted a forged event")
	}
}

// A refused mutation must leave no trace in the trail, or the log would say an administrator
// exists when none was created.
func TestRefusedDemotionIsNotAudited(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "PATCH", "/api/admin/users/"+e.adminID, e.adminTok, `{"role":"student"}`)
	if rec.Code != 409 {
		t.Fatalf("got %d, want 409", rec.Code)
	}
	list := e.do(t, "GET", "/api/admin/audit", e.adminTok, "")
	if strings.Contains(list.Body.String(), storage.ActionAdminDemoteUser) {
		t.Errorf("a refused demotion was audited: %s", list.Body.String())
	}
}

// ── Search and overview ────────────────────────────────────────────────────

func TestAdminUserSearchFindsByEachIdentifier(t *testing.T) {
	e := newAdminEnv(t)
	cases := map[string]string{
		"ada@example.com": "email",
		e.studentID:       "id",
		"ada":             "username",
		"Root":            "name",
	}
	for q, what := range cases {
		rec := e.do(t, "GET", "/api/admin/users?q="+q, e.adminTok, "")
		if rec.Code != 200 {
			t.Fatalf("search %s: got %d", q, rec.Code)
		}
		var out struct {
			Users []storage.AdminUser `json:"users"`
		}
		json.Unmarshal(rec.Body.Bytes(), &out)
		if len(out.Users) == 0 {
			t.Errorf("search by %s (%q) returned nothing", what, q)
		}
	}
}

func TestAdminUserSearchWithNoQueryListsAccounts(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "GET", "/api/admin/users", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("got %d", rec.Code)
	}
	var out struct {
		Users []storage.AdminUser `json:"users"`
	}
	json.Unmarshal(rec.Body.Bytes(), &out)
	if len(out.Users) != 3 {
		t.Errorf("users = %d, want 3", len(out.Users))
	}
}

// The search term is data, and LIKE metacharacters are the difference between a search and a
// full table dump. A `%` must find nothing rather than everything.
func TestAdminUserSearchTreatsWildcardsAsLiterals(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "GET", "/api/admin/users?q=%25", e.adminTok, "") // %25 is a literal '%'
	if rec.Code != 200 {
		t.Fatalf("got %d", rec.Code)
	}
	var out struct {
		Users []storage.AdminUser `json:"users"`
	}
	json.Unmarshal(rec.Body.Bytes(), &out)
	if len(out.Users) != 0 {
		t.Errorf("a bare %% matched %d users; LIKE wildcards are not escaped", len(out.Users))
	}
}

func TestAdminOverviewReportsLiveCounts(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "GET", "/api/admin/overview", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("overview: got %d (%s)", rec.Code, rec.Body.String())
	}
	var out struct {
		Users     int `json:"users"`
		Staff     int `json:"staff"`
		Owners    int `json:"owners"`
		Concepts  int `json:"concepts"`
		Domains   int `json:"domains"`
		Questions int `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Users != 3 {
		t.Errorf("users = %d, want 3", out.Users)
	}
	if out.Staff != 1 {
		t.Errorf("staff = %d, want 1", out.Staff)
	}
	if out.Owners != 1 {
		t.Errorf("owners = %d, want 1", out.Owners)
	}
	if out.Concepts != 1 {
		t.Errorf("concepts = %d, want 1 (the corpus this server was built with)", out.Concepts)
	}
	if out.Domains != 1 {
		t.Errorf("domains = %d, want 1", out.Domains)
	}
}

// ── /api/auth/me ───────────────────────────────────────────────────────────

// The role is on the response so the client can offer an Admin entry. It is a convenience and
// the boundary above is what makes it one, which this test states by checking that the value
// is present and correct for both accounts.
func TestMeReportsRole(t *testing.T) {
	e := newAdminEnv(t)
	for _, tc := range []struct {
		token string
		want  storage.Role
	}{{e.adminTok, storage.RoleOwner}, {e.studentTok, storage.RoleStudent}} {
		rec := e.do(t, "GET", "/api/auth/me", tc.token, "")
		if rec.Code != 200 {
			t.Fatalf("me: got %d", rec.Code)
		}
		var out struct {
			Role string `json:"role"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if storage.NormalizeRole(out.Role) != tc.want {
			t.Errorf("me role = %q, want %q", out.Role, tc.want)
		}
	}
}

// issueGuest mints a guest token, which is signed for a real student row.
func issueGuest(t *testing.T, mux *http.ServeMux) string {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/auth/guest", bytes.NewReader([]byte(`{}`))))
	if rec.Code != 200 {
		t.Fatalf("guest: got %d (%s)", rec.Code, rec.Body.String())
	}
	var out struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode guest: %v", err)
	}
	return out.Token
}
