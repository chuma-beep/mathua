> Content sourced from [Algebrica](https://algebrica.org/exponential-distribution/) — CC BY-NC 4.0

## What is the exponential probability distribution

The **exponential distribution** characterizes the time elapsed between random, independent events occurring at a constant average rate. It is often applied to model waiting times, such as the [interval](<../intervals/>) before the next customer arrives or a machine experiences a failure, providing a simple and intuitive way to describe processes governed by randomness and continuity over time. In formal terms, given a random variable \\( X \\) and a positive real parameter \\( \lambda \\), the exponential distribution is defined by the following probability density function:

\\[f(x; \lambda) = \begin{cases} \lambda e^{-\lambda x} & \text{for } x > 0 \\\\[0.6em] 0 & \text{for } x \le 0 \end{cases} \\]

  * \\( \lambda \\) is the rate parameter, which determines how frequently events occur.
  * A larger value of \\( \lambda \\) corresponds to events happening more often.
  * A smaller value indicates longer expected waiting times between events.


This function represents a continuous probability model that quantifies the likelihood of waiting a certain amount of time before an event occurs in a process with a constant rate of occurrence.

![](/diagrams/algebrica/exponential-distribution-1.png)

From the definition of the exponential distribution, and as illustrated in the figure above, we can observe that the total area under the curve of \\( f(x; \lambda) \\) equals one. In other words, the probability density function is normalized, satisfying the fundamental property

\\[\int_{-\infty}^{+\infty} f(x; \lambda)\,dx = 1 \\]

Since \\( f(x; \lambda) = 0 \\) for \\( x \le 0 \\), the [definite integral](<../definite-integrals/>) effectively reduces to

\\[\int_{0}^{+\infty} \lambda e^{-\lambda x}\,dx = \Bigl[-e^{-\lambda x}\Bigr]_{0}^{+\infty} = 1 \\]

Therefore, the probability density function of the exponential distribution satisfies the normalization condition

\\[\int_{-\infty}^{+\infty} f(x; \lambda)\,dx = 1 \\]

which ensures that the total probability over the entire [domain](<../determining-the-domain-of-a-function/>) is equal to one, as required for any valid continuous probability distribution.

##### The sum of many independent exponential random variables approaches the [normal distribution](<../normal-distribution>), as stated by the Central Limit Theorem.

The exponential distribution derives its name from the [exponential function](<../exponential-function/>) that shapes its probability density. The decreasing pattern of the curve \\( \lambda e^{-\lambda x} \\) shows how the probability of observing longer waiting times diminishes exponentially, highlighting the intrinsic link between the concept of probability and the mathematical properties of the exponential function.

## Key features

  * \\[\text{1. } \quad f(x) = \lambda\, e^{-\lambda x} \quad x \ge 0 \\]

  * \\[\text{2. } \quad \mu = E(X) = \frac{1}{\lambda} \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = \frac{1}{\lambda^{2}} \\]

  * \\[\text{4. } \quad \sigma = \frac{1}{\lambda} \\]


##### Each expression highlights a key property of the exponential distribution, showing how it models waiting times between events, how its mean and variability depend on the rate parameter \\(\lambda\\), and how its memoryless nature distinguishes it from other continuous distributions.

## Expected value and variance of the exponential distribution

As introduced in the section on continuous random variables, the [expected value of a continuous random variable](<../mean-or-expected-value-of-a-random-variable/>) provides a measure of the central tendency of the distribution and is defined as:

\\[\mu = E(X) = \int_{-\infty}^{+\infty} x\,f(x)\,dx \\]

This expression is general and applies to any continuous distribution, where \\( f(x) \\) denotes the probability density function of the random variable \\( X \\). In the specific case of the exponential distribution, the expected value is obtained as

\\[\mu = E(X) = \int_{0}^{+\infty} x\,\lambda e^{-\lambda x}\,dx = \frac{1}{\lambda} \\]

This result represents the average waiting time between two consecutive events in a process that occurs at a constant rate \\( \lambda \\).


The variance of the exponential distribution can be obtained from the general definition of [variance for continuous random variables](<../variance-and-covariance-of-a-random-variable/>):

\\[\sigma^2 = E(X^2) - [E(X)]^2 \\]

Using the corresponding integral expression and substituting the exponential density function, we have:

\\[E(X^2) = \int_{0}^{+\infty} x^2 \lambda e^{-\lambda x}\,dx = \frac{2}{\lambda^2} \\]

Therefore, the variance is

\\[\sigma^2 = \frac{2}{\lambda^2} - \left(\frac{1}{\lambda}\right)^2 = \frac{1}{\lambda^2} \\]

This shows that the variability of the exponential distribution decreases as the rate parameter \\( \lambda \\) increases, meaning that higher event frequencies correspond to shorter and more consistent waiting times.

## The survival function of the exponential distribution

The survival function of the exponential distribution expresses the probability that the random variable \\( X \\) takes on a value greater than a given threshold \\( x \\). It is defined as the complement of the cumulative distribution function:

\\[S(x) = P(X > x) = 1 - F(x) \\]

For the exponential distribution with rate parameter \\( \lambda > 0 \\), the survival function takes the form

\\[S(x) = e^{-\lambda x} \quad x \ge 0 \\]

This function shows that the probability of “surviving” beyond a certain time decreases exponentially as \\( x \\) increases. In other words, the longer we wait, the smaller the chance that the event has not yet occurred.

## The memoryless property

An important property of a random variable \\( X \\) that follows an exponential distribution is the so-called memoryless property. This characteristic can be expressed by the equality

\\[P(X > s + t \mid X > s) = P(X > t) \\]

which means that the probability of waiting an additional time \\( t \\) does not depend on how much time \\( s \\) has already elapsed. In other words, the exponential distribution “forgets” past events, making it unique among continuous probability distributions.

Starting from the definition of conditional probability, we can write:

\\[P(X > s + t \mid X > s) = \frac{P(X > s + t)}{P(X > s)} \\]

By substituting the survival function of the exponential distribution, we obtain:

\\[\begin{align} P(X > s + t \mid X > s) &= \frac{\int_{s+t}^{+\infty} \lambda e^{-\lambda x}\,dx} {\int_{s}^{+\infty} \lambda e^{-\lambda x}\,dx} \\\\[0.6em] &= \frac{e^{-\lambda (s + t)}}{e^{-\lambda s}} \\\\[0.6em] &= e^{-\lambda t} \end{align} \\]

Finally, observing that \\( e^{-\lambda t} = P(X > t) \\), we can conclude that

\\[P(X > s + t \mid X > s) = P(X > t) \\]

This confirms that the exponential distribution has no memory of past events. The probability of waiting an additional time \\( t \\) remains the same, regardless of how long \\( s \\) has already elapsed.

## Example 1

A certain type of light bulb has a lifetime that follows an exponential distribution with a rate parameter \\( \lambda = 0.002 \\) failures per hour.

  * Find the probability that a bulb lasts more than 400 hours.
  * Determine the probability that a bulb lasts between 200 and 600 hours.
  * Compute the expected lifetime and the variance of the bulb’s duration.


To find the probability that the bulb operates for more than 400 hours, we use the survival function of the exponential distribution:

\\[P(X > x) = e^{-\lambda x} \\]

Substituting the given value of \\( \lambda = 0.002 \\) and \\( x = 400 \\), we obtain

\\[P(X > 400) = e^{-0.002 \times 400} = e^{-0.8} \approx 0.4493 \\]

This means there is approximately a 44.9% probability that the bulb will last beyond 400 hours.


To determine the probability that the bulb’s lifetime falls between 200 and 600 hours, we calculate the difference between the probabilities of lasting more than 200 hours and more than 600 hours:

\\[\begin{align} P(200 < X < 600) &= P(X > 200) - P(X > 600) \\\\[0.6em] &= e^{-0.002 \times 200} - e^{-0.002 \times 600} \\\\[0.6em] &= e^{-0.4} - e^{-1.2} \approx 0.6703 - 0.3010 = 0.3693 \end{align} \\]

Thus, there is approximately a 36.9% probability that the bulb operates between 200 and 600 hours.


The expected lifetime of a bulb that follows an exponential distribution is given by the inverse of the rate parameter:

\\[E(X) = \frac{1}{\lambda} = \frac{1}{0.002} = 500 \quad \text{hours} \\]

This means that, on average, a bulb is expected to last 500 hours before failing.

The variance, which measures the spread of the distribution, is obtained as the square of the mean lifetime:

\\[\sigma^2 = \frac{1}{\lambda^2} = \frac{1}{(0.002)^2} = 250000 \\]

Consequently, the standard deviation is:

\\[\sigma = \sqrt{\sigma^2} = 500 \ \text{hours} \\]

Recapping the results of the exercise, we have:

\\[\begin{aligned} &P(X > 400) = 0.4493 \\\\[6px] &P(200 < X < 600) = 0.3693 \\\\[6px] &E(X) = 500 \ \text{hours} \\\\[6px] &\sigma^2 = \mathrm{Var}(X) = 250000 \\\\[6px] &\sigma = 500 \ \text{hours} \end{aligned} \\]

## Connection to the gamma distribution

The exponential distribution can be viewed as one of the simplest special cases of the [gamma distribution](<../gamma-distribution>). In particular, by setting the shape parameter of the gamma distribution to \\( \alpha = 1 \\), the gamma density:

\\[G(x;\alpha,\beta) = \frac{1}{\beta^{\alpha}\Gamma(\alpha)}\, x^{\alpha - 1} e^{-x/\beta} \quad x>0 \\]

simplifies considerably. When \\( \alpha = 1 \\), we have \\( \Gamma(1) = 1 \\) and the term \\( x^{\alpha - 1} \\) becomes \\( x^{0} = 1 \\). The density therefore reduces to

\\[G(x;1,\beta)=\frac{1}{\beta} e^{-x/\beta} \quad x>0 \\]

which is exactly the exponential distribution with scale parameter \\( \beta \\). If we switch to the rate parametrization by setting \\( \lambda = 1/\beta \\), the same expression becomes \\[G(x;\lambda)=\lambda\, e^{-\lambda x} \quad x>0 \\]

This shows that the exponential distribution is simply the gamma distribution with shape equal to one.

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

938

[Sampling Distributions](https://algebrica.org/sampling-distributions/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
