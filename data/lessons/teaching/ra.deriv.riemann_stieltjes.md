# The Riemann-Stieltjes Integral

## Riemann-Stieltjes sums

The Riemann-Stieltjes integral replaces the length increment of a subinterval by the increment of an integrator function g. A sum uses values of f at sample points multiplied by the corresponding changes in g. When g is the identity, the construction reduces to the ordinary Riemann integral.

The integrator need not be differentiable. It may be a step function, in which case the integral picks out weighted jumps, or a monotone function, in which case the increments are nonnegative. This flexibility lets one express sums, weighted averages and distribution functions in a single framework.

## Integration by parts

Integration by parts has a symmetric form for Riemann-Stieltjes integrals: the integral of f with respect to g plus the integral of g with respect to f equals the boundary term. The formula holds when the integrals exist and the boundary values are interpreted correctly, and it is often the easiest way to evaluate one of the two integrals.

The formula is especially useful when one integrator is smooth and the other is a step function. It transfers the jump from one factor to the other and can turn a sum into an ordinary integral. The boundary term must not be forgotten, since it carries the contribution of the endpoints.

## When the integral exists

If f is continuous and g has bounded variation, the Riemann-Stieltjes integral exists. Continuity of f controls the sampling error, while bounded variation of g controls the total increment. Either hypothesis alone is not enough in every formulation, so the pair is the standard sufficient condition.

If g is continuously differentiable, the Stieltjes integral equals the ordinary integral of f times the derivative of g. If g is a step function, the integral becomes a finite sum of values of f times the jumps. These two reductions cover most computations and explain why the Stieltjes notation is a unifying device rather than a new difficulty.
