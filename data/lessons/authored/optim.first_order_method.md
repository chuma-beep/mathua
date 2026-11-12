# First-order methods at scale

A first-order method uses only gradients, so each iteration costs linear time and constant memory beyond the variables.

## Cheap iterations

One gradient costs about as much as one matrix-vector product. No factorisation is formed, so the cost per step stays linear in the problem size.

## Memory

Beyond the variables and one or two vectors the method needs little storage. That is what lets it run on models too large to factorise.

## Convergence

The iteration count grows with the condition number. A well-conditioned problem converges quickly and a badly conditioned one needs acceleration or preconditioning.

To check your understanding, cost one gradient, count vectors and read the condition number.

First-order methods trade iteration count for iteration cost.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
