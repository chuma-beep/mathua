package odes

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 21 adds four differential-equation concepts: coupled oscillators,
// eigenfunction expansions, resolvent kernels and the Hopf bifurcation.

func registerWave21(reg *generator.Registry) {
	reg.Register("ode.systems.coupled", &w21Gen{kind: "coupled"})
	reg.Register("ode.boundary.eigenfunction", &w21Gen{kind: "eigenfunction"})
	reg.Register("ode.integral.resolvent", &w21Gen{kind: "resolvent"})
	reg.Register("ode.nonlinear.hopf", &w21Gen{kind: "hopf"})
}

type w21Gen struct{ kind string }

func w21N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w21P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w21I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func (g *w21Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "coupled":
		return w21Coupled(kp)
	case "eigenfunction":
		return w21Eigenfunction(kp)
	case "resolvent":
		return w21Resolvent(kp)
	default:
		return w21Hopf(kp)
	}
}

func w21Coupled(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w21N(2, 30)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("A system has %d coupled oscillators, each with a position and a velocity. What is the size of the first-order system? (enter a number)", n),
			strconv.Itoa(2*n),
			fmt.Sprintf("Each of the %d oscillators contributes two states, so the system has %d equations.", n, 2*n), f)
	case 2:
		n := w21N(2, 20)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("The stiffness matrix of a chain of %d masses has how many entries? (enter a number)", n),
			strconv.Itoa(n*n),
			fmt.Sprintf("The symmetric stiffness matrix is %d by %d, so it has %d entries.", n, n, n*n), f)
	default:
		k := w21N(1, 30)
		w21I(f, "k", k)
		return w21P(kp,
			fmt.Sprintf("Two equal masses are joined by a spring of stiffness %d. What is the squared frequency of the relative mode? (enter a number)", k),
			strconv.Itoa(2*k),
			fmt.Sprintf("The relative mode has squared frequency 2k, which is %d.", 2*k), f)
	}
}

func w21Eigenfunction(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w21N(1, 30)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("On an interval of length L, how many interior zeros does the %d-th sine eigenfunction have? (enter a number)", n),
			strconv.Itoa(n-1),
			fmt.Sprintf("The %d-th sine mode completes %d half periods, so it has %d interior zeros.", n, n, n-1), f)
	case 2:
		terms := w21N(2, 40)
		w21I(f, "terms", terms)
		return w21P(kp,
			fmt.Sprintf("An eigenfunction expansion is truncated after %d terms. How many coefficients must be computed? (enter a number)", terms),
			strconv.Itoa(terms),
			fmt.Sprintf("One coefficient belongs to each term, so %d are needed.", terms), f)
	default:
		val := w21N(1, 30)
		w21I(f, "val", val)
		return w21P(kp,
			fmt.Sprintf("For sin(n pi x/L) on [0,L], the integral of the square equals %d. The integral is L/2. What is L? (enter a number)", val),
			strconv.Itoa(2*val),
			fmt.Sprintf("The integral is L/2 = %d, so L = %d.", val, 2*val), f)
	}
}

func w21Resolvent(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w21N(2, 30)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("A constant kernel has lambda = 1/%d. The resolvent factor is n/(n-1). What is n? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("The parameter n is the denominator of lambda, so n = %d.", n), f)
	case 2:
		terms := w21N(1, 40)
		w21I(f, "terms", terms)
		return w21P(kp,
			fmt.Sprintf("A Neumann series is truncated after %d terms. How many applications of the integral operator does the last term contain? (enter a number)", terms),
			strconv.Itoa(terms),
			fmt.Sprintf("The N-th term applies the operator N times, so the last term has %d applications.", terms), f)
	default:
		n := w21N(2, 30)
		w21I(f, "n", n)
		return w21P(kp,
			fmt.Sprintf("The Fredholm determinant is 1 - 1/%d. What is its numerator over the denominator %d? (enter a number)", n, n),
			strconv.Itoa(n-1),
			fmt.Sprintf("1 - 1/%d = %d/%d, so the numerator is %d.", n, n-1, n, n-1), f)
	}
}

func w21Hopf(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		k := w21N(1, 30)
		w21I(f, "k", k)
		return w21P(kp,
			fmt.Sprintf("A Hopf bifurcation occurs when the parameter reaches %d. At what parameter value does the steady state lose stability? (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("The eigenvalue pair crosses the imaginary axis at the parameter value %d.", k), f)
	case 2:
		sq := w21N(1, 12)
		freq2 := sq * sq
		w21I(f, "sq", sq)
		return w21P(kp,
			fmt.Sprintf("At the Hopf crossing the squared frequency is %d. What is the integer frequency? (enter a number)", freq2),
			strconv.Itoa(sq),
			fmt.Sprintf("The squared frequency is %d, so the frequency is %d.", freq2, sq), f)
	default:
		d := w21N(1, 30)
		w21I(f, "d", d)
		return w21P(kp,
			fmt.Sprintf("The parameter is %d above the Hopf value. The squared amplitude equals that distance. What is the squared amplitude? (enter a number)", d),
			strconv.Itoa(d),
			fmt.Sprintf("The squared amplitude equals the parameter gap, which is %d.", d), f)
	}
}
