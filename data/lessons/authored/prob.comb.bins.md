# Balls into bins

Balls into bins asks how full the fullest bin becomes when balls are thrown independently and uniformly.

## Maximum load

With \(n\) balls and \(n\) bins the maximum load is about \(\log n/\log\log n\) with high probability.

## Two choices

Giving each ball two random choices and placing it in the lighter bin reduces the maximum load to about \(\log\log n\).

## Concentration

The maximum load concentrates because each bin's load is a sum of independent indicators.

To check your understanding, throw n balls, sample two bins and write a bin load.

Random placement leaves a logarithmically full bin, and two choices almost flatten it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
