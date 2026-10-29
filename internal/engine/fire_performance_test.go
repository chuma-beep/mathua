package engine

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/fire"
	"github.com/chuma-beep/mathua/internal/generator/all"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/repetition"
	"github.com/chuma-beep/mathua/internal/storage"
)

// TestFIRePerformance is the opt-in, production-shaped benchmark for the active
// application path and the scheduler read path. It is skipped in normal CI.
//
//	MATHUA_BENCH=1 go test ./internal/engine -run TestFIRePerformance -v
//	MATHUA_BENCH_DSN=postgres://... MATHUA_BENCH=1 go test ./internal/engine -run TestFIRePerformance -v
func TestFIRePerformance(t *testing.T) {
	if os.Getenv("MATHUA_BENCH") == "" {
		t.Skip("set MATHUA_BENCH=1 to run the targeted performance test")
	}
	students := envInt("MATHUA_BENCH_STUDENTS", 120)
	workers := envInt("MATHUA_BENCH_WORKERS", 16)
	if students < 30 {
		t.Fatalf("MATHUA_BENCH_STUDENTS must be at least 30, got %d", students)
	}
	if workers < 1 {
		t.Fatalf("MATHUA_BENCH_WORKERS must be positive, got %d", workers)
	}
	oldLog := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(oldLog)
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	store, engineName, closeStore := benchStore(t)
	defer closeStore()
	e := New(store, dag, all.Registry(), nil, nil)
	now := time.Now().UTC()
	if _, err := e.SetFIReActivation(storage.FireActivation{
		Mode: string(fire.Active), Cohort: "bench", AcknowledgedBy: "bench",
		AcknowledgedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("activate: %v", err)
	}
	var root *concepts.Concept
	for _, c := range dag.Order() {
		if len(c.Prerequisites) == 0 {
			if _, err := all.Registry().Get(c.ID); err == nil {
				root = c
				break
			}
		}
	}
	if root == nil {
		t.Fatal("no root concept with a generator")
	}
	seedTarget := func(st *storage.Student, conceptID string) {
		due := now.Add(-time.Hour)
		last := now.AddDate(0, 0, -10)
		stability := 10.0
		difficulty := 5.0
		if err := e.repo.UpsertProgress(&storage.ConceptProgress{
			StudentID: st.ID, ConceptID: conceptID, Status: string(mastery.StatusMastered),
			Streak: 1, BestStreak: 5, Attempts: 6, LastAttempted: &last, LastReviewed: &last,
			NextReviewDue: &due, SM2Repetitions: 4, SM2Interval: 30, SM2EFactor: 2.5,
			FSRSStability: &stability, FSRSDifficulty: &difficulty, FSRSLastReview: &last,
		}); err != nil {
			t.Fatalf("seed target: %v", err)
		}
	}
	makeStudents := func(n int, prefix string) []*storage.Student {
		out := make([]*storage.Student, 0, n)
		for i := 0; i < n; i++ {
			st, err := e.CreateStudent(fmt.Sprintf("%s-%03d", prefix, i))
			if err != nil {
				t.Fatalf("create student: %v", err)
			}
			out = append(out, st)
		}
		return out
	}
	applyStudents := makeStudents(students, "fire-apply")
	for i, st := range applyStudents {
		seedTarget(st, root.ID)
		if _, _, err := e.repo.(storage.CreditStore).AllocateCredit(st.ID, fmt.Sprintf("bench-apply-%d", i), 1, []storage.CreditCandidate{{ConceptID: root.ID, RawUnits: repetition.Micro}}); err != nil {
			t.Fatalf("allocate apply: %v", err)
		}
	}
	applyResult := runWorkload(workers, len(applyStudents), func(i int) error {
		st := applyStudents[i]
		_, err := e.repo.(storage.CreditApplyStore).ApplyCredit(storage.CreditApplyRequest{
			StudentID: st.ID, TargetConceptID: root.ID, SourceConceptID: "",
			SourceAttemptID: fmt.Sprintf("bench-apply-direct-%d", i), Epoch: 1,
			Cohort: "bench", Now: now,
		})
		return err
	})
	report(t, "fire direct apply", applyResult)

	allow := make(map[string]struct{}, len(applyStudents))
	for _, st := range applyStudents {
		allow[st.ID] = struct{}{}
	}
	e.fireConfig = fire.Config{Mode: fire.Active, Cohort: "bench", Allow: allow}
	e.fireMode = fire.Active
	readStudents := makeStudents(students, "fire-read")
	for _, st := range readStudents {
		allow[st.ID] = struct{}{}
	}
	for i, st := range readStudents {
		seedTarget(st, root.ID)
		if _, _, err := e.repo.(storage.CreditStore).AllocateCredit(st.ID, fmt.Sprintf("bench-read-%d", i), 1, []storage.CreditCandidate{{ConceptID: root.ID, RawUnits: repetition.Micro}}); err != nil {
			t.Fatalf("allocate read: %v", err)
		}
	}
	readApplyResult := runWorkload(workers, len(readStudents), func(i int) error {
		return e.SatisfyImplicitReviews(readStudents[i].ID)
	})
	report(t, "scheduler read with one application", readApplyResult)

	readOnlyResult := runWorkload(workers, len(readStudents), func(i int) error {
		_, err := e.ReviewsDue(readStudents[i].ID)
		return err
	})
	report(t, "scheduler read with no banked work", readOnlyResult)
	t.Logf("dataset: engine=%s students_per_group=%d workers=%d", engineName, students, workers)
}
