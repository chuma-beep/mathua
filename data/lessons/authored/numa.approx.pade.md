# Padé approximants

A Padé approximant matches a function and its derivatives at one point using a ratio of polynomials. It often approximates better than a truncated Taylor series.

## Matching the Taylor series

A \([m/n]\) Padé approximant matches the Taylor series through degree \(m+n\). The numerator has degree \(m\) and the denominator degree \(n\).

The free coefficients number \(m+n+1\) after normalization. The matching conditions determine them.

Normalizing the denominator constant to 1 removes one redundant degree of freedom.

## The 1/1 approximant of the exponential

For \(e^x\), the \([1/1]\) approximant is \((1+x/2)/(1-x/2)\). It matches the series through the \(x^2\) term.

Unlike a truncated series, it can capture a pole-like behavior and may remain accurate farther from the expansion point.

The approximant is not the exponential. It agrees with it only near the expansion point.

## Evaluating the approximant

Substitute the value of \(x\) into numerator and denominator. The result is a rational number, not a polynomial term.

The denominator must not be zero at the evaluation point. A nearby zero limits the useful range of the approximation.

Rational functions can model poles that polynomials cannot, which is the main reason to use a Padé approximant.
