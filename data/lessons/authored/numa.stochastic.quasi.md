# Quasi-Monte Carlo

Quasi-Monte Carlo replaces random points with a low-discrepancy sequence. The points fill the domain more evenly, so the error falls faster than the random rate in low dimension.

## Discrepancy

Discrepancy measures the largest gap between the empirical and true fraction of points in a box. A low-discrepancy sequence keeps it small.

## Error rate

For a smooth integrand the error is about \(O((\log N)^{d}/N)\), better than \(N^{-1/2}\) when d is modest. The logarithm eventually dominates in high dimension.

## Randomisation

Scrambling restores unbiasedness and gives a computable error estimate. The scrambled sequence keeps its low discrepancy.

To check your understanding, compare counts in a box, compare with the random rate and scramble the sequence.

Evenly spread points beat random ones until the dimension catches up.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
