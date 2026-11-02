# Row Reduction and Echelon Forms

Row reduction is the standard algorithm for solving a linear system. It replaces the original equations with simpler equivalent equations, so the solution set does not change.

## Elementary Row Operations

Three operations preserve the solution set of a linear system. Swap two rows, multiply a row by a nonzero scalar, and add a multiple of one row to another row. Each operation can be undone, so every intermediate system is equivalent to the original one.

On an augmented matrix these operations act on whole rows. The vertical bar separates the coefficient matrix from the constants, and it is only a visual reminder that the last column is not a coefficient.

## Row Echelon Form

A matrix is in row echelon form when every nonzero row starts with a leading one, the leading one of a lower row is to the right of the leading one above it, and zero rows sit at the bottom. The leading entries are the pivots. The number of pivots is the rank of the coefficient matrix.

Gaussian elimination reaches row echelon form. Continuing to clear the entries above each pivot gives reduced row echelon form, where every pivot is the only nonzero entry in its column.

## Reading Solutions from a Reduced Matrix

In reduced row echelon form, pivot columns name basic variables and the other columns name free variables. Move each free variable to the right side and read the basic variables from the rows. If a row reads `0 = 1`, the system is inconsistent. If there are no free variables, the solution is unique.
