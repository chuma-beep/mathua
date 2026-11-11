# Step size and convergence

The step size controls how far gradient descent moves. Too small and progress is slow, too large and the iterates oscillate or diverge. For a quadratic there is an exact interval of stable sizes, and a natural choice inside it, which makes the trade-off quantitative rather than a matter of taste.

## Testing a step size

On a quadratic with curvature, the iteration is stable exactly when the step size is below twice the reciprocal of the largest curvature. The fastest convergence for a simple quadratic uses the reciprocal of the curvature, and for a general quadratic the reciprocal of the largest eigenvalue. Testing a proposed size means comparing it with that bound.

## Counting iterations

For a quadratic, the error contracts by a fixed factor per step, so the number of steps needed to reach a tolerance follows from a logarithm. If the contraction factor is r, the iteration count is the logarithm of the tolerance ratio divided by the logarithm of r. Counting this way shows why a badly conditioned problem is slow even with a perfectly chosen size.

## Choosing a stable step

A stable step for an unknown curvature can be obtained by a line search, which reduces the size until decrease is observed. A fixed step is safe only when a bound on the curvature is known. The two strategies differ in what they assume, and a method that reports divergence usually had neither a bound nor a search.

The step size is where stability lives. Understanding the quadratic case gives the reference against which backtracking and exact line searches can be read.
