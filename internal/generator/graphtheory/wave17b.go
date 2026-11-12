package graphtheory

// Wave 17, second file: colouring refinements and spectral invariants.

import (
	"fmt"
	"math"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave17Rest(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"graph.choosability":         "choosability",
		"graph.grundy":               "grundy",
		"graph.achromatic":           "achromatic",
		"graph.fractional_coloring":  "fractional",
		"graph.circular_coloring":    "circular",
		"graph.defective_coloring":   "defective",
		"graph.acyclic_coloring":     "acyclic",
		"graph.total_coloring":       "total",
		"graph.hadwiger_number":      "hadwiger",
		"graph.complete_coloring":    "complete_coloring",
		"graph.incidence":            "incidence",
		"graph.normalized_laplacian": "normalized",
		"graph.signless_laplacian":   "signless",
		"graph.interlacing":          "interlacing",
		"graph.closed_walks":         "closed_walks",
		"graph.expander_mixing":      "expander",
		"graph.cheeger":              "cheeger",
		"graph.heat_kernel":          "heat",
		"graph.distance_matrix":      "distance",
		"graph.energy":               "energy",
		"graph.rotation_system":      "rotation",
		"graph.nonorientable_genus":  "nonorientable",
		"graph.toroidal":             "toroidal",
		"graph.one_planar":           "one_planar",
		"graph.book_embedding":       "book",
		"graph.page_number":          "page",
		"graph.queue_number":         "queue",
		"graph.bandwidth":            "bandwidth",
		"graph.planar_separator":     "separator",
		"graph.apex":                 "apex",
	} {
		reg.Register(id, &w17Gen{kind: kind})
	}
}

func w17RestGenerate(kind string, kp int) generator.Problem {
	switch kind {
	case "choosability":
		return w17Choosability(kp)
	case "grundy":
		return w17Grundy(kp)
	case "achromatic":
		return w17Achromatic(kp)
	case "fractional":
		return w17Fractional(kp)
	case "circular":
		return w17Circular(kp)
	case "defective":
		return w17Defective(kp)
	case "acyclic":
		return w17Acyclic(kp)
	case "total":
		return w17Total(kp)
	case "hadwiger":
		return w17Hadwiger(kp)
	case "complete_coloring":
		return w17CompleteColoring(kp)
	case "incidence":
		return w17Incidence(kp)
	case "normalized":
		return w17Normalized(kp)
	case "signless":
		return w17Signless(kp)
	case "interlacing":
		return w17Interlacing(kp)
	case "closed_walks":
		return w17ClosedWalks(kp)
	case "expander":
		return w17Expander(kp)
	case "cheeger":
		return w17Cheeger(kp)
	case "heat":
		return w17Heat(kp)
	case "distance":
		return w17Distance(kp)
	case "energy":
		return w17Energy(kp)
	case "rotation":
		return w17Rotation(kp)
	case "nonorientable":
		return w17Nonorientable(kp)
	case "toroidal":
		return w17Toroidal(kp)
	case "one_planar":
		return w17OnePlanar(kp)
	case "book":
		return w17Book(kp)
	case "page":
		return w17Page(kp)
	case "queue":
		return w17Queue(kp)
	case "bandwidth":
		return w17Bandwidth(kp)
	case "separator":
		return w17Separator(kp)
	case "apex":
		return w17Apex(kp)
	}
	return w17Handshake(kp)
}

// ---------------------------------------------------------------------------
// Colouring refinements
// ---------------------------------------------------------------------------

