package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// The last-admin rule.
//
// The interesting case is the one no sequential test can reach: two administrators demoting
// each other at the same instant. A read of COUNT(*) followed by an UPDATE loses it — both
// read two, both write, and the system reaches zero through the very feature meant to prevent
// it. That is why the count lives inside the UPDATE statement. These tests use a file-backed
// store for the concurrent case, and the reason is not incidental; see the note on
// `:memory:` in NewSQLiteStore.

func newAdminStore(t *testing.T) (*SQLiteStore, []string) {
	t.Helper()
	s, err := NewSQLiteStore(filepath.Join(t.TempDir(), "admin.db"))
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	ids := make([]string, 0, 2)
	for _, n := range []string{"A", "B"} {
		st, err := s.CreateStudent(n)
		if err != nil {
			t.Fatalf("create %s: %v", n, err)
		}
		ids = append(ids, st.ID)
	}
	return s, ids
}

func promote(t *testing.T, s *SQLiteStore, ids ...string) {
	t.Helper()
	for _, id := range ids {
		if err := s.SetStudentRole(id, RoleAdmin); err != nil {
			t.Fatalf("promote %s: %v", id, err)
		}
	}
}

// A default is the migration's real promise: an account created before the role column existed
// is a learner, and an account created after it with no role supplied is a learner.
func TestRoleDefaultsToStudent(t *testing.T) {
	s, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	defer s.Close()
	st, err := s.CreateStudent("New")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	role, err := s.GetStudentRole(st.ID)
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	if role != RoleStudent {
		t.Errorf("new account role = %q, want student", role)
	}
}

// Reading a role must never fail the request it is on. A missing row is "no privilege", not an
// error: the call sits in the authentication path of every admin request, and a deleted id in
// the last second should not take down the middleware.
func TestGetStudentRoleOnMissingAccountIsStudent(t *testing.T) {
	s, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	defer s.Close()
	role, err := s.GetStudentRole("no-such-id")
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if role != RoleStudent {
		t.Errorf("role = %q, want student", role)
	}
}

// A row written by a future version with a role this build has never heard of must read as the
// least privilege, not as something unknown-but-powerful. `NormalizeRole` is the whole reason.
func TestNormalizeRoleFailsClosed(t *testing.T) {
	for _, in := range []string{"", "root", "superuser", "Admin", "ADMIN", "teacher", "admin ", " null"} {
		if got := NormalizeRole(in); got != RoleStudent {
			t.Errorf("NormalizeRole(%q) = %q, want student", in, got)
		}
	}
	if got := NormalizeRole("admin"); got != RoleAdmin {
		t.Errorf(`NormalizeRole("admin") = %q, want admin`, got)
	}
	// ParseRole is stricter on purpose: assignment goes through it, so a typo is refused rather
	// than silently becoming "student" while the caller believes it promoted someone.
	for _, bad := range []string{"", "Admin", "root", "admin "} {
		if _, ok := ParseRole(bad); ok {
			t.Errorf("ParseRole(%q) accepted", bad)
		}
	}
	for _, good := range []string{"admin", "student"} {
		if _, ok := ParseRole(good); !ok {
			t.Errorf("ParseRole(%q) rejected", good)
		}
	}
}

func TestSetStudentRoleOnMissingAccountReportsNoRows(t *testing.T) {
	s, err := NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	defer s.Close()
	if err := s.SetStudentRole("no-such-id", RoleAdmin); err == nil {
		t.Error("promoting a missing account reported success")
	}
}

func TestDemoteLastAdminIsRefused(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0])
	if n, _ := s.CountAdmins(); n != 1 {
		t.Fatalf("admins = %d, want 1", n)
	}
	err := s.DemoteAdminSafely(ids[0])
	if !errors.Is(err, ErrLastAdmin) {
		t.Fatalf("err = %v, want ErrLastAdmin", err)
	}
	if n, _ := s.CountAdmins(); n != 1 {
		t.Errorf("the refused demotion applied: admins = %d", n)
	}
}

func TestDemoteAnotherAdminIsAllowed(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0], ids[1])
	if err := s.DemoteAdminSafely(ids[1]); err != nil {
		t.Fatalf("demote: %v", err)
	}
	if n, _ := s.CountAdmins(); n != 1 {
		t.Errorf("admins = %d, want 1", n)
	}
	role, _ := s.GetStudentRole(ids[1])
	if role != RoleStudent {
		t.Errorf("role = %q, want student", role)
	}
}

