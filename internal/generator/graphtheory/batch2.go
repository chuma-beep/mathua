package graphtheory

// The second graph-theory batch: twenty concepts that take the domain past its
// introductory material into planarity and minors, spectral graph theory,
// extremal and random graphs, and flows and digraphs.
//
// Ten concepts are multiple choice and ten are numeric. Every numeric generator
// publishes the graph parameters its answer follows from, so the tests recompute
// the answer from those facts rather than trusting the author's arithmetic.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Kuratowski's theorem (multiple choice)
// ---------------------------------------------------------------------------

type kuratowskiGen struct{}

func (g *kuratowskiGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return kuratowskiProblem()
}

// ---------------------------------------------------------------------------
// Graph minors (multiple choice)
// ---------------------------------------------------------------------------

type minorGen struct{}

func (g *minorGen) Generate(ctx generator.GeneratorContext) generator.Problem { return minorProblem() }

// ---------------------------------------------------------------------------
// Treewidth (numeric)
// ---------------------------------------------------------------------------

type treewidthGen struct{}

// Generate draws the graph whose treewidth is asked for: trees, paths and cycles
// at the first ordinal, the worked cycles at the second, and grids, complete
// bipartite graphs and complete graphs at the third, where a bounded-width
// dynamic program is the subject.
func (g *treewidthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	treeClaim := func() {
		n := graphRange(4, 15)
		kind, ans = "tree", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the treewidth of a tree on %d vertices? (enter a number)", n)
		reason = fmt.Sprintf("Bags that each hold one edge, chained along the tree, form a tree decomposition of width 1, and the %d edges of the tree force width at least 1.", n-1)
	}
	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the treewidth of the path P_%d? (enter a number)", n)
		reason = fmt.Sprintf("A path is a tree on %d vertices, so chaining one bag per edge gives width 1, and a graph with an edge never has width 0.", n)
	}
	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", 2
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the treewidth of the cycle C_%d? (enter a number)", n)
		reason = fmt.Sprintf("The bags {1, 2, 3}, {1, 3, 4} up to {1, %d - 1, %d} hold every edge of C_%d, chain along a path, and keep the bags containing each vertex consecutive, so the width is 2; C_%d is not a forest, so the width cannot be 1.", n, n, n, n)
	}
	completeClaim := func() {
		n := graphRange(3, 12)
		kind, ans = "complete", n-1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the treewidth of the complete graph K_%d? (enter a number)", n)
		reason = fmt.Sprintf("Every pair of the %d vertices of K_%d is adjacent, so one bag must hold all of them and the width is %d - 1 = %d.", n, n, n, ans)
	}
	gridClaim := func() {
		r := graphRange(3, 8)
		c := graphRange(3, 8)
		m := graphMin(r, c)
		kind, ans = "grid", m
		facts["rows"], facts["cols"] = strconv.Itoa(r), strconv.Itoa(c)
		question = fmt.Sprintf("What is the treewidth of the %d-by-%d grid graph? (enter a number)", r, c)
		reason = fmt.Sprintf("An r-by-c grid has treewidth min(r, c), and the bags down one column form a decomposition attaining it, so the treewidth of the %d-by-%d grid is %d.", r, c, ans)
	}
	bipartiteClaim := func() {
		a := graphRange(2, 6)
		b := graphRange(2, 7)
		m := graphMin(a, b)
		kind, ans = "bipartite", m
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question = fmt.Sprintf("What is the treewidth of the complete bipartite graph K_{%d,%d}? (enter a number)", a, b)
		reason = fmt.Sprintf("K_{a,b} has treewidth min(a, b) for positive a and b, and the smaller side of K_{%d,%d} has %d vertices, so the treewidth is %d.", a, b, m, ans)
	}

	switch o {
	case 1:
		switch rand.Intn(3) {
		case 0:
			treeClaim()
		case 1:
			pathClaim()
		default:
			cycleClaim()
		}
	case 2:
		if rand.Intn(3) == 0 {
			pathClaim()
		} else {
			cycleClaim()
		}
	default:
		switch rand.Intn(3) {
		case 0:
			gridClaim()
		case 1:
			bipartiteClaim()
		default:
			completeClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// Graph genus (numeric)
// ---------------------------------------------------------------------------

type genusGen struct{}

// Generate draws the graph whose genus is asked for: cycles and planar graphs at
// the first ordinal, the complete graphs that leave the sphere at the second,
// and the complete and complete bipartite graphs of the genus formulas at the
// third.
func (g *genusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", 0
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the genus of the cycle C_%d? (enter a number)", n)
		reason = fmt.Sprintf("C_%d has a drawing with no crossings in the plane, and the plane is the sphere of genus 0, so its genus is 0.", n)
	}
	planarClaim := func() {
		v := graphRange(4, 15)
		e := v - 1 + rand.Intn(3*v-6-(v-1)+1)
		kind, ans = "planar", 0
		facts["n"] = strconv.Itoa(v)
		facts["m"] = strconv.Itoa(e)
		question = fmt.Sprintf("A connected planar graph has %d vertices and %d edges. What is its genus? (enter a number)", v, e)
		reason = fmt.Sprintf("The graph is planar, so it embeds in the sphere, whose genus is 0. Its %d edges are within the planar bound 3V - 6 = %d for %d vertices, as a planar graph must be.", e, 3*v-6, v)
	}
	completeClaim := func() {
		n := graphRange(3, 14)
		kind, ans = "complete", graphCeil((n-3)*(n-4), 12)
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the genus of the complete graph K_%d? (enter a number)", n)
		reason = fmt.Sprintf("Ringel's formula gives genus(K_n) = ceil((n-3)(n-4)/12), so K_%d has genus ceil(%d * %d / 12) = %d.", n, n-3, n-4, ans)
	}
	bipartiteClaim := func() {
		a := graphRange(2, 7)
		b := graphRange(2, 8)
		kind, ans = "bipartite", graphCeil((a-2)*(b-2), 4)
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question = fmt.Sprintf("What is the genus of the complete bipartite graph K_{%d,%d}? (enter a number)", a, b)
		reason = fmt.Sprintf("Ringel's formula for complete bipartite graphs gives ceil((a-2)(b-2)/4), so K_{%d,%d} has genus ceil(%d * %d / 4) = %d.", a, b, a-2, b-2, ans)
	}

	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			cycleClaim()
		} else {
			planarClaim()
		}
	case 2:
		if rand.Intn(2) == 0 {
			completeClaim()
		} else {
			cycleClaim()
		}
	default:
		if rand.Intn(2) == 0 {
			completeClaim()
		} else {
			bipartiteClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// graphCeil divides a non-negative numerator by a positive denominator and
// rounds up, which is what the genus formulas ask for.
func graphCeil(a, b int) int {
	if a <= 0 {
		return 0
	}
	return (a + b - 1) / b
}

// ---------------------------------------------------------------------------
// Outerplanar graphs (numeric)
// ---------------------------------------------------------------------------

type outerplanarGen struct{}

// Generate draws the polygon whose triangulation bounds the edge count. The
// three ordinals ask the same arithmetic through the three routes the section
// takes: the edge bound itself, the polygon-and-diagonals count and the maximal
// outerplanar graph.
func (g *outerplanarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	n := graphRange(3, 15)
	ans := 2*n - 3
	facts := map[string]string{"n": strconv.Itoa(n)}

	var question, reason string
	switch o {
	case 1:
		question = fmt.Sprintf("A simple graph is outerplanar when it can be drawn in the plane with every vertex on the outer face. What is the maximum number of edges of an outerplanar graph on %d vertices? (enter a number)", n)
		reason = fmt.Sprintf("A maximal outerplanar graph on n vertices has 2n - 3 edges, so with n = %d the maximum is 2 * %d - 3 = %d.", n, n, ans)
	case 2:
		question = fmt.Sprintf("A triangulated polygon has %d vertices. It uses its %d sides and then adds non-crossing diagonals until no more fit. How many edges does the triangulation have? (enter a number)", n, n)
		reason = fmt.Sprintf("A triangulation of an n-gon has n sides and n - 3 diagonals, so with n = %d there are %d + %d = %d edges.", n, n, n-3, ans)
	default:
		question = fmt.Sprintf("A maximal outerplanar graph on %d vertices is a triangulated polygon. How many edges does it have? (enter a number)", n)
		reason = fmt.Sprintf("Adding any edge to a maximal outerplanar graph would break outerplanarity, and the triangulated polygon on %d vertices has 2 * %d - 3 = %d edges.", n, n, ans)
	}

	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// Adjacency matrix (numeric)
// ---------------------------------------------------------------------------

type adjacencyGen struct{}

func (g *adjacencyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	walkClaim := func(kindName, symbol string, n, length int) {
		from := rand.Intn(n)
		to := rand.Intn(n)
		kind = kindName
		ans = graphWalks(kindName, n, from, to, length)
		facts["n"] = strconv.Itoa(n)
		facts["i"] = strconv.Itoa(from + 1)
		facts["j"] = strconv.Itoa(to + 1)
		facts["length"] = strconv.Itoa(length)
		label := graphGraphName(kindName, n)
		question = fmt.Sprintf("In %s, how many walks of length %d go from vertex %d to vertex %d? (enter a number)", label, length, from+1, to+1)
		if length == 1 {
			reason = fmt.Sprintf("A walk of length 1 is a single edge, so the count is the (%d,%d) entry of the adjacency matrix of %s_%d, which is %d.", from+1, to+1, symbol, n, ans)
		} else {
			reason = fmt.Sprintf("The (%d,%d) entry of A^%d counts the walks of length %d from vertex %d to vertex %d in %s_%d, and adding the counts one step at a time gives %d.", from+1, to+1, length, length, from+1, to+1, symbol, n, ans)
		}
	}

	// A closed walk of length two is the row sum of the adjacency matrix, so the
	// first ordinal asks for the matrix and its row sums through one count.
	degreeClaim := func(kindName, symbol string, n int) {
		v := rand.Intn(n)
		kind = kindName
		ans = graphWalks(kindName, n, v, v, 2)
		facts["n"] = strconv.Itoa(n)
		facts["i"] = strconv.Itoa(v + 1)
		facts["j"] = strconv.Itoa(v + 1)
		facts["length"] = "2"
		question = fmt.Sprintf("In %s, how many walks of length 2 start and end at vertex %d? (enter a number)", graphGraphName(kindName, n), v+1)
		reason = fmt.Sprintf("Each neighbour of vertex %d gives one walk out and back, so the count is its degree, which is also the (%d,%d) entry of A^2 and the sum of row %d of A: %d.", v+1, v+1, v+1, v+1, ans)
	}

	switch o {
	case 1:
		switch rand.Intn(3) {
		case 0:
			degreeClaim("path", "P", graphRange(3, 12))
		case 1:
			degreeClaim("cycle", "C", graphRange(3, 12))
		default:
			degreeClaim("complete", "K", graphRange(3, 12))
		}
	case 2:
		switch rand.Intn(3) {
		case 0:
			walkClaim("path", "P", graphRange(3, 12), 2)
		case 1:
			walkClaim("cycle", "C", graphRange(3, 12), 2)
		default:
			walkClaim("complete", "K", graphRange(3, 12), 2)
		}
	default:
		if rand.Intn(2) == 0 {
			walkClaim("cycle", "C", graphRange(3, 12), 3)
		} else {
			walkClaim("complete", "K", graphRange(3, 10), 2)
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// graphGraphName renders a named graph family for a question stem.
func graphGraphName(kindName string, n int) string {
	switch kindName {
	case "path":
		return fmt.Sprintf("the path P_%d", n)
	case "cycle":
		return fmt.Sprintf("the cycle C_%d", n)
	default:
		return fmt.Sprintf("the complete graph K_%d", n)
	}
}

// ---------------------------------------------------------------------------
// Laplacian matrix (numeric)
// ---------------------------------------------------------------------------

type laplacianGen struct{}

// Generate draws the graph whose Laplacian is inspected. The zero eigenvalue
// count is the number of connected components, so a connected graph answers 1
// and a disjoint union answers the number of pieces. The first ordinal asks for
// a connected graph, the second for a disjoint union, and the third for the
// connected families the later section uses.
func (g *laplacianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many eigenvalues of the Laplacian of the path P_%d equal zero? (enter a number)", n)
		reason = fmt.Sprintf("P_%d is connected, and a connected graph has exactly one zero Laplacian eigenvalue, so the multiplicity is 1.", n)
	}
	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many eigenvalues of the Laplacian of the cycle C_%d equal zero? (enter a number)", n)
		reason = fmt.Sprintf("C_%d is connected, so the kernel of its Laplacian is spanned by the all-ones vector and the zero eigenvalue has multiplicity 1.", n)
	}
	treeClaim := func() {
		n := graphRange(4, 15)
		kind, ans = "tree", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many eigenvalues of the Laplacian of a tree on %d vertices equal zero? (enter a number)", n)
		reason = fmt.Sprintf("Every tree is connected, and on %d vertices the Laplacian has exactly one zero eigenvalue.", n)
	}
	completeClaim := func() {
		n := graphRange(3, 12)
		kind, ans = "complete", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many eigenvalues of the Laplacian of the complete graph K_%d equal zero? (enter a number)", n)
		reason = fmt.Sprintf("K_%d is connected with Laplacian eigenvalues 0 and %d repeated %d times, so exactly one eigenvalue is zero.", n, n, n-1)
	}
	disjointClaim := func() {
		k := graphRange(2, 7)
		m := graphRange(3, 10)
		kind, ans = "disjoint", k
		facts["components"] = strconv.Itoa(k)
		facts["m"] = strconv.Itoa(m)
		question = fmt.Sprintf("How many eigenvalues of the Laplacian of %d disjoint copies of the cycle C_%d equal zero? (enter a number)", k, m)
		reason = fmt.Sprintf("Each connected component contributes exactly one zero eigenvalue, and %d disjoint cycles have %d components, so the multiplicity is %d.", k, k, ans)
	}

	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			pathClaim()
		} else {
			completeClaim()
		}
	case 2:
		disjointClaim()
	default:
		switch rand.Intn(3) {
		case 0:
			treeClaim()
		case 1:
			cycleClaim()
		default:
			completeClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// Matrix-tree theorem (numeric)
// ---------------------------------------------------------------------------

type matrixTreeGen struct{}

// Generate draws the graph whose spanning-tree count the matrix-tree theorem
// supplies: cycles and paths at the first ordinal, complete graphs under
// Cayley's formula at the second, and complete and complete bipartite graphs at
// the third.
func (g *matrixTreeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", n
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many spanning trees does the cycle C_%d have? (enter a number)", n)
		reason = fmt.Sprintf("Every spanning tree of C_%d is the cycle with one edge deleted, so the %d edges give %d spanning trees.", n, n, ans)
	}
	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", 1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many spanning trees does the path P_%d have? (enter a number)", n)
		reason = fmt.Sprintf("P_%d is already a tree and spans itself, so it has exactly 1 spanning tree.", n)
	}
	completeClaim := func() {
		n := graphRange(3, 9)
		kind, ans = "complete", graphIntPow(n, n-2)
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("How many spanning trees does the complete graph K_%d have? (enter a number)", n)
		reason = fmt.Sprintf("Cayley's formula gives n^(n-2) spanning trees, so K_%d has %d^%d = %d.", n, n, n-2, ans)
	}
	bipartiteClaim := func() {
		a := graphRange(2, 6)
		b := graphRange(2, 7)
		kind, ans = "bipartite", graphIntPow(a, b-1)*graphIntPow(b, a-1)
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question = fmt.Sprintf("How many spanning trees does the complete bipartite graph K_{%d,%d} have? (enter a number)", a, b)
		reason = fmt.Sprintf("K_{a,b} has a^(b-1) b^(a-1) spanning trees, so K_{%d,%d} has %d^%d * %d^%d = %d.", a, b, a, b-1, b, a-1, ans)
	}

	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			cycleClaim()
		} else {
			pathClaim()
		}
	case 2:
		if rand.Intn(2) == 0 {
			completeClaim()
		} else {
			cycleClaim()
		}
	default:
		if rand.Intn(2) == 0 {
			completeClaim()
		} else {
			bipartiteClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// Spectral radius (numeric)
// ---------------------------------------------------------------------------

type spectralRadiusGen struct{}

// Generate draws the graph whose spectral radius is asked for. The first ordinal
// takes the bound through complete graphs, the second through the regular
// families the section works, and the third through the Perron-Frobenius bounds
// on regular graphs of a drawn degree.
func (g *spectralRadiusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	completeClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "complete", n-1
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the spectral radius of the complete graph K_%d? (enter a number)", n)
		reason = fmt.Sprintf("K_%d is %d-regular, and the all-ones vector is an eigenvector of every regular graph with eigenvalue its degree, so the spectral radius is %d.", n, n-1, ans)
	}
	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", 2
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the spectral radius of the cycle C_%d? (enter a number)", n)
		reason = fmt.Sprintf("C_%d is 2-regular, so the all-ones vector is an eigenvector with eigenvalue 2 and the spectral radius is 2.", n)
	}
	bipartiteClaim := func() {
		d := graphRange(2, 13)
		kind, ans = "regular", d
		facts["d"] = strconv.Itoa(d)
		question = fmt.Sprintf("What is the spectral radius of the complete bipartite graph K_{%d,%d}? (enter a number)", d, d)
		reason = fmt.Sprintf("K_{%d,%d} is %d-regular, so the all-ones vector is an eigenvector with eigenvalue %d and the spectral radius is %d.", d, d, d, d, ans)
	}
	hypercubeClaim := func() {
		d := graphRange(2, 12)
		kind, ans = "regular", d
		facts["d"] = strconv.Itoa(d)
		question = fmt.Sprintf("What is the spectral radius of the %d-dimensional hypercube Q_%d? (enter a number)", d, d)
		reason = fmt.Sprintf("Q_%d is %d-regular: every vertex flips one of its %d coordinates, so it has %d neighbours and the spectral radius is %d.", d, d, d, d, ans)
	}

	switch o {
	case 1:
		completeClaim()
	case 2:
		if rand.Intn(2) == 0 {
			bipartiteClaim()
		} else {
			cycleClaim()
		}
	default:
		if rand.Intn(2) == 0 {
			hypercubeClaim()
		} else {
			bipartiteClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// Algebraic connectivity (numeric)
// ---------------------------------------------------------------------------

type spectralGapGen struct{}

// Generate draws the graph whose algebraic connectivity is asked for: complete
// graphs at the first ordinal, the cycles whose cosine formula is exact at
// lengths 4 and 6 together with the stars that bottleneck at 1 at the second,
// and complete bipartite graphs at the third.
//
// The cycle question is restricted to lengths 4 and 6 on purpose. For other
// lengths 2 - 2cos(2pi/n) is not an integer, and an integer answer would be a
// rounding of the truth rather than the truth.
func (g *spectralGapGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	completeClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "complete", n
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the algebraic connectivity of the complete graph K_%d? (enter a number)", n)
		reason = fmt.Sprintf("The Laplacian eigenvalues of K_n are 0 and n repeated n - 1 times, so the second-smallest eigenvalue of K_%d is %d.", n, ans)
	}
	cycleClaim := func() {
		n := 4
		if rand.Intn(2) == 0 {
			n = 6
		}
		kind, ans = "cycle", 1
		if n == 4 {
			ans = 2
		}
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the algebraic connectivity of the cycle C_%d? (enter a number)", n)
		reason = fmt.Sprintf("For C_n the second-smallest Laplacian eigenvalue is 2 - 2cos(2pi/n), and at n = %d it is %d.", n, ans)
	}
	starClaim := func() {
		k := graphRange(3, 14)
		kind, ans = "star", 1
		facts["n"] = strconv.Itoa(k)
		question = fmt.Sprintf("What is the algebraic connectivity of the star K_{1,%d}? (enter a number)", k)
		reason = fmt.Sprintf("The Laplacian eigenvalues of K_{1,%d} are 0, 1 repeated %d times, and %d, so the second-smallest is 1.", k, k-1, k+1)
	}
	bipartiteClaim := func() {
		a := graphRange(2, 7)
		b := graphRange(2, 8)
		m := graphMin(a, b)
		kind, ans = "bipartite", m
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question = fmt.Sprintf("What is the algebraic connectivity of the complete bipartite graph K_{%d,%d}? (enter a number)", a, b)
		reason = fmt.Sprintf("The Laplacian eigenvalues of K_{a,b} are 0, a, b and a + b, so the second-smallest is min(a, b) = %d here.", m)
	}

	switch o {
	case 1:
		completeClaim()
	case 2:
		if rand.Intn(2) == 0 {
			cycleClaim()
		} else {
			starClaim()
		}
	default:
		bipartiteClaim()
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// Mantel's theorem (numeric)
// ---------------------------------------------------------------------------

type mantelGen struct{}

func (g *mantelGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 3 + rand.Intn(22)
	max := n * n / 4
	left := n / 2
	right := n - left
	answer := strconv.Itoa(max)
	return generator.Problem{
		Question: fmt.Sprintf("By Mantel's theorem, what is the maximum number of edges in a triangle-free graph on %d vertices? (enter a number)", n),
		Answer:   answer,
		Explanation: fmt.Sprintf("The bound is floor(n^2/4), attained by the balanced complete bipartite graph: floor(%d^2/4) = %d, with parts of sizes %d and %d, so %d times %d edges.",
			n, max, left, right, left, right),
		Facts: map[string]string{"n": strconv.Itoa(n), "answer": answer},
	}
}

// ---------------------------------------------------------------------------
// Turán's theorem (numeric)
// ---------------------------------------------------------------------------

type turanGen struct{}

// Generate draws n and r for the Turán graph. The first ordinal takes r = 2,
// where the bound is Mantel's, the second r = 3, the case the worked example
// uses, and the third the general r of the quadratic formula. r stays at most 6
// because the balanced part sizes must give exactly (r-1)n^2/(2r) edges, and
// that identity holds for every n only up to r = 7.
func (g *turanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	var n, r int
	switch o {
	case 1:
		r = 2
		n = graphRange(4, 16)
	case 2:
		r = 3
		n = graphRange(5, 17)
	default:
		r = graphRange(4, 6)
		n = graphRange(r+1, 19)
	}
	// The balanced complete r-partite graph has parts as equal as possible; its
	// edge count is what Turán's theorem asserts is maximal.
	base, extra := n/r, n%r
	edges := 0
	for i := 0; i < r; i++ {
		for j := i + 1; j < r; j++ {
			left, right := base, base
			if i < extra {
				left++
			}
			if j < extra {
				right++
			}
			edges += left * right
		}
	}
	parts := make([]string, 0, r)
	for i := 0; i < r; i++ {
		size := base
		if i < extra {
			size++
		}
		parts = append(parts, strconv.Itoa(size))
	}
	question := fmt.Sprintf("By Turán's theorem, what is the maximum number of edges in a graph on %d vertices that contains no K_%d? (enter a number)", n, r+1)
	reason := fmt.Sprintf("The Turán graph T_%d(%d) splits the %d vertices into %d parts of sizes %s and joins every pair of vertices in different parts, giving %d edges, and no K_%d-free graph on %d vertices has more.", r, n, n, r, strings.Join(parts, ", "), edges, r+1, n)
	return graphNumeric(o, question, edges, reason, map[string]string{"n": strconv.Itoa(n), "r": strconv.Itoa(r)})
}

// ---------------------------------------------------------------------------
// Ramsey numbers (numeric)
// ---------------------------------------------------------------------------

// graphRamseyValues holds the Ramsey numbers whose exact value is known for
// every pair this generator draws. The two-argument family is the elementary
// identity R(2, k) = k, computed rather than listed, so only the pairs with both
// arguments at least 3 need a table.
var graphRamseyValues = map[[2]int]int{
	{3, 3}: 6, {3, 4}: 9, {3, 5}: 14, {3, 6}: 18, {3, 7}: 23, {3, 8}: 28,
	{4, 4}: 18, {4, 5}: 25,
}

// graphRamseyPairs lists the pairs above in a fixed order, so the draw is
// reproducible without sorting a map's keys.
var graphRamseyPairs = [][2]int{{3, 3}, {3, 4}, {3, 5}, {3, 6}, {3, 7}, {3, 8}, {4, 4}, {4, 5}}

type ramseyGen struct{}

// Generate draws a Ramsey question in one of three families: the identity
// R(2, k) = k, the small exact values of R(3, k) and R(4, k), and the
// Erdős-Szekeres upper bound R(a, b) <= C(a+b-2, a-1), which is what a learner
// can actually compute for the pairs whose exact value is still open.
func (g *ramseyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}

	switch o {
	case 1:
		k := graphRange(3, 14)
		facts["kind"] = "two"
		facts["a"], facts["b"] = "2", strconv.Itoa(k)
		question := fmt.Sprintf("What is the Ramsey number R(2,%d), the least n such that every red-blue coloring of the edges of K_n contains a red K_2 or a blue K_%d? (enter a number)", k, k)
		reason := fmt.Sprintf("R(2,k) = k: a colouring of K_%d with no red edge is entirely blue, and then it contains a blue K_%d, while K_%d can still be coloured all blue without a blue K_%d, so %d is the least such n.", k, k, k-1, k, k)
		return graphNumeric(o, question, k, reason, facts)
	case 2:
		// R(a,b) = R(b,a): swapping the two colours turns a red K_a into a blue
		// K_a, so both orders of the same pair are the same number and the two
		// orders are two different questions.
		pair := graphRamseyPairs[rand.Intn(len(graphRamseyPairs))]
		a, b := pair[0], pair[1]
		if a != b && rand.Intn(2) == 0 {
			a, b = b, a
		}
		answer := graphRamseyValues[pair]
		facts["kind"] = "exact"
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question := fmt.Sprintf("What is the Ramsey number R(%d,%d), the least n such that every red-blue coloring of the edges of K_n contains a red K_%d or a blue K_%d? (enter a number)", a, b, a, b)
		reason := fmt.Sprintf("R(%d,%d) = %d: some red-blue colouring of K_%d avoids both a red K_%d and a blue K_%d, and every colouring of K_%d contains one of them, so %d is the least such n. The value is symmetric, R(a,b) = R(b,a), because swapping the two colours exchanges the two cliques.", a, b, answer, answer-1, a, b, answer, answer)
		return graphNumeric(o, question, answer, reason, facts)
	default:
		a := graphRange(3, 6)
		b := graphRange(3, 6)
		bound := graphChoose(a+b-2, a-1)
		facts["kind"] = "bound"
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question := fmt.Sprintf("By the Erdős-Szekeres bound R(a,b) <= C(a+b-2, a-1), what upper bound does it give for R(%d,%d)? (enter a number)", a, b)
		reason := fmt.Sprintf("Substitute a = %d and b = %d: C(%d + %d - 2, %d - 1) = C(%d, %d) = %d, so R(%d,%d) is at most %d.", a, b, a, b, a, a+b-2, a-1, bound, a, b, bound)
		return graphNumeric(o, question, bound, reason, facts)
	}
}

// ---------------------------------------------------------------------------
// Random graph thresholds (multiple choice)
// ---------------------------------------------------------------------------

type randomThresholdGen struct{}

func (g *randomThresholdGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return randomThresholdProblem()
}

// ---------------------------------------------------------------------------
// Moore bound (numeric)
// ---------------------------------------------------------------------------

type mooreBoundGen struct{}

// Generate draws the degree the Moore bound is evaluated at. The three ordinals
// ask it as the degree-and-diameter count, as the Petersen graph worked example
// and as the definition of a Moore graph.
func (g *mooreBoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	var question, reason string
	var d int
	switch o {
	case 1:
		d = graphRange(1, 14)
		question = fmt.Sprintf("What is the Moore bound d^2 + 1 on the number of vertices of a graph with maximum degree %d and diameter 2? (enter a number)", d)
		reason = fmt.Sprintf("Count the vertex itself, its %d neighbours and the at most %d * %d vertices at distance two: 1 + %d + %d = %d.", d, d, d-1, d, d*(d-1), d*d+1)
	case 2:
		d = graphRange(2, 14)
		question = fmt.Sprintf("A graph has maximum degree %d and diameter 2. How many vertices can it have at most, by the Moore bound? (enter a number)", d)
		reason = fmt.Sprintf("At distance 0 there is 1 vertex, at distance 1 at most %d, and at distance 2 at most %d * %d, so the bound is 1 + %d + %d = %d.", d, d, d-1, d, d*(d-1), d*d+1)
	default:
		d = graphRange(2, 14)
		question = fmt.Sprintf("What is the Moore bound for a graph of maximum degree %d and diameter 2? A Moore graph is a graph that meets the bound exactly. (enter a number)", d)
		reason = fmt.Sprintf("The bound counts the vertex itself, its %d neighbours and at most %d * %d vertices at distance two: 1 + %d + %d = %d. A Moore graph meets this bound with equality.", d, d, d-1, d, d*(d-1), d*d+1)
	}
	return graphNumeric(o, question, d*d+1, reason, map[string]string{"d": strconv.Itoa(d)})
}

// ---------------------------------------------------------------------------
// Max-flow min-cut (numeric)
// ---------------------------------------------------------------------------

type maxFlowGen struct{}

// Generate draws a small layered network and computes its maximum flow with
// Edmonds-Karp, so the answer is derived from the drawn capacities rather than
// read from a table. The three ordinals draw networks of different shapes: a
// chain with few chords, a pair of parallel paths and a denser layered network.
func (g *maxFlowGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	var nodes, chords, minCap, maxCap int
	switch o {
	case 1:
		nodes, chords, minCap, maxCap = graphRange(3, 5), graphRange(0, 2), 2, 9
	case 2:
		nodes, chords, minCap, maxCap = graphRange(3, 4), graphRange(2, 3), 1, 9
	default:
		nodes, chords, minCap, maxCap = graphRange(4, 6), graphRange(3, 6), 1, 9
	}
	sink := nodes - 1
	edges := make([][3]int, 0, nodes+chords)
	seen := map[[2]int]bool{}
	put := func(u, v, c int) {
		if u >= v || seen[[2]int{u, v}] {
			return
		}
		seen[[2]int{u, v}] = true
		edges = append(edges, [3]int{u, v, c})
	}
	// A chain from source to sink keeps the flow positive and the network
	// connected; the chords add the alternative routes that make the cut
	// interesting.
	for i := 0; i < sink; i++ {
		put(i, i+1, graphRange(minCap, maxCap))
	}
	for i := 0; i < chords; i++ {
		u := rand.Intn(nodes)
		v := rand.Intn(nodes)
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		put(u, v, graphRange(minCap, maxCap))
	}
	answer := graphMaxFlowValue(nodes, edges)

	parts := make([]string, 0, len(edges))
	specs := make([]string, 0, len(edges))
	for _, e := range edges {
		parts = append(parts, fmt.Sprintf("%d->%d (%d)", e[0], e[1], e[2]))
		specs = append(specs, fmt.Sprintf("%d,%d,%d", e[0], e[1], e[2]))
	}
	question := fmt.Sprintf("A network has source 0, sink %d and edges %s. What is the maximum flow? (enter a number)", sink, strings.Join(parts, ", "))
	reason := fmt.Sprintf("Augmenting along residual paths until none remains sends %d units from source 0 to sink %d, and by the max-flow min-cut theorem no cut of this network has smaller capacity, so the maximum flow is %d.", answer, sink, answer)
	return graphNumeric(o, question, answer, reason, map[string]string{
		"nodes": strconv.Itoa(nodes),
		"edges": strings.Join(specs, ";"),
		"s":     "0",
		"t":     strconv.Itoa(sink),
	})
}

// ---------------------------------------------------------------------------
// Tournaments (numeric)
// ---------------------------------------------------------------------------

type tournamentGen struct{}

// Generate draws the order of the tournament. The first ordinal counts its
// games, the second orients them and raises two to that count, and the third
// reads the score sequence, whose total is the same game count.
func (g *tournamentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	var question, reason, kind string
	var n, answer int
	switch o {
	case 1:
		n = graphRange(3, 15)
		kind, answer = "games", n*(n-1)/2
		question = fmt.Sprintf("A tournament on %d vertices orients every edge of K_%d. How many games does it decide? (enter a number)", n, n)
		reason = fmt.Sprintf("Every one of the %d games is an edge of K_%d, and K_%d has %d * %d / 2 = %d edges.", answer, n, n, n, n-1, answer)
	case 2:
		if rand.Intn(2) == 0 {
			n = graphRange(1, 8)
			kind, answer = "count", graphIntPow(2, n*(n-1)/2)
			question = fmt.Sprintf("A tournament orients every edge of K_n. How many distinct tournaments are there on %d labeled vertices? (enter a number)", n)
			reason = fmt.Sprintf("K_%d has %d * %d / 2 = %d games and each is oriented one of two ways, so there are 2^%d = %d tournaments.", n, n, n-1, n*(n-1)/2, n*(n-1)/2, answer)
		} else {
			n = graphRange(3, 6)
			kind, answer = "transitive", graphFactorial(n)
			question = fmt.Sprintf("A tournament is transitive when its vertices can be listed so that every edge points from an earlier vertex to a later one. How many transitive tournaments are there on %d labeled vertices? (enter a number)", n)
			reason = fmt.Sprintf("A transitive tournament is determined by the linear order of its %d vertices, and no two orders give the same orientation, so there are %d! = %d of them.", n, n, answer)
		}
	default:
		n = graphRange(3, 15)
		if rand.Intn(2) == 0 {
			kind, answer = "scores", n*(n-1)/2
			question = fmt.Sprintf("A tournament on %d vertices has one out-degree per vertex, called its score. What is the sum of all %d scores? (enter a number)", n, n)
			reason = fmt.Sprintf("Each of the %d games adds exactly one to exactly one score, so the scores sum to the number of games, %d * %d / 2 = %d.", n*(n-1)/2, n, n-1, answer)
		} else {
			kind, answer = "path", n-1
			question = fmt.Sprintf("Every tournament has a Hamiltonian path, a directed path that visits every vertex once. How many edges does such a path have in a tournament on %d vertices? (enter a number)", n)
			reason = fmt.Sprintf("A path that visits %d vertices once uses one edge between each consecutive pair, so it has %d - 1 = %d edges.", n, n, answer)
		}
	}
	return graphNumeric(o, question, answer, reason, map[string]string{"n": strconv.Itoa(n), "kind": kind})
}

// ---------------------------------------------------------------------------
// Strong connectivity (multiple choice)
// ---------------------------------------------------------------------------

type strongConnectivityGen struct{}

func (g *strongConnectivityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return strongConnectivityProblem()
}

// ---------------------------------------------------------------------------
// Topological sort (multiple choice)
// ---------------------------------------------------------------------------

type topologicalSortGen struct{}

func (g *topologicalSortGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return topologicalSortProblem()
}

// ---------------------------------------------------------------------------
// Counting paths in a DAG (numeric)
// ---------------------------------------------------------------------------

type dagPathsGen struct{}

func (g *dagPathsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// A layered DAG: edges only run from one layer to the next, ending at a
	// single sink, so the path count is computed by dynamic programming rather
	// than read from a table.
	layers := 2 + rand.Intn(3)
	width := 1 + rand.Intn(3)
	nodes := layers*width + 1
	sink := nodes - 1
	counts := make([]int, width)
	for i := range counts {
		counts[i] = 1
	}
	edges := [][2]int{}
	for l := 0; l < layers-1; l++ {
		next := make([]int, width)
		for j := 0; j < width; j++ {
			for i := 0; i < width; i++ {
				if rand.Intn(4) == 0 {
					continue
				}
				edges = append(edges, [2]int{l*width + i, (l+1)*width + j})
				next[j] += counts[i]
			}
			if next[j] == 0 {
				i := rand.Intn(width)
				edges = append(edges, [2]int{l*width + i, (l+1)*width + j})
				next[j] += counts[i]
			}
		}
		counts = next
	}
	total := 0
	for j := 0; j < width; j++ {
		edges = append(edges, [2]int{(layers-1)*width + j, sink})
		total += counts[j]
	}
	parts := make([]string, 0, len(edges))
	for _, e := range edges {
		parts = append(parts, strconv.Itoa(e[0])+","+strconv.Itoa(e[1]))
	}
	// Count only paths that start at vertex 0. Edges were appended layer by
	// layer, so a single forward sweep is a topological order.
	paths := make([]int, nodes)
	paths[0] = 1
	for _, e := range edges {
		paths[e[1]] += paths[e[0]]
	}
	total = paths[sink]
	return generator.Problem{
		Question: fmt.Sprintf("In the DAG with vertices 0..%d and edges %s, how many directed paths go from 0 to %d? (enter a number)", sink, edgeList(parts), sink),
		Answer:   strconv.Itoa(total),
		Explanation: "A dynamic program over a topological order adds the path counts of each vertex's predecessors, " +
			"so every path is counted exactly once.",
		Facts: map[string]string{
			"nodes": strconv.Itoa(nodes), "edges": strings.Join(parts, ";"),
			"s": "0", "t": strconv.Itoa(sink), "answer": strconv.Itoa(total),
		},
	}
}

// edgeList renders edges for the question text, e.g. 0->1, 1->2.
func edgeList(parts []string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.Replace(p, ",", "->", 1))
	}
	return strings.Join(out, ", ")
}
