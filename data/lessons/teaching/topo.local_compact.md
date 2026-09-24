# Local Compactness

**Locally compact:** $X$ is locally compact if every $x$ has a compact neighbourhood, meaning a compact set containing a neighbourhood of $x$. Every compact space is locally compact; $\mathbb{R}^n$ is locally compact via closed balls.

## Worked: the line is locally compact

Find a compact neighbourhood of $x$ in $\mathbb{R}$:

1. Take the closed interval $[x-e, x+e]$ for any small positive $e$.
2. It contains the open interval $(x-e, x+e)$, a neighbourhood of $x$.
3. It is closed and bounded, hence compact by Heine-Borel.

So every real point sits inside a compact neighbourhood, and the same closed-ball argument works in $\mathbb{R}^n$.

## Worked: the rationals are not

Try to find a compact neighbourhood of 0 in $\mathbb{Q}$:

1. Any neighbourhood in $\mathbb{Q}$ looks like $(-e, e)$ intersected with $\mathbb{Q}$.
2. Its closure in $\mathbb{R}$ is $[-e, e]$, which adds irrationals outside $\mathbb{Q}$.
3. So no $\mathbb{Q}$-neighbourhood has compact closure inside $\mathbb{Q}$: local compactness fails at every point.

So $\mathbb{Q}$ is nowhere locally compact; completeness of the ambient space matters.

## Worked: one-point compactification

Compactify a locally compact Hausdorff $X$:

1. Add one point at infinity, with neighbourhoods of infinity being complements of compact sets.
2. The result $X$ union $\{\infty\}$ is compact Hausdorff containing $X$ densely.
3. For instance $\mathbb{R}$ union $\{\infty\}$ is homeomorphic to $S^1$.

So local compactness is exactly what lets a space grow by one point into a compact Hausdorff space.
