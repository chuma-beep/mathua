# Gauss-Jordan elimination

Gauss-Jordan elimination reduces an augmented matrix to reduced row echelon form, so the solution is read directly from the last column.

## The augmented matrix

Write each equation as a row, with the constants in an extra column. Row operations preserve the solution set.

The goal is a pivot of 1 in each variable column and zeros elsewhere in that column.

A row swap is a valid operation and may be needed to bring a nonzero entry into the pivot position.

## Elimination and back substitution

Gauss-Jordan elimination removes variables both above and below each pivot. The result needs no separate back substitution.

The determinant is the product of the pivots before normalization, with a sign change for each row swap.

The extra work above the pivots is what distinguishes Gauss-Jordan from ordinary Gaussian elimination.

## Singular systems

A zero pivot with a nonzero constant in its row means the system is inconsistent. A zero row of zeros means infinitely many solutions.

The coefficient determinant is zero exactly when the system does not have a unique solution.

A nearly zero pivot can cause instability even when the matrix is not exactly singular.
