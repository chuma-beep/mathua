# Confidence interval for a mean with t

When the population standard deviation is unknown, a confidence interval for a mean uses the sample standard deviation and a t critical value. The t distribution has heavier tails than the normal distribution, so the interval is wider. That extra width is the price of estimating the spread from the same sample that supplies the mean.

## Why t replaces z

The standard error of the sample mean is estimated by

\[
SE = \frac{s}{\sqrt{n}},
\]

where \(s\) is the sample standard deviation. If the population standard deviation \(\sigma\) were known, the standardized mean would be normal. With \(s\) in its place the ratio

\[
t = \frac{\bar{x} - \mu}{s/\sqrt{n}}
\]

has a t distribution with \(n-1\) degrees of freedom, provided the underlying population is roughly normal or the sample is large enough for the central limit theorem to help.

The degrees of freedom are \(n-1\), not \(n\), because the sample standard deviation is computed around the sample mean and one degree of freedom is spent on that mean. As the degrees of freedom grow the t distribution approaches the standard normal, and the critical value falls toward 1.96 for 95 percent confidence.

## The margin of error

The margin of error is

\[
E = t^* \frac{s}{\sqrt{n}},
\]

where \(t^*\) is the two-sided critical value for the chosen confidence level and \(n-1\) degrees of freedom. For 95 percent confidence and 8 degrees of freedom the critical value is 2.306. A sample with standard deviation 10 and size 9 therefore has standard error \(10/3\) and margin \(2.306 \times 10/3 \approx 7.687\).

The critical value must match the degrees of freedom. Using the normal value 1.96 with a small sample understates the margin and produces an interval that is too narrow. That error is common because software often reports both values and the labels are easy to mix up.

The margin is a half-width, not a full width. The interval extends the same distance above and below the sample mean, so its total width is twice the margin.

## Endpoints of the interval

The interval is

\[
\left(\bar{x} - t^* \frac{s}{\sqrt{n}}, \; \bar{x} + t^* \frac{s}{\sqrt{n}}\right).
\]

If the sample mean is 52 and the margin is 7.687, the interval runs from 44.313 to 59.687. The interpretation is about the procedure: if many samples were drawn and the same calculation were repeated, about 95 percent of the resulting intervals would contain the true mean. It is not a statement that the true mean has a 95 percent probability of lying in this particular interval.

The interval is wider when the sample is small, when the standard deviation is large or when the confidence level is higher. Those three effects all enter through the same product, which is why the margin is the quantity to inspect first.