func w17Choosability(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the choosability of the complete graph K_%d?", n),
			n,
			"All lists may be identical, so the choosability equals the chromatic number n.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		k := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("What is the choosability of the even cycle C_%d?", 2*k),
			2,
			"Two colours always suffice from two-element lists around an even cycle.",
			map[string]string{"k": strconv.Itoa(k)})
	default:
		k := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("What is the choosability of the odd cycle C_%d?", 2*k+1),
			3,
			"Two-element lists can fail on an odd cycle, and three colours always suffice.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Grundy(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the Grundy number of the complete graph K_%d?", n),
			n,
			"Whatever the order, each vertex sees every earlier colour, so a new colour is used every time.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(3, 16)
		return w17P(kp,
			fmt.Sprintf("What is the Grundy number of a path on %d vertices?", n),
			2,
			"An order that visits one internal vertex before its neighbours can force two colours, and no order forces a third.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(4, 12)
		return w17P(kp,
			fmt.Sprintf("What is the Grundy number of K_%d with one edge deleted?", n),
			n-1,
			"A greedy order can force n-1 colours, and the missing edge prevents the last one.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Achromatic(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the achromatic number of the complete graph K_%d?", n),
			n,
			"Every colour class is a single vertex, and every pair of classes meets in an edge.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("What is the achromatic number of a path on %d vertices?", n),
			(n+1)/2,
			"Pair vertices into colour classes of size at most two, and every pair of classes can be arranged to meet in an edge.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("What is the achromatic number of the star K(1,%d)?", n),
			2,
			"The centre and one leaf form one class, the remaining leaves another, and the two classes are joined by an edge.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Fractional(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the fractional chromatic number of the complete graph K_%d?", n),
			n,
			"Every independent set is a single vertex, so covering each vertex costs one unit.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		k := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("The odd cycle C_%d has fractional chromatic number (2k+1)/k. What is the denominator in lowest terms?", 2*k+1),
			k,
			"The fraction (2k+1)/k is already reduced because the numerator exceeds twice the denominator by one.",
			map[string]string{"k": strconv.Itoa(k)})
	default:
		m := w17N(2, 10)
		n := w17N(2, 10)
		return w17P(kp,
			fmt.Sprintf("What is the fractional chromatic number of the complete bipartite graph K(%d,%d)?", m, n),
			2,
			"The two sides are independent sets, and covering both costs exactly two.",
			map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n)})
	}
}

func w17Circular(kp int) generator.Problem {
	switch kp {
	case 1:
		k := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("The circular chromatic number of C_%d is (2k+1)/k. What is the numerator?", 2*k+1),
			2*k+1,
			"The circumference is the numerator and the gap is the denominator.",
			map[string]string{"k": strconv.Itoa(k)})
	case 2:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the circular chromatic number of the complete graph K_%d?", n),
			n,
			"Adjacent vertices need every gap, so the ratio reaches the integer chromatic number.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		k := w17N(2, 10)
		return w17P(kp,
			fmt.Sprintf("What is the circular chromatic number of the even cycle C_%d?", 2*k),
			2,
			"An even cycle is bipartite and has circular chromatic number two.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Defective(kp int) generator.Problem {
	delta := w17N(3, 12)
	d := w17N(1, 4)
	switch kp {
	case 1:
		return w17P(kp,
			fmt.Sprintf("A graph has maximum degree %d. How many colours does the greedy (k,d)-colouring with defect %d need?", delta, d),
			w17CeilDiv(delta+1, d+1),
			"Choose k so that k(d+1) reaches delta+1.",
			map[string]string{"delta": strconv.Itoa(delta), "d": strconv.Itoa(d)})
	case 2:
		n := 2 * w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("A defective colouring of %d vertices allows at most %d same-colour neighbours per vertex. What is the largest possible number of monochromatic edges?", n, d),
			n*d/2,
			"Each colour class has maximum degree d, so it contributes at most n_c d/2 edges, and the sum is nd/2.",
			map[string]string{"n": strconv.Itoa(n), "d": strconv.Itoa(d)})
	default:
		k := w17N(2, 6)
		return w17P(kp,
			fmt.Sprintf("With %d colours and defect %d, what is the largest maximum degree the greedy bound covers?", k, d),
			k*(d+1)-1,
			"The condition is k(d+1) at least delta+1, so the largest delta is k(d+1)-1.",
			map[string]string{"k": strconv.Itoa(k), "d": strconv.Itoa(d)})
	}
}

