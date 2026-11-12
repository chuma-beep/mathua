# Finite-difference gradients

A finite difference approximates a derivative from function values at nearby points. The step size trades truncation error against rounding error, and the two errors pull in opposite directions.

## The error balance

A forward difference has truncation error proportional to the step and rounding error proportional to the reciprocal of the step. The best choice is near the square root of machine precision.

## Central differences

A central difference cancels the first-order term and leaves a quadratic truncation error. It costs two evaluations and widens the usable step range, which is why it is the default for a single derivative.

## Gradients in many dimensions

A full gradient costs one evaluation per coordinate with forward differences and two with central. That cost is what motivates automatic differentiation and adjoint methods.

To check your understanding, name truncation error, cancel the first-order term and count evaluations per coordinate.

Finite differences are the fallback when no analytic derivative exists, and their accuracy is bounded by the square root of machine precision.
