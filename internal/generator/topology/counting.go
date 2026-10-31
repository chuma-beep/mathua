package topology

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// topoVerdict is one candidate in a counting question. yes records whether the
// candidate has the property the question asks about, and note says why in one
// short clause. Every pool holds eight yes cases and eight no cases, so drawing
// eight candidates gives a count anywhere from zero to eight and the answer
// stays a small integer while the stem changes with every draw.
type topoVerdict struct {
	label string
	yes   bool
	note  string
}

const topoDraw = 8

// countTopoVerdicts draws eight distinct candidates from pool and asks how many
// have the property. The answer is the count, so the question is production
// rather than recognition; the explanation still names every verdict.
func countTopoVerdicts(lead string, pool []topoVerdict) generator.Problem {
	var yes, no []topoVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	if len(yes) > topoDraw {
		yes = yes[:topoDraw]
	}
	if len(no) > topoDraw {
		no = no[:topoDraw]
	}
	// Choose the answer first, then fill the list to match it. Sampling eight
	// candidates and counting would reach 0 and 8 only when every draw lands on
	// the same side, which is too rare to keep all nine answers alive.
	var chosen []topoVerdict
	var question string
	count := 0
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(topoDraw + 1)
		if count > len(yes) {
			count = len(yes)
		}
		if topoDraw-count > len(no) {
			count = topoDraw - len(no)
		}
		rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
		rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
		chosen = append([]topoVerdict{}, yes[:count]...)
		chosen = append(chosen, no[:topoDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

		parts := make([]string, 0, topoDraw)
		for i, v := range chosen {
			parts = append(parts, fmt.Sprintf("(%s) %s", string(rune('a'+i)), v.label))
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		// A label such as S^2 or Z/2 can contain the count itself. Redraw the
		// count until the stem does not hand the answer away.
		if !strings.Contains(question, strconv.Itoa(count)) {
			break
		}
	}

	notes := make([]string, 0, topoDraw)
	for _, v := range chosen {
		verdict := "no"
		if v.yes {
			verdict = "yes"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note))
	}
	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(count),
		Explanation: strings.Join(notes, "; ") + ".",
		Facts:       map[string]string{"count": strconv.Itoa(count)},
	}
}

func openClosedProblem() generator.Problem {
	a := rand.Intn(20) - 10
	b := a + 1 + rand.Intn(8)
	c := b + 1 + rand.Intn(5)
	d := c + 1 + rand.Intn(5)
	kind := "open"
	if rand.Intn(2) == 0 {
		kind = "closed"
	}
	pick := func(open, closed bool) bool {
		if kind == "open" {
			return open
		}
		return closed
	}
	pool := []topoVerdict{
		{fmt.Sprintf("(%d, %d)", a, b), pick(true, false), "an open interval is open and not closed"},
		{fmt.Sprintf("(-inf, %d)", b), pick(true, false), "an open ray is open and not closed"},
		{fmt.Sprintf("R \\ {%d}", a), pick(true, false), "removing a point leaves an open set that is not closed"},
		{fmt.Sprintf("(%d, %d) union (%d, %d)", a, b, c, d), pick(true, false), "a union of open intervals is open and not closed"},
		{fmt.Sprintf("(%d, inf)", b), pick(true, false), "an open ray is open and not closed"},
		{fmt.Sprintf("(-inf, %d) union (%d, inf)", a, b), pick(true, false), "a union of open rays is open and not closed"},
		{fmt.Sprintf("[%d, %d]", a, b), pick(false, true), "a closed interval is closed and not open"},
		{fmt.Sprintf("[%d, inf)", a), pick(false, true), "a closed ray is closed and not open"},
		{fmt.Sprintf("{%d, %d}", a, b), pick(false, true), "a finite set is closed and not open"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, c, d), pick(false, true), "a finite union of closed intervals is closed and not open"},
		{"N", pick(false, true), "the natural numbers form a closed set that is not open"},
		{fmt.Sprintf("{%d, %d, %d}", a, b, c), pick(false, true), "a finite set is closed and not open"},
		{"R", pick(true, true), "R is both open and closed"},
		{"the empty set", pick(true, true), "the empty set is both open and closed"},
		{fmt.Sprintf("(%d, %d]", a, b), pick(false, false), "a half-open interval is neither open nor closed"},
		{fmt.Sprintf("[%d, %d)", a, b), pick(false, false), "a half-open interval is neither open nor closed"},
		{"Q", pick(false, false), "the rationals have empty interior and closure R"},
		{"R minus Q", pick(false, false), "the irrationals have empty interior and closure R"},
	}
	return countTopoVerdicts(fmt.Sprintf("How many of these subsets of R are %s?", kind), pool)
}

func compactProblem() generator.Problem {
	a := rand.Intn(12) - 6
	b := a + 1 + rand.Intn(6)
	c := b + 1 + rand.Intn(4)
	d := c + 1 + rand.Intn(4)
	pool := []topoVerdict{
		{fmt.Sprintf("[%d, %d]", a, b), true, "a closed bounded interval is compact"},
		{fmt.Sprintf("{%d, %d}", a, b), true, "a finite set is compact"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, c, d), true, "a finite union of compact sets is compact"},
		{"the empty set", true, "the empty set is compact"},
		{fmt.Sprintf("{%d}", a), true, "a one-point set is compact"},
		{fmt.Sprintf("[%d, %d] union {%d}", a, b, c), true, "a finite union of compact sets is compact"},
		{"{0} union {1/n : n >= 1}", true, "the set is closed and bounded"},
		{fmt.Sprintf("[%d, %d] intersect Z", a, b), true, "a finite set is compact"},
		{fmt.Sprintf("(%d, %d)", a, b), false, "the open interval is bounded but not closed"},
		{fmt.Sprintf("[%d, inf)", a), false, "the ray is closed but not bounded"},
		{fmt.Sprintf("(%d, %d]", a, b), false, "the half-open interval is not closed"},
		{"Z", false, "the integers are closed but not bounded"},
		{"R", false, "R is closed but not bounded"},
		{"Q", false, "the rationals are not closed in R"},
		{"R minus Q", false, "the irrationals are not closed in R"},
		{fmt.Sprintf("(%d, inf) intersect Q", a), false, "the set is not closed in R"},
	}
	return countTopoVerdicts("How many of these subsets of R are compact?", pool)
}

