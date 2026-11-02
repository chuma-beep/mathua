# The two-sample t test

The two-sample t test compares the means of two independent groups. It asks whether the observed difference is larger than the difference expected from sampling variation alone.

## Two independent samples

The null hypothesis says the two population means are equal. The alternative says they differ, or says one is larger than the other.

Independence matters. A measurement in the first group must not be paired with a measurement in the second group. Paired observations need the paired t test instead.

## The standard error of the difference

Each sample mean has its own standard error. For sample standard deviations $s_1$ and $s_2$ with sample sizes $n_1$ and $n_2$, the standard error of the difference is

$$SE = \sqrt{\frac{s_1^2}{n_1} + \frac{s_2^2}{n_2}}.$$

The variances add because the two samples are independent. Larger samples make the standard error smaller.

## The two-sample t statistic

The test statistic is

$$t = \frac{\bar{x}_1 - \bar{x}_2}{SE}.$$

The pooled degrees of freedom for the equal-variance form are $n_1 + n_2 - 2$. A large absolute t statistic is evidence against the null hypothesis.
