# The Cramer-Rao bound

The Cramer-Rao bound limits the variance of an unbiased estimator by the reciprocal of the Fisher information.

## Statement

For an unbiased estimator, the variance is at least one over the Fisher information of the sample.

## Fisher information

Information is the expected squared score, the derivative of the log likelihood with respect to the parameter.

## Equality

Equality holds when the score is proportional to the estimation error, which happens for exponential families in natural parameters.

To check your understanding, assume unbiasedness, differentiate the log likelihood and relate score to error.

Information sets a floor on the variance of unbiased estimators.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
