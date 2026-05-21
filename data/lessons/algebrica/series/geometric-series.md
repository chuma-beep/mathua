> Content sourced from [Algebrica](https://algebrica.org/geometric-series/) — CC BY-NC 4.0

## What is a geometric series

The **geometric series** is defined as the infinite sum:

\\[\sum_{n=0}^{\infty} r^n \\]

\\( r \\) is the common ratio, which defines how each term is obtained by multiplying the previous one by \\( r \\). If \\( r = 0 \\), the geometric series converges and its sum is \\( 1 \\), since all terms after the first are zero: \\[\sum_{n=0}^{\infty} 0^n = 1 + 0 + 0 + \dots = 1 \\]

If \\( r = 1 \\), the series diverges, because the sequence of partial sums grows without bound: \\[s_n = 1 + 1 + \dots + 1 = n + 1 \quad \rightarrow \quad \lim_{n \to \infty} s_n = \infty \\]

##### While a geometric series focuses on the sum of terms, it is built upon the underlying structure of a [geometric sequence](<../geometric-sequance>), an ordered list where each term is obtained by multiplying the previous one by a constant ratio.

The geometric series is fundamental because its behavior is well understood and precisely characterized. As a result, it serves as a reference model for analyzing more complex series. In particular, for [series with positive terms](<../series-with-positive-terms/>), the geometric series can be used in comparison tests to determine convergence or divergence.

##### If a given series can be bounded above or below by a geometric series, we can often deduce the nature of the series by analogy.


The \\( n \\)-th partial sum of a geometric [series](<../series>) with \\( r \ne 0 \\) and \\( r \ne 1 \\) is given by:

\\[S_n = \sum_{k=0}^{n} r^k = \frac{1 - r^{n+1}}{1 - r} \\]

This result is obtained starting from:

\\[s_n = 1 + r + r^2 + r^3 + \dots + r^n \\]

Multiplying both sides by \\( r \\), we get:

\\[\begin{align} r \cdot s_n &= r(1 + r + r^2 + \dots + r^n) \\\\[0.5em] &= r + r^2 + r^3 + \dots + r^n + r^{n+1} \\\ &= 1+ r + r^2 + r^3 + \dots + r^n + r^{n+1} -1 \\\ \end{align} \\]

In the last term of the equation, we added and subtracted \\( +1 \\). This way, the expression from \\( 1 + \dots + r^n \\) becomes \\( s_n \\). By substituting the values, we obtain:

\\[\begin{align} r \cdot s_n &= s_n + r^{n+1} - 1 \\\\[0.5em] s_n (1 - r) &= 1 - r^{n+1} \end{align} \\]

and finally:

\\[s_n = \frac{1 - r^{n+1}}{1 - r} \quad \text{for } r \ne 1 \\]

##### The condition \\( r \ne 1 \\) is essential; otherwise, the denominator becomes zero and the formula is undefined.


Let us now consider the case where \\( |r| < 1 \\). In this case, the geometric series converges, because the terms get smaller and smaller in absolute value. The infinite sum can be computed using the formula:

\\[\sum_{n=0}^{\infty} r^n = \frac{1}{1 - r} \\]

This result holds because \\( \lim_{n \to \infty} r^{n+1} = 0 \\) when \\( |r| < 1 \\), which simplifies the formula for the partial sum:

\\[s_n = \frac{1 - r^{n+1}}{1 - r} \quad \rightarrow \quad \lim_{n \to \infty} s_n = \frac{1}{1 - r} \\]

In general, when the index starts at a value \\( n > 0 \\), for example \\( \alpha \\), the formula for the sum of the geometric series becomes:

\\[\sum_{n=\alpha}^{\infty} r^n = \frac{r^{\alpha}}{1 - r} \\]


In the case where \\( r > 1 \\), the geometric series diverges. This happens because the terms \\( r^n \\) grow larger and larger as \\( n \to \infty \\), and so does the sequence of partial sums:

\\[\lim_{n \to \infty} s_n = \lim_{n \to \infty} \sum_{k=0}^{n} r^k = \infty \\]

Since the terms do not tend to zero and their sum grows without bound, the series does not have a finite limit and therefore diverges.


Finally, in the case where \\( r \leq -1 \\), the limit of \\( s_n \\) does not exist, and the series is indeterminate. This happens because the terms \\( r^n \\) alternate in sign and grow in magnitude, causing the partial sums \\( s_n \\) to oscillate without approaching any finite value. As a result, the series does not converge nor diverge to infinity, but instead has no limit at all. For example, if \\( r = -2 \\), the terms become:

\\[1 - 2 + 4 - 8 + 16 - \dots \\]

and the sequence of partial sums fluctuates increasingly without stabilizing, making the series divergent in an oscillatory and unbounded way.


Therefore, we can conclude that the series:

  * converges if \\( |r| < 1 \\) and its sum is \\( \dfrac{1}{1 - r} \\) when the index starts at a value \\(n = 0\\), or \\( \dfrac{r^\alpha}{1 - r}\\) when the index starts at a value \\(n > 0\\);
  * diverges to \\( +\infty \\) if \\( r \geq 1 \\);
  * is indeterminate if \\( r \leq -1 \\).


## Example

Let’s study the behavior of the following geometric series and compute its sum.

\\[\sum_{n=0}^{\infty} \frac{3^n + 4^n}{5^n}\\]


By the linearity of summation, we can rewrite the series as the sum of two geometric series:

\\[\sum_{n=0}^{\infty} \left( \frac{3}{5} \right)^n + \sum_{n=0}^{\infty} \left( \frac{4}{5} \right)^n\\]

We observe that the common ratio of both series is \\(< 1\\), therefore, by the properties of geometric series, both series converge. Let us now compute the sum. Let us recall that the sum of a geometric series is given by:

\\[\sum_{n=0}^{\infty} r^n = \frac{1}{1 - r} \\]

We have:

\\[\sum_{n=0}^{\infty} \left( \frac{3}{5} \right)^n = \frac{1}{1 - \frac{3}{5}} = \frac{1}{\frac{2}{5}} = \frac{5}{2}\\]

\\[\sum_{n=0}^{\infty} \left( \frac{4}{5} \right)^n = \frac{1}{1 - \frac{4}{5}} = \frac{1}{\frac{1}{5}} = 5\\]

By adding the two values, we obtain:

\\[\frac{5}{2} + 5 = \frac{15}{2}\\]

Therefore, the series converges and its sum is: \\[\frac{15}{2} \\]

##### Geometric series are among the easiest to work with, since their convergence is straightforward to determine, and their sum can be computed when they converge.

## Determine the nature of the following series and compute their sum.

  * \\[\text{1. } \quad \sum_{n=1}^{\infty} \left( \sqrt{3} \right)^n\\] [solution](<../tests/geometric-series/#s1>)

  * \\[\text{2. } \quad \sum_{n=1}^{\infty} \left( 3k + 2 \right)^n \\] [solution](<../tests/geometric-series/#s2>)

  * \\[\text{3. } \quad \sum_{n=1}^{\infty} \left( 1 - 2 \cos x \right)^n \\] [solution](<../tests/geometric-series/#s3>)


##### The proposed series are carefully designed to help you consolidate your understanding of infinite series. Try analyzing their behavior and computing their sums on your own before checking the provided solutions.

## Glossary

  * Geometric series: a sum of terms where each term is obtained by multiplying the previous one by a constant ratio.

  * Geometric sequence: an ordered list of numbers where each term after the first is found by multiplying the previous one by a fixed, non-zero number called the common ratio.

  * Common ratio \\( r \\): the constant value by which each term in a geometric series or sequence is multiplied to obtain the next term.

  * Comparison tests: methods used to determine the convergence or divergence of a series by comparing it to a known convergent or divergent series, such as a geometric series.

  * Oscillation: the behavior of a sequence of partial sums that fluctuates without settling on a specific value.
