package storage

import "testing"

// A closed league week is read from its snapshot. A later XP award must not
// rewrite the table a learner already saw, and the stored order must survive.
func TestLeagueRosterFinal_ReadsTheStoredSnapshot(t *testing.T) {
	store := newTestStore(t)
	low := seedBoardStudent(t, store, "Low", "low", 10)
	high := seedBoardStudent(t, store, "High", "high", 20)
	execSeed(t, store, `INSERT INTO league_members
		(week, student_id, tier, league_no, xp, rank, movement, from_tier)
		VALUES ('2026-11-09', ?, 'bronze', 1, 10, 2, 'stayed', 'bronze')`, low)
	execSeed(t, store, `INSERT INTO league_members
		(week, student_id, tier, league_no, xp, rank, movement, from_tier)
		VALUES ('2026-11-09', ?, 'bronze', 1, 20, 1, 'stayed', 'bronze')`, high)
	if err := store.AwardXP(low, 999, "study", "after-the-week"); err != nil {
		t.Fatal(err)
	}

	rows, err := store.LeagueRosterFinal("2026-11-09", "bronze", 1)
	if err != nil {
		t.Fatalf("final roster: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
	if rows[0].StudentID != high || rows[0].XP != 20 || rows[0].Rank != 1 {
		t.Fatalf("first row = %+v, want High at 20 and rank 1", rows[0])
	}
	if rows[1].StudentID != low || rows[1].XP != 10 || rows[1].Rank != 2 {
		t.Fatalf("second row = %+v, want Low still at 10 and rank 2", rows[1])
	}
}

// The transition writes the result of week W onto W+1. The outcome reader must
// ask for the new seat, not the closed one.
func TestLeagueOutcome_ReadsTheNextSeat(t *testing.T) {
	store := newTestStore(t)
	sid := seedBoardStudent(t, store, "Mover", "mover", 30)
	execSeed(t, store, `INSERT INTO league_members
		(week, student_id, tier, league_no, xp, rank, movement, from_tier)
		VALUES ('2026-11-09', ?, 'silver', 1, 30, 1, 'promoted', 'bronze')`, sid)

	movement, fromTier, toTier, found, err := store.LeagueOutcome(sid, "2026-11-09")
	if err != nil {
		t.Fatalf("outcome: %v", err)
	}
	if !found || movement != "promoted" || fromTier != "bronze" || toTier != "silver" {
		t.Fatalf("outcome = %q %q %q found=%v", movement, fromTier, toTier, found)
	}

	_, _, _, found, err = store.LeagueOutcome(sid, "2026-11-02")
	if err != nil {
		t.Fatalf("missing outcome: %v", err)
	}
	if found {
		t.Fatal("a week with no seat must not invent an outcome")
	}
}

func TestEarliestLeagueWeek(t *testing.T) {
	store := newTestStore(t)
	sid := seedBoardStudent(t, store, "First", "first", 1)
	execSeed(t, store, `INSERT INTO league_members
		(week, student_id, tier, league_no, xp, rank, movement, from_tier)
		VALUES ('2026-10-26', ?, 'bronze', 1, 1, 1, 'new', '')`, sid)
	execSeed(t, store, `INSERT INTO league_members
		(week, student_id, tier, league_no, xp, rank, movement, from_tier)
		VALUES ('2026-11-02', ?, 'bronze', 1, 2, 1, 'stayed', 'bronze')`, sid)

	week, ok, err := store.EarliestLeagueWeek()
	if err != nil {
		t.Fatalf("earliest week: %v", err)
	}
	if !ok || week != "2026-10-26" {
		t.Fatalf("earliest week = %q ok=%v, want 2026-10-26", week, ok)
	}
}
