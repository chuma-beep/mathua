# Sylow Theorems

**Sylow $p$-subgroup:** For $|G| = p^k m$ with $p$ not dividing $m$, a subgroup of order $p^k$. Sylow guarantees one exists, all are conjugate, and their count $n_p$ satisfies $n_p \equiv 1 \bmod p$ with $n_p$ dividing $m$.

## Worked: Sylow 2-subgroups of order 12

For $|G| = 12 = 2^2 \cdot 3$:
1. A Sylow 2-subgroup has order 4: exists by Sylow.
2. $n_2$ divides 3 and is odd: $n_2 \in \{1, 3\}$.
3. Either one normal Sylow 2-subgroup or three conjugate ones — the arithmetic narrows it to two cases.

So Sylow turns "what subgroups exist?" into a short finite checklist.

## Worked: n_5 in a group of order 15

For $|G| = 15 = 3 \cdot 5$:
1. A Sylow 5-subgroup has order 5.
2. $n_5$ divides 3 and $n_5 \equiv 1 \bmod 5$: candidates 1, 3 — only 1 works.
3. So $n_5 = 1$: the Sylow 5-subgroup is unique, hence normal.

So congruence plus divisibility can force normality with no further work.

## Worked: n_3 possibilities in order 12

For $|G| = 12$ again: $n_3$ divides 4 and $n_3 \equiv 1 \bmod 3$.
1. Divisors of 4: 1, 2, 4. Only 1 and 4 are 1 mod 3.
2. So $n_3 \in \{1, 4\}$: either a normal Sylow 3-subgroup or four of them.
3. Both occur across different groups of order 12 — the constraints are sharp.

So the Sylow conditions bound the answer tightly but need not pin it uniquely.
