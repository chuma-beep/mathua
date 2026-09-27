package repair

import (
	"bytes"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

func testRepo(t *testing.T) (*storage.SQLiteStore, *concepts.DAG, string, string) {
	t.Helper()
	dag, err := concepts.Build([]concepts.Concept{{
		ID: "t", Label: "T", Domain: "d", GradingType: "numeric",
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
	}})
	if err != nil {
		t.Fatalf("build dag: %v", err)
	}
	repo, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	stu, err := repo.CreateStudent("S")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	sess, err := repo.CreateSession(stu.ID)
	if err != nil {
		t.Fatalf("session: %v", err)
	}
	return repo, dag, stu.ID, sess.ID
}

func TestRepair_FlipsFalseNegative(t *testing.T) {
	repo, dag, sid, sess := testRepo(t)
	base := time.Now().UTC().Add(-time.Hour)

	// False negative: the old numeric grader rejected the spoken mixed number.
	must(t, repo.RecordAttempt(storage.AttemptEntry{
		SessionID: sess, StudentID: sid, ConceptID: "t",
		Answer: "4 and 1/10", Expected: "4 1/10", Correct: false,
		ElapsedSeconds: 5, Timestamp: base,
	}))
	// Genuine miss: must stay incorrect after repair.
	must(t, repo.RecordAttempt(storage.AttemptEntry{
		SessionID: sess, StudentID: sid, ConceptID: "t",
		Answer: "5", Expected: "4 1/10", Correct: false,
		ElapsedSeconds: 9, Timestamp: base.Add(time.Minute),
	}))

	reg := generator.NewRegistry()
	eng := engine.New(repo, dag, reg, nil, nil)
	defer eng.Close()

	var buf bytes.Buffer
	sum, err := Run(eng, repo, dag, true, &buf)
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}
	if sum.Repaired != 1 {
		t.Fatalf("dry repaired=%d, want 1\n%s", sum.Repaired, buf.String())
	}
	if all, _ := repo.GetAllAttempts(); all[0].Correct {
		t.Fatal("dry run must not write")
	}

	buf.Reset()
	sum, err = Run(eng, repo, dag, false, &buf)
	if err != nil {
		t.Fatalf("repair: %v", err)
	}
	if sum.Repaired != 1 || sum.ProgressWrites != 1 {
		t.Fatalf("repair repaired=%d writes=%d, want 1/1\n%s", sum.Repaired, sum.ProgressWrites, buf.String())
	}

	all, err := repo.GetAllAttempts()
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	var flipped, wrong int
	for _, a := range all {
		if a.Answer == "4 and 1/10" && a.Correct {
			flipped++
		}
		if a.Answer == "5" && !a.Correct {
			wrong++
		}
	}
	if flipped != 1 {
		t.Errorf("false negative not flipped: %+v", all)
	}
	if wrong != 1 {
		t.Errorf("genuine miss was changed: %+v", all)
	}

	// Rebuilt progress reflects the corrected history (attempts counted).
	p, err := repo.GetProgress(sid, "t")
	if err != nil {
		t.Fatalf("progress: %v", err)
	}
	if p == nil || p.Attempts != 2 {
		t.Fatalf("progress not rebuilt: %+v", p)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
}
