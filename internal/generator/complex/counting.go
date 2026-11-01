package complex

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// A complex analysis statement is either true or false, and the learner counts
// how many of eight are true. The answer is a number the grader checks, so the
// exercise tests the property rather than the recognition of a wording.
type complexVerdict struct {
	label string
	yes   bool
	note  string
}

const complexDraw = 8

func countComplexVerdicts(lead string, pool []complexVerdict) generator.Problem {
	var yes, no []complexVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	for attempt := 0; attempt < 40; attempt++ {
		count := rand.Intn(complexDraw + 1)
		if count > len(yes) {
			count = len(yes)
		}
		if complexDraw-count > len(no) {
			count = complexDraw - len(no)
		}
		ys := append([]complexVerdict(nil), yes...)
		ns := append([]complexVerdict(nil), no...)
		rand.Shuffle(len(ys), func(i, j int) { ys[i], ys[j] = ys[j], ys[i] })
		rand.Shuffle(len(ns), func(i, j int) { ns[i], ns[j] = ns[j], ns[i] })
		chosen := append([]complexVerdict(nil), ys[:count]...)
		chosen = append(chosen, ns[:complexDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

		parts := make([]string, len(chosen))
		notes := make([]string, len(chosen))
		for i, v := range chosen {
			parts[i] = fmt.Sprintf("(%c) %s", 'a'+i, v.label)
			verdict := "false"
			if v.yes {
				verdict = "true"
			}
			notes[i] = fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note)
		}
		q := fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " "))
		if strings.Contains(q, strconv.Itoa(count)) {
			continue
		}
		return generator.Problem{
			Question:    q,
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("Exactly %d of the 8 statements are true. %s", count, strings.Join(notes, "; ")),
			Facts:       map[string]string{"count": strconv.Itoa(count)},
		}
	}
	return generator.Problem{
		Question:    lead + " (enter a number)",
		Answer:      "0",
		Explanation: "No statement in this draw is true.",
		Facts:       map[string]string{"count": "0"},
	}
}

func eulerIdentityProblem() generator.Problem {
	pool := []complexVerdict{
		{"e^(i pi) = -1", true, "Euler's identity at theta = pi"},
		{"e^(i pi/2) = i", true, "a quarter turn"},
		{"e^(i pi/3) = 1/2 + (sqrt 3/2) i", true, "cos and sin of pi/3"},
		{"e^(i pi/4) = (sqrt 2/2) + (sqrt 2/2) i", true, "cos and sin of pi/4"},
		{"e^(i pi/6) = (sqrt 3/2) + (1/2) i", true, "cos and sin of pi/6"},
		{"e^(i 3pi/2) = -i", true, "three quarter turns"},
		{"e^(i 2pi) = 1", true, "a full turn"},
		{"|e^(i theta)| = 1 for every real theta", true, "points on the unit circle"},
		{"e^(i pi) = 1", false, "it is -1"},
		{"e^(i pi/2) = -1", false, "it is i"},
		{"e^(i pi/4) = 1 + i", false, "the modulus would be sqrt 2"},
		{"e^(i pi/3) = 1/2 + (1/2) i", false, "the sine is sqrt 3/2"},
		{"e^(i pi/6) = (sqrt 3/2) + (sqrt 3/2) i", false, "the cosine and sine are different"},
		{"e^(i 3pi/2) = i", false, "it is -i"},
		{"e^(i 2pi) = -1", false, "a full turn returns 1"},
		{"e^(i theta) is real for every real theta", false, "only at multiples of pi"},
	}
	return countComplexVerdicts("How many of these statements about Euler's formula are true?", pool)
}

func cauchyRiemannProblem() generator.Problem {
	pool := []complexVerdict{
		{"u_x = v_y", true, "first Cauchy-Riemann equation"},
		{"u_y = -v_x", true, "second Cauchy-Riemann equation"},
		{"the equations are necessary for differentiability", true, "the derivative would not exist otherwise"},
		{"the equations plus continuous partials are sufficient", true, "the standard sufficient condition"},
		{"an analytic function satisfies the equations", true, "analytic means complex differentiable"},
		{"f(z) = z^2 satisfies the equations", true, "u = x^2 - y^2 and v = 2xy"},
		{"f(z) = e^z satisfies the equations", true, "u = e^x cos y and v = e^x sin y"},
		{"real and imaginary parts of an analytic function are harmonic", true, "a consequence of the equations"},
		{"u_x = -v_y", false, "the sign is wrong"},
		{"u_y = v_x", false, "the second equation has a minus"},
		{"the equations alone are sufficient", false, "the partials must also be continuous"},
		{"f(z) = z bar satisfies the equations", false, "its v has the opposite sign"},
		{"f(z) = |z|^2 satisfies the equations everywhere", false, "only at the origin"},
		{"every continuous function satisfies the equations", false, "continuity is far too weak"},
		{"the equations imply the function is a polynomial", false, "e^z is not a polynomial"},
		{"the equations are needed only at the origin", false, "they hold on the domain of differentiability"},
	}
	return countComplexVerdicts("How many of these statements about the Cauchy-Riemann equations are true?", pool)
}

