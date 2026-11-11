# The Lagrange dual function

The Lagrange dual function takes multipliers and returns the smallest value of the Lagrangian over the primal variables. For any multipliers that keep the constraints penalized in the right direction, that value is a lower bound on the primal optimum. Maximizing the bound over the multipliers is the dual problem.

## Evaluating the dual function

Form the Lagrangian by adding each constraint multiplied by its multiplier to the objective. For fixed multipliers, minimize over the primal variables. For a quadratic objective the minimizer is available in closed form, so the dual value can be computed exactly. The result is a number for each multiplier vector.

## Minimizing over a primal variable

The minimization may be unbounded below, which makes the dual value minus infinity and that multiplier useless as a bound. Otherwise the stationary point of the Lagrangian gives the minimum, and substituting it back removes the primal variables. The remaining expression depends only on the multipliers.

## Finding a dual bound

Any feasible multiplier vector gives a bound. Evaluating the dual function at several vectors and keeping the largest gives the best bound found. The dual problem is the search for the largest one, and weak duality guarantees it never exceeds the primal optimum.

The dual function turns constraints into prices and produces bounds that can be evaluated without solving the primal. It is the object that both duality theory and decomposition methods act on.
