package graphtheory

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10bArcs(s string) [][2]int {
	var out [][2]int
	for _, part := range strings.Split(s, ", ") {
		var u, v int
		if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d>%d", &u, &v); err == nil {
			out = append(out, [2]int{u, v})
		}
	}
	return out
}

func w10bArcSet(s string) map[[2]int]bool {
	out := map[[2]int]bool{}
	for _, e := range w10bArcs(s) {
		out[e] = true
	}
	return out
}

func w10bDirectedCycles(edges [][2]int, removed map[[2]int]bool) int {
	adj := map[int]int{}
	for _, e := range edges {
		if !removed[e] {
			adj[e[0]] = e[1]
		}
	}
	seen := map[int]bool{}
	cycles := 0
	for start := range adj {
		if seen[start] {
			continue
		}
		cur := start
		path := map[int]bool{}
		for {
			if seen[cur] {
				break
			}
			if path[cur] {
				cycles++
				break
			}
			path[cur] = true
			next, ok := adj[cur]
			if !ok {
				break
			}
			cur = next
		}
		for v := range path {
			seen[v] = true
		}
	}
	return cycles
}

func TestWave10RestGraphAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 40; draw++ {
			p := (&w10gFeedbackArcGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			edges := w10bArcs(p.Facts["edges"])
			cycles := w10bDirectedCycles(edges, nil)
			want := cycles
			if kp == 2 {
				want = w10bDirectedCycles(edges, w10bArcSet(p.Facts["proposed"]))
			} else if kp == 3 {
				want = len(edges) - cycles
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("feedback arc kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gArboricityGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n, m, bound := w10TInt(t, p, "a"), w10TInt(t, p, "b"), w10TInt(t, p, "c")
			want = (m + n - 2) / (n - 1)
			if kp == 2 {
				want = (n + 1) / 2
			} else if kp == 3 {
				want = n - 1
			}
			if p.Answer != strconv.Itoa(want) || (kp == 1 && bound != want) {
				t.Fatalf("arboricity kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gThicknessGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n, m, bound = w10TInt(t, p, "a"), w10TInt(t, p, "b"), w10TInt(t, p, "c")
			want = (m + 3*n - 7) / (3*n - 6)
			if kp == 2 {
				want = 2
				if n >= 9 {
					want = 3
				}
			} else if kp == 3 {
				want = 3*n - 6
			}
			if p.Answer != strconv.Itoa(want) || (kp == 1 && bound != want) {
				t.Fatalf("thickness kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gCrossingGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n, cr, claimed := w10TInt(t, p, "a"), w10TInt(t, p, "b"), w10TInt(t, p, "c")
			table := map[int]int{4: 0, 5: 1, 6: 3, 7: 9, 8: 18}
			want = table[n]
			if kp == 2 {
				want = 1
			} else if kp == 3 {
				want = claimed - cr
			}
			if p.Answer != strconv.Itoa(want) || (kp == 1 && cr != table[n]) {
				t.Fatalf("crossing kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10gCliqueCoverGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			n = w10TInt(t, p, "a")
			want = 1
			if kp == 2 {
				want = n
			} else if kp == 3 {
				want = 2
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("clique cover kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}
		}
	}
}

func TestWave10RestGraphDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"feedback-arc", &w10gFeedbackArcGen{}}, {"arboricity", &w10gArboricityGen{}},
		{"thickness", &w10gThicknessGen{}}, {"crossing", &w10gCrossingGen{}},
		{"clique-cover", &w10gCliqueCoverGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			rand.Seed(int64(20261110 + kp))
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
