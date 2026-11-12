# Learning-rate decay

A constant step settles at a noise floor. Decaying the step over time lowers that floor at the cost of slower progress later in training.

## Step decay

The rate is multiplied by a factor at fixed milestones. It is simple to schedule and easy to reason about, but the drops are abrupt and the milestones are chosen in advance.

## Polynomial decay

The rate follows a power law in the step count, often with an exponent between one half and one. It gives a smooth curve and matches the rates used in convergence proofs.

## Exponential decay

The rate falls by a constant factor each epoch. It reaches a target accuracy quickly and then stalls, which is why it is often paired with a floor or a warm restart.

To check your understanding, pick milestones, choose an exponent and decay each epoch.

The schedule encodes a belief about how the noise and the curvature change as the iterate approaches the solution, and the proof and the practice should agree on that belief.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
