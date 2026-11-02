package realanalysis

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The 2026 content expansion adds ten topics that the course never reached:
// order-theoretic tools, measure-zero sets, total boundedness, the classical
// convergence tests for function series and the approximation theorems. Each
// branch asks for a number that can be checked without a proof, and each of the
// three knowledge points has its own question shape.
var raBridgeKinds = map[string]string{
	"ra.order.zorn":               "zorn",
	"ra.order.schroder_bernstein": "schroder_bernstein",
	"ra.measure.null":             "measure_null",
	"ra.measure.cantor_set":       "cantor_set",
	"ra.metric.total_bounded":     "total_bounded",
	"ra.seq.weierstrass_m":        "weierstrass_m",
	"ra.deriv.bounded_variation":  "bounded_variation",
	"ra.deriv.riemann_stieltjes":  "riemann_stieltjes",
	"ra.deriv.arzela_ascoli":      "arzela_ascoli",
	"ra.deriv.weierstrass_approx": "weierstrass_approx",
}

type raBridgeGen struct {
	kind string
}

func raBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func raBridgeProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer},
		KPOrdinal:   o,
	}
}

func raBridgeDecimal(v float64, places int) string {
	s := strconv.FormatFloat(v, 'f', places, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "-0" {
		return "0"
	}
	return s
}

func raBridgeFactorial(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func raBridgeDistinctPrimes(n int) int {
	count := 0
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			count++
			for n%p == 0 {
				n /= p
			}
		}
	}
	if n > 1 {
		count++
	}
	return count
}

// raBridgeItem is one statement in a count question. Exactly half of every
// pool is true, and raBridgeCount draws the answer before it chooses items, so
// every count in 0..draw is reachable with equal probability.
type raBridgeItem struct {
	label string
	yes   bool
	note  string
}

func raBridgeCount(o int, lead string, pool []raBridgeItem, draw int) generator.Problem {
	var yes, no []raBridgeItem
	for _, it := range pool {
		if it.yes {
			yes = append(yes, it)
		} else {
			no = append(no, it)
		}
	}
	count := rand.Intn(draw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if draw-count > len(no) {
		count = draw - len(no)
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := make([]raBridgeItem, 0, draw)
	chosen = append(chosen, yes[:count]...)
	chosen = append(chosen, no[:draw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	parts := make([]string, len(chosen))
	notes := make([]string, len(chosen))
	for i, it := range chosen {
		parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, it.label)
		verdict := "no"
		if it.yes {
			verdict = "yes"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s", it.label, verdict, it.note)
	}
	q := fmt.Sprintf("%s\n%s\n(enter a number)", lead, strings.Join(parts, "\n"))
	e := fmt.Sprintf("Exactly %d of the %d statements are true. %s", count, draw, strings.Join(notes, " "))
	return raBridgeProblem(o, q, strconv.Itoa(count), e)
}

func (g *raBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := raBridgeOrdinal(ctx)
	switch g.kind {
	case "zorn":
		return zornProblem(o)
	case "schroder_bernstein":
		return schroderBernsteinProblem(o)
	case "measure_null":
		return measureNullProblem(o)
	case "cantor_set":
		return cantorSetProblem(o)
	case "total_bounded":
		return totalBoundedProblem(o)
	case "weierstrass_m":
		return weierstrassMProblem(o)
	case "bounded_variation":
		return boundedVariationProblem(o)
	case "riemann_stieltjes":
		return riemannStieltjesProblem(o)
	case "arzela_ascoli":
		return arzelaAscoliProblem(o)
	case "weierstrass_approx":
		return weierstrassApproxProblem(o)
	default:
		return raBridgeProblem(o, "What is 2 + 2? (enter a number)", "4", "The sum is 4.")
	}
}

func zornProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 12 + rand.Intn(109)
		k := raBridgeDistinctPrimes(n)
		return raBridgeProblem(o, fmt.Sprintf("Consider the proper divisors of %d ordered by divisibility. How many maximal elements does this poset have? (enter a number)", n), strconv.Itoa(k), fmt.Sprintf("A proper divisor is maximal when it is %d divided by one of its distinct prime factors, so there are %d maximal elements.", n, k))
	case 2:
		n := 3 + rand.Intn(5)
		ans := raBridgeFactorial(n)
		return raBridgeProblem(o, fmt.Sprintf("How many maximal chains run from the empty set to a full set of %d elements in its power set, adding one element at each step? (enter a number)", n), strconv.Itoa(ans), fmt.Sprintf("A maximal chain is an ordering of the %d elements, so there are %d! = %d chains.", n, n, ans))
	default:
		n := 3 + rand.Intn(10)
		ans := n + 1
		return raBridgeProblem(o, fmt.Sprintf("A maximal chain in the power set of an %d-element set runs from the empty set to the full set, adding one element at each step. How many sets does the chain contain? (enter a number)", n), strconv.Itoa(ans), fmt.Sprintf("The chain contains the empty set, the full set and one set at each of the %d intermediate sizes, so %d sets in total.", n-1, ans))
	}
}

func schroderBernsteinProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 2 + rand.Intn(9)
		ans := 1 << n
		return raBridgeProblem(o, fmt.Sprintf("How many subsets does a set with %d elements have? (enter a number)", n), strconv.Itoa(ans), fmt.Sprintf("Each of the %d elements is either in or out, so there are 2^%d = %d subsets.", n, n, ans))
	case 2:
		n := 4 + rand.Intn(6)
		k := 1 + rand.Intn(n)
		ans := 1
		for i := 0; i < k; i++ {
			ans *= n - i
		}
		return raBridgeProblem(o, fmt.Sprintf("How many injections are there from a %d-element set into an %d-element set? (enter a number)", k, n), strconv.Itoa(ans), fmt.Sprintf("Choose distinct images in order: %d choices, then %d, and so on for %d steps, giving %d injections.", n, n-1, k, ans))
	default:
		n := 2 + rand.Intn(7)
		ans := raBridgeFactorial(n)
		return raBridgeProblem(o, fmt.Sprintf("How many bijections are there between two sets that each have %d elements? (enter a number)", n), strconv.Itoa(ans), fmt.Sprintf("A bijection is a permutation of the %d elements, so there are %d! = %d bijections.", n, n, ans))
	}
}

