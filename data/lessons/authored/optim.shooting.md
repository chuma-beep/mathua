# Single shooting

Single shooting treats the control sequence as the only decision variable and simulates the dynamics to obtain the state. The optimisation therefore lives in a low-dimensional space but sees a nonlinear map.

## The parameterisation

The control is discretised into a fixed number of parameters. For any parameter vector the dynamics are integrated forward, which produces the entire state trajectory and the cost.

## The sensitivity

The gradient of the cost with respect to the controls is computed by propagating derivatives through the integrator, either by finite differences or by an adjoint pass. The adjoint costs one backward integration regardless of the number of controls.

## The weakness

A small change in an early control can move a late state a long way, so the problem is badly conditioned and sensitive. That is why shooting struggles with unstable dynamics and long horizons.

To check your understanding, discretise the control, propagate derivatives and amplify early changes.

Shooting is simple to state and hard to condition, and the conditioning is a property of the dynamics rather than of the solver.
