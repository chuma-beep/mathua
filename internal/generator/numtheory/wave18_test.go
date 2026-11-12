package numtheory

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The Wave 18 families all answer with an integer, and every one publishes the
// parameters it drew. These tests re-derive each answer from those facts, so a
// generator that changes its arithmetic without changing its question fails
// here even though the question still reads plausibly.

func wave18IDs() []string {
	return []string{
		"nt.analytic.dirichlet_series",
		"nt.analytic.euler_product",
		"nt.analytic.continuation",
		"nt.analytic.functional_equation",
		"nt.analytic.zero_count",
		"nt.analytic.explicit_formula",
		"nt.analytic.mertens",
		"nt.analytic.abel",
		"nt.analytic.dirichlet_l",
		"nt.analytic.bias",
		"nt.app.rsa_keygen",
		"nt.app.elgamal",
		"nt.app.secret_sharing",
		"nt.app.checksum",
		"nt.app.zero_knowledge",
		"nt.app.ecc",
		"nt.app.primality_cert",
		"nt.app.reed_solomon",
		"nt.app.factorization",
		"nt.app.universal_hashing",
		"nt.dio.three_squares",
		"nt.dio.waring",
		"nt.dio.markov",
		"nt.dio.mordell",
		"nt.dio.catalan",
		"nt.dio.abc",
		"nt.dio.taxicab",
		"nt.dio.frobenius",
		"nt.dio.fermat_last",
		"nt.dio.four_cubes",
		"nt.comp.trial_division",
		"nt.comp.pollard_rho",
		"nt.comp.quadratic_sieve",
		"nt.comp.aks",
		"nt.comp.bsgs",
		"nt.comp.pohlig_hellman",
		"nt.comp.lll",
		"nt.comp.pratt",
		"nt.comp.cf_factor",
		"nt.comp.nfs",
		"nt.add.partitions",
		"nt.add.vinogradov",
		"nt.add.circle_method",
		"nt.add.schnirelmann",
		"nt.add.sumsets",
		"nt.add.cauchy_davenport",
		"nt.add.freiman",
		"nt.add.additive_bases",
		"nt.add.goldbach_count",
		"nt.add.additive_energy",
	}
}

func wave18Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave18Draw(t *testing.T, reg *generator.Registry, id string, ordinal int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{KPOrdinal: ordinal})
	if err != nil {
		t.Fatalf("%s ordinal %d: %v", id, ordinal, err)
	}
	return p
}

func wave18Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%s: missing fact %q\nfacts: %v", p.Question, key, p.Facts)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, raw)
	}
	return v
}

func wave18List(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%s: missing fact %q\nfacts: %v", p.Question, key, p.Facts)
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out []int
	for _, part := range strings.Split(raw, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("%s: list fact %q has entry %q", p.Question, key, part)
		}
		out = append(out, v)
	}
	return out
}

func TestWave18AnswersAreIntegers(t *testing.T) {
	reg := wave18Registry(t)
	for _, id := range wave18IDs() {
		if !reg.Has(id) {
			t.Fatalf("%s is not registered", id)
		}
		for ordinal := 1; ordinal <= 3; ordinal++ {
			for draw := 0; draw < 8; draw++ {
				p := wave18Draw(t, reg, id, ordinal)
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s ordinal %d: answer %q is not an integer\n%s", id, ordinal, p.Answer, p.Question)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s ordinal %d: facts answer %q disagrees with answer %q", id, ordinal, p.Facts["answer"], p.Answer)
				}
				if strings.TrimSpace(p.Explanation) == "" {
					t.Fatalf("%s ordinal %d: empty explanation\n%s", id, ordinal, p.Question)
				}
				if p.KPOrdinal != ordinal {
					t.Fatalf("%s: ordinal = %d, want %d", id, p.KPOrdinal, ordinal)
				}
			}
		}
	}
}

