# Pseudorandom generators

A pseudorandom generator is a deterministic recurrence whose output passes statistical tests for randomness. Its period and correlation structure decide whether a simulation is trustworthy.

## Linear congruential

\(x_{n+1}=(ax_n+c)\bmod m\) has period at most m, and much less for poor multipliers. Its low-dimensional structure is visible in lattice tests.

## Modern generators

The Mersenne Twister has period \(2^{19937}-1\) and passes most tests, but it is not cryptographically secure. Counter-based generators are parallel-friendly and reproducible.

## Seeding

Independent streams need seeds that are provably separated. Reusing a seed across parallel workers correlates the results.

To check your understanding, apply the recurrence, state the Mersenne period and separate the streams.

A simulation inherits every defect of the stream that drives it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
