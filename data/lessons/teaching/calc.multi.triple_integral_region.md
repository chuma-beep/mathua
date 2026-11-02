# Triple Integrals over General Regions

## Describing a General Region

A general three-dimensional region is usually described by inequalities rather than a single box. A common example is the tetrahedron cut from the first octant by a plane such as \(x/a+y/b+z/c=1\). The outer integral runs over the range of one variable, and the inner limits are functions of the variables that remain. Writing the region as a set of inequalities is the step that turns a picture into an iterated integral.

The order of integration is a choice. A region that is awkward in the order \(dz\,dy\,dx\) may become simple after projecting onto a different coordinate plane. Sketch the projection, identify the curved boundary, and let the inner variable run between the lower and upper surfaces above that projection.

## Setting Up the Integral

For a region between a lower surface \(z=g(x,y)\) and an upper surface \(z=h(x,y)\), the inner integral runs from \(g(x,y)\) to \(h(x,y)\). The remaining double integral runs over the projection of the region in the \(xy\)-plane. When the region is a tetrahedron, the projection is a triangle and the upper surface is the plane that cuts the tetrahedron from the box.

A volume is the triple integral of the constant function one. A mass is the triple integral of density, and an average value is the integral of the function divided by the volume. Keeping the meaning in view prevents the common mistake of integrating the wrong quantity after the limits have been set.

## Evaluating and Checking

Evaluate the inner integral first, treating the other variables as constants. The result is a double integral, which is then evaluated by the same procedure as before. Check the answer by comparing it with a simpler volume or by reversing the order of integration when both orders are possible. A negative volume, or a volume larger than a surrounding box, is a signal that a limit was reversed.

For a tetrahedron with intercepts \(a\), \(b\) and \(c\), the volume is \(abc/6\). That formula is useful as a check even when the exercise asks for a different integrand. If the integrand is a constant density, the mass is the density times the volume, so the two calculations should agree.
