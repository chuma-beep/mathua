# The Pohlig-Hellman algorithm

Pohlig-Hellman solves a discrete logarithm one prime factor of the group order at a time. A smooth group order makes the problem easy.

## The reduction

The exponent is recovered modulo each prime power dividing the group order. The Chinese remainder theorem then combines the residues.

## The cost

The work is dominated by the largest prime factor, not by the full order. A group order with only small factors offers almost no security.

## The defence

Cryptographic groups are chosen with a large prime factor in their order. That choice forces the largest subproblem to stay hard.

To check your understanding, factor the order, find the largest factor and choose a large prime factor.

Pohlig-Hellman exploits the factorisation of the group order, which is why key generation insists on a large prime factor.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
