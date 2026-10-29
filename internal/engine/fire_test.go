package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/repetition"
	"github.com/chuma-beep/mathua/internal/storage"
)

func fireEngine(t *testing.T, opts ...Option) (*Engine, storage.Repository) {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{
			ID:               "prereq",
			Label:            "Prerequisite Concept",
			Domain:           "algebra",
			GradingType:      "numeric",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
		{
			ID:               "adv",
			Label:            "Advanced Concept",
			Domain:           "calculus",
			GradingType:      "numeric",
			Prerequisites:    []string{"prereq"},
			Encompasses:      []string{"prereq"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
		{
			ID:               "standalone",
			Label:            "Standalone Concept",
			Domain:           "geometry",
			GradingType:      "numeric",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}

	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	reg := generator.NewRegistry()
	reg.Register("prereq", &testGen{answer: "10"})
	reg.Register("adv", &testGen{answer: "20"})
	reg.Register("standalone", &testGen{answer: "30"})

	return New(store, d, reg, nil, nil, opts...), store
}

func TestEngine_FIRe_CompletedAttemptAllocatesImplicitCredit(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("bob")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	// Mark prereq and standalone as mastered so adv is the single available concept
	now := time.Now().UTC()
	err = repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "prereq",
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		LastReviewed:   &now,
		SM2Interval:    10,
		SM2EFactor:     2.5,
		SM2Repetitions: 2,
	})
	if err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}
	err = repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "standalone",
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		LastReviewed:   &now,
		SM2Interval:    10,
		SM2EFactor:     2.5,
		SM2Repetitions: 2,
	})
	if err != nil {
		t.Fatalf("upsert standalone: %v", err)
	}

	sess, err := repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	// Fetch question on adv and answer correctly
	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if q.ConceptID != "adv" {
		t.Fatalf("expected adv concept, got %s", q.ConceptID)
	}

	res, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "20", 5.0)
	if err != nil {
		t.Fatalf("submit answer: %v", err)
	}
	if !res.Correct {
		t.Fatalf("expected correct answer")
	}

	// Verify that implicit credit was allocated and banked for prereq
	cs := repo.(storage.CreditStore)
	acc, err := cs.AccountOf(st.ID, "prereq", 0)
	if err != nil {
		t.Fatalf("account of prereq: %v", err)
	}
	if acc.RetainedUnits != repetition.Micro {
		t.Fatalf("expected %d retained units, got %d", repetition.Micro, acc.RetainedUnits)
	}
	if acc.Pending() != 1 {
		t.Fatalf("expected 1 pending banked ordinal, got %d", acc.Pending())
	}
}

func TestEngine_FIRe_BankedCreditDoesNotSuppressOrAdvance(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("charlie")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	now := time.Now().UTC()
	past := now.AddDate(0, 0, -5)
	if err := repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "prereq",
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		LastReviewed:   &past,
		NextReviewDue:  &past,
		SM2Interval:    5,
		SM2EFactor:     2.5,
		SM2Repetitions: 1,
	}); err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}

	sess, err := repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	cs := repo.(storage.CreditStore)
	if _, _, err := cs.AllocateCredit(st.ID, "attempt:fire:1", 0, []storage.CreditCandidate{
		{ConceptID: "prereq", RawUnits: repetition.Micro},
	}); err != nil {
		t.Fatalf("allocate credit: %v", err)
	}

	due, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if due != 1 {
		t.Fatalf("shadow credit changed the due count: got %d, want 1", due)
	}
	rq, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next review: %v", err)
	}
	if rq == nil || rq.ConceptID != "prereq" {
		t.Fatalf("shadow credit suppressed the explicit review: got %+v", rq)
	}
	p, err := repo.GetProgress(st.ID, "prereq")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if p.NextReviewDue == nil || p.NextReviewDue.Unix() != past.Unix() {
		t.Fatalf("shadow credit moved the due date: got %v, want %v", p.NextReviewDue, past)
	}
	if p.SM2Interval != 5 {
		t.Fatalf("shadow credit changed the interval: got %d, want 5", p.SM2Interval)
	}
}