func measureNullProblem(o int) generator.Problem {
	switch o {
	case 1:
		a := rand.Intn(6)
		b := a + 1 + rand.Intn(10)
		ans := b - a
		return raBridgeProblem(o, fmt.Sprintf("What is the Lebesgue measure of the interval [%d, %d]? (enter a number)", a, b), strconv.Itoa(ans), fmt.Sprintf("The measure of a closed interval is its length, and %d - %d = %d.", b, a, ans))
	case 2:
		pool := []raBridgeItem{
			{"a finite set", true, "finite sets are countable and have measure zero"},
			{"a countable set", true, "countable sets are null"},
			{"the Cantor set", true, "its total length is zero"},
			{"{1/n : n >= 1} together with 0", true, "a countable set"},
			{"the rationals", true, "Q is countable"},
			{"the algebraic numbers", true, "a countable set"},
			{"the integers", true, "Z is countable"},
			{"a two-point set", true, "a finite set"},
			{"the real line", false, "it has infinite measure"},
			{"the interval [0, 1]", false, "it has measure 1"},
			{"the interval (0, 1)", false, "it has measure 1"},
			{"the irrationals in [0, 1]", false, "they have measure 1"},
			{"the complement of the Cantor set in [0, 1]", false, "it has measure 1"},
			{"[0, 1] union [2, 3]", false, "it has measure 2"},
			{"a fat Cantor set", false, "its measure is positive"},
			{"the interval [0, 2]", false, "it has measure 2"},
		}
		return raBridgeCount(o, "How many of these sets have Lebesgue measure zero?", pool, 8)
	default:
		k := 1 + rand.Intn(5)
		n := 2 + rand.Intn(8)
		return raBridgeProblem(o, fmt.Sprintf("A set is formed by removing %d distinct points from the interval [0, %d]. What is the Lebesgue measure of the remaining set? (enter a number)", k, n), strconv.Itoa(n), fmt.Sprintf("Removing finitely many points changes no measure, so the remaining set still has measure %d.", n))
	}
}

func cantorSetProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 1 + rand.Intn(6)
		removed := 1 - math.Pow(2.0/3.0, float64(k))
		return raBridgeProblem(o, fmt.Sprintf("Starting from [0, 1], the Cantor construction removes the middle third at each stage. What total length has been removed after %d stages? (enter a decimal to 6 places)", k), raBridgeDecimal(removed, 6), fmt.Sprintf("The remaining length after %d stages is (2/3)^%d, so the removed length is 1 - (2/3)^%d = %s.", k, k, k, raBridgeDecimal(removed, 6)))
	case 2:
		k := 1 + rand.Intn(8)
		ans := 1 << k
		return raBridgeProblem(o, fmt.Sprintf("How many closed intervals remain after %d stages of the Cantor construction? (enter a number)", k), strconv.Itoa(ans), fmt.Sprintf("Each interval splits into two, so the count doubles %d times: 2^%d = %d intervals.", k, k, ans))
	default:
		k := 1 + rand.Intn(6)
		rem := math.Pow(2.0/3.0, float64(k))
		return raBridgeProblem(o, fmt.Sprintf("What is the total length of the intervals remaining after %d stages of the Cantor construction? (enter a decimal to 6 places)", k), raBridgeDecimal(rem, 6), fmt.Sprintf("Each stage keeps two thirds of the length, so after %d stages the length is (2/3)^%d = %s.", k, k, raBridgeDecimal(rem, 6)))
	}
}

func totalBoundedProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(18)
		return raBridgeProblem(o, fmt.Sprintf("How many balls of radius 1/3 are needed at minimum to cover %d points that are pairwise more than 1 apart? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("No ball of radius 1/3 can contain two of the points, so one ball is needed for each of the %d points.", n))
	case 2:
		n := 2 + rand.Intn(11)
		return raBridgeProblem(o, fmt.Sprintf("How many intervals of length 1/%d are needed at minimum to cover [0, 1]? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("The total length is 1 and each interval has length 1/%d, so at least %d intervals are needed and %d suffice.", n, n, n))
	default:
		pool := []raBridgeItem{
			{"a finite set of points", true, "finite sets are totally bounded"},
			{"the interval [0, 1]", true, "compact subsets of R are totally bounded"},
			{"the interval (0, 1)", true, "it is bounded in R"},
			{"{1/n : n >= 1} together with 0", true, "a compact countable set"},
			{"the Cantor set", true, "compact and hence totally bounded"},
			{"the rationals in [0, 1]", true, "a bounded subset of R"},
			{"the algebraic numbers in [0, 1]", true, "a bounded subset of R"},
			{"a finite union of bounded intervals", true, "bounded in R"},
			{"the integers", false, "unbounded in R"},
			{"the real line", false, "unbounded"},
			{"the natural numbers", false, "unbounded"},
			{"the set {n : n >= 1}", false, "unbounded"},
			{"the interval [0, infinity)", false, "unbounded"},
			{"the unit ball of an infinite-dimensional Hilbert space", false, "not totally bounded"},
			{"an infinite set with all pairwise distances equal to 1", false, "it needs infinitely many balls of radius 1/3"},
			{"the set of all integers", false, "unbounded"},
		}
		return raBridgeCount(o, "How many of these metric spaces are totally bounded?", pool, 8)
	}
}

func weierstrassMProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 2 + rand.Intn(8)
		c := rand.Intn(6)
		a := rand.Intn(4)
		b := a + 1 + rand.Intn(5)
		ans := k*b + c
		return raBridgeProblem(o, fmt.Sprintf("For f(x) = %d x + %d on [%d, %d], what is the maximum of |f(x)|? (enter a number)", k, c, a, b), strconv.Itoa(ans), fmt.Sprintf("The function is increasing, so the maximum is at x = %d: %d(%d) + %d = %d.", b, k, b, c, ans))
	case 2:
		m := 1 + rand.Intn(9)
		r := []float64{0.5, 0.25, 0.2, 0.1}[rand.Intn(4)]
		sum := float64(m) / (1 - r)
		return raBridgeProblem(o, fmt.Sprintf("What is the sum of the geometric series with first term %d and ratio %s? (enter a decimal to 4 places)", m, raBridgeDecimal(r, 2)), raBridgeDecimal(sum, 4), fmt.Sprintf("The sum is a/(1-r) = %d/(1-%s) = %s.", m, raBridgeDecimal(r, 2), raBridgeDecimal(sum, 4)))
	default:
		pool := []raBridgeItem{
			{"p = 2", true, "the p-series converges for p > 1"},
			{"p = 3", true, "the p-series converges for p > 1"},
			{"p = 4", true, "the p-series converges for p > 1"},
			{"p = 5", true, "the p-series converges for p > 1"},
			{"p = 6", true, "the p-series converges for p > 1"},
			{"p = 7", true, "the p-series converges for p > 1"},
			{"p = 8", true, "the p-series converges for p > 1"},
			{"p = 9", true, "the p-series converges for p > 1"},
			{"p = 1", false, "the harmonic series diverges"},
			{"p = 0", false, "the terms do not tend to zero"},
			{"p = -1", false, "the terms grow"},
			{"p = -2", false, "the terms grow"},
			{"p = -3", false, "the terms grow"},
			{"p = -4", false, "the terms grow"},
			{"p = -5", false, "the terms grow"},
			{"p = -6", false, "the terms grow"},
		}
		return raBridgeCount(o, "For how many of these exponents p does the p-series with terms 1/n^p converge?", pool, 8)
	}
}

func boundedVariationProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 2 + rand.Intn(8)
		a := rand.Intn(5)
		b := a + 1 + rand.Intn(6)
		ans := k * (b - a)
		return raBridgeProblem(o, fmt.Sprintf("What is the total variation of f(x) = %d x on [%d, %d]? (enter a number)", k, a, b), strconv.Itoa(ans), fmt.Sprintf("A linear function increases steadily, so its variation is f(%d) - f(%d) = %d.", b, a, ans))
	case 2:
		n := 3 + rand.Intn(4)
		jumps := make([]string, n)
		total := 0
		for i := 0; i < n; i++ {
			d := 1 + rand.Intn(9)
			total += d
			jumps[i] = strconv.Itoa(d)
		}
		return raBridgeProblem(o, fmt.Sprintf("A step function jumps upward by %s at successive points in an interval and is constant between jumps. What is its total variation on the interval? (enter a number)", strings.Join(jumps, ", ")), strconv.Itoa(total), fmt.Sprintf("The total variation of a step function is the sum of the jump sizes: %s = %d.", strings.Join(jumps, " + "), total))
	default:
		pool := []raBridgeItem{
			{"x on [0, 1]", true, "monotone functions have bounded variation"},
			{"x^2 on [0, 1]", true, "monotone increasing"},
			{"e^x on [0, 1]", true, "monotone increasing"},
			{"arctan x on [0, 1]", true, "monotone increasing"},
			{"a step function with finitely many jumps on [0, 1]", true, "the variation is the sum of the jumps"},
			{"|x| on [-1, 1]", true, "monotone on each side of zero"},
			{"sqrt(x) on [0, 1]", true, "monotone increasing"},
			{"sin x on [0, 1]", true, "monotone on this interval"},
			{"sin(1/x) on (0, 1]", false, "it oscillates infinitely often near zero"},
			{"cos(1/x) on (0, 1]", false, "it oscillates infinitely often near zero"},
			{"x sin(1/x) on (0, 1]", false, "its variation is unbounded near zero"},
			{"x^2 sin(1/x^2) on (0, 1]", false, "its variation is unbounded near zero"},
			{"1/x on (0, 1]", false, "it is unbounded"},
			{"tan x on (0, 1]", false, "it is unbounded near pi/2"},
			{"e^(1/x) on (0, 1]", false, "it is unbounded near zero"},
			{"1/(1-x) on (0, 1)", false, "it is unbounded near 1"},
		}
		return raBridgeCount(o, "How many of these functions have bounded variation on the stated interval?", pool, 8)
	}
}

func riemannStieltjesProblem(o int) generator.Problem {
	switch o {
	case 1:
		c := 2 + rand.Intn(8)
		ga := 1 + rand.Intn(5)
		gb := ga + 1 + rand.Intn(8)
		ans := c * (gb - ga)
		return raBridgeProblem(o, fmt.Sprintf("What is the Riemann-Stieltjes integral of the constant %d with respect to g over [%d, %d] if g(%d) = %d and g(%d) = %d? (enter a number)", c, ga, gb, ga, ga, gb, gb), strconv.Itoa(ans), fmt.Sprintf("The integral of a constant c with respect to g is c(g(b)-g(a)) = %d(%d-%d) = %d.", c, gb, ga, ans))
	case 2:
		a := rand.Intn(6)
		b := a + 1 + rand.Intn(8)
		ans := float64(b*b-a*a) / 2
		return raBridgeProblem(o, fmt.Sprintf("What is the Riemann-Stieltjes integral of x with respect to x over [%d, %d]? (enter a decimal to 4 places)", a, b), raBridgeDecimal(ans, 4), fmt.Sprintf("This integral equals the ordinary integral of x, which is (b^2-a^2)/2 = (%d^2-%d^2)/2 = %s.", b, a, raBridgeDecimal(ans, 4)))
	default:
		a := rand.Intn(4)
		b := a + 1 + rand.Intn(4)
		ans := 2 * float64(b*b*b-a*a*a) / 3
		return raBridgeProblem(o, fmt.Sprintf("What is the Riemann-Stieltjes integral of x with respect to x^2 over [%d, %d]? (enter a decimal to 4 places)", a, b), raBridgeDecimal(ans, 4), fmt.Sprintf("With g(x)=x^2, the integral of x dg equals the integral of 2x^2 dx, which is 2(b^3-a^3)/3 = %s.", raBridgeDecimal(ans, 4)))
	}
}

