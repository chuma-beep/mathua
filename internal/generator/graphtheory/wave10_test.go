package graphtheory

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10TInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w10TEdges(s string) [][2]int {
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

func w10TSet(s string) map[int]bool {
	out := map[int]bool{}
	s = strings.TrimSpace(s)
	if s == "" || s == "none" {
		return out
	}
	for _, part := range strings.Split(s, ",") {
		if v, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			out[v] = true
		}
	}
	return out
}

func w10TCycleLengths(edges [][2]int, removed map[int]bool) []int {
	adj := map[int]map[int]bool{}
	for _, e := range edges {
		if removed[e[0]] || removed[e[1]] {
			continue
		}
		if adj[e[0]] == nil {
			adj[e[0]] = map[int]bool{}
		}
		if adj[e[1]] == nil {
			adj[e[1]] = map[int]bool{}
		}
		adj[e[0]][e[1]] = true
		adj[e[1]][e[0]] = true
	}
	seen := map[int]bool{}
	var lengths []int
	for v := range adj {
		if seen[v] {
			continue
		}
		stack := []int{v}
		seen[v] = true
		size := 0
		cycle := true
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			size++
			if len(adj[u]) != 2 {
				cycle = false
			}
			for w := range adj[u] {
				if !seen[w] {
					seen[w] = true
					stack = append(stack, w)
				}
			}
		}
		if cycle && size >= 3 {
			lengths = append(lengths, size)
		}
	}
	return lengths
}

func w10TMatch(edges [][2]int, n int) int {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for _, e := range edges {
		adj[e[0]][e[1]], adj[e[1]][e[0]] = true, true
	}
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
				if c := 1 + dfs(mask|(1<<i)|(1<<j)); c > best {
					best = c
				}
			}
		}
		memo[mask] = best
		return best
	}
	return dfs(0)
}

func TestWave10GraphAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 60; draw++ {
			p := (&w10FeedbackVertexSetGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			edges := w10TEdges(p.Facts["edges"])
			cycles := w10TCycleLengths(edges, nil)
			if kp == 1 {
				if p.Answer != strconv.Itoa(len(cycles)) {
					t.Fatalf("fvs kp1: answer = %s, want %d\nfacts: %v", p.Answer, len(cycles), p.Facts)
				}
			}
			if kp == 2 {
				removed := w10TSet(p.Facts["proposed"])
				remaining := w10TCycleLengths(edges, removed)
				if p.Answer != strconv.Itoa(len(remaining)) {
					t.Fatalf("fvs kp2: answer = %s, want %d\nfacts: %v", p.Answer, len(remaining), p.Facts)
				}
				edgeCount := 0
				for _, e := range edges {
					if !removed[e[0]] && !removed[e[1]] {
						edgeCount++
					}
				}
				if w10TInt(t, p, "remaining_edges") != edgeCount {
					t.Fatalf("fvs kp2: remaining_edges fact %d, want %d", w10TInt(t, p, "remaining_edges"), edgeCount)
				}
			}
			if kp == 3 {
				removed := w10TSet(p.Facts["proposed"])
				edgeCount := 0
				for _, e := range edges {
					if !removed[e[0]] && !removed[e[1]] {
						edgeCount++
					}
				}
				if p.Answer != strconv.Itoa(edgeCount) {
					t.Fatalf("fvs kp3: answer = %s, want %d\nfacts: %v", p.Answer, edgeCount, p.Facts)
				}
				if len(removed) != len(cycles) {
					t.Fatalf("fvs kp3: removed set size %d, want one per cycle %d", len(removed), len(cycles))
				}
			}

			p = (&w10EdgeCoverGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n := w10TInt(t, p, "n")
			edges = w10TEdges(p.Facts["edges"])
			matching := w10TMatch(edges, n)
			if matching != w10TInt(t, p, "max_matching") {
				t.Fatalf("edge cover kp%d: max_matching fact %d, want %d", kp, w10TInt(t, p, "max_matching"), matching)
			}
			cover := n - matching
			if cover != w10TInt(t, p, "edge_cover") {
				t.Fatalf("edge cover kp%d: edge_cover fact %d, want %d", kp, w10TInt(t, p, "edge_cover"), cover)
			}
			switch kp {
			case 1:
				if p.Answer != strconv.Itoa(cover) {
					t.Fatalf("edge cover kp1: answer = %s, want %d", p.Answer, cover)
				}
			case 2:
				proposed := w10TEdges(p.Facts["proposed"])
				covered := make([]bool, n)
				for _, e := range proposed {
					covered[e[0]], covered[e[1]] = true, true
				}
				uncovered := 0
				for _, c := range covered {
					if !c {
						uncovered++
					}
				}
				if p.Answer != strconv.Itoa(uncovered) {
					t.Fatalf("edge cover kp2: answer = %s, want %d\nfacts: %v", p.Answer, uncovered, p.Facts)
				}
			case 3:
				if p.Answer != strconv.Itoa(matching) {
					t.Fatalf("edge cover kp3: answer = %s, want %d", p.Answer, matching)
				}
			}
		}
	}
}

func TestWave10GraphQuestionDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"feedback-vertex-set", &w10FeedbackVertexSetGen{}},
		{"edge-cover", &w10EdgeCoverGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				seen[tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
			}
			if len(seen) < 10 {
				t.Fatalf("%s kp%d: only %d distinct questions in 200 draws", tc.name, kp, len(seen))
			}
			short := map[string]bool{}
			for i := 0; i < 12; i++ {
				short[tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
			}
			if len(short) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", tc.name, kp, len(short))
			}
		}
	}
}