func residueProblem() generator.Problem {
	pool := []complexVerdict{
		{"Res(1/(z-a), a) = 1", true, "a simple pole with unit numerator"},
		{"Res(1/(z-a)^2, a) = 0", true, "no 1/(z-a) term"},
		{"Res(e^z, 0) = 1", true, "the constant term"},
		{"Res(sin z / z, 0) = 1", true, "the removable value"},
		{"the sum of all residues on the Riemann sphere is 0", true, "the residue at infinity balances them"},
		{"Res(1/(z^2+1), i) = 1/(2i)", true, "a simple pole at i"},
		{"a simple pole has residue lim (z-a) f(z)", true, "the standard formula"},
		{"the residue theorem multiplies the sum by 2 pi i", true, "for a positively oriented contour"},
		{"Res(1/(z-a), a) = 0", false, "the coefficient of 1/(z-a) is 1"},
		{"Res(1/(z-a)^2, a) = 1", false, "the coefficient is 0"},
		{"Res(e^z, 0) = 0", false, "e^z has constant term 1"},
		{"the residue theorem omits the factor 2 pi i", false, "the factor is essential"},
		{"residues depend on the chosen contour", false, "the residue is a local quantity"},
		{"Res(1/z, 0) = 0", false, "it is 1"},
		{"the residue is the coefficient a_1", false, "it is a_{-1}"},
		{"a removable singularity has residue 1", false, "it has residue 0"},
	}
	return countComplexVerdicts("How many of these statements about residues are true?", pool)
}

func analyticProblem() generator.Problem {
	pool := []complexVerdict{
		{"every polynomial is entire", true, "polynomials are differentiable everywhere"},
		{"e^z is entire", true, "the exponential series converges everywhere"},
		{"sin z is entire", true, "its power series has infinite radius"},
		{"cos z is entire", true, "its power series has infinite radius"},
		{"1/z is analytic on C minus {0}", true, "the quotient rule applies away from 0"},
		{"a rational function is analytic where its denominator is nonzero", true, "quotients of analytic functions"},
		{"z^3 + 2z is entire", true, "a polynomial"},
		{"a composition of analytic functions is analytic", true, "the chain rule"},
		{"z bar is analytic", false, "it reverses orientation"},
		{"|z|^2 is analytic", false, "it satisfies Cauchy-Riemann only at 0"},
		{"Re z is analytic", false, "its imaginary part is 0 and Cauchy-Riemann fails"},
		{"1/z is analytic at 0", false, "it has a pole there"},
		{"sqrt z is analytic on the whole plane", false, "it needs a branch cut"},
		{"arg z is analytic on C minus {0}", false, "it is not even continuous around the origin"},
		{"the conjugate function is entire", false, "it is not complex differentiable"},
		{"every continuous function is analytic", false, "analyticity is much stronger"},
	}
	return countComplexVerdicts("How many of these functions are analytic on the stated domain?", pool)
}

func laurentProblem() generator.Problem {
	pool := []complexVerdict{
		{"a Laurent series converges on an annulus", true, "the natural domain of a Laurent expansion"},
		{"the principal part determines the type of singularity", true, "negative powers carry the information"},
		{"a removable singularity has no negative powers", true, "the series is a Taylor series"},
		{"a pole has a finite principal part", true, "the negative powers stop"},
		{"an essential singularity has infinitely many negative powers", true, "the principal part never stops"},
		{"the residue is the coefficient a_{-1}", true, "the definition"},
		{"a Laurent expansion is unique on its annulus", true, "coefficients are determined"},
		{"an annulus gives an annular region of convergence", true, "the series converges there"},
		{"a Taylor series covers every singularity", false, "it fails outside the largest disc"},
		{"positive powers form the principal part", false, "negative powers do"},
		{"a removable singularity has a pole", false, "these are different types"},
		{"an essential singularity has a finite principal part", false, "it has infinitely many terms"},
		{"the residue is the coefficient a_0", false, "it is a_{-1}"},
		{"a Laurent series is unique only on a disc", false, "uniqueness holds on the annulus"},
		{"the pole order is read from a_0", false, "it is read from the first nonzero negative power"},
		{"every Laurent series defines an entire function", false, "negative powers diverge at the centre"},
	}
	return countComplexVerdicts("How many of these statements about Laurent series are true?", pool)
}

