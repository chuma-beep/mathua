# Galois Groups of Polynomials

**Galois group of a polynomial:** The Galois group of its splitting field over the base field: permutations of the roots preserving all algebraic relations. It encodes exactly how much symmetry the roots carry.

## Worked: x^2 - 2 gives Z_2

Find $\mathrm{Gal}$ of $x^2 - 2$ over $\mathbb Q$:
1. Splitting field is $\mathbb Q(\sqrt2)$, degree 2.
2. Automorphisms are fixed by the image of $\sqrt2$: $\sqrt2 \mapsto \pm\sqrt2$ — two choices.
3. So the Galois group has order 2: $Z_2$.

So for a single quadratic, the Galois group just flips the root sign.

## Worked: x^3 - 2 gives S_3

Find the Galois group of $x^3 - 2$ over $\mathbb Q$:
1. Roots: $\sqrt[3]{2}$, $\omega\sqrt[3]{2}$, $\omega^2\sqrt[3]{2}$ — adjoining one real root misses the complex pair.
2. The splitting field needs $\sqrt[3]{2}$ and $\omega$: degree $3 \cdot 2 = 6$.
3. The group permutes 3 roots faithfully with both a transposition (complex conjugation) and a 3-cycle present: $S_3$, order 6.

So degree 6 plus a transposition plus a 3-cycle forces the full symmetric group.

## Discriminants detect A_n

Square discriminant puts the Galois group inside $A_n$: parity of root permutations becomes visible in the coefficients. Generic quintics land on $S_5$ — the source of unsolvability.
