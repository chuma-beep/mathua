package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/repetition"
)

func creditApplyStore(t *testing.T, repo Repository) (CreditStore, CreditApplyStore, CreditPreviewStore) {
	t.Helper()
	alloc, ok := repo.(CreditStore)
	if !ok {
		t.Fatalf("%T does not implement CreditStore", repo)
	}
	apply, ok := repo.(CreditApplyStore)
	if !ok {
		t.Fatalf("%T does not implement CreditApplyStore", repo)
	}
	preview, ok := repo.(CreditPreviewStore)
	if !ok {
		t.Fatalf("%T does not implement CreditPreviewStore", repo)
	}
	return alloc, apply, preview
}

func activateCreditTest(t *testing.T, repo Repository, cohort string) {
	t.Helper()
	store, ok := repo.(FireActivationStore)
	if !ok {
		t.Fatalf("%T does not implement FireActivationStore", repo)
	}
	if _, err := store.SetFireActivation(FireActivation{
		Mode: "active", Cohort: cohort, AcknowledgedBy: "test",
		AcknowledgedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("activate fire: %v", err)
	}
}

func seedMasteredTarget(t *testing.T, repo Repository, studentID, conceptID string, due time.Time) *ConceptProgress {
	t.Helper()
	last := due.AddDate(0, 0, -10)
	stability := 4.0
	difficulty := 5.0
	p := &ConceptProgress{
		StudentID:      studentID,
		ConceptID:      conceptID,
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		BestStreak:     4,
		Attempts:       9,
		LastAttempted:  &last,
		LastReviewed:   &last,
		NextReviewDue:  &due,
		SM2Repetitions: 3,
		SM2Interval:    10,
		SM2EFactor:     2.5,
		MasteredAt:     &last,
		FSRSStability:  &stability,
		FSRSDifficulty: &difficulty,
		FSRSLastReview: &last,
	}
	if err := repo.UpsertProgress(p); err != nil {
		t.Fatalf("seed progress: %v", err)
	}
	return p
}

func TestApplyCreditMovesOnlyTheSchedule(t *testing.T) {
	for name, repo := range creditStores(t) {
		t.Run(name, func(t *testing.T) {
			alloc, apply, preview := creditApplyStore(t, repo)
			activateCreditTest(t, repo, "test")
			st := seedStudent(t, repo)
			now := time.Now().UTC()
			due := now.AddDate(0, 0, -1)
			before := seedMasteredTarget(t, repo, st.ID, "target", due)
			if _, _, err := alloc.AllocateCredit(st.ID, "attempt-source", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
				t.Fatalf("allocate: %v", err)
			}
			xpBefore, _, err := repo.GetXP(st.ID)
			if err != nil {
				t.Fatalf("xp before: %v", err)
			}

			pv, err := preview.PreviewCredit(CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-source", Epoch: 1, Cohort: "test", Now: now,
			})
			if err != nil {
				t.Fatalf("preview: %v", err)
			}
			if !pv.Applied || pv.ProposedDue == "" {
				t.Fatalf("preview = %+v", pv)
			}
			after, err := repo.GetProgress(st.ID, "target")
			if err != nil {
				t.Fatalf("progress after preview: %v", err)
			}
			if after.NextReviewDue == nil || after.NextReviewDue.Unix() != due.Truncate(time.Second).Unix() {
				t.Fatalf("preview moved the due date: %v", after.NextReviewDue)
			}

			app, err := apply.ApplyCredit(CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-source", Epoch: 1, Cohort: "test", Now: now,
			})
			if err != nil {
				t.Fatalf("apply: %v", err)
			}
			if !app.Applied || app.Replayed {
				t.Fatalf("application = %+v", app)
			}
			after, err = repo.GetProgress(st.ID, "target")
			if err != nil {
				t.Fatalf("progress: %v", err)
			}
			if after.Status != before.Status || after.Attempts != before.Attempts ||
				after.Streak != before.Streak || after.BestStreak != before.BestStreak ||
				after.PlacementSeeded || after.IsProvisional {
				t.Fatalf("application changed learner state: before=%+v after=%+v", before, after)
			}
			if after.NextReviewDue == nil || !after.NextReviewDue.After(due.Truncate(time.Second)) {
				t.Fatalf("due did not move forward: %v", after.NextReviewDue)
			}
			if after.FSRSStability == nil || *after.FSRSStability <= 0 {
				t.Fatal("stability was not written")
			}
			xpAfter, _, err := repo.GetXP(st.ID)
			if err != nil {
				t.Fatalf("xp after: %v", err)
			}
			if xpAfter != xpBefore {
				t.Fatalf("xp changed: %d -> %d", xpBefore, xpAfter)
			}

			replay, err := apply.ApplyCredit(CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-source", Epoch: 1, Cohort: "test", Now: now.Add(time.Minute),
			})
			if err != nil {
				t.Fatalf("replay: %v", err)
			}
			if !replay.Replayed {
				t.Fatalf("replay = %+v", replay)
			}
			apps, err := apply.CreditApplications(st.ID, "target", 10)
			if err != nil {
				t.Fatalf("applications: %v", err)
			}
			if len(apps) != 1 {
				t.Fatalf("applications = %d, want 1", len(apps))
			}
		})
	}
}

