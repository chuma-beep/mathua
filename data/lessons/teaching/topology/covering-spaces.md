# Covering Spaces

A covering space unwraps a space the way a spiral staircase unwraps a circle:
locally it looks identical, globally it can be much simpler. The model example
is the exponential map \(p: \mathbb{R} \to S^1\) given by
\(p(t) = (\cos 2\pi t, \sin 2\pi t)\), which winds the line around the circle
once per unit interval.

## Evenly Covered Neighborhoods

A map \(p: E \to B\) is a **covering** if every \(x \in B\) has a neighborhood
\(U\) whose preimage \(p^{-1}(U)\) is a disjoint union of open sets, each mapped
homeomorphically onto \(U\) by \(p\). Each such sheet looks exactly like \(U\).

For \(p(t) = (\cos 2\pi t, \sin 2\pi t)\), any arc \(U \subset S^1\) shorter than
the full circle lifts to disjoint intervals \(\dots, (a-1, b-1), (a, b),
(a+1, b+1), \dots\), one per winding. The fiber \(p^{-1}(x)\) over any point is
a copy of \(\mathbb{Z}\). A covering with \(n\)-point fibers is called
**\(n\)-sheeted**; the circle also has a connected 2-sheeted cover
\(z \mapsto z^2\).

## Path Lifting

Coverings are valued for one property: **paths and homotopies lift uniquely**.
Given a path \(\gamma\) in \(B\) and a starting lift \(\tilde{x}_0 \in E\), there
is exactly one lifted path \(\tilde{\gamma}\) in \(E\) with
\(p \circ \tilde{\gamma} = \gamma\). Homotopies of paths lift too, so the
endpoint of a lifted loop depends only on the loop's homotopy class.

Consequence: a loop in \(B\) lifts to a loop in \(E\) exactly when its class
lies in the subgroup \(p_*(\pi_1(E)) \subseteq \pi_1(B)\). For the exponential
covering, the loop winding \(n\) times lifts to the path from \(0\) to \(n\),
which closes up only for \(n = 0\). This is the standard proof that
\(\pi_1(S^1) \cong \mathbb{Z}\).

## Deck Transformations and Classification

A **deck transformation** is a homeomorphism \(d: E \to E\) with
\(p \circ d = p\): a symmetry permuting the sheets. For
\(p: \mathbb{R} \to S^1\) these are the integer translations
\(t \mapsto t + n\), so the deck group is \(\mathbb{Z}\), matching
\(\pi_1(S^1)\).

In general, connected coverings of a nice space \(B\) are classified by
subgroups of \(\pi_1(B)\) up to conjugacy: every subgroup occurs, and the
simply connected **universal cover** corresponds to the trivial subgroup. The
2-sphere \(S^2\) is simply connected, so every covering of \(S^2\) is a
homeomorphism. Real projective space \(\mathbb{RP}^n\) has universal cover
\(S^n\) with deck group \(\mathbb{Z}/2\).
