package calculus

import "github.com/chuma-beep/mathua/internal/generator"

func integralStrategyProblem() generator.Problem {
	return countCalcVerdicts("How many of these integrals call for the stated method?", []calcVerdict{
		{"The integral of x e^(x^2) dx calls for substitution u = x^2", true, "the derivative of x^2 supplies the factor x"},
		{"The integral of x/(x^2+1) dx calls for substitution u = x^2+1", true, "the numerator is half the derivative of the denominator"},
		{"The integral of 1/(x^2+9) dx calls for an arctangent", true, "it matches the form 1/(x^2+a^2)"},
		{"The integral of sqrt(9-x^2) dx calls for x = 3 sin t", true, "the radical becomes 3 cos t"},
		{"The integral of x^2/(x+1) dx calls for polynomial division first", true, "the numerator degree exceeds the denominator degree"},
		{"The integral of 1/(x^2-1) dx calls for partial fractions", true, "the denominator factors as (x-1)(x+1)"},
		{"The integral of x ln x dx calls for integration by parts", true, "differentiating ln x simplifies it"},
		{"The integral of sin^3 x cos^2 x dx calls for a Pythagorean rewrite", true, "an odd power of sine leaves one factor for the differential"},
		{"The integral of x e^(x^2) dx calls for partial fractions", false, "there is no rational denominator to split"},
		{"The integral of x/(x^2+1) dx calls for a trig substitution", false, "substitution handles it directly"},
		{"The integral of 1/(x^2+9) dx calls for polynomial division", false, "the numerator already has lower degree"},
		{"The integral of sqrt(9-x^2) dx calls for partial fractions", false, "the integrand is not rational"},
		{"The integral of x^2/(x+1) dx calls for a trig substitution", false, "division reduces it to a polynomial plus a simple fraction"},
		{"The integral of 1/(x^2-1) dx calls for integration by parts", false, "partial fractions is the direct route"},
		{"The integral of x ln x dx calls for substitution u = x^2", false, "that leaves ln x behind"},
		{"The integral of sin^3 x cos^2 x dx calls for partial fractions", false, "the integrand is trigonometric, not rational"},
	})
}

func seriesOperationsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about series operations are true?", []calcVerdict{
		{"If sum a_n and sum b_n converge, then sum (a_n+b_n) converges", true, "the limit of a sum is the sum of the limits"},
		{"If sum a_n converges, then sum c a_n converges for a constant c", true, "scaling multiplies the sum by c"},
		{"If sum a_n diverges and c is nonzero, then sum c a_n diverges", true, "scaling cannot repair divergence"},
		{"Changing finitely many terms does not affect convergence", true, "the tail decides convergence"},
		{"Absolute convergence implies convergence", true, "the Cauchy criterion applies to the absolute series"},
		{"The sum of two divergent series can converge", true, "sum 1/n and sum -1/n give the zero series"},
		{"Termwise addition is valid for convergent series", true, "this is the linearity of convergent sums"},
		{"Multiplying a convergent series by zero gives the zero series", true, "every term becomes zero"},
		{"If both series diverge, then their sum diverges", false, "a series and its negative both diverge but sum to zero"},
		{"If sum a_n converges, then sum |a_n| converges", false, "the alternating harmonic series is a counterexample"},
		{"Rearranging a conditionally convergent series preserves its sum", false, "Riemann's theorem allows any sum"},
		{"If a_n tends to zero, then sum a_n converges", false, "the harmonic series has terms tending to zero and diverges"},
		{"If sum a_n and sum b_n converge, then sum a_n b_n converges", false, "termwise products need not converge"},
		{"The sum of a convergent and a divergent series converges", false, "it diverges"},
		{"Termwise multiplication is valid for convergent series", false, "products of convergent series are not termwise sums"},
		{"If sum a_n converges, then sum (-1)^n a_n converges absolutely", false, "the signs do not make a conditionally convergent series absolute"},
	})
}

