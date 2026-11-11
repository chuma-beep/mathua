# Equality-constrained optimization

An equality constraint removes directions from the feasible set, so the gradient no longer has to vanish at a constrained minimum. Instead the gradient must be a linear combination of the constraint gradients. The coefficients are the Lagrange multipliers, and solving the resulting system gives the candidates.

## Solving a two-variable problem

Write the Lagrangian, differentiate with respect to each variable and each multiplier, and set the derivatives to zero. The multiplier equations reproduce the constraint. Solving the system gives the stationary points of the Lagrangian, which are the constrained candidates. Substituting them back into the objective ranks them.

## Testing a multiplier

A multiplier is valid when the stationarity equations hold with it and the constraint is satisfied. Its sign is free for an equality constraint, which distinguishes it from an inequality multiplier. A computed multiplier that fails the equations indicates an algebra error, not a different optimum.

## Classifying a constrained point

The second-order condition involves the Hessian of the Lagrangian restricted to the tangent space of the constraints. Positive definiteness on that subspace gives a constrained minimum. Ignoring the restriction and testing the full Hessian can misclassify a point, because directions that leave the constraint are irrelevant.

Equality constraints are the cleanest constrained case. The multiplier rule and the restricted second-order test are the reference for the inequality case that follows.
