# The duality gap

The duality gap is the difference between a primal objective value and a dual objective value at a feasible pair. It is always non-negative and it bounds the suboptimality of both. A small gap is the standard stopping certificate in convex optimization.

## Computing a gap

Evaluate the primal objective at a feasible point and the dual function at a dual-feasible multiplier vector, then subtract. Absolute values are unnecessary because the sign is guaranteed. The result is a number in objective units, which makes it directly comparable with a tolerance.

## Testing a stopping rule

Stop when the gap falls below a tolerance. The gap bounds how much better the true optimum can be than the current primal point, so the tolerance has a clear meaning. A relative tolerance divides the gap by the objective magnitude and is safer across scales.

## Comparing two pairs

Two feasible pairs can be ranked by their gaps, and the smaller gap certifies a better bound. The primal values and dual values each improve in their own direction, so the comparison should report both. A pair with a smaller gap but a worse primal value is still a better certificate, not a better solution.

The gap is the number that makes a solver's stopping rule meaningful. It converts an abstract convergence claim into a bound on how far the answer can be from optimal.
