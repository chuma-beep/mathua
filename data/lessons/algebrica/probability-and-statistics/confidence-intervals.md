> Content sourced from [Algebrica](https://algebrica.org/confidence-intervals/) — CC BY-NC 4.0

## Point estimation

When we rely on a sample to learn something about an unknown population parameter, say the [mean](<../introduction-to-the-mean/>) \\( \mu \\), a natural first step is to use a point estimator. In practice, this means choosing a statistic \\( \hat{\theta} \\) computed from the data and hoping that it captures, in a reasonable way, the value of the parameter we are trying to estimate. A property of a good estimator is unbiasedness: on average, across all possible samples we could have drawn, the estimator should hit the correct target. Formally, this requirement is written as: \\[\mu_{\hat{\theta}} = E(\hat{\theta}) = \theta \\] meaning that the [sampling distribution](<../sampling-distributions/>) of \\( \hat{\theta} \\) is centered exactly at the true parameter.


Of course, being centered on the right value does not mean that every single estimate will be perfect. Each sample carries its own variability, and this inevitably propagates to the estimator. The degree of this fluctuation is described by the sampling [variance](<../variance/>) \\[\operatorname{Var}(\hat{\theta}) \\] which quantifies how much the estimator tends to move around from one sample to the next. For instance, when estimating a population mean, the sample mean \\( \bar{X} \\) has the properties  
\\[E(\bar{X})=\mu \quad \operatorname{Var}(\bar{X})=\frac{\sigma^{2}}{n} \\] so its accuracy improves naturally as the sample size \\( n \\) grows. A natural consequence of this variability is that any point estimate will generally differ from the true parameter by some amount. This difference is often referred to as the estimation error, and it can be written formally as  
\\[\delta = \hat{\theta} - \theta \\] Because different samples lead to different values of \\( \hat{\theta} \\), the estimation error is itself a random quantity. While we cannot eliminate this error \\(\delta\\), we can study its behavior, quantify how large it is likely to be, and design estimators whose variability is as small as possible.


Despite this, it remains rather unlikely that the specific value of \\( \hat{\theta} \\) obtained from a single sample will match the true parameter exactly. Every sample is just one of many possible snapshots of the population, each with its own quirks and randomness. This unavoidable sample-to-sample variation is precisely why statisticians often go beyond point estimates. To convey uncertainty in a more faithful and transparent way, it becomes essential to complement any point estimate with a range of plausible values, a principle that lies at the foundation of **confidence intervals**.

##### All these considerations take place within the framework of [sampling distributions](https://algebrica.org/sampling-distributions/). Different samples drawn from the same population can produce different results, and examining the distribution of a statistic across all such samples helps us understand its variability and the uncertainty involved when inferring population parameters from sample data.

## Interval estimation

Since a point estimate will almost never coincide exactly with the true value of a population parameter, it is often more meaningful to report a range of values within which the parameter is likely to fall. This leads to the idea of _interval estimation_. The intuition is simple: rather than committing to a single number, we use the information in the sample to build two bounds that frame the parameter. In symbolic form, we write  
\\[\hat{\theta}_{L} < \theta < \hat{\theta}_{U} \\] where \\( \hat{\theta}_{L} \\) and \\( \hat{\theta}_{U} \\) represent the lower and upper limits derived from the data. In this way, the estimate acknowledges the uncertainty inherent in sampling and provides a realistic range of plausible values for the population parameter, such as the mean \\( \mu \\). The interval can also be interpreted as providing bounds on the estimation error: the deviation \\( \delta = \hat{\theta} - \theta \\) is confined between the lower and upper limits, ensuring that the error cannot exceed the specified range.


To make this idea more precise, it helps to remember that the interval itself is built from the sample. This means that its endpoints, \\( \hat{\theta}_{L} \\) and \\( \hat{\theta}_{U} \\), change from one sample to another—just as any statistic does. As a result, the interval is not fixed but every time we draw a new sample, we would obtain a slightly different pair of bounds.

What we can consider is the probability that the interval constructed from the sample will contain the true parameter \\( \theta \\). This probability is called the confidence level, and it is written as: \\[P\\!\left( \hat{\theta}_{L} < \theta < \hat{\theta}_{U} \right) = 1 - \alpha \\] In simpler terms, if we were to repeat the sampling process many times, roughly a fraction \\( 1 - \alpha \\) of the intervals we compute would succeed in capturing the actual value of the parameter. In summary:

  * The interval \\(\hat{\theta}_{L} < \theta < \hat{\theta}_{U}\\) defines the \\(100(1-\alpha)%\\) confidence interval.
  * The quantity \\(1 - \alpha\\) is the confidence level.
  * The bounds \\( \hat{\theta}_{L} \\) and \\( \hat{\theta}_{U} \\) are the lower and upper confidence limits.


In practice, the value of \\( \alpha \\) is often chosen so that the confidence level \\( 1 - \alpha \\) falls between 95% and 99%. These levels strike a balance between having an interval that is sufficiently narrow to be informative and sufficiently wide to capture the true parameter with high probability. As the confidence level increases, the interval necessarily widens, reflecting the fact that a broader range gives a higher chance of containing the unknown parameter.

## Interval estimation of the mean with known \\(\sigma\\)

Consider the problem of estimating the population mean by constructing an interval from a single sample drawn from a normally distributed population with known \\(sigma\\). In this setting, a confidence interval for the mean \\( \mu \\) can be obtained by using the [Central Limit Theorem](<../normal-distribution/>) and the fact that the standardized sample mean follows the [standard normal distribution](<../normal-distribution/>). These results describe the behavior of the sampling distribution of \\( \bar{X} \\) and make it possible to quantify the uncertainty associated with the estimate.

To construct the confidence interval, we introduce the standardized variable  
\\[Z = \frac{\bar{X} - \mu}{\sigma / \sqrt{n}} \\]  
which follows the standard normal distribution when the population variance \\( \sigma^{2} \\) is known. Using this variable, we build an interval that captures the central portion of the distribution:

\\[P\\!\left( -z_{\alpha/2} < Z < z_{\alpha/2} \right) = 1 - \alpha \\]

where \\( z_{\alpha/2} \\) is the critical value such that the area in the two tails of the standard normal curve sums to \\( \alpha \\), and each tail individually accounts for \\( \alpha/2 \\). In other words, \\( z_{\alpha/2} \\) is chosen so that the area to its right under the standard normal density is exactly \\( \alpha/2 \\).

![Confidence intervals.](/diagrams/algebrica/confidence-intervals.png)

To solve for the population mean \\( \mu \\), we substitute the expression for \\( Z \\) into the probability statement above, obtaining:

\\[P\\!\left( -z_{\alpha/2} < \frac{\bar{X} - \mu}{\sigma / \sqrt{n}} < z_{\alpha/2} \right) = 1 - \alpha \\]

By isolating the parameter \\( \mu \\) within the [inequality](<../linear-inequalities/>), we obtain an expression that directly describes the range of plausible values for the population mean. Solving for \\( \mu \\) yields:

\\[P!\left( \bar{X} - z_{\alpha/2},\frac{\sigma}{\sqrt{n}} \;<\; \mu \;<\; \bar{X} + z_{\alpha/2}\,\frac{\sigma}{\sqrt{n}} \right) = 1 - \alpha \\]

which represents the \\( 100(1-\alpha)% \\) confidence interval for the mean.


A confidence interval with level \\(100(1-\alpha)%\\) offers a quantitative measure of the reliability of a point estimate. If the true mean \\( \mu \\) were positioned exactly at the midpoint of the interval, then \\( \bar{X} \\) would coincide with \\( \mu \\) and no estimation error would arise. In general, though, the sample mean will not match the population mean perfectly, and the resulting estimate inevitably deviates from the true value.

This deviation can be described through the absolute difference \\( |\bar{X} - \mu| \\) and with probability \\( 1-\alpha \\) the error is bounded above by: \\[z_{\alpha/2}\,\frac{\sigma}{\sqrt{n}} \\]

A common misunderstanding is to think that a 95% confidence interval means there is a 95% chance that \\( \mu \\) is inside this interval. This is not the right interpretation. The true value \\( \mu \\) is fixed, it does not move and it does not come with a probability attached to it. What actually changes from sample to sample is the interval itself.

If we were to collect many samples and build a new interval each time, about 95% of those intervals would include the true mean. So the 95% refers to how well the method works in the long run, not to the chance that \\( \mu \\)lies in the specific interval we obtained from our single sample.

## Sample size

One practical question that often arises when working with interval estimation is how large the sample needs to be in order to keep the estimation error below a chosen threshold \\( \delta \\). In other words, we want a sample size that allows us to be confident at the \\(100(1-\alpha)%\\) level that the difference between the sample mean and the true population mean will not exceed \\( \delta \\).

This formula applies specifically to the case in which the population standard deviation \\( \sigma \\) is known. Under this assumption, the minimum sample size required is:

\\[n = \left( \frac{z_{\alpha/2}\,\sigma}{\delta} \right)^{2} \\]

## Example

A company tests the battery life of a new model of smartphone. From a sample of 25 phones, the average battery duration recorded is \\(\bar{x} = 11.8\\) hours. Assume that previous studies indicate a known population standard deviation of \\(\sigma = 1.5\\) hours.

  * Compute the 95% confidence interval for the true mean battery life \\( \mu \\).
  * Determine how large the sample must be to ensure, with 95% confidence, that the estimation error does not exceed \\( \delta = 0.20 \\) hours.


We are dealing with a standard example in which the task is to determine a confidence interval for a sample drawn from a population with known variance. Since the population standard deviation \\( \sigma \\) is known, we use the standard normal distribution. For a 95% confidence level, the critical value is:

\\[z_{0.025} = 1.96 \\]

###### The value \\(1.96\\)is obtained from the [standard normal Z Table](<../standard-normal-z-table/>).


To construct the interval, we first compute the standard error of the sample mean, which measures how much \\( \bar{X} \\) is expected to vary from sample to sample. Using the known population standard deviation, we obtain: \\[\frac{\sigma}{\sqrt{n}} = \frac{1.5}{\sqrt{25}} = \frac{1.5}{5} = 0.30 \\]

Next, we determine the margin of error by multiplying the standard error by the critical value \\( z_{0.025} = 1.96 \\): \\[1.96 \times 0.30 = 0.588 \\]

This gives the half-width of the confidence interval. Therefore, the 95% confidence interval for the population mean \\( \mu \\) is \\(11.8 \pm 0.588\\) which corresponds to the range \\((11.212,\; 12.388)\\).


Using the standard formula for the sample size when the population standard deviation is known, we write: \\[n = \left( \frac{z_{\alpha/2},\sigma}{\delta} \right)^{2} \\]

We now substitute the numerical values into the expression: \\[n = \left( \frac{1.96 \times 1.5}{0.20} \right)^{2} = \left( \frac{2.94}{0.20} \right)^{2} = 216.09 \\]

Since the sample size must be an [integer](<../integers/>) and we always round upward to preserve the desired confidence level, the required sample size is \\( n = 217 \\).

In summary, the two results we obtain are:

  * 95% confidence interval for the mean \\((11.212,\; 12.388)\\).
  * Required sample size for an error of \\( \delta = 0.20 \\) is \\(n = 217\\).
