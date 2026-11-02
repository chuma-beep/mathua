package graphtheory

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// graphBridgeKinds maps the concepts added by the 2026 content expansion to
// the question family their generator serves.
var graphBridgeKinds = map[string]string{
	"graph.hamiltonian":   "hamiltonian",
	"graph.blocks":        "blocks",
	"graph.planar_dual":   "planar_dual",
	"graph.automorphism":  "automorphism",
	"graph.girth":         "girth",
	"graph.four_color":    "four_color",
	"graph.matroid":       "matroid",
	"graph.shortest_path": "shortest_path",
	"graph.degeneracy":    "degeneracy",
	"graph.mst":           "mst",
}

type graphBridgeGen struct{ kind string }

func graphBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func graphBridgeProblem(ordinal int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		KPOrdinal:   ordinal,
		Facts:       map[string]string{"answer": answer},
	}
}

func graphFactorial(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func graphChoose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

func graphMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *graphBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := graphBridgeOrdinal(ctx)
	switch g.kind {
	case "hamiltonian":
		return hamiltonianProblem(o)
	case "blocks":
		return blocksProblem(o)
	case "planar_dual":
		return planarDualProblem(o)
	case "automorphism":
		return automorphismProblem(o)
	case "girth":
		return girthProblem(o)
	case "four_color":
		return fourColorProblem(o)
	case "matroid":
		return matroidProblem(o)
	case "shortest_path":
		return shortestPathProblem(o)
	case "degeneracy":
		return degeneracyProblem(o)
	default:
		return mstProblem(o)
	}
}

func hamiltonianProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(7)
		ans := graphFactorial(n-1) / 2
		return graphBridgeProblem(o,
			fmt.Sprintf("How many distinct Hamiltonian cycles does the complete graph K_%d have? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Fix one starting vertex, order the other %d vertices, and divide by two for direction: %d!/2 = %d.", n-1, n-1, ans))
	case 2:
		n := 2 * (2 + rand.Intn(5))
		ans := n / 2
		return graphBridgeProblem(o,
			fmt.Sprintf("Dirac's condition guarantees a Hamiltonian cycle in a simple graph on %d vertices when every vertex has degree at least what value? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Dirac's bound is half the vertex count: %d/2 = %d.", n, ans))
	default:
		n := 3 + rand.Intn(9)
		return graphBridgeProblem(o,
			fmt.Sprintf("The cycle graph C_%d is itself a Hamiltonian cycle. How many edges does that Hamiltonian cycle contain? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("A Hamiltonian cycle on %d vertices uses all %d edges of C_%d.", n, n, n))
	}
}

func blocksProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(10)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many blocks does a tree on %d vertices have? (enter a number)", n),
			strconv.Itoa(n-1),
			fmt.Sprintf("Every edge of a tree is a bridge and therefore its own block, so the %d edges give %d blocks.", n-1, n-1))
	case 2:
		n := 4 + rand.Intn(9)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many cut vertices does the path graph P_%d have? (enter a number)", n),
			strconv.Itoa(n-2),
			fmt.Sprintf("The two endpoints are not cut vertices, and removing any internal vertex disconnects the path: %d - 2 = %d.", n, n-2))
	default:
		n := 3 + rand.Intn(10)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many blocks does the cycle graph C_%d have? (enter a number)", n),
			"1",
			fmt.Sprintf("C_%d is biconnected and has no cut vertex, so the whole cycle is a single block.", n))
	}
}

func planarDualProblem(o int) generator.Problem {
	v := 4 + rand.Intn(7)
	maxEdges := 3*v - 6
	e := v - 1 + rand.Intn(maxEdges-(v-1)+1)
	faces := e - v + 2
	switch o {
	case 1:
		return graphBridgeProblem(o,
			fmt.Sprintf("A connected planar graph has %d vertices and %d edges. How many faces does a planar drawing have? (enter a number)", v, e),
			strconv.Itoa(faces),
			fmt.Sprintf("Euler's formula gives F = E - V + 2 = %d - %d + 2 = %d.", e, v, faces))
	case 2:
		return graphBridgeProblem(o,
			fmt.Sprintf("A connected planar graph has %d vertices and %d edges. How many vertices does its planar dual have? (enter a number)", v, e),
			strconv.Itoa(faces),
			fmt.Sprintf("Each face becomes one dual vertex, and F = E - V + 2 = %d - %d + 2 = %d.", e, v, faces))
	default:
		return graphBridgeProblem(o,
			fmt.Sprintf("A connected planar graph has %d vertices and %d edges. How many edges does its planar dual have? (enter a number)", v, e),
			strconv.Itoa(e),
			fmt.Sprintf("Each edge of the original graph corresponds to exactly one edge of the dual, so the dual also has %d edges.", e))
	}
}

func automorphismProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(5)
		ans := graphFactorial(n)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many automorphisms does the complete graph K_%d have? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("Every permutation of the %d vertices preserves adjacency, so there are %d! = %d automorphisms.", n, n, ans))
	case 2:
		n := 3 + rand.Intn(8)
		ans := 2 * n
		return graphBridgeProblem(o,
			fmt.Sprintf("How many automorphisms does the cycle graph C_%d have? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("There are %d rotations and %d reflections, giving 2 x %d = %d automorphisms.", n, n, n, ans))
	default:
		m := 2 + rand.Intn(4)
		n := 2 + rand.Intn(4)
		for n == m {
			n = 2 + rand.Intn(4)
		}
		ans := graphFactorial(m) * graphFactorial(n)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many automorphisms does the complete bipartite graph K_%d,%d have? (enter a number)", m, n),
			strconv.Itoa(ans),
			fmt.Sprintf("The parts have different sizes, so each part maps onto itself: %d! x %d! = %d.", m, n, ans))
	}
}

func girthProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(8)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the girth of the cycle graph C_%d? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("The only cycle in C_%d has length %d, so the girth is %d.", n, n, n))
	case 2:
		m := 2 + rand.Intn(5)
		n := 2 + rand.Intn(5)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the girth of the complete bipartite graph K_%d,%d? (enter a number)", m, n),
			"4",
			"Two vertices in one part and two in the other form a 4-cycle, and bipartite graphs have no triangles.")
	default:
		n := 3 + rand.Intn(6)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the girth of the complete graph K_%d? (enter a number)", n),
			"3",
			fmt.Sprintf("Any three vertices of K_%d form a triangle, so the shortest cycle has length 3.", n))
	}
}

func fourColorProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(10)
		ans := graphMin(n, 4)
		return graphBridgeProblem(o,
			fmt.Sprintf("What number of colors is guaranteed to color every simple planar graph on %d vertices? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The four color theorem guarantees four colors, but a graph on %d vertices cannot need more than %d.", n, ans))
	case 2:
		n := 3 + rand.Intn(5)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the chromatic number of the complete graph K_%d? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("Every pair of vertices in K_%d is adjacent, so all %d colors must be distinct.", n, n))
	default:
		n := 3 + rand.Intn(10)
		ans := 3
		if n%2 == 0 {
			ans = 2
		}
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the chromatic number of the cycle graph C_%d? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("An even cycle is bipartite and needs 2 colors; an odd cycle needs 3. C_%d is %s, so the answer is %d.", n, map[bool]string{true: "even", false: "odd"}[n%2 == 0], ans))
	}
}

func matroidProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 4 + rand.Intn(6)
		r := 2 + rand.Intn(n-3)
		ans := graphChoose(n, r)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many bases does the uniform matroid U_%d,%d have? (enter a number)", r, n),
			strconv.Itoa(ans),
			fmt.Sprintf("A basis is any %d-element subset of the %d-element ground set, so there are C(%d,%d) = %d.", r, n, n, r, ans))
	case 2:
		n := 2 + rand.Intn(9)
		return graphBridgeProblem(o,
			fmt.Sprintf("A connected graph has %d vertices. What is the rank of its graphic matroid? (enter a number)", n),
			strconv.Itoa(n-1),
			fmt.Sprintf("The rank is the size of a spanning tree, which is %d - 1 = %d.", n, n-1))
	default:
		r := 1 + rand.Intn(6)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the size of every circuit in the uniform matroid U_%d,%d? (enter a number)", r, r+4),
			strconv.Itoa(r+1),
			fmt.Sprintf("Every set of at most %d elements is independent, so a minimal dependent set has %d elements.", r, r+1))
	}
}

