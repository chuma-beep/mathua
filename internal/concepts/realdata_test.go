package concepts

import "testing"

// The curriculum graph is only useful if a collision in it is a build failure.
// LoadDir already checks ids, prerequisite resolution and acyclicity; this adds
// the title-level check that caught 19 real cross-domain duplicates.
func TestRealCurriculumGraphIsValid(t *testing.T) {
	dag, err := LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load curriculum: %v", err)
	}
	if dag.Count() < 1000 {
		t.Fatalf("only %d concepts loaded; the graph check would pass vacuously", dag.Count())
	}
	allow := ReviewedDuplicateLabels()
	if err := CheckReviewedDuplicateLabelsExist(dag, allow); err != nil {
		t.Fatal(err)
	}
	if err := CheckDuplicateLabels(dag, allow); err != nil {
		t.Fatal(err)
	}
}
