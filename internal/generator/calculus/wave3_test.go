package calculus

import (
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Independent arithmetic used by the checks below. The generator has its own
// helpers; this file rebuilds the values through math/big and through direct
// series evaluation so a shared mistake cannot make both sides agree.
// ---------------------------------------------------------------------------

func w3TestRat(num, den int64) *big.Rat {
	return new(big.Rat).SetFrac(big.NewInt(num), big.NewInt(den))
}

func w3TestRatString(r *big.Rat) string {
	if r.IsInt() {
		return r.Num().String()
	}
	return r.Num().String() + "/" + r.Denom().String()
}

func w3TestPow(base int64, exp int) int64 {
	res := int64(1)
	for i := 0; i < exp; i++ {
		res *= base
	}
	return res
}

func w3TestFactorial(n int) int64 {
	res := int64(1)
	for i := 2; i <= n; i++ {
		res *= int64(i)
	}
	return res
}

func w3TestBinom(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	res := int64(1)
	for i := 1; i <= k; i++ {
		res = res * int64(n-k+i) / int64(i)
	}
	return res
}

func w3TestFloatFactorial(n int) float64 {
	res := 1.0
	for i := 2; i <= n; i++ {
		res *= float64(i)
	}
	return res
}

func w3TestFresnel(x float64, sine bool) float64 {
	sum := 0.0
	sign := 1.0
	for n := 0; n < 30; n++ {
		var term float64
		if sine {
			term = sign * math.Pow(math.Pi/2, float64(2*n+1)) * math.Pow(x, float64(4*n+3)) /
				(w3TestFloatFactorial(2*n+1) * float64(4*n+3))
		} else {
			term = sign * math.Pow(math.Pi/2, float64(2*n)) * math.Pow(x, float64(4*n+1)) /
				(w3TestFloatFactorial(2*n) * float64(4*n+1))
		}
		sum += term
		sign = -sign
	}
	return sum
}

func w3TestEllipticCoefficient(n int) float64 {
	num := w3TestFloatFactorial(2 * n)
	den := math.Pow(2, float64(2*n)) * math.Pow(w3TestFloatFactorial(n), 2)
	a := num / den
	return a * a
}

func w3TestEllipticK(k float64) float64 {
	sum := 0.0
	for n := 0; n < 30; n++ {
		sum += w3TestEllipticCoefficient(n) * math.Pow(k, float64(2*n))
	}
	return math.Pi / 2 * sum
}

func w3TestEllipticE(k float64) float64 {
	sum := 1.0
	for n := 1; n < 30; n++ {
		sum -= w3TestEllipticCoefficient(n) * math.Pow(k, float64(2*n)) / float64(2*n-1)
	}
	return math.Pi / 2 * sum
}

func w3TestAGM(a, b float64) float64 {
	for i := 0; i < 100; i++ {
		nextA := (a + b) / 2
		nextB := math.Sqrt(a * b)
		if math.Abs(nextA-nextB) < 1e-15 {
			return (nextA + nextB) / 2
		}
		a, b = nextA, nextB
	}
	return (a + b) / 2
}

func w3TestFixed(x float64, places int) string {
	scale := math.Pow(10, float64(places))
	return strconv.FormatFloat(math.Round(x*scale)/scale, 'f', places, 64)
}

func w3TestRatFromInts(num, den int) string {
	return w3TestRatString(w3TestRat(int64(num), int64(den)))
}

func w3TestFloatFact(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%q: fact %q = %q is not a number", p.Question, key, p.Facts[key])
	}
	return v
}

func w3TestWant(t *testing.T, p generator.Problem, want string) {
	t.Helper()
	if p.Answer != want {
		t.Fatalf("answer = %q, want %q", p.Answer, want)
	}
}

// w3TestRaabeLimit evaluates n(1 - a_{n+1}/a_n) at a large n. The terms are
// simple enough for a float evaluation to settle on the integer limit.
func w3TestRaabeLimit(a func(int) float64) int {
	const n = 10000000
	return int(math.Round(float64(n) * (1 - a(n+1)/a(n))))
}

