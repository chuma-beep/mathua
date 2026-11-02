# Inference for the regression slope

The fitted slope is an estimate from one sample. Inference for the slope asks whether the population slope could be zero, which would mean no straight-line relationship.

## The slope as an estimate

The null hypothesis is usually $H_0: \beta_1 = 0$. The alternative says the population slope is positive, negative or nonzero.

A nonzero sample slope is not enough on its own. Random samples from a population with zero slope can still produce a sloping least-squares line.

## Standard error of the slope

The standard error of the slope measures how much the estimate would vary across samples. It grows with the residual spread and shrinks as the values of the explanatory variable spread out.

More observations and a wider range of $x$ values both make the slope more precise. Extrapolating far beyond the observed range remains risky even when the slope is precise.

## Testing the slope

The test statistic is

$$t = \frac{b_1}{SE(b_1)},$$

with $n - 2$ degrees of freedom for a simple linear regression. For a simple regression, testing a zero slope gives the same p-value as testing a zero correlation.

A confidence interval for the slope combines the estimate with its standard error. If the interval excludes zero, the data are inconsistent with a zero population slope at that confidence level.
