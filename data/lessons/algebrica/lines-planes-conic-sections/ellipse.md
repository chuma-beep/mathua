> Content sourced from [Algebrica](https://algebrica.org/ellipse/) — CC BY-NC 4.0

## Introduction to conic sections

When introducing the [parabola](<../parabola>), we saw that when a plane intersects a cone, the resulting shape, when projected onto the plane, can be a [circumference](<../circumference>), a parabola, an ellipse, or a [hyperbola](<../hyperbola>). These curves are collectively referred to as conics. More formally, a conic is a second-degree algebraic curve in the plane. It is defined as the set of points \\( (x, y) \in \mathbb{R}^2 \\) that satisfy a general quadratic equation in the variables \\( x \\) and \\( y \\):

\\[f(x, y) = a_{11}x^2 + 2a_{12}xy + a_{22}y^2 + 2a_{13}x + 2a_{23}y + a_{33} = 0 \\]

Here, the coefficients \\( a_{ij} \in \mathbb{R} \\), and to ensure the curve is truly quadratic, we require that both \\( a_{11} \\) and \\( a_{22} \\) are nonzero.

## What is an ellipse

Given two fixed points in the plane, \\( F_1 \\) and \\( F_2 \\), an ellipse is defined as the set of all points \\( P \\) in the plane such that the sum of the distances from \\( P \\) to each focus is constant.

\\[\overline{PF_1} + \overline{PF_2} = \text{constant} \\]

![](https://algebrica.org/wp-content/uploads/resources/images/ellipse-1.png)

\\( F_1 \\) and \\( F_2 \\) are the foci of the ellipse. Assuming that the focus \\( F_1 \\) has coordinates \\( (-c, 0) \\) and the focus \\( F_2 \\) has coordinates \\( (c, 0) \\), the distance between \\( F_1 \\) and \\( F_2 \\) is called the focal distance and is equal to \\( 2c \\). The midpoint of the segment \\( \overline{F_1F_2} \\) is the center of the ellipse.


We define the major axis and minor axis of the ellipse, respectively, as the longest and shortest diameters passing through its center. The major axis lies along the direction of maximum extension of the ellipse and passes through both foci. Its total length is \\( 2a \\), where \\( a \\) is the semi-major axis. The minor axis is perpendicular to the major axis and also passes through the center of the ellipse. Its total length is \\( 2b \\), where \\( b \\) is the semi-minor axis.


By choosing a point \\( P = (a, 0) \\), located at the right endpoint of the major axis, we consider the case where the ellipse intersects the \\( x \\)-axis at its farthest horizontal extent. In this configuration, the distance from the left focus \\( F_1 = (-c, 0) \\) to the point \\( P \\) is \\(\overline{F_1P} = a + c\\). Similarly, we have \\(\overline{F_2P} = a + c\\).

![](https://algebrica.org/wp-content/uploads/resources/images/ellipse-2.png)

From this, we conclude that the constant sum of the distances from any point on the ellipse to the two foci is:

\\[\overline{F_1P} + \overline{F_2P} = (a + c) + (a - c) = 2a \\]

In the standard form, an ellipse centered at the origin with horizontal major axis is described by the equation:

\\[\frac{x^2}{a^2} + \frac{y^2}{b^2} = 1 \\]

where \\( b^2 = a^2 - c^2 \\), with \\( b > 0 \\) and \\( a > b \\), from which it follows that:

\\[c = \sqrt{a^2 - b^2} \\]

## Why is the sum of distances to the foci always constant in an ellipse?

Because that’s what defines it. An ellipse is the set of all points for which the sum of the distances to the two foci is exactly \\( 2a \\). Any point not satisfying this condition lies outside or inside the curve.

## Vertices

An ellipse intersects the coordinate axes at four key points, its vertices. The two on the major axis represent the farthest horizontal reach.

![](https://algebrica.org/wp-content/uploads/resources/images/ellipse-3.png)

The two on the minor axis define the vertical extent. All are symmetric with respect to the center and capture the ellipse’s full geometric footprint.

## Eccentricity

The ratio between the focal distance and the length of the major axis of an ellipse is called its eccentricity. It is denoted by \\( e \\) and satisfies the condition:

\\[0 \leq e < 1 \\]

![](https://algebrica.org/wp-content/uploads/resources/images/ellipse-4.png)

The closer \\( e \\) is to 0, the more circular the ellipse appears. As \\( e \\) approaches 1, the ellipse becomes increasingly elongated. The value of the eccentricity \\( e \\) is given by:

\\[e = \frac{c}{a} = \frac{\sqrt{a^2 - b^2}}{a} \\]

##### Eccentricity describes how “stretched” an ellipse is. When \\( e = 0 \\), the ellipse is indistinguishable from a circle and its foci coincide at the center. As \\( e \\) increases, the foci move apart and the shape elongates along the major axis. What matters is not the size, but the ratio: eccentricity is a pure measure of shape.

## Example

Let us determine the equation of the ellipse with foci \\( F_1 = (1, 0) \\) and \\( F_2 = (-1, 0) \\), such that the sum of the distances from any point on the ellipse to the two foci is equal to 6.


A point \\( P(x, y) \\) belongs to the ellipse if it satisfies the condition:

\\[\sqrt{(x - 1)^2 + y^2} + \sqrt{(x + 1)^2 + y^2} = 6 \\]

Each square root represents the Euclidean distance between \\( P(x, y) \\) and one of the two foci. This equation expresses the geometric definition of the ellipse: the set of all points in the plane such that the sum of their distances to the two foci is constant. Hence, determining the equation of the ellipse simply requires performing the necessary calculations. We have:

\\[\begin{align} &(x - 1)^2 + y^2 = 36 + (x + 1)^2 + y^2 - 12\sqrt{(x + 1)^2 + y^2}\\\\[0.5em] &x^2 - 2x + 1 + y^2 = 36 + x^2 +2x +1 +y^2 - 12\sqrt{(x + 1)^2 + y^2}\\\\[0.5em] &-4x -36 = - 12\sqrt{(x + 1)^2 + y^2}\\\\[0.5em] &x + 9 = 3\sqrt{(x + 1)^2 + y^2}\\\\[0.5em] \end{align} \\]


By squaring both sides of the equation, we obtain:

\\[\begin{align} &(x + 9)^2 = (3\sqrt{(x + 1)^2 + y^2})^2\\\\[0.5em] &x^2 + 18x +81 = 9x^2 +18x +9 + 9y^2\\\\[0.5em] &8x^2 + 9y^2 = 72\\\\[0.5em] \end{align} \\]


Dividing both sides of the equation by 72, we obtain:

\\[\frac{8x^2}{72} + \frac{9y^2}{72} = 1 \\]

Therefore, the equation of the ellipse passing through the points mentioned above, and such that the sum of the distances from any point ( P ) on the curve to the two foci is equal to 6, is:

\\[\frac{x^2}{9} + \frac{y^2}{8} = 1 \\]

## Glossary

  * Ellipse: the set of all points in a plane where the sum of the distances to two fixed points (foci) is constant.

  * Foci: the two fixed points in the plane used to define an ellipse.

  * Conic section: a curve formed by the intersection of a plane and a cone; includes circumferences, parabolas, ellipses, and hyperbolas.

  * Major axis: the longest diameter of an ellipse, passing through the foci and the center. Its length is \\( 2a \\).

  * Semi-najor axis: half the length of the major axis, denoted by \\( a \\).

  * Minor axis: the shortest diameter of an ellipse, perpendicular to the major axis and passing through the center. Its length is \\( 2b \\).

  * Semi-minor axis: half the length of the minor axis, denoted by \\( b \\).

  * Center: the midpoint of the segment connecting the two foci.

  * Focal distance: the distance between the two foci, equal to \\( 2c \\).

  * Vertices: the four points where an ellipse intersects the coordinate axes.

  * Eccentricity \\(e\\): the ratio of the focal distance \\( 2c \\) to the length of the major axis \\( 2a \\), \\( e = c/a \\). It is a measure of how stretched the ellipse is, with \\( 0 \leq e < 1 \\).


Lines, Planes and Conic Sections

Conic sections are curves formed by intersecting a plane with a cone.

1.7k

[Lines](https://algebrica.org/lines/)

2k

[Vector and Parametric Equations of a Line](https://algebrica.org/vector-and-parametric-equations-of-a-line/)

1.9k

[Polar Coordinates](https://algebrica.org/polar-coordinates/)

1.8k

[Parabola](https://algebrica.org/parabola/)

4.3k

[Circumference](https://algebrica.org/circumference/)

3 comments

1.8k

[Hyperbola](https://algebrica.org/hyperbola/)
