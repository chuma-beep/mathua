> Content sourced from [Algebrica](https://algebrica.org/sampling-distributions/) — CC BY-NC 4.0

## From Populations to samples

A **sampling distribution** represents the distribution of a statistic obtained from all possible samples of a given size drawn from a population. In statistics, some problems, due to the size or complexity of the observable data, do not allow for a direct analysis of every element in a population.

For this reason, it becomes necessary to select a sample, that is, a smaller set of observations drawn from the population. A sample is a representative subset of the population, chosen randomly to minimize potential distortions or bias arising from other selection criteria. The process of selecting a sample is called sampling.

Inferential statistics uses samples to infer, from their observed characteristics, the properties or parameters of the entire population. A statistic is the value of a variable computed from the data of a sample. Examples of statistics include the sample mean, the sample variance, and the sample proportion. Sampling distributions describe how these values vary from one sample to another when the sampling process is repeated from the same population.

## Mean, mode, and median of a sampling distribution

Also for sampling distributions, it is possible to define the mean, mode, and median, each describing in a different way the central tendency of the distribution. The sample mean differs from the [simple arithmetic mean](<../introduction-to-the-mean/>) in that it represents the average of the observed values in a sample drawn from a population and is used as an estimate of the population mean. Unlike the arithmetic mean, which describes a fixed set of known data, the sample mean plays an inferential role, as it varies from one sample to another and follows its own sampling distribution.


In formal terms, the sample mean is defined as

\\[\overline{X} = \frac{1}{n} \sum_{i=1}^{n} X_i \\]

  * \\( \overline{X} \\) represents the sample mean.
  * \\( n \\) is the sample size.
  * \\( X_i \\) denotes the value of the \\( i \\)-th observation in the sample.


For example, consider a sample \\( X \\) consisting of \\(5\\) observed values:

\\( i \\)| 1| 2| 3| 4| 5  
---|---|---|---|---|---  
\\( X_i \\)| 4| 6| 5| 7| 8  
  
The sample mean is computed as

\\[\overline{X} = \frac{1}{n} \sum_{i=1}^{n} X_i = \frac{1}{5}(4 + 6 + 5 + 7 + 8) \\]

\\[\overline{X} = \frac{30}{5} = 6 \\]

Therefore, the sample mean is \\( \overline{X} = 6 \\).

##### It is worth noting that this mean is representative only of the specific sample considered and not of all the possible samples that could be drawn from the population, nor of the population itself.


The sample mode, on the other hand, is the value that occurs most frequently within the sample. It represents the most common observation and provides an indication of where the sample values tend to cluster. Formally, the sample mode can be expressed as

\\[\mathrm{Mode}(X) = x_k \, : \, f(x_k) = \max_{x_i} f(x_i) \\]

  * \\( X \\) denotes the random variable representing the sample data.
  * \\( f(x_i) \\) denotes the frequency of the value \\( x_i \\) in the sample
  * \\( x_k \\) is the observation that occurs most frequently.


Consider a sample \\( X \\) consisting of 7 observed values:

\\( i \\)| 1| 2| 3| 4| 5| 6| 7  
---|---|---|---|---|---|---|---  
\\( X_i \\)| 4| 6| 5| 6| 8| 6| 7  
  
The sample mode is the value that occurs most frequently within the sample. In this case, the value \\( 6 \\) appears three times, more than any other observation. Therefore, the sample mode is \\( 6 \\).

It is important to note that the mode is sensitive to the frequency of individual observations and may not be unique — multiple modes can exist if two or more values occur with the same highest frequency. In such cases, the distribution is referred to as multimodal.


The sample median, on the other hand, is defined as the middle value of the ordered sample.  
It divides the data set into two equal parts, with 50% of the observations below and 50% above this value. Formally, the sample median can be expressed as

\\[\tilde{x} = \begin{cases} x_{\frac{n+1}{2}}, & \text{if } n \text{ is odd} \\\\[6pt] \frac{1}{2}\left(x_{\frac{n}{2}} + x_{\frac{n}{2}+1}\right), & \text{if } n \text{ is even} \end{cases} \\]

  * \\( x_i \\) are the ordered observations of the sample.
  * \\( n \\) is the sample size.


Consider a sample \\( X \\) consisting of 6 observed values:

\\( i \\)| 1| 2| 3| 4| 5| 6  
---|---|---|---|---|---|---  
\\( X_i \\)| 4| 5| 6| 8| 9| 10  
  
To find the sample median, the observations must first be ordered from smallest to largest (which they already are in this case). When the sample size \\( n \\) is even, the median is calculated as the average of the two central values:

\\[\begin{align} \tilde{x} = \frac{1}{2}\left(x_{\frac{n}{2}} + x_{\frac{n}{2}+1}\right) \end{align} \\]

Substituting the corresponding values we have:

\\[\begin{align} \tilde{x} = \frac{1}{2}(x_3 + x_4) = \frac{1}{2}(6 + 8) = 7 \end{align} \\]

Therefore, the sample median is \\( \tilde{x} = 7 \\).

##### The median is less sensitive to extreme values (outliers) compared to the mean, making it a more robust measure of central tendency in samples that contain skewed or non-symmetric data.

## Sample variance and sample standard deviation

The sample variance measures how the observations in a sample are distributed with respect to the sample mean. Larger values of variance indicate that the observations are spread out more widely around the mean, while smaller values indicate that they are more tightly clustered.  
Formally, given \\( X_1, X_2, \dots, X_n \\) [random variables](<../discrete-random-variables/>), the sample variance is defined as

\\[\begin{align} S^2 = \frac{1}{n - 1} \sum_{i=1}^{n} (X_i - \overline{X})^2 \end{align} \\]

  * \\( S^2 \\) denotes the sample variance
  * \\( \overline{X} \\) is the sample mean
  * \\( n \\) is the sample size.


The sample variance differs from the [simple variance](<../variance/>) because it divides by \\( n - 1 \\) instead of \\( n \\), where \\( n - 1 \\) represents the degrees of freedom. This adjustment accounts for the fact that the sample mean, being used in the calculation, constrains one observation and leaves only \\( n - 1 \\) values free to vary. Dividing by \\( n - 1 \\) compensates for this constraint and ensures that the sample variance is an unbiased estimator of the population variance.


The sample variance can also be written in an equivalent computational form that avoids explicitly calculating the sample mean. This alternative expression is obtained by [expanding the squared differences](<../notable-products/>) and simplifying the terms of the variance formula:

\\[S^2 = \frac{1}{n(n - 1)} \left[ n \sum_{i=1}^{n} X_i^2 - \left( \sum_{i=1}^{n} X_i \right)^2 \right] \\]

This formulation is particularly useful for manual calculations or when working with large data sets, as it reduces the need for repeated subtraction of the mean and minimizes rounding errors.


The sample standard deviation is defined as the [square root](<../radicals/>) of the sample variance \\( S^2 \\). It provides a measure of dispersion expressed in the same units as the data, indicating how much, on average, the observations deviate from the sample mean.  
Formally, it is expressed as

\\[S = \sqrt{S^2} = \sqrt{\frac{1}{n - 1} \sum_{i=1}^{n} (X_i - \overline{X})^2} \\]

A smaller value of ( S ) indicates that the data points are closely clustered around the mean, whereas a larger value suggests greater variability within the sample.

## Sample range

The sample range is defined as the difference between the largest and smallest observed values within a sample. It provides a simple measure of dispersion that indicates the total spread of the data, though it is highly sensitive to extreme values. Formally, it can be expressed as

\\[R = X_{\text{max}} - X_{\text{min}} \\]

where \\( X_{\text{max}} \\) and \\( X_{\text{min}} \\) represent, respectively, the maximum and minimum observations in the sample.

## Connection with the normal distribution

As a consequence of the [Central Limit Theorem](<../normal-distribution/>), when the sample size \\( n \\) becomes large, the distribution of the sample mean \\( \overline{X} \\) approaches a normal distribution with mean \\( \mu \\) and variance \\( \frac{\sigma^2}{n} \\). This can be expressed through the standardized variable

\\[Z = \frac{\overline{X} - \mu}{\sigma / \sqrt{n}} \\]

which, for sufficiently large \\( n \\), follows approximately the [standard normal distribution](<../standard-normal-z-table/>)

\\[Z \sim \mathcal{N}(x; 0, 1) \\]

regardless of the shape of the population distribution. There is no exact value of \\( n \\) that guarantees normality, since the Central Limit Theorem describes an asymptotic behavior.  
However, in practice, the sample mean tends to be approximately normal when:

  * \\( n \ge 30 \\) for most population distributions with finite variance.
  * \\( n < 30 \\) if the population is already close to normal.
  * \\( n > 50 \\) or \\( n > 100 \\) if the population is highly skewed or contains outliers.


##### In general, the larger the sample size, the closer the sampling distribution of the mean is to the normal distribution.

## Example 1

To show in practice how sampling distributions are connected to the normal distribution, consider a company that manufactures industrial machines whose operating lifespan follows a normal distribution with a mean of 5,000 hours and a standard deviation of 200 hours.  
We want to calculate the probability that a random sample of 25 machines has an average lifespan of less than 4,950 hours.


According to the Central Limit Theorem, the sampling distribution of \\( \overline{X} \\) will be approximately normal, with mean

\\[\mu_{\overline{X}} = 5000 \\]

and standard deviation

\\[\sigma_{\overline{X}} = \frac{\sigma}{\sqrt{n}} = \frac{200}{\sqrt{25}} = 40 \\]

For the observed sample mean \\( \overline{X} = 4950 \\), we can apply the standardization formula to convert it into a corresponding \\( z \\)-score. The transformation from a sample mean to a standard normal variable is given by:

\\[Z = \frac{\overline{X} - \mu_{\overline{X}}}{\sigma_{\overline{X}}} \\]

Substituting the known values, we obtain

\\[z = \frac{4{,}950 - 5{,}000}{40} = \frac{-50}{40} = -1.25 \\]

This means that the sample mean of 4950 hours is 1.25 standard deviations below the expected population mean. We can therefore express the probability associated with the sample mean in terms of the standardized variable \\( Z \\). By substituting the corresponding \\( z \\)-score, the probability becomes

\\[P(\overline{X} < 4{,}950) = P(Z < -1.25) \\]

Using the [standard normal Z table](<../standard-normal-z-table/>), we find that

\\[P(Z < -1.25) = 0.1056 \\]

##### This means there is approximately a 10.56% probability that the sample of 25 machines will have an average lifespan shorter than 4,950 hours.

Probability and Statistics

Probability and statistics provide a framework for analyzing uncertainty, modeling random phenomena.

1.5k

[Introduction to the Mean](https://algebrica.org/introduction-to-the-mean/)

1.2k

[Arithmetic Mean](https://algebrica.org/arithmetic-mean/)

1.6k

[Geometric Mean](https://algebrica.org/geometric-mean/)

1.2k

[Harmonic Mean](https://algebrica.org/harmonic-mean/)

1k

[Root Mean Square](https://algebrica.org/root-mean-square/)

1k

[Median and Quantiles](https://algebrica.org/median/)

1.1k

[Variance](https://algebrica.org/variance/)

946

[Discrete Random Variables](https://algebrica.org/discrete-random-variables/)

1.2k

[Continuous Random Variables](https://algebrica.org/continuous-random-variables/)

1.3k

[Mean or Expected Value of a Random Variable](https://algebrica.org/mean-or-expected-value-of-a-random-variable/)

1.1k

[Variance and Covariance of a Random Variable](https://algebrica.org/variance-and-covariance-of-a-random-variable/)

1.3k

[Bernoulli Distribution](https://algebrica.org/bernoulli-distribution/)

1.4k

[Binomial Distribution](https://algebrica.org/binomial-distribution/)

1.5k

[Hypergeometric Distribution](https://algebrica.org/hypergeometric-distribution/)

1.8k

[Geometric Distribution](https://algebrica.org/geometric-distribution/)

1.4k

[Poisson Distribution](https://algebrica.org/poisson-distribution/)

1k

[Uniform Distribution](https://algebrica.org/uniform-distribution/)

1.7k

[Beta Distribution](https://algebrica.org/beta-distribution/)

6.4k

[Normal Distribution](https://algebrica.org/normal-distribution/)

1.3k

[Standard Normal Z Table](https://algebrica.org/standard-normal-z-table/)

1.3k

[Gamma Distribution](https://algebrica.org/gamma-distribution/)

1.4k

[Chi-square Distribution](https://algebrica.org/chi-square-distribution/)

1k

[Student’s t Distribution](https://algebrica.org/student-t-distribution/)

1.3k

[Exponential Distribution](https://algebrica.org/exponential-distribution/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
