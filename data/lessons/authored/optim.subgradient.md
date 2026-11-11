# Subgradients

A subgradient generalizes the derivative to a convex function that has a corner. At a differentiable point the subgradient is the gradient, and at a corner it is any slope of a line that stays below the graph. The set of all such slopes is the subdifferential, and it is what a nonsmooth method actually uses.

## Testing a subgradient

A vector is a subgradient at a point when the linear function through that point with that slope never exceeds the function anywhere. For the absolute value at zero, every slope between minus one and one passes the test, so the subdifferential is an interval. Testing a claimed slope means checking the global inequality, not just the local left and right derivatives.

## The subdifferential of an absolute value

For the absolute value, the subdifferential is a single number away from zero, equal to the sign, and the interval from minus one to one at zero. For a weighted absolute value the interval scales by the weight. For a sum of absolute values, the subdifferential is the sum of the intervals, which is again an interval whose endpoints are the sums of the endpoints.

## Supporting lines

A subgradient defines a supporting line that touches the graph at the point and lies below it everywhere. The line is the first-order model a nonsmooth method trusts, and a step is taken along its slope. At a corner the choice of slope inside the subdifferential determines the direction, which is why the method can stall or zigzag when the choice is arbitrary.

Nonsmooth objectives are common once absolute values, maximums or hinge losses appear. The subgradient is the smallest object that still supports a descent argument, and the subdifferential describes the ambiguity at a corner.
