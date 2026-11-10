package graphtheory

// Wave 10 graph theory, part two. Every branch publishes a, b, c and answer.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10Rest(reg *generator.Registry) {
	reg.Register("graph.feedback_arc_set", &w10gFeedbackArcGen{})
	reg.Register("graph.arboricity", &w10gArboricityGen{})
	reg.Register("graph.thickness", &w10gThicknessGen{})
	reg.Register("graph.crossing_number", &w10gCrossingGen{})
	reg.Register("graph.clique_cover", &w10gCliqueCoverGen{})
}

func w10gP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10gI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w10gArcDisplay(edges [][2]int) string {
	if len(edges) == 0 {
		return "empty"
	}
	return w10gArcString(edges)
}

func w10gArcString(edges [][2]int) string {
	parts := make([]string, len(edges))
	for i, e := range edges {
		parts[i] = fmt.Sprintf("%d>%d", e[0], e[1])
	}
	return strings.Join(parts, ", ")
}

type w10gFeedbackArcGen struct{}

func (g *w10gFeedbackArcGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	k := 1 + rand.Intn(3)
	var edges [][2]int
	var cycles [][]int
	next := 0
	totalEdges := 0
	for c := 0; c < k; c++ {
		length := 3 + rand.Intn(3)
		vs := make([]int, length)
		for i := 0; i < length; i++ {
			vs[i] = next
			next++
		}
		cycles = append(cycles, vs)
		for i := 0; i < length; i++ {
			edges = append(edges, [2]int{vs[i], vs[(i+1)%length]})
			totalEdges++
		}
	}
	proposed := make([][2]int, 0, k)
	remainingCycles := 0
	for _, vs := range cycles {
		if o != 2 || rand.Intn(2) == 0 {
			i := rand.Intn(len(vs))
			proposed = append(proposed, [2]int{vs[i], vs[(i+1)%len(vs)]})
		} else {
			remainingCycles++
		}
	}
	facts := map[string]string{}
	w10gI(facts, "a", k)
	w10gI(facts, "b", remainingCycles)
	w10gI(facts, "c", totalEdges-k)
	facts["edges"] = w10gArcString(edges)
	facts["proposed"] = w10gArcString(proposed)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A directed graph is a disjoint union of %d directed cycles with arcs %s. What is the minimum number of arcs to remove to make it acyclic?", k, w10gArcString(edges)), strconv.Itoa(k), fmt.Sprintf("Each directed cycle needs at least one removed arc, and one per cycle suffices. The minimum is %d.", k), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A directed graph is a disjoint union of %d directed cycles with arcs %s. A proposed feedback arc set is %s. How many directed cycles remain?", k, w10gArcString(edges), w10gArcDisplay(proposed)), strconv.Itoa(remainingCycles), fmt.Sprintf("Check each directed cycle for a removed arc. The number left intact is %d.", remainingCycles), facts)
	default:
		return w10gP(o, fmt.Sprintf("A directed graph is a disjoint union of %d directed cycles with %d arcs. After removing one arc from each cycle, how many arcs remain?", k, totalEdges), strconv.Itoa(totalEdges-k), fmt.Sprintf("Remove one arc per cycle: %d - %d = %d.", totalEdges, k, totalEdges-k), facts)
	}
}

type w10gArboricityGen struct{}

func (g *w10gArboricityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	n := 2 + rand.Intn(14)
	m := n - 1 + rand.Intn(3*n)
	bound := (m + n - 2) / (n - 1)
	complete := (n + 1) / 2
	facts := map[string]string{}
	w10gI(facts, "a", n)
	w10gI(facts, "b", m)
	w10gI(facts, "c", bound)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A graph has %d vertices and %d edges. What is the arboricity lower bound ceil(m/(n-1))?", n, m), strconv.Itoa(bound), fmt.Sprintf("Divide edges by vertices minus one and round up: %d/%d rounds up to %d.", m, n-1, bound), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("What is the arboricity of the complete graph K_%d?", n), strconv.Itoa(complete), fmt.Sprintf("The arboricity of K_n is ceil(n/2): ceil(%d/2) = %d.", n, complete), facts)
	default:
		return w10gP(o, fmt.Sprintf("A forest on %d vertices has how many edges at most?", n), strconv.Itoa(n-1), fmt.Sprintf("A forest with n vertices has at most n-1 edges: %d - 1 = %d.", n, n-1), facts)
	}
}

