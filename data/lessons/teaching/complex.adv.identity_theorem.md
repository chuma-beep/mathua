# The Identity Theorem

## Zeros of analytic functions

Zeros of a nonzero analytic function on a domain are isolated. If the function vanishes at a point and is not identically zero, its power series has a first nonzero coefficient, so the zero has a finite order. In a small punctured neighbourhood the function is nonzero.

This isolation is a strong rigidity property. A nonzero analytic function cannot vanish on an open set, on a curve or on any set with an accumulation point inside the domain. If it did, the accumulation point would force all Taylor coefficients to vanish, and the function would be zero on a neighbourhood, then on the whole connected domain.

The order of a zero is the first index with a nonzero Taylor coefficient. It measures how flat the function is at the point. A zero of order one is simple, order two is double, and so on. The order controls the local behaviour and appears in the argument principle and in residue computations.

## Agreement on a set

If two analytic functions on a domain agree on a set with an accumulation point inside the domain, then they agree everywhere. Apply the theorem to their difference. The difference is analytic, vanishes on the set, and therefore has a non-isolated zero, so it must be identically zero.

The accumulation point must lie inside the domain, not merely on the boundary. Agreement on a boundary arc does not force agreement inside unless additional hypotheses such as continuity up to the boundary are present. This distinction is a common source of errors in uniqueness arguments.

The theorem justifies analytic continuation. If a function is defined by a power series on one disc and by another expression on an overlapping disc, and the two agree on the overlap, then they define the same analytic function on the union. The continuation is unique whenever the union is connected.

## Analytic continuation

Analytic continuation extends a function along paths by repeatedly expanding in power series. The identity theorem guarantees that different continuations along homotopic paths agree. Along non-homotopic paths, the continuation can differ, which produces monodromy and multivalued functions such as the logarithm.

The logarithm illustrates the phenomenon. Continuing the principal logarithm around the origin once returns a value shifted by two pi i. The function is not single-valued on the punctured plane, but it is single-valued on any simply connected domain that avoids the origin.

Uniqueness also underlies the standard method of proving identities. To show two analytic expressions are equal, prove they agree on a set with an accumulation point, such as an interval of real numbers. Polynomial identities, trigonometric identities and series identities can all be established this way.

