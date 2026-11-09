# Euler's transformation of series

Euler's transformation converts an alternating series into a new series built from forward differences. The new series often converges faster, and for geometric terms its sum can be written in closed form.

## Forward differences

For a sequence \(a_n\), the forward difference is \(\Delta a_n=a_n-a_{n+1}\). Higher differences are applied repeatedly. Euler's transformation is
\[
\sum_{n=0}^{\infty}(-1)^n a_n=\sum_{k=0}^{\infty}\frac{\Delta^k a_0}{2^{k+1}}.
\]
The factor \(2^{k+1}\) is the price of replacing the alternating signs by one half at each step.

For a geometric sequence \(a_n=r^n\), the first difference is \((1-r)r^n\). Repeating the operation gives \(\Delta^k a_0=(1-r)^k\). The transformed coefficient is therefore \((1-r)^k/2^{k+1}\).

## Geometric terms

A geometric alternating series has a known sum. Euler's transformation reproduces it. The transformed series is geometric with ratio \((1-r)/2\), so its sum is
\[
\frac{1}{1-(1-r)/2}=\frac{1}{1+r}.
\]
The original alternating geometric series sums to \(1/(1+r)\) as well. The transformation has not changed the value, only the representation.

For \(r=1/2\) the transformed ratio is \(1/4\). For \(r=3/4\) it is \(1/8\). A larger \(r\) makes the transformed series converge faster.

## The harmonic series

For the alternating harmonic series, \(a_n=1/(n+1)\). Its forward differences are \(\Delta^k a_0=1/(k+1)\). The transformed coefficient is
\[
c_k=\frac{1}{(k+1)2^{k+1}}.
\]
The first coefficients are \(1/2\), \(1/8\) and \(1/24\). They decrease rapidly. The transformed series is a practical way to compute the sum, which is \(\ln 2\).