func TestApplyCreditRejectsIneligibleTargetsWithoutConsumingTheOrdinal(t *testing.T) {
	cases := []struct {
		name   string
		status string
		due    func(time.Time) time.Time
		reason string
	}{
		{"not mastered", string(mastery.StatusPracticing), func(now time.Time) time.Time { return now.AddDate(0, 0, -1) }, "not_mastered"},
		{"not due", string(mastery.StatusMastered), func(now time.Time) time.Time { return now.AddDate(0, 0, 2) }, "not_due"},
	}
	for name, repo := range creditStores(t) {
		for _, tc := range cases {
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				alloc, apply, _ := creditApplyStore(t, repo)
				activateCreditTest(t, repo, "test")
				st := seedStudent(t, repo)
				now := time.Now().UTC()
				due := tc.due(now)
				p := seedMasteredTarget(t, repo, st.ID, "target", due)
				p.Status = tc.status
				if err := repo.UpsertProgress(p); err != nil {
					t.Fatalf("seed: %v", err)
				}
				if _, _, err := alloc.AllocateCredit(st.ID, "attempt-"+tc.name, 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
					t.Fatalf("allocate: %v", err)
				}
				res, err := apply.ApplyCredit(CreditApplyRequest{
					StudentID: st.ID, TargetConceptID: "target", SourceAttemptID: "attempt-" + tc.name, Epoch: 1, Cohort: "test", Now: now,
				})
				if err != nil {
					t.Fatalf("apply: %v", err)
				}
				if res.Applied || res.Reason != tc.reason {
					t.Fatalf("result = %+v, want reason %q", res, tc.reason)
				}
				ordinals, err := repo.(CreditOrdinalStore).BankedOrdinalDetails(st.ID)
				if err != nil {
					t.Fatalf("ordinals: %v", err)
				}
				if len(ordinals) != 1 {
					t.Fatalf("ordinal consumed on rejection: %+v", ordinals)
				}
			})
		}
	}
}

func TestApplyCreditConcurrentReplayWritesOneApplication(t *testing.T) {
	for name, repo := range creditStores(t) {
		t.Run(name, func(t *testing.T) {
			alloc, apply, _ := creditApplyStore(t, repo)
			activateCreditTest(t, repo, "test")
			st := seedStudent(t, repo)
			now := time.Now().UTC()
			seedMasteredTarget(t, repo, st.ID, "target", now.AddDate(0, 0, -1))
			if _, _, err := alloc.AllocateCredit(st.ID, "attempt-race", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
				t.Fatalf("allocate: %v", err)
			}
			req := CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-race", Epoch: 1, Cohort: "test", Now: now,
			}
			type result struct {
				app CreditApplication
				err error
			}
			results := make(chan result, 2)
			for i := 0; i < 2; i++ {
				go func() {
					app, err := apply.ApplyCredit(req)
					results <- result{app, err}
				}()
			}
			applied, replayed := 0, 0
			for i := 0; i < 2; i++ {
				res := <-results
				if res.err != nil {
					t.Fatalf("apply: %v", res.err)
				}
				if res.app.Replayed {
					replayed++
				} else if res.app.Applied {
					applied++
				}
			}
			if applied != 1 || replayed != 1 {
				t.Fatalf("applied=%d replayed=%d, want 1 and 1", applied, replayed)
			}
			apps, err := apply.CreditApplications(st.ID, "target", 10)
			if err != nil {
				t.Fatalf("applications: %v", err)
			}
			if len(apps) != 1 {
				t.Fatalf("applications = %d, want 1", len(apps))
			}
		})
	}
}