func connectedProblem() generator.Problem {
	a := rand.Intn(12) - 6
	b := a + 1 + rand.Intn(6)
	c := b + 1 + rand.Intn(4)
	d := c + 1 + rand.Intn(4)
	pool := []topoVerdict{
		{fmt.Sprintf("(%d, %d)", a, b), true, "an interval is connected"},
		{fmt.Sprintf("{%d}", a), true, "a one-point set is connected"},
		{fmt.Sprintf("[%d, %d]", a, b), true, "an interval is connected"},
		{"R", true, "R is connected"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, b, c), true, "the intervals overlap, so the union is connected"},
		{fmt.Sprintf("(%d, inf)", a), true, "a ray is an interval and is connected"},
		{fmt.Sprintf("(-inf, %d)", b), true, "a ray is an interval and is connected"},
		{fmt.Sprintf("[%d, inf)", a), true, "a ray is an interval and is connected"},
		{fmt.Sprintf("(%d, %d) union (%d, %d)", a, b, c, d), false, "the two intervals are separated"},
		{fmt.Sprintf("{%d, %d}", a, b), false, "a two-point set is disconnected"},
		{"Q", false, "the rationals are totally disconnected"},
		{"Z", false, "the integers are totally disconnected"},
		{"R minus Q", false, "the irrationals are totally disconnected"},
		{"the Cantor set", false, "the Cantor set is totally disconnected"},
		{fmt.Sprintf("{%d, %d, %d}", a, b, c), false, "a three-point set is disconnected"},
		{fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, d, d+1), false, "the two closed intervals are separated"},
	}
	return countTopoVerdicts("How many of these subsets of R are connected?", pool)
}

func continuousProblem() generator.Problem {
	a := rand.Intn(12) - 6
	k := 1 + rand.Intn(5)
	m := 1 + rand.Intn(5)
	pool := []topoVerdict{
		{fmt.Sprintf("x^2 + %dx + %d", k, m), true, "a polynomial is continuous"},
		{fmt.Sprintf("|x - %d|", a), true, "the absolute value of a continuous function is continuous"},
		{fmt.Sprintf("1/(x^2 + %d)", k), true, "the denominator never vanishes"},
		{"sin x", true, "the sine function is continuous"},
		{"e^x", true, "the exponential is continuous"},
		{"cos x", true, "the cosine function is continuous"},
		{"x^3", true, "a polynomial is continuous"},
		{fmt.Sprintf("x + 1/(x^2 + %d)", k), true, "a sum of continuous functions is continuous"},
		{"floor(x)", false, "the floor function jumps at every integer"},
		{fmt.Sprintf("1/(x - %d)", a), false, "the function has a pole at x = a"},
		{"1/x", false, "the function has a pole at x = 0"},
		{fmt.Sprintf("sqrt(x - %d)", a), false, "the function is not defined for x < a"},
		{"tan x", false, "the tangent has poles at odd multiples of pi/2"},
		{fmt.Sprintf("1/(x - %d)^2", a), false, "the function has a pole at x = a"},
		{"the step function", false, "the step function jumps at the origin"},
		{fmt.Sprintf("ln|x - %d|", a), false, "the logarithm has a pole at x = a"},
	}
	return countTopoVerdicts("How many of these functions are continuous on R?", pool)
}

func limitPointProblem() generator.Problem {
	pool := []topoVerdict{
		{"{1/n : n >= 1}", true, "the sequence 1/n converges to 0"},
		{"(-1, 1)", true, "every neighbourhood of 0 meets the interval"},
		{"(0, 1)", true, "every neighbourhood of 0 meets the interval"},
		{"Q", true, "the rationals are dense, so they accumulate at 0"},
		{"R minus Q", true, "the irrationals are dense, so they accumulate at 0"},
		{"the Cantor set", true, "0 is a non-isolated point of the Cantor set"},
		{"{1/n : n >= 1} union {0}", true, "0 is the limit of the sequence 1/n"},
		{"R", true, "every point of R is a limit point of R"},
		{"Z", false, "each integer has a ball of radius 1/2 that misses the rest of Z"},
		{"N", false, "each natural number is isolated in N"},
		{"{1}", false, "the singleton {1} has no other points near 0"},
		{"{1/2, 1/3}", false, "the set is finite and stays away from 0"},
		{"(1, 2)", false, "the interval is bounded away from 0"},
		{"{2^n : n >= 1}", false, "the set is unbounded and stays away from 0"},
		{"[1, 2]", false, "the interval is bounded away from 0"},
		{"{0}", false, "a neighbourhood of 0 meets the set only at 0 itself"},
	}
	return countTopoVerdicts("In R with the usual topology, for how many of these sets is 0 a limit point?", pool)
}

func completeMetricProblem() generator.Problem {
	pool := []topoVerdict{
		{"R", true, "every Cauchy sequence of reals converges"},
		{"C", true, "the complex plane is complete"},
		{"R^2", true, "Euclidean space is complete"},
		{"[0, 1]", true, "a closed subset of a complete space is complete"},
		{"Z", true, "Z is closed in R and therefore complete"},
		{"the Cantor set", true, "the Cantor set is compact and therefore complete"},
		{"a finite subset of R", true, "a finite set has no Cauchy sequence without a constant tail"},
		{"[0, 1] union [2, 3]", true, "a finite union of closed subsets of R is closed"},
		{"Q", false, "a Cauchy sequence of rationals can converge to sqrt(2), which is not rational"},
		{"(0, 1)", false, "1/n is Cauchy in (0, 1) but its limit 0 is missing"},
		{"R minus Q", false, "a sequence of irrationals can converge to a rational"},
		{"R minus {0}", false, "1/n is Cauchy in R minus {0} but its limit 0 is missing"},
		{"[0, 1)", false, "the sequence 1 - 1/n is Cauchy with limit 1 outside the set"},
		{"Q intersect [0, 1]", false, "the rationals in [0, 1] miss limits such as sqrt(2)/2"},
		{"(0, inf)", false, "the sequence 1/n is Cauchy with limit 0 outside the set"},
		{"{1/n : n >= 1}", false, "the sequence 1/n is Cauchy in the set but its limit 0 is missing"},
	}
	return countTopoVerdicts("How many of these metric spaces are complete with the usual metric?", pool)
}

