# Orthogonal projection

Orthogonal projection finds the point of a closed subspace closest to a given vector. In a Hilbert space the closest point is unique and the error is orthogonal to the subspace.

## Closest points

Take the subspace \(M\) spanned by \((1,1)\) in \(\mathbb{R}^2\) and the point \(x = (2,0)\). The closest point of \(M\) is \((1,1)\), and the distance is \(\sqrt{2} \approx 1.414\). Any other multiple of \((1,1)\) is farther away, so the nearest point is unique.

## The projection theorem

For a closed subspace \(M\) of a Hilbert space every \(x\) splits uniquely as \(Px + (x - Px)\), with \(Px\) in \(M\) and \(x - Px\) orthogonal to \(M\). Here \(x - Px = (1,-1)\), whose inner product with \((1,1)\) is 0. The map \(P\) is linear and has norm at most 1.

## Orthonormal expansions

With an orthonormal basis \(e_n\), the projection onto the span of the first \(N\) vectors is \(\sum_{n \le N} \langle x, e_n\rangle e_n\). For \(x = (2,0)\) and \(e_1 = (1,1)/\sqrt{2}\), the coefficient is \(\sqrt{2}\) and the projection is \((1,1)\). Best approximation and Fourier expansion are the same operation.

Projection replaces a vector by its shadow, and orthogonality is what makes the shadow the closest point.
