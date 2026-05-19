> Content sourced from [Algebrica](https://algebrica.org/series-with-positive-terms/) — CC BY-NC 4.0

## What is a series with positive terms?

A series with positive terms is one in which every term \\( a_k \\) satisfies \\( a_k > 0 \\) for all \\( k \in \mathbb{N} \\). As a result, the sequence of partial sums:

\\[S_n = \sum_{k=1}^{n} a_k \\]

is strictly increasing, since each new term contributes a positive quantity. This guarantees that the [series](<../series>) cannot oscillate or decrease. Such a series has only two possible behaviors: it either converges to a finite value if the partial sums are bounded above, or it diverges to infinity if they are not. It is never indeterminate or conditionally convergent. More formally, for a series of the form:

\\[\sum_{k=1}^{\infty} a_k \quad \text{with } a_k > 0 \\]

the series converges if and only if the sequence of partial sums \\( (S_n) \\) is bounded. This property makes positive-term series especially suitable for convergence tests such as the comparison test, the integral test, and the ratio test, all of which require non-negative terms.


The [harmonic series](<../harmonic-series/>) is an example of a positive-term series that diverges:

\\[\sum_{k=1}^{\infty} \frac{1}{k} \\]

In fact, despite the fact that the terms \\( \frac{1}{k} \\) tend to zero, the sequence of partial sums increases without limit.

In general, we refer to a constant-sign series when the terms of the sequence \\( {a_n} \\) all have the same sign for every \\( n \in \mathbb{N} \\); that is, they are either all positive or all negative.

## Comparison test

The **comparison test** is a method used to determine whether a series converges or diverges by comparing it to another series whose behavior is already known. It is particularly useful when dealing with series with positive terms, where direct evaluation of convergence is difficult.

Let \\( \sum a_k \\) and \\( \sum b_k \\) be two series with positive terms. Suppose there exists an [integer](<../integers/>) \\( N \in \mathbb{N} \\) such that:

\\[0 \leq a_k \leq b_k \quad \text{for all } k \geq N \\]

In other words, each term \\( a_k \\) is non-negative and is less than or equal to the corresponding term \\( b_k \\) (this condition is essential for applying the comparison test correctly). Then:

  * If \\( \sum b_k \\) converges, then \\( \sum a_k \\) also converges.
  * If \\( \sum a_k \\) diverges, then \\( \sum b_k \\) also diverges.


Let’s consider the partial sums of each series:

\\[S_n = \sum_{k=1}^{n} a_k, \quad T_n = \sum_{k=1}^{n} b_k \\]

Since both sequences \\( a_k \\) and \\( b_k \\) are made of non-negative terms, both \\( S_n \\) and \\( T_n \\) are non-decreasing. Now, because \\( a_k \leq b_k \\) for all \\( k \geq N \\), we have:

\\[S_n \leq T_n \quad \text{for all } n \geq N \\]

If the series \\( \sum b_k \\) converges, that means \\( T_n \\) has a finite limit — it is bounded above. Since \\( S_n \leq T_n \\), the sequence of partial sums \\( S_n \\) is also bounded above. And since \\( S_n \\) is non-decreasing and bounded, it must converge. Therefore, \\( \sum a_k \\) also converges.

If \\( \sum a_k \\) diverges, then \\( S_n \to \infty \\). But since \\( S_n \leq T_n \\), the only way this inequality can hold is if \\( T_n \\) also grows without bound. Therefore, \\( \sum b_k \\) also diverges.

## Example

Using the comparison test, we determine the nature of the following series:

\\[\sum_{n=1}^{\infty} \frac{1}{n^2 + 2n + 1} \\]

The series is a positive-term series, since the denominator is a polynomial and all terms are positive. Therefore, its nature can be determined using the comparison test.

##### It is important to verify this point, since the comparison test is only valid for series in which all terms are positive.


First, we check whether the necessary condition for convergence is satisfied:

\\[\lim_{n \to +\infty} \frac{1}{n^2 + 2n + 1} = 0 \\]

It is evident that the [limit](<../limits>) has the form \\( \frac{1}{\infty} \\), which implies that it equals zero. The necessary condition for convergence is therefore satisfied. At this point, using the comparison test, we can state that:

\\[\frac{1}{n^2 + 2n + 1} < \frac{1}{n^2} \\]

since the denominator in the first expression is greater than the denominator in the second. Let

\\[a_n = \frac{1}{n^2 + 2n + 1} \quad \text{and} \quad b_n = \frac{1}{n^2} \\]

using the comparison test, we observe that

\\[a_n < b_n \\]

The series

\\[\sum_{n=1}^{\infty} b_n = \sum_{n=1}^{\infty} \frac{1}{n^2} \\]

is a [generalized harmonic series](<../harmonic-series>), which is known to converge when the exponent in the denominator satisfies the condition \\( p > 1 \\).

Hence, by the comparison test, since the series \\( \sum b_n \\) converges, the series \\( \sum a_n \\) also converges.

##### Determining the nature of a positive-term series using the comparison test is relatively straightforward, but it requires plenty of practice to choose the right comparison and justify the inequality correctly.

## Glossary

  * Series with positive terms: a series where every term \\( a_k \\) is greater than zero for all indices \\( k \\).

  * Sequence of partial sums \\( S_n \\): the sequence formed by the sum of the first \\( n \\) terms of a series, \\( S_n = \sum_{k=1}^{n} a_k \\).

  * Strictly increasing sequence: a sequence where each term is greater than the previous term.

  * Bounded above: a sequence is bounded above if there exists a number \\( M \\) such that every term in the sequence is less than or equal to \\( M \\).

  * Constant-sign series: a series where all terms have the same sign (either all positive or all negative).

  * Comparison test: a method used to determine the convergence or divergence of a series by comparing it term-by-term to another series whose behavior is already known.

  * Harmonic series: The series \\( \sum_{k=1}^{\infty} \frac{1}{k} \\), which is a known example of a positive-term series that diverges.

  * Generalized harmonic series: A series of the form \\( \sum_{n=1}^{\infty} \frac{1}{n^p} \\), which converges if \\( p > 1 \\) and diverges if \\( p \leq 1 \\).

  * Necessary condition for convergence: For a series \\( \sum a_k \\) to converge, it is necessary that \\( \lim_{k \to \infty} a_k = 0 \\). However, this condition is not sufficient for convergence.


Series

Series are infinite sums of sequence terms, typically indexed by natural numbers.

1.3k

[Series](https://algebrica.org/series/)

2.4k

[Cauchy’s Convergence Criterion for Series](https://algebrica.org/cauchy-convergence-criterion-series/)

1.7k

[Harmonic Series](https://algebrica.org/harmonic-series/)

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
