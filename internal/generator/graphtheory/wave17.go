package graphtheory

// Wave 17 graph theory: extremal bounds, flow duality and the classical
// threshold theorems. Every generator publishes the parameters it drew, so
// wave17_test.go can recompute each answer instead of trusting the arithmetic.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave17(reg *generator.Registry) {
	reg.Register("graph.handshake", &w17Gen{kind: "handshake"})
	reg.Register("graph.dirac", &w17Gen{kind: "dirac"})
	reg.Register("graph.ore", &w17Gen{kind: "ore"})
	reg.Register("graph.erdos_stone", &w17Gen{kind: "erdos_stone"})
	reg.Register("graph.cage", &w17Gen{kind: "cage"})
	reg.Register("graph.zarankiewicz", &w17Gen{kind: "zarankiewicz"})
	reg.Register("graph.critical", &w17Gen{kind: "critical"})
	reg.Register("graph.goodman", &w17Gen{kind: "goodman"})
	reg.Register("graph.friendship", &w17Gen{kind: "friendship"})
	reg.Register("graph.hajnal_szemeredi", &w17Gen{kind: "hajnal"})
	reg.Register("graph.hall", &w17Gen{kind: "hall"})
	reg.Register("graph.konig", &w17Gen{kind: "konig"})
	reg.Register("graph.min_cost_flow", &w17Gen{kind: "min_cost"})
	reg.Register("graph.circulation", &w17Gen{kind: "circulation"})
	reg.Register("graph.gomory_hu", &w17Gen{kind: "gomory_hu"})
	reg.Register("graph.multiway_cut", &w17Gen{kind: "multiway_cut"})
	reg.Register("graph.flow_integrality", &w17Gen{kind: "flow_integrality"})
	reg.Register("graph.edge_disjoint", &w17Gen{kind: "edge_disjoint"})
	reg.Register("graph.dinic", &w17Gen{kind: "dinic"})
	reg.Register("graph.push_relabel", &w17Gen{kind: "push_relabel"})
	registerWave17Rest(reg)
}

type w17Gen struct{ kind string }

func (g *w17Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := w17Ord(ctx)
	switch g.kind {
	case "handshake":
		return w17Handshake(kp)
	case "dirac":
		return w17Dirac(kp)
	case "ore":
		return w17Ore(kp)
	case "erdos_stone":
		return w17ErdosStone(kp)
	case "cage":
		return w17Cage(kp)
	case "zarankiewicz":
		return w17Zarankiewicz(kp)
	case "critical":
		return w17Critical(kp)
	case "goodman":
		return w17Goodman(kp)
	case "friendship":
		return w17Friendship(kp)
	case "hajnal":
		return w17Hajnal(kp)
	case "hall":
		return w17Hall(kp)
	case "konig":
		return w17Konig(kp)
	case "min_cost":
		return w17MinCost(kp)
	case "circulation":
		return w17Circulation(kp)
	case "gomory_hu":
		return w17GomoryHu(kp)
	case "multiway_cut":
		return w17MultiwayCut(kp)
	case "flow_integrality":
		return w17FlowIntegrality(kp)
	case "edge_disjoint":
		return w17EdgeDisjoint(kp)
	case "dinic":
		return w17Dinic(kp)
	case "push_relabel":
		return w17PushRelabel(kp)
	}
	return w17RestGenerate(g.kind, kp)
}

