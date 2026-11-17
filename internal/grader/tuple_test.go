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

func TestTupleAcceptsRearrangedParametricEquations(t *testing.T) {
	g := tupleGrader{}
	expected := "x=5+3t,y=-4+t"
	for _, answer := range []string{
		"x = 3t + 5, y = t -4",
		"x=5+3*t,y=-4+t",
		"x = 5 + 3t, y = t - 4",
	} {
		if r := g.grade(expected, answer); !r.Correct {
			t.Errorf("grade(%q) = incorrect, want correct: %s", answer, r.Feedback)
		}
	}
	if r := g.grade(expected, "y = 3t + 5, x = t - 4"); r.Correct {
		t.Error("equations with swapped labels graded correct")
	}
	if r := g.grade(expected, "x = 3t + 5, y = t + 4"); r.Correct {
		t.Error("a different parametric line graded correct")
	}
}
