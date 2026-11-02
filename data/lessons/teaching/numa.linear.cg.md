# The Conjugate Gradient Method

## Quadratic Forms

A symmetric positive definite system is equivalent to minimizing a quadratic function. The gradient of that function is the residual of the linear system.

The level sets of the quadratic are ellipses whose axes are the eigenvectors of the matrix. The condition number measures how stretched those ellipses are.

Conjugate gradient chooses search directions that are orthogonal with respect to the matrix. That condition prevents the method from undoing work done in earlier directions.

## Conjugate Directions

Two directions are conjugate when their inner product through the matrix is zero. A set of conjugate directions can be followed one at a time without revisiting earlier components of the solution.

The method builds these directions from residuals. Each new direction is a combination of the current residual and the previous direction, which keeps storage to a few vectors.

In exact arithmetic the method reaches the solution in at most n steps for an n by n matrix. Rounding errors usually require a few extra iterations.

## The Algorithm

Each iteration computes a step length, updates the solution and residual, and then computes a new direction. The cost is one matrix-vector product plus a few vector operations.

The residual norm decreases monotonically, which gives a practical stopping test. A preconditioner can be applied to improve the distribution of eigenvalues and reduce the iteration count.

Conjugate gradient is the method of choice for large sparse symmetric positive definite systems. It never forms a factorization, so memory use stays linear in the number of nonzeros.
