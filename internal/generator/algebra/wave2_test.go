package algebra

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func factInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("missing fact %q", key)
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("fact %q is not an integer: %q", key, raw)
	}
	return n
}

func intPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

func factInts(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	raw := p.Facts[key]
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			t.Fatalf("fact %q has a non-integer part: %q", key, part)
		}
		out[i] = n
	}
	return out
}

func TestWave2AnswersMatchFacts(t *testing.T) {
	cases := []struct {
		name  string
		gen   generator.Generator
		check func(t *testing.T, p generator.Problem, kp int)
	}{
		{
			name: "binomial",
			gen:  &binomialGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				n := factInt(t, p, "n")
				a := factInt(t, p, "a")
				k := factInt(t, p, "k")
				switch kp {
				case 1:
					want := binomCoef(n, k) * intPow(a, n-k)
					if got := factInt(t, p, "coefficient"); got != want {
						t.Fatalf("coefficient fact = %d, want %d", got, want)
					}
					if p.Answer != strconv.Itoa(want) {
						t.Fatalf("answer = %s, want %d", p.Answer, want)
					}
				case 2:
					want := intPow(a, n)
					if got := factInt(t, p, "constant"); got != want {
						t.Fatalf("constant fact = %d, want %d", got, want)
					}
					if p.Answer != strconv.Itoa(want) {
						t.Fatalf("answer = %s, want %d", p.Answer, want)
					}
				default:
					want := intPow(a+1, n)
					if got := factInt(t, p, "sum"); got != want {
						t.Fatalf("sum fact = %d, want %d", got, want)
					}
					if p.Answer != strconv.Itoa(want) {
						t.Fatalf("answer = %s, want %d", p.Answer, want)
					}
				}
			},
		},
		{
			name: "infinite geometric",
			gen:  &infiniteGeometricGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				first := factInt(t, p, "first")
				d := factInt(t, p, "denominator")
				sum := factInt(t, p, "sum")
				if sum*(d-1) != first*d {
					t.Fatalf("sum %d does not satisfy S(1-r)=a", sum)
				}
				want := []int{sum, d, first}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "exponential from points",
			gen:  &expFromPointsGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				y3 := factInt(t, p, "y3")
				if y3 != a*intPow(b, 3) {
					t.Fatalf("y3 = %d, want a*b^3", y3)
				}
				want := []int{b, a, y3}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "log scale",
			gen:  &logScaleGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				exp := factInt(t, p, "exponent")
				value := factInt(t, p, "value")
				if value != intPow(10, exp) {
					t.Fatalf("value = %d, want 10^%d", value, exp)
				}
				want := exp
				if kp == 2 {
					want = value
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "linear regression",
			gen:  &linearRegressionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				var xs, ys []float64
				for i := 1; i <= 4; i++ {
					xs = append(xs, float64(factInt(t, p, "x"+strconv.Itoa(i))))
					ys = append(ys, float64(factInt(t, p, "y"+strconv.Itoa(i))))
				}
				var mx, my float64
				for i := range xs {
					mx += xs[i] / 4
					my += ys[i] / 4
				}
				var num, den float64
				for i := range xs {
					num += (xs[i] - mx) * (ys[i] - my)
					den += (xs[i] - mx) * (xs[i] - mx)
				}
				slope := int(math.Round(num / den))
				intercept := int(math.Round(my - float64(slope)*mx))
				if got := factInt(t, p, "slope"); got != slope {
					t.Fatalf("slope fact = %d, want %d", got, slope)
				}
				if got := factInt(t, p, "intercept"); got != intercept {
					t.Fatalf("intercept fact = %d, want %d", got, intercept)
				}
				want := []int{slope, intercept, 3*slope + intercept}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "parabola focus",
			gen:  &parabolaFocusGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				pp := factInt(t, p, "p")
				want := []int{pp, -pp, 4 * pp}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "ellipse foci",
			gen:  &ellipseFociGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := factInt(t, p, "a")
				b := factInt(t, p, "b")
				c := factInt(t, p, "c")
				if c*c != a*a-b*b {
					t.Fatalf("c^2 = %d, want a^2-b^2 = %d", c*c, a*a-b*b)
				}
				want := []int{c, 2 * a, c}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "absolute inequality",
			gen:  &absoluteInequalityGraphGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				center := factInt(t, p, "center")
				radius := factInt(t, p, "radius")
				if got := factInt(t, p, "count"); got != 2*radius+1 {
					t.Fatalf("count = %d, want 2r+1", got)
				}
				want := []int{center, radius, 2*radius + 1}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "zeros and multiplicity",
			gen:  &zerosMultiplicityGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				mults := factInts(t, p, "multiplicities")
				roots := factInts(t, p, "roots")
				degree := 0
				for _, m := range mults {
					degree += m
				}
				if got := factInt(t, p, "degree"); got != degree {
					t.Fatalf("degree = %d, want %d", got, degree)
				}
				if got := factInt(t, p, "distinct"); got != len(roots) {
					t.Fatalf("distinct = %d, want %d", got, len(roots))
				}
				want := 0
				switch kp {
				case 1:
					root := factInt(t, p, "asked_root")
					mult := factInt(t, p, "asked_multiplicity")
					for i, r := range roots {
						if r == root {
							want = mults[i]
							if mult != want {
								t.Fatalf("asked multiplicity = %d, want %d", mult, want)
							}
						}
					}
				case 2:
					want = degree
				default:
					want = len(roots)
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
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
					if p.Answer == "" {
						t.Fatal("empty answer")
					}
					tc.check(t, p, kp)
				}
			}
		})
	}
}
