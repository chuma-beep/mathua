# Line and conic intersections

A line and a parabola can meet in zero, one or two points. Substituting the line into the parabola turns the geometry into a quadratic equation, and the discriminant decides how many real intersections exist.

## Substitution

Take the parabola \(y = x^2\) and the line \(y = mx + b\). Replacing \(y\) gives
\[
x^2 = mx + b,
\]
or
\[
x^2 - mx - b = 0.
\]
Every intersection point has an \(x\)-coordinate that satisfies this quadratic. Once an \(x\)-coordinate is known, the line gives the corresponding \(y\)-coordinate.

## The discriminant

For a quadratic \(Ax^2 + Bx + C = 0\), the discriminant is \(B^2 - 4AC\). Here \(A = 1\), \(B = -m\) and \(C = -b\), so
\[
D = m^2 + 4b.
\]
If \(D > 0\), the line cuts the parabola twice. If \(D = 0\), the line is tangent and touches it once. If \(D < 0\), there is no real intersection.

## Sums of intersection coordinates

When there are two intersections, their \(x\)-coordinates are the two roots of \(x^2 - mx - b = 0\). By Vieta, their sum is \(m\) and their product is \(-b\). The sum does not depend on the discriminant, but the points are real only when \(D \ge 0\). The midpoint of the two intersections has \(x\)-coordinate \(m/2\).
