package storage

import (
	"fmt"
	"testing"
	"time"
)

// leagueFixture creates n learners with descending lifetime XP, all earned
// inside the week starting prevFrom, and returns them in XP order. The first
// week is bootstrapped from lifetime XP; later weeks use the ledger window.
func leagueFixture(t *testing.T, store Repository, n int, prevFrom time.Time) []*Student {
	t.Helper()
	out := make([]*Student, 0, n)
	for i := 0; i < n; i++ {
		st, err := store.CreateStudent(fmt.Sprintf("learner%02d", i))
		if err != nil {
			t.Fatalf("create learner %d: %v", i, err)
		}
		if err := store.AwardXP(st.ID, 1000-i, "test", fmt.Sprintf("seed:%d", i)); err != nil {
			t.Fatalf("award learner %d: %v", i, err)
		}
		out = append(out, st)
	}
	// Move every seed event into the previous week's window so the transition
	// reads them the way it would in production.
	execSeed(t, store, `UPDATE xp_events SET created_at = ?`, prevFrom.Add(24*time.Hour).UnixNano())
	return out
}

func weekMonday(t *testing.T, value string) time.Time {
	t.Helper()
	day, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("parse week: %v", err)
	}
	return day.UTC()
}

func countWeekRows(t *testing.T, store Repository, week string) int {
	t.Helper()
	return seedInt(t, store, `SELECT COUNT(*) FROM league_members WHERE week = ?`, week)
}

func tierCounts(t *testing.T, store Repository, week string) map[string]int {
	t.Helper()
	rows := querySeed(t, store,
		`SELECT tier, COUNT(*) FROM league_members WHERE week = ? GROUP BY tier`, week)
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var tier string
		var n int
		if err := rows.Scan(&tier, &n); err != nil {
			t.Fatalf("scan tier: %v", err)
		}
		out[tier] = n
	}
	return out
}

// TestResolveLeagueWeek_BootstrapsThenPromotesAndRelegates pins the whole
// transition: the first week seats everyone from lifetime XP, the second moves
// the top five of each thirty-person league up and the bottom five down, and a
// repeat call does not move anybody a second time.
func TestResolveLeagueWeek_BootstrapsThenPromotesAndRelegates(t *testing.T) {
	store := newTestStore(t)
	week1 := weekMonday(t, "2026-10-12")
	week2 := week1.AddDate(0, 0, 7)
	learners := leagueFixture(t, store, 40, week1)

	if err := store.ResolveLeagueWeek("2026-10-12", week1.AddDate(0, 0, -7), week1); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if got := countWeekRows(t, store, "2026-10-12"); got != 40 {
		t.Fatalf("bootstrap seats = %d, want 40", got)
	}
	// League 1 is full at thirty; the remaining ten form league 2.
	league1 := seedInt(t, store, `SELECT COUNT(*) FROM league_members WHERE week='2026-10-12' AND tier='bronze' AND league_no=1`)
	league2 := seedInt(t, store, `SELECT COUNT(*) FROM league_members WHERE week='2026-10-12' AND tier='bronze' AND league_no=2`)
	if league1 != LeagueSize || league2 != 10 {
		t.Errorf("bootstrap split = %d/%d, want %d/10", league1, league2, LeagueSize)
	}

	if err := store.ResolveLeagueWeek("2026-10-19", week1, week2); err != nil {
		t.Fatalf("transition: %v", err)
	}
	counts := tierCounts(t, store, "2026-10-19")
	// Five from league 1 and three from league 2 (min(5, 10/3)) move up.
	if counts["silver"] != 8 {
		t.Errorf("silver = %d, want 8", counts["silver"])
	}
	if counts["bronze"] != 32 {
		t.Errorf("bronze = %d, want 32", counts["bronze"])
	}

	// The five highest earners are silver; the lowest five of league 1 stay
	// bronze rather than being carried up.
	for i := 0; i < 5; i++ {
		tier, _, found, err := store.LeagueMembershipOf(learners[i].ID, "2026-10-19")
		if err != nil || !found {
			t.Fatalf("membership of top learner %d: %v found=%v", i, err, found)
		}
		if tier != "silver" {
			t.Errorf("top learner %d tier = %q, want silver", i, tier)
		}
	}
	for i := 25; i < 30; i++ {
		tier, _, _, _ := store.LeagueMembershipOf(learners[i].ID, "2026-10-19")
		if tier != "bronze" {
			t.Errorf("learner %d (bottom of league 1) tier = %q, want bronze", i, tier)
		}
	}

	// Retry: the week marker makes the second call a no-op.
	before := countWeekRows(t, store, "2026-10-19")
	if err := store.ResolveLeagueWeek("2026-10-19", week1, week2); err != nil {
		t.Fatalf("retry: %v", err)
	}
	if after := countWeekRows(t, store, "2026-10-19"); after != before {
		t.Errorf("retry wrote %d seats, was %d", after, before)
	}
}

