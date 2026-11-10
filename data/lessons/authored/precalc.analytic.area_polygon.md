# Area of a polygon from coordinates

The shoelace formula uses the coordinates of a polygon directly. It avoids splitting the figure into triangles.

## The shoelace formula

List the vertices in order and repeat the first vertex. Sum the cross products of consecutive coordinates. Half the absolute value of the difference is the area.

## Triangles and rectangles

For a triangle the formula reduces to half the absolute determinant. For a rectangle aligned with the axes it reduces to width times height. The general formula handles every simple polygon.

## Missing coordinates

A missing coordinate can be recovered from a known area. Substitute the known vertices into the shoelace expression and solve the resulting linear equation. Check that the recovered point keeps the polygon simple.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving on.
The same three moves apply to every generated instance of this concept. Keep the arithmetic exact and check the sign or range before submitting.
