# The Weierstrass M-test

## Majorants for function series

A series of functions converges uniformly when its partial sums converge uniformly. The M-test replaces the functions by nonnegative numbers that bound them: if the absolute value of the n-th function is at most M_n everywhere, and the numerical series with terms M_n converges, then the function series converges uniformly and absolutely.

The bound must be uniform in the variable. A pointwise bound that changes with x is not enough. The test is a comparison test transplanted to function space, and the convergence of the numerical majorant is what makes the comparison useful.

## The M-test

Apply the Cauchy criterion to the tail of the series. The tail is bounded by the tail of the numerical series, and the numerical tail can be made small independently of the variable. That independence is exactly uniform convergence. No continuity of the individual terms is required for the convergence conclusion.

If the terms are continuous and the series converges uniformly, the sum is continuous. Integrating a uniformly convergent series term by term is also valid on a bounded interval. Differentiating term by term needs a stronger hypothesis, usually uniform convergence of the differentiated series, so the M-test alone does not justify differentiation.

## Uniform convergence consequences

Uniform convergence lets limits pass through integrals and through continuous functions. It does not generally let limits pass through derivatives. A standard counterexample is a sequence of smooth functions converging uniformly to a function with a corner, which shows that the derivative of the limit need not be the limit of the derivatives.

The M-test is often used to build continuous functions with unusual properties. A series of sawtooth or trigonometric terms with rapidly decreasing amplitudes converges uniformly, and the resulting function can be continuous while failing to be differentiable anywhere. Such constructions rely on the uniformity that the M-test supplies.
