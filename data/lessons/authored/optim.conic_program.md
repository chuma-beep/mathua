# Conic programs

A conic program minimises a linear function over the intersection of an affine set and a convex cone. Linear, quadratic and semidefinite programs are conic programs.

## The standard form

The primal minimises c transposed x subject to Ax equal to b and x in a closed convex cone. The cone replaces the nonnegative orthant of linear programming.

## The dual cone

The dual cone holds the vectors with nonnegative inner product against every cone member. The dual program maximises b transposed y subject to a transposed y plus a slack in the dual cone.

## Self-dual cones

The nonnegative orthant, the second-order cone and the semidefinite cone are self-dual. Self-duality makes the primal and dual forms look alike, which is what interior-point methods exploit.

To check your understanding, write the conic constraint, write the dual objective and recognise a self-dual cone.

A conic program is a linear program whose variables live in a cone rather than an orthant.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
