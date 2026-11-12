# Pareto optimality tests

A Pareto optimality test asks whether any feasible direction improves every objective at once. If none does, the point is locally optimal.

## Sufficient conditions

For convex objectives over a convex set, a point supported by a weighted sum is globally Pareto optimal. Convexity is what upgrades a local condition to a global one.

## Convex fronts

A convex front contains every weighted-sum optimum. A nonconvex front leaves a gap that no weight reaches.

## Testing a candidate

Compare the candidate with its neighbours. If a neighbour is better in every objective, the candidate is dominated and the search continues.

To check your understanding, check convexity, recognise a convex front and compare with neighbours.

Optimality in several objectives is a statement about directions, not about a single number.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
