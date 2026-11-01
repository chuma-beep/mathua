package odes

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The differential-equations batch used to ask recognition questions. These
// builders replace them with count problems over a pool of statements. Each
// pool holds sixteen statements, eight true and eight false, and the answer is
// drawn uniformly from zero to eight before the statements are chosen.

type odeVerdict struct {
	label string
	yes   bool
	note  string
}

const odeDraw = 8

func countOdeVerdicts(lead string, pool []odeVerdict) generator.Problem {
	var yes, no []odeVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}

	question := ""
	chosen := make([]odeVerdict, 0, odeDraw)
	count := 0
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(odeDraw + 1)
		yesCopy := append([]odeVerdict(nil), yes...)
		noCopy := append([]odeVerdict(nil), no...)
		rand.Shuffle(len(yesCopy), func(i, j int) { yesCopy[i], yesCopy[j] = yesCopy[j], yesCopy[i] })
		rand.Shuffle(len(noCopy), func(i, j int) { noCopy[i], noCopy[j] = noCopy[j], noCopy[i] })
		chosen = append(chosen[:0], yesCopy[:count]...)
		chosen = append(chosen, noCopy[:odeDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

		parts := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
		}
		question = fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if !strings.Contains(question, fmt.Sprintf("%d", count)) {
			break
		}
	}

	notes := make([]string, len(chosen))
	for i, v := range chosen {
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes[i] = fmt.Sprintf("%s: %s. %s", v.label, verdict, v.note)
	}

	return generator.Problem{
		Question:    question,
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the %d statements are true. %s", count, odeDraw, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func bernoulliProblem() generator.Problem {
	pool := []odeVerdict{
		{"A Bernoulli equation has the form y' + P(x)y = Q(x)y^n with n not equal to 0 or 1.", true, "The power n is what makes the equation nonlinear."},
		{"The substitution v = y^(1-n) linearizes a Bernoulli equation.", true, "Dividing by y^n and substituting turns the nonlinear term into a linear one."},
		{"For n = 2 the Bernoulli substitution is v = 1/y.", true, "One minus n is minus one, so v is the reciprocal of y."},
		{"After the Bernoulli substitution the equation becomes first-order linear.", true, "That is the purpose of the substitution."},
		{"The equation y' + y = y^2 has Bernoulli exponent n = 2.", true, "The nonlinear factor is y squared."},
		{"A Bernoulli equation is nonlinear when n is neither 0 nor 1.", true, "Only those two exponents collapse to a linear equation."},
		{"Solving a Bernoulli equation usually divides by y^n.", true, "The division exposes the derivative of v."},
		{"The solution of a Bernoulli equation may be given implicitly.", true, "The linear equation in v can still be hard to invert."},
		{"Every first-order equation is a Bernoulli equation.", false, "Bernoulli form requires the right side to be Q(x) times a power of y."},
		{"A Bernoulli equation with n = 1 is nonlinear.", false, "At n = 1 the equation is already linear."},
		{"The Bernoulli substitution is v = y^n.", false, "The correct exponent is one minus n."},
		{"The Bernoulli substitution produces a second-order equation.", false, "It produces a first-order linear equation."},
		{"The equation y' + xy = sin x is a Bernoulli equation.", false, "The right side has no power of y."},
		{"A Bernoulli equation with n = 0 is nonlinear.", false, "At n = 0 the equation is linear."},
		{"The Bernoulli substitution is v = e^y.", false, "That substitution belongs to a different equation family."},
		{"The Bernoulli substitution removes the dependence on x.", false, "It removes the nonlinear power of y, not the x dependence."},
	}
	return countOdeVerdicts("How many of these statements about Bernoulli equations are true?", pool)
}

func reductionProblem() generator.Problem {
	pool := []odeVerdict{
		{"Reduction of order needs one known nonzero solution.", true, "The method trades that solution for a procedure that finds the second."},
		{"Reduction of order sets y = y1 v.", true, "The known solution multiplies an unknown function."},
		{"Reduction of order produces a first-order equation for v'.", true, "The y1 v term cancels because y1 solves the homogeneous equation."},
		{"For y'' - y = 0 with y1 = e^x, the second solution is proportional to e^{-x}.", true, "The two independent solutions are e^x and e^{-x}."},
		{"The second solution formula is y2 = y1 times the integral of exp(-integral P) divided by y1 squared.", true, "This is the standard reduction-of-order formula."},
		{"The Wronskian of two independent solutions is nonzero on the interval.", true, "A zero Wronskian would make the solutions dependent."},
		{"The known solution y1 must be nonzero on the interval.", true, "The reduced equation divides by y1."},
		{"Reduction of order applies to linear homogeneous second-order equations.", true, "That is the setting in which the y1 v substitution works."},
		{"Reduction of order needs two known solutions.", false, "One known solution is enough."},
		{"Reduction of order sets y = y1 + v.", false, "The unknown function multiplies y1."},
		{"Reduction of order leaves a second-order equation in v.", false, "The equation drops to first order in v'."},
		{"Reduction of order applies directly to nonlinear equations.", false, "It relies on linearity and superposition."},
		{"Reduction of order finds the characteristic roots.", false, "That is the constant-coefficient method."},
		{"The second solution can always be taken equal to y1.", false, "The two solutions must be independent."},
		{"The Wronskian of two independent solutions is zero.", false, "Independence is equivalent to a nonzero Wronskian."},
		{"The known solution y1 may vanish on the interval.", false, "Division by y1 would fail where it vanishes."},
	}
	return countOdeVerdicts("How many of these statements about reduction of order are true?", pool)
}

