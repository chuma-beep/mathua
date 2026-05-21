> Content sourced from [Algebrica](https://algebrica.org/leibniz-criterion/) — CC BY-NC 4.0

## What is Leibniz’s criterion used for

Leibniz’s criterion is used to study the convergence of alternating [series](<../series>), those composed of an infinite sequence of positive and negative terms that alternate in sign. Consider the alternating series \\[\sum_{n=1}^{+\infty} (-1)^n a_n \\]

where \\( a_n \geq 0 \forall n \in \mathbb{N} \\). Leibniz’s criterion states that the series converges if the following conditions are satisfied:

  * The [sequence](<../sequence>) \\( {a_n} \\) is infinitesimal, that is the [limit](<../limits>) exists and is finite: \\[\lim_{n \to +\infty} a_n = 0 \\]

  * \\( {a_n} \\) is eventually non-increasing, that is there exists an index \\( n_0 \\) such that for every \\( n \geq n_0 \\), it holds that: \\[a_{n+1} \leq a_n \\]


In practice, Leibniz’s criterion allows one to immediately establish the convergence of a series based solely on the verification of the above conditions. However, the test is not generalizable and applies only to alternating series.


If the series

\\[\sum_{n=1}^{+\infty} (-1)^k a_k \\]

is convergent and \\( S \\) is its sum, then for every \\( n \in \mathbb{N} \\), the error made by truncating the series at the \\( n \\)-th term is at most equal to the next term \\( a_{n+1} \\):

\\[\left| S - \sum_{k=1}^{n} (-1)^k a_k \right| \leq a_{n+1}. \\]

This means that if you are summing the alternating [harmonic series](<../harmonic-series>):

\\[\sum_{n=1}^{+\infty} \frac{(-1)^n}{n} \\]

and you stop at the term \\( n = 10 \\), then the maximum error you are making is:

\\[\left| S - \sum_{k=1}^{10} \frac{(-1)^k}{k} \right| \leq \frac{1}{11} \\]

Leibniz’s criterion only tells us whether an alternating series converges; it provides no information about the actual sum of the series. To compute the sum, one must use other tools such as Taylor or Maclaurin series expansions, when available.

## Proof

Let \\( {a_n} \\) be a sequence of [positive real numbers](<../series-with-positive-terms/>) such that \\( a_n \geq a_{n+1} \forall \, n \\), and \\( \lim_{n \to \infty} a_n = 0 \\). We consider the alternating series:

\\[\sum_{n=1}^{\infty} (-1)^{n+1} a_n \\]

and aim to prove that it converges.


To do this, we examine the sequence of partial sums:

\\[S_n = \sum_{k=1}^{n} (-1)^{k+1} a_k \\]

We analyze this sequence by distinguishing between the even and odd partial sums. When \\( n \\) is even, say \\( n = 2m \\), the sum becomes:

\\[S_{2m} = a_1 - a_2 + a_3 - a_4 + \cdots + a_{2m-1} - a_{2m}. \\]

Grouping the terms in pairs gives:

\\[S_{2m} = (a_1 - a_2) + (a_3 - a_4) + \cdots + (a_{2m-1} - a_{2m}). \\]

Since the sequence \\( {a_n} \\) is decreasing, each difference \\( a_{k} - a_{k+1} \\) is non-negative. Therefore, each term in the grouped sum is positive, and the sequence \\( \lbrace S_{2m} \rbrace \\) is increasing. Additionally, since each \\( a_n > 0 \\), the total sum is bounded above by \\( a_1 \\). Thus the subsequence \\( \lbrace S_{2m} \rbrace \\) converges.


Next, we observe that the odd partial sums can be written as

\\[S_{2m+1} = S_{2m} + a_{2m+1}. \\]

Since \\( \lim_{n \to \infty} a_n = 0 \\), the difference between \\( S_{2m+1} \\) and \\( S_{2m} \\) tends to zero as \\( m \to \infty \\). Therefore, both subsequences \\( \lbrace S_{2m} \rbrace \\) and \\( \lbrace S_{2m+1} \rbrace \\) converge to the same limit.

It follows that the full sequence \\( \lbrace S_n \rbrace \\) of partial sums converges, and hence the alternating series is convergent.

## Example

Let us consider the following alternating series:

\\[\sum_{n=1}^{+\infty} \frac{(-1)^n}{n!} \\]

We define \\( a_n = \frac{1}{n!} \\), in this way we focus on analyzing the behavior of the [absolute values](<../absolute-value>) of the terms, which is essential when applying Leibniz’s criterion.


To apply Leibniz’s criterion, we need to verify three conditions. First, the terms \\( a_n \\) are all positive for every \\( n \in \mathbb{N} \\). Second, the sequence \\( {a_n} \\) is decreasing. In fact, since the [factorial](<../factorial>) function grows rapidly, we have:

\\[a_{n+1} = \frac{1}{(n+1)!} < \frac{1}{n!} = a_n \\]

This confirms that the sequence is strictly decreasing. Third, the limit of the general term is zero:

\\[\lim_{n \to +\infty} \frac{1}{n!} = 0 \\]

Since all three conditions are satisfied, Leibniz’s criterion ensures that the series

\\[\sum_{n=1}^{+\infty} \frac{(-1)^n}{n!} \\]

converges.

## Determine the nature of the following series.

  * \\[\text{1. } \quad \sum_{n=3}^{\infty} (-1)^n \frac{1}{\log n}\\] [solution](<#>)

  * \\[\text{2. } \quad \sum_{n=1}^{\infty} (-1)^n \sin\left(\frac{1}{n}\right) \\] [solution](<#>)

  * \\[\text{3. } \quad \sum_{n=1}^{\infty} (-1)^n \frac{\log n}{n e^n} \\] [solution](<#>)


##### The proposed alternating series are selected to help you strengthen your understanding of convergence using Leibniz’s criterion. Try analyzing each series to determine whether the conditions are satisfied before checking the provided solutions.

## Glossary

  * Alternating series: a series in which the terms alternate in sign.

  * Convergence: the property of an infinite series whose partial sums approach a finite limit as the number of terms increases indefinitely.

  * Leibniz’s Criterion: a test used to determine the convergence of alternating series based on the properties of the sequence of the absolute values of the terms.

  * Infinitesimal sequence: a sequence whose limit as the index approaches infinity is zero.

  * Sum of a series: the finite limit that the partial sums of a convergent series approach. Partial Sum: The sum of the first n terms of an infinite series.
