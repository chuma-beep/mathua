# The paired t test

The paired t test compares two measurements taken on the same unit. Each pair produces one difference, and the test asks whether the mean difference is zero.

## Paired observations

Before-and-after measurements are paired. A reading from the same person under two treatments is paired. The pairing removes the variation between units, which often makes the test more sensitive.

The pairs must stay together. Treating paired observations as two independent samples discards the pairing and gives the wrong standard error.

## Difference scores

For each pair, subtract one measurement from the other. If the before value is $b$ and the after value is $a$, use $d = a - b$.

The sample mean of the differences is $\bar{d}$. A positive mean difference means the second measurement tends to be larger. A negative mean difference means the first tends to be larger.

## The paired t statistic

Let $s_d$ be the standard deviation of the differences and $n$ the number of pairs. The test statistic is

$$t = \frac{\bar{d}}{s_d / \sqrt{n}}.$$

The degrees of freedom are $n - 1$. The calculation is the one-sample t test applied to the difference scores.