// TestEnsureLeagueMembership_SeatsLateJoinersInBronze covers the learner who
// starts earning XP after the transition has already run: they join the
// lowest-numbered bronze league with room, and asking twice does not move them.
func TestEnsureLeagueMembership_SeatsLateJoinersInBronze(t *testing.T) {
	store := newTestStore(t)
	week := weekMonday(t, "2026-10-19")
	leagueFixture(t, store, 35, week)

	if err := store.ResolveLeagueWeek("2026-10-19", week.AddDate(0, 0, -7), week); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	late, err := store.CreateUser("Late", "late_learner", "hash")
	if err != nil {
		t.Fatalf("create late: %v", err)
	}
	if err := store.AwardXP(late.ID, 12, "test", "late:1"); err != nil {
		t.Fatalf("award late: %v", err)
	}

	if err := store.EnsureLeagueMembership(late.ID, "2026-10-19"); err != nil {
		t.Fatalf("seat late: %v", err)
	}
	tier, no, found, err := store.LeagueMembershipOf(late.ID, "2026-10-19")
	if err != nil || !found {
		t.Fatalf("late membership: %v found=%v", err, found)
	}
	if tier != "bronze" || no != 2 {
		t.Errorf("late seat = %s/%d, want bronze/2", tier, no)
	}
	// Idempotent: a second call leaves the seat alone.
	if err := store.EnsureLeagueMembership(late.ID, "2026-10-19"); err != nil {
		t.Fatalf("re-seat: %v", err)
	}
	tier2, no2, _, _ := store.LeagueMembershipOf(late.ID, "2026-10-19")
	if tier2 != tier || no2 != no {
		t.Errorf("re-seat moved late learner to %s/%d", tier2, no2)
	}
}

// TestEnsureLeagueMembership_SkipsGuestsAndZeroXP pins the eligibility rule:
// a guest credential and an account that has earned nothing do not occupy a
// seat, so a week's leagues are made of learners who actually competed.
func TestEnsureLeagueMembership_SkipsGuestsAndZeroXP(t *testing.T) {
	store := newTestStore(t)
	week := weekMonday(t, "2026-10-19")

	guest, err := store.ClaimGuestStudent("guest_league_test", "Guest")
	if err != nil {
		t.Fatalf("claim guest: %v", err)
	}
	if err := store.EnsureLeagueMembership(guest.ID, "2026-10-19"); err != nil {
		t.Fatalf("guest seat: %v", err)
	}
	if _, _, found, _ := store.LeagueMembershipOf(guest.ID, "2026-10-19"); found {
		t.Error("guest was seated in a league")
	}

	quiet, err := store.CreateUser("Quiet", "quiet_learner", "hash")
	if err != nil {
		t.Fatalf("create quiet: %v", err)
	}
	if err := store.EnsureLeagueMembership(quiet.ID, "2026-10-19"); err != nil {
		t.Fatalf("zero-xp seat: %v", err)
	}
	if _, _, found, _ := store.LeagueMembershipOf(quiet.ID, "2026-10-19"); found {
		t.Error("zero-XP account was seated in a league")
	}

	// A rejected week is not a Monday.
	if err := store.ResolveLeagueWeek("2026-10-20", week, week); err == nil {
		t.Error("non-Monday week accepted")
	}
}

// TestLeagueRoster_OrdersByXPThenNameAndID fixes the tie-break the pagination
// and promotion rules rely on: equal XP is ordered by name, then by account id,
// so two learners on the same score never swap places between reads.
func TestLeagueRoster_OrdersByXPThenNameAndID(t *testing.T) {
	store := newTestStore(t)
	week := weekMonday(t, "2026-10-19")
	ids := []string{}
	for _, name := range []string{"zoe", "mia", "abe"} {
		st, err := store.CreateStudent(name)
		if err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
		if err := store.AwardXP(st.ID, 10, "test", "tie:"+name); err != nil {
			t.Fatalf("award %s: %v", name, err)
		}
		execSeed(t, store, `UPDATE xp_events SET created_at = ? WHERE student_id = ?`,
			week.Add(time.Hour).UnixNano(), st.ID)
		ids = append(ids, st.ID)
	}
	if err := store.ResolveLeagueWeek("2026-10-19", week.AddDate(0, 0, -7), week); err != nil {
		t.Fatalf("seat: %v", err)
	}
	tier, no, _, _ := store.LeagueMembershipOf(ids[0], "2026-10-19")
	rows, err := store.LeagueRoster("2026-10-19", tier, no, week, week.AddDate(0, 0, 7))
	if err != nil {
		t.Fatalf("roster: %v", err)
	}
	want := []string{"abe", "mia", "zoe"}
	if len(rows) != len(want) {
		t.Fatalf("roster size = %d, want %d", len(rows), len(want))
	}
	for i, name := range want {
		if rows[i].Name != name {
			t.Errorf("row %d = %s, want %s", i, rows[i].Name, name)
		}
	}
	// Equal scores share the competition rank; the next distinct score does
	// not skip a place it never held.
	if rows[0].Rank != 1 || rows[1].Rank != 1 || rows[2].Rank != 1 {
		t.Errorf("tie ranks = %d/%d/%d, want 1/1/1", rows[0].Rank, rows[1].Rank, rows[2].Rank)
	}
}

