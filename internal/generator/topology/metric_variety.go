package topology

import (
	"fmt"
	"math/rand"
)

// topoClaim is the shared shape of a recognition claim: a question, the one
// right answer, three plausible wrong answers and the reasoning that explains
// it. The generators in batch.go declared this type inside each function while
// each concept carried three claims; computed claims are appended to those same
// pools, so the type has to be visible to both.
type topoClaim struct {
	question    string
	correct     string
	distractors []string
	reason      string
}

// topoDistinct returns three distractors that differ from the answer and from
// each other.
//
// generator.MultipleChoice dedupes its options and caps them at the four
// letters a recognition question can use, so a distractor that repeats the
// answer silently drops the question to three options. Topping up from a fixed
// pool keeps every claim answerable whatever numbers were drawn.
func topoDistinct(correct string, candidates ...string) []string {
	seen := map[string]bool{correct: true}
	out := make([]string, 0, 3)
	for _, c := range candidates {
		if c == "" || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	for _, filler := range []string{"0", "1", "2", "none of these"} {
		if len(out) == 3 {
			break
		}
		if !seen[filler] {
			seen[filler] = true
			out = append(out, filler)
		}
	}
	return out[:3]
}

// computedTopoClaims returns claims whose numbers are drawn, so the pool keeps
// growing instead of reciting the same sentences.
func computedTopoClaims(id string) []topoClaim {
	switch id {
	case "topo.metric.basis":
		n := 2 + rand.Intn(6) // 2..7
		power := 1 << n
		return []topoClaim{
			{
				question:    fmt.Sprintf("A set X has %d points and carries the discrete topology. How many open sets does X have?", n),
				correct:     fmt.Sprintf("%d", power),
				distractors: topoDistinct(fmt.Sprintf("%d", power), fmt.Sprintf("%d", n*n), fmt.Sprintf("%d", 2*n), fmt.Sprintf("%d", n)),
				reason:      fmt.Sprintf("Every subset is open in the discrete topology, and a set with %d points has %d subsets.", n, power),
			},
			{
				question:    fmt.Sprintf("How many sets does a basis need at minimum to generate the discrete topology on a set with %d points?", n),
				correct:     fmt.Sprintf("%d", n),
				distractors: topoDistinct(fmt.Sprintf("%d", n), fmt.Sprintf("%d", power), fmt.Sprintf("%d", n-1), fmt.Sprintf("%d", 2*n)),
				reason:      fmt.Sprintf("The singletons form a basis, and there are %d of them; any coarser family cannot separate the points.", n),
			},
			{
				question:    "The open intervals with rational endpoints form a basis for which topology on the real line?",
				correct:     "the usual topology",
				distractors: topoDistinct("the usual topology", "the discrete topology", "the Sorgenfrey topology", "the cofinite topology"),
				reason:      "Every usual open interval is a union of intervals with rational endpoints, because the rationals are dense and the endpoints can be approximated from inside.",
			},
		}
	case "topo.metric.dense":
		n := 2 + rand.Intn(6) // 2..7
		power := 1 << n
		return []topoClaim{
			{
				question:    "Is the set of integers dense in the real line with the usual topology?",
				correct:     "no",
				distractors: topoDistinct("no", "yes", "only the negative integers are", "it depends on the metric"),
				reason:      "The open interval (0, 1) contains no integer, so the integers miss an open set and cannot be dense.",
			},
			{
				question:    fmt.Sprintf("A set has %d points and carries the discrete topology. How many of its subsets are dense?", n),
				correct:     fmt.Sprintf("%d", power),
				distractors: topoDistinct(fmt.Sprintf("%d", power), fmt.Sprintf("%d", n), fmt.Sprintf("%d", 2*n), "1"),
				reason:      fmt.Sprintf("In the discrete topology every subset is open and closed, so each of the %d subsets is its own closure and is dense in itself.", power),
			},
			{
				question:    "Every nonempty open interval of the real line contains how many rational numbers?",
				correct:     "infinitely many",
				distractors: topoDistinct("infinitely many", "exactly one", "finitely many but at least one", "none"),
				reason:      "Between any two reals there is a rational, and repeating that inside the interval produces infinitely many.",
			},
		}
	case "topo.metric.subspace":
		n := 2 + rand.Intn(6) // 2..7
		power := 1 << n
		return []topoClaim{
			{
				question:    "The interval [0, 1] carries the subspace topology inherited from the real line. Is [0, 1] open in that subspace?",
				correct:     "yes",
				distractors: topoDistinct("yes", "no", "only if it is also closed", "only at the endpoints"),
				reason:      "The subspace opens are the intersections of ambient opens with [0, 1], and the whole space is always open in itself.",
			},
			{
				question:    "Is the open interval (0, 1), with the subspace topology from the real line, homeomorphic to the real line?",
				correct:     "yes",
				distractors: topoDistinct("yes", "no", "only after adding the endpoints", "only if both carry the discrete topology"),
				reason:      "For example x is sent to tan(pi(x - 1/2)), a continuous bijection with continuous inverse between (0, 1) and the real line.",
			},
			{
				question:    fmt.Sprintf("A set with %d points sits inside the real line and carries the subspace topology. How many open sets does it have?", n),
				correct:     fmt.Sprintf("%d", power),
				distractors: topoDistinct(fmt.Sprintf("%d", power), fmt.Sprintf("%d", n), fmt.Sprintf("%d", 2*n), fmt.Sprintf("%d", n*n)),
				reason:      fmt.Sprintf("Finitely many points in the real line are separated by small intervals, so the subspace is discrete and has %d open sets.", power),
			},
		}
	}
	return nil
}
