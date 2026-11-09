package graphtheory

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w3TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w3TestGraph(t *testing.T, p generator.Problem) (int, [][3]int, bool) {
	t.Helper()
	n := w3TestInt(t, p, "n")
	dir, _ := strconv.ParseBool(p.Facts["directed"])
	var edges [][3]int
	if strings.TrimSpace(p.Facts["edges"]) == "" {
		return n, edges, dir
	}
	for _, part := range strings.Split(p.Facts["edges"], ", ") {
		part = strings.TrimSpace(part)
		var u, v, w int
		if dir {
			if _, err := fmt.Sscanf(part, "%d->%d (%d)", &u, &v, &w); err != nil {
				t.Fatalf("bad directed edge %q: %v", part, err)
			}
		} else {
			if _, err := fmt.Sscanf(part, "%d-%d", &u, &v); err != nil {
				t.Fatalf("bad undirected edge %q: %v", part, err)
			}
			w = 1
		}
		edges = append(edges, [3]int{u, v, w})
	}
	return n, edges, dir
}

func w3TestAdj(n int, edges [][3]int, dir bool) [][]int {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		if !dir {
			adj[e[1]] = append(adj[e[1]], e[0])
		}
	}
	for i := range adj {
		sort.Ints(adj[i])
	}
	return adj
}

func w3TestParseList(s string) []int {
	s = strings.TrimSpace(strings.Trim(s, "[]"))
	if s == "" {
		return nil
	}
	fields := strings.Fields(s)
	out := make([]int, len(fields))
	for i, f := range fields {
		out[i], _ = strconv.Atoi(f)
	}
	return out
}

func w3TestParsePrefs(s string) [][]int {
	s = strings.TrimSpace(strings.Trim(s, "[]"))
	rows := strings.Split(s, "] [")
	out := make([][]int, len(rows))
	for i, row := range rows {
		out[i] = w3TestParseList(row)
	}
	return out
}