// The shipping contract asks for at least five distinct questions in twelve
// draws at each ordinal. A family that names no varying quantity fails here
// before it reaches the corpus.
func TestWave18QuestionVariety(t *testing.T) {
	reg := wave18Registry(t)
	for _, id := range wave18IDs() {
		for ordinal := 1; ordinal <= 3; ordinal++ {
			seen := map[string]bool{}
			for draw := 0; draw < 12; draw++ {
				seen[wave18Draw(t, reg, id, ordinal).Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s ordinal %d: %d distinct questions in 12 draws", id, ordinal, len(seen))
			}
		}
	}
}

type wave18Case struct {
	name string
	id   string
	kp   int
	want func(t *testing.T, p generator.Problem) int
}

func wave18Radical(t *testing.T, p generator.Problem) int {
	counts := map[int]bool{}
	for _, key := range []string{"a", "b", "c"} {
		for prime := range w18FactorCounts(wave18Fact(t, p, key)) {
			counts[prime] = true
		}
	}
	rad := 1
	for prime := range counts {
		rad *= prime
	}
	return rad
}

func wave18PartitionCount(n int) int {
	p := make([]int, n+1)
	p[0] = 1
	for part := 1; part <= n; part++ {
		for total := part; total <= n; total++ {
			p[total] += p[total-part]
		}
	}
	return p[n]
}

func wave18DistinctPartitions(n int) int {
	q := make([]int, n+1)
	q[0] = 1
	for part := 1; part <= n; part++ {
		for total := n; total >= part; total-- {
			q[total] += q[total-part]
		}
	}
	return q[n]
}

func wave18PrimeTriples(n int) int {
	primes := w18PrimesUpTo(n)
	count := 0
	for i := 0; i < len(primes); i++ {
		for j := i; j < len(primes); j++ {
			for k := j; k < len(primes); k++ {
				if primes[i]+primes[j]+primes[k] == n {
					count++
				}
			}
		}
	}
	return count
}

func wave18Energy(t *testing.T, p generator.Problem) int {
	m := wave18Fact(t, p, "m")
	set := wave18List(t, p, "a")
	counts := map[int]int{}
	for _, x := range set {
		for _, y := range set {
			counts[(x+y)%m]++
		}
	}
	energy := 0
	for _, c := range counts {
		energy += c * c
	}
	return energy
}

func wave18SetSumSize(t *testing.T, p generator.Problem) int {
	m := wave18Fact(t, p, "m")
	setA := wave18List(t, p, "a")
	setB := wave18List(t, p, "b")
	seen := map[int]bool{}
	for _, x := range setA {
		for _, y := range setB {
			seen[(x+y)%m] = true
		}
	}
	return len(seen)
}

func TestWave18AnswersRecompute(t *testing.T) {
	reg := wave18Registry(t)
	cases := []wave18Case{
		{name: "rsa modulus", id: "nt.app.rsa_keygen", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return wave18Fact(t, p, "p") * wave18Fact(t, p, "q")
		}},
		{name: "rsa totient", id: "nt.app.rsa_keygen", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return (wave18Fact(t, p, "p") - 1) * (wave18Fact(t, p, "q") - 1)
		}},
		{name: "rsa private exponent", id: "nt.app.rsa_keygen", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			phi := (wave18Fact(t, p, "p") - 1) * (wave18Fact(t, p, "q") - 1)
			return w18ModInverse(wave18Fact(t, p, "e"), phi)
		}},
		{name: "elgamal public value", id: "nt.app.elgamal", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return w18PowMod(wave18Fact(t, p, "g"), wave18Fact(t, p, "x"), wave18Fact(t, p, "p"))
		}},
		{name: "shamir interpolation", id: "nt.app.secret_sharing", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			mod := wave18Fact(t, p, "p")
			return ((2*wave18Fact(t, p, "f1")-wave18Fact(t, p, "f2"))%mod + mod) % mod
		}},
		{name: "three squares residue", id: "nt.dio.three_squares", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return wave18Fact(t, p, "n") % 8
		}},
		{name: "waring function", id: "nt.dio.waring", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			k := wave18Fact(t, p, "k")
			return w18IPow(2, k) + w18IPow(3, k)/w18IPow(2, k) - 2
		}},
		{name: "markov square sum", id: "nt.dio.markov", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			x, y, z := wave18Fact(t, p, "x"), wave18Fact(t, p, "y"), wave18Fact(t, p, "z")
			return x*x + y*y + z*z
		}},
		{name: "markov vieta jump", id: "nt.dio.markov", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			return 3*wave18Fact(t, p, "x")*wave18Fact(t, p, "y") - wave18Fact(t, p, "z")
		}},
		{name: "mordell right hand side", id: "nt.dio.mordell", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			x, k := wave18Fact(t, p, "x"), wave18Fact(t, p, "k")
			return x*x*x + k
		}},
		{name: "catalan difference", id: "nt.dio.catalan", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			return w18IPow(wave18Fact(t, p, "a"), 2) - w18IPow(wave18Fact(t, p, "b"), 3)
		}},
		{name: "abc radical", id: "nt.dio.abc", kp: 1, want: wave18Radical},
		{name: "abc c minus radical", id: "nt.dio.abc", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			return wave18Fact(t, p, "c") - wave18Radical(t, p)
		}},
		{name: "taxicab cube sum", id: "nt.dio.taxicab", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			a, b := wave18Fact(t, p, "a"), wave18Fact(t, p, "b")
			return a*a*a + b*b*b
		}},
		{name: "frobenius number", id: "nt.dio.frobenius", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			a, b := wave18Fact(t, p, "a"), wave18Fact(t, p, "b")
			return a*b - a - b
		}},
		{name: "frobenius gap count", id: "nt.dio.frobenius", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			a, b := wave18Fact(t, p, "a"), wave18Fact(t, p, "b")
			return (a - 1) * (b - 1) / 2
		}},
		{name: "frobenius conductor", id: "nt.dio.frobenius", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			a, b := wave18Fact(t, p, "a"), wave18Fact(t, p, "b")
			return a*b - a - b + 1
		}},
		{name: "four cubes residue", id: "nt.dio.four_cubes", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return wave18Fact(t, p, "n") % 9
		}},
		{name: "trial division smallest", id: "nt.comp.trial_division", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return w18SmallestFactor(wave18Fact(t, p, "n"))
		}},
		{name: "trial division largest", id: "nt.comp.trial_division", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			return w18LargestFactor(wave18Fact(t, p, "n"))
		}},
		{name: "pollard rho step", id: "nt.comp.pollard_rho", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			x, c, n := wave18Fact(t, p, "x0"), wave18Fact(t, p, "c"), wave18Fact(t, p, "n")
			return (x*x + c) % n
		}},
		{name: "quadratic sieve prime count", id: "nt.comp.quadratic_sieve", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return len(w18PrimesUpTo(wave18Fact(t, p, "b")))
		}},
		{name: "aks totient", id: "nt.comp.aks", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return w18Totient(wave18Fact(t, p, "r"))
		}},
		{name: "bsgs table size", id: "nt.comp.bsgs", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return w18CeilSqrt(wave18Fact(t, p, "n"))
		}},
		{name: "pohlig hellman omega", id: "nt.comp.pohlig_hellman", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return w18Omega(wave18Fact(t, p, "n"))
		}},
		{name: "partition count", id: "nt.add.partitions", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return wave18PartitionCount(wave18Fact(t, p, "n"))
		}},
		{name: "distinct partition count", id: "nt.add.partitions", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return wave18DistinctPartitions(wave18Fact(t, p, "n"))
		}},
		{name: "three prime representations", id: "nt.add.vinogradov", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return wave18PrimeTriples(wave18Fact(t, p, "n"))
		}},
		{name: "circle method totient", id: "nt.add.circle_method", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return w18Totient(wave18Fact(t, p, "q"))
		}},
		{name: "farey length", id: "nt.add.circle_method", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			q := wave18Fact(t, p, "q")
			total := 1
			for i := 1; i <= q; i++ {
				total += w18Totient(i)
			}
			return total
		}},
		{name: "integer sumset", id: "nt.add.sumsets", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return wave18Fact(t, p, "a") + wave18Fact(t, p, "b") - 1
		}},
		{name: "doubled interval", id: "nt.add.sumsets", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return 2*wave18Fact(t, p, "a") - 1
		}},
		{name: "modular sumset", id: "nt.add.sumsets", kp: 3, want: wave18SetSumSize},
		{name: "cauchy davenport cap", id: "nt.add.cauchy_davenport", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			prime := wave18Fact(t, p, "p")
			ans := wave18Fact(t, p, "a") + wave18Fact(t, p, "b") - 1
			if ans > prime {
				ans = prime
			}
			return ans
		}},
		{name: "cauchy davenport raw", id: "nt.add.cauchy_davenport", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return wave18Fact(t, p, "a") + wave18Fact(t, p, "b") - 1
		}},
		{name: "freiman bound", id: "nt.add.freiman", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			return 3*wave18Fact(t, p, "k") - 4
		}},
		{name: "freiman interval", id: "nt.add.freiman", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			return 2*wave18Fact(t, p, "k") - 3
		}},
		{name: "freiman sumset floor", id: "nt.add.freiman", kp: 3, want: func(t *testing.T, p generator.Problem) int {
			return 2*wave18Fact(t, p, "k") - 1
		}},
		{name: "goldbach count", id: "nt.add.goldbach_count", kp: 1, want: func(t *testing.T, p generator.Problem) int {
			n := wave18Fact(t, p, "n")
			count := 0
			for _, prime := range w18PrimesUpTo(n) {
				if prime > n-prime {
					break
				}
				if w18IsPrime(n - prime) {
					count++
				}
			}
			return count
		}},
		{name: "additive energy", id: "nt.add.additive_energy", kp: 1, want: wave18Energy},
		{name: "additive energy cube", id: "nt.add.additive_energy", kp: 2, want: func(t *testing.T, p generator.Problem) int {
			size := len(wave18List(t, p, "a"))
			return size * size * size
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for draw := 0; draw < 8; draw++ {
				p := wave18Draw(t, reg, tc.id, tc.kp)
				want := tc.want(t, p)
				got, err := strconv.Atoi(p.Answer)
				if err != nil {
					t.Fatalf("answer %q is not an integer", p.Answer)
				}
				if got != want {
					t.Fatalf("answer = %d, recomputed %d\n%s\nfacts: %v", got, want, p.Question, p.Facts)
				}
			}
		})
	}
}