func TestEngine_FIRe_NonEncompassedConceptIsNotSuppressed(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("diana")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	now := time.Now().UTC()
	past := now.AddDate(0, 0, -5)

	// Standalone concept has no encompasses relationships
	err = repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "standalone",
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		LastReviewed:   &past,
		NextReviewDue:  &past,
		SM2Interval:    5,
		SM2EFactor:     2.5,
		SM2Repetitions: 1,
	})
	if err != nil {
		t.Fatalf("upsert standalone: %v", err)
	}

	// Bank credit for standalone
	cs := repo.(storage.CreditStore)
	_, _, err = cs.AllocateCredit(st.ID, "attempt:fire:standalone", 0, []storage.CreditCandidate{
		{ConceptID: "standalone", RawUnits: repetition.Micro},
	})
	if err != nil {
		t.Fatalf("allocate credit: %v", err)
	}

	// Since standalone is not encompassed, explicit review must not be suppressed
	due, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if due != 1 {
		t.Fatalf("expected 1 due review for non-encompassed concept, got %d", due)
	}

	sess, err := repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	rq, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil {
		t.Fatalf("next review: %v", err)
	}
	if rq == nil || rq.ConceptID != "standalone" {
		t.Fatalf("expected standalone review question served, got %+v", rq)
	}
}

func TestEngine_FIRe_SubmitAnswerDoesNotMoveEncompassedReview(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("diana")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}

	now := time.Now().UTC()
	past := now.AddDate(0, 0, -10)

	// Prereq is mastered and currently due for review
	err = repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "prereq",
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		LastReviewed:   &past,
		NextReviewDue:  &past,
		SM2Interval:    5,
		SM2EFactor:     2.5,
		SM2Repetitions: 1,
	})
	if err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}

	// Standalone is mastered and not due
	err = repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:      st.ID,
		ConceptID:      "standalone",
		Status:         string(mastery.StatusMastered),
		Streak:         3,
		LastReviewed:   &now,
		SM2Interval:    10,
		SM2EFactor:     2.5,
		SM2Repetitions: 2,
	})
	if err != nil {
		t.Fatalf("upsert standalone: %v", err)
	}

	// Submit a correct attempt on adv (which encompasses prereq)
	diff := 0.5
	res, err := e.SubmitQuizAnswer(st.ID, "adv", "20", "20", 5.0, "what is 20?", &diff)
	if err != nil {
		t.Fatalf("submit quiz answer: %v", err)
	}
	if !res.Correct {
		t.Fatalf("expected correct answer")
	}

	// Shadow records the credit but must not move the encompassed review.
	p, err := repo.GetProgress(st.ID, "prereq")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if p.NextReviewDue == nil || p.NextReviewDue.Unix() != past.Unix() {
		t.Fatalf("shadow credit moved the due date: got %v, want %v", p.NextReviewDue, past)
	}
	if p.SM2Interval != 5 {
		t.Fatalf("shadow credit changed the interval: got %d, want 5", p.SM2Interval)
	}
}

func TestEngine_FIRe_OffDoesNotAllocateOrQueueCredit(t *testing.T) {
	e, repo := fireEngine(t, WithFIReMode(fire.Off))
	if e.fireMode != fire.Off {
		t.Fatalf("fire mode = %s, want off", e.fireMode)
	}
	for _, effect := range e.studyEffectsFor("practice", "adv", "attempt-off", true, false, nil) {
		if effect.Effect == studyEffectImplicitCredit {
			t.Fatal("off mode queued an implicit_credit effect")
		}
	}

	st, err := e.CreateStudent("off-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	if err := repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "prereq", Status: string(mastery.StatusMastered),
		Streak: 3, LastReviewed: &now, SM2Interval: 10, SM2EFactor: 2.5, SM2Repetitions: 2,
	}); err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}
	if err := repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "standalone", Status: string(mastery.StatusMastered),
		Streak: 3, LastReviewed: &now, SM2Interval: 10, SM2EFactor: 2.5, SM2Repetitions: 2,
	}); err != nil {
		t.Fatalf("upsert standalone: %v", err)
	}
	sess, err := repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextQuestion(sess.ID, st.ID)
	if err != nil || q == nil || q.ConceptID != "adv" {
		t.Fatalf("next question = %+v, %v; want adv", q, err)
	}
	if _, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "20", 5.0); err != nil {
		t.Fatalf("submit: %v", err)
	}
	cs := repo.(storage.CreditStore)
	acc, err := cs.AccountOf(st.ID, "prereq", 0)
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	if acc.RetainedUnits != 0 {
		t.Fatalf("off mode retained %d units, want 0", acc.RetainedUnits)
	}
}

