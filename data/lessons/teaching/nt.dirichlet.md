# Dirichlet's Theorem on Arithmetic Progressions

**Dirichlet's theorem:** The progression \(a, a+m, a+2m\), onward holds infinitely many primes exactly when \(a\) and \(m\) are coprime. The proof studies character sums \(L(s,\chi)\) mod \(m\) and their non-vanishing at \(s = 1\).

## Worked: primes 1 mod 4

Take \(a = 1\), \(m = 4\) with \(\gcd(1,4) = 1\):
1. The progression reads 1, 5, 9, 13, 17, 21, 29, and so on.
2. Primes appear throughout: 5, 13, 17, 29 are all 1 mod 4.
3. Dirichlet guarantees this never stops: infinitely many primes sit in this lane.

So coprimality opens the lane, and the theorem says traffic never ends.

## Worked: the even lane fails

Take \(a = 2\), \(m = 2\) with \(\gcd(2,2) = 2\):
1. The progression reads 2, 4, 6, 8, and so on: every term even.
2. The only even prime is 2 itself, appearing once at the start.
3. No theorem can help: the common divisor 2 blocks every later term.

So the coprime hypothesis is sharp: sharing a factor strangles the progression to finitely many primes.

## Worked: primes 3 mod 4

Take \(a = 3\), \(m = 4\) with \(\gcd(3,4) = 1\):
1. The progression reads 3, 7, 11, 15, 19, 23, and so on.
2. Primes 3, 7, 11, 19, 23 appear; 15 factors as \(3 \cdot 5\).
3. Again infinite: Euclid-style arguments show this case directly, and Dirichlet covers every coprime lane uniformly.

So 3 mod 4 behaves like 1 mod 4: coprime lanes always deliver infinitely many primes.
