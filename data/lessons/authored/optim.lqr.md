# The linear quadratic regulator

The linear quadratic regulator is the optimal control problem with linear dynamics and a quadratic cost. It is the one optimal control problem with a closed-form solution, and that solution is a state feedback law.

## The problem

The state evolves under a linear map plus a linear control input. The cost penalises state deviation and control effort with positive semidefinite and positive definite weights, and the horizon is often infinite.

## The feedback law

The optimal control is a constant matrix times the state, so the control is computed from the current state rather than from a precomputed trajectory. The matrix comes from the solution of an algebraic Riccati equation.

## The weights as design

The state weight and the control weight trade settling speed against effort. Changing them changes the closed-loop poles, so tuning the regulator is a pole placement exercise expressed in cost terms.

To check your understanding, use linear dynamics, multiply the state by a matrix and trade speed against effort.

The regulator is the reference case for optimal control: a closed-form feedback law whose only design freedom is the pair of weights.
