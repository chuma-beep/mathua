> Content sourced from [Algebrica](https://algebrica.org/gamma-distribution/) — CC BY-NC 4.0

## Introduction to the gamma distribution

The gamma distribution is a continuous probability distribution defined on the positive half-line. It is used to model waiting times, event durations, and phenomena where independent contributions accumulate over time. It originates from the gamma function, defined as

\\[\Gamma(\alpha) = \int_{0}^{\infty} x^{\,\alpha - 1} e^{-x}\, dx \\]

which extends the [factorial](<../factorial/>) to the real [domain](<../determining-the-domain-of-a-function/>) through the identity \\( \Gamma(n) = (n - 1)! \\) for every positive [integer](<../integers/>) \\( n \\). In general, a [continuous random variable](<../continuous-random-variables/>) \\( X \\) is said to follow a gamma distribution with parameters \\( \alpha \\) and \\( \beta \\) when its probability density function is given by

\\[G(x;\alpha,\beta)= \begin{cases} \dfrac{1}{\beta^{\alpha}\,\Gamma(\alpha)}\, x^{\alpha - 1}\, e^{-x/\beta} & x>0\\\\[10pt] 0 & x\leq 0 \end{cases} \\]

  * \\( \alpha \\) is the shape parameter: it controls how quickly the density rises near the origin and determines the degree of skewness. Larger values make the distribution more symmetric and shift the peak to the right.
  * \\( \beta \\) is the scale parameter: it stretches the distribution horizontally. Increasing \\( \beta \\) produces longer average waiting times and a broader spread.


The support of the distribution is the positive half-line, reflecting the fact that it models durations, waiting times, or other quantities that cannot take negative values. The interaction between \\( \alpha \\) and \\( \beta \\) shapes the overall behavior: \\( \alpha \\) governs the internal structure, while \\( \beta \\) sets the scale.

![Plot of the gamma distribution for different parameter values.](/diagrams/algebrica/gamma-distribution.png)

When \\(\alpha\\) grows beyond \\(1\\), the gamma density no longer peaks at zero but forms a maximum at a positive value of \\( x \\). As \\(\alpha\\) increases, this peak moves to the right and the overall shape becomes smoother and less skewed.


As with any continuous distribution, the total area under the density curve must equal 1. This is the same principle that holds for the [normal distribution](<../normal-distribution/>), whose density [integrates](<../definite-integrals/>) to \\(1\\) over the entire real line. The gamma distribution follows the same requirement: its density is defined so that the integral over the positive half-line is exactly equal to \\(1\\). Formally, we have

\\[\int_{0}^{+\infty} \frac{1}{\beta^{\alpha}\,\Gamma(\alpha)}\, x^{\alpha - 1}\, e^{-x/\beta}\, dx = 1 \\]

Solving the integral by applying the [substitution](<../integration-by-substitution/>) \\( x = \beta t \\), which gives \\( dx = \beta\, dt \\), we can rewrite the expression as:

\\[\int_{0}^{+\infty} \frac{1}{\beta^{\alpha}\,\Gamma(\alpha)}\, (\beta t)^{\alpha - 1} e^{-t}\, \beta\, dt \\]

After collecting the powers of \\( \beta \\), this becomes:

\\[\frac{1}{\Gamma(\alpha)} \int_{0}^{+\infty} t^{\alpha - 1} e^{-t}\, dt \\]

The integral on the right-hand side is exactly the definition of the gamma function, so we obtain

\\[\frac{1}{\Gamma(\alpha)} \cdot \Gamma(\alpha) = 1 \\]

## Key features

  * \\[\text{1. } \quad f(x) = \frac{1}{\Gamma(\alpha)\,\beta^{\alpha}} \, x^{\alpha - 1} e^{-x/\beta} \quad x > 0 \\]

  * \\[\text{2. } \quad \mu = E(X) = \alpha\,\beta \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = \alpha\,\beta ^{2} \\]

  * \\[\text{4. } \quad \sigma = \beta\,\sqrt{\alpha} \\]


##### Each expression highlights a key property of the Gamma distribution, whose density is defined through the Gamma function \\(\Gamma(\alpha)\\). Its mean and variability depend jointly on the shape parameter \\(\alpha\\) and the scale parameter \\(\beta\\), determining how the distribution models waiting times and positively skewed processes.

## Expected value of the gamma distribution

As introduced in the section on continuous random variables, the [expected value of a continuous random variable](<../mean-or-expected-value-of-a-random-variable/>) describes the central tendency of its distribution and is defined as

\\[\mu = E(X) = \int_{-\infty}^{+\infty} x\, f(x)\, dx \\]

This general formula applies to any continuous distribution, where \\( f(x) \\) denotes the probability density function of the random variable \\( X \\). For the gamma distribution with shape parameter \\( \alpha \\) and scale parameter \\( \beta \\), the density is:

\\[f(x;\alpha,\beta)=\frac{1}{\beta^{\alpha}\Gamma(\alpha)}\, x^{\alpha - 1} e^{-x/\beta} \quad x>0 \\]

so the expected value is computed as:

\\[\mu = E(X) = \int_{0}^{+\infty} x\, \frac{1}{\beta^{\alpha}\Gamma(\alpha)}\, x^{\alpha - 1} e^{-x/\beta}\, dx \\]

Combining the powers of \\( x \\), we obtain:

\\[E(X) = \frac{1}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} x^{\alpha} e^{-x/\beta}\, dx \\]

