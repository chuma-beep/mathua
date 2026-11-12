# The Hamilton-Jacobi-Bellman equation

Dynamic programming defines a value function as the best remaining cost from a state and time. The Hamilton-Jacobi-Bellman equation is the partial differential equation that the value function satisfies.

## The value function

The value function records the optimal cost to the horizon from every state and time. Its gradient with respect to the state is the costate, so the value function encodes the whole family of optimal trajectories at once.

## The equation

The equation sets the time derivative of the value plus the minimum over the control of the Hamiltonian to zero, with a terminal condition at the horizon. The minimum is taken pointwise, exactly as in the Hamiltonian formulation.

## Why it is hard

The equation lives in the state space, which grows with the dimension, and its solution needs a grid or a parametric approximation. Curse of dimensionality is the reason direct methods often beat dynamic programming in practice.

To check your understanding, record the remaining cost, differentiate in time and face state-space growth.

The Bellman equation is exact and global, and its cost is the state-space grid that makes it impractical above a few dimensions.
