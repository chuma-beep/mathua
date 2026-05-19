> Content sourced from [Algebrica](https://algebrica.org/circumference/) — CC BY-NC 4.0

## Introduction to conic sections

When introducing the [parabola](<../parabola>), we observed that the intersection of a plane with a cone produces, once projected onto the plane, one of four possible curves: a circumference, a parabola, an [ellipse](<../ellipse>), or a [hyperbola](<../hyperbola>). These four curves are collectively known as conics, and they share a common algebraic nature that allows us to treat them within a single unified framework. From a more formal standpoint, a conic is a second-degree algebraic curve in the plane, defined as the set of points \\((x, y) \in \mathbb{R}^2\\) whose coordinates satisfy a general quadratic equation in the variables \\(x\\) and \\(y\\). This [equation](<../equations/>) takes the following form:

\\[f(x, y) = a_{11}x^2 + 2a_{12}xy + a_{22}y^2 + 2a_{13}x + 2a_{23}y + a_{33} = 0 \\]

  * The coefficients \\(a_{ij}\\) are real numbers.
  * The factor \\(2\\) appearing in front of the mixed and linear terms is a conventional choice that simplifies the matrix representation of the conic.


In order for the equation to describe a quadratic curve, at least one among \\(a_{11}\\), \\(a_{12}\\), and \\(a_{22}\\) must be different from zero otherwise the expression reduces to a [linear equation](<../linear-equations/>) and the locus degenerates into a line.


The specific type of conic obtained depends on the relative values of these coefficients, and in particular on the sign of the [discriminant](<../quadratic-formula/>) associated with the quadratic part. The circumference corresponds to the simplest and most symmetric case, in which the coefficients of \\(x^2\\) and \\(y^2\\) are equal and the mixed term \\(a_{12}xy\\) vanishes.

![](https://algebrica.org/wp-content/uploads/resources/images/conic-circle-1.png)

From a geometric standpoint, this algebraic symmetry reflects the way the plane intersects the cone. When the cutting plane is perpendicular to the axis of the cone, the resulting section is a circumference, whose points are all equidistant from the axis itself. Any deviation from this perpendicularity breaks the symmetry between the two quadratic coefficients and produces one of the other conic sections.

## What is a circumference

Given a point \\(C\\) in the plane, called the center, the circumference is defined as the set of all points in the plane that are equidistant from \\(C\\). Consider a circumference \\(\Gamma\\) with center \\(C\\) and radius \\(r\\). For every point \\(P\\) belonging to \\(\Gamma\\) the following relation holds:

\\[d(P, C) = r \qquad r \in \mathbb{R} \quad \ r > 0 \\]

![](https://algebrica.org/wp-content/uploads/resources/images/circumference-5.png)

Let \\(C = (x_0, y_0)\\) be the coordinates of the center and \\(P = (x, y)\\) a generic point of the circumference. Using the formula for the distance between two points in the plane, the condition \\(d(P,C) = r\\) becomes:

\\[\sqrt{(x - x_0)^2 + (y - y_0)^2} = r \\]

Since both sides are non-negative, squaring does not introduce spurious solutions and we obtain the standard form of the equation of a circumference:

\\[(x - x_0)^2 + (y - y_0)^2 = r^2 \\]

Expanding the squares of the two binomials, we obtain:

\\[x^2 - 2x\,x_0 + x_0^2 + y^2 - 2y\,y_0 + y_0^2 = r^2 \\]

We now move all terms to the left-hand side and introduce the following substitutions:

\\[\begin{aligned} a &= -2x_0 \\\\[6pt] b &= -2y_0 \\\\[6pt] c &= x_0^2 + y_0^2 - r^2 \end{aligned} \\]

With this choice of coefficients, the equation takes the compact form known as the general form of the equation of a circumference:

\\[x^2 + y^2 + ax + by + c = 0 \\]

## Center and radius from the general form

The general form \\(x^2 + y^2 + ax + by + c = 0\\) is algebraically compact, but it does not expose the geometric quantities of the circumference as transparently as the standard form does. Given the coefficients \\(a\\), \\(b\\), and \\(c\\), it is natural to ask whether the center and the radius can be recovered directly from them, and under which conditions the equation actually describes a circumference. Comparing the general form with the expressions of the coefficients obtained earlier:

\\[\begin{aligned} a &= -2x_0 \\\\[6pt] b &= -2y_0 \\\\[6pt] c &= x_0^2 + y_0^2 - r^2 \end{aligned} \\]

The first two relations give immediately the coordinates of the center:

\\[x_0 = -\frac{a}{2} \qquad y_0 = -\frac{b}{2} \\]

Substituting these values into the third relation and solving for \\(r^2\\), we find:

\\[r^2 = \frac{a^2}{4} + \frac{b^2}{4} - c \\]

Since \\(r\\) is a length, it must be a non-negative real number, and therefore the right-hand side of the previous equation cannot be negative. To make the discussion more compact, let us denote the quantity on the right-hand side by \\(\rho\\):

\\[\rho = \frac{a^2}{4} + \frac{b^2}{4} - c \\]

The nature of the locus described by the equation \\(x^2 + y^2 + ax + by + c = 0\\) depends entirely on the sign of \\(\rho\\), and three distinct situations arise.

  * When \\(\rho > 0\\), the equation represents a proper circumference. Its center has coordinates \\(\left(-a/2,\, -b/2\right)\\) and its radius is \\(r = \sqrt{\rho}\\). This is the only case in which the equation describes an actual curve in the plane.

  * When \\(\rho = 0\\), the radius collapses to zero and the equation is satisfied only by the point \\(\left(-a/2,\, -b/2\right)\\). The locus degenerates into a single point, which is sometimes referred to as a degenerate circumference of zero radius.

  * When \\(\rho < 0\\), no real pair \\((x, y)\\) can satisfy the equation, because the sum of two squares can never equal a negative quantity. The locus is therefore empty: the equation still has a second-degree algebraic form, but it does not represent any real curve in the plane.


The condition \\(\rho > 0\\) is thus the condition of existence of a real circumference expressed in general form. Whenever it is satisfied, the equation \\(x^2 + y^2 + ax + by + c = 0\\) can be converted back into the standard form by completing the square, and the center and radius can be read off directly from the formulas above.

## Circle and circumference

The terms circle and circumference are often used interchangeably but they refer to two distinct objects. The following points summarize the terminology and the formulas associated with each of them.

  * The circumference is the curve formed by all points of the plane equidistant from a fixed point called the center. It is a one-dimensional object, a closed curve with no interior.
  * The circle is the region of the plane bounded by a circumference, including all the points in its interior. It is a two-dimensional object, a surface.
  * The length of a circumference of radius \\(r\\) is given by the formula \\(L = 2\pi r\\). The constant \\(\pi\\) expresses the fixed ratio between the length of any circumference and its diameter, a ratio that does not depend on the size of the circumference itself.
  * The area of a circle of radius \\(r\\) is given by the formula \\(A = \pi r^2\\). It would be incorrect to refer to the area of a circumference, since the circumference is a curve and, as such, has no area.


## Radius, chord, diameter, and circle

Once the circumference has been defined, it is useful to introduce the basic geometric elements associated with it: the radius, the chord, the diameter, and the circle itself. These elements appear in the study of conics and in trigonometric applications and each of them identifies a specific relationship between points, segments, and regions of the plane.

![](https://algebrica.org/wp-content/uploads/resources/images/circumference-2.png)

  * A radius is any [line](<../lines>) segment that connects the center of the circumference to a point on the circumference itself.
  * A chord is any segment whose endpoints both lie on the circumference.
  * A diameter is any chord that passes through the center of the circumference; it is therefore the longest possible chord, and its length is exactly twice the radius.
  * A circle is the set containing all the points on the circumference together with all the points in its interior.


Among these elements, the chord plays a particularly important role, because its length can be expressed directly in terms of the radius and of the angle it subtends at the center. If \\(\theta\\) denotes the central angle subtended by the chord, measured in radians, the length \\(c\\) of the chord is given by the following formula, which involves the [sine function](<../sine-function/>):

\\[c = 2r \cdot \sin\\!\left(\frac{\theta}{2}\right) \\]

![](https://algebrica.org/wp-content/uploads/resources/images/circumference-10.png)

The formula has a geometric interpretation: the chord, together with the two radii joining its endpoints to the center, forms an isosceles triangle whose apex angle is \\(\theta\\), and the expression \\(2r\sin(\theta/2)\\) is precisely the length of the base of that triangle.

## Arcs and circular sectors

An arc is a portion of a circumference bounded by two of its points. The endpoints of a chord divide the circumference into two arcs, and we say that the chord subtends the two arcs, or equivalently that each arc is subtended by the chord.

![](https://algebrica.org/wp-content/uploads/resources/images/circumference-3-1.png)

If \\(\theta\\) denotes the central angle corresponding to the arc, measured in radians, and \\(r\\) the radius of the circumference, the length of the arc is given by the formula:

\\[\text{Arc length} = r \cdot \theta \\]


A circular sector is the portion of a circle enclosed between an arc and the two radii connecting the center to the endpoints of the arc. Using the same notation as above, the area \\(A\\) of a circular sector is given by the formula:

\\[A = \frac{1}{2} r^2 \theta \\]

![](https://algebrica.org/wp-content/uploads/resources/images/circumference-6.png)

Arcs represent measurable portions of a circumference, and are directly linked to central angles. Sectors, defined by an arc and two radii, allow us to compute areas and relate angular measures to linear distances. These concepts are essential for deriving arc length formulas and solving problems in trigonometry.

## Example 1

Consider the circumference centered at the origin with radius \\(r = 1\\). Substituting \\(x_0 = 0\\), \\(y_0 = 0\\), and \\(r = 1\\) into the standard form \\((x - x_0)^2 + (y - y_0)^2 = r^2\\), we obtain:

\\[x^2 + y^2 = 1 \\]

To express the same circumference in general form, we compare this equation with:

\\[x^2 + y^2 + ax + by + c = 0 \\]

By direct inspection, the coefficients are:

\\[a = 0 \qquad b = 0 \qquad c = -1 \\]

The equation of the unit circumference in general form is therefore:

\\[x^2 + y^2 - 1 = 0 \\]

This is the simplest non-trivial instance of the general form, and it coincides with the equation of the [unit circle](<../unit-circle>), a reference configuration widely used in trigonometry to define the sine and cosine of an angle.

## Circumferences and lines

The position of a line with respect to a circumference depends on the distance between the line and the center. A [line](<../lines/>) can be classified as secant, tangent, or external to the circumference, and denoting by \\(D\\) the distance from the center to the line and by \\(r\\) the radius, the three cases are mutually exclusive and cover all possibilities.

![](https://algebrica.org/wp-content/uploads/resources/images/circumference-4.png)

  * A line is secant when it intersects the circumference at two distinct points, and this happens when \\(D < r\\).
  * A line is tangent when it touches the circumference at exactly one point, and this happens when \\(D = r\\).
  * A line is external when it does not meet the circumference at all, and this happens when \\(D > r\\).


This geometric classification admits a purely algebraic counterpart. The intersection points between a line and a circumference are the solutions of the system formed by their equations:

\\[\begin{cases} x^2 + y^2 + ax + by + c = 0 \\\\[6pt] y = mx + q \end{cases} \\]

Substituting the expression for \\(y\\) from the second equation into the first leads to a quadratic equation in the unknown \\(x\\), whose number of real solutions coincides with the number of intersection points. Denoting by \\(\Delta\\) the discriminant of this quadratic equation, the three cases discussed above correspond respectively to \\(\Delta > 0\\), \\(\Delta = 0\\), and \\(\Delta < 0\\).

  * When \\(\Delta > 0\\), the quadratic equation admits two real and distinct solutions, and the line is secant to the circumference.
  * When \\(\Delta = 0\\), the quadratic equation admits a single real solution, counted with multiplicity two, and the line is tangent to the circumference.
  * When \\(\Delta < 0\\), the quadratic equation has no real solutions, and the line is external to the circumference.


## Example 2

Consider the circumference of equation \\(x^2 + y^2 - 2x - 2y - 3 = 0\\) and the line of equation \\(y = -x + 3\\). We want to determine whether the line is secant, tangent, or external to the circumference, and in the first two cases find the coordinates of the intersection points.

The intersection points are the solutions of the system formed by the two equations:

\\[\begin{cases} x^2 + y^2 - 2x - 2y - 3 = 0 \\\\[6pt] y = -x + 3 \end{cases} \\]

To simplify the calculations and to reveal the geometric content of the first equation, we rewrite it using the method of [completing the square](<completing-the-square/>), which transforms each quadratic expression in \\(x\\) and \\(y\\) into a perfect square plus a constant:

\\[(x - 1)^2 - 1 + (y - 1)^2 - 1 - 3 = 0 \\]

> The method of completing the square rewrites a quadratic expression in a more structured and manageable form, making it easier to analyze and solve. The goal is to transform a quadratic [polynomial](<../polynomials>) or equation into an equivalent form in which a perfect square [trinomial](<../trinomials>) appears isolated on one side of the equation.


Moving the constants to the right-hand side, the system becomes:

\\[\begin{cases} (x - 1)^2 + (y - 1)^2 = 5 \\\\[6pt] y = -x + 3 \end{cases} \\]

The first equation now reveals that the given curve is a circumference with center \\((1, 1)\\) and radius \\(r = \sqrt{5}\\). Substituting the expression for \\(y\\) from the second equation into the first, we obtain:

\\[(x - 1)^2 + (-x + 2)^2 = 5 \\]

Expanding both squares:

\\[x^2 - 2x + 1 + x^2 - 4x + 4 = 5 \\]

Combining like terms and moving all terms to the left-hand side, the equation reduces to:

\\[2x^2 - 6x = 0 \\]

Factoring out \\(2x\\), we obtain the equivalent [quadratic equation](<../quadratic-equation>):

\\[2x(x - 3) = 0 \\]

This equation admits two distinct real solutions:

\\[x = 0 \qquad \text{or} \qquad x = 3 \\]

The corresponding values of \\(y\\) are obtained by substituting into the equation of the line \\(y = -x + 3\\). For \\(x = 0\\) we get \\(y = 3\\), and for \\(x = 3\\) we get \\(y = 0\\).

Therefore, the points of intersection are:

\\[P_1 = (0,\ 3) \quad \text{and} \quad P_2 = (3,\ 0) \\]

Since there are two distinct solutions, the line is a secant to the circle.

Concept

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Easy

1

Requires

1

Enables

The following concepts, [Quadratic Equations](https://algebrica.org/quadratic-equations/), are required as prerequisites for this entry.

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

1.8k

[Ellipse](https://algebrica.org/ellipse/)

1.8k

[Hyperbola](https://algebrica.org/hyperbola/)
