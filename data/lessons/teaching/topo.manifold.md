# Topological Manifolds

**Topological manifold:** A Hausdorff, second-countable space that is locally Euclidean: every point has a neighbourhood homeomorphic to $\mathbb{R}^n$. The dimension $n$ is a topological invariant.

## Worked: the circle as a 1-manifold

Cover $S^1$ by Euclidean neighbourhoods:

1. Remove the point 1: the remaining arc maps homeomorphically to $\mathbb{R}$ by stereographic projection.
2. Rotate: removing any other point gives a second chart covering the missing point 1.
3. Every point lies in some chart, so $S^1$ is locally Euclidean of dimension 1.

So the circle is a compact 1-manifold: locally an interval, globally a loop.

## Worked: the sphere as a 2-manifold

Chart $S^2$ the same way:

1. Remove the north pole: stereographic projection gives a homeomorphism with the plane $\mathbb{R}^2$.
2. Remove the south pole for a second chart covering the north pole.
3. Both charts are Euclidean patches, so $S^2$ is locally $\mathbb{R}^2$.

So the sphere is a compact 2-manifold, and the torus $T^2=S^1\times S^1$ is another: products of manifolds are manifolds.

## When local Euclideanness fails

Test the cross (two lines meeting at the origin) in the plane:

1. Delete the crossing point: the remainder has four components.
2. Deleting any point of $\mathbb{R}$ leaves at most two components.
3. No neighbourhood of the crossing is homeomorphic to $\mathbb{R}$, so the cross is not a 1-manifold.

So branch points obstruct manifolds: locally Euclidean means no junctions allowed.
