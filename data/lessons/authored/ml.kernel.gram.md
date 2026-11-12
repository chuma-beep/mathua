# Gram matrices

A Gram matrix collects all pairwise kernel values on a data set. It is square and symmetric, and it is positive semidefinite for every valid kernel.

## Size and symmetry

For \(n\) points the matrix has \(n^2\) entries. Symmetry leaves \(n(n+1)/2\) entries to store.

## Positive semidefiniteness

Every eigenvalue of a Gram matrix is non-negative. That property is exactly what the kernel trick requires.

## Rank

The rank is at most the number of distinct feature vectors. A linear kernel on \(n\) points in \(d\) dimensions has rank at most \(d\).

To check your understanding, count entries, check eigenvalues and bound the rank.

The Gram matrix is the kernel viewed as a finite quadratic form.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
