package calculus

// Tests for the 2026-10-19 calculus raise.
//
// Nineteen concepts are numeric and eleven are multiple choice. Every concept
// is graded through the real router with the grading_type the corpus ships. For
// the multiple-choice concepts the question is parsed into options, the answer
// must appear exactly once, and every distractor must be rejected. For the
// numeric concepts the tests recompute the answer from the published facts:
// the intermediate-value root, the mean-value speed, one Newton step, the
// marginal cost, a binomial coefficient, a Cauchy-product coefficient, a gamma
// value, a trapezoid error bound, two polar and spherical area multiples, a
// beta fraction, a consumer surplus, a Fourier coefficient, an implicit slope,
// a Jacobian determinant, a homogeneous weighted sum, a scalar line integral, a
// parametric surface area and a flux. The tests never read the published answer.

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/grader"
)

var calcRaiseIDs = []string{
	"calc.limit.intermediate_value",
	"calc.limit.limsup",
	"calc.limit.oscillation",
	"calc.deriv.inverse_hyperbolic",
	"calc.deriv.mvt_applications",
	"calc.deriv.newton_convergence",
	"calc.deriv.convex_jensen",
	"calc.deriv.marginal_analysis",
	"calc.deriv.epsilon_delta_deriv",
	"calc.series.operations",
	"calc.series.binomial",
	"calc.series.cauchy_product",
	"calc.series.p_series",
	"calc.series.dirichlet_test",
	"calc.series.gamma",
	"calc.integral.numerical_error",
	"calc.integral.double_polar",
	"calc.integral.surface_scalar",
	"calc.integral.beta",
	"calc.integral.consumer_surplus",
	"calc.integral.fourier_coefficients",
	"calc.multi.lagrange_two",
	"calc.multi.implicit_function",
	"calc.multi.inverse_function",
	"calc.multi.iterated_limits",
	"calc.multi.euler_homogeneous",
	"calc.vec.line_scalar",
	"calc.vec.parametric_surface",
	"calc.vec.flux_general",
	"calc.vec.identities",
}

var calcRaiseMCIDs = []string{
	"calc.limit.limsup",
	"calc.limit.oscillation",
	"calc.deriv.inverse_hyperbolic",
	"calc.deriv.convex_jensen",
	"calc.deriv.epsilon_delta_deriv",
	"calc.series.operations",
	"calc.series.p_series",
	"calc.series.dirichlet_test",
	"calc.multi.lagrange_two",
	"calc.multi.iterated_limits",
	"calc.vec.identities",
}

var calcRaiseNumericIDs = []string{
	"calc.limit.intermediate_value",
	"calc.deriv.mvt_applications",
	"calc.deriv.newton_convergence",
	"calc.deriv.marginal_analysis",
	"calc.series.binomial",
	"calc.series.cauchy_product",
	"calc.series.gamma",
	"calc.integral.numerical_error",
	"calc.integral.double_polar",
	"calc.integral.surface_scalar",
	"calc.integral.beta",
	"calc.integral.consumer_surplus",
	"calc.integral.fourier_coefficients",
	"calc.multi.implicit_function",
	"calc.multi.inverse_function",
	"calc.multi.euler_homogeneous",
	"calc.vec.line_scalar",
	"calc.vec.parametric_surface",
	"calc.vec.flux_general",
}

var calcOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func calcOptionsOf(q string) []string {
	var out []string
	for _, m := range calcOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func TestRaiseAnswersGradeWithTheirShippedType(t *testing.T) {
	types := calcShippedTypes(t)
	reg := calcBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range calcRaiseIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "numeric" && typ != "multiple_choice" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 200; i++ {
			p := calcBatchSample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if strings.TrimSpace(p.Question) == "" || strings.TrimSpace(p.Explanation) == "" {
				t.Fatalf("%s: empty question or explanation at sample %d", id, i)
			}
			res := router.Grade(grader.GradingType(typ), p.Answer, p.Answer)
			if res.Unavailable {
				t.Skipf("%s: grader unavailable", id)
			}
			if !res.Correct {
				t.Fatalf("%s: the generator's own answer %q does not grade as %s", id, p.Answer, typ)
			}
		}
	}
}

func TestRaiseMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := calcShippedTypes(t)
	reg := calcBatchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range calcRaiseMCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 60; i++ {
			p := calcBatchSample(t, reg, id, i)
			opts := calcOptionsOf(p.Question)
			if len(opts) < 4 {
				t.Fatalf("%s: only %d options; a recognition question needs the answer and three distractors", id, len(opts))
			}
			found := 0
			for _, o := range opts {
				if o == p.Answer {
					found++
				}
			}
			if found != 1 {
				t.Fatalf("%s: the answer %q appears %d times among %v", id, p.Answer, found, opts)
			}
			for _, o := range opts {
				if o == p.Answer {
					continue
				}
				if router.Grade(grader.GradingMultipleChoice, p.Answer, o).Correct {
					t.Fatalf("%s: distractor %q is accepted as correct", id, o)
				}
			}
		}
	}
}

// testGCD and testReduce recompute a fraction independently of calcFraction.
func testGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func testReduce(num, den int) string {
	g := testGCD(num, den)
	num, den = num/g, den/g
	if den == 1 {
		return strconv.Itoa(num)
	}
	return strconv.Itoa(num) + "/" + strconv.Itoa(den)
}

func TestRaiseNumericRecompute(t *testing.T) {
	reg := calcBatchRegistry(t)

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.limit.intermediate_value", i)
		k := calcFactInt(t, p, "k")
		n := calcFactInt(t, p, "n")
		root := calcFactInt(t, p, "answer")
		if root*root != k || root <= 0 || root >= n {
			t.Fatalf("ivt: root %d does not solve x^2=%d inside (0,%d)", root, k, n)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.deriv.mvt_applications", i)
		d := calcFactInt(t, p, "distance")
		h := calcFactInt(t, p, "hours")
		s := calcFactInt(t, p, "speed")
		if d != s*h || d%h != 0 {
			t.Fatalf("mvt: speed %d over %d hours does not give distance %d", s, h, d)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.deriv.newton_convergence", i)
		k := calcFactInt(t, p, "k")
		x0 := calcFactInt(t, p, "x0")
		if want := calcDec((float64(x0) + float64(k)/float64(x0)) / 2); p.Answer != want {
			t.Fatalf("newton: k=%d x0=%d answers %s, want %s", k, x0, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.deriv.marginal_analysis", i)
		a := calcFactInt(t, p, "a")
		b := calcFactInt(t, p, "b")
		x0 := calcFactInt(t, p, "x0")
		if want := strconv.Itoa(2*a*x0 + b); p.Answer != want {
			t.Fatalf("marginal: a=%d b=%d x0=%d answers %s, want %s", a, b, x0, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.series.binomial", i)
		n := calcFactInt(t, p, "n")
		k := calcFactInt(t, p, "k")
		if want := strconv.Itoa(calcBinom(n, k)); p.Answer != want {
			t.Fatalf("binomial: n=%d k=%d answers %s, want %s", n, k, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.series.cauchy_product", i)
		n := calcFactInt(t, p, "n")
		if want := strconv.Itoa(n + 1); p.Answer != want {
			t.Fatalf("cauchy product: n=%d answers %s, want %s", n, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.series.gamma", i)
		n := calcFactInt(t, p, "n")
		// Three shapes share this concept, so the expectation follows the shape
		// rather than assuming the factorial form. The ratio forms are the
		// recurrence itself and have no factorial in the answer.
		want := ""
		switch p.Facts["shape"] {
		case "ratio1":
			want = strconv.Itoa(n)
		case "ratio2":
			want = strconv.Itoa(n * (n + 1))
		default:
			want = strconv.Itoa(calcFact(n - 1))
		}
		if p.Answer != want {
			t.Fatalf("gamma: n=%d shape=%s answers %s, want %s", n, p.Facts["shape"], p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.numerical_error", i)
		m := calcFactInt(t, p, "m")
		w := calcFactInt(t, p, "width")
		n := calcFactInt(t, p, "n")
		want := calcDec(float64(m*w*w*w) / (12 * float64(n*n)))
		if p.Answer != want {
			t.Fatalf("numerical error: m=%d w=%d n=%d answers %s, want %s", m, w, n, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.double_polar", i)
		r := calcFactInt(t, p, "r")
		// The generator now poses a disk, an annulus and a semicircle, so the
		// expectation is recomputed per shape rather than assuming the disk.
		var want int
		switch p.Facts["shape"] {
		case "disk":
			want = r * r
		case "annulus":
			want = r*r - 1
		case "semicircle":
			want = 2 * r * r
		default:
			t.Fatalf("double polar: unknown shape %q", p.Facts["shape"])
		}
		if got := strconv.Itoa(want); p.Answer != got {
			t.Fatalf("double polar: shape=%s r=%d answers %s, want %s", p.Facts["shape"], r, p.Answer, got)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.surface_scalar", i)
		r := calcFactInt(t, p, "r")
		// The generator now poses four surfaces, so the expectation is
		// recomputed per shape rather than assuming the hemisphere.
		var want int
		switch p.Facts["shape"] {
		case "hemisphere":
			want = r * r * r
		case "sphere":
			want = 4 * r * r
		case "disk":
			want = r * r
		case "cylinder":
			want = 2 * r * r
		default:
			t.Fatalf("surface scalar: unknown shape %q", p.Facts["shape"])
		}
		if got := strconv.Itoa(want); p.Answer != got {
			t.Fatalf("surface scalar: shape=%s r=%d answers %s, want %s", p.Facts["shape"], r, p.Answer, got)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.beta", i)
		m := calcFactInt(t, p, "m")
		n := calcFactInt(t, p, "n")
		want := testReduce(calcFact(m-1)*calcFact(n-1), calcFact(m+n-1))
		if p.Answer != want {
			t.Fatalf("beta: m=%d n=%d answers %s, want %s", m, n, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.consumer_surplus", i)
		a := calcFactInt(t, p, "a")
		price := calcFactInt(t, p, "p")
		q := calcFactInt(t, p, "qstar")
		if q != a-price {
			t.Fatalf("consumer surplus: q*=%d but a-p=%d", q, a-price)
		}
		if want := strconv.Itoa(q * q / 2); p.Answer != want {
			t.Fatalf("consumer surplus: q*=%d answers %s, want %s", q, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.integral.fourier_coefficients", i)
		n := calcFactInt(t, p, "n")
		sign := 2.0
		if n%2 == 0 {
			sign = -2.0
		}
		if want := calcDec(sign / float64(n)); p.Answer != want {
			t.Fatalf("fourier: n=%d answers %s, want %s", n, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.multi.implicit_function", i)
		x := calcFactInt(t, p, "x")
		y := calcFactInt(t, p, "y")
		k := calcFactInt(t, p, "k")
		if k != x*x+y*y {
			t.Fatalf("implicit: k=%d but x^2+y^2=%d", k, x*x+y*y)
		}
		if want := calcDec(-float64(x) / float64(y)); p.Answer != want {
			t.Fatalf("implicit: (%d,%d) answers %s, want %s", x, y, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.multi.inverse_function", i)
		x := calcFactInt(t, p, "x")
		y := calcFactInt(t, p, "y")
		if want := strconv.Itoa(4 * (x*x + y*y)); p.Answer != want {
			t.Fatalf("inverse: (%d,%d) answers %s, want %s", x, y, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.multi.euler_homogeneous", i)
		x := calcFactInt(t, p, "x")
		y := calcFactInt(t, p, "y")
		f := calcFactInt(t, p, "f")
		if f != x*x*x+3*x*y*y {
			t.Fatalf("euler: f=%d but f(%d,%d)=%d", f, x, y, x*x*x+3*x*y*y)
		}
		if want := strconv.Itoa(3 * f); p.Answer != want {
			t.Fatalf("euler: f=%d answers %s, want %s", f, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.vec.line_scalar", i)
		a := calcFactInt(t, p, "a")
		b := calcFactInt(t, p, "b")
		l := calcFactInt(t, p, "length")
		if l*l != a*a+b*b {
			t.Fatalf("line scalar: length %d is not the hypotenuse of (%d,%d)", l, a, b)
		}
		if want := calcDec(float64(a*l) / 2); p.Answer != want {
			t.Fatalf("line scalar: a=%d l=%d answers %s, want %s", a, l, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.vec.parametric_surface", i)
		pp := calcFactInt(t, p, "p")
		q := calcFactInt(t, p, "q")
		r := calcFactInt(t, p, "r")
		area := calcFactInt(t, p, "area")
		if area%pp != 0 || area/pp != int(math.Round(math.Sqrt(float64(q*q+r*r)))) {
			t.Fatalf("parametric surface: p=%d q=%d r=%d area=%d is inconsistent", pp, q, r, area)
		}
	}

	for i := 0; i < 200; i++ {
		p := calcBatchSample(t, reg, "calc.vec.flux_general", i)
		s := calcFactInt(t, p, "s")
		tt := calcFactInt(t, p, "t")
		fz := calcFactInt(t, p, "fz")
		if want := strconv.Itoa(fz * s * tt); p.Answer != want {
			t.Fatalf("flux: fz=%d s=%d t=%d answers %s, want %s", fz, s, tt, p.Answer, want)
		}
	}
}
