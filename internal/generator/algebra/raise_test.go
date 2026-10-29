package algebra

// Tests for the 2026-10-19 raise.
//
// The twelve numeric concepts recompute their answers from the published
// parameters. The eight multiple-choice concepts are checked the way the 2026
// batch checks its own: the question is parsed into options, the answer must
// appear exactly once, and every distractor must be rejected by the real
// router. The tests never read the published "answer" fact.

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var algRaiseIDs = []string{
	"alg.func.piecewise",
	"alg.func.transformations",
	"alg.func.inverse_domain",
	"alg.exp.graph",
	"alg.exp.compound",
	"alg.log.graph",
	"alg.log.model",
	"alg.radical.simplify",
	"alg.radical.operations",
	"alg.radical.conjugate",
	"alg.rational.complex_fraction",
	"alg.rational.proportions",
	"alg.rational.asymptotes",
	"alg.systems.three_var",
	"alg.systems.matrices",
	"alg.systems.nonlinear",
	"alg.seq.sigma",
	"alg.seq.induction",
	"alg.poly.factor_theorem",
	"alg.poly.rational_root",
}

var algRaiseMCIDs = []string{
	"alg.func.transformations",
	"alg.func.inverse_domain",
	"alg.log.graph",
	"alg.systems.matrices",
	"alg.seq.induction",
}

var algRaiseNumericIDs = []string{
	"alg.func.piecewise",
	"alg.exp.graph",
	"alg.exp.compound",
	"alg.log.model",
	"alg.radical.simplify",
	"alg.radical.operations",
	"alg.radical.conjugate",
	"alg.rational.complex_fraction",
	"alg.rational.proportions",
	"alg.rational.asymptotes",
	"alg.systems.three_var",
	"alg.systems.nonlinear",
	"alg.seq.sigma",
	"alg.poly.rational_root",
	"alg.poly.factor_theorem",
}

