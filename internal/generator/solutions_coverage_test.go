package generator_test

import (
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/solutions"
)

// learnerDomains are the six domains a learner is expected to work through.
// Every concept in them that computes from parameters must ship an authored
// schema, so a missed answer is never a bare derivation trace.
var learnerDomains = map[string]bool{
	"arithmetic": true, "fractions": true, "prealgebra": true,
	"algebra": true, "geometry": true, "trigonometry": true,
}

// generatorExplained are the concepts deliberately left on their generator's
// own explanation. Every entry is table-driven: the question is a fixed
// string, the answer is a fixed string, and the explanation is already a
// complete sentence, "the sine function oscillates between -1 and 1; the
// maximum is 1". There is no arithmetic to interpolate and no instance to
// address, so a schema would relocate text and buy nothing.
//
// The list is explicit rather than inferred, and that is the point: a new
// table-driven concept in these domains has to be added here deliberately,
// with a reason, instead of passing because nobody looked. A concept that
// *does* compute from parameters belongs in data/lessons/solutions, not here.
var generatorExplained = map[string]string{
	"arith.div.zero":           "numeric count question: per-statement explanation from a generated pool",
	"arith.factor.composite":   "numeric count question: per-statement explanation from a generated pool",
	"arith.factor.prime":       "numeric count question: per-statement explanation from a generated pool",
	"frac.compare.word":        "numeric count question: per-statement explanation from a generated pool",
	"frac.ops.benchmark":       "numeric count question: per-statement explanation from a generated pool",
	"prealg.ineq.graph":        "numeric conversion: explanation is computed from the drawn inequality",
	"alg.conic.classify":       "numeric count question: per-statement explanation from a generated pool",
	"alg.ineq.two_var":         "numeric conversion: explanation is computed from the drawn inequality",
	"alg.linear.standard_form": "numeric conversion: explanation is computed from the substituted point",
	"alg.poly.end_behavior":    "numeric conversion: explanation is computed from degree and leading coefficient",
	"geo.basic.angles":         "numeric count question: per-statement explanation from a generated pool",
	"geo.basic.points_lines":   "numeric count question: per-statement explanation from a generated pool",
	"geo.circle.parts":         "numeric count question: per-statement explanation from a generated pool",
	"geo.quad.types":           "numeric count question: per-statement explanation from a generated pool",
	"geo.triangle.types":       "numeric count question: per-statement explanation from a generated pool",
	// trigonometry, table-driven, complete-sentence explanations
	"trig.adv.arctan":                 "table-driven: fixed question, fixed answer",
	"trig.adv.inverse":                "table-driven: fixed question, fixed answer",
	"trig.eq.basic":                   "table-driven: fixed question, fixed answer",
	"trig.eq.general_solution":        "table-driven: fixed question, fixed answer",
	"trig.eq.homogeneous":             "table-driven: fixed question, fixed answer",
	"trig.func.arccosine":             "table-driven: fixed question, fixed answer",
	"trig.func.arccotangent":          "table-driven: fixed question, fixed answer",
	"trig.func.arcsine":               "table-driven: fixed question, fixed answer",
	"trig.func.arctangent":            "table-driven: fixed question, fixed answer",
	"trig.graph.cos":                  "table-driven: fixed question, fixed answer",
	"trig.graph.inverse":              "table-driven: fixed question, fixed answer",
	"trig.graph.reciprocal":           "table-driven: fixed question, fixed answer",
	"trig.graph.sin":                  "table-driven: fixed question, fixed answer",
	"trig.graph.tan":                  "table-driven: fixed question, fixed answer",
	"trig.graph.transformations":      "table-driven: fixed question, fixed answer",
	"trig.hyperbolic.derivative":      "table-driven: fixed question, fixed answer",
	"trig.hyperbolic.identities":      "table-driven: fixed question, fixed answer",
	"trig.hyperbolic.sinh_cosh":       "table-driven: fixed question, fixed answer",
	"trig.hyperbolic.tanh_coth":       "table-driven: fixed question, fixed answer",
	"trig.ident.identities":           "table-driven: fixed question, fixed answer",
	"trig.ident.product_to_sum":       "table-driven: fixed question, fixed answer",
	"trig.ident.pythagorean_variants": "table-driven: fixed question, fixed answer",
	"trig.ident.sum_to_product":       "table-driven: fixed question, fixed answer",
	"trig.ineq.advanced":              "table-driven: fixed question, fixed answer",
	// complex analysis, count questions explained per statement
	"trig.basics.radians":         "count question: per-statement explanation from a generated pool",
	"trig.basics.reference_angle": "count question: per-statement explanation from a generated pool",
	"trig.basics.special_angles":  "count question: per-statement explanation from a generated pool",
	"trig.basics.unit_circle":     "count question: per-statement explanation from a generated pool",
	"trig.ineq.basic":             "table-driven: fixed question, fixed answer",
	// prealgebra, table-driven
	"prealg.real.concept":    "table-driven: fixed question, fixed answer",
	"prealg.real.properties": "table-driven: fixed question, fixed answer",
	"prealg.types":           "table-driven: fixed question, fixed answer",
	// algebra, table-driven
	"alg.func.dirichlet":       "table-driven: fixed question, fixed answer",
	"alg.func.even_odd":        "table-driven: fixed question, fixed answer",
	"alg.func.inverse_domain":  "table-driven: fixed question, fixed answer",
	"alg.func.sigmoid":         "table-driven: fixed question, fixed answer",
	"alg.func.transformations": "table-driven: fixed question, fixed answer",
	"alg.log.graph":            "table-driven: fixed question, fixed answer",
	"alg.seq.induction":        "table-driven: fixed question, fixed answer",
	"alg.systems.matrices":     "table-driven: fixed question, fixed answer",
}

