# Bounds on gradient variance

A variance bound is the number that decides the step size and the convergence rate of a stochastic method. It can come from bounded gradients or from bounded deviation at an optimum.

## Bounded second moment

If the expected squared gradient norm is bounded by a constant, a constant step converges to a neighbourhood whose radius scales with the step and the bound. A decreasing step removes the neighbourhood at the cost of a slower rate.

## Bounded variance at the optimum

A stronger assumption bounds the variance at the solution rather than everywhere. The neighbourhood then shrinks with the step in the same way, but the constant is the variance at the optimum, which is often much smaller.

## The relaxation term

Some analyses use a relaxed bound that grows with the gradient norm. It permits unbounded noise far from the solution and still yields a neighbourhood whose size depends on the noise at the optimum.

To check your understanding, bound the squared norm, bound at the solution and grow with the norm.

Choosing between these bounds is not cosmetic. The step size, the expected rate and the achievable accuracy all follow from which bound the analysis assumes.
