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
	reg.Register("topo.compact", &compactGen{})
	reg.Register("topo.connected", &connectedGen{})
	reg.Register("topo.hausdorff", &hausdorffGen{})
	reg.Register("topo.quotient", &quotientGen{})
	reg.Register("topo.product", &productGen{})
	reg.Register("topo.urysohn", &urysohnGen{})
	reg.Register("topo.metrization", &metrizationGen{})
	reg.Register("topo.connected.path", &pathConnectedGen{})
	reg.Register("topo.compact.local", &localCompactGen{})
	reg.Register("topo.separation.tietze", &tietzeGen{})
	reg.Register("topo.homotopy", &homotopyGen{})
	reg.Register("topo.manifold", &manifoldGen{})
	reg.Register("topo.connected.local_path", &localPathGen{})
	reg.Register("topo.compact.one_point", &onePointGen{})
	reg.Register("topo.separation.completely_regular", &completelyRegularGen{})
	reg.Register("topo.homotopy.fundamental", &fundamentalGen{})
	reg.Register("topo.manifold.orientability", &orientabilityGen{})
	reg.Register("topo.embedding", &embeddingGen{})
	reg.Register("topo.knot", &knotGen{})
	reg.Register("topo.cohomology", &cohomologyGen{})
	reg.Register("topo.compactification.stone_cech", &stoneCechGen{})
	reg.Register("topo.dimension", &dimensionGen{})
}

// ----- 1. metric -----

type metricGen struct{}