func matrixExpProblem() generator.Problem {
	pool := []odeVerdict{
		{"The matrix exponential is defined by the power series sum of (At)^n divided by n factorial.", true, "The definition copies the scalar exponential series."},
		{"For a diagonal matrix the exponential exponentiates each diagonal entry.", true, "Powers of a diagonal matrix act entry by entry."},
		{"The derivative of e^(At) with respect to t is A e^(At).", true, "Differentiating the series produces A times the series."},
		{"e^(At) solves x' = Ax with x(0) equal to the identity.", true, "The derivative is A e^(At), and the value at zero is the identity."},
		{"If A is nilpotent the exponential series terminates.", true, "Some power of A is zero, so only finitely many terms survive."},
		{"If A and B commute then e^(A+B) equals e^A e^B.", true, "The usual exponent law needs commutativity."},
		{"The inverse of e^A is e^(-A).", true, "The two exponentials multiply to the identity."},
		{"The determinant of e^A equals exp of the trace of A.", true, "This is Jacobi's formula for the matrix exponential."},
		{"The matrix exponential is computed by exponentiating each entry.", false, "That only works for diagonal matrices."},
		{"Every matrix exponential is diagonalizable.", false, "A nilpotent matrix gives a non-diagonal exponential."},
		{"e^(A+B) always equals e^A e^B.", false, "The identity can fail when A and B do not commute."},
		{"The derivative of e^(At) is e^(At) A^(-1).", false, "The derivative is A e^(At)."},
		{"The matrix exponential may be singular.", false, "Its inverse is e^(-A), so it is always invertible."},
		{"If A squared is zero then e^A is zero.", false, "The series becomes the identity plus A."},
		{"e^(At) solves x' = A^(-1)x.", false, "It solves x' = Ax."},
		{"The determinant of e^A equals the trace of A.", false, "The determinant is the exponential of the trace."},
	}
	return countOdeVerdicts("How many of these statements about the matrix exponential are true?", pool)
}

func stepProblem() generator.Problem {
	pool := []odeVerdict{
		{"The Laplace transform of the unit step at a is e^(-a s)/s.", true, "This is the time-shift rule applied to 1/s."},
		{"The Laplace transform of the Dirac delta at a is e^(-a s).", true, "The delta samples the exponential kernel at t = a."},
		{"The derivative of the unit step is the Dirac delta.", true, "The step jumps at a, and the delta is the distributional derivative of that jump."},
		{"A step function models a switch turned on at time a.", true, "It is zero before a and one after."},
		{"The impulse response is the solution when the input is a delta.", true, "Its transform is the transfer function."},
		{"The Laplace transform of the delta at zero is one.", true, "Setting a equal to zero gives e^0."},
		{"The unit step at a is zero before a and one after.", true, "That is the definition of the Heaviside step."},
		{"The response to a delayed impulse is the impulse response shifted in time.", true, "Time invariance moves the response with the input."},
		{"The Laplace transform of every unit step is 1/s.", false, "A step at a carries the factor e^(-a s)."},
		{"The Laplace transform of the Dirac delta is 1/s.", false, "The delta transforms to e^(-a s), which is one at a equal to zero."},
		{"The derivative of the unit step is zero everywhere.", false, "The jump produces a delta at a."},
		{"A step function models a smooth sinusoid.", false, "It is piecewise constant with a jump."},
		{"The impulse response is the solution with zero input.", false, "It is the response to a unit impulse."},
		{"The Dirac delta is an ordinary function.", false, "It is a distribution defined by its action on test functions."},
		{"The step at zero transforms to e^(-s).", false, "The step at zero transforms to 1/s."},
		{"The derivative of the Dirac delta is the unit step.", false, "The relation goes the other way."},
	}
	return countOdeVerdicts("How many of these statements about steps and impulses are true?", pool)
}
