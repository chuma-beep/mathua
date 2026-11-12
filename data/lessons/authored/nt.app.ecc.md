# Elliptic curve cryptography

Elliptic curve cryptography replaces modular multiplication with addition on a curve. The group law is geometric, and the discrete logarithm problem becomes harder for the same key size.

## The curve equation

A short Weierstrass curve has the form y squared equals x cubed plus a x plus b. Points satisfy the equation modulo a prime, together with a point at infinity.

## The group law

Two points add by drawing a line, finding the third intersection and reflecting. The slope is a rational function of the coordinates, computed modulo the prime.

## Key size

The number of points is close to the prime by Hasse's bound. That bound is why a much smaller elliptic curve matches the security of a large multiplicative group.

To check your understanding, evaluate the right side, compute the slope numerator and state Hasse's bound.

Elliptic curve cryptography is the same discrete logarithm idea on a richer group, and the group law is where the arithmetic lives.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
