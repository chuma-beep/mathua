package graphtheory

// The first graph-theory batch: twenty concepts that open the domain with
// isomorphism, basic constructions, connectivity, matching and coloring.
//
// The discrete-math domain already has graphs, paths, trees, Eulerian and
// Hamiltonian graphs, coloring, bipartite matchings, planarity, Cayley and
// degree sequences. This domain starts after those and goes deeper: it treats
// isomorphism and homomorphism as structure-preserving maps, counts distances
// and cuts, proves the matching/cover min-max pair, and develops chromatic and
// list coloring. The second batch will add planarity, spectral graph theory,
// extremal and random graphs, and flows.
//
// Ten concepts are numeric and ten are multiple choice. Every numeric
// generator publishes the graph parameters its answer follows from, so the
// tests recompute the answer from those facts.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

// ---------------------------------------------------------------------------
// 1. Graph isomorphism
// ---------------------------------------------------------------------------

type isomorphismGen struct{}

func (g *isomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return isomorphismProblem()
}

// ---------------------------------------------------------------------------
// 2. The complement of a graph (numeric)
// ---------------------------------------------------------------------------

type complementGen struct{}

func (g *complementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 4 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 6) // 4..12
	totalEdges := n * (n - 1) / 2
	m := 1 + rand.Intn(totalEdges-1)
	ans := totalEdges - m

	prompt := fmt.Sprintf("A simple graph has %d vertices and %d edges. How many edges does its complement have?", n, m)
	expl := fmt.Sprintf("Subgoal 1: State governing formula. The complement of a simple graph \\(G = (V, E)\\) on \\(n\\) vertices contains all edges not present in \\(G\\). A complete graph \\(K_n\\) has \\(\\binom{n}{2} = \\frac{n(n - 1)}{2}\\) edges, so the complement has \\(|E(\\overline{G})| = \\frac{n(n - 1)}{2} - |E(G)|\\) edges.\nSubgoal 2: Substitute parameters. For \\(n = %d\\) and \\(|E(G)| = %d\\), calculate the total possible edges: \\(\\frac{%d \\times %d}{2} = %d\\).\nSubgoal 3: Evaluate and simplify. Subtract the existing edges: %d - %d = %d.", n, m, n, n-1, totalEdges, totalEdges, m, ans)

	return generator.Problem{
		Question:    prompt,
		Answer:      strconv.Itoa(ans),
		Explanation: expl,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"m":      strconv.Itoa(m),
			"answer": strconv.Itoa(ans),
		},
	}
}

// ---------------------------------------------------------------------------
// 3. Subgraphs and induced subgraphs
// ---------------------------------------------------------------------------

type subgraphGen struct{}

func (g *subgraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return subgraphProblem()
}

// ---------------------------------------------------------------------------
// 4. Graph homomorphisms
// ---------------------------------------------------------------------------

type homomorphismGen struct{}

func (g *homomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return homomorphismProblem()
}

// ---------------------------------------------------------------------------
// 5. Graph products (numeric)
// ---------------------------------------------------------------------------

type productGen struct{}

func (g *productGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n1 := 3 + rand.Intn(2) + mathutil.ScaleInt(ctx.Difficulty, 0, 3)
	maxM1 := n1 * (n1 - 1) / 2
	m1 := n1 - 1 + rand.Intn(maxM1-n1+2)

	n2 := 2 + rand.Intn(2) + mathutil.ScaleInt(ctx.Difficulty, 0, 3)
	maxM2 := n2 * (n2 - 1) / 2
	m2 := n2 - 1 + rand.Intn(maxM2-n2+2)

	ans := m1*n2 + n1*m2

	prompt := fmt.Sprintf("Graphs G and H have (%d vertices, %d edges) and (%d vertices, %d edges). How many edges does their Cartesian product have?", n1, m1, n2, m2)
	expl := fmt.Sprintf("Subgoal 1: State governing formula. The number of edges in the Cartesian product \\(G \\square H\\) of two graphs is given by \\(|E(G \\square H)| = |E(G)| \\cdot |V(H)| + |V(G)| \\cdot |E(H)| = m_1 n_2 + n_1 m_2\\).\nSubgoal 2: Substitute parameters. For \\(n_1 = %d\\), \\(m_1 = %d\\), \\(n_2 = %d\\) and \\(m_2 = %d\\), form the expression: (%d \\times %d) + (%d \\times %d).\nSubgoal 3: Evaluate and simplify. Multiply and sum the terms: %d + %d = %d.", n1, m1, n2, m2, m1, n2, n1, m2, m1*n2, n1*m2, ans)

	return generator.Problem{
		Question:    prompt,
		Answer:      strconv.Itoa(ans),
		Explanation: expl,
		Facts: map[string]string{
			"n1":     strconv.Itoa(n1),
			"m1":     strconv.Itoa(m1),
			"n2":     strconv.Itoa(n2),
			"m2":     strconv.Itoa(m2),
			"answer": strconv.Itoa(ans),
		},
	}
}

