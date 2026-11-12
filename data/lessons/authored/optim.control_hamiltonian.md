# The control Hamiltonian

The Hamiltonian packages the running cost and the dynamics into one scalar function of state, control and costate. Minimising it over the control at each instant gives the optimal control as a pointwise decision.

## The definition

The Hamiltonian is the running cost plus the inner product of the costate with the state derivative. The costate plays the role of a multiplier attached to the dynamics at every instant.

## Pointwise minimisation

At an interior optimum the derivative of the Hamiltonian with respect to the control vanishes. With a control bound the minimisation becomes a constrained problem at each instant, and the solution may sit on the boundary.

## The costate equation

The costate obeys a differential equation driven by the derivative of the Hamiltonian with respect to the state, integrated backward from the terminal condition. State and costate together form a two-point boundary value problem.

To check your understanding, add the running cost, differentiate in the control and integrate backward.

The Hamiltonian converts a calculus of variations problem into a pointwise minimisation plus a boundary value problem, which is the form every indirect method solves.