func w17Acyclic(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("What is the acyclic chromatic number of a tree on %d vertices?", n),
			2,
			"Two colours make every colour class independent, and their union is the whole tree, which is a forest.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the acyclic chromatic number of the complete graph K_%d?", n),
			n,
			"Any two colour classes meet in at most one edge, so every vertex needs its own colour.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(4, 16)
		return w17P(kp,
			fmt.Sprintf("What is the acyclic chromatic number of the cycle C_%d?", n),
			3,
			"Two colours make the whole cycle bichromatic, so a third colour is needed to break it.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Total(kp int) generator.Problem {
	switch kp {
	case 1:
		delta := w17N(2, 10)
		return w17P(kp,
			fmt.Sprintf("A graph has maximum degree %d. What is the lower bound on its total chromatic number?", delta),
			delta+1,
			"The vertex and its incident edges form a conflict clique of size delta+1.",
			map[string]string{"delta": strconv.Itoa(delta)})
	case 2:
		n := w17N(4, 12)
		answer := n
		if n%2 == 0 {
			answer = n + 1
		}
		return w17P(kp,
			fmt.Sprintf("What is the total chromatic number of the complete graph K_%d?", n),
			answer,
			"An odd complete graph reaches n, and an even one needs n+1 because the edge colours cannot be paired.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(3, 16)
		return w17P(kp,
			fmt.Sprintf("What is the total chromatic number of a path on %d vertices?", n),
			3,
			"The two end vertices and their incident edges already force three colours, and three suffice.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Hadwiger(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("What is the Hadwiger number of the complete graph K_%d?", n),
			n,
			"The graph is already a complete graph, so its largest complete minor is itself.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 16)
		return w17P(kp,
			fmt.Sprintf("What is the Hadwiger number of the cycle C_%d?", n),
			3,
			"Contracting a path leaves a triangle, and no cycle has a K_4 minor.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(4, 18)
		return w17P(kp,
			fmt.Sprintf("What is the Hadwiger number of a tree on %d vertices?", n),
			2,
			"Contracting edges of a tree produces another tree, and the largest complete minor is a single edge.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17CompleteColoring(kp int) generator.Problem {
	switch kp {
	case 1:
		k := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("A complete colouring has %d colour classes. How many edges between classes are needed at least?", k),
			w17C2(k),
			"Every pair of classes needs one edge, so count the pairs.",
			map[string]string{"k": strconv.Itoa(k)})
	case 2:
		m := w17N(6, 60)
		k := 0
		for (k+1)*k/2 <= m {
			k++
		}
		return w17P(kp,
			fmt.Sprintf("A graph has %d edges. What is the largest number of classes a complete colouring can have?", m),
			k,
			"Each pair of classes needs an edge, so the number of classes is limited by k choose two at most m.",
			map[string]string{"m": strconv.Itoa(m)})
	default:
		k := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("Refining a colouring to %d classes requires at least how many edges inside the graph?", k),
			w17C2(k),
			"Completeness demands one edge for each pair of classes.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

// ---------------------------------------------------------------------------
// Spectral graph theory
// ---------------------------------------------------------------------------

func w17Incidence(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 12)
		m := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("A graph has %d vertices and %d edges. How many entries has its incidence matrix?", n, m),
			n*m,
			"The incidence matrix has one row per vertex and one column per edge.",
			map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m)})
	case 2:
		m := w17N(3, 24)
		return w17P(kp,
			fmt.Sprintf("A graph has %d edges. How many darts does it have?", m),
			2*m,
			"Every edge has two darts, one in each direction.",
			map[string]string{"m": strconv.Itoa(m)})
	default:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("A connected graph has %d vertices. What is the rank of its incidence matrix over the reals?", n),
			n-1,
			"The incidence matrix of a connected graph has rank n-1, with the all-ones vector in its null space.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Normalized(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("The normalized Laplacian of a graph on %d vertices has trace n. What is the sum of its eigenvalues?", n),
			n,
			"The trace is the sum of the eigenvalues, and it equals the number of vertices.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		c := w17N(1, 9)
		return w17P(kp,
			fmt.Sprintf("A graph has %d connected components. How many zero eigenvalues has its normalized Laplacian?", c),
			c,
			"Each component contributes exactly one zero eigenvalue.",
			map[string]string{"c": strconv.Itoa(c)})
	default:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("The complete graph K_%d has second normalized eigenvalue n/(n-1). What is (n-1) times that eigenvalue?", n),
			n,
			"Multiply n/(n-1) by n-1 to recover n.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Signless(kp int) generator.Problem {
	switch kp {
	case 1:
		m := w17N(3, 24)
		return w17P(kp,
			fmt.Sprintf("A graph has %d edges. What is the trace of its signless Laplacian?", m),
			2*m,
			"The trace is the degree sum, which is twice the number of edges.",
			map[string]string{"m": strconv.Itoa(m)})
	case 2:
		d := w17N(1, 12)
		return w17P(kp,
			fmt.Sprintf("A vertex has degree %d. What is the corresponding row sum of the signless Laplacian?", d),
			d,
			"The signless Laplacian is D+A, and each row sums to the degree.",
			map[string]string{"d": strconv.Itoa(d)})
	default:
		n := w17N(4, 16)
		return w17P(kp,
			fmt.Sprintf("A connected bipartite graph has %d vertices. How many zero eigenvalues has its signless Laplacian?", n),
			1,
			"Bipartiteness makes the smallest signless eigenvalue zero, and connectivity makes it simple.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Interlacing(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("How many adjacency eigenvalues has a graph on %d vertices?", n),
			n,
			"The adjacency matrix is n by n, so it has n eigenvalues counted with multiplicity.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(5, 16)
		k := w17N(1, n-2)
		return w17P(kp,
			fmt.Sprintf("Deleting %d vertices from a graph on %d vertices leaves how many adjacency eigenvalues?", k, n),
			n-k,
			"The induced subgraph has one eigenvalue per remaining vertex.",
			map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k)})
	default:
		k := w17N(1, 8)
		return w17P(kp,
			fmt.Sprintf("After deleting %d vertices, the i-th eigenvalue of the subgraph lies between which two indices of the parent spectrum?", k),
			k+1,
			"Interlacing shifts the upper index by k, so lambda_i of the subgraph is at most lambda_(i+k) of the parent.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17ClosedWalks(kp int) generator.Problem {
	switch kp {
	case 1:
		m := w17N(3, 24)
		return w17P(kp,
			fmt.Sprintf("A graph has %d edges. How many closed walks of length two does it have?", m),
			2*m,
			"Every edge supports two closed walks of length two, one starting in each direction.",
			map[string]string{"m": strconv.Itoa(m)})
	case 2:
		t := w17N(1, 10)
		return w17P(kp,
			fmt.Sprintf("A graph has %d triangles. How many closed walks of length three does it have?", t),
			6*t,
			"Each triangle supports three starting vertices and two directions.",
			map[string]string{"t": strconv.Itoa(t)})
	default:
		d := w17N(1, 12)
		return w17P(kp,
			fmt.Sprintf("A vertex has degree %d. How many walks of length two start at it?", d),
			d,
			"A two-step walk returns to the starting vertex through one of its d neighbours.",
			map[string]string{"d": strconv.Itoa(d)})
	}
}

func w17Expander(kp int) generator.Problem {
	switch kp {
	case 1:
		n := 4 * w17N(2, 8)
		d := 2 * w17N(1, 6)
		return w17P(kp,
			fmt.Sprintf("A %d-regular graph on %d vertices has sets S and T of size %d each. What is the expected number of edges between them?", d, n, n/2),
			d*n/4,
			"The expected count is d times the product of the set sizes divided by n.",
			map[string]string{"n": strconv.Itoa(n), "d": strconv.Itoa(d)})
	case 2:
		lambda := w17N(1, 4)
		s := w17N(2, 10)
		return w17P(kp,
			fmt.Sprintf("Two sets of size %d each meet in at most lambda times the square root of their product of sizes. With lambda %d, what is that bound?", s, lambda),
			lambda*s,
			"The square root of s squared is s, so the bound is lambda times s.",
			map[string]string{"lambda": strconv.Itoa(lambda), "s": strconv.Itoa(s)})
	default:
		n := 4 * w17N(2, 8)
		d := 4 * w17N(1, 5)
		return w17P(kp,
			fmt.Sprintf("A %d-regular graph on %d vertices has a set T of %d vertices. How many edges per vertex leave the set in expectation?", d, n, n/4),
			d/4,
			"The expected degree into T is d times the fraction of vertices in T.",
			map[string]string{"n": strconv.Itoa(n), "d": strconv.Itoa(d)})
	}
}

func w17Cheeger(kp int) generator.Problem {
	switch kp {
	case 1:
		n := 2 * w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("The complete graph K_%d has Cheeger constant n/2. What is twice that constant?", n),
			n,
			"Multiply n/2 by two.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := 2 * w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("The set realising the Cheeger constant of K_%d has what size?", n),
			n/2,
			"Cut off half the vertices, so the smaller side has n/2 vertices.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		k := w17N(1, 8)
		s := w17N(2, 10)
		return w17P(kp,
			fmt.Sprintf("A set of %d vertices has expansion %d. How many boundary edges does it have?", s, k),
			k*s,
			"Expansion is the boundary size divided by the set size, so multiply back.",
			map[string]string{"k": strconv.Itoa(k), "s": strconv.Itoa(s)})
	}
}

func w17Heat(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("At time zero the heat kernel on a graph with %d vertices has what trace?", n),
			n,
			"The heat kernel is the identity at time zero, whose trace is the number of vertices.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		h := w17N(5, 60)
		return w17P(kp,
			fmt.Sprintf("A function on a graph has total heat %d. The heat equation conserves total heat. What is the total after a long time?", h),
			h,
			"Diffusion moves heat between vertices but never changes the total.",
			map[string]string{"h": strconv.Itoa(h)})
	default:
		c := w17N(1, 6)
		return w17P(kp,
			fmt.Sprintf("A graph has %d connected components. How many eigenvalues of the Laplacian are zero, and therefore how many stationary heat modes are there?", c),
			c,
			"Each connected component contributes one constant stationary mode.",
			map[string]string{"c": strconv.Itoa(c)})
	}
}

func w17Distance(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("A graph has %d vertices. How many entries has its distance matrix?", n),
			n*n,
			"The distance matrix is square with one row and one column per vertex.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("What is the sum of all pairwise distances on a path with %d vertices?", n),
			n*(n*n-1)/6,
			"The classical path sum is n(n squared minus one) divided by six.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(4, 24)
		return w17P(kp,
			fmt.Sprintf("What is the diameter of a path on %d vertices?", n),
			n-1,
			"The endpoints are as far apart as possible, at distance n-1.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Energy(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("What is the energy of the complete graph K_%d?", n),
			2*(n-1),
			"The eigenvalues are n-1 once and -1 with multiplicity n-1, so the absolute sum is 2(n-1).",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		a := w17N(2, 6)
		b := w17N(2, 6)
		return w17P(kp,
			fmt.Sprintf("What is the energy of the complete bipartite graph K(%d,%d), whose nonzero eigenvalues are plus and minus the square root of mn?", a*a, b*b),
			2*a*b,
			"The square root of a squared times b squared is ab, and both signs contribute.",
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)})
	default:
		a := w17N(2, 8)
		b := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("A graph is the disjoint union of K_%d and K_%d. What is its energy?", a, b),
			2*(a-1)+2*(b-1),
			"Energy is additive over disjoint unions, and a complete graph contributes twice one less than its order.",
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b)})
	}
}

