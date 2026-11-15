package calculus

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestWave2AnswersMatchFacts(t *testing.T) {
	cases := []struct {
		name  string
		gen   generator.Generator
		check func(t *testing.T, p generator.Problem, kp int)
	}{
		{
			name: "conjugate limits",
			gen:  &conjugateLimitGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				m := calcFactInt(t, p, "m")
				base := calcFactInt(t, p, "base")
				want := []int{m, 2 * m, 2 * base}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "complex fraction limits",
			gen:  &complexFractionLimitGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				b := calcFactInt(t, p, "b")
				want := a * a
				if b > 0 {
					want = b * b * b / 2
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "absolute value limits",
			gen:  &absoluteValueLimitGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				want := 0
				switch kp {
				case 1:
					want = -1
					if p.Facts["mirrored"] == "true" {
						want = 1
					}
				case 2:
					want = 1
					if p.Facts["mirrored"] == "true" {
						want = -1
					}
				default:
					if p.Facts["form"] == "ratio" {
						want = 1
					}
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "normal lines",
			gen:  &normalLineGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				var want int
				switch p.Facts["function"] {
				case "x^3":
					want = 3 * a * a
				case "x^2+x":
					want = 2*a + 1
				case "x^4":
					want = 4 * a * a * a
				default:
					want = 2 * a
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "piecewise derivatives",
			gen:  &piecewiseDerivativeGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				b := calcFactInt(t, p, "b")
				c := calcFactInt(t, p, "c")
				want := []int{a, c, b}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "fence optimization",
			gen:  &fenceOptimizationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				perimeter := calcFactInt(t, p, "perimeter")
				width := perimeter / 4
				length := perimeter - 2*width
				if p.Facts["divided"] == "true" {
					length = width
				}
				area := width * length
				want := []int{width, length, area}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "revenue optimization",
			gen:  &revenueOptimizationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				b := calcFactInt(t, p, "b")
				intercept := calcFactInt(t, p, "intercept")
				price := intercept / (2 * b)
				quantity := intercept - b*price
				revenue := price * quantity
				want := []int{price, quantity, revenue}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "ladder rate",
			gen:  &ladderRateGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				x := calcFactInt(t, p, "x")
				y := calcFactInt(t, p, "y")
				l := calcFactInt(t, p, "ladder")
				factor := calcFactInt(t, p, "factor")
				if x*x+y*y != l*l {
					t.Fatalf("%d-%d-%d is not a right triangle", x, y, l)
				}
				want := y * factor
				switch kp {
				case 2:
					want = calcFactInt(t, p, "sign")
				case 3:
					want = x * factor
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "cone rate",
			gen:  &coneRateGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				h := calcFactInt(t, p, "h")
				c := calcFactInt(t, p, "c")
				dh := calcFactInt(t, p, "dh")
				ratio := calcFactInt(t, p, "ratio")
				if ratio*ratio*c != h*h*dh {
					t.Fatalf("c=%d does not satisfy h^2 dh / ratio^2", c)
				}
				want := []int{dh, c, h}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "total distance",
			gen:  &totalDistanceGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				b := calcFactInt(t, p, "b")
				left := a * a / 2
				right := (b - a) * (b - a) / 2
				want := []int{a, left, left + right}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "piecewise integral",
			gen:  &piecewiseIntegralGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				b := calcFactInt(t, p, "b")
				c := calcFactInt(t, p, "c")
				d := calcFactInt(t, p, "d")
				left := a * c
				right := b * (d - c)
				want := []int{left, right, left + right}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "absolute integral",
			gen:  &absoluteIntegralGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				shift := calcFactInt(t, p, "shift")
				want := []int{a * a, a*a/2 + shift*a, a*a/2 + a*shift}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "accumulation",
			gen:  &accumulationGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				b := calcFactInt(t, p, "b")
				solveX := calcFactInt(t, p, "solve_x")
				target := calcFactInt(t, p, "target")
				want := []int{b*b - a*a, 2 * b, solveX}[kp-1]
				if kp == 3 {
					if solveX <= a {
						t.Fatalf("solve_x = %d, want > a = %d", solveX, a)
					}
					if solveX*solveX-a*a != target {
						t.Fatalf("solve_x does not satisfy F(x)=%d", target)
					}
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "motion",
			gen:  &motionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := calcFactInt(t, p, "a")
				b := calcFactInt(t, p, "b")
				c := calcFactInt(t, p, "c")
				velocity := 2*c + b
				displacement := (c*c - a*a) + b*(c-a)
				average := displacement / (c - a)
				want := []int{velocity, displacement, average}[kp-1]
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
