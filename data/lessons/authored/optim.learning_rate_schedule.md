# Step-size schedules

A step-size schedule fixes the step as a function of the iteration, by division, by a staircase or by a smooth curve.

## Polynomial decay

The step is the initial step divided by one plus the iteration raised to an exponent. The exponent sets how fast the step falls, and an exponent of one is the common default.

## Staircase decay

The step is multiplied by a factor at fixed milestones. Between milestones the step is constant, which makes the schedule easy to reason about.

## Warmup

A short increasing phase avoids large early steps when the estimates are poor. Warmup is common in large models and rare in small convex problems.

To check your understanding, evaluate a polynomial step, apply a milestone factor and ramp the step.

The schedule is part of the algorithm: the same update with different steps is a different method.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
