# Integrals of hyperbolic functions

The hyperbolic functions have elementary antiderivatives that follow directly from their derivatives.

## The basic antiderivatives

Reading the derivative table backwards gives \(\int\sinh x\,dx = \cosh x + C\),
\(\int\cosh x\,dx = \sinh x + C\) and \(\int \operatorname{sech}^{2}x\,dx = \tanh x + C\). At
\(x = 1\) the antiderivative \(\sinh x\) equals 1.1752, while its derivative there,
\(\cosh 1 = 1.5431\), is the integrand. No sign appears in these integrals, because the derivative
of \(\cosh\) is plus \(\sinh\). A table of derivatives read from right to left is all the table of
integrals that is needed.

## Definite integrals

Over the interval from 0 to 1, the integral of \(\cosh x\) is \(\sinh 1 \approx 1.1752\), while
the integral of \(\sinh x\) is \(\cosh 1 - \cosh 0 \approx 1.5431 - 1 = 0.5431\). The second is
the smaller because the hyperbolic sine starts at 0 and the hyperbolic cosine starts at 1. Both
answers are positive, as the graphs demand.

## The chain rule case

A scaled argument brings a reciprocal factor. For \(\cosh 3x\) the antiderivative is
\(\frac{1}{3}\sinh 3x\), so the integral from 0 to 1 is \(\frac{1}{3}\sinh 3 \approx 3.339\). For
\(\sinh 2x\) it is \(\frac{1}{2}\cosh 2x\), so the integral from 0 to 1 is
\(\frac{1}{2}(\cosh 2 - 1) \approx 1.381\).

Sketch the integrand to check the signs, then differentiate the answer to check the factor in
front of it.
