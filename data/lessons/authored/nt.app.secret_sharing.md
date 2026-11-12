# Shamir secret sharing

Secret sharing splits a secret into shares so that any threshold of them reconstructs it. Shamir's scheme uses polynomial interpolation modulo a prime.

## The polynomial

A random polynomial of degree one less than the threshold hides the secret in its constant term. Fewer than threshold shares leave every constant term possible.

## Interpolation

Any threshold shares determine the polynomial uniquely. Lagrange interpolation evaluates it at zero without recovering the other coefficients.

## The prime

All arithmetic happens modulo a prime larger than the secret and the share indices. Inverses exist because every nonzero residue is a unit.

To check your understanding, choose the degree, take threshold shares and pick a large prime.

Shamir sharing is interpolation in disguise, and the threshold is exactly the number of points that pins down a polynomial.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