func pSeriesProblem() generator.Problem {
	return countCalcVerdicts("How many of these p-series statements are true?", []calcVerdict{
		{"sum 1/n^2 converges", true, "p = 2 is greater than 1"},
		{"sum 1/n^(3/2) converges", true, "p = 3/2 is greater than 1"},
		{"sum 1/n diverges", true, "p = 1 sits on the boundary"},
		{"sum 1/sqrt(n) diverges", true, "p = 1/2 is at most 1"},
		{"sum 1/n^p converges exactly when p > 1", true, "this is the p-series test"},
		{"sum 1/n^(1/2) diverges", true, "the exponent is 1/2"},
		{"sum 1/n^4 converges", true, "p = 4 is greater than 1"},
		{"sum 1/n^(1.1) converges", true, "p = 1.1 is greater than 1"},
		{"sum 1/n^(0.9) converges", false, "the exponent is below 1"},
		{"sum 1/n^2 diverges", false, "p = 2 gives convergence"},
		{"A p-series converges for every p at least 1", false, "it diverges at p = 1"},
		{"sum 1/n^(2/3) converges", false, "the exponent is below 1"},
		{"sum 1/n^3 diverges", false, "p = 3 gives convergence"},
		{"A p-series converges for every p > 0", false, "it diverges for 0 < p <= 1"},
		{"sum 1/n^(1.5) diverges", false, "p = 1.5 gives convergence"},
		{"sum 1/n^(0.1) converges", false, "the exponent is below 1"},
	})
}

func dirichletTestProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about Dirichlet's test are true?", []calcVerdict{
		{"Dirichlet: bounded partial sums of a_n and b_n decreasing to 0 give convergence of sum a_n b_n", true, "this is the test"},
		{"Abel: sum a_n converges and b_n monotone bounded gives convergence of sum a_n b_n", true, "this is Abel's test"},
		{"sum sin(n)/n converges by Dirichlet", true, "partial sums of sin n are bounded and 1/n decreases to zero"},
		{"sum (-1)^n/sqrt(n) converges by Dirichlet", true, "the alternating signs have bounded partial sums"},
		{"Dirichlet needs bounded partial sums, not convergence of sum a_n", true, "boundedness is the weaker hypothesis"},
		{"sum cos(n)/n converges by Dirichlet", true, "partial sums of cos n are bounded"},
		{"The alternating series test is a special case of Dirichlet", true, "take a_n = (-1)^n"},
		{"sum sin(n)/n^2 converges absolutely", true, "the absolute series is bounded by 1/n^2"},
		{"Dirichlet needs sum a_n to converge", false, "bounded partial sums suffice"},
		{"Dirichlet requires a_n to decrease", false, "the monotone hypothesis is on b_n"},
		{"sum sin(n) converges by Dirichlet", false, "b_n = 1 does not tend to zero"},
		{"Dirichlet proves absolute convergence", false, "it proves conditional convergence in general"},
		{"Abel's test requires b_n to decrease to zero", false, "it requires monotone and bounded"},
		{"sum cos(n)/sqrt(n) converges absolutely", false, "the absolute series behaves like 1/sqrt(n)"},
		{"Dirichlet applies only to alternating series", false, "sine and cosine series also qualify"},
		{"sum sin(n)/n diverges", false, "Dirichlet gives convergence"},
	})
}

func lagrangeStatementsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about Lagrange multipliers are true?", []calcVerdict{
		{"The maximum of x+y on x^2+y^2=r^2 is r sqrt(2)", true, "the diagonal point gives r sqrt(2)"},
		{"The minimum of x^2+y^2 given x+y=s is s^2/2", true, "the symmetric point is closest to the origin"},
		{"The maximum of xy given x+y=s is s^2/4", true, "the symmetric point gives s^2/4"},
		{"Lagrange's condition is grad f = lambda grad g", true, "the gradients are parallel"},
		{"The constraint gradient must be nonzero for the multiplier to exist", true, "otherwise the condition says nothing"},
		{"Extrema occur where the gradients are parallel or the constraint is singular", true, "these are the two cases of the theorem"},
		{"The maximum of x+y on x^2+y^2=1 is sqrt(2)", true, "substitute r = 1"},
		{"The minimum of x^2+y^2 on x+y=2 is 2", true, "the point (1,1) gives 2"},
		{"The maximum of x+y on x^2+y^2=1 is 1", false, "the diagonal point gives sqrt(2)"},
		{"Lagrange's condition is grad f = 0", false, "that is unconstrained critical point theory"},
		{"The maximum of xy given x+y=4 is 8", false, "the maximum is 4"},
		{"The multiplier always equals the maximum value", false, "it measures sensitivity, not the value"},
		{"The constraint gradient may be zero without special handling", false, "the theorem needs a regular constraint"},
		{"The maximum of x+y on a circle of radius 2 is 2", false, "it is 2 sqrt(2)"},
		{"Lagrange multipliers apply only to linear constraints", false, "smooth constraints are allowed"},
		{"The minimum of x^2+y^2 on x+y=2 is 4", false, "the minimum is 2"},
	})
}

