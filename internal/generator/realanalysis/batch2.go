package realanalysis

// The second real-analysis batch: the rest of limits and continuity, the
// derivative and the integral made rigorous, uniform convergence, and the first
// look at metric spaces.
//
// The first batch built the real line from completeness and proved the theorems
// about sequences, series and the first two-thirds of continuity. This batch
// finishes continuity with connectedness, the classification of discontinuities
// and monotone functions; proves the mean value theorem, Taylor's theorem and
// the Riemann integral with the fundamental theorem, integration by parts and
// improper integrals; treats uniform convergence and the interchange of limits;
// and opens the metric-space picture with open and closed sets, closure,
// completeness, compactness, continuity, the contraction mapping principle,
// Baire's theorem and equivalent metrics.
//
// Ten concepts are numeric, so the definitions can be exercised, and fifteen are
// multiple choice, because the domain is largely about statements and proofs.
// Every numeric generator publishes the parameters its answer follows from, and
// the tests recompute the answer rather than trusting the author's arithmetic.

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func raIntPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// ---------------------------------------------------------------------------
// Connectedness (multiple choice)
// ---------------------------------------------------------------------------

type connectedGen struct{}

func (g *connectedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A subset of the real line is connected exactly when...",
			correct:     "it is an interval",
			distractors: []string{"it is open", "it is closed and bounded", "it has no isolated points"},
			reason:      "The connected subsets of the line are exactly the intervals, bounded or not and open, closed or half-open.",
		},
		{
			q:           "Continuity preserves which property of connected sets when the image is taken?",
			correct:     "connectedness",
			distractors: []string{"open", "compact", "closed"},
			reason:      "A continuous function cannot split a connected set into two separated pieces.",
		},
		{
			q:           "The union of the intervals from 0 to 1 and from 2 to 3 is...",
			correct:     "disconnected",
			distractors: []string{"connected", "compact but not closed", "a single interval"},
			reason:      "The gap between 1 and 2 splits the union into two nonempty pieces separated in the subspace.",
		},
	}
	return pickMC(append(claims, computedLimitClaims("ra.limit.connected")...))
}

// ---------------------------------------------------------------------------
// Classification of discontinuities (multiple choice)
// ---------------------------------------------------------------------------

type discontinuityGen struct{}

func (g *discontinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A removable discontinuity at a point is one where...",
			correct:     "the two-sided limit exists but differs from the value, or the value is undefined",
			distractors: []string{"the one-sided limits differ", "at least one one-sided limit fails to exist", "the function is unbounded near the point"},
			reason:      "Removable means the gap can be closed by redefining the single value, because the limit itself exists.",
		},
		{
			q:           "A jump discontinuity at a point is one where...",
			correct:     "both one-sided limits exist but they are different",
			distractors: []string{"the two-sided limit exists but differs from the value", "the function is unbounded near the point", "the function is continuous from the left only in value"},
			reason:      "A jump is a finite step: the left and right limits exist separately and disagree.",
		},
		{
			q:           "The function equal to 0 for every x other than 0 and equal to 1 at 0 has...",
			correct:     "a removable discontinuity at 0",
			distractors: []string{"a jump discontinuity at 0", "an essential discontinuity at 0", "no discontinuity at 0"},
			reason:      "The limit at 0 is 0 while the value is 1, so redefining the single value removes the discontinuity.",
		},
	}
	return pickMC(append(claims, computedLimitClaims("ra.limit.discontinuity")...))
}

// ---------------------------------------------------------------------------
// Monotone functions (multiple choice)
// ---------------------------------------------------------------------------

type monotoneFnGen struct{}

