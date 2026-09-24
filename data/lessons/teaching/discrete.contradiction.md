# Proof by Contradiction

**Contradiction:** To prove $P$, assume $\lnot P$ and derive something false, $Q\land\lnot Q$; then $\lnot P$ is impossible, so $P$ holds. This differs from contrapositive, which proves $\lnot Q\to\lnot P$ directly without deriving anything false.

## Worked: sqrt(2) is irrational

Assume $\sqrt2=p/q$ with $p,q$ positive integers in lowest terms:
1. Squaring gives $2q^2=p^2$, so $p^2$ is even and hence $p$ is even, say $p=2r$.
2. Then $2q^2=4r^2$, so $q^2=2r^2$ and $q$ is even too.
3. Both $p$ and $q$ even contradicts lowest terms, so the assumption is false.

So $\sqrt2$ cannot be written as a fraction: it is irrational.

## Worked: infinitely many primes

Assume only finitely many primes exist, $p_1,\dots,p_k$:
1. Form $N=p_1\cdots p_k+1$, one more than their product.
2. $N>1$ has some prime divisor $p$, and $p$ cannot be any $p_i$ since $N$ leaves remainder 1 upon division by each $p_i$.
3. That $p$ is a prime off the list, contradicting completeness.

So no finite list holds all primes: there are infinitely many.

## Contradiction versus contrapositive

To show $P\to Q$ by contrapositive, assume $\lnot Q$ and deduce $\lnot P$: every step is a true implication, nothing false is derived. By contradiction, assume $P\land\lnot Q$ and force $Q\land\lnot Q$. In general, reach for contradiction when $\lnot P$ hands you a concrete equation or object to push until it breaks.
