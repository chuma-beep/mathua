package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// The ledger is only worth having if the uniqueness rule actually holds. These
// tests exercise it at the store, where the database constraint lives, rather
// than at the handler, where a mistaken caller could bypass it.

func TestAwardXP_IsIdempotentPerEventKey(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)

	for i := 0; i < 3; i++ {
		if err := store.AwardXP(st.ID, 12, "study", "attempt:41"); err != nil {
			t.Fatalf("award %d: %v", i, err)
		}
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 12 {
		t.Fatalf("xp_total = %d, want 12 after three awards of one event", total)
	}
	events := seedInt(t, store, `SELECT COUNT(*) FROM xp_events WHERE student_id = ?`, st.ID)
	if events != 1 {
		t.Fatalf("ledger rows = %d, want 1", events)
	}
}

func TestAwardXP_ConcurrentDuplicatesCannotPayTwice(t *testing.T) {
	// A file-backed store, not ":memory:", because the pool opens more than one
	// connection and each ":memory:" connection would be its own empty database.
	store, err := NewSQLiteStore(t.TempDir() + "/xp.db")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	st := seedStudent(t, store)

	const racers = 16
	var wg sync.WaitGroup
	errs := make([]error, racers)
	for i := 0; i < racers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = store.AwardXP(st.ID, 9, "study", "attempt:concurrent")
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("racer %d: %v", i, err)
		}
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 9 {
		t.Fatalf("xp_total = %d, want 9: %d concurrent awards of one event paid %d times",
			total, racers, total/9)
	}
}

func TestAwardXP_AllowsNegativeLifetimeButKeepsTodayNonNegative(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)

	if err := store.AwardXP(st.ID, -5, "study", "negative:1"); err != nil {
		t.Fatalf("negative award: %v", err)
	}
	total, today, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp after negative: %v", err)
	}
	if total != -5 || today != 0 {
		t.Fatalf("after -5: total=%d today=%d, want -5 and 0", total, today)
	}
	if err := store.AwardXP(st.ID, 4, "study", "negative:2"); err != nil {
		t.Fatalf("recovery award: %v", err)
	}
	total, today, err = store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp after recovery: %v", err)
	}
	if total != -1 || today != 4 {
		t.Fatalf("after +4: total=%d today=%d, want -1 and 4", total, today)
	}
	if err := store.AwardXP(st.ID, 10, "study", "negative:3"); err != nil {
		t.Fatalf("positive award: %v", err)
	}
	total, today, err = store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp after positive: %v", err)
	}
	if total != 9 || today != 14 {
		t.Fatalf("after +10: total=%d today=%d, want 9 and 14", total, today)
	}
}

func TestAwardXP_DifferentEventsAccumulateAndNegativeAmountsAreRecorded(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)

	if err := store.AwardXP(st.ID, 10, "study", "attempt:1"); err != nil {
		t.Fatalf("award 1: %v", err)
	}
	if err := store.AwardXP(st.ID, -3, "rush_penalty", "penalty:1"); err != nil {
		t.Fatalf("penalty: %v", err)
	}
	total, _, err := store.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	if total != 7 {
		t.Fatalf("xp_total = %d, want 7", total)
	}
}

func TestAwardXP_RefusesEmptyIdentityOrKey(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	if err := store.AwardXP("", 5, "study", "attempt:1"); err == nil {
		t.Error("empty student accepted")
	}
	if err := store.AwardXP(st.ID, 5, "study", ""); err == nil {
		t.Error("empty event key accepted")
	}
}

func TestWeeklyXPByStudent_CountsOnlyTheWindow(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	if err := store.AwardXP(st.ID, 20, "study", "attempt:this-week"); err != nil {
		t.Fatalf("award: %v", err)
	}
	// Backdate the row: the query is a time window, not a label on the row, so
	// this is the only way to prove the window is applied.
	execSeed(t, store,
		`UPDATE xp_events SET created_at = ? WHERE event_key = 'attempt:this-week'`,
		time.Now().UTC().Add(-14*24*time.Hour).UnixNano())
	if err := store.AwardXP(st.ID, 5, "study", "attempt:now"); err != nil {
		t.Fatalf("award now: %v", err)
	}

	monday := time.Now().UTC().Add(-24 * time.Hour)
	got, err := store.WeeklyXPByStudent(monday)
	if err != nil {
		t.Fatalf("weekly xp: %v", err)
	}
	if got[st.ID] != 5 {
		t.Fatalf("weekly xp = %d, want 5 (the old row must not count)", got[st.ID])
	}
}

// seedStudent creates a real account. It is not a guest on purpose: the
// leaderboard excludes guests, and a guest here would make a later eligibility
// test pass for the wrong reason.
func seedStudent(t *testing.T, store Repository) *Student {
	t.Helper()
	st, err := store.CreateStudent(fmt.Sprintf("Learner %d", time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	return st
}
