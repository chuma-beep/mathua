package discrete

// Tests for the 2026-10-17 batch. Each concept is graded through the real router with the
// grading_type the corpus ships for it, so a generator and its declared type cannot drift apart.

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
	"discrete.logic.inference",
	"discrete.logic.predicates",
	"discrete.logic.satisfiability",
	"discrete.proof.contrapositive",
	"discrete.proof.counterexample",
	"discrete.proof.structural",
	"discrete.sets.powerset",
	"discrete.sets.partial_order",
	"discrete.sets.countable",
	"discrete.combinatorics.complement",
	"discrete.graphs.bipartite",
	"discrete.graphs.planar",
	"discrete.graphs.cayley",
	"discrete.graphs.degree_sequence",
	"discrete.sequences.linear_recurrence",
	"discrete.sequences.growth_rates",
}

type corpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

// shippedTypes reads the grading_type the corpus declares, so the test cannot pass by assuming one.
func shippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/discrete_math.json")
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
		Seed:       int64(9000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var batchOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

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
				t.Fatalf("%s: no facts published; schemas would fall back", id)
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
			var opts []string
			for _, m := range batchOptionLine.FindAllStringSubmatch(p.Question, -1) {
				opts = append(opts, m[2])
			}
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
				t.Fatalf("%s: answer %q appears %d times among %v", id, p.Answer, found, opts)
			}
			for _, o := range opts {
				if o == p.Answer {
					continue
				}
				if router.Grade(grader.GradingType(types[id]), p.Answer, o).Correct {
					t.Fatalf("%s: distractor %q is accepted as correct", id, o)
				}
			}
			checked++
		}
	}
	if checked == 0 {
		t.Skip("the discrete batch now serves computed answers")
	}
}

func TestBatchTupleAcceptsSpacingAndRejectsReordering(t *testing.T) {
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()
	for i := 0; i < 40; i++ {
		p := sample(t, reg, "discrete.graphs.degree_sequence", i)
		spaced := strings.ReplaceAll(p.Answer, ",", ", ")
		if !router.Grade(grader.GradingTuple, p.Answer, spaced).Correct {
			t.Fatalf("tuple: %q should be accepted for %q", spaced, p.Answer)
		}
		parts := strings.Split(p.Answer, ",")
		for a, b := 0, len(parts)-1; a < b; a, b = a+1, b-1 {
			parts[a], parts[b] = parts[b], parts[a]
		}
		reversed := strings.Join(parts, ",")
		if reversed == p.Answer {
			continue // a palindrome (K(3,3) gives 3,3,3,3,3,3) has nothing to reverse
		}
		if router.Grade(grader.GradingTuple, p.Answer, reversed).Correct {
			t.Fatalf("tuple: the reordered answer %q was accepted for %q", reversed, p.Answer)
		}
	}
}

func TestBatchSymbolicAcceptsEquivalentFormsAndRejectsWrongOnes(t *testing.T) {
	reg := batchRegistry(t)
	router := grader.NewRouter()
	defer router.Close()
	checked := 0
	for i := 0; i < 20; i++ {
		p := sample(t, reg, "discrete.sequences.linear_recurrence", i)
		// The same closed form written with its terms swapped.
		parts := strings.Split(p.Answer, " + ")
		if len(parts) != 2 {
			t.Fatalf("expected two terms, got %q", p.Answer)
		}
		swapped := parts[1] + " + " + parts[0]
		res := router.Grade(grader.GradingSymbolic, p.Answer, swapped)
		if res.Unavailable {
			t.Skip("SymPy runtime missing; symbolic equivalence cannot be checked here")
		}
		if !res.Correct {
			t.Fatalf("symbolic: the commutative form %q was rejected for %q", swapped, p.Answer)
		}
		// A form with the right roots but the wrong coefficients must fail.
		if router.Grade(grader.GradingSymbolic, p.Answer, p.Facts["term1"]).Correct && p.Facts["term2"] != "" {
			t.Fatalf("symbolic: a single-term answer was accepted for %q", p.Answer)
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no symbolic instances were checked")
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

// TestBatchCounterexampleIsMinimal re-derives the claim from the published facts: the value really is
// composite at n, and every smaller value from 0 upward is prime. Without this the generator could
// publish any n and the self-grading test would still pass.
func TestBatchCounterexampleIsMinimal(t *testing.T) {
	reg := batchRegistry(t)
	checked := 0
	for i := 0; i < 40; i++ {
		p := sample(t, reg, "discrete.proof.counterexample", i)
		c, err1 := strconv.Atoi(p.Facts["c"])
		n, err2 := strconv.Atoi(p.Facts["n"])
		if err1 != nil || err2 != nil {
			t.Fatalf("facts are not integers: %v", p.Facts)
		}
		if n < 1 || n > c-1 {
			t.Fatalf("n = %d is outside 1..%d", n, c-1)
		}
		if isPrime(n*n + n + c) {
			t.Fatalf("f(%d) = %d is prime, so it is not a counterexample", n, n*n+n+c)
		}
		for k := 0; k < n; k++ {
			if !isPrime(k*k + k + c) {
				t.Fatalf("f(%d) = %d is composite, so %d is not the smallest", k, k*k+k+c, n)
			}
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no instances checked")
	}
}
