# Finite elements

The finite element method approximates the solution by piecewise polynomials on a mesh. It enforces the weak form element by element, which suits irregular geometry.

## Weak form

Multiplying by a test function and integrating by parts lowers the derivative order. Boundary conditions enter the integral rather than the approximation space.

## Basis functions

Each node carries a hat function that equals one there and zero at neighbouring nodes. In one dimension n elements give n+1 nodes.

## Assembly

Element matrices are added into the global matrix at shared nodes. The global matrix is sparse because each basis function overlaps only its neighbours.

To check your understanding, multiply by a test function, draw a hat function and assemble element matrices.

The weak form moves differentiation onto test functions and leaves a sparse system.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
