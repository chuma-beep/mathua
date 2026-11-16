package realanalysis

// Wave 2h replaces the fixed-claim tables in fourteen real-analysis
// generators. Each old table offered six precomputed instances, so a learner
// reviewing any single knowledge point met the same six questions and then saw
// them repeat. Every family here branches on KPOrdinal and draws its operands,
// so one ordinal alone supplies at least ten distinct problems. The fact keys
// that the older tests and any authored schema rely on are kept: lo, hi, pn,
// pd, a, b, n, x0, m, en, ed, c and answer.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func raW2HOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func raW2HProblem(o int, question, answer, explanation string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = answer
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   o,
	}
}

func raW2HRange(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func raW2HGCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func raW2HReduce(pn, pd int) (int, int) {
	g := raW2HGCD(pn, pd)
	if g == 0 {
		return pn, pd
	}
	return pn / g, pd / g
}

// raW2HLinear writes mx + b with the sign of b folded into the expression.
func raW2HLinear(m, b int) string {
	term := fmt.Sprintf("%dx", m)
	switch m {
	case 1:
		term = "x"
	case -1:
		term = "−x"
	}
	switch {
	case b == 0:
		return term
	case b > 0:
		return fmt.Sprintf("%s + %d", term, b)
	default:
		return fmt.Sprintf("%s − %d", term, -b)
	}
}

// raW2HXShift writes (x + b), folding a negative shift into a minus sign.
func raW2HXShift(b int) string {
	switch {
	case b == 0:
		return "x"
	case b > 0:
		return fmt.Sprintf("(x + %d)", b)
	default:
		return fmt.Sprintf("(x − %d)", -b)
	}
}

// ---------------------------------------------------------------------------
// Supremum and infimum
// ---------------------------------------------------------------------------

func raW2HSupremum(kp int) generator.Problem {
	switch kp {
	case 1:
		lo := raW2HRange(-12, 8)
		hi := lo + raW2HRange(1, 12)
		var set string
		switch rand.Intn(3) {
		case 0:
			set = fmt.Sprintf("the open interval of all x with %d < x < %d", lo, hi)
		case 1:
			set = fmt.Sprintf("the half-open interval of all x with %d < x ≤ %d", lo, hi)
		default:
			set = fmt.Sprintf("the closed interval of all x with %d ≤ x ≤ %d", lo, hi)
		}
		answer := strconv.Itoa(hi)
		return raW2HProblem(kp,
			"What is the supremum of "+set+"? (enter a decimal)",
			answer,
			"The supremum is the least upper bound. Every element of "+set+" is at most "+answer+
				", and elements come arbitrarily close to it from below, so no smaller number is an upper bound.",
			map[string]string{"lo": strconv.Itoa(lo), "hi": answer, "kind": "interval"})
	case 2:
		lo := raW2HRange(-9, 7)
		hi := lo + raW2HRange(1, 12)
		var set string
		if rand.Intn(2) == 0 {
			set = fmt.Sprintf("the half-open interval of all x with %d ≤ x < %d", lo, hi)
		} else {
			set = fmt.Sprintf("the half-open interval of all x with %d < x ≤ %d", lo, hi)
		}
		answer := strconv.Itoa(hi)
		return raW2HProblem(kp,
			"What is the supremum of "+set+"? (enter a decimal)",
			answer,
			"The upper endpoint "+answer+" is an upper bound. Numbers just below it are exceeded by elements of the half-open interval, so no smaller number is an upper bound.",
			map[string]string{"lo": strconv.Itoa(lo), "hi": answer, "kind": "half_open"})
	default:
		a := raW2HRange(-10, 10)
		answer := a + 1
		return raW2HProblem(kp,
			fmt.Sprintf("The set S = {%d + (n − 1)/n : n ≥ 1} is bounded above but has no maximum. What is the supremum of S? (enter a decimal)", a),
			strconv.Itoa(answer),
			fmt.Sprintf("The fractions (n − 1)/n increase towards 1 without reaching it, so every element of S is below %d and elements come arbitrarily close to %d. The supremum is %d, although no element of S equals it.", answer, answer, answer),
			map[string]string{"lo": strconv.Itoa(a), "hi": strconv.Itoa(answer), "kind": "no_maximum"})
	}
}

// ---------------------------------------------------------------------------
// The Archimedean property
// ---------------------------------------------------------------------------

func raW2HArchimedean(kp int) generator.Problem {
	pd := raW2HRange(3, 60)
	pn := raW2HRange(1, pd-1)
	pn, pd = raW2HReduce(pn, pd)
	switch kp {
	case 1:
		answer := pd/pn + 1
		return raW2HProblem(kp,
			fmt.Sprintf("By the Archimedean property there is a positive integer N with 1/N < ε. For ε = %d/%d, what is the smallest such N? (enter a whole number)", pn, pd),
			strconv.Itoa(answer),
			fmt.Sprintf("The inequality 1/N < %d/%d is the same as N > %d/%d. The smallest integer above %d/%d is %d.", pn, pd, pd, pn, pd, pn, answer),
			map[string]string{"pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "kind": "smallest_reciprocal"})
	case 2:
		answer := pd/pn + 1
		return raW2HProblem(kp,
			fmt.Sprintf("The Archimedean property says that for every real number B there is an integer N with N > B. For B = %d/%d, what is the smallest such N? (enter a whole number)", pd, pn),
			strconv.Itoa(answer),
			fmt.Sprintf("The real number B is %d/%d. Taking the floor and adding one gives the smallest integer strictly greater than B, which is %d.", pd, pn, answer),
			map[string]string{"pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "kind": "integer_above"})
	default:
		answer := pd / pn
		return raW2HProblem(kp,
			fmt.Sprintf("For ε = %d/%d, what is the largest positive integer N with 1/N ≥ ε? (enter a whole number)", pn, pd),
			strconv.Itoa(answer),
			fmt.Sprintf("The inequality 1/N ≥ %d/%d is the same as N ≤ %d/%d. The largest integer at most %d/%d is %d, and for that N the reciprocal is still at least ε.", pn, pd, pd, pn, pd, pn, answer),
			map[string]string{"pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "kind": "largest_reciprocal"})
	}
}

// ---------------------------------------------------------------------------
// Monotone convergence
// ---------------------------------------------------------------------------

func raW2HMonotone(kp int) generator.Problem {
	switch kp {
	case 1:
		limit := raW2HRange(2, 20)
		a := limit * (limit - 1)
		return raW2HProblem(kp,
			fmt.Sprintf("The increasing sequence x_1 = 0, x_{n+1} = sqrt(%d + x_n) is bounded above and converges. Its limit L satisfies L = sqrt(%d + L). What is L? (enter a decimal)", a, a),
			strconv.Itoa(limit),
			fmt.Sprintf("Squaring L = sqrt(%d + L) gives L² − L − %d = 0, whose positive root is (1 + sqrt(1 + %d))/2 = %d. The increasing sequence starting at zero converges to that root.", a, a, 4*a, limit),
			map[string]string{"a": strconv.Itoa(a), "kind": "sqrt_linear"})
	case 2:
		limit := raW2HRange(2, 12)
		b := raW2HRange(1, 3)
		a := limit*limit - b*limit
		if a <= 0 {
			b = 1
			a = limit * (limit - 1)
		}
		return raW2HProblem(kp,
			fmt.Sprintf("The increasing sequence x_1 = 0, x_{n+1} = sqrt(%d + %d·x_n) is bounded above and converges. Its limit L satisfies L = sqrt(%d + %d·L). What is L? (enter a decimal)", a, b, a, b),
			strconv.Itoa(limit),
			fmt.Sprintf("Squaring gives L² − %d·L − %d = 0. The positive root is (%d + sqrt(%d + %d))/2 = %d, and the increasing sequence converges to it.", b, a, b, b*b, 4*a, limit),
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "kind": "sqrt_affine"})
	default:
		limit := raW2HRange(2, 20)
		q := raW2HRange(3, 5)
		c := (q - 1) * limit
		start := limit - raW2HRange(1, 8)
		return raW2HProblem(kp,
			fmt.Sprintf("The sequence x_1 = %d, x_{n+1} = (x_n + %d)/%d is increasing and bounded above. What is its limit? (enter a decimal)", start, c, q),
			strconv.Itoa(limit),
			fmt.Sprintf("A limit L must satisfy L = (L + %d)/%d, so %d·L = %d and L = %d. Monotonicity and the upper bound show the sequence really converges to it.", c, q, q-1, c, limit),
			map[string]string{"start": strconv.Itoa(start), "c": strconv.Itoa(c), "q": strconv.Itoa(q), "limit": strconv.Itoa(limit), "kind": "affine_average"})
	}
}

// ---------------------------------------------------------------------------
// Series and partial sums
// ---------------------------------------------------------------------------

func raW2HSeries(kp int) generator.Problem {
	switch kp {
	case 1:
		a := raW2HRange(1, 6)
		pd := raW2HRange(2, 9)
		pn, pd := raW2HReduce(raW2HRange(1, pd-1), pd)
		n := raW2HRange(3, 8)
		r := float64(pn) / float64(pd)
		sum := float64(a) * (1 - math.Pow(r, float64(n))) / (1 - r)
		answer := decimal6(sum)
		return raW2HProblem(kp,
			fmt.Sprintf("For the geometric series with first term %d and ratio r = %d/%d, what is the sum of its first %d terms? (enter a decimal)", a, pn, pd, n),
			answer,
			fmt.Sprintf("The partial sum is a(1 − r^n)/(1 − r). Here a = %d, r = %d/%d and n = %d, which gives %s.", a, pn, pd, n, answer),
			map[string]string{"a": strconv.Itoa(a), "pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "n": strconv.Itoa(n), "kind": "partial_sum"})
	case 2:
		a := raW2HRange(1, 8)
		pd := raW2HRange(2, 10)
		pn, pd := raW2HReduce(raW2HRange(1, pd-1), pd)
		answer := decimal6(float64(a) * float64(pd) / float64(pd-pn))
		return raW2HProblem(kp,
			fmt.Sprintf("The geometric series with first term %d and ratio r = %d/%d converges. What is its sum? (enter a decimal)", a, pn, pd),
			answer,
			fmt.Sprintf("A geometric series with |r| < 1 sums to a/(1 − r). Here that is %d/(1 − %d/%d) = %s.", a, pn, pd, answer),
			map[string]string{"a": strconv.Itoa(a), "pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "kind": "infinite_sum"})
	default:
		a := raW2HRange(1, 9)
		c := raW2HRange(1, 9)
		b := raW2HRange(0, 9)
		d := raW2HRange(1, 9)
		answer := decimal6(float64(a) / float64(c))
		return raW2HProblem(kp,
			fmt.Sprintf("The series has terms a_n = (%d·n + %d)/(%d·n + %d). The terms do not tend to zero, so the divergence test applies. What is the limit of a_n? (enter a decimal)", a, b, c, d),
			answer,
			fmt.Sprintf("Dividing numerator and denominator by n leaves the ratio of the leading coefficients, %d/%d = %s. A series whose terms tend to a nonzero limit diverges.", a, c, answer),
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "kind": "term_limit"})
	}
}

// ---------------------------------------------------------------------------
// The limit of a function
// ---------------------------------------------------------------------------

func raW2HFunctionLimit(kp int) generator.Problem {
	switch kp {
	case 1:
		m := raW2HRange(1, 9)
		en := raW2HRange(1, 5)
		ed := raW2HRange(1, 10)
		answer := decimal6((float64(en) / float64(ed)) / float64(m))
		return raW2HProblem(kp,
			fmt.Sprintf("For f(x) = %s and its limit at x = 0, the largest δ that makes |f(x) − f(0)| < ε whenever |x − 0| < δ is ε divided by the slope. For ε = %d/%d, what is δ? (enter a decimal)", raW2HLinear(m, 0), en, ed),
			answer,
			fmt.Sprintf("Here |f(x) − f(0)| = %d|x|, so the condition is %d|x| < %d/%d. The largest δ is ε/%d = %s.", m, m, en, ed, m, answer),
			map[string]string{"m": strconv.Itoa(m), "en": strconv.Itoa(en), "ed": strconv.Itoa(ed), "b": "0", "c": "0", "kind": "linear_at_zero"})
	case 2:
		m := raW2HRange(1, 9)
		b := raW2HRange(-6, 6)
		c := raW2HRange(-5, 5)
		en := raW2HRange(1, 5)
		ed := raW2HRange(1, 10)
		answer := decimal6((float64(en) / float64(ed)) / float64(m))
		return raW2HProblem(kp,
			fmt.Sprintf("For f(x) = %s and its limit at x = %d, the largest δ that makes |f(x) − f(%d)| < ε whenever |x − %d| < δ is ε divided by the slope. For ε = %d/%d, what is δ? (enter a decimal)", raW2HLinear(m, b), c, c, c, en, ed),
			answer,
			fmt.Sprintf("Because f is linear, |f(x) − f(%d)| = %d|x − %d|. The condition is therefore |x − %d| < ε/%d, and the largest δ is %s.", c, m, c, c, m, answer),
			map[string]string{"m": strconv.Itoa(m), "en": strconv.Itoa(en), "ed": strconv.Itoa(ed), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "kind": "linear_at_point"})
	default:
		a := raW2HRange(4, 12)
		if rand.Intn(2) == 0 {
			answer := 2 * a
			return raW2HProblem(kp,
				fmt.Sprintf("The function f(x) = (x² − %d)/(x − %d) is undefined at x = %d but has a limit there. What is the limit? (enter a decimal)", a*a, a, a),
				strconv.Itoa(answer),
				fmt.Sprintf("Factor the numerator as (x − %d)(x + %d). For x ≠ %d the quotient is x + %d, which tends to %d as x tends to %d. Changing the value at %d cannot change the limit.", a, a, a, a, answer, a, a),
				map[string]string{"a": strconv.Itoa(a), "degree": "2", "kind": "factor_cancel"})
		}
		answer := 3 * a * a
		return raW2HProblem(kp,
			fmt.Sprintf("The function f(x) = (x³ − %d)/(x − %d) is undefined at x = %d but has a limit there. What is the limit? (enter a decimal)", a*a*a, a, a),
			strconv.Itoa(answer),
			fmt.Sprintf("Factor the numerator as (x − %d)(x² + %dx + %d). For x ≠ %d the quotient tends to %d as x tends to %d. The value at %d is irrelevant to the limit.", a, a, a*a, a, answer, a, a),
			map[string]string{"a": strconv.Itoa(a), "degree": "3", "kind": "factor_cancel"})
	}
}

// ---------------------------------------------------------------------------
// The mean value theorem
// ---------------------------------------------------------------------------

func raW2HMVT(kp int) generator.Problem {
	switch kp {
	case 1:
		b := raW2HRange(2, 20)
		answer := decimal6(float64(b) / 2)
		return raW2HProblem(kp,
			fmt.Sprintf("For f(x) = x² on the interval from 0 to %d, the mean value theorem gives a point c with f'(c) equal to the average slope. What is c? (enter a decimal)", b),
			answer,
			fmt.Sprintf("The average slope is (f(%d) − f(0))/%d = %d. Since f'(c) = 2c, the equation 2c = %d gives c = %s.", b, b, b, b, answer),
			map[string]string{"b": strconv.Itoa(b), "kind": "square_from_zero"})
	case 2:
		a := raW2HRange(-8, 8)
		b := a + raW2HRange(1, 12)
		answer := decimal6(float64(a+b) / 2)
		return raW2HProblem(kp,
			fmt.Sprintf("For f(x) = x² on the interval from %d to %d, the mean value theorem gives a point c with f'(c) equal to the average slope. What is c? (enter a decimal)", a, b),
			answer,
			fmt.Sprintf("The average slope is (f(%d) − f(%d))/(%d − %d) = %d. Since f'(c) = 2c, the equation 2c = %d gives c = %s.", b, a, b, a, a+b, a+b, answer),
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "kind": "square_general"})
	default:
		m := raW2HRange(1, 9)
		a := raW2HRange(-6, 6)
		b := a + raW2HRange(1, 10)
		answer := m * (b - a)
		return raW2HProblem(kp,
			fmt.Sprintf("A differentiable function on [%d, %d] has |f'(x)| ≤ %d throughout. By the mean value theorem, what is the largest possible value of |f(%d) − f(%d)|?", a, b, m, b, a),
			strconv.Itoa(answer),
			fmt.Sprintf("The mean value theorem gives f(%d) − f(%d) = f'(c)(%d − %d) for some c. Taking absolute values bounds the change by %d·(%d − %d) = %d.", b, a, b, a, m, b, a, answer),
			map[string]string{"m": strconv.Itoa(m), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "kind": "derivative_bound"})
	}
}

// ---------------------------------------------------------------------------
// The chain rule
// ---------------------------------------------------------------------------

func raW2HChain(kp int) generator.Problem {
	switch kp {
	case 1:
		a := raW2HRange(1, 5)
		b := raW2HRange(-4, 4)
		n := raW2HRange(2, 5)
		x0 := raW2HRange(-3, 3)
		inner := a*x0 + b
		answer := decimal6(float64(n * a * raIntPow(inner, n-1)))
		return raW2HProblem(kp,
			fmt.Sprintf("By the chain rule, the derivative of f(x) = (%s)^%d is f'(x) = %d·%d·(%s)^%d. What is f'(x) at x = %d? (enter a decimal)", raW2HLinear(a, b), n, n, a, raW2HLinear(a, b), n-1, x0),
			answer,
			fmt.Sprintf("The inner value is %d·%d + (%d) = %d. The derivative is %d·%d·%d^%d = %s.", a, x0, b, inner, n, a, inner, n-1, answer),
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "n": strconv.Itoa(n), "x0": strconv.Itoa(x0), "kind": "power_of_linear"})
	case 2:
		s := raW2HRange(2, 10)
		a := raW2HRange(1, 6)
		x0 := raW2HRange(-2, 3)
		b := s*s - a*x0
		answer := decimal6(float64(a) / float64(2*s))
		return raW2HProblem(kp,
			fmt.Sprintf("For f(x) = sqrt(%s), what is f'(x) at x = %d? (enter a decimal)", raW2HLinear(a, b), x0),
			answer,
			fmt.Sprintf("The chain rule gives f'(x) = %d/(2·sqrt(%s)). At x = %d the inner value is %d² = %d, so the derivative is %d/(2·%d) = %s.", a, raW2HLinear(a, b), x0, s, s*s, a, s, answer),
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "x0": strconv.Itoa(x0), "inner": strconv.Itoa(s * s), "kind": "square_root"})
	default:
		b := raW2HRange(0, 4)
		n := raW2HRange(2, 3)
		c := raW2HRange(1, 5)
		m := raW2HRange(2, 3)
		x0 := raW2HRange(-2, 2)
		u := x0 + b
		v := raIntPow(u, n) + c
		answer := m * raIntPow(v, m-1) * n * raIntPow(u, n-1)
		return raW2HProblem(kp,
			fmt.Sprintf("Apply the chain rule twice to h(x) = ((%s)^%d + %d)^%d. What is h'(x) at x = %d? (enter a whole number)", raW2HXShift(b), n, c, m, x0),
			strconv.Itoa(answer),
			fmt.Sprintf("At x = %d the inner base is %d, so the inner power is %d and the outer base is %d. The two rates are %d·%d^%d for the outer power and %d·%d^%d for the inner power, whose product is %d.", x0, u, raIntPow(u, n), v, m, v, m-1, n, u, n-1, answer),
			map[string]string{"b": strconv.Itoa(b), "n": strconv.Itoa(n), "c": strconv.Itoa(c), "m": strconv.Itoa(m), "x0": strconv.Itoa(x0), "kind": "double_chain"})
	}
}

// ---------------------------------------------------------------------------
// Taylor's theorem
// ---------------------------------------------------------------------------

func raW2HTaylor(kp int) generator.Problem {
	switch kp {
	case 1:
		a := raW2HRange(1, 6)
		k := raW2HRange(2, 5)
		answer := decimal6(float64(raIntPow(a, k)) / float64(raBridgeFactorial(k)))
		return raW2HProblem(kp,
			fmt.Sprintf("The Taylor series at 0 of e^(%d·x) has x^%d coefficient a^k/k!. For a = %d, what is the coefficient of x^%d? (enter a decimal)", a, k, a, k),
			answer,
			fmt.Sprintf("The %d-th derivative of e^(%d·x) at 0 is %d^%d, so the coefficient is %d/%d! = %s.", k, a, a, k, raIntPow(a, k), k, answer),
			map[string]string{"a": strconv.Itoa(a), "k": strconv.Itoa(k), "kind": "exponential_coefficient"})
	case 2:
		a := raW2HRange(2, 20)
		c := raW2HRange(0, 5)
		answer := decimal6(1 / float64(a))
		center := "x"
		if c > 0 {
			center = fmt.Sprintf("(x − %d)", c)
		}
		return raW2HProblem(kp,
			fmt.Sprintf("What is the radius of convergence of the power series with terms %d^n·%s^n? (enter a decimal)", a, center),
			answer,
			fmt.Sprintf("The ratio of successive coefficients has magnitude %d, so the root test gives radius 1/%d = %s. Shifting the centre does not change the radius.", a, a, answer),
			map[string]string{"a": strconv.Itoa(a), "c": strconv.Itoa(c), "kind": "radius"})
	default:
		m := raW2HRange(6, 16)
		n := raW2HRange(2, 3)
		answer := raBridgeFactorial(m) / (raBridgeFactorial(n) * raBridgeFactorial(m-n))
		return raW2HProblem(kp,
			fmt.Sprintf("In the Taylor expansion of (1 + x)^%d at 0, the coefficient of x^%d is the binomial coefficient m choose n. For m = %d and n = %d, what is that coefficient? (enter a whole number)", m, n, m, n),
			strconv.Itoa(answer),
			fmt.Sprintf("Differentiating (1 + x)^%d %d times at 0 gives %d!/(%d − %d)!, and dividing by %d! leaves the binomial coefficient %d.", m, n, m, m, n, n, answer),
			map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n), "kind": "binomial_coefficient"})
	}
}

// ---------------------------------------------------------------------------
// L'Hopital's rule
// ---------------------------------------------------------------------------

func raW2HLHopital(kp int) generator.Problem {
	a := raW2HRange(3, 15)
	switch kp {
	case 1:
		answer := decimal6(float64(a*a) / 2)
		return raW2HProblem(kp,
			fmt.Sprintf("The limit of (1 − cos(%d·x))/x² as x tends to 0 is a²/2, by L'Hopital's rule applied twice. For a = %d, what is the limit? (enter a decimal)", a, a),
			answer,
			fmt.Sprintf("Both numerator and denominator tend to zero. Differentiating twice gives %d²·cos(%d·x)/2, and cos(0) = 1, so the limit is %d/2 = %s.", a, a, a*a, answer),
			map[string]string{"a": strconv.Itoa(a), "kind": "cosine"})
	case 2:
		answer := decimal6(float64(a*a) / 2)
		return raW2HProblem(kp,
			fmt.Sprintf("Apply L'Hopital's rule twice to the limit of (e^(%d·x) − 1 − %d·x)/x² as x tends to 0. What is the limit? (enter a decimal)", a, a),
			answer,
			fmt.Sprintf("The first differentiation gives (%d·e^(%d·x) − %d)/(2x), still zero over zero. The second gives %d²·e^(%d·x)/2, whose value at zero is %d/2 = %s.", a, a, a, a, a, a*a, answer),
			map[string]string{"a": strconv.Itoa(a), "kind": "exponential_remainder"})
	default:
		answer := decimal6(float64(a*a*a) / 6)
		return raW2HProblem(kp,
			fmt.Sprintf("Applying L'Hopital's rule twice to (%d·x − sin(%d·x))/x³ leaves an indeterminate form. Continuing once more, what is the limit as x tends to 0? (enter a decimal)", a, a),
			answer,
			fmt.Sprintf("The first two derivatives of the numerator vanish at zero and the third is %d³, while the third derivative of x³ is 6. The limit is therefore %d/6 = %s.", a, a*a*a, answer),
			map[string]string{"a": strconv.Itoa(a), "kind": "sine_remainder"})
	}
}

// ---------------------------------------------------------------------------
// The Riemann integral
// ---------------------------------------------------------------------------

func raW2HRiemann(kp int) generator.Problem {
	switch kp {
	case 1:
		b := raW2HRange(2, 15)
		n := raW2HRange(2, 20)
		answer := decimal6(float64(b*b) * float64(n+1) / (2 * float64(n)))
		return raW2HProblem(kp,
			fmt.Sprintf("Use the right-endpoint Riemann sum with n = %d equal subintervals to approximate the integral of f(x) = x on the interval from 0 to %d. What is the value of the sum? (enter a decimal)", n, b),
			answer,
			fmt.Sprintf("With step %d/%d the sum is Σ (i·%d/%d)(%d/%d) for i from 1 to %d. Using the triangular sum gives %d²(%d + 1)/(2·%d) = %s.", b, n, b, n, b, n, n, b, n, n, answer),
			map[string]string{"b": strconv.Itoa(b), "n": strconv.Itoa(n), "kind": "right_endpoint_linear"})
	case 2:
		a := raW2HRange(-5, 5)
		b := a + raW2HRange(1, 10)
		n := raW2HRange(2, 12)
		answer := decimal6(float64((b-a)*(b-a)) / float64(n))
		return raW2HProblem(kp,
			fmt.Sprintf("For f(x) = x on [%d, %d] with n = %d equal subintervals, what is the difference between the upper sum and the lower sum? (enter a decimal)", a, b, n),
			answer,
			fmt.Sprintf("On each subinterval the maximum exceeds the minimum by the step (%d − %d)/%d. Summing that gap over %d subintervals gives (%d − %d)²/%d = %s.", b, a, n, n, b, a, n, answer),
			map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "n": strconv.Itoa(n), "kind": "upper_lower_gap"})
	default:
		b := raW2HRange(2, 10)
		n := raW2HRange(2, 15)
		answer := decimal6(float64(b*b*b) * float64(n+1) * float64(2*n+1) / (6 * float64(n*n)))
		return raW2HProblem(kp,
			fmt.Sprintf("Use the right-endpoint Riemann sum with n = %d equal subintervals to approximate the integral of f(x) = x² on the interval from 0 to %d. What is the value of the sum? (enter a decimal)", n, b),
			answer,
			fmt.Sprintf("The sum is Σ (i·%d/%d)²(%d/%d). Using the square sum gives %d³(%d + 1)(2·%d + 1)/(6·%d²) = %s.", b, n, b, n, b, n, n, n, answer),
			map[string]string{"b": strconv.Itoa(b), "n": strconv.Itoa(n), "kind": "right_endpoint_square"})
	}
}

// ---------------------------------------------------------------------------
// Integration by parts
// ---------------------------------------------------------------------------

func raW2HParts(kp int) generator.Problem {
	switch kp {
	case 1:
		b := raW2HRange(1, 20)
		answer := decimal4(float64(b-1)*math.Exp(float64(b)) + 1)
		return raW2HProblem(kp,
			fmt.Sprintf("Use integration by parts to evaluate the integral from 0 to %d of x·e^x with respect to x. What is the value? (enter a decimal)", b),
			answer,
			fmt.Sprintf("With u = x and dv = e^x dx, the boundary term is [x·e^x] from 0 to %d and the remaining integral is ∫ e^x dx. The value is (%d − 1)e^%d + 1 = %s.", b, b, b, answer),
			map[string]string{"b": strconv.Itoa(b), "kind": "linear_times_exponential"})
	case 2:
		b := raW2HRange(1, 15)
		answer := decimal4(math.Exp(float64(b))*(float64(b*b)-2*float64(b)+2) - 2)
		return raW2HProblem(kp,
			fmt.Sprintf("Use integration by parts twice to evaluate the integral from 0 to %d of x²·e^x with respect to x. What is the value? (enter a decimal)", b),
			answer,
			fmt.Sprintf("Two rounds of parts leave the antiderivative e^x(x² − 2x + 2). Evaluating between 0 and %d gives e^%d(%d² − 2·%d + 2) − 2 = %s.", b, b, b, b, answer),
			map[string]string{"b": strconv.Itoa(b), "kind": "quadratic_times_exponential"})
	default:
		b := raW2HRange(1, 12)
		answer := decimal4(math.Exp(float64(b))*(float64(b*b*b)-3*float64(b*b)+6*float64(b)-6) + 6)
		return raW2HProblem(kp,
			fmt.Sprintf("Use tabular integration by parts to evaluate the integral from 0 to %d of x³·e^x with respect to x. What is the value? (enter a decimal)", b),
			answer,
			fmt.Sprintf("The tabular antiderivative is e^x(x³ − 3x² + 6x − 6). Evaluating between 0 and %d gives e^%d(%d³ − 3·%d² + 6·%d − 6) + 6 = %s.", b, b, b, b, b, answer),
			map[string]string{"b": strconv.Itoa(b), "kind": "cubic_times_exponential"})
	}
}

// ---------------------------------------------------------------------------
// Power series
// ---------------------------------------------------------------------------

func raW2HPowerSeries(kp int) generator.Problem {
	c := raW2HRange(2, 25)
	answer := decimal6(1 / float64(c))
	switch kp {
	case 1:
		return raW2HProblem(kp,
			fmt.Sprintf("Find the radius of convergence R of the power series with terms %d^n·x^n. (enter a decimal)", c),
			answer,
			fmt.Sprintf("By the root test, the n-th root of the coefficient magnitude is %d, so the series converges for |x| < 1/%d and R = %s.", c, c, answer),
			map[string]string{"c": strconv.Itoa(c), "kind": "root_test"})
	case 2:
		return raW2HProblem(kp,
			fmt.Sprintf("The power series with terms %d^n·x^n is differentiated term by term. What is the radius of convergence of the derivative, whose terms are n·%d^n·x^(n−1)? (enter a decimal)", c, c),
			answer,
			fmt.Sprintf("Differentiating term by term does not change the radius. The root test still gives %d for the coefficients, so R = 1/%d = %s.", c, c, answer),
			map[string]string{"c": strconv.Itoa(c), "kind": "term_by_term"})
	default:
		a := raW2HRange(1, 9)
		return raW2HProblem(kp,
			fmt.Sprintf("Find the radius of convergence R of the power series with terms %d^n·(x − %d)^n. (enter a decimal)", c, a),
			answer,
			fmt.Sprintf("Shifting the centre does not change the radius. The ratio test gives |%d(x − %d)| < 1, so |x − %d| < 1/%d and R = %s.", c, a, a, c, answer),
			map[string]string{"c": strconv.Itoa(c), "a": strconv.Itoa(a), "kind": "shifted"})
	}
}

// ---------------------------------------------------------------------------
// The contraction mapping principle
// ---------------------------------------------------------------------------

func raW2HContraction(kp int) generator.Problem {
	pd := raW2HRange(3, 12)
	pn, pd := raW2HReduce(raW2HRange(1, pd-1), pd)
	a := float64(pn) / float64(pd)
	switch kp {
	case 1:
		b := raW2HRange(-9, 9)
		answer := decimal6(float64(b) / (1 - a))
		return raW2HProblem(kp,
			fmt.Sprintf("The contraction T(x) = a·x + b with a = %d/%d and b = %d has a unique fixed point, the limit of its iterates from any start. What is the fixed point? (enter a decimal)", pn, pd, b),
			answer,
			fmt.Sprintf("The fixed point solves x = a·x + b, so x = b/(1 − a) = %d/(1 − %d/%d) = %s.", b, pn, pd, answer),
			map[string]string{"pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "b": strconv.Itoa(b), "kind": "fixed_point"})
	case 2:
		k := raW2HRange(2, 5)
		answer := decimal6(math.Pow(a, float64(k)))
		return raW2HProblem(kp,
			fmt.Sprintf("A contraction T has ratio a = %d/%d. Its %d-th iterate is again a contraction. What is the contraction ratio of T^%d? (enter a decimal)", pn, pd, k, k),
			answer,
			fmt.Sprintf("The k-th iterate multiplies the distance bound by the same factor k times. Here the factor is (%d/%d)^%d = %s.", pn, pd, k, answer),
			map[string]string{"pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "k": strconv.Itoa(k), "kind": "error_rate"})
	default:
		b := raW2HRange(1, 12)
		start := raW2HRange(-10, 10)
		answer := decimal6(float64(b) / (1 - a))
		return raW2HProblem(kp,
			fmt.Sprintf("Starting from x_0 = %d, the iterates of T(x) = a·x + b with a = %d/%d and b = %d converge to the fixed point. What is that limit? (enter a decimal)", start, pn, pd, b),
			answer,
			fmt.Sprintf("The limit L solves L = a·L + b, so L = b/(1 − a) = %d/(1 − %d/%d) = %s. The starting value does not affect the limit because %d/%d is below one.", b, pn, pd, answer, pn, pd),
			map[string]string{"pn": strconv.Itoa(pn), "pd": strconv.Itoa(pd), "b": strconv.Itoa(b), "start": strconv.Itoa(start), "kind": "iterate_limit"})
	}
}

// ---------------------------------------------------------------------------
// The Cantor set
// ---------------------------------------------------------------------------

func raW2HCantorSet(kp int) generator.Problem {
	switch kp {
	case 1:
		k := raW2HRange(1, 12)
		if rand.Intn(2) == 0 {
			answer := 1 << k
			return raW2HProblem(kp,
				fmt.Sprintf("Starting from [0, 1], the Cantor construction removes the middle third of every remaining interval at each stage. How many closed intervals remain after %d stages? (enter a number)", k),
				strconv.Itoa(answer),
				fmt.Sprintf("Each interval splits into two, so the count doubles %d times: 2^%d = %d intervals.", k, k, answer),
				map[string]string{"k": strconv.Itoa(k), "kind": "remaining_intervals"})
		}
		answer := 1 << (k - 1)
		return raW2HProblem(kp,
			fmt.Sprintf("In the Cantor construction, how many open middle-third intervals are removed at stage %d, counting only the intervals removed at that stage? (enter a number)", k),
			strconv.Itoa(answer),
			fmt.Sprintf("At stage %d there are 2^%d intervals present before the removal, and each loses its middle third, so 2^%d = %d intervals are removed.", k, k-1, k-1, answer),
			map[string]string{"k": strconv.Itoa(k), "kind": "removed_intervals"})
	case 2:
		k := raW2HRange(1, 15)
		removed := 1 - math.Pow(2.0/3.0, float64(k))
		answer := decimal6(removed)
		return raW2HProblem(kp,
			fmt.Sprintf("Starting from [0, 1], the Cantor construction removes the middle third at each stage. What total length has been removed after %d stages? (enter a decimal to 6 places)", k),
			answer,
			fmt.Sprintf("The remaining length after %d stages is (2/3)^%d, so the removed length is 1 − (2/3)^%d = %s.", k, k, k, answer),
			map[string]string{"k": strconv.Itoa(k), "kind": "removed_length"})
	default:
		k := raW2HRange(1, 15)
		remaining := math.Pow(2.0/3.0, float64(k))
		answer := decimal6(remaining)
		return raW2HProblem(kp,
			fmt.Sprintf("What is the total length of the intervals remaining after %d stages of the Cantor construction? (enter a decimal to 6 places)", k),
			answer,
			fmt.Sprintf("Each stage keeps two thirds of the length, so after %d stages the length is (2/3)^%d = %s.", k, k, answer),
			map[string]string{"k": strconv.Itoa(k), "kind": "remaining_length"})
	}
}
