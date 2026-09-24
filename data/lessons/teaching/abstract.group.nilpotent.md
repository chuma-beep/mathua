# Nilpotent Groups

**Nilpotent group:** $G$ is nilpotent if its upper central series reaches $G$ in finitely many steps — equivalently, it is the direct product of its Sylow subgroups. Nilpotent sits strictly between abelian and solvable.

## Worked: a group of order 8

Take $|G| = 8 = 2^3$:
1. $G$ is a 2-group: its only prime is 2.
2. Finite $p$-groups are nilpotent (nontrivial center, induct upward).
3. So every group of order 8 is nilpotent — no further analysis needed.

So prime-power order alone forces nilpotence; the Sylow structure is trivial.

## Worked: counting Sylows in order 12

A nilpotent group of order $12 = 2^2 \cdot 3$:
1. Nilpotent means every Sylow is normal, hence unique: one Sylow 2-subgroup, one Sylow 3-subgroup.
2. Two Sylow subgroups total — versus up to $3 + 4 = 7$ in a non-nilpotent group of order 12.
3. So nilpotence collapses the Sylow combinatorics to the minimum.

So uniqueness of Sylows is the usable face of nilpotence: count them, and there are exactly two.

## Strictly between

Abelian implies nilpotent (class 1: $A_3$ has class 1), nilpotent implies solvable. $S_3$ is solvable but not nilpotent — the inclusions are strict both ways.
