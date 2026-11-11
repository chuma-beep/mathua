# Quadratic penalty method

A penalty method adds the squared constraint violation to the objective with a large weight and solves the resulting unconstrained problem. As the weight grows, the minimizers approach feasibility. The method is simple and robust, and it trades exact feasibility for a sequence of easier problems.

## Evaluating a penalized objective

Compute the original objective and add the weight times the square of each constraint violation. A feasible point has no addition. The penalized value is always at least the original value, and the gap measures how much the point is being charged for infeasibility.

## Taking a penalty step

For a fixed weight, minimize the penalized objective by any unconstrained method. The minimizer is the next iterate. Increasing the weight and repeating drives the violation down, at the cost of worse conditioning because the curvature grows with the weight.

## Choosing a penalty parameter

A small initial weight keeps the problem well conditioned, and a schedule that multiplies it by a constant balances the two effects. The stopping rule compares the violation with a tolerance. A single huge weight is numerically poor even though it seems to enforce the constraint.

The quadratic penalty is the simplest bridge from constrained to unconstrained optimization. Its weakness, finite accuracy at finite weight and growing conditioning, motivates the augmented Lagrangian.
