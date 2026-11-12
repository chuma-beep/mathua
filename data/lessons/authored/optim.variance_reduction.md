# Variance reduction

Variance-reduced methods keep a snapshot gradient and correct it with a cheap per-sample gradient. The correction removes most of the noise.

## The estimator

The corrected gradient is the snapshot plus the difference between the current and snapshot per-sample gradients. Its expectation is the full gradient at the current point.

## The cost

One expensive snapshot pass is amortised over many cheap steps. The variance falls as the snapshot ages.

## A simple form

One common method recomputes the snapshot every epoch. The step can then be constant, which removes the shrinking-step requirement.

To check your understanding, form the corrected gradient, amortise a snapshot and use a constant step.

Variance reduction makes constant-step stochastic methods behave like full-gradient descent near the solution.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
