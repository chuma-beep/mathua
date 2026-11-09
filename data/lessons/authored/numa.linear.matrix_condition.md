# The condition number of a matrix

The condition number measures how strongly a linear system can amplify relative error. A large condition number means an ill-conditioned problem.

## Singular values

For the 2-norm, the condition number is the largest singular value divided by the smallest. It is at least 1.

A singular matrix has a zero singular value, so its condition number is infinite.

The condition number is a property of the matrix, not of the algorithm used to solve the system.

## The diagonal case

For a diagonal matrix, the singular values are the absolute values of the diagonal entries. The condition number is the largest magnitude divided by the smallest.

Scaling a row changes the matrix but not the underlying difficulty in a simple way. The condition number captures that difficulty.

A diagonal matrix with entries of very different sizes is ill-conditioned.

## Digits of accuracy

A condition number of \(10^k\) can cost about \(k\) decimal digits of accuracy. This is a rule of thumb, not a guarantee.

Ill-conditioning is a property of the problem. A stable algorithm cannot remove it, though it can avoid adding to it.

A backward-stable algorithm solves a nearby problem exactly, which is the best guarantee available in floating point.
