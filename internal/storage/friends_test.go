package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func friendStore(t *testing.T, repo Repository) FriendStore {
	t.Helper()
	store, ok := repo.(FriendStore)
	if !ok {
		t.Fatalf("store %T does not implement FriendStore", repo)
	}
	return store
}

func createFriendUser(t *testing.T, repo Repository, name string, verified bool) *Student {
	t.Helper()
	st, err := repo.CreateUser(name, strings.ToLower(name)+"_"+newUUID()[:8], "hash")
	if err != nil {
		t.Fatalf("create friend user %s: %v", name, err)
	}
	if verified {
		if err := repo.SetEmailVerified(st.ID, true); err != nil {
			t.Fatalf("verify friend user %s: %v", name, err)
		}
	}
	return st
}

func friendIdentity(t *testing.T, store FriendStore, studentID string) FriendIdentity {
	t.Helper()
	id, err := store.GetFriendIdentity(studentID)
	if err != nil {
		t.Fatalf("get friend identity: %v", err)
	}
	return id
}

func friendCount(t *testing.T, repo Repository, query string, args ...any) int {
	t.Helper()
	return seedInt(t, repo, query, args...)
}

func TestFriendIdentityIsCreatedAndRotatable(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	st := createFriendUser(t, repo, "ada", true)
	identity := friendIdentity(t, store, st.ID)

	if !strings.HasPrefix(identity.ProfileID, "p_") || len(identity.ProfileID) != 18 {
		t.Fatalf("profile id %q is not p_ plus 16 characters", identity.ProfileID)
	}
	if len(identity.FriendCode) != 8 {
		t.Fatalf("friend code %q is not eight characters", identity.FriendCode)
	}

	got, err := store.ResolveFriendTarget(identity.ProfileID)
	if err != nil || got != identity.ProfileID {
		t.Fatalf("resolve by profile id: %q, %v", got, err)
	}
	got, err = store.ResolveFriendTarget(strings.ToLower(identity.FriendCode))
	if err != nil || got != identity.ProfileID {
		t.Fatalf("resolve by lower-case friend code: %q, %v", got, err)
	}
	got, err = store.ResolveFriendTarget(strings.ToUpper(st.Username))
	if err != nil || got != identity.ProfileID {
		t.Fatalf("resolve by username: %q, %v", got, err)
	}
	if _, err := store.ResolveFriendTarget("missing-target"); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("unknown target error = %v, want ErrFriendNotFound", err)
	}

	rotated, err := store.RotateFriendCode(st.ID, time.Now().UTC())
	if err != nil {
		t.Fatalf("rotate friend code: %v", err)
	}
	if rotated == identity.FriendCode || len(rotated) != 8 {
		t.Fatalf("rotated code %q is invalid", rotated)
	}
	if _, err := store.ResolveFriendTarget(identity.FriendCode); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("old friend code still resolves: %v", err)
	}
	got, err = store.ResolveFriendTarget(rotated)
	if err != nil || got != identity.ProfileID {
		t.Fatalf("resolve by rotated code: %q, %v", got, err)
	}
}

