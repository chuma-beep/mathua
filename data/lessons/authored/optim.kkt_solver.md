# KKT linear solvers

A constrained Newton step is a linear system with the KKT matrix. The matrix is symmetric and indefinite, so a Cholesky factorisation does not apply and a direct solver must tolerate a zero diagonal block.

## Block elimination

The primal block is eliminated first when it is invertible, which leaves a smaller system in the multipliers. That system is the Schur complement, and it is positive definite when the constraints are independent and the primal block is positive definite.

## Symmetric indefinite factorisation

An LDL factorisation with symmetric pivoting handles the full KKT matrix without elimination. The pivoting keeps the factorisation stable, and the inertia of the diagonal block reveals whether the reduced Hessian is correct.

## Which to choose

Block elimination is faster when the primal block factorises cheaply and the multiplier count is small. The indefinite factorisation is more robust when the primal block is singular or the constraints are nearly dependent.

To check your understanding, eliminate the primal block, factor the full matrix and count multipliers.

The KKT solve is the inner loop of every second-order constrained method, and the choice between elimination and indefinite factorisation decides its robustness.
