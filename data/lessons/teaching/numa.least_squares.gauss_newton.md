# The Gauss-Newton Method

## Nonlinear Least Squares

Nonlinear least squares minimizes a sum of squared residuals that depend nonlinearly on the parameters. The Gauss-Newton method linearizes the residuals at the current estimate.

The linearization turns each iteration into a linear least-squares problem for the parameter step. The step is then added to the current estimate and the process repeats.

The method works best when the residuals are small at the solution. In that case the neglected second-order terms are unimportant and convergence is fast.

## The Normal Equations

The step solves the normal equations involving the Jacobian of the residuals. The matrix of the system is the transpose of the Jacobian times the Jacobian.

That matrix is square with one row and one column per parameter. It is symmetric and positive semidefinite, and it is positive definite when the Jacobian has full column rank.

The right-hand side is the Jacobian transpose applied to the current residual. Solving the system gives the parameter change that best reduces the linearized residual.

## Convergence

For a linear model the Jacobian is constant, so the first step solves the problem exactly. For a nonlinear model several iterations are normally required.

A line search or trust region is often added because a full step can increase the residual. The damped version is known as Levenberg-Marquardt.

The method can stall when the Jacobian loses rank or when the model is badly scaled. Column scaling and a rank-revealing solver make it more robust in practice.
