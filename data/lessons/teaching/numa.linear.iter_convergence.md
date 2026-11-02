# Convergence of Iterative Methods

## Iteration Matrices

A stationary iteration rewrites the system as x_new = T x_old + c. The matrix T is the iteration matrix, and its powers determine how the error changes from step to step.

The error after k steps is T^k times the initial error. Convergence therefore depends on how quickly powers of T shrink, not on the size of the right-hand side.

Jacobi and Gauss-Seidel both fit this pattern. Their iteration matrices differ, which is why one can converge while the other diverges on the same system.

## Spectral Radius

The spectral radius is the largest absolute eigenvalue of the iteration matrix. The iteration converges for every starting vector exactly when the spectral radius is strictly less than one.

A radius below one gives geometric decay, with the error multiplied by roughly the radius at each step. A radius above one amplifies the error until the computation overflows or loses all accuracy.

The number of iterations needed for a given tolerance grows as the logarithm of the tolerance divided by the logarithm of the radius. Small improvements in the radius can therefore produce large savings.

## Practical Tests

Diagonal dominance is a convenient sufficient condition. If every row has a diagonal entry larger than the sum of the off-diagonal magnitudes, Jacobi and Gauss-Seidel both converge.

For symmetric positive definite matrices, Gauss-Seidel converges even without diagonal dominance. Jacobi may fail, so the structure of the matrix matters.

When no theorem applies, the spectral radius can be estimated numerically or observed during a trial run. A radius close to one warns that convergence will be slow.
