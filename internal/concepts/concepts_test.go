package concepts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Valid(t *testing.T) {
	d, err := LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if d.Count() != 438 {
		t.Errorf("expected 438 concepts, got %d", d.Count())
	}

	if len(d.Order()) != 438 {
		t.Errorf("expected 438 in topo order, got %d", len(d.Order()))
	}
}

func TestLoad_Empty(t *testing.T) {
	path := writeTemp(t, `[]`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty concepts")
	}
}

func TestLoad_DuplicateID(t *testing.T) {
	path := writeTemp(t, `[
		{"id":"a","label":"A","domain":"d","prerequisites":[],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}},
		{"id":"a","label":"A again","domain":"d","prerequisites":[],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}}
	]`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate concept ID")
	}
}

func TestLoad_EmptyID(t *testing.T) {
	path := writeTemp(t, `[
		{"id":"","label":"no id","domain":"d","prerequisites":[],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}}
	]`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty concept ID")
	}
}

func TestLoad_OrphanPrerequisite(t *testing.T) {
	path := writeTemp(t, `[
		{"id":"a","label":"A","domain":"d","prerequisites":["nonexistent"],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}}
	]`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for orphan prerequisite")
	}
}

func TestLoad_Cycle(t *testing.T) {
	path := writeTemp(t, `[
		{"id":"a","label":"A","domain":"d","prerequisites":["b"],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}},
		{"id":"b","label":"B","domain":"d","prerequisites":["a"],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}}
	]`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for cycle")
	}
}

func TestDAG_Concepts(t *testing.T) {
	d := buildTestDAG(t)
	cs := d.Concepts()
	if len(cs) != 4 {
		t.Errorf("expected 4 concepts, got %d", len(cs))
	}
}

func TestDAG_Order(t *testing.T) {
	d := buildTestDAG(t)
	order := d.Order()
	if len(order) != 4 {
		t.Errorf("expected 4 in order, got %d", len(order))
	}
	for _, c := range order {
		for _, pid := range c.Prerequisites {
			found := false
			for _, pc := range order {
				if pc.ID == pid {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}
}

func TestDAG_Domains(t *testing.T) {
	d := buildTestDAG(t)
	domains := d.Domains()
	if len(domains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(domains))
	}
}

func TestDAG_PrereqsOf(t *testing.T) {
	d := buildTestDAG(t)
	ps := d.PrereqsOf("arith.add.single")
	if len(ps) != 1 || ps[0].ID != "count.basics.objects" {
		t.Errorf("expected 1 prereq count.objects, got %v", ps)
	}
	ps = d.PrereqsOf("count.basics.objects")
	if len(ps) != 0 {
		t.Errorf("expected 0 prereqs for root, got %d", len(ps))
	}
}

func TestDAG_DependentsOf(t *testing.T) {
	d := buildTestDAG(t)
	deps := d.DependentsOf("count.basics.objects")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependents of count.objects, got %d", len(deps))
	}
	ids := make(map[string]bool)
	for _, dep := range deps {
		ids[dep.ID] = true
	}
	if !ids["count.basics.cardinality"] || !ids["arith.add.single"] {
		t.Errorf("unexpected dependents: %v", deps)
	}
}

func TestDAG_Available(t *testing.T) {
	d := buildTestDAG(t)

	// Nothing mastered — only concepts with no prereqs should be available.
	avail := d.Available(map[string]bool{})
	if len(avail) != 1 || avail[0].ID != "count.basics.objects" {
		t.Errorf("expected only count.objects available, got %v", avail)
	}

	// count.objects mastered — its dependents become available.
	avail = d.Available(map[string]bool{"count.basics.objects": true})
	if len(avail) != 2 {
		t.Errorf("expected 2 available after count.objects mastered, got %d", len(avail))
	}

	// Everything mastered — nothing available.
	all := map[string]bool{
		"count.basics.objects":     true,
		"count.basics.cardinality": true,
		"arith.add.single":  true,
		"arith.sub.single":  true,
	}
	avail = d.Available(all)
	if len(avail) != 0 {
		t.Errorf("expected 0 available when all mastered, got %d", len(avail))
	}
}

func TestDAG_Concept(t *testing.T) {
	d := buildTestDAG(t)
	c := d.Concept("arith.add.single")
	if c == nil {
		t.Fatal("expected concept to exist")
	}
	if c.Label != "Single-digit addition" {
		t.Errorf("unexpected label: %q", c.Label)
	}
	if d.Concept("bogus") != nil {
		t.Error("expected nil for nonexistent concept")
	}
}

func writeTemp(t *testing.T, jsonContent string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test_concepts.json")
	if err := os.WriteFile(path, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func buildTestDAG(t *testing.T) *DAG {
	t.Helper()
	raw := []Concept{
		{
			ID: "count.basics.objects", Label: "Count objects 1-10",
			Domain: "counting", Prerequisites: []string{},
			MasteryThreshold: MasteryThreshold{Streak: 5, AvgTimeSeconds: 10},
		},
		{
			ID: "count.basics.cardinality", Label: "Cardinality",
			Domain: "counting", Prerequisites: []string{"count.basics.objects"},
			MasteryThreshold: MasteryThreshold{Streak: 5, AvgTimeSeconds: 10},
		},
		{
			ID: "arith.add.single", Label: "Single-digit addition",
			Domain: "arithmetic", Prerequisites: []string{"count.basics.objects"},
			MasteryThreshold: MasteryThreshold{Streak: 5, AvgTimeSeconds: 6},
		},
		{
			ID: "arith.sub.single", Label: "Single-digit subtraction",
			Domain: "arithmetic", Prerequisites: []string{"arith.add.single"},
			MasteryThreshold: MasteryThreshold{Streak: 5, AvgTimeSeconds: 6},
		},
	}
	// Validate and build
	if err := validate(raw); err != nil {
		t.Fatalf("test DAG should be valid: %v", err)
	}
	return build(raw)
}

func TestLoad_SingleConcept(t *testing.T) {
	path := writeTemp(t, `[
		{"id":"c","label":"C","domain":"d","prerequisites":[],
		 "mastery_threshold":{"streak":1,"avg_time_seconds":1}}
	]`)
	d, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Count() != 1 {
		t.Errorf("expected 1 concept, got %d", d.Count())
	}
}

func Test_MasteryThresholdUnmarshal(t *testing.T) {
	input := `{"streak":7,"avg_time_seconds":6.0}`
	var mt MasteryThreshold
	if err := json.Unmarshal([]byte(input), &mt); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if mt.Streak != 7 || mt.AvgTimeSeconds != 6.0 {
		t.Errorf("unexpected: %+v", mt)
	}
}
