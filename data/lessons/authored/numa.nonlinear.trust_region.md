# Trust-region methods

A trust-region method minimises a local model inside a ball whose radius it controls. The radius grows when the model predicts well and shrinks when it does not.

## The ratio

The agreement ratio is the actual reduction divided by the predicted reduction. A ratio above 0.75 suggests a larger radius; below 0.25 suggests a smaller one.

## The subproblem

Minimising a quadratic inside a ball is a constrained problem solved by a dogleg or conjugate-gradient iteration. The solution lies on the boundary when the model is unbounded below.

## Robustness

Because every step is bounded, convergence can be proved from remote starting points. A line search needs more care to achieve the same guarantee.

To check your understanding, compute actual reduction, solve in the ball and bound every step.

The trust radius is a measured claim about how far the model can be believed.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
