package graphtheory

// Wave 3 graph theory: nine exact-computation concepts. Every generator draws a
// small graph, publishes its vertices and edges in Facts, and asks for one
// integer: a discovery time, a distance, a component count, an optimum, a
// clique number or a domination number. wave3_test.go recomputes each answer
// from the published graph.

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave3(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"graph.dfs":             "dfs",
		"graph.bfs":             "bfs",
		"graph.bellman_ford":    "bellman_ford",
		"graph.floyd_warshall":  "floyd_warshall",
		"graph.union_find":      "union_find",
		"graph.assignment":      "assignment",
		"graph.stable_matching": "stable_matching",
		"graph.clique":          "clique",
		"graph.dominating_set":  "dominating_set",
	} {
		reg.Register(id, &w3GraphGen{kind: kind})
	}
}

type w3GraphGen struct{ kind string }

func (g *w3GraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "dfs":
		return w3DFS(kp)
	case "bfs":
		return w3BFS(kp)
	case "bellman_ford":
		return w3BellmanFord(kp)
	case "floyd_warshall":
		return w3FloydWarshall(kp)
	case "union_find":
		return w3UnionFind(kp)
	case "assignment":
		return w3Assignment(kp)
	case "stable_matching":
		return w3StableMatching(kp)
	case "clique":
		return w3Clique(kp)
	case "dominating_set":
		return w3DominatingSet(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w3G(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w3GI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w3GR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

type w3Edge struct {
	u, v, w int
}

type w3Graph struct {
	n     int
	edges []w3Edge
	dir   bool
}

func (g w3Graph) label() string {
	parts := make([]string, len(g.edges))
	for i, e := range g.edges {
		if g.dir {
			parts[i] = fmt.Sprintf("%d->%d (%d)", e.u, e.v, e.w)
		} else {
			parts[i] = fmt.Sprintf("%d-%d", e.u, e.v)
		}
	}
	return strings.Join(parts, ", ")
}

func (g w3Graph) facts(m map[string]string) {
	w3GI(m, "n", g.n)
	m["edges"] = g.label()
	m["directed"] = strconv.FormatBool(g.dir)
}

func (g w3Graph) adj() [][]int {
	out := make([][]int, g.n)
	for _, e := range g.edges {
		out[e.u] = append(out[e.u], e.v)
		if !g.dir {
			out[e.v] = append(out[e.v], e.u)
		}
	}
	for i := range out {
		sort.Ints(out[i])
	}
	return out
}

func w3RandomConnected(n int) w3Graph {
	var edges []w3Edge
	seen := map[[2]int]bool{}
	add := func(u, v int) {
		if u > v {
			u, v = v, u
		}
		if seen[[2]int{u, v}] {
			return
		}
		seen[[2]int{u, v}] = true
		edges = append(edges, w3Edge{u: u, v: v, w: 1})
	}
	perm := rand.Perm(n)
	for i := 1; i < n; i++ {
		add(perm[i], perm[rand.Intn(i)])
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(4) == 0 {
				add(i, j)
			}
		}
	}
	return w3Graph{n: n, edges: edges}
}

func w3DFS(kp int) generator.Problem {
	g := w3RandomConnected(5)
	adj := g.adj()
	disc := make([]int, g.n)
	fin := make([]int, g.n)
	parent := make([]int, g.n)
	for i := range parent {
		parent[i] = -1
	}
	clock := 0
	var visit func(int)
	visit = func(u int) {
		clock++
		disc[u] = clock
		for _, v := range adj[u] {
			if disc[v] == 0 {
				parent[v] = u
				visit(v)
			}
		}
		clock++
		fin[u] = clock
	}
	visit(0)
	target := w3GR(0, g.n-1)
	f := map[string]string{}
	g.facts(f)
	w3GI(f, "start", 0)
	w3GI(f, "target", target)
	w3GI(f, "discovery", disc[target])
	w3GI(f, "parent", parent[target])
	w3GI(f, "finish", fin[target])
	switch kp {
	case 1:
		return w3G(kp, fmt.Sprintf("A depth-first search starts at vertex 0 and visits neighbours in increasing order. The graph has vertices 0..%d and edges %s. What is the discovery time of vertex %d? (enter a number)", g.n-1, g.label(), target), strconv.Itoa(disc[target]),
			fmt.Sprintf("The clock increments on entry. Vertex %d is discovered at time %d.", target, disc[target]), f)
	case 2:
		p := parent[target]
		return w3G(kp, fmt.Sprintf("For the same DFS on vertices 0..%d with edges %s starting at 0, what is the parent of vertex %d? (enter -1 if it is the root)", g.n-1, g.label(), target), strconv.Itoa(p),
			fmt.Sprintf("The DFS tree records the vertex from which %d was first reached. Its parent is %d.", target, p), f)
	default:
		return w3G(kp, fmt.Sprintf("For the same DFS on vertices 0..%d with edges %s starting at 0, what is the finish time of vertex %d? (enter a number)", g.n-1, g.label(), target), strconv.Itoa(fin[target]),
			fmt.Sprintf("A vertex finishes after every vertex reachable from it has finished. The finish time of %d is %d.", target, fin[target]), f)
	}
}

func w3BFS(kp int) generator.Problem {
	g := w3RandomConnected(5)
	adj := g.adj()
	dist := make([]int, g.n)
	order := make([]int, 0, g.n)
	seen := make([]bool, g.n)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	seen[0] = true
	queue := []int{0}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)
		for _, v := range adj[u] {
			if !seen[v] {
				seen[v] = true
				dist[v] = dist[u] + 1
				queue = append(queue, v)
			}
		}
	}
	target := w3GR(0, g.n-1)
	f := map[string]string{}
	g.facts(f)
	w3GI(f, "start", 0)
	w3GI(f, "target", target)
	w3GI(f, "distance", dist[target])
	w3GI(f, "at_distance", dist[target])
	for i, v := range order {
		w3GI(f, fmt.Sprintf("order_%d", v), i+1)
	}
	switch kp {
	case 1:
		return w3G(kp, fmt.Sprintf("A breadth-first search starts at vertex 0. The graph has vertices 0..%d and edges %s. What is the distance from 0 to vertex %d? (enter a number)", g.n-1, g.label(), target), strconv.Itoa(dist[target]),
			fmt.Sprintf("BFS explores in layers. The distance from 0 to %d is %d.", target, dist[target]), f)
	case 2:
		count := 0
		for _, d := range dist {
			if d == dist[target] {
				count++
			}
		}
		if count == 1 {
			return w3G(kp, fmt.Sprintf("For the same BFS on vertices 0..%d with edges %s starting at 0, how many vertices are at the same distance as vertex %d?", g.n-1, g.label(), target), "1",
				fmt.Sprintf("Vertex %d is the only vertex at distance %d.", target, dist[target]), f)
		}
		return w3G(kp, fmt.Sprintf("For the same BFS on vertices 0..%d with edges %s starting at 0, how many vertices are at the same distance as vertex %d?", g.n-1, g.label(), target), strconv.Itoa(count),
			fmt.Sprintf("The vertices at distance %d are counted in the BFS layer. There are %d.", dist[target], count), f)
	default:
		pos := 0
		for i, v := range order {
			if v == target {
				pos = i + 1
				break
			}
		}
		return w3G(kp, fmt.Sprintf("For the same BFS on vertices 0..%d with edges %s starting at 0, what is the 1-based position of vertex %d in the BFS visit order?", g.n-1, g.label(), target), strconv.Itoa(pos),
			fmt.Sprintf("BFS visits vertices in nondecreasing distance and increasing neighbour order. Vertex %d is visited at position %d.", target, pos), f)
	}
}

