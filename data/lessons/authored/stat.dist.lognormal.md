# Lognormal distribution

A positive random variable is lognormal when its logarithm is normally distributed. The name describes the transformation, not the shape of the density. Prices, incomes, particle sizes and waiting times that grow by percentages rather than by fixed amounts often fit this model, because a sum of many small multiplicative effects becomes normal on the log scale.

## The logarithm is normal

Write \(X\) for the positive variable and let

\[
\ln X \sim N(\mu, \sigma^2).
\]

The parameters \(\mu\) and \(\sigma\) describe the log scale, not the original scale. The median of \(X\) is \(e^{\mu}\), because the normal distribution is symmetric about \(\mu\) and the exponential is increasing. Half of the probability lies below \(e^{\mu}\) and half above it.

The mode is \(e^{\mu - \sigma^2}\), which is below the median whenever \(\sigma > 0\). The density is skewed to the right on the original scale even though it is symmetric on the log scale. That asymmetry is the visible signature of a lognormal variable.

A common error is to read \(\mu\) as the mean of \(X\). It is the mean of \(\ln X\). The mean of \(X\) is larger than \(e^{\mu}\), and the gap grows with \(\sigma\).

## Median and mean

The mean of a lognormal variable is

\[
E[X] = e^{\mu + \sigma^2/2}.
\]

The extra factor \(e^{\sigma^2/2}\) is the effect of the right skew. When \(\sigma\) is small the factor is close to one, and the mean is close to the median. When \(\sigma\) is large the mean can be far above the median, because a few very large values pull the average up.

For example, if \(\ln X\) has mean 1.2 and standard deviation 0.4, then the median is \(e^{1.2} \approx 3.3201\). The mean is \(e^{1.2 + 0.16/2} = e^{1.28} \approx 3.5966\). The mean exceeds the median, as it must for a right-skewed distribution.

The variance is \(e^{2\mu + \sigma^2}(e^{\sigma^2} - 1)\). It is not needed for every question, but it shows the same pattern: both the mean and the variance are controlled by the log-scale parameters through exponential factors.

## The mean-to-median ratio

Dividing the mean by the median cancels the \(e^{\mu}\) factor and leaves

\[
\frac{E[X]}{\text{median}(X)} = e^{\sigma^2/2}.
\]

This ratio depends only on the log-scale standard deviation. It is a compact measure of the skew: a ratio of 1 says the mean and median agree, and a ratio of 2 says the mean is twice the median.

Because the ratio is unitless, it can be compared across variables measured in different units. That is useful when two income distributions have different currencies but a similar shape. The ratio reports the shape without the scale.

The same factor appears in the mean formula, so the two results are one fact seen from two directions. Computing either one checks the other, and that consistency is worth using when a result looks surprising.