func shortestPathProblem(o int) generator.Problem {
	switch o {
	case 1:
		a := 1 + rand.Intn(8)
		b := 1 + rand.Intn(8)
		ans := a + b
		return graphBridgeProblem(o,
			fmt.Sprintf("On a grid with unit horizontal and vertical steps, what is the shortest path length from (0,0) to (%d,%d)? (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The path needs %d horizontal and %d vertical steps, for a total of %d.", a, b, ans))
	case 2:
		n := 2 + rand.Intn(5)
		weights := make([]string, n)
		total := 0
		for i := 0; i < n; i++ {
			w := 1 + rand.Intn(5)
			total += w
			weights[i] = strconv.Itoa(w)
		}
		return graphBridgeProblem(o,
			fmt.Sprintf("A path graph has consecutive edge weights %s. What is the shortest path distance between its two endpoints? (enter a number)", strings.Join(weights, ", ")),
			strconv.Itoa(total),
			fmt.Sprintf("The only path between the endpoints uses every edge, so its length is the sum %s = %d.", strings.Join(weights, " + "), total))
	default:
		a := 1 + rand.Intn(5)
		b := 1 + rand.Intn(5)
		ans := graphChoose(a+b, a)
		return graphBridgeProblem(o,
			fmt.Sprintf("How many shortest monotone lattice paths run from (0,0) to (%d,%d)? (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("Each path has %d steps, of which %d are horizontal, so there are C(%d,%d) = %d paths.", a+b, a, a+b, a, ans))
	}
}

func degeneracyProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 2 + rand.Intn(7)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the degeneracy of the complete graph K_%d? (enter a number)", n),
			strconv.Itoa(n-1),
			fmt.Sprintf("Every subgraph of K_%d is complete, so its minimum degree is at most %d, and the whole graph attains %d.", n, n-1, n-1))
	case 2:
		n := 3 + rand.Intn(8)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the degeneracy of a tree on %d vertices? (enter a number)", n),
			"1",
			"Every subgraph of a tree is a forest and has a leaf, so the minimum degree is at most 1, and an edge attains 1.")
	default:
		n := 3 + rand.Intn(8)
		return graphBridgeProblem(o,
			fmt.Sprintf("What is the degeneracy of the cycle graph C_%d? (enter a number)", n),
			"2",
			fmt.Sprintf("C_%d has minimum degree 2, and removing vertices leaves paths whose minimum degree is at most 2, so the degeneracy is 2.", n))
	}
}

func mstProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(8)
		return graphBridgeProblem(o,
			fmt.Sprintf("A star graph on %d vertices has unit edge weights. What is the weight of a minimum spanning tree? (enter a number)", n),
			strconv.Itoa(n-1),
			fmt.Sprintf("Every spanning tree of the star uses all %d edges, each of weight 1, so the minimum weight is %d.", n-1, n-1))
	case 2:
		n := 2 + rand.Intn(9)
		return graphBridgeProblem(o,
			fmt.Sprintf("A connected graph has %d vertices. How many edges does every minimum spanning tree have? (enter a number)", n),
			strconv.Itoa(n-1),
			fmt.Sprintf("Every spanning tree on %d vertices has exactly %d edges.", n, n-1))
	default:
		n := 3 + rand.Intn(6)
		weights := make([]string, n-1)
		total := 0
		for i := range weights {
			w := 1 + rand.Intn(3)
			total += w
			weights[i] = strconv.Itoa(w)
		}
		return graphBridgeProblem(o,
			fmt.Sprintf("A path graph on %d vertices has consecutive edge weights %s. What is the weight of its minimum spanning tree? (enter a number)", n, strings.Join(weights, ", ")),
			strconv.Itoa(total),
			fmt.Sprintf("The path is already a tree, so its only spanning tree uses every edge: %s = %d.", strings.Join(weights, " + "), total))
	}
}
