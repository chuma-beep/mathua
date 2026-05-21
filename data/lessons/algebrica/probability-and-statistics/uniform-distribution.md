> Content sourced from [Algebrica](https://algebrica.org/uniform-distribution/) — CC BY-NC 4.0

## Definition of the uniform distribution

The uniform distribution is one of the simplest continuous distributions to describe. It models a random variable that can take any value within a specified [interval](<../intervals/>), assigning the same probability to all subintervals of equal length. In other words, the possible outcomes are spread evenly across the interval, with no region being more likely than another.

In formal terms, a [continuous random variable](<../continuous-random-variables/>) \\( X \\) is said to follow a uniform distribution on the interval \\(A = [a, b]\\) if its probability density function is constant on that interval. Formally we have:

\\[U(x;a,b) = \begin{cases} \dfrac{1}{\,b - a\,} & x \in A \\\\[0.6em] 0 & x \notin A \end{cases} \\]

For any two subintervals \\( I_1 \\) and \\( I_2 \\) contained in \\([a, b]\\) and having the same length, the uniform distribution assigns them the same probability. If \\( w \\) denotes the common width of these intervals, then:

\\[P(I_1) = P(I_2) = \frac{w}{\,b - a\,} \\]

The image shows the density of the uniform distribution over the interval \\([a, b]\\). The flat horizontal line indicates that every value between \\(a\\) and \\(b\\) is assigned the same likelihood. The endpoints are drawn as open circles to highlight a key property of continuous distributions: individual points have zero probability and we have:

\\[P(X = x_0) = 0 \\]

![The figure illustrates the density of the uniform distribution on the interval \[a,b\], showing its constant height and the equal likelihood assigned to all values within the range.](/diagrams/algebrica/uniform-distribution.png)

## Key features

  * \\[\text{1. } \quad U(x;a,b) = \frac{1}{b - a} \quad a \le x \le b \\]

  * \\[\text{2. } \quad \mu = E(X) = \frac{a + b}{2} \\]

  * \\[\text{3. } \quad \sigma^{2} = \mathrm{Var}(X) = \frac{(b - a)^{2}}{12} \\]

  * \\[\text{4. } \quad \sigma = \frac{b - a}{2\sqrt{3}} \\]


##### Each expression highlights a key property of the uniform distribution, showing how probability is spread evenly across the interval and how its mean and variability depend solely on the endpoints \\(a\\) and \\(b\\).

## Mean of the uniform distribution

The [mean](<../introduction-to-the-mean/>), or [expected value](<../mean-or-expected-value-of-a-random-variable/>), of a uniform distribution describes the average outcome we would expect from a random variable that can take any value between \\(a\\) and \\(b\\) with equal likelihood. Because the distribution is completely flat over this interval, computing the expected value is straightforward and follows directly from the general definition:

\\[\mu = E(X) = \int_{a}^{b} x\, U(x;a,b) \, dx \\]

For the uniform distribution, the probability density function is constant:

\\[f(x) = \frac{1}{b - a} \qquad \text{for } a \le x \le b \\]

Substituting this expression into the [integral](<../definite-integrals/>) gives us:

\\[E(X) = \int_{a}^{b} x \frac{1}{b - a} \, dx \\]

Since the density does not vary across the interval, we can simply factor it out of the integral:

\\[E(X) = \frac{1}{b - a} \int_{a}^{b} x \, dx \\]

The remaining integral is elementary. The antiderivative of \\(x\\) is \\(x^{2}/2\\), so evaluating it between \\(a\\) and \\(b\\) yields:

\\[\int_{a}^{b} x \, dx = \frac{b^{2}}{2} - \frac{a^{2}}{2} \\]

Plugging this back into the expression for the expected value, we obtain:

\\[E(X) = \frac{1}{b - a} \left( \frac{b^{2} - a^{2}}{2} \right) \\]

Noting that \\( b^{2} - a^{2} = (b - a)(b + a) \\), the expression simplifies neatly to:

\\[E(X) = \frac{b + a}{2} \\]

##### This shows that the mean of the uniform distribution is exactly the midpoint of the interval \\([a, b]\\).

## Variance of the uniform distribution

The [variance](<../variance-and-covariance-of-a-random-variable/>) of a uniform distribution describes how much the random variable is expected to spread out around its mean. While the mean identifies the central value of the interval, the variance quantifies how concentrated or dispersed the outcomes are across \\([a, b]\\). Because the uniform distribution assigns the same density to every point in the interval, its variability depends entirely on the length of that interval. Formally, the variance for a continuous random variable is defined as:

\\[\sigma^{2} = \mathrm{Var}(X) = E(X^{2}) - [E(X)]^{2} \\]

Starting from the probability density function of the uniform distribution we have:

\\[E(X^{2}) = \int_{a}^{b} x^{2} U(x;a,b)\, dx \\]

For the uniform distribution, the density is:

\\[U(x;a,b) = \frac{1}{b - a} \\]

Substituting this expression into the formula gives:

\\[E(X^{2}) = \int_{a}^{b} x^{2} \frac{1}{b - a} \, dx = \frac{1}{b - a} \int_{a}^{b} x^{2}\, dx \\]

The integral of \\(x^{2}\\) is straightforward to evaluate:

\\[\int_{a}^{b} x^{2}\, dx = \frac{b^{3}}{3} - \frac{a^{3}}{3} \\]

Thus, we have:

\\[E(X^{2}) = \frac{1}{b - a} \left( \frac{b^{3} - a^{3}}{3} \right) \\]

Using the [factorization](<../notable-products/>) \\( b^{3} - a^{3} = (b - a)(b^{2} + ab + a^{2}) \\), the expression simplifies to:

\\[E(X^{2}) = \frac{b^{2} + ab + a^{2}}{3} \\]

We now combine this with the mean of the uniform distribution,

\\[E(X) = \frac{a + b}{2} \\]

so that:

\\[[E(X)]^{2} = \left( \frac{a + b}{2} \right)^{2} = \frac{a^{2} + 2ab + b^{2}}{4} \\]

Putting everything together, we obtain:

\\[\mathrm{Var}(X) = \frac{b^{2} + ab + a^{2}}{3} \- \frac{a^{2} + 2ab + b^{2}}{4} \\]

After simplifying the expression, the variance reduces to the form:

\\[\sigma^{2} = \mathrm{Var}(X) = \frac{(b - a)^{2}}{12} \\]

## Relationship between the uniform and the beta distribution

There is a connection between the uniform distribution and the [beta distribution](<../beta-distribution>) \\(B(x),\alpha,\beta\\). When the two shape parameters are both equal to one, that is \\( \alpha = 1 \\) and \\( \beta = 1 \\), the density of the beta distribution becomes perfectly flat over the interval \\((0,1)\\). The probability density function of a beta distribution is:

\\[B(x;\alpha,\beta) = \frac{x^{\alpha - 1} (1 - x)^{\beta - 1}}{B(\alpha,\beta)} \\] \\[0 < x < 1 \\]

Setting \\( \alpha = 1 \\) and \\( \beta = 1 \\) makes both exponents in the numerator equal to zero, and the beta function evaluates to \\( B(1,1) = 1 \\). As a result, the density simplifies to:

\\[B(x;1,1) = 1 \\]

which is exactly the density of a uniform random variable on \\((0,1)\\).

## Example 1

An industrial cutting machine completes a full cycle in a time that varies slightly due to mechanical tolerances and temperature fluctuations. Measurements show that the cycle time is equally likely to take any value between \\(4.8\\) and \\(5.4\\) seconds. Let \\(X\\) denote the cycle time in seconds, and assume

\\[X \sim \mathrm{U}(x; 4.8, 5.4) \\]


Compute the probability that a randomly selected cycle lasts less than 5 seconds. The density of a uniform distribution on \\((a, b)\\) is constant and equal to \\(1/(b - a)\\). Therefore,

\\[P(X < 5) = \frac{5 - 4.8}{5.4 - 4.8} \\]

We obtain:

\\[P(X < 5) = \frac{0.2}{0.6} = \frac{1}{3} \\]

So the probability is approximately \\(0.333\\).


Determine the probability that the cycle time lies between 5.1 and 5.3 seconds. We have:

\\[P(5.1 \le X \le 5.3) = \frac{5.3 - 5.1}{5.4 - 4.8} \\]

This gives:

\\[P(5.1 \le X \le 5.3) = \frac{0.2}{0.6} = \frac{1}{3} \\]

So the probability is again approximately \\(0.333\\).


Find the expected cycle time \\(E(X)\\). For a uniform distribution on \\((a, b)\\) we have:

\\[E(X) = \frac{a + b}{2} \\]

Thus:

\\[E(X) = \frac{4.8 + 5.4}{2} = \frac{10.2}{2} = 5.1 \\]

The expected cycle time is 5.1 seconds.


Calculate the variance \\(\mathrm{Var}(X)\\). The variance of a uniform distribution on \\((a, b)\\) is:

\\[\mathrm{Var}(X) = \frac{(b - a)^{2}}{12} \\]

Substituting the values we obtain:

\\[\mathrm{Var}(X) = \frac{(5.4 - 4.8)^{2}}{12} = \frac{0.6^{2}}{12} = \frac{0.36}{12} = 0.03 \\]

So the variance is \\(0.03\\).

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
