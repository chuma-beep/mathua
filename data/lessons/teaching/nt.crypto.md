# Cryptography: The RSA Algorithm

**RSA** (Rivest–Shamir–Adleman) is a public-key cryptosystem that uses number theory to secure communication.

## Key Generation

1. Choose two large distinct primes $p$ and $q$.
2. Compute $n = pq$ and $\phi(n) = (p-1)(q-1)$.
3. Choose $e$ such that $1 < e < \phi(n)$ and $\gcd(e, \phi(n)) = 1$.
4. Compute $d \equiv e^{-1} \pmod{\phi(n)}$ (the modular inverse).
5. **Public key:** $(n, e)$. **Private key:** $d$.

## Encryption

To send a message $m$ (as an integer with $0 \le m < n$):
$$c \equiv m^e \pmod{n}$$

## Decryption

To recover $m$ from $c$:
$$m \equiv c^d \pmod{n}$$

## Why it works

By Euler's theorem, $m^{\phi(n)} \equiv 1 \pmod{n}$ for $\gcd(m, n) = 1$.
Since $ed \equiv 1 \pmod{\phi(n)}$, we have $ed = 1 + k\phi(n)$, so:
$$c^d \equiv (m^e)^d = m^{ed} = m^{1 + k\phi(n)} = m \cdot (m^{\phi(n)})^k \equiv m \cdot 1^k = m \pmod{n}$$

## Security

The security of RSA relies on the difficulty of factoring $n$ to recover $p$ and $q$. With sufficiently large primes (e.g., 2048-bit $n$), factoring is computationally infeasible.

**Simple example:** Let $p = 5$, $q = 11$.
$n = 55$, $\phi(55) = 40$.
Choose $e = 3$ (since $\gcd(3, 40) = 1$).
$d \equiv 3^{-1} \pmod{40} = 27$ (since $3 \times 27 = 81 \equiv 1 \pmod{40}$).
Encrypt $m = 7$: $c = 7^3 = 343 \equiv 13 \pmod{55}$.
Decrypt $c = 13$: $m = 13^{27} \pmod{55} = 7$. ✓
