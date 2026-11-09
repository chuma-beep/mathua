# Fresnel integrals

The Fresnel integrals describe the amplitude of a diffracted wave. They are defined by oscillatory integrals whose limits at infinity are finite.

## The cosine integral

The Fresnel cosine integral is
\[
C(x)=\int_0^x\cos\left(\frac{\pi t^2}{2}\right)dt.
\]
The integrand oscillates faster as \(t\) grows. The oscillations cancel, so \(C(x)\) approaches \(1/2\) as \(x\to\infty\).

At small \(x\) the integrand is close to 1 and \(C(x)\) is close to \(x\). The first correction is negative and of order \(x^5\), which is why the graph bends below the line \(y=x\).

## The sine integral

The Fresnel sine integral is
\[
S(x)=\int_0^x\sin\left(\frac{\pi t^2}{2}\right)dt.
\]
It also approaches \(1/2\) at infinity. Near zero the integrand behaves like \(\pi t^2/2\), so \(S(x)\) grows like \(\pi x^3/6\).

The two integrals are the components of the Cornu spiral. A point on the spiral has coordinates \((C(x),S(x))\), and the spiral winds toward \((1/2,1/2)\).

## Maclaurin coefficients

Substituting the cosine series and integrating term by term gives
\[
C(x)=\sum_{n=0}^{\infty}\frac{(-1)^n(\pi/2)^{2n}x^{4n+1}}{(2n)!(4n+1)}.
\]
The coefficient of \(x^{4n+1}\) is the displayed fraction. For \(n=0\) it is 1, for \(n=1\) it is \(-\pi^2/40\), and the signs alternate.

The numerical values in this lesson are approximations obtained by summing the series through enough terms. The series converges for every \(x\), but the number of useful terms grows with \(x\).
