# Bessel functions of the first kind

Bessel functions solve the radial part of the Laplacian in polar coordinates.

## The Bessel equation

The equation is \(x^{2}y''+xy'+(x^{2}-n^{2})y=0\) for a nonnegative integer order \(n\).

The equation is linear of second order. The coefficient of \(y'\) is singular at the origin. No elementary function solves the equation. Only one member of the two-parameter family stays finite there. For order 0 the equation becomes \(xy''+y'+xy=0\).

At \(x=1\) that equation reads \(y''+y'+y=0\), whose solutions oscillate with period \(4\pi/\sqrt{3}\), about 7.26.

## The series solution

The solution finite at the origin is \(J_{n}(x)=\sum_{k}(-1)^{k}\frac{1}{k!(k+n)!}(x/2)^{2k+n}\).

The factorials in the denominator make the series converge everywhere. Matching powers of \(x\) gives a two-term recurrence for the coefficients, which is where those factorials come from.

For order 0 the coefficient of \(x^{2}\) is \(-1/4\) and of \(x^{4}\) is \(1/64\).

## A worked value

At \(x=2\) the ratio \(x/2\) is 1, so \(J_{0}(2)\) is the alternating sum \(1-1+1/4-1/36+1/576\) and so on.

The partial sums are 0, 0.25, 0.2222 and 0.22396, settling near 0.2239. A few terms already give four correct digits. The first positive zero is at 2.4048.

The series gives values near the origin. The zeros give the frequencies of a drum.
