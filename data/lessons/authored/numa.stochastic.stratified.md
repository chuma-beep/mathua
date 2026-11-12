# Stratified Monte Carlo

Stratified sampling splits the domain into strata and samples each one separately. Removing between-stratum variation from the error is the source of the gain.

## The estimator

Each stratum mean is weighted by its probability, and the weights sum to one. The estimator is unbiased when the strata partition the domain.

## Variance

The variance is the weighted sum of within-stratum variances divided by the stratum sample counts. Between-stratum variance does not enter.

## Allocation

Neyman allocation gives more samples to strata with larger variance times size. Proportional allocation is simpler and needs no pilot estimate.

To check your understanding, weight the stratum means, sum the within variances and apply Neyman allocation.

Stratification removes variation that the strata already explain.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
