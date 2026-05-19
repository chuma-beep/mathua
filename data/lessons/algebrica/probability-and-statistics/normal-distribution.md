> Content sourced from [Algebrica](https://algebrica.org/normal-distribution/) — CC BY-NC 4.0

## Definition of the normal distribution

The **normal distribution** , also known as the Gaussian distribution, is one of the most important continuous probability distributions in both probability and statistics. It plays a central role in modeling real-world phenomena where values tend to cluster around a mean, following a characteristic bell-shaped curve. Mathematically, the normal distribution is defined by two parameters: the [mean](<../introduction-to-the-mean/>) \\( \mu \\), which determines the center of the distribution, and the [standard deviation](<../variance/>) \\( \sigma \\), which controls its spread. It is usually denoted as:

\\[\mathcal{N}(x; \mu, \sigma) \\]

As previously introduced, the normal distribution has a distinctive bell-shaped form and follows a set of well-defined mathematical properties that make it unique among continuous probability distributions.

![](https://algebrica.org/wp-content/uploads/resources/images/normal-distribution-1.png)

  * The total area under the curve equals \\(1\\). This means that the integral of its probability density function over the entire real line, from \\( -\infty \\) to \\( +\infty \\), is equal to \\(1\\).
  * The curve is symmetric around the mean \\( \mu \\). In other words, it looks the same on both sides of the mean, with half of the total probability lying to the left and the other half to the right.
  * The curve has two [inflection points](<../maximum-minimum-and-inflection-points/>), located at \\( x = \mu + \sigma \\) and \\( x = \mu - \sigma \\). At these points, the curvature of the graph [changes sign](<../sign-analysis-in-inequalities/>), marking the transition between the concave and convex regions of the distribution.
  * The curve is asymptotic to the x-axis for values of \\( x \\) that move farther away from the mean.


## Key features

  * \\[\text{1. } \quad \mathcal{N}(x; \mu, \sigma) = \frac{1}{\sigma \sqrt{2\pi}} \exp\\!\left( -\,\frac{(x - \mu)^{2}}{2\sigma^{2}} \right) \\]

  * \\[\text{2. } \quad \mu = E(X) = \mu \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = \sigma^{2} \\]

  * \\[\text{4. } \quad \sigma = \sigma \\]


##### Each expression highlights a key property of the normal distribution, showing how its bell shape is fully determined by the mean \\(\mu\\) and the standard deviation \\(\sigma\\), which control its center and spread.

## Probability density function of the normal distribution

The random variable \\( X \\) that follows a normal distribution is known as a normal random variable. It represents a [continuous random variable](<../continuous-random-variables/>) whose probabilities are described by the normal probability density [function](<../functions>), defined as:

\\[\mathcal{N}(x; \mu, \sigma) = \frac{1}{\sqrt{2\pi} \,\sigma} \, e^{-\frac{1}{2\sigma^2}(x - \mu)^2} \\]

\\[-\infty < x < +\infty\\]

This function describes how probability is distributed over the possible values of the continuous random variable \\( X \\), depending on the mean \\( \mu \\) and the standard deviation \\( \sigma \\). To explore in greater depth the concepts of the mean (or expected value), variance, and standard deviation of a random variable ([discrete](<../discrete-random-variables/>) or [continuous](<continuous-random-variables/>)), refer to the following related topics:

  * 1.3k [Mean or Expected Value of a Random Variable](https://algebrica.org/mean-or-expected-value-of-a-random-variable/)
  * 1.1k [Variance and Covariance of a Random Variable](https://algebrica.org/variance-and-covariance-of-a-random-variable/)


As discussed above, the integral of the normal density function over the entire real line, from \\( -\infty \\) to \\( +\infty \\), is equal to 1:

\\[\frac{1}{\sqrt{2\pi}\,\sigma} \, \int_{-\infty}^{+\infty} e^{-\frac{1}{2\sigma^2}(x - \mu)^2} \, dx = 1 \\]

From the integral it follows that, if we want to compute the area under the curve between two points \\( x_0 \\) and \\( x_1 \\), we must evaluate the [definite integral](<../definite-integrals/>) of the normal density function within that [interval](<../intervals/>).

![](https://algebrica.org/wp-content/uploads/resources/images/normal-distribution-2.png)

It follows intuitively that the shaded area under the curve represents the probability that the random variable \\( X \\) assumes a value within the interval \\([x_0, x_1]\\). In other words, the integral of the normal density function over this range quantifies the likelihood of observing \\( X \\) between those two limits. Formally, the probability that the random variable \\( X \\) takes a value between \\( x_0 \\) and \\( x_1 \\) is given by:

\\[\begin{align} P(x_0 < X < x_1) &= \int_{x_0}^{x_1} n(x; \mu, \sigma)\,dx \\\\[6pt] &= \frac{1}{\sqrt{2\pi}\,\sigma} \int_{x_0}^{x_1} e^{-\frac{1}{2\sigma^2}(x - \mu)^2}\,dx \end{align} \\]

## Standard normal distribution

To make probability calculations easier and more general, the normal distribution is often rewritten in a standardized form. In this process, the original variable \\( X \\) is transformed into a new variable \\( Z \\), defined as:

\\[Z = \frac{X - \mu}{\sigma} \\]

This new variable \\( Z \\) follows what is called the **standard normal distribution** , a special case where the mean is \\( 0 \\) and the standard deviation is \\( 1 \\). By standardizing, we can work with a single universal curve and use [the standard normal Z table](<../standard-normal-z-table/>) to find probabilities, instead of computing the integral for each specific distribution. In practice, every normal distribution can be converted into the standard one, making comparisons and calculations much simpler.

##### The values reported in Z-tables represent the cumulative area under the standard normal curve to the left of a given \\( Z \\) value.


Let us reconsider the case of the probability defined over a generic interval \\([x_0, x_1]\\). It can be expressed as:

\\[\begin{align} P(x_0 < X < x_1) &= \frac{1}{\sqrt{2\pi}\,\sigma} \int_{x_0}^{x_1} e^{-\frac{1}{2\sigma^2}(x - \mu)^2}\,dx \end{align} \\]

If we transform the variable \\( X \\) into its standardized form, the interval \\( P(x_0 < X < x_1) \\) can be rewritten in terms of the standard variable \\( Z \\) as:

\\[P(x_0 < X < x_1) = P\left( \frac{x_0 - \mu}{\sigma} < Z < \frac{x_1 - \mu}{\sigma} \right) \\]

where the variable \\( X \\) has been replaced by its standardized form \\( Z \\), and the limits \\( x_0 \\) and \\( x_1 \\) have been replaced by their corresponding standardized values. This transformation allows us to express the probability in the standard normal framework, where \\( Z \\) follows a distribution with mean \\( 0 \\) and standard deviation \\( 1 \\). Starting from the general form of the probability over an interval \\([x_0, x_1]\\) we have:

\\[\begin{align} P(x_0 < X < x_1) &= \frac{1}{\sqrt{2\pi}\,\sigma} \int_{x_0}^{x_1} e^{-\frac{1}{2\sigma^2}(x - \mu)^2}\,dx\\\\[6pt] &=\frac{1}{\sqrt{2\pi}\,\sigma} \int_{z_0}^{z_1} e^{-\frac{1}{2}(z)^2}\,dx = P(z_0 < Z < z_1) \end{align} \\]

This formulation highlights how the process of standardization provides a direct link between any normal distribution and the standard normal curve, enabling probabilities to be determined through universal reference values of \\( Z \\).

## Three-sigma rule

In a normal distribution, probabilities are symmetrically arranged around the mean.  
There exists a fundamental relationship between these probabilities and their distance from the mean, known as the 68–95–99.7 rule or three-sigma rule, which describes how most of the probability mass is concentrated near the center of the distribution.

![](https://algebrica.org/wp-content/uploads/resources/images/normal-distribution-3.png)

  * Approximately 68% of all values fall within one standard deviation of the mean, with about 34.1% on each side.
  * Expanding the range to two standard deviations includes roughly 95% of the data
  * Expanding the range to three standard deviations cover about 99.7% of all possible values.


## Central limit theorem

Let \\( X_1, X_2, \dots, X_n \\) be a sequence of independent and identically distributed random variables, each having an expected value \\( E[X_i] = \mu \\) and a finite variance \\( \mathrm{Var}(X_i) = \sigma^2 > 0 \\). We define their sample mean as:

\\[\bar{X}_{n} = \frac{1}{n} \sum_{i=1}^{n} X_i \\]

As the number of observations \\( n \\) grows larger, the distribution of the standardized variable \\(Z\\)

\\[Z = \frac{\bar{X}_n - \mu}{\sigma / \sqrt{n}} \\]

gradually approaches the standard normal distribution in law, according to the relation:

\\[\frac{\bar{X}_n - \mu}{\sigma / \sqrt{n}} \xrightarrow{d} \mathcal{N}(x; 0,1) \quad \text{as } n \to \infty \\]

In other words, regardless of the original distribution of the random variables the distribution of their mean tends to become approximately normal as the sample size increases. This result explains why the normal distribution appears so frequently in statistics: it acts as a limiting model for the behavior of averages when the number of observations is sufficiently large.

##### The notation \\( \xrightarrow{d} \\) denotes “convergence in distribution”, meaning that the probability distribution of the standardized variable gradually approaches the standard normal distribution as \\( n \\) increases.

## Selected references

  * **L. Wasserman**. [All of Statistics (free PDF draft)](https://www.stat.cmu.edu/~larry/all-of-statistics/)

  * **J. Blitzstein, J. Hwang**. [Introduction to Probability](https://projects.iq.harvard.edu/stat110/home)

  * **MIT OpenCourseWare**. [Introduction to Probability and Statistics](https://ocw.mit.edu/courses/18-05-introduction-to-probability-and-statistics-spring-2014/)

  * **OpenStax**. [Introductory Statistics 2e](https://openstax.org/details/books/introductory-statistics-2e)

  * **Penn State (STAT 414)**. [Probability Theory – Normal Distribution](https://online.stat.psu.edu/stat414/)

  * **University of Berkeley**. [Statistics](https://www.stat.berkeley.edu/~stark/SticiGui/)

  * **University of Oxford**. [Statistics Lecture Notes – Normal Distribution](https://www.stats.ox.ac.uk/~reinert/time/notesht10short.pdf)

  * **NIST/SEMATECH**. [Engineering Statistics Handbook – Normal Distribution](https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm)


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

938

[Sampling Distributions](https://algebrica.org/sampling-distributions/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
