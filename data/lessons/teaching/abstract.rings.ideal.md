# Ideals in Rings

**Ideal:** $I \subseteq R$ is an ideal if it is closed under addition and absorbs multiplication: $r \cdot x \in I$ for every $r \in R$, $x \in I$. Ideals are exactly the kernels of ring homomorphisms — and exactly what quotients require.

## Worked: {0,1,2} fails in Z_6

Test closure under addition mod 6:
1. $2 + 2 = 4$ mod 6.
2. 4 lies outside $\{0, 1, 2\}$.
3. So the set is not closed, hence not an ideal — one computation kills it.

So ideal checks fail fast: a single escaped sum is a complete disproof.

## Worked: 3Z passes in Z

Test both axioms for multiples of 3:
1. Closure: $3a + 3b = 3(a+b)$ stays a multiple of 3.
2. Absorption: $r \cdot 3a = 3(ra)$ stays a multiple of 3, for any integer $r$.
3. So $3Z$ is an ideal — and $Z/3Z \cong Z_3$ is the quotient it builds.

So the two checks together certify an ideal; the quotient is the reward.

## The zero-constant-term ideal

Polynomials in $R[x]$ with zero constant term form an ideal (sums and multiples keep zero constant term): it is the kernel of evaluation at 0, with quotient $R$.