func w17Ord(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

func w17N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w17P(o int, q string, a int, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	s := strconv.Itoa(a)
	facts["answer"] = s
	return generator.Problem{Question: q, Answer: s, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w17I(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w17C2(n int) int { return n * (n - 1) / 2 }

func w17CeilDiv(a, b int) int { return (a + b - 1) / b }

func w17Height(n int) int {
	if n%2 == 0 {
		return n / 2
	}
	return n/2 + 1
}

// ---------------------------------------------------------------------------
// Extremal graph theory
// ---------------------------------------------------------------------------

func w17Handshake(kp int) generator.Problem {
	switch kp {
	case 1:
		m := w17N(6, 30)
		return w17P(kp,
			fmt.Sprintf("A graph has %d edges. What is the sum of the degrees of its vertices?", m),
			2*m,
			fmt.Sprintf("Each edge contributes two to the degree sum, so the total is 2 times %d.", m),
			map[string]string{"m": strconv.Itoa(m)})
	case 2:
		sum := 2 * w17N(8, 30)
		return w17P(kp,
			fmt.Sprintf("A degree list sums to %d. How many edges does the graph have?", sum),
			sum/2,
			fmt.Sprintf("The handshake lemma gives edges equal to half the degree sum, so %d/2.", sum),
			map[string]string{"sum": strconv.Itoa(sum)})
	default:
		n := 2 * w17N(3, 9)
		avg := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("A graph on %d vertices has average degree %d. How many edges does it have?", n, avg),
			n*avg/2,
			fmt.Sprintf("The edge count is half the product of order and average degree, so %d times %d divided by two.", n, avg),
			map[string]string{"n": strconv.Itoa(n), "avg": strconv.Itoa(avg)})
	}
}

func w17Dirac(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 21)
		return w17P(kp,
			fmt.Sprintf("Dirac's theorem applies to a graph on %d vertices. What minimum degree does it require?", n),
			w17Height(n),
			"Dirac requires minimum degree at least half the order, rounded up.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 21)
		deg := w17Height(n)
		total := n * deg
		if total%2 != 0 {
			total++
		}
		return w17P(kp,
			fmt.Sprintf("A graph on %d vertices meets Dirac's minimum degree of %d. What is the smallest possible number of edges?", n, deg),
			total/2,
			"The degree sum is at least the order times the minimum degree, and each edge supplies two to that sum.",
			map[string]string{"n": strconv.Itoa(n), "deg": strconv.Itoa(deg)})
	default:
		k := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("Two disjoint cliques each have %d vertices. What is the minimum degree of their union?", k),
			k-1,
			"Inside a clique of size k every vertex meets the other k-1 vertices, and no edge joins the two cliques.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Ore(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 18)
		return w17P(kp,
			fmt.Sprintf("Ore's theorem is applied to a graph on %d vertices. What degree sum must every nonadjacent pair reach?", n),
			n,
			"Ore requires the two degrees to sum to at least the order.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		d1 := w17N(2, 10)
		d2 := w17N(2, 10)
		return w17P(kp,
			fmt.Sprintf("Two nonadjacent vertices have degrees %d and %d. What is their degree sum?", d1, d2),
			d1+d2,
			"Add the two degrees, since the pair is nonadjacent and no edge already joins them.",
			map[string]string{"d1": strconv.Itoa(d1), "d2": strconv.Itoa(d2)})
	default:
		n := w17N(4, 18)
		return w17P(kp,
			fmt.Sprintf("Dirac's condition on %d vertices is the special case of Ore's condition with which minimum degree?", n),
			w17Height(n),
			"Two degrees each at least half the order sum to at least the order, so Dirac's bound implies Ore's.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17ErdosStone(kp int) generator.Problem {
	switch kp {
	case 1:
		r := w17N(2, 7)
		return w17P(kp,
			fmt.Sprintf("A forbidden graph has chromatic number %d. In the Erdos-Stone density factor the leading fraction is (r-1)/r. What is r-1?", r+1),
			r-1,
			"The density factor for a forbidden graph of chromatic number r+1 is 1 minus 1/r, whose numerator is r-1.",
			map[string]string{"r": strconv.Itoa(r)})
	case 2:
		n := w17N(6, 24)
		r := w17N(2, 5)
		parts := make([]int, r)
		for i := 0; i < n; i++ {
			parts[i%r]++
		}
		edges := 0
		for i := 0; i < r; i++ {
			for j := i + 1; j < r; j++ {
				edges += parts[i] * parts[j]
			}
		}
		return w17P(kp,
			fmt.Sprintf("The balanced complete %d-partite graph has %d vertices. How many edges does it have?", r, n),
			edges,
			"Count every pair of vertices in different parts, which is what the Turán graph keeps.",
			map[string]string{"n": strconv.Itoa(n), "r": strconv.Itoa(r)})
	default:
		r := w17N(2, 7)
		return w17P(kp,
			fmt.Sprintf("How many colours does the complete graph K_(r+1) need when r is %d?", r),
			r+1,
			"Every pair of vertices in a complete graph is adjacent, so each needs its own colour.",
			map[string]string{"r": strconv.Itoa(r)})
	}
}

func w17Cage(kp int) generator.Problem {
	switch kp {
	case 1:
		r := w17N(3, 5)
		g := w17N(5, 12)
		levels := (g - 1) / 2
		total := 1
		term := r
		for i := 1; i <= levels; i++ {
			total += term
			term *= r - 1
		}
		return w17P(kp,
			fmt.Sprintf("What is the Moore bound for a graph of minimum degree %d and girth %d?", r, g),
			total,
			"Sum the breadth-first levels around a vertex until the girth closes the tree.",
			map[string]string{"r": strconv.Itoa(r), "g": strconv.Itoa(g), "levels": strconv.Itoa(levels)})
	case 2:
		v := 2 * w17N(5, 20)
		r := w17N(3, 5)
		return w17P(kp,
			fmt.Sprintf("A %d-regular graph has %d vertices. How many edges does it have?", r, v),
			v*r/2,
			"Handshake lemma: the degree sum is v times r, and each edge counts twice.",
			map[string]string{"v": strconv.Itoa(v), "r": strconv.Itoa(r)})
	default:
		type cage struct {
			r, g, v int
		}
		cages := []cage{
			{3, 5, 10}, {3, 6, 14}, {3, 7, 24}, {3, 8, 30},
			{3, 9, 58}, {3, 10, 70}, {3, 11, 112}, {3, 12, 126},
			{4, 5, 19}, {4, 6, 26}, {5, 5, 30},
		}
		pick := cages[rand.Intn(len(cages))]
		return w17P(kp,
			fmt.Sprintf("How many vertices has the (%d,%d)-cage?", pick.r, pick.g),
			pick.v,
			"The cage is the smallest regular graph with the given degree and girth, and its order is classically known.",
			map[string]string{"r": strconv.Itoa(pick.r), "g": strconv.Itoa(pick.g), "v": strconv.Itoa(pick.v)})
	}
}

func w17Zarankiewicz(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(5, 40)
		bound := int(math.Floor(float64(n) / 4 * (1 + math.Sqrt(float64(4*n-3)))))
		return w17P(kp,
			fmt.Sprintf("The Reiman bound limits the edges of a K(2,2)-free graph on %d vertices. What is that bound?", n),
			bound,
			"Count length-two paths and apply convexity to get n/4 times one plus the square root of 4n-3.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		s := w17N(2, 9)
		t := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("How many edges has the complete bipartite graph K(%d,%d)?", s, t),
			s*t,
			"Every vertex on one side joins every vertex on the other, so multiply the side sizes.",
			map[string]string{"s": strconv.Itoa(s), "t": strconv.Itoa(t)})
	default:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("How many four-cycles does K(%d,%d) contain?", n, n),
			w17C2(n)*w17C2(n),
			"A four-cycle is a choice of two vertices on each side of the bipartition.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Critical(kp int) generator.Problem {
	switch kp {
	case 1:
		k := w17N(3, 9)
		return w17P(kp,
			fmt.Sprintf("A graph is %d-critical. What is its minimum possible degree?", k),
			k-1,
			"A vertex of degree below k-1 could be deleted, the rest coloured with k-1 colours, and the vertex recoloured.",
			map[string]string{"k": strconv.Itoa(k)})
	case 2:
		k := w17N(3, 9)
		n := w17N(k, 24)
		total := n * (k - 1)
		if total%2 != 0 {
			total++
		}
		return w17P(kp,
			fmt.Sprintf("A %d-critical graph has %d vertices. What is the smallest possible number of edges?", k, n),
			total/2,
			"Minimum degree k-1 gives a degree sum of at least n(k-1), and each edge contributes two.",
			map[string]string{"k": strconv.Itoa(k), "n": strconv.Itoa(n)})
	default:
		k := w17N(3, 9)
		return w17P(kp,
			fmt.Sprintf("How many edges has the complete graph K_%d, the standard %d-critical example?", k, k),
			w17C2(k),
			"Every pair of vertices is joined, so the edge count is k choose two.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Goodman(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 12)
		return w17P(kp,
			fmt.Sprintf("How many triangles has the complete graph on %d vertices?", n),
			n*(n-1)*(n-2)/6,
			"A triangle is a choice of three vertices.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 14)
		d := w17N(1, n-1)
		return w17P(kp,
			fmt.Sprintf("In a graph on %d vertices, a vertex has degree %d. How many mixed pairs of incident edges and non-edges pass through it?", n, d),
			d*(n-1-d),
			"Choose one of the d neighbours and one of the n-1-d non-neighbours, so multiply the two counts.",
			map[string]string{"n": strconv.Itoa(n), "d": strconv.Itoa(d)})
	default:
		n := w17N(5, 12)
		mixed := 2 * w17N(2, 20)
		return w17P(kp,
			fmt.Sprintf("A two-colouring of K_%d has %d mixed pairs counted over all vertices. How many monochromatic triangles does Goodman's formula give?", n, mixed),
			n*(n-1)*(n-2)/6-mixed/2,
			"Subtract half the mixed pairs from the total number of triangles.",
			map[string]string{"n": strconv.Itoa(n), "mixed": strconv.Itoa(mixed)})
	}
}

func w17Friendship(kp int) generator.Problem {
	k := w17N(2, 12)
	switch kp {
	case 1:
		return w17P(kp,
			fmt.Sprintf("The windmill graph F_%d is made of %d triangles sharing one vertex. How many vertices has it?", k, k),
			2*k+1,
			"Each triangle adds two new vertices around the shared one, giving 2k+1.",
			map[string]string{"k": strconv.Itoa(k)})
	case 2:
		return w17P(kp,
			fmt.Sprintf("The windmill graph F_%d has %d triangles sharing one vertex. How many edges has it?", k, k),
			3*k,
			"Each triangle contributes three edges and the triangles meet only at the shared vertex.",
			map[string]string{"k": strconv.Itoa(k)})
	default:
		return w17P(kp,
			fmt.Sprintf("How many triangles has the windmill graph F_%d?", k),
			k,
			"Every triangle of the windmill contains the shared vertex, and there are k of them.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Hajnal(kp int) generator.Problem {
	delta := w17N(2, 9)
	k := delta + 1
	n := w17N(k, 40)
	switch kp {
	case 1:
		return w17P(kp,
			fmt.Sprintf("A graph has maximum degree %d. How many colours does the Hajnal-Szemeredi theorem guarantee for an equitable colouring?", delta),
			k,
			"One more colour than the maximum degree is always enough.",
			map[string]string{"delta": strconv.Itoa(delta)})
	case 2:
		return w17P(kp,
			fmt.Sprintf("An equitable colouring of %d vertices uses %d colours. What is the largest class size?", n, k),
			w17CeilDiv(n, k),
			"Distribute the vertices as evenly as possible, so the largest class is n divided by k rounded up.",
			map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k)})
	default:
		return w17P(kp,
			fmt.Sprintf("An equitable colouring of %d vertices uses %d colours. What is the smallest class size?", n, k),
			n/k,
			"Distribute the vertices as evenly as possible, so the smallest class is n divided by k rounded down.",
			map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k)})
	}
}

// ---------------------------------------------------------------------------
// Flows and matching duality
// ---------------------------------------------------------------------------

func w17Hall(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 10)
		k := w17N(2, 5)
		return w17P(kp,
			fmt.Sprintf("A %d-regular bipartite graph has %d vertices on each side. What size is the perfect matching guaranteed by Hall's theorem?", k, n),
			n,
			"Hall's condition holds in a regular bipartite graph, so every vertex on one side is matched.",
			map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k)})
	case 2:
		n := w17N(4, 12)
		best := w17N(2, n-1)
		return w17P(kp,
			fmt.Sprintf("One side of a bipartite graph has %d vertices and the largest neighbourhood of any subset has %d. What matching size does Hall's deficiency leave?", n, best),
			n-best,
			"The deficiency is the worst shortfall between a subset and its neighbourhood.",
			map[string]string{"n": strconv.Itoa(n), "best": strconv.Itoa(best)})
	default:
		n := w17N(3, 7)
		f := 1
		for i := 2; i <= n; i++ {
			f *= i
		}
		return w17P(kp,
			fmt.Sprintf("How many perfect matchings has K(%d,%d)?", n, n),
			f,
			"A perfect matching is a permutation of the n vertices on one side.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Konig(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 24)
		return w17P(kp,
			fmt.Sprintf("A path on %d vertices has how many edges in a maximum matching, by Konig's theorem equal to its minimum vertex cover?", n),
			n/2,
			"A path can be matched along alternating edges, leaving at most one vertex.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 24)
		return w17P(kp,
			fmt.Sprintf("A cycle on %d vertices has what maximum matching size?", n),
			n/2,
			"Pair consecutive vertices around the cycle; for odd n one vertex remains.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		m := w17N(2, 10)
		n := w17N(2, 10)
		best := m
		if n < best {
			best = n
		}
		return w17P(kp,
			fmt.Sprintf("What is the maximum matching size in K(%d,%d)?", m, n),
			best,
			"A matching pairs distinct vertices across the sides, so it is limited by the smaller side.",
			map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n)})
	}
}

func w17MinCost(kp int) generator.Problem {
	switch kp {
	case 1:
		f := w17N(3, 40)
		c := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("A flow of %d units travels along an edge of unit cost %d. What is the total cost?", f, c),
			f*c,
			"Multiply the flow by the cost per unit.",
			map[string]string{"f": strconv.Itoa(f), "c": strconv.Itoa(c)})
	case 2:
		f1 := w17N(2, 20)
		c1 := w17N(2, 8)
		f2 := w17N(2, 20)
		c2 := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("A flow sends %d units at cost %d and %d units at cost %d. What is the total cost?", f1, c1, f2, c2),
			f1*c1+f2*c2,
			"Sum the two edge costs, each flow times its own unit price.",
			map[string]string{"f1": strconv.Itoa(f1), "c1": strconv.Itoa(c1), "f2": strconv.Itoa(f2), "c2": strconv.Itoa(c2)})
	default:
		f := w17N(3, 30)
		c := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("Cancelling %d units of flow on an edge of unit cost %d saves how much cost?", f, c),
			f*c,
			"The reverse residual edge carries the negative of the original cost, so cancelling recovers the same amount.",
			map[string]string{"f": strconv.Itoa(f), "c": strconv.Itoa(c)})
	}
}

