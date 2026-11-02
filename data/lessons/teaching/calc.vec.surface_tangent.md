# Tangent Planes to Parameterized Surfaces

## The Normal Vector

For a surface given as a graph \(z=f(x,y)\), a normal vector is \(\langle f_x,f_y,-1\rangle\). The first two components describe how the surface tilts in the coordinate directions, and the last component fixes the vertical scale. At a point \((x_0,y_0)\), substitute the point into the partial derivatives to obtain a numerical normal vector.

For a parameterized surface \(\mathbf r(u,v)\), the normal is the cross product \(\mathbf r_u\times\mathbf r_v\). The two tangent vectors span the tangent plane, and their cross product is perpendicular to both. The order of the cross product changes the sign of the normal but not the plane it defines.

## The Tangent Plane

The tangent plane at a point is the plane through that point with the normal vector just computed. For a graph, its equation can be written \(z=f(x_0,y_0)+f_x(x_0,y_0)(x-x_0)+f_y(x_0,y_0)(y-y_0)\). This is the linear approximation of the surface near the point, and it is the surface analogue of the tangent line.

The tangent plane is useful for estimation. Moving a small distance in the \(x\)-direction changes the height by approximately \(f_x\) times that distance, and similarly for \(y\). The approximation is exact only for a plane, and the error grows as the point moves away from the point of tangency.

## Parameterized Surfaces

A parameterized surface uses two parameters, often written \(u\) and \(v\). The partial derivatives with respect to those parameters are tangent to the surface, and their cross product gives a normal. If the cross product vanishes at a point, the parameterization is singular there and the tangent plane may not be defined by this method.

When the surface is given implicitly by \(F(x,y,z)=0\), a normal is the gradient \(\nabla F\). All three descriptions, graph, parameterization and implicit equation, lead to the same geometric plane. Choosing the description that makes the derivatives simplest is part of solving the problem efficiently.