func w3BellmanFord(kp int) generator.Problem {
	n := 4
	w1, w2, w3, w4, w5 := w3GR(1, 9), w3GR(1, 9), w3GR(1, 9), w3GR(1, 9), w3GR(1, 9)
	edges := []w3Edge{{0, 1, w1}, {1, 2, w2}, {2, 3, w3}, {0, 2, w4}, {1, 3, w5}}
	hasNeg := false
	if kp == 3 && rand.Intn(2) == 0 {
		edges = append(edges, w3Edge{3, 0, -(w1 + w2 + w3) - w3GR(1, 4)})
		hasNeg = true
	}
	g := w3Graph{n: n, edges: edges, dir: true}
	dist := make([]int, n)
	prev := make([]int, n)
	for i := range dist {
		dist[i] = 1 << 30
		prev[i] = -1
	}
	dist[0] = 0
	for i := 0; i < n-1; i++ {
		for _, e := range edges {
			if dist[e.u] < 1<<30 && dist[e.u]+e.w < dist[e.v] {
				dist[e.v] = dist[e.u] + e.w
				prev[e.v] = e.u
			}
		}
	}
	f := map[string]string{}
	g.facts(f)
	w3GI(f, "start", 0)
	if kp == 3 {
		is := 0
		if hasNeg {
			is = 1
		}
		w3GI(f, "negative_cycle", is)
		return w3G(kp, fmt.Sprintf("Run Bellman-Ford on the directed graph with vertices 0..%d and edges %s. Does it contain a negative-weight cycle reachable from 0? Enter 1 for yes and 0 for no.", n-1, g.label()), strconv.Itoa(is),
			fmt.Sprintf("Relax every edge %d times. A further improvement on pass %d would prove a negative cycle. The answer is %d.", n, n, is), f)
	}
	target := w3GR(1, n-1)
	pathLen := 0
	for v := target; prev[v] != -1; v = prev[v] {
		pathLen++
	}
	w3GI(f, "target", target)
	w3GI(f, "distance", dist[target])
	w3GI(f, "path_edges", pathLen)
	if kp == 1 {
		return w3G(kp, fmt.Sprintf("Run Bellman-Ford from vertex 0 on the directed graph with vertices 0..%d and edges %s. What is the shortest distance to vertex %d?", n-1, g.label(), target), strconv.Itoa(dist[target]),
			fmt.Sprintf("Relax the edges repeatedly. The shortest distance from 0 to %d is %d.", target, dist[target]), f)
	}
	return w3G(kp, fmt.Sprintf("For the directed graph with vertices 0..%d and edges %s, how many edges does the Bellman-Ford shortest path from 0 to vertex %d use?", n-1, g.label(), target), strconv.Itoa(pathLen),
		fmt.Sprintf("Follow the predecessor pointers from %d back to 0. The path uses %d edges.", target, pathLen), f)
}