func w17Circulation(kp int) generator.Problem {
	switch kp {
	case 1:
		in := w17N(4, 40)
		return w17P(kp,
			fmt.Sprintf("At a vertex of a circulation the inflow is %d. What is the outflow?", in),
			in,
			"Conservation forces inflow and outflow to balance at every vertex.",
			map[string]string{"in": strconv.Itoa(in)})
	case 2:
		k := w17N(3, 9)
		l := w17N(1, 6)
		return w17P(kp,
			fmt.Sprintf("A circulation has %d edges, each with lower bound %d. What is the total of the lower bounds?", k, l),
			k*l,
			"Add the lower bound once for each edge.",
			map[string]string{"k": strconv.Itoa(k), "l": strconv.Itoa(l)})
	default:
		k := w17N(2, 6)
		c := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("A cut consists of %d edges of capacity %d. What is the cut capacity?", k, c),
			k*c,
			"The capacity of a cut is the sum of the capacities of its edges.",
			map[string]string{"k": strconv.Itoa(k), "c": strconv.Itoa(c)})
	}
}

func w17GomoryHu(kp int) generator.Problem {
	n := w17N(3, 14)
	switch kp {
	case 1:
		return w17P(kp,
			fmt.Sprintf("A Gomory-Hu tree has how many edges for a graph on %d vertices?", n),
			n-1,
			"A tree on n vertices always has n-1 edges.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		return w17P(kp,
			fmt.Sprintf("A Gomory-Hu tree answers all pairwise minimum cuts of %d vertices. How many vertex pairs is that?", n),
			w17C2(n),
			"Every unordered pair of vertices is one minimum-cut query.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		return w17P(kp,
			fmt.Sprintf("Constructing a Gomory-Hu tree for %d vertices needs how many maximum-flow computations?", n),
			n-1,
			"Each computation adds one tree edge, and the tree needs n-1 edges.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17MultiwayCut(kp int) generator.Problem {
	k := w17N(3, 10)
	switch kp {
	case 1:
		return w17P(kp,
			fmt.Sprintf("A star has %d leaves, all terminals. Deleting all spokes gives a multiway cut. How many edges is that?", k),
			k,
			"One spoke per leaf, and no two leaves remain connected once every spoke is gone.",
			map[string]string{"k": strconv.Itoa(k)})
	case 2:
		return w17P(kp,
			fmt.Sprintf("The isolating-cut heuristic for %d terminals is within a factor of (2k-2)/k of the optimum. What is the numerator 2k-2?", k),
			2*k-2,
			"Substitute k into 2-2/k and read the numerator over k.",
			map[string]string{"k": strconv.Itoa(k)})
	default:
		return w17P(kp,
			fmt.Sprintf("How many edges has the complete graph on %d terminals?", k),
			w17C2(k),
			"Each pair of terminals is joined by one edge.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17FlowIntegrality(kp int) generator.Problem {
	switch kp {
	case 1:
		c := w17N(3, 30)
		return w17P(kp,
			fmt.Sprintf("A single edge with integer capacity %d carries the maximum flow. What is that flow value?", c),
			c,
			"The edge is the only route, so the maximum flow equals its capacity.",
			map[string]string{"c": strconv.Itoa(c)})
	case 2:
		k := w17N(2, 6)
		c := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("A cut has %d edges of integer capacity %d. What is its capacity?", k, c),
			k*c,
			"Sum the integer capacities along the cut.",
			map[string]string{"k": strconv.Itoa(k), "c": strconv.Itoa(c)})
	default:
		d := w17N(2, 6)
		c := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("The source has %d outgoing edges of capacity %d. What is the total capacity leaving the source?", d, c),
			d*c,
			"No flow can exceed the total capacity out of the source.",
			map[string]string{"d": strconv.Itoa(d), "c": strconv.Itoa(c)})
	}
}

func w17EdgeDisjoint(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 12)
		return w17P(kp,
			fmt.Sprintf("The edge connectivity of K_%d equals the number of edge-disjoint paths between any two vertices. What is it?", n),
			n-1,
			"Each vertex has n-1 incident edges, and that is both the minimum cut and the path packing.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("A cycle on %d vertices has what edge connectivity?", n),
			2,
			"Two edges must be removed to break a cycle, and two edge-disjoint paths join any pair of vertices.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("A tree has %d vertices. How many edge-disjoint paths join two vertices in different branches?", n),
			1,
			"Every edge of a tree is a bridge, so only one path can join the two vertices.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Dinic(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("Dinic's algorithm on a graph with %d vertices has at most how many phases?", n),
			n-1,
			"The distance from source to sink strictly increases, and it cannot exceed n-1.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		m := w17N(4, 60)
		root := int(math.Floor(math.Sqrt(float64(m))))
		return w17P(kp,
			fmt.Sprintf("On a unit-capacity network with %d edges, Dinic's algorithm needs about the square root of the edge count in phases. What is that integer?", m),
			root,
			"The unit-capacity bound is O(sqrt(E)), so take the integer square root.",
			map[string]string{"m": strconv.Itoa(m)})
	default:
		n := w17N(4, 10)
		m := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("The general Dinic bound is O(V^2 E). For V=%d and E=%d, what is V squared times E?", n, m),
			n*n*m,
			"Square the vertex count and multiply by the edge count.",
			map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m)})
	}
}

func w17PushRelabel(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("Push-relabel relabels each vertex at most 2V-1 times. For V=%d, what is that bound?", n),
			2*n-1,
			"Double the vertex count and subtract one.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		d := w17N(2, 6)
		c := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("The source has %d outgoing edges of capacity %d, all saturated at the start. What is the initial preflow value?", d, c),
			d*c,
			"Saturating every source edge sends the sum of those capacities into the network.",
			map[string]string{"d": strconv.Itoa(d), "c": strconv.Itoa(c)})
	default:
		l := w17N(2, 12)
		return w17P(kp,
			fmt.Sprintf("A path from source to sink has %d edges. How many pushes does it take to send one unit along it?", l),
			l,
			"Each push advances one edge, and the unit ends at the sink after the last edge.",
			map[string]string{"l": strconv.Itoa(l)})
	}
}
