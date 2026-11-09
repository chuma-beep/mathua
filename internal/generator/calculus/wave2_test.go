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
				if p.Answer != strconv.Itoa(a*a) {
					t.Fatalf("answer = %s, want %d", p.Answer, a*a)
				}
			},
		},
		{
			name: "absolute value limits",
			gen:  &absoluteValueLimitGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				want := []int{-1, 1, 0}[kp-1]
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
				want := []int{2 * a, 2 * a, 3 * a * a}[kp-1]
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
				if x*x+y*y != l*l {
					t.Fatalf("%d-%d-%d is not a right triangle", x, y, l)
				}
				want := []int{y, 1, x}[kp-1]
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
				if 4*c != h*h*dh {
					t.Fatalf("c=%d does not satisfy h^2 dh / 4", c)
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
				want := []int{a * a, a * a / 2, a * a / 2}[kp-1]
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
