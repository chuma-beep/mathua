# VC dimension

The VC dimension is the largest set a hypothesis class can shatter. It measures capacity without reference to a data distribution.

## Intervals on a line

Intervals shatter two points but not three. Their VC dimension is two.

## Half-spaces

Affine half-spaces in \(d\) dimensions have VC dimension \(d+1\). The extra dimension comes from the bias term.

## Rectangles

Axis-aligned rectangles in the plane shatter four points but not five. Their VC dimension is four.

To check your understanding, shatter two points, count the parameters and shatter four points.

VC dimension counts the points a class can label in every possible way.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
