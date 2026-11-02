package complex

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The 2026 content expansion adds nine analysis concepts to the complex
// numbers domain. Each generator is section aware: KP 1, 2 and 3 ask
// different questions even when they use the same underlying idea.

var complexBridgeKinds = map[string]string{
	"complex.integral.cauchy_formula": "cauchy_formula",
	"complex.adv.liouville":           "liouville",
	"complex.adv.maximum_modulus":     "maximum_modulus",
	"complex.integral.real_integrals": "real_integrals",
	"complex.adv.argument_principle":  "argument_principle",
	"complex.adv.winding":             "winding",
	"complex.adv.riemann_mapping":     "riemann_mapping",
	"complex.adv.identity_theorem":    "identity_theorem",
	"complex.adv.morera":              "morera",
}

type complexBridgeGen struct{ kind string }

func (g *complexBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "cauchy_formula":
		return cauchyFormulaProblem(o)
	case "liouville":
		return liouvilleProblem(o)
	case "maximum_modulus":
		return maximumModulusProblem(o)
	case "real_integrals":
		return realIntegralsProblem(o)
	case "argument_principle":
		return argumentPrincipleProblem(o)
	case "winding":
		return windingProblem(o)
	case "riemann_mapping":
		return riemannMappingProblem(o)
	case "identity_theorem":
		return identityTheoremProblem(o)
	case "morera":
		return moreraProblem(o)
	}
	return complexBridgeProblem(o, "Compute 2 + 2. (enter a number)", "4", "Two plus two is four.")
}

func complexBridgeProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer},
		KPOrdinal:   o,
	}
}

func complexBridgeDecimal(v float64) string {
	s := strconv.FormatFloat(v, 'f', 4, 64)
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	if s == "-0" {
		return "0"
	}
	return s
}

func complexBridgePow(a, n int) int {
	r := 1
	for i := 0; i < n; i++ {
		r *= a
	}
	return r
}

func complexBridgeCount(o int, lead string, pool []complexVerdict) generator.Problem {
	p := countComplexVerdicts(lead, pool)
	p.KPOrdinal = o
	return p
}

func complexPool(yes, no []string) []complexVerdict {
	pool := make([]complexVerdict, 0, len(yes)+len(no))
	for _, s := range yes {
		pool = append(pool, complexVerdict{s, true, "standard example"})
	}
	for _, s := range no {
		pool = append(pool, complexVerdict{s, false, "standard counterexample"})
	}
	return pool
}

var liouvilleFunctionsPool = complexPool(
	[]string{"f(z) = 3", "f(z) = 0", "f(z) = -2 + i", "f(z) = 7", "f(z) = pi", "f(z) = -i", "f(z) = 4 - 2i", "f(z) = 1/2"},
	[]string{"f(z) = z", "f(z) = z^2", "f(z) = e^z", "f(z) = sin z", "f(z) = cos z", "f(z) = z^3 + z", "f(z) = e^z + 1", "f(z) = sinh z"},
)

var riemannDomainsPool = complexPool(
	[]string{"the open unit disc", "the upper half-plane", "the open square (0,1) x (0,1)", "the strip 0 < Im z < 1", "the plane slit along the negative real axis", "the interior of a triangle", "the half-plane Re z > 0", "the open ball of radius 2"},
	[]string{"the punctured unit disc", "the annulus 1 < |z| < 2", "the plane minus two points", "the unit disc minus 0", "the plane minus three points", "the annulus 2 < |z| < 5", "the punctured plane", "the plane minus the integers"},
)

var riemannStatementsPool = complexPool(
	[]string{"Every simply connected proper domain is conformally equivalent to the unit disc", "The map is unique up to an automorphism of the disc", "Automorphisms of the disc are Mobius transformations", "The upper half-plane is conformally equivalent to the disc", "The theorem requires the domain to be simply connected", "A conformal map preserves angles", "The unit disc is simply connected", "A simply connected domain has no holes"},
	[]string{"Every domain is conformally equivalent to the disc", "The map is unique without normalization", "An annulus is conformally equivalent to the disc", "The theorem applies to the whole plane", "Conformal maps need not be injective", "A punctured disc is simply connected", "The theorem gives an explicit formula for every domain", "Conformal equivalence preserves area"},
)