// ---------------------------------------------------------------------------
// Planarity and layouts
// ---------------------------------------------------------------------------

func w17Rotation(kp int) generator.Problem {
	switch kp {
	case 1:
		m := w17N(4, 24)
		return w17P(kp,
			fmt.Sprintf("A graph has %d edges. How many darts does its rotation system order?", m),
			2*m,
			"Each edge has two darts, and the rotation system orders the darts at every vertex.",
			map[string]string{"m": strconv.Itoa(m)})
	case 2:
		n := w17N(4, 16)
		m := w17N(n-1, 3*n-6)
		return w17P(kp,
			fmt.Sprintf("A connected planar graph has %d vertices and %d edges. How many faces has the embedding?", n, m),
			m-n+2,
			"Euler's formula for the plane gives faces equal to edges minus vertices plus two.",
			map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m)})
	default:
		d := w17N(1, 12)
		return w17P(kp,
			fmt.Sprintf("A vertex has degree %d. How many darts leave it in the rotation system?", d),
			d,
			"One dart per incident edge, ordered cyclically around the vertex.",
			map[string]string{"d": strconv.Itoa(d)})
	}
}

func w17Nonorientable(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 16)
		m := w17N(n, 3*n)
		f := w17N(2, 12)
		return w17P(kp,
			fmt.Sprintf("A connected graph on a nonorientable surface has %d vertices, %d edges and %d faces. How many crosscaps does the surface have?", n, m, f),
			2-n+m-f,
			"Solve n-m+f=2-k for the number of crosscaps k.",
			map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m), "f": strconv.Itoa(f)})
	case 2:
		k := w17N(2, 9)
		return w17P(kp,
			fmt.Sprintf("The projective plane has crosscap number 1, and the connected sum of %d projective planes adds the crosscap numbers. What is the crosscap number of the sum?", k),
			k,
			"Crosscap numbers add under connected sum, so k copies of the projective plane give crosscap number k.",
			map[string]string{"k": strconv.Itoa(k)})
	default:
		k := w17N(1, 3)
		n := w17N(4, 14)
		m := w17N(n, 3*n)
		return w17P(kp,
			fmt.Sprintf("A connected graph on a surface with %d crosscaps has %d vertices and %d edges. How many faces has the embedding?", k, n, m),
			m-n+2-k,
			"Euler's formula on a nonorientable surface gives f = m-n+2-k.",
			map[string]string{"k": strconv.Itoa(k), "n": strconv.Itoa(n), "m": strconv.Itoa(m)})
	}
}