// ---------------------------------------------------------------------------
// 6. Distance and diameter (numeric)
// ---------------------------------------------------------------------------

type distanceGen struct{}

func (g *distanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(4)
	var prompt, expl, kind string
	var ans int
	facts := map[string]string{}

	switch mode {
	case 0:
		kind = "path"
		n := 4 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 4..11
		ans = n - 1
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("What is the diameter of the path graph P_%d?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. The diameter of a connected graph is the maximum shortest-path distance between any pair of vertices. For a path graph \\(P_n\\), the maximum distance occurs between the two leaf endpoints and equals \\(n - 1\\).\nSubgoal 2: Substitute parameters. With \\(n = %d\\), the path has %d edges between endpoints.\nSubgoal 3: Evaluate and simplify. The diameter is %d - 1 = %d.", n, n-1, n, ans)

	case 1:
		kind = "cycle"
		n := 4 + rand.Intn(4) + mathutil.ScaleInt(ctx.Difficulty, 0, 6) // 4..13
		ans = n / 2
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("What is the diameter of the cycle graph C_%d?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. For a cycle graph \\(C_n\\), the shortest path between any two vertices cannot exceed half the circumference of the cycle: \\(\\text{diam}(C_n) = \\lfloor n / 2 \\rfloor\\).\nSubgoal 2: Substitute parameters. For \\(n = %d\\), compute \\(\\lfloor %d / 2 \\rfloor\\).\nSubgoal 3: Evaluate and simplify. Evaluate integer division: %d / 2 = %d.", n, n, n, ans)

	case 2:
		kind = "complete"
		n := 3 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..10
		ans = 1
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("What is the diameter of the complete graph K_%d?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. In a complete graph \\(K_n\\) with \\(n \\ge 2\\), every pair of distinct vertices is connected by an edge of length 1. Hence the distance between every distinct pair is 1, and \\(\\text{diam}(K_n) = 1\\).\nSubgoal 2: Substitute parameters. For \\(K_%d\\), all %d vertices are mutually adjacent.\nSubgoal 3: Evaluate and simplify. The diameter of any complete graph on at least 2 vertices is 1.", n, n)

	case 3:
		kind = "grid"
		r := 2 + rand.Intn(2) + mathutil.ScaleInt(ctx.Difficulty, 0, 3) // 2..6
		c := 2 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 3) // 2..7
		ans = r + c - 2
		facts["rows"] = strconv.Itoa(r)
		facts["cols"] = strconv.Itoa(c)
		prompt = fmt.Sprintf("What is the diameter of the %d-by-%d grid graph?", r, c)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A grid graph of dimensions \\(r \\times c\\) is the Cartesian product \\(P_r \\square P_c\\). The diameter is the distance between opposite corners: \\(\\text{diam}(P_r \\square P_c) = (r - 1) + (c - 1) = r + c - 2\\).\nSubgoal 2: Substitute parameters. Substitute \\(r = %d\\) and \\(c = %d\\): (%d - 1) + (%d - 1).\nSubgoal 3: Evaluate and simplify. Sum the coordinate distances: %d + %d = %d.", r, c, r, c, r-1, c-1, ans)
	}

	facts["kind"] = kind
	facts["answer"] = strconv.Itoa(ans)

	return generator.Problem{
		Question:    prompt,
		Answer:      strconv.Itoa(ans),
		Explanation: expl,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// 7. Cut vertices (numeric)
// ---------------------------------------------------------------------------

type cutVertexGen struct{}

func (g *cutVertexGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(4)
	var prompt, expl, kind string
	var ans int
	facts := map[string]string{}

	switch mode {
	case 0:
		kind = "path"
		n := 4 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 4..11
		ans = n - 2
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many cut vertices does the path graph P_%d have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A cut vertex is a vertex whose removal increases the number of connected components. In a path graph \\(P_n\\), removing an endpoint leaves a connected path \\(P_{n-1}\\), while removing any of the internal vertices disconnects the path into two components. Total cut vertices: \\(n - 2\\).\nSubgoal 2: Substitute parameters. For \\(n = %d\\), there are 2 endpoints and \\(%d - 2\\) internal vertices.\nSubgoal 3: Evaluate and simplify. Compute the internal vertex count: %d - 2 = %d.", n, n, n, ans)

	case 1:
		kind = "star"
		n := 3 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..10
		ans = 1
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many cut vertices does the star graph K_{1,%d} have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A star graph \\(K_{1,k}\\) consists of one central vertex connected to \\(k\\) leaves. Removing any leaf leaves the remaining star connected. Removing the central vertex isolates all \\(k\\) leaves into disconnected components.\nSubgoal 2: Substitute parameters. In \\(K_{1,%d}\\), only the single central vertex disconnects the graph.\nSubgoal 3: Evaluate and simplify. The star graph has exactly 1 cut vertex.", n)

	case 2:
		kind = "cycle"
		n := 3 + rand.Intn(4) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..11
		ans = 0
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many cut vertices does the cycle graph C_%d have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. In a cycle graph \\(C_n\\) with \\(n \\ge 3\\), removing any single vertex leaves a connected path graph \\(P_{n-1}\\). Because no single vertex removal disconnects the graph, cycles have no cut vertices.\nSubgoal 2: Substitute parameters. For \\(C_%d\\), removing any vertex leaves a single connected component \\(P_%d\\).\nSubgoal 3: Evaluate and simplify. The number of cut vertices in \\(C_%d\\) is 0.", n, n-1, n)

	case 3:
		kind = "bowtie"
		ans = 1
		prompt = "How many cut vertices does the bowtie graph (two triangles sharing exactly one common vertex) have?"
		expl = "Subgoal 1: State governing formula. A cut vertex is a vertex whose removal strictly increases the number of connected components.\nSubgoal 2: Substitute parameters. The shared central vertex belongs to both triangles. Removing any of the four non-shared vertices leaves the rest connected. Removing the shared central vertex disconnects the graph into two independent edges (components).\nSubgoal 3: Evaluate and simplify. Exactly 1 vertex (the central bridge junction) is a cut vertex."
	}

	facts["kind"] = kind
	facts["answer"] = strconv.Itoa(ans)

	return generator.Problem{
		Question:    prompt,
		Answer:      strconv.Itoa(ans),
		Explanation: expl,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// 8. Bridges (numeric)
// ---------------------------------------------------------------------------

type bridgeGen struct{}

func (g *bridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(6)
	var prompt, expl, kind string
	var ans int
	facts := map[string]string{}

	switch mode {
	case 0:
		kind = "path"
		n := 3 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..10
		ans = n - 1
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many bridges does the path graph P_%d have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A bridge (cut edge) is an edge whose removal increases the number of connected components. In any tree or path graph \\(P_n\\), removing any edge disconnects the graph into two components. Total bridges: \\(n - 1\\).\nSubgoal 2: Substitute parameters. A path on \\(n = %d\\) vertices has %d total edges.\nSubgoal 3: Evaluate and simplify. Since every edge in a path is a bridge, the total number of bridges is %d - 1 = %d.", n, n-1, n, ans)

	case 1:
		kind = "cycle"
		n := 3 + rand.Intn(4) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..11
		ans = 0
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many bridges does the cycle graph C_%d have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. By Whitney's bridge theorem, an edge is a bridge if and only if it belongs to no cycle. In a cycle graph \\(C_n\\), every edge lies on the cycle itself.\nSubgoal 2: Substitute parameters. For \\(C_%d\\), every edge is part of a cycle.\nSubgoal 3: Evaluate and simplify. The cycle graph contains 0 bridges.", n)

	case 2:
		kind = "star"
		n := 3 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..10
		ans = n
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many bridges does the star graph K_{1,%d} have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A star graph \\(K_{1,k}\\) is a tree with \\(k\\) leaves connected to a central vertex. Removing any leaf edge isolates that leaf.\nSubgoal 2: Substitute parameters. In \\(K_{1,%d}\\), there are %d leaf edges incident to the central vertex.\nSubgoal 3: Evaluate and simplify. Since every leaf edge is a bridge, there are %d bridges.", n, n, ans)

	case 3:
		kind = "tree"
		n := 4 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 6) // 4..12
		ans = n - 1
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("How many bridges does a tree on %d vertices have?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A tree on \\(n\\) vertices is a connected acyclic graph containing exactly \\(n - 1\\) edges. Because it contains no cycles, every edge in a tree is a bridge.\nSubgoal 2: Substitute parameters. For \\(n = %d\\) vertices, the tree contains \\(%d - 1\\) edges.\nSubgoal 3: Evaluate and simplify. All %d edges are bridges: %d - 1 = %d.", n, n, n-1, n, ans)

	case 4:
		kind = "join"
		ans = 1
		prompt = "How many bridges does a graph made of two complete graphs K_3 joined by a single edge have?"
		expl = "Subgoal 1: State governing formula. An edge is a bridge if and only if its removal disconnects the graph (or equivalently, if it belongs to no cycle).\nSubgoal 2: Substitute parameters. Within each \\(K_3\\), every edge belongs to a 3-cycle and cannot be a bridge. The single edge connecting the two cliques belongs to no cycle and connects two disjoint components.\nSubgoal 3: Evaluate and simplify. Only the connecting edge is a bridge, so the graph has exactly 1 bridge."

	case 5:
		kind = "chord"
		ans = 0
		prompt = "How many bridges does a square cycle graph with an added diagonal chord have?"
		expl = "Subgoal 1: State governing formula. An edge is a bridge if and only if it belongs to no cycle.\nSubgoal 2: Substitute parameters. Adding a diagonal chord to a 4-cycle creates two 3-cycles sharing the chord. Every edge in the graph belongs to at least one 3-cycle or the outer 4-cycle.\nSubgoal 3: Evaluate and simplify. Since every edge belongs to a cycle, the graph contains 0 bridges."
	}

	facts["kind"] = kind
	facts["answer"] = strconv.Itoa(ans)

	return generator.Problem{
		Question:    prompt,
		Answer:      strconv.Itoa(ans),
		Explanation: expl,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// 9. Vertex and edge connectivity (numeric)
// ---------------------------------------------------------------------------

type connectivityGen struct{}

// Generate varies the graph and the invariant. The first ordinal asks for the
// vertex connectivity of a cycle, a path or a disconnected graph, the second for
// the complete and complete bipartite graphs the section works through, and the
// third for edge connectivity, where the chain kappa <= lambda <= delta is the
// subject.
func (g *connectivityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	switch o {
	case 1:
		switch rand.Intn(3) {
		case 0:
			n := graphRange(3, 14)
			kind, ans = "cycle", 2
			facts["n"] = strconv.Itoa(n)
			question = fmt.Sprintf("What is the vertex connectivity kappa(C_%d)? (enter a number)", n)
			reason = fmt.Sprintf("Deleting one vertex of C_%d leaves the connected path P_%d, while deleting two vertices that are not adjacent disconnects the cycle, so kappa(C_%d) = %d.", n, n-1, n, ans)
		case 1:
			n := graphRange(2, 14)
			kind, ans = "path", 1
			facts["n"] = strconv.Itoa(n)
			question = fmt.Sprintf("What is the vertex connectivity kappa(P_%d)? (enter a number)", n)
			reason = fmt.Sprintf("Every internal vertex of P_%d is a cut vertex, so one deletion disconnects the path and kappa(P_%d) = %d.", n, n, ans)
		default:
			k := graphRange(2, 5)
			n := 3*k + graphRange(0, 9)
			kind, ans = "disjoint", 0
			facts["n"] = strconv.Itoa(n)
			facts["components"] = strconv.Itoa(k)
			question = fmt.Sprintf("A graph on %d vertices is the disjoint union of %d cycles. What is its vertex connectivity? (enter a number)", n, k)
			reason = fmt.Sprintf("The graph already has %d components, and deleting vertices cannot rejoin them, so the vertex connectivity is %d.", k, ans)
		}
	case 2:
		if rand.Intn(2) == 0 {
			n := graphRange(3, 13)
			kind, ans = "complete", n-1
			facts["a"] = strconv.Itoa(n)
			question = fmt.Sprintf("What is the vertex connectivity kappa(K_%d)? (enter a number)", n)
			reason = fmt.Sprintf("Deleting any %d vertices of K_%d leaves a single vertex, and fewer deletions leave a complete graph on at least two vertices, so kappa(K_%d) = %d - 1 = %d.", n-1, n, n, n, ans)
		} else {
			a := graphRange(2, 6)
			b := graphRange(2, 7)
			m := graphMin(a, b)
			kind, ans = "bipartite", m
			facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
			question = fmt.Sprintf("What is the vertex connectivity kappa(K_{%d,%d})? (enter a number)", a, b)
			reason = fmt.Sprintf("Deleting the whole side of size %d isolates every vertex on the other side, and no smaller set disconnects K_{%d,%d}, so kappa = %d.", m, a, b, ans)
		}
	default:
		switch rand.Intn(3) {
		case 0:
			n := graphRange(3, 13)
			kind, ans = "complete", n-1
			facts["a"] = strconv.Itoa(n)
			question = fmt.Sprintf("What is the edge connectivity lambda(K_%d)? (enter a number)", n)
			reason = fmt.Sprintf("All %d edges at a vertex must be cut to isolate it, and kappa <= lambda <= delta forces equality at %d for K_%d, so lambda(K_%d) = %d.", n-1, n-1, n, n, ans)
		case 1:
			n := graphRange(3, 15)
			kind, ans = "path", 1
			facts["n"] = strconv.Itoa(n)
			question = fmt.Sprintf("What is the edge connectivity lambda of a tree on %d vertices? (enter a number)", n)
			reason = "Every edge of a tree is a bridge, so cutting one edge disconnects it and lambda = 1."
		default:
			a := graphRange(2, 6)
			b := graphRange(2, 7)
			m := graphMin(a, b)
			kind, ans = "bipartite", m
			facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
			question = fmt.Sprintf("What is the edge connectivity lambda(K_{%d,%d})? (enter a number)", a, b)
			reason = fmt.Sprintf("K_{%d,%d} has minimum degree %d and the minimum degree is an upper bound on lambda, while the %d edges at one vertex of the smaller side form a cut, so lambda = %d.", a, b, m, m, ans)
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// 10. Menger's theorem
// ---------------------------------------------------------------------------

type mengerGen struct{}

func (g *mengerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return mengerProblem()
}

// ---------------------------------------------------------------------------
// 11. Matchings (numeric)
// ---------------------------------------------------------------------------

type matchingGen struct{}

func (g *matchingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(3)
	var prompt, expl, kind string
	var ans int
	facts := map[string]string{}

	switch mode {
	case 0:
		kind = "path"
		n := 3 + rand.Intn(4) + mathutil.ScaleInt(ctx.Difficulty, 0, 6) // 3..12
		ans = n / 2
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("What is the maximum matching size in the path graph P_%d?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A matching is a set of pairwise non-adjacent edges. In a path graph \\(P_n\\), each matched edge covers 2 vertices, and alternating edges yield the maximum matching: \\(\\nu(P_n) = \\lfloor n / 2 \\rfloor\\).\nSubgoal 2: Substitute parameters. For \\(n = %d\\), compute \\(\\lfloor %d / 2 \\rfloor\\).\nSubgoal 3: Evaluate and simplify. Evaluate integer division: %d / 2 = %d.", n, n, n, ans)

	case 1:
		kind = "cycle"
		n := 3 + rand.Intn(4) + mathutil.ScaleInt(ctx.Difficulty, 0, 6) // 3..12
		ans = n / 2
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("What is the maximum matching size in the cycle graph C_%d?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. In a cycle graph \\(C_n\\), choosing alternating edges around the perimeter yields the maximum matching size: \\(\\nu(C_n) = \\lfloor n / 2 \\rfloor\\).\nSubgoal 2: Substitute parameters. For \\(n = %d\\), compute \\(\\lfloor %d / 2 \\rfloor\\).\nSubgoal 3: Evaluate and simplify. Evaluate integer division: %d / 2 = %d.", n, n, n, ans)

	case 2:
		kind = "star"
		n := 3 + rand.Intn(3) + mathutil.ScaleInt(ctx.Difficulty, 0, 5) // 3..10
		ans = 1
		facts["n"] = strconv.Itoa(n)
		prompt = fmt.Sprintf("What is the maximum matching size in the star graph K_{1,%d}?", n)
		expl = fmt.Sprintf("Subgoal 1: State governing formula. A matching requires edges that do not share any common vertices. In a star graph \\(K_{1,k}\\), every edge is incident to the single central vertex.\nSubgoal 2: Substitute parameters. Since any two edges share the central vertex, no two edges can be included in the same matching.\nSubgoal 3: Evaluate and simplify. The maximum matching size of \\(K_{1,%d}\\) is 1.", n)
	}

	facts["kind"] = kind
	facts["answer"] = strconv.Itoa(ans)

	return generator.Problem{
		Question:    prompt,
		Answer:      strconv.Itoa(ans),
		Explanation: expl,
		Facts:       facts,
	}
}

// ---------------------------------------------------------------------------
// 12. Perfect matchings
// ---------------------------------------------------------------------------

type perfectMatchingGen struct{}

func (g *perfectMatchingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return perfectMatchingProblem()
}

// ---------------------------------------------------------------------------
// 13. Tutte's theorem
// ---------------------------------------------------------------------------

type tutteGen struct{}

func (g *tutteGen) Generate(ctx generator.GeneratorContext) generator.Problem { return tutteProblem() }

// ---------------------------------------------------------------------------
// 14. Vertex covers (numeric)
// ---------------------------------------------------------------------------

type vertexCoverGen struct{}

// Generate draws the graph the cover belongs to: a path or a cycle at the first
// ordinal, the worked paths, stars and cycles at the second, and complete
// bipartite graphs, where Koenig's theorem gives the smaller side, at the third.
func (g *vertexCoverGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", n/2
		facts["a"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the minimum vertex cover size of the path P_%d? (enter a number)", n)
		reason = fmt.Sprintf("Choosing vertices 2, 4 and so on covers every edge of P_%d, and each chosen vertex covers two edges, so the minimum is %d / 2 = %d.", n, n, ans)
	}
	cycleClaim := func() {
		n := graphRange(4, 15)
		// A cover is the complement of an independent set, and the largest
		// independent set of C_n has floor(n/2) vertices, so the smallest cover
		// has n - floor(n/2) = ceil(n/2) vertices. On an odd cycle that is one
		// more than half.
		kind, ans = "cycle", (n+1)/2
		facts["a"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the minimum vertex cover size of the cycle C_%d? (enter a number)", n)
		reason = fmt.Sprintf("A cover is the complement of an independent set, and the largest independent set of C_%d has floor(%d / 2) = %d vertices, so the smallest cover has %d - %d = %d.", n, n, n/2, n, n/2, ans)
	}
	starClaim := func() {
		k := graphRange(3, 14)
		kind, ans = "star", 1
		facts["a"] = strconv.Itoa(k)
		question = fmt.Sprintf("What is the minimum vertex cover size of the star K_{1,%d}? (enter a number)", k)
		reason = fmt.Sprintf("The centre meets all %d leaf edges, and no cover can be smaller than one vertex, so the minimum is 1.", k)
	}
	bipartiteClaim := func() {
		a := graphRange(2, 7)
		b := graphRange(2, 8)
		m := graphMin(a, b)
		kind, ans = "bipartite", m
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question = fmt.Sprintf("What is the minimum vertex cover size of K_{%d,%d}? (enter a number)", a, b)
		reason = fmt.Sprintf("Koenig's theorem equates the cover with the maximum matching, and a matching of K_{%d,%d} saturates the smaller side, so the minimum is %d.", a, b, m)
	}

	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			pathClaim()
		} else {
			cycleClaim()
		}
	case 2:
		switch rand.Intn(3) {
		case 0:
			pathClaim()
		case 1:
			starClaim()
		default:
			cycleClaim()
		}
	default:
		if rand.Intn(3) == 0 {
			pathClaim()
		} else {
			bipartiteClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// 15. Independent sets (numeric)
// ---------------------------------------------------------------------------

type independentSetGen struct{}

// Generate draws the graph whose independence number is asked for: a path or a
// cycle at the first ordinal, the worked path, star and complete bipartite
// examples at the second, and the star, cycle and bipartite families that
// complement a cover at the third.
func (g *independentSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", (n+1)/2
		facts["a"] = strconv.Itoa(n)
		question = fmt.Sprintf("An independent set has no two adjacent vertices. What is the maximum independent set size of the path P_%d? (enter a number)", n)
		reason = fmt.Sprintf("Vertices 1, 3, 5 and so on are pairwise non-adjacent, and no independent set of P_%d is larger, so the maximum is ceil(%d / 2) = %d.", n, n, ans)
	}
	cycleClaim := func() {
		n := graphRange(4, 15)
		kind, ans = "cycle", n/2
		facts["a"] = strconv.Itoa(n)
		question = fmt.Sprintf("An independent set has no two adjacent vertices. What is the maximum independent set size of the cycle C_%d? (enter a number)", n)
		reason = fmt.Sprintf("No independent set of C_%d can hold two adjacent vertices, and alternating around the ring gives floor(%d / 2) = %d.", n, n, ans)
	}
	starClaim := func() {
		k := graphRange(3, 14)
		kind, ans = "star", k
		facts["a"] = strconv.Itoa(k)
		question = fmt.Sprintf("An independent set has no two adjacent vertices. What is the maximum independent set size of the star K_{1,%d}? (enter a number)", k)
		reason = fmt.Sprintf("All %d leaves are pairwise non-adjacent, so they form an independent set of size %d, and the centre cannot join them.", k, ans)
	}
	bipartiteClaim := func() {
		a := graphRange(2, 7)
		b := graphRange(2, 8)
		m := a
		if b > m {
			m = b
		}
		kind, ans = "bipartite", m
		facts["a"], facts["b"] = strconv.Itoa(a), strconv.Itoa(b)
		question = fmt.Sprintf("An independent set has no two adjacent vertices. What is the maximum independent set size of K_{%d,%d}? (enter a number)", a, b)
		reason = fmt.Sprintf("Every vertex of one side is non-adjacent to the rest of its own side, so the larger side of size %d is independent and no larger set is.", m)
	}

	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			pathClaim()
		} else {
			cycleClaim()
		}
	case 2:
		switch rand.Intn(3) {
		case 0:
			pathClaim()
		case 1:
			starClaim()
		default:
			bipartiteClaim()
		}
	default:
		switch rand.Intn(3) {
		case 0:
			starClaim()
		case 1:
			cycleClaim()
		default:
			bipartiteClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// 16. The chromatic polynomial (numeric)
// ---------------------------------------------------------------------------

type chromaticPolynomialGen struct{}

// Generate evaluates the chromatic polynomial at a drawn k for a drawn graph:
// a path at the first ordinal where the counting argument is the subject, the
// worked paths, cycles and complete graphs at the second, and complete graphs
// and trees at the third, where the polynomial is evaluated and factored.
func (g *chromaticPolynomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	pathClaim := func() {
		n := graphRange(3, 10)
		k := graphRange(2, 5)
		kind, ans = "path", k*graphIntPow(k-1, n-1)
		facts["n"], facts["k"] = strconv.Itoa(n), strconv.Itoa(k)
		question = fmt.Sprintf("What is P(P_%d, %d), the chromatic polynomial of the path P_%d evaluated at %d? (enter a number)", n, k, n, k)
		reason = fmt.Sprintf("A path is coloured from one end: the first vertex has %d choices and each of the other %d vertices has %d once its predecessor is fixed, so P(P_%d, %d) = %d * %d^%d = %d.", k, n-1, k-1, n, k, k, k-1, n-1, ans)
	}
	cycleClaim := func() {
		n := graphRange(3, 9)
		k := graphRange(2, 4)
		sign := 1
		if n%2 == 1 {
			sign = -1
		}
		op := "+"
		if sign < 0 {
			op = "-"
		}
		kind, ans = "cycle", graphIntPow(k-1, n)+sign*(k-1)
		facts["n"], facts["k"] = strconv.Itoa(n), strconv.Itoa(k)
		question = fmt.Sprintf("What is P(C_%d, %d), the chromatic polynomial of the cycle C_%d evaluated at %d? (enter a number)", n, k, n, k)
		reason = fmt.Sprintf("P(C_n, k) = (k-1)^n + (-1)^n (k-1), so P(C_%d, %d) = %d^%d %s %d = %d.", n, k, k-1, n, op, k-1, ans)
	}
	completeClaim := func() {
		n := graphRange(3, 6)
		k := graphRange(n, n+4)
		term := 1
		parts := make([]string, 0, n)
		for j := 0; j < n; j++ {
			term *= k - j
			parts = append(parts, strconv.Itoa(k-j))
		}
		kind, ans = "complete", term
		facts["n"], facts["k"] = strconv.Itoa(n), strconv.Itoa(k)
		question = fmt.Sprintf("What is P(K_%d, %d), the chromatic polynomial of K_%d evaluated at %d? (enter a number)", n, k, n, k)
		reason = fmt.Sprintf("Every pair of vertices of K_%d is adjacent, so the colours must all differ: %s = %d.", n, strings.Join(parts, " * "), ans)
	}
	treeClaim := func() {
		n := graphRange(3, 14)
		k := graphRange(2, 4)
		kind, ans = "tree", k*graphIntPow(k-1, n-1)
		facts["n"], facts["k"] = strconv.Itoa(n), strconv.Itoa(k)
		question = fmt.Sprintf("What is P(T_%d, %d), the chromatic polynomial of a tree on %d vertices evaluated at %d? (enter a number)", n, k, n, k)
		reason = fmt.Sprintf("A tree on %d vertices has %d edges and no cycle, so the first vertex has %d choices and each of the other %d vertices has %d: %d * %d^%d = %d.", n, n-1, k, n-1, k-1, k, k-1, n-1, ans)
	}

	switch o {
	case 1:
		pathClaim()
	case 2:
		if rand.Intn(2) == 0 {
			cycleClaim()
		} else {
			completeClaim()
		}
	default:
		if rand.Intn(2) == 0 {
			completeClaim()
		} else {
			treeClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// 17. Brooks' theorem
// ---------------------------------------------------------------------------

type brooksGen struct{}

func (g *brooksGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return brooksProblem()
}

// ---------------------------------------------------------------------------
// 18. Edge coloring (numeric)
// ---------------------------------------------------------------------------

type edgeColoringGen struct{}

// Generate draws the graph whose edge chromatic number is asked for: a cycle or
// a path at the first ordinal, the worked complete graphs and cycles at the
// second, and the stars and paths of Vizing's classification at the third.
func (g *edgeColoringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", 2
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the edge chromatic number of the path P_%d? (enter a number)", n)
		reason = fmt.Sprintf("Consecutive edges of P_%d share a vertex and must differ, and alternating two colours does that, so the edge chromatic number is 2.", n)
	}
	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", 3
		if n%2 == 0 {
			ans = 2
			reason = fmt.Sprintf("C_%d is even, so two colours alternate around it and one colour cannot colour two adjacent edges: the edge chromatic number is 2.", n)
		} else {
			reason = fmt.Sprintf("C_%d is odd, so two colours cannot alternate all the way round and three are needed: the edge chromatic number is 3.", n)
		}
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the edge chromatic number of the cycle C_%d? (enter a number)", n)
	}
	completeClaim := func() {
		n := graphRange(3, 14)
		kind = "complete"
		if n%2 == 0 {
			ans = n - 1
			reason = fmt.Sprintf("K_%d with %d even is 1-factorable: its %d edges split into %d perfect matchings, so the edge chromatic number is %d - 1 = %d.", n, n, n*(n-1)/2, n-1, n, ans)
		} else {
			ans = n
			reason = fmt.Sprintf("K_%d with %d odd has %d edges, and each colour class is a matching of at most %d edges, so %d colours cover at most %d edges and %d are needed; %d suffice, so the answer is %d.", n, n, n*(n-1)/2, (n-1)/2, n-1, (n-1)*(n-1)/2, n, n, ans)
		}
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the edge chromatic number of K_%d? (enter a number)", n)
	}
	starClaim := func() {
		k := graphRange(2, 14)
		kind, ans = "star", k
		facts["n"] = strconv.Itoa(k)
		question = fmt.Sprintf("What is the edge chromatic number of the star K_{1,%d}? (enter a number)", k)
		reason = fmt.Sprintf("All %d edges meet at the centre, so no two of them may share a colour and the edge chromatic number is %d.", k, ans)
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
			starClaim()
		} else {
			pathClaim()
		}
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}

// ---------------------------------------------------------------------------
// 19. Perfect graphs
// ---------------------------------------------------------------------------

type perfectGraphGen struct{}

func (g *perfectGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return perfectGraphProblem()
}

// ---------------------------------------------------------------------------
// 20. List coloring (numeric)
// ---------------------------------------------------------------------------

type listColoringGen struct{}

// Generate draws the graph whose list chromatic number is asked for: paths and
// cycles at the first ordinal, the worked even and odd cycles at the second, and
// paths at the third.
func (g *listColoringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphOrdinal(ctx)
	facts := map[string]string{}
	var kind, question, reason string
	ans := 0

	pathClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "path", 2
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the list chromatic number of the path P_%d? (enter a number)", n)
		reason = fmt.Sprintf("Colour P_%d from one end, always choosing a colour from the current list other than the previous vertex's colour; each list has two colours and only one is forbidden, so two suffice and one does not.", n)
	}
	cycleClaim := func() {
		n := graphRange(3, 15)
		kind, ans = "cycle", 3
		if n%2 == 0 {
			ans = 2
			reason = fmt.Sprintf("C_%d is even and every even cycle is 2-choosable: alternate colours around the ring, and one colour cannot cover two adjacent lists. The list chromatic number is 2.", n)
		} else {
			reason = fmt.Sprintf("C_%d is odd, and an odd cycle is not 2-choosable, so its list chromatic number is 3.", n)
		}
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("What is the list chromatic number of the cycle C_%d? (enter a number)", n)
	}
	evenOddClaim := func() {
		n := 2 * graphRange(2, 7)
		if rand.Intn(2) == 0 {
			n++
		}
		kind = "cycle"
		parity := "odd"
		if n%2 == 0 {
			parity = "even"
			ans = 2
			reason = fmt.Sprintf("C_%d is an even cycle, and every even cycle is 2-choosable, so its list chromatic number is 2.", n)
		} else {
			ans = 3
			reason = fmt.Sprintf("C_%d is an odd cycle, which is not 2-choosable, so its list chromatic number is 3.", n)
		}
		facts["n"] = strconv.Itoa(n)
		question = fmt.Sprintf("An even cycle is 2-choosable and an odd cycle is not. What is the list chromatic number of C_%d, whose length is %s? (enter a number)", n, parity)
	}

	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			pathClaim()
		} else {
			cycleClaim()
		}
	case 2:
		evenOddClaim()
	default:
		pathClaim()
	}

	facts["kind"] = kind
	return graphNumeric(o, question, ans, reason, facts)
}
