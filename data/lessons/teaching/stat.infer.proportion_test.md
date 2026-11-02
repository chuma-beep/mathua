# The one-proportion z test

The one-proportion z test compares a sample proportion with a claimed population proportion. It is the proportion analogue of the one-sample z test for a mean.

## A single population proportion

The null hypothesis names a value $p_0$. The sample proportion is

$$\hat{p} = \frac{x}{n},$$

where $x$ is the number of successes and $n$ is the sample size. For a two-sided test the alternative says $p \ne p_0$.

## Standard error under the null

The standard error uses the null proportion, not the sample proportion:

$$SE = \sqrt{\frac{p_0(1-p_0)}{n}}.$$

Using the null value keeps the test honest: the hypothesised model determines the variability expected under that model.

## The one-proportion z statistic

The test statistic is

$$z = \frac{\hat{p} - p_0}{SE}.$$

The normal approximation is reasonable when $n p_0$ and $n(1-p_0)$ are both at least about 10. Otherwise an exact binomial test is safer.

## Reading the result

A large positive $z$ means the sample proportion sits far above the claimed value in the direction the null model measures. A large negative $z$ means it sits far below. The p-value turns that distance into a probability by asking how often a $z$ at least that extreme would appear if $p_0$ were true.

The test assumes independent observations and a sample that is small relative to the population. A survey that samples the same household twice, or that draws from a list with duplicate entries, breaks independence and makes the standard error too small, so the test overstates its evidence. Report the sample proportion, the null value, the standard error and the p-value together, because the statistic alone does not say how large the difference is in practical terms.