func creditFileHash(t *testing.T, store *SQLiteStore, path string) string {
	t.Helper()
	if _, err := store.db.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		t.Fatalf("checkpoint: %v", err)
	}
	sum, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read database: %v", err)
	}
	hash := sha256.Sum256(sum)
	return hex.EncodeToString(hash[:])
}

func TestPreviewCreditLeavesTheSQLiteFileUnchanged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fire.db")
	store, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	defer store.Close()
	alloc, _, preview := creditApplyStore(t, store)
	st := seedStudent(t, store)
	now := time.Now().UTC()
	seedMasteredTarget(t, store, st.ID, "target", now.AddDate(0, 0, -1))
	if _, _, err := alloc.AllocateCredit(st.ID, "attempt-file", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	before := creditFileHash(t, store, path)
	res, err := preview.PreviewCredit(CreditApplyRequest{
		StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
		SourceAttemptID: "attempt-file", Epoch: 1, Now: now,
	})
	if err != nil {
		t.Fatalf("preview: %v", err)
	}
	if !res.Applied {
		t.Fatalf("preview = %+v, want an applicable row", res)
	}
	after := creditFileHash(t, store, path)
	if before != after {
		t.Fatalf("preview changed the database file: %s -> %s", before, after)
	}
}

func TestExistingCreditLedgerSurvivesMigration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fire-ledger.db")
	store, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	alloc, _, _ := creditApplyStore(t, store)
	st := seedStudent(t, store)
	now := time.Now().UTC()
	seedMasteredTarget(t, store, st.ID, "target", now.AddDate(0, 0, -1))
	if _, _, err := alloc.AllocateCredit(st.ID, "attempt-ledger", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	reopened, err := NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer reopened.Close()
	ordinals, err := reopened.AllBankedOrdinalDetails()
	if err != nil {
		t.Fatalf("ordinals: %v", err)
	}
	if len(ordinals) != 1 || ordinals[0].StudentID != st.ID || ordinals[0].ConceptID != "target" {
		t.Fatalf("existing ledger row did not survive migration: %+v", ordinals)
	}
	apps, err := reopened.CreditApplications(st.ID, "target", 10)
	if err != nil {
		t.Fatalf("applications: %v", err)
	}
	if len(apps) != 0 {
		t.Fatalf("migration created applications: %+v", apps)
	}
	res, err := reopened.PreviewCredit(CreditApplyRequest{
		StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
		SourceAttemptID: "attempt-ledger", Epoch: 1, Now: now,
	})
	if err != nil {
		t.Fatalf("preview after migration: %v", err)
	}
	if !res.Applied || res.Ordinal != ordinals[0].Ordinal {
		t.Fatalf("preview after migration = %+v, ordinals %+v", res, ordinals)
	}
}

func TestDryRunFeedDoesNotWrite(t *testing.T) {
	for name, repo := range creditStores(t) {
		t.Run(name, func(t *testing.T) {
			alloc, apply, _ := creditApplyStore(t, repo)
			st := seedStudent(t, repo)
			now := time.Now().UTC()
			seedMasteredTarget(t, repo, st.ID, "target", now.AddDate(0, 0, -1))
			if _, _, err := alloc.AllocateCredit(st.ID, "attempt-preview", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
				t.Fatalf("allocate: %v", err)
			}
			ordinals, err := repo.(CreditOrdinalStore).AllBankedOrdinalDetails()
			if err != nil {
				t.Fatalf("ordinals: %v", err)
			}
			if len(ordinals) == 0 {
				t.Fatal("expected a banked ordinal")
			}
			apps, err := apply.CreditApplications(st.ID, "target", 10)
			if err != nil {
				t.Fatalf("applications: %v", err)
			}
			if len(apps) != 0 {
				t.Fatalf("dry-run feed wrote applications: %+v", apps)
			}
		})
	}
}

