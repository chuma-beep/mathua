# The Hessian matrix

The Hessian collects every second partial derivative into a symmetric matrix. Its diagonal holds the curvature along each coordinate and its off-diagonal entries hold the mixed curvature. Almost every second-order claim about a multivariate objective is a statement about this matrix rather than about its entries one at a time.

## Building the matrix

Differentiate once with respect to each variable to get the gradient, then differentiate each gradient component with respect to each variable. The entry in row i and column j is the derivative with respect to variable i and then variable j. Symmetry follows from the equality of mixed partials when the second derivatives are continuous, so only the upper triangle needs computing.

## Evaluating at a point

The entries are functions, and the Hessian at a point is the numeric matrix obtained by substituting the point into each entry. Two points of the same function can have Hessians of different definiteness, which is why a curvature claim always names a point. Evaluating first and inspecting afterwards keeps the algebra separate from the classification.

## Testing positive definiteness

A symmetric matrix is positive definite when every eigenvalue is positive, equivalently when every leading principal minor is positive. In two variables the determinant and the top-left entry suffice. Positive semidefiniteness allows zero eigenvalues and requires checking all principal minors, not just the leading ones. The test is what turns a Hessian into a conclusion about convexity or a minimum.

The Hessian is the object that second-order methods read, assemble and approximate. A correct build plus a definiteness test answers most questions about curvature in a single pass.
