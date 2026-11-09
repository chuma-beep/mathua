# Elliptic integrals

Elliptic integrals arise when arc lengths and periods cannot be written with elementary functions. The two complete integrals of Legendre are the most common.

## The first kind

The complete elliptic integral of the first kind is
\[
K(k)=\int_0^{\pi/2}\frac{d\theta}{\sqrt{1-k^2\sin^2\theta}}.
\]
The parameter \(k\) is the modulus. At \(k=0\) the integrand is 1 and \(K(0)=\pi/2\). As \(k\) approaches 1 the integral grows without bound.

The Maclaurin series is
\[
K(k)=\frac{\pi}{2}\sum_{n=0}^{\infty}\left(\frac{(2n)!}{2^{2n}(n!)^2}\right)^2k^{2n}.
\]
The first coefficients are 1, \(1/4\), \(9/64\) and \(25/256\). The series converges for \(|k|<1\).

## The second kind

The complete elliptic integral of the second kind is
\[
E(k)=\int_0^{\pi/2}\sqrt{1-k^2\sin^2\theta}\,d\theta.
\]
It measures the quarter circumference of an ellipse with eccentricity \(k\). At \(k=0\) it equals \(\pi/2\), and at \(k=1\) it equals 1.

Its series is
\[
E(k)=\frac{\pi}{2}\left[1-\sum_{n=1}^{\infty}\left(\frac{(2n)!}{2^{2n}(n!)^2}\right)^2\frac{k^{2n}}{2n-1}\right].
\]
The coefficients are negative after the constant term. The first correction is \(-k^2/4\), the second is \(-3k^4/64\).

## The arithmetic-geometric mean

The arithmetic-geometric mean gives a rapidly convergent way to compute \(K\). Start with \(a_0=1\) and \(b_0=\sqrt{1-k^2}\). Iterate
\[
a_{n+1}=\frac{a_n+b_n}{2},\qquad b_{n+1}=\sqrt{a_nb_n}.
\]
Both sequences converge to the same limit, the AGM. Then
\[
K(k)=\frac{\pi}{2\,\mathrm{AGM}(1,\sqrt{1-k^2})}.
\]
Each iteration roughly doubles the number of correct digits. Ten iterations are usually enough for ten decimal places.
