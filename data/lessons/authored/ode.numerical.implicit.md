# Implicit Euler and stiff equations

Implicit Euler uses the slope at the next point. It is more stable than explicit Euler on stiff equations.

## Apply the implicit step

The update is y next equals y plus h f of t plus h and y next. The unknown appears on both sides. This implicit equation must be solved at every step.

## Solve the implicit equation

For a linear equation the solve is a division by 1 minus h times the coefficient. For a nonlinear equation a Newton iteration is typical. The extra work per step buys stability.

## Test stiffness

A stiff equation has a fast decaying mode and a slow mode. Explicit Euler needs a very small step for stability. Implicit Euler can take a much larger step and still remain bounded.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