func (g *monotoneFnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A monotone function on an interval can have...",
			correct:     "at most countably many discontinuities, and all of them are jumps",
			distractors: []string{"uncountably many discontinuities", "only removable discontinuities", "no discontinuities at all"},
			reason:      "Each jump straddles a rational, so the jumps can be counted, and monotonicity forbids other kinds.",
		},
		{
			q:           "A monotone function on a closed bounded interval is...",
			correct:     "Riemann integrable",
			distractors: []string{"continuous", "differentiable almost everywhere is false in general", "uniformly continuous"},
			reason:      "Its only discontinuities are jumps, which form a set of measure zero, so it is Riemann integrable.",
		},
		{
			q:           "If f is increasing on an interval ending at b, then the left limit of f at b...",
			correct:     "exists and equals the supremum of the values on the interval",
			distractors: []string{"always equals f(b)", "need not exist", "equals the infimum of the values"},
			reason:      "The values are bounded above by f(b) and increase, so their supremum is approached from the left.",
		},
	}
	return pickMC(append(claims, computedLimitClaims("ra.limit.monotone_fn")...))
}

// ---------------------------------------------------------------------------
// Darboux's theorem (multiple choice)
// ---------------------------------------------------------------------------

type darbouxGen struct{}

func (g *darbouxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "Darboux's theorem says that a derivative...",
			correct:     "has the intermediate value property even though it need not be continuous",
			distractors: []string{"is always continuous", "is always bounded", "is always monotone"},
			reason:      "Derivatives take every intermediate value between two of their values, a property continuity implies but that does not imply continuity.",
		},
		{
			q:           "Darboux's theorem rules out which behaviour in a derivative?",
			correct:     "a jump discontinuity",
			distractors: []string{"an unbounded oscillation", "a removable discontinuity", "a discontinuity at a single point"},
			reason:      "A jump skips values between the two one-sided limits, which the intermediate value property forbids.",
		},
		{
			q:           "The function f(x) = x² sin(1/x) for x other than 0, with f(0) = 0, is differentiable at 0, and its derivative...",
			correct:     "is not continuous at 0 but still has the intermediate value property",
			distractors: []string{"is continuous at 0", "does not exist at points near 0", "is monotone near 0"},
			reason:      "The derivative oscillates without a limit at 0, which Darboux's theorem permits because it is not a jump.",
		},
	}
	return pickMC(append(claims, computedDerivClaims("ra.deriv.darboux")...))
}

// ---------------------------------------------------------------------------
// Uniform convergence (multiple choice)
// ---------------------------------------------------------------------------

type uniformConvGen struct{}

func (g *uniformConvGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A sequence of functions converges uniformly on a set when...",
			correct:     "one index N works for every point of the set at once",
			distractors: []string{"each point has its own index but the indices are bounded", "it converges at every point", "the functions are continuous"},
			reason:      "Uniform convergence demands a single threshold that controls the whole domain, not one threshold per point.",
		},
		{
			q:           "Pointwise convergence of a sequence of functions...",
			correct:     "does not imply uniform convergence",
			distractors: []string{"implies uniform convergence", "is the same as uniform convergence", "implies the limit is continuous"},
			reason:      "The index that works can drift without bound as the point moves, as x^n on the unit interval shows.",
		},
		{
			q:           "Which property survives a uniform limit of continuous functions?",
			correct:     "continuity",
			distractors: []string{"differentiable", "monotone", "bounded but not necessarily continuous"},
			reason:      "Uniform control lets the limit be exchanged with the pointwise limit, which preserves continuity.",
		},
	}
	return pickMC(append(claims, computedDerivClaims("ra.deriv.uniform_conv")...))
}

// ---------------------------------------------------------------------------
// Interchange of limits (multiple choice)
// ---------------------------------------------------------------------------

type interchangeGen struct{}

