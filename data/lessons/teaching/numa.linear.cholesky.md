# Cholesky Factorization

## Symmetric Positive Definite Matrices

A symmetric matrix is positive definite when every nonzero vector v satisfies v^T A v > 0. This condition guarantees that all eigenvalues are positive and that the matrix has a square-root factor with real entries.

Positive definiteness is what makes the Cholesky factorization reliable. Gaussian elimination can still be applied, but the special structure lets us halve the work and avoid searching for pivots.

In practice a matrix is often positive definite because it is a Gram matrix or a covariance matrix. Those constructions make the quadratic form a sum of squares, which is positive unless the vector is zero.

## The Cholesky Factor

The factorization writes A = L L^T with L lower triangular and positive diagonal entries. For a two by two matrix the first diagonal entry is the square root of a, and the off-diagonal entry is b divided by that root.

The remaining diagonal entry is the square root of c minus the square of the off-diagonal factor. Every quantity under a square root is positive when A is positive definite, so the algorithm needs no special pivot strategy.

Because only one triangle is stored, the method uses about half the memory and half the arithmetic of a general LU factorization. The factor is unique once the diagonal is required to be positive.

## Solving with the Factor

A linear system is solved in two triangular passes. First solve L y = b by forward substitution, then solve L^T x = y by back substitution.

Each pass costs O(n^2) operations, so the factorization dominates the total cost. Reusing the factor is therefore valuable when several right-hand sides share the same matrix.

The same factor also gives a cheap determinant, the product of the squared diagonal entries, and a test for positive definiteness. If a diagonal entry fails to be positive, the matrix is not positive definite.
