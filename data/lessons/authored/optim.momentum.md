# Momentum

Momentum accumulates past steps in a velocity, which damps oscillation across a narrow valley and accelerates along it.

## The recursion

The velocity is the decay factor times the old velocity minus the step times the gradient, and the point moves by the velocity. The decay factor sets how long memory lasts.

## The effective step

Along a consistent direction the velocity settles at the step times the gradient divided by one minus the decay. A decay near one multiplies the step by a large factor.

## Oscillation

Across a direction where the gradient alternates, the accumulated terms cancel. That cancellation is why momentum tolerates larger steps.

To check your understanding, update the velocity, compute the steady factor and cancel alternating terms.

Momentum is an exponential average of gradients, used as the step direction.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