func basisProblem() generator.Problem {
	pool := []topoVerdict{
		{"open intervals (a, b) with a, b real", true, "every open set is a union of open intervals"},
		{"open intervals (p, q) with p, q rational", true, "rational endpoints refine every open interval"},
		{"open balls B(q, r) with q, r rational", true, "rational centres and radii refine every ball"},
		{"open intervals with dyadic rational endpoints", true, "the dyadic rationals are dense in R"},
		{"open intervals with terminating decimal endpoints", true, "terminating decimals are dense in R"},
		{"open balls B(x, r) with x real and r rational", true, "a rational radius refines every ball"},
		{"open intervals with algebraic endpoints", true, "the algebraic numbers are countable and dense"},
		{"rays (-inf, q) and (q, inf) with q rational", true, "their finite intersections give rational open intervals"},
		{"closed intervals [a, b]", false, "a union of closed intervals need not be open"},
		{"singletons {x}", false, "unions of singletons give only discrete sets"},
		{"half-open intervals [a, b)", false, "the half-open end point cannot be produced by a union"},
		{"half-open intervals (a, b]", false, "the half-open end point cannot be produced by a union"},
		{"finite subsets of R", false, "finite unions stay finite and cannot give an open interval"},
		{"the integers Z", false, "unions of integers are never open intervals"},
		{"the singletons {n} for integers n", false, "the collection generates only subsets of Z"},
		{"closed balls with rational centres and radii", false, "unions of closed balls need not be open"},
	}
	return countTopoVerdicts("How many of these collections form a basis for the Euclidean topology on R?", pool)
}

func subspaceProblem() generator.Problem {
	pool := []topoVerdict{
		{"(0, 1)", true, "it is the trace of the open interval (0, 1)"},
		{"[0, 1/2)", true, "it is (-1, 1/2) intersect [0, 1]"},
		{"(1/2, 1]", true, "it is (1/2, 2) intersect [0, 1]"},
		{"[0, 1]", true, "the whole subspace is open in itself"},
		{"the empty set", true, "the empty set is open in every space"},
		{"[0, 1/3) union (2/3, 1]", true, "each piece is a trace of an ambient open set"},
		{"(1/4, 3/4)", true, "it is already open in R and lies inside Y"},
		{"[0, 1/4) union (3/4, 1]", true, "each end piece is a trace of an ambient open set"},
		{"{0}", false, "every neighbourhood of 0 in Y also contains points of (0, 1)"},
		{"{1}", false, "every neighbourhood of 1 in Y also contains points below 1"},
		{"{1/2}", false, "a singleton inside Y is not open in the subspace"},
		{"[0, 1/2]", false, "its right end point cannot be separated from points above 1/2"},
		{"(0, 1/2]", false, "its right end point cannot be separated from points above 1/2"},
		{"[1/3, 2/3]", false, "the closed interval contains its end points and is not open"},
		{"Q intersect [0, 1]", false, "the rationals are not open in the subspace"},
		{"Z intersect [0, 1]", false, "the finite set {0, 1} is not open in Y"},
	}
	return countTopoVerdicts("In the subspace Y = [0, 1] of R, how many of these subsets are open in Y?", pool)
}

func denseProblem() generator.Problem {
	pool := []topoVerdict{
		{"Q", true, "the rationals meet every non-empty open interval"},
		{"R minus Q", true, "the irrationals meet every non-empty open interval"},
		{"R", true, "R is closed and contains every point"},
		{"the dyadic rationals", true, "dyadic rationals approximate every real"},
		{"the terminating decimals", true, "terminating decimals approximate every real"},
		{"Q minus {0}", true, "removing one point leaves the rationals dense"},
		{"the algebraic numbers", true, "the algebraic numbers are countable and contain Q"},
		{"Q + sqrt(2)", true, "a translate of Q is still dense"},
		{"Z", false, "the integers are closed and discrete in R"},
		{"N", false, "the natural numbers are closed and discrete in R"},
		{"{0, 1, 2}", false, "a finite set is closed and not dense"},
		{"(0, 1)", false, "the closure of (0, 1) is [0, 1], not R"},
		{"{1/n : n >= 1}", false, "the set accumulates only at 0"},
		{"the Cantor set", false, "the Cantor set is closed and has empty interior"},
		{"[0, 1]", false, "the closure of [0, 1] is [0, 1], not R"},
		{"Q intersect [0, 1]", false, "the closure is [0, 1], not R"},
	}
	return countTopoVerdicts("How many of these subsets are dense in R with the Euclidean topology?", pool)
}

func homeomorphismProblem() generator.Problem {
	pool := []topoVerdict{
		{"connectedness", true, "the continuous image of a connected space is connected"},
		{"compactness", true, "the continuous image of a compact space is compact"},
		{"being Hausdorff", true, "Hausdorff separation is preserved by homeomorphism"},
		{"path-connectedness", true, "paths transport along a homeomorphism"},
		{"separability", true, "the image of a countable dense set stays countable and dense"},
		{"metrizability", true, "transporting the metric along a homeomorphism gives a metric"},
		{"second-countability", true, "the image of a countable basis is a countable basis"},
		{"the Euler characteristic", true, "it is computed from the same cell structure"},
		{"boundedness", false, "R and (0, 1) are homeomorphic but only one is bounded"},
		{"the diameter", false, "the diameter depends on the chosen metric"},
		{"completeness", false, "R and (0, 1) are homeomorphic but only R is complete"},
		{"total boundedness", false, "a bounded metric can be transported to an unbounded one"},
		{"being a Cauchy sequence", false, "Cauchyness depends on the metric, not the topology"},
		{"being open in R", false, "a homeomorphism of R need not preserve the ambient embedding"},
		{"being closed in R", false, "a homeomorphism of R need not preserve the ambient embedding"},
		{"the distance between two points", false, "distances are not topological data"},
	}
	return countTopoVerdicts("How many of these properties are preserved by every homeomorphism?", pool)
}

