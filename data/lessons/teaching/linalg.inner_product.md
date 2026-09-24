# Inner Products

**Inner product:** A symmetric positive-definite bilinear form $\langle u, v\rangle$: linear in each argument, $\langle u, v\rangle = \langle v, u\rangle$, and $\langle u, u\rangle \ge 0$ with equality only for $u = 0$. It turns a vector space into a geometry with lengths and angles.

## Worked: (1,2)·(3,4) = 11

Compute the dot product:
1. Multiply componentwise: $1 \cdot 3 = 3$ and $2 \cdot 4 = 8$.
2. Add: $3 + 8 = 11$.
3. So $(1,2)\cdot(3,4) = 11$ — and $\|(3,4)\| = 5$ since $\sqrt{9+16} = 5$.

So the dot product is multiply-then-add; the norm is its square root at $u = v$.

## Worked: weighted inner product on R^2

Test $\langle u,v\rangle = 2u_1v_1 + u_2v_2$:
1. Bilinear and symmetric: clear from the formula.
2. Positive definite: $2u_1^2 + u_2^2 \ge 0$, zero only at $u = 0$.
3. So it is an inner product — positive weights keep definiteness.

So weights are allowed as long as every weight stays positive; a zero or negative weight breaks it.

## Cauchy-Schwarz

$|\langle u,v\rangle| \le \|u\|\cdot\|v\|$: the angle between vectors has a cosine in $[-1, 1]$. Every inner product space inherits this bound.
