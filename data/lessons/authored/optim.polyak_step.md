# The Polyak step size

The Polyak step uses an estimate of the optimal value to set the step. It removes the tuning that a plain diminishing schedule needs.

## The unknown optimum

The step formula divides the gap to the estimated optimum by the squared norm of the subgradient. The estimate can be the best value seen so far.

## The step formula

A larger gap gives a longer step and a larger subgradient gives a shorter one. The step is exactly the distance that would reach the target along the chosen direction.

## The stopping test

When the best value is within the tolerance of the estimate the method stops. The gap also supplies a certificate of near optimality.

To check your understanding, compute the value gap, compare the new value with the best and subtract the best value from the estimate.

The estimate must not be below the true optimum. A lower estimate makes the step too long and can destabilise the method.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
