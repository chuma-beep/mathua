> Content sourced from [Algebrica](https://algebrica.org/parabola/) — CC BY-NC 4.0

## Introduction to Conic Sections

When a plane intersects a cone, the shape formed at the intersection, when projected onto the plane, can be a [circumference](<../circumference>), a parabola, an [ellipse](<../ellipse>), or a [hyperbola](<../hyperbola>). These curves are collectively known as conic sections, or simply conics. In more formal terms, a conic is a second-degree plane algebraic curve, defined as the set of points \\( (x, y) \in \mathbb{R}^2 \\) satisfying a quadratic equation in the variables \\( x \\) and \\( y \\):

\\[f(x, y) = a_{11}x^2 + 2a_{12}xy + a_{22}y^2 + 2a_{13}x + 2a_{23}y + a_{33} = 0 \\]

where the coefficients \\( a_{ij} \in \mathbb{R} \\), with \\( a_{11} \\) and \\( a_{22} \\) both nonzero to ensure the equation is genuinely quadratic.

## Each coefficient carries specific geometric information about the conic.

  * \\( a_{11} \\), \\( a_{22} \\) determine the curvature along the x-axe and y-axe, respectively. They influence the shape (ellipse, hyperbola, parabola) and orientation of the conic.
  * \\( a_{12} \\) controls the rotation of the conic. If \\( a_{12} = 0 \\), the conic is aligned with the coordinate axes.
  * \\( a_{13} \\), \\( a_{23} \\) affect the position of the conic in the plane. They act as translation terms along the x-axis and y-axis.
  * \\( a_{33} 1\\)determines the overall position of the curve with respect to the origin; it can be viewed as a constant term shifting the graph up/down or left/right depending on context.


## Degenerate conic

If the polynomial \\( f(x, y)\ \\) can be factored as a product of two linear polynomials:

\\[f(x, y) = (ax + by + c)(a’x + b’y + c’) = 0 \\]

where \\(a, b, c, a’, b’, c’ \in \mathbb{C}\\), then the conic is said to be degenerate.


##### A degenerate conic does not represent a proper curved figure such as a parabola, ellipse, or hyperbola. Instead, it corresponds to simpler geometric objects like a pair of lines, a single line, or in some cases the empty set.

## The Parabola

The **parabola** is a plane curve defined as the set of all points that are equidistant from a fixed point \\(F \\) called the focus, and a fixed line \\( d \\) called the directrix.

![Parabola](https://algebrica.org/wp-content/uploads/resources/images/parabola-1.png)

The length of the segment \\( \overline{FP} \\) is the same as that of the segment \\( \overline{PD} \\). The parabola is one of the so-called conic sections which includes the circle, ellipse, and hyperbola. These curves can be obtained by intersecting a conical surface with a plane. The specific curve formed depends on the angle at which the plane intersects the cone.


The line passing through the focus and perpendicular to the directrix is called the axis of the parabola. The point \\( V \\) where the parabola intersects this axis is known as the vertex.

![](https://algebrica.org/wp-content/uploads/resources/images/parabola-2.png)

The equation of a parabola with its vertex at the origin and its axis coinciding with the y-axis of the Cartesian plane is given by:

\\[y = ax^2, \quad a \neq 0\\]

Every parabola that satisfies this equation is symmetric with respect to the y-axis. In the case where the coefficient \\( a = 0 \\), the parabola is said to be degenerate, and the equation becomes \\( y = 0 \\). The coordinates of the focus are:

\\[F = \left( 0, \frac{1}{4a} \right)\\]

The equation of the directrix is:

\\[y = - \frac{1}{4a} \\]


When the coefficient \\( a > 0 \\), the parabola opens upward. Therefore, \\( y \geq 0 \\) for every value of \\( x \\). The focus is also located on the positive half of the y-axis.

![](https://algebrica.org/wp-content/uploads/resources/images/parabola-3.png)


The coefficient \\( a \\) determines another characteristic of the parabola, namely its width or opening. If \\( a > 0 \\), as the value of \\( a \\) increases, the opening becomes narrower. Similarly, if \\( a < 0 \\) as the [absolute value](<../absolute-value>) of \\( a \\) increases, the opening also becomes narrower.

![](https://algebrica.org/wp-content/uploads/resources/images/parabola-4.png)

## The parabola in standard quadratic form

Let’s now consider the general case of the equation of a parabola with its axis parallel to the y-axis. The equation is given by:

\\[y = ax^2 + bx + c \quad \text{with} \quad a \neq 0\\]

The equation that describes a parabola is a [second-degree equation](<../quadratic-equation>).


The equation of the axis is given by:

\\[x = -\frac{b}{2a}\\]


The coordinates of the vertex are given by:

\\[V \left( -\frac{b}{2a}, \, -\frac{\Delta}{4a} \right)\\]

where \\( \Delta = b^2 - 4ac \\) is the [discriminant](<../quadratic-formula>) of the quadratic equation.


The coordinates of the focus are:

\\[F \left( -\frac{b}{2a}, \, \frac{1 - \Delta}{4a} \right)\\]


The equation of the directrix is:

\\[y = -\frac{1 + \Delta}{4a}\\]

Graphically, a generic \\(y = ax^2 + bx + c\\) parabola with its axis parallel to the y-axis looks like the following:

![](https://algebrica.org/wp-content/uploads/resources/images/parabola-5-1.png)

  * When \\( b = 0 \\) and \\( c \neq 0 \\) the equation becomes \\(y = ax^2 + c\\). The parabola has its vertex at \\( V(0, c) \\), and its axis of symmetry is the y-axis.

  * When Case: \\( c = 0 \\) and \\( b \neq 0 \\) the equation becomes \\(y = ax^2 + bx\\). The parabola has its vertex at: \\[V \left( -\frac{b}{2a}, -\frac{b^2}{4a} \right) \\] The parabola always passes through the origin ( 0, 0 ).


Let’s now see how to determine the intersection points, if they exist, between a parabola and a generic line with the equation \\( y = mx + q \\). To do this, we need to solve the system between the equation of the line and that of the parabola. We obtain:

\\[\begin{cases} y = ax^2 + bx + c \\\\[0.5 em] y = mx + q \end{cases} \\]

Performing the calculations, we obtain:

\\[\begin{align} &ax^2 + bx + c = mx + q\\\\[0.5 em] &ax^2 + bx + c - mx - q = 0\\\\[0.5 em] &ax^2 + x(b- m) +c - q = 0 \end{align} \\]

The solutions of equation \\((5)\\) represent the \\( x \\)-coordinates of the intersection points between the parabola and the line. Since we are dealing with a quadratic equation, there can be at most two distinct solutions. More precisely, the number of solutions depends on the value of the discriminant \\( \Delta \\), as follows:

  * If \\( \Delta > 0 \\) the solutions are real and distinct, and the line intersects the parabola at two points. In this case, the line is called a secant.
  * If \\( \Delta = 0 \\) the solutions are real and coincident, and the line is tangent to the parabola at a single point.
  * If \\( \Delta < 0 \\) there are no solutions, and the line is external to the parabola.


Below is the representation of a secant line intersecting the parabola at two points.

![](https://algebrica.org/wp-content/uploads/resources/images/parabola-6-2.png)


Now, let’s consider a point \\( P \\) on the plane and determine the lines passing through this point that are tangent to the parabola. There are three possible cases:

  * Both lines are tangent to the parabola: in this case, the point \\(P\\) is external to the parabola.
  * Only one line is tangent to the parabola: in this case, the point \\(P\\) is on the parabola.
  * There is no line tangent to the parabola: in this case, the point \\(P\\) is internal to the parabola.


![](https://algebrica.org/wp-content/uploads/resources/images/parabola-7.png)

To find the equations of the lines passing through a given point \\( P(x_0, y_0) \\) and tangent to the parabola described by the equation \\( y = ax^2 + bx + c \\), we need to formulate and solve the system consisting of the parabola’s equation and the equation of the pencil of lines passing through \\( P \\). We have:

\\[\begin{cases} y - y_0 = m(x - x_0) \\\\[0.5em] y = ax^2 + bx + c \end{cases} \\]

To impose the condition of tangency, we need to set the discriminant \\( \Delta \\) of the equation obtained from the system equal to zero. Then, we solve with respect to \\( m \\) and substitute the resulting values into the equation of the pencil of lines.

## Example

Determine the equations of any lines passing through \\( P(3, -6) \\) and tangent to the parabola described by the equation \\( y = x^2 - 4 \\).


The general equation of the pencil of lines is:

\\[y - y_0 = m(x - x_0)\\]

By substituting the values \\( x_0 = 3 \\) and \\( y_0 = -6 \\) of the point \\( P(3, -6) \\), we obtain:

\\[y + 6 = m(x - 3)\\]


The system becomes:

\\[\begin{cases} y = x^2 - 4 \\\\[0.5em] y + 6 = m(x - 3) \end{cases} \\]

We obtain: \\[x^2 - 4 = m(x-3) - 6 \rightarrow x^2 -mx +3m -2 = 0 \\]


Let’s now determine the value of \\( \Delta \\) considering that \\( a = 1 \\), \\( b = -m \\), and \\( c = 3m - 2 \\). We have:

\\[\Delta = m^2 -12m - 8\\]


We apply the condition of tangency by setting \\( \Delta = 0 \\) and solving the second degree equation.

\\[m^2 -12m - 8 = 0 \rightarrow m_1, m_2 = \frac{12 \pm \sqrt{144+32}}{2}\\]

By simplifying the expression, we obtain:

\\[m_1 = 6 - 2\sqrt{11}, \quad m_2 = 6 + 2\sqrt{11} \\]


By substituting these values into the equation of the pencil of lines, we obtain:

\\[y = \left(6 - 2\sqrt{11}\right)\left(x - 3\right) -6 \\] \\[y = \left(6 + 2\sqrt{11}\right)\left(x - 3\right) - 6\\]

Therefore, the equations of the two tangent lines are given by:

\\[\begin{align*} y &= \left(6 - 2\sqrt{11}\right)x + 6\sqrt{11} - 24 \\\\[0.5em] y &= \left(6 + 2\sqrt{11}\right)x - 6\sqrt{11} - 24 \end{align*} \\]

Lines, Planes and Conic Sections

Conic sections are curves formed by intersecting a plane with a cone.

1.7k

[Lines](https://algebrica.org/lines/)

2k

[Vector and Parametric Equations of a Line](https://algebrica.org/vector-and-parametric-equations-of-a-line/)

1.9k

[Polar Coordinates](https://algebrica.org/polar-coordinates/)

4.3k

[Circumference](https://algebrica.org/circumference/)

3 comments

1.8k

[Ellipse](https://algebrica.org/ellipse/)

1.8k

[Hyperbola](https://algebrica.org/hyperbola/)
