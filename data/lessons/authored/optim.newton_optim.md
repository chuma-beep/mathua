# Newton's method for minimization

Newton's method for minimization solves the stationarity equation by linearizing the gradient. Each step multiplies the inverse Hessian by the gradient and subtracts, which is the same formula as a root-finding Newton step applied to the gradient. Near a minimum with positive curvature the convergence is quadratic, and away from it the method can fail.

## Taking a Newton step

Evaluate the gradient and the Hessian at the current point, solve the linear system whose matrix is the Hessian and whose right side is the gradient, then subtract the solution from the point. Solving rather than inverting is both cheaper and more accurate. The step is the minimizer of the local quadratic model, which is the geometric reading of the algebra.

## Testing convergence

For a quadratic, one Newton step reaches the exact minimizer, so a nonzero remaining gradient means a mistake in the Hessian or the solve. For a general smooth function, the gradient norm should fall very fast near the solution. A step that increases the objective signals that the current point is outside the region where the quadratic model is trustworthy.

## Comparing with gradient descent

Gradient descent uses a scalar step and ignores curvature, while Newton uses the full curvature matrix and rescales the space. On an ill-conditioned quadratic, Newton is exact in one step and gradient descent is slow. The cost is assembling and factoring the Hessian, which is why large problems use approximations instead.

Newton's method is the reference for fast local convergence. Its step, its failure modes and its cost all come from the local quadratic model written with the Hessian.