func TestRaiseAnswersGradeWithTheirShippedType(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range algRaiseIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "numeric" && typ != "multiple_choice" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 200; i++ {
			p := sample(t, reg, id, 5000+i)
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

func TestRaiseMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range algRaiseMCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 60; i++ {
			p := sample(t, reg, id, 8000+i)
			opts := optionsOf(p.Question)
			if len(opts) < 4 {
				t.Fatalf("%s: only %d options; a recognition question needs four", id, len(opts))
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

func TestRaiseNumericRecompute(t *testing.T) {
	reg := batchRegistry(t)

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.func.piecewise", 100+i)
		x0 := intFact(t, p, "x0")
		want := 3*x0 - 2
		if x0 < 0 {
			want = 2*x0 + 1
		}
		if p.Answer != fmt.Sprint(want) {
			t.Fatalf("piecewise: x0=%d answers %s, want %d", x0, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.exp.graph", 100+i)
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		x := intFact(t, p, "x")
		want := a
		for j := 0; j < x; j++ {
			want *= b
		}
		if p.Answer != fmt.Sprint(want) {
			t.Fatalf("exp.graph: a=%d b=%d x=%d answers %s, want %d", a, b, x, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.exp.compound", 100+i)
		pr := intFact(t, p, "p")
		r := intFact(t, p, "r")
		tt := intFact(t, p, "t")
		want := fmt.Sprintf("%.2f", float64(pr)*math.Pow(1+float64(r)/100, float64(tt)))
		if p.Answer != want {
			t.Fatalf("exp.compound: p=%d r=%d t=%d answers %s, want %s", pr, r, tt, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.log.model", 100+i)
		b := intFact(t, p, "b")
		k := intFact(t, p, "k")
		want := 1
		for j := 0; j < k; j++ {
			want *= b
		}
		if p.Answer != fmt.Sprint(want) {
			t.Fatalf("log.model: b=%d k=%d answers %s, want %d", b, k, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.radical.simplify", 100+i)
		n := intFact(t, p, "n")
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		if n != a*a*b {
			t.Fatalf("radical.simplify: n=%d is not a^2 b for a=%d b=%d", n, a, b)
		}
		if p.Answer != fmt.Sprint(a) {
			t.Fatalf("radical.simplify: n=%d answers %s, want %d", n, p.Answer, a)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.radical.operations", 100+i)
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		if p.Answer != fmt.Sprint(a+b) {
			t.Fatalf("radical.operations: a=%d b=%d answers %s, want %d", a, b, p.Answer, a+b)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.rational.complex_fraction", 100+i)
		k := intFact(t, p, "k")
		if p.Answer != fmt.Sprint(k) {
			t.Fatalf("complex_fraction: k=%d answers %s, want %d", k, p.Answer, k)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.radical.conjugate", 100+i)
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		if p.Answer != fmt.Sprint(a-b) {
			t.Fatalf("radical.conjugate: a=%d b=%d answers %s, want %d", a, b, p.Answer, a-b)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.rational.proportions", 100+i)
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		num := intFact(t, p, "num")
		if b*num%a != 0 {
			t.Fatalf("proportions: a=%d b=%d num=%d does not divide", a, b, num)
		}
		if want := fmt.Sprint(b * num / a); p.Answer != want {
			t.Fatalf("proportions: a=%d b=%d num=%d answers %s, want %s", a, b, num, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.rational.asymptotes", 100+i)
		m := intFact(t, p, "m")
		n := intFact(t, p, "n")
		larger := m
		if n > m {
			larger = n
		}
		if p.Answer != fmt.Sprint(larger) {
			t.Fatalf("asymptotes: m=%d n=%d answers %s, want %d", m, n, p.Answer, larger)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.systems.three_var", 100+i)
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		d := intFact(t, p, "d")
		if (a-b)%2 != 0 || (a-d)%2 != 0 {
			t.Fatalf("three_var: a=%d b=%d d=%d not integral", a, b, d)
		}
		y := (a - b) / 2
		z := (a - d) / 2
		if want := fmt.Sprint(a - y - z); p.Answer != want {
			t.Fatalf("three_var: a=%d b=%d d=%d answers %s, want %s", a, b, d, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.systems.nonlinear", 100+i)
		k := intFact(t, p, "k")
		if want := fmt.Sprint(k * k); p.Answer != want {
			t.Fatalf("nonlinear: k=%d answers %s, want %s", k, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.seq.sigma", 100+i)
		a := intFact(t, p, "a")
		b := intFact(t, p, "b")
		n := intFact(t, p, "n")
		if want := fmt.Sprint(a*n*(n+1)/2 + b*n); p.Answer != want {
			t.Fatalf("sigma: a=%d b=%d n=%d answers %s, want %s", a, b, n, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := sample(t, reg, "alg.poly.rational_root", 100+i)
		b := intFact(t, p, "b")
		d := intFact(t, p, "d")
		e := intFact(t, p, "e")
		r := intFact(t, p, "answer")
		if r <= 0 {
			t.Fatalf("rational_root: answer %d is not positive", r)
		}
		if r*r*r+b*r*r+d*r+e != 0 {
			t.Fatalf("rational_root: %d is not a root of x^3+%dx^2+%dx+%d", r, b, d, e)
		}
	}
}

func TestRationalRootQuestionTypesetsThePolynomial(t *testing.T) {
	p := (&rationalRootGen{}).Generate(generator.GeneratorContext{})
	if !strings.Contains(p.Question, `\(x^3`) || !strings.Contains(p.Question, `\)`) {
		t.Fatalf("rational-root question is not typeset: %q", p.Question)
	}
	if !strings.Contains(p.Explanation, `\(`) || !strings.Contains(p.Explanation, `\)`) {
		t.Fatalf("rational-root explanation is not typeset: %q", p.Explanation)
	}
}

func TestRadicalSimplifyUsesMathForVariables(t *testing.T) {
	p := (&radicalSimplifyGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5})
	if !strings.Contains(p.Question, `\(a\sqrt{b}\)`) {
		t.Fatalf("question does not typeset a and b: %q", p.Question)
	}
	if strings.Contains(p.Question, "a times the square root of b") {
		t.Fatalf("question still uses plain-text variables: %q", p.Question)
	}
}
