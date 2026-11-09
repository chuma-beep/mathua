package graphtheory

// Wave 6 graph theory: chordal graphs and cographs. The generators use small
// graphs and publish the edge list, ordering or cotree expression so the tests
// can recompute every answer.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"graph.chordal": "chordal",
		"graph.cograph": "cograph",
	} {
		reg.Register(id, &w6GraphGen{kind: kind})
	}
}

type w6GraphGen struct{ kind string }

func (g *w6GraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "chordal":
		return w6Chordal(kp)
	case "cograph":
		return w6Cograph(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6G(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w6GI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6GR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6GJoin(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w6GEdges(adj [][]bool) string {
	var parts []string
	for i := 0; i < len(adj); i++ {
		for j := i + 1; j < len(adj); j++ {
			if adj[i][j] {
				parts = append(parts, fmt.Sprintf("%d-%d", i, j))
			}
		}
	}
	return strings.Join(parts, ", ")
}

func w6GRandomGraph(n int, p int) [][]bool {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(10) < p {
				adj[i][j], adj[j][i] = true, true
			}
		}
	}
	return adj
}

func w6GChordalGraph(n int) [][]bool {
	// A tree is chordal, and adding edges preserves chordality.
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	perm := rand.Perm(n)
	for i := 1; i < n; i++ {
		j := perm[rand.Intn(i)]
		u, v := perm[i], j
		adj[u][v], adj[v][u] = true, true
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(4) == 0 {
				adj[i][j], adj[j][i] = true, true
			}
		}
	}
	return adj
}

func w6GIsPEO(adj [][]bool, order []int) bool {
	n := len(adj)
	pos := make([]int, n)
	for i, v := range order {
		pos[v] = i
	}
	for i, u := range order {
		var later []int
		for v := 0; v < n; v++ {
			if adj[u][v] && pos[v] > i {
				later = append(later, v)
			}
		}
		for a := 0; a < len(later); a++ {
			for b := a + 1; b < len(later); b++ {
				if !adj[later[a]][later[b]] {
					return false
				}
			}
		}
	}
	return true
}

func w6GHardcodedPEO(adj [][]bool) []int {
	n := len(adj)
	if n <= 7 {
		var search func([]int, []bool) []int
		search = func(prefix []int, used []bool) []int {
			if len(prefix) == n {
				if w6GIsPEO(adj, prefix) {
					return append([]int(nil), prefix...)
				}
				return nil
			}
			for v := 0; v < n; v++ {
				if used[v] {
					continue
				}
				used[v] = true
				if out := search(append(prefix, v), used); out != nil {
					return out
				}
				used[v] = false
			}
			return nil
		}
		return search(nil, make([]bool, n))
	}
	return nil
}

func w6GHasChordlessCycle(adj [][]bool) int {
	n := len(adj)
	for mask := 0; mask < 1<<n; mask++ {
		if popcount6(mask) < 4 {
			continue
		}
		deg := make([]int, n)
		edges := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < n; j++ {
				if mask&(1<<j) == 0 {
					continue
				}
				if adj[i][j] {
					deg[i]++
					deg[j]++
					edges++
				}
			}
		}
		ok := true
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 && deg[i] != 2 {
				ok = false
			}
		}
		if ok && edges == popcount6(mask) {
			// A union of disjoint cycles also has degree 2 everywhere. The
			// selected subgraph must be connected to be one chordless cycle.
			start := -1
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					start = i
					break
				}
			}
			seen := map[int]bool{start: true}
			queue := []int{start}
			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]
				for v := 0; v < n; v++ {
					if mask&(1<<v) != 0 && adj[u][v] && !seen[v] {
						seen[v] = true
						queue = append(queue, v)
					}
				}
			}
			if len(seen) == popcount6(mask) {
				return 1
			}
		}
	}
	return 0
}

func w6GMaximalCliques(adj [][]bool) int {
	n := len(adj)
	cliques := make([]int, 0, 64)
	for mask := 1; mask < 1<<n; mask++ {
		ok := true
		for i := 0; i < n && ok; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < n; j++ {
				if mask&(1<<j) != 0 && !adj[i][j] {
					ok = false
					break
				}
			}
		}
		if ok {
			cliques = append(cliques, mask)
		}
	}
	count := 0
	for _, a := range cliques {
		maximal := true
		for _, b := range cliques {
			if a != b && a&b == a {
				maximal = false
				break
			}
		}
		if maximal {
			count++
		}
	}
	return count
}

func popcount6(x int) int {
	c := 0
	for x > 0 {
		c += x & 1
		x >>= 1
	}
	return c
}

