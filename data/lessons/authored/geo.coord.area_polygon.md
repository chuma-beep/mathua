# The shoelace formula

The shoelace formula computes a polygon's area from its vertex coordinates by summing cross products around the boundary.

## The formula

Sum x_i y_{i+1}-x_{i+1} y_i around the vertices, then halve the absolute value.

## Orientation

Reversing the vertex order changes the sign of the sum but not the area, so the absolute value removes the orientation.

## Triangles

For a triangle the shoelace sum is twice the area, which matches the determinant formula.

To check your understanding, list the vertices in order, reverse the order and take three vertices.

The formula works for any simple polygon whose vertices are listed in order.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
