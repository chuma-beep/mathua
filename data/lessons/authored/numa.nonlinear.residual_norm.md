# Residual norms

The residual measures how far the current iterate is from satisfying the equation. Its norm gives a computable stopping test, though it bounds the error only through the conditioning.

## Definition

For \(F(x)=0\) the residual is \(F(x)\); for \(Ax=b\) it is \(b-Ax\). The Euclidean norm is the square root of the sum of squares.

## Relative test

Dividing by the norm of the right-hand side gives a scale-free measure. A tolerance of \(10^{-8}\) asks for eight significant digits in the residual.

## Relation to error

The error is bounded by the inverse Jacobian times the residual. An ill-conditioned problem can have a small residual and a large error.

To check your understanding, form the residual, divide by the right-hand side and relate residual to error.

A small residual is evidence of convergence, not proof of accuracy.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