func w6Chordal(kp int) generator.Problem {
	n := w6GR(5, 6)
	var adj [][]bool
	if kp == 2 {
		adj = w6GRandomGraph(n, 4)
	} else {
		adj = w6GChordalGraph(n)
	}
	f := map[string]string{}
	w6GI(f, "n", n)
	f["edges"] = w6GEdges(adj)
	switch kp {
	case 1:
		order := rand.Perm(n)
		if rand.Intn(2) == 0 {
			if peo := w6GHardcodedPEO(adj); peo != nil {
				order = peo
			}
		}
		f["ordering"] = w6GJoin(order)
		is := 0
		if w6GIsPEO(adj, order) {
			is = 1
		}
		w6GI(f, "peo", is)
		return w6G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, is the ordering %s a perfect elimination ordering? Enter 1 for yes and 0 for no.", n-1, f["edges"], w6GJoin(order)), strconv.Itoa(is),
			fmt.Sprintf("A perfect elimination ordering requires the later neighbours of every vertex to form a clique. The answer is %d.", is), f)
	case 2:
		has := w6GHasChordlessCycle(adj)
		w6GI(f, "chordless_cycle", has)
		return w6G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, does it contain a chordless cycle of length at least 4? Enter 1 for yes and 0 for no.", n-1, f["edges"]), strconv.Itoa(has),
			fmt.Sprintf("A chordless cycle is an induced cycle. The answer is %d.", has), f)
	default:
		count := w6GMaximalCliques(adj)
		w6GI(f, "maximal_cliques", count)
		return w6G(kp, fmt.Sprintf("For the chordal graph with vertices 0..%d and edges %s, how many maximal cliques does it have?", n-1, f["edges"]), strconv.Itoa(count),
			fmt.Sprintf("List every maximal clique and count them. There are %d.", count), f)
	}
}

type w6CNode struct {
	leaf  int
	op    byte
	left  *w6CNode
	right *w6CNode
	verts []int
}

func w6CBuild(vs []int) *w6CNode {
	if len(vs) == 1 {
		return &w6CNode{leaf: vs[0], verts: append([]int(nil), vs...)}
	}
	cut := 1 + rand.Intn(len(vs)-1)
	leftV := append([]int(nil), vs[:cut]...)
	rightV := append([]int(nil), vs[cut:]...)
	op := byte('U')
	if rand.Intn(2) == 0 {
		op = 'J'
	}
	return &w6CNode{op: op, left: w6CBuild(leftV), right: w6CBuild(rightV), verts: append([]int(nil), vs...)}
}

func w6CExpr(t *w6CNode) string {
	if t == nil {
		return ""
	}
	if t.left == nil && t.right == nil {
		return strconv.Itoa(t.leaf)
	}
	return fmt.Sprintf("%c(%s,%s)", t.op, w6CExpr(t.left), w6CExpr(t.right))
}

func w6CGraph(t *w6CNode, n int) [][]bool {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	var fill func(*w6CNode)
	fill = func(node *w6CNode) {
		if node == nil || node.left == nil {
			return
		}
		fill(node.left)
		fill(node.right)
		if node.op == 'J' {
			for _, u := range node.left.verts {
				for _, v := range node.right.verts {
					adj[u][v], adj[v][u] = true, true
				}
			}
		}
	}
	fill(t)
	return adj
}

func w6CInternal(t *w6CNode) int {
	if t == nil || t.left == nil {
		return 0
	}
	return 1 + w6CInternal(t.left) + w6CInternal(t.right)
}

func w6CP4Free(adj [][]bool) int {
	n := len(adj)
	for a := 0; a < n; a++ {
		for b := a + 1; b < n; b++ {
			for c := b + 1; c < n; c++ {
				for d := c + 1; d < n; d++ {
					sub := []int{a, b, c, d}
					deg := make([]int, 4)
					edges := 0
					for i := 0; i < 4; i++ {
						for j := i + 1; j < 4; j++ {
							if adj[sub[i]][sub[j]] {
								deg[i]++
								deg[j]++
								edges++
							}
						}
					}
					if edges == 3 {
						ones, twos := 0, 0
						for _, x := range deg {
							if x == 1 {
								ones++
							}
							if x == 2 {
								twos++
							}
						}
						if ones == 2 && twos == 2 {
							return 0
						}
					}
				}
			}
		}
	}
	return 1
}

func w6Cograph(kp int) generator.Problem {
	n := w6GR(4, 6)
	vs := make([]int, n)
	for i := range vs {
		vs[i] = i
	}
	t := w6CBuild(vs)
	adj := w6CGraph(t, n)
	f := map[string]string{}
	w6GI(f, "n", n)
	f["edges"] = w6GEdges(adj)
	f["expression"] = w6CExpr(t)
	switch kp {
	case 1:
		count := w6CInternal(t)
		w6GI(f, "internal_nodes", count)
		return w6G(kp, fmt.Sprintf("A cograph has cotree expression %s, where U is disjoint union and J is join. How many internal nodes does the cotree have?", f["expression"]), strconv.Itoa(count),
			fmt.Sprintf("A binary cotree with %d leaves has %d internal nodes.", n, count), f)
	case 2:
		m := len(strings.Split(strings.TrimSpace(f["edges"]), ","))
		if f["edges"] == "" {
			m = 0
		}
		ans := n*(n-1)/2 - m
		w6GI(f, "complement_edges", ans)
		return w6G(kp, fmt.Sprintf("For the cograph with cotree expression %s and edges %s, how many edges does its complement have?", f["expression"], f["edges"]), strconv.Itoa(ans),
			fmt.Sprintf("The complete graph has %d edges and the cograph has %d, so the complement has %d.", n*(n-1)/2, m, ans), f)
	default:
		if rand.Intn(2) == 0 {
			adj = w6GRandomGraph(n, 4)
			f["edges"] = w6GEdges(adj)
		}
		is := w6CP4Free(adj)
		w6GI(f, "p4_free", is)
		return w6G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, is it P4-free? Enter 1 for yes and 0 for no.", n-1, f["edges"]), strconv.Itoa(is),
			fmt.Sprintf("A graph is P4-free exactly when it is a cograph. The answer is %d.", is), f)
	}
}
