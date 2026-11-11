# The KKT matrix

The KKT matrix assembles the Hessian of the Lagrangian and the constraint Jacobian into one block matrix. Solving a linear system with it gives the Newton step for an equality-constrained problem together with the new multipliers. Its structure determines whether the step is unique.

## Assembling the matrix

Place the Hessian of the Lagrangian in the top-left block, the constraint Jacobian in the top-right and its transpose in the bottom-left, with zeros in the bottom-right. The right-hand side combines the negative gradient and the constraint residual. The sign conventions must match the definition of the Lagrangian.

## Testing nonsingularity

The KKT matrix is nonsingular when the Hessian is positive definite on the null space of the Jacobian and the Jacobian has full row rank. A rank-deficient Jacobian means the constraints are dependent and the multipliers are not unique. Checking the rank is part of checking the model.

## Solving a small system

For a small problem, write the block equations and eliminate the step to get a reduced system in the multipliers, then back-substitute. The reduced system has the size of the number of constraints, which is often much smaller than the number of variables. This elimination is the Schur complement step.

The KKT matrix is the linear algebra of constrained optimization. Its blocks, rank and factorization explain both Newton steps and the conditions under which multipliers are unique.
