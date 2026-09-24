# One-Point Compactification

**One-point (Alexandroff) compactification:** For locally compact Hausdorff $X$, form $X^*=X$ union $\{\infty\}$ where neighbourhoods of $\infty$ are complements of compact sets. Then $X^*$ is compact Hausdorff containing $X$ densely.

## Worked: the line closes into a circle

Compactify $\mathbb{R}$ by one point:

1. Compact sets in $\mathbb{R}$ are exactly the closed bounded sets.
2. A neighbourhood of $\infty$ is therefore $(-\infty, a)$ union $(b, \infty)$ plus $\infty$: both tails at once.
3. Sending $x$ to its stereographic image and $\infty$ to the north pole matches these with circle arcs.

So $\mathbb{R}$ union $\{\infty\}$ is homeomorphic to $S^1$: the two ends join at infinity.

## Worked: the plane closes into a sphere

Repeat for $\mathbb{R}^2$:

1. Complements of large closed discs, plus $\infty$, form the neighbourhoods of infinity.
2. Stereographic projection sends each such neighbourhood to a cap around the north pole.
3. The bijection extends continuously both ways.

So $\mathbb{R}^2$ union $\{\infty\}$ is homeomorphic to $S^2$; one point compactifies the plane into a sphere.

## When compactification fails

Try the procedure on $\mathbb{Q}$:

1. $\mathbb{Q}$ is not locally compact: no point has a compact neighbourhood inside $\mathbb{Q}$.
2. The complement-of-compact neighbourhoods of $\infty$ then fail to separate $\infty$ from points of $\mathbb{Q}$.
3. The result is not Hausdorff compact, so Alexandroff needs local compactness.

So the construction is sharp: locally compact Hausdorff in, compact Hausdorff out; already-compact spaces gain nothing.
