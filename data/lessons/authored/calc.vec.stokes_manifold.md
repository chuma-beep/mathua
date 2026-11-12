# Stokes' theorem on manifolds

Stokes' theorem on manifolds says the integral of an exact form over a region equals the integral over its boundary.

## The boundary operator

The boundary of a region is a manifold of one lower dimension, which inherits an orientation.

The region lies to the left of the boundary direction. The boundary of a boundary is empty. That fact, written \(\partial\partial=0\), matches \(d^{2}=0\) for the derivative.

The unit disk has the unit circle as its boundary. The circle itself has empty boundary.

## The general statement

The integral of \(d\omega\) over a region equals the integral of \(\omega\) over its boundary.

Nothing else is needed: no metric, no coordinates and no special shape. Green's theorem, the divergence theorem and the classical Stokes theorem all follow by choosing the degree.

On the unit disk the form \(-y\,dx+x\,dy\) has derivative \(2\,dx\wedge dy\), so both sides give \(2\pi\).

## A worked boundary integral

Take the unit square and the one-form \(\omega=x\,dy\), whose derivative is \(dx\wedge dy\).

The integral of the derivative over the square is 1. On the boundary only the right edge contributes, where \(x\) is 1 and \(y\) runs from 0 to 1.

The other three edges give nothing, since on them \(dy\) vanishes or \(x\) does. A boundary computation has replaced an area computation.