// Awarding XP is the moment an account becomes eligible for a league. The
// profile and My League used to be the only places that seated somebody, so an
// active learner who never opened either had no tier and no avatar frame.
func TestAwardXPSeatsAWeeklyLeague(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateUser("Earner", "earner", "hash")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := store.AwardXP(st.ID, 40, "test", "seat:1"); err != nil {
		t.Fatalf("award: %v", err)
	}
	week := LeagueWeekKey(time.Now())
	tier, _, found, err := store.LeagueMembershipOf(st.ID, week)
	if err != nil {
		t.Fatalf("membership: %v", err)
	}
	if !found || tier != "bronze" {
		t.Fatalf("membership = %q found=%v, want bronze", tier, found)
	}
}

func TestBackfillLeagueMembershipsSeatsEligibleAccounts(t *testing.T) {
	store := newTestStore(t)
	earner, err := store.CreateUser("Legacy Earner", "legacy_earner", "hash")
	if err != nil {
		t.Fatalf("create earner: %v", err)
	}
	if err := store.AwardXP(earner.ID, 25, "test", "legacy:1"); err != nil {
		t.Fatalf("award: %v", err)
	}
	quiet, err := store.CreateUser("Quiet", "quiet_learner", "hash")
	if err != nil {
		t.Fatalf("create quiet: %v", err)
	}
	week := LeagueWeekKey(time.Now())
	// Simulate the old lazy behaviour: the XP exists, the seat does not.
	execSeed(t, store, `DELETE FROM league_members WHERE student_id = ?`, earner.ID)

	backfiller, ok := store.(interface {
		BackfillLeagueMemberships(string) (int, error)
	})
	if !ok {
		t.Fatal("store does not support league backfill")
	}
	seated, err := backfiller.BackfillLeagueMemberships(week)
	if err != nil {
		t.Fatalf("backfill: %v", err)
	}
	if seated != 1 {
		t.Fatalf("seated = %d, want 1", seated)
	}
	tier, _, found, err := store.LeagueMembershipOf(earner.ID, week)
	if err != nil || !found || tier != "bronze" {
		t.Fatalf("earner membership = %q found=%v err=%v, want bronze", tier, found, err)
	}
	if _, _, found, err := store.LeagueMembershipOf(quiet.ID, week); err != nil || found {
		t.Fatalf("quiet account found=%v err=%v, want no seat", found, err)
	}
}

// Adopted guests keep their guest_ id but have credentials. The id prefix is
// not the account type, so they must be seated and framed like any learner.
func TestAdoptedGuestIsSeatedByAward(t *testing.T) {
	store := newTestStore(t)
	guest, err := store.ClaimGuestStudent("guest_adopted_league", "Guest")
	if err != nil {
		t.Fatalf("claim guest: %v", err)
	}
	adopted, err := store.AdoptGuestAsUser(guest.ID, "Adopted", "adopted", "hash")
	if err != nil {
		t.Fatalf("adopt: %v", err)
	}
	if err := store.AwardXP(adopted.ID, 30, "test", "adopt:1"); err != nil {
		t.Fatalf("award: %v", err)
	}
	week := LeagueWeekKey(time.Now())
	tier, _, found, err := store.LeagueMembershipOf(adopted.ID, week)
	if err != nil || !found || tier != "bronze" {
		t.Fatalf("adopted membership = %q found=%v err=%v, want bronze", tier, found, err)
	}
}

func TestLeagueRangeAggregatesStoredMovement(t *testing.T) {
	store := newTestStore(t)
	a, err := store.CreateUser("Range A", "range_a", "hash")
	if err != nil {
		t.Fatalf("create a: %v", err)
	}
	b, err := store.CreateUser("Range B", "range_b", "hash")
	if err != nil {
		t.Fatalf("create b: %v", err)
	}
	execSeed(t, store, `INSERT INTO league_members (week, student_id, tier, league_no, xp, rank, movement) VALUES
		('2026-11-02', ?, 'bronze', 1, 10, 1, 'new'),
		('2026-11-02', ?, 'bronze', 1, 5, 2, 'new'),
		('2026-11-09', ?, 'silver', 1, 12, 1, 'up'),
		('2026-11-09', ?, 'bronze', 1, 8, 2, 'stay')`, a.ID, b.ID, a.ID, b.ID)
	got, err := store.(interface {
		LeagueRange(string, string, string) (*LeagueRange, error)
	}).LeagueRange("2026-11-02", "2026-11-09", "week")
	if err != nil {
		t.Fatalf("range: %v", err)
	}
	if got.Weeks != 2 || got.Seats != 4 || got.Unique != 2 || got.Up != 1 || got.Stay != 1 || got.New != 2 {
		t.Fatalf("range totals = %+v", got)
	}
	if len(got.Buckets) != 2 || len(got.Tiers) != 2 {
		t.Fatalf("buckets=%d tiers=%d", len(got.Buckets), len(got.Tiers))
	}
}
