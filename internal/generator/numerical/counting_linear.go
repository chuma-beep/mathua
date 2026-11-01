package numerical

import "github.com/chuma-beep/mathua/internal/generator"

func orthogonalProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about orthogonal polynomials and least squares.", []numaVerdict{
		{"orthogonal polynomials are orthogonal with respect to a weight", true, "the inner product carries the weight"},
		{"Legendre polynomials are orthogonal on [-1,1] with weight 1", true, "the standard family"},
		{"least squares can use orthogonal basis functions", true, "the basis simplifies the system"},
		{"normal equations become diagonal for an orthogonal basis", true, "cross terms vanish"},
		{"Gram-Schmidt constructs orthogonal functions", true, "it removes projections"},
		{"Chebyshev polynomials are orthogonal with weight 1/sqrt(1-x^2)", true, "the Chebyshev weight"},
		{"an orthogonal basis improves conditioning", true, "the system is better scaled"},
		{"Fourier series uses orthogonal sine and cosine functions", true, "the classical orthogonal system"},
		{"orthogonal polynomials are orthogonal only at the data nodes", false, "orthogonality is an integral"},
		{"Legendre polynomials are orthogonal on [0,1]", false, "the interval is [-1,1]"},
		{"normal equations stay dense for an orthogonal basis", false, "they become diagonal"},
		{"Gram-Schmidt destroys orthogonality", false, "it creates it"},
		{"the Chebyshev weight is 1", false, "it is 1/sqrt(1-x^2)"},
		{"orthogonal bases worsen conditioning", false, "they improve it"},
		{"Fourier series uses nonorthogonal functions", false, "sines and cosines are orthogonal"},
		{"orthogonal polynomials cannot be generated recursively", false, "three-term recurrences exist"},
	})
}

func extrapolationProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about extrapolation.", []numaVerdict{
		{"Richardson extrapolation combines approximations to cancel the leading error", true, "the error terms are matched"},
		{"Romberg integration applies Richardson extrapolation to the trapezoid rule", true, "that is its construction"},
		{"extrapolation assumes an asymptotic error expansion", true, "the cancellation needs the form"},
		{"extrapolation can improve the order of a method", true, "the leading term is removed"},
		{"two O(h) estimates can combine into an O(h^2) estimate", true, "the linear terms cancel"},
		{"the extrapolation tableau is triangular", true, "each column improves the previous"},
		{"extrapolation can amplify roundoff at high order", true, "differences of close values"},
		{"step sizes are often halved between estimates", true, "the usual Richardson choice"},
		{"extrapolation cancels all error terms at once", false, "only the modelled terms"},
		{"Romberg uses Simpson estimates only", false, "it starts from trapezoid"},
		{"extrapolation needs no error model", false, "the expansion is the assumption"},
		{"extrapolation lowers the order", false, "it raises it"},
		{"two O(h) estimates give an O(h^3) estimate", false, "they give O(h^2)"},
		{"the extrapolation tableau is square and dense", false, "it is triangular"},
		{"extrapolation reduces roundoff at every order", false, "it can amplify it"},
		{"step sizes must double between estimates", false, "halving is the usual choice"},
	})
}

func integrationErrorProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about quadrature error.", []numaVerdict{
		{"the composite trapezoid error is O(h^2) globally", true, "second order"},
		{"the composite Simpson error is O(h^4) globally", true, "fourth order"},
		{"the composite midpoint error is O(h^2) globally", true, "second order"},
		{"the error involves a derivative of the integrand", true, "the Euler-Maclaurin form"},
		{"halving h reduces trapezoid error by about four", true, "the O(h^2) scaling"},
		{"Simpson requires an even number of subintervals", true, "pairs of intervals"},
		{"a bound on the derivative gives an error bound", true, "the standard estimate"},
		{"composite rules sum the local errors", true, "each panel contributes"},
		{"the composite trapezoid error is O(h) globally", false, "it is O(h^2)"},
		{"the composite Simpson error is O(h^2) globally", false, "it is O(h^4)"},
		{"the composite midpoint error is O(h^4) globally", false, "it is O(h^2)"},
		{"the error is independent of the integrand", false, "derivatives set the size"},
		{"halving h doubles trapezoid error", false, "it quarters it"},
		{"Simpson requires an odd number of subintervals", false, "it requires an even number"},
		{"composite rules do not accumulate local error", false, "the local errors add"},
		{"the error bound does not involve the interval length", false, "the length scales it"},
	})
}

func rombergProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about Romberg integration.", []numaVerdict{
		{"Romberg combines trapezoid estimates", true, "the first column is trapezoid"},
		{"Romberg uses Richardson extrapolation", true, "each column cancels an error term"},
		{"Romberg builds a triangular tableau", true, "one row per refinement"},
		{"Romberg can reach high order with few evaluations", true, "the reuse of points helps"},
		{"the first column is the composite trapezoid rule", true, "the base estimates"},
		{"halving the step supplies the next row", true, "the standard refinement"},
		{"Romberg assumes smoothness", true, "the error expansion needs it"},
		{"roundoff can limit the achievable order", true, "differences amplify noise"},
		{"Romberg uses Simpson estimates only", false, "it starts from trapezoid"},
		{"Romberg needs no extrapolation", false, "extrapolation is the method"},
		{"Romberg's tableau is a single row", false, "it is triangular"},
		{"Romberg cannot improve on the trapezoid rule", false, "it raises the order"},
		{"the first column is Gauss quadrature", false, "it is trapezoid"},
		{"doubling the step supplies the next row", false, "the step is halved"},
		{"Romberg works for discontinuous integrands without penalty", false, "smoothness is assumed"},
		{"roundoff improves at higher order", false, "it can worsen"},
	})
}

func adaptiveProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about adaptive quadrature.", []numaVerdict{
		{"adaptive quadrature refines where the error is large", true, "the estimate drives the split"},
		{"it compares two estimates on a subinterval", true, "the difference estimates the error"},
		{"it can save work on smooth regions", true, "those panels stay coarse"},
		{"it uses a tolerance to decide subdivision", true, "the local test is tolerance-based"},
		{"recursion splits intervals", true, "the usual implementation"},
		{"adaptive methods handle localised features", true, "they concentrate effort"},
		{"the total error is controlled by local tests", true, "the local bounds sum"},
		{"it can be more efficient than a uniform rule", true, "fewer wasted evaluations"},
		{"adaptive quadrature refines uniformly", false, "that is the non-adaptive rule"},
		{"it ignores local error estimates", false, "they drive the refinement"},
		{"it always uses more work than a uniform rule", false, "it often uses less"},
		{"it needs no tolerance", false, "the tolerance is the stopping rule"},
		{"it cannot subdivide intervals", false, "subdivision is the method"},
		{"it fails on localised features", false, "that is its strength"},
		{"local tests guarantee the global error exactly", false, "they estimate it"},
		{"it is restricted to polynomials", false, "it applies to general integrands"},
	})
}

func improperProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about improper integrals.", []numaVerdict{
		{"an improper integral can have an infinite limit", true, "one or both endpoints are infinite"},
		{"it can have an integrand singularity", true, "the integrand can blow up"},
		{"convergence is defined by a limit", true, "the limit of proper integrals"},
		{"substitution can map an infinite interval to a finite one", true, "a standard transformation"},
		{"the p-test gives convergence for p>1 on [1,inf)", true, "the tail is integrable"},
		{"the p-test gives convergence for p<1 on (0,1]", true, "the singularity is integrable"},
		{"splitting at a singularity is often needed", true, "each side is handled separately"},
		{"numerical quadrature needs a transformed finite integral", true, "standard rules assume finite panels"},
		{"improper integrals always converge", false, "many diverge"},
		{"an infinite limit makes the integral diverge automatically", false, "decay can make it converge"},
		{"convergence is decided by the integrand at one point", false, "the whole tail matters"},
		{"substitution cannot remove an infinite interval", false, "it often maps it to a finite one"},
		{"the p-test gives convergence for p<1 on [1,inf)", false, "that diverges"},
		{"the p-test gives convergence for p>1 on (0,1]", false, "that diverges"},
		{"a singularity never needs splitting", false, "splitting is common"},
		{"standard quadrature handles infinite limits directly without transformation", false, "a transform or truncation is needed"},
	})
}

func pivotingProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about pivoting in Gaussian elimination.", []numaVerdict{
		{"partial pivoting chooses the largest pivot in the column", true, "the row with the largest magnitude"},
		{"pivoting reduces the growth of roundoff", true, "it avoids small divisors"},
		{"a zero pivot requires a row exchange", true, "elimination cannot proceed"},
		{"complete pivoting searches rows and columns", true, "the largest entry in the remaining block"},
		{"pivoting changes the order of equations, not the solution", true, "the system is equivalent"},
		{"scaled pivoting compares rows relative to their size", true, "it normalises by row scale"},
		{"pivoting improves stability", true, "the multiplier sizes stay bounded"},
		{"a tiny pivot can amplify error", true, "division by it magnifies rounding"},
		{"partial pivoting chooses the smallest pivot", false, "it chooses the largest"},
		{"pivoting increases roundoff", false, "it controls it"},
		{"a zero pivot is ignored", false, "it must be exchanged"},
		{"complete pivoting searches rows only", false, "rows and columns"},
		{"pivoting changes the solution", false, "the solution is unchanged"},
		{"scaled pivoting ignores row size", false, "row size is the scale"},
		{"pivoting worsens stability", false, "it improves it"},
		{"a tiny pivot is harmless", false, "it amplifies error"},
	})
}

func gaussSeidelProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about the Gauss-Seidel method.", []numaVerdict{
		{"Gauss-Seidel updates each component using the newest values", true, "updates are used immediately"},
		{"it is an iterative method for linear systems", true, "it refines an initial guess"},
		{"it converges for diagonally dominant matrices", true, "a standard sufficient condition"},
		{"it often converges faster than Jacobi", true, "the newest values help"},
		{"the iteration matrix uses the lower triangle", true, "the sweep is forward"},
		{"a good initial guess can reduce iterations", true, "the start matters"},
		{"it can diverge for some matrices", true, "diagonal dominance is sufficient, not necessary"},
		{"it is a stationary iteration", true, "the iteration matrix is fixed"},
		{"Gauss-Seidel updates using only old values", false, "that is Jacobi"},
		{"it is a direct method", false, "it is iterative"},
		{"it always converges", false, "it can diverge"},
		{"it is always slower than Jacobi", false, "it is often faster"},
		{"the iteration matrix uses the upper triangle only", false, "the lower triangle"},
		{"the initial guess does not affect iterations", false, "it affects the path"},
		{"it cannot diverge", false, "it can"},
		{"it is a nonstationary Krylov method", false, "it is stationary"},
	})
}

func leastSquaresQRProblem() generator.Problem {
	return countNumaVerdicts("Consider these statements about QR least squares.", []numaVerdict{
		{"QR solves least squares without forming normal equations", true, "it factors the design matrix"},
		{"Q has orthonormal columns", true, "the defining property"},
		{"R is upper triangular", true, "the triangular factor"},
		{"the QR least-squares solution solves R x = Q^T b", true, "the reduced system"},
		{"QR is more stable than normal equations", true, "it avoids squaring the condition number"},
		{"Householder reflections produce QR", true, "a standard construction"},
		{"Givens rotations produce QR", true, "another standard construction"},
		{"normal equations square the condition number", true, "that is their weakness"},
		{"QR forms normal equations explicitly", false, "it avoids them"},
		{"Q has orthogonal rows only", false, "its columns are orthonormal"},
		{"R is lower triangular", false, "it is upper triangular"},
		{"the solution solves Q x = R b", false, "it solves R x = Q^T b"},
		{"QR is less stable than normal equations", false, "it is more stable"},
		{"Householder reflections produce LU", false, "they produce QR"},
		{"Givens rotations produce Cholesky", false, "they produce QR"},
		{"normal equations improve the condition number", false, "they square it"},
	})
}
