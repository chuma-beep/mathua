# Residuals of iterative solvers

The residual of an approximate solution measures how far the system is from being satisfied. It is computable without knowing the error.

## Definition

For an approximate solution x, the residual is b minus Ax. Its norm bounds the error through the condition number of A.

## Recursion

Iterative methods often update the residual recursively. This avoids a fresh matrix-vector product at every step but accumulates rounding error.

## Stopping

A relative residual below a tolerance is the usual stopping rule. The tolerance must be read against the scale of b and A.

To check your understanding, compute b minus Ax, update the residual and set a tolerance.

The residual is the only error measure an iterative solver can see directly.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
