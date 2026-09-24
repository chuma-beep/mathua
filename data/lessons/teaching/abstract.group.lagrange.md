# Lagrange's Theorem

**Lagrange's theorem:** For finite $G$ and subgroup $H$, $|H|$ divides $|G|$. The cosets of $H$ partition $G$ into equal pieces of size $|H|$, so the index $[G:H] = |G|/|H|$ counts them.

## Worked: cosets of {0, 3} in Z_6

Partition $Z_6$ by the subgroup $\{0, 3\}$:
1. Coset $0 + H = \{0, 3\}$; coset $1 + H = \{1, 4\}$; coset $2 + H = \{2, 5\}$.
2. Three cosets, each of size 2, covering all 6 elements with no overlap.
3. So $[Z_6 : H] = 3$ and $2 \cdot 3 = 6$: the subgroup order divides the group order.

So cosets tile the group evenly — divisibility is just counting tiles.

## Worked: no subgroup of order 4 in a group of order 6

Suppose $|G| = 6$ with $H \le G$ of order 4:
1. Lagrange forces $|H|$ to divide 6; 4 does not divide 6.
2. The cosets would need to tile 6 elements with pieces of size 4 — impossible without overlap.
3. So no such $H$ exists, in any group of order 6.

So Lagrange rules subgroups out by arithmetic alone, before any structural analysis.

## Converse fails

Divisibility is necessary, not sufficient: $A_4$ has order 12 but no subgroup of order 6, even though 6 divides 12. Lagrange constrains; it does not construct.
