# Quadratic regression

Quadratic regression fits a parabola \(y = ax^2 + bx + c\) to a set of points. The three coefficients are chosen so that the squared vertical distances from the points to the curve are as small as possible.

## The leading coefficient

The coefficient \(a\) controls the curvature and the direction of opening. A positive \(a\) opens upward, a negative \(a\) opens downward, and a larger magnitude makes the parabola narrower.

In a least-squares fit the coefficient is found by solving the normal equations. The equations come from setting the partial derivatives of the squared error to zero, which produces one linear equation for each coefficient.

When the data lie exactly on a parabola, the fit recovers the coefficients without error. Substituting any three points that are not collinear determines a unique parabola, and the remaining points must agree with it.

## The constant term

The constant term is the fitted value at \(x = 0\). If the data include the point \((0, y_0)\), the constant term is \(y_0\) in an exact fit.

The constant term does not describe the vertex or the curvature. It only fixes the height where the parabola crosses the y-axis, and shifting every y-value by a constant shifts \(c\) by the same amount.

For the points \((-1, 0)\), \((0, -3)\) and \((1, -4)\), the exact fit has \(c = -3\). The other coefficients then follow from the two remaining equations.

## A predicted value

Prediction substitutes a new input into the fitted equation. For a quadratic, the arithmetic is one square, one product and one addition.

With \(y = 2x^2 - 3x + 1\), the prediction at \(x = 4\) is \(2 \times 16 - 3 \times 4 + 1 = 21\). The same equation gives 1 at \(x = 0\) and 0 at \(x = 1\).

A prediction is only as good as the model. A parabola can fit a curved pattern well inside the data range and still mislead outside it, so extrapolation needs a reason beyond the coefficient values.