func TestWave3GraphAnswerRecomputation(t *testing.T) {
	kinds := []string{"dfs", "bfs", "bellman_ford", "floyd_warshall", "union_find", "assignment", "stable_matching", "clique", "dominating_set"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w3GraphGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w3TestGraphAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w3TestGraphAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	n, edges, dir := w3TestGraph(t, p)
	switch kind {
	case "dfs":
		adj := w3TestAdj(n, edges, dir)
		disc := make([]int, n)
		fin := make([]int, n)
		parent := make([]int, n)
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
		visit(w3TestInt(t, p, "start"))
		target := w3TestInt(t, p, "target")
		if kp == 1 {
			return disc[target]
		}
		if kp == 2 {
			return parent[target]
		}
		return fin[target]
	case "bfs":
		adj := w3TestAdj(n, edges, dir)
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		dist[0] = 0
		seen := make([]bool, n)
		seen[0] = true
		queue := []int{0}
		var order []int
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
		target := w3TestInt(t, p, "target")
		if kp == 1 {
			return dist[target]
		}
		if kp == 2 {
			c := 0
			for _, d := range dist {
				if d == dist[target] {
					c++
				}
			}
			return c
		}
		for i, v := range order {
			if v == target {
				return i + 1
			}
		}
		t.Fatalf("target %d is not in the BFS order", target)
	case "bellman_ford":
		src := 0
		dist := make([]int, n)
		prev := make([]int, n)
		for i := range dist {
			dist[i] = 1 << 30
			prev[i] = -1
		}
		dist[src] = 0
		for i := 0; i < n-1; i++ {
			for _, e := range edges {
				if dist[e[0]] < 1<<30 && dist[e[0]]+e[2] < dist[e[1]] {
					dist[e[1]] = dist[e[0]] + e[2]
					prev[e[1]] = e[0]
				}
			}
		}
		if kp == 3 {
			neg := 0
			for _, e := range edges {
				if dist[e[0]] < 1<<30 && dist[e[0]]+e[2] < dist[e[1]] {
					neg = 1
				}
			}
			return neg
		}
		target := w3TestInt(t, p, "target")
		if kp == 1 {
			return dist[target]
		}
		l := 0
		for v := target; prev[v] != -1; v = prev[v] {
			l++
		}
		return l
	case "floyd_warshall":
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
			if e[2] < d[e[0]][e[1]] {
				d[e[0]][e[1]] = e[2]
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
		if kp == 1 {
			i, j := w3TestInt(t, p, "from"), w3TestInt(t, p, "to")
			if d[i][j] >= inf {
				return -1
			}
			return d[i][j]
		}
		if kp == 2 {
			diam := 0
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if d[i][j] < inf && d[i][j] > diam {
						diam = d[i][j]
					}
				}
			}
			return diam
		}
		c := 0
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && d[i][j] < inf {
					c++
				}
			}
		}
		return c
	case "union_find":
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
		for _, op := range strings.Split(p.Facts["operations"], "; ") {
			var a, b int
			if _, err := fmt.Sscanf(op, "union(%d,%d)", &a, &b); err != nil {
				t.Fatalf("bad union op %q: %v", op, err)
			}
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
		if kp == 1 {
			return comps
		}
		target := w3TestInt(t, p, "target")
		root := find(target)
		if kp == 2 {
			return root
		}
		return size[root]
	case "assignment":
		var cost [3][3]int
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				cost[i][j] = w3TestInt(t, p, fmt.Sprintf("c%d%d", i+1, j+1))
			}
		}
		perms := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
		best := 1 << 30
		bestPerm := [3]int{}
		for _, perm := range perms {
			s := cost[0][perm[0]] + cost[1][perm[1]] + cost[2][perm[2]]
			if s < best {
				best = s
				bestPerm = perm
			}
		}
		if kp == 1 {
			return best
		}
		if kp == 2 {
			return bestPerm[0] + 1
		}
		parts := strings.Split(p.Facts["candidate"], ",")
		perm := [3]int{}
		for i, s := range parts {
			perm[i], _ = strconv.Atoi(s)
			perm[i]--
		}
		if cost[0][perm[0]]+cost[1][perm[1]]+cost[2][perm[2]] == best {
			return 1
		}
		return 0
	case "stable_matching":
		men := w3TestParsePrefs(p.Facts["men"])
		women := w3TestParsePrefs(p.Facts["women"])
		nn := w3TestInt(t, p, "n")
		next := make([]int, nn)
		matchM := make([]int, nn)
		matchW := make([]int, nn)
		for i := range matchM {
			matchM[i] = -1
			matchW[i] = -1
		}
		rankW := make([][]int, nn)
		for w := 0; w < nn; w++ {
			rankW[w] = make([]int, nn)
			for r, m := range women[w] {
				rankW[w][m] = r
			}
		}
		queue := make([]int, nn)
		for i := range queue {
			queue[i] = i
		}
		proposals := 0
		for len(queue) > 0 {
			m := queue[0]
			queue = queue[1:]
			if matchM[m] != -1 || next[m] >= nn {
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
		if kp == 1 {
			return matchM[w3TestInt(t, p, "asked_man")-1] + 1
		}
		if kp == 2 {
			return proposals
		}
		if matchM[w3TestInt(t, p, "asked_man")-1] == w3TestInt(t, p, "asked_woman")-1 {
			return 1
		}
		return 0
	case "clique":
		adj := w3TestAdj(n, edges, dir)
		isClique := func(sub []int) bool {
			for i := 0; i < len(sub); i++ {
				for j := i + 1; j < len(sub); j++ {
					found := false
					for _, v := range adj[sub[i]] {
						if v == sub[j] {
							found = true
							break
						}
					}
					if !found {
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
		if kp == 1 {
			return best
		}
		if kp == 2 {
			if isClique(w3TestParseList(p.Facts["subset"])) {
				return 1
			}
			return 0
		}
		return count
	case "dominating_set":
		adj := make([][]bool, n)
		for i := range adj {
			adj[i] = make([]bool, n)
			adj[i][i] = true
		}
		for _, e := range edges {
			adj[e[0]][e[1]] = true
			adj[e[1]][e[0]] = true
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
		if kp == 2 {
			if dominates(w3TestParseList(p.Facts["subset"])) {
				return 1
			}
			return 0
		}
		v := -1
		if kp == 3 {
			v = w3TestInt(t, p, "vertex")
		}
		best := n + 1
		for mask := 1; mask < 1<<n; mask++ {
			if kp == 3 && mask&(1<<v) == 0 {
				continue
			}
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
		return best
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave3GraphQuestionDiversity(t *testing.T) {
	kinds := []string{"dfs", "bfs", "bellman_ford", "floyd_warshall", "union_find", "assignment", "stable_matching", "clique", "dominating_set"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w3GraphGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
