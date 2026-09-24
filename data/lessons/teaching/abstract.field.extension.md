# Field Extensions

**Field extension:** $K/F$ adjoins new elements to $F$ while keeping field operations. The degree $[K:F]$ is the dimension of $K$ as an $F$-vector space — it measures how much was added.

## Worked: [Q(√2):Q] = 2

Build $\mathbb Q(\sqrt2)$:
1. Adjoin $\sqrt2$: elements are $a + b\sqrt2$ with $a, b \in \mathbb Q$.
2. $\{1, \sqrt2\}$ is linearly independent over $\mathbb Q$ ($\sqrt2$ is irrational).
3. Basis of size 2: the degree is 2, matching the minimal polynomial $x^2 - 2$.

So the degree equals the minimal-polynomial degree: adjoining a root costs exactly that dimension.

## Worked: splitting x^2 - 2 costs nothing extra

Split $x^2 - 2$ over $\mathbb Q$:
1. Roots are $\pm\sqrt2$; adjoining $\sqrt2$ gets $-\sqrt2$ free.
2. So $\mathbb Q(\sqrt2)$ already contains both roots: it is the splitting field.
3. Degree 2 total — no tower needed when one root drags in the other.

So splitting fields can coincide with simple extensions; the degree counts distinct new directions only.

## Towers multiply

$[K:\mathbb Q] = [K:F][F:\mathbb Q]$: degrees multiply up towers. $\mathbb Q(\sqrt2, \sqrt3)$ has degree $2 \cdot 2 = 4$ — two independent quadratic steps.
