# Orthogonal projection

The orthogonal projection returns the closest point of a set to a given point. For a convex set it is unique and nonexpansive.

## The definition

The projection minimises the distance to the input over the set. The residual is orthogonal to every feasible direction at the projection.

## A convex set

For a convex set the projection is unique. A nonconvex set can have several nearest points and the projection may jump as the input moves.

## Nonexpansiveness

Projections onto convex sets never increase distances between points. That property is what makes projected methods stable.

To check your understanding, project a point, project onto a convex set and compare two projections.

Projection is how a constrained step becomes feasible without a constrained solve.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
