# Second-order cone programs

The second-order cone in n dimensions holds the pairs of a scalar and a vector whose norm is at most the scalar. Optimisation over it models norms, robust constraints and trust regions.

## The cone

Membership is a norm inequality. The cone is convex, pointed and self-dual, and its boundary is a curved surface rather than a facet.

## The dual

The dual cone is the same cone. A dual pair has a scalar at least as large as the norm of its vector, and the complementarity product is the scalar product minus the vector inner product.

## A simple model

Minimising a scalar subject to a norm bound and a linear relation is a second-order cone program. The optimal scalar is the norm of the point the relation fixes.

To check your understanding, test a point for membership, write a dual pair and fix the linear relation.

A second-order cone program is a linear objective over norm bounds, solved by the same machinery as linear programming.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
