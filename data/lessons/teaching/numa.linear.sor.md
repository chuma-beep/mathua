# Successive Over-Relaxation

## Relaxation

Successive over-relaxation updates each component using a weighted combination of the old value and the Gauss-Seidel value. The weight is called the relaxation parameter and is usually written omega.

When omega equals one the method is exactly Gauss-Seidel. Values above one extrapolate in the direction of the update, while values below one damp the step.

The point of extrapolation is speed. A well-chosen omega can cut the number of iterations substantially, especially for the large sparse systems that arise from discretized differential equations.

## Choosing Omega

For a tridiagonal matrix there is a classical formula for the optimal parameter. It depends only on the spectral radius of the Jacobi iteration matrix.

The formula gives a value between one and two when the Jacobi radius is below one. As the radius approaches one, the optimal omega approaches two and the gain over Gauss-Seidel grows.

In practice the spectral radius is not known exactly. It is estimated during the iteration, or a small set of candidate values is tried and the fastest is kept.

## Convergence

The method converges for symmetric positive definite matrices whenever omega lies strictly between zero and two. Outside that interval the iteration generally diverges, even for very well-conditioned problems.

The convergence rate is governed by the spectral radius of the iteration matrix. Minimizing that radius is exactly what the optimal parameter does.

Over-relaxation is cheap to add to an existing Gauss-Seidel code. It requires no extra matrix storage and only a few extra multiplications per component.
