# Cholesky decomposition

The Cholesky decomposition writes a positive-definite matrix as L times L transpose, with L lower triangular. It is the matrix square root used in least squares and simulation.

## The Cholesky factor

For a 2x2 positive-definite matrix, the factor L has the form p, 0, q, r. The (1,1) entry of A is p squared. The (2,1) entry is p times q, and the (2,2) entry is q squared plus r squared.

## Reading the factor

The diagonal entries of L are positive. The first diagonal entry is the square root of the (1,1) entry of A. The lower entry is the (2,1) entry divided by the first diagonal entry.

## Determinant

The determinant of A is the square of the determinant of L. For the 2x2 factor this is p squared times r squared. The determinant is positive because A is positive definite.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
