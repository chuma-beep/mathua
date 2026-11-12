# Concentration inequalities

Concentration inequalities bound how far an average deviates from its mean. Hoeffding needs bounded variables, and Bernstein adds a variance term.

## Hoeffding

The exponent is \(-2n\epsilon^2/(b-a)^2\) for variables in an interval of width \(b-a\). More samples shrink the deviation exponentially.

## Chernoff

The multiplicative form bounds a sum relative to its mean. It is sharper than Hoeffding for small deviations.

## Union bounds

A union over \(k\) events multiplies the failure probability by \(k\). The extra \(\log k\) appears in every sample complexity bound.

To check your understanding, set the interval width, use the multiplicative form and count the events.

Concentration turns an average into a guarantee.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
