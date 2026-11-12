# Semidefinite relaxations

A quadratic program in a vector becomes linear in a matrix once the vector is lifted. Dropping the rank-one condition leaves a semidefinite program that bounds the original.

## The lift

The product of the vector with itself is a symmetric matrix of rank one. Every quadratic form in the vector becomes a linear function of that matrix.

## Why the relaxation is looser

The lifted matrix is required only to be positive semidefinite, not rank one. That larger set contains the original points, so its optimum can only improve.

## Recovering a vector

A rank-one solution factorises back into a vector. A higher-rank solution needs rounding, and the rounding step is where the approximation factor enters.

To check your understanding, count the entries of a symmetric matrix, take the trace of the lifted matrix and check the leading minors.

The lift replaces a nonconvex rank condition with a convex cone condition. The price is a gap between the relaxed value and the true one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
