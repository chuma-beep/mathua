# Subsampling

Subsampling estimates a sum over a large data set from a random part of it. The estimate is unbiased when the sample is uniform.

## Sampling rate

A rate of one in k means each term is included with probability one over k. The estimate scales the sample average by the inverse rate.

## Cost

The cost of the estimate is proportional to the sample size rather than the data size. That is the whole reason for subsampling.

## Bias

A uniform sample gives an unbiased estimate of the mean. A sample drawn from a subset of the data gives a biased estimate of the full sum.

To check your understanding, scale a sample average, cost a sample and test unbiasedness.

Subsampling is how a sum over millions of terms becomes a gradient in constant time.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
