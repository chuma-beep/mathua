# Lagrangian relaxation for integer programs

Lagrangian relaxation moves hard constraints into the objective with multipliers. The relaxed integer problem is often easy, and its value bounds the original problem for any nonnegative multipliers.

## The relaxation

Each relaxed constraint contributes its residual times a multiplier to the objective. For a minimisation with a less-than constraint the multiplier is nonnegative, and the relaxed value is a lower bound.

## Choosing what to relax

The constraints that couple the variables are the ones worth relaxing, because the remaining problem usually decomposes. Relaxing too much weakens the bound and relaxing too little leaves the problem hard.

## Recovering feasibility

The relaxed solution usually violates some relaxed constraint. A repair heuristic turns it into a feasible integer point, and the gap between that point and the bound measures the relaxation quality.

To check your understanding, price each relaxed constraint, target coupling constraints and repair the relaxed point.

Lagrangian relaxation buys an easy subproblem and a bound, and the multiplier search is what converts the bound into a usable one.
