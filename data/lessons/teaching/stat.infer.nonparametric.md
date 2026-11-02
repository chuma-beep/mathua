# Nonparametric Tests

**Nonparametric test:** A test that does not assume a particular distribution for the population, often working with signs or ranks instead of the raw values. The sign test and the Wilcoxon signed-rank test are the standard paired examples.

## When normality is not available

Small samples, ordinal data and heavy outliers all undermine tests that assume normality. A nonparametric test replaces the distributional assumption with a weaker one, usually only that the observations are independent and come from a continuous distribution under the null.

## Worked: the sign test

1. For paired data, record only whether each difference is positive or negative.
2. Under the null of no effect, each sign is equally likely, so the number of positive signs is binomial with $n$ trials and $p = 1/2$.
3. Compare the observed count with that binomial distribution; an extreme count is evidence against the null.

## Ranks and the cost of fewer assumptions

The Wilcoxon signed-rank test ranks the absolute differences before comparing the positive and negative sums, which recovers some information about magnitude. Suppose eight pairs give differences

$$2,\ -1,\ 5,\ 8,\ -3,\ 4,\ 6,\ -2.$$

Their absolute values are 2, 1, 5, 8, 3, 4, 6 and 2. Ranking from smallest to largest gives the ranks 2, 1, 5, 8, 3, 4, 7 and 6 in the original order; the two differences of size 2 share the tied ranks 2 and 3, so each receives 2.5. The positive differences then contribute

$$2.5 + 6 + 8 + 5 + 7 = 28.5,$$

and the negative differences contribute

$$1 + 4 + 2.5 = 7.5.$$

The test compares these two sums. A large imbalance is evidence that the median difference is not zero, while a small imbalance is consistent with no effect.

Dropping the magnitudes costs power when the parametric assumptions actually hold. The signed-rank test also assumes the differences are symmetric about their median, and ties need a correction. Choose the nonparametric route when the assumptions are doubtful, not merely because it sounds safer. With clearly normal paired data and no serious outliers, the paired $t$ test is usually the better choice.
