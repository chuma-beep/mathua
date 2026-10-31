package realanalysis

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// raVerdict is one candidate in a counting question. The generator draws a
// fixed number of candidates and asks how many satisfy the condition, so a
// learner has to decide every case instead of recognising one option.
type raVerdict struct {
	label string
	yes   bool
	note  string
}

const raDraw = 8

// countRAVerdicts draws the answer first, then selects that many true cases
// and the rest false ones. Drawing first keeps the whole range 0..8 in play;
// sampling the pool directly makes the extreme counts vanishingly rare.
func countRAVerdicts(lead string, pool []raVerdict) generator.Problem {
	var yes, no []raVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}

	count := rand.Intn(raDraw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if raDraw-count > len(no) {
		count = raDraw - len(no)
	}

	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })

	chosen := append([]raVerdict{}, yes[:count]...)
	chosen = append(chosen, no[:raDraw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	parts := make([]string, len(chosen))
	notes := make([]string, len(chosen))
	for i, v := range chosen {
		parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		verdict := "no"
		if v.yes {
			verdict = "yes"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
	}

	return generator.Problem{
		Question:    fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " ")),
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the %d items satisfy the condition. %s", count, raDraw, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func countableProblem() generator.Problem {
	return countRAVerdicts("How many of these sets are countable?", []raVerdict{
		{label: "N", yes: true, note: "the natural numbers are the model enumeration"},
		{label: "Z", yes: true, note: "list zero, then each positive and negative pair"},
		{label: "Q", yes: true, note: "rationals are pairs of integers up to a diagonal enumeration"},
		{label: "Q intersect [0, 1]", yes: true, note: "a subset of a countable set is countable"},
		{label: "the algebraic numbers", yes: true, note: "countably many integer polynomials, finitely many roots each"},
		{label: "Z x Z", yes: true, note: "a countable product of countable sets is countable"},
		{label: "{1, 2, 3}", yes: true, note: "a finite set is countable"},
		{label: "Q union {sqrt 2}", yes: true, note: "a countable set plus one point stays countable"},
		{label: "R", yes: false, note: "the reals are uncountable by Cantor's diagonal argument"},
		{label: "[0, 1]", yes: false, note: "an interval of positive length is uncountable"},
		{label: "the irrational numbers", yes: false, note: "if they were countable, R would be a union of two countable sets"},
		{label: "the Cantor set", yes: false, note: "it is in bijection with the power set of N"},
		{label: "the power set of N", yes: false, note: "Cantor's theorem gives a strictly larger cardinality"},
		{label: "the transcendental numbers", yes: false, note: "the algebraic numbers are countable, so the complement is not"},
		{label: "(0, 1)", yes: false, note: "it has the cardinality of R"},
		{label: "R minus Z", yes: false, note: "removing a countable set from R leaves an uncountable set"},
	})
}

func openClosedProblem() generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 2 + rand.Intn(6)
	c := b + 2
	d := c + 2

	kind := "open"
	if rand.Intn(2) == 1 {
		kind = "closed"
	}

	pool := []raVerdict{
		{label: fmt.Sprintf("(%d, %d)", a, b), yes: kind == "open", note: "an open interval is open and its endpoints are limit points outside it"},
		{label: fmt.Sprintf("(%d, inf)", a), yes: kind == "open", note: "an open ray is open and unbounded, so it is not closed"},
		{label: fmt.Sprintf("(-inf, %d)", b), yes: kind == "open", note: "an open ray is open and unbounded, so it is not closed"},
		{label: fmt.Sprintf("R minus {%d}", a), yes: kind == "open", note: "the complement is a singleton, which is closed"},
		{label: fmt.Sprintf("(%d, %d) union (%d, %d)", a, b, c, d), yes: kind == "open", note: "a union of open intervals is open and misses the limit points"},
		{label: fmt.Sprintf("(-inf, %d) union (%d, inf)", a, b), yes: kind == "open", note: "the complement is the closed interval [a, b]"},
		{label: fmt.Sprintf("[%d, %d]", a, b), yes: kind == "closed", note: "a closed interval contains its endpoints and is bounded"},
		{label: fmt.Sprintf("[%d, inf)", a), yes: kind == "closed", note: "a closed ray is closed and unbounded, so it is not open"},
		{label: fmt.Sprintf("{%d, %d}", a, b), yes: kind == "closed", note: "finite sets in a metric space are closed"},
		{label: fmt.Sprintf("[%d, %d] union [%d, %d]", a, b, c, d), yes: kind == "closed", note: "a finite union of closed sets is closed"},
		{label: "Z", yes: kind == "closed", note: "Z has no finite limit points in R, so it is closed and not open"},
		{label: "{1/n} union {0}", yes: kind == "closed", note: "the only limit point is zero, which the set contains"},
		{label: "R", yes: true, note: "R is both open and closed"},
		{label: "the empty set", yes: true, note: "the empty set is both open and closed"},
		{label: fmt.Sprintf("(%d, %d]", a, b), yes: false, note: "it contains a with no open ball inside it and omits the limit point b"},
		{label: fmt.Sprintf("[%d, %d)", a, b), yes: false, note: "it contains a endpoint with no open ball inside it and omits the limit point b"},
		{label: "Q", yes: false, note: "Q has interior empty, so it is not open, and its closure is R, so it is not closed"},
		{label: "the irrational numbers", yes: false, note: "the complement Q is dense and the set has empty interior"},
	}

	return countRAVerdicts(fmt.Sprintf("How many of these subsets of R are %s?", kind), pool)
}

func uniformProblem() generator.Problem {
	return countRAVerdicts("How many of these functions are uniformly continuous on the stated domain?", []raVerdict{
		{label: "x on R", yes: true, note: "a linear function has delta equal to epsilon"},
		{label: "sin x on R", yes: true, note: "the derivative is bounded by one"},
		{label: "cos x on R", yes: true, note: "the derivative is bounded by one"},
		{label: "arctan x on R", yes: true, note: "the derivative 1/(1+x^2) is bounded by one"},
		{label: "x/(1+x^2) on R", yes: true, note: "the derivative is bounded on the whole line"},
		{label: "e^(-x^2) on R", yes: true, note: "the function is continuous and tends to zero at both ends"},
		{label: "sqrt(x) on [0, 1]", yes: true, note: "a continuous function on a compact interval is uniformly continuous"},
		{label: "1/(1+x^2) on R", yes: true, note: "the derivative is bounded on the whole line"},
		{label: "x^2 on R", yes: false, note: "the same delta cannot serve points arbitrarily far apart"},
		{label: "1/x on (0, 1)", yes: false, note: "the function blows up near zero"},
		{label: "sin(1/x) on (0, 1)", yes: false, note: "oscillation near zero defeats any single delta"},
		{label: "e^x on R", yes: false, note: "the derivative is unbounded"},
		{label: "ln x on (0, 1)", yes: false, note: "the function tends to minus infinity at zero"},
		{label: "x sin x on R", yes: false, note: "the derivative is unbounded"},
		{label: "tan x on (-pi/2, pi/2)", yes: false, note: "the function blows up at both endpoints"},
		{label: "1/(x-1) on (1, inf)", yes: false, note: "the function blows up at the missing endpoint one"},
	})
}

