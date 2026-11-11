# Active and inactive constraints

An inequality constraint is active when it holds with equality and inactive when it holds with slack. At an optimum only active constraints can have positive multipliers, and inactive ones can be ignored locally. Identifying the active set reduces a constrained problem to an equality-constrained one.

## Identifying active constraints

Evaluate each constraint at the point and compare its left side with its bound. Equality means active, a strict inequality means inactive. The comparison must use the same tolerance as the solver, because a value a hair below the bound is effectively active numerically.

## Testing a candidate

At a candidate, check which constraints are active, then solve the equality-constrained stationarity system for those. The multipliers of the inactive constraints are zero by complementary slackness. A positive multiplier on an inactive constraint disqualifies the candidate.

## Counting active inequalities

Count the constraints that hold with equality at the point. In two variables, two independent active constraints pin a vertex, one leaves an edge, and none leaves an interior point. The count is a quick check that the geometry matches the algebra.

The active set is the local equality structure of a constrained problem. Finding it correctly is most of the work in solving an inequality-constrained model by hand.
