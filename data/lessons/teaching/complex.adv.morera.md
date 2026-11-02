# Morera's Theorem

## The converse of Cauchy

Morera's theorem is the converse of Cauchy's theorem. If a function is continuous on a domain and its integral over every triangle inside the domain is zero, then the function is analytic. Cauchy's theorem says analytic functions have this property; Morera says the property characterizes them.

The triangle hypothesis is convenient because triangles are convex and the proof can use a primitive. Define a primitive by integrating from a fixed base point along a path. The triangle condition makes the integral path-independent, and the resulting primitive has derivative equal to the original function.

Once a primitive exists, the original function is the derivative of an analytic function, hence analytic. This argument also shows that Morera's theorem implies the existence of primitives on simply connected domains. The triangle condition is a local substitute for simple connectivity.

## Triangles and analyticity

The theorem only needs the triangle condition, not the integral over every closed curve. On a convex domain, every closed polygonal curve can be decomposed into triangles, so the triangle condition extends to all closed curves. This is why triangles are the right building blocks.

Continuity is essential. A function that is not continuous can have zero triangle integrals and still fail to be analytic. The proof uses continuity to compare the primitive with a difference quotient, and without it the limiting argument breaks down.

Morera is often used to prove that a limit of analytic functions is analytic. If a sequence converges uniformly on compact sets, the limit is continuous, and its triangle integrals are limits of the integrals of the approximants. Each approximant has zero triangle integral by Cauchy's theorem, so the limit does too.

## Using Morera

To apply Morera, verify continuity first, then show that the integral over every triangle vanishes. In practice, one shows that the function is a uniform limit of analytic functions, or that it has a primitive. Both routes are common in convergence theorems and in constructions of analytic functions by series.

Morera also proves that a continuous function with continuous partial derivatives satisfying the Cauchy-Riemann equations is analytic. The triangle integral can be converted into a boundary integral, and the Cauchy-Riemann equations make the integrand exact, so the integral vanishes.

A typical exam question gives a piecewise definition and asks whether the function is analytic at a junction. Morera can establish analyticity on each piece, but continuity and the matching of derivatives at the junction require separate checks. The theorem does not replace those checks.

