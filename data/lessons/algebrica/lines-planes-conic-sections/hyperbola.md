> Content sourced from [Algebrica](https://algebrica.org/hyperbola/) — CC BY-NC 4.0

## Introduction to conic sections

When introducing the [parabola](<../parabola>), we saw that when a plane intersects a cone, the resulting shape, when projected onto the plane, can be a [circumference](<../circumference>), a [parabola](<../parabola>), an [ellipse](<../ellipse>), or a hyperbola. These curves are collectively referred to as conics. More formally, a conic is a second-degree algebraic curve in the plane. It is defined as the set of points \\( (x, y) \in \mathbb{R}^2 \\) that satisfy a general quadratic equation in the variables \\( x \\) and \\( y \\):

\\[f(x, y) = a_{11}x^2 + 2a_{12}xy + a_{22}y^2 + 2a_{13}x + 2a_{23}y + a_{33} = 0 \\]

Here, the coefficients \\( a_{ij} \in \mathbb{R} \\), and to ensure the curve is truly quadratic, we require that both \\( a_{11} \\) and \\( a_{22} \\) are nonzero.

## What is a hyperbola

Given two fixed points in the plane, \\( F_1 \\) and \\( F_2 \\), a hyperbola is defined as the set of all points \\( P \\) such that the [absolute value](<../absolute-value/>) of the difference between the distances from \\( P \\) to each focus is constant. In other words:

\\[\left| PF_1 - PF_2 \right| = \text{constant} \\]

![Standard chart of a hyperbola.](/diagrams/algebrica/hyperbola-1.png)

\\(F_1\\) and \\(F_2\\) are the foci of the hyperbola. The midpoint of the segment \\( \overline{F_1F_2} \\) is called the center (which, in the figure, coincides with the origin of the Cartesian axes). The midpoint of the segment \\( \overline{F_1F_2} \\) is called the center.


Only the x-axis (also called the focal or transverse axis) intersects the hyperbola at two real points: \\( A(a, 0) \\) and \\( A’(-a, 0) \\), known as the vertices. The y-axis does not intersect the hyperbola and is referred to as the non-transverse axis.

![Asymptotes of a hyperbola.](/diagrams/algebrica/hyperbola-2.png)

The [asymptotes](<../asymptotes/>) of a hyperbola are straight lines that the curve approaches but never intersects. They represent the directions along which the branches of the hyperbola extend infinitely. In the case of a standard hyperbola centered at the origin, the asymptotes are given by the equations:

\\[y = \pm \frac{b}{a}x \\]

As \\(|x|\\) increases, \\(|y|\\) also increases, and the curve gets infinitely closer to the asymptotes.


When \\( P \\) lies on one of the two vertices, for example, \\( (a, 0) \\), the difference between the distances from \\( F_1 \\) and \\( F_2 \\) is exactly \\( 2a \\), and this value remains constant for all points on the hyperbola. Therefore, we have:

\\[\left| \overline{PF_1} - \overline{PF_2} \right| = 2a \\]

In the standard form, a hyperbola centered at the origin with a horizontal transverse axis is described by the equation:

\\[\frac{x^2}{a^2} - \frac{y^2}{b^2} = 1 \\]

where \\(b^2 = c^2 - a^2\\), with \\(b > 0\\) and \\(c > a\\) from which it follows that:

\\[c = \sqrt{a^2 + b^2} \\]

## Why is the difference of distances to the foci always constant in a hyperbola?

Because that’s what defines it. A hyperbola is the set of all points for which the absolute difference of the distances to the two foci is exactly \\( 2a \\). Any point not satisfying this condition lies outside the curve and does not belong to the hyperbola.

## Rectangular hyperbola

If in the canonical equation of a hyperbola we have \\( a = b \\), the hyperbola is called a rectangular hyperbola. This condition makes the asymptotes perpendicular, forming right angles. When the foci lie on the \\( x \\)-axis, the equation of the rectangular hyperbola becomes:

\\[\frac{x^2}{a^2} - \frac{y^2}{a^2} = 1 \quad \rightarrow \quad x^2 - y^2 = a^2 \\]

![Rectangular hyperbola graph.](/diagrams/algebrica/hyperbola-4.png)

##### In the Cartesian plane, the bisectors of the quadrants are the two lines \\( y = x \\) and \\( y = -x \\), which symmetrically divide the space with respect to the axes.

## Eccentricity

In a hyperbola, the eccentricity is defined as the ratio between the focal distance \\( c \\) and the semi-transverse axis \\( a \\). This value characterizes the opening of the hyperbola and is always greater than 1, so \\( e > 1 \\). That is:

\\[e = \frac{c}{a} = \frac{\sqrt {a^2 + b^2}}{a} \\]

![Eccentricity of a hyperbola.
](/diagrams/algebrica/hyperbola-3.png)

##### Eccentricity describes how open a hyperbola is. When \\( e = 1 \\), the branches of the hyperbola are relatively narrow. As \\( e \\) increases, the foci move farther from the center, and the branches open wider. The eccentricity does not depend on the size of the hyperbola, but on the ratio between distances: it is a pure measure of shape.

## A bridge between circular and hyperbolic trigonometry

The equilateral hyperbola plays a central role in hyperbolic trigonometry. Just as the circular [sine and cosine](<../sine-and-cosine/>) are defined using the unit circle, the [hyperbolic sine and cosine](<../hyperbolic-sine-and-cosine/>) arise from the geometry of the hyperbola: \\[x^{2} - y^{2} = 1 \\] Here a hyperbolic sector determines a parameter \\(x\\), and the point \\(P\\)on the hyperbola associated with this sector has coordinates: \\[P_{x} = \cosh(x) = \frac{e^{x} + e^{-x}}{2} \\] \\[P_{y} = \sinh(x) = \frac{e^{x} - e^{-x}}{2} \\] This parallel between the circle and the hyperbola makes it clear that each curve gives rise to its own kind of trigonometric behaviour. The familiar circular functions have their hyperbolic counterparts, and the two settings fit together in a way that highlights the shared geometric idea behind both constructions.
