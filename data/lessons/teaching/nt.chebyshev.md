# Chebyshev Functions

**Chebyshev functions:** \(\theta(x) = \sum_{p \le x} \log p\) weights each prime by its log, and \(\psi(x) = \sum_{n \le x} \Lambda(n)\) weights prime powers via von Mangoldt \(\Lambda\). The prime number theorem is exactly \(\psi(x) \sim x\).

## Worked: theta at 10

Add \(\log p\) over the primes up to 10:
1. Primes: 2, 3, 5, 7, four terms.
2. \(\theta(10) = \log 2+\log 3+\log 5+\log 7\).
3. Numerically about \(0.69+1.10+1.61+1.95 = 5.35\), below 10 since the asymptote \(\theta(x) \sim x\) only bites for large \(x\).

So \(\theta\) turns a prime list into one weighted sum: four primes become a single number.

## Worked: von Mangoldt at 8, 9, 12

\(\Lambda(n) = \log p\) when \(n = p^k\), else 0. Evaluate three inputs:
1. \(8 = 2^3\): a prime power, so \(\Lambda(8) = \log 2\).
2. \(9 = 3^2\): a prime power, so \(\Lambda(9) = \log 3\).
3. \(12 = 2^2 \cdot 3\): two distinct primes, so \(\Lambda(12) = 0\).

So \(\Lambda\) detects prime powers and names the base prime's log: \(\psi\) sums exactly these weights.

## Worked: psi at 10

Sum \(\Lambda(n)\) for \(n\) up to 10:
1. Nonzero at the primes 2, 3, 5, 7 plus \(4 = 2^2\), \(8 = 2^3\) (each \(\log 2\)), and \(9 = 3^2\) (\(\log 3\)).
2. Total: \(\theta(10)\) plus \(2\log 2+\log 3\), about 5.35 + 2.48 = 7.83.
3. So \(\psi(10)\) exceeds \(\theta(10)\) by the prime-power tail, and \(\psi(x) \sim x\) is the sharpest Chebyshev form of PNT.

So \(\psi\) refines \(\theta\) with prime powers: same asymptote, tighter error.
