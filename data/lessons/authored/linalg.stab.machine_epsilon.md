# Machine epsilon

Machine epsilon is the gap between one and the next representable floating-point number. It sets the scale of rounding errors in every operation.

## Value

In double precision the significand has fifty-three bits, so the gap at one is two to the minus fifty-two. The unit roundoff is half that gap.

## Digits

Machine epsilon corresponds to about sixteen decimal digits. Every subtraction can lose digits in proportion to cancellation.

## Arithmetic

The relative error of a single rounded operation is bounded by the unit roundoff. Longer computations accumulate these bounds.

To check your understanding, count significand bits, convert to decimals and bound one operation.

Machine epsilon is the exchange rate between arithmetic operations and lost accuracy.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
