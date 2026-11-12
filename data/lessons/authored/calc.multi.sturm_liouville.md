# Sturm-Liouville problems

A Sturm-Liouville problem is an eigenvalue problem for a second-order equation in self-adjoint form.

## The Sturm-Liouville form

The equation is \(-(py')'+qy=\lambda wy\), with \(p\) and the weight \(w\) positive on the interval.

Writing it this way makes the self-adjointness visible, so integration by parts moves derivatives from one eigenfunction to the other without a boundary term.

With \(p=1\), \(q=0\) and \(w=1\) the equation is \(-y''=\lambda y\). Then \(\lambda=4\) has the solution \(\sin 2x\).

## Eigenvalues and eigenfunctions

The eigenvalues are real, increasing and tending to infinity. Eigenfunctions for different eigenvalues are orthogonal under the weighted integral.

Each eigenvalue has a one-dimensional space of eigenfunctions. On an interval of length \(\pi\) the eigenvalues are \(n^{2}\), so the first three are 1, 4 and 9.

## An example on an interval

Take \(p=1\), \(q=0\) and \(w=1\) on the interval from 0 to 2 with zero end values.

The eigenvalues are \((n\pi/2)^{2}\) and the functions are \(\sin(n\pi x/2)\). The first eigenvalue is \(\pi^{2}/4\), about 2.467. The first function has weighted norm squared equal to 1.

Its integral against the second function is zero, as orthogonality promises.

The weight changes the numbers, not the structure. The general theory promises as much.
