# The extragradient method

A plain gradient step on a saddle problem can diverge because the two blocks push the iterate around the equilibrium. The extragradient method takes a trial step first, evaluates the gradient there, and steps from the original point with the trial gradient.

## The two evaluations

The method computes a gradient at the current iterate and moves to a trial point. It then recomputes the gradient at the trial point and applies that direction to the original iterate. The trial point is never kept.

## The stability argument

For a monotone operator the extra evaluation removes the rotational overshoot that makes a single step diverge. Convergence holds for any step below the reciprocal of the Lipschitz constant, without a smaller step for the saddle case.

## The cost

Every iteration pays two gradient evaluations and one extra point in memory. In a distributed or expensive setting that doubling matters, and optimistic variants recover much of the stability with one evaluation.

To check your understanding, compute the current gradient, assume a monotone operator and pay two evaluations.

Extragradient buys stability for one extra gradient per step, and it is the reference method against which cheaper saddle solvers are measured.
