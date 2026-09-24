# Homotopy and Fundamental Group

**Homotopy:** A continuous deformation $H:X\times[0, 1]\to Y$ with $H(x, 0)=f(x)$ and $H(x, 1)=g(x)$. Loops up to homotopy form the fundamental group $\pi_1(X)$ under concatenation.

## Worked: winding number on the circle

Classify a loop $w$ on $S^1$ going twice around:

1. Lift $w$ to the universal cover $p:\mathbb{R}\to S^1$ starting at 0; the lift ends at 2.
2. A loop going once lifts to end at 1; going $n$ times ends at $n$.
3. Homotopic loops lift to the same endpoint, so the endpoint integer is the homotopy invariant.

So $\pi_1(S^1)$ is $\mathbb{Z}$: each loop carries an integer winding number, and homotopic loops share it.

## Worked: the sphere is simply connected

Contract a loop $w$ on $S^2$:

1. $w$ misses at least one point $p$ after a small perturbation (a loop is one-dimensional).
2. $S^2$ minus $p$ is homeomorphic to the plane, where the loop shrinks to a point.
3. Hence $w$ is null-homotopic in $S^2$.

So $\pi_1(S^2)$ is trivial: every loop contracts, unlike on the circle.

## Worked: maps into Euclidean space

Compare two maps $f, g:X\to\mathbb{R}^n$:

1. Set $H(x, t)=(1-t)f(x)+tg(x)$, the straight-line homotopy.
2. Each $H(x, t)$ stays in $\mathbb{R}^n$ since it is a vector space.
3. $H$ runs continuously from $f$ to $g$.

So any two maps into $\mathbb{R}^n$ are homotopic: Euclidean targets see no homotopy distinction.