type w10gThicknessGen struct{}

func (g *w10gThicknessGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	n := 4 + rand.Intn(17)
	if o == 2 {
		n = 5 + rand.Intn(6)
	}
	m := 3*n + rand.Intn(5*n)
	bound := (m + 3*n - 7) / (3*n - 6)
	complete := 2
	if n >= 9 {
		complete = 3
	}
	facts := map[string]string{}
	w10gI(facts, "a", n)
	w10gI(facts, "b", m)
	w10gI(facts, "c", bound)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A graph has %d vertices and %d edges. What is the planar thickness lower bound ceil(m/(3n-6))?", n, m), strconv.Itoa(bound), fmt.Sprintf("Use 3n-6 = %d and round the quotient up to %d.", 3*n-6, bound), facts)
	case 2:
		claimed := complete + rand.Intn(4)
		facts["claimed"] = strconv.Itoa(claimed)
		return w10gP(o, fmt.Sprintf("A claimed decomposition of K_%d uses %d planar layers. What is its thickness?", n, claimed), strconv.Itoa(complete), fmt.Sprintf("For this complete graph the thickness is %d, so the claimed decomposition is not minimal.", complete), facts)
	default:
		return w10gP(o, fmt.Sprintf("What is the maximum number of edges in a planar graph with %d vertices?", n), strconv.Itoa(3*n-6), fmt.Sprintf("Use 3n-6 = 3(%d) - 6 = %d.", n, 3*n-6), facts)
	}
}

type w10gCrossingGen struct{}

func (g *w10gCrossingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	n := 4 + rand.Intn(5)
	cr := map[int]int{5: 1, 6: 3, 7: 9, 8: 18}[n]
	claimed := cr + 1 + rand.Intn(8)
	facts := map[string]string{}
	w10gI(facts, "a", n)
	w10gI(facts, "b", cr)
	w10gI(facts, "c", claimed)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("A drawing of K_%d is claimed to have %d crossings. What is its crossing number?", n, claimed), strconv.Itoa(cr), fmt.Sprintf("The crossing number of K_%d is %d.", n, cr), facts)
	case 2:
		return w10gP(o, fmt.Sprintf("A drawing of K_{3,3} is claimed to have %d crossings. What is its crossing number?", claimed), "1", "K_{3,3} has crossing number 1, so the claimed drawing is not minimal.", facts)
	default:
		return w10gP(o, fmt.Sprintf("A claimed drawing of K_%d has %d crossings. How many crossings above the minimum is that?", n, claimed), strconv.Itoa(claimed-cr), fmt.Sprintf("Subtract the crossing number %d: %d - %d = %d.", cr, claimed, cr, claimed-cr), facts)
	}
}

type w10gCliqueCoverGen struct{}

func (g *w10gCliqueCoverGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	n := 3 + rand.Intn(18)
	m := 2 + rand.Intn(5)
	facts := map[string]string{}
	w10gI(facts, "a", n)
	w10gI(facts, "b", m)
	switch o {
	case 1:
		return w10gP(o, fmt.Sprintf("What is the minimum clique cover of K_%d?", n), "1", "The whole complete graph is one clique, so the cover number is 1.", facts)
	case 2:
		return w10gP(o, fmt.Sprintf("What is the minimum clique cover of the empty graph on %d vertices?", n), strconv.Itoa(n), "No two vertices are adjacent, so every vertex needs its own clique.", facts)
	default:
		return w10gP(o, fmt.Sprintf("What is the minimum clique cover of the complete bipartite graph K_{%d,%d}?", m, n), "2", "Each side is an independent set, and the complement is two cliques, so the cover number is 2.", facts)
	}
}