func w3FloydWarshall(kp int) generator.Problem {
	n := 4
	edges := []w3Edge{{0, 1, w3GR(1, 9)}, {1, 2, w3GR(1, 9)}, {2, 3, w3GR(1, 9)}, {3, 0, w3GR(1, 9)}}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j && rand.Intn(5) == 0 {
				edges = append(edges, w3Edge{i, j, w3GR(1, 9)})
			}
		}
	}
	bestEdge := map[[2]int]int{}
	var order [][2]int
	for _, e := range edges {
		key := [2]int{e.u, e.v}
		if _, ok := bestEdge[key]; !ok {
			order = append(order, key)
		}
		if old, ok := bestEdge[key]; !ok || e.w < old {
			bestEdge[key] = e.w
		}
	}
	edges = edges[:0]
	for _, key := range order {
		edges = append(edges, w3Edge{key[0], key[1], bestEdge[key]})
	}
	g := w3Graph{n: n, edges: edges, dir: true}
	const inf = 1 << 29
	d := make([][]int, n)
	for i := range d {
		d[i] = make([]int, n)
		for j := range d[i] {
			if i == j {
				d[i][j] = 0
			} else {
				d[i][j] = inf
			}
		}
	}
	for _, e := range edges {
		if e.w < d[e.u][e.v] {
			d[e.u][e.v] = e.w
		}
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if d[i][k]+d[k][j] < d[i][j] {
					d[i][j] = d[i][k] + d[k][j]
				}
			}
		}
	}
	f := map[string]string{}
	g.facts(f)
	switch kp {
	case 1:
		i, j := w3GR(0, n-1), w3GR(0, n-1)
		for i == j {
			j = w3GR(0, n-1)
		}
		w3GI(f, "from", i)
		w3GI(f, "to", j)
		ans := d[i][j]
		if ans >= inf {
			ans = -1
		}
		return w3G(kp, fmt.Sprintf("Use Floyd-Warshall on the directed graph with vertices 0..%d and edges %s. What is the shortest distance from %d to %d? (enter -1 if unreachable)", n-1, g.label(), i, j), strconv.Itoa(ans),
			fmt.Sprintf("After considering intermediate vertices in order, the shortest distance from %d to %d is %d.", i, j, ans), f)
	case 2:
		diam := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if d[i][j] < inf && d[i][j] > diam {
					diam = d[i][j]
				}
			}
		}
		return w3G(kp, fmt.Sprintf("For the same directed graph with vertices 0..%d and edges %s, what is the diameter (the largest finite shortest-path distance between distinct vertices)?", n-1, g.label()), strconv.Itoa(diam),
			fmt.Sprintf("Take the maximum of the finite shortest-path distances. The diameter is %d.", diam), f)
	default:
		count := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && d[i][j] < inf {
					count++
				}
			}
		}
		return w3G(kp, fmt.Sprintf("For the same directed graph with vertices 0..%d and edges %s, how many ordered pairs of distinct vertices have a finite shortest path?", n-1, g.label()), strconv.Itoa(count),
			fmt.Sprintf("Count the finite off-diagonal entries of the Floyd-Warshall matrix. There are %d.", count), f)
	}
}

