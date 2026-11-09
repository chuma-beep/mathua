package algebra

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// w3TestParsePoly reads the rendered polynomial back into coefficients indexed
// by degree, so a test can count sign changes rather than trust the generator's
// own count.
func w3TestParsePoly(poly string) []int {
	s := strings.ReplaceAll(poly, " ", "")
	terms := make([]string, 0, 8)
	start := 0
	for i := 1; i < len(s); i++ {
		if (s[i] == '+' || s[i] == '-') && i > start {
			terms = append(terms, s[start:i])
			start = i
		}
	}
	terms = append(terms, s[start:])
	byDegree := map[int]int{}
	maxDegree := 0
	for _, term := range terms {
		if term == "" {
			continue
		}
		coef, degree := 0, 0
		if idx := strings.Index(term, "x"); idx >= 0 {
			head := term[:idx]
			switch head {
			case "", "+":
				coef = 1
			case "-":
				coef = -1
			default:
				coef, _ = strconv.Atoi(head)
			}
			rest := term[idx+1:]
			if strings.HasPrefix(rest, "^") {
				degree, _ = strconv.Atoi(rest[1:])
			} else {
				degree = 1
			}
		} else {
			coef, _ = strconv.Atoi(term)
			degree = 0
		}
		byDegree[degree] += coef
		if degree > maxDegree {
			maxDegree = degree
		}
	}
	coefs := make([]int, maxDegree+1)
	for degree, coef := range byDegree {
		coefs[degree] = coef
	}
	return coefs
}

func w3TestSignChanges(coefs []int) int {
	changes, prev := 0, 0
	for _, c := range coefs {
		if c == 0 {
			continue
		}
		if prev != 0 && (c > 0) != (prev > 0) {
			changes++
		}
		prev = c
	}
	return changes
}

func w3TestDivisors(n int) int {
	if n < 0 {
		n = -n
	}
	count := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			count++
		}
	}
	return count
}

func w3TestLinear(s string) (int, int) {
	coefs := w3TestParsePoly(s)
	if len(coefs) < 2 {
		return coefs[0], 0
	}
	return coefs[1], coefs[0]
}

func w3TestFib(n int) int {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a
}

func w3TestIntRoot(value, power int) int {
	for r := 1; r <= value; r++ {
		if intPow(r, power) == value {
			return r
		}
	}
	return 0
}

func w3TestBaseFrom(argument, exponent string) int {
	if exponent == "1/2" {
		n, _ := strconv.Atoi(argument)
		return n * n
	}
	e, _ := strconv.Atoi(exponent)
	if e > 0 {
		n, _ := strconv.Atoi(argument)
		return int(math.Round(math.Pow(float64(n), 1/float64(e))))
	}
	parts := strings.Split(argument, "/")
	d, _ := strconv.Atoi(parts[len(parts)-1])
	return int(math.Round(math.Pow(float64(d), 1/float64(-e))))
}

func w3CheckInt(t *testing.T, p generator.Problem, want int) {
	t.Helper()
	if p.Answer != strconv.Itoa(want) {
		t.Fatalf("answer = %s, want %d\nquestion: %s\nfacts: %v", p.Answer, want, p.Question, p.Facts)
	}
}

func w3TestShift(t *testing.T, equation, variable string) int {
	t.Helper()
	re := regexp.MustCompile(`\(` + variable + ` ([+-]) (\d+)\)\^\{2\}`)
	m := re.FindStringSubmatch(equation)
	if m == nil {
		t.Fatalf("equation %q has no shifted square in %s", equation, variable)
	}
	n, err := strconv.Atoi(m[2])
	if err != nil {
		t.Fatal(err)
	}
	if m[1] == "+" {
		return -n
	}
	return n
}

