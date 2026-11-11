# LP relaxation

The LP relaxation drops the integrality restriction and solves the remaining linear program. Its optimal value bounds the integer optimum, providing a lower bound for a minimization and an upper bound for a maximization. The gap between the two measures how much the integer restriction costs.

## Solving a relaxation

Remove the integer domains, keep the linear constraints and solve the resulting linear program. If the optimal solution happens to be integral, it is also optimal for the integer problem. Otherwise it supplies a bound and a starting point for branching.

## Comparing bounds

For a minimization, the relaxation value is at most the integer value. The difference is the integrality gap for that instance. Comparing the relaxation of two formulations of the same problem shows which one is tighter, and the tighter one prunes more in a branch-and-bound search.

## Testing an integrality gap

Compute the relaxation value and an integer feasible value. The gap is their difference, and a zero gap certifies that the integer solution is optimal. A positive gap says only that the integer solution may be improvable, not that it is.

The relaxation is the main source of bounds in discrete optimization. Every branch-and-bound node, every rounding argument and every integrality gap claim starts from it.
