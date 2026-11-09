# Gauss's test

Gauss's test reads the coefficient of \(1/n\) in the ratio of consecutive terms. It is a systematic way to handle terms built from gamma functions and rational factors.

## The Gauss expansion

Suppose the terms are positive and the ratio has the form
\[
\frac{a_n}{a_{n+1}}=1+\frac{h}{n}+\frac{B(n)}{n^r},
\]
where \(r>1\) and \(B(n)\) stays bounded. If \(h>1\) the series converges. If \(h\le1\) it diverges. The remainder is small enough that the single number \(h\) decides the question.

The test is a refined ratio test. A ratio that tends to 1 too slowly for the ratio test is still described by the size of its first correction.

## Gamma ratios

For \(a_n=\Gamma(n+\alpha)/\Gamma(n+\beta)\),
\[
\frac{a_n}{a_{n+1}}=\frac{n+\beta}{n+\alpha}=1+\frac{\beta-\alpha}{n+\alpha}.
\]
Expanding the denominator gives \(1+(\beta-\alpha)/n+O(1/n^2)\). Hence \(h=\beta-\alpha\).

The gamma functions need not be evaluated. Only their shift parameters matter, and the difference of those parameters is the coefficient the test asks for.

## Rational terms

A rational expression with denominator degree larger than numerator degree also fits the expansion. For
\[
a_n=\frac{n+\alpha}{(n+1)^d},
\]
the numerator contributes \(-1/n\) to the logarithm of the ratio and the denominator contributes \(d/n\). The coefficient of \(1/n\) is therefore \(d-1\).

This family shows why the test is stated with a general remainder. The remainder is smaller than \(1/n\), so it cannot change \(h\) even when the rational expression has several factors.
