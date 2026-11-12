# KKT conditions without smoothness

Without differentiability the KKT gradient condition becomes an inclusion. Zero must lie in the sum of the objective subdifferential and the constraint normals.

## Subdifferentials in KKT

Each nonsmooth term contributes a set rather than a vector. The stationarity condition asks whether the sets can combine to contain zero.

## The inclusion

The sum of the subdifferentials plus the multiplier-weighted constraint normals must contain zero. A single choice from each set is enough to witness it.

## Verifying a point

A candidate is checked by finding one subgradient from each set that satisfies the inclusion. Failure means no such choice exists, so the point is not optimal.

To check your understanding, list the subgradients of each term, choose one element from each set and check each active constraint.

The inclusion is harder to check than an equation because it is existential. Small instances are checked by enumerating the active pieces.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
