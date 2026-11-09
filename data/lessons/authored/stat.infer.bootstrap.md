# Bootstrap confidence intervals

The bootstrap estimates the sampling distribution of a statistic by resampling the observed data.

## Resample a statistic

Draw a sample with replacement from the observed data. Compute the statistic on the resample. Repeat many times to build the bootstrap distribution.

## Build a percentile interval

Order the bootstrap statistics. Take the lower and upper percentiles. The interval between them is the percentile confidence interval.

## Estimate the standard error

The standard deviation of the bootstrap distribution estimates the standard error. The variance is the square of that standard deviation. More resamples reduce Monte Carlo noise but not the data noise.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