var identityFunctionsPool = complexPool(
	[]string{"zeros at 1/n for all n", "zeros at 1/n^2 for all n", "zeros at every rational point", "zero on an open subset", "zeros with an accumulation point in the domain", "zero on a line segment inside the domain", "zero at every point of a convergent sequence", "zero on a small disc inside the domain"},
	[]string{"exactly one zero at 0", "zeros at 1 and 2", "zeros at the integers", "zero at a single point on the boundary", "zeros at 2 and 3 only", "zero at 0 of order 2", "zeros at 1/n for n = 1 to 5 only", "zero on the boundary circle only"},
)

var moreraFunctionsPool = complexPool(
	[]string{"z^2", "e^z", "sin z", "z^3 + z", "1/(z-5) on |z| < 2", "cos z", "z/(z-5) on |z| < 2", "e^z + z^2"},
	[]string{"|z|^2", "conjugate z", "Re z", "Im z", "1/z on |z| < 1", "|z|", "conjugate z squared", "Re(z^2)"},
)

var moreraStatementsPool = complexPool(
	[]string{"Morera's theorem is the converse of Cauchy's theorem", "Continuity and zero triangle integrals imply analyticity", "Analytic functions have zero integrals over triangles", "Morera can prove that a uniform limit of analytic functions is analytic", "The triangle condition implies a primitive on a convex domain", "Morera's theorem assumes continuity", "A continuous function with zero triangle integrals is analytic", "Morera is used to prove convergence theorems"},
	[]string{"Morera requires differentiability", "Morera applies only to polynomials", "The triangle may be replaced by any set", "Morera is the converse of Liouville's theorem", "The function must be bounded", "Morera requires the domain to be simply connected", "Zero triangle integrals alone imply boundedness", "Morera characterizes entire functions only"},
)

var identityStatementsPool = complexPool(
	[]string{"A nonzero analytic function has isolated zeros", "Agreement on a set with an accumulation point implies equality", "Taylor coefficients determine an analytic function", "Analytic continuation along homotopic paths is unique", "The zero set of a nonzero analytic function has empty interior", "An analytic function vanishing on an open set is identically zero", "The identity theorem applies on connected domains", "A power series is determined by its values on any set with an accumulation point"},
	[]string{"A nonzero analytic function can vanish on an open set", "Two analytic functions agreeing at one point are equal", "The identity theorem holds only for polynomials", "Zeros of a nonzero analytic function can accumulate inside the domain", "An analytic function is determined by its boundary values alone", "Agreement on a boundary arc always forces agreement inside", "The identity theorem applies to disconnected domains without extra conditions", "A nonzero analytic function can have uncountably many zeros in a compact set"},
)

func cauchyFormulaProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 2 + rand.Intn(5)
		a := 1 + rand.Intn(4)
		v := complexBridgePow(a, n)
		return complexBridgeProblem(o, fmt.Sprintf("Let f(z) = z^%d. What is f(%d)? (enter a number)", n, a), strconv.Itoa(v), fmt.Sprintf("%d^%d = %d.", a, n, v))
	case 2:
		n := 1 + rand.Intn(5)
		a := 1 + rand.Intn(3)
		r := a + 1 + rand.Intn(3)
		v := complexBridgePow(a, n)
		return complexBridgeProblem(o, fmt.Sprintf("Evaluate (1/(2 pi i)) times the integral of z^%d/(z-%d) around |z| = %d. (enter a number)", n, a, r), strconv.Itoa(v), fmt.Sprintf("The Cauchy formula gives f(%d) = %d^%d = %d.", a, a, n, v))
	default:
		n := 2 + rand.Intn(5)
		a := 1 + rand.Intn(4)
		v := n * complexBridgePow(a, n-1)
		return complexBridgeProblem(o, fmt.Sprintf("Let f(z) = z^%d. What is f'(%d)? (enter a number)", n, a), strconv.Itoa(v), fmt.Sprintf("f'(z) = %d z^%d, so f'(%d) = %d.", n, n-1, a, v))
	}
}

func liouvilleProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(8)
		c := 2 + rand.Intn(8)
		return complexBridgeProblem(o, fmt.Sprintf("An entire function satisfies |f(z)| <= %d|z|^%d for large |z|. What is the largest possible degree of f? (enter a number)", c, n), strconv.Itoa(n), fmt.Sprintf("Polynomial growth of order %d forces a polynomial of degree at most %d.", n, n))
	case 2:
		return complexBridgeCount(o, "How many of these functions are bounded entire functions?", liouvilleFunctionsPool)
	default:
		a := rand.Intn(7) - 3
		b := rand.Intn(9) - 4
		c := rand.Intn(7) - 3
		return complexBridgeProblem(o, fmt.Sprintf("A bounded entire function satisfies f(%d) = %d. What is f(%d)? (enter a number)", a, b, c), strconv.Itoa(b), "A bounded entire function is constant, so every value equals the given value.")
	}
}

func maximumModulusProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(5)
		r := 2 + rand.Intn(4)
		v := complexBridgePow(r, n)
		return complexBridgeProblem(o, fmt.Sprintf("What is the maximum of |z^%d| on the circle |z| = %d? (enter a number)", n, r), strconv.Itoa(v), fmt.Sprintf("On |z| = %d, |z^%d| = %d^%d = %d.", r, n, r, n, v))
	case 2:
		a := 1 + rand.Intn(4)
		b := rand.Intn(6)
		n := 1 + rand.Intn(4)
		r := 2 + rand.Intn(3)
		v := a*complexBridgePow(r, n) + b
		return complexBridgeProblem(o, fmt.Sprintf("What is the maximum of |%d z^%d + %d| on the disc |z| <= %d? (enter a number)", a, n, b, r), strconv.Itoa(v), fmt.Sprintf("The triangle inequality is sharp on the boundary: %d times %d^%d plus %d equals %d.", a, r, n, b, v))
	default:
		r := 1 + rand.Intn(10)
		v := math.Exp(float64(r))
		return complexBridgeProblem(o, fmt.Sprintf("What is the maximum of |e^z| on the disc |z| <= %d? (enter a decimal)", r), complexBridgeDecimal(v), fmt.Sprintf("The maximum occurs at z = %d, and e^%d is about %s.", r, r, complexBridgeDecimal(v)))
	}
}

func realIntegralsProblem(o int) generator.Problem {
	switch o {
	case 1:
		type triple struct{ a, b, h int }
		ts := []triple{{5, 3, 4}, {5, 4, 3}, {13, 5, 12}, {13, 12, 5}, {10, 6, 8}, {17, 8, 15}}
		t := ts[rand.Intn(len(ts))]
		v := 2 / float64(t.h)
		return complexBridgeProblem(o, fmt.Sprintf("The integral from 0 to 2 pi of d theta/(%d + %d cos theta) equals k pi. What is k? (enter a decimal)", t.a, t.b), complexBridgeDecimal(v), fmt.Sprintf("The standard formula gives 2 pi divided by the square root of %d^2 - %d^2 = %d, so k = 2/%d = %s.", t.a, t.b, t.h, t.h, complexBridgeDecimal(v)))
	case 2:
		a := 1 + rand.Intn(10)
		v := 1 / float64(a)
		return complexBridgeProblem(o, fmt.Sprintf("The integral from minus infinity to infinity of dx/(x^2 + %d^2) equals k pi. What is k? (enter a decimal)", a), complexBridgeDecimal(v), fmt.Sprintf("The integral is pi/%d, so k = 1/%d = %s.", a, a, complexBridgeDecimal(v)))
	default:
		a := 1 + rand.Intn(5)
		v := 1 / (2 * float64(a*a*a))
		return complexBridgeProblem(o, fmt.Sprintf("The integral from minus infinity to infinity of dx/(x^2 + %d^2)^2 equals k pi. What is k? (enter a decimal)", a), complexBridgeDecimal(v), fmt.Sprintf("The standard formula is pi/(2*%d^3), so k = %s.", a, complexBridgeDecimal(v)))
	}
}

func argumentPrincipleProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(10)
		return complexBridgeProblem(o, fmt.Sprintf("How many zeros does z^%d have inside the unit disc, counted with multiplicity? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("The only zero is 0, with multiplicity %d.", n))
	case 2:
		n := 1 + rand.Intn(5)
		m := 1 + rand.Intn(5)
		v := n - m
		return complexBridgeProblem(o, fmt.Sprintf("For f(z) = z^%d/(z-1)^%d inside |z| < 2, what is the number of zeros minus the number of poles? (enter a number)", n, m), strconv.Itoa(v), fmt.Sprintf("There are %d zeros at 0 and %d poles at 1, so the difference is %d.", n, m, v))
	default:
		n := 1 + rand.Intn(8)
		r := 1 + rand.Intn(5)
		return complexBridgeProblem(o, fmt.Sprintf("How many zeros does z^%d e^z have inside |z| < %d? (enter a number)", n, r), strconv.Itoa(n), fmt.Sprintf("The exponential has no zeros, so the count is the multiplicity %d at 0.", n))
	}
}

func windingProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(9)
		return complexBridgeProblem(o, fmt.Sprintf("What is the winding number of z^%d around 0 on the unit circle? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("The image winds %d times around the origin.", n))
	case 2:
		a := rand.Intn(9) - 4
		v := 0
		if a > -2 && a < 2 {
			v = 1
		}
		return complexBridgeProblem(o, fmt.Sprintf("What is the winding number of the circle |z| = 2 around the point %d? (enter a number)", a), strconv.Itoa(v), "The point lies inside the circle exactly when its distance from 0 is less than 2.")
	default:
		n := 1 + rand.Intn(3)
		r := 2 + rand.Intn(2)
		lim := complexBridgePow(r, n)
		a := 1 + rand.Intn(lim+2)
		v := 0
		if a < lim {
			v = n
		}
		return complexBridgeProblem(o, fmt.Sprintf("What is the winding number of z^%d on |z| = %d around the point %d? (enter a number)", n, r, a), strconv.Itoa(v), fmt.Sprintf("The image circle has radius %d^%d = %d; the winding is %d when the point lies inside it.", r, n, lim, n))
	}
}

func riemannMappingProblem(o int) generator.Problem {
	switch o {
	case 1:
		return complexBridgeCount(o, "How many of these domains are simply connected?", riemannDomainsPool)
	case 2:
		return complexBridgeCount(o, "How many of these statements about the Riemann mapping theorem are true?", riemannStatementsPool)
	default:
		theta := 15 * (1 + rand.Intn(6))
		return complexBridgeProblem(o, fmt.Sprintf("A conformal map sends two curves that meet at %d degrees to two image curves. At what angle do the image curves meet? (enter a number of degrees)", theta), strconv.Itoa(theta), "Conformal maps preserve angles, including orientation.")
	}
}

func identityTheoremProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(8)
		a := rand.Intn(7) - 3
		return complexBridgeProblem(o, fmt.Sprintf("What is the order of the zero of (z-%d)^%d at z = %d? (enter a number)", a, n, a), strconv.Itoa(n), fmt.Sprintf("The factor appears %d times, so the zero has order %d.", n, n))
	case 2:
		return complexBridgeCount(o, "How many of these zero patterns force an analytic function on the unit disc to be identically zero?", identityFunctionsPool)
	default:
		return complexBridgeCount(o, "How many of these statements about the identity theorem are true?", identityStatementsPool)
	}
}

func moreraProblem(o int) generator.Problem {
	switch o {
	case 1:
		return complexBridgeCount(o, "How many of these functions satisfy the hypothesis of Morera's theorem on the stated domain?", moreraFunctionsPool)
	case 2:
		return complexBridgeCount(o, "How many of these statements about Morera's theorem are true?", moreraStatementsPool)
	default:
		return complexBridgeCount(o, "How many of these functions are analytic on the stated domain?", moreraFunctionsPool)
	}
}
