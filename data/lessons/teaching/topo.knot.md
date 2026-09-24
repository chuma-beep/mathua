# Knot Theory Basics

**Knot:** An embedding of $S^1$ in $\mathbb{R}^3$ (or $S^3$). The unknot bounds a disc; the trefoil with 3 crossings is the simplest nontrivial knot; the knot group $\pi_1$ of the complement is a knot invariant.

## Worked: the trefoil is knotted

Show the trefoil $T$ is not the unknot:

1. Count crossings in the standard diagram: 3, minimal over all diagrams of $T$.
2. The unknot has a diagram with 0 crossings, and crossing number is a knot invariant.
3. Since 3 differs from 0, $T$ cannot be deformed into the unknot.

So crossing number already separates the trefoil from the unknot: the simplest invariant does real work.

## Worked: the knot group

Compute the invariant for the unknot versus the trefoil:

1. The unknot complement deformation-retracts to a solid torus, so its group is $\mathbb{Z}$.
2. The trefoil complement has group with presentation $\langle x,y\mid x^2=y^3\rangle$, nonabelian.
3. Nonabelian differs from $\mathbb{Z}$, confirming the trefoil is knotted independently of diagrams.

So the knot group sees knottedness algebraically: abelian $\mathbb{Z}$ for the unknot, nonabelian for the trefoil.

## Worked: polynomials distinguish

Tell the trefoil from the figure-eight knot:

1. The trefoil is $3_1$ (3 crossings); the figure-eight is $4_1$ (4 crossings).
2. Their Jones polynomials differ: the trefoil's is nontrivially asymmetric while the figure-eight's is symmetric.
3. Gordon-Luecke adds the capstone: the complement (hence its group) determines the knot completely.

So invariants stack up: crossing numbers, Jones polynomials, and finally the whole complement classify knots.
