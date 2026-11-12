# Jacobi iteration

Jacobi iteration updates every block from the previous iterate. The updates are independent, so the method parallelises across blocks.

## Simultaneous updates

Each block is solved using the values from the end of the previous sweep. No block sees a value produced during the current sweep.

## The iteration matrix

A linear system splits into a diagonal part and the rest. The iteration matrix is the inverse of the diagonal times the remainder.

## Spectral radius

The method converges exactly when the spectral radius of the iteration matrix is below one. Diagonal dominance is a cheap sufficient condition.

To check your understanding, sum the off-diagonal magnitudes in a row, apply one sweep and read the previous iterate.

The price of independence is slower convergence than a sequential sweep. The gain is that every block can run at once.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
