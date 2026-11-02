# The Riemann Mapping Theorem

## Conformal equivalence

Two domains are conformally equivalent when a one-to-one analytic map sends one onto the other with an analytic inverse. Conformal equivalence preserves angles, harmonic functions and the local structure of analytic functions. The unit disc is the model domain for the simply connected case.

The Riemann mapping theorem states that every simply connected proper domain in the complex plane is conformally equivalent to the open unit disc. Proper means not the whole plane. Simply connected means no holes, equivalently that every closed curve can be contracted to a point inside the domain.

The theorem is an existence result. It does not usually give a formula for the map, and explicit maps are known only for special domains such as half-planes, strips and polygons. For polygons, the Schwarz-Christoffel formula provides the map in terms of an integral.

## The mapping theorem

The proof has two parts: a normal family argument shows that a suitable extremal map exists, and a boundary argument shows that its image cannot omit a point. The extremal map maximizes the derivative at a chosen interior point among maps into the disc that send that point to zero.

Uniqueness follows from the Schwarz lemma. If two maps have the same normalization, their composition with the inverse of the other is an automorphism of the disc fixing the origin. Such an automorphism is a rotation, and the derivative normalization forces the rotation to be the identity.

The theorem fails without simple connectivity. An annulus is not conformally equivalent to the disc, because it has a hole and the fundamental group is infinite cyclic. The classification of annuli up to conformal equivalence is by a real modulus, which is a much finer invariant than topological type.

## Boundary behaviour

The interior map does not automatically extend nicely to the boundary. For domains bounded by Jordan curves, Caratheodory's theorem guarantees that the map extends to a homeomorphism of the closures. For more complicated boundaries, the extension can fail at a dense set of points.

Boundary behaviour matters for applications because boundary values determine harmonic functions. If the map extends continuously, solving a boundary value problem on a complicated domain can be reduced to solving it on the disc and pulling back. This is the classical method for the Dirichlet problem on simply connected domains.

The theorem is also a warning about invariants. Angles, harmonicity and analytic structure are preserved, but distances and areas are not. A conformal map can distort lengths dramatically near the boundary. Quantitative control of that distortion is the subject of distortion theorems and quasiconformal mapping.