func absoluteProblem() generator.Problem {
	return countRAVerdicts("How many of these series converge absolutely?", []raVerdict{
		{label: "sum 1/n^2", yes: true, note: "a p-series with p greater than one converges"},
		{label: "sum (-1)^n/n^2", yes: true, note: "the absolute series is sum 1/n^2"},
		{label: "sum 1/2^n", yes: true, note: "a geometric series with ratio less than one"},
		{label: "sum (-1)^n/2^n", yes: true, note: "the absolute series is geometric with ratio one half"},
		{label: "sum 1/n^3", yes: true, note: "a p-series with p greater than one converges"},
		{label: "sum (-1)^n/(n^2+1)", yes: true, note: "the absolute series is comparable to sum 1/n^2"},
		{label: "sum 1/(n^2+n)", yes: true, note: "it is bounded by sum 1/n^2"},
		{label: "sum (-1)^n/n!", yes: true, note: "the factorial series converges absolutely"},
		{label: "sum 1/n", yes: false, note: "the harmonic series diverges"},
		{label: "sum (-1)^n/n", yes: false, note: "it converges conditionally, not absolutely"},
		{label: "sum 1/sqrt(n)", yes: false, note: "a p-series with p at most one diverges"},
		{label: "sum (-1)^n/sqrt(n)", yes: false, note: "the absolute series is sum 1/sqrt(n), which diverges"},
		{label: "sum 1/(n+1)", yes: false, note: "it is the harmonic series shifted"},
		{label: "sum (-1)^n/(n+1)", yes: false, note: "it converges conditionally, not absolutely"},
		{label: "sum 1/(n ln n)", yes: false, note: "the integral test gives divergence"},
		{label: "sum (-1)^n/ln n", yes: false, note: "the terms do not even tend to zero, so the series diverges"},
	})
}

func uncountableProblem() generator.Problem {
	return countRAVerdicts("How many of these sets are uncountable?", []raVerdict{
		{label: "R", yes: true, note: "Cantor's diagonal argument applies to the whole line"},
		{label: "(0, 1)", yes: true, note: "an open interval is in bijection with the whole line"},
		{label: "[0, 1]", yes: true, note: "a closed interval has the cardinality of the continuum"},
		{label: "the irrational numbers", yes: true, note: "the rationals are countable and the reals are not"},
		{label: "the Cantor set", yes: true, note: "it is in bijection with the power set of the naturals"},
		{label: "the power set of N", yes: true, note: "its cardinality is strictly larger than that of N"},
		{label: "the transcendental numbers", yes: true, note: "the algebraic numbers are countable"},
		{label: "R x R", yes: true, note: "the plane has the same cardinality as the line"},
		{label: "N", yes: false, note: "the naturals can be listed by definition"},
		{label: "Z", yes: false, note: "the integers can be listed as 0, 1, -1, 2, -2, ..."},
		{label: "Q", yes: false, note: "the rationals form a countable set"},
		{label: "the algebraic numbers", yes: false, note: "there are countably many polynomials with rational coefficients"},
		{label: "Q intersect [0, 1]", yes: false, note: "a subset of a countable set is countable"},
		{label: "Z x Z", yes: false, note: "a countable product of countable sets is countable"},
		{label: "{1, 2, 3}", yes: false, note: "a finite set is countable"},
		{label: "Q union {sqrt 2}", yes: false, note: "adding one point to a countable set keeps it countable"},
	})
}

func densityProblem() generator.Problem {
	return countRAVerdicts("How many of these sets are dense in R?", []raVerdict{
		{label: "Q", yes: true, note: "every open interval contains a rational"},
		{label: "the irrational numbers", yes: true, note: "every open interval contains an irrational"},
		{label: "R", yes: true, note: "the whole line meets every open interval"},
		{label: "the dyadic rationals", yes: true, note: "grid points with denominator a power of two fall in every interval"},
		{label: "the terminating decimals", yes: true, note: "they contain the dyadic rationals"},
		{label: "Q minus {0}", yes: true, note: "removing one point leaves a dense set"},
		{label: "the algebraic numbers", yes: true, note: "they contain the rationals"},
		{label: "Q + sqrt 2", yes: true, note: "a translate of a dense set stays dense"},
		{label: "Z", yes: false, note: "the interval (0, 1) contains no integer"},
		{label: "N", yes: false, note: "the interval (0, 1) contains no natural number"},
		{label: "{0, 1, 2}", yes: false, note: "a finite set misses most open intervals"},
		{label: "(0, 1)", yes: false, note: "it misses the interval (2, 3)"},
		{label: "{1/n : n >= 1}", yes: false, note: "it misses the interval (2, 3)"},
		{label: "the Cantor set", yes: false, note: "it misses the middle third of every removed interval"},
		{label: "[0, 1]", yes: false, note: "it misses the interval (2, 3)"},
		{label: "Q intersect [0, 1]", yes: false, note: "its closure is [0, 1], not the whole line"},
	})
}

func cauchyProblem() generator.Problem {
	return countRAVerdicts("How many of these sequences are Cauchy?", []raVerdict{
		{label: "the constant sequence 1", yes: true, note: "all terms are equal"},
		{label: "1/n", yes: true, note: "it converges to zero"},
		{label: "1/2^n", yes: true, note: "it converges to zero"},
		{label: "(-1)^n/n", yes: true, note: "it converges to zero"},
		{label: "the partial sums of 1/k^2", yes: true, note: "the series converges, so its partial sums form a Cauchy sequence"},
		{label: "1/n^2", yes: true, note: "it converges to zero"},
		{label: "(n+1)/n", yes: true, note: "it converges to one"},
		{label: "sin(n)/n", yes: true, note: "it is squeezed between -1/n and 1/n"},
		{label: "(-1)^n", yes: false, note: "successive terms stay two apart"},
		{label: "n", yes: false, note: "the terms grow without bound"},
		{label: "ln n", yes: false, note: "the terms grow without bound"},
		{label: "the partial sums of 1/k", yes: false, note: "the harmonic series diverges"},
		{label: "sin n", yes: false, note: "it has no limit, so it is not Cauchy"},
		{label: "n^2", yes: false, note: "the terms grow without bound"},
		{label: "(-1)^n n", yes: false, note: "successive terms are far apart"},
		{label: "2^n", yes: false, note: "the terms grow without bound"},
	})
}

func continuityProblem() generator.Problem {
	return countRAVerdicts("How many of these functions are continuous at x = 0?", []raVerdict{
		{label: "x", yes: true, note: "a polynomial is continuous everywhere"},
		{label: "x^2", yes: true, note: "a polynomial is continuous everywhere"},
		{label: "sin x", yes: true, note: "the sine function is continuous everywhere"},
		{label: "cos x", yes: true, note: "the cosine function is continuous everywhere"},
		{label: "e^x", yes: true, note: "the exponential is continuous everywhere"},
		{label: "|x|", yes: true, note: "the two one-sided limits both equal zero"},
		{label: "x sin(1/x), with f(0) = 0", yes: true, note: "the function is squeezed by |x|"},
		{label: "1/(1+x^2)", yes: true, note: "the denominator is nonzero at zero"},
		{label: "1/x", yes: false, note: "the function is unbounded near zero"},
		{label: "floor x", yes: false, note: "the left limit is -1 while the right limit is 0"},
		{label: "sign x, with sign(0) = 0", yes: false, note: "the left limit is -1 and the right limit is 1"},
		{label: "sin(1/x), with f(0) = 0", yes: false, note: "the function oscillates near zero"},
		{label: "|x|/x, with f(0) = 0", yes: false, note: "the left limit is -1 and the right limit is 1"},
		{label: "1/x^2", yes: false, note: "the function is unbounded near zero"},
		{label: "cos(1/x), with f(0) = 0", yes: false, note: "the function oscillates near zero"},
		{label: "the step function with f(0) = 1", yes: false, note: "the left limit is 0 while the value is 1"},
	})
}

