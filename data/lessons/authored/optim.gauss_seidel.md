# Gauss-Seidel iteration

Gauss-Seidel uses each new value as soon as it is computed. The sequential order usually halves the iteration count of Jacobi.

## Sequential updates

The first block uses the previous sweep, the second uses the new first block, and so on. Information travels through the sweep.

## The triangular split

The splitting keeps the lower triangle on the left and the strict upper triangle on the right. The iteration matrix is the inverse of the lower triangle times the upper one.

## Comparing with Jacobi

Both methods converge for the same diagonally dominant systems. Gauss-Seidel typically needs fewer sweeps but cannot be parallelised across blocks.

To check your understanding, use the newest values from the current sweep, apply sweeps until the residual is small and count the Jacobi sweeps.

The order of the sweep is part of the algorithm. Reordering the blocks changes the iterates even though the fixed point is the same.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
