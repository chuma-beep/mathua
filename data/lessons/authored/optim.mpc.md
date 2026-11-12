# Model predictive control

Model predictive control solves a finite-horizon optimal control problem at every sample, applies the first control move, and repeats. The horizon recedes, so the plan is always recomputed from the measured state.

## The receding horizon

At each step the controller solves a problem over a window of future steps using the current state as the initial condition. Only the first input is applied, and the next measurement restarts the problem.

## Constraints

State and input bounds enter the finite-horizon problem directly, which is the main advantage over a fixed feedback law. The bounds make the problem a quadratic program with inequality constraints at every step of the window.

## Feasibility and stability

A terminal constraint or a terminal cost is what turns repeated feasibility into a stability proof. Without one the controller can be feasible at every step and still fail to converge.

To check your understanding, solve over a window, bound the state and add a terminal constraint.

Predictive control is a strategy for embedding constraints, and a terminal condition is what makes the repeated solve a stabilising controller rather than a sequence of good decisions.
