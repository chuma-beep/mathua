# The Wolfe conditions

The Wolfe conditions pair sufficient decrease with a curvature condition. The first ensures progress, the second ensures the step is long enough that the slope has flattened appreciably. Together they guarantee that a quasi-Newton update based on the step keeps its curvature approximation positive definite.

## Testing the curvature condition

The curvature condition compares the slope at the trial point with the slope at the current point, scaled by a constant near one. It requires the new slope to be at least that fraction of the old slope, which rules out steps that stop in a steep part of the descent. The test needs the gradient at the trial point.

## Testing both conditions

A step is Wolfe-acceptable when it satisfies the decrease inequality and the curvature inequality at the same time. The intervals of acceptable sizes can be empty for some constants, which is why the constants are chosen with a gap between them. Checking both is two inequalities and one extra gradient.

## Choosing a valid step

A bracketing search finds an interval containing acceptable sizes, then shrinks it while maintaining the bracket. The result satisfies both conditions by construction. For a quadratic, the exact minimizing step satisfies them for any admissible constants, which is a useful check on an implementation.

The curvature condition is what makes a line search compatible with BFGS. It costs one extra gradient and buys a stable curvature update.
