# Testing cone membership

A second-order cone holds a scalar and a vector when the scalar dominates the vector's norm. Membership is a square-root comparison, so it is cheap to evaluate.

## The second-order cone

The cone is the set of pairs whose first entry is at least the Euclidean norm of the remaining entries. It is convex, pointed and closed under positive scaling.

## Testing a point

Square both sides to compare the scalar with the squared norm. Squaring is safe here because both sides are nonnegative.

## The boundary

The difference between the scalar and the norm is the cone margin. A zero margin is the boundary and a negative margin is outside.

To check your understanding, square each tail entry, compare the scalar with the norm and take the norm of the tail.

The test needs no square root once both sides are squared. The margin still uses the norm itself when its size matters.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
