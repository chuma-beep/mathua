# The Mayer-Vietoris sequence

## The Exact Sequence

Mayer-Vietoris relates the homology of a union to the homology of its pieces and their intersection. For open sets A and B covering X, there is a long exact sequence involving the homology of A, B and their intersection.

The sequence is natural and works for any homology theory satisfying the Eilenberg-Steenrod axioms. It is the main computational tool for splitting a space into simpler pieces, and it is often faster than a direct cellular calculation.

The hypothesis that the interiors cover X is essential. For closed sets the sequence can fail unless the pair is excisive. In practice one chooses open covers or good pairs to satisfy the hypothesis.

## Connecting Homomorphisms

The connecting homomorphism lowers degree by one. It measures the failure of a cycle in the intersection to be a boundary in either piece. Its image and kernel are exactly the terms that correct the naive direct sum.

The map is defined by a diagram chase. A cycle in the intersection is a cycle in both A and B; if it is a boundary in X, the chain witnessing the boundary splits into two parts, and the difference is a cycle in the intersection of one degree lower.

Naturality of the connecting map is what makes Mayer-Vietoris compatible with maps of pairs. It also makes the sequence computable: one can chase explicit cycles through the diagram for small examples.

## Computing Homology

The torus can be decomposed into two overlapping cylinders whose intersection is two disjoint annuli. Mayer-Vietoris then reproduces the homology of the torus without triangulating it, and the same method handles surfaces of higher genus.

The sequence also computes the homology of suspensions and of wedges. For a wedge of two spaces with a common point, the reduced homology is the direct sum of the reduced homology of the pieces, and Mayer-Vietoris recovers that statement.

In cohomology the sequence runs in the opposite direction, and the connecting map raises degree. The two versions are equivalent under universal coefficients when the coefficient ring is a field, but they differ integrally in the presence of torsion.

