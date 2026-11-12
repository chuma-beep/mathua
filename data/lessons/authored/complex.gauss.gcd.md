# The Gaussian gcd

The Gaussian gcd is produced by the Euclidean algorithm and is unique up to multiplication by a unit.

## Euclidean steps

Each step replaces \(z,w\) by \(w,r\) with \(r=z-qw\). The final nonzero remainder is a gcd.

## Common divisors

A common divisor of \(z\) and \(w\) divides every integer combination, so it divides the gcd produced by the algorithm.

## Coprime pairs

Two Gaussian integers are coprime when their gcd is a unit. Conjugate pairs with coprime coordinates are an example.

To check your understanding, run one Euclidean step, test a common divisor and test for a unit gcd.

Uniqueness up to a unit is the Gaussian form of uniqueness of the gcd.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
