# Feasible and infeasible sets

The feasible set is every point that satisfies all constraints at once. It can be a single point, a segment, a polygon, a curve or the empty set, and the last case is not a failure of the solver but a property of the model. Testing a point is the cheapest operation in optimization and the one that catches the most mistakes.

## Testing a point

Substitute the coordinates into every constraint and record whether each one holds. A point is feasible only when all of them hold, including the variable bounds, which are easy to forget because they are often written once at the top of a model. Report the first violated constraint with its left and right values, since that pair is what makes a rejection actionable.

## Finding the feasible interval

In one variable, each linear constraint becomes a bound, and the feasible interval is the intersection of those bounds. Solve each inequality for the variable, keep the tightest lower and upper bound, and the interval is feasible exactly when the lower bound does not exceed the upper bound. Two constraints that point away from each other produce an empty interval, which is the one-dimensional picture of infeasibility.

## Deciding emptiness

A set defined by linear inequalities is empty when the inequalities contradict each other. In one variable this is a bound comparison. In two variables it is a question of whether the half-planes share a point, which can be decided by checking the vertices or by a certificate of infeasibility. Equality constraints reduce the dimension, so a pair of parallel equalities with different constants is immediately empty.

Feasibility is a separate question from optimality. A solver that reports no solution may have found an empty feasible set, and distinguishing that from a numerical failure changes what the modeller should fix.
