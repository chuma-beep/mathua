# Conjugate gradients

Conjugate gradients solves a symmetric positive definite system by minimising the quadratic form along mutually conjugate directions. It needs one matrix-vector product per step.

## Directions

Search directions are chosen to be conjugate with respect to the matrix. Each step is optimal along its direction and never spoils earlier directions.

## Short recurrence

The method needs only the current residual and direction, not the whole Krylov basis. This keeps storage constant.

## Finite termination

In exact arithmetic the method reaches the solution in at most n steps. Rounding usually delays termination, so a residual tolerance is used.

To check your understanding, choose conjugate directions, update the residual and bound by n steps.

Conjugate gradients is the Krylov method that fits in a fixed amount of memory.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
