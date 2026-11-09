# Kummer's test

Kummer's test compares a series against a chosen positive sequence. It contains the ratio test and Raabe's test as special cases and makes the comparison explicit.

## Kummer's comparison

Choose a positive sequence \(b_n\) whose reciprocal series diverges. For a series of positive terms \(a_n\) define
\[
L=\lim_{n\to\infty}\left(b_n\frac{a_n}{a_{n+1}}-b_{n+1}\right).
\]
If \(L>0\) the series converges. If \(L<0\) it diverges. The comparison sequence is a tool, not part of the answer, so a well-chosen \(b_n\) turns a hard limit into a short expansion.

The term \(b_n a_n/a_{n+1}\) is designed to cancel the leading growth of \(b_{n+1}\). What remains is the first correction, and its sign decides the test.

## Linear comparison sequences

Take \(b_n=cn\) and \(a_n=1/n^p\). Then
\[
b_n\frac{a_n}{a_{n+1}}=cn\left(1+\frac1n\right)^p=cn+cp+O\left(\frac1n\right).
\]
Subtracting \(b_{n+1}=c(n+1)\) leaves \(c(p-1)\). Thus \(L=c(p-1)\), which is positive exactly when \(p>1\).

A logarithmic factor in \(b_n\) can be used as well. With \(b_n=cn\ln n\) and \(a_n=1/(n(\ln n)^p)\), the same cancellation leaves \(c(p-1)\).

## Rational-term comparison

A term with an extra factor \(n+1\) in the denominator behaves like \(1/n^{p+1}\). With \(b_n=cn\),
\[
b_n\frac{a_n}{a_{n+1}}=c\left(1+\frac1n\right)^{p-1}(n+2).
\]
Subtracting \(c(n+1)\) leaves \(cp\), so the limit is \(cp\). The shift raises the effective exponent by one, and Kummer's test reports that increase directly.
