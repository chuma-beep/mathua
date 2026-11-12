package complex

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

var wave19IDs = []string{
	"complex.gauss.norm", "complex.gauss.units", "complex.gauss.division",
	"complex.gauss.gcd", "complex.gauss.primes", "complex.gauss.two_squares",
	"complex.gauss.associates", "complex.gauss.pythagorean", "complex.gauss.theta",
	"complex.gauss.eisenstein",
	"complex.geom.rotation", "complex.geom.area", "complex.geom.circle",
	"complex.geom.reflection", "complex.geom.spiral", "complex.geom.ptolemy",
	"complex.geom.triangle", "complex.geom.centroid", "complex.geom.perpendicular",
	"complex.geom.lattice",
	"complex.poly.fta", "complex.poly.vieta_sum", "complex.poly.vieta_product",
	"complex.poly.unity_sum", "complex.poly.cyclotomic", "complex.poly.chebyshev",
	"complex.poly.discriminant", "complex.poly.conjugate_pairs",
	"complex.poly.multiplicity", "complex.poly.resultant",
	"complex.meth.simple_pole", "complex.meth.higher_pole",
	"complex.meth.partial_fractions", "complex.meth.residue_sum",
	"complex.meth.winding", "complex.meth.pole_count", "complex.meth.zero_count",
	"complex.meth.order", "complex.meth.jordan", "complex.meth.branch_cut",
	"complex.spec.gamma_half", "complex.spec.beta", "complex.spec.zeta_even",
	"complex.spec.theta", "complex.spec.bessel", "complex.spec.airy",
	"complex.spec.hypergeometric", "complex.spec.polygamma", "complex.spec.erf",
	"complex.spec.polylog",
}