// Already a learner is not an error and not ErrLastAdmin: the caller needs to tell "nothing to
// do" from "refused", and collapsing them would report a rule violation where none happened.
func TestDemoteAStudentReportsNoRows(t *testing.T) {
	s, ids := newAdminStore(t)
	err := s.DemoteAdminSafely(ids[0])
	if err == nil {
		t.Fatal("demoting a learner reported success")
	}
	if errors.Is(err, ErrLastAdmin) {
		t.Error("demoting a learner reported ErrLastAdmin")
	}
}

// The race the rule exists to prevent. Repeated because a lost update is a timing artefact and
// one run proves nothing; with a read-then-write implementation this fails intermittently,
// which is worse than failing loudly.
func TestConcurrentDemotionsCannotReachZeroAdmins(t *testing.T) {
	for attempt := 0; attempt < 20; attempt++ {
		s, ids := newAdminStore(t)
		promote(t, s, ids[0], ids[1])

		var wg sync.WaitGroup
		errs := make([]error, 2)
		wg.Add(2)
		go func() { defer wg.Done(); errs[0] = s.DemoteAdminSafely(ids[0]) }()
		go func() { defer wg.Done(); errs[1] = s.DemoteAdminSafely(ids[1]) }()
		wg.Wait()

		n, err := s.CountAdmins()
		if err != nil {
			t.Fatalf("count: %v", err)
		}
		if n < 1 {
			t.Fatalf("attempt %d: admins = %d — the system reached zero administrators", attempt, n)
		}
		refusals := 0
		for _, e := range errs {
			if errors.Is(e, ErrLastAdmin) {
				refusals++
			} else if e != nil {
				t.Fatalf("attempt %d: unexpected error %v", attempt, e)
			}
		}
		if refusals != 1 {
			t.Fatalf("attempt %d: %d refusals, want exactly 1 (admins=%d)", attempt, refusals, n)
		}
	}
}

// Promotions may race freely: the rule constrains demotions, and two administrators promoting
// each other is not a way to zero anything.
func TestConcurrentPromotionsAreAllApplied(t *testing.T) {
	s, ids := newAdminStore(t)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); _ = s.SetStudentRole(ids[0], RoleAdmin) }()
	go func() { defer wg.Done(); _ = s.SetStudentRole(ids[1], RoleAdmin) }()
	wg.Wait()
	if n, _ := s.CountAdmins(); n != 2 {
		t.Errorf("admins = %d, want 2", n)
	}
}

// ── The user projection ────────────────────────────────────────────────────

