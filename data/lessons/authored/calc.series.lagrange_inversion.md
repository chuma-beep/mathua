# Lagrange inversion

Lagrange inversion extracts the coefficients of an inverse function without solving for the function explicitly. It is the standard tool for counting problems stated by an implicit equation.

## The inversion formula

Suppose \(w=z\phi(w)\) and \(\phi(0)\ne0\). The coefficient of \(z^n\) in the solution \(w(z)\) is
\[
[z^n]w(z)=\frac1n[t^{n-1}]\phi(t)^n.
\]
Only one coefficient of one power of \(\phi\) is needed. For small \(n\) the binomial theorem finishes the calculation.

The formula is valid for \(n\ge1\). The constant term is zero because \(w(0)=0\). Each application therefore starts by expanding \(\phi(t)^n\) and reading the term one degree below \(n\).

## The exponential case

Take \(w=z e^w\), so \(\phi(t)=e^t\). Then \(\phi(t)^n=e^{nt}\), and the coefficient of \(t^{n-1}\) is \(n^{n-1}/(n-1)!\). Dividing by \(n\) gives
\[
[z^n]w=\frac{n^{n-1}}{n!}.
\]
These are the coefficients of the tree function. For \(n=1\) the coefficient is 1 and for \(n=2\) it is 1. For \(n=3\) it is \(3/2\).

## Catalan coefficients

Take \(w=z/(1-w)\), which is equivalent to \(w=z(1+w+\cdots)\). The solution is the generating function of the Catalan numbers. The inversion formula with \(\phi(t)=1/(1-t)\) gives
\[
[z^n]w=\frac{1}{n+1}\binom{2n}{n}.
\]
The first values are 1, 1, 2, 5, 14 and 42. They count binary trees, balanced parentheses and many other structures. The coefficient is always an integer even though the formula contains a division.
