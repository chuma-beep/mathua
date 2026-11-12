# Stochastic gradient steps

Stochastic gradient descent replaces the full gradient with an unbiased estimate from one sample or a small batch. The step is noisy but cheap.

## The update

The new point is the old point minus the step times the sampled gradient. The sampled gradient has the true gradient as its expectation, so the noise averages out along the trajectory.

## Step sizes

A fixed step converges to a neighbourhood whose radius scales with the step and the noise. A diminishing step shrinks that radius toward zero.

## What is guaranteed

With a diminishing step and bounded variance the average iterate approaches a stationary point. Convexity strengthens the rate but is not needed for stationarity.

To check your understanding, apply one sampled step, evaluate a schedule and state the stationarity target.

Stochastic descent buys cheap steps with noise, and the step schedule decides how much noise remains.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
