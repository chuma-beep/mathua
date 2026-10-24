package topology

// The 2026-10-18 topology batch: twelve concepts that take the domain from 30
// to the target of 42.
//
// The survey found that the domain jumps from three basics straight to named
// theorems (Urysohn, Tietze, Stone-Cech, Van Kampen, de Rham) with almost no
// computational middle, and that 29 of its 30 concepts are multiple choice.
// The batch therefore builds the metric-space toolkit the later lessons assume
// (balls, interior and boundary, limit points, convergence, completeness,
// bases, subspaces, dense sets) and adds the computable invariants the
// algebraic-topology lessons describe but never name (homeomorphism, Euler
// characteristic, genus, Betti numbers).
//
// Six of the twelve are numeric, and every numeric generator publishes the
// facts its answer follows from, so the tests recompute the answer instead of
// comparing a string with itself. The six recognition concepts use
// generator.MultipleChoice, so every distractor is a real option the router
// must reject. The new `topology.invariants` subdomain holds the four invariant
// concepts; the others join the existing `basics` and `constructions` homes.

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("topo.metric.ball", &ballGen{})
	reg.Register("topo.metric.interior", &interiorGen{})
	reg.Register("topo.metric.limit_point", &limitPointGen{})
	reg.Register("topo.metric.convergence", &convergenceGen{})
	reg.Register("topo.metric.complete", &completeGen{})
	reg.Register("topo.metric.basis", &basisGen{})
	reg.Register("topo.metric.subspace", &subspaceGen{})
	reg.Register("topo.metric.dense", &denseGen{})
	reg.Register("topo.invariant.homeomorphism", &homeomorphismGen{})
	reg.Register("topo.invariant.euler", &eulerGen{})
	reg.Register("topo.invariant.genus", &genusGen{})
	reg.Register("topo.invariant.betti", &bettiGen{})
}

// ---------------------------------------------------------------------------
// 1. Open balls (numeric)
// ---------------------------------------------------------------------------

type ballGen struct{}

