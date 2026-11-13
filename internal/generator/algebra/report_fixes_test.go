package algebra

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func TestOrdinalSuffixes(t *testing.T) {
	cases := map[int]string{
		1: "1st", 2: "2nd", 3: "3rd", 4: "4th", 5: "5th",
		11: "11th", 12: "12th", 13: "13th",
		21: "21st", 22: "22nd", 23: "23rd",
	}
	for n, want := range cases {
		if got := ordinal(n); got != want {
			t.Errorf("ordinal(%d) = %q, want %q", n, got, want)
		}
	}
}

// Report 453: the question printed the first three terms and asked for term 3,
// so the answer was already on screen. Ask for a term the learner must compute.
func TestSeqGeomAsksBeyondPrintedTerms(t *testing.T) {
	g := &seqGeomGen{}
	re := regexp.MustCompile(`Find term (\d+)`)
	for i := 0; i < 300; i++ {
		p := g.Generate(generator.GeneratorContext{})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected sequence question: %q", p.Question)
		}
		n, _ := strconv.Atoi(m[1])
		if n < 4 {
			t.Fatalf("term %d is already printed in %q", n, p.Question)
		}
		a1, _ := strconv.Atoi(p.Facts["a1"])
		r, _ := strconv.Atoi(p.Facts["r"])
		pow := 1
		for j := 1; j < n; j++ {
			pow *= r
		}
		if p.Answer != strconv.Itoa(a1*pow) {
			t.Fatalf("%q answered %s, want %d", p.Question, p.Answer, a1*pow)
		}
	}
}

// The printed equation and the stored answer must agree. An integer division
// of a*x by b used to truncate, so the answer solved a different equation from
// the one the learner could see.
func TestEqRationalPrintedEquationMatchesAnswer(t *testing.T) {
	g := &eqRationalGen{}
	re := regexp.MustCompile(`^Solve: \\\((\d+)x/(\d+) = (\d+)\\\)$`)
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected question shape: %q", p.Question)
		}
		a, _ := strconv.Atoi(m[1])
		b, _ := strconv.Atoi(m[2])
		rhs, _ := strconv.Atoi(m[3])
		x, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("answer is not an integer: %q", p.Answer)
		}
		if a*x != rhs*b {
			t.Fatalf("answer %d does not solve %q", x, p.Question)
		}
		if x < 2 {
			t.Fatalf("answer %d is below the intended range for %q", x, p.Question)
		}
	}
}
