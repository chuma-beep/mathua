# Sieve of Eratosthenes and Brun's Sieve

**Sieve of Eratosthenes:** List 2 through \(n\), then for each prime \(p\) cross out multiples starting at \(p^2\); survivors are prime. It runs in \(n \log \log n\) time. Brun's sieve truncates the same idea to prove upper bounds, e.g. that reciprocal twin primes converge.

## Worked: sieving to 30

Sift with \(p = 2, 3, 5\) (the primes up to the square root of 30):
1. Cross out multiples of 2 from 4 upward: 4, 6, 8, and so on.
2. Cross out multiples of 3 from 9 upward: 9, 15, 21, and so on (6 and 12 already gone).
3. Cross out multiples of 5 from 25 upward: only 25 is new; survivors are 2, 3, 5, 7, 11, 13, 17, 19, 23, 29.

So 10 primes up to 30 fall out of 3 passes: each composite's smallest factor was already handled.

## Worked: why start at p squared

Every composite below \(p^2\) was already crossed out by a smaller prime:
1. Take \(p = 5\): candidates below 25 like 15 or 21 carry a factor 3 or below.
2. In general a composite \(m\) below \(p^2\) has a prime factor at most the square root of \(m\), hence below \(p\).
3. So starting at \(p^2\) skips redundant work with no misses.

So the \(p^2\) start is bookkeeping, not a trick: smaller primes already covered the range below.

## Worked: Brun's bound on twins

Brun's sieve caps twin-prime counts where Eratosthenes would overcount:
1. Naive inclusion-exclusion over all primes up to 100 has over 33 million terms (25 primes) and is infeasible.
2. Brun truncates the alternating sum after a few terms, keeping an upper bound: twins up to \(x\) number at most about \(x/(\log x)^2\) up to a constant.
3. Summing that bound shows reciprocal twin primes converge (Brun's constant), without deciding whether twins are infinite.

So truncation trades exactness for control: an upper bound strong enough to settle convergence.
