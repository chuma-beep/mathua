# Conjugate directions

Two directions are conjugate with respect to a matrix when their inner product weighted by the matrix is zero. For a quadratic, moving exactly along each of a set of conjugate directions reaches the minimizer in at most as many steps as there are variables. The construction turns a coupled problem into independent one-dimensional minimizations.

## Testing conjugacy

Multiply one direction by the matrix, take the inner product with the other and check that the result is zero. For a diagonal matrix, conjugate directions are the coordinate axes. For a general matrix, they are skewed. The test is one matrix-vector product and one inner product, which makes it easy to verify a proposed set.

## Building a conjugate direction

Given a direction and the gradient at the new point, a combination of the two can be made conjugate to the previous direction by choosing one scalar. The scalar is the ratio of the new gradient's matrix-weighted inner product to the previous direction's. Repeating this builds a full set without ever factorizing the matrix.

## Counting exact steps

For a quadratic in n variables, n conjugate directions suffice to reach the exact minimizer, because the error is expressed in a basis that the method eliminates one component at a time. In exact arithmetic the count is n, and rounding makes it slightly more in practice. This finite termination is the property gradient descent lacks.

Conjugate directions are the idea behind conjugate gradient methods. They show how a quadratic can be solved in a finite number of cheap steps when the directions are chosen well.
