# Residual bounds

A residual bound converts a computable residual into a bound on the error. The conversion uses the inverse of the matrix or its smallest singular value.

## Basic bound

The error norm is at most the norm of the inverse times the residual norm. This follows directly from multiplying the residual by the inverse.

## Relative form

Dividing by the solution norm gives a relative error bound in terms of the relative residual and the condition number.

## Backward stability

A backward stable solver makes the residual small relative to the data. The bound then explains the observed accuracy.

To check your understanding, multiply by the inverse, divide by the solution and bound the residual.

The residual is computable; the error is not, and the bound is the bridge.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
