# The sign test

The sign test compares paired observations by counting how many differences are positive, ignoring their sizes.

## Signs of differences

Subtract one measurement from the other for each pair and record only whether the difference is positive, negative or zero. For ten pairs with differences 3, minus 1, 4, 2, minus 2, 5, 1, 3, 2 and 6 there are 8 positive signs and 2 negative signs. The magnitudes 6 and 1 count the same, which is both the price and the protection of the method.

## The test statistic

Under the null hypothesis each non-zero difference is equally likely to be positive or negative, so the positive count is binomial with n equal to the number of non-zero differences and \(p=0.5\). The statistic is that count; the p-value is the probability of a split at least as uneven. With 8 positives out of 10 the one-sided p-value is \((45 + 10 + 1)/1024 = 56/1024\), about 0.055.

## A worked count

Zeros are dropped before counting, so 12 pairs with one zero leave 11 tests. If 9 of 11 are positive, the one-sided p-value is \((55 + 11 + 1)/2048 = 67/2048\), about 0.033, which is significant at the 5 percent level, while 8 of 11 gives \((165 + 55 + 11 + 1)/2048 = 232/2048\), about 0.113, which is not.

The sign test is the weakest of the paired tests and the most robust, so it is the right choice when the differences are ordinal or their distribution is unknown.
