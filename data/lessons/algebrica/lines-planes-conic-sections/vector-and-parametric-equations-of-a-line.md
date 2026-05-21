> Content sourced from [Algebrica](https://algebrica.org/vector-and-parametric-equations-of-a-line/) — CC BY-NC 4.0

## From Cartesian to vector form

We’ve seen that the equation of a line can be written in [Cartesian form](<../lines>), such as the explicit equation: \\[y = mx + q \\]

where the line is described using coordinates in the plane. However, lines can also be expressed in **vector form** , using points and direction [vectors](<../vectors/>) to represent their position and orientation.


Let’s consider a directed vector \\( \vec{v} \\). The line passing through a point \\( P_0 \\) and parallel to \\( \vec{v} \\) consists of all points \\(P\\) such that the vector \\(P - P_0 \\) is parallel to \\( \vec{v} \\). The **vector equation** of a line is given by:

\\[P - P_0 = t \vec{v} \quad \text{with} \quad t \in \mathbb{R} \\]

![Vector equations of a line.](/diagrams/algebrica/vector-equation-line-1.png)

Here, \\( P_0 \\) is a fixed point on the line, \\( \vec{v} \\) is a direction vector, and \\( t \\) is a real parameter. The point \\( P \\) varies along the line as \\( t \\) changes.

##### The vector representation of a line expresses all its points in a compact and geometric form, making the direction and position immediately clear, and serves as a natural bridge between coordinate geometry and linear algebra.

## Parametric form

Let’s now fix an orthonormal basis, that is a reference system where the axes are perpendicular to each other and the vectors that define them have unit length. In practice, this means we’re working in the standard Cartesian plane, where the x-axis is aligned with \\( \vec{i} = (1, 0) \\) and the y-axis with \\( \vec{j} = (0, 1) \\).

Let’s represent the fixed point \\( P_0 \\) and a generic point \\( P \\) using coordinates:

\\[P_0 = (x_0, y_0) \quad \text{and} \quad P = (x, y) \\]

Using the vector equation of the line, we obtain:

\\[\begin{align} P - P_0 &= (P - O) - (P_0 - O)\\\\[0.5em] &= (x - x_0)\vec{i} + (y - y_0)\vec{j} \end{align} \\]

Let’s now take into account the fact that in an orthonormal basis, any vector can be decomposed along the x and y axes. In other words:

\\[\vec{v} = \text{(movement along x)} \cdot \vec{i} + \text{(movement along y)} \cdot \vec{j} \\]

Therefore, the term \\( t\vec{v} \\) on the right-hand side of the vector equation of the line can be written as:

\\[\vec{v} = k\vec{i} + h\vec{j} \\]

so the expression becomes:

\\[t\vec{v} = t(k\vec{i} + h\vec{j}) \\]

We can now rewrite the vector equation of the line as:

\\[(x - x_0)\vec{i} + (y - y_0)\vec{j} = tk\vec{i} + th\vec{j} \\]

By equating the components along the \\( \vec{i} \\) and \\( \vec{j} \\) directions, we obtain the following system of **parametric equations** :

\\[\begin{cases} x = x_0 + kt \\\\[0.5em] y = y_0 + ht \end{cases} \\]

## Example

Let’s find the parametric equations of the line passing through the points:

\\[P_0 = (2, 1) \quad \text{and} \quad P = (5, 4) \\]


To define a line parametrically, we need a point and a direction. Since both \\( P_0 \\) and \\( P \\) lie on the line, we can compute the direction vector \\( \vec{v} \\) by subtracting their coordinates:

\\[\vec{v} = P - P_0 = (5 - 2, 4 - 1) = (3, 3) \\]

This vector tells us how to move along the line starting from \\( P_0 \\): for every 3 units in the x-direction, we move 3 units in the y-direction. We now use the parametric form of the line, which expresses each coordinate as a function of a parameter \\( t \\):

\\[\begin{cases} x = x_0 + kt \\\\[0.5em] y = y_0 + ht \end{cases} \quad \text{with } t \in \mathbb{R} \\]

In our case we have:

  * Coordinates of \\( P_0 \\): \\( x_0 = 2 \\), \\( y_0 = 1 \\)
  * components of the direction vector \\( \vec{v} \\): \\( k = 3 \\), \\( h = 3 \\)


Substituting into the equations, we obtain:

\\( \begin{cases} x = 2 + 3t \\\\[0.5em] y = 1 + 3t \end{cases} \quad \text{with } t \in \mathbb{R} \\)

This system describes all the points on the line that passes through \\( (2, 1) \\) and \\( (5, 4) \\). By varying \\( t \\), we can generate every point along the line in both directions.

##### This method works for any two points: by finding the direction vector and using the parametric form, we can always describe the entire line.

## How are the vector and parametric forms of a line different?

The vector form describes the line geometrically: it builds every point \\( P \\) by starting from a fixed point \\( P_0 \\) and moving in the direction of a vector \\( \vec{v} \\). It looks like this:

\\[P = P_0 + t\vec{v} \\]

The parametric form takes that same rule and writes it out in terms of coordinates, one equation for each axis, using the components of the direction vector:

\\[\begin{cases} x = x_0 + kt \\\\[0.5em] y = y_0 + ht \end{cases} \quad t \in \mathbb{R} \\]

Both forms describe the same set of points. The vector form is compact and geometric; the parametric form is explicit and ready to calculate with. Use whichever makes more sense for the problem you’re solving, they’re mathematically identical.

Lines, Planes and Conic Sections

Conic sections are curves formed by intersecting a plane with a cone.

1.7k

[Lines](https://algebrica.org/lines/)

1.9k

[Polar Coordinates](https://algebrica.org/polar-coordinates/)

1.8k

[Parabola](https://algebrica.org/parabola/)

4.3k

[Circumference](https://algebrica.org/circumference/)

3 comments

1.8k

[Ellipse](https://algebrica.org/ellipse/)

1.8k

[Hyperbola](https://algebrica.org/hyperbola/)
