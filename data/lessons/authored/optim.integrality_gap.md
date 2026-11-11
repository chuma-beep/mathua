# The integrality gap

The integrality gap is the ratio or difference between the integer optimum and the LP relaxation optimum for a problem family. It measures how much the linear relaxation loses by ignoring integrality. A small gap justifies a rounding heuristic, and a large one warns that branching will be expensive.

## Computing a gap

Solve the relaxation and the integer problem, or bound the integer value, then compare. For a minimization the gap is the integer value minus the relaxation value, and its ratio is often reported. Both forms are useful, and the report must say which one it uses.

## Comparing two formulations

Two formulations of the same integer problem can have different relaxation values. The one with the higher relaxation value for a minimization is tighter and generally better for branch and bound. Comparing them requires the same objective sense and the same feasible integer set.

## Testing a rounding bound

If rounding a relaxed solution increases the objective by at most a known amount, that amount bounds the gap for every instance of the family. A constant-factor rounding argument gives a multiplicative gap. The bound is a statement about the algorithm, not about a single instance.

The gap connects a continuous relaxation to a discrete answer. It is the number that decides whether rounding is good enough or branching is unavoidable.
