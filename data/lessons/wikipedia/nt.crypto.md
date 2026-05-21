> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/RSA_%28cryptosystem%29) — CC BY-SA 4.0

# Introduction to RSA cryptography

**RSA** (Rivest–Shamir–Adleman) is a public-key cryptosystem widely used for secure data transmission. It relies on the practical difficulty of factoring the product of two large prime numbers.

### Key generation

1. Choose two distinct large primes \(p\) and \(q\)
2. Compute \(n = pq\)
3. Compute \(\phi(n) = (p-1)(q-1)\)
4. Choose an encryption exponent \(e\) coprime with \(\phi(n)\)
5. Compute the decryption exponent \(d\) such that \(ed \equiv 1 \pmod{\phi(n)}\)

The **public key** is \((n, e)\); the **private key** is \(d\).

### Encryption and decryption

To encrypt a message \(m\) (represented as an integer \(0 \le m < n\)):

\[
c = m^e \bmod n
\]

To decrypt:

\[
m = c^d \bmod n
\]

The security of RSA depends on the difficulty of factoring \(n\) and computing \(d\) without knowing \(p\) and \(q\).