func TestWave3QuestionDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"reciprocal", &w3ReciprocalGen{}},
		{"bounds", &w3PolyBoundsGen{}},
		{"operations", &w3FuncOperationsGen{}},
		{"inverse graph", &w3InverseGraphGen{}},
		{"composition", &w3CompositionGraphGen{}},
		{"exponential regression", &w3ExpRegressionGen{}},
		{"log domain", &w3LogDomainGen{}},
		{"log base equation", &w3LogBaseEquationGen{}},
		{"arithmetic means", &w3ArithmeticMeanGen{}},
		{"geometric means", &w3GeometricMeanGen{}},
		{"fibonacci", &w3FibonacciGen{}},
		{"conic intersection", &w3ConicIntersectionGen{}},
		{"conic translation", &w3ConicTranslationGen{}},
		{"conic tangent", &w3ConicTangentGen{}},
		{"matrix inverse", &w3MatrixInverseGen{}},
		{"quadratic regression", &w3QuadraticRegressionGen{}},
		{"power regression", &w3PowerRegressionGen{}},
		{"logarithmic regression", &w3LogarithmicRegressionGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", tc.name, kp, len(seen))
			}
		}
	}
}

func TestWave3FirstQuestionsDifferByOrdinal(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"reciprocal", &w3ReciprocalGen{}},
		{"bounds", &w3PolyBoundsGen{}},
		{"operations", &w3FuncOperationsGen{}},
		{"inverse graph", &w3InverseGraphGen{}},
		{"composition", &w3CompositionGraphGen{}},
		{"exponential regression", &w3ExpRegressionGen{}},
		{"log domain", &w3LogDomainGen{}},
		{"log base equation", &w3LogBaseEquationGen{}},
		{"arithmetic means", &w3ArithmeticMeanGen{}},
		{"geometric means", &w3GeometricMeanGen{}},
		{"fibonacci", &w3FibonacciGen{}},
		{"conic intersection", &w3ConicIntersectionGen{}},
		{"conic translation", &w3ConicTranslationGen{}},
		{"conic tangent", &w3ConicTangentGen{}},
		{"matrix inverse", &w3MatrixInverseGen{}},
		{"quadratic regression", &w3QuadraticRegressionGen{}},
		{"power regression", &w3PowerRegressionGen{}},
		{"logarithmic regression", &w3LogarithmicRegressionGen{}},
	}
	for _, tc := range gens {
		first := make([]string, 3)
		for kp := 1; kp <= 3; kp++ {
			first[kp-1] = tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question
		}
		if first[0] == first[1] || first[1] == first[2] || first[0] == first[2] {
			t.Fatalf("%s: two ordinals produced the same question", tc.name)
		}
	}
}

