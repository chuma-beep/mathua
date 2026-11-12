# Gradient aggregation

Aggregation combines several sample gradients before a step. A larger batch lowers variance, and the way the batch is formed decides whether the saving comes from averaging or from correlation.

## Batch averaging

The variance of a mean of independent samples falls with the batch size. Doubling the batch halves the variance, so the standard deviation falls by the square root of two.

## Correlated samples

When samples are correlated the variance falls more slowly than the batch size suggests. The effective sample size is the batch size divided by one plus the average correlation, which is the number that matters for the step.

## Cost accounting

Aggregation spends the same sample budget as several single steps but produces one direction. It is worthwhile when the step is expensive, for example when the update touches a large matrix.

To check your understanding, average independent samples, account for correlation and spend one budget.

Aggregation is the simplest variance reduction and the baseline against which snapshot and control-variate methods are measured.