// TestLearnerDomainsHaveSchemas is the gate that keeps the explanation
// contract from decaying. It is report-shaped on purpose: the exempt list is
// named and reasoned, and anything new shows up in the failure text so it gets
// a decision rather than a shrug.
func TestLearnerDomainsHaveSchemas(t *testing.T) {
	l, err := solutions.Load("../../data/lessons/solutions")
	if err != nil {
		t.Fatalf("load schemas: %v", err)
	}
	have := map[string]bool{}
	for _, id := range l.Concepts() {
		have[id] = true
	}
	// The same loader production uses, so the gate cannot disagree with what
	// the server actually serves.
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load concepts: %v", err)
	}
	inDomain := map[string]bool{}
	var missing []string
	total := 0
	for id, c := range dag.Concepts() {
		if c == nil || !learnerDomains[c.Domain] {
			continue
		}
		inDomain[id] = true
		total++
		if !have[id] && generatorExplained[id] == "" {
			missing = append(missing, id)
		}
	}
	// An exemption for a concept outside the six domains, or one that now has a
	// schema, is dead weight and hides a real change.
	for id := range generatorExplained {
		if !inDomain[id] {
			t.Errorf("generatorExplained lists %q, which is not in a learner domain", id)
		}
		if have[id] {
			t.Errorf("generatorExplained lists %q, but it now has a schema, drop the exemption", id)
		}
	}
	if len(missing) > 0 {
		t.Errorf("%d concepts in the six learner domains have no authored schema and no exemption:\n  %s\n"+
			"either author data/lessons/solutions/%s.json, or add it to generatorExplained with a reason",
			len(missing), strings.Join(missing, "\n  "), missing[0])
	}
	authored, exempt := 0, 0
	for id := range inDomain {
		if have[id] {
			authored++
		} else {
			exempt++
		}
	}
	t.Logf("learner domains: %d concepts, %d authored schemas, %d generator-explained", total, authored, exempt)
}
