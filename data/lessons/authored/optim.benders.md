# Benders decomposition

Benders decomposition keeps the hard integer variables in a master and the continuous ones in a subproblem. Each subproblem solve adds a cut to the master.

## Complicating variables

Fixing the complicating variables leaves a subproblem that is easy to solve. Its optimal value is a function of the fixed variables.

## The cut

A subproblem dual supplies a supporting hyperplane of that value function. The master adds the hyperplane as one linear cut.

## The master sequence

The master minimises the cut model, then the subproblem evaluates the chosen point and adds another cut. The process ends when the model value meets the subproblem value.

To check your understanding, fix the master variables, evaluate the dual multipliers and compare the master value with the subproblem value.

The cuts are valid for every feasible master point. Feasibility cuts keep the master away from points where the subproblem has no solution.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