func TestApplyCreditRequiresMatchingDurableActivation(t *testing.T) {
	for name, repo := range creditStores(t) {
		t.Run(name, func(t *testing.T) {
			alloc, apply, _ := creditApplyStore(t, repo)
			st := seedStudent(t, repo)
			now := time.Now().UTC()
			seedMasteredTarget(t, repo, st.ID, "target", now.AddDate(0, 0, -1))
			if _, _, err := alloc.AllocateCredit(st.ID, "attempt-gate", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
				t.Fatalf("allocate: %v", err)
			}
			// Activation is global, and the PostgreSQL store is shared across
			// tests in this package. Start from a known shadow row so this test
			// does not inherit an active cohort written by an earlier test.
			if _, err := repo.(FireActivationStore).SetFireActivation(FireActivation{
				Mode: "shadow", AcknowledgedBy: "test", UpdatedAt: now,
			}); err != nil {
				t.Fatalf("reset activation: %v", err)
			}
			req := CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-gate", Epoch: 1, Cohort: "test", Now: now,
			}
			res, err := apply.ApplyCredit(req)
			if err != nil {
				t.Fatalf("apply without activation: %v", err)
			}
			if res.Applied || res.Reason != "not_active" {
				t.Fatalf("result = %+v, want not_active", res)
			}
			activateCreditTest(t, repo, "other")
			res, err = apply.ApplyCredit(req)
			if err != nil {
				t.Fatalf("apply wrong cohort: %v", err)
			}
			if res.Applied || res.Reason != "cohort_mismatch" {
				t.Fatalf("result = %+v, want cohort_mismatch", res)
			}
			activateCreditTest(t, repo, "test")
			res, err = apply.ApplyCredit(req)
			if err != nil {
				t.Fatalf("apply active: %v", err)
			}
			if !res.Applied {
				t.Fatalf("result = %+v, want applied", res)
			}
		})
	}
}

func TestApplyCreditStopsAfterDeactivationAndStillReplays(t *testing.T) {
	for name, repo := range creditStores(t) {
		t.Run(name, func(t *testing.T) {
			alloc, apply, _ := creditApplyStore(t, repo)
			activateCreditTest(t, repo, "test")
			st := seedStudent(t, repo)
			now := time.Now().UTC()
			seedMasteredTarget(t, repo, st.ID, "target", now.AddDate(0, 0, -1))
			if _, _, err := alloc.AllocateCredit(st.ID, "attempt-1", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
				t.Fatalf("allocate first: %v", err)
			}
			first, err := apply.ApplyCredit(CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-1", Epoch: 1, Cohort: "test", Now: now,
			})
			if err != nil || !first.Applied {
				t.Fatalf("first apply = %+v, %v", first, err)
			}
			if _, _, err := alloc.AllocateCredit(st.ID, "attempt-2", 1, []CreditCandidate{{ConceptID: "target", RawUnits: repetition.Micro}}); err != nil {
				t.Fatalf("allocate second: %v", err)
			}
			if _, err := repo.(FireActivationStore).SetFireActivation(FireActivation{
				Mode: "shadow", AcknowledgedBy: "test", UpdatedAt: now,
			}); err != nil {
				t.Fatalf("deactivate: %v", err)
			}
			blocked, err := apply.ApplyCredit(CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-2", Epoch: 1, Cohort: "test", Now: now,
			})
			if err != nil {
				t.Fatalf("apply after deactivation: %v", err)
			}
			if blocked.Applied || blocked.Reason != "not_active" {
				t.Fatalf("result = %+v, want not_active", blocked)
			}
			replay, err := apply.ApplyCredit(CreditApplyRequest{
				StudentID: st.ID, TargetConceptID: "target", SourceConceptID: "source",
				SourceAttemptID: "attempt-1", Epoch: 1, Cohort: "test", Now: now,
			})
			if err != nil || !replay.Replayed {
				t.Fatalf("replay = %+v, %v", replay, err)
			}
			apps, err := apply.CreditApplications(st.ID, "target", 10)
			if err != nil {
				t.Fatalf("applications: %v", err)
			}
			if len(apps) != 1 {
				t.Fatalf("applications = %d, want 1", len(apps))
			}
		})
	}
}

func TestCreditBaselineComparesInstantsNotStringShapes(t *testing.T) {
	cases := []struct {
		name     string
		baseline string
		actual   string
		want     bool
	}{
		{"empty baseline", "", "2026-10-29T09:00:00Z", true},
		{"same RFC3339", "2026-10-29T09:00:00Z", "2026-10-29T09:00:00Z", true},
		{"legacy SQLite shape", "2026-10-29 09:00:00", "2026-10-29T09:00:00Z", true},
		{"same instant with offset", "2026-10-29T11:00:00+02:00", "2026-10-29T09:00:00Z", true},
		{"different instant", "2026-10-29T09:00:00Z", "2026-10-30T09:00:00Z", false},
		{"unparseable falls back to text", "not-a-date", "2026-10-29T09:00:00Z", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := creditBaselineMatches(tc.baseline, tc.actual); got != tc.want {
				t.Fatalf("creditBaselineMatches(%q, %q) = %v, want %v", tc.baseline, tc.actual, got, tc.want)
			}
		})
	}
}
