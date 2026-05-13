package planning

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
)

func testPlanner(t *testing.T) *Planner {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Domain: "d", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "b", Domain: "d", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "c", Domain: "d", Prerequisites: []string{"a"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
		{ID: "d", Domain: "d", Prerequisites: []string{"b", "c"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10}},
	})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	dir := t.TempDir()
	coursesPath := filepath.Join(dir, "courses.json")
	os.WriteFile(coursesPath, []byte(`[
		{"id":"test","name":"Test","grade":"1","description":"d","targets":["d"]}
	]`), 0644)
	p, err := Load(coursesPath, d)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return p
}

func TestPlanner_Courses(t *testing.T) {
	p := testPlanner(t)
	if len(p.Courses()) != 1 {
		t.Errorf("expected 1 course, got %d", len(p.Courses()))
	}
}

func TestPrerequisiteChain(t *testing.T) {
	p := testPlanner(t)
	path, err := p.PrerequisiteChain([]string{"d"})
	if err != nil {
		t.Fatalf("chain: %v", err)
	}
	ids := make([]string, len(path.Concepts))
	for i, c := range path.Concepts {
		ids[i] = c.ID
	}
	expected := []string{"a", "b", "c", "d"}
	if len(ids) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, ids)
	}
	for i, e := range expected {
		if ids[i] != e {
			t.Errorf("position %d: expected %q, got %q", i, e, ids[i])
		}
	}
}

func TestPrerequisiteChain_Subset(t *testing.T) {
	p := testPlanner(t)
	path, err := p.PrerequisiteChain([]string{"b"})
	if err != nil {
		t.Fatalf("chain: %v", err)
	}
	ids := make([]string, len(path.Concepts))
	for i, c := range path.Concepts {
		ids[i] = c.ID
	}
	// Should include "a" and "b" but not "c" or "d"
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Errorf("expected [a b], got %v", ids)
	}
}

func TestPathForCourse(t *testing.T) {
	p := testPlanner(t)
	path, err := p.PathForCourse("test")
	if err != nil {
		t.Fatalf("path: %v", err)
	}
	if len(path.Concepts) != 4 {
		t.Errorf("expected 4 concepts in path, got %d", len(path.Concepts))
	}
	if path.Course == nil || path.Course.Name != "Test" {
		t.Errorf("expected course info, got %v", path.Course)
	}
}

func TestPathForCourse_NotFound(t *testing.T) {
	p := testPlanner(t)
	_, err := p.PathForCourse("nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestLoad_RealData(t *testing.T) {
	d, err := concepts.Load("../../data/concepts.json")
	if err != nil {
		t.Skip("concepts not available")
	}
	p, err := Load("../../data/courses.json", d)
	if err != nil {
		t.Fatalf("load real data: %v", err)
	}
	if len(p.Courses()) != 10 {
		t.Errorf("expected 10 courses, got %d", len(p.Courses()))
	}
	// Check that 4th grade chain has reasonable size
	path, err := p.PathForCourse("4")
	if err != nil {
		t.Fatalf("4th grade path: %v", err)
	}
	if len(path.Concepts) < 10 || len(path.Concepts) > 60 {
		t.Errorf("expected 10-60 concepts in 4th grade chain, got %d", len(path.Concepts))
	}
	t.Logf("4th grade chain: %d concepts", len(path.Concepts))
}