// AdminUser cannot express a credential, which is what makes "the admin list leaks nothing
// sensitive" a property of a type rather than a promise about a query. The projection is
// asserted field by field so adding a field to it is a deliberate act.
func TestAdminUserProjectionHasNoCredentialFields(t *testing.T) {
	s, _ := NewSQLiteStore(":memory:")
	defer s.Close()
	if err := s.SetPasswordHash((&Student{ID: "x"}).ID, ""); err == nil {
		t.Log("hash set on a nonexistent row is tolerated")
	}
	st, err := s.CreateUser("Ada", "ada", "$argon2id$v=19$m=65536,t=3$c2FsdA$aGFzaA")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := s.SetShareToken(st.ID, "share-secret"); err != nil {
		t.Fatalf("set share token: %v", err)
	}
	if err := s.AddXP(st.ID, 4200); err != nil {
		t.Fatalf("add xp: %v", err)
	}

	rows, err := s.SearchAdminUsers("ada", 10, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	got := rows[0]
	if !got.HasPassword {
		t.Error("HasPassword = false for an account with a password")
	}
	// The projection has no field that could carry these, and the marshalled form is checked
	// too because that is what actually crosses the wire.
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, secret := range []string{"argon2id", "$argon2", "share-secret", "4200"} {
		if strings.Contains(string(encoded), secret) {
			t.Errorf("projection leaked %q: %s", secret, encoded)
		}
	}
}

// LIKE metacharacters in a search term are the difference between a search and a table dump.
func TestSearchAdminUsersEscapesWildcards(t *testing.T) {
	s, _ := NewSQLiteStore(":memory:")
	defer s.Close()
	if _, err := s.CreateUser("Ada", "ada", "hash"); err != nil {
		t.Fatalf("create: %v", err)
	}
	rows, err := s.SearchAdminUsers("%", 10, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("a bare %% matched %d users", len(rows))
	}
	rows, err = s.SearchAdminUsers("_", 10, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("a bare _ matched %d users", len(rows))
	}
}

func TestSearchAdminUsersEmptyQueryListsAll(t *testing.T) {
	s, _ := NewSQLiteStore(":memory:")
	defer s.Close()
	for i, n := range []string{"A", "B", "C"} {
		if _, err := s.CreateUser(n, string(rune('a'+i)), "hash"); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	rows, err := s.SearchAdminUsers("  ", 10, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("rows = %d, want 3", len(rows))
	}
}

// The audit log is append-only. There is no update and no delete method on either store, so the
// table cannot be rewritten through the repository — and these tests say so where a future
// method would be added.
func TestAdminAuditIsAppendOnly(t *testing.T) {
	s, _ := NewSQLiteStore(":memory:")
	defer s.Close()
	for i := 0; i < 3; i++ {
		if _, err := s.RecordAdminAudit(AdminAuditRecord{
			ActorID:    "actor",
			Action:     ActionAdminPromoteUser,
			EntityType: "user",
			EntityID:   "target",
			BeforeJSON: `{"role":"student"}`,
			AfterJSON:  `{"role":"admin"}`,
		}); err != nil {
			t.Fatalf("record: %v", err)
		}
	}
	rows, err := s.ListAdminAudit(10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows = %d, want 3", len(rows))
	}
	// Newest first, which is the only order an audit page can usefully render.
	if rows[0].ID <= rows[2].ID {
		t.Errorf("rows are not newest-first: %d then %d", rows[0].ID, rows[2].ID)
	}
	if rows[0].CreatedAt == "" {
		t.Error("event has no timestamp")
	}
	// Paging is bounded by the caller, and an offset past the end is empty rather than an error.
	empty, err := s.ListAdminAudit(10, 100)
	if err != nil {
		t.Fatalf("paged list: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("paged past the end returned %d rows", len(empty))
	}
}

func TestCountHelpers(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0])
	if n, err := s.CountAdminUsers(); err != nil || n != 2 {
		t.Errorf("users = %d (err %v), want 2", n, err)
	}
	if n, err := s.CountAdmins(); err != nil || n != 1 {
		t.Errorf("admins = %d (err %v), want 1", n, err)
	}
	if n, err := s.CountStoredQuestions(); err != nil || n != 0 {
		t.Errorf("questions = %d (err %v), want 0", n, err)
	}
}

// ── The mutation and its event are one unit of work ────────────────────────

// blockAuditInserts makes every write to admin_audit fail, without touching anything else.
//
// A trigger rather than a fake repository, because the property under test is about what the
// *database* does when the insert fails. Substituting a stub would let the code pass a mock and
// still leave a real deployment with an unaudited role change.
func blockAuditInserts(t *testing.T, s *SQLiteStore) {
	t.Helper()
	if _, err := s.db.Exec(`CREATE TRIGGER block_admin_audit BEFORE INSERT ON admin_audit
		BEGIN SELECT RAISE(ABORT, 'audit blocked'); END;`); err != nil {
		t.Fatalf("install trigger: %v", err)
	}
	t.Cleanup(func() {
		_, _ = s.db.Exec(`DROP TRIGGER IF EXISTS block_admin_audit`)
	})
}

func roleOf(t *testing.T, s *SQLiteStore, id string) Role {
	t.Helper()
	r, err := s.GetStudentRole(id)
	if err != nil {
		t.Fatalf("role: %v", err)
	}
	return r
}

// The whole point of the transaction. A failed audit insert must not leave the role changed:
// otherwise the system holds a privileged change that no record accounts for, which is the
// failure an audit trail exists to make impossible.
func TestFailedAuditInsertRollsBackPromotion(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0]) // one administrator, so the promotion below is not a no-op
	blockAuditInserts(t, s)

	if _, err := s.SetRoleAudited(ids[0], ids[1], RoleAdmin); err == nil {
		t.Fatal("promotion succeeded despite the audit insert failing")
	}

	if got := roleOf(t, s, ids[1]); got != RoleStudent {
		t.Errorf("role = %q after a failed audit insert, want student: the mutation was not rolled back", got)
	}
	if n, _ := s.CountAdmins(); n != 1 {
		t.Errorf("admins = %d, want 1", n)
	}
	events, err := s.ListAdminAudit(10, 0)
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("audit log holds %d events; the blocked insert should have written none", len(events))
	}
}

