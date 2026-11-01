package complex

// Tests for the 2026-10-17 batch. Every concept is graded through the real router with the
// grading_type the corpus ships for it, so a generator and its declared type cannot drift apart.
//
// Seven concepts answer with the `complex` grading type, which needs the SymPy runtime. When it is
// missing those cases skip with that stated: Result.Unavailable means "the grader could not run",
// which the engine does not record as a miss (TestUngradedAnswerCarriesNoEvidence).

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

var batchIDs = []string{
	"complex.basics.parts",
	"complex.basics.powers",
	"complex.basics.reciprocal",
	"complex.basics.geometry",
	"complex.basics.loci",
	"complex.basics.arg_arithmetic",
	"complex.basics.equations",
	"complex.adv.conjugate_roots",
	"complex.adv.trig_exponential",
	"complex.adv.hyperbolic",
	"complex.adv.geometric_series",
	"complex.adv.mobius_fixed",
	"complex.adv.residue_simple",
	"complex.adv.cross_ratio",
	"complex.topology.open_sets",
	"complex.topology.domain",
}

type corpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func shippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/complex_numbers.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []corpusConcept
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatalf("parse corpus: %v", err)
	}
	out := map[string]string{}
	for _, r := range rows {
		if r.ID != "" {
			out[r.ID] = r.GradingType
		}
	}
	return out
}

func batchRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func sample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(11000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var optionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func optionsOf(q string) []string {
	var out []string
	for _, m := range optionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range batchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		for i := 0; i < 120; i++ {
			p := sample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if p.Facts == nil {
				t.Fatalf("%s: no facts published; its schema would fall back", id)
			}
			res := router.Grade(grader.GradingType(typ), p.Answer, p.Answer)
			if res.Unavailable {
				t.Skipf("%s: grader unavailable (SymPy runtime missing)", id)
			}
			if !res.Correct {
				t.Fatalf("%s: the generator's own answer %q does not grade as %s", id, p.Answer, typ)
			}
		}
	}
}

func TestBatchMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range batchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := sample(t, reg, id, i)
			opts := optionsOf(p.Question)
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
			checked++
		}
	}
	if checked == 0 {
		t.Skip("the complex batch now serves computed count problems")
	}
}

func TestBatchComplexFormatsAndNearMisses(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range batchIDs {
		if types[id] != "complex" {
			continue
		}
		for i := 0; i < 30; i++ {
			p := sample(t, reg, id, i)
			// Spacing is formatting, not mathematics: a spaced copy must still be accepted.
			spaced := strings.ReplaceAll(p.Answer, "+", " + ")
			spaced = strings.ReplaceAll(spaced, "-", " - ")
			res := router.Grade(grader.GradingComplex, p.Answer, spaced)
			if res.Unavailable {
				t.Skip("SymPy runtime missing; the complex format cases cannot run here")
			}
			if !res.Correct {
				t.Fatalf("%s: the spaced form %q was rejected for %q", id, spaced, p.Answer)
			}
			// A near miss must be rejected: add 1 to the expression.
			if router.Grade(grader.GradingComplex, p.Answer, "("+p.Answer+")+1").Correct {
				t.Fatalf("%s: adding 1 to %q was accepted", id, p.Answer)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no complex-typed concepts were checked")
	}
}

func TestBatchNumericRejectsNearMisses(t *testing.T) {
	types := shippedTypes(t)
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range batchIDs {
		if types[id] != "numeric" {
			continue
		}
		for i := 0; i < 30; i++ {
			p := sample(t, reg, id, i)
			n, err := strconv.Atoi(p.Answer)
			if err != nil {
				t.Fatalf("%s: answer %q is not an integer", id, p.Answer)
			}
			if router.Grade(grader.GradingNumeric, p.Answer, strconv.Itoa(n+1)).Correct {
				t.Fatalf("%s: %d was accepted for %q", id, n+1, p.Answer)
			}
		}
	}
}

// TestBatchMobiusPicksTheLargerRoot pins the selection rule the concept promises: the answer is the
// fixed point with the larger real part, not merely one of the two roots.
func TestBatchMobiusPicksTheLargerRoot(t *testing.T) {
	reg := batchRegistry(t)
	for i := 0; i < 60; i++ {
		p := sample(t, reg, "complex.adv.mobius_fixed", i)
		lo, err1 := strconv.Atoi(p.Facts["root1"])
		hi, err2 := strconv.Atoi(p.Facts["root2"])
		if err1 != nil || err2 != nil {
			t.Fatalf("facts are not integers: %v", p.Facts)
		}
		if lo >= hi {
			t.Fatalf("roots are not ordered: root1 = %d, root2 = %d", lo, hi)
		}
		if p.Answer != strconv.Itoa(hi) {
			t.Fatalf("answer %q is not the larger root %d", p.Answer, hi)
		}
		// The pair must actually solve z = f(z): substitute both into cz^2 + (d-a)z - b = 0.
		a, _ := strconv.Atoi(p.Facts["a"])
		b, _ := strconv.Atoi(p.Facts["b"])
		c, _ := strconv.Atoi(p.Facts["c"])
		d, _ := strconv.Atoi(p.Facts["d"])
		for _, r := range []int{lo, hi} {
			if c*r*r+(d-a)*r-b != 0 {
				t.Fatalf("root %d does not satisfy the fixed-point equation for a=%d b=%d c=%d d=%d", r, a, b, c, d)
			}
		}
	}
}
