> Content sourced from [Algebrica](https://algebrica.org/root-test-for-series-convergence/) — CC BY-NC 4.0

## What is the root test

The root test is a method used to determine whether an infinite [series](<../series>) converges or diverges. It is particularly useful when each term of the series involves an expression raised to the \\( n \\)-th [power](<../powers>), such as [exponentials](<../exponential-function>) or [roots](<../radicals>). Suppose we have a series with positive terms of the form:

\\[\sum_{n=1}^{+\infty} a_n \\]

Assume that the [limit](<../limits>) exists and is finite:

\\[\limsup_{n \to \infty} \sqrt[n]{|a_n|} = L \\]

We use \\( \limsup \\) because it is always defined for real sequences that are bounded from below. This makes the root test reliable even when the usual limit \\( \lim_{n \to \infty} \sqrt[n]{|a_n|} \\) does not exist. By taking the limit superior, the test can still detect the long-term behavior of the sequence and determine whether the series converges or diverges.

If the limit \\( L \\) exists, three cases can occur:

  * If \\( L < 1 \\), the series converges absolutely.
  * If \\( L > 1 \\) or \\( L = \infty \\), the series diverges.
  * If \\( L = 1 \\), the test is inconclusive.


##### We use the [absolute value](<../absolute-value>) \\( |a_n| \\) to apply the root test to the absolute convergence of the series. This ensures the test works even if the terms \\( a_n \\) are negative or alternate in sign, allowing us to focus only on their magnitude.

## How to recognize when to apply the root test

The root test is especially useful when the general term of a series is expressed in the form \\( a_n = (b_n)^n \\), that is, when the entire term is raised to the power of \\( n \\). In such cases, taking the \\( n \\)-th root of \\( a_n \\) simplifies the expression significantly and often leads directly to a manageable limit.

By contrast, other tests may become more complicated in this context, especially when the ratio \\( \frac{a_{n+1}}{a_n} \\) does not simplify easily or when factorials or [exponential](<../exponential-function>) terms are involved in a way that makes limits more difficult to compute.

In summary, the Root Test is particularly effective in the following situations:

  * When \\( a_n = (f(n))^n \\), with \\( f(n) > 0 \\)
  * When the terms involve exponential-like growth or decay
  * When the limit \\( \sqrt[n]{|a_n|} \\) is easier to compute than \\( \frac{a_{n+1}}{a_n} \\)


In all other cases, if the general term is not raised to the \\( n \\)-th power other tests may be more suitable.

## Proof

Let us consider the case of absolute convergence where \\( L < 1 \\). Since \\( L < 1 \\), there exists a [real number](<../types-of-numbers>) \\( r \\) such that:

\\[L < r < 1 \\]

By the definition of \\( \limsup \\), there exists an [integer](<../integers/>) \\( N \\) such that for all \\( n \geq N \\):

\\[\sqrt[n]{|a_n|} < r \quad \rightarrow \quad |a_n| < r^n \\]

So the tail of the series satisfies:

\\[\sum_{n=N}^{\infty} |a_n| < \sum_{n=N}^{\infty} r^n \\]

Since \\( 0 < r < 1 \\), the [geometric series](<../geometric-series>) \\( \sum r^n \\) converges. Therefore, by the [comparison test](<../series-with-positive-terms/>), the tail \\( \sum_{n=N}^{\infty} |a_n| \\) converges, and so does the entire series \\( \sum |a_n| \\).


Let us now consider the case in which the series diverges, assuming \\( L > 1 \\). Then, by definition of \\( \limsup \\), for infinitely many indices \\( n \\), we have:

\\[\sqrt[n]{|a_n|} > r \quad \text{for some } r > 1 \\]

Thus, \\( |a_n| > r^n \\) infinitely often. But since \\( r^n \to \infty \\), we know that \\( |a_n| \nrightarrow 0 \\).

Therefore, the necessary condition for convergence of \\( \sum a_n \\) fails. So, the series diverges.


Let us consider the final case, where nothing can be concluded about the convergence of the series, namely, when \\( L = 1 \\). If \\( L = 1 \\), then for every \\( \varepsilon > 0 \\), we can find infinitely many \\( n \\) such that:

\\[\sqrt[n]{|a_n|} > 1 - \varepsilon \quad \text{and} \quad \sqrt[n]{|a_n|} < 1 + \varepsilon \\]

This range includes both convergent and divergent behaviors. For example, the [harmonic series](<../harmonic-series>) \\( a_n = \frac{1}{n} \\) has \\( \sqrt[n]{|a_n|} \to 1 \\) and diverges. The [p-series](<../harmonic-series>) ( a_n = \frac{1}{n^2} ) also has \\( \sqrt[n]{|a_n|} \to 1 \\), but it converges. So, the test is inconclusive when \\( L = 1 \\).

## Example

Consider the series:

\\[\sum_{n=1}^{\infty} \left( \frac{3n}{5n + 2} \right)^n \\]

We want to determine whether this series converges or diverges. In this case, we choose to apply the root test because each term of the series is given in the form \\( a_n = (\text{expression})^n \\). This structure makes the root test especially effective and simpler to use than other methods.


Let us define the [sequence](<../sequences>):

\\[a_n = \left( \frac{3n}{5n + 2} \right)^n \\]

To apply the root test, we evaluate the limit superior of the \\( n \\)-th root of \\( a_n \\):

\\[\limsup_{n \to \infty} \sqrt[n]{a_n} = \lim_{n \to \infty} \left( \frac{3n}{5n + 2} \right) \\]

Since the terms are positive, we omit the absolute value. Now we simplify the expression:

\\[\frac{3n}{5n + 2} = \frac{3}{5 + \frac{2}{n}} \longrightarrow \frac{3}{5} \quad \text{as } n \to \infty \\]

Therefore, we find:

\\[\limsup_{n \to \infty} \sqrt[n]{a_n} = \frac{3}{5} < 1 \\]

Since the limit is less than 1, the Root Test tells us that the series converges absolutely.

## Glossary

  * Infinite series: the sum of an infinite sequence of numbers, typically represented as \\( \sum_{n=1}^{+\infty} a_n.\\)

  * [Convergence of a series](<../series>): an infinite series converges if the sequence of its partial sums approaches a finite limit.

  * Absolute convergence: a series \\( \sum a_n \\) converges absolutely if the series of the absolute values of its terms, \\( \sum |a_n| \\), converges. Absolute convergence implies convergence.

  * Root test: a method for determining the convergence or divergence of an infinite series by analyzing the limit of the \\( n \\)-th root of the absolute value of its terms.

  * Limit superior: for a sequence, the largest limit point of the sequence. It is always defined for bounded sequences and provides a way to analyze the long-term behavior even if a standard limit does not exist.

  * General term \\( a_n \\): the formula or expression that defines the \\( n \\)-th term of a sequence or series.

  * [Geometric series](<../geometric-series>): a series of the form \\( \sum_{n=0}^{\infty} ar^n \\), which converges if \\( |r| < 1 \\) and diverges if \\( |r| \geq 1.\\)

  * [Comparison test](<../series-with-positive-terms/>): a test for the convergence or divergence of a series by comparing it to another series whose convergence or divergence is known.