func nestedIntervalsProblem() generator.Problem {
	return countRAVerdicts("How many of these nested families of intervals have a nonempty intersection?", []raVerdict{
		{label: "I_n = [0, 1/n], n >= 1", yes: true, note: "the intersection is {0}"},
		{label: "I_n = [-1/n, 1/n], n >= 1", yes: true, note: "the intersection is {0}"},
		{label: "I_n = [0, 1 + 1/n], n >= 1", yes: true, note: "the intersection is [0, 1]"},
		{label: "I_n = [-1, 1/n], n >= 1", yes: true, note: "the intersection is [-1, 0]"},
		{label: "I_n = [0, 1], n >= 1", yes: true, note: "every interval is [0, 1]"},
		{label: "I_n = [-1/n, 1 + 1/n], n >= 1", yes: true, note: "the intersection is [0, 1]"},
		{label: "I_n = [0, 1 - 1/n], n >= 1", yes: true, note: "the intersection is {0}"},
		{label: "I_n = [-1 - 1/n, 1/n], n >= 1", yes: true, note: "the intersection is [-1, 0]"},
		{label: "I_n = (0, 1/n), n >= 1", yes: false, note: "no positive number is below every 1/n"},
		{label: "I_n = (0, 1/n], n >= 1", yes: false, note: "zero is excluded from every interval"},
		{label: "I_n = [1/n, 1/n + 1/n^2], n >= 1", yes: false, note: "the intervals shrink towards zero, which they do not contain"},
		{label: "I_n = [n, n + 1], n >= 1", yes: false, note: "the intervals move off to infinity"},
		{label: "I_n = [1/n, 1/(n-1)], n >= 2", yes: false, note: "the intervals shrink towards zero, which they do not contain"},
		{label: "I_n = (1/n, 1/n + 1/n^2), n >= 1", yes: false, note: "the intervals shrink towards zero, which they do not contain"},
		{label: "I_n = [n, infinity), n >= 1", yes: false, note: "the intervals move off to infinity"},
		{label: "I_n = [1/n, 1), n >= 1", yes: false, note: "no positive number is below every 1/n"},
	})
}

func subsequenceProblem() generator.Problem {
	return countRAVerdicts("For the sequence a_n = 1/n, how many of these are subsequences?", []raVerdict{
		{label: "a_1, a_2, a_3, ...", yes: true, note: "this is the sequence itself"},
		{label: "a_2, a_4, a_6, ...", yes: true, note: "the indices are strictly increasing"},
		{label: "a_1, a_3, a_5, ...", yes: true, note: "the indices are strictly increasing"},
		{label: "a_2, a_3, a_4, ...", yes: true, note: "the first term is omitted"},
		{label: "a_1, a_4, a_9, ...", yes: true, note: "the indices are the squares"},
		{label: "a_2, a_3, a_5, a_7, ...", yes: true, note: "the indices are the primes"},
		{label: "a_1, a_2, a_4, a_8, ...", yes: true, note: "the indices are the powers of two"},
		{label: "a_3, a_4, a_5, ...", yes: true, note: "the first two terms are omitted"},
		{label: "0, a_1, a_2, ...", yes: false, note: "zero is not a term of the sequence"},
		{label: "a_1, a_1, a_2, ...", yes: false, note: "the index 1 is repeated"},
		{label: "a_2, a_1, a_3, ...", yes: false, note: "the indices are not increasing"},
		{label: "2, a_1, a_2, ...", yes: false, note: "two is not a term of the sequence"},
		{label: "a_1, a_3, a_2, a_4, ...", yes: false, note: "the indices are not increasing"},
		{label: "a_1, a_2, a_2, a_3, ...", yes: false, note: "a term is repeated"},
		{label: "a_1, a_2, a_4, a_3, ...", yes: false, note: "the indices are not increasing"},
		{label: "a_1, a_2, a_3, ..., 0", yes: false, note: "zero is the limit, not a term"},
	})
}

func completeProblem() generator.Problem {
	return countRAVerdicts("How many of these spaces are complete with the stated metric?", []raVerdict{
		{label: "R with the usual metric", yes: true, note: "every Cauchy sequence of reals converges"},
		{label: "C with the usual metric", yes: true, note: "the complex plane is complete"},
		{label: "R^2 with the Euclidean metric", yes: true, note: "finite-dimensional Euclidean space is complete"},
		{label: "[0, 1] with the usual metric", yes: true, note: "a closed bounded interval is complete"},
		{label: "Z with the usual metric", yes: true, note: "a closed subset of a complete space is complete"},
		{label: "the Cantor set with the usual metric", yes: true, note: "it is closed in R"},
		{label: "a finite set with the discrete metric", yes: true, note: "every Cauchy sequence is eventually constant"},
		{label: "[0, 1] union [2, 3] with the usual metric", yes: true, note: "a finite union of closed sets is closed"},
		{label: "Q with the usual metric", yes: false, note: "a sequence of rationals can converge to an irrational"},
		{label: "(0, 1) with the usual metric", yes: false, note: "the sequence 1/n has no limit in the space"},
		{label: "the irrationals with the usual metric", yes: false, note: "a sequence of irrationals can converge to a rational"},
		{label: "R minus {0} with the usual metric", yes: false, note: "the sequence 1/n has no limit in the space"},
		{label: "[0, 1) with the usual metric", yes: false, note: "the sequence 1 - 1/n has no limit in the space"},
		{label: "Q intersect [0, 1] with the usual metric", yes: false, note: "it is not closed in R"},
		{label: "(0, infinity) with the usual metric", yes: false, note: "the sequence 1/n has no limit in the space"},
		{label: "{1/n : n >= 1} with the usual metric", yes: false, note: "the sequence is Cauchy but its limit zero is missing"},
	})
}

func metricConvergenceProblem() generator.Problem {
	return countRAVerdicts("In the discrete metric, how many of these sequences converge?", []raVerdict{
		{label: "1, 1, 1, 1, ...", yes: true, note: "a constant sequence converges"},
		{label: "0, 0, 0, 0, ...", yes: true, note: "a constant sequence converges"},
		{label: "2, 2, 2, 2, ...", yes: true, note: "a constant sequence converges"},
		{label: "3, 3, 3, 3, ...", yes: true, note: "a constant sequence converges"},
		{label: "1, 2, 2, 2, ...", yes: true, note: "the sequence is eventually constant"},
		{label: "0, 0, 1, 1, ...", yes: true, note: "the sequence is eventually constant"},
		{label: "2, 2, 1, 1, ...", yes: true, note: "the sequence is eventually constant"},
		{label: "3, 3, 3, 1, ...", yes: true, note: "the sequence is eventually constant"},
		{label: "1/n", yes: false, note: "the terms never become constant"},
		{label: "(-1)^n", yes: false, note: "the terms keep changing"},
		{label: "n", yes: false, note: "the terms never become constant"},
		{label: "sin n", yes: false, note: "the terms keep changing"},
		{label: "0, 1, 0, 1, ...", yes: false, note: "the terms keep changing"},
		{label: "1, 1/2, 1/3, ...", yes: false, note: "the terms never become constant"},
		{label: "n mod 2", yes: false, note: "the terms keep changing"},
		{label: "2, 1, 1/2, 1/3, ...", yes: false, note: "the terms never become constant"},
	})
}

