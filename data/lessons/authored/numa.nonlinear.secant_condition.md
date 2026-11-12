# The secant condition

The secant condition requires the updated model to reproduce the observed change along the last step. It is one equation for a whole matrix, so many updates satisfy it.

## Statement

\(B_{k+1}s_k=y_k\), with \(s_k=x_{k+1}-x_k\) and \(y_k=F(x_{k+1})-F(x_k)\). In one dimension it determines the slope exactly.

## Underdetermined in higher dimensions

The condition fixes one direction of an \(n\times n\) matrix. A minimum-change principle selects the rank-one update among all solutions.

## Role in convergence

Satisfying the condition at every step is what produces superlinear convergence. Without it the iteration can stall.

To check your understanding, write the step and change, count equations and unknowns and enforce the condition.

One scalar equation leaves most of the approximate Jacobian free.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
