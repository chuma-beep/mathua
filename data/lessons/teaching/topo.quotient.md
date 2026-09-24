# Quotient Topology

**Quotient topology:** For a surjection $q:X\to Y$ where $Y=X/\sim$ collects equivalence classes, a set $U$ in $Y$ is open exactly when its preimage $q^{-1}(U)$ is open in $X$. It is the finest topology making $q$ continuous.

## Worked: the circle from an interval

Glue the endpoints of $X=[0, 1]$ by $0\sim 1$:

1. Points strictly inside map to singletons; the class $\{0, 1\}$ becomes one point $p$.
2. A neighbourhood of $p$ pulls back to $[0, e)$ union $(1-e, 1]$ for some small $e$, open in the class direction at both ends.
3. Sending $t$ to the angle $2\pi t$ matches these neighbourhoods with arcs around the circle.

So the quotient $X/\sim$ is homeomorphic to $S^1$: glueing endpoints closes the interval into a loop.

## Worked: the quotient map is continuous

Check $q:X\to X/\sim$ against the definition:

1. Take any open $U$ in the quotient $X/\sim$.
2. By the defining rule, $q^{-1}(U)$ is open in $X$ exactly when $U$ counts as open.
3. Hence preimages of opens are open, which is continuity.

So $q$ is continuous by construction: the topology on $Y$ is chosen as the finest one making $q$ continuous.

## When the quotient is not Hausdorff

Collapse the non-closed set $(0, 1]$ inside $X=[0, 2]$ to a point $p$:

1. The class $p$ has neighbourhoods pulling back to opens containing all of $(0, 1]$.
2. Any such neighbourhood also reaches toward 0, so every neighbourhood of $p$ meets every neighbourhood of the class $\{0\}$.
3. Thus $p$ and $\{0\}$ admit no disjoint neighbourhoods: the quotient is not Hausdorff.

So quotients of nice spaces can be badly behaved: collapsing a non-closed set destroys separation.
