package realanalysis

import (
	"fmt"
	"math/rand"
)

// computedMetricClaims adds questions whose answers are decided rather than
// recited. Each of the metric-space generators shipped three definitional
// claims and nothing else, so a learner who met all three had seen the whole
// concept. These claims vary the space and the numbers, which is where the
// variation has to come from when the definition itself is one sentence.
func computedMetricClaims(id string) []mcClaim {
	switch id {
	case "ra.metric.open_closed":
		a := 1 + rand.Intn(8)
		b := a + 1 + rand.Intn(8)
		return []mcClaim{
			{
				q:           fmt.Sprintf("In R with the usual metric, is the interval (%d, %d) open?", a, b),
				correct:     "yes",
				distractors: []string{"no", "only if it is bounded", "only if its complement is closed"},
				reason:      fmt.Sprintf("Every point x of (%d, %d) has the ball of radius min(x-%d, %d-x) inside the interval.", a, b, a, b),
			},
			{
				q:           fmt.Sprintf("Is the interval [%d, %d] open in R with the usual metric?", a, b),
				correct:     "no",
				distractors: []string{"yes", "yes, because it is bounded", "yes, because it is closed"},
				reason:      fmt.Sprintf("The endpoint %d has no ball that stays inside [%d, %d], so the set is not open.", a, a, b),
			},
			{
				q:           fmt.Sprintf("Is the interval [%d, %d) closed in R with the usual metric?", a, b),
				correct:     "no",
				distractors: []string{"yes", "yes, because it contains an endpoint", "yes, because it is bounded"},
				reason:      fmt.Sprintf("%d is a limit point of the interval that the interval does not contain, so its complement is not open.", b),
			},
		}
	case "ra.metric.closure":
		a := 1 + rand.Intn(8)
		b := a + 1 + rand.Intn(8)
		return []mcClaim{
			{
				q:           fmt.Sprintf("What is the closure of the open interval (%d, %d) in R?", a, b),
				correct:     fmt.Sprintf("[%d, %d]", a, b),
				distractors: []string{fmt.Sprintf("(%d, %d)", a, b), fmt.Sprintf("[%d, %d)", a, b), fmt.Sprintf("{%d, %d}", a, b)},
				reason:      fmt.Sprintf("The endpoints %d and %d are limit points of the interval that it omits, and adding them makes the set closed.", a, b),
			},
			{
				q:           fmt.Sprintf("What is the closure of the half-open interval [%d, %d) in R?", a, b),
				correct:     fmt.Sprintf("[%d, %d]", a, b),
				distractors: []string{fmt.Sprintf("[%d, %d)", a, b), fmt.Sprintf("(%d, %d)", a, b), fmt.Sprintf("(%d, %d]", a, b)},
				reason:      fmt.Sprintf("Only the missing right endpoint %d is a limit point outside the set, so the closure adds it.", b),
			},
			{
				q:           fmt.Sprintf("What is the closure of the singleton {%d} in R?", a),
				correct:     fmt.Sprintf("{%d}", a),
				distractors: []string{"the empty set", "R", fmt.Sprintf("(%d, %d)", a-1, a+1)},
				reason:      "A singleton contains all of its own limit points, because its only candidate limit point is the point itself.",
			},
		}
	case "ra.metric.convergence":
		p := 1 + rand.Intn(3)
		k := 2 + rand.Intn(6) // k must differ from 1, which is the limit itself
		return []mcClaim{
			{
				q:           fmt.Sprintf("In R with the usual metric, does the sequence 1/n^%d converge?", p),
				correct:     "yes, to 0",
				distractors: []string{"no", "yes, to 1", "only if it is monotone"},
				reason:      fmt.Sprintf("For every epsilon the term 1/n^%d is below epsilon once n exceeds epsilon^(-1/%d).", p, p),
			},
			{
				q:           fmt.Sprintf("In R with the usual metric, what is the limit of n/(n+%d)?", k),
				correct:     "1",
				distractors: []string{"0", fmt.Sprintf("%d", k), "the sequence does not converge"},
				reason:      fmt.Sprintf("Dividing through by n gives 1/(1+%d/n), which tends to 1.", k),
			},
			{
				q:           "In R with the discrete metric, does the sequence 1/n converge?",
				correct:     "no",
				distractors: []string{"yes, to 0", "yes, to 1", "yes, because it is Cauchy"},
				reason:      "The discrete metric makes every pair of distinct points distance 1 apart, and 1/n is not eventually constant.",
			},
		}
	case "ra.metric.complete":
		a := 1 + rand.Intn(6)
		b := a + 1 + rand.Intn(6)
		return []mcClaim{
			{
				q:           fmt.Sprintf("Is the open interval (%d, %d) complete with the usual metric?", a, b),
				correct:     "no",
				distractors: []string{"yes", "yes, because it is bounded", "yes, because it is a subset of a complete space"},
				reason:      fmt.Sprintf("A sequence inside (%d, %d) can converge to %d, which lies outside the space, so the space is not complete.", a, b, b),
			},
			{
				q:           fmt.Sprintf("Is the closed interval [%d, %d] complete with the usual metric?", a, b),
				correct:     "yes",
				distractors: []string{"no", "only if it is compact in every metric", "only if it is open"},
				reason:      "A closed subset of a complete metric space is complete in the induced metric.",
			},
			{
				q:           fmt.Sprintf("Is the set of integers Z with the metric d(x,y)=%d|x-y| complete?", a),
				correct:     "yes",
				distractors: []string{"no", "only if it is bounded", "only if it is dense in R"},
				reason:      fmt.Sprintf("Scaling a metric by %d does not change its Cauchy sequences, and Z is complete with the usual metric.", a),
			},
		}
	case "ra.metric.compact":
		a := 1 + rand.Intn(6)
		b := a + 1 + rand.Intn(6)
		return []mcClaim{
			{
				q:           fmt.Sprintf("Is the closed interval [%d, %d] compact with the usual metric?", a, b),
				correct:     "yes",
				distractors: []string{"no", "only if it is open", "only if it is countable"},
				reason:      "A closed bounded interval in R is compact by the Heine-Borel theorem.",
			},
			{
				q:           fmt.Sprintf("Is the interval (%d, %d] compact with the usual metric?", a, b),
				correct:     "no",
				distractors: []string{"yes", "yes, because it is bounded", "yes, because it is connected"},
				reason:      fmt.Sprintf("The cover by the sets (%d + 1/n, %d] has no finite subcover, because %d is not in the space.", a, b, a),
			},
			{
				q:           fmt.Sprintf("Is the interval [%d, infinity) compact with the usual metric?", a),
				correct:     "no",
				distractors: []string{"yes", "yes, because it is closed", "yes, because it has a least element"},
				reason:      fmt.Sprintf("The cover by [%d, %d+n) over n has no finite subcover.", a, a),
			},
		}
	case "ra.metric.continuity":
		p := 2 + rand.Intn(3)
		a := 1 + rand.Intn(5)
		c := 2 + rand.Intn(8)
		return []mcClaim{
			{
				q:           fmt.Sprintf("Is f(x)=x^%d uniformly continuous on [0, %d]?", p, a),
				correct:     "yes",
				distractors: []string{"no", "only if it is injective", "only if the interval is open"},
				reason:      "A continuous function on a compact interval is uniformly continuous.",
			},
			{
				q:           fmt.Sprintf("Is f(x)=1/x uniformly continuous on (0, %d]?", a),
				correct:     "no",
				distractors: []string{"yes", "yes, because it is continuous", "yes, because the interval is bounded"},
				reason:      "Near 0 the function changes by arbitrarily large amounts over arbitrarily small distances, so one delta cannot serve every point.",
			},
			{
				q:           fmt.Sprintf("Is f(x)=%d x + %d uniformly continuous on R?", c, a),
				correct:     "yes",
				distractors: []string{"no", "only on a bounded interval", "only if it is bounded"},
				reason:      fmt.Sprintf("The map is linear with slope %d, so delta equal to epsilon divided by %d works at every point.", c, c),
			},
		}
	case "ra.metric.equivalent":
		k := 2 + rand.Intn(5)
		return []mcClaim{
			{
				q:           fmt.Sprintf("Are the metrics d and %d d on the same set topologically equivalent?", k),
				correct:     "yes",
				distractors: []string{"no", "only if the set is bounded", "only if d is complete"},
				reason:      fmt.Sprintf("Scaling by %d changes every distance but not which sequences converge, so the open sets agree.", k),
			},
			{
				q:           "Are the metrics d and d/(1+d) topologically equivalent?",
				correct:     "yes",
				distractors: []string{"no", "only if d is bounded", "only if the space is compact"},
				reason:      "The map t to t/(1+t) is increasing and continuous at 0, so small distances stay small in both directions.",
			},
			{
				q:           "Is the usual metric on R topologically equivalent to the discrete metric?",
				correct:     "no",
				distractors: []string{"yes", "yes, because both are complete", "yes, because both are unbounded"},
				reason:      "The discrete metric makes every singleton open, which the usual metric does not.",
			},
		}
	case "ra.metric.baire":
		k := 2 + rand.Intn(7)
		return []mcClaim{
			{
				q:           "Are the rationals of first category in R?",
				correct:     "yes",
				distractors: []string{"no", "only if they are bounded", "only if R is compact"},
				reason:      "The rationals are a countable union of singletons, and each singleton is nowhere dense in R.",
			},
			{
				q:           fmt.Sprintf("Can R be written as a union of %d nowhere dense sets?", k),
				correct:     "no",
				distractors: []string{"yes, for every finite number of them", "yes, because R is complete", "yes, because R is uncountable"},
				reason:      "A finite union of nowhere dense sets is nowhere dense, and R is not nowhere dense in itself.",
			},
			{
				q:           "Is a countable union of nowhere dense sets meagre?",
				correct:     "yes",
				distractors: []string{"no", "only if the union is closed", "only if the space is compact"},
				reason:      "Meagre means exactly a countable union of nowhere dense sets.",
			},
		}
	}
	return nil
}