func (g *ballGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string // distance | endpoint
		x1, y1   int
		x2, y2   int
		center   int
		radius   int
		side     string // lower | upper
		reason   string
	}
	claims := []claim{
		{
			question: "In the Euclidean plane, what is the distance d((0,0),(3,4))? (enter a number)",
			answer:   "5", kind: "distance", x1: 0, y1: 0, x2: 3, y2: 4,
			reason: "The Euclidean distance is sqrt(3^2 + 4^2) = sqrt(25) = 5.",
		},
		{
			question: "In the Euclidean plane, what is the distance d((0,0),(6,8))? (enter a number)",
			answer:   "10", kind: "distance", x1: 0, y1: 0, x2: 6, y2: 8,
			reason: "The Euclidean distance is sqrt(6^2 + 8^2) = sqrt(100) = 10.",
		},
		{
			question: "In the Euclidean plane, what is the distance d((1,2),(4,6))? (enter a number)",
			answer:   "5", kind: "distance", x1: 1, y1: 2, x2: 4, y2: 6,
			reason: "The differences are 3 and 4, so the distance is sqrt(9 + 16) = 5.",
		},
		{
			question: "In R with the usual metric, the open ball of radius 2 centred at 3 is the interval (a, b). What is a? (enter a number)",
			answer:   "1", kind: "endpoint", center: 3, radius: 2, side: "lower",
			reason: "The ball is (3 - 2, 3 + 2) = (1, 5), so the lower endpoint is 1.",
		},
		{
			question: "In R with the usual metric, the open ball of radius 2 centred at 3 is the interval (a, b). What is b? (enter a number)",
			answer:   "5", kind: "endpoint", center: 3, radius: 2, side: "upper",
			reason: "The ball is (3 - 2, 3 + 2) = (1, 5), so the upper endpoint is 5.",
		},
		{
			question: "In R with the usual metric, the open ball of radius 3 centred at 1 is the interval (a, b). What is a? (enter a number)",
			answer:   "-2", kind: "endpoint", center: 1, radius: 3, side: "lower",
			reason: "The ball is (1 - 3, 1 + 3) = (-2, 4), so the lower endpoint is -2.",
		},
		{
			question: "In R with the usual metric, the open ball of radius 4 centred at 0 is the interval (a, b). What is b? (enter a number)",
			answer:   "4", kind: "endpoint", center: 0, radius: 4, side: "upper",
			reason: "The ball is (-4, 4), so the upper endpoint is 4.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.kind == "distance" {
		facts["x1"] = strconv.Itoa(c.x1)
		facts["y1"] = strconv.Itoa(c.y1)
		facts["x2"] = strconv.Itoa(c.x2)
		facts["y2"] = strconv.Itoa(c.y2)
	} else {
		facts["center"] = strconv.Itoa(c.center)
		facts["radius"] = strconv.Itoa(c.radius)
		facts["side"] = c.side
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// 2. Interior, closure and boundary (numeric)
// ---------------------------------------------------------------------------

type interiorGen struct{}

func (g *interiorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		points   []int
		kind     string // boundary | interior
		reason   string
	}
	claims := []claim{
		{
			question: "In R with the usual topology, how many boundary points does the set [0, 2] have? (enter a number)",
			answer:   "2", points: []int{0, 2}, kind: "boundary",
			reason: "The boundary is {0, 2}: every neighbourhood of either endpoint meets both the set and its complement.",
		},
		{
			question: "In R with the usual topology, how many boundary points does the set (0, 1) union (1, 2) have? (enter a number)",
			answer:   "3", points: []int{0, 1, 2}, kind: "boundary",
			reason: "The boundary is {0, 1, 2}: the endpoints and the missing point 1 all have neighbourhoods meeting both sides.",
		},
		{
			question: "In R with the usual topology, how many interior points does the set {1, 2, 3} have? (enter a number)",
			answer:   "0", points: []int{}, kind: "interior",
			reason: "A finite set of isolated points has empty interior: no ball around 1, 2 or 3 stays inside the set.",
		},
		{
			question: "In R with the usual topology, how many boundary points does the set [0, 1] union [2, 3] have? (enter a number)",
			answer:   "4", points: []int{0, 1, 2, 3}, kind: "boundary",
			reason: "The boundary is {0, 1, 2, 3}: the four endpoints each have neighbourhoods meeting the set and its complement.",
		},
		{
			question: "In R with the usual topology, how many boundary points does the set (0, infinity) have? (enter a number)",
			answer:   "1", points: []int{0}, kind: "boundary",
			reason: "The only boundary point is 0; every other real is either inside the ray or has a ball outside it.",
		},
		{
			question: "In R with the usual topology, how many boundary points does the set (0, 1) have? (enter a number)",
			answer:   "2", points: []int{0, 1}, kind: "boundary",
			reason: "The boundary is {0, 1}: both endpoints are limits from inside and outside the open interval.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	parts := make([]string, len(c.points))
	for i, p := range c.points {
		parts[i] = strconv.Itoa(p)
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"points": strings.Join(parts, ","), "kind": c.kind, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 3. Limit points
// ---------------------------------------------------------------------------

type limitPointGen struct{}

func (g *limitPointGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "In R with the usual topology, which point is a limit point of the set {1/n : n at least 1}?",
			correct:     "0",
			distractors: []string{"1", "1/2", "2"},
			reason:      "The sequence 1/n converges to 0, and every neighbourhood of 0 contains some 1/n, so 0 is a limit point.",
		},
		{
			question:    "Is every point of the closed interval [0, 1] a limit point of [0, 1]? (yes/no)",
			correct:     "yes",
			distractors: []string{"no", "only the endpoints", "only the interior points"},
			reason:      "Every neighbourhood of any point of [0, 1] meets the interval in infinitely many other points, endpoints included.",
		},
		{
			question:    "Which subset of R has no limit points?",
			correct:     "the integers",
			distractors: []string{"the rationals", "the open interval (0, 1)", "the set {1/n : n at least 1}"},
			reason:      "Each integer has a ball of radius 1/2 containing no other integer, so no point is a limit of the integers.",
		},
		{
			question:    "A point x is a limit point of a set A exactly when...",
			correct:     "every neighbourhood of x meets A at a point other than x",
			distractors: []string{"x belongs to A", "A contains a neighbourhood of x", "x is an endpoint of A"},
			reason:      "The definition is about arbitrary closeness, not membership: the point x itself may or may not lie in A.",
		},
		{
			question:    "In the subspace {0} union {1/n : n at least 1} of R, how many limit points does the subspace have?",
			correct:     "1",
			distractors: []string{"0", "2", "infinitely many"},
			reason:      "The point 0 is the only limit point; each 1/n is isolated because a small ball around it misses the rest of the set.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 4. Convergence (numeric)
// ---------------------------------------------------------------------------

type convergenceGen struct{}

func (g *convergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		a, b, c  int
		d        int
		reason   string
	}
	claims := []claim{
		{
			question: "The sequence a_n = (1)/(n) converges to what limit? (enter a number)",
			answer:   "0", a: 0, b: 1, c: 1, d: 0,
			reason: "As n grows, 1/n tends to 0: every ball around 0 eventually contains the whole tail of the sequence.",
		},
		{
			question: "The sequence a_n = (3n + 1)/(n) converges to what limit? (enter a number)",
			answer:   "3", a: 3, b: 1, c: 1, d: 0,
			reason: "Split the fraction as 3 + 1/n; the correction 1/n tends to 0, leaving the limit 3.",
		},
		{
			question: "The sequence a_n = (2n + 1)/(n) converges to what limit? (enter a number)",
			answer:   "2", a: 2, b: 1, c: 1, d: 0,
			reason: "Split the fraction as 2 + 1/n; the correction vanishes, so the limit is 2.",
		},
		{
			question: "The sequence a_n = (5n - 1)/(n) converges to what limit? (enter a number)",
			answer:   "5", a: 5, b: -1, c: 1, d: 0,
			reason: "Split the fraction as 5 - 1/n; the correction vanishes, so the limit is 5.",
		},
		{
			question: "The sequence a_n = (n)/(n + 1) converges to what limit? (enter a number)",
			answer:   "1", a: 1, b: 0, c: 1, d: 1,
			reason: "Divide numerator and denominator by n: the fraction becomes 1/(1 + 1/n), which tends to 1.",
		},
		{
			question: "The sequence a_n = (7n)/(2n + 1) converges to what limit? (enter a decimal)",
			answer:   "3.5", a: 7, b: 0, c: 2, d: 1,
			reason: "Divide by n: the fraction becomes 7/(2 + 1/n), which tends to 7/2 = 3.5.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b),
			"c": strconv.Itoa(c.c), "d": strconv.Itoa(c.d), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Completeness and Cauchy sequences
// ---------------------------------------------------------------------------

type completeGen struct{}

func (g *completeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Which metric space is complete?",
			correct:     "the real line with the usual metric",
			distractors: []string{"the rationals with the usual metric", "the open interval (0, 1) with the usual metric", "the rationals in [0, 1] with the usual metric"},
			reason:      "Every Cauchy sequence of reals converges to a real; the other spaces lose limits such as sqrt(2) or 0.",
		},
		{
			question:    "A Cauchy sequence in a complete metric space...",
			correct:     "converges to a point of the space",
			distractors: []string{"has a convergent subsequence only", "must be eventually constant", "must diverge"},
			reason:      "Completeness is exactly the statement that every Cauchy sequence has a limit inside the space.",
		},
		{
			question:    "The completion of the rationals with the usual metric is...",
			correct:     "the real line",
			distractors: []string{"the rationals again", "the complex plane", "the empty set"},
			reason:      "Filling the missing limits of Cauchy sequences of rationals produces the reals.",
		},
		{
			question:    "Every compact metric space is...",
			correct:     "complete",
			distractors: []string{"discrete", "not Hausdorff", "unbounded"},
			reason:      "A Cauchy sequence in a compact space has a convergent subsequence, and the whole Cauchy sequence then converges to the same limit.",
		},
		{
			question:    "Which sequence in (0, 1) converges to 0, showing that (0, 1) is not complete?",
			correct:     "1/n",
			distractors: []string{"1/2, 1/2, 1/2, ...", "1/2 + 1/n", "1 - 1/n"},
			reason:      "The sequence 1/n is Cauchy but its only limit is 0, which is not in (0, 1); the other options converge inside or to 1.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 6. Bases and subbases
// ---------------------------------------------------------------------------

type basisGen struct{}

func (g *basisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []topoClaim{
		{
			question:    "A basis for a topology is a collection of open sets such that...",
			correct:     "every open set is a union of basis elements",
			distractors: []string{"every basis element is compact", "the basis is closed under complements", "the basis has finitely many sets"},
			reason:      "The basis generates the topology: each open set is assembled from basis elements, and the basis axioms make that assembly consistent.",
		},
		{
			question:    "Which collection is a basis for the usual topology on R?",
			correct:     "the open intervals (a, b)",
			distractors: []string{"the closed intervals [a, b]", "the singletons {x}", "the half-open intervals [a, b)"},
			reason:      "Open intervals generate the usual opens; closed intervals generate the discrete topology and half-open intervals generate the Sorgenfrey topology.",
		},
		{
			question:    "In a metric space, the collection of open balls is a basis for...",
			correct:     "the metric topology",
			distractors: []string{"the discrete topology", "the cofinite topology", "the indiscrete topology"},
			reason:      "Open sets in a metric space are unions of balls, which is the defining property of a basis for the metric topology.",
		},
		{
			question:    "A subbasis differs from a basis because...",
			correct:     "finite intersections of subbasis elements form a basis",
			distractors: []string{"it must be countable", "it contains the empty set", "it is closed under arbitrary unions"},
			reason:      "A subbasis is the raw generating family; taking finite intersections produces a basis, and unions of that basis produce the topology.",
		},
	}
	claims = append(claims, computedTopoClaims("topo.metric.basis")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 7. The subspace topology
// ---------------------------------------------------------------------------

type subspaceGen struct{}

func (g *subspaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []topoClaim{
		{
			question:    "For a subset A of a space X, a set U is open in the subspace topology exactly when...",
			correct:     "U is the intersection of A with an open set of X",
			distractors: []string{"U is open in X", "U is closed in X", "U contains every limit point of A"},
			reason:      "The subspace opens are the traces A intersect V of the ambient opens V, by definition.",
		},
		{
			question:    "In the subspace [0, 1] of R, the set [0, 1/2) is...",
			correct:     "open in the subspace but not in R",
			distractors: []string{"open in R", "closed in the subspace", "neither open nor closed"},
			reason:      "It is [0, 1] intersect (-1, 1/2), so it is a subspace open; in R it is not open because no ball around 0 stays inside it.",
		},
		{
			question:    "The subspace topology on A is the...",
			correct:     "coarsest topology making the inclusion of A into X continuous",
			distractors: []string{"finest topology making the inclusion continuous", "quotient topology of A", "product topology on A"},
			reason:      "The inclusion is continuous exactly when each subspace open is a trace of an ambient open, which is the coarsest such choice.",
		},
		{
			question:    "A subspace of a Hausdorff space is...",
			correct:     "Hausdorff",
			distractors: []string{"not necessarily Hausdorff", "always compact", "always connected"},
			reason:      "Separating two points by disjoint ambient opens and intersecting both with the subspace gives disjoint subspace opens.",
		},
	}
	claims = append(claims, computedTopoClaims("topo.metric.subspace")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 8. Dense sets and separability
// ---------------------------------------------------------------------------

type denseGen struct{}

func (g *denseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []topoClaim{
		{
			question:    "A subset D of a space X is dense exactly when...",
			correct:     "every nonempty open set meets D",
			distractors: []string{"D is closed", "D is countable", "D contains every limit point"},
			reason:      "D is dense when its closure is X, which is the same as every nonempty open set meeting D.",
		},
		{
			question:    "Which subset is dense in R with the usual topology?",
			correct:     "the rationals",
			distractors: []string{"the integers", "the integers between 0 and 1", "the set {0}"},
			reason:      "Every open interval contains a rational, so the rationals meet every nonempty open set; the other sets leave gaps.",
		},
		{
			question:    "A separable space is one with...",
			correct:     "a countable dense subset",
			distractors: []string{"a countable basis", "finitely many points", "no limit points"},
			reason:      "Separability is defined by the existence of a countable dense subset; a countable basis is a different, stronger condition.",
		},
		{
			question:    "The closure of a dense set is...",
			correct:     "the whole space",
			distractors: []string{"the empty set", "the set itself", "a proper subset of the space"},
			reason:      "Dense means closure equal to the whole space; the set itself need not be closed.",
		},
	}
	claims = append(claims, computedTopoClaims("topo.metric.dense")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 9. Homeomorphisms
// ---------------------------------------------------------------------------

type homeomorphismGen struct{}

func (g *homeomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A homeomorphism is a bijection that...",
			correct:     "is continuous with a continuous inverse",
			distractors: []string{"is continuous", "preserves distances", "is linear"},
			reason:      "Topological sameness needs both directions continuous; a continuous bijection alone need not have a continuous inverse.",
		},
		{
			question:    "Which property is a topological invariant?",
			correct:     "connectedness",
			distractors: []string{"boundedness", "the metric", "the diameter"},
			reason:      "A homeomorphism carries connected sets to connected sets, while boundedness and diameter depend on the chosen metric.",
		},
		{
			question:    "Are the open interval (0, 1) and the real line homeomorphic? (yes/no)",
			correct:     "yes",
			distractors: []string{"no", "only if both are compact", "only under the discrete metric"},
			reason:      "The map x to tan(pi(x - 1/2)) is a continuous bijection with a continuous inverse from (0, 1) to R.",
		},
		{
			question:    "A continuous bijection is a homeomorphism when...",
			correct:     "its inverse is continuous",
			distractors: []string{"the spaces are compact", "the spaces are metric", "it preserves distances"},
			reason:      "Continuity of the inverse is the missing direction; compactness or a metric can supply it under extra hypotheses, but not by themselves.",
		},
		{
			question:    "Which map is a homeomorphism from (0, 1) to R?",
			correct:     "x to tan(pi(x - 1/2))",
			distractors: []string{"x to x", "x to x^2", "x to sin(x)"},
			reason:      "The tangent map sends the open interval onto the whole line continuously with a continuous inverse; the others are not onto R.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 10. Euler characteristic (numeric)
// ---------------------------------------------------------------------------

type eulerGen struct{}

func (g *eulerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		v, e, f  int
		reason   string
	}
	claims := []claim{
		{
			question: "A tetrahedron has 4 vertices, 6 edges and 4 faces. What is its Euler characteristic V - E + F? (enter a number)",
			answer:   "2", v: 4, e: 6, f: 4,
			reason: "4 - 6 + 4 = 2, the Euler characteristic of a sphere.",
		},
		{
			question: "A cube has 8 vertices, 12 edges and 6 faces. What is its Euler characteristic V - E + F? (enter a number)",
			answer:   "2", v: 8, e: 12, f: 6,
			reason: "8 - 12 + 6 = 2, as for every polyhedron homeomorphic to a sphere.",
		},
		{
			question: "A square pyramid has 5 vertices, 8 edges and 5 faces. What is its Euler characteristic V - E + F? (enter a number)",
			answer:   "2", v: 5, e: 8, f: 5,
			reason: "5 - 8 + 5 = 2; the pyramid is also a sphere topologically.",
		},
		{
			question: "A triangulated torus has 9 vertices, 27 edges and 18 faces. What is its Euler characteristic V - E + F? (enter a number)",
			answer:   "0", v: 9, e: 27, f: 18,
			reason: "9 - 27 + 18 = 0, the Euler characteristic of a torus.",
		},
		{
			question: "A triangulated projective plane has 6 vertices, 15 edges and 10 faces. What is its Euler characteristic V - E + F? (enter a number)",
			answer:   "1", v: 6, e: 15, f: 10,
			reason: "6 - 15 + 10 = 1, the Euler characteristic of the projective plane.",
		},
		{
			question: "An octahedron has 6 vertices, 12 edges and 8 faces. What is its Euler characteristic V - E + F? (enter a number)",
			answer:   "2", v: 6, e: 12, f: 8,
			reason: "6 - 12 + 8 = 2, again the sphere value.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"v": strconv.Itoa(c.v), "e": strconv.Itoa(c.e), "f": strconv.Itoa(c.f),
			"answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. Genus of a surface (numeric)
// ---------------------------------------------------------------------------

type genusGen struct{}

func (g *genusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string // genus | chi
		chi      int
		genus    int
		reason   string
	}
	claims := []claim{
		{
			question: "A closed orientable surface has Euler characteristic 0. What is its genus? (enter a number)",
			answer:   "1", kind: "genus", chi: 0,
			reason: "For a closed orientable surface chi = 2 - 2g, so g = (2 - 0)/2 = 1: the torus.",
		},
		{
			question: "A closed orientable surface has Euler characteristic -2. What is its genus? (enter a number)",
			answer:   "2", kind: "genus", chi: -2,
			reason: "Solving 2 - 2g = -2 gives g = 2, the double torus.",
		},
		{
			question: "A closed orientable surface has Euler characteristic 2. What is its genus? (enter a number)",
			answer:   "0", kind: "genus", chi: 2,
			reason: "Solving 2 - 2g = 2 gives g = 0: the sphere.",
		},
		{
			question: "A closed orientable surface has Euler characteristic -4. What is its genus? (enter a number)",
			answer:   "3", kind: "genus", chi: -4,
			reason: "Solving 2 - 2g = -4 gives g = 3, the triple torus.",
		},
		{
			question: "A closed orientable surface has genus 3. What is its Euler characteristic 2 - 2g? (enter a number)",
			answer:   "-4", kind: "chi", genus: 3,
			reason: "2 - 2(3) = -4.",
		},
		{
			question: "A closed orientable surface has genus 4. What is its Euler characteristic 2 - 2g? (enter a number)",
			answer:   "-6", kind: "chi", genus: 4,
			reason: "2 - 2(4) = -6.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.kind == "genus" {
		facts["chi"] = strconv.Itoa(c.chi)
	} else {
		facts["genus"] = strconv.Itoa(c.genus)
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// 12. Betti numbers (numeric)
// ---------------------------------------------------------------------------

type bettiGen struct{}

func (g *bettiGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		param    int
		reason   string
	}
	claims := []claim{
		{
			question: "The first Betti number of a circle is what? (enter a number)",
			answer:   "1", kind: "circle",
			reason: "A circle has one independent loop, so H_1 has rank 1.",
		},
		{
			question: "A closed orientable surface of genus 2 has first Betti number b1 = 2g. What is b1? (enter a number)",
			answer:   "4", kind: "surface", param: 2,
			reason: "b1 = 2g = 2(2) = 4 for the double torus.",
		},
		{
			question: "A wedge of 3 circles has first Betti number equal to the number of circles. What is b1? (enter a number)",
			answer:   "3", kind: "wedge", param: 3,
			reason: "Each circle contributes one independent loop, so a wedge of 3 circles has b1 = 3.",
		},
		{
			question: "A space has 4 connected components. What is its zeroth Betti number b0? (enter a number)",
			answer:   "4", kind: "components", param: 4,
			reason: "b0 counts connected components, so four components give b0 = 4.",
		},
		{
			question: "A tree is a connected graph with no cycles. What is its first Betti number b1? (enter a number)",
			answer:   "0", kind: "tree",
			reason: "A tree has no independent loop, so H_1 has rank 0.",
		},
		{
			question: "For a sphere, b0 = 1, b1 = 0 and b2 = 1. What is b0 - b1 + b2? (enter a number)",
			answer:   "2", kind: "sphere_euler",
			reason: "1 - 0 + 1 = 2, which is also the Euler characteristic of the sphere.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"kind": c.kind, "param": strconv.Itoa(c.param), "answer": c.answer,
		},
	}
}
