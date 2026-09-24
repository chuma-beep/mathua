# Unique Factorization Domains

**UFD:** An integral domain where every nonzero non-unit factors into irreducibles, uniquely up to order and units. UFD is the precise sense in which "prime factorization works" — and it can fail.

## Worked: 12 has three prime factors in Z

Factor 12 in the integers:
1. $12 = 2 \cdot 6 = 2 \cdot 2 \cdot 3$.
2. 2 and 3 are irreducible (hence prime) in $\mathbb Z$.
3. Three prime factors with multiplicity — and any other factorization regroups these same primes.

So unique factorization means the multiset of irreducibles is an invariant of the number.

## Worked: 6 breaks in Z[√-5]

Factor 6 two ways in $\mathbb Z[\sqrt{-5}]$:
1. $6 = 2 \cdot 3$: both factors are irreducible (no proper factorization exists).
2. $6 = (1+\sqrt{-5})(1-\sqrt{-5})$: also irreducible, genuinely different (norms: $N(2) = 4$ vs $N(1\pm\sqrt{-5}) = 6$).
3. Two distinct factorizations: $\mathbb Z[\sqrt{-5}]$ is not a UFD.

So norms certify the break: same element, incompatible irreducible pieces.

## PID implies UFD

Every PID is a UFD (irreducible elements are prime once principal ideals give Bézout). $\mathbb Z$ and $F[x]$ inherit uniqueness from principality.
