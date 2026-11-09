# Bertrand's test

Bertrand's test refines Raabe's test at the boundary where the Raabe limit equals 1. It measures the next logarithmic correction to the ratio of consecutive terms.

## The Bertrand limit

For a series of positive terms define
\[
L=\lim_{n\to\infty}\ln n\left[n\left(1-\frac{a_{n+1}}{a_n}\right)-1\right].
\]
If \(L>1\) the series converges. If \(L<1\) it diverges. The test is designed for terms that shrink like \(1/n\) up to a slowly varying factor, which is where Raabe's test is silent.

The logarithmic factor turns a correction of size \(1/\ln n\) into a finite number. Without that factor the correction would vanish in the limit and carry no information.

## Logarithmic corrections

Take \(a_n=1/(n(\ln n)^p)\) with \(n\ge2\). Expanding the ratio gives
\[
\frac{a_{n+1}}{a_n}=1-\frac1n+\frac{p}{n\ln n}+O\left(\frac{1}{n^2}\right).
\]
After multiplying by \(n\) and subtracting 1, the remaining term is \(p/\ln n\). Multiplication by \(\ln n\) leaves \(p\), so \(L=p\).

This agrees with the known rule for the logarithmic p-series. The test converges for \(p>1\) and diverges for \(p\le1\).

## A constant factor

A constant multiplier in the term cancels from the ratio. If \(a_n=c/(n(\ln n)^p)\), the same expansion gives \(L=p\). The constant changes the size of every term and the value of the series, but it cannot change convergence.

When the logarithm appears in the numerator instead, the sign reverses. For \(a_n=(\ln n)^p/n\) the limit is \(-p\), which is below 1 for every positive \(p\) and correctly reports divergence.
