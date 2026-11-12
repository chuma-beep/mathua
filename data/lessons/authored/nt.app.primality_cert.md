# Primality certificates

A primality certificate is a short proof that a number is prime. It is verified faster than the number can be tested from scratch.

## Pratt certificates

A Pratt certificate exhibits a primitive root and factors p minus one recursively. Each prime factor comes with its own certificate.

## Certificate size

The size of a Pratt certificate is governed by the factorisation of p minus one. Numbers whose predecessor is smooth have short certificates.

## Verification

The verifier checks the order of the primitive root and recurses on the factors. Every check is a modular exponentiation or a multiplication.

To check your understanding, exhibit a primitive root, count prime factors and check the order.

A certificate shifts work from the verifier to the prover, and its length measures how complicated the predecessor of p is.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
