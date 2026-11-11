# The trust-region subproblem

The trust-region subproblem minimizes a quadratic model inside a ball. It is a small constrained problem that is solved at every outer iteration. Its solution either lies strictly inside the ball, where it is an unconstrained minimizer, or on the boundary, where a multiplier enforces the radius.

## Testing a boundary solution

A solution on the boundary satisfies the ball constraint with equality. The stationarity condition includes the multiplier times the point, so the gradient of the model must be parallel to the point. A solution inside the ball must have a zero multiplier and a vanishing model gradient.

## Solving a two-variable subproblem

For a diagonal model, the interior candidate is the unconstrained minimizer and the boundary case reduces to one equation in the multiplier. Solving it gives the step. For a general model the same structure holds with an eigendecomposition, and the multiplier is found by a scalar root search.

## Finding the multiplier

The multiplier is non-negative and makes the shifted matrix positive semidefinite. Its value is the smallest number that both keeps the shifted curvature non-negative and makes the step fit the radius. This characterization is what makes the subproblem solvable without enumerating boundary points.

The trust-region subproblem is where curvature and the radius meet. Its multiplier is the same Lagrange object as in the constrained theory, applied to a single quadratic ball constraint.
