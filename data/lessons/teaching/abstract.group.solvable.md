# Solvable Groups

**Solvable group:** $G$ is solvable if it has a subnormal series down to 1 with abelian quotients — equivalently, the derived series $G \rhd G' \rhd G'' \rhd \cdots$ reaches 1. Solvability is what makes polynomials solvable by radicals.

## Worked: S_3 is solvable in two steps

Run the derived series of $S_3$:
1. $S_3' = A_3$: commutators of permutations are exactly the even ones.
2. $A_3$ is abelian (cyclic of order 3), so $A_3' = 1$.
3. Two steps reach 1: $S_3$ is solvable, with series $1 \lhd A_3 \lhd S_3$ and abelian quotients of orders 3 and 2.

So solvability is witnessed by an explicit short series, not an abstract property.

## Worked: S_5 is not solvable

Suppose $S_5$ had an abelian-quotient series:
1. The only nontrivial proper normal subgroup of $S_5$ is $A_5$ (simplicity).
2. Any series must pass through $A_5$, whose own only quotients are itself (non-abelian) and 1.
3. No abelian quotient can appear: $S_5$ is not solvable — and this is why quintics resist radicals.

So simplicity blocks solvability: with nowhere abelian to step, the series cannot exist.

## p-groups are solvable

Finite $p$-groups are nilpotent, hence solvable. The hard direction of Galois theory is that $S_n$ for $n \ge 5$ breaks the pattern.
