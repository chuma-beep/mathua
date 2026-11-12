# The number field sieve

The number field sieve is the fastest known classical factorisation method for large numbers. It replaces the quadratic polynomial with a number field.

## The idea

Two polynomials with a common root modulo n produce pairs of smooth values in two different rings. The pairs combine into a congruence of squares.

## The complexity

The heuristic running time is exponential in the cube root of the logarithm times a power of the log log. The constant is the cube root of sixty-four ninths.

## The record

The method set the records for factoring RSA challenge numbers. Its subexponential growth is why RSA key sizes grew with computing power.

To check your understanding, choose two polynomials, state the exponent and name a record.

The number field sieve is the quadratic sieve lifted into a number field, and the lift buys a cube root in the exponent.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
