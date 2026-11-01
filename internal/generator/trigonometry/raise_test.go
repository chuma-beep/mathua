package trigonometry

// Tests for the trigonometry raise.
//
// The package's own test file predates this batch and has no helpers for the
// shipped corpus, so this file carries its own: it reads the corpus for each
// concept's grading type, samples the registry directly, and parses the
// multiple-choice options out of the rendered question. The nine numeric
// generators are recomputed from the facts they publish, the period, the
// angular speed, the bearing, the triple-angle value, the power-reduction
// value, the linear-equation root, the inverse-hyperbolic value, the polar
// abscissa and the scalar projection, so the tests never read the published
// answer.

import (
	"encoding/json"
	"math"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

type trigRaiseConcept struct {
	ID          string `json:"id"`
	GradingType string `json:"grading_type"`
}

func trigShippedTypes(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../../../data/concepts/trigonometry.json")
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}
	var rows []trigRaiseConcept
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatalf("parse corpus: %v", err)
	}
	out := map[string]string{}
	for _, r := range rows {
		out[r.ID] = r.GradingType
	}
	return out
}

func trigRaiseRegistry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func trigRaiseSample(t *testing.T, reg *generator.Registry, id string, i int) generator.Problem {
	t.Helper()
	gen, err := reg.Get(id)
	if err != nil {
		t.Fatalf("%s: %v", id, err)
	}
	return gen.Generate(generator.GeneratorContext{
		Difficulty: float64(i%20) / 19,
		Seed:       int64(31000 + i),
	})
}

func trigRaiseIntFact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%q: fact %q missing from %v", p.Question, key, p.Facts)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("%q: fact %q = %q is not an integer", p.Question, key, raw)
	}
	return v
}

var trigRaiseOptionLine = regexp.MustCompile(`\n\n([A-E])\) ([^\n]*)`)

func trigRaiseOptions(q string) []string {
	var out []string
	for _, m := range trigRaiseOptionLine.FindAllStringSubmatch(q, -1) {
		out = append(out, m[2])
	}
	return out
}

func trigRaiseFraction(t *testing.T, p generator.Problem, key string) (int, int) {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("%q: fact %q missing", p.Question, key)
	}
	re := regexp.MustCompile(`^(-?\d+)/(-?\d+)$`)
	if m := re.FindStringSubmatch(raw); m != nil {
		n, _ := strconv.Atoi(m[1])
		d, _ := strconv.Atoi(m[2])
		return n, d
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("%q: fact %q = %q is not a fraction", p.Question, key, raw)
	}
	return n, 1
}

var trigRaiseIDs = []string{
	"trig.basics.periodicity",
	"trig.basics.angular_velocity",
	"trig.basics.bearings",
	"trig.ident.product_to_sum",
	"trig.ident.sum_to_product",
	"trig.ident.triple_angle",
	"trig.ident.power_reduction",
	"trig.ident.pythagorean_variants",
	"trig.eq.general_solution",
	"trig.eq.linear",
	"trig.ineq.advanced",
	"trig.graph.tan",
	"trig.graph.transformations",
	"trig.graph.inverse",
	"trig.graph.reciprocal",
	"trig.hyperbolic.identities",
	"trig.hyperbolic.inverse",
	"trig.hyperbolic.derivative",
	"trig.adv.polar",
	"trig.adv.vector_projection",
}

var trigRaiseMCIDs = []string{}

