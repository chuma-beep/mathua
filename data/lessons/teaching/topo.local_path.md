# Locally Path-Connected Spaces

**Locally path-connected:** $X$ is locally path-connected if every point has a neighbourhood base of path-connected sets. Locally path-connected plus connected implies path-connected.

## Worked: Euclidean balls do the job

Check $\mathbb{R}^n$ at a point $x$:

1. Take any ball $B$ around $x$; balls form a neighbourhood base.
2. Any two points of $B$ join by the straight segment, which stays in $B$ by convexity.
3. Hence each basic neighbourhood is path-connected.

So manifolds are locally path-connected: they are locally Euclidean, hence locally ball-like.

## Worked: the bridge theorem

Prove connected plus locally path-connected gives path-connected:

1. Fix $x$; the set $U$ of points reachable from $x$ by paths is open, since path-connected neighbourhoods extend reachability.
2. The complement is likewise open: points outside $U$ have path-connected neighbourhoods missing $U$.
3. In a connected space one of the two opens is empty, so $U$ is everything.

So local path data assembles into global paths exactly when the space holds together.

## When local paths fail

Inspect the comb space at a baseline point:

1. Every neighbourhood of the baseline contains teeth segments plus gaps between teeth.
2. No such neighbourhood is path-connected: paths cannot cross the gaps near the baseline.
3. Yet the whole comb is path-connected via routes through the top handle.

So path-connected does not imply locally path-connected; the comb is the standard counterexample, and $\mathbb{Q}$ is neither.
