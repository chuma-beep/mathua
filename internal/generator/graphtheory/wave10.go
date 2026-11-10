package graphtheory

// Wave 10 graph theory: minimum feedback vertex sets and edge covers. Both
// generators publish the full edge list, so a test can rebuild the graph and
// recompute the answer without trusting the generator's arithmetic.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10(reg *generator.Registry) {
	reg.Register("graph.feedback_vertex_set", &w10FeedbackVertexSetGen{})
	reg.Register("graph.edge_cover", &w10EdgeCoverGen{})
	registerWave10Rest(reg)
}

func w10P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10I(m map[string]string, k string, v int) {
	m[k] = strconv.Itoa(v)
}

func w10Ord(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

func w10EdgesString(edges [][2]int) string {
	parts := make([]string, len(edges))
	for i, e := range edges {
		parts[i] = fmt.Sprintf("%d-%d", e[0], e[1])
	}
	return strings.Join(parts, ", ")
}

func w10ParseEdges(s string) [][2]int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out [][2]int
	for _, part := range strings.Split(s, ", ") {
		var u, v int
		if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d-%d", &u, &v); err == nil {
			out = append(out, [2]int{u, v})
		}
	}
	return out
}

func w10IntSet(s string) map[int]bool {
	out := map[int]bool{}
	s = strings.TrimSpace(s)
	if s == "" || s == "none" {
		return out
	}
	s = strings.Trim(s, "{}")
	for _, part := range strings.Split(s, ",") {
		if v, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			out[v] = true
		}
	}
	return out
}

func w10SetString(xs []int) string {
	if len(xs) == 0 {
		return "none"
	}
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w10SetDisplay(xs []int) string {
	if len(xs) == 0 {
		return "empty"
	}
	return "{" + strings.ReplaceAll(w10SetString(xs), ",", ", ") + "}"
}

func w10VertexWord(n int) string {
	if n == 1 {
		return "vertex"
	}
	return "vertices"
}

// ---------------------------------------------------------------------------
// 1. Minimum feedback vertex sets
// ---------------------------------------------------------------------------

type w10FeedbackVertexSetGen struct{}

func (g *w10FeedbackVertexSetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	k := 1 + rand.Intn(3)
	lengths := make([]int, k)
	var edges [][2]int
	var cycles [][]int
	next := 0
	for c := 0; c < k; c++ {
		length := 3 + rand.Intn(3)
		lengths[c] = length
		vertices := make([]int, length)
		for i := 0; i < length; i++ {
			vertices[i] = next
			next++
		}
		cycles = append(cycles, vertices)
		for i := 0; i < length; i++ {
			edges = append(edges, [2]int{vertices[i], vertices[(i+1)%length]})
		}
	}
	forestEdges := 0
	for _, length := range lengths {
		forestEdges += length - 2
	}
	proposed := []int{}
	remainingCycles := k
	switch o {
	case 1:
		// The minimum is one vertex from each cycle.
	case 2:
		for _, vertices := range cycles {
			if rand.Intn(10) < 7 {
				proposed = append(proposed, vertices[rand.Intn(len(vertices))])
			}
		}
		remainingCycles = 0
		for _, vertices := range cycles {
			hit := false
			for _, v := range vertices {
				for _, p := range proposed {
					if v == p {
						hit = true
					}
				}
			}
			if !hit {
				remainingCycles++
			}
		}
	default:
		for _, vertices := range cycles {
			proposed = append(proposed, vertices[rand.Intn(len(vertices))])
		}
		remainingCycles = 0
	}
	remainingEdges := forestEdges
	if o == 2 {
		removed := map[int]bool{}
		for _, v := range proposed {
			removed[v] = true
		}
		remainingEdges = 0
		for _, e := range edges {
			if !removed[e[0]] && !removed[e[1]] {
				remainingEdges++
			}
		}
	}
	facts := map[string]string{}
	w10I(facts, "cycles", k)
	w10I(facts, "remaining_cycles", remainingCycles)
	w10I(facts, "remaining_edges", remainingEdges)
	facts["lengths"] = w10SetString(lengths)
	facts["edges"] = w10EdgesString(edges)
	facts["proposed"] = w10SetString(proposed)
	switch o {
	case 1:
		return w10P(o,
			fmt.Sprintf("A graph is a disjoint union of cycles whose lengths are %s. What is the minimum number of vertices that must be removed to make the graph acyclic?", strings.ReplaceAll(w10SetString(lengths), ",", ", ")),
			strconv.Itoa(k),
			fmt.Sprintf("Each cycle needs at least one removed vertex, and removing one vertex from each cycle leaves a forest. The minimum is %d.", k),
			facts)
	case 2:
		return w10P(o,
			fmt.Sprintf("A graph is a disjoint union of cycles with edge list %s. A proposed feedback vertex set is %s. How many cycles remain after removing it?", w10EdgesString(edges), w10SetDisplay(proposed)),
			strconv.Itoa(remainingCycles),
			fmt.Sprintf("Check each cycle for a removed vertex. The number of cycles that contain none of the proposed vertices is %d.", remainingCycles),
			facts)
	default:
		return w10P(o,
			fmt.Sprintf("A graph is a disjoint union of cycles with lengths %s. Remove one vertex from each cycle. How many edges remain?", strings.ReplaceAll(w10SetString(lengths), ",", ", ")),
			strconv.Itoa(forestEdges),
			fmt.Sprintf("A cycle of length L becomes a path with L-2 edges after one vertex is removed. The total is %d.", forestEdges),
			facts)
	}
}

// ---------------------------------------------------------------------------
// 2. Edge covers
// ---------------------------------------------------------------------------

type w10EdgeCoverGen struct{}

func w10NewGraph(n int) ([][2]int, [][]bool) {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	var edges [][2]int
	add := func(u, v int) {
		if u == v || adj[u][v] {
			return
		}
		adj[u][v], adj[v][u] = true, true
		edges = append(edges, [2]int{u, v})
	}
	for i := 1; i < n; i++ {
		add(i, rand.Intn(i))
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(10) < 3 {
				add(i, j)
			}
		}
	}
	return edges, adj
}