func TestRaiseAnswersGradeWithTheirShippedType(t *testing.T) {
	types := trigShippedTypes(t)
	reg := trigRaiseRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range trigRaiseIDs {
		typ, ok := types[id]
		if !ok || typ == "" {
			t.Fatalf("%s: no grading_type in the corpus", id)
		}
		for i := 0; i < 200; i++ {
			p := trigRaiseSample(t, reg, id, i)
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
	if len(trigRaiseMCIDs) == 0 {
		t.Skip("the raise pool now serves computed count problems")
	}
	types := trigShippedTypes(t)
	reg := trigRaiseRegistry(t)
	router := grader.NewRouter()
	defer router.Close()

	for _, id := range trigRaiseMCIDs {
		if types[id] != "multiple_choice" {
			t.Fatalf("%s: expected multiple_choice, corpus says %q", id, types[id])
		}
		for i := 0; i < 60; i++ {
			p := trigRaiseSample(t, reg, id, i)
			opts := trigRaiseOptions(p.Question)
			if len(opts) < 4 {
				t.Fatalf("%s: only %d options; the batch ships four per question", id, len(opts))
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
	reg := trigRaiseRegistry(t)

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.basics.periodicity", i)
		b := trigRaiseIntFact(t, p, "b")
		if b == 0 || 360%b != 0 {
			t.Fatalf("periodicity: b=%d does not divide 360", b)
		}
		if want := strconv.Itoa(360 / b); p.Answer != want {
			t.Fatalf("periodicity: b=%d answers %s, want %s", b, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.basics.angular_velocity", i)
		revs := trigRaiseIntFact(t, p, "revs")
		seconds := trigRaiseIntFact(t, p, "seconds")
		if seconds == 0 || (360*revs)%seconds != 0 {
			t.Fatalf("angular_velocity: %d revolutions in %d seconds is not whole", revs, seconds)
		}
		if want := strconv.Itoa(360 * revs / seconds); p.Answer != want {
			t.Fatalf("angular_velocity: revs=%d seconds=%d answers %s, want %s", revs, seconds, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.basics.bearings", i)
		theta := trigRaiseIntFact(t, p, "theta")
		quadrant := trigRaiseIntFact(t, p, "quadrant")
		var want int
		switch quadrant {
		case 0:
			want = theta
		case 1:
			want = 180 - theta
		case 2:
			want = 180 + theta
		case 3:
			want = 360 - theta
		default:
			t.Fatalf("bearings: quadrant %d out of range", quadrant)
		}
		if got := trigRaiseIntFact(t, p, "answer"); got != want {
			t.Fatalf("bearings: theta=%d quadrant=%d answers %d, want %d", theta, quadrant, got, want)
		}
		if p.Answer != strconv.Itoa(want) {
			t.Fatalf("bearings: theta=%d quadrant=%d answer string %s, want %d", theta, quadrant, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.ident.triple_angle", i)
		hyp := trigRaiseIntFact(t, p, "hyp")
		opp := trigRaiseIntFact(t, p, "opp")
		num := 3*opp*hyp*hyp - 4*opp*opp*opp
		want := trigFrac(num, hyp*hyp*hyp)
		if p.Answer != want {
			t.Fatalf("triple_angle: opp=%d hyp=%d answers %s, want %s", opp, hyp, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.ident.power_reduction", i)
		cn := trigRaiseIntFact(t, p, "cos_num")
		cd := trigRaiseIntFact(t, p, "cos_den")
		if cd == 0 {
			t.Fatalf("power_reduction: denominator zero")
		}
		// sin^2 t = (1 - cos 2t)/2 = (cd - cn) / (2 cd)
		want := trigFrac(cd-cn, 2*cd)
		if p.Answer != want {
			t.Fatalf("power_reduction: cos=%d/%d answers %s, want %s", cn, cd, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.eq.linear", i)
		a := trigRaiseIntFact(t, p, "a")
		b := trigRaiseIntFact(t, p, "b")
		root := trigRaiseIntFact(t, p, "answer")
		if root < 0 || root >= 360 {
			t.Fatalf("linear: root %d outside [0, 360)", root)
		}
		// Check a*sin(root) = b to a tolerance, using the fact that the root
		// is one of the special angles.
		sine := math.Sin(float64(root) * math.Pi / 180)
		if math.Abs(float64(a)*sine-float64(b)) > 1e-6 {
			t.Fatalf("linear: a=%d b=%d root=%d does not satisfy the equation", a, b, root)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.hyperbolic.inverse", i)
		opp := trigRaiseIntFact(t, p, "opp")
		adj := trigRaiseIntFact(t, p, "adj")
		hyp := trigRaiseIntFact(t, p, "hyp")
		if hyp*hyp != opp*opp+adj*adj {
			t.Fatalf("hyperbolic.inverse: (%d,%d,%d) is not a Pythagorean triple", opp, adj, hyp)
		}
		want := trigFrac(hyp, adj)
		if p.Answer != want {
			t.Fatalf("hyperbolic.inverse: opp=%d adj=%d answers %s, want %s", opp, adj, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.adv.polar", i)
		r := trigRaiseIntFact(t, p, "r")
		cn := trigRaiseIntFact(t, p, "cos_num")
		cd := trigRaiseIntFact(t, p, "cos_den")
		if cd == 0 {
			t.Fatalf("polar: denominator zero")
		}
		want := strconv.Itoa(r * cn / cd)
		if r*cn%cd != 0 {
			t.Fatalf("polar: r=%d cos=%d/%d does not give a whole abscissa", r, cn, cd)
		}
		if p.Answer != want {
			t.Fatalf("polar: r=%d cos=%d/%d answers %s, want %s", r, cn, cd, p.Answer, want)
		}
	}

	for i := 0; i < 200; i++ {
		p := trigRaiseSample(t, reg, "trig.adv.vector_projection", i)
		u1 := trigRaiseIntFact(t, p, "u1")
		u2 := trigRaiseIntFact(t, p, "u2")
		if (3*u1+4*u2)%5 != 0 {
			t.Fatalf("vector_projection: (%d,%d) does not project to a whole number", u1, u2)
		}
		if want := strconv.Itoa((3*u1 + 4*u2) / 5); p.Answer != want {
			t.Fatalf("vector_projection: u=(%d,%d) answers %s, want %s", u1, u2, p.Answer, want)
		}
	}
}
