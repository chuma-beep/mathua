package algebra

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

// Report 494: every draw had the drain faster than the fill, so the net rate
// was negative and the tank never filled. The learner's -2/15 was right.
func TestRationalWordDrainIsSlowerThanFill(t *testing.T) {
	g := &algBridgeGen{kind: "rational_word"}
	re := regexp.MustCompile(`fills a tank in (\d+) hours and a drain empties it in (\d+) hours`)
	for i := 0; i < 500; i++ {
		p := g.Generate(generator.GeneratorContext{KPOrdinal: 3})
		m := re.FindStringSubmatch(p.Question)
		if m == nil {
			t.Fatalf("unexpected tank question: %q", p.Question)
		}
		fill, _ := strconv.Atoi(m[1])
		drain, _ := strconv.Atoi(m[2])
		if drain <= fill {
			t.Fatalf("drain %d h is not slower than fill %d h in %q", drain, fill, p.Question)
		}
		if fill*drain%(drain-fill) != 0 {
			t.Fatalf("%d and %d do not give a whole number of hours", fill, drain)
		}
		want := fill * drain / (drain - fill)
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("%q answered %s, want %d", p.Question, p.Answer, want)
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

// Report 500: the learner factored the trinomial correctly but separated the
// two binomials with a comma and spaces, and the exact-string grader rejected
// it. Order is irrelevant too: multiplication commutes.
func TestFactorTrinomialGradeIgnoresSeparatorsAndOrder(t *testing.T) {
	g := &factorTrinomialGen{}
	for i := 0; i < 40; i++ {
		p := g.Generate(generator.GeneratorContext{})
		if !strings.Contains(p.Question, "product of two binomials") {
			t.Fatalf("question does not name the expected form: %q", p.Question)
		}
		compact := strings.ReplaceAll(p.Answer, " ", "")
		if compact == p.Answer {
			t.Fatalf("answer has no spaces to vary: %q", p.Answer)
		}
		comma := strings.Replace(compact, ")(", ") , (", 1)
		swapped := swapFactorOrder(compact)
		for _, answer := range []string{compact, comma, swapped, "  " + comma + "  "} {
			if res := g.Grade(p.Answer, answer); !res.Correct {
				t.Fatalf("Grade(%q, %q) rejected a correct factorisation: %+v", p.Answer, answer, res)
			}
		}
		r1, _ := strconv.Atoi(p.Facts["r1"])
		wrong := fmt.Sprintf("(x + %d)(x + %d)", r1+1, r1+1)
		if res := g.Grade(p.Answer, wrong); res.Correct {
			t.Fatalf("Grade accepted %q for %q", wrong, p.Question)
		}
	}
}
