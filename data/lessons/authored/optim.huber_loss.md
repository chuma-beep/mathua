# The Huber loss

The Huber loss is quadratic near zero and linear in the tails. It keeps the smoothness of least squares while bounding the influence of an outlier.

## Quadratic near zero

Inside the threshold the loss is half the squared residual. Small errors are treated exactly as least squares treats them.

## Linear in the tail

Outside the threshold the loss grows linearly and its derivative is constant. One large residual can no longer dominate the gradient.

## The transition point

The threshold is where the two pieces meet with the same value and the same slope. Differentiability at that point is what makes the loss easy to optimise.

To check your understanding, compute the residual, check the residual size and set the threshold.

The loss is convex and twice differentiable except at the threshold, where the second derivative jumps. The jump is harmless for gradient methods.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