func w3UnionFind(kp int) generator.Problem {
	n := 6
	parent := make([]int, n)
	size := make([]int, n)
	for i := range parent {
		parent[i] = i
		size[i] = 1
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	var ops []string
	for k := 0; k < 4; k++ {
		a, b := w3GR(0, n-1), w3GR(0, n-1)
		ops = append(ops, fmt.Sprintf("union(%d,%d)", a, b))
		ra, rb := find(a), find(b)
		if ra != rb {
			if size[ra] < size[rb] {
				ra, rb = rb, ra
			}
			parent[rb] = ra
			size[ra] += size[rb]
		}
	}
	comps := 0
	for i := range parent {
		if find(i) == i {
			comps++
		}
	}
	target := w3GR(0, n-1)
	root := find(target)
	f := map[string]string{}
	w3GI(f, "n", n)
	f["operations"] = strings.Join(ops, "; ")
	w3GI(f, "components", comps)
	w3GI(f, "target", target)
	w3GI(f, "root", root)
	w3GI(f, "size", size[root])
	switch kp {
	case 1:
		return w3G(kp, fmt.Sprintf("A disjoint-set structure starts with %d singleton sets. Apply %s. How many components remain? (enter a number)", n, strings.Join(ops, ", ")), strconv.Itoa(comps),
			fmt.Sprintf("Each successful union reduces the component count by one. After the operations, %d components remain.", comps), f)
	case 2:
		return w3G(kp, fmt.Sprintf("After the unions %s on vertices 0..%d, what is the representative (root) of vertex %d?", strings.Join(ops, ", "), n-1, target), strconv.Itoa(root),
			fmt.Sprintf("Follow parent pointers from %d until a self-parent is reached. The root is %d.", target, root), f)
	default:
		return w3G(kp, fmt.Sprintf("After the unions %s on vertices 0..%d, what is the size of the component containing vertex %d?", strings.Join(ops, ", "), n-1, target), strconv.Itoa(size[root]),
			fmt.Sprintf("The component rooted at %d contains %d vertices.", root, size[root]), f)
	}
}

func w3Assignment(kp int) generator.Problem {
	perms := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	var cost [3][3]int
	var best int
	var bestPerm [3]int
	for {
		for i := range cost {
			for j := range cost[i] {
				cost[i][j] = w3GR(1, 9)
			}
		}
		best = 1 << 30
		optimal := 0
		for _, p := range perms {
			s := cost[0][p[0]] + cost[1][p[1]] + cost[2][p[2]]
			if s < best {
				best = s
				bestPerm = p
				optimal = 1
			} else if s == best {
				optimal++
			}
		}
		if optimal == 1 {
			break
		}
	}
	f := map[string]string{}
	w3GI(f, "n", 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			w3GI(f, fmt.Sprintf("c%d%d", i+1, j+1), cost[i][j])
		}
	}
	w3GI(f, "minimum", best)
	for i := 0; i < 3; i++ {
		w3GI(f, fmt.Sprintf("assign_%d", i+1), bestPerm[i]+1)
	}
	f["matrix"] = fmt.Sprintf("[[%d,%d,%d],[%d,%d,%d],[%d,%d,%d]]", cost[0][0], cost[0][1], cost[0][2], cost[1][0], cost[1][1], cost[1][2], cost[2][0], cost[2][1], cost[2][2])
	switch kp {
	case 1:
		return w3G(kp, fmt.Sprintf("For the 3x3 assignment cost matrix %s, what is the minimum total cost of a perfect matching? (enter a number)", f["matrix"]), strconv.Itoa(best),
			fmt.Sprintf("Check all six permutations of the three columns. The minimum total cost is %d.", best), f)
	case 2:
		return w3G(kp, fmt.Sprintf("For the assignment matrix %s, in an optimal matching, which column is assigned to row 1? (enter 1, 2 or 3)", f["matrix"]), strconv.Itoa(bestPerm[0]+1),
			fmt.Sprintf("An optimal assignment gives row 1 column %d.", bestPerm[0]+1), f)
	default:
		p := perms[rand.Intn(len(perms))]
		total := cost[0][p[0]] + cost[1][p[1]] + cost[2][p[2]]
		ok := 0
		if total == best {
			ok = 1
		}
		f["candidate"] = fmt.Sprintf("%d,%d,%d", p[0]+1, p[1]+1, p[2]+1)
		return w3G(kp, fmt.Sprintf("For the assignment matrix %s, is the assignment (row1,row2,row3) = (%d,%d,%d) optimal? Enter 1 for yes and 0 for no.", f["matrix"], p[0]+1, p[1]+1, p[2]+1), strconv.Itoa(ok),
			fmt.Sprintf("Its cost is %d, while the optimum is %d, so the answer is %d.", total, best, ok), f)
	}
}