func (g *interchangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "If a sequence of continuous functions converges uniformly on a closed bounded interval, then the limit of the integrals is...",
			correct:     "the integral of the limit function",
			distractors: []string{"the limit of the values at the midpoint", "not defined", "the integral of the pointwise limit only when the functions are monotone"},
			reason:      "Uniform convergence justifies passing the limit through the integral sign.",
		},
		{
			q:           "Term-by-term differentiation of a series of functions is valid when...",
			correct:     "the series of derivatives converges uniformly",
			distractors: []string{"the series of values converges at one point", "the terms are continuous", "the terms tend to zero uniformly"},
			reason:      "Differentiation is the delicate operation; it needs uniform convergence of the derived series, not of the original.",
		},
		{
			q:           "Term-by-term integration of a uniformly convergent series on a closed bounded interval is...",
			correct:     "always valid",
			distractors: []string{"valid only when the series converges at the endpoints", "valid only for power series", "never valid without differentiability"},
			reason:      "Integrals are bounded by the interval length, so uniform convergence controls the tail of the integrated series.",
		},
	}
	return pickMC(append(claims, computedDerivClaims("ra.deriv.interchange")...))
}

// ---------------------------------------------------------------------------
// Open and closed sets (multiple choice)
// ---------------------------------------------------------------------------

type openClosedGen struct{}

func (g *openClosedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A subset of a metric space is open when...",
			correct:     "every point of it has an open ball contained in the set",
			distractors: []string{"it contains its limit points", "its complement is bounded", "it is an interval"},
			reason:      "An open set gives every point some room entirely inside the set.",
		},
		{
			q:           "A subset of a metric space is closed exactly when...",
			correct:     "its complement is open",
			distractors: []string{"it is bounded", "it is an interval", "it is not open"},
			reason:      "Closed is defined as the complement of open, from which the limit-point description follows.",
		},
		{
			q:           "In a metric space, an arbitrary union of open sets is open, and a finite intersection of open sets is open; an arbitrary intersection of open sets is...",
			correct:     "not necessarily open",
			distractors: []string{"always open", "always closed", "always empty"},
			reason:      "Intersecting the balls of radius 1/n about a point shrinks to a point, which is not open.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.open_closed")...))
}

// ---------------------------------------------------------------------------
// Closure and interior (multiple choice)
// ---------------------------------------------------------------------------

type closureGen struct{}

func (g *closureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "The closure of a set is...",
			correct:     "the set together with all of its limit points",
			distractors: []string{"the set without its boundary", "the largest open set it contains", "the set of its isolated points"},
			reason:      "Adding every point that can be approached from the set makes the smallest closed set containing it.",
		},
		{
			q:           "A point lies in the closure of a set exactly when...",
			correct:     "every open ball around it meets the set",
			distractors: []string{"it belongs to the set", "some open ball around it misses the set", "it is a limit of a constant sequence in the set"},
			reason:      "The closure is characterised by points that cannot be separated from the set by any ball.",
		},
		{
			q:           "The boundary of a set consists of...",
			correct:     "the points of the closure that are not in the interior",
			distractors: []string{"the points that belong to the set and to its complement", "the isolated points of the set", "the complement of the set"},
			reason:      "A boundary point is one every ball around which meets the set and its complement.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.closure")...))
}

// ---------------------------------------------------------------------------
// Convergence in metric spaces (multiple choice)
// ---------------------------------------------------------------------------

type metricConvGen struct{}

func (g *metricConvGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A sequence in a metric space converges to a point when...",
			correct:     "for every positive epsilon there is an index beyond which the distance to the point is below epsilon",
			distractors: []string{"the distances between consecutive terms tend to zero", "the sequence is bounded", "the sequence is eventually constant"},
			reason:      "The epsilon definition is the same as on the line, with the absolute value replaced by the metric.",
		},
		{
			q:           "A convergent sequence in a metric space is...",
			correct:     "Cauchy",
			distractors: []string{"bounded but not necessarily Cauchy", "eventually constant", "a Cauchy sequence only when the space is complete"},
			reason:      "Two terms beyond the threshold are each close to the limit, so they are close to each other.",
		},
		{
			q:           "In a general metric space, a Cauchy sequence...",
			correct:     "need not converge",
			distractors: []string{"always converges", "always has a convergent subsequence", "converges exactly when it is bounded"},
			reason:      "Convergence of every Cauchy sequence is completeness, which is an extra assumption.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.convergence")...))
}

