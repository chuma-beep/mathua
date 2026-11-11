# First-order optimality

At an interior local minimum of a differentiable function, the gradient is zero. The condition is necessary, not sufficient, so it produces candidates rather than answers. In several variables it is a system of equations, one per partial derivative, and solving that system is the first-order method in its purest form.

## Testing a stationary point

Evaluate every partial derivative at the point and check that all are zero, within the tolerance the problem states. A nonzero component means the point is not stationary, and the negative of that component is a descent direction. The test is local and cheap, and it says nothing about whether the point is a minimum, a maximum or a saddle.

## Solving the gradient condition

For a quadratic objective the gradient is a linear system, so the stationary point is found by solving it. For a separable objective the equations decouple into one equation per variable. For a general smooth objective the system is nonlinear and is solved iteratively. Writing the system explicitly is what allows a check by substitution afterward.

## Using a descent direction

A direction is a descent direction at a point when its inner product with the gradient is negative. The negative gradient is always one, and so is any direction within ninety degrees of it. A step along a descent direction lowers the objective for a small enough step, which is the fact every line search relies on. A direction with a zero inner product is orthogonal to the gradient and needs a second-order argument.

The first-order condition is the entry point to every smooth method. It identifies candidates, supplies the descent direction, and is the condition that algorithms monitor when they stop.
