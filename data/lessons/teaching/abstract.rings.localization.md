# Localization

**Localization:** Given $R$ and a multiplicative set $S$ (closed under multiplication, $0 \notin S$), $S^{-1}R$ adjoins inverses for everything in $S$: fractions $r/s$ with the usual cross-multiplication rules. Localizing focuses the ring at the primes avoiding $S$.

## Worked: halving odd denominators

In $\mathbb Z[1/2]$, add $1/2 + 1/4$:
1. Common denominator 4: $2/4 + 1/4$.
2. Sum: $3/4$.
3. So $1/2 + 1/4 = 3/4$ — ordinary fraction arithmetic, now legal inside the ring.

So localization just permits the divisions in $S$; everything else computes normally.

## Worked: Z_(p) has one maximal ideal

Localize $\mathbb Z$ at the complement of $(p)$:
1. Allowed denominators are integers not divisible by $p$.
2. Every element outside $p\mathbb Z_{(p)}$ is a unit (numerator and denominator both prime to $p$).
3. Non-units form the single ideal $p\mathbb Z_{(p)}$ — hence local, with one maximal ideal.

So throwing away the primes outside $(p)$ collapses the ideal structure to a single point.

## Fractions as localization

The field of fractions of a domain is localization at all nonzero elements: $S = R \setminus \{0\}$. Inverting everything possible leaves a field.
