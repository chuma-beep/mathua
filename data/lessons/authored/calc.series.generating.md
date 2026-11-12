# Ordinary generating functions as series

## From sequence to series

A generating function packs a sequence into a power series: \(G(x) = \sum a_n x^n\). For \(a_n = n\) the series is \(x + 2x^2 + 3x^3 + \cdots\), which sums to \(x/(1-x)^2\) inside the unit interval. At \(x = 0.1\) the first three terms give 0.12 against the closed value 0.1235.

## Extracting a coefficient

The coefficient is recovered by differentiation: \(a_n = G^{(n)}(0)/n!\). For \(G(x) = 1/(1-2x) = \sum 2^n x^n\) the coefficient of \(x^3\) is 8. Three derivatives at zero give \(48/6 = 8\). Division in closed form is often the shorter route when the pattern is visible.

## The geometric generating function

The geometric generating function is the base case: \(\sum x^n = 1/(1-x)\) has every coefficient equal to 1. At \(x = 1/2\) the series totals 2. Differentiating and multiplying by \(x\) gives \(x/(1-x)^2 = \sum n x^n\), which at \(x = 1/2\) also totals 2, so \(\sum n/2^n = 2\).

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
