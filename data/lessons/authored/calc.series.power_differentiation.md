# Differentiating power series

A power series may be differentiated term by term inside its interval of convergence. Differentiation multiplies each coefficient by its exponent and lowers the exponent by one.

## Term-by-term differentiation

If \(f(x)=\sum a_n x^n\), then
\[
f'(x)=\sum_{n=1}^{\infty}n a_n x^{n-1}.
\]
The coefficient of \(x^{n-1}\) is \(n a_n\). The radius of convergence does not change, although the behaviour at the endpoints may.

A differentiation question usually asks for one coefficient. Read the exponent, multiply by the matching coefficient and stop.

## The geometric series

Differentiating the geometric series gives
\[
\frac{1}{(1-x)^2}=\sum_{n=1}^{\infty}n x^{n-1}.
\]
The coefficient of \(x^k\) is \(k+1\). The coefficients grow linearly, which is why the radius of convergence is still 1 and why the series diverges at both endpoints.

A second differentiation produces the coefficients \(n(n-1)\) and the series for \(2/(1-x)^3\). Each differentiation raises the polynomial degree of the coefficient.

## The even series

The even geometric series is
\[
\sum_{n=0}^{\infty}x^{2n}=\frac{1}{1-x^2}.
\]
Differentiating gives
\[
\sum_{n=1}^{\infty}2n x^{2n-1}=\frac{2x}{(1-x^2)^2}.
\]
Only odd powers appear. The coefficient of \(x^{2k-1}\) is \(2k\), so the coefficients are the positive even numbers in order.

The missing even powers are not an accident. The original series had only even powers, and differentiation lowers every exponent by one.
