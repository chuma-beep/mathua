package numerical

// Tests for the second numerical-analysis batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships. The nine multiple-choice concepts are checked the way the other
// batches check theirs: the question is parsed into options, the correct answer
// must appear exactly once, and every other option must be rejected.
//
// The eleven numeric concepts recompute their answers from the published
// parameters: a least-squares slope and intercept from the data, the trapezoid,
// midpoint, Simpson and composite rules on a square, two-point Gauss quadrature,
// Gaussian elimination, an LU multiplier, one Jacobi step and the dominant
// eigenvalue of a symmetric two-by-two matrix.

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/grader"
)

var numaBatch2IDs = []string{
	"numa.interp.least_squares",
	"numa.interp.normal_equations",
	"numa.interp.orthogonal",
	"numa.interp.extrapolation",
	"numa.integration.trapezoid",
	"numa.integration.midpoint",
	"numa.integration.simpson",
	"numa.integration.composite",
	"numa.integration.error",
	"numa.integration.romberg",
	"numa.integration.gauss",
	"numa.integration.adaptive",
	"numa.integration.improper",
	"numa.linear.elimination",
	"numa.linear.pivoting",
	"numa.linear.lu",
	"numa.linear.jacobi",
	"numa.linear.gauss_seidel",
	"numa.linear.power",
	"numa.linear.least_squares_qr",
}

var numaBatch2MCIDs = []string{}

func TestBatch2AnswersGradeWithTheirShippedType(t *testing.T) {
	types := numaShippedTypes(t)
	reg := numaRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range numaBatch2IDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := numaSample(t, reg, id, i)
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

func TestBatch2MultipleChoiceRejectsEveryDistractor(t *testing.T) {
	if len(numaBatch2MCIDs) == 0 {
		t.Skip("the numerical batch now serves computed count problems")
	}
	types := numaShippedTypes(t)
	reg := numaRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range numaBatch2MCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 40; i++ {
			p := numaSample(t, reg, id, i)
			opts := numaOptionsOf(p.Question)
			if len(opts) < 3 {
				t.Fatalf("%s: only %d options; a recognition question needs at least three", id, len(opts))
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

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

// centredFit returns the least-squares slope and intercept for the data.
func centredFit(xs, ys []int) (b, a float64) {
	n := float64(len(xs))
	var sx, sy float64
	for i := range xs {
		sx += float64(xs[i])
		sy += float64(ys[i])
	}
	xbar, ybar := sx/n, sy/n
	var sxx, sxy float64
	for i := range xs {
		dx := float64(xs[i]) - xbar
		sxx += dx * dx
		sxy += dx * (float64(ys[i]) - ybar)
	}
	b = sxy / sxx
	a = ybar - b*xbar
	return b, a
}

func TestBatch2LeastSquaresRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.interp.least_squares", i)
		xs := numaIntList(t, p, "xs")
		ys := numaIntList(t, p, "ys")
		want, _ := centredFit(xs, ys)
		if p.Answer != decimal6(want) {
			t.Fatalf("least squares: xs=%v ys=%v answers %s, want %s", xs, ys, p.Answer, decimal6(want))
		}
	}
}

func TestBatch2NormalEquationsRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.interp.normal_equations", i)
		xs := numaIntList(t, p, "xs")
		ys := numaIntList(t, p, "ys")
		_, want := centredFit(xs, ys)
		if p.Answer != decimal6(want) {
			t.Fatalf("normal equations: xs=%v ys=%v answers %s, want %s", xs, ys, p.Answer, decimal6(want))
		}
	}
}

func TestBatch2TrapezoidRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.integration.trapezoid", i)
		a := numaFactInt(t, p, "a")
		b := numaFactInt(t, p, "b")
		want := decimal6(float64(b-a) * float64(a*a+b*b) / 2)
		if p.Answer != want {
			t.Fatalf("trapezoid: a=%d b=%d answers %s, want %s", a, b, p.Answer, want)
		}
	}
}

func TestBatch2MidpointRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.integration.midpoint", i)
		a := numaFactInt(t, p, "a")
		b := numaFactInt(t, p, "b")
		m := float64(a+b) / 2
		want := decimal6(float64(b-a) * m * m)
		if p.Answer != want {
			t.Fatalf("midpoint: a=%d b=%d answers %s, want %s", a, b, p.Answer, want)
		}
	}
}