func w17Toroidal(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 24)
		return w17P(kp,
			fmt.Sprintf("A simple toroidal graph has %d vertices. What is the largest possible number of edges?", n),
			3*n,
			"The torus gives an edge bound of 3n, one more factor of n than the planar bound.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("A triangulation of the torus has %d vertices and meets the bound of 3n edges. How many edges has it?", n),
			3*n,
			"A toroidal triangulation spends the full 3n edge budget.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(3, 16)
		return w17P(kp,
			fmt.Sprintf("A toroidal graph has %d vertices and meets the bound of 3n edges. How many faces has the embedding?", n),
			2*n,
			"Euler's formula on the torus reads n-3n+f=0, so f=2n.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17OnePlanar(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 24)
		return w17P(kp,
			fmt.Sprintf("A simple 1-planar graph has %d vertices. What is its maximum number of edges?", n),
			4*n-8,
			"Planarising each crossing adds a degree-four vertex and yields the bound 4n-8.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("For %d vertices, by how much does the 1-planar edge bound exceed the planar bound?", n),
			n-2,
			"Subtract 3n-6 from 4n-8 to get n-2.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(4, 12)
		return w17P(kp,
			fmt.Sprintf("How many edges has K_%d, the graph used to test the 1-planar bound?", n),
			w17C2(n),
			"Count the pairs of vertices in a complete graph.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Book(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("What is the book thickness of the complete graph K_%d?", n),
			w17CeilDiv(n, 2),
			"The lower bound is half the order rounded up, and the same number of pages suffices.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("What is the book thickness of an outerplanar graph on %d vertices?", n),
			1,
			"An outerplanar graph fits in one page with every vertex on the spine.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		k := w17N(2, 8)
		return w17P(kp,
			fmt.Sprintf("For an even order 2k with k=%d, how many pages does K_2k need?", k),
			k,
			"The book thickness of K_n is half the order rounded up, which is k for even n.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Page(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("What is the page number of a path on %d vertices?", n),
			1,
			"A path is outerplanar, so one page suffices.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("What upper bound does the page number of a planar graph on %d vertices satisfy?", n),
			4,
			"Every planar graph has page number at most four, independent of the order.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("What is the page number of the complete graph K_%d?", n),
			w17CeilDiv(n, 2),
			"The complete graph needs half its order rounded up in pages.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Queue(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("What is the queue number of a path on %d vertices?", n),
			1,
			"A path can be laid out on a line with every edge in a single queue.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 20)
		return w17P(kp,
			fmt.Sprintf("What is the queue number of a tree on %d vertices?", n),
			1,
			"Trees admit a one-queue layout by ordering vertices by a depth-first traversal.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("What is the queue number of the complete graph K_%d?", n),
			n/2,
			"The complete graph needs half its order, rounded down, in queues.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Bandwidth(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("What is the bandwidth of a path on %d vertices?", n),
			1,
			"Ordering the vertices along the path makes every edge join consecutive positions.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(3, 20)
		return w17P(kp,
			fmt.Sprintf("What is the bandwidth of a cycle on %d vertices?", n),
			2,
			"Two opposite arcs around the cycle force some edge to span two positions, and two suffice.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(3, 14)
		return w17P(kp,
			fmt.Sprintf("What is the bandwidth of the complete graph K_%d?", n),
			n-1,
			"The first and last vertices are adjacent, so every edge spans at most n-1 positions.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}

func w17Separator(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(4, 60)
		bound := int(math.Ceil(2 * math.Sqrt(2*float64(n))))
		return w17P(kp,
			fmt.Sprintf("The planar separator theorem gives a separator of at most 2 times the square root of 2n. For n=%d, what is that bound rounded up?", n),
			bound,
			"Evaluate 2 sqrt(2n) and round up to an integer.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(4, 60)
		return w17P(kp,
			fmt.Sprintf("A planar separator leaves components of at most 2n/3 vertices. For n=%d, what is that bound?", n),
			2*n/3,
			"Multiply the order by two and divide by three, rounding down.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		k := w17N(3, 12)
		return w17P(kp,
			fmt.Sprintf("A square grid has k by k vertices with k=%d. What separator size does its middle row give?", k),
			k,
			"Remove one complete row of the grid, which has k vertices.",
			map[string]string{"k": strconv.Itoa(k)})
	}
}

func w17Apex(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w17N(5, 16)
		return w17P(kp,
			fmt.Sprintf("How many vertices must be deleted from K_%d to leave a planar graph?", n),
			n-4,
			"K_5 is the smallest nonplanar complete graph, so deleting n-4 vertices leaves K_4.",
			map[string]string{"n": strconv.Itoa(n)})
	case 2:
		n := w17N(5, 16)
		return w17P(kp,
			fmt.Sprintf("Deleting one vertex from K_%d leaves how many edges?", n),
			w17C2(n-1),
			"The remaining graph is K_(n-1), which has (n-1)(n-2)/2 edges.",
			map[string]string{"n": strconv.Itoa(n)})
	default:
		n := w17N(5, 20)
		return w17P(kp,
			fmt.Sprintf("An apex graph on %d vertices deletes one vertex to become planar. What is its maximum number of edges?", n),
			3*(n-1)-6,
			"The remaining planar graph has at most 3(n-1)-6 edges.",
			map[string]string{"n": strconv.Itoa(n)})
	}
}
