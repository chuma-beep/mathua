# Quotient Topology

**Quotient topology:** For a surjection $q:X\to Y$ where $Y=X/\sim$ is the set of equivalence classes, $U\subseteq Y$ is open iff $q^{-1}(U)$ is open in $X$. It is the finest topology making $q$ continuous.

## Glueing Spaces

### The Quotient Map
$q$ collapses each class to a point. Openness in the quotient is pulled back to $X$: you check preimages.

### Circles from Intervals
Take $X=[0,1]$ with $0\sim1$. The quotient identifies the endpoints, yielding a space homeomorphic to $S^{1}$. The loop closes because neighborhoods of the glued point come from neighborhoods of both $0$ and $1$.

## Example

Is the quotient map $q:X\to X/\sim$ continuous? Yes by definition — the topology on $Y$ is chosen exactly to make $q$ continuous and final.
