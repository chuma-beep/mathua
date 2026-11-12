# Convex quadratics as cone programs

A convex quadratic becomes a second-order cone constraint once a square root is moved into a norm. The move is exact when the quadratic is positive semidefinite.

## The epigraph step

For a positive definite matrix the quadratic form is a squared norm under its square root. Bounding the quadratic by t is then the same as bounding that norm by the square root of t.

## The rotated cone

A product of two nonnegative scalars and a squared norm forms a rotated cone. It turns a term such as s squared over t into a single cone membership.

## When the reformulation is exact

The cone form represents the same feasible set when the matrix is positive semidefinite. An indefinite matrix has no such representation, because its quadratic takes both signs.

To check your understanding, introduce one epigraph variable per quadratic, multiply the two scalars and check the smallest eigenvalue.

Check positive semidefiniteness before writing the cone. A factor of the matrix supplies the norm, and the bound supplies the right-hand side.
