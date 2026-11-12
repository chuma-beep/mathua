# The Dirac delta in several variables

The Dirac delta in several variables is a unit source concentrated at a point.

## The defining property

Integrating \(\delta(x-a)\) against a test function returns the value at \(a\).

In two variables the delta is a product of one-dimensional deltas, so the total mass is 1. Integrating \(\delta(x-a)\delta(y-b)\phi(x,y)\) over the plane gives \(\phi(a,b)\), the value at the chosen point.

For \(\phi=1+x+y\) at the point \((1,2)\) the result is 4.

## Scaling and translation

Scaling brings the absolute value of the Jacobian determinant, so in \(n\) dimensions \(\delta(cx)\) equals \(\delta(x)/|c|^{n}\).

Translation moves the point without changing the mass. In the plane \(\delta(2x)\) contributes a factor \(1/4\), so its integral against the constant 8 is 2.

## The delta as a measure

The delta is a measure that gives mass 1 to any set containing the origin and 0 to any set avoiding it.

Composing with a smooth function \(g\) produces a measure on the level set, weighted by \(1/|\nabla g|\). For \(g=x^{2}+y^{2}-1\) the level set is the unit circle, where \(|\nabla g|=2\).

The mass there is the length \(2\pi\) divided by 2, which is \(\pi\).

The measure reading is the one that survives approximation.
