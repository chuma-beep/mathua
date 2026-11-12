# The central path

The central path is the curve of points that minimise a barrier-penalised objective as the barrier parameter falls. It turns an inequality-constrained problem into a sequence of smooth ones.

## The barrier family

For a problem with m inequalities the barrier objective is f(x) minus mu times the sum of the logarithms of the slacks. Each positive mu gives one point on the path, and the point moves as mu changes.

## Following the parameter

Reducing mu pushes the point toward the boundary and toward the constrained optimum. The duality gap along the path is bounded by m mu, so a small parameter certifies a small gap.

## Reading a point

At a path point the barrier gradient vanishes. The multiplier estimate is mu divided by the matching slack, which is why multipliers grow as slacks shrink.

To check your understanding, write the barrier objective, predict the direction of movement and solve the stationarity equation.

Central-path methods trade one hard constrained problem for a sequence of smooth unconstrained ones.