func wave19Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave19FactInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func wave19Fact64(t *testing.T, p generator.Problem, key string) int64 {
	t.Helper()
	v, err := strconv.ParseInt(p.Facts[key], 10, 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an int64\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func wave19Answer(t *testing.T, p generator.Problem) int {
	t.Helper()
	v, err := strconv.Atoi(p.Answer)
	if err != nil {
		t.Fatalf("%s: answer %q is not an integer", p.Question, p.Answer)
	}
	return v
}

func wave19Answer64(t *testing.T, p generator.Problem) int64 {
	t.Helper()
	v, err := strconv.ParseInt(p.Answer, 10, 64)
	if err != nil {
		t.Fatalf("%s: answer %q is not an int64", p.Question, p.Answer)
	}
	return v
}

func wave19Check(t *testing.T, p generator.Problem, want int) {
	t.Helper()
	if got := wave19Answer(t, p); got != want {
		t.Fatalf("%s: answer = %d, want %d\nfacts: %v", p.Question, got, want, p.Facts)
	}
	if got := wave19FactInt(t, p, "answer"); got != want {
		t.Fatalf("%s: facts[answer] = %d, want %d", p.Question, got, want)
	}
	if p.Explanation == "" {
		t.Fatalf("%s: empty explanation", p.Question)
	}
	if p.KPOrdinal < 1 || p.KPOrdinal > 3 {
		t.Fatalf("%s: ordinal %d out of range", p.Question, p.KPOrdinal)
	}
}

func wave19Check64(t *testing.T, p generator.Problem, want int64) {
	t.Helper()
	if got := wave19Answer64(t, p); got != want {
		t.Fatalf("%s: answer = %d, want %d\nfacts: %v", p.Question, got, want, p.Facts)
	}
	if got := wave19Fact64(t, p, "answer"); got != want {
		t.Fatalf("%s: facts[answer] = %d, want %d", p.Question, got, want)
	}
	if p.Explanation == "" {
		t.Fatalf("%s: empty explanation", p.Question)
	}
}

func wave19Factorial(t *testing.T, n int) int64 {
	t.Helper()
	out := int64(1)
	for i := 2; i <= n; i++ {
		out *= int64(i)
	}
	return out
}

func wave19Totient(n int) int {
	count := 0
	for k := 1; k <= n; k++ {
		if wave19GCD(k, n) == 1 {
			count++
		}
	}
	return count
}

func wave19GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func wave19DivisorCounts(n int) (int, int) {
	var one, three int
	for d := 1; d <= n; d++ {
		if n%d != 0 {
			continue
		}
		switch d % 4 {
		case 1:
			one++
		case 3:
			three++
		}
	}
	return one, three
}

func wave19MinSquares(n int) int {
	for a := 0; a*a <= n; a++ {
		if a*a == n {
			return 1
		}
	}
	for a := 0; a*a <= n; a++ {
		for b := 0; a*a+b*b <= n; b++ {
			if a*a+b*b == n {
				return 2
			}
		}
	}
	for a := 0; a*a <= n; a++ {
		for b := 0; a*a+b*b <= n; b++ {
			for c := 0; a*a+b*b+c*c <= n; c++ {
				if a*a+b*b+c*c == n {
					return 3
				}
			}
		}
	}
	return 4
}

func TestWave19AnswersAreIntegers(t *testing.T) {
	reg := wave19Registry(t)
	for _, id := range wave19IDs {
		g, err := reg.Get(id)
		if err != nil {
			t.Fatalf("%s: no generator: %v", id, err)
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 8; i++ {
				p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
				if p.KPOrdinal != kp {
					t.Fatalf("%s: ordinal = %d, want %d", id, p.KPOrdinal, kp)
				}
				if _, err := strconv.Atoi(p.Answer); err != nil {
					t.Fatalf("%s: answer %q is not an integer", id, p.Answer)
				}
				if p.Facts == nil || p.Facts["answer"] != p.Answer {
					t.Fatalf("%s: facts[answer] = %q, answer = %q", id, p.Facts["answer"], p.Answer)
				}
				if p.Explanation == "" {
					t.Fatalf("%s: empty explanation", id)
				}
			}
		}
	}
}

func TestWave19QuestionVariety(t *testing.T) {
	reg := wave19Registry(t)
	for _, id := range wave19IDs {
		g, _ := reg.Get(id)
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 12; i++ {
				p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", id, kp, len(seen))
			}
		}
	}
}

type wave19Case struct {
	id string
	kp int
}

func TestWave19AnswersRecompute(t *testing.T) {
	reg := wave19Registry(t)
	for _, id := range wave19IDs {
		g, _ := reg.Get(id)
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 6; draw++ {
				p := g.Generate(generator.GeneratorContext{KPOrdinal: kp})
				switch id {
				case "complex.gauss.norm":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					switch kp {
					case 1:
						wave19Check(t, p, a*a+b*b)
					case 2:
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						wave19Check(t, p, (a*a+b*b)*(c*c+d*d))
					default:
						wave19Check(t, p, 2*(a*a+b*b))
					}
				case "complex.gauss.units", "complex.geom.rotation":
					a, b, k := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b"), wave19FactInt(t, p, "k")
					re := [5]int{a, -b, -a, b, a}
					im := [5]int{b, a, -b, -a, b}
					switch kp {
					case 1:
						wave19Check(t, p, re[k])
					case 2:
						wave19Check(t, p, im[k])
					default:
						if id == "complex.gauss.units" {
							wave19Check(t, p, a*a+b*b)
						} else {
							wave19Check(t, p, re[k])
						}
					}
				case "complex.gauss.division":
					n := wave19FactInt(t, p, "n")
					switch kp {
					case 1:
						a := wave19FactInt(t, p, "a")
						q := int(roundHalf(float64(a) / float64(n)))
						wave19Check(t, p, q)
					case 2:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						q := int(roundHalf(float64(a) / float64(n)))
						wave19Check(t, p, (a-q*n)*(a-q*n)+b*b)
					default:
						wave19Check(t, p, n/2)
					}
				case "complex.gauss.gcd":
					switch kp {
					case 1:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						wave19Check(t, p, a*a+b*b)
					case 2:
						q := wave19FactInt(t, p, "p")
						wave19Check(t, p, q*q)
					default:
						wave19Check(t, p, 1)
					}
				case "complex.gauss.primes":
					switch kp {
					case 1:
						n := wave19FactInt(t, p, "n")
						count := 0
						for q := 3; q <= n; q += 4 {
							if isPrime(q) {
								count++
							}
						}
						wave19Check(t, p, count)
					case 2:
						q := wave19FactInt(t, p, "p")
						if q%4 == 1 {
							wave19Check(t, p, 1)
						} else {
							wave19Check(t, p, 0)
						}
					default:
						wave19Check(t, p, wave19FactInt(t, p, "p"))
					}
				case "complex.gauss.two_squares":
					switch kp {
					case 1:
						n := wave19FactInt(t, p, "n")
						one, three := wave19DivisorCounts(n)
						wave19Check(t, p, 4*(one-three))
					case 2:
						wave19Check(t, p, wave19MinSquares(wave19FactInt(t, p, "n")))
					default:
						wave19Check(t, p, 8)
					}
				case "complex.gauss.associates":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					switch kp {
					case 1:
						wave19Check(t, p, 4)
					case 2:
						wave19Check(t, p, a*a+b*b)
					default:
						wave19Check(t, p, a-b)
					}
				case "complex.gauss.pythagorean":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					leg := a*a - b*b
					if leg < 0 {
						leg = -leg
					}
					switch kp {
					case 1:
						wave19Check(t, p, a*a+b*b)
					case 2:
						wave19Check(t, p, leg)
					default:
						wave19Check(t, p, a*a+b*b+leg+2*a*b)
					}
				case "complex.gauss.theta", "complex.spec.theta":
					switch kp {
					case 1:
						wave19Check(t, p, 2*wave19FactInt(t, p, "k")+1)
					case 2:
						n := wave19FactInt(t, p, "n")
						wave19Check(t, p, n*n)
					default:
						n := wave19FactInt(t, p, "n")
						one, three := wave19DivisorCounts(n)
						wave19Check(t, p, 4*(one-three))
					}
				case "complex.gauss.eisenstein":
					switch kp {
					case 1:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						wave19Check(t, p, a*a-a*b+b*b)
					case 2:
						wave19Check(t, p, 6)
					default:
						a := wave19FactInt(t, p, "aa")
						wave19Check(t, p, a*a)
					}
				case "complex.geom.area":
					switch kp {
					case 1:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						e, f := wave19FactInt(t, p, "e"), wave19FactInt(t, p, "f")
						det := (c-a)*(f-b) - (d-b)*(e-a)
						if det < 0 {
							det = -det
						}
						wave19Check(t, p, det)
					case 2:
						wave19Check(t, p, wave19FactInt(t, p, "p")*wave19FactInt(t, p, "q")/2)
					default:
						wave19Check(t, p, 0)
					}
				case "complex.geom.circle":
					switch kp {
					case 1, 3:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						wave19Check(t, p, (c-a)*(c-a)+(d-b)*(d-b))
					default:
						wave19Check(t, p, wave19FactInt(t, p, "p")+wave19FactInt(t, p, "r"))
					}
				case "complex.geom.reflection":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					switch kp {
					case 1:
						wave19Check(t, p, -b)
					case 2:
						wave19Check(t, p, -a)
					default:
						wave19Check(t, p, b)
					}
				case "complex.geom.spiral":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "k"))
					case 2:
						wave19Check(t, p, 90*wave19FactInt(t, p, "k"))
					default:
						wave19Check(t, p, wave19FactInt(t, p, "k1")*wave19FactInt(t, p, "k2"))
					}
				case "complex.geom.ptolemy":
					p1, q := wave19FactInt(t, p, "p"), wave19FactInt(t, p, "q")
					switch kp {
					case 1:
						wave19Check(t, p, p1*p1+q*q)
					case 2:
						wave19Check(t, p, 2*p1*q)
					default:
						wave19Check(t, p, (p1-q)*(p1-q))
					}
				case "complex.geom.triangle":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					switch kp {
					case 1:
						wave19Check(t, p, a+b)
					case 2:
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						wave19Check(t, p, (a+c)*(a+c)+(b+d)*(b+d))
					default:
						wave19Check(t, p, a*a+b*b)
					}
				case "complex.geom.centroid":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
					e, f := wave19FactInt(t, p, "e"), wave19FactInt(t, p, "f")
					switch kp {
					case 1:
						wave19Check(t, p, a+c+e)
					case 2:
						wave19Check(t, p, b+d+f)
					default:
						wave19Check(t, p, a+c+e)
					}
				case "complex.geom.perpendicular":
					switch kp {
					case 1:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						wave19Check(t, p, a*c+b*d)
					case 2:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						wave19Check(t, p, a*d-b*c)
					default:
						wave19Check(t, p, 4)
					}
				case "complex.geom.lattice":
					switch kp {
					case 1:
						n := wave19FactInt(t, p, "n")
						wave19Check(t, p, (n+1)*(n+1))
					case 2:
						a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
						c, d := wave19FactInt(t, p, "c"), wave19FactInt(t, p, "d")
						wave19Check(t, p, (c-a)*(c-a)+(d-b)*(d-b))
					default:
						wave19Check(t, p, wave19FactInt(t, p, "w")*wave19FactInt(t, p, "h"))
					}
				case "complex.poly.fta":
					d := wave19FactInt(t, p, "d")
					if kp == 2 {
						wave19Check(t, p, 1)
					} else {
						wave19Check(t, p, d)
					}
				case "complex.poly.vieta_sum":
					switch kp {
					case 1:
						wave19Check(t, p, -wave19FactInt(t, p, "b"))
					case 2:
						wave19Check(t, p, -wave19FactInt(t, p, "b")/wave19FactInt(t, p, "a"))
					default:
						wave19Check(t, p, 0)
					}
				case "complex.poly.vieta_product":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "c"))
					case 2:
						wave19Check(t, p, -wave19FactInt(t, p, "d"))
					default:
						n := wave19FactInt(t, p, "n")
						if n%2 == 0 {
							wave19Check(t, p, -1)
						} else {
							wave19Check(t, p, 1)
						}
					}
				case "complex.poly.unity_sum":
					n := wave19FactInt(t, p, "n")
					switch kp {
					case 1:
						wave19Check(t, p, n)
					case 2:
						wave19Check(t, p, 0)
					default:
						j := wave19FactInt(t, p, "j")
						if j%n == 0 {
							wave19Check(t, p, n)
						} else {
							wave19Check(t, p, 0)
						}
					}
				case "complex.poly.cyclotomic":
					switch kp {
					case 1:
						wave19Check(t, p, wave19Totient(wave19FactInt(t, p, "n")))
					case 2:
						wave19Check(t, p, wave19FactInt(t, p, "p")-1)
					default:
						wave19Check(t, p, wave19FactInt(t, p, "n"))
					}
				case "complex.poly.chebyshev":
					n := wave19FactInt(t, p, "n")
					switch kp {
					case 1:
						wave19Check(t, p, 1)
					case 2:
						if n%2 == 1 {
							wave19Check(t, p, -1)
						} else {
							wave19Check(t, p, 1)
						}
					default:
						if n%2 == 1 {
							wave19Check(t, p, 0)
						} else if (n/2)%2 == 1 {
							wave19Check(t, p, -1)
						} else {
							wave19Check(t, p, 1)
						}
					}
				case "complex.poly.discriminant":
					switch kp {
					case 1:
						a, b, c := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b"), wave19FactInt(t, p, "c")
						wave19Check(t, p, b*b-4*a*c)
					case 2:
						wave19Check(t, p, 1)
					default:
						b, c := wave19FactInt(t, p, "b"), wave19FactInt(t, p, "c")
						disc := b*b - 4*c
						if disc > 0 {
							wave19Check(t, p, 2)
						} else if disc == 0 {
							wave19Check(t, p, 1)
						} else {
							wave19Check(t, p, 0)
						}
					}
				case "complex.poly.conjugate_pairs":
					d := wave19FactInt(t, p, "d")
					switch kp {
					case 1:
						wave19Check(t, p, d-wave19FactInt(t, p, "r"))
					case 2:
						wave19Check(t, p, 1)
					default:
						r := wave19FactInt(t, p, "r")
						wave19Check(t, p, (d-r)/2)
					}
				case "complex.poly.multiplicity":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "k"))
					case 2:
						wave19Check(t, p, wave19FactInt(t, p, "d")-wave19FactInt(t, p, "k")+1)
					default:
						wave19Check(t, p, wave19FactInt(t, p, "d"))
					}
				case "complex.poly.resultant":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					switch kp {
					case 1:
						wave19Check(t, p, (a-wave19FactInt(t, p, "c"))*(b-wave19FactInt(t, p, "c")))
					case 2:
						wave19Check(t, p, a-b)
					default:
						c := wave19FactInt(t, p, "c")
						if c == a || c == b {
							wave19Check(t, p, 1)
						} else {
							wave19Check(t, p, 0)
						}
					}
				case "complex.meth.simple_pole":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "a")+wave19FactInt(t, p, "b"))
					case 2:
						d, c := wave19FactInt(t, p, "d"), wave19FactInt(t, p, "c")
						wave19Check(t, p, d*d+c)
					default:
						wave19Check(t, p, 2)
					}
				case "complex.meth.higher_pole":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "m"))
					case 2:
						wave19Check(t, p, 1)
					default:
						wave19Check(t, p, -1)
					}
				case "complex.meth.partial_fractions":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "coefA"))
					case 2:
						wave19Check(t, p, wave19FactInt(t, p, "coefB"))
					default:
						wave19Check(t, p, wave19FactInt(t, p, "coefA")+wave19FactInt(t, p, "coefB"))
					}
				case "complex.meth.residue_sum":
					a, b := wave19FactInt(t, p, "a"), wave19FactInt(t, p, "b")
					switch kp {
					case 1:
						wave19Check(t, p, a+b)
					case 2:
						wave19Check(t, p, -(a + b))
					default:
						wave19Check(t, p, 0)
					}
				case "complex.meth.winding":
					k := wave19FactInt(t, p, "k")
					switch kp {
					case 1:
						wave19Check(t, p, k)
					case 2:
						wave19Check(t, p, -k)
					default:
						wave19Check(t, p, 0)
					}
				case "complex.meth.pole_count":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "n"))
					case 2:
						wave19Check(t, p, 2)
					default:
						wave19Check(t, p, 0)
					}
				case "complex.meth.zero_count":
					switch kp {
					case 1, 2:
						wave19Check(t, p, wave19FactInt(t, p, "n"))
					default:
						wave19Check(t, p, 0)
					}
				case "complex.meth.order":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "k"))
					case 2:
						wave19Check(t, p, wave19FactInt(t, p, "k1")+wave19FactInt(t, p, "k2"))
					default:
						wave19Check(t, p, wave19FactInt(t, p, "k1")-wave19FactInt(t, p, "k2"))
					}
				case "complex.meth.jordan":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "dn")-wave19FactInt(t, p, "nn"))
					case 2:
						wave19Check(t, p, 2)
					default:
						wave19Check(t, p, 1)
					}
				case "complex.meth.branch_cut":
					if kp == 3 {
						wave19Check(t, p, 1)
					} else {
						wave19Check(t, p, 2)
					}
				case "complex.spec.gamma_half":
					switch kp {
					case 1:
						wave19Check64(t, p, wave19Factorial(t, wave19FactInt(t, p, "n")))
					case 2:
						n := wave19FactInt(t, p, "n")
						num := wave19Factorial(t, 2*n)
						den := int64(1)
						for i := 0; i < n; i++ {
							den *= 4
						}
						den *= wave19Factorial(t, n)
						g := wave19GCD64(num, den)
						wave19Check64(t, p, den/g)
					default:
						wave19Check64(t, p, wave19Factorial(t, wave19FactInt(t, p, "n")-1))
					}
				case "complex.spec.beta":
					m, n := wave19FactInt(t, p, "m"), wave19FactInt(t, p, "n")
					num := wave19Factorial(t, m-1) * wave19Factorial(t, n-1)
					den := wave19Factorial(t, m+n-1)
					g := wave19GCD64(num, den)
					switch kp {
					case 1:
						wave19Check64(t, p, num)
					case 2:
						wave19Check64(t, p, den)
					default:
						wave19Check64(t, p, den/g)
					}
				case "complex.spec.zeta_even":
					switch kp {
					case 1:
						k := wave19FactInt(t, p, "k")
						m := int64(wave19FactInt(t, p, "m"))
						dens := []int64{6, 90, 945, 9450, 93555, 638512875, 18243225, 325641566250}
						wave19Check64(t, p, m*dens[k-1])
					case 2:
						wave19Check(t, p, -wave19FactInt(t, p, "k")*wave19FactInt(t, p, "m"))
					default:
						wave19Check(t, p, 0)
					}
				case "complex.spec.bessel":
					switch kp {
					case 1:
						n := wave19FactInt(t, p, "n")
						if n == 0 {
							wave19Check(t, p, 1)
						} else {
							wave19Check(t, p, 0)
						}
					case 2:
						wave19Check(t, p, 0)
					default:
						n := wave19FactInt(t, p, "n")
						den := int64(1)
						for i := 0; i < n; i++ {
							den *= 2
						}
						den *= wave19Factorial(t, n)
						wave19Check64(t, p, den)
					}
				case "complex.spec.airy":
					if kp == 3 {
						wave19Check(t, p, 0)
					} else {
						wave19Check(t, p, 2)
					}
				case "complex.spec.hypergeometric":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "k"))
					case 2:
						wave19Check(t, p, wave19FactInt(t, p, "p"))
					default:
						wave19Check(t, p, 1)
					}
				case "complex.spec.polygamma":
					switch kp {
					case 1:
						wave19Check(t, p, wave19FactInt(t, p, "n"))
					case 2:
						wave19Check(t, p, -1)
					default:
						wave19Check(t, p, 1)
					}
				case "complex.spec.erf":
					if kp == 1 {
						wave19Check(t, p, 0)
					} else {
						wave19Check(t, p, 1)
					}
				case "complex.spec.polylog":
					switch kp {
					case 1:
						wave19Check(t, p, 0)
					case 2:
						wave19Check(t, p, 1)
					default:
						k := wave19FactInt(t, p, "k")
						m := int64(wave19FactInt(t, p, "m"))
						dens := []int64{6, 90, 945, 9450, 93555, 638512875, 18243225, 325641566250}
						wave19Check64(t, p, m*dens[k-1])
					}
				default:
					t.Fatalf("%s: no recompute case", id)
				}
			}
		}
	}
}

func roundHalf(v float64) float64 {
	if v < 0 {
		return -float64(int(-v + 0.5))
	}
	return float64(int(v + 0.5))
}

func wave19GCD64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
