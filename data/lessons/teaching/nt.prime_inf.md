# Infinitude of Primes

**Theorem (Euclid):** There are infinitely many primes.

## Euclid's Construction

### The Product Plus One
Given any finite set of primes $\{p_1,\dots,p_k\}$, form

$$N = p_1p_2\cdots p_k + 1$$

### A New Prime Factor
$N$ is not divisible by any $p_i$ (it leaves remainder $1$ upon division by each $p_i$). Hence either $N$ is prime, or $N$ has a prime divisor not in the set. In either case a new prime exists.

## Example

Take primes $2,3,5$: $N=2\cdot3\cdot5+1=31$. $31$ is not divisible by $2,3,$ or $5$; it is prime and is a new prime.

Take primes $2,3,5,7$: $N=2\cdot3\cdot5\cdot7+1=211$. $211$ is prime.

Take primes $2,3,5,7,11,13$: $N=30031=59\times509$. Both $59$ and $509$ are new primes not in the original list.

## Why This Proves Infinitude
If there were only finitely many primes, Euclid's construction applied to the complete finite list would produce a prime outside the list — contradiction.
