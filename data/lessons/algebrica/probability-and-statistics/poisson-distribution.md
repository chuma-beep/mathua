> Content sourced from [Algebrica](https://algebrica.org/poisson-distribution/) — CC BY-NC 4.0

## Introduction to the Poisson distribution

The **Poisson distribution** is a discrete probability distribution that describes how many times a specific event may occur within a fixed period of time or space. It is applied when events occur randomly, independently of each other, and with a constant average rate over the observed [interval](<../intervals/>). This distribution is particularly useful for modeling rare or random occurrences, such as the number of customer arrivals per hour or the frequency of network failures in a day.

Formally, the Poisson distribution is expressed as

\\[p(x; \lambda) = P(X = x) = \frac{e^{-\lambda} \lambda^x}{x!} \\]

where:

  * \\( X \\) represents the discrete random variable counting the number of events occurring in a given interval.
  * \\( x \\) is the specific number of occurrences being observed.
  * \\( \lambda \\) is the average rate of occurrence, that is, the expected number of events per interval.
  * \\( e \\) is [Euler’s number](<../euler-number-limit-sequence/>).
  * \\( x! \\) denotes the [factorial](<../factorial>) of \\( x \\).


The parameter \\( \lambda \\) fully characterizes the Poisson distribution and must be positive (\\( \lambda > 0 \\)).

## The Poisson process

The Poisson distribution is closely related to the Poisson process, which models the random occurrence and accumulation of events over time or space, while the Poisson distribution expresses the probability of observing a specific number of such events within a fixed interval. A process can be defined as a Poisson process if it satisfies the following conditions:

  * Events occur one at a time, never simultaneously.
  * The events are independent, meaning that the occurrence of one does not influence the probability of another.
  * The average rate of occurrence (λ) remains constant over time or space.
  * The probability of more than one event occurring in an infinitesimally small interval is negligible.


Under these assumptions, the number of events observed in any fixed interval of duration \\( t \\) follows a Poisson distribution with parameter \\( \lambda t \\), which represents the expected number of occurrences during that interval. Formally, the probability that exactly \\( x \\) events occur in time \\( t \\) is given by:

\\[p(x; \lambda t) = P(X = x) = \frac{e^{-\lambda t}(\lambda t)^x}{x!} \\] \\[x = 0, 1, 2, \ldots \\]

When \\( t = 1 \\), the expression reduces to the standard form of the Poisson distribution with parameter \\( \lambda \\).

## Key features

  * \\[\text{1. } \quad P(X = x) = \dfrac{\lambda^{x} e^{-\lambda}}{x!} \quad x = 0,1,2,\dots \\]

  * \\[\text{2. } \quad \mu = E(X) = \lambda \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = \lambda \\]

  * \\[\text{4. } \quad \sigma = \sqrt{\lambda} \\]


##### Each expression highlights a key property of the Poisson distribution, describing how it models the frequency of rare events, where its mean value lies, and how its variability is directly linked to the rate parameter \\(\lambda\\).

## Mean of the Poisson distribution

The [mean](<../introduction-to-the-mean/>), or [expected value](<../mean-or-expected-value-of-a-random-variable/>), of a Poisson distribution represents the average number of events that can be expected to occur in a fixed interval of time or space. It provides a direct measure of the central tendency of the distribution. Formally, the expected value is defined as:

\\[\mu = E(X) = \sum_{x=0}^{\infty} x \, P(X = x) = \lambda \\]

Substituting the probability mass function of the Poisson distribution we have:

\\[E(X) = \sum_{x=0}^{\infty} x \frac{e^{-\lambda}\lambda^x}{x!} \\]

Since the term with \\( x = 0 \\) equals zero, we can start the summation from \\( x = 1 \\):

\\[E(X) = e^{-\lambda} \sum_{x=1}^{\infty} \frac{x \lambda^x}{x!}. \\]

Using the [factorial](<../factorial>) identity ( x/x! = 1/(x - 1)! ), the expression can be rewritten as:

\\[E(X) = e^{-\lambda} \lambda \sum_{x=1}^{\infty} \frac{\lambda^{x-1}}{(x-1)!}. \\]

By changing the index of summation \\( x-1 = k \\), we obtain:

\\[E(X) = e^{-\lambda} \lambda \sum_{k=0}^{\infty} \frac{\lambda^{k}}{k!}. \\]

The infinite sum equals \\( e^{\lambda} \\), which cancels out the \\(e^{-\lambda} \\) term, giving:

\\[\mu = E(X) = \lambda \\]

This result shows that for a Poisson distribution, the mean value is equal to the parameter \\( \lambda \\), which also represents the expected number of events occurring within the given interval. In other words, \\( \lambda \\) characterizes both the average and the rate at which events occur.

## Variance of the Poisson distribution

The [variance](<../variance-and-covariance-of-a-random-variable/>) of a Poisson distribution describes how much the number of observed events tends to fluctuate around the mean value \\( \mu = \lambda \\). While the mean expresses the expected frequency of occurrence, the variance indicates how widely the actual counts may differ from this average when the experiment is repeated many times. By definition, the variance is written as:

\\[\sigma^2 = \mathrm{Var}(X) = E(X^2) - [E(X)]^2 = \lambda \\]

We can rewrite \\( x^2 \\) as \\( x[(x-1)+1] \\), which gives:

\\[E(X^2) = \lambda e^{-\lambda} \left[ \sum_{x=1}^{\infty} \frac{\lambda^{x-1}}{(x-1)!} \+ \sum_{x=2}^{\infty} \frac{\lambda^{x-2}}{(x-2)!} \right] \\]

Since these sums correspond to exponential series, the expression can be rewritten as:

\\[E(X^2) = \lambda e^{-\lambda} (\lambda e^{\lambda} + e^{\lambda}) = \lambda(\lambda + 1) \\]

Since the mean of the Poisson distribution is \\( E(X) = \lambda \\), the variance is obtained as:

\\[\mathrm{Var}(X) = E(X^2) - [E(X)]^2 = \lambda(\lambda + 1) - \lambda^2 = \lambda \\]

Thus, the variance of a Poisson distribution is equal to its mean:

\\[\sigma^2 = \lambda \\]

##### This equality between mean and variance is a defining property of the Poisson model, indicating that the expected frequency of events \\( \lambda \\) also determines the variability observed within any fixed interval of time or space.

## Cumulative Poisson distribution

In many applications, the focus is not on finding the probability of observing exactly \\( x \\) events, but rather on determining the probability of observing no more than a certain number of events within a fixed time or spatial interval. This concept is described by the cumulative Poisson distribution, obtained by adding together the probabilities of all outcomes from \\( x = 0 \\) up to a chosen value \\( x = r \\):

\\[F(r; \lambda) = \sum_{x=0}^{r} \frac{e^{-\lambda}\lambda^x}{x!}. \\]

Hence, the cumulative probability that the [random variable](<../discrete-random-variables/>) \\( X \\) takes a value less than or equal to \\( r \\) is:

\\[P(X \le r) = e^{-\lambda}\sum_{x=0}^{r}\frac{\lambda^x}{x!} \\]

Each term in this summation represents the probability of exactly \\( x \\) events occurring, while the total sum expresses the probability of observing at most \\( r \\) events in the given interval. These probabilities are usually obtained from Poisson cumulative tables that provide pre-calculated values of \\( P(X \le r) \\) for different values of the parameter \\( \lambda \\). Their structure is similar to that of the [standard normal Z table](<../standard-normal-z-table/>):

λ| x = 0| 1| 2| 3| 4| …  
---|---|---|---|---|---|---  
0.02| 0.980| 1.000| | | | …  
0.04| 0.961| 0.999| 1.000| | | …  
0.06| 0.942| 0.998| 1.000| | | …  
0.08| 0.923| 0.997| 1.000| | | …  
0.10| 0.905| 0.995| 1.000| | | …  
0.15| 0.861| 0.990| 0.999| 1.000| | …  
…| …| …| …| …| …| …  
  
In this table, each value represents a cumulative probability of the Poisson distribution. To use it, locate the value of \\(\lambda\\) along the rows, which indicates the average expected number of events, and find the value of x along the columns, which represents the maximum number of events considered. The cell at the intersection of these two values gives the cumulative probability \\( P(X \le x) \\), meaning the probability that the observed number of events is less than or equal to x for the chosen \\(\lambda\\).

##### Where no numerical value is displayed, the probability is either extremely small or practically equal to one, and therefore omitted for clarity.

## Example 1

Consider a call center that receives customer service calls throughout the day. Historical data show that, on average, five calls are received every ten minutes. We want to determine the probability that exactly eight calls are received during a ten-minute interval.


Let \\( X \\) be the random variable representing the number of calls received in ten minutes. Assuming that calls arrive independently and at a constant average rate, \\( X \\) follows a Poisson distribution with parameter \\( \lambda = 5 \\). The probability mass function is:

\\[P(X = x) = \frac{e^{-\lambda}\lambda^x}{x!} \\]

Substituting \\( x = 8 \\) and \\( \lambda = 5 \\) we obtain:

\\[P(X = 8) = \frac{e^{-5}5^8}{8!} \\]

We have:

\\[P(X = 8) = \frac{e^{-5} \cdot 390625}{40320} \approx 0.0653 \\]

Thus, the probability that exactly eight calls are received in ten minutes is approximately 0.065, or 6.5%.


However, the calculations are not always straightforward, and in such cases it is often more convenient to use cumulative Poisson tables. In these tables, each entry represents a cumulative probability up to a certain value of \\( x \\). Therefore, to find the probability of exactly eight calls, we subtract the cumulative probability up to \\( x = 7 \\) from the cumulative probability up to \\( x = 8 \\), since both values include all outcomes up to their respective limits.

\\[P(X \le 8) = 0.9319, \quad P(X \le 7) = 0.8666 \\]

The difference between them gives:

\\[P(X = 8) = 0.9319 - 0.8666 = 0.0653 \\]

Therefore, the probability that exactly eight calls are received during a ten-minute interval is approximately \\(0.0653\\), meaning there is about a \\(6.5%\\) chance of observing this specific count of calls under the assumed conditions.

## From the binomial to the Poisson distribution

The Poisson distribution can be derived as a limiting form of the [binomial distribution](<../binomial-distribution>) when the number of trials becomes very large and the probability of success in each trial becomes very small. Imagine dividing a fixed time interval \\( T \\) into \\( n \\) smaller subintervals, each of length \\( T/n \\). We assume that:

  * In each subinterval at most one event may occur.
  * The probability of occurrence in any subinterval is proportional to its length: \\[P(E_k) = \frac{\lambda}{n} \quad (k=1,2, \ldots, n) \\] where \\( \lambda = cT \\) represents the expected number of events in the whole interval.
  * All events are independent of one another.


Under these conditions, the number of observed events \\( X \\) follows a binomial distribution: \\[P(X = x) = \binom{n}{x} p^x (1 - p)^{n - x} \quad \text{with } p = \frac{\lambda}{n} \\]

Now consider the [limit](<../limits>) as \\( n \\) grows indefinitely large and \\( p \\) becomes very small, while keeping their product \\( n p = \lambda \\) constant. In this limit, the binomial expression approaches:

\\[p(x; \lambda) = P(X = x) = \frac{e^{-\lambda} \lambda^x}{x!} \\]

which is exactly the Poisson distribution. Let’s see the derivation formally:

\\[P(X = x) = \binom{n}{x} \left(\frac{\lambda}{n}\right)^x \left(1 - \frac{\lambda}{n}\right)^{n - x} \\]

\\[= \frac{n!}{x!(n - x)!} \frac{\lambda^x}{n^x} \left(1 - \frac{\lambda}{n}\right)^n \left(1 - \frac{\lambda}{n}\right)^{-x} \\]

As \\( n \\) increases:

\\[\left(1 - \frac{\lambda}{n}\right)^n \to e^{-\lambda}, \quad \frac{n!}{(n - x)! n^x} \to 1 \\]

and thus we obtain the limit:

\\[p(x; \lambda) = P(X = x) = \frac{e^{-\lambda} \lambda^x}{x!} \\]

##### This limiting relationship shows that the Poisson distribution can be seen as an approximation of the binomial law for rare events that occur independently and infrequently over time or space.

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

938

[Sampling Distributions](https://algebrica.org/sampling-distributions/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
