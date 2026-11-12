# Pollard rho

Pollard's rho finds a factor by iterating a polynomial modulo n and looking for a cycle. A gcd exposes the factor when the cycle closes modulo one prime but not the other.

## The iteration

The sequence starts at a small value and maps x to x squared plus c modulo n. The constant c is chosen to avoid degenerate behaviour.

## The gcd step

At each step the algorithm takes the difference of two sequence values and a gcd with n. A nontrivial gcd is a factor.

## The expected cost

A cycle modulo a prime factor p appears after about the square root of p steps. The cost therefore depends on the smallest factor, not on n.

To check your understanding, square and add, take a difference and take a square root.

Pollard rho converts a factorisation problem into a collision problem, and collisions arrive after a square-root number of steps.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
