# Gauss-Legendre quadrature

Gauss-Legendre quadrature chooses nodes and weights to integrate polynomials exactly. An \(n\)-point rule is exact for every polynomial of degree at most \(2n-1\).

## The two-point rule

On the interval \([-1,1]\) the two-point rule is
\[
\int_{-1}^{1}f(x)\,dx\approx f\left(-\frac{1}{\sqrt3}\right)+f\left(\frac{1}{\sqrt3}\right).
\]
The nodes are symmetric and the weights are both 1. The rule is exact for polynomials up to degree 3.

To use another interval, scale the nodes by the half-length and multiply the whole rule by the half-length. On \([-a,a]\) the nodes become \(\pm a/\sqrt3\) and each weight becomes \(a\).

For \(f(x)=x^2\) on \([-a,a]\), each node contributes \(a^3/3\), so the rule returns \(2a^3/3\), which is the exact integral.

## The three-point rule

The three-point rule adds the center node:
\[
\int_{-1}^{1}f(x)\,dx\approx\frac59 f\left(-\sqrt{\frac35}\right)+\frac89 f(0)+\frac59 f\left(\sqrt{\frac35}\right).
\]
The weights sum to 2, the length of the interval. The rule is exact for polynomials up to degree 5.

For \(f(x)=x^4\) on \([-a,a]\), the center term vanishes and each outer node contributes \(a^5/5\). The rule returns \(2a^5/5\). A four-point rule would be needed to integrate \(x^6\) exactly.

## Changing the interval

Every rule is stated on \([-1,1]\) and moved by an affine substitution. If the target interval is \([a,b]\), the midpoint is \((a+b)/2\) and the half-length is \((b-a)/2\).

The nodes are midpoint plus half-length times the standard nodes. Every weight is multiplied by the half-length. This scaling preserves the exactness degree.

For a linear function the two-point rule returns the exact integral, because the midpoint rule and the trapezoid rule agree on a line. For a quadratic it also returns the exact integral, which is the first nontrivial gain over the trapezoid rule.