func w3StableMatching(kp int) generator.Problem {
	n := 3
	if rand.Intn(2) == 0 {
		n = 4
	}
	men := make([][]int, n)
	women := make([][]int, n)
	for i := 0; i < n; i++ {
		men[i] = rand.Perm(n)
		women[i] = rand.Perm(n)
	}
	// Gale-Shapley with men proposing.
	next := make([]int, n)
	matchM := make([]int, n)
	matchW := make([]int, n)
	for i := range matchM {
		matchM[i] = -1
		matchW[i] = -1
	}
	rankW := make([][]int, n)
	for w := 0; w < n; w++ {
		rankW[w] = make([]int, n)
		for r, m := range women[w] {
			rankW[w][m] = r
		}
	}
	queue := make([]int, n)
	for i := range queue {
		queue[i] = i
	}
	proposals := 0
	for len(queue) > 0 {
		m := queue[0]
		queue = queue[1:]
		if matchM[m] != -1 || next[m] >= n {
			continue
		}
		w := men[m][next[m]]
		next[m]++
		proposals++
		if matchW[w] == -1 {
			matchM[m], matchW[w] = w, m
		} else if rankW[w][m] < rankW[w][matchW[w]] {
			old := matchW[w]
			matchM[old] = -1
			queue = append(queue, old)
			matchM[m], matchW[w] = w, m
		} else {
			queue = append(queue, m)
		}
	}
	f := map[string]string{}
	w3GI(f, "n", n)
	f["men"] = fmt.Sprintf("%v", men)
	f["women"] = fmt.Sprintf("%v", women)
	w3GI(f, "proposals", proposals)
	for i := 0; i < n; i++ {
		w3GI(f, fmt.Sprintf("partner_%d", i+1), matchM[i]+1)
	}
	switch kp {
	case 1:
		i := w3GR(0, n-1)
		w3GI(f, "asked_man", i+1)
		return w3G(kp, fmt.Sprintf("Run man-proposing Gale-Shapley with %d men and %d women. The men's preferences are %v and the women's are %v (columns list men in preference order). Who is the partner of man %d? (enter 1..%d)", n, n, men, women, i+1, n), strconv.Itoa(matchM[i]+1),
			fmt.Sprintf("Trace the proposals. Man %d ends with woman %d.", i+1, matchM[i]+1), f)
	case 2:
		return w3G(kp, fmt.Sprintf("For the stable-matching instance with %d men and %d women, men's preferences %v and women's preferences %v, how many proposals are made before the algorithm stops?", n, n, men, women), strconv.Itoa(proposals),
			fmt.Sprintf("Count every proposal, including rejected ones. The algorithm makes %d proposals.", proposals), f)
	default:
		i, j := w3GR(0, n-1), w3GR(0, n-1)
		w3GI(f, "asked_man", i+1)
		w3GI(f, "asked_woman", j+1)
		ok := 0
		if matchM[i] == j {
			ok = 1
		}
		return w3G(kp, fmt.Sprintf("For the stable matching with men's preferences %v and women's preferences %v, is man %d matched with woman %d? Enter 1 for yes and 0 for no.", men, women, i+1, j+1), strconv.Itoa(ok),
			fmt.Sprintf("The stable matching gives man %d partner %d, so the answer is %d.", i+1, matchM[i]+1, ok), f)
	}
}

