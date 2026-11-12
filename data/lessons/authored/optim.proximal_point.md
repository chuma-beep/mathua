# The proximal point step

A proximal step trades a little objective for a shorter move. It turns an unstable gradient step into a firmly contractive one.

## The proximal step

The update minimises the objective plus half the squared distance from the current point, scaled by the step. The penalty keeps the new point near the old one.

## The fixed point

A point is optimal exactly when it equals its own proximal image. This is why proximal iterations converge to minimisers rather than to stationary points of a penalty.

## Choosing the step

A larger step makes the method more like a projection and a smaller step makes it more like gradient descent. The contraction factor depends on that choice.

To check your understanding, add the quadratic to the distance penalty, apply the proximal map and apply the map repeatedly.

The proximal step is implicit, so it is often solved in closed form. A quadratic objective gives a simple scaling of the gradient.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
