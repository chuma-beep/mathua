# Piecewise linear interpolation

Piecewise linear interpolation joins consecutive data points with straight segments. It is simple, local and easy to evaluate.

## Between two nodes

On an interval between \((x_0,y_0)\) and \((x_1,y_1)\), the interpolant is the straight line through the endpoints.

The midpoint value is the average of the endpoint values. This follows from linearity.

The formula is the same as point-slope form with the segment slope.

## Slope on a segment

The slope is \((y_1-y_0)/(x_1-x_0)\). It is constant on the segment and can change abruptly at a node.

A large slope change means a visible corner. Piecewise linear interpolation is continuous but not differentiable at the nodes.

The slope is a difference quotient, so it is exact for the linear segment.

## Finding the interval

To evaluate at a query point, first locate the interval that contains it. The left endpoint identifies the segment.

The interpolant is local: changing one data value affects only the two adjacent segments.

A binary search makes interval lookup fast when there are many nodes.
