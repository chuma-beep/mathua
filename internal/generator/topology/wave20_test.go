package topology

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave20IDs() []string {
	return []string{
		"topo.surface.connected_sum", "topo.surface.crosscap",
		"topo.surface.handle", "topo.surface.euler",
		"topo.surface.punctured", "topo.surface.boundary",
		"topo.surface.double_cover", "topo.surface.genus_formula",
		"topo.surface.boundary_connect", "topo.surface.triangulation",
		"topo.diff.tangent_bundle", "topo.diff.vector_field",
		"topo.diff.critical_point", "topo.diff.morse",
		"topo.diff.sard", "topo.diff.transversality",
		"topo.diff.whitney", "topo.diff.tubular",
		"topo.diff.degree", "topo.diff.flow",
		"topo.hom.chain_complex", "topo.hom.exact_sequence",
		"topo.hom.snake_lemma", "topo.hom.five_lemma",
		"topo.hom.kunneth", "topo.hom.universal_coefficient",
		"topo.hom.euler_poincare", "topo.hom.sphere",
		"topo.hom.projective", "topo.hom.torus",
		"topo.knot.crossing", "topo.knot.unknotting",
		"topo.knot.seifert_genus", "topo.knot.alexander",
		"topo.knot.skein", "topo.knot.torus_knot",
		"topo.knot.connected_sum", "topo.knot.bracket",
		"topo.knot.writhe", "topo.knot.linking",
		"topo.geo.cayley", "topo.geo.growth",
		"topo.geo.quasi_isometry", "topo.geo.ends",
		"topo.geo.svarc_milnor", "topo.geo.hyperbolic",
		"topo.geo.boundary", "topo.geo.dehn",
		"topo.geo.amenable", "topo.geo.property_t",
	}
}

func wave20Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave20Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func wave20TBinom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

func wave20TPow2(n int) int { return 1 << n }

func TestWave20AnswersAreIntegers(t *testing.T) {
	reg := wave20Registry(t)
	for _, id := range wave20IDs() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registry has no generator for %s: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 8; draw++ {
				p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s kp%d: answer %q is not an integer\nquestion: %s", id, kp, p.Answer, p.Question)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: facts answer %q differs from answer %q", id, kp, p.Facts["answer"], p.Answer)
				}
				if p.Explanation == "" {
					t.Fatalf("%s kp%d: empty explanation", id, kp)
				}
				if p.KPOrdinal != kp {
					t.Fatalf("%s: ordinal = %d, want %d", id, p.KPOrdinal, kp)
				}
			}
		}
	}
}

func TestWave20QuestionVariety(t *testing.T) {
	reg := wave20Registry(t)
	for _, id := range wave20IDs() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registry has no generator for %s: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for draw := 0; draw < 12; draw++ {
				p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", id, kp, len(seen))
			}
		}
	}
}

type wave20Case struct {
	name string
	id   string
	kp   int
	want func(t *testing.T, p generator.Problem) int
}

