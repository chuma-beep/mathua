# Group Actions

**Group action:** $G$ acts on $X$ when each $g$ permutes $X$ compatibly: $e$ fixes everything and $(gh)x = g(hx)$. The orbit of $x$ is everything reachable from $x$; the stabilizer is what fixes $x$.

## Worked: orbits partition {1, 2, 3} under S_3

Act by permuting positions:
1. From 1, some permutation sends $1 \to 2$ and another sends $1 \to 3$: the orbit of 1 is all of $\{1, 2, 3\}$.
2. One orbit covers everything, so the action is transitive — a single orbit.
3. The stabilizer of 3 is the permutations fixing 3: $\{e, (1\,2)\}$, of order 2.

So orbit-stabilizer checks out: $3 \cdot 2 = 6 = |S_3|$.

## Worked: |D_4| from orbit-stabilizer

$D_4$ acts on the 4 square vertices in one orbit:
1. One orbit of size 4: the action is transitive.
2. The stabilizer of a vertex has order 2 (identity plus reflection across its diagonal).
3. So $|D_4| = 4 \cdot 2 = 8$ — the group order falls out of the action.

So orbit-stabilizer turns geometry (one orbit, small stabilizer) into the group order.

## Orbit sizes vary

Fixed points are orbits of size 1; other orbits are larger. Orbit sizes always divide $|G|$, but different orbits of one action need not match — only the partition property is guaranteed.
