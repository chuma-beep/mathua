# Gradient noise

The gradient noise in stochastic descent comes from sampling and is usually assumed bounded or sub-Gaussian.

## Variance

The variance of a one-sample estimate is the variance of the per-sample gradient. A batch divides it by the batch size.

## Bias

Sampling must be uniform or otherwise unbiased. A biased sampler turns the method into descent on a different objective.

## Heavy tails

Bounded variance is weaker than a sub-Gaussian assumption. Heavy-tailed noise needs clipping or robust averaging for the same guarantees.

To check your understanding, compute a batch variance, compute an expectation and clip an outlier.

Noise is a modelling assumption, and the guarantee depends on which one holds.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