func TestWave3FactsDriveAnswers(t *testing.T) {
	cases := []struct {
		name  string
		gen   generator.Generator
		mc    bool
		check func(t *testing.T, p generator.Problem, kp int)
	}{
		{
			name: "cauchy condensation",
			gen:  &w3CauchyCondensationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				pp := calcFactInt(t, p, "p")
				switch kp {
				case 1:
					w3TestWant(t, p, strconv.Itoa(1-pp))
				case 2:
					w3TestWant(t, p, strconv.Itoa(int(w3TestPow(2, pp-1))))
				default:
					num := w3TestPow(2, pp-1)
					w3TestWant(t, p, w3TestRatString(w3TestRat(num, num-1)))
				}
			},
		},
		{
			name: "raabe",
			gen:  &w3RaabeGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					pp := calcFactInt(t, p, "p")
					got := w3TestRaabeLimit(func(n int) float64 {
						return 1 / math.Pow(float64(n), float64(pp))
					})
					w3TestWant(t, p, strconv.Itoa(got))
				case 2:
					pp := calcFactInt(t, p, "p")
					q := calcFactInt(t, p, "q")
					got := w3TestRaabeLimit(func(n int) float64 {
						return math.Pow(float64(n+1), float64(q)) / math.Pow(float64(n), float64(pp))
					})
					w3TestWant(t, p, strconv.Itoa(got))
				default:
					pp := calcFactInt(t, p, "p")
					got := w3TestRaabeLimit(func(n int) float64 {
						return 1 / (math.Pow(float64(n), float64(pp)) * float64(n+1))
					})
					w3TestWant(t, p, strconv.Itoa(got))
				}
			},
		},
		{
			name: "bertrand",
			gen:  &w3BertrandGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				const n = 10000000
				pp := calcFactInt(t, p, "p")
				var a func(int) float64
				switch kp {
				case 1:
					a = func(n int) float64 {
						return 1 / (float64(n) * math.Pow(math.Log(float64(n)), float64(pp)))
					}
				case 2:
					a = func(n int) float64 {
						return math.Pow(math.Log(float64(n)), float64(pp)) / float64(n)
					}
				default:
					c := calcFactInt(t, p, "c")
					a = func(n int) float64 {
						return float64(c) / (float64(n) * math.Pow(math.Log(float64(n)), float64(pp)))
					}
				}
				inner := float64(n)*(1-a(n+1)/a(n)) - 1
				got := int(math.Round(math.Log(float64(n)) * inner))
				w3TestWant(t, p, strconv.Itoa(got))
			},
		},
		{
			name: "kummer",
			gen:  &w3KummerGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				const n = 10000000
				pp := calcFactInt(t, p, "p")
				c := calcFactInt(t, p, "c")
				var a, b func(int) float64
				switch kp {
				case 1:
					a = func(n int) float64 { return 1 / math.Pow(float64(n), float64(pp)) }
					b = func(n int) float64 { return float64(c * n) }
				case 2:
					a = func(n int) float64 {
						return 1 / (float64(n) * math.Pow(math.Log(float64(n)), float64(pp)))
					}
					b = func(n int) float64 { return float64(c) * float64(n) * math.Log(float64(n)) }
				default:
					a = func(n int) float64 {
						return 1 / (math.Pow(float64(n), float64(pp)) * float64(n+1))
					}
					b = func(n int) float64 { return float64(c * n) }
				}
				got := int(math.Round(b(n)*a(n)/a(n+1) - b(n+1)))
				w3TestWant(t, p, strconv.Itoa(got))
			},
		},
		{
			name: "gauss",
			gen:  &w3GaussGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				const n = 10000000
				switch kp {
				case 1:
					pp := calcFactInt(t, p, "p")
					ratio := math.Pow(float64(n+1)/float64(n), float64(pp))
					w3TestWant(t, p, strconv.Itoa(int(math.Round(float64(n)*(ratio-1)))))
				case 2:
					alpha := calcFactInt(t, p, "alpha")
					beta := calcFactInt(t, p, "beta")
					ratio := float64(n+beta) / float64(n+alpha)
					w3TestWant(t, p, strconv.Itoa(int(math.Round(float64(n)*(ratio-1)))))
				default:
					d := calcFactInt(t, p, "d")
					alpha := calcFactInt(t, p, "alpha")
					ratio := float64(n+alpha) / float64(n+1+alpha) *
						math.Pow(float64(n+2)/float64(n+1), float64(d))
					w3TestWant(t, p, strconv.Itoa(int(math.Round(float64(n)*(ratio-1)))))
				}
			},
		},
		{
			name: "abel",
			gen:  &w3AbelGen{},
			mc:   true,
			check: func(t *testing.T, p generator.Problem, kp int) {
				if kp <= 2 {
					if calcFactInt(t, p, "a_converges") != 1 || calcFactInt(t, p, "b_bounded") != 1 || calcFactInt(t, p, "b_monotone") != 1 {
						t.Fatalf("facts do not satisfy Abel's hypothesis: %v", p.Facts)
					}
					w3TestWant(t, p, "converges")
					return
				}
				want := "no"
				if calcFactInt(t, p, "bounded") == 1 && calcFactInt(t, p, "monotone") == 1 {
					want = "yes"
				}
				w3TestWant(t, p, want)
			},
		},
		{
			name: "uniform convergence",
			gen:  &w3UniformConvergenceGen{},
			mc:   true,
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					want := "no"
					if p.Facts["sup_limit"] == "0" {
						want = "yes"
					}
					w3TestWant(t, p, want)
				case 2:
					uniform := map[string]bool{
						"power_half": true, "power_third": true, "power_unit": false, "p2": true,
						"harmonic_open": false, "alt_harmonic": true, "p3": true, "p4": true,
					}
					want := "no"
					if uniform[p.Facts["kind"]] {
						want = "yes"
					}
					w3TestWant(t, p, want)
				default:
					want := "no"
					if p.Facts["m_converges"] == "1" {
						want = "yes"
					}
					w3TestWant(t, p, want)
				}
			},
		},
		{
			name: "euler transform",
			gen:  &w3EulerTransformGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					rn := calcFactInt(t, p, "r_num")
					rd := calcFactInt(t, p, "r_den")
					k := calcFactInt(t, p, "k")
					num := w3TestPow(int64(rd-rn), k)
					den := w3TestPow(int64(rd), k) * w3TestPow(2, k+1)
					w3TestWant(t, p, w3TestRatString(w3TestRat(num, den)))
				case 2:
					rn := calcFactInt(t, p, "r_num")
					rd := calcFactInt(t, p, "r_den")
					w3TestWant(t, p, w3TestRatString(w3TestRat(int64(rd), int64(rd+rn))))
				default:
					k := calcFactInt(t, p, "k")
					w3TestWant(t, p, w3TestRatString(w3TestRat(1, int64(k+1)*w3TestPow(2, k+1))))
				}
			},
		},
		{
			name: "power integration",
			gen:  &w3PowerIntegrationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					n := calcFactInt(t, p, "n")
					c := calcFactInt(t, p, "c")
					w3TestWant(t, p, strconv.Itoa(c/(n+1)))
				case 2:
					m := calcFactInt(t, p, "m")
					w3TestWant(t, p, w3TestRatFromInts(1, m))
				default:
					k := calcFactInt(t, p, "k")
					w3TestWant(t, p, w3TestRatFromInts(1, 2*k+1))
				}
			},
		},
		{
			name: "power differentiation",
			gen:  &w3PowerDifferentiationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					n := calcFactInt(t, p, "n")
					c := calcFactInt(t, p, "c")
					w3TestWant(t, p, strconv.Itoa(n*c))
				case 2:
					k := calcFactInt(t, p, "k")
					w3TestWant(t, p, strconv.Itoa(k+1))
				default:
					k := calcFactInt(t, p, "k")
					w3TestWant(t, p, strconv.Itoa(2*k))
				}
			},
		},
		{
			name: "lagrange inversion",
			gen:  &w3LagrangeInversionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					n := calcFactInt(t, p, "n")
					pp := calcFactInt(t, p, "p")
					w3TestWant(t, p, strconv.FormatInt(w3TestBinom(pp*n, n-1)/int64(n), 10))
				case 2:
					n := calcFactInt(t, p, "n")
					w3TestWant(t, p, w3TestRatString(w3TestRat(w3TestPow(int64(n), n-1), w3TestFactorial(n))))
				default:
					n := calcFactInt(t, p, "n")
					w3TestWant(t, p, strconv.FormatInt(w3TestBinom(2*n, n)/int64(n+1), 10))
				}
			},
		},
		{
			name: "romberg",
			gen:  &w3RombergGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					t1 := calcFactInt(t, p, "t_h")
					t2 := calcFactInt(t, p, "t_half")
					w3TestWant(t, p, w3TestRatString(w3TestRat(int64(4*t2-t1), 3)))
				case 2:
					r1 := calcFactInt(t, p, "r_prev")
					r2 := calcFactInt(t, p, "r_cur")
					w3TestWant(t, p, w3TestRatString(w3TestRat(int64(16*r2-r1), 15)))
				default:
					t0 := calcFactInt(t, p, "t_h")
					t1 := calcFactInt(t, p, "t_half")
					t2 := calcFactInt(t, p, "t_quarter")
					r11 := w3TestRat(int64(4*t1-t0), 3)
					r21 := w3TestRat(int64(4*t2-t1), 3)
					diff := new(big.Rat).Sub(r21, r11)
					r22 := new(big.Rat).Add(r21, new(big.Rat).Quo(diff, w3TestRat(15, 1)))
					w3TestWant(t, p, w3TestRatString(r22))
				}
			},
		},
		{
			name: "monte carlo",
			gen:  &w3MonteCarloGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					a := calcFactInt(t, p, "a")
					b := calcFactInt(t, p, "b")
					n := calcFactInt(t, p, "n")
					sum := 0
					for _, part := range strings.Split(p.Facts["samples"], ",") {
						v, err := strconv.Atoi(strings.TrimSpace(part))
						if err != nil {
							t.Fatalf("sample %q is not an integer", part)
						}
						sum += v
					}
					w3TestWant(t, p, w3TestRatFromInts((b-a)*sum, n))
				case 2:
					m := calcFactInt(t, p, "m")
					n := calcFactInt(t, p, "n")
					w3TestWant(t, p, w3TestFixed(4*float64(m)/float64(n), 3))
				default:
					s := calcFactInt(t, p, "s")
					n := calcFactInt(t, p, "n")
					w3TestWant(t, p, w3TestFixed(float64(s)/math.Sqrt(float64(n)), 3))
				}
			},
		},
		{
			name: "gauss legendre",
			gen:  &w3GaussLegendreGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					a := calcFactInt(t, p, "a")
					w3TestWant(t, p, w3TestRatFromInts(2*a*a*a, 3))
				case 2:
					a := calcFactInt(t, p, "a")
					w3TestWant(t, p, w3TestRatFromInts(2*a*a*a*a*a, 5))
				default:
					a := calcFactInt(t, p, "a")
					b := calcFactInt(t, p, "b")
					w3TestWant(t, p, w3TestRatFromInts((b-a)*(a+b), 2))
				}
			},
		},
		{
			name: "fresnel",
			gen:  &w3FresnelGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					x := w3TestFloatFact(t, p, "x")
					w3TestWant(t, p, w3TestFixed(w3TestFresnel(x, false), 4))
				case 2:
					x := w3TestFloatFact(t, p, "x")
					w3TestWant(t, p, w3TestFixed(w3TestFresnel(x, true), 4))
				default:
					n := calcFactInt(t, p, "n")
					v := math.Pow(-1, float64(n)) * math.Pow(math.Pi/2, float64(2*n)) /
						(w3TestFloatFactorial(2*n) * float64(4*n+1))
					w3TestWant(t, p, w3TestFixed(v, 6))
				}
			},
		},
		{
			name: "elliptic",
			gen:  &w3EllipticGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				k := w3TestFloatFact(t, p, "k")
				switch kp {
				case 1:
					w3TestWant(t, p, w3TestFixed(w3TestEllipticK(k), 4))
				case 2:
					w3TestWant(t, p, w3TestFixed(w3TestEllipticE(k), 4))
				default:
					agm := w3TestAGM(1, math.Sqrt(1-k*k))
					w3TestWant(t, p, w3TestFixed(math.Pi/(2*agm), 4))
				}
			},
		},
		{
			name: "wallis",
			gen:  &w3WallisGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					m := calcFactInt(t, p, "m")
					value := w3TestRat(1, 1)
					for j := 1; j <= m; j++ {
						value.Mul(value, w3TestRat(int64(2*j), int64(2*j+1)))
					}
					w3TestWant(t, p, w3TestRatString(value))
				case 2:
					m := calcFactInt(t, p, "m")
					value := w3TestRat(1, 2)
					for j := 1; j <= m; j++ {
						value.Mul(value, w3TestRat(int64(2*j-1), int64(2*j)))
					}
					w3TestWant(t, p, w3TestRatString(value))
				default:
					n := calcFactInt(t, p, "n")
					w3TestWant(t, p, w3TestRatFromInts(n-1, n))
				}
			},
		},
		{
			name: "centroid of a curve",
			gen:  &w3CentroidCurveGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				switch kp {
				case 1:
					x1 := calcFactInt(t, p, "x1")
					x2 := calcFactInt(t, p, "x2")
					w3TestWant(t, p, w3TestRatFromInts(x1+x2, 2))
				case 2:
					y1 := calcFactInt(t, p, "y1")
					y2 := calcFactInt(t, p, "y2")
					w3TestWant(t, p, w3TestRatFromInts(y1+y2, 2))
				default:
					dx := calcFactInt(t, p, "dx")
					dy := calcFactInt(t, p, "dy")
					w3TestWant(t, p, strconv.Itoa(int(math.Round(math.Hypot(float64(dx), float64(dy))))))
				}
			},
		},
	}

	if len(cases) != 18 {
		t.Fatalf("the table covers %d concepts, want 18", len(cases))
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for kp := 1; kp <= 3; kp++ {
				for seed := int64(0); seed < 300; seed++ {
					rand.Seed(seed)
					p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp, Seed: seed})
					if p.KPOrdinal != kp {
						t.Fatalf("kp%d seed%d: ordinal = %d", kp, seed, p.KPOrdinal)
					}
					if p.Question == "" || p.Answer == "" {
						t.Fatalf("kp%d seed%d: blank question or answer", kp, seed)
					}
					if p.Explanation == "" {
						t.Fatalf("kp%d seed%d: empty explanation", kp, seed)
					}
					if p.Facts == nil {
						t.Fatalf("kp%d seed%d: empty facts", kp, seed)
					}
					if !tc.mc && !strings.Contains(p.Explanation, p.Answer) {
						t.Fatalf("kp%d seed%d: explanation does not state %q: %s", kp, seed, p.Answer, p.Explanation)
					}
					tc.check(t, p, kp)
				}
			}
		})
	}
}