func sequentialProblem() generator.Problem {
	return countRAVerdicts("How many of these functions have a limit at x = 0?", []raVerdict{
		{label: "f(x) = x^2", yes: true, note: "a polynomial is continuous at 0"},
		{label: "f(x) = sin x", yes: true, note: "the sine function is continuous at 0"},
		{label: "f(x) = cos x", yes: true, note: "the cosine function is continuous at 0"},
		{label: "f(x) = e^x", yes: true, note: "the exponential function is continuous at 0"},
		{label: "f(x) = |x|", yes: true, note: "the absolute value is continuous at 0"},
		{label: "f(x) = x sin(1/x) for x not 0", yes: true, note: "the squeeze theorem gives limit 0"},
		{label: "f(x) = (e^x - 1)/x for x not 0", yes: true, note: "the limit is 1"},
		{label: "f(x) = sqrt(|x|)", yes: true, note: "the limit is 0"},
		{label: "f(x) = 1/x for x not 0", yes: false, note: "the values escape to infinity near 0"},
		{label: "f(x) = sin(1/x) for x not 0", yes: false, note: "the values oscillate between -1 and 1"},
		{label: "f(x) = cos(1/x) for x not 0", yes: false, note: "the values oscillate between -1 and 1"},
		{label: "f(x) = sign x", yes: false, note: "the left and right limits differ"},
		{label: "f(x) = floor x", yes: false, note: "the left and right limits differ"},
		{label: "f(x) = |x|/x for x not 0", yes: false, note: "the left and right limits differ"},
		{label: "f(x) = 1/x^2 for x not 0", yes: false, note: "the values escape to infinity near 0"},
		{label: "f(x) = 0 for x < 0 and 1 for x >= 0", yes: false, note: "the left and right limits differ"},
	})
}

func ivtProblem() generator.Problem {
	return countRAVerdicts("For how many of these does the intermediate value theorem guarantee a zero from the endpoint signs?", []raVerdict{
		{label: "f(x) = x^2 - 4 on [0, 3]", yes: true, note: "f(0) = -4 and f(3) = 5 have opposite signs"},
		{label: "f(x) = x^3 - 2 on [0, 2]", yes: true, note: "f(0) = -2 and f(2) = 6 have opposite signs"},
		{label: "f(x) = e^x - 2 on [0, 1]", yes: true, note: "f(0) = -1 and f(1) = e - 2 have opposite signs"},
		{label: "f(x) = cos x on [0, 3]", yes: true, note: "cos 0 = 1 and cos 3 is negative"},
		{label: "f(x) = sin x on [1, 4]", yes: true, note: "sin 1 is positive and sin 4 is negative"},
		{label: "f(x) = ln x on [1/2, 2]", yes: true, note: "ln(1/2) is negative and ln 2 is positive"},
		{label: "f(x) = x^5 + x - 1 on [0, 1]", yes: true, note: "f(0) = -1 and f(1) = 1 have opposite signs"},
		{label: "f(x) = x^3 - 3x + 1 on [0, 1]", yes: true, note: "f(0) = 1 and f(1) = -1 have opposite signs"},
		{label: "f(x) = x^2 + 1 on [-1, 1]", yes: false, note: "both endpoint values are positive"},
		{label: "f(x) = e^x on [0, 1]", yes: false, note: "both endpoint values are positive"},
		{label: "f(x) = 1/x on [-1, 1]", yes: false, note: "the function is not continuous on the interval"},
		{label: "f(x) = x^2 - 4 on [2, 3]", yes: false, note: "both endpoint values are positive"},
		{label: "f(x) = sin x on [0, 1]", yes: false, note: "both endpoint values are positive"},
		{label: "f(x) = |x| on [-1, 1]", yes: false, note: "both endpoint values are positive"},
		{label: "f(x) = floor x on [0, 1]", yes: false, note: "the function is not continuous on the interval"},
		{label: "f(x) = x^2 on [-1, 1]", yes: false, note: "both endpoint values are positive"},
	})
}

func evtProblem() generator.Problem {
	return countRAVerdicts("How many of these functions attain both a maximum and a minimum on the stated domain?", []raVerdict{
		{label: "f(x) = x^2 on [0, 1]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = sin x on [0, pi]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = e^x on [0, 1]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = |x| on [-1, 1]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = x^3 - x on [-1, 1]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = cos x on [0, 2]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = 1/(1 + x^2) on [0, 2]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = x on [0, 1]", yes: true, note: "a continuous function on a compact interval attains its bounds"},
		{label: "f(x) = x on (0, 1)", yes: false, note: "the endpoint values are approached but never reached"},
		{label: "f(x) = 1/x on (0, 1]", yes: false, note: "the function is unbounded above"},
		{label: "f(x) = x^2 on R", yes: false, note: "the function is unbounded above"},
		{label: "f(x) = e^x on R", yes: false, note: "the function is unbounded above"},
		{label: "f(x) = sin x on (0, pi)", yes: false, note: "the maximum is approached but never reached"},
		{label: "f(x) = 1/(1 + x^2) on R", yes: false, note: "the minimum is approached but never reached"},
		{label: "f(x) = tan x on (-pi/2, pi/2)", yes: false, note: "the function is unbounded"},
		{label: "f(x) = ln x on (0, 1]", yes: false, note: "the function is unbounded below"},
	})
}

func closureProblem() generator.Problem {
	return countRAVerdicts("How many of these sets have closure equal to [0, 1]?", []raVerdict{
		{label: "(0, 1)", yes: true, note: "the endpoints are limit points"},
		{label: "[0, 1)", yes: true, note: "the missing endpoint 1 is a limit point"},
		{label: "(0, 1]", yes: true, note: "the missing endpoint 0 is a limit point"},
		{label: "Q intersect [0, 1]", yes: true, note: "the rationals are dense in the interval"},
		{label: "the irrationals intersect [0, 1]", yes: true, note: "the irrationals are dense in the interval"},
		{label: "[0, 1] minus {1/2}", yes: true, note: "the removed point is still a limit point"},
		{label: "the dyadic rationals in [0, 1]", yes: true, note: "the dyadic rationals are dense in the interval"},
		{label: "(0, 1) union {0, 1}", yes: true, note: "the set already contains both endpoints"},
		{label: "{1/n : n >= 1}", yes: false, note: "the closure also contains 0"},
		{label: "Z", yes: false, note: "the closure is Z itself"},
		{label: "{0}", yes: false, note: "the closure is {0}"},
		{label: "(0, 1) union (2, 3)", yes: false, note: "the closure also contains 1 and 2"},
		{label: "[0, 1] union {2}", yes: false, note: "the closure contains 2"},
		{label: "the Cantor set", yes: false, note: "the closure is the Cantor set itself"},
		{label: "(0, 1/2)", yes: false, note: "the closure is [0, 1/2]"},
		{label: "[0, 1) union {2}", yes: false, note: "the closure contains 2"},
	})
}

