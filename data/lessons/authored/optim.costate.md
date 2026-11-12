# Costate equations

The costate is the multiplier attached to the dynamics, and it measures the marginal cost of the state. Its equation is the dual of the state equation and runs backward in time.

## The economic reading

The costate at a time is the sensitivity of the optimal cost to a small change in the state at that time. A large costate means that the state is expensive and worth correcting.

## The backward equation

The costate derivative is the negative derivative of the Hamiltonian with respect to the state. The terminal cost provides the terminal costate, so the equation is integrated from the horizon back to the start.

## The jump conditions

A state constraint that becomes active introduces a jump in the costate. The jump is the multiplier of the constraint, and ignoring it produces a trajectory that satisfies the dynamics but not optimality.

To check your understanding, read a marginal cost, differentiate the Hamiltonian and detect an active state constraint.

The costate prices the state, and its backward equation is the shadow-price system of an optimal control problem.
