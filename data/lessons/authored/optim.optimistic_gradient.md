# Optimistic gradient

Optimistic gradient predicts the next gradient and uses the prediction in the current step. It keeps the last gradient, adds it with a negative sign to the current one and thus approximates the extragradient correction with one evaluation per step.

## The extrapolation

The update direction is twice the current gradient minus the previous gradient. When the gradient changes slowly the extrapolation is close to the trial-point gradient, and the two-step recurrence damps oscillation.

## Convergence conditions

The method converges for monotone operators under a step bound, like extragradient, and it keeps a single gradient in memory. The guarantee is asymptotic for the last iterate and in the average for general monotone problems.

## Where it is used

Generative adversarial training and primal-dual imaging are the common settings, because both are saddle problems with cheap gradients and a strong oscillation.

To check your understanding, double the current gradient, assume monotonicity and name adversarial training.

Optimistic gradient is the cheap relative of extragradient: one gradient per step, one extra vector in memory and the same rotational damping.
