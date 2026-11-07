package storage

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// These tests require a real Postgres instance. They deliberately skip when
// TEST_POSTGRES_DSN is unset; the SQLite suite is not a substitute for the row
// lock that protects concurrent deductions.
func TestPostgresConcurrentXPAdjustments(t *testing.T) {
	pg := newPostgresStore(t)
	es := EffortStore(pg)

	t.Run("deductions share one balance and may go below zero", func(t *testing.T) {
		st := seedStudent(t, pg)
		if err := pg.AwardXP(st.ID, 10, "seed", "pg-floor-seed-"+st.ID); err != nil {
			t.Fatalf("seed: %v", err)
		}
		const racers = 12
		var applied int64
		var wg sync.WaitGroup
		for i := 0; i < racers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				amount, err := es.ApplyXPAdjustment(st.ID, -2, "accountability:poor_performance",
					fmt.Sprintf("pg-floor-%s-%d", st.ID, i))
				if err != nil {
					t.Errorf("deduction %d: %v", i, err)
					return
				}
				atomic.AddInt64(&applied, int64(amount))
			}(i)
		}
		wg.Wait()
		total, _, err := pg.GetXP(st.ID)
		if err != nil {
			t.Fatalf("get xp: %v", err)
		}
		if total != -14 {
			t.Fatalf("balance = %d, want -14 after all deductions", total)
		}
		if applied != -24 {
			t.Fatalf("applied = %d, want -24 (the full calculated deductions)", applied)
		}
		if ledger := pgLedgerSum(t, pg, st.ID); ledger != total {
			t.Fatalf("ledger sum = %d, balance = %d; they must agree", ledger, total)
		}
	})

	t.Run("one duplicate event applies once", func(t *testing.T) {
		st := seedStudent(t, pg)
		if err := pg.AwardXP(st.ID, 10, "seed", "pg-dup-seed-"+st.ID); err != nil {
			t.Fatalf("seed: %v", err)
		}
		const racers = 12
		var applied int64
		var wg sync.WaitGroup
		for i := 0; i < racers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				amount, err := es.ApplyXPAdjustment(st.ID, -3, "accountability:repeated_guessing",
					"pg-duplicate-"+st.ID)
				if err != nil {
					t.Errorf("duplicate deduction: %v", err)
					return
				}
				atomic.AddInt64(&applied, int64(amount))
			}()
		}
		wg.Wait()
		total, _, err := pg.GetXP(st.ID)
		if err != nil {
			t.Fatalf("get xp: %v", err)
		}
		if total != 7 || applied != -3 {
			t.Fatalf("balance = %d, applied = %d; want 7 and -3", total, applied)
		}
		var rows int
		if err := pg.db.QueryRow(`SELECT COUNT(*) FROM xp_events WHERE student_id = $1 AND event_key = $2`,
			st.ID, "pg-duplicate-"+st.ID).Scan(&rows); err != nil {
			t.Fatalf("count duplicate event: %v", err)
		}
		if rows != 1 {
			t.Fatalf("duplicate event rows = %d, want 1", rows)
		}
		if ledger := pgLedgerSum(t, pg, st.ID); ledger != total {
			t.Fatalf("ledger sum = %d, balance = %d", ledger, total)
		}
	})

	t.Run("concurrent awards and deductions keep ledger and balance equal", func(t *testing.T) {
		st := seedStudent(t, pg)
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(2)
			go func(i int) {
				defer wg.Done()
				if err := pg.AwardXP(st.ID, 5, "study", fmt.Sprintf("pg-award-%s-%d", st.ID, i)); err != nil {
					t.Errorf("award %d: %v", i, err)
				}
			}(i)
			go func(i int) {
				defer wg.Done()
				if _, err := es.ApplyXPAdjustment(st.ID, -4, "accountability:poor_performance",
					fmt.Sprintf("pg-mixed-%s-%d", st.ID, i)); err != nil {
					t.Errorf("deduction %d: %v", i, err)
				}
			}(i)
		}
		wg.Wait()
		total, _, err := pg.GetXP(st.ID)
		if err != nil {
			t.Fatalf("get xp: %v", err)
		}
		if total != 10 {
			t.Fatalf("balance = %d, want 10 after 10 awards of 5 and 10 deductions of 4", total)
		}
		if ledger := pgLedgerSum(t, pg, st.ID); ledger != total {
			t.Fatalf("ledger sum = %d, balance = %d; concurrent writes diverged", ledger, total)
		}
	})

	t.Run("positive XP recovers a negative balance", func(t *testing.T) {
		st := seedStudent(t, pg)
		if err := pg.AwardXP(st.ID, 10, "seed", "pg-recovery-seed-"+st.ID); err != nil {
			t.Fatalf("seed: %v", err)
		}
		if _, err := es.ApplyXPAdjustment(st.ID, -15, "accountability:extremely_poor_performance", "pg-recovery-down-"+st.ID); err != nil {
			t.Fatalf("deduction: %v", err)
		}
		if err := pg.AwardXP(st.ID, 12, "study", "pg-recovery-up-"+st.ID); err != nil {
			t.Fatalf("recovery award: %v", err)
		}
		total, _, err := pg.GetXP(st.ID)
		if err != nil {
			t.Fatalf("get xp: %v", err)
		}
		if total != 7 {
			t.Fatalf("balance = %d, want 7 after recovering from -5", total)
		}
		if ledger := pgLedgerSum(t, pg, st.ID); ledger != total {
			t.Fatalf("ledger sum = %d, balance = %d", ledger, total)
		}
	})
}

func pgLedgerSum(t *testing.T, pg *PostgresStore, studentID string) int {
	t.Helper()
	var total int
	if err := pg.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM xp_events WHERE student_id = $1`,
		studentID).Scan(&total); err != nil {
		t.Fatalf("ledger sum: %v", err)
	}
	return total
}