func hausdorffProblem() generator.Problem {
	pool := []topoVerdict{
		{"R with the usual topology", true, "metric spaces are Hausdorff"},
		{"[0, 1] with the subspace topology", true, "a subspace of a Hausdorff space is Hausdorff"},
		{"any metric space", true, "open balls separate distinct points"},
		{"the discrete topology on a two-point set", true, "singletons are disjoint open sets"},
		{"R^2 with the Euclidean topology", true, "the Euclidean metric is Hausdorff"},
		{"the Cantor set", true, "it is a subspace of R"},
		{"a product of two Hausdorff spaces", true, "products of Hausdorff spaces are Hausdorff"},
		{"Z with the discrete topology", true, "every discrete space is Hausdorff"},
		{"the cofinite topology on an infinite set", false, "any two nonempty open sets intersect"},
		{"the indiscrete topology on a two-point set", false, "the only open sets are the empty set and the whole space"},
		{"the Sierpinski space", false, "its two points cannot be separated"},
		{"the Zariski topology on R", false, "nonempty open sets have finite complements"},
		{"the cofinite topology on R", false, "any two nonempty open sets intersect"},
		{"the particular-point topology on R", false, "every nonempty open set contains the particular point"},
		{"the excluded-point topology on R", false, "the open sets overlap at every other point"},
		{"the cocountable topology on an uncountable set", false, "any two nonempty open sets intersect"},
	}
	return countTopoVerdicts("How many of these spaces are Hausdorff?", pool)
}

func urysohnProblem() generator.Problem {
	pool := []topoVerdict{
		{"a normal space with disjoint closed sets A and B has a continuous f with f = 0 on A and f = 1 on B", true, "this is Urysohn's lemma"},
		{"the separating function can be chosen continuous", true, "continuity is the point of the lemma"},
		{"the separating function takes values in [0, 1]", true, "the lemma produces a map into the unit interval"},
		{"the lemma applies to metric spaces", true, "every metric space is normal"},
		{"the lemma applies to compact Hausdorff spaces", true, "compact Hausdorff spaces are normal"},
		{"the lemma needs A and B closed and disjoint", true, "these are its hypotheses"},
		{"Tietze extension extends an Urysohn function", true, "Tietze generalises the lemma"},
		{"the lemma gives f(A) = {0} and f(B) = {1}", true, "the sets are separated exactly"},
		{"the lemma needs A and B open and disjoint", false, "the hypotheses use closed sets"},
		{"the lemma applies in every Hausdorff space", false, "Hausdorff alone does not give disjoint closed neighbourhoods"},
		{"the lemma makes f(A) = {1} and f(B) = {1}", false, "that does not separate the sets"},
		{"the lemma produces a homeomorphism", false, "it produces a continuous function, not a bijection"},
		{"the lemma requires X to be connected", false, "connectedness is irrelevant"},
		{"the lemma requires X to be second countable", false, "normality is the hypothesis"},
		{"the lemma only applies to finite spaces", false, "it applies to every normal space"},
		{"the separating function can be unbounded", false, "its values lie in [0, 1]"},
	}
	return countTopoVerdicts("How many of these statements about Urysohn's lemma are true?", pool)
}

func metrizationProblem() generator.Problem {
	pool := []topoVerdict{
		{"every second-countable regular Hausdorff space is metrizable", true, "this is the Urysohn metrization theorem"},
		{"R^n is metrizable", true, "the Euclidean metric induces its topology"},
		{"every metric space is normal", true, "metric spaces are perfectly normal"},
		{"every metrizable space is Hausdorff", true, "its metric separates points"},
		{"a subspace of a metrizable space is metrizable", true, "restrict the metric"},
		{"a countable product of metrizable spaces is metrizable", true, "use a bounded sum of metrics"},
		{"every discrete space is metrizable", true, "the discrete metric induces the topology"},
		{"every compact Hausdorff space with a countable base is metrizable", true, "it is regular and second countable"},
		{"every Hausdorff space is metrizable", false, "the cofinite topology on an infinite set is a counterexample"},
		{"every normal space is metrizable", false, "normality does not give a countable base"},
		{"second-countability alone implies metrizable", false, "regularity is also needed"},
		{"every metrizable space is compact", false, "R is metrizable and not compact"},
		{"the cofinite topology on R is metrizable", false, "it is not Hausdorff"},
		{"the Sorgenfrey line is metrizable", false, "it is separable but not second countable"},
		{"every separable space is metrizable", false, "separability alone is far too weak"},
		{"every connected space is metrizable", false, "connectedness says nothing about a metric"},
	}
	return countTopoVerdicts("How many of these statements about metrizability are true?", pool)
}

func quotientProblem() generator.Problem {
	pool := []topoVerdict{
		{"compactness", true, "the continuous image of a compact space is compact"},
		{"connectedness", true, "the continuous image of a connected space is connected"},
		{"path-connectedness", true, "a path lifts through the quotient map"},
		{"separability", true, "the image of a countable dense set stays dense"},
		{"the Lindelof property", true, "continuous images of Lindelof spaces are Lindelof"},
		{"continuity of the quotient map", true, "the quotient topology is defined to make it continuous"},
		{"density of the image of a dense subset", true, "continuity preserves dense images"},
		{"the quotient topology being the finest one that makes the map continuous", true, "that is its defining property"},
		{"the Hausdorff property", false, "R with two points identified has no disjoint neighbourhoods"},
		{"metrizability", false, "quotients of metric spaces need not be metrizable"},
		{"first countability", false, "the quotient of a first-countable space need not be first countable"},
		{"local compactness", false, "a quotient of a locally compact space need not be locally compact"},
		{"the quotient map being open", false, "the quotient map is not always open"},
		{"the quotient map being closed", false, "the quotient map is not always closed"},
		{"the quotient map being injective", false, "it is surjective by construction"},
		{"closedness of subsets", false, "a closed set can have a non-closed image"},
	}
	return countTopoVerdicts("How many of these properties are preserved by every quotient map?", pool)
}