func (g *metricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		reason   string
	}
	table := []entry{
		{
			"Is the usual distance \\(d(x,y)=|x-y|\\) a metric on \\(\\mathbb{R}\\)? (yes/no)",
			"yes",
			"d(x,y)=|x-y| satisfies all metric axioms: non-negativity, identity of indiscernibles, symmetry, and triangle inequality.",
		},
		{
			"In a metric space, is it true that \\(d(x,y)=0\\) if and only if \\(x=y\\)? (yes/no)",
			"yes",
			"This is the identity of indiscernibles axiom: d(x,y)=0 ⇔ x=y.",
		},
		{
			"In a metric space, is \\(d(x,z) \\leq d(x,y)+d(y,z)\\)? (yes/no)",
			"yes",
			"This is the triangle inequality, one of the defining properties of a metric.",
		},
		{
			"Is \\(d(x,y) = (x-y)^{2}\\) a metric on \\(\\mathbb{R}\\)? (yes/no)",
			"no",
			"It violates the triangle inequality. For example, d(0,2)=4 but d(0,1)+d(1,2)=1+1=2 < 4.",
		},
		{
			"In a metric space, is \\(d(x,y)\\) always non-negative? (yes/no)",
			"yes",
			"Non-negativity is one of the metric axioms: d(x,y) ≥ 0 for all x,y.",
		},
		{
			"Is the discrete metric \\((d(x,y)=1\\) if \\(x\\neq y\\), \\(0\\) if \\(x=y)\\) a metric? (yes/no)",
			"yes",
			"The discrete metric satisfies all metric axioms including the triangle inequality.",
		},
		{
			"In a metric space, is \\(d(x,y) = d(y,x)\\)? (yes/no)",
			"yes",
			"Symmetry is one of the metric axioms: d(x,y) = d(y,x) for all x,y.",
		},
		{
			"Is \\(d(x,y) = |x^{2}-y^{2}|\\) a metric on \\(\\mathbb{R}\\)? (yes/no)",
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

func (g *openClosedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
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

func (g *continuousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		f          string
		domain     string
		continuous string
		reason     string
	}
	table := []entry{
		{"\\(f(x)=x^{2}\\)", "\\(\\mathbb{R}\\)", "yes", "\\(f(x)=x^{2}\\) is a polynomial, and all polynomials are continuous on \\(\\mathbb{R}\\)."},
		{"\\(f(x)=1/x\\)", "\\(\\mathbb{R}\\)", "no", "\\(f(x)=1/x\\) is discontinuous at \\(x=0\\) because the limit as \\(x\\to 0\\) does not exist (approaches \\(\\pm \\infty\\))."},
		{"the floor function \\(\\lfloor x \\rfloor\\)", "\\(\\mathbb{R}\\)", "no", "The floor function is discontinuous at every integer because the left and right limits differ."},
		{"\\(f(x)=\\sin(x)\\)", "\\(\\mathbb{R}\\)", "yes", "\\(\\sin(x)\\) is continuous on \\(\\mathbb{R}\\) — it is a trigonometric function with no breaks."},
		{"\\(f(x)=e^{x}\\)", "\\(\\mathbb{R}\\)", "yes", "\\(f(x)=e^{x}\\) is continuous on \\(\\mathbb{R}\\) as an exponential function."},
		{"\\(f(x)=|x|\\)", "\\(\\mathbb{R}\\)", "yes", "\\(f(x)=|x|\\) is continuous on \\(\\mathbb{R}\\) (although it is not differentiable at \\(x=0\\), continuity is satisfied)."},
		{"\\(f(x)=1/(x^{2}-1)\\)", "\\(\\mathbb{R}\\)", "no", "\\(f(x)=1/(x^{2}-1)\\) is discontinuous at \\(x=\\pm 1\\) where the denominator is zero."},
		{"\\(f(x)=\\sqrt{x}\\)", "\\([0,\\infty)\\)", "yes", "\\(f(x)=\\sqrt{x}\\) is continuous on its domain \\([0,\\infty)\\)."},
		{"\\(f(x)=\\tan(x)\\)", "\\(\\mathbb{R}\\)", "no", "\\(\\tan(x)\\) is discontinuous at \\(x=\\pi/2 + n\\pi\\) where the function approaches \\(\\pm \\infty\\)."},
		{"\\(f(x)=x^{3}-3x+1\\)", "\\(\\mathbb{R}\\)", "yes", "\\(f(x)=x^{3}-3x+1\\) is a polynomial, and all polynomials are continuous on \\(\\mathbb{R}\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\(%s\\) continuous on \\(%s\\)? (yes/no)", e.f, e.domain),
		Answer:      e.continuous,
		Explanation: e.reason,
	}
}

// ----- 4. compact -----

type compactGen struct{}

func (g *compactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set     string
		compact string
		reason  string
	}
	table := []entry{
		{"\\([0,1]\\)", "yes", "Heine-Borel: [0,1] is closed and bounded in R, hence compact."},
		{"\\((0,1)\\)", "no", "(0,1) is not closed (0 and 1 are limit points outside), so not compact."},
		{"\\(\\mathbb{R}\\)", "no", "R is not bounded, so not compact (Heine-Borel)."},
		{"\\(\\{0\\} \\cup \\{1/n : n\\in\\mathbb{N}\\}\\)", "yes", "This set is closed and bounded, and every open cover has a finite subcover."},
		{"\\((0,1]\\)", "no", "(0,1] is not closed (0 missing), so not compact."},
		{"\\([0,1] \\cup [2,3]\\)", "yes", "Finite union of compact sets is compact; both intervals are compact."},
		{"\\(\\{x : x\\geq 0\\}\\)", "no", "Unbounded above, so not compact."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is the set %s compact in \\(\\mathbb{R}\\)? (yes/no)", e.set),
		Answer:      e.compact,
		Explanation: e.reason,
	}
}

// ----- 5. connected -----

type connectedGen struct{}

func (g *connectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set       string
		connected string
		reason    string
	}
	table := []entry{
		{"\\((0,1)\\)", "yes", "(0,1) is an interval, hence connected in R."},
		{"\\((0,1) \\cup (2,3)\\)", "no", "Two disjoint open intervals can be separated by open sets, so not connected."},
		{"\\(\\mathbb{R}\\)", "yes", "R is connected (intervals are connected)."},
		{"\\(\\mathbb{Q}\\)", "no", "Q is totally disconnected: between any two rationals there is an irrational separating them."},
		{"\\([0,1]\\)", "yes", "Closed intervals are connected."},
		{"\\(\\{0,1\\}\\)", "no", "Two isolated points can be separated by disjoint opens, so not connected."},
		{"\\((0,1] \\)", "yes", "(0,1] is an interval, hence connected."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is the set %s connected in \\(\\mathbb{R}\\)? (yes/no)", e.set),
		Answer:      e.connected,
		Explanation: e.reason,
	}
}

// ----- 6. hausdorff -----

type hausdorffGen struct{}

func (g *hausdorffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		space     string
		hausdorff string
		reason    string
	}
	table := []entry{
		{"\\(\\mathbb{R}\\) with the usual metric", "yes", "Every metric space is Hausdorff: distinct points have disjoint epsilon-balls."},
		{"the discrete topology on any set", "yes", "Discrete spaces are metric (discrete metric), hence Hausdorff."},
		{"the indiscrete topology on \\(\\{a,b\\}\\)", "no", "In indiscrete topology only {} and the whole set are open, so distinct points cannot be separated."},
		{"the cofinite topology on \\(\\mathbb{R}\\)", "no", "Any two nonempty opens intersect (their complements are finite), so separation fails."},
		{"\\(\\mathbb{R}^2\\) with the Euclidean metric", "yes", "Metric spaces are Hausdorff."},
		{"the Sierpinski space", "no", "Points cannot be separated by disjoint opens in the Sierpinski topology."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s Hausdorff? (yes/no)", e.space),
		Answer:      e.hausdorff,
		Explanation: e.reason,
	}
}

// ----- 7. quotient -----

type quotientGen struct{}

func (g *quotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		reason   string
	}
	table := []entry{
		{"Is the quotient map \\(q: X \\to X/\\sim\\) always continuous by definition? (yes/no)", "yes", "The quotient topology is the finest topology making q continuous."},
		{"If \\(X=[0,1]\\) and \\(0\\sim1\\), is \\(X/\\sim\\) homeomorphic to \\(S^{1}\\)? (yes/no)", "yes", "Glueing the endpoints of an interval yields a circle."},
		{"Is the quotient of a compact space by any equivalence relation always Hausdorff? (yes/no)", "no", "The quotient of compact need not be Hausdorff; e.g. collapsing a non-closed set."},
		{"Does the quotient topology make the projection an open map? (yes/no)", "no", "Quotient maps need not be open; they are continuous and closed under compact-Hausdorff hypotheses."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 8. product -----

type productGen struct{}

func (g *productGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		reason   string
	}
	table := []entry{
		{"Is the product of two compact spaces compact (Tychonoff for finite products)? (yes/no)", "yes", "Finite products of compact spaces are compact."},
		{"Is \\(\\mathbb{R}^{\\mathbb{N}}\\) with product topology metrizable? (yes/no)", "yes", "Countable product of metrizable spaces is metrizable."},
		{"Are projections \\(\\pi_i: \\prod X_i \\to X_i\\) continuous in the product topology? (yes/no)", "yes", "Product topology is coarsest making all projections continuous."},
		{"Is the box topology finer than the product topology on infinite products? (yes/no)", "yes", "Box topology has more opens than product topology when infinitely many factors."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 9. urysohn -----

type urysohnGen struct{}

func (g *urysohnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		reason   string
	}
	table := []entry{
		{"Does Urysohn's lemma guarantee a continuous \\(f: X \\to [0,1]\\) separating disjoint closed sets in a normal space? (yes/no)", "yes", "In normal spaces, disjoint closed sets can be separated by a continuous function."},
		{"Is every metric space normal (hence Urysohn's lemma applies)? (yes/no)", "yes", "Metric spaces are normal, so the lemma holds."},
		{"Does Urysohn's lemma require the space to be Hausdorff? (yes/no)", "no", "It requires normality, which implies Hausdorff only in some formulations; normality is the key."},
		{"Can Urysohn's function be chosen with \\(f(A)=0\\) and \\(f(B)=1\\) for disjoint closed \\(A,B\\)? (yes/no)", "yes", "The lemma constructs f with exactly those values."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 10. metrization -----

type metrizationGen struct{}

func (g *metrizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		reason   string
	}
	table := []entry{
		{"Does Urysohn metrization state that a second-countable regular Hausdorff space is metrizable? (yes/no)", "yes", "Regular + second-countable + Hausdorff implies metrizable."},
		{"Is every compact Hausdorff second-countable space metrizable? (yes/no)", "yes", "Compact Hausdorff + second-countable → metrizable by Urysohn."},
		{"Does metrizability imply the space is Hausdorff? (yes/no)", "yes", "Every metric space is Hausdorff, so metrizable spaces are Hausdorff."},
		{"Is the Sorgenfrey line metrizable? (yes/no)", "no", "The Sorgenfrey line is regular and Hausdorff but not second-countable in the metrizable sense; it is not metrizable."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 11. path-connected -----

type pathConnectedGen struct{}

func (g *pathConnectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is [0,1] path-connected? (yes/no)", "yes", "[0,1] is an interval, hence path-connected via linear paths."},
		{"Is (0,1) ∪ (2,3) path-connected? (yes/no)", "no", "No path can jump between the two disjoint intervals."},
		{"Does path-connected imply connected? (yes/no)", "yes", "Path-connected spaces are always connected."},
		{"Is Q path-connected? (yes/no)", "no", "Q is totally disconnected, hence not path-connected."},
	}
	hard := []entry{
		{"Is the topologist's sine curve connected but not path-connected? (yes/no)", "yes", "Classic example: connected but no path between the curve and the limit segment."},
		{"Does path-connectedness imply local path-connectedness? (yes/no)", "no", "Comb space is path-connected but not locally path-connected."},
		{"Is R^2 minus a point path-connected? (yes/no)", "yes", "Punctured plane remains path-connected (go around the hole)."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 12. local compact -----

type localCompactGen struct{}

func (g *localCompactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is R locally compact? (yes/no)", "yes", "Every point has a compact neighbourhood [x-ε,x+ε]."},
		{"Is [0,1] locally compact? (yes/no)", "yes", "Compact spaces are locally compact."},
		{"Is Q locally compact? (yes/no)", "no", "Q has no compact neighbourhood of any point in the subspace topology."},
		{"Does locally compact Hausdorff imply Tychonoff? (yes/no)", "yes", "Locally compact Hausdorff spaces are completely regular."},
	}
	hard := []entry{
		{"Is R^n locally compact? (yes/no)", "yes", "Euclidean spaces are locally compact."},
		{"Is an infinite-dimensional Hilbert space locally compact? (yes/no)", "no", "Unit ball is not compact in infinite dimensions, so not locally compact."},
		{"Does one-point compactification apply to locally compact Hausdorff spaces? (yes/no)", "yes", "One-point compactification gives a compact Hausdorff space."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 13. Tietze -----

type tietzeGen struct{}

func (g *tietzeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Does Tietze extension extend continuous f: A→[0,1] from closed A⊂X normal to all X? (yes/no)", "yes", "Tietze: closed subset of normal space, bounded continuous function extends."},
		{"Is every normal space one where Tietze holds for [0,1]-valued functions? (yes/no)", "yes", "Characterization of normality via Tietze for [0,1]."},
		{"Does Tietze imply Urysohn's lemma? (yes/no)", "yes", "Urysohn is special case of Tietze with two closed sets."},
	}
	hard := []entry{
		{"Can any continuous f: A→R from closed A in normal X be extended to X? (yes/no)", "yes", "Full Tietze for real-valued (possibly unbounded) functions on normal spaces."},
		{"Does Tietze fail if A is not closed? (yes/no)", "yes", "Closedness is essential; e.g. Q ⊂ R not closed cannot extend all functions."},
		{"Is Tietze equivalent to normality? (yes/no)", "yes", "Normal iff every bounded continuous function from a closed subset extends."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 14. homotopy -----

type homotopyGen struct{}

func (g *homotopyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is fundamental group π_1(S^1) ≅ Z? (yes/no)", "yes", "Loops around circle classified by winding number."},
		{"Is simply connected equivalent to π_1 trivial? (yes/no)", "yes", "Every loop contracts to a point iff fundamental group is trivial."},
		{"Are homotopic maps into R^n always homotopic? (yes/no)", "yes", "R^n is contractible, so any two maps are homotopic."},
		{"Does homotopy define an equivalence relation? (yes/no)", "yes", "Reflexive, symmetric, transitive via glueing homotopies."},
	}
	hard := []entry{
		{"Is π_1(S^2) trivial? (yes/no)", "yes", "2-sphere is simply connected."},
		{"Is π_1(T^2) ≅ Z×Z? (yes/no)", "yes", "Torus fundamental group is product of two circle factors."},
		{"Does covering space p: R→S^1 give universal cover of circle? (yes/no)", "yes", "R is simply connected cover of S^1."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 15. manifold -----

type manifoldGen struct{}

func (g *manifoldGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is S^1 a 1-manifold? (yes/no)", "yes", "Circle locally looks like an interval (R)."},
		{"Is S^2 a 2-manifold? (yes/no)", "yes", "Sphere locally looks like R^2 via stereographic charts."},
		{"Is R^n a manifold of dimension n? (yes/no)", "yes", "Euclidean space is the model manifold."},
		{"Is a manifold required to be Hausdorff and second-countable? (yes/no)", "yes", "Standard definition includes Hausdorff and second-countable plus locally Euclidean."},
	}
	hard := []entry{
		{"Is torus T^2 = S^1×S^1 a compact 2-manifold? (yes/no)", "yes", "Product of compact manifolds is compact manifold."},
		{"Is every manifold metrizable? (yes/no)", "yes", "Urysohn plus manifold hypotheses imply metrizable."},
		{"Is R with discrete topology a 0-manifold? (yes/no)", "yes", "Discrete spaces are 0-dimensional manifolds."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 16. locally path-connected -----

type localPathGen struct{}

func (g *localPathGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is R locally path-connected? (yes/no)", "yes", "Intervals are path-connected neighbourhoods."},
		{"Is Q locally path-connected? (yes/no)", "no", "No neighbourhood in Q contains a path-connected open set."},
		{"Does locally path-connected + connected imply path-connected? (yes/no)", "yes", "Classic theorem for locally path-connected spaces."},
	}
	hard := []entry{
		{"Is comb space locally path-connected at the base? (yes/no)", "no", "The baseline has no path-connected neighbourhood."},
		{"Is topologist's sine curve locally connected? (yes/no)", "no", "No small connected neighbourhood of points on the limit segment."},
		{"Does manifold ⇒ locally path-connected? (yes/no)", "yes", "Manifolds are locally Euclidean, hence locally path-connected."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 17. one-point compactification -----

type onePointGen struct{}

func (g *onePointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is one-point compactification of R homeomorphic to S^1? (yes/no)", "yes", "R ∪ {∞} ≅ S^1."},
		{"Is one-point compactification of R^2 homeomorphic to S^2? (yes/no)", "yes", "R^2 ∪ {∞} ≅ S^2 via stereographic."},
		{"Is locally compact Hausdorff required for one-point compactification? (yes/no)", "yes", "Alexandroff compactification needs locally compact Hausdorff."},
	}
	hard := []entry{
		{"Is one-point compactification of Q compact? (yes/no)", "no", "Q is not locally compact, so Alexandroff does not give Hausdorff compact."},
		{"Does one-point compactification add exactly one point? (yes/no)", "yes", "By definition."},
		{"Is compact Hausdorff space its own one-point compactification? (yes/no)", "no", "Already compact, no point added."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 18. completely regular -----

type completelyRegularGen struct{}

func (g *completelyRegularGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Does Tietze imply completely regular for normal spaces? (yes/no)", "yes", "Normal ⇒ completely regular (Tychonoff)."},
		{"Is every metric space completely regular? (yes/no)", "yes", "Metric ⇒ normal ⇒ completely regular."},
		{"Does completely regular mean points and closed sets separated by functions? (yes/no)", "yes", "For each closed A and x∉A, ∃f:X→[0,1] with f(x)=1, f|_A=0."},
	}
	hard := []entry{
		{"Is Tychonoff = completely regular + T1? (yes/no)", "yes", "Tychonoff spaces are T1 and completely regular."},
		{"Does product of completely regular spaces remain completely regular? (yes/no)", "yes", "Product preserves complete regularity."},
		{"Is Niemytzki plane completely regular? (yes/no)", "yes", "It is Tychonoff."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 19. fundamental group calculations -----

type fundamentalGen struct{}

func (g *fundamentalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is π_1(S^1) ≅ Z? (yes/no)", "yes", "Winding number."},
		{"Is π_1(S^2) trivial? (yes/no)", "yes", "S^2 simply connected."},
		{"Is π_1(T^2) ≅ Z×Z? (yes/no)", "yes", "Torus product of circles."},
	}
	hard := []entry{
		{"Is π_1(R^2 minus n points) free group on n generators? (yes/no)", "yes", "Punctured plane."},
		{"Is π_1(Klein bottle) non-abelian? (yes/no)", "yes", "Presentation ⟨a,b|aba^{-1}=b^{-1}⟩."},
		{"Does Van Kampen compute π_1 of wedge of circles? (yes/no)", "yes", "Free group."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 20. orientability -----

type orientabilityGen struct{}

func (g *orientabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is S^2 orientable? (yes/no)", "yes", "Sphere is orientable."},
		{"Is Möbius band orientable? (yes/no)", "no", "Möbius has one side, non-orientable."},
		{"Is torus T^2 orientable? (yes/no)", "yes", "Torus is orientable."},
	}
	hard := []entry{
		{"Is Klein bottle non-orientable? (yes/no)", "yes", "Contains Möbius band."},
		{"Does orientability mean consistent choice of local orientation? (yes/no)", "yes", "Atlas with positive Jacobian transitions."},
		{"Is R^n orientable? (yes/no)", "yes", "Euclidean space orientable."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

type embeddingGen struct{}
func (g *embeddingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{q,a,e string}
	easy:=[]entry{{"Does Whitney embedding embed n-manifold in R^{2n}? (yes/no)","yes","Whitney."},{"Is embedding injective immersion homeomorphic onto image? (yes/no)","yes","Embedding."},{"Does Urysohn embedding use countable product of [0,1]? (yes/no)","yes","Metrization proof."}}
	hard:=[]entry{{"Does Nash embed Riemannian manifold isometrically in R^n? (yes/no)","yes","Nash embedding."},{"Is Alexander horned sphere embedding wild? (yes/no)","yes","Wild embedding."},{"Does Schoenflies generalize Jordan curve? (yes/no)","yes","Sphere."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.q,Answer:e.a,Explanation:e.e}
}
type knotGen struct{}
func (g *knotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{q,a,e string}
	easy:=[]entry{{"Is unknot trivial? (yes/no)","yes","Unknotted circle."},{"Does trefoil have crossing number 3? (yes/no)","yes","Simplest nontrivial."},{"Is knot complement fundamental group knot invariant? (yes/no)","yes","Knot group."}}
	hard:=[]entry{{"Is Jones polynomial knot invariant? (yes/no)","yes","Jones."},{"Does Gordon-Luecke show knot determined by complement? (yes/no)","yes","Knot complement."},{"Is figure-eight knot 4_1? (yes/no)","yes","Four crossings."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.q,Answer:e.a,Explanation:e.e}
}
type cohomologyGen struct{}
func (g *cohomologyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{q,a,e string}
	easy:=[]entry{{"Is cohomology dual to homology? (yes/no)","yes","Universal coefficients."},{"Does de Rham cohomology use differential forms? (yes/no)","yes","Smooth manifolds."},{"Is H^0 connected components? (yes/no)","yes","H^0."}}
	hard:=[]entry{{"Does cohomology ring have cup product? (yes/no)","yes","Ring structure."},{"Is singular cohomology homotopy invariant? (yes/no)","yes","Homotopy."},{"Does Poincaré duality hold for compact orientable manifold? (yes/no)","yes","Duality."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.q,Answer:e.a,Explanation:e.e}
}
type stoneCechGen struct{}
func (g *stoneCechGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{q,a,e string}
	easy:=[]entry{{"Is Stone-Čech compactification largest compactification? (yes/no)","yes","Universal."},{"Does βN contain remainder N*? (yes/no)","yes","Growth."},{"Is Tychonoff exactly embeddable in compact Hausdorff? (yes/no)","yes","Characterization."}}
	hard:=[]entry{{"Is βR not metrizable? (yes/no)","yes","Huge."},{"Does Stone-Čech of discrete is ultrafilters? (yes/no)","yes","βD."},{"Is completely regular iff subspace of compact Hausdorff? (yes/no)","yes","Tychonoff."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.q,Answer:e.a,Explanation:e.e}
}
type dimensionGen struct{}
func (g *dimensionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale:=int(1+ctx.Difficulty*4)
	type entry struct{q,a,e string}
	easy:=[]entry{{"Is covering dimension of R^n equal n? (yes/no)","yes","Lebesgue dimension."},{"Is dimension of Cantor set 0? (yes/no)","yes","Totally disconnected."},{"Does dimension not increase under subspace? (yes/no)","yes","Subspace."}}
	hard:=[]entry{{"Is inductive dimension equal covering for separable metric? (yes/no)","yes","Coincidence."},{"Does Brouwer show dimension is topological invariant? (yes/no)","yes","Invariance of domain."},{"Is Menger-Nöbeling embed n-dim compact in R^{2n+1}? (yes/no)","yes","Embedding."}}
	pool:=easy
	if scale>3{pool=append(easy,hard...)}
	e:=pool[rand.Intn(len(pool))]
	return generator.Problem{Question:e.q,Answer:e.a,Explanation:e.e}
}
