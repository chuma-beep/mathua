# The bootstrap

The bootstrap estimates sampling variability by resampling the observed data. It replaces an analytic standard error with the spread of many resampled statistics.

## The procedure

Draw B samples of size n with replacement from the data, compute the statistic on each and use the spread of the replicates. The original sample stands in for the population.

## Standard error

The bootstrap standard error is the standard deviation of the replicates. It converges to the usual estimate as B grows.

## Limits

The method fails for extreme order statistics and for dependent data without a suitable resampling scheme. A block bootstrap preserves dependence.

To check your understanding, resample with replacement, take the replicate spread and recognise an order statistic.

The bootstrap treats the sample as a stand-in population and resamples it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
