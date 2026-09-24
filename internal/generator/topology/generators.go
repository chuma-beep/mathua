package topology

import (
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
	reg.Register("topo.covering", &coveringGen{})
	reg.Register("topo.fundamental.van_kampen", &vanKampenGen{})
	reg.Register("topo.fiber_bundle", &fiberBundleGen{})
	reg.Register("topo.homology.cellular", &cellularHomologyGen{})
	reg.Register("topo.cohomology.de_rham", &deRhamGen{})
}

// ----- 1. metric -----

type metricGen struct{}

func (g *metricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"What is d(x,x) in any metric space?", "0", "Identity of indiscernibles: d(x,y)=0 iff x=y, so d(x,x)=0."},
		{"In a metric space, what does d(x,y) = 0 force about x and y? (type like x=y)", "x=y", "Identity of indiscernibles: distance zero holds exactly for equal points."},
		{"Which axiom fails for d(x,y) = (x-y)^2 on R: triangle or symmetry? (type one word)", "triangle", "Squares break the triangle inequality: d(0,2)=4 exceeds d(0,1)+d(1,2)=2."},
		{"Is the usual distance \\(d(x,y)=|x-y|\\) a metric on \\(\\mathbb{R}\\)? (yes/no)", "yes", "d(x,y)=|x-y| satisfies all metric axioms: non-negativity, identity of indiscernibles, symmetry, and triangle inequality."},
		{"Is \\(d(x,y) = |x^{2}-y^{2}|\\) a metric on \\(\\mathbb{R}\\)? (yes/no)", "no", "It is not a metric because d(1,-1)=0 but 1≠-1, violating the identity of indiscernibles."},
		{"Is the discrete metric \\((d(x,y)=1\\) if \\(x\\neq y\\), \\(0\\) if \\(x=y)\\) a metric? (yes/no)", "yes", "The discrete metric satisfies all metric axioms including the triangle inequality."},
	}
	hard := []entry{
		{"In a metric space, is it true that \\(d(x,y)=0\\) if and only if \\(x=y\\)? (yes/no)", "yes", "This is the identity of indiscernibles axiom: d(x,y)=0 ⇔ x=y."},
		{"In a metric space, is \\(d(x,z) \\leq d(x,y)+d(y,z)\\)? (yes/no)", "yes", "This is the triangle inequality, one of the defining properties of a metric."},
		{"In a metric space, is \\(d(x,y)\\) always non-negative? (yes/no)", "yes", "Non-negativity is one of the metric axioms: d(x,y) ≥ 0 for all x,y."},
		{"In a metric space, is \\(d(x,y) = d(y,x)\\)? (yes/no)", "yes", "Symmetry is one of the metric axioms: d(x,y) = d(y,x) for all x,y."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 2. open_closed -----

type openClosedGen struct{}

func (g *openClosedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Is (0,1) open, closed, both, or neither? (type one word)", "open", "Every point of (0,1) has an ε-neighborhood inside it, but 0 and 1 are limit points outside, so open and not closed."},
		{"Is [0,1] open, closed, both, or neither? (type one word)", "closed", "[0,1] contains all its limit points, so closed; endpoints have no ε-neighborhood inside, so not open."},
		{"Is R open, closed, both, or neither? (type one word)", "both", "R is open by definition and closed because its complement (empty set) is open."},
		{"Is (0,1] open, closed, both, or neither? (type one word)", "neither", "(0,1] misses limit point 0 so not closed, and 1 has no ε-neighborhood inside so not open."},
		{"Is the interval/set (0,1) open in R? (yes/no)", "yes", "(0,1) is open in R because every point has an ε-neighborhood contained in the interval."},
		{"Is the interval/set (0,1) closed in R? (yes/no)", "no", "(0,1) is not closed because 0 and 1 are limit points not in the set."},
	}
	hard := []entry{
		{"Is the interval/set {0} closed in R? (yes/no)", "yes", "A singleton {0} is closed in R (it contains its only limit point)."},
		{"Is the interval/set {1/n : n∈N} ∪ {0} closed in R? (yes/no)", "yes", "This set is closed because it contains its limit point 0."},
		{"Is the interval/set (0,∞) open in R? (yes/no)", "yes", "(0,∞) is open in R because every positive x has an ε-neighborhood in (0,∞)."},
		{"Is the interval/set [0,∞) closed in R? (yes/no)", "yes", "[0,∞) is closed in R because it contains all its limit points."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 3. continuous -----

type continuousGen struct{}

func (g *continuousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"At which point is f(x)=1/x discontinuous? (type like x=0)", "x=0", "The limit as x→0 does not exist (approaches ±∞), so 1/x breaks at x=0."},
		{"Which function is discontinuous at x=0: x^2 or 1/x? (type it)", "1/x", "x^2 is a polynomial hence continuous everywhere; 1/x blows up at x=0."},
		{"Is \\(f(x)=x^{2}\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "yes", "f(x)=x^2 is a polynomial, and all polynomials are continuous on R."},
		{"Is the floor function \\(\\lfloor x \\rfloor\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "no", "The floor function is discontinuous at every integer because the left and right limits differ."},
		{"Is \\(f(x)=|x|\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "yes", "f(x)=|x| is continuous on R (although not differentiable at x=0, continuity holds)."},
		{"Is \\(f(x)=\\tan(x)\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "no", "tan(x) is discontinuous at x=π/2 + nπ where the function approaches ±∞."},
	}
	hard := []entry{
		{"Is \\(f(x)=\\sin(x)\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "yes", "sin(x) is continuous on R — a trigonometric function with no breaks."},
		{"Is \\(f(x)=e^{x}\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "yes", "f(x)=e^x is continuous on R as an exponential function."},
		{"Is \\(f(x)=\\sqrt{x}\\) continuous on \\([0,\\infty)\\)? (yes/no)", "yes", "f(x)=√x is continuous on its domain [0,∞)."},
		{"Is \\(f(x)=1/(x^{2}-1)\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "no", "f(x)=1/(x^2-1) is discontinuous at x=±1 where the denominator is zero."},
		{"Is \\(f(x)=x^{3}-3x+1\\) continuous on \\(\\mathbb{R}\\)? (yes/no)", "yes", "f(x)=x^3-3x+1 is a polynomial, and all polynomials are continuous on R."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 4. compact -----

type compactGen struct{}

func (g *compactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Which Heine-Borel condition fails for (0,1): closed or bounded? (type one word)", "closed", "(0,1) is bounded but misses limit points 0 and 1, so it fails closedness."},
		{"Which Heine-Borel condition fails for R: closed or bounded? (type one word)", "bounded", "R is closed but unbounded, so it fails boundedness and is not compact."},
		{"Is the set [0,1] compact in \\(\\mathbb{R}\\)? (yes/no)", "yes", "Heine-Borel: [0,1] is closed and bounded in R, hence compact."},
		{"Is the set (0,1) compact in \\(\\mathbb{R}\\)? (yes/no)", "no", "(0,1) is not closed (0 and 1 are limit points outside), so not compact."},
		{"Is the set R compact in \\(\\mathbb{R}\\)? (yes/no)", "no", "R is not bounded, so not compact (Heine-Borel)."},
	}
	hard := []entry{
		{"Is the set \\(\\{0\\} \\cup \\{1/n : n\\in\\mathbb{N}\\}\\) compact in \\(\\mathbb{R}\\)? (yes/no)", "yes", "This set is closed and bounded, and every open cover has a finite subcover."},
		{"Is the set \\((0,1]\\) compact in \\(\\mathbb{R}\\)? (yes/no)", "no", "(0,1] is not closed (0 missing), so not compact."},
		{"Is the set \\([0,1] \\cup [2,3]\\) compact in \\(\\mathbb{R}\\)? (yes/no)", "yes", "Finite union of compact sets is compact; both intervals are compact."},
		{"Is the set \\(\\{x : x\\geq 0\\}\\) compact in \\(\\mathbb{R}\\)? (yes/no)", "no", "Unbounded above, so not compact."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 5. connected -----

type connectedGen struct{}

func (g *connectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"How many connected components does (0,1) union (2,3) have?", "2", "The two disjoint open intervals separate, giving one component each."},
		{"How many connected components does R have?", "1", "R is connected (intervals are connected), so a single component."},
		{"Is the set (0,1) connected in \\(\\mathbb{R}\\)? (yes/no)", "yes", "(0,1) is an interval, hence connected in R."},
		{"Is the set Q connected in \\(\\mathbb{R}\\)? (yes/no)", "no", "Q is totally disconnected: between any two rationals an irrational separates them."},
		{"Is the set \\(\\{0,1\\}\\) connected in \\(\\mathbb{R}\\)? (yes/no)", "no", "Two isolated points can be separated by disjoint opens, so not connected."},
	}
	hard := []entry{
		{"Is the set \\(\\mathbb{R}\\) connected in \\(\\mathbb{R}\\)? (yes/no)", "yes", "R is connected (intervals are connected)."},
		{"Is the set \\([0,1]\\) connected in \\(\\mathbb{R}\\)? (yes/no)", "yes", "Closed intervals are connected."},
		{"Is the set \\((0,1]\\) connected in \\(\\mathbb{R}\\)? (yes/no)", "yes", "(0,1] is an interval, hence connected."},
		{"Is the set \\((0,1) \\cup (2,3)\\) connected in \\(\\mathbb{R}\\)? (yes/no)", "no", "Two disjoint open intervals can be separated by open sets, so not connected."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 6. hausdorff -----

type hausdorffGen struct{}

func (g *hausdorffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Which topology on {a,b} is not Hausdorff: discrete or indiscrete? (type one word)", "indiscrete", "In the indiscrete topology only {} and the whole set are open, so distinct points cannot be separated."},
		{"Which topology on R is not Hausdorff: usual, discrete, or cofinite? (type one word)", "cofinite", "Any two nonempty cofinite opens intersect (complements are finite), so separation fails."},
		{"Is R with the usual metric Hausdorff? (yes/no)", "yes", "Every metric space is Hausdorff: distinct points have disjoint epsilon-balls."},
		{"Is the cofinite topology on R Hausdorff? (yes/no)", "no", "Any two nonempty opens intersect (their complements are finite), so separation fails."},
	}
	hard := []entry{
		{"Is the discrete topology on any set Hausdorff? (yes/no)", "yes", "Discrete spaces are metric (discrete metric), hence Hausdorff."},
		{"Is R^2 with the Euclidean metric Hausdorff? (yes/no)", "yes", "Metric spaces are Hausdorff."},
		{"Is the Sierpinski space Hausdorff? (yes/no)", "no", "Points cannot be separated by disjoint opens in the Sierpinski topology."},
		{"Is the indiscrete topology on {a,b} Hausdorff? (yes/no)", "no", "In indiscrete topology only {} and the whole set are open, so distinct points cannot be separated."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 7. quotient -----

type quotientGen struct{}

func (g *quotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"What space is [0,1] with 0~1 homeomorphic to? (type like S^1)", "S^1", "Glueing the endpoints of an interval closes the loop into a circle."},
		{"Is the quotient map \\(q: X \\to X/\\sim\\) always continuous by definition? (yes/no)", "yes", "The quotient topology is the finest topology making q continuous."},
		{"Is the quotient of a compact space by any equivalence relation always Hausdorff? (yes/no)", "no", "The quotient of compact need not be Hausdorff; e.g. collapsing a non-closed set."},
		{"Does the quotient topology make the projection an open map? (yes/no)", "no", "Quotient maps need not be open; they are continuous and closed under compact-Hausdorff hypotheses."},
	}
	hard := []entry{
		{"What space is the disc D^2 with its boundary collapsed to a point homeomorphic to? (type like S^2)", "S^2", "Collapsing the boundary circle of a 2-disc closes it into a 2-sphere."},
		{"Is the quotient of a compact Hausdorff space by a closed equivalence relation Hausdorff? (yes/no)", "yes", "Closed saturated sets keep the quotient Hausdorff in the compact Hausdorff setting."},
		{"If \\(X=[0,1]\\) and \\(0\\sim1\\), is \\(X/\\sim\\) homeomorphic to \\(S^{1}\\)? (yes/no)", "yes", "Glueing the endpoints of an interval yields a circle."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 8. product -----

type productGen struct{}

func (g *productGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Which topology is finer on infinite products: box or product? (type one word)", "box", "The box topology constrains every factor, giving strictly more opens than the product topology."},
		{"Are projections \\(\\pi_i\\) from a product continuous in the product topology? (yes/no)", "yes", "The product topology is the coarsest making all projections continuous."},
		{"Is the product of two compact spaces compact? (yes/no)", "yes", "Finite products of compact spaces are compact (finite Tychonoff)."},
		{"Does a basic product open constrain infinitely many coordinates? (yes/no)", "no", "A basic open constrains only finitely many factors; the rest are whole spaces."},
	}
	hard := []entry{
		{"Is an arbitrary product of compact spaces compact (Tychonoff)? (yes/no)", "yes", "Tychonoff: arbitrary products of compact spaces are compact."},
		{"Is the box topology the same as the product topology on infinite products? (yes/no)", "no", "Box has more opens; they agree only for finite products."},
		{"Is \\(\\mathbb{R}^{\\mathbb{N}}\\) with product topology metrizable? (yes/no)", "yes", "Countable product of metrizable spaces is metrizable."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 9. urysohn -----

type urysohnGen struct{}

func (g *urysohnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"What interval does the Urysohn function map into? (type like [0,1])", "[0,1]", "Urysohn builds f: X→[0,1] with f=0 on A and f=1 on B."},
		{"What value does the Urysohn function take on the closed set A?", "0", "The lemma constructs f with f(A)=0 and f(B)=1 for disjoint closed A,B."},
		{"Does Urysohn's lemma separate disjoint closed sets in a normal space? (yes/no)", "yes", "In normal spaces, disjoint closed sets can be separated by a continuous function."},
		{"Does Urysohn's lemma apply without normality? (yes/no)", "no", "Normality is the key hypothesis; without it the separating function may not exist."},
	}
	hard := []entry{
		{"Is every metric space normal (hence Urysohn's lemma applies)? (yes/no)", "yes", "Metric spaces are normal, so the lemma holds."},
		{"Does Urysohn's lemma require the space to be Hausdorff? (yes/no)", "no", "It requires normality, which implies Hausdorff only in some formulations; normality is the key."},
		{"Can Urysohn's function be chosen with \\(f(A)=0\\) and \\(f(B)=1\\) for disjoint closed \\(A,B\\)? (yes/no)", "yes", "The lemma constructs f with exactly those values."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 10. metrization -----

type metrizationGen struct{}

func (g *metrizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Which line is not metrizable: real or Sorgenfrey? (type one word)", "Sorgenfrey", "The Sorgenfrey line is regular and Hausdorff but fails the countability hypotheses for metrizability."},
		{"How many hypotheses does Urysohn metrization combine: 2 or 3? (type a number)", "3", "Regular + Hausdorff + second-countable: three hypotheses imply metrizable."},
		{"Is every compact Hausdorff second-countable space metrizable? (yes/no)", "yes", "Compact Hausdorff + second-countable → metrizable by Urysohn."},
		{"Is the Sorgenfrey line metrizable? (yes/no)", "no", "The Sorgenfrey line is regular and Hausdorff but not second-countable in the metrizable sense; it is not metrizable."},
	}
	hard := []entry{
		{"Does Urysohn metrization state that a second-countable regular Hausdorff space is metrizable? (yes/no)", "yes", "Regular + second-countable + Hausdorff implies metrizable."},
		{"Does metrizability imply the space is Hausdorff? (yes/no)", "yes", "Every metric space is Hausdorff, so metrizable spaces are Hausdorff."},
		{"Are metrizable spaces normal? (yes/no)", "yes", "Metric spaces are normal, so metrizable spaces inherit normality."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
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
		{"How many path components does (0,1) union (2,3) have?", "2", "No path can jump between the disjoint intervals, so each is its own component."},
		{"Does path-connected imply connected? (yes/no)", "yes", "Path-connected spaces are always connected."},
		{"Is [0,1] path-connected? (yes/no)", "yes", "[0,1] is an interval, hence path-connected via linear paths."},
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
		{"Which space is not locally compact: R or Q? (type one word)", "Q", "Q has no compact neighbourhood of any point in the subspace topology."},
		{"Is R locally compact? (yes/no)", "yes", "Every point has a compact neighbourhood [x-ε,x+ε]."},
		{"Is Q locally compact? (yes/no)", "no", "Q has no compact neighbourhood of any point in the subspace topology."},
		{"Is [0,1] locally compact? (yes/no)", "yes", "Compact spaces are locally compact."},
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
		{"From which subsets does Tietze extend functions: open or closed? (type one word)", "closed", "Tietze extends continuous functions from closed subsets of normal spaces."},
		{"Does Tietze extension extend continuous f: A→[0,1] from closed A⊂X normal to all X? (yes/no)", "yes", "Tietze: closed subset of normal space, bounded continuous function extends."},
		{"Does Tietze fail if A is not closed? (yes/no)", "yes", "Closedness is essential; e.g. Q ⊂ R not closed cannot extend all functions."},
		{"Does Tietze extend from non-closed sets in general? (yes/no)", "no", "Closedness is essential; extension from arbitrary subsets fails in general."},
	}
	hard := []entry{
		{"Can any continuous f: A→R from closed A in normal X be extended to X? (yes/no)", "yes", "Full Tietze for real-valued (possibly unbounded) functions on normal spaces."},
		{"Does Tietze imply Urysohn's lemma? (yes/no)", "yes", "Urysohn is special case of Tietze with two closed sets."},
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
		{"What is pi_1(S^1)? (type like Z)", "Z", "Loops around the circle are classified by winding number, an integer."},
		{"What is pi_1(S^2)? (type 0 for trivial)", "0", "The 2-sphere is simply connected: every loop contracts."},
		{"Is simply connected equivalent to pi_1 trivial? (yes/no)", "yes", "Every loop contracts to a point iff fundamental group is trivial."},
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
		{"What dimension is S^2 as a manifold?", "2", "Each point of the sphere has an R^2 neighbourhood via stereographic charts."},
		{"What dimension is S^1 as a manifold?", "1", "Each point of the circle has an interval (R) neighbourhood."},
		{"Is S^1 a 1-manifold? (yes/no)", "yes", "Circle locally looks like an interval (R)."},
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
		{"Which space is path-connected but not locally path-connected: interval or comb? (type one word)", "comb", "The comb space is path-connected but has no path-connected neighbourhood at the baseline."},
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
		{"What is the one-point compactification of R? (type like S^1)", "S^1", "R ∪ {∞} ≅ S^1 via stereographic projection."},
		{"What is the one-point compactification of R^2? (type like S^2)", "S^2", "R^2 ∪ {∞} ≅ S^2 via stereographic projection."},
		{"Is locally compact Hausdorff required for one-point compactification? (yes/no)", "yes", "Alexandroff compactification needs locally compact Hausdorff."},
		{"Is one-point compactification of R homeomorphic to S^1? (yes/no)", "yes", "R ∪ {∞} ≅ S^1."},
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
		{"What interval do the separating functions map into? (type like [0,1])", "[0,1]", "Complete regularity separates x from closed A by f: X→[0,1] with f(x)=1, f|_A=0."},
		{"Is every metric space completely regular? (yes/no)", "yes", "Metric ⇒ normal ⇒ completely regular."},
		{"Is every Hausdorff space completely regular? (yes/no)", "no", "Hausdorff alone does not give separating functions; complete regularity is strictly stronger."},
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
		{"What is pi_1(T^2)? (type like ZxZ)", "ZxZ", "The torus is a product of two circles, so loops carry two winding numbers."},
		{"Is pi_1(S^1) = Z? (yes/no)", "yes", "Winding number classifies loops on the circle."},
		{"Is pi_1(S^2) trivial? (yes/no)", "yes", "S^2 simply connected: every loop contracts."},
		{"Is pi_1(T^2) = ZxZ? (yes/no)", "yes", "Torus product of circles gives product of winding numbers."},
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
		{"Which band is non-orientable: cylinder or Mobius? (type one word)", "Mobius", "The Mobius band has one side: a loop around it flips orientation."},
		{"Is S^2 orientable? (yes/no)", "yes", "Sphere is orientable (outward normal field)."},
		{"Is Mobius band orientable? (yes/no)", "no", "Möbius has one side, non-orientable."},
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
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What is the Whitney target for a closed n-manifold? (type like R^2n)", "R^2n", "Whitney embeds closed n-manifolds in R^2n."},
		{"Is an embedding an injective immersion homeomorphic onto its image? (yes/no)", "yes", "Embedding definition."},
		{"Is the Alexander horned sphere tame? (yes/no)", "no", "It is the classic wild embedding of S^2 in R^3."},
		{"Does Whitney embedding embed n-manifolds in R^2n? (yes/no)", "yes", "Whitney theorem."},
	}
	hard := []entry{{"Does Nash embed Riemannian manifold isometrically in R^n? (yes/no)", "yes", "Nash embedding."}, {"Is Alexander horned sphere embedding wild? (yes/no)", "yes", "Wild embedding."}, {"Does Schoenflies generalize Jordan curve? (yes/no)", "yes", "Sphere."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type knotGen struct{}

func (g *knotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What is the crossing number of the trefoil?", "3", "Three crossings: the simplest nontrivial knot."},
		{"Is the unknot trivial? (yes/no)", "yes", "Unknotted circle bounds a disc."},
		{"Is the figure-eight knot the same as the trefoil? (yes/no)", "no", "Figure-eight is 4_1 with four crossings; trefoil is 3_1."},
		{"Does the trefoil have crossing number 3? (yes/no)", "yes", "Simplest nontrivial knot."},
	}
	hard := []entry{{"Is Jones polynomial knot invariant? (yes/no)", "yes", "Jones."}, {"Does Gordon-Luecke show knot determined by complement? (yes/no)", "yes", "Knot complement."}, {"Is figure-eight knot 4_1? (yes/no)", "yes", "Four crossings."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type cohomologyGen struct{}

func (g *cohomologyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What is H^0 of a connected space? (type like Z)", "Z", "H^0 counts components: one copy of Z per component."},
		{"Is cohomology dual to homology? (yes/no)", "yes", "Universal coefficients pair them."},
		{"Does de Rham cohomology use differential forms? (yes/no)", "yes", "Smooth manifolds: closed forms modulo exact."},
		{"Is H^0 the group counting connected components? (yes/no)", "yes", "H^0 records components."},
	}
	hard := []entry{{"Does cohomology ring have cup product? (yes/no)", "yes", "Ring structure."}, {"Is singular cohomology homotopy invariant? (yes/no)", "yes", "Homotopy."}, {"Does Poincaré duality hold for compact orientable manifold? (yes/no)", "yes", "Duality."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type stoneCechGen struct{}

func (g *stoneCechGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"Which compactification is largest: one-point or Stone-Cech? (type with hyphen)", "Stone-Cech", "Stone-Cech is the universal (largest) compactification of Tychonoff spaces."},
		{"Does beta N contain a remainder N*? (yes/no)", "yes", "The growth beta N minus N is the remainder."},
		{"Is beta R metrizable? (yes/no)", "no", "Stone-Cech remainders are huge and non-metrizable."},
		{"Is Tychonoff exactly embeddable in compact Hausdorff? (yes/no)", "yes", "Characterization of Tychonoff spaces."},
	}
	hard := []entry{{"Is βR not metrizable? (yes/no)", "yes", "Huge."}, {"Does Stone-Čech of discrete is ultrafilters? (yes/no)", "yes", "βD."}, {"Is completely regular iff subspace of compact Hausdorff? (yes/no)", "yes", "Tychonoff."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type dimensionGen struct{}

func (g *dimensionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What is the covering dimension of R^2?", "2", "dim R^n = n, so the plane has dimension 2."},
		{"What is the covering dimension of the Cantor set?", "0", "Totally disconnected compacta have dimension 0."},
		{"Is covering dimension of R^n equal n? (yes/no)", "yes", "Lebesgue dimension of Euclidean n-space is n."},
		{"Does dimension increase under subspaces? (yes/no)", "no", "Subspaces never exceed the ambient dimension."},
	}
	hard := []entry{{"Is inductive dimension equal covering for separable metric? (yes/no)", "yes", "Coincidence."}, {"Does Brouwer show dimension is topological invariant? (yes/no)", "yes", "Invariance of domain."}, {"Is Menger-Nöbeling embed n-dim compact in R^{2n+1}? (yes/no)", "yes", "Embedding."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type coveringGen struct{}

func (g *coveringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"How many sheets does z-map-z-squared have over the circle?", "2", "Each nonzero w has two square roots, so two sheets."},
		{"Is a covering map a local homeomorphism? (yes/no)", "yes", "Every point has an evenly covered neighbourhood."},
		{"Is the universal cover simply connected? (yes/no)", "yes", "Universal covers are simply connected by definition."},
		{"Does every covering have unique path lifting? (yes/no)", "yes", "Paths lift uniquely from a chosen starting sheet."},
	}
	hard := []entry{{"Does covering correspond to subgroup of π_1? (yes/no)", "yes", "Classification."}, {"Is deck transformation group π_1? (yes/no)", "yes", "Deck."}, {"Does covering of path-connected is path-connected? (yes/no)", "yes", "Connected."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type vanKampenGen struct{}

func (g *vanKampenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What is pi_1 of a wedge of 2 circles? (type like F2)", "F2", "A wedge of n circles has free fundamental group on n generators."},
		{"Does Van Kampen compute pi_1 of a union? (yes/no)", "yes", "Van Kampen glues pi_1 from path-connected opens."},
		{"Is pi_1 of a wedge of circles free? (yes/no)", "yes", "Wedge of circles gives a free group."},
		{"Does Van Kampen need path-connected intersection? (yes/no)", "yes", "The intersection must be path-connected for the gluing."},
	}
	hard := []entry{{"Is π_1 torus via Van Kampen Z×Z? (yes/no)", "yes", "Torus."}, {"Does Van Kampen give amalgamated product? (yes/no)", "yes", "Amalgam."}, {"Is Seifert-Van Kampen for open cover? (yes/no)", "yes", "Open."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type fiberBundleGen struct{}

func (g *fiberBundleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What is the fiber of the Hopf fibration? (type like S^1)", "S^1", "Hopf S^3→S^2 has circle fibers."},
		{"Is a covering a fiber bundle with discrete fiber? (yes/no)", "yes", "Coverings are bundles with discrete fibers."},
		{"Does a fiber bundle have local trivializations? (yes/no)", "yes", "Local product structure is the definition."},
		{"Is the tangent bundle of S^1 trivial? (yes/no)", "yes", "The circle has a nonvanishing vector field."},
	}
	hard := []entry{{"Is Hopf fibration S^3→S^2 fiber bundle? (yes/no)", "yes", "Hopf."}, {"Does Möbius band fiber bundle over S^1? (yes/no)", "yes", "Möbius."}, {"Is vector bundle fiber is vector space? (yes/no)", "yes", "Vector."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type cellularHomologyGen struct{}

func (g *cellularHomologyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What kind of complex does cellular homology use? (type like CW)", "CW", "Cellular chains come from the CW skeleton."},
		{"Does cellular homology use the CW skeleton? (yes/no)", "yes", "Cells in each dimension give the chain groups."},
		{"Is the cellular chain built from cells? (yes/no)", "yes", "One generator per cell."},
		{"Does cellular homology agree with singular homology? (yes/no)", "yes", "They are isomorphic for CW complexes."},
	}
	hard := []entry{{"Is cellular boundary via degree? (yes/no)", "yes", "Degree."}, {"Does CW complex have cellular homology? (yes/no)", "yes", "CW."}, {"Is cellular homology finitely generated for finite CW? (yes/no)", "yes", "Finite."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type deRhamGen struct{}

func (g *deRhamGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"What does de Rham cohomology use: forms or chains? (type one word)", "forms", "Closed differential forms modulo exact forms."},
		{"Is de Rham cohomology homotopy invariant? (yes/no)", "yes", "Homotopic maps induce the same map on de Rham cohomology."},
		{"Does the de Rham theorem relate it to singular cohomology? (yes/no)", "yes", "Integration gives the isomorphism with real coefficients."},
		{"Does de Rham cohomology vanish above the manifold dimension? (yes/no)", "yes", "No nonzero k-forms exist above the dimension."},
	}
	hard := []entry{{"Does Poincaré lemma give H^n of R^n? (yes/no)", "yes", "Poincaré."}, {"Is de Rham cohomology ring via wedge product? (yes/no)", "yes", "Ring."}, {"Does de Rham cohomology vanish above dimension? (yes/no)", "yes", "Dimension."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
