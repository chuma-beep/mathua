# Pratt certificates

A Pratt certificate proves primality by exhibiting a primitive root and recursively certifying the factors of p minus one. It turns primality into a small proof tree.

## The certificate

The prover gives a primitive root g and the factorisation of p minus one. For each prime factor q, the prover also gives a certificate that q is prime.

## The verification

The verifier checks that g has order p minus one by testing the powers at each prime factor. The recursion bottoms out at small primes.

## The size

The certificate length is governed by the number of prime factors of the predecessors. Smooth chains give short proofs.

To check your understanding, exhibit a primitive root, check the order and count factors.

Pratt certificates make primality verifiable in polynomial time, with the recursion carrying the proof down to the base primes.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