To simplify the integral, we apply the change of variable \\( x = \beta t \\), which gives \\( dx = \beta\, dt \\). Substituting, we have:

\\[E(X) = \frac{1}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} (\beta t)^{\alpha} e^{-t}\, \beta\, dt \\]

Collecting the powers of \\( \beta \\):

\\[\begin{align} E(X) &= \frac{1}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} (\beta t)^{\alpha} e^{-t}, \beta, dt \\\\[0.6em] &= \frac{\beta^{\alpha}\,\beta}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} t^{\alpha} e^{-t}\, dt \\\\[0.6em] &= \frac{\beta}{\Gamma(\alpha)} \int_{0}^{+\infty} t^{\alpha} e^{-t}\, dt \end{align} \\]

The remaining integral is recognized as the gamma function evaluated at \\( \alpha + 1 \\):

\\[\int_{0}^{+\infty} t^{\alpha} e^{-t}\, dt = \Gamma(\alpha + 1) \\]

Using the identity \\( \Gamma(\alpha + 1) = \alpha \Gamma(\alpha) \\), we find

\\[\mu = \alpha \beta \\]

This shows that the mean of the gamma distribution depends on both parameters: \\( \alpha \\) shapes how the mass is distributed along the positive axis, while \\( \beta \\) stretches the distribution horizontally, and their product gives the average value of the variable.


The gamma distribution is sometimes written in an alternative form that uses a rate parameter instead of a scale parameter. In that case, the density is expressed as:

\\[f(x;\alpha,\beta)=\frac{\beta^{\alpha}}{\Gamma(\alpha)}, x^{\alpha - 1} e^{-\beta x} \quad x>0 \\]

where \\( \beta \\) now plays the role of the rate which is the inverse of the scale. Under this parametrization, the expected value becomes:

\\[E(X) = \frac{\alpha}{\beta} \\]

##### The two expressions for the mean are completely consistent: with a scale parameter the mean is \\( \alpha\beta \\), while with a rate parameter it is \\( \alpha/\beta \\). The difference comes solely from the choice of parametrization, not from the distribution itself.

## Variance of the gamma distribution

The variance of the gamma distribution can be obtained from the general definition of [variance for continuous random variables](<../variance-and-covariance-of-a-random-variable/>):

\\[\sigma^{2} = E(X^{2}) - [E(X)]^{2} \\]

Using the corresponding integral expression and substituting the gamma density, we have:

\\[E(X^{2}) = \int_{0}^{+\infty} x^{2}\, \frac{1}{\beta^{\alpha}\Gamma(\alpha)}\, x^{\alpha - 1} e^{-x/\beta}\, dx \\]

Combining the powers of \\( x \\), this becomes:

\\[E(X^{2}) = \frac{1}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} x^{\alpha + 1} e^{-x/\beta}\, dx \\]

Applying the change of variable \\( x = \beta t \\), we may rewrite the integral as:

\\[\begin{align} E(X^{2}) &= \frac{1}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} (\beta t)^{\alpha + 1} e^{-t}\, \beta\, dt \\\\[0.6em] &= \frac{\beta^{\alpha + 2}}{\beta^{\alpha}\Gamma(\alpha)} \int_{0}^{+\infty} t^{\alpha + 1} e^{-t}\, dt \\\\[0.6em] &= \frac{\beta^{2}}{\Gamma(\alpha)} \Gamma(\alpha + 2) \end{align} \\]

Using the identity:

\\[\Gamma(\alpha + 2) = (\alpha + 1)\alpha\, \Gamma(\alpha) \\]

we obtain:

\\[E(X^{2}) = \beta^{2} \alpha(\alpha + 1) \\]

Since the mean of the gamma distribution is

\\[E(X) = \alpha \beta \\]

the variance becomes:

\\[\sigma^{2} = \beta^{2}\alpha(\alpha + 1) - (\alpha\beta)^{2} = \alpha\beta^{2} \\]


As with the expected value, when the gamma distribution is written using a rate parameter instead of a scale parameter, the expression for the variance also changes. Under this form, the variance of the gamma distribution becomes:

\\[\mathrm{Var}(X)=\frac{\alpha}{\beta^{2}} \\]

##### This result is fully consistent with the scale–parameter version. It is simply a different way of writing the same distribution.

## Special cases of the gamma distribution

The gamma distribution includes several important special cases, each obtained by choosing specific values for its parameters. One of the most notable is the [exponential distribution](<../exponential-distribution/>), which arises when \\( \alpha = 1 \\). In this situation, the density takes the simpler form

\\[f(x;\lambda)= \begin{cases} \lambda e^{-\lambda x} & x>0\\\\[0.4em] 0 & x \le 0 \end{cases} \\]

where \\( \lambda \\) is the rate parameter that determines how quickly the distribution decays. The exponential distribution is often used to model waiting times between successive events that occur independently and at a constant average rate, such as the time between arrivals in a [Poisson process](<../poisson-distribution/>).

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