// TestWave3Diversity checks the same floor the shipping contract applies: each
// ordinal supplies at least five different questions in twelve draws.
func TestWave3Diversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"cauchy condensation", &w3CauchyCondensationGen{}},
		{"raabe", &w3RaabeGen{}},
		{"bertrand", &w3BertrandGen{}},
		{"kummer", &w3KummerGen{}},
		{"gauss", &w3GaussGen{}},
		{"abel", &w3AbelGen{}},
		{"uniform convergence", &w3UniformConvergenceGen{}},
		{"euler transform", &w3EulerTransformGen{}},
		{"power integration", &w3PowerIntegrationGen{}},
		{"power differentiation", &w3PowerDifferentiationGen{}},
		{"lagrange inversion", &w3LagrangeInversionGen{}},
		{"romberg", &w3RombergGen{}},
		{"monte carlo", &w3MonteCarloGen{}},
		{"gauss legendre", &w3GaussLegendreGen{}},
		{"fresnel", &w3FresnelGen{}},
		{"elliptic", &w3EllipticGen{}},
		{"wallis", &w3WallisGen{}},
		{"centroid of a curve", &w3CentroidCurveGen{}},
	}
	for _, tc := range gens {
		t.Run(tc.name, func(t *testing.T) {
			rand.Seed(4242)
			seeded := make([]string, 3)
			for kp := 1; kp <= 3; kp++ {
				seeded[kp-1] = tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, Seed: 4242, KPOrdinal: kp}).Question
			}
			if seeded[0] == seeded[1] || seeded[1] == seeded[2] || seeded[0] == seeded[2] {
				t.Fatalf("the same seeded question appears at two ordinals")
			}
			for kp := 1; kp <= 3; kp++ {
				best := 0
				for attempt := 0; attempt < 4 && best < 5; attempt++ {
					seen := map[string]bool{}
					for i := 0; i < 12; i++ {
						p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
						seen[p.Question] = true
					}
					if len(seen) > best {
						best = len(seen)
					}
				}
				if best < 5 {
					t.Fatalf("kp%d: only %d distinct questions in 12 draws", kp, best)
				}
			}
		})
	}
}
