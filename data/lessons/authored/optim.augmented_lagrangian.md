# Augmented Lagrangian

The augmented Lagrangian adds a multiplier term to the quadratic penalty. The multiplier carries the first-order information that the penalty alone lacks, so the method can converge without sending the weight to infinity. Each outer iteration updates the multiplier and the penalty weight, and each inner iteration minimizes the augmented objective.

## Updating a multiplier

After minimizing the augmented objective for the current multiplier, update it by adding the weight times the constraint value. A feasible constraint leaves it unchanged. The update is a dual ascent step, and it is what drives the multiplier toward its optimal value.

## Evaluating the augmented objective

Take the original objective, add the multiplier times the constraint value and add half the weight times the squared constraint value. At a feasible point the last two terms collapse to the multiplier times zero. The value is used only through minimization, so an additive constant does not matter.

## Testing a convergence step

Convergence is declared when the constraint violation is below one tolerance and the change between multipliers is below another. Both are needed, because a small violation with a still-moving multiplier means the dual has not settled. The weight is increased only when the violation fails to fall fast enough.

The augmented Lagrangian fixes the accuracy problem of a pure penalty while keeping the unconstrained inner solve. It is the standard method behind many constrained solvers.
