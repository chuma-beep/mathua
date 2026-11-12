# Gradient penalties

A gradient penalty adds the norm of a gradient to the objective. The penalty discourages a steep or unstable response and can enforce a constraint on the derivative without a hard projection.

## The two uses

In generative models the penalty keeps the discriminator smooth so its gradient carries signal. In constrained optimisation the penalty is a soft way to enforce a Lipschitz condition on a learned map.

## The discretisation

The gradient norm is evaluated at sampled points, often interpolated between real and generated data. The penalty is therefore a Monte Carlo estimate and has its own variance.

## The weight

The penalty weight trades constraint satisfaction against the original objective. Too small a weight leaves the constraint violated, and too large a weight flattens the model.

To check your understanding, smooth a discriminator, sample points and trade constraint against objective.

A gradient penalty is a soft Lipschitz constraint, and its weight is the only lever between fidelity and smoothness.
