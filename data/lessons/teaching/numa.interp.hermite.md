# Hermite Interpolation

## Matching Values and Slopes

Hermite interpolation matches both the value and the derivative of a function at each node. Two nodes therefore supply four conditions, enough to determine a cubic polynomial uniquely.

The extra derivative information makes the interpolant follow the local direction of the data instead of merely passing through the points. That is why Hermite curves look smooth at joins in computer graphics and in motion planning.

The data are written as repeated nodes in a divided-difference table. Repeating a node lets the table encode a derivative in the same framework used for ordinary interpolation.

## Divided Differences with Repetition

The divided-difference table is built exactly as before, except that a repeated node uses the derivative in place of the usual quotient. Each repetition raises the order of the derivative that can be matched.

For two nodes with one derivative each, the table has four entries along the diagonal. The resulting polynomial has degree at most three and reproduces all four conditions.

The error term contains a product over the nodes, with each node appearing as many times as the number of conditions imposed there. This explains why derivative information improves accuracy near the nodes.

## Error and Examples

The interpolation error is proportional to the next derivative of the function and to the product of the distances to the nodes. Repeating a node effectively squares its distance factor, which usually reduces the error nearby.

A common example matches position and velocity at the start and end of an interval. The cubic that results is the same curve used for smooth camera paths and keyframe animation.

Hermite interpolation is also the building block of cubic splines. A spline adds continuity of second derivatives across the joins, which turns a collection of local Hermite pieces into a globally smooth curve.
