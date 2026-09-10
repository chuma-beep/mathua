# De Rham Cohomology

De Rham cohomology detects holes with calculus: integrate differential forms
over cycles and ask which integrals are forced by topology rather than by the
choice of form. On a smooth manifold \(M\), a \(k\)-form \(\omega\) assigns to
each point an alternating \(k\)-linear map on tangent vectors; the exterior
derivative \(d: \Omega^k \to \Omega^{k+1}\) extends gradient, curl, and
divergence to all degrees with the golden rule \(d^2 = 0\).

## Closed Forms and Exact Forms

A form is **closed** if \(d\omega = 0\) and **exact** if \(\omega = d\eta\) for
some \(\eta\). Every exact form is closed (\(d^2 = 0\)), but the converse fails
exactly where the manifold has holes. The **Poincaré lemma** says the converse
holds locally: on any contractible open set (a ball, all of
\(\mathbb{R}^n\)), every closed form is exact.

## A Hole You Can Integrate

The angular form \(d\theta = (-y\,dx + x\,dy)/(x^2+y^2)\) on
\(\mathbb{R}^2 \setminus \{0\}\) is closed but not exact: its integral over the
unit circle is \(2\pi\), while the integral of any exact form over a closed
loop is 0 by Stokes' theorem. The puncture is what makes the difference.

## Cohomology Groups

The **\(k\)-th de Rham group** is \(H^k_{dR}(M) = \{\text{closed } k\text{-forms}\} /
\{\text{exact } k\text{-forms}\}\): closed forms modulo the trivial ones.
\(H^0_{dR}(M) \cong \mathbb{R}^c\) where \(c\) is the number of connected
components (closed 0-forms are locally constant functions). For the circle,
\([d\theta]\) generates \(H^1_{dR}(S^1) \cong \mathbb{R}\): one independent
hole, measured by winding-number integrals. De Rham's theorem identifies these
groups with singular cohomology over \(\mathbb{R}\), so calculus sees the same
holes as combinatorics.
