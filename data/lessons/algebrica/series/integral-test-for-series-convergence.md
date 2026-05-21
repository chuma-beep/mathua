> Content sourced from [Algebrica](https://algebrica.org/integral-test-for-series-convergence/) — CC BY-NC 4.0

## What is the integral test

Determining the sum of an infinite [series](<../series>) and assessing its convergence or divergence is not always straightforward. Several methods are available to study convergence, one of which involves comparing the series to an [improper integral](<../improper-integrals>). This test applies to [series with positive terms](<../series-with-positive-terms/>) and relies on the principle that the convergence of the series can be determined by comparing it to the behavior of an associated improper integral.


Let \\( f \\) be a positive, [decreasing function](<../increasing-and-decreasing-functions/>) defined on \\( [1, +\infty) \\), such as a [rational](<../rational-functions>) or [polynomial function](<../polynomial-function/>). Then the series

\\[\sum_{n=1}^{\infty} f(n) \\]

converges or diverges if and only if the improper integral

\\[\int_{1}^{\infty} f(x) \, dx \\]

does the same, assuming that \\( f \\) is [continuous](<../continuous-functions/>) on \\( [1, +\infty). \\)


![](/diagrams/algebrica/integral-test-series-1.png)

The graph illustrates the connection between a series and an improper integral as stated by the Integral Test.

  * The curve \\( f(x) \\) represents the continuous [function](<../functions>).
  * The gray area shows a portion of the improper integral (the area under the curve from \\( x = 1 \\) to some \\( x = n \\)).
  * The vertical rectangles represent the terms of the series \\( f(n) \\), each with base 1 and height \\( f(n) \\).


##### This visual helps compare the discrete sum (the series) and the continuous accumulation (the integral). Since the rectangles overestimate or underestimate the area depending on the function’s behavior, the integral can be used to determine the convergence of the series.

## Proof

Let us consider the partial sum of the series:

\\[s_k = \sum_{n=1}^{k} f(n) \quad k \in \mathbb{N} \\]

This represents the sum of the first \\( k \\) terms of the series \\( \sum f(n) \\). Since the series has positive terms, the [sequence](<../sequences>) of partial sums \\( {s_k} \\) is increasing and admits a limit as \\( k \to \infty \\):

\\[\lim_{k \to +\infty} s_k = s \in [0, +\infty] \\]

Just as the series is defined by the limit of its partial sums, the improper integral is defined as the limit of the definite integral as the upper bound tends to infinity:

\\[\lim_{k \to +\infty} \int_1^k f(x)\, dx = \int_1^{+\infty} f(x)\, dx \\]

By the linearity of the integral, and in particular its additivity over adjacent intervals, we can write:

\\[\int_1^k f(x)\,dx = \sum_{n=1}^{k-1} \int_n^{n+1} f(x)\,dx \\]

This holds because the definite integral over \\([1, k]\\) can be decomposed into a sum of integrals over the unit-length subintervals \\([n, n+1]\\), which are disjoint and consecutive. Since \\( f \\) is assumed to be decreasing, we obtain the following inequality for all \\( x \in [n, n+1] \\):

\\[f(n+1) \leq f(x) \leq f(n) \\]

By applying the inequality within the integral, we obtain:

\\[\int_n^{n+1} f(n+1)\,dx \leq \int_n^{n+1} f(x)\,dx \leq \int_n^{n+1} f(n)\,dx \\]

By the properties of definite integrals, the first and third terms represent integrals of constant functions. Therefore, the constants can be factored out of the integrals, giving:

\\[f(n+1) \leq \int_n^{n+1} f(x), dx \leq f(n) \\]

Now summing these inequalities from \\( n = 1 \\) to \\( k-1 \\):

\\[\sum_{n=1}^{k-1} f(n+1) \leq \sum_{n=1}^{k-1} \int_n^{n+1} f(x)\, dx \leq \sum_{n=1}^{k-1} f(n) \\]

By taking the limit as \\( k \to \infty \\), we obtain:

\\[\sum_{n=2}^{\infty} f(n) \leq \int_1^{\infty} f(x)\, dx \leq \sum_{n=1}^{\infty} f(n) \\]

which shows that the improper integral is bounded between two versions of the series differing only by the first term \\( f(1) \\). Because the integral lies between two versions of the series that differ only by the first term, if the integral converges, so does the series, and if the integral diverges, the series diverges as well.

## Example

Determine whether the following series converges or diverges using the integral test:

\\[\sum_{n=2}^{\infty} \frac{1}{n \log n} \\]


First, consider the associated function:

\\[f(x) = \frac{1}{x \log x} \\]

defined on the interval \\( x \geq 2 \\). This function is positive, continuous, and decreasing on \\( [2, +\infty) ,\\) so the conditions for using the integral test are satisfied.


Now we evaluate the improper integral:

\\[\int_2^{\infty} \frac{1}{x \log x} \, dx \\]

To compute this, use the [substitution](<../integration-by-substitution/>) \\( u = \log x \\), which implies \\( du = \frac{1}{x} dx \\). The integral becomes:

\\[\int_{\log 2}^{\infty} \frac{1}{u} \, du = \lim_{t \to \infty} \int_{\log 2}^{t} \frac{1}{u} \, du = \lim_{t \to \infty} [\log u]_{\log 2}^{t} = \infty \\]

Since the integral diverges, the integral test tells us that the series also diverges.

## Glossary

  * Infinite Series: an expression of the form \\( \sum_{n=1}^{\infty} a_n = a_1 + a_2 + a_3 + \dots \\), where \\( a_n \\) are the terms of the series.

  * Convergence of a series: an infinite series converges if its sequence of partial sums approaches a finite limit.

  * Divergence of a series: an infinite series diverges if its sequence of partial sums does not approach a finite limit (either it goes to infinity or oscillates).

  * Improper integral: a definite integral where at least one of the limits of integration is infinite, or the integrand has a discontinuity within the interval of integration.

  * Decreasing function: a function \\( f(x) \\) is decreasing on an interval if for any \\( x_1 < x_2 \\) in that interval, \\( f(x_1) \geq f(x_2) .\\)

  * Continuous function: a function whose graph can be drawn without lifting the pen, meaning there are no abrupt jumps or breaks.

  * Partial sum \\( s_k \\): the sum of the first \\( k \\) terms of an infinite series, denoted as \\( s_k = \sum_{n=1}^{k} a_n \\).

  * Limit of a sequence: the value that the terms of a sequence approach as the index tends to infinity.

  * Additivity of integrals: the property that the definite integral over a compound interval is the sum of the definite integrals over the disjoint subintervals that make up the compound interval.
