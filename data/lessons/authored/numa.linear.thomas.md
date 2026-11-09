# The Thomas algorithm

The Thomas algorithm solves a tridiagonal system in linear time by eliminating the lower diagonal and then substituting backward.

## Tridiagonal structure

Each equation involves only its own variable and its immediate neighbours. The matrix has one main diagonal and two off-diagonals.

Storing only these three diagonals avoids the cost of a full matrix.

The first and last equations have only one neighbour, so their off-diagonal entries are zero.

## Forward elimination

Eliminate the lower diagonal from top to bottom. Each step updates the main diagonal and the right-hand side.

The updated coefficients are sometimes called the modified pivots. They must remain nonzero for the algorithm to continue.

The work per row is constant, so the sweep is linear in the number of unknowns.

## Back substitution

After the forward sweep, the system is upper bidiagonal. Solve from the last equation upward.

The total work is linear in the number of unknowns, which is why the algorithm is preferred over general elimination for tridiagonal systems.

A diagonal pivot can be small even when the matrix is nonsingular, so pivoting may still be needed for stability.