func TestFriendRequestLifecycle(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "ada", true)
	b := createFriendUser(t, repo, "bob", true)
	ai := friendIdentity(t, store, a.ID)
	bi := friendIdentity(t, store, b.ID)

	if _, err := store.SendFriendRequest(a.ID, ai.ProfileID, time.Now().UTC()); !errors.Is(err, ErrFriendSelf) {
		t.Fatalf("self request error = %v, want ErrFriendSelf", err)
	}
	unverified := createFriendUser(t, repo, "carol", false)
	if _, err := store.SendFriendRequest(unverified.ID, ai.ProfileID, time.Now().UTC()); !errors.Is(err, ErrFriendNotVerified) {
		t.Fatalf("unverified request error = %v, want ErrFriendNotVerified", err)
	}

	f, err := store.SendFriendRequest(a.ID, bi.ProfileID, time.Now().UTC())
	if err != nil {
		t.Fatalf("send request: %v", err)
	}
	if f.Status != FriendshipPending || f.RequestedBy != a.ID {
		t.Fatalf("request = %+v, want pending from ada", f)
	}
	incoming, outgoing, err := store.ListFriendRequests(b.ID)
	if err != nil {
		t.Fatalf("list requests: %v", err)
	}
	if len(incoming) != 1 || incoming[0].StudentID != a.ID || incoming[0].Direction != FriendDirectionIncoming {
		t.Fatalf("bob incoming = %+v", incoming)
	}
	if len(outgoing) != 0 {
		t.Fatalf("bob outgoing = %+v, want empty", outgoing)
	}

	f, err = store.AcceptFriendRequest(b.ID, ai.ProfileID, time.Now().UTC())
	if err != nil {
		t.Fatalf("accept request: %v", err)
	}
	if f.Status != FriendshipAccepted || f.RequestedBy != a.ID || f.RespondedAt.IsZero() {
		t.Fatalf("accepted = %+v", f)
	}
	friends, err := store.ListFriends(a.ID)
	if err != nil || len(friends) != 1 || friends[0].StudentID != b.ID {
		t.Fatalf("ada friends = %+v, %v", friends, err)
	}
	friends, err = store.ListFriends(b.ID)
	if err != nil || len(friends) != 1 || friends[0].StudentID != a.ID {
		t.Fatalf("bob friends = %+v, %v", friends, err)
	}

	if err := store.RemoveFriend(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("remove friend: %v", err)
	}
	if friends, _ := store.ListFriends(a.ID); len(friends) != 0 {
		t.Fatalf("friends after removal = %+v", friends)
	}

	if _, err := store.SendFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("bob sends to ada: %v", err)
	}
	if err := store.DeclineFriendRequest(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("decline request: %v", err)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friendships"); n != 0 {
		t.Fatalf("friendship rows after decline = %d", n)
	}

	if _, err := store.SendFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("bob sends again: %v", err)
	}
	if err := store.CancelFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("cancel request: %v", err)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friendships"); n != 0 {
		t.Fatalf("friendship rows after cancel = %d", n)
	}
}

func TestFriendReciprocalRequestsConverge(t *testing.T) {
	repo := newFriendConcurrencyStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "ada", true)
	b := createFriendUser(t, repo, "bob", true)
	ai := friendIdentity(t, store, a.ID)
	bi := friendIdentity(t, store, b.ID)

	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_, err := store.SendFriendRequest(a.ID, bi.ProfileID, time.Now().UTC())
		errs <- err
	}()
	go func() {
		defer wg.Done()
		<-start
		_, err := store.SendFriendRequest(b.ID, ai.ProfileID, time.Now().UTC())
		errs <- err
	}()
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("reciprocal request failed: %v", err)
		}
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friendships WHERE status = 'accepted'"); n != 1 {
		t.Fatalf("accepted rows = %d, want 1", n)
	}
	if friends, err := store.ListFriends(a.ID); err != nil || len(friends) != 1 || friends[0].StudentID != b.ID {
		t.Fatalf("ada friends = %+v, %v", friends, err)
	}
	if friends, err := store.ListFriends(b.ID); err != nil || len(friends) != 1 || friends[0].StudentID != a.ID {
		t.Fatalf("bob friends = %+v, %v", friends, err)
	}
}

