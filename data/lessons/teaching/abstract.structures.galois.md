# Galois Correspondence

**Galois group:** $\operatorname{Gal}(K/F)$ is automorphisms of $K$ fixing $F$. For Galois (normal + separable) extension, intermediate fields $F\subseteq E\subseteq K$ correspond bijectively to subgroups $H\le\operatorname{Gal}(K/F)$ (order-reversing), with $[K^{H}:F]=|G|/|H|$.

## The Correspondence

### Fields ↔ Subgroups
$E\mapsto \operatorname{Gal}(K/E)$, $H\mapsto K^{H}=\{x\in K:h\cdot x=x\ \forall h\in H\}$. Normal subgroups ↔ Galois subextensions.

### Solvability by Radicals
Polynomial solvable by radicals iff its Galois group is solvable. $x^{3}-2$ has Galois group $S_3$ (non-abelian but solvable), so solvable; general quintic $S_5$ is not solvable.

## Example

$K=\mathbb Q(\sqrt2)$, $F=\mathbb Q$: $\operatorname{Gal}(K/F)=\{id,\sigma\}$ with $\sigma(\sqrt2)=-\sqrt2$, isomorphic to $\mathbb Z_2$. Only intermediate fields are $F$ and $K$.
