# Regularised Newton systems

A Newton system with a singular or indefinite Hessian has no useful solution. Adding a multiple of the identity makes it positive definite and bounds the step, which is the principle behind Levenberg-Marquardt.

## The damping

The shifted system solves for a step that trades the Newton direction against a gradient direction. As the damping grows the step shrinks toward the scaled negative gradient.

## The update rule

If the step reduces the model as predicted, the damping is reduced. If it does not, the damping is increased and the step is recomputed, which is a trust-region update expressed as a penalty.

## The Levenberg-Marquardt reading

For a least-squares objective the damped system interpolates between Gauss-Newton and gradient descent. The interpolation is smooth and monotone in the damping parameter.

To check your understanding, shift the Hessian, compare predicted and actual and start from Gauss-Newton.

Regularisation is how a Newton method survives a bad Hessian, and the damping update is the same trade a trust region makes explicitly.