func TestFriendVisibilityAndBlocking(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "ada", true)
	b := createFriendUser(t, repo, "bob", true)
	ai := friendIdentity(t, store, a.ID)
	bi := friendIdentity(t, store, b.ID)
	if _, err := store.SendFriendRequest(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("send request: %v", err)
	}
	if _, err := store.AcceptFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("accept request: %v", err)
	}

	if err := repo.UpdateSettings(b.ID, `{"friends_visibility":"nobody"}`); err != nil {
		t.Fatalf("hide bob: %v", err)
	}
	if friends, _ := store.ListFriends(a.ID); len(friends) != 0 {
		t.Fatalf("ada sees hidden bob: %+v", friends)
	}
	if _, err := store.FriendProfile(a.ID, bi.ProfileID); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("hidden friend profile error = %v", err)
	}
	if _, err := store.FriendGraph(a.ID, bi.ProfileID); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("hidden friend graph error = %v", err)
	}
	if _, err := store.SendFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); !errors.Is(err, ErrFriendVisibilityOff) {
		t.Fatalf("request from hidden user error = %v", err)
	}
	if err := repo.UpdateSettings(b.ID, `{}`); err != nil {
		t.Fatalf("restore bob: %v", err)
	}
	if friends, _ := store.ListFriends(a.ID); len(friends) != 1 {
		t.Fatalf("ada friends after restore = %+v", friends)
	}

	if err := store.BlockUser(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("block bob: %v", err)
	}
	if friends, _ := store.ListFriends(a.ID); len(friends) != 0 {
		t.Fatalf("ada sees blocked bob: %+v", friends)
	}
	if _, err := store.SendFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); !errors.Is(err, ErrFriendBlocked) {
		t.Fatalf("blocked request error = %v", err)
	}
	if _, err := store.FriendProfile(a.ID, bi.ProfileID); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("blocked friend profile error = %v", err)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friendships"); n != 0 {
		t.Fatalf("friendship rows after block = %d", n)
	}

	if err := store.UnblockUser(a.ID, bi.ProfileID); err != nil {
		t.Fatalf("unblock bob: %v", err)
	}
	if _, err := store.SendFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("request after unblock: %v", err)
	}

	if err := repo.UpdateSettings(a.ID, `{"friends_visibility":"nobody"}`); err != nil {
		t.Fatalf("hide ada: %v", err)
	}
	if friends, _ := store.ListFriends(a.ID); len(friends) != 0 {
		t.Fatalf("hidden ada list = %+v", friends)
	}
	if _, err := store.FriendProfile(a.ID, bi.ProfileID); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("hidden caller profile error = %v", err)
	}
}

func TestFriendProfileAndGraphExposeOnlySafeFields(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "ada", true)
	b := createFriendUser(t, repo, "bob", true)
	ai := friendIdentity(t, store, a.ID)
	bi := friendIdentity(t, store, b.ID)
	if _, err := store.SendFriendRequest(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("send request: %v", err)
	}
	if _, err := store.AcceptFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("accept request: %v", err)
	}

	masteredAt := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
	execSeed(t, repo,
		"INSERT INTO concept_progress (student_id, concept_id, status, attempts, weakness_score, mastered_at) VALUES (?, ?, 'MASTERED', 9, 0.9, ?)",
		b.ID, "arith.add.single", masteredAt.Format(time.RFC3339),
	)
	execSeed(t, repo,
		"INSERT INTO concept_progress (student_id, concept_id, status, attempts, weakness_score) VALUES (?, ?, 'LEARNING', 3, 0.7)",
		b.ID, "arith.add.double",
	)

	profile, err := store.FriendProfile(a.ID, bi.ProfileID)
	if err != nil {
		t.Fatalf("friend profile: %v", err)
	}
	if profile.StudentID != b.ID || profile.Name != b.Name || profile.ProfileID != bi.ProfileID {
		t.Fatalf("profile = %+v", profile)
	}
	if profile.Mastered != 1 {
		t.Fatalf("mastered = %d, want 1", profile.Mastered)
	}
	if profile.FriendsSince.IsZero() {
		t.Fatal("friends since is zero")
	}

	graph, err := store.FriendGraph(a.ID, bi.ProfileID)
	if err != nil {
		t.Fatalf("friend graph: %v", err)
	}
	if len(graph) != 2 {
		t.Fatalf("graph nodes = %+v, want 2", graph)
	}
	if graph[0].ConceptID != "arith.add.double" || graph[0].Status != "LEARNING" {
		t.Fatalf("first node = %+v", graph[0])
	}
	if graph[1].ConceptID != "arith.add.single" || graph[1].Status != "MASTERED" || graph[1].MasteredAt == nil {
		t.Fatalf("second node = %+v", graph[1])
	}
	if !graph[1].MasteredAt.Equal(masteredAt) {
		t.Fatalf("mastered at = %v, want %v", graph[1].MasteredAt, masteredAt)
	}

	c := createFriendUser(t, repo, "carol", true)
	ci := friendIdentity(t, store, c.ID)
	if _, err := store.FriendProfile(c.ID, bi.ProfileID); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("non-friend profile error = %v", err)
	}
	if _, err := store.FriendProfile(a.ID, ci.ProfileID); !errors.Is(err, ErrFriendNotFound) {
		t.Fatalf("unknown friend profile error = %v", err)
	}
}

