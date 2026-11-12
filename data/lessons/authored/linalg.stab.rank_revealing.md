# Rank-revealing factorisations

A rank-revealing factorisation exposes the numerical rank of a matrix. Small singular values are separated from large ones by a threshold.

## Numerical rank

The numerical rank counts singular values above a tolerance. The tolerance is usually a multiple of machine epsilon times the largest singular value.

## Factorisations

QR with column pivoting and the singular value decomposition both reveal rank. The SVD gives the most reliable separation.

## Truncation

Dropping small singular values gives a low-rank approximation. The error is the largest dropped singular value.

To check your understanding, choose a tolerance, pivot the columns and drop small values.

Rank is a discrete idea, but numerical rank depends on a threshold and therefore on the data scale.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
