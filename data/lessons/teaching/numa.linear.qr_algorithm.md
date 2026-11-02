# The QR Algorithm

## The Basic Iteration

The QR algorithm repeatedly factors a matrix into an orthogonal factor and an upper triangular factor. It then multiplies the factors in the reverse order.

Each step preserves eigenvalues because the new matrix is orthogonally similar to the old one. Under mild conditions the off-diagonal entries shrink and the matrix approaches triangular form.

The diagonal entries of the limiting triangular matrix are the eigenvalues. For a symmetric matrix the limit is diagonal, and the orthogonal factors accumulate the eigenvectors.

## Shifts

A shift subtracts a scalar from the diagonal before the factorization and adds it back afterwards. A good shift greatly accelerates the disappearance of the subdiagonal entries.

The Rayleigh quotient shift uses the bottom-right entry. The Wilkinson shift chooses the eigenvalue of the trailing two by two block that is closer to that entry.

Deflation is applied once a subdiagonal entry becomes negligible. The problem then splits into smaller blocks, and the work shrinks at every stage.

## Reading Eigenvalues

For an upper triangular matrix the eigenvalues are simply the diagonal entries. The QR algorithm drives a general matrix toward that form through orthogonal similarity.

A real matrix may have complex eigenvalues, which appear as two by two blocks in the real Schur form. Those blocks are solved as quadratics rather than forced to become diagonal.

The algorithm is the standard dense eigensolver because it is backward stable and its cost is cubic in the matrix size. The preliminary reduction to Hessenberg form makes each QR step much cheaper.
