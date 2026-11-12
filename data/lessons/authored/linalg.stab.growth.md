# Growth factors

The growth factor of Gaussian elimination is the largest multiplier ratio produced during the factorisation. It controls the size of the backward error.

## Definition

At each elimination step compare the largest entry of the remaining submatrix with the largest entry of the original matrix. The maximum ratio is the growth factor.

## Bound

Without pivoting the growth factor can be as large as two to the n minus one. Partial pivoting keeps multipliers bounded but does not bound growth in theory.

## Practice

Large growth is rare for random matrices but can be constructed. The bound is worst case, not a typical estimate.

To check your understanding, track the submatrix maximum, raise two to n minus one and sample random matrices.

Growth measures how much elimination amplifies the entries, and pivoting is the usual control.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
