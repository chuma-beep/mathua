# Saddle points

A saddle point is a pair where neither player can improve by moving alone. The first block is minimal in its own direction and the second is maximal in its own, so the point is a minimum in one subspace and a maximum in another.

## The defining inequalities

For every admissible first coordinate and every admissible second coordinate, the value at the pair is no larger than the value after changing only the second, and no smaller than the value after changing only the first. The two inequalities hold together and neither implies the other.

## First-order form

When the sets are open the gradient with respect to the first block vanishes and the gradient with respect to the second vanishes too. The Hessian is then indefinite: it has both negative and positive curvature, which is why a saddle is stationary but not an extremum.

## Existence

A convex-concave function on compact convex sets has a saddle point, and the minimax value equals the maximin value. Removing compactness or convexity can leave the inequality strict, so the value exists only in the limit.

To check your understanding, compare along the first block, set both block gradients to zero and require convex-concave structure.

A saddle point is the equilibrium of a two-block problem, and the indefinite Hessian is its signature rather than a defect.
