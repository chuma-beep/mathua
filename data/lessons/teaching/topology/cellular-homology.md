# Cellular Homology

Cellular homology counts holes by building a space from balls and doing linear
algebra on how they attach. A **CW complex** is assembled inductively: start
with 0-cells (points), attach 1-cells (intervals) by gluing endpoints to the
0-skeleton, attach 2-cells (disks) along loops in the 1-skeleton, and so on.
Spheres, tori, projective spaces, and most spaces in practice admit small CW
structures.

## The Cellular Chain Complex

Let \(C_n(X)\) be the free abelian group on the \(n\)-cells of \(X\). The
**cellular boundary** \(d_n: C_n \to C_{n-1}\) records, for each \(n\)-cell,
the degree with which its attaching map wraps each \((n-1)\)-cell. These maps
satisfy \(d_{n-1} \circ d_n = 0\), forming a chain complex
\(\dots \to C_2 \xrightarrow{d_2} C_1 \xrightarrow{d_1} C_0 \to 0\).

## Cycles Modulo Boundaries

Cellular homology is \(H_n(X) = \ker d_n / \operatorname{im} d_{n+1}\): cycles
modulo boundaries. A class survives exactly when it wraps a hole no higher
cell fills. Because \(d^2 = 0\) is automatic, computing homology is Smith
normal form on integer matrices, which is why this is the computable version
of the theory.

## Sample Computations

The \(n\)-sphere \(S^n\) has one 0-cell and one \(n\)-cell with constant
attaching map (degree 0), so \(H_n(S^n) \cong \mathbb{Z}\) and all other
reduced groups vanish: one \(n\)-dimensional hole, nothing else.

The torus \(T\) has one 0-cell, two 1-cells \(a, b\), and one 2-cell attached
along the commutator \(aba^{-1}b^{-1}\). Abelianizing kills the commutator, so
\(d_2 = 0\) and \(H_1(T) \cong \mathbb{Z} \oplus \mathbb{Z}\),
\(H_2(T) \cong \mathbb{Z}\): two independent loops and one enclosed void. The
Klein bottle differs by one sign in the attaching word, giving
\(H_1 \cong \mathbb{Z} \oplus \mathbb{Z}/2\): torsion records the twist.
