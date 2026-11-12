# Choosing a smoothing parameter

The smoothing parameter sets both the accuracy and the conditioning. A schedule is often better than one fixed value.

## The accuracy trade

A small parameter keeps the solution accurate but makes the gradient change quickly. A large parameter conditions the problem and blurs the answer.

## Choosing one value

A fixed value must balance the two errors. The best fixed choice often scales with the square root of the target accuracy.

## The schedule

A decreasing schedule starts smooth and sharpens over time. It reaches high accuracy without paying the conditioning cost from the first iteration.

To check your understanding, multiply the parameter by the Lipschitz constant, estimate the conditioning and halve the parameter.

The error from smoothing is linear in the parameter, while the iteration cost grows as its inverse. The balance is what a schedule exploits.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