func iteratedLimitsProblem() generator.Problem {
	return countCalcVerdicts("How many of these statements about iterated limits are true?", []calcVerdict{
		{"If the double limit and one iterated limit exist, they are equal", true, "the iterated limit is a restriction of the double limit"},
		{"Iterated limits may differ when the double limit does not exist", true, "different paths can give different values"},
		{"x^2 y^2/(x^2+y^2) has iterated limits 0 and double limit 0", true, "the bound by y^2 gives the double limit"},
		{"Uniform convergence in one variable justifies swapping the limits", true, "uniformity makes the interchange valid"},
		{"xy/(x^2+y^2) has iterated limits 0 but no double limit at the origin", true, "the diagonal path gives 1/2"},
		{"A continuous function on a compact rectangle allows limits to be interchanged", true, "continuity on a compact set is enough"},
		{"Iterated limits along the axes can differ from the limit along y = x", true, "path dependence exposes the difference"},
		{"A double limit implies both iterated limits when they exist", true, "restrict the double limit to each axis"},
		{"Iterated limits always equal each other", false, "the function xy/(x^2+y^2) separates them"},
		{"Equal iterated limits imply the double limit exists", false, "path limits can still disagree"},
		{"xy/(x^2+y^2) has double limit 0 at the origin", false, "the diagonal path gives 1/2"},
		{"Limits can always be interchanged", false, "uniformity or compactness is needed"},
		{"Double limit existence requires equal iterated limits", false, "the iterated limits may fail to exist"},
		{"(x-y)/(x+y) has equal iterated limits at the origin", false, "the two axis limits are 1 and -1"},
		{"Iterated limits exist only for continuous functions", false, "many discontinuous functions have them"},
		{"A path limit is the same as a double limit", false, "one path cannot control all paths"},
	})
}

func vectorIdentitiesProblem() generator.Problem {
	return countCalcVerdicts("How many of these vector identities are true?", []calcVerdict{
		{"curl grad f = 0", true, "mixed partial derivatives commute"},
		{"div curl F = 0", true, "the divergence of a curl vanishes"},
		{"div(fF) = f div F + grad f dot F", true, "this is the product rule for divergence"},
		{"curl(fF) = f curl F + grad f cross F", true, "this is the product rule for curl"},
		{"div(F cross G) = G dot curl F - F dot curl G", true, "this is the divergence of a cross product"},
		{"The curl of a gradient of any smooth scalar field vanishes", true, "the identity is universal"},
		{"The divergence of a curl of any smooth vector field vanishes", true, "the identity is universal"},
		{"grad(div F) is not generally equal to the Laplacian of F", true, "the missing term is curl curl F"},
		{"div curl F = div F", false, "the left side is always zero"},
		{"curl grad f = grad f", false, "the left side is always zero"},
		{"div(fF) = f div F", false, "the gradient term is missing"},
		{"curl(F cross G) = F div G - G div F", false, "the correct identity has more terms"},
		{"grad div F = Laplacian F always", false, "the curl curl term need not vanish"},
		{"div(F cross G) = F dot curl G - G dot curl F", false, "the signs are reversed"},
		{"curl curl F = 0 always", false, "it vanishes only for irrotational fields"},
		{"grad(fg) = grad f grad g", false, "the product rule gives f grad g + g grad f"},
	})
}
