# Preconditioning

A preconditioner multiplies the system by an approximate inverse. A good one reduces the condition number and with it the iteration count.

## The idea

Solving a transformed system with the same solution can be easier. The preconditioner should be cheap to apply and close to the inverse of the matrix.

## Condition number

The iteration count of a first-order or Krylov method grows with the condition number. A preconditioner that lowers it lowers the count.

## Cost

Each iteration pays for one preconditioner application. A strong preconditioner that is expensive can lose to a weak cheap one.

To check your understanding, apply a preconditioner, compute a condition number and cost an application.

Preconditioning is the art of changing the problem without changing its answer.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
