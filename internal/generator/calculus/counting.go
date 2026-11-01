package calculus

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// calcCountGen serves a computed count question instead of a recognition
// table. The nine calculus concepts it covers are conceptual, so the honest
// production form is a claim pool: the learner decides how many of the listed
// statements are true rather than picking a definition.
type calcCountGen struct {
	build func() generator.Problem
}

func (g *calcCountGen) Generate(generator.GeneratorContext) generator.Problem {
	return g.build()
}

type calcVerdict struct {
	label string
	yes   bool
	note  string
}

const calcDraw = 8

func countCalcVerdicts(lead string, pool []calcVerdict) generator.Problem {
	var yes, no []calcVerdict
	for _, c := range pool {
		if c.yes {
			yes = append(yes, c)
		} else {
			no = append(no, c)
		}
	}
	var chosen []calcVerdict
	count := 0
	var question string
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(calcDraw + 1)
		rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
		rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
		chosen = chosen[:0]
		chosen = append(chosen, yes[:count]...)
		chosen = append(chosen, no[:calcDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		parts := make([]string, len(chosen))
		for i, c := range chosen {
			parts[i] = fmt.Sprintf("(%s) %s", string(rune('a'+i)), c.label)
		}
		question = fmt.Sprintf("%s\n%s\n(enter a number)", lead, strings.Join(parts, "\n"))
		if !strings.Contains(question, fmt.Sprintf("%d", count)) {
			break
		}
	}
	notes := make([]string, len(chosen))
	for i, c := range chosen {
		verdict := "false"
		if c.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s, %s", c.label, verdict, c.note)
	}
	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the %d statements are true. %s", count, calcDraw, strings.Join(notes, "\n")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func continuityStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about continuity are true?", []calcVerdict{
		{"Every polynomial is continuous on R", true, "polynomials are built from sums and products of continuous functions"},
		{"sin x is continuous on R", true, "the sine function is continuous everywhere"},
		{"|x| is continuous at 0", true, "the two one-sided limits both equal zero"},
		{"1/(x^2+1) is continuous on R", true, "the denominator never vanishes"},
		{"Sums and products of continuous functions are continuous", true, "continuity is preserved by the usual operations"},
		{"Differentiability implies continuity", true, "a differentiable function has a limit equal to its value"},
		{"A function continuous at a point has its limit equal to its value there", true, "this is the definition of continuity at a point"},
		{"sqrt(x) is continuous on [0,infinity)", true, "the square root is continuous on its domain"},
		{"Every function that has a limit at a point is continuous there", false, "the limit may differ from the value or the value may be missing"},
		{"1/x is continuous at 0", false, "the function is not defined at zero"},
		{"The floor function is continuous at every integer", false, "it jumps at each integer"},
		{"Continuity implies differentiability", false, "|x| is continuous but not differentiable at zero"},
		{"A continuous function must be bounded", false, "x is continuous and unbounded on R"},
		{"A removable discontinuity makes a function discontinuous forever", false, "defining the missing value repairs it"},
		{"|x|/x is continuous at 0", false, "it has different one-sided limits"},
		{"tan x is continuous at pi/2", false, "tan is undefined there"},
	})
}

func littleOStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about little-o are true?", []calcVerdict{
		{"x^2 = o(x) as x tends to 0", true, "x^2/x tends to zero"},
		{"x = o(x^2) as x tends to infinity", true, "x/x^2 tends to zero"},
		{"x^3 = o(x^2) as x tends to 0", true, "x^3/x^2 tends to zero"},
		{"1/x = o(1) as x tends to infinity", true, "1/x tends to zero"},
		{"x^2 = o(x^3) as x tends to infinity", true, "x^2/x^3 tends to zero"},
		{"f = o(g) implies f = O(g)", true, "little-o is stronger than big-O"},
		{"f = o(g) means f/g tends to zero", true, "this is the definition"},
		{"x log x = o(x^2) as x tends to infinity", true, "log x grows slower than x"},
		{"x = o(x)", false, "the ratio is always one"},
		{"x^2 = o(x) as x tends to infinity", false, "the ratio x grows without bound"},
		{"Little-o and big-O mean the same thing", false, "little-o excludes a nonzero limiting ratio"},
		{"f = o(g) implies g = o(f)", false, "the relation is not symmetric"},
		{"A nonzero constant is o(1)", false, "its ratio to one does not tend to zero"},
		{"log x = o(1) as x tends to infinity", false, "log x grows without bound"},
		{"e^x = o(x) as x tends to infinity", false, "the exponential dominates every polynomial"},
		{"sin x = o(x) as x tends to 0", false, "sin x / x tends to one, not zero"},
	})
}

func discontinuityCountProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about discontinuities are true?", []calcVerdict{
		{"1/x has an infinite discontinuity at 0", true, "the function grows without bound near zero"},
		{"The floor function has a jump discontinuity at every integer", true, "the one-sided limits differ by one"},
		{"A removable discontinuity occurs when the limit exists but differs from the value", true, "the hole can be filled"},
		{"|x|/x has a jump discontinuity at 0", true, "the one-sided limits are -1 and 1"},
		{"sin(1/x) has an oscillatory discontinuity at 0", true, "it oscillates between -1 and 1"},
		{"x^2 is continuous everywhere", true, "it is a polynomial"},
		{"A discontinuity can be removed by redefining one value", true, "this is a removable discontinuity"},
		{"1/(x-2) has an infinite discontinuity at 2", true, "the denominator vanishes there"},
		{"Every discontinuity is removable", false, "jumps and infinite discontinuities are not"},
		{"A jump discontinuity has equal one-sided limits", false, "that would make it removable"},
		{"1/x is continuous at 0", false, "it is undefined at zero"},
		{"A removable discontinuity cannot be repaired", false, "the missing value can be supplied"},
		{"The floor function is continuous at every integer", false, "it jumps there"},
		{"sin(1/x) has a limit at 0", false, "it oscillates without settling"},
		{"A discontinuity always means the function is unbounded", false, "a jump discontinuity can be bounded"},
		{"A polynomial can have a jump discontinuity", false, "polynomials are continuous everywhere"},
	})
}

func seqConvergenceStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about sequence convergence are true?", []calcVerdict{
		{"1/n converges to 0", true, "the terms become arbitrarily small"},
		{"Every constant sequence converges", true, "it converges to the constant"},
		{"Every convergent sequence is bounded", true, "convergence forces the tail to stay near the limit"},
		{"(-1)^n diverges", true, "it alternates between -1 and 1"},
		{"(1+1/n)^n converges to e", true, "this is a standard limit"},
		{"A bounded monotone sequence converges", true, "the monotone convergence theorem"},
		{"A subsequence of a convergent sequence converges to the same limit", true, "subsequences inherit the limit"},
		{"n/(n+1) converges to 1", true, "divide numerator and denominator by n"},
		{"Every bounded sequence converges", false, "(-1)^n is bounded and divergent"},
		{"(-1)^n converges", false, "it has two subsequential limits"},
		{"1/n converges to 1", false, "it converges to zero"},
		{"Every convergent sequence is unbounded", false, "convergent sequences are bounded"},
		{"n/(n+1) converges to 0", false, "it converges to one"},
		{"Every divergent sequence is unbounded", false, "(-1)^n diverges and is bounded"},
		{"A sequence can converge to two different limits", false, "limits are unique"},
		{"sin n converges to 0", false, "sin n does not converge"},
	})
}

func limsupStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about limsup are true?", []calcVerdict{
		{"limsup is at least liminf", true, "the largest subsequential limit dominates the smallest"},
		{"The limsup of (-1)^n is 1", true, "one subsequence tends to one"},
		{"The limsup of 1/n is 0", true, "the sequence tends to zero"},
		{"The limsup of n is infinity", true, "the sequence is unbounded above"},
		{"A bounded sequence has a finite limsup", true, "boundedness bounds every subsequential limit"},
		{"A limsup is a subsequential limit", true, "some subsequence attains it"},
		{"The limsup of sin n is 1", true, "the values become arbitrarily close to one"},
		{"The limsup of -1/n is 0", true, "the sequence increases to zero"},
		{"The limsup always equals the limit", false, "they differ when the sequence oscillates"},
		{"The limsup of (-1)^n is -1", false, "the upper subsequential limit is one"},
		{"limsup is at most liminf", false, "the inequality goes the other way"},
		{"Every sequence has a finite limsup", false, "unbounded sequences have infinite limsup"},
		{"The limsup of 1/n is 1", false, "it is zero"},
		{"Limsup ignores subsequences", false, "it is defined through subsequential limits"},
		{"The limsup of n is 0", false, "it is infinity"},
		{"The limsup of sin n is -1", false, "the upper limit is one"},
	})
}

func oscillationStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about oscillation are true?", []calcVerdict{
		{"The oscillation of sin(1/x) at 0 is 2", true, "the values fill the interval from -1 to 1"},
		{"A function is continuous at a point exactly when its oscillation there is 0", true, "small oscillation is the epsilon-delta condition"},
		{"The oscillation of 1/x at 0 is infinite", true, "the function is unbounded near zero"},
		{"The oscillation of x^2 at 2 is 0", true, "x^2 is continuous at 2"},
		{"Oscillation is always nonnegative", true, "it is a supremum minus an infimum"},
		{"The oscillation of the floor function at an integer is 1", true, "the jump has size one"},
		{"A bounded function can have positive oscillation", true, "sin(1/x) is bounded with oscillation 2"},
		{"Oscillation measures the spread of values near a point", true, "it is the sup minus the inf"},
		{"Oscillation is always zero", false, "discontinuous functions have positive oscillation"},
		{"Continuous functions have positive oscillation", false, "continuity means zero oscillation"},
		{"The oscillation of sin(1/x) at 0 is 0", false, "the function oscillates between -1 and 1"},
		{"Oscillation can be negative", false, "a supremum minus an infimum cannot be negative"},
		{"The oscillation of x at 0 is 1", false, "x is continuous at zero"},
		{"Every discontinuity has infinite oscillation", false, "a jump discontinuity has finite oscillation"},
		{"Oscillation ignores nearby values", false, "it is defined by values near the point"},
		{"The oscillation of a constant function is 1", false, "it is zero"},
	})
}

func inverseHyperbolicStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about inverse hyperbolic functions are true?", []calcVerdict{
		{"The derivative of arcsinh x is 1/sqrt(x^2+1)", true, "differentiate the logarithmic form"},
		{"The domain of arccosh is [1,infinity)", true, "cosh only reaches values at least one"},
		{"The domain of arctanh is (-1,1)", true, "tanh only takes values between -1 and 1"},
		{"The derivative of arctanh x is 1/(1-x^2)", true, "this is the standard derivative"},
		{"arcsinh 0 = 0", true, "sinh 0 = 0"},
		{"arccosh 1 = 0", true, "cosh 0 = 1"},
		{"arctanh 0 = 0", true, "tanh 0 = 0"},
		{"Inverse hyperbolic functions have logarithmic forms", true, "each can be written with logarithms"},
		{"The derivative of arcsinh x is 1/sqrt(1-x^2)", false, "that is the derivative of arcsin"},
		{"The domain of arccosh is all real numbers", false, "cosh never falls below one"},
		{"The domain of arctanh is all real numbers", false, "it excludes the closed outside of (-1,1)"},
		{"The derivative of arctanh x is 1/(1+x^2)", false, "that is the derivative of arctan"},
		{"arcsinh 1 = 1", false, "arcsinh 1 is log(1+sqrt 2)"},
		{"arccosh 0 = 0", false, "zero is outside the domain"},
		{"arctanh 1 = 0", false, "arctanh is not defined at one"},
		{"arcsinh is bounded", false, "it grows like log x and is unbounded"},
	})
}

func convexJensenStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about convexity and Jensen's inequality are true?", []calcVerdict{
		{"f'' > 0 implies f is convex", true, "positive curvature means convex"},
		{"Jensen gives f((x+y)/2) <= (f(x)+f(y))/2 for convex f", true, "the midpoint value lies below the chord midpoint"},
		{"e^x is convex", true, "its second derivative is e^x, which is positive"},
		{"x^2 is convex", true, "its second derivative is 2"},
		{"Linear functions are both convex and concave", true, "the defining inequality holds with equality"},
		{"A convex function lies below the chord between two graph points", true, "this is the chord definition"},
		{"|x| is convex", true, "it is a piecewise linear function with a minimum"},
		{"log x is concave", true, "its second derivative is negative"},
		{"f'' < 0 implies f is convex", false, "negative curvature means concave"},
		{"Every convex function is linear", false, "x^2 is convex and not linear"},
		{"x^3 is convex on R", false, "its second derivative changes sign at zero"},
		{"Every convex function is bounded", false, "x^2 is convex and unbounded"},
		{"Jensen's inequality reverses for convex functions", false, "the inequality holds in the stated direction"},
		{"e^x is concave", false, "its second derivative is positive"},
		{"|x| is concave", false, "it is convex"},
		{"Convexity requires f'' > 0 everywhere", false, "a convex function need not be twice differentiable"},
	})
}

func epsilonDeltaStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about the epsilon-delta definition of the derivative are true?", []calcVerdict{
		{"Differentiability means the difference quotient has a limit", true, "the derivative is that limit"},
		{"The epsilon-delta form controls the difference quotient near the point", true, "it bounds the error from the derivative"},
		{"Differentiability implies continuity", true, "a differentiable function has a limit equal to its value"},
		{"The derivative of x^2 at 3 is 6", true, "the power rule gives 2x, and 2 times 3 is 6"},
		{"The derivative gives the linear approximation f(a+h) = f(a) + Lh + o(h)", true, "this is the first-order expansion"},
		{"Epsilon-delta estimates can prove differentiability", true, "they verify the limit definition"},
		{"|x| is not differentiable at 0", true, "the one-sided difference quotients differ"},
		{"The derivative at a point is unique", true, "a limit is unique"},
		{"Continuity implies differentiability", false, "|x| is continuous but not differentiable at zero"},
		{"The epsilon-delta derivative definition only tests continuity", false, "it controls the difference quotient"},
		{"|x| is differentiable at 0", false, "the left and right slopes disagree"},
		{"The derivative of x^2 at 3 is 9", false, "9 is the function value, not the derivative"},
		{"Differentiability means the function is linear", false, "it means a linear approximation exists locally"},
		{"Epsilon-delta methods cannot prove differentiability", false, "they are the definition"},
		{"A derivative can take two different values at one point", false, "limits are unique"},
		{"The difference quotient limit ignores the point", false, "it is taken as x approaches that point"},
	})
}
