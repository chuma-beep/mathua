# Diffie-Hellman Key Exchange

## Public and Private Values

Diffie-Hellman key exchange lets two parties agree on a shared secret over a public channel. They agree on a large prime p and a generator g modulo p. Each party chooses a private exponent, publishes g raised to that exponent modulo p, and keeps the exponent secret.

## The Shared Secret

The shared secret is the other party's public value raised to one's own private exponent. Both sides compute the same number because exponentiation commutes: g to the a b equals g to the b a modulo p. An eavesdropper sees g, p, g to the a and g to the b, and recovering the secret requires solving a discrete logarithm problem.

## Choosing Parameters

Security depends on choosing p and g carefully and on keeping the private exponents secret. Small parameters are fine for illustrating the arithmetic but provide no security. Real implementations use large primes, ephemeral exponents and authenticated messages so that an active attacker cannot substitute public values.
