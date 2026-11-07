package storage

import (
	"fmt"
	"testing"
	"time"
)

// seedBoardStudent creates an account, optionally names it, and awards xp in
// this week. The key carries the name so one test can seed fifty learners
// without inventing fifty keys by hand.
func seedBoardStudent(t *testing.T, store Repository, name, username string, xp int) string {
	t.Helper()
	st, err := store.CreateUser(name, username, "hash")
	if err != nil {
		t.Fatalf("create student %q: %v", name, err)
	}
	if xp != 0 {
		if err := store.AwardXP(st.ID, xp, "study", "test:"+name); err != nil {
			t.Fatalf("award xp to %q: %v", name, err)
		}
	}
	return st.ID
}

// The weekly board is one page of one ranking. The old endpoint returned every
// row and ranked them in Go; this pins the two numbers a client can no longer
// compute for itself: the page slice and the total.
// TestLeaderboardPage_ExcludesUUIDGuestsAndIncludesAdoptedGuest pins the real
// identity rule. A guest created through CreateStudent has a UUID id and no
// credentials, so the old guest_ prefix test missed it. An adopted guest keeps
// the guest_ id but gains credentials, so it must be ranked.
func TestLeaderboardPage_ExcludesUUIDGuestsAndIncludesAdoptedGuest(t *testing.T) {
	store := newTestStore(t)
	uuidGuest, err := store.CreateStudent("UUID Guest")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AwardXP(uuidGuest.ID, 500, "study", "uuid-guest"); err != nil {
		t.Fatal(err)
	}
	namedGuest, err := store.CreateStudent("Guest")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SetUsername(namedGuest.ID, "ivory_heron55"); err != nil {
		t.Fatal(err)
	}
	if err := store.AwardXP(namedGuest.ID, 400, "study", "named-guest"); err != nil {
		t.Fatal(err)
	}

	adopted, err := store.ClaimGuestStudent("guest_adopted", "Adopted")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AdoptGuestAsUser(adopted.ID, "Adopted", "adopted", "hash"); err != nil {
		t.Fatal(err)
	}
	if err := store.AwardXP(adopted.ID, 300, "study", "adopted"); err != nil {
		t.Fatal(err)
	}

	real, err := store.CreateUser("Real", "real", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AwardXP(real.ID, 200, "study", "real"); err != nil {
		t.Fatal(err)
	}

	week := time.Now().UTC().Add(-24 * time.Hour)
	for _, view := range []string{"weekly", "all_time"} {
		rows, total, err := store.LeaderboardPage(view, week, 25, 0)
		if err != nil {
			t.Fatalf("%s page: %v", view, err)
		}
		if total != 2 {
			t.Fatalf("%s total = %d, want 2", view, total)
		}
		seen := map[string]bool{}
		for _, row := range rows {
			seen[row.StudentID] = true
		}
		if seen[uuidGuest.ID] || seen[namedGuest.ID] {
			t.Fatalf("%s ranked a guest: %+v", view, rows)
		}
		if !seen[adopted.ID] || !seen[real.ID] {
			t.Fatalf("%s missing a real account: %+v", view, rows)
		}
	}

	rank, _, found, err := store.LeaderboardPosition(uuidGuest.ID, "weekly", week)
	if err != nil {
		t.Fatal(err)
	}
	if found || rank != 0 {
		t.Fatalf("UUID guest position = %d, found %v", rank, found)
	}
	rank, _, found, err = store.LeaderboardPosition(adopted.ID, "weekly", week)
	if err != nil {
		t.Fatal(err)
	}
	if !found || rank != 1 {
		t.Fatalf("adopted guest position = %d, found %v, want rank 1", rank, found)
	}
}

func TestLeaderboardPage_PaginatesAndCounts(t *testing.T) {
	store := newTestStore(t)
	for i := 1; i <= 30; i++ {
		seedBoardStudent(t, store, fmt.Sprintf("Learner %02d", i), fmt.Sprintf("learner%02d", i), i)
	}

	week := time.Now().UTC().Add(-24 * time.Hour)
	page1, total1, err := store.LeaderboardPage("weekly", week, 25, 0)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if total1 != 30 {
		t.Fatalf("total = %d, want 30", total1)
	}
	if len(page1) != 25 {
		t.Fatalf("page 1 rows = %d, want 25", len(page1))
	}
	if page1[0].XP != 30 {
		t.Errorf("first row xp = %d, want 30 (highest first)", page1[0].XP)
	}
	for i := 1; i < len(page1); i++ {
		if page1[i-1].XP < page1[i].XP {
			t.Fatalf("page not sorted at %d: %d then %d", i, page1[i-1].XP, page1[i].XP)
		}
	}

	page2, total2, err := store.LeaderboardPage("weekly", week, 25, 25)
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if total2 != 30 {
		t.Errorf("page 2 total = %d, want 30", total2)
	}
	if len(page2) != 5 {
		t.Fatalf("page 2 rows = %d, want 5", len(page2))
	}
	if page2[4].XP != 1 {
		t.Errorf("last row xp = %d, want 1", page2[4].XP)
	}
}

// An offset past the end is not an error and not a 404: it is an empty page with
// the real total, so the client can show "page 5 of 2" as empty rather than
// pretending the board vanished.
func TestLeaderboardPage_OutOfRangeIsEmptyWithTotal(t *testing.T) {
	store := newTestStore(t)
	seedBoardStudent(t, store, "Solo", "solo", 4)

	rows, total, err := store.LeaderboardPage("weekly", time.Now().UTC().Add(-24*time.Hour), 25, 500)
	if err != nil {
		t.Fatalf("out of range: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("rows = %d, want 0", len(rows))
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
}

// Guests are rows in students, minted on ordinary visits. They are credentials,
// so they must never be ranked publicly; a zero-XP account has nothing to rank
// yet and must not pad the board either.
func TestLeaderboardPage_ExcludesGuestsAndZeroXP(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.ClaimGuestStudent("guest_abcdef", "Guest"); err != nil {
		t.Fatalf("create guest: %v", err)
	}
	if err := store.AwardXP("guest_abcdef", 999, "study", "guest:1"); err != nil {
		t.Fatalf("award guest xp: %v", err)
	}
	seedBoardStudent(t, store, "Idle", "idle", 0)
	seedBoardStudent(t, store, "Negative", "negative", -5)
	seedBoardStudent(t, store, "Active", "active", 7)

	week := time.Now().UTC().Add(-24 * time.Hour)
	rows, total, err := store.LeaderboardPage("weekly", week, 25, 0)
	if err != nil {
		t.Fatalf("weekly: %v", err)
	}
	if total != 1 || len(rows) != 1 || rows[0].XP != 7 {
		t.Fatalf("weekly rows = %+v total=%d, want only the 7-XP account", rows, total)
	}

	rows, total, err = store.LeaderboardPage("all_time", week, 25, 0)
	if err != nil {
		t.Fatalf("all-time: %v", err)
	}
	if total != 1 || len(rows) != 1 || rows[0].XPTotal != 7 {
		t.Fatalf("all-time rows = %+v total=%d, want only the 7-XP account", rows, total)
	}
}

// The two views must draw from different sources. A learner who earned 10 XP
// last week is off this week's board and still on the all-time one.
func TestLeaderboardPage_AllTimeOutlivesTheWeek(t *testing.T) {
	store := newTestStore(t)
	sid := seedBoardStudent(t, store, "Last Week", "lastweek", 10)
	old := time.Now().UTC().Add(-8 * 24 * time.Hour).UnixNano()
	execSeed(t, store, `UPDATE xp_events SET created_at = ? WHERE student_id = ?`, old, sid)

	week := time.Now().UTC().Add(-24 * time.Hour)
	weekly, wTotal, err := store.LeaderboardPage("weekly", week, 25, 0)
	if err != nil {
		t.Fatalf("weekly: %v", err)
	}
	if wTotal != 0 || len(weekly) != 0 {
		t.Errorf("weekly = %+v total=%d, want empty: the XP was last week", weekly, wTotal)
	}

	all, aTotal, err := store.LeaderboardPage("all_time", week, 25, 0)
	if err != nil {
		t.Fatalf("all-time: %v", err)
	}
	if aTotal != 1 || len(all) != 1 || all[0].XP != 10 || all[0].XPTotal != 10 {
		t.Errorf("all-time = %+v total=%d, want the 10 XP account", all, aTotal)
	}
}

// Ties are canonical, not incidental: equal XP falls back to lifetime XP, then
// account age, then id. Without the last one Postgres and SQLite could disagree
// about which of two identical rows comes first, and so could two pages.
func TestLeaderboardPage_TiesAreDeterministic(t *testing.T) {
	store := newTestStore(t)
	var ids []string
	for i := 0; i < 5; i++ {
		ids = append(ids, seedBoardStudent(t, store, fmt.Sprintf("Tie %d", i), "", 20))
	}
	// Pin created_at to an identical value across all candidates so the tie reaches id.
	sameCreated := time.Now().UTC().Format(time.RFC3339)
	switch s := store.(type) {
	case *SQLiteStore:
		if _, err := s.db.Exec("UPDATE students SET created_at = ?", sameCreated); err != nil {
			t.Fatalf("pin created_at: %v", err)
		}
	case *PostgresStore:
		if _, err := s.db.Exec("UPDATE students SET created_at = $1", sameCreated); err != nil {
			t.Fatalf("pin created_at: %v", err)
		}
	}
	week := time.Now().UTC().Add(-24 * time.Hour)

	first, _, err := store.LeaderboardPage("weekly", week, 25, 0)
	if err != nil {
		t.Fatalf("first read: %v", err)
	}
	second, _, err := store.LeaderboardPage("weekly", week, 25, 0)
	if err != nil {
		t.Fatalf("second read: %v", err)
	}
	for i := range first {
		if first[i].StudentID != second[i].StudentID {
			t.Fatalf("tie order changed between reads at %d: %s then %s", i, first[i].StudentID, second[i].StudentID)
		}
	}
	for i := 1; i < len(first); i++ {
		if first[i-1].StudentID > first[i].StudentID {
			t.Errorf("tie not broken by id: %s before %s", first[i-1].StudentID, first[i].StudentID)
		}
	}
}

func TestLeaderboardPosition_ReportsRankOutsideThePage(t *testing.T) {
	store := newTestStore(t)
	var last string
	for i := 1; i <= 30; i++ {
		sid := seedBoardStudent(t, store, fmt.Sprintf("L%02d", i), fmt.Sprintf("l%02d", i), i)
		if i == 1 {
			last = sid
		}
	}

	rank, xp, found, err := store.LeaderboardPosition(last, "weekly", time.Now().UTC().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("position: %v", err)
	}
	if !found || rank != 30 || xp != 1 {
		t.Fatalf("position = (%d, %d, %v), want (30, 1, true)", rank, xp, found)
	}

	// The same learner is not on page 1, which is the case the UI must handle
	// with a separate row rather than by appending them to the page.
	page1, _, err := store.LeaderboardPage("weekly", time.Now().UTC().Add(-24*time.Hour), 25, 0)
	if err != nil {
		t.Fatalf("page: %v", err)
	}
	for _, row := range page1 {
		if row.StudentID == last {
			t.Fatal("rank 30 appeared on page 1")
		}
	}
}

func TestLeaderboardPage_UnknownViewIsRefused(t *testing.T) {
	store := newTestStore(t)
	if _, _, err := store.LeaderboardPage("monthly", time.Now().UTC(), 25, 0); err == nil {
		t.Fatal("unknown view accepted")
	}
}

// TestMigrate_KeepsExistingXP is the assurance an existing learner needs before
// a leaderboard rewrite: the schema change is additive. Lifetime XP lives on
// students.xp_total and is read by the all-time board; xp_events is a new table
// beside it, not a replacement. Nothing deletes or rescales a counter, and a
// second migration cannot subtract from one.
func TestMigrate_KeepsExistingXP(t *testing.T) {
	store := newTestStore(t)
	st, err := store.CreateUser("veteran", "veteran", "hash")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := store.AwardXP(st.ID, 175, "legacy", "legacy:1"); err != nil {
		t.Fatalf("award: %v", err)
	}
	// Re-running the migration must not touch what is already stored.
	if err := store.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 175 {
		t.Errorf("lifetime XP = %d after migration, want 175", total)
	}
	// And it is what the all-time board ranks.
	rows, _, err := store.LeaderboardPage("all_time", time.Now().UTC(), 25, 0)
	if err != nil {
		t.Fatalf("all-time page: %v", err)
	}
	if len(rows) != 1 || rows[0].XP != 175 {
		t.Errorf("all-time board = %+v, want the 175 XP learner", rows)
	}
}
