# Semidefinite programs

A semidefinite program optimises a linear function of a symmetric matrix subject to a linear matrix inequality. The cone is the set of positive semidefinite matrices.

## The cone

A symmetric matrix is positive semidefinite when every eigenvalue is nonnegative. The cone is convex, and its interior is the positive definite matrices.

## The constraint

A linear matrix inequality writes a symmetric matrix that depends affinely on the variables as positive semidefinite. It generalises a scalar inequality to a matrix order.

## The dual

The dual variable is a positive semidefinite matrix, and the duality gap is the trace of a product of two such matrices. Strong duality holds under a strict feasibility condition.

To check your understanding, test a diagonal matrix, evaluate an affine matrix and form the dual matrix.

Semidefinite programming turns eigenvalue conditions into convex constraints.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
