# RSA key generation

RSA builds a trapdoor from multiplication of two large primes. Publishing the product hides the factorisation, and the totient supplies the inverse exponent.

## The modulus

Choose distinct primes p and q and publish their product n. The modulus is public, and its factorisation is the secret.

## The totient

The totient of n is the product of one less than each prime. It counts the units modulo n and therefore the exponents invertible modulo it.

## The private exponent

The public exponent e is invertible modulo the totient. The private exponent d is its inverse, so exponentiation by d undoes exponentiation by e.

To check your understanding, multiply the primes, subtract one from each prime and invert e.

RSA security rests on the gap between multiplying two primes and factoring their product, and the totient is the bridge between the two exponents.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
