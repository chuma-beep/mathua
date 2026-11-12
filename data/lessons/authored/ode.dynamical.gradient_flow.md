# Gradient flows

A gradient flow moves against the gradient of a potential, so the potential decreases along every trajectory. The dynamics cannot oscillate and every bounded orbit converges to a critical point.

## The equation

The velocity is minus the gradient of the potential. Stationary points of the flow are exactly the critical points of the potential, so the fixed points are read from an algebraic equation.

## Monotone decrease

The time derivative of the potential along the flow is minus the squared gradient. It is negative wherever the gradient is non-zero, so the potential is a Lyapunov function with no oscillation possible.

## Convergence

LaSalle's invariance principle says a bounded orbit converges to the largest invariant set inside the zero set of the derivative. For a gradient flow that set is the critical points, and the local curvature sets the rate.

To check your understanding, write the gradient flow, differentiate the potential and apply LaSalle's principle.

A gradient flow is the simplest dynamical system with a Lyapunov function, and it is the continuous limit of gradient descent.
