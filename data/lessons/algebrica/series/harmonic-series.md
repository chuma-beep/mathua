> Content sourced from [Algebrica](https://algebrica.org/harmonic-series/) — CC BY-NC 4.0

## What is a harmonic series

The **harmonic series** is defined as the infinite sum:

\\[\sum_{k=1}^{\infty} \frac{1}{k} \\]

where each term is the reciprocal of a [natural number](<../natural-numbers>). Despite the terms approaching zero, the [series](<../series>) diverges, meaning the sum grows without bound. In particular, since the harmonic series is a series with [positive terms](<../series-with-positive-terms/>), it can be shown that it diverges to positive infinity. In fact, since the series has positive terms, the limit of the sequence of its partial sums exists:

\\[S = \lim_{n \to +\infty} s_n = \lim_{n \to +\infty} \sum_{k=1}^{n} \frac{1}{k} \\]

At first glance, one might think that as \\( n \to \infty \\), the terms of the harmonic series tend to zero and therefore the series itself might converge. However, this is a logical error: although the terms \\( \frac{1}{n} \\) do approach zero, they do not do so fast enough for the series to converge.

![](https://algebrica.org/wp-content/uploads/resources/images/harmonic-series-1-2.png)

##### Here is the graph of the partial sums of the harmonic series up to \\(n = 100.\\) As can be seen, the curve rises slowly yet unceasingly, confirming that the series is divergent.


There are several ways to demonstrate that the harmonic series diverges. One approach involves using [partial sums](<../series/>), as shown below:

\\[\sum_{k=1}^{\infty} \frac{1}{k} = 1 + \frac{1}{2} + \left( \frac{1}{3} + \frac{1}{4} \right) + \left( \frac{1}{5} + \frac{1}{6} + \frac{1}{7} + \frac{1}{8} \right) + \cdots \\]

Each group contains twice as many terms as the previous one. By observing that each group adds at least \\( \frac{1}{2} \\) to the sum, we conclude that the overall series grows without bound:

\\[\sum_{k=1}^{\infty} \frac{1}{k} = \infty \\]

Hence, the harmonic series diverges.

Knowing how the harmonic series and its variants behave is useful, since the comparison test often lets us relate a complex series to a harmonic one to check whether it converges.

## Generalized harmonic series (p-series)

Let’s consider a generalized version of the harmonic series, where the denominator is raised to a [power](<../powers>) \\( a \\):

\\[\sum_{k=1}^{\infty} \frac{1}{k^a} \\]

This series, known as **generalized harmonic series** , has an additional feature compared to the standard harmonic series: its convergence or divergence depends on the value of the exponent \\( a \\). We have:

  * If \\( a > 1 \\), the series converges.
  * If \\( a \leq 1 \\), the series diverges.


A necessary condition for the convergence of a series \\( \sum a_k \\) is that the general term tends to zero:

\\[\lim_{k \to \infty} a_k = 0. \\]

If this condition is not satisfied, that is, if \\( \lim_{k \to \infty} a_k \neq 0 \\), then the series diverges. In the case of the generalized harmonic series with exponent \\( a \leq 1 \\), this condition fails. For example, when \\( a = 0 \\), the terms become constant \\( a_k = 1 \\), and:

\\[\lim_{k \to \infty} \frac{1}{k^0} = \lim_{k \to \infty} 1 = 1 \neq 0. \\]

Hence, the series diverges because its general term does not tend to zero.


When \\( a > 1 \\), we can apply the [integral test](<../integral-test-for-series-convergence/>) to determine the convergence of the series. Consider the corresponding improper integral, we have:

\\[\int_1^{\infty} \frac{1}{x^a} \, dx \\]

Since \\( a > 1 \\), we have:

\\[\int_1^{\infty} \frac{1}{x^a} \, dx = \lim_{t \to \infty} \int_1^t x^{-a} \, dx \\]

By evaluating the integral at the endpoints, we obtain:

\\[\lim_{t \to \infty} \left[ \frac{x^{1-a}}{1 - a} \right]_1^t = \lim_{t \to \infty} \left( \frac{t^{1 - a}}{1 - a} - \frac{1}{1 - a} \right) \\]

Because \\( a > 1 \\), the exponent \\( 1 - a < 0 \\), so \\( t^{1 - a} \to 0 \\). Therefore:

\\[\int_1^{\infty} \frac{1}{x^a} \, dx = \frac{1}{a - 1} \\]

which is finite. Hence the series converges for all \\( a > 1 \\).

## Logarithmically modified harmonic series

Finally, let us consider the following series:

\\[\sum_{k=2}^{\infty} \frac{1}{k (\log^\alpha k)} \\]

This is a so-called **logarithmically modified** series, where the presence of the [logarithmic](<../logarithms>) term affects the rate at which the series converges or diverges. The convergence of the series depends on the exponent \\( \alpha \\) in the logarithmic term:

  * If \\( \alpha > 1 \\), the series converges.
  * If \\( \alpha \leq 1 \\), the series diverges.


The summation starts at \\( k = 2 \\) to avoid the singularities at \\( k = 0 \\) (where \\( \log 0 \\) is undefined) and \\( k = 1 \\) (where \\( \log 1 = 0 \\), causing division by zero).


To demonstrate the convergence, we apply the improper integral test. We consider the function:

\\[f(x) = \frac{1}{x (\log x)^\alpha} \\]

which is positive, [continuous](<../continuous-functions/>), and [decreasing](<../increasing-and-decreasing-functions/>) for \\( x \geq 2 \\). We then evaluate the improper integral:

\\[\int_2^{\infty} \frac{1}{x (\log x)^\alpha} \, dx \\]

Using the substitution \\( u = \log x \\), we get:

\\[\int_2^{\infty} \frac{1}{x (\log x)^\alpha} \, dx = \int_{\log 2}^{\infty} \frac{1}{u^\alpha} \, du \\]

This integral converges if and only if \\( \alpha > 1 \\). Therefore, by the integral test, the series converges if and only if \\( \alpha > 1 \\).

## Example

Let us determine the nature of the following series:

\\[\sum_{n=2}^{\infty} \frac{1}{n \sqrt{\log n}} \\]


We observe that this series has the general form:

\\[\sum_{n=2}^{\infty} \frac{1}{n (\log n)^\alpha} \\]

with \\( \alpha = \frac{1}{2} \\). This is known as a logarithmically modified harmonic series. According to known results, the series converges if and only if \\( \alpha > 1 \\).

Since \\( \alpha = \frac{1}{2} < 1 \\), the series diverges.

We conclude that the series diverges by comparison with a divergent harmonic series modified by a logarithmic term.

Series

Series are infinite sums of sequence terms, typically indexed by natural numbers.

1.3k

[Series](https://algebrica.org/series/)

2.4k

[Cauchy’s Convergence Criterion for Series](https://algebrica.org/cauchy-convergence-criterion-series/)

1.4k

[Series with Positive Terms](https://algebrica.org/series-with-positive-terms/)

1.6k

[Geometric Series](https://algebrica.org/geometric-series/)

1.4k

[Integral Test for Series Convergence](https://algebrica.org/integral-test-for-series-convergence/)

1.1k

[Root Test for Series Convergence](https://algebrica.org/root-test-for-series-convergence/)

1.9k

[Leibniz’s Criterion](https://algebrica.org/leibniz-criterion/)

640

[Function Series](https://algebrica.org/function-series/)

635

[Power Series](https://algebrica.org/power-series/)

1.7k

[Taylor Series](https://algebrica.org/taylor-series/)

4 comments

1.2k

[Fourier Series](https://algebrica.org/fourier-series/)