func productProblem() generator.Problem {
	pool := []topoVerdict{
		{"the Hausdorff property", true, "points differ in some coordinate"},
		{"compactness", true, "this is Tychonoff's theorem"},
		{"connectedness", true, "products of connected spaces are connected"},
		{"path-connectedness", true, "concatenate coordinate paths"},
		{"the T1 property", true, "points are closed in each factor"},
		{"regularity", true, "regularity is productive"},
		{"complete regularity", true, "continuous functions separate points from closed sets"},
		{"second countability of a countable product", true, "use the countable product of countable bases"},
		{"local compactness of an arbitrary product", false, "an infinite product of compact spaces need not be locally compact"},
		{"first countability of an uncountable product", false, "uncountably many coordinates give uncountably many neighbourhoods"},
		{"metrizability of an uncountable product", false, "such a product is not first countable"},
		{"normality of an uncountable product", false, "products of normal spaces need not be normal"},
		{"the Lindelof property of an uncountable product", false, "uncountable products of compact spaces need not be Lindelof"},
		{"discreteness of an infinite product of nontrivial spaces", false, "basic open sets restrict only finitely many coordinates"},
		{"compactness in the box topology", false, "Tychonoff's theorem uses the product topology"},
		{"separability of every product of separable spaces", false, "an uncountable product of separable spaces can fail to be separable"},
	}
	return countTopoVerdicts("How many of these properties are preserved by arbitrary products?", pool)
}

func pathConnectedProblem() generator.Problem {
	pool := []topoVerdict{
		{"R with the usual topology", true, "the line is path-connected"},
		{"R^n for n >= 2", true, "segments join any two points"},
		{"the closed interval [0, 1]", true, "it is an interval"},
		{"the circle S^1", true, "an arc joins any two points"},
		{"the sphere S^2", true, "great circle arcs join any two points"},
		{"the torus T^2", true, "it is a product of circles"},
		{"an open ball in R^n", true, "it is convex"},
		{"R^2 minus the origin", true, "two points can be joined around the puncture"},
		{"R minus the origin", false, "the two rays are separate path components"},
		{"the rationals Q", false, "a path from 0 to 1 would have to cross irrational values"},
		{"two disjoint open intervals", false, "no path crosses the gap"},
		{"the set {1/n : n >= 1}", false, "it is totally disconnected"},
		{"Z with the usual topology", false, "it is discrete"},
		{"the Cantor set", false, "it is totally disconnected"},
		{"the topologist's sine curve", false, "it is connected but not path-connected"},
		{"the union of two disjoint circles", false, "the components are separated"},
	}
	return countTopoVerdicts("How many of these spaces are path-connected?", pool)
}

func localCompactProblem() generator.Problem {
	pool := []topoVerdict{
		{"R", true, "closed bounded intervals are compact neighbourhoods"},
		{"R^n", true, "closed balls are compact by Heine-Borel"},
		{"[0, 1]", true, "it is compact Hausdorff"},
		{"Z with the discrete topology", true, "singletons are compact open neighbourhoods"},
		{"a finite discrete space", true, "every singleton is compact and open"},
		{"an open interval (0, 1)", true, "it is an open subspace of R"},
		{"the closed unit ball in R^n", true, "it is compact"},
		{"any compact Hausdorff space", true, "the whole space is a compact neighbourhood of each point"},
		{"Q", false, "every neighbourhood has non-compact closure in Q"},
		{"the irrationals", false, "closures in the irrationals are never compact"},
		{"Q intersect [0, 1]", false, "missing irrational limits prevent compact closures"},
		{"Q^2", false, "no point has a compact neighbourhood"},
		{"(R minus Q) intersect [0, 1]", false, "its closures are not compact"},
		{"an infinite-dimensional Hilbert space", false, "the closed unit ball is not compact"},
		{"the Sorgenfrey line", false, "its compact subsets are countable"},
		{"the rational points on the unit circle", false, "they are dense but have no compact neighbourhoods"},
	}
	return countTopoVerdicts("How many of these spaces are locally compact?", pool)
}

func tietzeProblem() generator.Problem {
	pool := []topoVerdict{
		{"a continuous real-valued function on a closed subset of a normal space extends continuously", true, "this is the Tietze extension theorem"},
		{"R^n is normal", true, "metric spaces are normal"},
		{"every metric space is normal", true, "metric spaces are perfectly normal"},
		{"an unbounded continuous function on a closed subset of R^n extends to R^n", true, "the full theorem has no boundedness hypothesis"},
		{"Urysohn's lemma is the special case for two disjoint closed sets", true, "map the sets to 0 and 1"},
		{"compact Hausdorff spaces satisfy the hypotheses", true, "they are normal"},
		{"a finite subset of R is closed", true, "finite sets in a metric space are closed"},
		{"a bounded continuous function into [a, b] on a closed subset of R extends into [a, b]", true, "the interval-valued form of Tietze"},
		{"the subset A must be open", false, "Tietze requires A to be closed"},
		{"A must be compact", false, "closedness is the hypothesis"},
		{"the target must be compact", false, "the target can be all of R"},
		{"the theorem holds in every Hausdorff space", false, "normality is needed"},
		{"the extension is unique", false, "extensions are rarely unique"},
		{"the space X must be compact", false, "normality is enough"},
		{"the function f must be bounded", false, "unbounded functions also extend"},
		{"the extension is always linear", false, "Tietze gives continuity, not linearity"},
	}
	return countTopoVerdicts("How many of these statements about the Tietze extension theorem are true?", pool)
}

func homotopyProblem() generator.Problem {
	pool := []topoVerdict{
		{"R^2", true, "it is contractible"},
		{"R^3", true, "it is contractible"},
		{"S^2", true, "it is simply connected"},
		{"S^3", true, "it is simply connected"},
		{"the closed unit disc D^2", true, "it is contractible"},
		{"the interval [0, 1]", true, "it is contractible"},
		{"a single point", true, "its fundamental group is trivial"},
		{"any contractible space", true, "contractible spaces have trivial fundamental group"},
		{"S^1", false, "its fundamental group is Z"},
		{"R^2 minus the origin", false, "it deformation retracts onto S^1"},
		{"T^2", false, "its fundamental group is Z x Z"},
		{"RP^2", false, "its fundamental group is Z_2"},
		{"R^2 minus two points", false, "its fundamental group is the free group on two generators"},
		{"a wedge of two circles", false, "its fundamental group is the free group on two generators"},
		{"R^3 minus a circle", false, "it deformation retracts onto S^1"},
		{"the Mobius band", false, "it deformation retracts onto S^1"},
	}
	return countTopoVerdicts("How many of these spaces are simply connected?", pool)
}

