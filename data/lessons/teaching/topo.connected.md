# Connectedness

**Connected space:** A space $X$ is connected if it cannot be written as the disjoint union of two nonempty open sets. Equivalently, the only subsets that are both open and closed are the empty set and $X$ itself.

## Definition

Decide whether $(0, 1)$ union $(2, 3)$ is connected:

1. Take $U = (0, 1)$ and $V = (2, 3)$; both are open in the subspace.
2. They are disjoint, nonempty, and their union is the whole space.
3. So the space splits as two disjoint nonempty opens: it is disconnected.

So the definition is a splitting test: one successful split into disjoint nonempty opens proves disconnectedness.

## Examples and characterizations

Show the rationals $\mathbb{Q}$ are totally disconnected:

1. Pick rationals $p$ below $q$ and an irrational $r$ strictly between them.
2. Then $(-\infty, r)$ and $(r, \infty)$ are disjoint opens in $\mathbb{Q}$ separating $p$ from $q$.
3. Since every pair splits this way, no two rationals share a connected piece.

So intervals in $\mathbb{R}$ are connected, while $\mathbb{Q}$ shatters completely: irrationals cut between any two rationals.

## Separation by open sets

Prove [0, 1] is connected by contradiction:

1. Suppose disjoint nonempty opens $U, V$ in [0, 1] cover it, with 0 in $U$.
2. Let $s$ be the least point of $V$; then $s$ is a limit of points of $U$ below it.
3. Every neighbourhood of $s$ meets both $U$ and $V$, so $s$ lies in neither open set, contradicting the cover.

So closed intervals are connected: to disconnect a space you must exhibit the separating opens, and here none exist.