func TestWave3AnswersMatchFacts(t *testing.T) {
	cases := []struct {
		name  string
		gen   generator.Generator
		check func(t *testing.T, p generator.Problem, kp int)
	}{
		{
			name: "reciprocal",
			gen:  &w3ReciprocalGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				coefs := w3TestParsePoly(p.Facts["poly"])
				if len(coefs) != 5 {
					t.Fatalf("degree = %d, want 4", len(coefs)-1)
				}
				switch kp {
				case 1:
					w3CheckInt(t, p, coefs[4])
				case 2:
					if coefs[0] == 0 || coefs[1]%coefs[0] != 0 {
						t.Fatalf("sum of reciprocals is not an integer: a1=%d a0=%d", coefs[1], coefs[0])
					}
					w3CheckInt(t, p, -coefs[1]/coefs[0])
				default:
					if coefs[0] == 0 || coefs[4]%coefs[0] != 0 {
						t.Fatalf("product of reciprocals is not an integer: a4=%d a0=%d", coefs[4], coefs[0])
					}
					w3CheckInt(t, p, coefs[4]/coefs[0])
				}
			},
		},
		{
			name: "bounds",
			gen:  &w3PolyBoundsGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				leading := factInt(t, p, "leading")
				constant := factInt(t, p, "constant")
				divisors := w3TestDivisors(constant)
				maxAbs := 0
				for _, c := range w3TestParsePoly(p.Facts["poly"]) {
					if w3Abs(c) > maxAbs {
						maxAbs = w3Abs(c)
					}
				}
				coefs := w3TestParsePoly(p.Facts["poly"])
				if coefs[len(coefs)-1] != leading || coefs[0] != constant {
					t.Fatalf("poly %q does not have leading %d and constant %d", p.Facts["poly"], leading, constant)
				}
				cauchy := fracOrInt(leading+maxAbs, leading)
				candidates := 2 * divisors * w3TestDivisors(leading)
				switch kp {
				case 1:
					w3CheckInt(t, p, divisors)
				case 2:
					if p.Answer != cauchy {
						t.Fatalf("answer = %s, want cauchy bound %s", p.Answer, cauchy)
					}
				default:
					w3CheckInt(t, p, candidates)
				}
			},
		},
		{
			name: "combining functions",
			gen:  &w3FuncOperationsGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				fv := factInt(t, p, "f_value")
				gv := factInt(t, p, "g_value")
				want := []int{fv + gv, fv - gv, fv * gv}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "inverse graph",
			gen:  &w3InverseGraphGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				x0 := factInt(t, p, "point_x")
				y0 := factInt(t, p, "point_y")
				want := []int{x0, y0, x0 - y0}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "composition from graphs",
			gen:  &w3CompositionGraphGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a, b := w3TestLinear(p.Facts["f"])
				c, d := w3TestLinear(p.Facts["g"])
				x := factInt(t, p, "x")
				fv := factInt(t, p, "f_value")
				gv := factInt(t, p, "g_value")
				if fv != a*x+b || gv != c*x+d {
					t.Fatalf("published values do not match the functions: f(%d)=%d, g(%d)=%d", x, fv, x, gv)
				}
				want := []int{a*gv + b, c*fv + d, a*fv + b}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "exponential regression",
			gen:  &w3ExpRegressionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				for i := 0; i <= 3; i++ {
					if got, want := factInt(t, p, "y"+strconv.Itoa(i)), a*intPow(b, i); got != want {
						t.Fatalf("y%d = %d, want %d", i, got, want)
					}
				}
				want := []int{b, a, a * intPow(b, 4)}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "log domain",
			gen:  &w3LogDomainGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				threshold := factInt(t, p, "threshold")
				var want int
				switch kp {
				case 1:
					want = threshold + 1
				case 2:
					want = threshold - 1
				default:
					want = int(math.Floor(math.Sqrt(float64(threshold)))) + 1
				}
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "log base equation",
			gen:  &w3LogBaseEquationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				want := w3TestBaseFrom(p.Facts["argument"], p.Facts["exponent"])
				w3CheckInt(t, p, want)
				if factInt(t, p, "base") != want {
					t.Fatalf("base fact = %s, want %d", p.Facts["base"], want)
				}
			},
		},
		{
			name: "arithmetic means",
			gen:  &w3ArithmeticMeanGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				first := factInt(t, p, "first")
				last := factInt(t, p, "last")
				count := factInt(t, p, "count")
				difference := (last - first) / (count + 1)
				if got := factInt(t, p, "difference"); got != difference {
					t.Fatalf("difference fact = %d, want %d", got, difference)
				}
				var want int
				switch kp {
				case 1:
					want = (first + last) / 2
				case 2:
					want = difference
				default:
					want = first + factInt(t, p, "position")*difference
				}
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "geometric means",
			gen:  &w3GeometricMeanGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				first := factInt(t, p, "first")
				last := factInt(t, p, "last")
				count := factInt(t, p, "count")
				ratio := w3TestIntRoot(last/first, count+1)
				if ratio == 0 {
					t.Fatalf("last/first = %d is not a %d-th power", last/first, count+1)
				}
				var want int
				switch kp {
				case 1:
					want = first * ratio
				case 2:
					want = ratio
				default:
					want = first * intPow(ratio, factInt(t, p, "position"))
				}
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "fibonacci",
			gen:  &w3FibonacciGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				n := factInt(t, p, "n")
				var want int
				switch kp {
				case 1:
					want = w3TestFib(n)
				case 2:
					want = w3TestFib(n) + w3TestFib(n+1)
				default:
					want = w3TestFib(n+2) - 1
				}
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "conic intersection",
			gen:  &w3ConicIntersectionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				m := factInt(t, p, "m")
				b := factInt(t, p, "b")
				d := m*m + 4*b
				if factInt(t, p, "discriminant") != d {
					t.Fatalf("discriminant fact = %s, want %d", p.Facts["discriminant"], d)
				}
				count := 0
				if d > 0 {
					count = 2
				} else if d == 0 {
					count = 1
				}
				want := []int{count, d, m}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "conic translation",
			gen:  &w3ConicTranslationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				h := w3TestShift(t, p.Facts["equation"], "x")
				k := w3TestShift(t, p.Facts["equation"], "y")
				if factInt(t, p, "h") != h || factInt(t, p, "k") != k {
					t.Fatalf("facts h=%s k=%s, equation says %d and %d", p.Facts["h"], p.Facts["k"], h, k)
				}
				want := []int{h, k, h + k}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "conic tangent",
			gen:  &w3ConicTangentGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				px := factInt(t, p, "point_x")
				py := factInt(t, p, "point_y")
				switch kp {
				case 1:
					if p.Facts["curve"] != "circle" {
						t.Fatalf("curve = %q, want circle", p.Facts["curve"])
					}
					want := fracOrInt(-px, py)
					if p.Answer != want {
						t.Fatalf("answer = %s, want slope %s", p.Answer, want)
					}
				case 2:
					if py != px*px {
						t.Fatalf("point (%d, %d) is not on y = x^2", px, py)
					}
					w3CheckInt(t, p, 2*px)
				default:
					if (py*py)%(4*px) != 0 {
						t.Fatalf("point (%d, %d) is not on y^2 = 4px", px, py)
					}
					pCoef := py * py / (4 * px)
					if py%(2*pCoef) != 0 {
						t.Fatalf("t = %d/%d is not an integer", py, 2*pCoef)
					}
					tParam := py / (2 * pCoef)
					w3CheckInt(t, p, -pCoef*tParam*tParam)
				}
			},
		},
		{
			name: "matrix inverse",
			gen:  &w3MatrixInverseGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				c := factInt(t, p, "c")
				d := factInt(t, p, "d")
				e := factInt(t, p, "e")
				f := factInt(t, p, "f")
				det := a*d - b*c
				if det == 0 {
					t.Fatalf("singular matrix: a=%d b=%d c=%d d=%d", a, b, c, d)
				}
				x := (d*e - b*f) / det
				y := (a*f - c*e) / det
				if a*x+b*y != e || c*x+d*y != f {
					t.Fatalf("solution (%d, %d) does not satisfy the system", x, y)
				}
				want := []int{det, x, y}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "quadratic regression",
			gen:  &w3QuadraticRegressionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				c := factInt(t, p, "c")
				want := []int{a, c, 9*a + 3*b + c}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "power regression",
			gen:  &w3PowerRegressionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				if factInt(t, p, "y2") != a*intPow(2, b) || factInt(t, p, "y3") != a*intPow(3, b) {
					t.Fatalf("published values do not follow y = %d x^%d", a, b)
				}
				want := []int{b, a, a * intPow(4, b)}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
		{
			name: "logarithmic regression",
			gen:  &w3LogarithmicRegressionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				if factInt(t, p, "y2") != a+b || factInt(t, p, "y3") != a+2*b {
					t.Fatalf("published values do not follow y = %d + %d ln x", a, b)
				}
				want := []int{b, a, a + 3*b}[kp-1]
				w3CheckInt(t, p, want)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for kp := 1; kp <= 3; kp++ {
				for seed := int64(0); seed < 300; seed++ {
					p := tc.gen.Generate(generator.GeneratorContext{KPOrdinal: kp, Seed: seed, Difficulty: 0.5})
					if p.KPOrdinal != kp {
						t.Fatalf("ordinal = %d, want %d", p.KPOrdinal, kp)
					}
					if p.Question == "" || p.Answer == "" || p.Explanation == "" {
						t.Fatalf("blank question, answer or explanation: %+v", p)
					}
					tc.check(t, p, kp)
					if !strings.Contains(p.Explanation, p.Answer) {
						t.Fatalf("explanation does not state the answer %q:\n%s", p.Answer, p.Explanation)
					}
				}
			}
		})
	}
}
