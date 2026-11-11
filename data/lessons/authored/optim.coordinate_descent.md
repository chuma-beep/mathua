# Coordinate descent

Coordinate descent moves along one coordinate at a time, holding the others fixed. Each subproblem is one-dimensional, often with a closed-form solution, so the method is cheap per step and easy to implement. It can stall at a point that is minimal along every coordinate and yet not a local minimum.

## Taking a coordinate step

Fix every variable except one, then minimize the resulting one-dimensional function. For a quadratic the minimizer is minus the partial derivative divided by the diagonal curvature. Update that coordinate only, then move to the next. The objective never increases, which gives the method its stability.

## Choosing the next coordinate

Cyclic order is the simplest rule and the easiest to analyse. A greedy rule picks the coordinate with the largest partial derivative, which usually reduces the objective faster. Random order avoids the worst cycles and is common in large-scale settings. The choice changes the speed, not the fixed points.

## Testing a coordinate minimum

A point is coordinate-minimal when every partial derivative is zero, provided each one-dimensional subproblem is smooth. When a subproblem is nonsmooth, the test uses the subdifferential of that coordinate instead. Coordinate minimality is weaker than stationarity in the full gradient sense only in the nonsmooth case; for smooth functions the two coincide.

Coordinate descent turns a multivariate problem into a sequence of scalar ones. Its simplicity is real, and so is its weakness when coordinates interact strongly.
