# Normal Lines to Surfaces

## The Direction of the Normal

The normal line to a surface at a point is the line through the point in the direction of a normal vector. For a graph \(z=f(x,y)\), the normal direction is \(\langle f_x,f_y,-1\rangle\). The tangent plane is perpendicular to this line, so the two objects share the same normal direction.

At a point \((x_0,y_0)\), the normal vector is obtained by evaluating the partial derivatives there. For \(f(x,y)=ax^2+by^2\), the normal is \(\langle 2ax_0,2by_0,-1\rangle\). The negative one in the third component reflects the choice of writing the surface as \(z-f(x,y)=0\).

## Parametric Equations

A line through \((x_0,y_0,z_0)\) with direction \(\langle A,B,C\rangle\) has parametric equations \(x=x_0+At\), \(y=y_0+Bt\), \(z=z_0+Ct\). Substituting the normal direction gives the normal line. Different values of \(t\) give different points on the line, and the point of tangency corresponds to \(t=0\).

The direction vector can be scaled by any nonzero constant without changing the line. A normal of \(\langle 2,4,-1\rangle\) and a normal of \(\langle -2,-4,1\rangle\) describe the same line with opposite orientation. The parametric equations change, but the set of points does not.

## Points Along the Line

To find a point on the normal line, substitute a value of \(t\) into all three equations. At \(t=1\), each coordinate has moved one full direction vector from the point of tangency. This is often the quickest way to answer a question about a specific point rather than the whole line.

The normal line is useful for measuring distance from a surface and for constructing the tangent plane. It is not the same as the normal vector: the vector gives direction and magnitude, while the line is the infinite set of points through the surface point. Keeping that distinction clear prevents sign and scaling errors.
