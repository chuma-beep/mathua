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

func (g *treewidthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question            string
		answer              string
		kind                string
		n, rows, cols, a, b int
		reason              string
	}
	claims := []claim{
		{question: "What is the treewidth of a tree on 6 vertices? (enter a number)",
			answer: "1", kind: "tree", n: 6,
			reason: "A tree has a tree decomposition of width 1, and any connected graph with an edge has width at least 1."},
		{question: "What is the treewidth of the path P_5? (enter a number)",
			answer: "1", kind: "path", n: 5,
			reason: "A path is a tree, so its treewidth is 1."},
		{question: "What is the treewidth of the cycle C_6? (enter a number)",
			answer: "2", kind: "cycle", n: 6,
			reason: "A cycle has a tree decomposition whose bags are the edges, of width 2, and it is not a tree."},
		{question: "What is the treewidth of the complete graph K_5? (enter a number)",
			answer: "4", kind: "complete", n: 5,
			reason: "Every pair of vertices is adjacent, so one bag must hold all n vertices and the width is n - 1 = 4."},
		{question: "What is the treewidth of the 3-by-4 grid graph? (enter a number)",
			answer: "3", kind: "grid", rows: 3, cols: 4,
			reason: "An r-by-c grid has treewidth min(r, c), which is 3 here."},
		{question: "What is the treewidth of the complete bipartite graph K_{2,3}? (enter a number)",
			answer: "2", kind: "bipartite", a: 2, b: 3,
			reason: "K_{a,b} has treewidth min(a, b) for a, b at least 1, which is 2 here."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	switch c.kind {
	case "grid":
		facts["rows"] = strconv.Itoa(c.rows)
		facts["cols"] = strconv.Itoa(c.cols)
	case "bipartite":
		facts["a"] = strconv.Itoa(c.a)
		facts["b"] = strconv.Itoa(c.b)
	default:
		facts["n"] = strconv.Itoa(c.n)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// Graph genus (numeric)
// ---------------------------------------------------------------------------

type genusGen struct{}

func (g *genusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n, a, b  int
		reason   string
	}
	claims := []claim{
		{question: "What is the genus of the complete graph K_5? (enter a number)",
			answer: "1", kind: "complete", n: 5,
			reason: "ceil((5-3)(5-4)/12) = ceil(2/12) = 1: K_5 embeds on the torus but not the sphere."},
		{question: "What is the genus of the complete graph K_6? (enter a number)",
			answer: "1", kind: "complete", n: 6,
			reason: "ceil((6-3)(6-4)/12) = ceil(6/12) = 1."},
		{question: "What is the genus of the complete graph K_8? (enter a number)",
			answer: "2", kind: "complete", n: 8,
			reason: "ceil((8-3)(8-4)/12) = ceil(20/12) = 2."},
		{question: "What is the genus of the complete bipartite graph K_{3,3}? (enter a number)",
			answer: "1", kind: "bipartite", a: 3, b: 3,
			reason: "ceil((3-2)(3-2)/4) = ceil(1/4) = 1: K_{3,3} needs the torus."},
		{question: "What is the genus of a planar graph? (enter a number)",
			answer: "0", kind: "planar",
			reason: "A planar graph is embedded in the sphere, whose genus is 0."},
		{question: "What is the genus of the cycle C_6? (enter a number)",
			answer: "0", kind: "cycle", n: 6,
			reason: "Every cycle is planar, so its genus is 0."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	switch c.kind {
	case "complete", "cycle":
		facts["n"] = strconv.Itoa(c.n)
	case "bipartite":
		facts["a"] = strconv.Itoa(c.a)
		facts["b"] = strconv.Itoa(c.b)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// Outerplanar graphs (numeric)
// ---------------------------------------------------------------------------

type outerplanarGen struct{}

func (g *outerplanarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		n      int
		answer string
		reason string
	}
	claims := []claim{
		{n: 3, answer: "3", reason: "A maximal outerplanar graph on 3 vertices is the triangle, with 2*3 - 3 = 3 edges."},
		{n: 4, answer: "5", reason: "2*4 - 3 = 5 edges: a triangle plus one vertex joined to three others."},
		{n: 5, answer: "7", reason: "2*5 - 3 = 7 edges."},
		{n: 6, answer: "9", reason: "2*6 - 3 = 9 edges."},
		{n: 7, answer: "11", reason: "2*7 - 3 = 11 edges."},
		{n: 8, answer: "13", reason: "2*8 - 3 = 13 edges."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: "A simple graph is outerplanar when it can be drawn in the plane with every vertex on the " +
			"outer face. What is the maximum number of edges of an outerplanar graph on " +
			strconv.Itoa(c.n) + " vertices? (enter a number)",
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"n": strconv.Itoa(c.n), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// Adjacency matrix (numeric)
// ---------------------------------------------------------------------------

type adjacencyGen struct{}

func (g *adjacencyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n, i, j  int
		length   int
		reason   string
	}
	claims := []claim{
		{question: "In the path P_3, how many walks of length 2 go from vertex 1 to vertex 1? (enter a number)",
			answer: "1", kind: "path", n: 3, i: 1, j: 1, length: 2,
			reason: "(A^2)_{11} counts walks 1-2-1, and there is exactly one."},
		{question: "In the path P_3, how many walks of length 2 go from vertex 1 to vertex 3? (enter a number)",
			answer: "1", kind: "path", n: 3, i: 1, j: 3, length: 2,
			reason: "(A^2)_{13} counts walks 1-2-3, and there is exactly one."},
		{question: "In the cycle C_4, how many walks of length 2 go from vertex 1 to vertex 1? (enter a number)",
			answer: "2", kind: "cycle", n: 4, i: 1, j: 1, length: 2,
			reason: "(A^2)_{11} counts the two ways out and back: 1-2-1 and 1-4-1."},
		{question: "In the cycle C_4, how many walks of length 2 go from vertex 1 to vertex 3? (enter a number)",
			answer: "2", kind: "cycle", n: 4, i: 1, j: 3, length: 2,
			reason: "(A^2)_{13} counts 1-2-3 and 1-4-3."},
		{question: "In the complete graph K_3, how many walks of length 2 go from vertex 1 to vertex 1? (enter a number)",
			answer: "2", kind: "complete", n: 3, i: 1, j: 1, length: 2,
			reason: "(A^2)_{11} is the degree of vertex 1 in K_3, which is 2."},
		{question: "In the path P_4, how many walks of length 2 go from vertex 1 to vertex 3? (enter a number)",
			answer: "1", kind: "path", n: 4, i: 1, j: 3, length: 2,
			reason: "The only walk is 1-2-3; there is no shortcut between 1 and 3."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: c.question,
		Answer:   c.answer,
		Explanation: c.reason + " The (i,j) entry of A^k counts the walks of length k from i to j; " +
			"here the answer is the (i,j) entry of A^2.",
		Facts: map[string]string{
			"kind": c.kind, "n": strconv.Itoa(c.n), "i": strconv.Itoa(c.i), "j": strconv.Itoa(c.j),
			"length": strconv.Itoa(c.length), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Laplacian matrix (numeric)
// ---------------------------------------------------------------------------

type laplacianGen struct{}

func (g *laplacianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question      string
		answer        string
		kind          string
		n, components int
		reason        string
	}
	claims := []claim{
		{question: "How many eigenvalues of the Laplacian of the path P_4 equal zero? (enter a number)",
			answer: "1", kind: "path", n: 4,
			reason: "A connected graph has exactly one zero Laplacian eigenvalue, and P_4 is connected."},
		{question: "How many eigenvalues of the Laplacian of the cycle C_5 equal zero? (enter a number)",
			answer: "1", kind: "cycle", n: 5,
			reason: "A cycle is connected, so the multiplicity of the zero eigenvalue is 1."},
		{question: "How many eigenvalues of the Laplacian of a tree on 6 vertices equal zero? (enter a number)",
			answer: "1", kind: "tree", n: 6,
			reason: "Every tree is connected, so there is exactly one zero Laplacian eigenvalue."},
		{question: "How many eigenvalues of the Laplacian of the complete graph K_4 equal zero? (enter a number)",
			answer: "1", kind: "complete", n: 4,
			reason: "K_4 is connected; its Laplacian eigenvalues are 0, 4, 4, 4."},
		{question: "How many eigenvalues of the Laplacian of a graph with two disjoint edges equal zero? (enter a number)",
			answer: "2", kind: "disjoint", n: 4, components: 2,
			reason: "The multiplicity of 0 equals the number of connected components, and there are two."},
		{question: "How many eigenvalues of the Laplacian of three disjoint triangles equal zero? (enter a number)",
			answer: "3", kind: "disjoint", n: 9, components: 3,
			reason: "Each connected component contributes one zero eigenvalue, so the multiplicity is 3."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.kind == "disjoint" {
		facts["components"] = strconv.Itoa(c.components)
		facts["n"] = strconv.Itoa(c.n)
	} else {
		facts["n"] = strconv.Itoa(c.n)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// Matrix-tree theorem (numeric)
// ---------------------------------------------------------------------------

type matrixTreeGen struct{}

func (g *matrixTreeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n, a, b  int
		reason   string
	}
	claims := []claim{
		{question: "How many spanning trees does the complete graph K_4 have? (enter a number)",
			answer: "16", kind: "complete", n: 4,
			reason: "Cayley's formula gives n^{n-2} = 4^2 = 16."},
		{question: "How many spanning trees does the complete graph K_5 have? (enter a number)",
			answer: "125", kind: "complete", n: 5,
			reason: "n^{n-2} = 5^3 = 125."},
		{question: "How many spanning trees does the complete graph K_3 have? (enter a number)",
			answer: "3", kind: "complete", n: 3,
			reason: "n^{n-2} = 3^1 = 3: delete one of the three edges."},
		{question: "How many spanning trees does the cycle C_5 have? (enter a number)",
			answer: "5", kind: "cycle", n: 5,
			reason: "A spanning tree of C_n is obtained by deleting any one of its n edges, so there are 5."},
		{question: "How many spanning trees does the path P_6 have? (enter a number)",
			answer: "1", kind: "path", n: 6,
			reason: "A path is already a tree, so it has exactly one spanning tree, itself."},
		{question: "How many spanning trees does the complete bipartite graph K_{2,3} have? (enter a number)",
			answer: "12", kind: "bipartite", a: 2, b: 3,
			reason: "K_{a,b} has a^{b-1} b^{a-1} = 2^2 * 3^1 = 12 spanning trees."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	switch c.kind {
	case "complete", "cycle", "path":
		facts["n"] = strconv.Itoa(c.n)
	case "bipartite":
		facts["a"] = strconv.Itoa(c.a)
		facts["b"] = strconv.Itoa(c.b)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// Spectral radius (numeric)
// ---------------------------------------------------------------------------

type spectralRadiusGen struct{}

func (g *spectralRadiusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n, d     int
		reason   string
	}
	claims := []claim{
		{question: "What is the spectral radius of the complete graph K_5? (enter a number)",
			answer: "4", kind: "complete", n: 5,
			reason: "K_n is (n-1)-regular, and for a d-regular graph the spectral radius is d, so it is 4."},
		{question: "What is the spectral radius of the complete graph K_6? (enter a number)",
			answer: "5", kind: "complete", n: 6,
			reason: "K_6 is 5-regular, so the largest adjacency eigenvalue is 5."},
		{question: "What is the spectral radius of the complete graph K_7? (enter a number)",
			answer: "6", kind: "complete", n: 7,
			reason: "K_7 is 6-regular, so the spectral radius is 6."},
		{question: "What is the spectral radius of the cycle C_6? (enter a number)",
			answer: "2", kind: "cycle", n: 6,
			reason: "Every cycle is 2-regular, so its largest adjacency eigenvalue is 2."},
		{question: "What is the spectral radius of the Petersen graph? (enter a number)",
			answer: "3", kind: "regular", d: 3,
			reason: "The Petersen graph is 3-regular, so its spectral radius is 3."},
		{question: "What is the spectral radius of the complete bipartite graph K_{4,4}? (enter a number)",
			answer: "4", kind: "regular", d: 4,
			reason: "K_{4,4} is 4-regular, so its spectral radius is 4."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.kind == "regular" {
		facts["d"] = strconv.Itoa(c.d)
	} else {
		facts["n"] = strconv.Itoa(c.n)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// Algebraic connectivity (numeric)
// ---------------------------------------------------------------------------

type spectralGapGen struct{}

func (g *spectralGapGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n, a, b  int
		reason   string
	}
	claims := []claim{
		{question: "What is the algebraic connectivity of the complete graph K_4? (enter a number)",
			answer: "4", kind: "complete", n: 4,
			reason: "The Laplacian eigenvalues of K_n are 0 and n, repeated n - 1 times, so the second-smallest is n = 4."},
		{question: "What is the algebraic connectivity of the complete graph K_3? (enter a number)",
			answer: "3", kind: "complete", n: 3,
			reason: "The Laplacian eigenvalues of K_3 are 0, 3, 3."},
		{question: "What is the algebraic connectivity of the cycle C_4? (enter a number)",
			answer: "2", kind: "cycle", n: 4,
			reason: "The Laplacian eigenvalues of C_4 are 0, 2, 2, 4, so the second-smallest is 2."},
		{question: "What is the algebraic connectivity of the cycle C_6? (enter a number)",
			answer: "1", kind: "cycle", n: 6,
			reason: "For C_n the second-smallest Laplacian eigenvalue is 2 - 2cos(2pi/n), which is 1 at n = 6."},
		{question: "What is the algebraic connectivity of the complete bipartite graph K_{2,3}? (enter a number)",
			answer: "2", kind: "bipartite", a: 2, b: 3,
			reason: "K_{a,b} has algebraic connectivity min(a, b), which is 2 here."},
		{question: "What is the algebraic connectivity of the star K_{1,4}? (enter a number)",
			answer: "1", kind: "star",
			reason: "A star on at least two vertices has algebraic connectivity 1."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	switch c.kind {
	case "complete", "cycle":
		facts["n"] = strconv.Itoa(c.n)
	case "bipartite":
		facts["a"] = strconv.Itoa(c.a)
		facts["b"] = strconv.Itoa(c.b)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
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

func (g *turanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		n, r   int
		answer string
		reason string
	}
	claims := []claim{
		{n: 5, r: 2, answer: "6", reason: "Turán's graph T_2(5) = K_{2,3} has floor(5^2/4) = 6 edges."},
		{n: 6, r: 2, answer: "9", reason: "T_2(6) = K_{3,3} has floor(6^2/4) = 9 edges."},
		{n: 5, r: 3, answer: "8", reason: "T_3(5) has parts 2, 2, 1, so 2*2 + 2*1 + 2*1 = 8 edges."},
		{n: 6, r: 3, answer: "12", reason: "T_3(6) has parts 2, 2, 2, so 3 * (2*2) = 12 edges."},
		{n: 7, r: 3, answer: "16", reason: "T_3(7) has parts 3, 2, 2, so 6 + 6 + 4 = 16 edges."},
		{n: 8, r: 4, answer: "24", reason: "T_4(8) has four parts of size 2, so C(4,2) * 4 = 24 edges."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: "By Turán's theorem, what is the maximum number of edges in a graph on " +
			strconv.Itoa(c.n) + " vertices that contains no K_" + strconv.Itoa(c.r+1) +
			"? (enter a number)",
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"n": strconv.Itoa(c.n), "r": strconv.Itoa(c.r), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// Ramsey numbers (numeric)
// ---------------------------------------------------------------------------

type ramseyGen struct{}

func (g *ramseyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		a, b   int
		answer string
		reason string
	}
	claims := []claim{
		{a: 3, b: 3, answer: "6", reason: "R(3,3) = 6: five vertices admit a 5-cycle coloring with no monochromatic triangle, and six do not."},
		{a: 2, b: 4, answer: "4", reason: "R(2,4) = 4, since a red edge alone gives a red K_2 and four vertices force a blue K_4 or a red K_2."},
		{a: 2, b: 5, answer: "5", reason: "R(2,k) = k, so R(2,5) = 5."},
		{a: 3, b: 4, answer: "9", reason: "R(3,4) = 9."},
		{a: 3, b: 5, answer: "14", reason: "R(3,5) = 14."},
		{a: 4, b: 4, answer: "18", reason: "R(4,4) = 18."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: "What is the Ramsey number R(" + strconv.Itoa(c.a) + "," + strconv.Itoa(c.b) +
			"), the least n such that every red-blue coloring of the edges of K_n contains a red K_" +
			strconv.Itoa(c.a) + " or a blue K_" + strconv.Itoa(c.b) + "? (enter a number)",
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"a": strconv.Itoa(c.a), "b": strconv.Itoa(c.b), "answer": c.answer},
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

func (g *mooreBoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		d      int
		answer string
		reason string
	}
	claims := []claim{
		{d: 1, answer: "2", reason: "d^2 + 1 = 1 + 1 = 2, attained by K_2: one vertex plus one neighbour at distance 1."},
		{d: 2, answer: "5", reason: "d^2 + 1 = 4 + 1 = 5, attained by the cycle C_5."},
		{d: 3, answer: "10", reason: "d^2 + 1 = 9 + 1 = 10, attained by the Petersen graph."},
		{d: 4, answer: "17", reason: "d^2 + 1 = 16 + 1 = 17."},
		{d: 5, answer: "26", reason: "d^2 + 1 = 25 + 1 = 26."},
		{d: 6, answer: "37", reason: "d^2 + 1 = 36 + 1 = 37."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: "What is the Moore bound d^2 + 1 on the number of vertices of a graph with maximum degree " +
			strconv.Itoa(c.d) + " and diameter 2? (enter a number)",
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"d": strconv.Itoa(c.d), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// Max-flow min-cut (numeric)
// ---------------------------------------------------------------------------

type maxFlowGen struct{}

func (g *maxFlowGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		nodes    int
		edges    string
		reason   string
	}
	claims := []claim{
		{question: "A network has source 0, sink 1 and edges 0->1 with capacities 3 and 4 in parallel. What is the maximum flow? (enter a number)",
			answer: "7", nodes: 2, edges: "0,1,3;0,1,4",
			reason: "The two parallel edges carry 3 + 4 = 7, and the cut separating source from sink has capacity 7."},
		{question: "A network has source 0, sink 3 and edges 0->1 (3), 0->2 (3), 1->3 (2), 2->3 (3). What is the maximum flow? (enter a number)",
			answer: "5", nodes: 4, edges: "0,1,3;0,2,3;1,3,2;2,3,3",
			reason: "The two paths 0-1-3 and 0-2-3 carry 2 + 3 = 5, meeting the cut capacity 5."},
		{question: "A network has source 0, sink 3 and edges 0->1 (4), 1->3 (3), 0->3 (1). What is the maximum flow? (enter a number)",
			answer: "4", nodes: 4, edges: "0,1,4;1,3,3;0,3,1",
			reason: "The direct edge carries 1 and the path 0-1-3 carries 3, for 4; the cut separating vertex 1 has capacity 4."},
		{question: "A network has source 0, sink 3 and edges 0->1 (2), 0->2 (2), 1->3 (2), 2->3 (2), 1->2 (1). What is the maximum flow? (enter a number)",
			answer: "4", nodes: 4, edges: "0,1,2;0,2,2;1,3,2;2,3,2;1,2,1",
			reason: "The two disjoint paths carry 2 + 2 = 4, which equals the capacity of the cut around the source."},
		{question: "A network has source 0, sink 2 and edges 0->1 (5), 1->2 (2), 0->2 (2). What is the maximum flow? (enter a number)",
			answer: "4", nodes: 3, edges: "0,1,5;1,2,2;0,2,2",
			reason: "The path through vertex 1 carries 2 and the direct edge carries 2, for 4; the cut separating vertex 1 has capacity 4."},
		{question: "A network has source 0, sink 3 and edges 0->1 (3), 0->2 (3), 1->2 (2), 1->3 (2), 2->3 (3). What is the maximum flow? (enter a number)",
			answer: "5", nodes: 4, edges: "0,1,3;0,2,3;1,2,2;1,3,2;2,3,3",
			reason: "Send 2 along 0-1-3 and 3 along 0-2-3 for 5; the cut separating 1 and 2 from the sink has capacity 5."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason + " By the max-flow min-cut theorem this value equals the capacity of a minimum cut.",
		Facts: map[string]string{
			"nodes": strconv.Itoa(c.nodes), "edges": c.edges, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// Tournaments (numeric)
// ---------------------------------------------------------------------------

type tournamentGen struct{}

func (g *tournamentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		n      int
		answer string
		reason string
	}
	claims := []claim{
		{n: 1, answer: "1", reason: "A single vertex has no edges, so the empty orientation is the only tournament."},
		{n: 2, answer: "2", reason: "The single pair is oriented one of two ways, so there are 2 tournaments."},
		{n: 3, answer: "8", reason: "2^{C(3,2)} = 2^3 = 8."},
		{n: 4, answer: "64", reason: "2^{C(4,2)} = 2^6 = 64."},
		{n: 5, answer: "1024", reason: "2^{C(5,2)} = 2^10 = 1024."},
		{n: 6, answer: "32768", reason: "2^{C(6,2)} = 2^15 = 32768."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question: "A tournament orients every edge of K_n. How many distinct tournaments are there on " +
			strconv.Itoa(c.n) + " labeled vertices? (enter a number)",
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"n": strconv.Itoa(c.n), "answer": c.answer},
	}
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
