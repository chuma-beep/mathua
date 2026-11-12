# Quadratic variation

The quadratic variation of a process accumulates squared increments. Brownian motion has finite quadratic variation even though its paths have infinite total variation.

## Definition

For a partition with mesh tending to zero, the sum of squared increments converges in probability to the quadratic variation. For Brownian motion on \([0,t]\), that limit is \(t\).

## Why it matters

Ordinary calculus assumes squared increments vanish faster than the first power. Brownian paths violate that, so a second-order term survives in the Ito expansion.

## Computing it

For \(B_t\), the variation over \(n\) equal subintervals is a sum of \(n\) squared Gaussian increments whose expectation is the interval length.

To check your understanding, sum squared increments, compare with total variation and count the subintervals.

Quadratic variation is the correction term that makes stochastic calculus different.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
