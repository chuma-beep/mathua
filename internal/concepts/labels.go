package concepts

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// NormalizeLabel produces the comparison key for concept titles. Punctuation,
// casing and spacing vary between domains ("L'Hopital's rule" and
// "L'Hôpital's rule" are the same topic), so a duplicate check that compares raw
// labels misses the class of duplicate the curriculum actually has.
func NormalizeLabel(label string) string {
	var b strings.Builder
	space := false
	for _, r := range strings.ToLower(label) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			space = false
		default:
			if !space && b.Len() > 0 {
				b.WriteByte(' ')
				space = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func labelPairKey(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "|" + b
}

// ReviewedDuplicateLabels is the explicit allowlist for normalized title
// collisions that are deliberate. Each pair is two concept ids that teach the
// same idea at different levels or in different domains, such as calculus power
// series and the real-analysis treatment of the same series.
//
// A pair belongs here only after a reviewer has decided the duplication is
// intentional. The list is not a suppression for accidental duplicates: the
// graph test fails on any collision that is not named here.
func ReviewedDuplicateLabels() map[string]bool {
	pairs := [][2]string{
		{"ra.order.supremum", "calc.limit.supremum"},
		{"ra.seq.cauchy", "calc.seq.cauchy"},
		{"ra.seq.absolute", "calc.series.absolute"},
		{"ra.limit.uniform", "calc.limit.uniform_continuity"},
		{"ra.deriv.lhopital", "calc.limit.lhopital"},
		{"ra.deriv.parts", "calc.integral.parts"},
		{"ra.deriv.improper", "calc.integral.improper"},
		{"ra.deriv.power_series", "calc.series.power"},
		{"ra.metric.spaces", "topo.basics.metric"},
		{"ra.metric.open_closed", "topo.basics.open_closed"},
		{"ra.metric.convergence", "topo.metric.convergence"},
		{"alg.func.inverse", "precalc.func.inverse_value"},
		{"alg.exp.compound", "precalc.exp.compound_interest"},
		{"alg.systems.matrices", "linalg.systems.matrix"},
		{"alg.poly.rational_root", "precalc.poly.rational_roots"},
		{"calc.deriv.hyperbolic", "trig.hyperbolic.derivative"},
		{"discrete.proof.induction", "precalc.induction"},
		{"geo.coord.polar", "trig.adv.polar"},
		{"linalg.vector.concept", "precalc.vector.basics"},
	}
	out := make(map[string]bool, len(pairs))
	for _, p := range pairs {
		out[labelPairKey(p[0], p[1])] = true
	}
	return out
}

// CheckDuplicateLabels fails when two concepts normalize to the same label and
// their pair is not in the reviewed allowlist. A group larger than two is
// checked pair by pair, so one unreviewed collision cannot hide behind an
// allowlisted one.
func CheckDuplicateLabels(d *DAG, allow map[string]bool) error {
	groups := make(map[string][]string)
	for _, c := range d.Order() {
		key := NormalizeLabel(c.Label)
		if key == "" {
			continue
		}
		groups[key] = append(groups[key], c.ID)
	}
	var keys []string
	for key, ids := range groups {
		if len(ids) > 1 {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var problems []string
	for _, key := range keys {
		ids := groups[key]
		sort.Strings(ids)
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				if allow[labelPairKey(ids[i], ids[j])] {
					continue
				}
				problems = append(problems, fmt.Sprintf("%q: %s and %s", key, ids[i], ids[j]))
			}
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("duplicate concept labels (%d): %s", len(problems), strings.Join(problems, "; "))
	}
	return nil
}

// CheckReviewedDuplicateLabelsExist catches stale allowlist entries after a
// concept is renamed or removed.
func CheckReviewedDuplicateLabelsExist(d *DAG, allow map[string]bool) error {
	for key := range allow {
		parts := strings.SplitN(key, "|", 2)
		if len(parts) != 2 {
			return fmt.Errorf("malformed duplicate-label key %q", key)
		}
		for _, id := range parts {
			if d.Concept(id) == nil {
				return fmt.Errorf("duplicate-label allowlist names missing concept %q", id)
			}
		}
	}
	return nil
}
