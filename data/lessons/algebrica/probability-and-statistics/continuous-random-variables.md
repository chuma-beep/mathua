> Content sourced from [Algebrica](https://algebrica.org/continuous-random-variables/) — CC BY-NC 4.0

## Definition of a continuous random variable

A **continuous random variable** is a function that assigns a real number to each element of a continuous sample space. Unlike [discrete random variables](<../discrete-random-variables>), it can take on infinitely many values within a given interval. Formally, a continuous random variable is a function:

\\[X : \Omega \to \mathbb{R} \\]

where \\( \Omega \\) is a continuous sample space composed of infinitely many outcomes that can vary smoothly. In this case, probabilities are not assigned to individual values but to intervals of values, and they are determined through a probability density function.


A continuous random variable has the special property that the probability of it assuming any precise value is zero.

\\[P(X = x_0) = 0 \\]

What matters instead is the probability of the variable falling within a given range of values. In general, therefore, when \\( X \\) is a continuous random variable, the probability of it taking any specific value is zero, and we have:

\\[P(a < X \le b) = P(a < X < b) + P(X = b) \\]

Since \\(P(X = b) = 0\\) we have:

\\[P(a < X \le b) = P(a < X < b) \\]

##### This means that, for continuous random variables, the exact endpoints of the interval are insignificant. The probability remains the same whether the bounds are included or not.

## Probability from the density function

To determine the probability distribution of a continuous random variable \\( X \\) within a general interval \\((a, b)\\), it is necessary to compute the following [definite integral](<../definite-integrals/>):

\\[P(a < X < b) = \int_a^b f(x)\,dx \\]

The resulting probability corresponds to the area under the curve of the probability density function \\( f(x) \\) between the limits \\( a \\) and \\( b \\).

![](https://algebrica.org/wp-content/uploads/resources/images/continuous-random-vars-1.png)


To describe a continuous probability distribution, we use a **probability density function** \\( f(x) \\). This function defines how the probability is distributed across the possible values of the random variable \\( X \\). A function \\( f(x) \\) qualifies as a valid probability density if it satisfies the following conditions:

\\[\begin{align} & f(x) \ge 0, \quad \text{for all } x \in \mathbb{R} \\\\[12pt] & \int_{-\infty}^{+\infty} f(x),dx = 1 \\\\[6pt] & P(a < X < b) = \int_a^b f(x)\,dx \end{align} \\]

These conditions ensure that the density function is always non-negative, that the total probability over the entire real line equals one, and that probabilities for any interval can be obtained by integration.

Unlike [discrete random variables](<../discrete-random-variables>), which assign probabilities to individual outcomes through a probability mass function (PMF), continuous random variables are described by a probability density function (PDF). In this case, probabilities are not tied to specific points but to areas under the curve, representing the likelihood of the variable falling within a given interval.

## Cumulative distribution function

As with discrete random variables, it is also possible to define a cumulative distribution function \\( F(x) \\) for a continuous random variable \\( X \\) with a probability density function \\( f(x) \\). The cumulative distribution function represents the probability that the variable \\( X \\) takes on a value less than or equal to \\( x \\). Formally, for a variable \\( X \\), the cumulative distribution function is defined as:

\\[F(x) = P(X \le x) = \int_{-\infty}^{x} f(t)\,dt \\]

It describes how probability accumulates from the left tail of the distribution up to the point \\( x \\). Considering a general interval \\((a, b)\\) like the one shown in the previous figure, the probability that the variable \\( X \\) takes a value within that range is given by:

\\[P(a < X < b) = F(b) - F(a) \\]

where \\( F(b) \\) and \\( F(a) \\) represent the values of the cumulative distribution function that is the [antiderivative](<../definite-integrals/>) of the density function evaluated at the upper and lower limits of the interval.

## Joint probability density function

In the case of two continuous random variables \\( X \\) and \\( Y \\), we can describe how their probabilities are distributed together through a joint probability density function \\( f(x, y) \\). The term joint indicates that the function represents the combined behavior of both variables that is, how the probability is shared across different pairs of values \\((x, y)\\) in the two-dimensional plane. It allows us to study the relationship between \\( X \\) and \\( Y \\), including whether they are independent or correlated.

A function \\( f(x, y) \\) is a valid joint probability density function if it satisfies the following conditions:

\\[\begin{align} & f(x, y) \ge 0, \forall \, (x, y) \\\\[12pt] & \int_{-\infty}^{+\infty} \int_{-\infty}^{+\infty} f(x, y)\,dx\,dy = 1 \\\\[6pt] & P\big[(X, Y) \in A\big] = \int \int_A f(x, y)\,dx\,dy \end{align} \\]


Starting from the joint probability density function \\( f_{X,Y}(x, y) \\), it is possible to define the marginal distributions of \\( X \\) and \\( Y \\). In the continuous case, these functions describe the probability behavior of each variable individually, regardless of the other. In other words, they represent the total probability obtained by integrating out the other variable. The marginal density functions are given by:

\\[\begin{align} f_X(x) &= \int_{-\infty}^{+\infty} f_{X,Y}(x, y)\,dy \\\\[12pt] f_Y(y) &= \int_{-\infty}^{+\infty} f_{X,Y}(x, y)\,dx \end{align} \\]

where \\( f_X(x) \\) is the marginal density of \\( X \\), and \\( f_Y(y) \\) is the marginal density of \\( Y \\).

## Example 1

Consider a continuous random variable \\( X \\) that represents the lifetime in years of a light bulb. Suppose its probability density function is given by:

\\[f(x) = \begin{cases} \lambda e^{-\lambda x} & x \ge 0 \\\\[6pt] 0 & x < 0 \end{cases} \\]

where \\( \lambda > 0 \\) is a constant that determines how quickly the probability decreases. This is known as the exponential distribution, commonly used to model waiting times or lifetimes of components.

##### The constant \\( e \approx 2.71828 \\) naturally appears in this context through the [exponential function](<../exponential-function>). Its origin can be traced back to the idea of a [limit of a sequence](<../euler-number-limit-sequence/>), where it first emerges as a fundamental mathematical constant.


The probability that the light bulb lasts between \\( a = 1 \\) and \\( b = 3 \\) years is obtained by integrating the density function:

\\[P(1 < X < 3) = \int_1^3 \lambda e^{-\lambda x}\,dx = e^{-\lambda} - e^{-3\lambda} \\]

Graphically, this probability corresponds to the area under the curve of \\( f(x) \\) between \\( x = 1 \\) and \\( x = 3 \\).

If, for instance, \\( \lambda = 0.5 \\) then:

\\[P(1 < X < 3) = e^{-0.5} - e^{-1.5} \approx 0.383 \\]

meaning there is about a 38% chance that the bulb will last between one and three years.

## Mean of a continuous random variable

Also for continuous random variables, it is possible to define the concept of the mean or [expected value](<../mean-or-expected-value-of-a-random-variable/>). It represents the long-run average value that the variable would take if the experiment were repeated infinitely many times. For a continuous random variable \\( X \\) with a probability density function \\( f(x) \\), the mean is given by:

\\[\mu = E(X) = \int_{-\infty}^{+\infty} x f(x)\,dx \\]

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

938

[Sampling Distributions](https://algebrica.org/sampling-distributions/)

1k

[Bayes’ Theorem](https://algebrica.org/bayes-theorem/)

1.4k

[Confidence Intervals](https://algebrica.org/confidence-intervals/)
