# Uniqueness of power series coefficients

## Equal series and equal coefficients

If two power series agree on an interval around the centre, their coefficients agree term by term. The geometric series \(1/(1-x) = 1 + x + x^2 + \cdots\) forces every coefficient to be 1. Any other series equal to the same function on the interval from -1 to 1 must therefore have \(a_0 = 1\), \(a_1 = 1\) and \(a_2 = 1\).

## Differentiating at the centre

Differentiating \(n\) times and setting \(x = 0\) isolates one coefficient: \(a_n = f^{(n)}(0)/n!\). For \(f(x) = e^{2x}\) this gives \(a_3 = 8/6 = 4/3\), matching the expansion \(1 + 2x + 2x^2 + (4/3)x^3 + \cdots\). The centre matters, since the same function expanded at 1 has different coefficients.

## Consequences

Uniqueness turns functional identities into coefficient identities. From \(1/(1-x)^2 = \sum (n+1)x^n\) the coefficient of \(x^3\) is 4, so the series for \((1-x)^{-2}\) is fixed without any division. The same move solves recurrences, because a generating function expanded two ways must agree at every power of \(x\).

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
