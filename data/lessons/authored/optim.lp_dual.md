# The dual linear program

Every linear program has a dual, built from the same data with the roles of objective and right-hand side exchanged. The dual has one variable per primal constraint and one constraint per primal variable. Weak duality says every dual-feasible value bounds every primal-feasible value, which is the source of certificates of optimality.

## Writing the dual

For a primal with a minimizing objective, equality constraints and non-negative variables, the dual maximizes the product of the right-hand side with the dual variables subject to the transpose of the constraint matrix times those variables being at most the objective costs. The direction of each dual constraint follows from the sign restriction of the corresponding primal variable.

## Reading a dual variable

A dual variable is associated with one primal constraint and measures the sensitivity of the optimal value to that constraint's right-hand side. Its sign and size are meaningful only at optimality or within a range. Reading it correctly requires knowing the primal form and the sign conventions.

## Comparing dimensions

The primal has as many constraints as the dual has variables, and as many variables as the dual has constraints. The constraint matrix of the dual is the transpose of the primal matrix. Checking those counts after a conversion catches the most common transcription error.

The dual is not a second problem to solve from scratch. It is a certificate, a sensitivity report and a bound, all encoded in the same matrix.