// ---------------------------------------------------------------------------
// Completeness (multiple choice)
// ---------------------------------------------------------------------------

type completeGen struct{}

func (g *completeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A metric space is complete when...",
			correct:     "every Cauchy sequence in it converges to a point of the space",
			distractors: []string{"every sequence has a convergent subsequence", "it is closed and bounded", "every bounded sequence converges"},
			reason:      "Completeness is about the Cauchy sequences, and it is the metric form of the least upper bound property.",
		},
		{
			q:           "The real line with the usual metric is...",
			correct:     "complete",
			distractors: []string{"compact", "connected only if bounded", "not complete because it is unbounded"},
			reason:      "Every Cauchy sequence of reals converges, which is the completeness of the reals restated.",
		},
		{
			q:           "The rationals with the usual metric are...",
			correct:     "not complete",
			distractors: []string{"complete", "compact", "equivalent to the reals"},
			reason:      "Truncations of the square root of two are Cauchy in the rationals and have no rational limit.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.complete")...))
}

// ---------------------------------------------------------------------------
// Compactness (multiple choice)
// ---------------------------------------------------------------------------

type metricCompactGen struct{}

func (g *metricCompactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A metric space is compact when...",
			correct:     "every open cover has a finite subcover",
			distractors: []string{"it is bounded", "every sequence converges", "it is closed"},
			reason:      "The finite subcover property is the definition; closed and bounded is a theorem, not the definition.",
		},
		{
			q:           "In a metric space, compactness is equivalent to...",
			correct:     "every sequence having a convergent subsequence with limit in the space",
			distractors: []string{"every sequence being Cauchy", "every subset being closed", "being totally ordered"},
			reason:      "Sequential compactness and the finite subcover property agree in metric spaces.",
		},
		{
			q:           "A continuous real-valued function on a nonempty compact metric space...",
			correct:     "is bounded and attains its maximum and minimum",
			distractors: []string{"is uniformly continuous only if the space is an interval", "is injective", "has a continuous inverse"},
			reason:      "The image of a compact set is compact, and a compact subset of the reals is closed and bounded.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.compact")...))
}

// ---------------------------------------------------------------------------
// Continuity in metric spaces (multiple choice)
// ---------------------------------------------------------------------------

type metricContinuityGen struct{}

func (g *metricContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "A function between metric spaces is continuous at a point when...",
			correct:     "for every positive epsilon there is a positive delta that keeps images within epsilon whenever inputs are within delta",
			distractors: []string{"it maps open sets to open sets", "it maps convergent sequences to bounded sequences", "it is a distance-preserving map"},
			reason:      "The epsilon-delta condition is the definition, with two metrics in place of the usual distance.",
		},
		{
			q:           "A function between metric spaces is continuous exactly when...",
			correct:     "the preimage of every open set is open",
			distractors: []string{"the image of every open set is open", "the image of every closed set is closed", "the preimage of every bounded set is bounded"},
			reason:      "Continuity is characterised by open preimages; the forward image condition is a stronger, different property.",
		},
		{
			q:           "A continuous map from a compact metric space to a metric space is...",
			correct:     "uniformly continuous",
			distractors: []string{"a bijection", "an isometry", "Lipschitz"},
			reason:      "Compactness upgrades pointwise continuity to a single delta for the whole domain.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.continuity")...))
}

// ---------------------------------------------------------------------------
// Baire category (multiple choice)
// ---------------------------------------------------------------------------

type baireGen struct{}

