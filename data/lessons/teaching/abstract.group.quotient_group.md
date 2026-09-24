# Quotient Groups

**Quotient group:** For $N$ normal in $G$, the cosets $G/N$ form a group under $(aN)(bN) = abN$. Normality is what makes this multiplication well-defined — without it, the product depends on representatives.

## Worked: S_3/A_3 has order 2

Quotient $S_3$ by $A_3$:
1. $|S_3| = 6$ and $|A_3| = 3$, so there are exactly 2 cosets: $A_3$ and $(1\,2)A_3$.
2. A group of order 2 is forced: the non-identity coset squares to the identity coset.
3. So $S_3/A_3 \cong Z_2$: the quotient records only parity.

So the quotient forgets everything except the invariant the subgroup cannot see.

## Worked: Z_12/{0,4,8} has order 4

Quotient out the order-3 subgroup:
1. $12 / 3 = 4$ cosets: $\{0,4,8\}$, $1 + H$, $2 + H$, $3 + H$.
2. The coset $1 + H$ has order 4: $k + H = H$ exactly when $k \in \{0, 4, 8\}$, and 4 is the first positive hit.
3. A cyclic generator of order 4: the quotient is $Z_4$.

So quotienting $Z_{12}$ by its order-3 part leaves the order-4 part visible.

## Quotients of Z_n

$Z_n/H$ is cyclic of order $n/|H|$: quotients of cyclic groups stay cyclic, with order divided down.