// The same guarantee on the demotion path, which is the one carrying the last-admin guard — so
// this also confirms the guard is evaluated inside the transaction and does not survive a
// rollback as a partial state.
func TestFailedAuditInsertRollsBackDemotion(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0], ids[1])
	blockAuditInserts(t, s)

	if _, err := s.SetRoleAudited(ids[0], ids[1], RoleStudent); err == nil {
		t.Fatal("demotion succeeded despite the audit insert failing")
	}

	if got := roleOf(t, s, ids[1]); got != RoleAdmin {
		t.Errorf("role = %q after a failed audit insert, want admin: the mutation was not rolled back", got)
	}
	if n, _ := s.CountAdmins(); n != 2 {
		t.Errorf("admins = %d, want 2", n)
	}
	events, err := s.ListAdminAudit(10, 0)
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("audit log holds %d events after a rolled-back demotion", len(events))
	}
}

// A refused demotion writes nothing and leaves the role alone, and — the part worth asserting
// separately — a refusal is not an audit failure. If the last-owner guard fired *after* a
// successful insert, the transaction would have to roll the insert back too.
//
// The invariant is on the owner in V2: an admin can be demoted freely (another owner or, in
// principle, the system itself can mint another), but emptying the owner role is unrecoverable
// because nothing but an owner can create one.
func TestLastOwnerRefusalWritesNoEventAndKeepsTheRole(t *testing.T) {
	s, ids := newAdminStore(t)
	if err := s.SetStudentRole(ids[0], RoleOwner); err != nil {
		t.Fatalf("promote to owner: %v", err)
	}

	_, err := s.SetRoleAudited(ids[0], ids[0], RoleStudent)
	if !errors.Is(err, ErrLastOwner) {
		t.Fatalf("err = %v, want ErrLastOwner", err)
	}
	if got := roleOf(t, s, ids[0]); got != RoleOwner {
		t.Errorf("role = %q, want owner", got)
	}
	events, _ := s.ListAdminAudit(10, 0)
	if len(events) != 0 {
		t.Errorf("a refused demotion wrote %d events", len(events))
	}
}

// The happy path still writes both, in one commit. Without this the tests above could pass
// because nothing is ever written.
func TestSetRoleAuditedWritesBothRoleAndEvent(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0])

	res, err := s.SetRoleAudited(ids[0], ids[1], RoleAdmin)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}
	if !res.Changed {
		t.Error("Changed = false for a real promotion")
	}
	if res.Before.Role != RoleStudent || res.After.Role != RoleAdmin {
		t.Errorf("snapshots = %q → %q, want student → admin", res.Before.Role, res.After.Role)
	}
	events, _ := s.ListAdminAudit(10, 0)
	if len(events) != 1 {
		t.Fatalf("events = %d, want 1", len(events))
	}
	e := events[0]
	if e.Action != ActionAdminPromoteUser {
		t.Errorf("action = %q, want %q", e.Action, ActionAdminPromoteUser)
	}
	if e.ActorID != ids[0] {
		t.Errorf("actor = %q, want %q", e.ActorID, ids[0])
	}
	if e.EntityID != ids[1] || e.EntityType != "user" {
		t.Errorf("entity = %q/%q", e.EntityType, e.EntityID)
	}
	if !strings.Contains(e.BeforeJSON, `"role":"student"`) || !strings.Contains(e.AfterJSON, `"role":"admin"`) {
		t.Errorf("snapshots = %q → %q", e.BeforeJSON, e.AfterJSON)
	}
}

// A no-op reports itself and writes nothing, and — the part that matters now that the write is
// transactional — it must not be mistaken for a failed change.
func TestSetRoleAuditedNoOpReportsItself(t *testing.T) {
	s, ids := newAdminStore(t)
	promote(t, s, ids[0])

	res, err := s.SetRoleAudited(ids[0], ids[0], RoleAdmin)
	if !errors.Is(err, ErrRoleUnchanged) {
		t.Fatalf("err = %v, want ErrRoleUnchanged", err)
	}
	if res.Changed {
		t.Error("Changed = true for a no-op")
	}
	// The snapshots are still returned, so the caller can answer the request without re-reading.
	if res.After.Role != RoleAdmin {
		t.Errorf("After.Role = %q", res.After.Role)
	}
	events, _ := s.ListAdminAudit(10, 0)
	if len(events) != 0 {
		t.Errorf("a no-op wrote %d events", len(events))
	}
}

func TestSetRoleAuditedUnknownAccount(t *testing.T) {
	s, _ := newAdminStore(t)
	_, err := s.SetRoleAudited("actor", "no-such-id", RoleAdmin)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("err = %v, want sql.ErrNoRows", err)
	}
	events, _ := s.ListAdminAudit(10, 0)
	if len(events) != 0 {
		t.Errorf("an unknown account produced %d events", len(events))
	}
}
