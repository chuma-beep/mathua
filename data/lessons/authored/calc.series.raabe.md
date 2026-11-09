# Raabe's test

Raabe's test measures how far the ratio of consecutive terms sits from 1. It is sharper than the ratio test when the ratio tends to 1, which is exactly the case the ratio test cannot decide.

## The Raabe limit

For a series of positive terms define
\[
L=\lim_{n\to\infty} n\left(1-\frac{a_{n+1}}{a_n}\right).
\]
If \(L>1\) the series converges. If \(L<1\) it diverges. If \(L=1\) the test gives no information.

The expression inside the limit compares the ratio with 1 at the scale \(1/n\). A ratio that behaves as \(1-c/n\) produces the limit \(c\). The sign convention matters: a positive limit above 1 means the terms shrink slightly faster than the harmonic terms.

## Power-law terms

For \(a_n=1/n^p\),
\[
\frac{a_{n+1}}{a_n}=\left(\frac{n}{n+1}\right)^p=\left(1+\frac1n\right)^{-p}.
\]
The binomial expansion gives \(1-p/n+O(1/n^2)\). Multiplying the difference from 1 by \(n\) leaves \(p\), so \(L=p\) and the test recovers the p-series rule.

The same calculation works when the numerator is a constant multiple. A constant factor cancels in the ratio and does not change \(L\).

## Shifted power terms

A shift in the denominator changes the limit by the size of the shift. For \(a_n=1/(n^p(n+1))\),
\[
\frac{a_{n+1}}{a_n}=\left(1+\frac1n\right)^p\frac{n+1}{n+2}.
\]
The first factor contributes \(p/n\) to the expansion and the second contributes \(1/n\). Hence \(L=p+1\).

For a numerator that grows like \((n+1)^q\), the numerator contributes \(q/n\) with the opposite sign, so \(L=p-q\). The test therefore reads the difference between the growth rates of the numerator and the denominator.