func TestFriendsLeaderboardWeeklyAllTimeAndPrivacy(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "ada", true)
	b := createFriendUser(t, repo, "bob", true)
	ai := friendIdentity(t, store, a.ID)
	bi := friendIdentity(t, store, b.ID)
	if _, err := store.SendFriendRequest(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("send request: %v", err)
	}
	if _, err := store.AcceptFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("accept request: %v", err)
	}
	if err := repo.AwardXP(a.ID, 30, "test", "friend-a-xp"); err != nil {
		t.Fatalf("award ada xp: %v", err)
	}
	if err := repo.AwardXP(b.ID, 10, "test", "friend-b-xp"); err != nil {
		t.Fatalf("award bob xp: %v", err)
	}

	week := time.Now().UTC().Add(-time.Hour)
	rows, err := store.FriendsLeaderboard(a.ID, "weekly", week)
	if err != nil {
		t.Fatalf("weekly leaderboard: %v", err)
	}
	if len(rows) != 2 || rows[0].StudentID != a.ID || rows[0].XP != 30 || rows[1].StudentID != b.ID || rows[1].XP != 10 {
		t.Fatalf("weekly rows = %+v", rows)
	}
	if !rows[0].IsMe || rows[1].IsMe {
		t.Fatalf("me flags = %+v", rows)
	}

	rows, err = store.FriendsLeaderboard(a.ID, "all_time", week)
	if err != nil {
		t.Fatalf("all-time leaderboard: %v", err)
	}
	if len(rows) != 2 || rows[0].StudentID != a.ID || rows[0].XP != 30 || rows[1].XP != 10 {
		t.Fatalf("all-time rows = %+v", rows)
	}

	if err := repo.UpdateSettings(b.ID, `{"leaderboard_privacy":"hidden"}`); err != nil {
		t.Fatalf("hide bob from leaderboard: %v", err)
	}
	rows, err = store.FriendsLeaderboard(a.ID, "weekly", week)
	if err != nil {
		t.Fatalf("hidden leaderboard: %v", err)
	}
	if len(rows) != 1 || rows[0].StudentID != a.ID {
		t.Fatalf("hidden rows = %+v", rows)
	}

	if err := repo.UpdateSettings(a.ID, `{"friends_visibility":"nobody"}`); err != nil {
		t.Fatalf("hide ada: %v", err)
	}
	if rows, err := store.FriendsLeaderboard(a.ID, "weekly", week); err != nil || len(rows) != 0 {
		t.Fatalf("hidden caller rows = %+v, %v", rows, err)
	}
}

func TestFriendAccountDeletionAndReset(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "ada", true)
	b := createFriendUser(t, repo, "bob", true)
	c := createFriendUser(t, repo, "carol", true)
	ai := friendIdentity(t, store, a.ID)
	bi := friendIdentity(t, store, b.ID)
	if _, err := store.SendFriendRequest(a.ID, bi.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("send request: %v", err)
	}
	if _, err := store.AcceptFriendRequest(b.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("accept request: %v", err)
	}
	if err := store.BlockUser(c.ID, ai.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("block ada: %v", err)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friend_blocks"); n != 1 {
		t.Fatalf("block rows = %d, want 1", n)
	}
	if err := repo.DeleteAccount(b.ID); err != nil {
		t.Fatalf("delete bob: %v", err)
	}
	if err := repo.DeleteAccount(c.ID); err != nil {
		t.Fatalf("delete carol: %v", err)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friendships"); n != 0 {
		t.Fatalf("friendship rows after deletion = %d", n)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friend_blocks"); n != 0 {
		t.Fatalf("block rows after deletion = %d", n)
	}
	if friends, _ := store.ListFriends(a.ID); len(friends) != 0 {
		t.Fatalf("ada friends after deletion = %+v", friends)
	}

	d := createFriendUser(t, repo, "dave", true)
	e := createFriendUser(t, repo, "erin", true)
	di := friendIdentity(t, store, d.ID)
	ei := friendIdentity(t, store, e.ID)
	if _, err := store.SendFriendRequest(d.ID, ei.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("send request: %v", err)
	}
	if _, err := store.AcceptFriendRequest(e.ID, di.ProfileID, time.Now().UTC()); err != nil {
		t.Fatalf("accept request: %v", err)
	}
	if err := repo.ResetProgress(d.ID); err != nil {
		t.Fatalf("reset dave: %v", err)
	}
	if n := friendCount(t, repo, "SELECT COUNT(*) FROM friendships"); n != 1 {
		t.Fatalf("reset removed friendship, rows = %d", n)
	}
	if friends, err := store.ListFriends(e.ID); err != nil || len(friends) != 1 || friends[0].StudentID != d.ID {
		t.Fatalf("friends after reset = %+v, %v", friends, err)
	}
}

