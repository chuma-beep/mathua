package graphtheory

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6TestEdges(t *testing.T, n int, s string) [][]bool {
	t.Helper()
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	if strings.TrimSpace(s) == "" {
		return adj
	}
	for _, part := range strings.Split(s, ", ") {
		var u, v int
		if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d-%d", &u, &v); err != nil {
			t.Fatalf("bad edge %q: %v", part, err)
		}
		adj[u][v], adj[v][u] = true, true
	}
	return adj
}

func w6TestOrder(t *testing.T, s string) []int {
	t.Helper()
	var out []int
	for _, part := range strings.Split(s, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("bad ordering %q", s)
		}
		out = append(out, v)
	}
	return out
}

func w6TestPEO(adj [][]bool, order []int) bool {
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

func w6TestPopcount(x int) int {
	c := 0
	for x > 0 {
		c += x & 1
		x >>= 1
	}
	return c
}

func w6TestChordlessCycle(adj [][]bool) int {
	n := len(adj)
	for mask := 0; mask < 1<<n; mask++ {
		if w6TestPopcount(mask) < 4 {
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
		if ok && edges == w6TestPopcount(mask) {
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
			if len(seen) == w6TestPopcount(mask) {
				return 1
			}
		}
	}
	return 0
}

func w6TestMaxCliques(adj [][]bool) int {
	n := len(adj)
	var cliques []int
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

func w6TestP4Free(adj [][]bool) int {
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

func TestWave6GraphAnswerRecomputation(t *testing.T) {
	kinds := []string{"chordal", "cograph"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6GraphGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				n := w6TestInt(t, p, "n")
				adj := w6TestEdges(t, n, p.Facts["edges"])
				want := 0
				switch kind {
				case "chordal":
					switch kp {
					case 1:
						if w6TestPEO(adj, w6TestOrder(t, p.Facts["ordering"])) {
							want = 1
						}
					case 2:
						want = w6TestChordlessCycle(adj)
					default:
						want = w6TestMaxCliques(adj)
					}
				case "cograph":
					switch kp {
					case 1:
						want = strings.Count(p.Facts["expression"], "U(") + strings.Count(p.Facts["expression"], "J(")
					case 2:
						m := 0
						if strings.TrimSpace(p.Facts["edges"]) != "" {
							m = len(strings.Split(p.Facts["edges"], ", "))
						}
						want = n*(n-1)/2 - m
					default:
						want = w6TestP4Free(adj)
					}
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func TestWave6GraphDiversity(t *testing.T) {
	kinds := []string{"chordal", "cograph"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w6GraphGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws, below the 10 target", kind, kp, len(seen))
			}
		}
	}
}