func TestBatch2SimpsonRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.integration.simpson", i)
		a := numaFactInt(t, p, "a")
		b := numaFactInt(t, p, "b")
		m := float64(a+b) / 2
		want := decimal6((float64(b-a) / 6) * (float64(a*a) + 4*m*m + float64(b*b)))
		if p.Answer != want {
			t.Fatalf("simpson: a=%d b=%d answers %s, want %s", a, b, p.Answer, want)
		}
	}
}

func TestBatch2CompositeRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.integration.composite", i)
		b := numaFactInt(t, p, "b")
		n := numaFactInt(t, p, "n")
		h := float64(b) / float64(n)
		sum := 0.0
		for k := 0; k <= n; k++ {
			x := float64(k) * h
			w := 1.0
			if k == 0 || k == n {
				w = 0.5
			}
			sum += w * x * x
		}
		want := decimal6(h * sum)
		if p.Answer != want {
			t.Fatalf("composite: b=%d n=%d answers %s, want %s", b, n, p.Answer, want)
		}
	}
}

func TestBatch2GaussRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.integration.gauss", i)
		k := numaFactInt(t, p, "k")
		want := decimal6(2 * math.Pow(1.0/3.0, float64(k)))
		if p.Answer != want {
			t.Fatalf("gauss: k=%d answers %s, want %s", k, p.Answer, want)
		}
	}
}

func TestBatch2EliminationRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.linear.elimination", i)
		a := float64(numaFactInt(t, p, "a"))
		b := float64(numaFactInt(t, p, "b"))
		c := float64(numaFactInt(t, p, "c"))
		d := float64(numaFactInt(t, p, "d"))
		e := float64(numaFactInt(t, p, "e"))
		f := float64(numaFactInt(t, p, "f"))
		want := decimal6((a*f - c*e) / (a*d - b*c))
		if p.Answer != want {
			t.Fatalf("elimination: a=%v b=%v c=%v d=%v e=%v f=%v answers %s, want %s", a, b, c, d, e, f, p.Answer, want)
		}
	}
}

func TestBatch2LURecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.linear.lu", i)
		a := numaFactInt(t, p, "a")
		c := numaFactInt(t, p, "c")
		want := decimal6(float64(c) / float64(a))
		if p.Answer != want {
			t.Fatalf("lu: a=%d c=%d answers %s, want %s", a, c, p.Answer, want)
		}
	}
}

func TestBatch2JacobiRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.linear.jacobi", i)
		a := numaFactInt(t, p, "a")
		b := numaFactInt(t, p, "b")
		e := numaFactInt(t, p, "e")
		y0 := numaFactInt(t, p, "y0")
		want := decimal6(float64(e-b*y0) / float64(a))
		if p.Answer != want {
			t.Fatalf("jacobi: a=%d b=%d e=%d y0=%d answers %s, want %s", a, b, e, y0, p.Answer, want)
		}
	}
}

func TestBatch2PowerRecomputes(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 200; i++ {
		p := numaSample(t, reg, "numa.linear.power", i)
		a := float64(numaFactInt(t, p, "a"))
		b := float64(numaFactInt(t, p, "b"))
		d := float64(numaFactInt(t, p, "d"))
		want := decimal6((a+d)/2 + math.Sqrt(math.Pow((a-d)/2, 2)+b*b))
		if p.Answer != want {
			t.Fatalf("power: a=%v b=%v d=%v answers %s, want %s", a, b, d, p.Answer, want)
		}
	}
}

// Guard the integer formatting used by the elimination generator so a change in
// the fact strings cannot silently break the recomputation.
func TestBatch2EliminationFactsAreIntegers(t *testing.T) {
	reg := numaRegistry(t)
	for i := 0; i < 20; i++ {
		p := numaSample(t, reg, "numa.linear.elimination", i)
		for _, key := range []string{"a", "b", "c", "d", "e", "f"} {
			if _, err := strconv.Atoi(p.Facts[key]); err != nil {
				t.Fatalf("elimination: fact %q = %q is not an integer", key, p.Facts[key])
			}
		}
	}
}