func w3Clique(kp int) generator.Problem {
	n := 5
	edges := []w3Edge{}
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(5) < 2 {
				edges = append(edges, w3Edge{i, j, 1})
				adj[i][j], adj[j][i] = true, true
			}
		}
	}
	isClique := func(sub []int) bool {
		for i := 0; i < len(sub); i++ {
			for j := i + 1; j < len(sub); j++ {
				if !adj[sub[i]][sub[j]] {
					return false
				}
			}
		}
		return true
	}
	best, count := 0, 0
	for mask := 1; mask < 1<<n; mask++ {
		var sub []int
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sub = append(sub, i)
			}
		}
		if !isClique(sub) {
			continue
		}
		if len(sub) > best {
			best = len(sub)
			count = 1
		} else if len(sub) == best {
			count++
		}
	}
	g := w3Graph{n: n, edges: edges}
	f := map[string]string{}
	g.facts(f)
	w3GI(f, "clique_number", best)
	w3GI(f, "max_cliques", count)
	switch kp {
	case 1:
		return w3G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, what is the clique number? (enter a number)", n-1, g.label()), strconv.Itoa(best),
			fmt.Sprintf("The largest set of pairwise adjacent vertices has size %d.", best), f)
	case 2:
		size := w3GR(2, 4)
		sub := rand.Perm(n)[:size]
		ok := 0
		if isClique(sub) {
			ok = 1
		}
		f["subset"] = fmt.Sprintf("%v", sub)
		return w3G(kp, fmt.Sprintf("For the graph with edges %s, is the vertex set %v a clique? Enter 1 for yes and 0 for no.", g.label(), sub), strconv.Itoa(ok),
			fmt.Sprintf("Check every pair in the set. The answer is %d.", ok), f)
	default:
		return w3G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, how many maximum cliques are there?", n-1, g.label()), strconv.Itoa(count),
			fmt.Sprintf("Enumerate all cliques and count those of size %d. There are %d.", best, count), f)
	}
}

func w3DominatingSet(kp int) generator.Problem {
	n := 5
	edges := []w3Edge{}
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
		adj[i][i] = true
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(5) < 3 {
				edges = append(edges, w3Edge{i, j, 1})
				adj[i][j], adj[j][i] = true, true
			}
		}
	}
	dominates := func(sub []int) bool {
		covered := make([]bool, n)
		for _, u := range sub {
			for v := 0; v < n; v++ {
				if adj[u][v] {
					covered[v] = true
				}
			}
		}
		for _, c := range covered {
			if !c {
				return false
			}
		}
		return true
	}
	best := n + 1
	for mask := 1; mask < 1<<n; mask++ {
		var sub []int
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sub = append(sub, i)
			}
		}
		if len(sub) < best && dominates(sub) {
			best = len(sub)
		}
	}
	g := w3Graph{n: n, edges: edges}
	f := map[string]string{}
	g.facts(f)
	w3GI(f, "domination_number", best)
	switch kp {
	case 1:
		return w3G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, what is the domination number? (enter a number)", n-1, g.label()), strconv.Itoa(best),
			fmt.Sprintf("Find the smallest set whose closed neighbourhood is all vertices. Its size is %d.", best), f)
	case 2:
		size := w3GR(1, 3)
		sub := rand.Perm(n)[:size]
		ok := 0
		if dominates(sub) {
			ok = 1
		}
		f["subset"] = fmt.Sprintf("%v", sub)
		return w3G(kp, fmt.Sprintf("For the graph with edges %s, does the set %v dominate the graph? Enter 1 for yes and 0 for no.", g.label(), sub), strconv.Itoa(ok),
			fmt.Sprintf("Take the union of the closed neighbourhoods of the set. The answer is %d.", ok), f)
	default:
		v := w3GR(0, n-1)
		minWith := n + 1
		for mask := 1; mask < 1<<n; mask++ {
			if mask&(1<<v) == 0 {
				continue
			}
			var sub []int
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					sub = append(sub, i)
				}
			}
			if len(sub) < minWith && dominates(sub) {
				minWith = len(sub)
			}
		}
		f["vertex"] = strconv.Itoa(v)
		return w3G(kp, fmt.Sprintf("For the graph with edges %s, what is the smallest size of a dominating set that contains vertex %d?", g.label(), v), strconv.Itoa(minWith),
			fmt.Sprintf("Restrict the search to sets containing %d. The smallest has size %d.", v, minWith), f)
	}
}
