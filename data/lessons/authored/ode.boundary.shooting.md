# The shooting method

The shooting method turns a boundary-value problem into an initial-value problem with an unknown initial slope.

## Choose an initial slope

Replace the missing initial derivative by a parameter m. Solve the initial-value problem from the left boundary. The computed value at the right boundary is a function of m.

## Integrate to the boundary

The integration uses any accurate ODE solver. The residual is the computed right value minus the prescribed right value. A zero residual means the chosen slope is correct.

## Update the slope

Use a root-finding step such as the secant method to update m. Repeat until the residual is small. The final trajectory satisfies both boundary conditions.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
