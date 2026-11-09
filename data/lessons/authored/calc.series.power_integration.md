# Integrating power series

A power series may be integrated term by term inside its interval of convergence. The resulting series has the same radius of convergence, and its coefficients are the old coefficients divided by the new exponent.

## Term-by-term integration

If \(f(x)=\sum a_n x^n\) on an interval around zero, then
\[
\int_0^x f(t)\,dt=\sum_{n=0}^{\infty}\frac{a_n}{n+1}x^{n+1}.
\]
The coefficient of \(x^{n+1}\) is \(a_n/(n+1)\). An integration question therefore reduces to dividing one coefficient by one integer.

The constant of integration is fixed by the lower limit. Starting at zero removes the constant and keeps the answer a power series.

## The logarithmic series

The geometric series is
\[
\sum_{n=0}^{\infty}x^n=\frac{1}{1-x}.
\]
Integrating from zero gives
\[
\sum_{n=0}^{\infty}\frac{x^{n+1}}{n+1}=-\ln(1-x).
\]
The coefficient of \(x^m\) is \(1/m\) for every \(m\ge1\). This identity is the source of many logarithm series and of the alternating harmonic sum at \(x=-1\).

The radius of convergence is still 1. At the endpoint \(x=1\) the series diverges, while at \(x=-1\) it converges to \(\ln2\).

## The inverse hyperbolic series

Integrating the even geometric series gives another useful family. Starting from
\[
\sum_{n=0}^{\infty}x^{2n}=\frac{1}{1-x^2},
\]
term-by-term integration produces
\[
\sum_{n=0}^{\infty}\frac{x^{2n+1}}{2n+1}=\operatorname{artanh}x.
\]
The coefficient of \(x^{2k+1}\) is \(1/(2k+1)\). Only odd powers appear, and each coefficient is the reciprocal of the exponent.
