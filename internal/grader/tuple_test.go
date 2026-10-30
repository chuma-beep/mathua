package grader

import "testing"

func TestTupleAcceptsDegreeVariants(t *testing.T) {
	g := tupleGrader{}
	for _, answer := range []string{"(1,90°)", "(1, 90)", "(1,90 deg)", "(1,90 degrees)"} {
		if r := g.grade("(1,90°)", answer); !r.Correct {
			t.Errorf("grade(%q) = incorrect, want correct", answer)
		}
	}
	if r := g.grade("(1,90°)", "(1,91)"); r.Correct {
		t.Error("a different angle graded correct")
	}
	if r := g.grade("(north,east)", "(north, east)"); !r.Correct {
		t.Error("a text tuple with spaces graded incorrect")
	}
	if r := g.grade("(north,east)", "(south,east)"); r.Correct {
		t.Error("a different text tuple graded correct")
	}
}
