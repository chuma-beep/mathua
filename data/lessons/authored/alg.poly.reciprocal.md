# Reciprocal polynomials

A reciprocal polynomial reverses the coefficients of a polynomial. If \(P(x) = a_4x^4 + a_3x^3 + a_2x^2 + a_1x + a_0\), then its reciprocal polynomial is
\[
Q(x) = x^4P(1/x) = a_0x^4 + a_1x^3 + a_2x^2 + a_3x + a_4.
\]
The degree stays the same when the constant term of \(P\) is nonzero.

## Forming the reciprocal polynomial

Read the coefficients of \(P\) from the highest power down, then write them from the lowest power up. The leading coefficient of \(P\) becomes the constant term of \(Q\), and the constant term of \(P\) becomes the leading coefficient of \(Q\).

For example, if
\[
P(x) = 2x^4 - 3x^3 + 5x^2 + x - 4,
\]
then
\[
Q(x) = -4x^4 + x^3 + 5x^2 - 3x + 2.
\]
Every coefficient keeps its sign; only the order changes.

## Roots and reciprocals

If \(r\) is a nonzero root of \(P\), then \(1/r\) is a root of \(Q\). Substitute \(x = 1/r\):
\[
Q(1/r) = (1/r)^4 P(r) = 0.
\]
The roots of \(Q\) are therefore the reciprocals of the roots of \(P\), provided no root is zero. A root at zero would lower the degree of \(Q\), which is why the constant term must be nonzero.

## Coefficients and Vieta

Vieta's formulas connect the reciprocal roots to the coefficients of \(P\). For a quartic,
\[
\sum_i \frac{1}{r_i} = -\frac{a_1}{a_0}
\]
and
\[
\prod_i \frac{1}{r_i} = \frac{a_4}{a_0}.
\]
The product formula has no extra sign because the degree is even. These two identities let you answer questions about reciprocal roots without factoring the polynomial.
