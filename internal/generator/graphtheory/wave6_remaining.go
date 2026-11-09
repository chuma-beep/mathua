package graphtheory

// Wave 6 graph theory, remaining batch: interval graphs and split graphs. Both
// generators publish the intervals or the edge list so wave6_remaining_test.go
// can recompute the answer without trusting the generator.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"graph.interval": "interval",
		"graph.split":    "split",
	} {
		reg.Register(id, &w6GraphRemainingGen{kind: kind})
	}
}

type w6GraphRemainingGen struct{ kind string }

func (g *w6GraphRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "interval":
		return w6Interval(kp)
	case "split":
		return w6Split(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6Interval(kp int) generator.Problem {
	n := w6GR(4, 6)
	type interval struct{ l, r int }
	ivs := make([]interval, n)
	for i := range ivs {
		l := w6GR(0, 8)
		r := l + w6GR(1, 5)
		ivs[i] = interval{l, r}
	}
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	edgeCount := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if ivs[i].l <= ivs[j].r && ivs[j].l <= ivs[i].r {
				adj[i][j], adj[j][i] = true, true
				edgeCount++
			}
		}
	}
	parts := make([]string, n)
	for i, iv := range ivs {
		parts[i] = fmt.Sprintf("%d,%d", iv.l, iv.r)
	}
	intervalText := strings.Join(parts, ";")
	f := map[string]string{"n": strconv.Itoa(n), "intervals": intervalText}
	switch kp {
	case 1:
		best := 0
		for mask := 1; mask < 1<<n; mask++ {
			var sub []int
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					sub = append(sub, i)
				}
			}
			ok := true
			for a := 0; a < len(sub) && ok; a++ {
				for b := a + 1; b < len(sub); b++ {
					if !adj[sub[a]][sub[b]] {
						ok = false
						break
					}
				}
			}
			if ok && len(sub) > best {
				best = len(sub)
			}
		}
		return w6G(kp, fmt.Sprintf("The intervals are %s, where each pair is l,r. What is the clique number of their intersection graph?", intervalText), strconv.Itoa(best),
			fmt.Sprintf("The largest set of pairwise overlapping intervals has size %d.", best), f)
	case 2:
		return w6G(kp, fmt.Sprintf("The intervals are %s, where each pair is l,r. How many edges does their intersection graph have?", intervalText), strconv.Itoa(edgeCount),
			fmt.Sprintf("Count the pairs of intervals that overlap. There are %d such pairs.", edgeCount), f)
	default:
		size := w6GR(2, 4)
		perm := rand.Perm(n)[:size]
		ok := 0
		clique := true
		for a := 0; a < len(perm) && clique; a++ {
			for b := a + 1; b < len(perm); b++ {
				if !adj[perm[a]][perm[b]] {
					clique = false
					break
				}
			}
		}
		if clique {
			ok = 1
		}
		f["subset"] = w6GJoin(perm)
		return w6G(kp, fmt.Sprintf("For the intervals %s, is the set %v a clique in the intersection graph? Enter 1 for yes and 0 for no.", intervalText, perm), strconv.Itoa(ok),
			fmt.Sprintf("Check every pair in the set for overlap. The answer is %d.", ok), f)
	}
}

func w6Split(kp int) generator.Problem {
	n := w6GR(5, 6)
	c := w6GR(2, 4)
	perm := rand.Perm(n)
	clique := perm[:c]
	independent := perm[c:]
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	add := func(u, v int) { adj[u][v], adj[v][u] = true, true }
	for i := 0; i < len(clique); i++ {
		for j := i + 1; j < len(clique); j++ {
			add(clique[i], clique[j])
		}
	}
	for _, u := range independent {
		// Keep the maximum clique equal to c by making each independent vertex
		// miss at least one vertex of the designated clique.
		for i, v := range clique {
			if i == len(clique)-1 && len(clique) > 1 {
				continue
			}
			if rand.Intn(2) == 0 {
				add(u, v)
			}
		}
	}
	// For the split test, sometimes replace the graph with a C5 (plus an
	// isolated vertex when n=6), which is not split.
	if kp == 3 && rand.Intn(2) == 0 {
		adj = make([][]bool, n)
		for i := range adj {
			adj[i] = make([]bool, n)
		}
		for i := 0; i < 5; i++ {
			add(i, (i+1)%5)
		}
	}
	f := map[string]string{"n": strconv.Itoa(n), "edges": w6GEdges(adj)}
	switch kp {
	case 1:
		best := 0
		for mask := 1; mask < 1<<n; mask++ {
			var sub []int
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					sub = append(sub, i)
				}
			}
			ok := true
			for a := 0; a < len(sub) && ok; a++ {
				for b := a + 1; b < len(sub); b++ {
					if !adj[sub[a]][sub[b]] {
						ok = false
						break
					}
				}
			}
			if ok && len(sub) > best {
				best = len(sub)
			}
		}
		return w6G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, what is the clique number?", n-1, f["edges"]), strconv.Itoa(best),
			fmt.Sprintf("The largest all-adjacent set has size %d.", best), f)
	case 2:
		size := w6GR(2, 3)
		sub := rand.Perm(n)[:size]
		ok := 1
		for i := 0; i < len(sub); i++ {
			for j := i + 1; j < len(sub); j++ {
				if adj[sub[i]][sub[j]] {
					ok = 0
				}
			}
		}
		f["subset"] = w6GJoin(sub)
		return w6G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, is the set %v independent? Enter 1 for yes and 0 for no.", n-1, f["edges"], sub), strconv.Itoa(ok),
			fmt.Sprintf("Check that no two vertices in the set are joined. The answer is %d.", ok), f)
	default:
		isSplit := w6IsSplitGraph(adj)
		answer := 0
		if isSplit {
			answer = 1
		}
		return w6G(kp, fmt.Sprintf("For the graph with vertices 0..%d and edges %s, is the graph split? Enter 1 for yes and 0 for no.", n-1, f["edges"]), strconv.Itoa(answer),
			fmt.Sprintf("A graph is split when its vertices can be partitioned into a clique and an independent set. The answer is %d.", answer), f)
	}
}

func w6IsSplitGraph(adj [][]bool) bool {
	n := len(adj)
	for mask := 0; mask < 1<<n; mask++ {
		var clique, independent []int
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				clique = append(clique, i)
			} else {
				independent = append(independent, i)
			}
		}
		ok := true
		for i := 0; i < len(clique) && ok; i++ {
			for j := i + 1; j < len(clique); j++ {
				if !adj[clique[i]][clique[j]] {
					ok = false
					break
				}
			}
		}
		if !ok {
			continue
		}
		for i := 0; i < len(independent) && ok; i++ {
			for j := i + 1; j < len(independent); j++ {
				if adj[independent[i]][independent[j]] {
					ok = false
					break
				}
			}
		}
		if ok {
			return true
		}
	}
	return false
}
