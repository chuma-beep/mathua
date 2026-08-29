package leaderboard

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func leagueStore(t *testing.T) storage.Repository {
	t.Helper()
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func seedLeague(t *testing.T, repo storage.Repository, names []string) {
	t.Helper()
	for _, n := range names {
		st, err := repo.CreateStudent(n)
		if err != nil {
			t.Fatalf("create %s: %v", n, err)
		}
		_ = repo.SetLeague(st.ID, "bronze", "", 0)
	}
}

func TestWeekKey(t *testing.T) {
	if got := weekKey(time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)); got != "2026-W35" {
		t.Errorf("expected 2026-W35, got %q", got)
	}
}

func TestResolveWeeklyReset_PromotesTopAndDemotesBottom(t *testing.T) {
	repo := leagueStore(t)
	now := time.Now().UTC()
	ids := make(map[string]string)
	for _, n := range []string{"a", "b", "c", "d", "e", "f"} {
		st, _ := repo.CreateStudent(n)
		ids[n] = st.ID
		_ = repo.SetLeague(st.ID, "bronze", "", 0)
	}
	// Weekly mastery: a=3, b=2 this week; others 0. Bronze has no lower tier,
	// so nobody demotes; a and b promote to silver.
	for i, n := range []string{"a", "b"} {
		for k := 0; k <= i; k++ {
			_ = repo.UpsertProgress(&storage.ConceptProgress{
				StudentID: ids[n], ConceptID: string(rune('a' + k)), Status: "MASTERED", MasteredAt: &now,
			})
		}
	}
	if _, err := ResolveWeeklyReset(repo, now); err != nil {
		t.Fatalf("reset: %v", err)
	}
	for _, tc := range []struct{ name, want string }{{"a", "silver"}, {"b", "silver"}, {"c", "bronze"}} {
		st, _ := repo.GetStudent(ids[tc.name])
		if st.League != tc.want {
			t.Errorf("%s: expected %s, got %s", tc.name, tc.want, st.League)
		}
	}
	// a should carry the promoted badge.
	a, _ := repo.GetStudent(ids["a"])
	if a.LeagueMoved != 1 {
		t.Errorf("expected promoted badge on a, got %d", a.LeagueMoved)
	}
}

func TestStandings_BoardShape(t *testing.T) {
	repo := leagueStore(t)
	now := time.Now().UTC()
	key := weekKey(now)
	for _, n := range []string{"a", "b", "c"} {
		st, _ := repo.CreateStudent(n)
		// Seed at the current week so the reset leaves them in bronze.
		_ = repo.SetLeague(st.ID, "bronze", key, 0)
	}
	board, err := Standings(repo, now)
	if err != nil {
		t.Fatalf("standings: %v", err)
	}
	if board.Week == "" {
		t.Error("expected week key")
	}
	if len(board.Leagues) == 0 {
		t.Fatal("expected at least the bronze league")
	}
	if board.Leagues[0].Tier != "bronze" {
		t.Errorf("expected bronze first, got %s", board.Leagues[0].Tier)
	}
	if len(board.Leagues[0].Members) != 3 {
		t.Errorf("expected 3 bronze members, got %d", len(board.Leagues[0].Members))
	}
}

func TestResolveWeeklyReset_Idempotent(t *testing.T) {
	repo := leagueStore(t)
	st, _ := repo.CreateStudent("solo")
	_ = repo.SetLeague(st.ID, "bronze", "", 0)

	key1, err := ResolveWeeklyReset(repo, time.Now().UTC())
	if err != nil {
		t.Fatalf("first reset: %v", err)
	}
	got, _ := repo.GetStudent(st.ID)
	if got.LeagueWeek != key1 {
		t.Errorf("expected league_week %s, got %s", key1, got.LeagueWeek)
	}
	// Second call in the same week must not touch the student again.
	got2, _ := repo.GetStudent(st.ID)
	_ = got2
	key2, err := ResolveWeeklyReset(repo, time.Now().UTC())
	if err != nil {
		t.Fatalf("second reset: %v", err)
	}
	if key2 != key1 {
		t.Errorf("expected same week key, got %s vs %s", key1, key2)
	}
	final, _ := repo.GetStudent(st.ID)
	if final.LeagueMoved != 0 {
		t.Errorf("solo student must not move, got %d", final.LeagueMoved)
	}
}
