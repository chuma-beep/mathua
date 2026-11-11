# Slater's condition

Slater's condition requires a point that satisfies every inequality strictly and every equality exactly. For a convex problem it guarantees strong duality and the existence of multipliers. The condition is about the geometry of the feasible set rather than about the objective.

## Testing strict feasibility

Substitute a candidate point into every inequality and require a strict margin, then check the equalities exactly. A single point with all strict margins is a certificate. A constraint that is affine can be satisfied with equality in some formulations, which relaxes the requirement for that constraint.

## Deciding whether it holds

If the feasible set has an interior relative to the inequalities, Slater holds. If every feasible point makes some inequality tight, it fails. A common failure is a feasible set that is a single point or a lower-dimensional face, where no strict margin exists.

## Finding an interior point

For a polytope, the average of its vertices is interior when the polytope is full-dimensional. For a ball constraint, the centre works. For general convex sets, an interior point can be found by maximizing the minimum slack, which is itself a small optimization problem.

Slater's condition is the usual hypothesis that makes strong duality usable. It is geometric, checkable and often satisfied by the obvious interior point of a model.
