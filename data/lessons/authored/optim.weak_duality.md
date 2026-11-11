# Weak duality

Weak duality says the dual objective at any dual-feasible point is at most the primal objective at any primal-feasible point. It holds without convexity and without constraint qualifications, which makes it the most reliable inequality in optimization. Its practical use is as a stopping rule: a small gap certifies near-optimality.

## Comparing bounds

Evaluate the primal objective at a feasible point and the dual function at a feasible multiplier vector. Weak duality guarantees the second number is not larger. The difference is the duality gap at that pair, and it bounds how far the feasible point can be from optimal.

## Testing a dual-feasible pair

A multiplier vector is dual feasible when its sign restrictions match the primal constraint directions. Given such a vector, the dual value can be computed, and it must not exceed any feasible primal value. A violation means one of the two points is not feasible or the dual was transcribed incorrectly.

## Reading the sign of the gap

The gap is non-negative by construction. A negative computed gap therefore indicates an error, not a discovery. A zero gap with both points feasible certifies that both are optimal, which is the strongest conclusion available from a single pair.

Weak duality is a free certificate. It needs no convexity, no solver and no convergence theory, only feasibility on both sides and one evaluation.
