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

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

// ---------------------------------------------------------------------------
// 1. Graph isomorphism
// ---------------------------------------------------------------------------

type isomorphismGen struct{}

func (g *isomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []graphClaim{
		{
			question:    "Which quantity is an isomorphism invariant of a graph?",
			correct:     "the multiset of vertex degrees",
			distractors: []string{"the drawing of the graph", "the names of the vertices", "the order in which the edges are listed"},
			reason:      "An isomorphism permutes vertices, so each vertex keeps its degree and the degree multiset is unchanged.",
		},
		{
			question:    "Two graphs are isomorphic when...",
			correct:     "a bijection of the vertices preserves adjacency and non-adjacency",
			distractors: []string{"they have the same number of vertices", "they have the same number of edges", "they are drawn the same way"},
			reason:      "Matching vertex and edge counts are necessary but not sufficient; the bijection must preserve the whole adjacency relation.",
		},
		{
			question:    "Graphs with degree sequences (2,2,2,2) and (3,3,1,1) are...",
			correct:     "not isomorphic",
			distractors: []string{"isomorphic", "both trees", "both disconnected"},
			reason:      "The degree multisets differ, and the degree multiset is an isomorphism invariant.",
		},
		{
			question:    "A graph and its complement can be...",
			correct:     "isomorphic or not, depending on the graph",
			distractors: []string{"never isomorphic", "always isomorphic", "isomorphic only when complete"},
			reason:      "Self-complementary graphs exist (P_4 is the smallest), and many graphs are not self-complementary.",
		},
		{
			question:    "Isomorphism is an equivalence relation because...",
			correct:     "it is reflexive, symmetric and transitive",
			distractors: []string{"it preserves degrees", "it is a bijection", "it is a graph property"},
			reason:      "Identity, inverse and composition of adjacency-preserving bijections give the three relation laws.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.isomorphism")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
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
	claims := []graphClaim{
		{
			question:    "An induced subgraph on a vertex set S contains...",
			correct:     "all edges of G with both ends in S",
			distractors: []string{"all edges of G", "only edges chosen by S", "a spanning tree of G"},
			reason:      "Induction keeps every edge whose two endpoints survive, and removes every edge that touches a deleted vertex.",
		},
		{
			question:    "A spanning subgraph has...",
			correct:     "the same vertex set as G",
			distractors: []string{"the same edge set", "one vertex removed", "no cycles"},
			reason:      "Spanning means every vertex is kept; the edge set is allowed to shrink.",
		},
		{
			question:    "Deleting a vertex from a graph...",
			correct:     "deletes the vertex and every incident edge",
			distractors: []string{"deletes only the vertex", "deletes only the incident edges", "keeps all edges"},
			reason:      "The remaining object must still be a graph, so edges with a missing endpoint cannot stay.",
		},
		{
			question:    "Which graph is a subgraph of every graph?",
			correct:     "the graph on the same vertices with no edges",
			distractors: []string{"a spanning tree", "a Hamiltonian cycle", "the complement"},
			reason:      "The empty-edge graph is contained in every simple graph; the others need not exist.",
		},
		{
			question:    "An induced subgraph is determined by...",
			correct:     "its vertex set",
			distractors: []string{"its edge set", "a spanning tree", "its complement"},
			reason:      "Once the vertex set is fixed, the induced subgraph keeps exactly the edges inside it.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.subgraph")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 4. Graph homomorphisms
// ---------------------------------------------------------------------------

type homomorphismGen struct{}

func (g *homomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []graphClaim{
		{
			question:    "A graph homomorphism f: G -> H maps...",
			correct:     "every edge of G to an edge of H",
			distractors: []string{"every edge of H to an edge of G", "vertices to non-adjacent vertices", "only the degrees"},
			reason:      "The defining condition is that adjacency in G implies adjacency of the images in H.",
		},
		{
			question:    "A proper k-coloring is a homomorphism to...",
			correct:     "K_k",
			distractors: []string{"the path P_k", "the cycle C_k", "the complement of K_k"},
			reason:      "The color classes are the preimages of the k vertices of K_k, and adjacent vertices must get different colors.",
		},
		{
			question:    "A graph isomorphism is a homomorphism that is...",
			correct:     "bijective with a homomorphism inverse",
			distractors: []string{"injective only", "surjective only", "degree-preserving only"},
			reason:      "An isomorphism must preserve adjacency in both directions, which is exactly a bijective homomorphism whose inverse is a homomorphism.",
		},
		{
			question:    "A homomorphism to K_2 is...",
			correct:     "a proper 2-coloring",
			distractors: []string{"a spanning tree", "an Eulerian trail", "a matching"},
			reason:      "Two colors, and every edge forces different colors: the definition of a proper 2-coloring.",
		},
		{
			question:    "Homomorphisms are more flexible than isomorphisms because...",
			correct:     "they need not be injective",
			distractors: []string{"they preserve non-adjacency", "they preserve degrees", "they are always bijections"},
			reason:      "Several vertices may map to the same target, which is what allows a coloring to collapse a large graph onto a small one.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.homomorphism")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
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

func (g *connectivityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		a, b     int
		reason   string
	}
	claims := []claim{
		{question: "What is the vertex connectivity kappa(K_5)? (enter a number)", answer: "4", kind: "complete", a: 5, reason: "Removing any four vertices leaves one vertex, and removing fewer cannot disconnect K_5."},
		{question: "What is the vertex connectivity kappa(C_6)? (enter a number)", answer: "2", kind: "cycle", a: 6, reason: "Two vertices can disconnect a cycle, and no single vertex can."},
		{question: "What is the vertex connectivity kappa(P_5)? (enter a number)", answer: "1", kind: "path", a: 5, reason: "Any internal vertex disconnects a path."},
		{question: "What is the vertex connectivity kappa(K_{3,3})? (enter a number)", answer: "3", kind: "bipartite", a: 3, b: 3, reason: "Removing one side disconnects the graph, and no smaller set does."},
		{question: "What is the vertex connectivity kappa(K_{2,4})? (enter a number)", answer: "2", kind: "bipartite", a: 2, b: 4, reason: "Removing the two vertices of the smaller side disconnects the graph."},
		{question: "What is the vertex connectivity kappa of two disjoint cycles? (enter a number)", answer: "0", kind: "disjoint", reason: "A disconnected graph has vertex connectivity 0."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.a > 0 {
		facts["a"] = strconv.Itoa(c.a)
	}
	if c.b > 0 {
		facts["b"] = strconv.Itoa(c.b)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// 10. Menger's theorem
// ---------------------------------------------------------------------------

type mengerGen struct{}

func (g *mengerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []graphClaim{
		{
			question:    "Menger's theorem equates the maximum number of internally vertex-disjoint u-v paths with...",
			correct:     "the minimum size of a vertex cut separating u and v",
			distractors: []string{"the maximum degree", "the number of edges", "the diameter"},
			reason:      "The theorem is a min-max duality: the largest family of disjoint paths is limited by, and equals, the smallest separating cut.",
		},
		{
			question:    "The edge version of Menger's theorem uses...",
			correct:     "edge-disjoint paths and edge cuts",
			distractors: []string{"vertex cuts only", "matchings", "spanning trees"},
			reason:      "Replacing internally disjoint paths by edge-disjoint paths and vertex cuts by edge cuts gives the edge form.",
		},
		{
			question:    "Menger's theorem implies which inequality?",
			correct:     "kappa(G) <= lambda(G) <= delta(G)",
			distractors: []string{"lambda(G) <= kappa(G) <= delta(G)", "delta(G) <= kappa(G) <= lambda(G)", "kappa(G) = delta(G)"},
			reason:      "A vertex cut can be made from the edges at a vertex, and an edge cut can be made from the edges at a minimum-degree vertex.",
		},
		{
			question:    "A graph is k-connected exactly when...",
			correct:     "every pair of vertices has k internally disjoint paths",
			distractors: []string{"every vertex has degree at least k", "it has k edges", "it has no cut vertices"},
			reason:      "That is the path form of k-connectivity, equivalent to the minimum vertex cut being at least k.",
		},
		{
			question:    "Menger's theorem is called a min-max theorem because...",
			correct:     "a maximum packing equals a minimum cut",
			distractors: []string{"it counts vertices", "it uses Euler's formula", "it applies only to trees"},
			reason:      "Packing paths and cutting them are dual optimisation problems with the same optimum.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.menger")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
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
	claims := []graphClaim{
		{
			question:    "A perfect matching covers...",
			correct:     "every vertex exactly once",
			distractors: []string{"every edge exactly once", "every vertex at least twice", "only the leaves"},
			reason:      "A matching uses disjoint edges; perfect means the disjoint edges touch every vertex.",
		},
		{
			question:    "Which graph cannot have a perfect matching?",
			correct:     "a path with an odd number of vertices",
			distractors: []string{"a cycle with an even number of vertices", "a complete graph with an even number of vertices", "a star with an even number of leaves"},
			reason:      "Every perfect matching covers vertices in pairs, so an odd vertex count is impossible.",
		},
		{
			question:    "Every regular bipartite graph...",
			correct:     "has a perfect matching",
			distractors: []string{"has no perfect matching", "is a tree", "is complete"},
			reason:      "Hall's condition follows from regularity, so a perfect matching exists.",
		},
		{
			question:    "A graph with a perfect matching must have...",
			correct:     "an even number of vertices",
			distractors: []string{"an odd number of vertices", "a cut vertex", "a bridge"},
			reason:      "The matching pairs the vertices, so their number is twice the number of matched edges.",
		},
		{
			question:    "A perfect matching is also...",
			correct:     "a maximum matching",
			distractors: []string{"a minimum vertex cover", "a spanning tree", "an Eulerian circuit"},
			reason:      "No matching can use more edges than one that already covers every vertex.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.perfect_matching")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 13. Tutte's theorem
// ---------------------------------------------------------------------------

type tutteGen struct{}

func (g *tutteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Tutte's theorem characterizes...",
			correct:     "graphs with a perfect matching",
			distractors: []string{"graphs with an Eulerian circuit", "planar graphs", "bipartite graphs"},
			reason:      "Tutte gives a necessary and sufficient condition for the existence of a perfect matching.",
		},
		{
			question:    "Tutte's condition counts...",
			correct:     "odd components of G - S for every vertex set S",
			distractors: []string{"even components", "edges", "cycles"},
			reason:      "The condition is that every S leaves at most |S| odd components.",
		},
		{
			question:    "If |V| is odd, Tutte fails at S = empty because...",
			correct:     "the single component is odd and larger than |S|",
			distractors: []string{"there are no components", "the graph is disconnected", "the matching is perfect"},
			reason:      "An odd number of vertices cannot be perfectly matched, and the condition detects it with the empty set.",
		},
		{
			question:    "Tutte's theorem generalizes...",
			correct:     "Hall's marriage theorem",
			distractors: []string{"Euler's formula", "Brooks' theorem", "Menger's theorem"},
			reason:      "Hall handles bipartite graphs; Tutte removes the bipartite restriction.",
		},
		{
			question:    "A bridge in a graph with a perfect matching...",
			correct:     "can exist, but each side must contain an odd number of vertices",
			distractors: []string{"cannot exist", "forces the graph to be disconnected", "is a perfect matching"},
			reason:      "The matching must use the bridge itself, and then each side has one unmatched vertex.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 14. Vertex covers (numeric)
// ---------------------------------------------------------------------------

type vertexCoverGen struct{}

func (g *vertexCoverGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		a, b     int
		reason   string
	}
	claims := []claim{
		{question: "What is the minimum vertex cover size of the path P_4? (enter a number)", answer: "2", kind: "path", a: 4, reason: "The two internal vertices cover every edge."},
		{question: "What is the minimum vertex cover size of the cycle C_6? (enter a number)", answer: "3", kind: "cycle", a: 6, reason: "Alternate vertices cover all six edges, and two cannot."},
		{question: "What is the minimum vertex cover size of the star K_{1,5}? (enter a number)", answer: "1", kind: "star", a: 5, reason: "The centre covers every edge."},
		{question: "What is the minimum vertex cover size of K_{2,3}? (enter a number)", answer: "2", kind: "bipartite", a: 2, b: 3, reason: "Koenig's theorem gives the smaller side, 2."},
		{question: "What is the minimum vertex cover size of the path P_5? (enter a number)", answer: "2", kind: "path", a: 5, reason: "Two well-placed internal vertices cover all four edges."},
		{question: "What is the minimum vertex cover size of K_{3,3}? (enter a number)", answer: "3", kind: "bipartite", a: 3, b: 3, reason: "Koenig's theorem gives 3, one entire side."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.a > 0 {
		facts["a"] = strconv.Itoa(c.a)
	}
	if c.b > 0 {
		facts["b"] = strconv.Itoa(c.b)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// 15. Independent sets (numeric)
// ---------------------------------------------------------------------------

type independentSetGen struct{}

func (g *independentSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		a, b     int
		reason   string
	}
	claims := []claim{
		{question: "What is the maximum independent set size of the path P_4? (enter a number)", answer: "2", kind: "path", a: 4, reason: "The two endpoints and the two internal vertices split into two independent pairs."},
		{question: "What is the maximum independent set size of the cycle C_6? (enter a number)", answer: "3", kind: "cycle", a: 6, reason: "Alternate vertices are independent and there are three."},
		{question: "What is the maximum independent set size of the star K_{1,5}? (enter a number)", answer: "5", kind: "star", a: 5, reason: "All five leaves are pairwise non-adjacent."},
		{question: "What is the maximum independent set size of K_{2,3}? (enter a number)", answer: "3", kind: "bipartite", a: 2, b: 3, reason: "The larger side of a complete bipartite graph is independent."},
		{question: "What is the maximum independent set size of the path P_5? (enter a number)", answer: "3", kind: "path", a: 5, reason: "Vertices 1, 3 and 5 are pairwise non-adjacent."},
		{question: "What is the maximum independent set size of the cycle C_5? (enter a number)", answer: "2", kind: "cycle", a: 5, reason: "Two non-adjacent vertices is the maximum in an odd cycle."},
	}
	c := claims[rand.Intn(len(claims))]
	facts := map[string]string{"kind": c.kind, "answer": c.answer}
	if c.a > 0 {
		facts["a"] = strconv.Itoa(c.a)
	}
	if c.b > 0 {
		facts["b"] = strconv.Itoa(c.b)
	}
	return generator.Problem{Question: c.question, Answer: c.answer, Explanation: c.reason, Facts: facts}
}

// ---------------------------------------------------------------------------
// 16. The chromatic polynomial (numeric)
// ---------------------------------------------------------------------------

type chromaticPolynomialGen struct{}

func (g *chromaticPolynomialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n, k     int
		reason   string
	}
	claims := []claim{
		{question: "What is P(P_4, 3), the chromatic polynomial of the path P_4 evaluated at 3? (enter a number)", answer: "24", kind: "path", n: 4, k: 3, reason: "P(P_n,k) = k(k-1)^(n-1), so 3*2^3 = 24."},
		{question: "What is P(C_4, 3), the chromatic polynomial of the cycle C_4 evaluated at 3? (enter a number)", answer: "18", kind: "cycle", n: 4, k: 3, reason: "P(C_n,k) = (k-1)^n + (-1)^n(k-1), so 2^4 + 2 = 18."},
		{question: "What is P(K_3, 4), the chromatic polynomial of K_3 evaluated at 4? (enter a number)", answer: "24", kind: "complete", n: 3, k: 4, reason: "The colors must be distinct, so 4*3*2 = 24."},
		{question: "What is P(T_5, 2), the chromatic polynomial of a tree on 5 vertices evaluated at 2? (enter a number)", answer: "2", kind: "tree", n: 5, k: 2, reason: "A tree is bipartite, so with two colors the first vertex has 2 choices and each later vertex has 1."},
		{question: "What is P(C_3, 3), the chromatic polynomial of the triangle evaluated at 3? (enter a number)", answer: "6", kind: "cycle", n: 3, k: 3, reason: "P(C_3,k) = (k-1)^3 - (k-1), so 8 - 2 = 6."},
		{question: "What is P(P_3, 2), the chromatic polynomial of the path P_3 evaluated at 2? (enter a number)", answer: "2", kind: "path", n: 3, k: 2, reason: "P(P_3,k) = k(k-1)^2, so 2*1 = 2."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"kind": c.kind, "n": strconv.Itoa(c.n), "k": strconv.Itoa(c.k), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 17. Brooks' theorem
// ---------------------------------------------------------------------------

type brooksGen struct{}

func (g *brooksGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []graphClaim{
		{
			question:    "Brooks' theorem bounds the chromatic number by...",
			correct:     "the maximum degree, except for complete graphs and odd cycles",
			distractors: []string{"the minimum degree", "the number of edges", "the clique number"},
			reason:      "The greedy bound is Delta + 1; Brooks removes the extra color except for the two exceptional families.",
		},
		{
			question:    "The exceptions in Brooks' theorem are...",
			correct:     "complete graphs and odd cycles",
			distractors: []string{"trees and paths", "bipartite graphs", "planar graphs"},
			reason:      "K_n needs n = Delta + 1 colors and odd cycles need 3 = Delta + 1.",
		},
		{
			question:    "A connected graph with maximum degree 3 has chromatic number at most...",
			correct:     "3 unless it is K_4",
			distractors: []string{"2", "4 always", "5"},
			reason:      "Brooks gives at most Delta = 3 except when the graph is K_4, which needs 4.",
		},
		{
			question:    "Brooks strengthens which bound?",
			correct:     "the greedy bound chi <= Delta + 1",
			distractors: []string{"chi <= omega", "chi <= n", "chi <= alpha"},
			reason:      "Greedy coloring already gives Delta + 1; Brooks lowers it to Delta off the exceptions.",
		},
		{
			question:    "A connected graph with maximum degree 2 is...",
			correct:     "a path or a cycle",
			distractors: []string{"complete", "a tree with a branch", "bipartite but not a path"},
			reason:      "Every vertex has at most two neighbours, so the component is a path or a cycle.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.brooks")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 18. Edge coloring (numeric)
// ---------------------------------------------------------------------------

type edgeColoringGen struct{}

func (g *edgeColoringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n        int
		reason   string
	}
	claims := []claim{
		{question: "What is the edge chromatic number of K_4? (enter a number)", answer: "3", kind: "complete", n: 4, reason: "K_n with n even needs n - 1 = 3 colors."},
		{question: "What is the edge chromatic number of K_5? (enter a number)", answer: "5", kind: "complete", n: 5, reason: "K_n with n odd needs n = 5 colors."},
		{question: "What is the edge chromatic number of C_4? (enter a number)", answer: "2", kind: "cycle", n: 4, reason: "An even cycle alternates two colors."},
		{question: "What is the edge chromatic number of C_5? (enter a number)", answer: "3", kind: "cycle", n: 5, reason: "An odd cycle needs three colors."},
		{question: "What is the edge chromatic number of the path P_4? (enter a number)", answer: "2", kind: "path", n: 4, reason: "Adjacent edges must differ, so two colors suffice and one does not."},
		{question: "What is the edge chromatic number of the star K_{1,4}? (enter a number)", answer: "4", kind: "star", n: 4, reason: "All four edges meet at the centre, so all four colors are needed."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"kind": c.kind, "n": strconv.Itoa(c.n), "answer": c.answer},
	}
}

// ---------------------------------------------------------------------------
// 19. Perfect graphs
// ---------------------------------------------------------------------------

type perfectGraphGen struct{}

func (g *perfectGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []graphClaim{
		{
			question:    "A graph is perfect when...",
			correct:     "omega(H) = chi(H) for every induced subgraph H",
			distractors: []string{"chi(H) <= 2", "it is planar", "it is regular"},
			reason:      "The clique number must equal the chromatic number on the graph and on every induced subgraph.",
		},
		{
			question:    "Which graphs are perfect?",
			correct:     "bipartite graphs",
			distractors: []string{"odd cycles of length at least 5", "all graphs", "only complete graphs"},
			reason:      "Bipartite graphs have clique number at most 2 and chromatic number at most 2, and the same holds for induced subgraphs.",
		},
		{
			question:    "An odd cycle of length at least 5 is...",
			correct:     "not perfect",
			distractors: []string{"perfect", "bipartite", "a tree"},
			reason:      "Its clique number is 2 but its chromatic number is 3, so the equality fails.",
		},
		{
			question:    "The strong perfect graph theorem forbids...",
			correct:     "odd holes and odd antiholes",
			distractors: []string{"even cycles", "trees", "complete graphs"},
			reason:      "A graph is perfect exactly when neither it nor its complement contains an induced odd cycle of length at least 5.",
		},
		{
			question:    "Chordal graphs are...",
			correct:     "perfect",
			distractors: []string{"never perfect", "always bipartite", "always planar"},
			reason:      "Chordal graphs have a perfect elimination ordering, which makes clique number and chromatic number agree on every induced subgraph.",
		},
	}
	claims = append(claims, computedGraphClaims("graph.perfect_graph")...)
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 20. List coloring (numeric)
// ---------------------------------------------------------------------------

type listColoringGen struct{}

func (g *listColoringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		kind     string
		n        int
		reason   string
	}
	claims := []claim{
		{question: "What is the list chromatic number of C_4? (enter a number)", answer: "2", kind: "cycle", n: 4, reason: "Even cycles are 2-choosable."},
		{question: "What is the list chromatic number of C_5? (enter a number)", answer: "3", kind: "cycle", n: 5, reason: "Odd cycles are not 2-choosable and need three colors."},
		{question: "What is the list chromatic number of the path P_5? (enter a number)", answer: "2", kind: "path", n: 5, reason: "A path is 2-choosable: color from one end, avoiding the previous color."},
		{question: "What is the list chromatic number of C_6? (enter a number)", answer: "2", kind: "cycle", n: 6, reason: "Even cycles are 2-choosable."},
		{question: "What is the list chromatic number of C_7? (enter a number)", answer: "3", kind: "cycle", n: 7, reason: "Odd cycles need three colors even with arbitrary lists."},
		{question: "What is the list chromatic number of the path P_3? (enter a number)", answer: "2", kind: "path", n: 3, reason: "A path with at least one edge is 2-choosable."},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts:       map[string]string{"kind": c.kind, "n": strconv.Itoa(c.n), "answer": c.answer},
	}
}
