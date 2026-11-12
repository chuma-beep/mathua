# Jacobi splitting

Jacobi splitting solves a linear system by updating every coordinate from the previous iterate. The diagonal part is inverted while the off-diagonal part is applied explicitly.

## Splitting

Write the matrix as a diagonal part plus a remainder. The iteration applies the inverse diagonal to the right-hand side minus the remainder times the current vector.

## Updates

Every coordinate is updated simultaneously, so one sweep costs one matrix-vector product plus a diagonal solve. Old values are never overwritten during a sweep.

## Convergence

The iteration converges when the spectral radius of the iteration matrix is below one. Strict diagonal dominance is a sufficient condition.

To check your understanding, separate the diagonal, sweep the coordinates and compute the spectral radius.

Jacobi is the simplest splitting, and its convergence is governed entirely by a spectral radius.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
