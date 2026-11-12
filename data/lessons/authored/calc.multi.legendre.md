# Legendre polynomials and their generating function

Legendre polynomials are the orthogonal polynomials on the interval from -1 to 1 with unit weight.

## The generating function

The generating function is \((1-2xt+t^{2})^{-1/2}\), whose expansion in powers of \(t\) has the polynomials as coefficients.

The coefficient of \(t^{n}\) is the polynomial of degree \(n\). The first three are 1, \(x\) and \((3x^{2}-1)/2\).

For small \(t\) the square root expands as a power series. Each coefficient is a polynomial in \(x\). At \(x=0.5\) the third is \((0.75-1)/2\), which is -0.125.

## The recurrence relation

The recurrence is \((n+1)P_{n+1}=(2n+1)xP_{n}-nP_{n-1}\).

It needs only the two previous values at each step, which is how the polynomials are computed in practice. At \(n=1\) it gives \(2P_{2}=3x^{2}-1\).

At \(x=2\) the right side is 11, so \(P_{2}=5.5\). The recurrence is stable, unlike a direct expansion at large degree.

## Orthogonality

The relation is \(\int_{-1}^{1}P_{m}P_{n}\,dx=\frac{2}{2n+1}\) when the indices agree. It is zero otherwise.

For \(P_{1}\) and \(P_{2}\) the integral vanishes by symmetry, since one factor is odd. For \(P_{2}\) alone it is \(2/5\).

Integrating \((9x^{4}-6x^{2}+1)/4\) gives \((1/4)(18/5-4+2)\), which is 0.4.

The factor \(2/(2n+1)\) makes the polynomials a basis for the weighted space. Orthogonality is what makes a series of them usable for approximation.
