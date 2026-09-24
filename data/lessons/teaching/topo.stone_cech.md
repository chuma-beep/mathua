# Stone-Cech Compactification

**Stone-Cech compactification:** For Tychonoff $X$, the largest compactification $\beta X$: every continuous $f:X\to K$ into compact Hausdorff $K$ extends uniquely to $\beta X$. For discrete $X$, points of $\beta X$ are ultrafilters.

## Worked: the universal property

Extend $f:\mathbb{N}\to[0, 1]$ to $\beta\mathbb{N}$:

1. $\mathbb{N}$ is discrete, hence Tychonoff, so $\beta\mathbb{N}$ exists with $\mathbb{N}$ dense in it.
2. $[0, 1]$ is compact Hausdorff, so universality hands a unique continuous extension $F:\beta\mathbb{N}\to[0, 1]$.
3. Uniqueness follows from density: two extensions agreeing on $\mathbb{N}$ agree everywhere.

So $\beta X$ is largest in the sense that all maps to compact Hausdorff spaces factor through it.

## Worked: growth of beta N

Describe the remainder $\mathbb{N}^*=\beta\mathbb{N}$ minus $\mathbb{N}$:

1. Free ultrafilters on $\mathbb{N}$ are the new points; principal ultrafilters are the old integers.
2. The remainder is nonempty, compact, and huge: $\beta\mathbb{N}$ has doubly-exponentially many points.
3. Hence $\beta\mathbb{R}$ is even larger, and neither is metrizable.

So Stone-Cech compactifications are enormous: a countable space grows a vast corona.

## Worked: who gets compactified

Test which spaces admit a Hausdorff compactification:

1. If $X$ sits inside some compact Hausdorff $K$, it inherits complete regularity from $K$.
2. Conversely every Tychonoff $X$ embeds in its own $\beta X$.
3. Hence completely regular (Tychonoff) is exactly the embeddability condition.

So Tychonoff means compactifiable: $\beta X$ is the universal compact Hausdorff home.
