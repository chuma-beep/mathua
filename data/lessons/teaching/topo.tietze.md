# Tietze Extension Theorem

**Tietze extension:** If $X$ is normal and $A$ is closed in $X$, every continuous $f:A\to[0, 1]$ (or $\mathbb{R}$) extends to a continuous $F:X\to[0, 1]$ with $F=f$ on $A$. This extension property characterizes normality.

## Worked: extending from two points

Extend $f(0)=0$, $f(1)=1$ from $A=\{0, 1\}$ to all of $\mathbb{R}$:

1. $A$ is finite, hence closed, and $\mathbb{R}$ is normal, so Tietze applies.
2. On [0, 1] set $F(x)=x$, matching both prescribed values.
3. Extend constantly outside: $F=0$ left of 0 and $F=1$ right of 1, continuous at the joints.

So the two prescribed values grow into a global continuous function agreeing on $A$.

## Worked: Urysohn as a corollary

Derive Urysohn's lemma from Tietze:

1. Take disjoint closed $A_0, A_1$; their union $A=A_0$ union $A_1$ is closed.
2. Define $f=0$ on $A_0$ and $f=1$ on $A_1$, continuous on $A$ since the pieces are separated.
3. Tietze extends $f$ to $F$ on all of $X$: exactly the Urysohn separator.

So Urysohn's lemma is the two-set special case of Tietze extension.

## When the subset is not closed

Try extending from $A=\mathbb{Q}$ inside $\mathbb{R}$:

1. $A$ is not closed: irrationals are limit points outside it.
2. A bounded continuous $f$ on $\mathbb{Q}$ can oscillate toward different limits at an irrational boundary point.
3. No continuous global $F$ can match all such $f$ at once, so extension fails in general.

So closedness is essential: Tietze needs the subset to contain its limit points.