func (g *baireGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "The Baire category theorem says a complete metric space...",
			correct:     "cannot be written as a countable union of nowhere dense sets",
			distractors: []string{"is compact", "is countable", "has a countable dense subset"},
			reason:      "A complete space is of the second category in itself, so it is not a meagre union.",
		},
		{
			q:           "A consequence of the Baire category theorem is that...",
			correct:     "the reals are not a countable union of nowhere dense sets",
			distractors: []string{"every countable set is closed", "every complete space is compact", "the rationals are complete"},
			reason:      "Applying the theorem to the complete real line gives a statement about its category.",
		},
		{
			q:           "A set of first category is...",
			correct:     "a countable union of nowhere dense sets",
			distractors: []string{"a set that is not open", "a countable set", "a set containing no interval"},
			reason:      "Meagre or first category is exactly a countable union of sets with empty interior closure.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.baire")...))
}

// ---------------------------------------------------------------------------
// Equivalent metrics (multiple choice)
// ---------------------------------------------------------------------------

type equivalentGen struct{}

func (g *equivalentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []mcClaim{
		{
			q:           "Two metrics on a set are topologically equivalent when...",
			correct:     "they have the same open sets, hence the same convergent sequences",
			distractors: []string{"they give the same distance to every pair", "one is twice the other", "they are both bounded"},
			reason:      "Topological equivalence preserves every notion defined from open sets, including limits and continuity.",
		},
		{
			q:           "A metric d and the bounded metric d/(1+d) are...",
			correct:     "topologically equivalent, though the second has every distance below one",
			distractors: []string{"not equivalent, because the second is bounded", "equal", "equivalent only on compact sets"},
			reason:      "The two metrics have the same open balls up to size, so they induce the same topology.",
		},
		{
			q:           "Two metrics are Lipschitz equivalent when...",
			correct:     "each is bounded by a constant multiple of the other, which is stronger than topological equivalence",
			distractors: []string{"they have the same open sets, which is the same condition", "they are equal on a dense subset", "they both make the space complete"},
			reason:      "A two-sided constant bound forces the same topology and also controls Cauchy sequences and completeness.",
		},
	}
	return pickMC(append(claims, computedMetricClaims("ra.metric.equivalent")...))
}

// ---------------------------------------------------------------------------
// The mean value theorem (numeric)
// ---------------------------------------------------------------------------

type mvtGen struct{}

