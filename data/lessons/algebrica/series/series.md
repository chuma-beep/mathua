> Content sourced from [Algebrica](https://algebrica.org/series/) — CC BY-NC 4.0

## What is a series

The concept of a **series** is closely tied to infinite [sequences](<../sequences>) of [real numbers](<../types-of-numbers>), with the main goal of studying the behavior of their infinite sum. This involves determining whether the sum approaches a finite [limit](<../limits>) (convergence) or grows without bound (divergence). From a formal point of view, let \\(\lbrace a_n \rbrace_{n \in \mathbb{N}}\\) be a sequence of real numbers. We define the partial sums of the sequence as follows:

\\[\begin{align} s_1 &= a_1\\\\[0.5em] s_2 &= a_1 + a_2 \\\\[0.5em] s_3 &= a_1 + a_2 + a_3 \\\\[0.5em] \vdots \\\\[0.5em] s_n &= a_1 + a_2 + \cdots + a_n \end{align} \\]

The sequence \\({s_n}\\), with \\(n \geq 1 \\), is called the sequence of partial sums, and each term is given by:

\\[s_n = \sum_{k=1}^{n} a_k \\]

A series with general term \\(a_k\\) is defined as the formal expression:

\\[\sum_{k=1}^{\infty} a_k \\]

## Nature of a series

The nature of a series is determined by analyzing the limit of its sequence of partial sums:

\\[\lim_{n \to \infty} s_n = \lim_{n \to \infty} \sum_{k=1}^{n} a_k \\]

  * If the limit exists and is finite, the series \\(\sum a_n\\) is said to converge. In this case, the value of the limi is called the sum of the series.
  * If the limit is \\(+\infty\\) or \\(-\infty\\), the series \\(\sum a_n\\) is said to diverge.
  * In all other cases, such as when the limit does not exist or oscillates, the series \\(\sum a_n\\) is considered indeterminate.


A series \\(\sum a_k\\) is said to **converge absolutely** if the series of [absolute values](<../absolute-value>) \\(\sum |a_k|\\) also converges. That is, the convergence is not affected by the signs of the terms.

##### The expression “sum of a series” is a conventional term: since infinitely many terms are involved, it is not a finite sum in the traditional sense, but the limit of the sequence of partial sums.

Altering a finite number of terms in the series does not affect its convergence or divergence. That is, two series that differ only by finitely many terms share the same convergence nature. However, they do not necessarily have the same sum.

## Necessary condition for convergence

Suppose that the series \\(\sum_{n=1}^{\infty} a_n\\) is convergent. As discussed above, this means that the sequence of partial sums converges to a finite limit. A necessary condition for convergence is:

\\[\lim_{n \to \infty} a_n = 0 \\]

In other words, the general term of the series must tend to zero. However, this condition is not sufficient: the fact that \\(a_n \to 0\\) does not guarantee that the series converges.


This property can be proven by analyzing the sequence of partial sums \\({s_n}\\), where:

\\[s_n = \sum_{k=1}^{n} a_k \\]

By definition, if the series converges, this sequence must tend to a finite limit. Let \\(S\\) be the sum of the series; then:

\\[S = \lim_{n \to \infty} s_n = \lim_{n \to \infty} s_{n-1} \\]

Since both \\(s_n\\) and \\(s_{n-1}\\) converge to the same limit \\(S\\), the difference between consecutive partial sums must tend to zero:

\\[\lim_{n \to \infty} (s_n - s_{n-1}) = \lim_{n \to \infty} a_n = S - S = 0 \\]

## Linear properties of series

Let \\(\sum_{k=1}^{\infty} a_k\\) be a convergent series, and let \\(\lambda \in \mathbb{R}\\), \\(\lambda \ne 0\\). Then the series \\(\sum_{k=1}^{\infty} \lambda a_k\\) also converges, and its sum is:

\\[\sum_{k=1}^{\infty} \lambda a_k = \lambda \sum_{k=0}^{\infty} a_k \\]


If both series \\(\sum_{k=1}^{\infty} a_k\\) and \\(\sum_{k=1}^{\infty} b_k\\) converge, then their term-by-term sum also defines a convergent series:

\\[\sum_{k=1}^{\infty} (a_k + b_k) \\]

Moreover, the sum of the resulting series is equal to the sum of the individual series:

\\[\sum_{k=1}^{\infty} (a_k + b_k) = \sum_{k=1}^{\infty} a_k + \sum_{k=1}^{\infty} b_k \\]

## Well-known series

Consider the following series:

\\[\sum_{n=1}^{\infty} \frac{1}{n} = 1 + \frac{1}{2} + \frac{1}{3} + \cdots + \frac{1}{n} + \cdots \\]

The series is called the harmonic series, and it is divergent.


The following series is called the [generalized harmonic series](<../harmonic-series>):

\\[\sum_{n=1}^{\infty} \frac{1}{n^p} = 1 + \frac{1}{2^p} + \frac{1}{3^p} + \cdots + \frac{1}{n^p} + \cdots, \quad p \in \mathbb{R} \\]

The convergence of the series depends on the value of \\(p\\):

  * If \\(p > 1\\) the series converges.
  * If \\(p \leq 1\\) the series diverges.


Consider the [geometric series](<../geometric-series/>) of ratio \\(q\\):

\\[\sum_{n=0}^{\infty} q^n = 1 + q + q^2 + q^3 + \cdots + q^n + \cdots \\]

The convergence of the series depends on the value of \\(q\\):

  * If \\(-1 < q < 1\\), the series converges.
  * If \\(q \geq 1\\), the series diverges.
  * If \\(q \leq -1\\), the series is irregular (diverges or oscillates).


Consider the following series, known as a telescoping series:

\\[\sum_{n=1}^{\infty} \frac{1}{n(n+1)} = \frac{1}{2} + \frac{1}{6} + \cdots + \frac{1}{n(n+1)} + \cdots = 1 \\]

This series is convergent.

## Glossary

  * Series: the limit of the sum of the terms of a sequence.

  * Sequence: an ordered list of numbers.

  * Partial sums \\( s_n \\): the sum of the first \\( n \\) terms of a sequence.

  * Convergence: a series converges if the limit of its sequence of partial sums is a finite number.

  * Converge absolutely: a series \\( \sum a_k \\) converges absolutely if the series of absolute values \\( \sum |a_k| \\) converges.

  * Divergence: a series diverges if the limit of its sequence of partial sums is \\( +\infty \\) or \\( -\infty \\).

  * Indeterminate: a series is indeterminate if the limit of its sequence of partial sums does not exist or oscillates.

  * Sum of the series: the finite limit of the sequence of partial sums for a convergent series.

  * General term \\( a_n \\) or \\( a_k \\): the formula or expression that defines each term of a sequence or series.
