# Interpolation line search

A line search chooses a step length along a descent direction. An interpolation search fits a low-degree polynomial to the values already sampled and takes its minimiser as the next trial.

## The first trial

The initial step comes from the previous iteration or from a unit step. A quadratic fit through the value and slope at zero and the value at the trial gives the next candidate in closed form.

## The safeguards

The interpolation is bracketed by a minimum and maximum step, so a bad fit cannot produce a wild trial. A cubic fit replaces the quadratic once two points are available.

## The acceptance test

A trial is accepted when it satisfies sufficient decrease and, for a quasi-Newton method, the curvature condition. Interpolation continues until both hold or the interval collapses.

To check your understanding, reuse the previous step, bracket the step and check sufficient decrease.

An interpolation line search spends a few function values to find a step that satisfies the same conditions a bisection would reach more slowly.
