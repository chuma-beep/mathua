package numerical

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave2FactInts(t *testing.T, p generator.Problem, key string) []int {
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
			name: "aitken",
			gen:  &aitkenGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				x0 := numaFactInt(t, p, "x0")
				x1 := numaFactInt(t, p, "x1")
				x2 := numaFactInt(t, p, "x2")
				limit := numaFactInt(t, p, "limit")
				delta := x2 - 2*x1 + x0
				if limit != x0-(x1-x0)*(x1-x0)/delta {
					t.Fatalf("limit fact = %d, Aitken recomputation disagrees", limit)
				}
				want := []int{limit, x0 - limit, delta}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "muller",
			gen:  &mullerGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := numaFactInt(t, p, "a")
				b := numaFactInt(t, p, "b")
				c := numaFactInt(t, p, "c")
				want := []int{a, 4*a + 2*b + c, b*b - 4*a*c}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "barycentric",
			gen:  &barycentricGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				n := numaFactInt(t, p, "n")
				k := numaFactInt(t, p, "k")
				value := numaFactInt(t, p, "value")
				want := []int{factorial(k) * factorial(n-k), value, n + 1}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "runge",
			gen:  &rungeGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				nodes := numaFactInt(t, p, "nodes")
				intervals := numaFactInt(t, p, "intervals")
				if nodes != intervals+1 {
					t.Fatalf("nodes = %d, want intervals+1", nodes)
				}
				want := []int{nodes, nodes - 2, 1}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "piecewise linear",
			gen:  &piecewiseLinearGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				x0 := numaFactInt(t, p, "x0")
				x1 := numaFactInt(t, p, "x1")
				y0 := numaFactInt(t, p, "y0")
				y1 := numaFactInt(t, p, "y1")
				if (y1-y0)%(x1-x0) != 0 {
					t.Fatalf("slope is not an integer")
				}
				want := []int{(y0 + y1) / 2, (y1 - y0) / (x1 - x0), x0}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "boole",
			gen:  &booleGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				weighted := numaFactInt(t, p, "weighted")
				panels := numaFactInt(t, p, "panels")
				want := []int{weighted, 2 * weighted, 4*panels + 1}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "gauss kronrod",
			gen:  &gaussKronrodGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				nodes := numaFactInt(t, p, "nodes")
				gauss := numaFactInt(t, p, "gauss_nodes")
				panels := numaFactInt(t, p, "panels")
				want := nodes
				switch kp {
				case 2:
					want = gauss
				case 3:
					want = nodes * panels
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "gauss jordan",
			gen:  &gaussJordanGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				a := numaFactInt(t, p, "a")
				b := numaFactInt(t, p, "b")
				c := numaFactInt(t, p, "c")
				d := numaFactInt(t, p, "d")
				e := numaFactInt(t, p, "e")
				f := numaFactInt(t, p, "f")
				x := numaFactInt(t, p, "x")
				y := numaFactInt(t, p, "y")
				det := numaFactInt(t, p, "det")
				if a*x+b*y != e || c*x+d*y != f {
					t.Fatalf("x=%d y=%d do not solve the system", x, y)
				}
				if det != a*d-b*c {
					t.Fatalf("determinant fact = %d, want %d", det, a*d-b*c)
				}
				want := []int{x, y, det}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "thomas",
			gen:  &thomasGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				diag := numaFactInt(t, p, "diag")
				off := numaFactInt(t, p, "off")
				x1 := numaFactInt(t, p, "x1")
				x2 := numaFactInt(t, p, "x2")
				x3 := numaFactInt(t, p, "x3")
				if diag*x1+off*x2 != numaFactInt(t, p, "d1") {
					t.Fatal("first equation does not hold")
				}
				if off*x1+diag*x2+off*x3 != numaFactInt(t, p, "d2") {
					t.Fatal("second equation does not hold")
				}
				if off*x2+diag*x3 != numaFactInt(t, p, "d3") {
					t.Fatal("third equation does not hold")
				}
				want := []int{x1, x2, x3}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "matrix condition",
			gen:  &matrixConditionGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				condition := numaFactInt(t, p, "condition")
				exponent := numaFactInt(t, p, "exponent")
				want := condition
				if kp == 3 {
					want = exponent
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "inverse power",
			gen:  &inversePowerGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				smallest := numaFactInt(t, p, "smallest")
				middle := numaFactInt(t, p, "middle")
				largest := numaFactInt(t, p, "largest")
				shift := numaFactInt(t, p, "shift")
				closest := numaFactInt(t, p, "closest")
				best := 1 << 30
				for _, v := range []int{smallest, middle, largest} {
					if d := absInt(v - shift); d < best {
						best = d
						closest = v
					} else if d == best {
						t.Fatalf("shift %d is equidistant from two eigenvalues", shift)
					}
				}
				want := []int{smallest, closest, largest}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "dft",
			gen:  &dftGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				n := numaFactInt(t, p, "n")
				sum := numaFactInt(t, p, "sum")
				alt := numaFactInt(t, p, "alternating")
				want := []int{sum, alt, n * n}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "aliasing",
			gen:  &aliasingGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				fs := numaFactInt(t, p, "fs")
				f := numaFactInt(t, p, "f")
				alias := numaFactInt(t, p, "alias")
				if f-fs != alias {
					t.Fatalf("alias = %d, want f-fs = %d", alias, f-fs)
				}
				want := []int{fs / 2, alias, 2 * f}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "pade",
			gen:  &padeGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				m := numaFactInt(t, p, "m")
				n := numaFactInt(t, p, "n")
				x := numaFactInt(t, p, "x")
				want := []int{m + n + 1, m + n, 1 + x/2}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "finite differences",
			gen:  &numDiffGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				x := numaFactInt(t, p, "x")
				h := numaFactInt(t, p, "h")
				want := []int{2*x + h, 2 * x, 2 * h * h}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "richardson",
			gen:  &richardsonGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				tExt := numaFactInt(t, p, "extrapolated")
				d1 := numaFactInt(t, p, "d1")
				d2 := numaFactInt(t, p, "d2")
				exponent := numaFactInt(t, p, "p")
				if (4*d2-d1)/3 != tExt {
					t.Fatalf("extrapolated fact = %d, fourth-order formula disagrees", tExt)
				}
				want := []int{tExt, d2 - d1, 1 << exponent}[kp-1]
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("answer = %s, want %d", p.Answer, want)
				}
			},
		},
		{
			name: "line search",
			gen:  &lineSearchGen{},
			check: func(t *testing.T, p generator.Problem, kp int) {
				center := numaFactInt(t, p, "center")
				x0 := numaFactInt(t, p, "x0")
				alpha := numaFactInt(t, p, "alpha")
				if x0+alpha != center && x0-alpha != center {
					t.Fatalf("alpha = %d does not reach center %d from x0 = %d", alpha, center, x0)
				}
				want := []int{alpha, center, 2}[kp-1]
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
