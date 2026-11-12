# The Weibull distribution

The Weibull distribution is the standard lifetime family with a shape parameter that controls the hazard's trend.

## Survival

Survival is the exponential of minus \((t/\lambda)^k\), with scale \(\lambda\) and shape \(k\).

## Hazard

The hazard is proportional to \(t^{k-1}\), so it decreases for \(k<1\), stays constant at \(k=1\) and increases for \(k>1\).

## Special case

Shape one is the exponential distribution, the memoryless member of the family.

To check your understanding, scale the time, differentiate the log survival and set shape to one.

The Weibull shape parameter chooses the ageing story.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
