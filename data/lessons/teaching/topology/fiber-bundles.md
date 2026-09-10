# Fiber Bundles

A fiber bundle is a space that looks locally like a product but may twist
globally. It consists of a total space \(E\), a base \(B\), a fiber \(F\), and
a projection \(\pi: E \to B\) such that every \(x \in B\) has a neighborhood
\(U\) with \(\pi^{-1}(U)\) homeomorphic to \(U \times F\) via a
fiber-preserving homeomorphism (a **local trivialization**).

## Local Triviality and Sections

Local triviality means: zoom in anywhere on the base and the bundle is just a
product. A **section** is a continuous map \(s: B \to E\) with
\(\pi \circ s = \mathrm{id}_B\): a continuous choice of one point per fiber.
The cylinder \(S^1 \times [0,1]\) over \(S^1\) admits sections (pick a constant
height). Whether sections exist globally is a first measure of twisting.

A covering space is the special case of discrete fiber.

## Vector Bundles and Twisting

A **vector bundle** has vector-space fibers with linear transition maps; the
tangent bundle \(TM\) of a manifold assigns to each point its tangent space.
The Möbius strip is the nontrivial line bundle over \(S^1\): locally an
interval times an arc, globally given a half-twist, and it admits no
nowhere-zero section. The twist is visible in the transition data, not in any
single trivialization.

## Clutching and the Hopf Fibration

Bundles over spheres are built by **clutching**: glue trivial bundles over the
two hemispheres along the equator with a transition map
\(S^{n-1} \to G\) into the structure group. The Möbius strip is two trivial
interval-bundles over semicircles glued with a flip on one overlap arc.

The showcase is the **Hopf fibration** \(\pi: S^3 \to S^2\) with circle fibers:
write \(S^3 \subset \mathbb{C}^2\) as pairs \((z_1, z_2)\) with
\(|z_1|^2 + |z_2|^2 = 1\) and send \((z_1, z_2)\) to the ratio \(z_1/z_2\) in
\(\mathbb{CP}^1 \cong S^2\). Every fiber is a circle, any two fibers link once,
and the bundle is nontrivial: \(S^3\) is not \(S^2 \times S^1\).
