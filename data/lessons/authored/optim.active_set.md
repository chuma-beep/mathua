# Active-set methods

An active-set method maintains a working set of constraints treated as equalities. It solves the equality-constrained subproblem, checks the multipliers and updates the working set by adding or dropping a constraint. The name describes the data structure as much as the algorithm.

## Choosing the working set

Start from a feasible point and a set of active constraints. Solve the equality-constrained problem on that set. If a multiplier is negative, the corresponding constraint should leave the working set. If the step violates another constraint, that one joins the set.

## Taking an active-set step

Solve the KKT system for the current working set to get a direction, then take the largest step that keeps feasibility. A blocked step adds the blocking constraint to the set. A full step with a negative multiplier drops the offending constraint. Each iteration changes the set by one constraint.

## Testing optimality on a face

On the face defined by the working set, optimality means the reduced gradient vanishes and all multipliers are non-negative. The test is the KKT test restricted to the face. A negative multiplier means a direction into the feasible region lowers the objective, so the face is not optimal.

Active-set methods are finite and exact for quadratic programs. They make the working set explicit, which makes both the algorithm and its failure modes easy to inspect.
