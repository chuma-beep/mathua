# Effect size

A p-value says whether an effect is detectable, but it does not say how large the effect is. A very small difference can be significant in a large sample, and a large difference can be missed in a small one. Effect size measures the magnitude on a scale that does not depend on the sample size, which is why it belongs beside every test result.

## Standardized mean difference

Cohen's d expresses the distance between two means in standard deviation units:

\[
d = \frac{\bar{x} - \mu}{s}
\]

for a one-sample comparison, and

\[
d = \frac{\bar{x}_1 - \bar{x}_2}{s_{\text{pooled}}}
\]

for two groups. The denominator makes the result unitless, so a d of 0.5 means the means differ by half a standard deviation whether the measurement is in seconds or in kilograms.

The sign reports direction. A positive d means the first mean is larger, and a negative d means it is smaller. The magnitude reports size, and the usual rough labels are small near 0.2, medium near 0.5 and large near 0.8. Those labels are conventions from social science, not mathematical thresholds.

A d of zero means the means coincide. It does not mean the groups are identical, because the distributions can still differ in spread or shape. Effect size summarizes one aspect of a difference, not every aspect.

## Recovering the mean difference

The formula can be rearranged to recover the original units:

\[
\bar{x} - \mu = d \, s.
\]

If a study reports d = 0.8 and a standard deviation of 5 points, the means differ by 4 points. That conversion is the first thing a reader should do, because the units are what a decision maker can act on. A standardized effect is comparable across studies, while the raw difference is what a policy or a treatment plan needs.

The reverse direction also works. Given a raw mean difference and a standard deviation, dividing produces the standardized effect. That operation is how a meta-analysis puts studies with different scales on one axis.

The conversion assumes the standard deviation is the right yardstick for the comparison. For a skewed variable the standard deviation may not describe the spread well, and a different standardized measure may be preferable.

## Two-sample effect size

For two independent groups the denominator is the pooled standard deviation:

\[
s_{\text{pooled}} = \sqrt{\frac{(n_1-1)s_1^2 + (n_2-1)s_2^2}{n_1+n_2-2}}.
\]

The pooled value weights each group's variance by its degrees of freedom. When the groups are the same size and have the same standard deviation, the pooled value is just that common standard deviation.

A two-sample d uses the difference in group means in the numerator and the pooled standard deviation in the denominator. The result answers the same question as the one-sample version: how far apart are the groups, measured in the spread of the observations themselves?

The pooled standard deviation should not be used when the groups have very different spreads or very different sizes. In that case a standardized difference based on one group, or a raw difference with both standard deviations reported, is more honest than a single pooled number.
