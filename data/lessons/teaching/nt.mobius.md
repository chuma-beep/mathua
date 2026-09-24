# Möbius Function and Inversion

**Möbius function:** For \(n\) factoring as \(p_1^{e_1}\) through \(p_k^{e_k}\), set \(\mu(n) = 0\) if any \(e_i\) exceeds 1, else \((-1)^k\). It is 1 at \(n = 1\), -1 at each prime, and 0 wherever a square divides \(n\).

## Worked: values at 6, 12, 30

Factor each input and count distinct primes:
1. \(6 = 2 \cdot 3\): two distinct primes and no square divides it, so \(\mu(6) = (-1)^2 = 1\).
2. \(12 = 2^2 \cdot 3\): the square \(2^2\) divides 12, so \(\mu(12) = 0\).
3. \(30 = 2 \cdot 3 \cdot 5\): three distinct primes, so \(\mu(30) = (-1)^3 = -1\).

So a squared prime forces 0, and otherwise the sign records whether the prime count is even or odd.

## Worked: divisor sum over 6

The identity \(\sum_{d \mid n} \mu(d) = 0\) holds for every \(n\) above 1. Check \(n = 6\) with divisors 1, 2, 3, 6:
1. \(\mu(1) = 1\), \(\mu(2) = -1\), \(\mu(3) = -1\), \(\mu(6) = 1\).
2. Add them: 1 - 1 - 1 + 1 = 0.
3. The plus and minus ones cancel exactly.

So the \(\mu\) values over a divisor lattice always telescope to 0 past \(n = 1\).

## Worked: inverting to phi at 6

Möbius inversion recovers \(\varphi\) via \(\varphi(n) = n \sum_{d \mid n} \mu(d)/d\). Take \(n = 6\):
1. The four terms are 1, -1/2, -1/3, 1/6.
2. Sum: 1 - 1/2 - 1/3 + 1/6 = 1/3.
3. Multiply by 6: the result is 2.

So \(\varphi(6) = 2\) (namely 1 and 5), extracted from \(\mu\) alone: inversion turns divisor data back into the totient.