func (g *mvtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{2, 4, 6, 8, 10, 3}
	b := claims[rand.Intn(len(claims))]
	answer := decimal6(float64(b) / 2)
	return generator.Problem{
		Question: "For f(x) = x² on the interval from 0 to " + strconv.Itoa(b) +
			", the mean value theorem gives a point c with f'(c) equal to the average slope. What is c? (enter a decimal)",
		Answer:      answer,
		Explanation: "The average slope is (f(b) − f(0))/b = b²/b = b, and f'(c) = 2c, so 2c = b and c = b/2.",
		Facts: map[string]string{
			"b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The chain rule (numeric)
// ---------------------------------------------------------------------------

type chainGen struct{}

func (g *chainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		a, b, n, x0 int
	}
	claims := []claim{
		{a: 2, b: 1, n: 3, x0: 1},
		{a: 1, b: 2, n: 4, x0: 0},
		{a: 3, b: 1, n: 2, x0: 2},
		{a: 2, b: 3, n: 3, x0: 0},
		{a: 1, b: 1, n: 5, x0: 1},
		{a: 4, b: 2, n: 2, x0: 1},
	}
	c := claims[rand.Intn(len(claims))]
	inner := c.a*c.x0 + c.b
	answer := decimal6(float64(c.n * c.a * raIntPow(inner, c.n-1)))
	return generator.Problem{
		Question: "By the chain rule, the derivative of f(x) = (ax + b)^n is f'(x) = n·a·(ax + b)^(n−1). For a = " +
			strconv.Itoa(c.a) + ", b = " + strconv.Itoa(c.b) + ", n = " + strconv.Itoa(c.n) +
			", what is f'(x) at x = " + strconv.Itoa(c.x0) + "? (enter a decimal)",
		Answer: answer,
		Explanation: "The inner value is a·x + b = " + strconv.Itoa(inner) +
			", so f'(x) = n·a·(inner)^(n−1) = " + strconv.Itoa(c.n) + "·" + strconv.Itoa(c.a) + "·" +
			strconv.Itoa(raIntPow(inner, c.n-1)) + " = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "n": strconv.Itoa(c.n),
			"x0": strconv.Itoa(c.x0), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// L'Hopital's rule (numeric)
// ---------------------------------------------------------------------------

type lhopitalGen struct{}

func (g *lhopitalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{2, 3, 4, 5, 6, 1}
	a := claims[rand.Intn(len(claims))]
	answer := decimal6(float64(a*a) / 2)
	return generator.Problem{
		Question: "The limit of (1 − cos(ax))/x² as x tends to 0 is a²/2, by L'Hopital's rule applied twice. For a = " +
			strconv.Itoa(a) + ", what is the limit? (enter a decimal)",
		Answer:      answer,
		Explanation: "Both numerator and denominator tend to zero; differentiating twice gives a²cos(ax)/2, and cos(0) = 1, so the limit is a²/2 = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The Riemann integral (numeric)
// ---------------------------------------------------------------------------

type riemannGen struct{}

func (g *riemannGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		b, n int
	}
	claims := []claim{
		{b: 1, n: 4},
		{b: 2, n: 4},
		{b: 3, n: 5},
		{b: 2, n: 10},
		{b: 1, n: 5},
		{b: 4, n: 8},
	}
	c := claims[rand.Intn(len(claims))]
	sum := float64(c.b*c.b) * float64(c.n+1) / (2 * float64(c.n))
	answer := decimal6(sum)
	return generator.Problem{
		Question: "Use the right-endpoint Riemann sum with n = " + strconv.Itoa(c.n) + " equal subintervals to approximate the integral of f(x) = x on the interval from 0 to " +
			strconv.Itoa(c.b) + ". What is the value of the sum? (enter a decimal)",
		Answer:      answer,
		Explanation: "With step b/n the sum is Σ (i·b/n)(b/n) for i from 1 to n, which equals b²·n(n+1)/2/n² = b²(n+1)/(2n) = " + answer + ".",
		Facts: map[string]string{
			"b": strconv.Itoa(c.b), "n": strconv.Itoa(c.n), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The fundamental theorem of calculus (numeric)
// ---------------------------------------------------------------------------

type ftcGen struct{}

func (g *ftcGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(6)
	c := rand.Intn(6)
	b := 2 + rand.Intn(8)
	value := float64(a)*float64(b*b)/2 + float64(c*b)
	answer := decimal6(value)
	return generator.Problem{
		Question: "By the fundamental theorem of calculus, the integral from 0 to b of (a·t + c) with respect to t is a·b²/2 + c·b. For a = " +
			strconv.Itoa(a) + ", c = " + strconv.Itoa(c) + " and b = " + strconv.Itoa(b) + ", what is the value? (enter a decimal)",
		Answer:      answer,
		Explanation: "An antiderivative of a·t + c is a·t²/2 + c·t, and evaluating it at b and at 0 gives a·b²/2 + c·b = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "c": strconv.Itoa(c), "b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Integration by parts (numeric)
// ---------------------------------------------------------------------------

type partsGen struct{}

func (g *partsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Six upper limits instead of three. The value is computed from the closed
	// form rather than looked up, so each limit is a different problem.
	claims := []int{1, 2, 3, 4, 5, 6}
	b := claims[rand.Intn(len(claims))]
	value := float64(b-1)*math.Exp(float64(b)) + 1
	answer := decimal6(value)
	return generator.Problem{
		Question: "Use integration by parts to evaluate the integral from 0 to b of x·e^x with respect to x; its value is (b − 1)e^b + 1. For b = " +
			strconv.Itoa(b) + ", what is the value? (enter a decimal)",
		Answer:      answer,
		Explanation: "Integrating by parts with u = x and dv = e^x dx gives [x·e^x − e^x] from 0 to b, which is (b − 1)e^b + 1 = " + answer + ".",
		Facts: map[string]string{
			"b": strconv.Itoa(b), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Improper integrals (numeric)
// ---------------------------------------------------------------------------

type improperGen struct{}

func (g *improperGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Every exponent above one is convergent, so the pool is the range rather
	// than a short list of selected values.
	p := 2 + rand.Intn(38)
	answer := decimal6(1 / float64(p-1))
	return generator.Problem{
		Question: "The improper integral from 1 to infinity of x^(−p) converges exactly when p > 1, to 1/(p − 1). For p = " +
			strconv.Itoa(p) + ", what is its value? (enter a decimal)",
		Answer:      answer,
		Explanation: "An antiderivative is x^(1−p)/(1−p), which tends to 0 at infinity when p > 1, leaving 1/(p − 1) = " + answer + ".",
		Facts: map[string]string{
			"p": strconv.Itoa(p), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Power series (numeric)
// ---------------------------------------------------------------------------

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{2, 3, 4, 5, 8, 10}
	c := claims[rand.Intn(len(claims))]
	answer := decimal6(1 / float64(c))
	return generator.Problem{
		Question: "Find the radius of convergence \\(R\\) of the power series \\(\\sum_{n=1}^\\infty " +
			strconv.Itoa(c) + "^n x^n\\). (enter as a decimal)",
		Answer:      answer,
		Explanation: "By the root test (or ratio test), \\(\\limsup_{n\\to\\infty} |" + strconv.Itoa(c) + "^n|^{1/n} = " + strconv.Itoa(c) + "\\), so the radius of convergence is \\(R = \\frac{1}{" + strconv.Itoa(c) + "} = " + answer + "\\).",
		Facts: map[string]string{
			"c": strconv.Itoa(c), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Taylor's theorem (numeric)
// ---------------------------------------------------------------------------

type taylorGen struct{}

func (g *taylorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []int{1, 2, 3, 4, 5, 6}
	a := claims[rand.Intn(len(claims))]
	answer := decimal6(float64(a*a*a) / 6)
	return generator.Problem{
		Question: "The Taylor series at 0 of e^(a·x) has x³-coefficient a³/3! = a³/6. For a = " +
			strconv.Itoa(a) + ", what is that coefficient? (enter a decimal)",
		Answer:      answer,
		Explanation: "The n-th derivative of e^(a·x) at 0 is a^n, so the x³ coefficient is a³/6 = " + answer + ".",
		Facts: map[string]string{
			"a": strconv.Itoa(a), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// The contraction mapping principle (numeric)
// ---------------------------------------------------------------------------

type contractionGen struct{}

func (g *contractionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		pn, pd, b int
	}
	claims := []claim{
		{pn: 1, pd: 2, b: 1},
		{pn: 1, pd: 3, b: 2},
		{pn: 2, pd: 3, b: 1},
		{pn: 3, pd: 4, b: 2},
		{pn: 1, pd: 2, b: 3},
		{pn: 1, pd: 4, b: 1},
	}
	c := claims[rand.Intn(len(claims))]
	a := float64(c.pn) / float64(c.pd)
	fixed := float64(c.b) / (1 - a)
	answer := decimal6(fixed)
	return generator.Problem{
		Question: "The contraction T(x) = a·x + b with a = " + strconv.Itoa(c.pn) + "/" + strconv.Itoa(c.pd) +
			" and b = " + strconv.Itoa(c.b) + " has a unique fixed point, the limit of its iterates from any start. What is the fixed point? (enter a decimal)",
		Answer: answer,
		Explanation: "The fixed point solves x = a·x + b, so x = b/(1 − a) = " + strconv.Itoa(c.b) + "/(1 − " +
			strconv.Itoa(c.pn) + "/" + strconv.Itoa(c.pd) + ") = " + answer + ".",
		Facts: map[string]string{
			"pn": strconv.Itoa(c.pn), "pd": strconv.Itoa(c.pd), "b": strconv.Itoa(c.b), "answer": answer,
		},
	}
}
