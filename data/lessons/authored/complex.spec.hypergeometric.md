# The hypergeometric function

The hypergeometric function sums a ratio of rising factorials, and Gauss evaluated it at 1.

## Series

\({}_{2}F_{1}(a,b;c;z)=\sum_{n\ge0}\frac{(a)_{n}(b)_{n}}{(c)_{n}}\frac{z^{n}}{n!}\).

## Gauss value

At \(z=1\) the series equals \(\frac{\Gamma(c)\Gamma(c-a-b)}{\Gamma(c-a)\Gamma(c-b)}\) when the real parts allow it.

## Special cases

Many elementary functions are hypergeometric with one parameter suppressed, and the geometric series is \({}_{1}F_{0}(1;;z)\).

To check your understanding, read the parameters, apply the Gauss formula and reduce a parameter.

Hypergeometric series unify most of the classical special functions.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