func TestEngine_FIRe_ShadowServingDoesNotWriteProgress(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("serve-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	past := now.AddDate(0, 0, -5)
	if err := repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "prereq", Status: string(mastery.StatusMastered),
		Streak: 3, LastReviewed: &past, NextReviewDue: &past,
		SM2Interval: 5, SM2EFactor: 2.5, SM2Repetitions: 1,
	}); err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}
	cs := repo.(storage.CreditStore)
	if _, _, err := cs.AllocateCredit(st.ID, "attempt:serve", 0, []storage.CreditCandidate{
		{ConceptID: "prereq", RawUnits: repetition.Micro},
	}); err != nil {
		t.Fatalf("allocate credit: %v", err)
	}
	before, err := repo.GetProgress(st.ID, "prereq")
	if err != nil || before == nil || before.NextReviewDue == nil {
		t.Fatalf("read before serving: %v", err)
	}
	beforeDue := *before.NextReviewDue
	beforeInterval := before.SM2Interval

	sess, err := repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if _, err := e.NextQuestion(sess.ID, st.ID); err != nil {
		t.Fatalf("next question: %v", err)
	}
	if _, err := e.ReviewsDue(st.ID); err != nil {
		t.Fatalf("reviews due: %v", err)
	}
	if _, err := e.RecommendNext(st.ID, nil); err != nil {
		t.Fatalf("recommend next: %v", err)
	}

	after, err := repo.GetProgress(st.ID, "prereq")
	if err != nil || after == nil || after.NextReviewDue == nil {
		t.Fatalf("read after serving: %v", err)
	}
	if after.NextReviewDue.Unix() != beforeDue.Unix() {
		t.Fatalf("serving moved the due date: got %v, want %v", after.NextReviewDue, beforeDue)
	}
	if after.SM2Interval != beforeInterval {
		t.Fatalf("serving changed the interval: got %d, want %d", after.SM2Interval, beforeInterval)
	}
}

func TestEngine_FIRe_OffDrainsQueuedCreditEffectWithoutWriting(t *testing.T) {
	e, repo := fireEngine(t, WithFIReMode(fire.Off))
	st, err := e.CreateStudent("drain-off")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := e.applyStudyEffect(storage.StudyAnswerEffect{
		Effect:    studyEffectImplicitCredit,
		StudentID: st.ID,
		PayloadJSON: studyEffectJSON(studyEffectPayload{
			ConceptID: "adv", AttemptID: "queued-off", Correct: true,
		}),
	}); err != nil {
		t.Fatalf("drain queued effect: %v", err)
	}
	cs := repo.(storage.CreditStore)
	acc, err := cs.AccountOf(st.ID, "prereq", 0)
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	if acc.RetainedUnits != 0 {
		t.Fatalf("off mode wrote %d retained units, want 0", acc.RetainedUnits)
	}
}

func TestEngine_FIRe_ShadowRetryAllocatesOnceWithoutProjecting(t *testing.T) {
	e, repo := fireEngine(t)
	st, err := e.CreateStudent("retry-shadow")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	past := now.AddDate(0, 0, -5)
	if err := repo.UpsertProgress(&storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "prereq", Status: string(mastery.StatusMastered),
		Streak: 3, LastReviewed: &past, NextReviewDue: &past,
		SM2Interval: 5, SM2EFactor: 2.5, SM2Repetitions: 1,
	}); err != nil {
		t.Fatalf("upsert prereq: %v", err)
	}
	before, err := repo.GetProgress(st.ID, "prereq")
	if err != nil || before == nil || before.NextReviewDue == nil {
		t.Fatalf("read before: %v", err)
	}
	beforeDue := before.NextReviewDue.Unix()
	effect := storage.StudyAnswerEffect{
		Effect:    studyEffectImplicitCredit,
		StudentID: st.ID,
		PayloadJSON: studyEffectJSON(studyEffectPayload{
			ConceptID: "adv", AttemptID: "retry-shadow", Correct: true,
		}),
	}
	for i := 0; i < 3; i++ {
		if err := e.applyStudyEffect(effect); err != nil {
			t.Fatalf("retry %d: %v", i, err)
		}
	}
	cs := repo.(storage.CreditStore)
	acc, err := cs.AccountOf(st.ID, "prereq", 0)
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	if acc.RetainedUnits != repetition.Micro {
		t.Fatalf("retries retained %d units, want one Micro", acc.RetainedUnits)
	}
	after, err := repo.GetProgress(st.ID, "prereq")
	if err != nil || after == nil || after.NextReviewDue == nil {
		t.Fatalf("read after: %v", err)
	}
	if after.NextReviewDue.Unix() != beforeDue {
		t.Fatalf("retry moved the due date: got %v, want unix %d", after.NextReviewDue, beforeDue)
	}
}