func TestFriendIdentityBackfillFillsLegacyRows(t *testing.T) {
	repo := newTestStore(t)
	db := rawDB(t, repo)
	rebind := rebindQuestion
	if _, ok := repo.(*PostgresStore); ok {
		rebind = rebindDollars
	}
	now := time.Now().UTC().Format(time.RFC3339)
	first, second := newUUID(), newUUID()
	execSeed(t, repo, "INSERT INTO students (id, name, settings, created_at) VALUES (?, ?, '{}', ?)", first, "legacy one", now)
	execSeed(t, repo, "INSERT INTO students (id, name, settings, created_at) VALUES (?, ?, '{}', ?)", second, "legacy two", now)
	if err := backfillFriendIdentities(db, rebind); err != nil {
		t.Fatalf("backfill friends: %v", err)
	}
	store := friendStore(t, repo)
	firstID, err := store.GetFriendIdentity(first)
	if err != nil {
		t.Fatalf("first legacy identity: %v", err)
	}
	secondID, err := store.GetFriendIdentity(second)
	if err != nil {
		t.Fatalf("second legacy identity: %v", err)
	}
	if firstID.ProfileID == "" || firstID.FriendCode == "" || secondID.ProfileID == "" || secondID.FriendCode == "" {
		t.Fatalf("legacy identities not filled: %+v %+v", firstID, secondID)
	}
	if firstID.ProfileID == secondID.ProfileID || firstID.FriendCode == secondID.FriendCode {
		t.Fatalf("legacy identities collide: %+v %+v", firstID, secondID)
	}
	if err := backfillFriendIdentities(db, rebind); err != nil {
		t.Fatalf("second backfill: %v", err)
	}
	again, err := store.GetFriendIdentity(first)
	if err != nil || again != firstID {
		t.Fatalf("backfill is not idempotent: %+v, %v", again, err)
	}
}

func newFriendConcurrencyStore(t *testing.T) Repository {
	t.Helper()
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		return newTestStore(t)
	}
	store, err := NewSQLiteStore(filepath.Join(t.TempDir(), "friends.db"))
	if err != nil {
		t.Fatalf("create file sqlite store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

// A Google identity attests the mailbox even when the students.email_verified
// flag was never updated. The friends gate must accept that proof instead of
// asking a Google user to verify the same address again.
func TestFriendRequestAcceptsAProviderVerifiedIdentity(t *testing.T) {
	repo := newTestStore(t)
	store := friendStore(t, repo)
	a := createFriendUser(t, repo, "Ada", false)
	b := createFriendUser(t, repo, "Grace", true)
	if err := repo.CreateIdentity("google", "google-"+newUUID()[:8], a.ID, "ada@example.com", true); err != nil {
		t.Fatalf("create identity: %v", err)
	}
	target := friendIdentity(t, store, b.ID)
	if _, err := store.SendFriendRequest(a.ID, target.ProfileID, time.Now()); err != nil {
		t.Fatalf("verified identity could not send a request: %v", err)
	}
}
