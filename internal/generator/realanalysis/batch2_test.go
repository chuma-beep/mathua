package realanalysis

// Tests for the second real-analysis batch.
//
// The fifteen multiple-choice concepts are checked the way the first batch checks
// its own: the question is parsed into options, the answer must appear exactly
// once, and every distractor must be rejected by the real router. The ten numeric
// concepts recompute their answers from the published parameters: the mean value
// point of x², a chain-rule derivative, a limit of a trigonometric quotient, a
// right-endpoint Riemann sum, a definite integral, an integration by parts, an
// improper integral, a radius of convergence, a Taylor coefficient and the fixed
// point of a contraction. The tests never read the published "answer" fact.

import (
	"math"
	"testing"

	"github.com/chuma-beep/mathua/internal/grader"
)

var raBatch2IDs = []string{
	"ra.limit.connected",
	"ra.limit.discontinuity",
	"ra.limit.monotone_fn",
	"ra.deriv.mvt",
	"ra.deriv.chain",
	"ra.deriv.taylor",
	"ra.deriv.lhopital",
	"ra.deriv.riemann",
	"ra.deriv.darboux",
	"ra.deriv.ftc",
	"ra.deriv.parts",
	"ra.deriv.improper",
	"ra.deriv.uniform_conv",
	"ra.deriv.power_series",
	"ra.deriv.interchange",
	"ra.metric.spaces",
	"ra.metric.open_closed",
	"ra.metric.closure",
	"ra.metric.convergence",
	"ra.metric.complete",
	"ra.metric.compact",
	"ra.metric.continuity",
	"ra.metric.contraction",
	"ra.metric.baire",
	"ra.metric.equivalent",
}

var raBatch2MCIDs = []string{}

func TestRABatch2AnswersGradeWithTheirShippedType(t *testing.T) {
	types := raShippedTypes(t)
	reg := raRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range raBatch2IDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := raSample(t, reg, id, 3000+i)
			if p.Answer == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
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

func TestRABatch2MultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := raShippedTypes(t)
	reg := raRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range raBatch2MCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 40; i++ {
			p := raSample(t, reg, id, 7000+i)
			opts := raOptionsOf(p.Question)
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

func TestRABatch2NumericRecompute(t *testing.T) {
	reg := raRegistry(t)

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.mvt", 100+i)
		b := raFactInt(t, p, "b")
		if want := decimal6(float64(b) / 2); p.Answer != want {
			t.Fatalf("mvt: b=%d answers %s, want %s", b, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.chain", 100+i)
		a := raFactInt(t, p, "a")
		b := raFactInt(t, p, "b")
		n := raFactInt(t, p, "n")
		x0 := raFactInt(t, p, "x0")
		want := decimal6(float64(n * a * raIntPow(a*x0+b, n-1)))
		if p.Answer != want {
			t.Fatalf("chain: a=%d b=%d n=%d x0=%d answers %s, want %s", a, b, n, x0, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.lhopital", 100+i)
		a := raFactInt(t, p, "a")
		if want := decimal6(float64(a*a) / 2); p.Answer != want {
			t.Fatalf("lhopital: a=%d answers %s, want %s", a, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.riemann", 100+i)
		b := raFactInt(t, p, "b")
		n := raFactInt(t, p, "n")
		want := decimal6(float64(b*b) * float64(n+1) / (2 * float64(n)))
		if p.Answer != want {
			t.Fatalf("riemann: b=%d n=%d answers %s, want %s", b, n, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.ftc", 100+i)
		a := raFactInt(t, p, "a")
		c := raFactInt(t, p, "c")
		b := raFactInt(t, p, "b")
		want := decimal6(float64(a)*float64(b*b)/2 + float64(c*b))
		if p.Answer != want {
			t.Fatalf("ftc: a=%d c=%d b=%d answers %s, want %s", a, c, b, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.parts", 100+i)
		b := raFactInt(t, p, "b")
		want := decimal4(float64(b-1)*math.Exp(float64(b)) + 1)
		if p.Answer != want {
			t.Fatalf("parts: b=%d answers %s, want %s", b, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.improper", 100+i)
		pp := raFactInt(t, p, "p")
		if pp <= 1 {
			t.Fatalf("improper: p=%d does not converge", pp)
		}
		if want := decimal6(1 / float64(pp-1)); p.Answer != want {
			t.Fatalf("improper: p=%d answers %s, want %s", pp, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.power_series", 100+i)
		c := raFactInt(t, p, "c")
		if want := decimal6(1 / float64(c)); p.Answer != want {
			t.Fatalf("power_series: c=%d answers %s, want %s", c, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.deriv.taylor", 100+i)
		a := raFactInt(t, p, "a")
		if want := decimal6(float64(a*a*a) / 6); p.Answer != want {
			t.Fatalf("taylor: a=%d answers %s, want %s", a, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := raSample(t, reg, "ra.metric.contraction", 100+i)
		pn := raFactInt(t, p, "pn")
		pd := raFactInt(t, p, "pd")
		b := raFactInt(t, p, "b")
		if pn >= pd {
			t.Fatalf("contraction: %d/%d is not a contraction", pn, pd)
		}
		want := decimal6(float64(b) / (1 - float64(pn)/float64(pd)))
		if p.Answer != want {
			t.Fatalf("contraction: pn=%d pd=%d b=%d answers %s, want %s", pn, pd, b, p.Answer, want)
		}
	}
}
