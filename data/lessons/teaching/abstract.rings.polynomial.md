# Polynomial Rings and Field Extensions

**Polynomial ring:** For a ring $R$, $R[x]$ is the set of polynomials with coefficients in $R$, added and multiplied in the usual way. When $R$ is a field, $R[x]$ behaves like the integers: it is a PID, so every ideal is principal and each quotient is controlled by a single polynomial.

## Ideals by evaluation

The map "evaluate at 0" sends $R[x]$ onto $R$; its kernel is $(x)$, the polynomials with zero constant term. So $R[x]/(x)\cong R$: quotienting by $(x)$ just sets $x=0$ everywhere.

## Worked: R[x]/(x^2+1) is C

Take two cosets and multiply them:
1. $(a+bx)+(x^2+1)$ times $(c+dx)+(x^2+1)$ gives $ac+(ad+bc)x+bd\,x^2$.
2. Reduce with $x^2=-1$: $(ac-bd)+(ad+bc)x$.
3. Compare with $(a+bi)(c+di)=(ac-bd)+(ad+bc)i$ — identical bookkeeping. The class of $x$ behaves exactly as $i$.

So $x^2+1$, irreducible over $\mathbb R$, generates a maximal ideal, and the quotient is a field: $\mathbb R[x]/(x^2+1)\cong\mathbb C$.

## Worked: adjoining a square root

In $\mathbb Q[x]/(x^2-2)$, the class of $x$ squares to 2, so it behaves as $\sqrt2$: the quotient is $\mathbb Q(\sqrt2)$. Every element is $a+bx$ with $a,b\in\mathbb Q$ — a 2-dimensional $\mathbb Q$-vector space. In general, irreducible $p(x)$ of degree $n$ gives an $n$-dimensional extension.

## When the quotient is not a field

$x^2$ is not irreducible, so $(x^2)$ is not maximal: in $\mathbb R[x]/(x^2)$ the class of $x$ is nonzero but squares to zero — a zero divisor. Reducible modulus gives zero divisors, so the quotient is not a field.
