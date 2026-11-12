# Optimal control problems

An optimal control problem chooses a control signal that steers a dynamical system while minimising a cost. The state evolves under a differential equation, the control is bounded, and the cost accumulates over time.

## The three ingredients

A dynamics model maps the state and control to a state derivative. A running cost charges the state and control at each instant. A terminal cost charges the final state, and the horizon is either fixed or free.

## The decision variable

The unknown is a function of time, not a vector. Discretising the horizon turns it into a finite vector, and the discretisation is a modelling choice that changes the feasible set.

## Constraints

Control bounds, state bounds and path constraints all appear as pointwise inequalities. A path constraint couples every instant, which is what separates optimal control from a stack of independent optimisations.

To check your understanding, define the dynamics, choose a control function and bound the control.

An optimal control problem is an optimisation over functions, and every practical method is a choice about how to make that finite.
