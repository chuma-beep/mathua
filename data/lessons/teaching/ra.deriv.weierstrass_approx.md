# The Weierstrass Approximation Theorem

## Polynomial approximation

Every continuous function on a closed bounded interval can be approximated uniformly by polynomials. For every positive epsilon there is a polynomial whose distance from the function, measured in the supremum norm, is smaller than epsilon. The polynomial can be chosen with real coefficients and its degree depends on the function and on epsilon.

The theorem does not claim that the approximating polynomials are partial sums of a Taylor series. A continuous function need not be differentiable, so its Taylor series may not exist. The approximation is a global construction, not a local expansion around a point.

## Bernstein polynomials

One constructive proof uses Bernstein polynomials. The n-th Bernstein polynomial evaluates the function at the binomial probabilities associated with n independent trials. The polynomial is a weighted average of values of the function, with weights that concentrate near the evaluation point as n grows.

Bernstein polynomials converge uniformly for every continuous function on the unit interval. The proof uses the weak law of large numbers or a direct estimate of the binomial variance. The rate of convergence can be slow, but uniformity is enough for the theorem.

## Density of polynomials

The theorem says that polynomials are dense in the space of continuous functions with the supremum norm. Consequently every continuous function is a uniform limit of polynomials, and every uniformly convergent sequence of polynomials has a continuous limit. The space of polynomials is not closed, but its closure is the whole space.

Approximation extends from a closed interval to compact subsets of higher-dimensional space by a similar construction. It also extends to periodic functions through trigonometric polynomials, which is the content of the Stone-Weierstrass theorem. The Weierstrass theorem is the concrete starting point for that general density result.
