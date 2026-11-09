# Confidence interval for a proportion

A proportion is a mean in disguise. Code each success as one and each failure as zero, and the sample proportion is the average of those zeros and ones. The confidence interval follows the same pattern as the interval for a mean, with a standard error built from the proportion itself.

## The standard error of a proportion

For a sample of size \(n\) with sample proportion \(\hat{p}\), the estimated standard error is

\[
SE(\hat{p}) = \sqrt{\frac{\hat{p}(1-\hat{p})}{n}}.
\]

The numerator is the Bernoulli variance and the denominator shrinks it by the sample size. The largest standard error for a fixed \(n\) occurs at \(\hat{p} = 0.5\), where the product \(\hat{p}(1-\hat{p})\) reaches one quarter. That is why polls near a tied race need larger samples for the same margin.

The standard error has the same units as a proportion, so it lies between zero and one. A standard error of 0.05 says the sample proportion typically differs from the true proportion by about five percentage points. The approximation is based on the normal distribution and works best when both \(n\hat{p}\) and \(n(1-\hat{p})\) are at least about ten.

## The margin of error

The margin of error is the critical value multiplied by the standard error:

\[
E = z^* \sqrt{\frac{\hat{p}(1-\hat{p})}{n}}.
\]

For 95 percent confidence the critical value is 1.96. For 90 percent it is 1.645, and for 99 percent it is 2.576. The critical value comes from the standard normal distribution because the sample proportion is approximately normal for large samples.

A poll of 400 respondents with \(\hat{p} = 0.5\) has standard error \(\sqrt{0.25/400} = 0.025\). The 95 percent margin is \(1.96 \times 0.025 = 0.049\), or 4.9 percentage points. Doubling the sample size would divide the standard error by the square root of two, which is why reducing a margin by half costs four times the sample.

The margin uses the estimated proportion, not a guessed population value. If the true proportion is unknown, the sample value is the only available estimate, and the interval is built around it.

## Reading the interval

The interval is

\[
\hat{p} \pm z^* \sqrt{\frac{\hat{p}(1-\hat{p})}{n}}.
\]

With the poll above, the interval runs from 0.451 to 0.549. The correct interpretation is again about the procedure: about 95 percent of intervals built this way from repeated samples would contain the true proportion. The particular interval either contains it or does not.

The interval is not a prediction about the next sample. A new poll of the same size has its own sampling error, and the spread of new poll results is wider than the confidence interval for the mean. That distinction between a confidence interval and a prediction interval is the same one that appears in regression.

When the interval includes 0.5, a claim that the population is evenly split is compatible with the data. When it lies entirely above or below 0.5, the data point to a side. The width of the interval says how strong that evidence is, so the endpoints are more informative than a yes or no conclusion.
