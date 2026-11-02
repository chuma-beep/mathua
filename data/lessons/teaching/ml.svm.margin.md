# Support vector machines and margins

## The maximum-margin idea

A support vector machine separates two classes with a hyperplane and chooses the hyperplane whose margin is as wide as possible. The margin is the gap between the two parallel supporting hyperplanes, one touching the closest point of each class.

For a weight vector w and bias b, the decision boundary is the set of points with w dot x + b = 0. Scaling w and b does not change that boundary, so the scale is fixed by requiring the closest training points to satisfy |w dot x + b| = 1.

## Computing the margin

The distance from a point to the boundary is |w dot x + b| divided by the norm of w. With the closest points at distance one in the scaled expression, the full margin width becomes 2 divided by the norm of w.

Maximising the margin is therefore the same as minimising the squared norm of w subject to the constraint that every training point lies on the correct side and outside the margin.

## Support vectors and slack

Only the points that lie exactly on the margin, or violate it in the soft-margin version, influence the solution. These are the support vectors, and removing any other training point leaves the fitted boundary unchanged.

The soft-margin parameter C controls the tradeoff between a wide margin and the penalty paid for points that fall inside the margin or on the wrong side. A small C tolerates more violations, while a very large C approaches the hard-margin solution.
