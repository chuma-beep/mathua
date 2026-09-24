# Orientability of Manifolds

**Orientable manifold:** A manifold with a consistent choice of local orientation, equivalently an atlas whose transition maps all have positive Jacobian. $S^2$, $T^2$, and $\mathbb{R}^n$ are orientable; the Mobius band and Klein bottle are not.

## Worked: the sphere is orientable

Orient $S^2$ by its outward normal field:

1. At each point $p$ take the outward unit normal $n(p)$, varying continuously with $p$.
2. The normal picks one of the two local sides consistently across overlapping charts.
3. All transition maps preserve this choice, so the atlas is oriented.

So $S^2$ is orientable: a global continuous normal field is a global orientation.

## Worked: the Mobius band is not

Track orientation around the central circle of the Mobius band:

1. Start with a small oriented frame at one point of the central circle.
2. Slide the frame once around the circle; the half-twist mirrors it.
3. The frame returns flipped, so no consistent global choice exists.

So the Mobius band is non-orientable: one trip around reverses orientation, and its double cover is the orientable cylinder.

## Worked: the orientability test

Decide orientability of a surface in practice:

1. Look for an embedded Mobius band: finding one proves non-orientability.
2. Equivalently compute the first Stiefel-Whitney class $w_1$; it vanishes exactly for orientable bundles.
3. The torus contains no Mobius band and has $w_1=0$, so $T^2$ is orientable.

So orientability is detected by loops: orientation-reversing loops are the obstruction.