func wave20Cases() []wave20Case {
	cases := []wave20Case{}
	add := func(id string, want func(t *testing.T, p generator.Problem, kp int) int) {
		for kp := 1; kp <= 3; kp++ {
			k := kp
			cases = append(cases, wave20Case{name: id, id: id, kp: kp, want: func(t *testing.T, p generator.Problem) int { return want(t, p, k) }})
		}
	}
	add("topo.surface.connected_sum", func(t *testing.T, p generator.Problem, kp int) int {
		g1, g2 := wave20Fact(t, p, "g1"), wave20Fact(t, p, "g2")
		switch kp {
		case 1:
			return 2 - 2*(g1+g2)
		case 2:
			return wave20Fact(t, p, "chi1") + wave20Fact(t, p, "chi2") - 2
		default:
			return g1 + g2
		}
	})
	add("topo.surface.crosscap", func(t *testing.T, p generator.Problem, kp int) int {
		k := wave20Fact(t, p, "k")
		if kp == 1 {
			return 2 - k
		}
		if kp == 2 {
			return k
		}
		return k - 1
	})
	add("topo.surface.handle", func(t *testing.T, p generator.Problem, kp int) int {
		g := wave20Fact(t, p, "g")
		switch kp {
		case 1:
			return 2 * g
		case 2:
			return 2*g + 2
		default:
			return 1 - 2*g
		}
	})
	add("topo.surface.euler", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 2 {
			chi := wave20Fact(t, p, "chi")
			return (2 - chi) / 2
		}
		g := wave20Fact(t, p, "g")
		if kp == 1 {
			return 2 - 2*g
		}
		return 2 * g
	})
	add("topo.surface.punctured", func(t *testing.T, p generator.Problem, kp int) int {
		g, punctures := wave20Fact(t, p, "g"), wave20Fact(t, p, "p")
		if kp == 2 {
			return 2*g + punctures - 1
		}
		if kp == 3 {
			return 2 - punctures
		}
		return 2 - 2*g - punctures
	})
	add("topo.surface.boundary", func(t *testing.T, p generator.Problem, kp int) int {
		g, b := wave20Fact(t, p, "g"), wave20Fact(t, p, "b")
		if kp == 2 {
			return 2*g + b - 1
		}
		if kp == 3 {
			return 2 - b
		}
		return 2 - 2*g - b
	})
	add("topo.surface.double_cover", func(t *testing.T, p generator.Problem, kp int) int {
		k := wave20Fact(t, p, "k")
		switch kp {
		case 1:
			return k - 1
		case 2:
			return 4 - 2*k
		default:
			return k
		}
	})
	add("topo.surface.genus_formula", func(t *testing.T, p generator.Problem, kp int) int {
		g, k := wave20Fact(t, p, "g"), wave20Fact(t, p, "k")
		switch kp {
		case 1:
			return 2 - 2*g - k
		case 2:
			return 2*g + k
		default:
			return 2 - k
		}
	})
	add("topo.surface.boundary_connect", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20Fact(t, p, "chi1") + wave20Fact(t, p, "chi2") - 1
		case 2:
			return wave20Fact(t, p, "g1") + wave20Fact(t, p, "g2")
		default:
			return 2 - 2*wave20Fact(t, p, "g1") - wave20Fact(t, p, "b1")
		}
	})
	add("topo.surface.triangulation", func(t *testing.T, p generator.Problem, kp int) int {
		chi := wave20Fact(t, p, "chi")
		if kp == 3 {
			return wave20Fact(t, p, "E")/3 + chi
		}
		v := wave20Fact(t, p, "V")
		if kp == 1 {
			return 3 * (v - chi)
		}
		return 2 * (v - chi)
	})
	add("topo.diff.tangent_bundle", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return 2 * n
		}
		return 2*n - 1
	})
	add("topo.diff.vector_field", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			n := wave20Fact(t, p, "n")
			if n%2 == 0 {
				return 2
			}
			return 0
		case 2:
			return wave20Fact(t, p, "sources") - wave20Fact(t, p, "saddles")
		default:
			return 2 - 2*wave20Fact(t, p, "g")
		}
	})
	add("topo.diff.critical_point", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20TPow2(wave20Fact(t, p, "n"))
		case 2:
			return 2*wave20Fact(t, p, "g") + 2
		default:
			return 2
		}
	})
	add("topo.diff.morse", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20TPow2(wave20Fact(t, p, "n"))
		case 2:
			return 2 + 2*wave20Fact(t, p, "g")
		default:
			return wave20Fact(t, p, "n") + 1
		}
	})
	add("topo.diff.sard", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		switch kp {
		case 1:
			return wave20Fact(t, p, "m") - 1
		case 2:
			return n - wave20Fact(t, p, "m") + 1
		default:
			return n
		}
	})
	add("topo.diff.transversality", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		switch kp {
		case 1:
			dim := wave20Fact(t, p, "a") + wave20Fact(t, p, "b") - n
			if dim < 0 {
				return 0
			}
			return dim
		case 2:
			return n - wave20Fact(t, p, "a")
		default:
			return 2*n - wave20Fact(t, p, "a") - wave20Fact(t, p, "b")
		}
	})
	add("topo.diff.whitney", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		switch kp {
		case 1:
			return 2 * n
		case 2:
			return 2*n - 1
		default:
			return n
		}
	})
	add("topo.diff.tubular", func(t *testing.T, p generator.Problem, kp int) int {
		n, m := wave20Fact(t, p, "n"), wave20Fact(t, p, "m")
		if kp == 1 {
			return n - m
		}
		return n
	})
	add("topo.diff.degree", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20Fact(t, p, "k")
		case 2:
			if wave20Fact(t, p, "n")%2 == 0 {
				return -1
			}
			return 1
		default:
			return -wave20Fact(t, p, "k")
		}
	})
	add("topo.diff.flow", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		if kp == 1 {
			return n
		}
		k := wave20Fact(t, p, "k")
		if kp == 2 {
			return k
		}
		return n - k
	})
	add("topo.hom.chain_complex", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20Fact(t, p, "n")
		case 2:
			n := wave20Fact(t, p, "n")
			sum := 0
			for k := 0; k <= n; k++ {
				c := wave20Fact(t, p, "c"+strconv.Itoa(k))
				if k%2 == 0 {
					sum += c
				} else {
					sum -= c
				}
			}
			return sum
		default:
			return wave20Fact(t, p, "c") - wave20Fact(t, p, "r1") - wave20Fact(t, p, "r2")
		}
	})
	add("topo.hom.exact_sequence", func(t *testing.T, p generator.Problem, kp int) int {
		a, c := wave20Fact(t, p, "a"), wave20Fact(t, p, "c")
		switch kp {
		case 1:
			return a + c
		case 2:
			return c
		default:
			return a
		}
	})
	add("topo.hom.snake_lemma", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 2 {
			return wave20Fact(t, p, "n") - wave20Fact(t, p, "k")
		}
		k := wave20Fact(t, p, "k")
		if kp == 1 {
			return 3 * k
		}
		return k
	})
	add("topo.hom.five_lemma", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		if kp == 1 {
			return 2 * n
		}
		if kp == 2 {
			return n - 1
		}
		return 2 * (n - 1)
	})
	add("topo.hom.kunneth", func(t *testing.T, p generator.Problem, kp int) int {
		m, n := wave20Fact(t, p, "m"), wave20Fact(t, p, "n")
		switch kp {
		case 1:
			return m + n
		case 2:
			return wave20TPow2(m + n)
		default:
			return wave20TBinom(m+n, 2)
		}
	})
	add("topo.hom.universal_coefficient", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 1 {
			return wave20Fact(t, p, "r")
		}
		return wave20Fact(t, p, "k")
	})
	add("topo.hom.euler_poincare", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			if wave20Fact(t, p, "n")%2 == 0 {
				return 2
			}
			return 1
		case 2:
			return 0
		default:
			return 2 - 2*wave20Fact(t, p, "g")
		}
	})
	add("topo.hom.sphere", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		switch kp {
		case 1:
			return 2
		case 2:
			return 1
		default:
			if n%2 == 0 {
				return 2
			}
			return 1
		}
	})
	add("topo.hom.projective", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		switch kp {
		case 1:
			if n%2 == 1 {
				return 2
			}
			return 1
		case 2:
			if n%2 == 0 {
				return 1
			}
			return 0
		default:
			return n / 2
		}
	})
	add("topo.hom.torus", func(t *testing.T, p generator.Problem, kp int) int {
		n := wave20Fact(t, p, "n")
		switch kp {
		case 1:
			return n
		case 2:
			return wave20TPow2(n)
		default:
			return wave20TBinom(n, 2)
		}
	})
	add("topo.knot.crossing", func(t *testing.T, p generator.Problem, kp int) int {
		pv, qv := wave20Fact(t, p, "p"), wave20Fact(t, p, "q")
		if kp == 2 {
			return qv
		}
		if kp == 3 {
			return pv
		}
		return pv*qv - w20Max(pv, qv)
	})
	add("topo.knot.unknotting", func(t *testing.T, p generator.Problem, kp int) int {
		pv, qv := wave20Fact(t, p, "p"), wave20Fact(t, p, "q")
		if kp == 1 {
			return (pv - 1) * (qv - 1) / 2
		}
		if kp == 2 {
			return (qv - 1) / 2
		}
		return qv - 1
	})
	add("topo.knot.seifert_genus", func(t *testing.T, p generator.Problem, kp int) int {
		pv, qv := wave20Fact(t, p, "p"), wave20Fact(t, p, "q")
		switch kp {
		case 1:
			return (pv - 1) * (qv - 1) / 2
		case 2:
			return (qv - 1) / 2
		default:
			return (pv - 1) * (qv - 1)
		}
	})
	add("topo.knot.alexander", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return (wave20Fact(t, p, "p") - 1) * (wave20Fact(t, p, "q") - 1)
		case 2:
			return wave20Fact(t, p, "q") - 1
		default:
			return wave20Fact(t, p, "d1") + wave20Fact(t, p, "d2")
		}
	})
	add("topo.knot.skein", func(t *testing.T, p generator.Problem, kp int) int {
		c := wave20Fact(t, p, "c")
		switch kp {
		case 1:
			return wave20TPow2(c)
		case 2:
			return wave20TBinom(c, wave20Fact(t, p, "k"))
		default:
			return wave20TPow2(c) - 1
		}
	})
	add("topo.knot.torus_knot", func(t *testing.T, p generator.Problem, kp int) int {
		pv, qv := wave20Fact(t, p, "p"), wave20Fact(t, p, "q")
		switch kp {
		case 1:
			return w20Min(pv, qv)
		case 2:
			return pv*qv - w20Max(pv, qv)
		default:
			return (pv - 1) * (qv - 1) / 2
		}
	})
	add("topo.knot.connected_sum", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20Fact(t, p, "g1") + wave20Fact(t, p, "g2")
		case 2:
			return wave20Fact(t, p, "c1") + wave20Fact(t, p, "c2")
		default:
			return wave20Fact(t, p, "d1") + wave20Fact(t, p, "d2")
		}
	})
	add("topo.knot.bracket", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20TPow2(wave20Fact(t, p, "c"))
		case 2:
			return wave20Fact(t, p, "P") - wave20Fact(t, p, "N")
		default:
			return wave20Fact(t, p, "P") + wave20Fact(t, p, "N")
		}
	})
	add("topo.knot.writhe", func(t *testing.T, p generator.Problem, kp int) int {
		pv, nv := wave20Fact(t, p, "P"), wave20Fact(t, p, "N")
		switch kp {
		case 1:
			return pv - nv
		case 2:
			return pv + nv
		default:
			return pv + wave20Fact(t, p, "k") - nv
		}
	})
	add("topo.knot.linking", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 1 {
			return wave20Fact(t, p, "n")
		}
		if kp == 2 {
			return wave20Fact(t, p, "L")
		}
		return wave20Fact(t, p, "N") + 2*wave20Fact(t, p, "L")
	})
	add("topo.geo.cayley", func(t *testing.T, p generator.Problem, kp int) int {
		n, k := wave20Fact(t, p, "n"), wave20Fact(t, p, "k")
		switch kp {
		case 1:
			return 2 * k
		case 2:
			return n
		default:
			return n * k
		}
	})
	add("topo.geo.growth", func(t *testing.T, p generator.Problem, kp int) int {
		k := wave20Fact(t, p, "k")
		switch kp {
		case 1:
			return 2 * k
		case 2:
			return 2 * k * (2*k - 1)
		default:
			return 2*k - 1
		}
	})
	add("topo.geo.quasi_isometry", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 3 {
			return wave20Fact(t, p, "k")
		}
		return wave20Fact(t, p, "n")
	})
	add("topo.geo.ends", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			if wave20Fact(t, p, "n") == 1 {
				return 2
			}
			return 1
		case 2:
			return 2 * wave20Fact(t, p, "k")
		default:
			if wave20Fact(t, p, "m")+wave20Fact(t, p, "n") == 1 {
				return 2
			}
			return 1
		}
	})
	add("topo.geo.svarc_milnor", func(t *testing.T, p generator.Problem, kp int) int {
		m, n := wave20Fact(t, p, "m"), wave20Fact(t, p, "n")
		if kp == 1 || kp == 3 {
			return m
		}
		return m + n
	})
	add("topo.geo.hyperbolic", func(t *testing.T, p generator.Problem, _ int) int {
		return wave20Fact(t, p, "n") - 1
	})
	add("topo.geo.boundary", func(t *testing.T, p generator.Problem, kp int) int {
		switch kp {
		case 1:
			return wave20Fact(t, p, "n") - 1
		case 2:
			if wave20Fact(t, p, "n") == 1 {
				return 2
			}
			return 1
		default:
			if wave20Fact(t, p, "m")+wave20Fact(t, p, "n") == 1 {
				return 2
			}
			return 1
		}
	})
	add("topo.geo.dehn", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 1 {
			return 1
		}
		return 2
	})
	add("topo.geo.amenable", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 1 {
			return 1
		}
		return 0
	})
	add("topo.geo.property_t", func(t *testing.T, p generator.Problem, kp int) int {
		if kp == 1 {
			return 0
		}
		return 1
	})
	return cases
}

func TestWave20AnswersRecompute(t *testing.T) {
	reg := wave20Registry(t)
	for _, c := range wave20Cases() {
		gen, err := reg.Get(c.id)
		if err != nil {
			t.Fatalf("registry has no generator for %s: %v", c.id, err)
		}
		for draw := 0; draw < 8; draw++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: c.kp})
			want := c.want(t, p)
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", c.id, c.kp, p.Answer, want, p.Question, p.Facts)
			}
		}
	}
}
