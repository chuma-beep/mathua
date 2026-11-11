# Projected gradient descent

Projected gradient descent takes a gradient step and then projects the result back onto the feasible set. For a convex set the projection is unique and the method inherits the descent guarantee of the unconstrained step. It is the natural method when the feasible set is simple, such as a box or a ball.

## Projecting a point

The projection is the closest feasible point in the Euclidean distance. For a box it clamps each coordinate to its interval. For a ball it scales the point back to the boundary when it lies outside. For a half-space it subtracts the violation along the normal. Each case has a closed form, which is why the method is practical.

## Taking a projected step

Compute the gradient step, then project the result. The objective need not decrease along the raw step, but it decreases along the projected step when the step size is small enough and the set is convex. The projection is what restores feasibility without a constraint solve.

## Testing a bound constraint

At a bound-constrained optimum, a coordinate at its lower bound has a non-negative partial derivative, and one at its upper bound has a non-positive derivative. A coordinate in the interior has a zero derivative. These are the KKT conditions specialized to a box, and they give a cheap optimality test.

Projected gradient is the cheapest way to add constraints when projection is easy. Its optimality test is the KKT condition seen coordinate by coordinate.