func arzelaAscoliProblem(o int) generator.Problem {
	switch o {
	case 1:
		return raBridgeCount(o, "How many of these families of functions on [0, 1] are uniformly bounded and equicontinuous?", arzelaPool(), 8)
	case 2:
		l := 1 + rand.Intn(9)
		n := 2 + rand.Intn(8)
		ans := float64(l) / float64(n)
		return raBridgeProblem(o, fmt.Sprintf("For the sequence f_n(x) = x/n on [0, %d], what is the maximum of |f_n(x)|? (enter a decimal to 4 places)", l), raBridgeDecimal(ans, 4), fmt.Sprintf("The maximum is at x = %d, giving %d/%d = %s.", l, l, n, raBridgeDecimal(ans, 4)))
	default:
		return raBridgeCount(o, "How many of these sequences of functions on [0, 1] have a uniformly convergent subsequence?", arzelaPool(), 8)
	}
}

func arzelaPool() []raBridgeItem {
	return []raBridgeItem{
		{"f_n(x) = x/n", true, "bounded by 1 and equicontinuous"},
		{"f_n(x) = x^2 + 1/n", true, "bounded and uniformly equicontinuous"},
		{"f_n(x) = sin(x)/n", true, "bounded by 1/n and equicontinuous"},
		{"f_n(x) = e^x/n", true, "bounded by e/n and equicontinuous"},
		{"f_n(x) = 0", true, "the constant zero family"},
		{"f_n(x) = x + 1/n", true, "bounded and uniformly equicontinuous"},
		{"f_n(x) = cos(x)/n", true, "bounded by 1/n and equicontinuous"},
		{"f_n(x) = n/(n+1) x", true, "bounded by 1 and equicontinuous"},
		{"f_n(x) = x^n", false, "it is not equicontinuous near 1"},
		{"f_n(x) = sin(nx)", false, "the derivatives grow with n"},
		{"f_n(x) = cos(nx)", false, "the derivatives grow with n"},
		{"f_n(x) = n x(1-x)", false, "the maximum grows with n"},
		{"f_n(x) = arctan(nx)", false, "it is not equicontinuous near zero"},
		{"f_n(x) = e^(nx)", false, "the values grow with n"},
		{"f_n(x) = n x", false, "the values grow with n"},
		{"f_n(x) = sin(n^2 x)", false, "the derivatives grow with n"},
	}
}

func weierstrassApproxProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 2 + rand.Intn(8)
		ans := 1 / (2 * float64(n))
		return raBridgeProblem(o, fmt.Sprintf("For f(x)=x^2, the Bernstein bound gives |f(x)-B_n(f)(x)| <= M/(4n) with M the maximum of |f''|. What bound does this give at n = %d? (enter a decimal to 4 places)", n), raBridgeDecimal(ans, 4), fmt.Sprintf("Here f''(x)=2, so M=2 and the bound is 2/(4n) = 1/(2n) = %s.", raBridgeDecimal(ans, 4)))
	case 2:
		n := 2 + rand.Intn(8)
		ans := n + 1
		return raBridgeProblem(o, fmt.Sprintf("How many Bernstein basis polynomials B_{i,n} are there at degree n = %d? (enter a number)", n), strconv.Itoa(ans), fmt.Sprintf("The index i runs from 0 to n inclusive, so there are n+1 = %d basis polynomials.", ans))
	default:
		pool := []raBridgeItem{
			{"f(x) = x", true, "continuous on [0, 1]"},
			{"f(x) = sin x", true, "continuous on [0, 1]"},
			{"f(x) = |x|", true, "continuous on [0, 1]"},
			{"f(x) = e^x", true, "continuous on [0, 1]"},
			{"f(x) = cos x", true, "continuous on [0, 1]"},
			{"f(x) = sqrt(x)", true, "continuous on [0, 1]"},
			{"f(x) = x^3", true, "continuous on [0, 1]"},
			{"f(x) = 1/(1+x)", true, "continuous on [0, 1]"},
			{"the step function that is 0 on [0, 1/2) and 1 on [1/2, 1]", false, "it has a jump"},
			{"the sign function on [-1, 1]", false, "it has a jump at zero"},
			{"the Dirichlet function", false, "it is discontinuous everywhere"},
			{"the floor function on [0, 1]", false, "it jumps at 1"},
			{"sin(1/x) at x = 0, with value 0", false, "it is not continuous at zero"},
			{"e^(1/x) at x = 0, with value 0", false, "it is not continuous at zero"},
			{"1/(1-x) on [0, 1]", false, "it is unbounded at 1"},
			{"a function with a jump discontinuity", false, "uniform limits of polynomials are continuous"},
		}
		return raBridgeCount(o, "How many of these functions on [0, 1] are uniform limits of polynomials?", pool, 8)
	}
}