func manifoldProblem() generator.Problem {
	pool := []topoVerdict{
		{"R^n", true, "it is locally Euclidean of dimension n"},
		{"S^n", true, "it is locally Euclidean of dimension n"},
		{"T^2", true, "it is a closed surface"},
		{"RP^2", true, "it is a closed surface"},
		{"an open subset of R^n", true, "it inherits Euclidean charts"},
		{"a product of an m-manifold and an n-manifold", true, "its dimension is m + n"},
		{"the boundary of a compact (n+1)-manifold", true, "the boundary is an n-manifold"},
		{"GL_n(R)", true, "it is an open subset of R^(n^2)"},
		{"Q", false, "it is totally disconnected and not locally Euclidean"},
		{"the irrationals", false, "they are totally disconnected"},
		{"the Cantor set", false, "it is totally disconnected"},
		{"the union of two intersecting lines", false, "no neighbourhood of the crossing point is an interval"},
		{"the figure-eight curve", false, "the crossing point has no Euclidean neighbourhood"},
		{"the closed interval [0, 1]", false, "its endpoints are boundary points"},
		{"the closed disc D^2", false, "its boundary points are not interior"},
		{"the letter Y", false, "the branch point has no interval neighbourhood"},
	}
	return countTopoVerdicts("How many of these spaces are topological manifolds without boundary?", pool)
}

func localPathProblem() generator.Problem {
	pool := []topoVerdict{
		{"R^n", true, "open balls are path-connected neighbourhoods"},
		{"an open ball in R^n", true, "it is homeomorphic to R^n"},
		{"a disjoint union of open balls", true, "each ball is a path-connected neighbourhood"},
		{"S^1", true, "small arcs are path-connected"},
		{"T^2", true, "it has a basis of path-connected coordinate patches"},
		{"a discrete space", true, "singletons are path-connected and open"},
		{"an open subset of R^n", true, "it has a basis of open balls"},
		{"a CW complex", true, "its cells give path-connected neighbourhood bases"},
		{"Q", false, "every nonempty open set is totally disconnected"},
		{"the irrationals", false, "every nonempty open set is totally disconnected"},
		{"the comb space", false, "points on the spine have no path-connected small neighbourhoods"},
		{"the topologist's sine curve", false, "it is not locally path-connected at the origin"},
		{"the Cantor set", false, "it is totally disconnected"},
		{"the set {1/n}", false, "each point is isolated but the set is not open"},
		{"Q^2", false, "every nonempty open set is totally disconnected"},
		{"the rational points on the unit circle", false, "small arcs contain no rational paths"},
	}
	return countTopoVerdicts("How many of these spaces are locally path-connected?", pool)
}

func onePointProblem() generator.Problem {
	pool := []topoVerdict{
		{"the one-point compactification of R^n is S^n", true, "inverse stereographic projection adds the north pole"},
		{"the one-point compactification of an open interval is S^1", true, "every open interval is homeomorphic to R"},
		{"the one-point compactification of an open ball is S^n", true, "an open ball is homeomorphic to R^n"},
		{"X* is compact for every space X", true, "the one-point compactification is always compact"},
		{"X* is Hausdorff exactly when X is locally compact Hausdorff", true, "this is Alexandroff's theorem"},
		{"neighbourhoods of infinity are complements of compact closed sets", true, "this is the defining topology at infinity"},
		{"if X is already compact then infinity is isolated", true, "the complement of X is empty, hence compact"},
		{"X is dense in X*", true, "infinity is the only added point"},
		{"X* is always Hausdorff", false, "a non-locally-compact X gives a non-Hausdorff X*"},
		{"X* is always metrizable", false, "metrizability needs stronger hypotheses"},
		{"X* is always connected", false, "a two-point discrete space gives three isolated points"},
		{"the one-point compactification of Q is Hausdorff", false, "Q is not locally compact"},
		{"X* adds an open point to every space", false, "infinity is isolated only when X is compact"},
		{"the compactification adds two points", false, "the construction adds exactly one point"},
		{"X* is compact exactly when X is compact", false, "X* is always compact"},
		{"X* is always a manifold", false, "the point at infinity need not have a Euclidean neighbourhood"},
	}
	return countTopoVerdicts("How many of these statements about the one-point compactification are true?", pool)
}

func completelyRegularProblem() generator.Problem {
	pool := []topoVerdict{
		{"every metric space is completely regular", true, "metric spaces are normal and normality implies complete regularity"},
		{"every normal space is completely regular", true, "Urysohn's lemma supplies the separating function"},
		{"every completely regular T1 space is regular", true, "separating a point from a closed set separates them by open sets"},
		{"a product of completely regular spaces is completely regular", true, "complete regularity is productive"},
		{"a subspace of a completely regular space is completely regular", true, "complete regularity is hereditary"},
		{"R^n is completely regular", true, "it is a metric space"},
		{"every compact Hausdorff space is completely regular", true, "compact Hausdorff spaces are normal"},
		{"a T1 space is completely regular exactly when points separate from closed sets by maps to [0, 1]", true, "this is the defining function form"},
		{"every completely regular space is normal", false, "complete regularity is weaker than normality"},
		{"every Hausdorff space is completely regular", false, "there are Hausdorff spaces that are not T3.5"},
		{"complete regularity is preserved by quotient maps", false, "quotients can collapse the separating functions"},
		{"every regular space is completely regular", false, "regularity alone does not supply the functions"},
		{"every compact space is completely regular", false, "compactness alone does not give the separation"},
		{"complete regularity implies compactness", false, "R is completely regular and not compact"},
		{"only metric spaces are completely regular", false, "compact Hausdorff spaces are another family"},
		{"closed subsets of a completely regular space fail complete regularity", false, "complete regularity is hereditary"},
	}
	return countTopoVerdicts("How many of these statements about complete regularity are true?", pool)
}

