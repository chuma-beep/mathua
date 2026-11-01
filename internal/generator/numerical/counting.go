package numerical

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The numerical-analysis corpus used to label several concepts multiple
// choice only because their generators asked for a definition. These pools
// keep the idea but ask for a count, which is production work: the learner
// has to judge each statement instead of recognising one option.
type numaVerdict struct {
	label string
	yes   bool
	note  string
}

const numaDraw = 8

// countNumaVerdicts draws the answer uniformly from 0 to 8, then takes that
// many true statements and the rest false ones. Drawing the count first keeps
// every answer reachable; sampling eight from a pool with extra entries would
// make 0 and 8 vanishingly rare.
func countNumaVerdicts(lead string, pool []numaVerdict) generator.Problem {
	var yes, no []numaVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })

	var count int
	var chosen []numaVerdict
	for attempt := 0; attempt < 40; attempt++ {
		count = rand.Intn(numaDraw + 1)
		chosen = append([]numaVerdict{}, yes[:count]...)
		chosen = append(chosen, no[:numaDraw-count]...)
		rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })
		if !strings.Contains(lead, strconv.Itoa(count)) {
			break
		}
	}

	parts := make([]string, 0, len(chosen))
	notes := make([]string, 0, len(chosen))
	for i, v := range chosen {
		parts = append(parts, fmt.Sprintf("(%c) %s", rune('a'+i), v.label))
		verdict := "false"
		if v.yes {
			verdict = "true"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note))
	}

	return generator.Problem{
		Question:    fmt.Sprintf("%s (a) ... (h) below. How many are true? (enter a number)\n\n%s", lead, strings.Join(parts, "\n")),
		Answer:      strconv.Itoa(count),
		Explanation: fmt.Sprintf("Exactly %d of the eight statements are true.\n%s", count, strings.Join(notes, "\n")),
		Facts:       map[string]string{"count": strconv.Itoa(count)},
	}
}

func roundoffProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about floating point roundoff.", []numaVerdict{
		{"machine epsilon is the gap between 1 and the next representable number", true, "it bounds relative rounding at 1"},
		{"subtracting nearly equal numbers can cancel significant digits", true, "this is catastrophic cancellation"},
		{"relative error is |x - xhat|/|x| for x nonzero", true, "the usual normalisation"},
		{"roundoff error comes from finite precision arithmetic", true, "reals are stored with finitely many bits"},
		{"adding a small number to a large one can lose the small one", true, "its exponent is too far below"},
		{"IEEE double precision carries about 16 decimal digits", true, "53 bits is about 15.95 digits"},
		{"cancellation is worst when operands are close in value", true, "the leading digits agree and vanish"},
		{"error can grow when many operations are chained", true, "each step can add a rounding"},
		{"roundoff is eliminated by using floating point exactly", false, "floating point is not exact for most reals"},
		{"machine epsilon is the largest representable number", false, "it is a gap near 1"},
		{"relative error is always larger than absolute error", false, "the comparison depends on magnitude"},
		{"subtracting two close numbers is numerically safe", false, "it can cancel digits"},
		{"IEEE double precision carries about 32 decimal digits", false, "it carries about 16"},
		{"roundoff depends only on hardware speed", false, "it depends on precision and operations"},
		{"adding many small numbers in any order gives identical results", false, "order changes rounding"},
		{"absolute error bounds relative error for every number", false, "a tiny x breaks that claim"},
	})
}

func stabilityProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about numerical stability.", []numaVerdict{
		{"a stable algorithm keeps the effect of input perturbations bounded", true, "error does not blow up"},
		{"a backward stable algorithm solves a nearby problem exactly", true, "the computed answer is exact for perturbed data"},
		{"conditioning is a property of the problem, not the algorithm", true, "it measures sensitivity"},
		{"an ill-conditioned problem can defeat any algorithm", true, "the data alone loses information"},
		{"stability concerns error growth during computation", true, "rounding must not be amplified"},
		{"forward error measures the distance from the exact answer", true, "output error"},
		{"a stable algorithm can still be inaccurate on an ill-conditioned problem", true, "stability cannot create information"},
		{"pivoting improves stability in Gaussian elimination", true, "it avoids tiny pivots"},
		{"stability and conditioning are the same property", false, "one is the algorithm, the other the problem"},
		{"a well-conditioned problem guarantees every algorithm is accurate", false, "an unstable algorithm can still fail"},
		{"an unstable algorithm always gives a wrong sign", false, "the error can be small on some data"},
		{"backward stability implies small forward error for every problem", false, "conditioning intervenes"},
		{"conditioning depends on the algorithm used", false, "it is intrinsic to the problem"},
		{"unstable algorithms are always faster", false, "speed and stability are independent"},
		{"a stable algorithm cannot amplify rounding at all", false, "bounded amplification is allowed"},
		{"forward error and backward error are always equal", false, "they measure different things"},
	})
}

func taylorErrorProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about Taylor truncation error.", []numaVerdict{
		{"Taylor's theorem gives a remainder term for the truncation error", true, "the remainder bounds the tail"},
		{"truncating after the linear term leaves an O(h^2) error", true, "the quadratic term is first omitted"},
		{"the Lagrange remainder involves a derivative at an unknown point", true, "the mean-value form"},
		{"a smaller step size usually reduces truncation error", true, "the omitted terms shrink"},
		{"truncation error comes from replacing a function by a polynomial", true, "the series is cut short"},
		{"an n-th order Taylor method has local error O(h^(n+1))", true, "one more power than the order"},
		{"a larger order usually gives a smaller truncation error", true, "more terms are retained"},
		{"the remainder can be bounded by the maximum of the next derivative", true, "a standard estimate"},
		{"Taylor truncation error is caused by floating point rounding", false, "that is roundoff, not truncation"},
		{"the Taylor remainder is nonzero only at the expansion point", false, "it vanishes there and grows away"},
		{"increasing the step size always reduces truncation error", false, "larger h means a larger tail"},
		{"Taylor's theorem applies only to polynomials", false, "it applies to smooth functions"},
		{"the remainder does not depend on the step size", false, "it scales with h"},
		{"truncation error and roundoff error are the same", false, "one is mathematical, one is finite precision"},
		{"a first-order method has O(h^3) local error", false, "it has O(h^2)"},
		{"Taylor series converge for every smooth function on every interval", false, "the radius can be finite"},
	})
}

func falsePositionProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about the false position method.", []numaVerdict{
		{"regula falsi keeps a bracket with a sign change", true, "the bracket is the method's invariant"},
		{"the method uses a secant line through the bracket endpoints", true, "that is the false position"},
		{"convergence is usually linear", true, "the basic form is not quadratic"},
		{"the basic form always keeps one endpoint after an update", true, "it replaces only one side"},
		{"it is generally slower than the secant method", true, "the bracket can stall one side"},
		{"it needs f(a) and f(b) with opposite signs", true, "a bracket is required"},
		{"it can stall when one endpoint never updates", true, "the Illinois variant addresses this"},
		{"the Illinois variant speeds up the basic method", true, "it rescales the stale side"},
		{"regula falsi can be applied without a sign change", false, "then no bracket is guaranteed"},
		{"its convergence is quadratic in general", false, "that is the secant method's behaviour"},
		{"it always discards both endpoints each step", false, "one endpoint remains"},
		{"it uses the derivative of f", false, "it uses function values only"},
		{"it is identical to Newton's method", false, "Newton uses a tangent at one point"},
		{"it cannot stall", false, "the basic form can"},
		{"it requires a starting interval of width one", false, "any width with a sign change works"},
		{"it converges only for polynomials", false, "continuity is the usual requirement"},
	})
}

func multipleRootProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about multiple roots.", []numaVerdict{
		{"a root of multiplicity greater than one also makes f' zero", true, "the derivative shares the factor"},
		{"Newton's method converges linearly at a multiple root in general", true, "the quadratic rate is lost"},
		{"modified Newton with known multiplicity m uses x - m f/f'", true, "the multiplicity rescales the step"},
		{"deflation removes a found root factor", true, "divide out (x - r)"},
		{"a double root touches the axis without crossing", true, "the sign does not change"},
		{"the multiplicity is the order of the first nonzero derivative", true, "a standard characterisation"},
		{"bracketing methods can miss an even-multiplicity root", true, "there is no sign change"},
		{"multiple roots make the problem ill-conditioned", true, "small data changes move the root"},
		{"at a double root f' is nonzero", false, "f' vanishes there"},
		{"Newton's method is quadratic at a multiple root without modification", false, "it degrades to linear"},
		{"a multiple root always changes sign", false, "even multiplicity does not"},
		{"deflation increases the multiplicity", false, "it removes a factor"},
		{"multiplicity is the value of f at the root", false, "that value is zero"},
		{"bisection always finds even-multiplicity roots", false, "it needs a sign change"},
		{"multiple roots are well-conditioned", false, "they are sensitive"},
		{"a triple root makes f'' nonzero", false, "f'' also vanishes"},
	})
}

func bracketingProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about bracketing methods.", []numaVerdict{
		{"bisection needs a sign change on the interval", true, "that guarantees a root"},
		{"the intermediate value theorem justifies a bracket", true, "continuity plus opposite signs"},
		{"bisection halves the interval each step", true, "the midpoint replaces one side"},
		{"bisection converges linearly with rate 1/2", true, "the error halves every step"},
		{"an even-multiplicity root may not produce a sign change", true, "the function can stay on one side"},
		{"the bracket always contains a root if f is continuous", true, "by the intermediate value theorem"},
		{"bisection is robust but slow", true, "it ignores derivative information"},
		{"a bracket of width w needs about log2(w/tol) steps", true, "each step halves the width"},
		{"bisection requires the derivative", false, "it uses only function values"},
		{"bisection converges quadratically", false, "it is linear"},
		{"a sign change guarantees a simple root", false, "the root can be a crossing of odd multiplicity"},
		{"bisection doubles the interval each step", false, "it halves it"},
		{"bracketing works for discontinuous functions without assumptions", false, "continuity is required"},
		{"bisection always finds all roots at once", false, "it finds one root in a bracket"},
		{"a bracket must have width one", false, "any width is allowed"},
		{"bisection can diverge for continuous functions", false, "the bracket shrinks to a root"},
	})
}

func interpErrorProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about interpolation error.", []numaVerdict{
		{"the error has a factor equal to the product of (x - x_i)", true, "the nodal polynomial"},
		{"more nodes do not always reduce the maximum error", true, "Runge's phenomenon is the counterexample"},
		{"Runge's phenomenon shows high-degree interpolation can oscillate", true, "the edges oscillate most"},
		{"a polynomial of degree n is determined by n+1 points", true, "counting coefficients"},
		{"the error formula involves a derivative of order n+1", true, "one order above the interpolant"},
		{"equally spaced nodes can give large edge oscillations", true, "the Lebesgue constant grows"},
		{"the interpolant matches the data exactly at the nodes", true, "that is the interpolation condition"},
		{"interpolation error is zero at every node", true, "the interpolant passes through them"},
		{"the interpolation error is zero everywhere for any data", false, "only at the nodes"},
		{"adding nodes always reduces the error", false, "Runge's phenomenon"},
		{"Runge's phenomenon affects low-degree interpolation only", false, "it appears at high degree"},
		{"n+1 points determine a polynomial of degree n+1", false, "they determine degree at most n"},
		{"the error formula involves no derivative", false, "the (n+1)-st derivative appears"},
		{"equally spaced nodes always minimise the error", false, "Chebyshev nodes do better"},
		{"the interpolant need not pass through the data", false, "it must interpolate"},
		{"interpolation error is independent of node placement", false, "the nodal polynomial depends on it"},
	})
}

func chebyshevProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about Chebyshev interpolation.", []numaVerdict{
		{"Chebyshev nodes minimise the maximum interpolation error on [-1,1]", true, "they control the nodal polynomial"},
		{"the Chebyshev nodes are cos((2k+1)pi/(2n))", true, "the standard formula"},
		{"Chebyshev polynomials satisfy T_n(cos t) = cos(nt)", true, "the defining identity"},
		{"the extrema of T_n on [-1,1] have magnitude one", true, "equioscillation"},
		{"Chebyshev nodes cluster near the endpoints", true, "the cosine spacing"},
		{"T_n has n simple roots in (-1,1)", true, "the Chebyshev nodes"},
		{"Chebyshev interpolation reduces Runge oscillations", true, "the edge error is controlled"},
		{"T_0 = 1 and T_1 = x", true, "the first two polynomials"},
		{"Chebyshev nodes are equally spaced", false, "they cluster at the ends"},
		{"Chebyshev nodes maximise the interpolation error", false, "they minimise it"},
		{"T_n(cos t) = sin(nt)", false, "the identity uses cosine"},
		{"the extrema of T_n grow like 2^n on [-1,1]", false, "they stay at magnitude one"},
		{"Chebyshev nodes cluster at the midpoint", false, "they cluster at the endpoints"},
		{"T_n has n double roots", false, "the roots are simple"},
		{"Chebyshev interpolation increases Runge oscillations", false, "it reduces them"},
		{"T_1 = 2x", false, "T_1 = x"},
	})
}

func splinesProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about cubic splines.", []numaVerdict{
		{"a cubic spline is piecewise cubic and twice continuously differentiable", true, "the standard smoothness"},
		{"natural splines set the second derivative to zero at the ends", true, "the natural boundary condition"},
		{"clamped splines prescribe endpoint slopes", true, "the clamped condition"},
		{"the spline coefficients solve a tridiagonal system", true, "each knot couples neighbours"},
		{"a cubic spline interpolates all data points", true, "that is the interpolation condition"},
		{"splines avoid the Runge oscillation of high-degree polynomials", true, "low degree on each piece"},
		{"a cubic spline is not a single cubic on the whole interval", true, "it is piecewise"},
		{"knots are the data points where pieces join", true, "the usual choice"},
		{"a cubic spline is a single cubic polynomial everywhere", false, "it changes coefficients at knots"},
		{"a natural spline sets the first derivative to zero at the ends", false, "it sets the second derivative"},
		{"clamped splines prescribe endpoint curvature", false, "they prescribe slopes"},
		{"the spline system is dense and unstructured", false, "it is tridiagonal"},
		{"a cubic spline need not match the data", false, "it interpolates"},
		{"splines always oscillate more than high-degree interpolants", false, "they reduce oscillation"},
		{"a cubic spline is only once differentiable", false, "it is twice differentiable"},
		{"knots are the endpoints only", false, "they are the interior data points"},
	})
}
