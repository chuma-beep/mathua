# Cauchy condensation test

The Cauchy condensation test replaces a series of positive terms by a sparse series. The sparse series is often geometric, so a convergence question becomes an algebra question about one ratio.

## Condensing the series

Let \(a_n\) be nonnegative and nonincreasing. The original series and the condensed series share the same convergence behaviour:
\[
\sum_{n=1}^{\infty} a_n \quad \text{and} \quad \sum_{n=0}^{\infty} 2^n a_{2^n}.
\]
The condensed term is \(b_n=2^n a_{2^n}\). Each block of terms between \(2^n\) and \(2^{n+1}-1\) is bounded above and below by a multiple of \(b_n\), which is why the two series cannot disagree.

If a problem gives \(a_n\), the first step is to substitute \(2^n\) for \(n\) and simplify. The exponent of the condensed term is usually the whole question.

## The condensed geometric series

For the p-series \(a_n=1/n^p\) the condensed term is
\[
b_n=2^n\cdot\frac{1}{(2^n)^p}=2^{n(1-p)}.
\]
This is geometric with ratio \(2^{1-p}\). When \(p>1\) the ratio is below 1 and the original series converges. When \(p\le1\) the ratio is at least 1 and the series diverges.

The denominator of the ratio is a power of two. Writing the ratio as \(1/m\) makes \(m=2^{p-1}\), which is an integer for every integer \(p\ge2\).

## Summing the condensed series

When the condensed series is geometric, its sum follows from the first term and the ratio. The first term is \(b_0=a_1\). For \(a_n=1/n^p\) this term is 1 and the sum is
\[
S=\frac{1}{1-2^{1-p}}=\frac{2^{p-1}}{2^{p-1}-1}.
\]
A larger exponent makes the denominator closer to the numerator, so the sum approaches 1 from above. The condensed sum is not the sum of the original series. It is a comparison series, and only its convergence or divergence is shared.