func orientabilityProblem() generator.Problem {
	pool := []topoVerdict{
		{"R^2", true, "Euclidean space has a global orientation"},
		{"S^2", true, "the outward normal field orients the sphere"},
		{"T^2", true, "the torus is a closed orientable surface"},
		{"the cylinder S^1 x [0, 1]", true, "it is a product of orientable spaces"},
		{"R^3", true, "Euclidean space has a global orientation"},
		{"a closed orientable surface of genus 2", true, "every orientable surface has a consistent orientation"},
		{"CP^2", true, "complex projective space is orientable"},
		{"S^1", true, "every circle has a consistent direction"},
		{"the Mobius band", false, "a loop around the band reverses orientation"},
		{"the Klein bottle", false, "it contains a Mobius band"},
		{"RP^2", false, "the antipodal quotient reverses orientation"},
		{"the nonorientable surface N_3", false, "connected sums of projective planes are nonorientable"},
		{"the open Mobius band", false, "removing the boundary does not restore orientability"},
		{"RP^2 minus a disc", false, "it is a Mobius band"},
		{"the connected sum RP^2 # T^2", false, "a nonorientable summand makes the sum nonorientable"},
		{"the nonorientable surface N_4", false, "every N_k with k at least 1 is nonorientable"},
	}
	return countTopoVerdicts("How many of these spaces are orientable?", pool)
}

func embeddingProblem() generator.Problem {
	pool := []topoVerdict{
		{"every embedding is injective", true, "an embedding is a homeomorphism onto its image"},
		{"an embedding is a homeomorphism onto its image", true, "this is the definition"},
		{"Whitney embeds a closed n-manifold in R^{2n}", true, "the Whitney embedding theorem gives this bound"},
		{"Nash embeds a Riemannian manifold isometrically in some R^N", true, "the Nash embedding theorem"},
		{"the inclusion of a subspace is an embedding", true, "the subspace topology makes the inclusion a homeomorphism onto its image"},
		{"a knot is an embedding of S^1 in R^3", true, "a knot is a tame embedded circle"},
		{"the Alexander horned sphere is a wild embedding of S^2", true, "its complement is not simply connected"},
		{"every embedding is continuous", true, "an embedding is a continuous map by definition"},
		{"every injective continuous map is an embedding", false, "the inverse need not be continuous"},
		{"every immersion is an embedding", false, "an immersion may self-intersect"},
		{"the Whitney bound is R^n", false, "the bound is R^{2n}"},
		{"the Alexander horned sphere is tame", false, "it is the standard wild embedding"},
		{"Nash's theorem is about topological embeddings", false, "it is about isometric embeddings of Riemannian manifolds"},
		{"an embedding need not be injective", false, "injectivity is part of the definition"},
		{"every embedding is an open map", false, "the image need not be open in the ambient space"},
		{"the image of an embedding is always closed", false, "an open interval embeds in R with open image"},
	}
	return countTopoVerdicts("How many of these statements about embeddings are true?", pool)
}

func knotProblem() generator.Problem {
	pool := []topoVerdict{
		{"the trefoil has crossing number 3", true, "three crossings is the minimum"},
		{"the unknot is trivial", true, "it bounds an embedded disc"},
		{"the figure-eight knot has crossing number 4", true, "it is 4_1"},
		{"the Jones polynomial is a knot invariant", true, "Reidemeister moves preserve it"},
		{"a knot complement determines the knot", true, "this is the Gordon-Luecke theorem"},
		{"the trefoil is chiral", true, "it is not equivalent to its mirror image"},
		{"a torus knot T(p, q) with coprime p, q > 1 is nontrivial", true, "it does not bound a disc in the solid torus"},
		{"Reidemeister moves preserve knot type", true, "they are the local moves of knot diagrams"},
		{"the trefoil is the unknot", false, "the trefoil is the simplest nontrivial knot"},
		{"the figure-eight knot is the trefoil", false, "they are 4_1 and 3_1"},
		{"crossing number is a complete invariant", false, "distinct knots can share a crossing number"},
		{"every knot is alternating", false, "nonalternating knots exist"},
		{"the unknot has crossing number 1", false, "the unknot has crossing number 0"},
		{"the Jones polynomial distinguishes all knots", false, "some distinct knots share it"},
		{"every knot has genus 1", false, "the unknot has genus 0 and higher genus knots exist"},
		{"every knot is a torus knot", false, "the figure-eight knot is not a torus knot"},
	}
	return countTopoVerdicts("How many of these statements about knots are true?", pool)
}

func cohomologyProblem() generator.Problem {
	pool := []topoVerdict{
		{"H^0 counts the connected components", true, "each component contributes one copy of Z"},
		{"cohomology is dual to homology with field coefficients", true, "the universal coefficient theorem gives the duality"},
		{"de Rham cohomology uses differential forms", true, "closed forms modulo exact forms"},
		{"the cup product gives cohomology a ring structure", true, "the product is graded commutative"},
		{"singular cohomology is homotopy invariant", true, "homotopic maps induce the same map"},
		{"Poincare duality holds for closed orientable manifolds", true, "H^k is dual to H^{n-k}"},
		{"H^k of a point is Z for k = 0 and 0 otherwise", true, "a point has one component and no higher cells"},
		{"integral cohomology can have torsion", true, "RP^2 has H^2 equal to Z/2"},
		{"H^0 always vanishes", false, "it is Z for each connected component"},
		{"cohomology is always isomorphic to homology", false, "torsion can differ without field coefficients"},
		{"de Rham cohomology applies to every topological space", false, "it needs a smooth manifold"},
		{"the cup product is commutative on the nose", false, "it is graded commutative"},
		{"Poincare duality holds for nonorientable manifolds", false, "orientability is required, or a twisted version is used"},
		{"H^k of R^n is Z for every k", false, "R^n is contractible, so only H^0 is Z"},
		{"cohomology cannot detect torsion", false, "integral cohomology detects it"},
		{"singular cohomology depends on the chosen triangulation", false, "it is a topological invariant"},
	}
	return countTopoVerdicts("How many of these statements about cohomology are true?", pool)
}

