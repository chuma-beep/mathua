# Kruskal-Wallis test

The Kruskal-Wallis test compares three or more independent groups without assuming that the data are normal. It works with ranks instead of the original values, so a single extreme observation cannot dominate the result. The null hypothesis is that all groups come from the same distribution, and the alternative is that at least one group tends to produce larger values than another.

## Ranks instead of values

The first step is to pool every observation from every group, sort the pooled sample and assign ranks from 1 to \(N\), where \(N\) is the total number of observations. Ties receive the average of the ranks they would have occupied. The original values are then replaced by their ranks, and all later arithmetic uses the ranks.

Ranking removes the units and the shape of the original distribution. A group with a few very large values does not receive a huge rank sum unless those values really are the largest, because ranks are bounded by the sample size. That robustness is the main reason to choose this test over a one-way ANOVA.

The cost is a loss of information. Two data sets with the same ordering have the same ranks even if the gaps between the values are very different. When the normal assumptions hold, the ANOVA uses that extra information and has more power. The Kruskal-Wallis test trades some power for protection against outliers and skew.

## The H statistic

The test statistic is

\[
H = \frac{12}{N(N+1)} \sum_{i=1}^{k} \frac{R_i^2}{n_i} - 3(N+1),
\]

where \(k\) is the number of groups, \(n_i\) is the size of group \(i\), \(R_i\) is the sum of the ranks in group \(i\) and \(N\) is the total sample size. The first term measures how unevenly the ranks are distributed across groups, and the subtraction centers the statistic under the null hypothesis.

A large \(H\) means the rank sums are far from what equal-sized groups would produce under a random assignment of ranks. A small \(H\) means the groups look alike once the values are ranked. The statistic is nonnegative in exact arithmetic, and a negative value signals a calculation error.

When there are ties, a correction divides \(H\) by a factor slightly below one. The uncorrected form is what the questions here use, and the data are chosen so that all observations are distinct. That choice keeps the arithmetic exact and the interpretation clean.

## Degrees of freedom and interpretation

Under the null hypothesis and with reasonably large groups, \(H\) is approximately chi-square with

\[
df = k - 1
\]

degrees of freedom. The degrees of freedom depend on the number of groups, not on the number of observations. Four groups give three degrees of freedom, and ten groups give nine.

A significant result says that at least one group differs from the others. It does not say which group, and it does not say how large the difference is. A post-hoc comparison is needed to locate the difference, and an effect-size measure is needed to describe its size.

The test is a rank-based analogue of one-way ANOVA, so the same design logic applies. The groups must be independent, and the observations within each group should be exchangeable under the null. If the data are paired or repeated, a different rank test is required.
