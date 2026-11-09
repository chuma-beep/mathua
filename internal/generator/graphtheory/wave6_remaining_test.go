package graphtheory

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6RemTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6RemTestInts(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' })
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		v, _ := strconv.Atoi(part)
		out = append(out, v)
	}
	return out
}

func w6RemTestParseEdges(s string, n int) [][]bool {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	if strings.TrimSpace(s) == "" {
		return adj
	}
	for _, part := range strings.Split(s, ", ") {
		uv := strings.Split(part, "-")
		if len(uv) != 2 {
			continue
		}
		u, _ := strconv.Atoi(uv[0])
		v, _ := strconv.Atoi(uv[1])
		if u >= 0 && u < n && v >= 0 && v < n {
			adj[u][v], adj[v][u] = true, true
		}
	}
	return adj
}

func w6RemTestMaxClique(adj [][]bool) int {
	n := len(adj)
	best := 0
	for mask := 1; mask < 1<<n; mask++ {
		var sub []int
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sub = append(sub, i)
			}
		}
		ok := true
		for i := 0; i < len(sub) && ok; i++ {
			for j := i + 1; j < len(sub); j++ {
				if !adj[sub[i]][sub[j]] {
					ok = false
					break
				}
			}
		}
		if ok && len(sub) > best {
			best = len(sub)
		}
	}
	return best
}

func w6RemTestIsSplit(adj [][]bool) bool {
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

func TestWave6RemainingGraphAnswerRecomputation(t *testing.T) {
	for _, kind := range []string{"interval", "split"} {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6GraphRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w6RemTestGraphAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6RemTestGraphAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "interval":
		n := w6RemTestInt(t, p, "n")
		raw := strings.Split(p.Facts["intervals"], ";")
		adj := make([][]bool, n)
		for i := range adj {
			adj[i] = make([]bool, n)
		}
		for i := 0; i < n; i++ {
			parts := strings.Split(raw[i], ",")
			li, _ := strconv.Atoi(parts[0])
			ri, _ := strconv.Atoi(parts[1])
			for j := i + 1; j < n; j++ {
				parts2 := strings.Split(raw[j], ",")
				lj, _ := strconv.Atoi(parts2[0])
				rj, _ := strconv.Atoi(parts2[1])
				if li <= rj && lj <= ri {
					adj[i][j], adj[j][i] = true, true
				}
			}
		}
		if kp == 1 {
			return w6RemTestMaxClique(adj)
		}
		if kp == 2 {
			c := 0
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					if adj[i][j] {
						c++
					}
				}
			}
			return c
		}
		sub := w6RemTestInts(p.Facts["subset"])
		for i := 0; i < len(sub); i++ {
			for j := i + 1; j < len(sub); j++ {
				if !adj[sub[i]][sub[j]] {
					return 0
				}
			}
		}
		return 1
	case "split":
		n := w6RemTestInt(t, p, "n")
		adj := w6RemTestParseEdges(p.Facts["edges"], n)
		if kp == 1 {
			return w6RemTestMaxClique(adj)
		}
		if kp == 2 {
			sub := w6RemTestInts(p.Facts["subset"])
			for i := 0; i < len(sub); i++ {
				for j := i + 1; j < len(sub); j++ {
					if adj[sub[i]][sub[j]] {
						return 0
					}
				}
			}
			return 1
		}
		if w6RemTestIsSplit(adj) {
			return 1
		}
		return 0
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave6RemainingGraphDiversity(t *testing.T) {
	for _, kind := range []string{"interval", "split"} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w6GraphRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Errorf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