func stoneCechProblem() generator.Problem {
	pool := []topoVerdict{
		{"beta X is compact Hausdorff when X is Tychonoff", true, "the Stone-Cech compactification exists for Tychonoff spaces"},
		{"X is dense in beta X", true, "the embedding has dense image"},
		{"beta X is the largest compactification of X", true, "every other compactification factors through it"},
		{"beta N has cardinality 2^(2^aleph0)", true, "it is the space of ultrafilters on N"},
		{"beta R is not metrizable", true, "its remainder is too large"},
		{"beta X is unique up to homeomorphism", true, "the universal property determines it"},
		{"beta X equals X when X is compact Hausdorff", true, "a compact Hausdorff space is already its own compactification"},
		{"continuous bounded functions on X extend to beta X", true, "this is the defining universal property"},
		{"beta X is always metrizable", false, "beta R is not metrizable"},
		{"beta X is the smallest compactification", false, "it is the largest"},
		{"beta Q is metrizable", false, "its remainder is not metrizable"},
		{"beta N is countable", false, "it has the cardinality of the power set of the power set of N"},
		{"beta X is always connected", false, "beta of a two-point discrete space is disconnected"},
		{"beta X is always path-connected", false, "beta of a discrete space is totally disconnected"},
		{"beta X is always second countable", false, "beta N is not second countable"},
		{"beta X is the one-point compactification for every X", false, "they agree only in special cases"},
	}
	return countTopoVerdicts("How many of these statements about the Stone-Cech compactification are true?", pool)
}

func coveringProblem() generator.Problem {
	pool := []topoVerdict{
		{"a covering map is a local homeomorphism", true, "each point has an evenly covered neighbourhood"},
		{"the fibers of a covering map are discrete", true, "the fiber over a point has the discrete topology"},
		{"the universal cover is simply connected", true, "this is its defining property"},
		{"the number of sheets equals the index of the image subgroup", true, "the cosets of the image give the sheets"},
		{"coverings of S^1 are classified by degree", true, "each is z maps to z^d up to isomorphism"},
		{"the universal cover of T^2 is R^2", true, "the universal cover of a torus is the plane"},
		{"a covering of a connected space need not be connected", true, "a two-sheeted cover can be disconnected"},
		{"deck transformations form a group", true, "they compose and invert"},
		{"every local homeomorphism is a covering map", false, "it need not be evenly covered"},
		{"all fibers of a covering are finite", false, "the universal cover of S^1 has infinite fibers"},
		{"the universal cover of every space exists", false, "it needs semilocally simply connected hypotheses"},
		{"covering maps are always injective", false, "a nontrivial cover identifies distinct points"},
		{"the universal cover of S^1 is S^1", false, "it is R"},
		{"coverings preserve fundamental groups", false, "the induced map is injective but usually not surjective"},
		{"every covering is regular", false, "regular coverings are a proper subclass"},
		{"the universal cover of RP^2 is S^1", false, "it is S^2"},
	}
	return countTopoVerdicts("How many of these statements about covering spaces are true?", pool)
}

func vanKampenProblem() generator.Problem {
	pool := []topoVerdict{
		{"pi_1(S^1) is Z", true, "the circle has one loop generator"},
		{"pi_1(T^2) is Z x Z", true, "the torus has two independent loops"},
		{"pi_1(S^2) is trivial", true, "every loop on the sphere contracts"},
		{"van Kampen computes pi_1 from open sets with path-connected intersection", true, "this is the hypothesis of the theorem"},
		{"pi_1 of a wedge of two circles is free on two generators", true, "the Seifert-van Kampen theorem gives the free product"},
		{"the abelianization of pi_1 is H_1", true, "the Hurewicz theorem in degree one"},
		{"pi_1(RP^2) is Z/2", true, "the antipodal loop has order two"},
		{"pi_1 of a connected graph is free of rank E - V + 1", true, "the cycle rank gives the free rank"},
		{"pi_1(T^2) is Z", false, "the torus has two generators"},
		{"pi_1(S^2) is Z", false, "the sphere is simply connected"},
		{"van Kampen applies only to closed sets", false, "it applies to open covers"},
		{"pi_1 of the Klein bottle is Z x Z", false, "it is nonabelian"},
		{"pi_1 of a wedge is the direct product", false, "it is the free product in general"},
		{"H_1 determines pi_1", false, "abelianization loses information"},
		{"pi_1(RP^2) is Z", false, "it is Z/2"},
		{"pi_1 of every space is abelian", false, "the wedge of two circles is nonabelian"},
	}
	return countTopoVerdicts("How many of these statements about the fundamental group are true?", pool)
}

func fiberBundleProblem() generator.Problem {
	pool := []topoVerdict{
		{"a fiber bundle is locally a product", true, "this is the local triviality condition"},
		{"the Hopf fibration S^3 to S^2 has fiber S^1", true, "the fiber is a circle"},
		{"the tangent bundle of an n-manifold has fiber R^n", true, "each tangent space is n-dimensional"},
		{"a vector bundle has a vector space as fiber", true, "the fibers are vector spaces"},
		{"a covering space is a fiber bundle with discrete fiber", true, "the fiber is the discrete set of sheets"},
		{"the Mobius band is a line bundle over S^1", true, "it is the nontrivial real line bundle"},
		{"the Euler characteristic multiplies for surface bundles", true, "chi(E) equals chi(B) times chi(F)"},
		{"the projection of a smooth bundle is a submersion", true, "it has maximal rank"},
		{"every fiber bundle is trivial", false, "the Mobius band is nontrivial"},
		{"all fibers of a fiber bundle are connected", false, "a covering can have discrete fibers"},
		{"the Hopf fibration has fiber S^2", false, "its fiber is S^1"},
		{"the tangent bundle has fiber S^n", false, "its fiber is R^n"},
		{"fiber bundles are always coverings", false, "the fiber need not be discrete"},
		{"the base of a fiber bundle must be a manifold", false, "any topological space can be a base"},
		{"every bundle over a contractible base is nontrivial", false, "every such bundle is trivial"},
		{"the Euler characteristic is additive for bundles", false, "it is multiplicative for surface bundles"},
	}
	return countTopoVerdicts("How many of these statements about fiber bundles are true?", pool)
}
