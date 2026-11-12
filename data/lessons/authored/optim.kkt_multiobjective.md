# Multiobjective KKT conditions

At a smooth Pareto point a nonnegative combination of the objective gradients vanishes. The multipliers are the weights that make the point optimal for a weighted sum.

## Stationarity

The condition is a nonnegative combination of gradients equal to zero, with not all multipliers zero. It generalises the single-objective stationarity condition.

## The multiplier cone

The admissible multipliers form a cone. A point whose gradients admit such a combination is supported by some weighting.

## Weighted stationarity

Normalising the multipliers to sum to one gives the weights of a scalarised problem. The normalised weights are the local trade-off rates.

To check your understanding, form a gradient combination, test cone membership and normalise the weights.

The multiobjective KKT conditions are the bridge from geometry to weights.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
