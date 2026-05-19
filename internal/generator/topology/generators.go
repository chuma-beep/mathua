package topology

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("topo.basics.metric", &metricGen{})
	reg.Register("topo.basics.open_closed", &openClosedGen{})
	reg.Register("topo.basics.continuous", &continuousGen{})
}

// ----- 1. metric -----

type metricGen struct{}

func (g *metricGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		question string
		answer   string
		reason   string
	}
	table := []entry{
		{
			"Is the usual distance d(x,y)=|x-y| a metric on R? (yes/no)",
			"yes",
			"d(x,y)=|x-y| satisfies all metric axioms: non-negativity, identity of indiscernibles, symmetry, and triangle inequality.",
		},
		{
			"In a metric space, is it true that d(x,y)=0 if and only if x=y? (yes/no)",
			"yes",
			"This is the identity of indiscernibles axiom: d(x,y)=0 ⇔ x=y.",
		},
		{
			"In a metric space, is d(x,z) ≤ d(x,y)+d(y,z)? (yes/no)",
			"yes",
			"This is the triangle inequality, one of the defining properties of a metric.",
		},
		{
			"Is d(x,y) = (x-y)² a metric on R? (yes/no)",
			"no",
			"It violates the triangle inequality. For example, d(0,2)=4 but d(0,1)+d(1,2)=1+1=2 < 4.",
		},
		{
			"In a metric space, is d(x,y) always non-negative? (yes/no)",
			"yes",
			"Non-negativity is one of the metric axioms: d(x,y) ≥ 0 for all x,y.",
		},
		{
			"Is the discrete metric (d(x,y)=1 if x≠y, 0 if x=y) a metric? (yes/no)",
			"yes",
			"The discrete metric satisfies all metric axioms including the triangle inequality.",
		},
		{
			"In a metric space, is d(x,y) = d(y,x)? (yes/no)",
			"yes",
			"Symmetry is one of the metric axioms: d(x,y) = d(y,x) for all x,y.",
		},
		{
			"Is d(x,y) = |x²-y²| a metric on R? (yes/no)",
			"no",
			"It is not a metric because d(1,-1)=0 but 1≠-1, violating the identity of indiscernibles.",
		},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.reason,
	}
}

// ----- 2. open_closed -----

type openClosedGen struct{}

func (g *openClosedGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		set      string
		isOpen   string
		isClosed string
		reason   string
	}
	table := []entry{
		{"(0,1)", "yes", "no", "(0,1) is open in R because every point has an ε-neighborhood contained in the interval. It is not closed because 0 and 1 are limit points not in the set."},
		{"[0,1]", "no", "yes", "[0,1] is closed in R because it contains all its limit points. It is not open because points 0 and 1 have no ε-neighborhood within the set."},
		{"(0,1]", "no", "no", "(0,1] is neither open nor closed: 1 has no open neighborhood within (0,1] is tricky—it's not open because 0 is a limit point not in the set, and it's not closed either because its complement includes limit point 1."},
		{"R", "yes", "yes", "R is both open and closed in itself. It is open by definition and closed because its complement (empty set) is open."},
		{"∅ (empty set)", "yes", "yes", "The empty set is both open and closed by definition (vacuously satisfies both conditions)."},
		{"{0}", "no", "yes", "A singleton {0} is closed in R (it contains its only limit point), but not open because 0 has no ε-neighborhood contained in {0}."},
		{"(0,∞)", "yes", "no", "(0,∞) is open in R because every positive x has an ε-neighborhood in (0,∞). It is not closed because 0 is a limit point not in the set."},
		{"[0,∞)", "no", "yes", "[0,∞) is closed in R because it contains all its limit points. It is not open because 0 has no ε-neighborhood contained in the set."},
		{"{1/n : n∈N} ∪ {0}", "no", "yes", "This set is closed because it contains its limit point 0. It is not open because 0 has no ε-neighborhood contained in the set."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Is the interval/set %s open in R? (yes/no)", e.set),
			Answer:      e.isOpen,
			Explanation: e.reason,
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is the interval/set %s closed in R? (yes/no)", e.set),
		Answer:      e.isClosed,
		Explanation: e.reason,
	}
}

// ----- 3. continuous -----

type continuousGen struct{}

func (g *continuousGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		f          string
		domain     string
		continuous string
		reason     string
	}
	table := []entry{
		{"f(x)=x²", "R", "yes", "f(x)=x² is a polynomial, and all polynomials are continuous on R."},
		{"f(x)=1/x", "R", "no", "f(x)=1/x is discontinuous at x=0 because the limit as x→0 does not exist (approaches ±∞)."},
		{"the floor function ⌊x⌋", "R", "no", "The floor function is discontinuous at every integer because the left and right limits differ."},
		{"f(x)=sin(x)", "R", "yes", "sin(x) is continuous on R — it is a trigonometric function with no breaks."},
		{"f(x)=eˣ", "R", "yes", "f(x)=eˣ is continuous on R as an exponential function."},
		{"f(x)=|x|", "R", "yes", "f(x)=|x| is continuous on R (although it is not differentiable at x=0, continuity is satisfied)."},
		{"f(x)=1/(x²−1)", "R", "no", "f(x)=1/(x²−1) is discontinuous at x=±1 where the denominator is zero."},
		{"f(x)=√x", "[0,∞)", "yes", "f(x)=√x is continuous on its domain [0,∞)."},
		{"f(x)=tan(x)", "R", "no", "tan(x) is discontinuous at x=π/2 + nπ where the function approaches ±∞."},
		{"f(x)=x³−3x+1", "R", "yes", "f(x)=x³−3x+1 is a polynomial, and all polynomials are continuous on R."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s continuous on %s? (yes/no)", e.f, e.domain),
		Answer:      e.continuous,
		Explanation: e.reason,
	}
}
