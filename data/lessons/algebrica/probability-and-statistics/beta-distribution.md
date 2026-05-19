> Content sourced from [Algebrica](https://algebrica.org/beta-distribution/) — CC BY-NC 4.0

## Introduction to the beta distribution

The **beta distribution** is a continuous probability distribution defined over the open [interval](<../intervals/>) \\( (0, 1) \\). It depends on two positive numbers, \\( \alpha \\) and \\( \beta \\), which determine how the curve bends and how its mass is distributed along the interval. Because it only takes values between \\(0\\) and \\(1\\), it is often used to describe random quantities that represent proportions, ratios, or probabilities, situations where the outcomes are naturally limited within these bounds. In formal terms, the beta distribution is defined by the following probability density function:

\\[B(x; \alpha, \beta) = \frac{x^{\alpha - 1}(1 - x)^{\beta - 1}}{B(\alpha, \beta)} \quad 0 < x < 1 \\]

where \\( B(\alpha, \beta) \\) is the beta function, related to the Gamma function by:

\\[B(\alpha, \beta) = \frac{\Gamma(\alpha)\Gamma(\beta)}{\Gamma(\alpha + \beta)} \\]

Therefore, the probability density function can also be written explicitly in terms of the Gamma function as:

\\[B(x; \alpha, \beta) = \frac{\Gamma(\alpha + \beta)}{\Gamma(\alpha)\Gamma(\beta)} \, x^{\alpha - 1}(1 - x)^{\beta - 1} \\]

with \\( B(x; \alpha, \beta) = 0 \\) for \\( x \notin (0, 1) \\). The gamma function \\( \Gamma(c) \\) itself is defined, for every \\( c \in \mathbb{R}^+, \\) by the following [integral](<../definite-integrals/>) representation:

\\[\Gamma(c) = \int_{0}^{+\infty} x^{c - 1} e^{-x} \, dx \\]

The gamma function can be regarded as a continuous extension of the [factorial](<../factorial/>), which is defined only for [natural numbers](<../natural-numbers>), to all positive real values.

## The shape of the beta distribution

The shape of the beta distribution depends on the values of its parameters \\( \alpha \\) and \\( \beta \\). Depending on their magnitude, the distribution can take various forms: unimodal, U-shaped, or [monotonic](<../increasing-and-decreasing-functions/>). The distribution reaches its mode at the point

\\[x_0 = \frac{\alpha - 1}{\alpha + \beta - 2} \\]

  * If \\( \alpha > 1 \\) and \\( \beta > 1 \\), the distribution has a mode at \\( x_0 \\), corresponding to a maximum point of the density function.
  * If \\( \alpha < 1 \\) and \\( \beta < 1 \\), the function has a [minimum](<../maximum-minimum-and-inflection-points/>) at the same point.
  * In all other parameter combinations, the distribution is monotonic.
  * When \\( \alpha = \beta \\), the distribution is [symmetric](<../even-and-odd-functions/>) with respect to the vertical line \\( x = x_0 = \tfrac{1}{2} \\).


The figure illustrates one of the possible shapes of the beta distribution when both parameters \\( \alpha \\) and \\( \beta \\) are less than 1 and equal to each other. In this configuration, the distribution takes on a characteristic U-shaped form, with the density approaching infinity near the boundaries of the interval \\( (0, 1) \\). It is possible to observe a minimum point located at \\( x = x_0 \\), corresponding to the lowest value of the probability density within the [domain](<../determining-the-domain-of-a-function/>).

![Typical U-shaped form of the Beta distribution with α < 1, β < 1, and α = β.
](https://algebrica.org/wp-content/uploads/resources/images/beta-distribution-1.png)

An interesting case occurs when the two parameters are equal, that is \\( \alpha = \beta \\), and both are greater than \\(1\\). In this situation, the beta distribution becomes symmetric with respect to the vertical line \\( x = \tfrac{1}{2} \\) and takes on a unimodal (that is, a single-peaked curve), bell-shaped form with a single central peak. As the values of \\( \alpha \\) and \\( \beta \\) increase, the curve becomes progressively narrower and increasingly similar to a [normal distribution](<../normal-distribution>) centered around \\( x = 0.5 \\).

![As α and β increase, the Beta distribution approaches a normal curve centered at 0.5.](https://algebrica.org/wp-content/uploads/resources/images/beta-distribution-2.png)

To be more precise, this is an [asymptotic](<../asymptotes/>) approximation that holds for large values of \\( \alpha \\) and \\( \beta \\). For sufficiently large parameters, the Beta distribution can be approximated by a normal distribution with:

\\[\mu = \frac{\alpha}{\alpha + \beta} \\] \\[\sigma^2 = \frac{\alpha \beta}{(\alpha + \beta)^2 (\alpha + \beta + 1)} \\]

In the symmetric case, where \\( \alpha = \beta = k \\) we have:

\\[\mu = \tfrac{1}{2} \quad \sigma^2 = \frac{1}{8(2k + 1)} \approx \frac{1}{16k} \\]

Therefore, as \\( k \to \infty \\):

\\[\mathrm{B}(x; k, k) \approx \mathcal{N}\\!\left(\tfrac{1}{2}, \tfrac{1}{8(k + 1)}\right) \\]

## Key features

  * \\[\text{1. } \quad f(x) = \frac{x^{\alpha - 1}(1 - x)^{\beta - 1}}{B(\alpha, \beta)} \quad 0 \le x \le 1 \\]

  * \\[\text{2. } \quad \mu = E(X) = \frac{\alpha}{\alpha + \beta} \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = \frac{\alpha \beta}{(\alpha + \beta)^{2}(\alpha + \beta + 1)} \\]

  * \\[\text{4. } \quad \sigma = \sqrt{\frac{\alpha \beta}{(\alpha + \beta)^{2}(\alpha + \beta + 1)}} \\]


##### Each expression highlights a key property of the Beta distribution, showing how its shape depends on the parameters \\(\alpha\\) and \\(\beta\\), and how its mean and variability reflect the balance between these two shape parameters.

## Mean of the beta distribution

The [mean](<../introduction-to-the-mean/>), or [expected value](<../mean-or-expected-value-of-a-random-variable/>), of a beta distribution represents the average value of a random variable defined on the interval \\( (0, 1) \\), depending on the shape parameters \\( \alpha \\) and \\( \beta \\). Formally, the mean is obtained from the general definition of the expected value:

\\[\mu = E(X) = \int_{0}^{1} x \, B(x; \alpha, \beta) \, dx \\]

Substituting the probability density function of the Beta distribution we have:

\\[E(X) = \int_{0}^{1} x \, \frac{x^{\alpha - 1} (1 - x)^{\beta - 1}}{B(\alpha, \beta)} \, dx \\]

which simplifies to:

\\[E(X) = \frac{1}{B(\alpha, \beta)} \int_{0}^{1} x^{\alpha} (1 - x)^{\beta - 1} \, dx \\]

Recognizing that the integral on the right-hand side is itself the Beta function \\( B(\alpha + 1, \beta) \\), we obtain:

\\[E(X) = \frac{B(\alpha + 1, \beta)}{B(\alpha, \beta)} \\]

Using the identity that relates the Beta and Gamma functions we obtain:

\\[B(\alpha, \beta) = \frac{\Gamma(\alpha)\Gamma(\beta)}{\Gamma(\alpha + \beta)} \\]

Therefore, the mean can be expressed as:

\\[E(X) = \frac{\alpha}{\alpha + \beta} \\]

##### Hence, the mean of the Beta distribution depends only on the two shape parameters and expresses the balance between them.

## Variance of the beta distribution

The [variance](<../variance-and-covariance-of-a-random-variable/>) of the beta distribution measures how much the random variable is expected to vary around its mean value. While the mean describes the central tendency of the distribution, the variance quantifies its spread — that is, how concentrated or dispersed the possible outcomes are within the interval \\( (0, 1) \\). Formally, the variance is defined as:

\\[\sigma^2 = \mathrm{Var}(X) = E(X^2) - [E(X)]^2 \\]

Starting from the probability density function of the beta distribution:

\\[B(x; \alpha, \beta) = \frac{x^{\alpha - 1} (1 - x)^{\beta - 1}}{B(\alpha, \beta)} \\]

the expression can be rewritten as:

\\[\begin{align} E(X^2) &= \int_{0}^{1} x^2 f(x; \alpha, \beta) \, dx \\\ &= \frac{1}{B(\alpha, \beta)} \int_{0}^{1} x^{\alpha + 1} (1 - x)^{\beta - 1} \, dx \\\\[12 pt] &= \frac{B(\alpha + 2, \beta)}{B(\alpha, \beta)} \end{align} \\]

Substituting this expression and the mean into the formula gives:

\\[\sigma^2 = \frac{B(\alpha + 2, \beta)}{B(\alpha, \beta)} - \left(\frac{\alpha}{\alpha + \beta}\right)^2 \\]

Using the relationship between the beta and gamma functions:

\\[B(\alpha, \beta) = \frac{\Gamma(\alpha)\Gamma(\beta)}{\Gamma(\alpha + \beta)} \\]

we obtain the simplified expression for the variance:

\\[\sigma^2 = \frac{\alpha \beta} {(\alpha + \beta)^2 (\alpha + \beta + 1)} \\]

##### When \\( \alpha \\) and \\( \beta \\) increase together, the variance decreases, causing the distribution to become more concentrated around its mean.

## Relationship between the beta and uniform distribution

The [uniform distribution](<../uniform-distribution/>) can be regarded as a special case of the beta distribution. When both parameters are equal to one, that is \\( \alpha = \beta = 1 \\), the probability density function of the beta distribution becomes constant over the interval \\( (0, 1) \\). In the general case, the continuous uniform distribution defined over an interval \\( (a, b) \\) is given by:

\\[f(x) = \begin{cases} \dfrac{1}{b - a} & a < x < b \\\\[10pt] 0 & \text{otherwise} \end{cases} \\]

For \\( a = 0 \\) and \\( b = 1 \\), this expression reduces to \\( f(x) = 1 \\), which corresponds exactly to the \\(B(1, 1)\\) distribution. In this case, the two parameters of the beta distribution take the values \\( \alpha = 1 \\) and \\( \beta = 1 \\), producing a constant probability density across the interval \\(0, 1\\).

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

938

[Sampling Distributions](https://algebrica.org/sampling-distributions/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
