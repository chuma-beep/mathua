# Line search

A line search chooses a step length along a search direction. An exact search minimizes the one-dimensional restriction of the objective.

## The exact step for a quadratic

For \(f(x)=(x-c)^2\) and direction \(d\), the exact step is \(\alpha=(c-x_0)d\). It reaches the minimum in one step.

The derivative along the line is zero at the minimum. For a quadratic, that condition is linear in \(\alpha\).

A negative direction reverses the sign of the step.

## Inexact line searches

Large problems rarely afford an exact line search. The Armijo condition requires a sufficient decrease in the objective.

The curvature condition prevents a step that is too short. Together they form the Wolfe conditions.

The two conditions bracket a useful range of step lengths without solving a one-dimensional problem exactly.

## Strong Wolfe and cost

The strong Wolfe conditions add a bound on the magnitude of the directional derivative. They consist of two inequalities.

A line search trades extra function evaluations for a better step. The best choice depends on the cost of evaluating the objective.

Cheap objectives favor tighter searches, while expensive simulations favor looser ones.