func w10MaxMatching(adj [][]bool) int {
	n := len(adj)
	memo := map[int]int{}
	var dfs func(int) int
	dfs = func(mask int) int {
		if v, ok := memo[mask]; ok {
			return v
		}
		i := 0
		for i < n && mask&(1<<i) != 0 {
			i++
		}
		if i == n {
			return 0
		}
		best := dfs(mask | (1 << i))
		for j := 0; j < n; j++ {
			if j != i && adj[i][j] && mask&(1<<j) == 0 {
				if candidate := 1 + dfs(mask|(1<<i)|(1<<j)); candidate > best {
					best = candidate
				}
			}
		}
		memo[mask] = best
		return best
	}
	return dfs(0)
}

func (g *w10EdgeCoverGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	n := 4 + rand.Intn(4)
	edges, adj := w10NewGraph(n)
	matching := w10MaxMatching(adj)
	cover := n - matching
	facts := map[string]string{}
	w10I(facts, "n", n)
	w10I(facts, "max_matching", matching)
	w10I(facts, "edge_cover", cover)
	facts["edges"] = w10EdgesString(edges)
	switch o {
	case 1:
		w10I(facts, "uncovered", 0)
		facts["proposed"] = w10EdgesString(edges)
		return w10P(o,
			fmt.Sprintf("A graph has %d vertices and edge list %s. No vertex is isolated. What is the minimum size of an edge cover?", n, w10EdgesString(edges)),
			strconv.Itoa(cover),
			fmt.Sprintf("A maximum matching has size %d. The minimum edge cover is n - matching = %d - %d = %d.", matching, n, matching, cover),
			facts)
	case 2:
		shuffled := append([][2]int(nil), edges...)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		limit := len(shuffled)
		if limit > 1 {
			limit--
		}
		count := 1 + rand.Intn(limit)
		proposed := shuffled[:count]
		covered := make([]bool, n)
		for _, e := range proposed {
			covered[e[0]] = true
			covered[e[1]] = true
		}
		uncovered := 0
		for _, c := range covered {
			if !c {
				uncovered++
			}
		}
		if uncovered == 0 && rand.Intn(2) == 0 && len(proposed) > 1 {
			proposed = proposed[:len(proposed)-1]
			covered = make([]bool, n)
			for _, e := range proposed {
				covered[e[0]] = true
				covered[e[1]] = true
			}
			uncovered = 0
			for _, c := range covered {
				if !c {
					uncovered++
				}
			}
		}
		facts["proposed"] = w10EdgesString(proposed)
		w10I(facts, "uncovered", uncovered)
		return w10P(o,
			fmt.Sprintf("A graph has %d vertices and edge list %s. A proposed set of edges is %s. How many vertices are not covered by this proposed set?", n, w10EdgesString(edges), w10EdgesString(proposed)),
			strconv.Itoa(uncovered),
			fmt.Sprintf("Mark the endpoints of every proposed edge. %d %s remain uncovered.", uncovered, w10VertexWord(uncovered)),
			facts)
	default:
		w10I(facts, "uncovered", 0)
		facts["proposed"] = w10EdgesString(edges)
		return w10P(o,
			fmt.Sprintf("A graph has %d vertices, no isolated vertices and edge list %s. Its minimum edge cover has size %d. What is the size of a maximum matching?", n, w10EdgesString(edges), cover),
			strconv.Itoa(matching),
			fmt.Sprintf("For a graph with no isolated vertices, minimum edge cover + maximum matching = number of vertices. Thus %d - %d = %d.", n, cover, matching),
			facts)
	}
}
