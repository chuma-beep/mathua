# Sieve of Eratosthenes and Brun's Sieve

**Sieve of Eratosthenes:** List $2,\dots,n$, iteratively mark multiples of each prime $p$ starting at $p^{2}$; survivors are primes. Time $O(n\log\log n)$.

**Brun's sieve:** Modern sieve giving upper bounds, e.g. count of twin primes $\le C x/(\log x)^{2}$ and convergence of $\sum_{p,p+2\text{ prime}}1/p$ (Brun's constant).

## Sieving Primes

### Eratosthenes
Multiples below $p^{2}$ already marked by smaller primes; remaining candidates after sieving to $\sqrt{n}$ are prime. Segmented sieve processes intervals for large $n$.

### Inclusion-Exclusion
Counting via Möbius: $\pi(n)=\sum_{d}\mu(d)\lfloor n/d\rfloor$-like sums; Brun truncates inclusion-exclusion to get upper bounds without full exactness.

## Example

$n=30$: primes $2,3,5$ mark multiples; survivors $\{2,3,5,7,11,13,17,19,23,29\}$.
