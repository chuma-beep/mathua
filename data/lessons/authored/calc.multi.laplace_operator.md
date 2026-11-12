# The Laplacian

The Laplacian of a scalar field is the sum of its pure second partial derivatives.

## The sum of second partials

The two-dimensional Laplacian is \(f_{xx}+f_{yy}\), written \(\Delta f\), with a third term \(f_{zz}\) in three dimensions.

Each term differentiates in one variable and holds the others fixed. A value above its surrounding average has a negative Laplacian.

For \(f(x,y)=x^{3}y^{2}\) the second derivatives are \(6xy^{2}\) and \(2x^{3}\), so at \((1,2)\) the sum is 24 plus 2, which is 26.

## The Laplacian in polar coordinates

In polar coordinates the sum becomes \(f_{rr}+\frac{1}{r}f_{r}+\frac{1}{r^{2}}f_{\theta\theta}\).

The middle term carries the curvature of the coordinate lines, so a radial function still has a nonzero Laplacian. For \(f=r^{2}\) the formula gives 2 plus 2, which is 4, matching \(\Delta(x^{2}+y^{2})\).

For \(f=\ln r\) the terms cancel and the Laplacian is zero at every radius.

## A worked computation

Take the upward bowl \(f(x,y)=x^{2}+3y^{2}\), whose second derivatives are 2 and 6.

The Laplacian is their sum, so it is 8 at every point. At \((2,1)\) the value is still 8. The function \(x^{2}-y^{2}\) is harmonic instead, since 2 and -2 cancel.

The Laplacian is the divergence of the gradient, which is why it governs every quantity that spreads out.
