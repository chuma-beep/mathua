# Dual coordinate ascent

Dual coordinate ascent updates one dual variable at a time while holding the others fixed. Each update is cheap and often has a closed form.

## One coordinate at a time

Fixing every other multiplier reduces the dual to a function of one variable. Maximising that function gives the new value.

## The update

With box constraints the update is a projection of an unconstrained maximiser onto the box. Clamping the maximiser to the bounds gives the answer.

## Convergence

A sweep costs one pass over the coordinates. The method converges when the largest coordinate change falls below a tolerance.

To check your understanding, fix the other multipliers, compute the unconstrained value and measure a sweep.

Coordinate ascent is the dual counterpart of coordinate descent.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
