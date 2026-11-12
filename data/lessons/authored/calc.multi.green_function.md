# Green functions for the Laplacian

A Green function is the response of a boundary value problem to a unit point source.

## The defining property

The Green function of a region satisfies \(-\Delta G=\delta\) inside and \(G=0\) on the boundary.

The delta is a unit source at one point, so \(G\) is the field that source produces while the boundary stays at zero. The solution of \(-\Delta u=f\) is then \(u(x)=\int G(x,y)f(y)\,dy\).

In the plane a source of strength 4 at the origin gives \(u=-(2/\pi)\ln r\), since the kernel is \(-\frac{1}{2\pi}\ln r\).

## The fundamental solution

On the whole plane the fundamental solution is \(-\frac{1}{2\pi}\ln r\). In space it is \(\frac{1}{4\pi r}\).

Both are singular at the origin and harmonic elsewhere. The constant is fixed by the flux through a small sphere.

For the space kernel the gradient has length \(\frac{1}{4\pi r^{2}}\) while the sphere of radius \(r\) has area \(4\pi r^{2}\), so the product is 1.

## A Green function on a disk

On the unit disk the Green function is the fundamental solution minus a correction that vanishes on the circle.

At the centre the kernel is \(-\frac{1}{2\pi}\ln r\), which at \(r=0.5\) equals 0.1103. The problem \(-\Delta u=1\) has the solution \(u=(1-r^{2})/4\), equal to 0.25 at the centre.

A Green function is the continuous analogue of an inverse matrix.
