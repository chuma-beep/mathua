# Weighted mean

A weighted mean gives each observation a share of the total that matches its importance. A final grade is the standard example: a project worth one tenth of the course should not count as much as a final examination worth half. The arithmetic is a sum of products divided by a sum of weights, so the result is still an average but one whose centre is pulled toward the observations that carry more weight.

## The weighted average

For values \(x_1, x_2, \dots, x_n\) with weights \(w_1, w_2, \dots, w_n\), the weighted mean is

\[
\bar{x}_w = \frac{\sum_{i=1}^{n} w_i x_i}{\sum_{i=1}^{n} w_i}.
\]

The numerator multiplies every value by its weight, and the denominator adds the weights. If every weight equals one, the formula collapses to the ordinary mean. If one weight is twice another, the first value contributes twice as much to the numerator and the denominator accounts for that extra contribution.

A weight is a share, not a count. Doubling every weight leaves the weighted mean unchanged because both the numerator and the denominator double. That property is why the weights can be written as percentages, as points or as any other positive scale.

## Solving for a missing weight or value

The formula has four kinds of quantity: the values, the weights, the total weight and the weighted mean. When all but one value is known, the missing value follows from a linear equation. Multiply the target weighted mean by the total weight, subtract the weighted contributions of the known values, and divide by the weight attached to the unknown.

For example, suppose a weighted mean of 85 is required across components with weights 2, 3 and 5. If the first two values are 78 and 85, the known weighted total is \(78 \times 2 + 85 \times 3 = 411\). The target total is \(85 \times 10 = 850\), so the last component must supply \(850 - 411 = 439\). Its weight is 5, so the missing value is \(439/5 = 87.8\).

The same rearrangement works for a missing weight, but then the unknown appears in both the numerator and the denominator. That equation is still linear after multiplying both sides by the total weight, which is why a single missing weight can also be recovered exactly.

## Weighted and unweighted means compared

The unweighted mean treats every observation as equally important, while the weighted mean follows the weights. The difference between them measures how much the weighting has moved the centre of the data.

If the larger values carry the larger weights, the weighted mean lies above the unweighted mean. If the larger values carry the smaller weights, the weighted mean lies below it. The two means agree when the weights are equal or when the weights happen to be uncorrelated with the values.

That comparison is useful in practice. A weighted mean above the unweighted mean says the assessment scheme rewards the stronger performances. A weighted mean below it says the scheme rewards the weaker ones, which is worth noticing before a grade is published.
