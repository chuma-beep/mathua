# A-stability

A-stability asks that a method's stability region contain the whole left half-plane. That is exactly the requirement for a method to step through a stiff decaying mode without being forced to shrink the step.

## The test equation

The scalar equation with a complex rate is the model for every decaying mode. The step size and the rate combine into one complex number, and the method acts on it through a stability function.

## Stability function

The numerical solution after one step is the stability function evaluated at that complex number times the previous value. The method is stable on the mode when the modulus of the function is at most one.

## Stability on the half-plane

The condition must hold for every complex number with negative real part. Backward Euler and the trapezoidal rule satisfy it, and no explicit Runge-Kutta method can, because a polynomial cannot stay bounded on the whole half-plane.

To check your understanding, write the scalar test problem, apply one step and state the half-plane condition.

A-stability is the dividing line between methods that can take long steps through stiff decay and methods that cannot.
