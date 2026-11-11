# Karush-Kuhn-Tucker conditions

The KKT conditions extend the Lagrange multiplier rule to inequality constraints. They combine stationarity of the Lagrangian, primal feasibility, dual feasibility and complementary slackness. For a convex problem satisfying a qualification they are sufficient as well as necessary.

## Testing a KKT point

Check four things: the gradient of the Lagrangian vanishes, every constraint holds, every inequality multiplier is non-negative, and each slack times its multiplier is zero. A failure in any one disqualifies the point. The order of the checks matters only for readability, not for the conclusion.

## Solving a small system

Assume an active set, set the inactive multipliers to zero, solve the stationarity equations together with the active constraints, then verify the signs and the inactive feasibility. If a multiplier has the wrong sign, the assumed active set is wrong and another must be tried. For small problems this enumeration is practical.

## Identifying a missing condition

A candidate that satisfies stationarity and feasibility but not the multiplier signs is not a KKT point. The missing condition is almost always dual feasibility or complementary slackness. Naming which one fails is more useful than reporting a generic failure.

KKT is the standard optimality certificate for constrained problems. It unifies equality and inequality constraints and gives the multipliers that sensitivity analysis reads.
