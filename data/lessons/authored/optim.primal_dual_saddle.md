# The primal-dual saddle formulation

A constrained convex problem becomes a saddle problem through its Lagrangian. The primal variable minimises and the multiplier maximises, and the saddle point is exactly the primal-dual optimal pair when strong duality holds.

## The Lagrangian as a saddle function

The Lagrangian is convex in the primal variable for fixed multipliers and affine in the multipliers for a fixed primal variable. That makes it convex-concave, which is the structure the saddle theory needs.

## Recovering the optimum

At a saddle point the primal coordinate solves the constrained problem and the multiplier coordinate is an optimal dual certificate. Complementary slackness appears as the condition that the multiplier vanishes where the constraint is slack.

## Why the reformulation helps

The saddle form removes the constraint from the update and replaces it with a gradient step on the multiplier. That makes first-order methods applicable to constrained problems without an inner solve, at the cost of a two-timescale dynamic.

To check your understanding, keep convexity in the primal block, read the primal coordinate and move the constraint into an update.

The saddle formulation turns a constraint into a second variable, and the price is that convergence now depends on the balance between the two blocks.
