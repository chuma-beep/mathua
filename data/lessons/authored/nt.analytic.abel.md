# Abel summation

Abel summation converts a sum of coefficients times a smooth weight into an integral of the partial sums. It is the discrete analogue of integration by parts.

## The identity

The sum of a(n) f(n) equals the partial sum at the top times f at the top, minus an integral of the partial sums against the derivative of f.

## Polynomial weights

When the coefficients are all one, the partial sums are the identity. Summing a polynomial then reduces to an integral of a polynomial.

## Why it helps

The identity replaces pointwise knowledge of the coefficients with knowledge of their averages. Average information is usually easier to estimate.

To check your understanding, write the boundary term, sum unit coefficients and replace coefficients by averages.

Abel summation trades a sum for an integral, and the price is only a boundary term.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
