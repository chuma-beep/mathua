# Convergence of interior-point methods

Interior-point methods converge in a number of iterations that grows with the logarithm of the accuracy rather than with the problem size.

## The iteration count

Each step reduces the complementarity gap by a constant factor. Reaching a gap from an initial gap takes about the logarithm of the ratio divided by the logarithm of the reduction.

## Work per iteration

One factorisation of the KKT matrix dominates the cost. Its size grows with the variables and constraints, so total work is the iteration count times the factorisation.

## Polynomial bounds

Worst-case bounds are polynomial in the dimension and the number of bits of the data. The logarithmic factor in accuracy is what makes the methods competitive.

To check your understanding, apply a reduction factor, estimate a factorisation and read a worst-case bound.

Interior-point complexity is logarithmic in accuracy and polynomial in size.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
