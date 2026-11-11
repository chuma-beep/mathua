# Strong duality

Strong duality says the best dual value equals the primal optimum, so the gap closes. For convex problems it holds under a constraint qualification such as Slater's condition. When it holds, a dual solution certifies a primal one and the multipliers become meaningful prices.

## Testing a zero gap

Compute the best primal value and the best dual value, then compare. Equality certifies strong duality for that instance. A small nonzero gap may be numerical or may indicate that the qualification fails, and the two cannot be distinguished from the numbers alone.

## Finding a primal-dual pair

A pair is optimal when both are feasible and the gap is zero. In practice a solver returns both, and checking feasibility plus the gap is the verification. For a linear program the pair also satisfies complementary slackness, which gives a second check.

## Deciding applicability

Strong duality is guaranteed for linear programs and for convex problems satisfying a qualification. It can fail for nonconvex problems, and it can fail for convex problems without a qualification. Deciding applicability means checking the class of the problem and, for the convex case, checking the qualification rather than assuming it.

Strong duality converts a bound into an equality. It is what lets a dual solution stand in for a primal one and what makes shadow prices exact rather than approximate.
