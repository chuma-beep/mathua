# The BFGS update

BFGS builds an approximation to the inverse Hessian from observed steps and gradients, so a quasi-Newton method gets curvature information without assembling second derivatives. The update is a rank-two correction that keeps the approximation symmetric and positive definite when the curvature condition holds.

## Updating the inverse Hessian

Let the step be s and the gradient change be y. The update adds two outer products to the current inverse approximation, one along s and one along the product of the approximation with y. The denominators are the inner products of s with y and of y with the approximation applied to y. A positive denominator keeps the result positive definite.

## Testing the secant condition

The updated approximation should map the gradient change to the step, which is the secant equation. Testing it means multiplying the new matrix by y and comparing with s. The condition is what the update is constructed to satisfy, so a failure indicates an arithmetic error rather than a limitation of the method.

## Applying a rank-two update

A rank-two update changes the matrix along two directions rather than one. The BFGS correction is exactly of that form, which is why it can add curvature information without destroying the existing approximation. Applying it to a vector costs a few inner products and vector operations, never a matrix factorization.

BFGS is the practical compromise between gradient descent and Newton. It keeps the cost per iteration low while recovering enough curvature to converge quickly on smooth problems.