func completenessProblem() generator.Problem {
	return countRAVerdicts("How many of these subsets of R have a supremum in R?", []raVerdict{
		{label: "(0, 1)", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "[0, 1]", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "{1/n : n >= 1}", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "Q intersect [0, 1]", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "{x : x^2 < 2}", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "(0, 1) union {2}", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "[0, sqrt 2]", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "{1 - 1/n : n >= 1}", yes: true, note: "bounded above, so completeness gives a supremum"},
		{label: "N", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "Z", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "{2^n : n >= 0}", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "[0, infinity)", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "R", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "(0, infinity)", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "{n + 1/n : n >= 1}", yes: false, note: "unbounded above, so it has no supremum"},
		{label: "Q", yes: false, note: "unbounded above, so it has no supremum"},
	})
}

func bolzanoProblem() generator.Problem {
	return countRAVerdicts("How many of these sequences have a convergent subsequence?", []raVerdict{
		{label: "(-1)^n", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "sin n", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "cos n", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "1/n", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "(-1)^n/n", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "n mod 3", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "sin(n pi/4)", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "1 + (-1)^n/n", yes: true, note: "bounded, so Bolzano-Weierstrass applies"},
		{label: "n", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "(-1)^n n", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "2^n", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "n^2", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "ln n", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "n + sin n", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "(-2)^n", yes: false, note: "unbounded, with no convergent subsequence"},
		{label: "n!", yes: false, note: "unbounded, with no convergent subsequence"},
	})
}

func seriesTestsProblem() generator.Problem {
	return countRAVerdicts("How many of these series converge?", []raVerdict{
		{label: "sum 1/n^2", yes: true, note: "a p-series with p = 2"},
		{label: "sum 1/2^n", yes: true, note: "a geometric series with ratio 1/2"},
		{label: "sum (-1)^n/n", yes: true, note: "the alternating harmonic series"},
		{label: "sum n/2^n", yes: true, note: "the ratio test gives a limit of 1/2"},
		{label: "sum 1/(n^2 + n)", yes: true, note: "it telescopes to 1"},
		{label: "sum (-1)^n/n^2", yes: true, note: "it converges absolutely"},
		{label: "sum 1/n!", yes: true, note: "the ratio test gives a limit of 0"},
		{label: "sum 1/(n (ln n)^2)", yes: true, note: "the integral test with u = ln n converges"},
		{label: "sum 1/n", yes: false, note: "the harmonic series diverges"},
		{label: "sum 1/sqrt n", yes: false, note: "a p-series with p = 1/2"},
		{label: "sum (-1)^n", yes: false, note: "the terms do not tend to zero"},
		{label: "sum n/(n + 1)", yes: false, note: "the terms do not tend to zero"},
		{label: "sum 1/ln n", yes: false, note: "the terms do not tend to zero"},
		{label: "sum 2^n", yes: false, note: "a geometric series with ratio 2"},
		{label: "sum n", yes: false, note: "the terms do not tend to zero"},
		{label: "sum 1/(n ln n)", yes: false, note: "the integral test with u = ln n diverges"},
	})
}

func rearrangementProblem() generator.Problem {
	return countRAVerdicts("For how many of these series does every rearrangement converge to the same sum?", []raVerdict{
		{label: "sum 1/n^2", yes: true, note: "the series converges absolutely"},
		{label: "sum 1/2^n", yes: true, note: "the series converges absolutely"},
		{label: "sum (-1)^n/n^2", yes: true, note: "the series converges absolutely"},
		{label: "sum 1/n^3", yes: true, note: "the series converges absolutely"},
		{label: "sum 1/n!", yes: true, note: "the series converges absolutely"},
		{label: "sum (-1)^n/2^n", yes: true, note: "the series converges absolutely"},
		{label: "sum 1/(n^2 + 1)", yes: true, note: "the series converges absolutely"},
		{label: "sum 1/(n^2 + n)", yes: true, note: "the series converges absolutely"},
		{label: "sum (-1)^n/n", yes: false, note: "conditionally convergent, so Riemann's theorem applies"},
		{label: "sum 1/n", yes: false, note: "the series diverges"},
		{label: "sum 1/sqrt n", yes: false, note: "the series diverges"},
		{label: "sum (-1)^n/sqrt n", yes: false, note: "conditionally convergent, so Riemann's theorem applies"},
		{label: "sum 1/(n ln n)", yes: false, note: "the series diverges"},
		{label: "sum (-1)^n/ln n", yes: false, note: "conditionally convergent, so Riemann's theorem applies"},
		{label: "sum 1/(n + 1)", yes: false, note: "the series diverges"},
		{label: "sum (-1)^n/(n + 1)", yes: false, note: "conditionally convergent, so Riemann's theorem applies"},
	})
}

func compactStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about compactness are true?", []raVerdict{
		{label: "a subset of R is compact exactly when it is closed and bounded", yes: true, note: "the Heine-Borel theorem"},
		{label: "the continuous image of a compact set is compact", yes: true, note: "compactness is preserved by continuous maps"},
		{label: "a closed subset of a compact space is compact", yes: true, note: "closed subsets inherit every finite subcover"},
		{label: "a finite union of compact sets is compact", yes: true, note: "combine a finite subcover from each piece"},
		{label: "a continuous function on a compact set is uniformly continuous", yes: true, note: "compactness upgrades continuity"},
		{label: "a compact subset of R attains its maximum and minimum", yes: true, note: "the extreme value theorem"},
		{label: "a compact metric space is complete", yes: true, note: "every Cauchy sequence has a convergent subsequence"},
		{label: "every sequence in a compact metric space has a convergent subsequence", yes: true, note: "sequential compactness"},
		{label: "the open interval (0, 1) is compact", yes: false, note: "it is not closed"},
		{label: "every compact subset of R is open", yes: false, note: "compact subsets of R are closed"},
		{label: "Q intersect [0, 1] is compact", yes: false, note: "it is not closed in R"},
		{label: "an infinite compact set must be countable", yes: false, note: "[0, 1] is uncountable and compact"},
		{label: "every compact set is connected", yes: false, note: "a two-point set is compact and disconnected"},
		{label: "[0, infinity) is compact", yes: false, note: "it is unbounded"},
		{label: "every bounded subset of R is compact", yes: false, note: "Q intersect [0, 1] is bounded but not compact"},
		{label: "compact subsets of R are dense in R", yes: false, note: "a compact subset of R is closed and bounded"},
	})
}

func connectedStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about connectedness are true?", []raVerdict{
		{label: "the continuous image of a connected set is connected", yes: true, note: "connectedness is preserved by continuous maps"},
		{label: "R with the usual topology is connected", yes: true, note: "R is an interval"},
		{label: "every interval in R is connected", yes: true, note: "the connected subsets of R are exactly the intervals"},
		{label: "a connected subset of R is an interval", yes: true, note: "this is the interval theorem"},
		{label: "the union of connected sets that share a point is connected", yes: true, note: "the shared point joins the pieces"},
		{label: "the closure of a connected set is connected", yes: true, note: "closure preserves connectedness"},
		{label: "every path-connected space is connected", yes: true, note: "a path joins any two points"},
		{label: "a connected subset of R with more than one point is uncountable", yes: true, note: "it contains an interval"},
		{label: "Q is connected", yes: false, note: "Q is totally disconnected"},
		{label: "a two-point set is connected", yes: false, note: "each point is a separate component"},
		{label: "every connected space is compact", yes: false, note: "R is connected and not compact"},
		{label: "every connected space is path-connected", yes: false, note: "the topologist's sine curve is a counterexample"},
		{label: "the Cantor set is connected", yes: false, note: "the Cantor set is totally disconnected"},
		{label: "R minus {0} is connected", yes: false, note: "it has two components"},
		{label: "every connected subset of R is open", yes: false, note: "[0, 1] is connected and not open"},
		{label: "the union of two disjoint connected sets is connected", yes: false, note: "the two pieces separate the union"},
	})
}

func discontinuityStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about discontinuities are true?", []raVerdict{
		{label: "a removable discontinuity has a finite two-sided limit", yes: true, note: "the limit exists but the value differs or is undefined"},
		{label: "a jump discontinuity has unequal one-sided limits", yes: true, note: "both one-sided limits exist and differ"},
		{label: "sin(1/x) at 0 has an essential discontinuity when f(0) = 0", yes: true, note: "the function oscillates without a limit"},
		{label: "a monotone function has only jump discontinuities", yes: true, note: "one-sided limits always exist"},
		{label: "the discontinuity set of a monotone function is at most countable", yes: true, note: "each jump contains a disjoint rational interval"},
		{label: "1/x at 0 has an infinite discontinuity", yes: true, note: "the function is unbounded near zero"},
		{label: "the floor function is discontinuous at every integer", yes: true, note: "the one-sided limits differ by one"},
		{label: "a removable discontinuity can be removed by changing one value", yes: true, note: "define the value to equal the limit"},
		{label: "every discontinuity is removable", yes: false, note: "jumps and essential discontinuities are not"},
		{label: "a monotone function has finitely many discontinuities", yes: false, note: "the set can be countably infinite"},
		{label: "a jump discontinuity is essential", yes: false, note: "a jump has both one-sided limits"},
		{label: "a function discontinuous at a point cannot have a limit there", yes: false, note: "a removable discontinuity has a limit"},
		{label: "sin(1/x) at 0 has a removable discontinuity", yes: false, note: "the limit does not exist"},
		{label: "the discontinuities of a monotone function are uncountable", yes: false, note: "they are at most countable"},
		{label: "1/x has a jump discontinuity at 0", yes: false, note: "the discontinuity is infinite"},
		{label: "the floor function is continuous at every integer", yes: false, note: "it jumps at each integer"},
	})
}

func monotoneStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about monotone functions are true?", []raVerdict{
		{label: "one-sided limits exist at every interior point", yes: true, note: "monotonicity bounds each side"},
		{label: "the set of discontinuities is at most countable", yes: true, note: "each jump contains a disjoint rational interval"},
		{label: "a monotone function on a closed interval has bounded variation", yes: true, note: "the variation equals the endpoint difference"},
		{label: "a strictly increasing function with the intermediate value property is continuous", yes: true, note: "a jump would omit values"},
		{label: "the inverse of a continuous strictly increasing function is continuous", yes: true, note: "strict monotonicity preserves order"},
		{label: "a monotone function is differentiable almost everywhere", yes: true, note: "Lebesgue's theorem"},
		{label: "a monotone function has at most countably many jumps", yes: true, note: "the jump set is at most countable"},
		{label: "the sum of two increasing functions is increasing", yes: true, note: "add the two inequalities"},
		{label: "every monotone function is continuous", yes: false, note: "the floor function is monotone and discontinuous"},
		{label: "a monotone function has only finitely many discontinuities", yes: false, note: "the set can be countably infinite"},
		{label: "every monotone function is differentiable everywhere", yes: false, note: "it is differentiable only almost everywhere"},
		{label: "the inverse of a monotone function is always continuous", yes: false, note: "continuity requires an interval domain or a continuous function"},
		{label: "a monotone function on a closed interval is bounded only if it is continuous", yes: false, note: "a monotone function on a closed interval is always bounded"},
		{label: "every strictly increasing function has the intermediate value property", yes: false, note: "a jump can omit values"},
		{label: "a monotone function can have an essential discontinuity", yes: false, note: "one-sided limits always exist"},
		{label: "every monotone function is uniformly continuous", yes: false, note: "x^2 is increasing on [0, infinity) but not uniformly continuous"},
	})
}

func extendedStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about the extended real line are true?", []raVerdict{
		{label: "the extended real line adds plus infinity and minus infinity", yes: true, note: "the two infinite endpoints"},
		{label: "every subset of the extended reals has a supremum", yes: true, note: "the order is complete"},
		{label: "the extended real line is order complete", yes: true, note: "every nonempty subset has a least upper bound"},
		{label: "the supremum of a set unbounded above is plus infinity", yes: true, note: "the convention absorbs unbounded sets"},
		{label: "the infimum of a set unbounded below is minus infinity", yes: true, note: "the convention absorbs unbounded sets"},
		{label: "the extended real line is compact", yes: true, note: "the two-point compactification of R is compact"},
		{label: "infinity plus a finite number is infinity", yes: true, note: "the infinite endpoints absorb finite shifts"},
		{label: "a monotone sequence in the extended reals converges to its supremum", yes: true, note: "monotone convergence holds in the extended order"},
		{label: "the extended real line is a field", yes: false, note: "the arithmetic of infinity is not a field"},
		{label: "infinity minus infinity is zero", yes: false, note: "this expression is undefined"},
		{label: "zero times infinity is zero", yes: false, note: "this expression is undefined"},
		{label: "the extended real line is countable", yes: false, note: "it contains all real numbers"},
		{label: "every subset of the extended reals has a maximum", yes: false, note: "the supremum need not belong to the set"},
		{label: "infinity divided by infinity is one", yes: false, note: "this expression is undefined"},
		{label: "the extended reals contain a largest finite real", yes: false, note: "the finite reals are unbounded"},
		{label: "the extended real line is disconnected", yes: false, note: "it is a connected compact interval"},
	})
}

func cauchyCriterionStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about the Cauchy criterion for series are true?", []raVerdict{
		{label: "a series converges exactly when its tail sums can be made arbitrarily small", yes: true, note: "the partial sums form a Cauchy sequence"},
		{label: "the harmonic series diverges even though its terms tend to zero", yes: true, note: "blocks of equal length add a fixed amount"},
		{label: "terms tending to zero is necessary but not sufficient", yes: true, note: "the harmonic series is the standard counterexample"},
		{label: "the criterion is phrased with differences of partial sums", yes: true, note: "a block of consecutive terms must have small sum"},
		{label: "completeness makes the Cauchy criterion equivalent to convergence", yes: true, note: "Cauchy sequences converge in a complete space"},
		{label: "an absolutely convergent series satisfies the Cauchy criterion", yes: true, note: "absolute convergence implies convergence"},
		{label: "the criterion applies in every complete metric space", yes: true, note: "the proof uses only completeness"},
		{label: "a series whose tail sums do not tend to zero diverges", yes: true, note: "the contrapositive of the criterion"},
		{label: "terms tending to zero is sufficient for convergence", yes: false, note: "the harmonic series diverges"},
		{label: "the harmonic series converges", yes: false, note: "its partial sums grow without bound"},
		{label: "the criterion only applies to positive series", yes: false, note: "it applies to every series in a complete space"},
		{label: "a Cauchy sequence of partial sums may fail to converge in R", yes: false, note: "R is complete"},
		{label: "the criterion is about the ratio of consecutive terms", yes: false, note: "that is the ratio test"},
		{label: "bounded partial sums imply convergence", yes: false, note: "the partial sums of 1/n are unbounded but boundedness alone is not enough for other series"},
		{label: "the Cauchy criterion is equivalent to absolute convergence", yes: false, note: "conditional convergence still satisfies the criterion"},
		{label: "every series with bounded terms converges", yes: false, note: "the terms of 1/n are bounded and the series diverges"},
	})
}

func metricCompactStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about compact metric spaces are true?", []raVerdict{
		{label: "compact is equivalent to sequentially compact in metric spaces", yes: true, note: "sequences characterize compactness in metric spaces"},
		{label: "a compact metric space is complete", yes: true, note: "every Cauchy sequence has a convergent subsequence"},
		{label: "a compact metric space is totally bounded", yes: true, note: "finitely many balls of every radius cover it"},
		{label: "a continuous real function on a compact metric space attains its extrema", yes: true, note: "the extreme value theorem"},
		{label: "a closed subset of a compact metric space is compact", yes: true, note: "closed subsets inherit compactness"},
		{label: "the continuous image of a compact metric space is compact", yes: true, note: "compactness is preserved by continuous maps"},
		{label: "a compact metric space is separable", yes: true, note: "total boundedness gives a countable dense set"},
		{label: "a compact metric space is bounded", yes: true, note: "finitely many balls of radius one cover it"},
		{label: "every complete metric space is compact", yes: false, note: "R is complete and not compact"},
		{label: "every bounded metric space is compact", yes: false, note: "(0, 1) is bounded and not compact"},
		{label: "a compact metric space must be finite", yes: false, note: "[0, 1] is compact and infinite"},
		{label: "Q intersect [0, 1] is compact", yes: false, note: "it is not closed in R"},
		{label: "compact metric spaces are discrete", yes: false, note: "[0, 1] is compact and not discrete"},
		{label: "a compact metric space is isometric to an interval", yes: false, note: "a circle is compact and not an interval"},
		{label: "every closed and bounded subset of a metric space is compact", yes: false, note: "this holds in R^n but not in every metric space"},
		{label: "compactness depends only on the cardinality", yes: false, note: "the topology matters"},
	})
}

func metricContinuityStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about continuity in metric spaces are true?", []raVerdict{
		{label: "continuity at a point is equivalent to preserving convergent sequences", yes: true, note: "the sequential criterion"},
		{label: "continuity is equivalent to preimages of open sets being open", yes: true, note: "the open-set characterization"},
		{label: "continuity is equivalent to preimages of closed sets being closed", yes: true, note: "take complements of open sets"},
		{label: "a Lipschitz map is continuous", yes: true, note: "choose delta equal to epsilon over the Lipschitz constant"},
		{label: "an isometry is continuous", yes: true, note: "distances are preserved exactly"},
		{label: "the composition of continuous maps is continuous", yes: true, note: "continuity is preserved by composition"},
		{label: "a continuous map on a compact domain is uniformly continuous", yes: true, note: "compactness upgrades continuity"},
		{label: "a continuous map from a compact space to a Hausdorff space is closed", yes: true, note: "the image of a closed set is compact and hence closed"},
		{label: "every continuous map is Lipschitz", yes: false, note: "sqrt x on [0, 1] is continuous but not Lipschitz"},
		{label: "continuity implies uniform continuity on every metric space", yes: false, note: "x^2 on R is continuous but not uniformly continuous"},
		{label: "a continuous map sends open sets to open sets", yes: false, note: "only open maps have that property"},
		{label: "continuity is equivalent to boundedness", yes: false, note: "x is continuous and unbounded on R"},
		{label: "every continuous map is an isometry", yes: false, note: "x^2 is continuous and not an isometry"},
		{label: "a continuous map sends compact sets to closed sets in every codomain", yes: false, note: "the image is compact and closed only when the codomain is Hausdorff"},
		{label: "continuity at a point requires the function to be bounded", yes: false, note: "a continuous function need not be bounded near a point"},
		{label: "a discontinuous map cannot have a limit at a point", yes: false, note: "a removable discontinuity has a limit"},
	})
}

func darbouxStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about Darboux's theorem are true?", []raVerdict{
		{label: "every derivative has the intermediate value property", yes: true, note: "this is Darboux's theorem"},
		{label: "a derivative cannot have a jump discontinuity", yes: true, note: "a jump would omit intermediate values"},
		{label: "a derivative attains every value between two of its values", yes: true, note: "the intermediate value property"},
		{label: "the Darboux property does not imply continuity", yes: true, note: "a derivative may be discontinuous"},
		{label: "a derivative has the Darboux property even when it is discontinuous", yes: true, note: "Darboux's theorem needs only differentiability of the primitive"},
		{label: "if f is differentiable on [a, b] then f' takes every value between f'(a) and f'(b)", yes: true, note: "apply Darboux on the interval"},
		{label: "a derivative can be discontinuous but not with a jump", yes: true, note: "the discontinuity must be essential"},
		{label: "Darboux's theorem applies to derivatives rather than to arbitrary functions", yes: true, note: "the theorem is specific to derivatives"},
		{label: "every continuous function is a derivative", yes: false, note: "not every continuous function has a primitive in elementary terms"},
		{label: "every derivative is continuous", yes: false, note: "x^2 sin(1/x) has a discontinuous derivative at zero"},
		{label: "a derivative can have a jump discontinuity", yes: false, note: "Darboux's theorem forbids jumps"},
		{label: "the Darboux property implies differentiability", yes: false, note: "the property says nothing about a derivative existing"},
		{label: "a function with the Darboux property must be continuous", yes: false, note: "a discontinuous derivative still has the property"},
		{label: "derivatives are always monotone", yes: false, note: "the derivative of sin x changes sign"},
		{label: "Darboux's theorem says every derivative is Riemann integrable", yes: false, note: "a derivative need not be Riemann integrable"},
		{label: "every function with the intermediate value property is a derivative", yes: false, note: "the property alone does not produce a primitive"},
	})
}

func uniformConvStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about uniform convergence are true?", []raVerdict{
		{label: "uniform convergence implies pointwise convergence", yes: true, note: "one bound works for every point"},
		{label: "the uniform limit of continuous functions is continuous", yes: true, note: "the standard interchange theorem"},
		{label: "a uniformly convergent sequence of bounded functions is uniformly bounded", yes: true, note: "the tail is bounded by one constant"},
		{label: "uniform convergence lets a limit pass through a Riemann integral on [a, b]", yes: true, note: "the integrals converge uniformly"},
		{label: "the Cauchy criterion characterizes uniform convergence", yes: true, note: "the sup norm makes the function space complete"},
		{label: "the Weierstrass M-test gives uniform convergence of a series", yes: true, note: "dominate the terms by summable constants"},
		{label: "uniform convergence preserves continuity of the limit", yes: true, note: "the uniform limit theorem"},
		{label: "a uniformly Cauchy sequence of functions has a uniform limit", yes: true, note: "completeness of the sup norm"},
		{label: "pointwise convergence implies uniform convergence", yes: false, note: "x^n on [0, 1] converges pointwise but not uniformly"},
		{label: "the uniform limit of continuous functions can be discontinuous", yes: false, note: "uniform limits preserve continuity"},
		{label: "uniform convergence preserves differentiability", yes: false, note: "the limit of derivatives need not exist"},
		{label: "termwise differentiation is valid under uniform convergence alone", yes: false, note: "it needs uniform convergence of the derivatives"},
		{label: "the pointwise limit of continuous functions is continuous", yes: false, note: "x^n on [0, 1] has a discontinuous pointwise limit"},
		{label: "uniform convergence of a series is equivalent to absolute convergence", yes: false, note: "the two notions are independent"},
		{label: "uniform convergence is the same as pointwise convergence on compact sets", yes: false, note: "pointwise convergence on compact sets is weaker"},
		{label: "uniform convergence always allows termwise integration over an unbounded interval", yes: false, note: "the interval must have finite length"},
	})
}

func interchangeStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about interchanging limits are true?", []raVerdict{
		{label: "uniform convergence justifies swapping a limit and an integral on a closed interval", yes: true, note: "the uniform limit theorem for integrals"},
		{label: "uniform convergence of derivatives justifies swapping a limit and a derivative", yes: true, note: "the C1 uniform convergence theorem"},
		{label: "continuity of a function of two variables justifies swapping two limits under the right hypotheses", yes: true, note: "the double limit theorem"},
		{label: "uniform convergence justifies swapping a limit and an infinite sum", yes: true, note: "the series is uniformly Cauchy"},
		{label: "a uniformly convergent sequence can be integrated term by term", yes: true, note: "the integral of the limit is the limit of the integrals"},
		{label: "uniform convergence on compact sets lets a limit pass through a continuous functional", yes: true, note: "continuous functionals preserve uniform limits"},
		{label: "the dominated convergence theorem justifies swapping a limit and an integral", yes: true, note: "a dominating integrable function controls the tail"},
		{label: "the monotone convergence theorem justifies swapping a limit and an integral", yes: true, note: "an increasing sequence of nonnegative functions"},
		{label: "pointwise convergence alone justifies swapping a limit and an integral", yes: false, note: "spikes can escape to infinity"},
		{label: "pointwise convergence alone justifies swapping a limit and a derivative", yes: false, note: "the derivatives need uniform convergence"},
		{label: "the order of two limits can always be swapped", yes: false, note: "the iterated limits need not agree"},
		{label: "continuity alone justifies swapping a limit and an integral over an unbounded interval", yes: false, note: "integrability of the limit must be controlled"},
		{label: "a pointwise limit of continuous functions can always be integrated term by term", yes: false, note: "the limit need not be integrable"},
		{label: "uniform convergence is unnecessary for termwise differentiation", yes: false, note: "uniform convergence of derivatives is the standard hypothesis"},
		{label: "the double limit always equals the iterated limit", yes: false, note: "the iterated limits may differ"},
		{label: "every interchange of limit and sum is valid for a convergent series", yes: false, note: "uniform convergence is needed"},
	})
}

func baireStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about Baire category are true?", []raVerdict{
		{label: "every complete metric space is a Baire space", yes: true, note: "the Baire category theorem"},
		{label: "R is not a countable union of nowhere dense sets", yes: true, note: "R is complete"},
		{label: "Q is meager in R", yes: true, note: "Q is a countable union of singleton sets"},
		{label: "a countable intersection of dense open sets is dense in a complete metric space", yes: true, note: "the standard formulation"},
		{label: "the Baire category theorem holds for locally compact Hausdorff spaces", yes: true, note: "the second formulation"},
		{label: "a complete metric space is not a countable union of closed sets with empty interior", yes: true, note: "the closed-set formulation"},
		{label: "the complement of a meager set is comeager", yes: true, note: "this is the definition of comeager"},
		{label: "the irrationals are comeager in R", yes: true, note: "their complement Q is meager"},
		{label: "every metric space is a Baire space", yes: false, note: "Q with the usual metric is not Baire"},
		{label: "Q is a Baire space", yes: false, note: "Q is meager in itself"},
		{label: "a countable union of nowhere dense sets always has empty interior", yes: false, note: "Q is meager and dense in R"},
		{label: "R is meager in itself", yes: false, note: "R is a Baire space"},
		{label: "the irrationals are meager in R", yes: false, note: "the irrationals are comeager"},
		{label: "every Baire space is complete", yes: false, note: "local compactness also gives Baire spaces"},
		{label: "a countable union of closed sets is closed", yes: false, note: "countable unions of closed sets need not be closed"},
		{label: "the Baire category theorem requires compactness", yes: false, note: "completeness or local compactness is enough"},
	})
}

func equivalentStatementsProblem() generator.Problem {
	return countRAVerdicts("How many of these statements about equivalent metrics are true?", []raVerdict{
		{label: "equivalent metrics induce the same topology", yes: true, note: "each open ball contains a ball of the other metric"},
		{label: "uniformly equivalent metrics preserve uniform continuity", yes: true, note: "the identity map is uniformly continuous both ways"},
		{label: "Lipschitz equivalent metrics are uniformly equivalent", yes: true, note: "a linear bound controls the uniformity"},
		{label: "equivalent metrics have the same convergent sequences", yes: true, note: "convergence is a topological property"},
		{label: "equivalent metrics have the same continuous functions", yes: true, note: "continuity depends only on the topology"},
		{label: "equivalent metrics have the same compact sets", yes: true, note: "compactness is a topological property"},
		{label: "equivalent metrics have the same open sets", yes: true, note: "the topologies coincide"},
		{label: "the usual metric and the discrete metric on N are topologically equivalent", yes: true, note: "both give the discrete topology on N"},
		{label: "equivalent metrics always have the same Cauchy sequences", yes: false, note: "that needs uniform equivalence"},
		{label: "equivalent metrics always preserve boundedness", yes: false, note: "d and d/(1+d) are equivalent but change boundedness"},
		{label: "equivalent metrics are always Lipschitz equivalent", yes: false, note: "topological equivalence is weaker"},
		{label: "equivalent metrics always have the same completeness", yes: false, note: "completeness is not a topological property"},
		{label: "equivalent metrics always have the same diameter", yes: false, note: "the diameter is not topological"},
		{label: "every metric is equivalent to the discrete metric", yes: false, note: "R with its usual metric is not discrete"},
		{label: "equivalent metrics have identical distance values", yes: false, note: "only the induced topology must agree"},
		{label: "uniformly equivalent metrics are always Lipschitz equivalent", yes: false, note: "the modulus need not be linear"},
	})
}
