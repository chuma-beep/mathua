# The AKS primality test

The AKS test decides primality in deterministic polynomial time. It replaces random bases with an identity between polynomials.

## The identity

A number n is prime exactly when the polynomial X plus a raised to n is congruent to X to the n plus a modulo n for every a in a bounded range.

## The modulus

The test works modulo X to the r minus one for a carefully chosen r. The choice of r keeps the polynomial degree small and the number of checks bounded.

## The cost

The original analysis gave a sixth power of the logarithm. Later work reduced the exponent, but the test remains slower than probabilistic methods in practice.

To check your understanding, raise the polynomial, choose r and state the exponent.

AKS is a proof that primality is in deterministic polynomial time, and its polynomial identity is the certificate that randomness is unnecessary.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