func contourProblem() generator.Problem {
	pool := []complexVerdict{
		{"the integral of dz/z around |z|=1 is 2 pi i", true, "winding number one"},
		{"the integral of z^n around a closed contour is 0 for n not equal to -1", true, "an antiderivative exists"},
		{"the integral of dz/z^2 around |z|=1 is 0", true, "the exponent is not -1"},
		{"contour integrals are additive over adjacent arcs", true, "split the path"},
		{"reversing orientation changes the sign", true, "the integral is oriented"},
		{"Cauchy's theorem gives 0 for a function analytic inside the contour", true, "no singularities inside"},
		{"an integral is unchanged by deformation avoiding singularities", true, "Cauchy's deformation theorem"},
		{"the ML estimate bounds an integral by length times maximum modulus", true, "the standard estimate"},
		{"the integral of dz/z around |z|=1 is 0", false, "it is 2 pi i"},
		{"the integral of dz/z is independent of orientation", false, "orientation changes the sign"},
		{"the integral of z^2 around |z|=1 is 2 pi i", false, "it is 0"},
		{"a contour integral is always real", false, "it is generally complex"},
		{"Cauchy's theorem applies through a singularity", false, "the singularity must be outside"},
		{"every closed contour has integral 0", false, "only when the integrand is analytic inside"},
		{"the integral of 1/z around a circle of radius 2 is 0", false, "the winding number is still one"},
		{"deforming a contour across a pole preserves the integral", false, "the residue changes the value"},
	}
	return countComplexVerdicts("How many of these statements about contour integrals are true?", pool)
}

func harmonicProblem() generator.Problem {
	pool := []complexVerdict{
		{"u = x^2 - y^2 is harmonic", true, "the real part of z^2"},
		{"u = x^3 - 3xy^2 is harmonic", true, "the real part of z^3"},
		{"u = e^x cos y is harmonic", true, "the real part of e^z"},
		{"u = ln|z| is harmonic away from 0", true, "the real part of log z"},
		{"the real part of an analytic function is harmonic", true, "Cauchy-Riemann gives Laplace's equation"},
		{"a harmonic function satisfies u_xx + u_yy = 0", true, "the definition"},
		{"a harmonic function has a local harmonic conjugate", true, "on a simply connected domain"},
		{"u = 2xy is harmonic", true, "the imaginary part of z^2"},
		{"u = x^2 + y^2 is harmonic", false, "its Laplacian is 4"},
		{"u = |z|^2 is harmonic", false, "the same failure in another form"},
		{"u = x^2 y is harmonic", false, "u_xx + u_yy = 2y"},
		{"u = x^3 + y^3 is harmonic", false, "the Laplacian is 6x + 6y"},
		{"u = x y^2 is harmonic", false, "the Laplacian is 2x"},
		{"u = x^4 + y^4 is harmonic", false, "the Laplacian is 12x^2 + 12y^2"},
		{"every harmonic function is the real part of a polynomial", false, "ln|z| is not a polynomial"},
		{"a harmonic function must be bounded", false, "u = x is harmonic and unbounded"},
	}
	return countComplexVerdicts("How many of these statements about harmonic functions are true?", pool)
}

func riemannSphereProblem() generator.Problem {
	pool := []complexVerdict{
		{"the Riemann sphere is C together with one point at infinity", true, "the one-point compactification"},
		{"the Riemann sphere is compact", true, "it is a closed bounded surface"},
		{"the Riemann sphere is homeomorphic to S^2", true, "stereographic projection"},
		{"exactly one point is added to C", true, "the point at infinity"},
		{"the map 1/z sends 0 to infinity", true, "the reciprocal swaps the two points"},
		{"neighbourhoods of infinity are complements of compact sets", true, "the one-point topology"},
		{"every Mobius map extends to the sphere", true, "Mobius maps act on the extended plane"},
		{"the Riemann sphere is connected", true, "it is a sphere"},
		{"two points are added to C", false, "only one"},
		{"the Riemann sphere is not compact", false, "compactness is the point of the construction"},
		{"the Riemann sphere is homeomorphic to a torus", false, "the torus has a hole"},
		{"the map 1/z sends 0 to 0", false, "it sends 0 to infinity"},
		{"infinity is an isolated point of C", false, "infinity is not in C"},
		{"the Riemann sphere is disconnected", false, "it is a connected surface"},
		{"C itself is compact", false, "C is not compact"},
		{"a circle is added at infinity", false, "a single point is added"},
	}
	return countComplexVerdicts("How many of these statements about the Riemann sphere are true?", pool)
}
