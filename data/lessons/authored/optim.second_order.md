# Second-order optimality and saddle points

The second derivative decides what a stationary point is. In one variable a positive value means a minimum, a negative value a maximum and a zero value leaves the question open. In several variables the Hessian replaces the scalar, and definiteness replaces the sign. A saddle point is stationary with directions of both positive and negative curvature.

## Testing a Hessian

A symmetric Hessian at a stationary point is positive definite when every eigenvalue is positive, negative definite when every eigenvalue is negative and indefinite when the eigenvalues have both signs. In two variables the determinant and the first entry decide the common cases quickly. A zero eigenvalue is the degenerate case, where higher-order terms decide.

## Classifying a stationary point

Stationary plus positive definite is a strict local minimum. Stationary plus negative definite is a strict local maximum. Stationary plus indefinite is a saddle. Stationary plus semidefinite is inconclusive, and the classification needs either a higher-order expansion or a direct comparison of nearby values. The order of these checks matters, because definiteness is only meaningful once the gradient is zero.

## Finding a saddle direction

An eigenvector with a negative eigenvalue is a direction of negative curvature, and moving a small distance along it lowers the objective from a stationary point. A positive eigenvalue gives a direction of positive curvature, along which the point behaves like a minimum. At a saddle both exist, and the negative one is what a minimization method can exploit to escape.

First-order conditions cannot separate minima from saddles, and saddle points are common in modern objectives. The Hessian supplies the missing information, at the cost of assembling and factoring a matrix.
