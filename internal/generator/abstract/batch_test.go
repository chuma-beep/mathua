package abstract

// Tests for the 2026-10-18 abstract-algebra batch.
//
// Every concept is graded through the real router with the grading_type the
// corpus ships, so a generator and its declared type cannot drift apart. The
// multiple-choice concepts are checked the way the algebra batch checks its
// own: the question is parsed into options, the correct answer must appear
// exactly once, and every other option must be rejected by the router.
//
// The two numeric concepts are recomputed, not compared with themselves. The
// characteristic's answer is derived from the field kind, prime and order the
// generator publishes, and the p-group's concrete answers are checked against a
// table of centers kept in this test rather than read from the production
// generator. A generator that quietly changed a center would fail here.

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

var abstractBatchIDs = []string{
	"abstract.group.center",
	"abstract.group.automorphism",
	"abstract.group.simple",
	"abstract.group.p_group",
	"abstract.group.composition_series",
	"abstract.rings.zero_divisor",
	"abstract.rings.integral_domain",
	"abstract.rings.prime_ideal",
	"abstract.rings.euclidean",
	"abstract.field.characteristic",
	"abstract.field.algebraic",
	"abstract.structures.submodule",
	"abstract.structures.module_hom",
}

var abstractNumericIDs = []string{
	"abstract.group.p_group",
	"abstract.field.characteristic",
}

// abstractConcreteCenters is the independent table of centers for the p-group
// family that names a group. It is written here, not read from the generator,
// so the test is a check rather than an echo.
var abstractConcreteCenters = map[string]int{
	"Q_8":           2,
	"D_4":           2,
	"D_8":           2,
	"Z_9":           9,
	"V_4":           4,
	"Heisenberg(3)": 3,
}

type abstractCorpusConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func abstractShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/abstract_algebra.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []abstractCorpusConcept
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

func abstractRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func abstractSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	p, err := reg.GenerateContext(id, generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(71000 + i),
	})
	if err != nil {
		t.Fatalf("%s: generate: %v", id, err)
	}
	return p
}

var abstractOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func abstractOptionsOf(q string) []string {
	var out []string
	for _, m := range abstractOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

// ---------------------------------------------------------------------------
// Router contracts
// ---------------------------------------------------------------------------

func TestBatchAnswersGradeWithTheirShippedType(t *testing.T) {
	types := abstractShippedTypes(t)
	reg := abstractRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range abstractBatchIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		if typ != "multiple_choice" && typ != "numeric" {
			t.Fatalf("%s: unexpected grading_type %q", id, typ)
		}
		for i := 0; i < 120; i++ {
			p := abstractSample(t, reg, id, i)
			if strings.TrimSpace(p.Answer) == "" {
				t.Fatalf("%s: empty answer at sample %d", id, i)
			}
			if strings.TrimSpace(p.Question) == "" || strings.TrimSpace(p.Explanation) == "" {
				t.Fatalf("%s: empty question or explanation at sample %d: q=%q e=%q",
					id, i, p.Question, p.Explanation)
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

func TestBatchMultipleChoiceRejectsEveryDistractor(t *testing.T) {
	types := abstractShippedTypes(t)
	reg := abstractRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	checked := 0
	for _, id := range abstractBatchIDs {
		if types[id] != "multiple_choice" {
			continue
		}
		for i := 0; i < 40; i++ {
			p := abstractSample(t, reg, id, i)
			opts := abstractOptionsOf(p.Question)
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
		t.Skip("every batch concept now serves a computed answer")
	}
}

// ---------------------------------------------------------------------------
// Independent numeric checks
// ---------------------------------------------------------------------------

func TestBatchCharacteristicRecomputesFromFacts(t *testing.T) {
	reg := abstractRegistry(t)
	for i := 0; i < 400; i++ {
		p := abstractSample(t, reg, "abstract.field.characteristic", i)
		kind := p.Facts["kind"]
		prime, err := strconv.Atoi(p.Facts["prime"])
		if err != nil {
			t.Fatalf("characteristic: bad prime fact %q", p.Facts["prime"])
		}
		order, err := strconv.Atoi(p.Facts["order"])
		if err != nil {
			t.Fatalf("characteristic: bad order fact %q", p.Facts["order"])
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("characteristic: non-integer answer %q", p.Answer)
		}
		switch kind {
		case "char0":
			if answer != 0 {
				t.Fatalf("characteristic: char-0 field answers %d, want 0", answer)
			}
		case "F_p":
			if order != prime || answer != prime {
				t.Fatalf("characteristic: F_p with prime=%d order=%d answer=%d", prime, order, answer)
			}
		case "F_pn":
			if order <= prime || answer != prime {
				t.Fatalf("characteristic: F_pn with prime=%d order=%d answer=%d", prime, order, answer)
			}
			// order must be a power of the prime, with exponent at least 2.
			n := 0
			v := order
			for v%prime == 0 {
				v /= prime
				n++
			}
			if v != 1 || n < 2 {
				t.Fatalf("characteristic: order %d is not %d^n with n >= 2", order, prime)
			}
		default:
			t.Fatalf("characteristic: unknown kind %q", kind)
		}
	}
}

func TestBatchCharacteristicKeepsOrderDistinctFromCharacteristic(t *testing.T) {
	reg := abstractRegistry(t)
	seenFinitePower, seenCharZero, seenPrimeField := false, false, false
	for i := 0; i < 400; i++ {
		p := abstractSample(t, reg, "abstract.field.characteristic", i)
		switch p.Facts["kind"] {
		case "F_pn":
			seenFinitePower = true
		case "char0":
			seenCharZero = true
		case "F_p":
			seenPrimeField = true
		}
	}
	if !seenFinitePower {
		t.Fatal("characteristic: no F_{p^n} instance in 400 samples, so the order/characteristic distinction is never tested")
	}
	if !seenCharZero {
		t.Fatal("characteristic: no characteristic-0 instance in 400 samples")
	}
	if !seenPrimeField {
		t.Fatal("characteristic: no F_p instance in 400 samples")
	}
}

func TestBatchPGroupRecomputesFromFacts(t *testing.T) {
	reg := abstractRegistry(t)
	seenLeast, seenConcrete := false, false
	for i := 0; i < 400; i++ {
		p := abstractSample(t, reg, "abstract.group.p_group", i)
		kind := p.Facts["kind"]
		prime, err := strconv.Atoi(p.Facts["prime"])
		if err != nil {
			t.Fatalf("p_group: bad prime fact %q", p.Facts["prime"])
		}
		answer, err := strconv.Atoi(p.Answer)
		if err != nil {
			t.Fatalf("p_group: non-integer answer %q", p.Answer)
		}
		switch kind {
		case "least":
			seenLeast = true
			if answer != prime {
				t.Fatalf("p_group: least center asks p=%d but answers %d", prime, answer)
			}
		case "concrete":
			seenConcrete = true
			group := p.Facts["group"]
			want, ok := abstractConcreteCenters[group]
			if !ok {
				t.Fatalf("p_group: concrete group %q has no center in the test table", group)
			}
			if answer != want {
				t.Fatalf("p_group: %s has center order %d, generator answers %d", group, want, answer)
			}
		default:
			t.Fatalf("p_group: unknown kind %q", kind)
		}
	}
	if !seenLeast || !seenConcrete {
		t.Fatalf("p_group: branches not both exercised (least=%v concrete=%v)", seenLeast, seenConcrete)
	}
}

func TestBatchNumericAnswersAreCanonicalIntegers(t *testing.T) {
	reg := abstractRegistry(t)
	for _, id := range abstractNumericIDs {
		for i := 0; i < 120; i++ {
			p := abstractSample(t, reg, id, i)
			v, err := strconv.Atoi(p.Answer)
			if err != nil {
				t.Fatalf("%s: answer %q is not a plain integer", id, p.Answer)
			}
			if strconv.Itoa(v) != p.Answer {
				t.Fatalf("%s: answer %q is not canonical", id, p.Answer)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Pool quality
// ---------------------------------------------------------------------------

func TestBatchPoolsAreDiverse(t *testing.T) {
	reg := abstractRegistry(t)
	const samples = 150
	for _, id := range abstractBatchIDs {
		answers := map[string]bool{}
		questions := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{Difficulty: 0.12, Seed: int64(90000 + i)})
			if err != nil {
				t.Fatalf("%s: generate: %v", id, err)
			}
			answers[strings.ToLower(strings.TrimSpace(p.Answer))] = true
			questions[p.Question] = true
			if !strings.Contains(p.Question, "(yes/no)") {
				allYesNo = false
			}
		}
		if len(answers) < 2 {
			t.Fatalf("%s: only %d distinct answers over %d samples", id, len(answers), samples)
		}
		if len(questions) < 2 {
			t.Fatalf("%s: only %d distinct questions over %d samples", id, len(questions), samples)
		}
		if allYesNo {
			t.Fatalf("%s: every question is yes/no-only", id)
		}
	}
}
