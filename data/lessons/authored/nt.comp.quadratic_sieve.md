# The quadratic sieve

The quadratic sieve collects smooth values of a quadratic polynomial and combines them into a congruence of squares. It is the classical subexponential factorisation method.

## The factor base

The algorithm chooses a smoothness bound and keeps primes for which n is a quadratic residue. Those primes form the factor base.

## Relations

Values of the polynomial near the square root of n are tested for smoothness over the factor base. Each smooth value gives a relation among the exponents.

## Linear algebra

The exponent vectors are reduced modulo two. A dependency produces a product of relations whose exponents are all even, hence a square.

To check your understanding, choose a bound, evaluate the polynomial and reduce modulo two.

The quadratic sieve is a pipeline from smooth values to a linear dependency, and every stage is cheaper than trial division on large inputs.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
