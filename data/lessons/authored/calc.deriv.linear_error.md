# The error in linear approximation

A linear approximation leaves an error controlled by the second derivative. The bound is half the curvature times the square of the distance from the base point.

## The bound

If \(|f''|\le M\) and the step is \(d\), the error is at most \(Md^{2}/2\). The square is why a linear approximation is good near the base point and poor far from it.

## The curvature

Doubling the second-derivative bound doubles the error bound. The curvature enters linearly while the step enters quadratically.

## Scaling the step

Multiplying the step by \(k\) multiplies the error bound by \(k^{2}\). A step twice as large can be four times as wrong, which is the practical warning behind the formula.

To check your understanding, read the curvature bound, double the bound and scale the step.

The linear error is bounded by half the curvature times the squared step.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
