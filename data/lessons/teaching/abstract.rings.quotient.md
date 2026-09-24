# Quotient Rings

**Quotient ring:** For an ideal $I \subseteq R$, the cosets $R/I$ form a ring under $(a+I) + (b+I) = (a+b)+I$ and similarly for multiplication. The ideal absorbs exactly the ambiguity representatives introduce.

## Worked: Z/6Z has six elements

List the cosets of $6\mathbb Z$:
1. $0 + 6\mathbb Z, 1 + 6\mathbb Z, \dots, 5 + 6\mathbb Z$: six distinct cosets.
2. $6 + 6\mathbb Z = 0 + 6\mathbb Z$: the seventh wraps back.
3. So $Z/6Z = \{0, 1, 2, 3, 4, 5\}$ with mod-6 arithmetic — nothing more.

So quotienting the integers by $n\mathbb Z$ just builds clock arithmetic.

## Worked: 2Z/4Z has two elements

Quotient even integers by multiples of 4:
1. Cosets: $2\mathbb Z$ itself (evens congruent to 0 mod 4) and $2 + 4\mathbb Z$ (evens congruent to 2 mod 4).
2. Two cosets, and doubling the non-identity one returns: $2 + 2 = 4 \in 4\mathbb Z$.
3. So $2\mathbb Z/4\mathbb Z \cong Z_2$: even the quotient of a non-unital ring by an ideal behaves.

So coset counting works uniformly — count the pieces, read off the structure.

## Ideals control quotients

Bigger ideal means fewer cosets: $Z/2Z$ has 2 elements where $Z/6Z$ has 6. The ideal is the dial setting the quotient's size.
