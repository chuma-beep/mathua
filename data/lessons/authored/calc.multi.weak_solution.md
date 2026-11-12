# Weak solutions of elliptic problems

A weak solution satisfies the equation after multiplication by a test function and integration by parts.

## The weak formulation

For \(-\Delta u=f\) with zero boundary values the weak form asks that the integral of \(\nabla u\cdot\nabla v\) equals the integral of \(fv\), for every test \(v\) vanishing on the boundary.

The derivatives fall on the test function, so less regularity is needed. Every classical solution satisfies the identity, so on the unit interval with \(f=1\) the solution \(x(1-x)/2\) does so.

At the midpoint that solution equals 0.125.

## Existence by Lax-Milgram

The form is bounded and coercive on the Sobolev space with zero boundary values.

Coercivity follows from Poincare's inequality, which bounds the function by its gradient. Lax-Milgram then gives a unique solution for every \(f\), as an abstract minimization of energy.

For the unit interval the energy of \(x(1-x)/2\) is the integral of \((1/2-x)^{2}\), which is \(1/12\).

## A one-dimensional example

Take \(-u''=1\) on the interval from 0 to 1 with zero ends, whose solution is \(x(1-x)/2\).

Check the weak form with \(v=x(1-x)\). The integral of \(u'v'\) is the integral of \((1/2-x)(1-2x)\), which is \(1/6\).

The integral of \(fv\) is also \(1/6\), so the two sides agree.
