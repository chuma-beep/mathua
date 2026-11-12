# Quadratic programs

A quadratic program minimises a quadratic objective over linear constraints. It is the simplest constrained problem with a curved objective.

## The form

The objective is half x transposed H x plus c transposed x, and the constraints are linear. A positive semidefinite H makes the problem convex.

## The KKT system

At an optimum the gradient is balanced by a combination of active constraint gradients. The multipliers of inactive constraints are zero.

## Solving it

Active-set methods add and remove constraints. Interior-point methods solve a sequence of smooth problems and are often faster on large instances.

To check your understanding, evaluate the objective, form the stationarity condition and take an active-set step.

A quadratic program is the workhorse behind portfolio, control and fitting models.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
