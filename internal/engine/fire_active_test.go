package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/repetition"
	"github.com/chuma-beep/mathua/internal/storage"
)

func activeFireEngine(t *testing.T) (*Engine, storage.Repository) {
	t.Helper()
	e, repo := fireEngine(t)
	// The student is created by the caller; the config is replaced after New so
	// the allowlist can name the real ID.
	return e, repo
}

func TestFIReActiveAppliesOneOrdinalAndDryRunWritesNothing(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("active learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	e.fireConfig = fire.Config{
		Mode:   fire.Active,
		Cohort: "test",
		Allow:  map[string]struct{}{st.ID: {}},
	}
	e.fireMode = fire.Active

	now := time.Now().UTC()
	if _, err := e.SetFIReActivation(storage.FireActivation{
		Mode: string(fire.Active), Cohort: "test", AcknowledgedBy: "test",
		AcknowledgedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("activate: %v", err)
	}
	due := now.AddDate(0, 0, -1)
	stability := 4.0
	difficulty := 5.0
	last := due.AddDate(0, 0, -10)
	if err := repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "prereq", Status: string(mastery.StatusMastered),
		Streak: 3, Attempts: 5, LastReviewed: &last, NextReviewDue: &due,
		SM2Interval: 10, SM2EFactor: 2.5, SM2Repetitions: 3,
		FSRSStability: &stability, FSRSDifficulty: &difficulty, FSRSLastReview: &last,
	}); err != nil {
		t.Fatalf("seed target: %v", err)
	}
	credit := repo.(storage.CreditStore)
	if _, _, err := credit.AllocateCredit(st.ID, "attempt-active", 1, []storage.CreditCandidate{{ConceptID: "prereq", RawUnits: repetition.Micro}}); err != nil {
		t.Fatalf("allocate: %v", err)
	}

	report, err := e.FIReDryRun(now)
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}
	if report.Checksum == "" || report.Revision == "" {
		t.Fatalf("report missing revision or checksum: %+v", report)
	}
	if report.Totals.Targets != 1 || report.Totals.Eligible != 1 {
		t.Fatalf("report totals = %+v", report.Totals)
	}
	before, err := repo.GetProgress(st.ID, "prereq")
	if err != nil {
		t.Fatalf("progress: %v", err)
	}
	if before.NextReviewDue == nil || before.NextReviewDue.Unix() != due.Truncate(time.Second).Unix() {
		t.Fatalf("dry run moved the due date: %v", before.NextReviewDue)
	}
	apps, err := repo.(storage.CreditApplyStore).CreditApplications(st.ID, "prereq", 10)
	if err != nil {
		t.Fatalf("applications: %v", err)
	}
	if len(apps) != 0 {
		t.Fatalf("dry run wrote applications: %+v", apps)
	}

	if err := e.SatisfyImplicitReviews(st.ID); err != nil {
		t.Fatalf("satisfy: %v", err)
	}
	after, err := repo.GetProgress(st.ID, "prereq")
	if err != nil {
		t.Fatalf("progress after: %v", err)
	}
	if after.Status != string(mastery.StatusMastered) || after.Attempts != before.Attempts || after.Streak != before.Streak {
		t.Fatalf("application changed learner state: before=%+v after=%+v", before, after)
	}
	if after.NextReviewDue == nil || !after.NextReviewDue.After(due.Truncate(time.Second)) {
		t.Fatalf("due did not move: %v", after.NextReviewDue)
	}
	xp, _, err := repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("xp: %v", err)
	}
	if xp != 0 {
		t.Fatalf("implicit review granted xp: %d", xp)
	}
}

func TestFIReActiveRejectsFallbackSourceAndOutsideCohort(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("cohort learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	e.fireConfig = fire.Config{
		Mode:   fire.Active,
		Cohort: "other",
		Allow:  map[string]struct{}{"someone-else": {}},
	}
	e.fireMode = fire.Active
	if err := e.allocateImplicitCreditErr(gradedAttempt{
		studentID: st.ID, conceptID: "adv", correct: true, source: "practice",
	}); err != nil {
		t.Fatalf("allocate outside cohort: %v", err)
	}
	e.fireConfig.Allow[st.ID] = struct{}{}
	if err := e.allocateImplicitCreditErr(gradedAttempt{
		studentID: st.ID, conceptID: "adv", correct: true, source: "quiz",
	}); err != nil {
		t.Fatalf("allocate quiz source: %v", err)
	}
	ordinals, err := repo.(storage.CreditOrdinalStore).BankedOrdinalDetails(st.ID)
	if err != nil {
		t.Fatalf("ordinals: %v", err)
	}
	if len(ordinals) != 0 {
		t.Fatalf("ineligible source or cohort banked credit: %+v", ordinals)
	}
}

func TestFIReActiveDoesNotStarveADueTargetBehindFutureTargets(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("starved learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	e.fireConfig = fire.Config{
		Mode:   fire.Active,
		Cohort: "test",
		Allow:  map[string]struct{}{st.ID: {}},
	}
	e.fireMode = fire.Active
	now := time.Now().UTC()
	if _, err := e.SetFIReActivation(storage.FireActivation{
		Mode: string(fire.Active), Cohort: "test", AcknowledgedBy: "test",
		AcknowledgedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("activate: %v", err)
	}
	credit := repo.(storage.CreditStore)
	seed := func(conceptID string, due time.Time) {
		last := due.AddDate(0, 0, -10)
		stability := 4.0
		difficulty := 5.0
		if err := repo.UpsertProgress(&storage.ConceptProgress{
			StudentID: st.ID, ConceptID: conceptID, Status: string(mastery.StatusMastered),
			Attempts: 5, LastReviewed: &last, NextReviewDue: &due,
			SM2Interval: 10, SM2EFactor: 2.5, SM2Repetitions: 3,
			FSRSStability: &stability, FSRSDifficulty: &difficulty, FSRSLastReview: &last,
		}); err != nil {
			t.Fatalf("seed %s: %v", conceptID, err)
		}
		if _, _, err := credit.AllocateCredit(st.ID, "attempt-"+conceptID, 1, []storage.CreditCandidate{{ConceptID: conceptID, RawUnits: repetition.Micro}}); err != nil {
			t.Fatalf("allocate %s: %v", conceptID, err)
		}
	}
	// More future targets than the old per-call attempt cap, all sorting before
	// the one target that is actually due.
	for i := 0; i < 33; i++ {
		seed(fmt.Sprintf("future-%02d", i), now.AddDate(0, 0, 30))
	}
	due := now.AddDate(0, 0, -1)
	seed("zzz-due", due)

	if err := e.SatisfyImplicitReviews(st.ID); err != nil {
		t.Fatalf("satisfy: %v", err)
	}
	after, err := repo.GetProgress(st.ID, "zzz-due")
	if err != nil {
		t.Fatalf("due progress: %v", err)
	}
	if after.NextReviewDue == nil || !after.NextReviewDue.After(due.Truncate(time.Second)) {
		t.Fatalf("due target was starved: %v", after.NextReviewDue)
	}
	future, err := repo.GetProgress(st.ID, "future-00")
	if err != nil {
		t.Fatalf("future progress: %v", err)
	}
	if future.NextReviewDue == nil || !future.NextReviewDue.After(now.AddDate(0, 0, 29)) {
		t.Fatalf("future target was moved: %v", future.NextReviewDue)
	}
}
