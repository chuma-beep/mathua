> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Line_%28geometry%29) — CC BY-SA 4.0

# Points, lines, and planes

the graphical conceptLine (graphics)

1d

In geometry, a **straight line**, usually abbreviated **line**, is an infinitely long object with no width, depth, or curvature. It is a special case of a curve and an idealization of such physical objects as a straightedge, a taut string, or a ray of light. Lines are spaces of dimension one, which may be embedded in spaces of dimension two, three, or higher. The word *line* may also refer, in everyday life, to a line segment, which is a part of a line delimited by two points (its *endpoints*).

Euclid's *Elements* defines a straight line as a "breadthless length" that "lies evenly with respect to the points on itself", and introduced several postulates as basic unprovable properties on which the rest of geometry was established. ***Euclidean line*** and *Euclidean geometry* are terms introduced to avoid confusion with generalizations introduced since the end of the 19th century, such as non-Euclidean, projective, and affine geometry.

## Properties
In the Greek deductive geometry of Euclid's *Elements*, a general *line* (now called a *curve*) is defined as a "breadthless length", and a *straight line* (now called a line segment) was defined as a line "which lies evenly with the points on itself". These definitions appeal to readers' physical experience, relying on terms that are not themselves defined, and the definitions are never explicitly referenced in the remainder of the text. In modern geometry, a line is usually either taken as a primitive notion with properties given by axioms, or else defined as a set of points obeying a linear relationship, for instance when real numbers are taken to be primitive and geometry is established analytically in terms of numerical coordinates.

In an axiomatic formulation of Euclidean geometry, such as that of Hilbert (modern mathematicians added to Euclid's original axioms to fill perceived logical gaps), a line is stated to have certain properties that relate it to other lines and points. For example, for any two distinct points, there is a unique line containing them, and any two distinct lines intersect at most at one point. In two dimensions (i.e., the Euclidean plane), two lines that do not intersect are called parallel. In higher dimensions, two lines that do not intersect are parallel if they are contained in a plane, or skew if they are not.

On a Euclidean plane, a line can be represented as a boundary between two regions. Any collection of finitely many lines partitions the plane into convex polygons (possibly unbounded); this partition is known as an arrangement of lines.

### In higher dimensions
In three-dimensional space, a first degree equation in the variables *x*, *y*, and *z* defines a plane, so two such equations, provided the planes they give rise to are not parallel, define a line which is the intersection of the planes. More generally, in *n*-dimensional space *n*−1 first-degree equations in the *n* coordinate variables define a line under suitable conditions.

In more general Euclidean space, **R**\(^{*n*}\) (and analogously in every other affine space), the line *L* passing through two different points *a* and *b* is the subset

\[
L = \left\{ (1 - t) \, a + t b \mid t\in\mathbb{R}\right\}.
\]

The orientation or direction of the **oriented line** (or **directed line**) above is from a reference point *a* (*t* = 0) to a target point *b* (*t* = 1), or in other words, in the direction of the relative vector *b* − *a*. A directed line is also called an ***axis***, especially if it plays a distinctive role (e.g., an axis of rotation). Swapping points *a* and *b* yields the opposite directed line.

#### Collinear points


Three or more points are said to be *collinear* if they lie on the same line. If three points are not collinear, there is exactly one plane that contains them.

In affine coordinates, in *n*-dimensional space the points *X* = (*x*\(_{1}\), *x*\(_{2}\), ..., *x*\(_{*n*}\)), *Y* = (*y*\(_{1}\), *y*\(_{2}\), ..., *y*\(_{*n*}\)), and *Z* = (*z*\(_{1}\), *z*\(_{2}\), ..., *z*\(_{*n*}\)) are collinear if the matrix

\[
\begin{bmatrix}
 1 & x_1 & x_2 & \cdots & x_n \\
 1 & y_1 & y_2 & \cdots & y_n \\
 1 & z_1 & z_2 & \cdots & z_n
\end{bmatrix}
\]

has a rank less than 3. In particular, for three points in the plane (*n* = 2), the above matrix is square and the points are collinear if and only if its determinant is zero.

Equivalently for three points in a plane, the points are collinear if and only if the slope between one pair of points equals the slope between any other pair of points (in which case the slope between the remaining pair of points will equal the other slopes). By extension, *k* points in a plane are collinear if and only if any (*k*–1) pairs of points have the same pairwise slopes.

In Euclidean geometry, the Euclidean distance *d*(*a*,*b*) between two points *a* and *b* may be used to express the collinearity between three points by:
The points *a*, *b* and *c* are collinear if and only if *d*(*x*,*a*) = *d*(*c*,*a*) and *d*(*x*,*b*) = *d*(*c*,*b*) implies *x* = *c*.
However, there are other notions of distance (such as the Manhattan distance) for which this property is not true.

In the geometries where the concept of a line is a primitive notion, as may be the case in some synthetic geometries, other methods of determining collinearity are needed.

### Relationship with other figures

In Euclidean geometry, all lines are congruent, meaning that every line can be obtained by moving a specific line. However, lines may play special roles with respect to other geometric objects and can be classified according to that relationship.

For instance, with respect to a conic (a circle, ellipse, parabola, or hyperbola), lines can be:
* tangent lines, which touch the conic at a single point;
* secant lines, which intersect the conic at two points and pass through its interior;
* exterior lines, which do not meet the conic at any point of the Euclidean plane; or
* a directrix, whose distance from a point helps to establish whether the point is on the conic.
* a coordinate line, a linear coordinate dimension

In the context of determining parallelism in Euclidean geometry, a transversal is a line that intersects two other lines that may or not be parallel to each other.

For more general algebraic curves, lines could also be:
* *i*-secant lines, meeting the curve in *i* points counted without multiplicity, or
* asymptotes, which a curve approaches arbitrarily closely without touching it.
With respect to triangles we have:
* the Euler line,
* the Simson lines, and
* central lines.

For a convex quadrilateral with at most two parallel sides, the Newton line is the line that connects the midpoints of the two diagonals.

For a hexagon with vertices lying on a conic we have the Pascal line and, in the special case where the conic is a pair of lines, we have the Pappus line.

Parallel lines are lines in the same plane that never cross. Intersecting lines share a single point in common. Coincidental lines coincide with each other—every point that is on either one of them is also on the other.

Perpendicular lines are lines that intersect at right angles.

In three-dimensional space, skew lines are lines that are not in the same plane and thus do not intersect each other.

### In axiomatic systems
In synthetic geometry, the concept of a line is often considered as a primitive notion, meaning it is not being defined by using other concepts, but  it is defined by the properties, called  axioms, that it must satisfy.

However, the axiomatic definition of a line does not explain the relevance of the concept and is often too abstract for beginners. So, the definition is often replaced or completed by a *mental image* or *intuitive description* that allows understanding of what a line is. Such descriptions are sometimes referred to as definitions, but are not true definitions since they cannot be used in mathematical proofs. The "definition" of a line in Euclid's Elements falls into this category; and is never used in proofs of theorems.

## Definition


### Linear equation


Lines in a Cartesian plane or, more generally, in affine coordinates, are characterized by linear equations. More precisely, every line \(L\) (including vertical lines) is the set of all points whose coordinates (*x*, *y*) satisfy a linear equation; that is,

\[
L = \{(x,y)\mid ax+by=c\},
\]

where *a*, *b* and *c* are fixed real numbers (called coefficients) such that *a* and *b* are not both zero. Using this form, vertical lines correspond to equations with *b* = 0.

One can further suppose either *c* = 1 or *c* = 0, by dividing everything by \(c\) if it is not zero.

There are many variant ways to write the equation of a line which can all be converted from one to another by algebraic manipulation. The above form is sometimes called the *standard form*. If the constant term is put on the left, the equation becomes

\[
ax + by - c = 0,
\]

and this is sometimes called the *general form* of the equation. However, this terminology is not universally accepted, and many authors do not distinguish these two forms.

These forms are generally named by the type of information (data) about the line that is needed to write down the form. Some of the important data of a line is its slope, x-intercept, known points on the line and y-intercept.

The equation of the line passing through two different points \(P_0( x_0, y_0 )\) and \(P_1(x_1, y_1)\) may be written as

\[
(y - y_0)(x_1 - x_0) = (y_1 - y_0)(x - x_0).
\]

If *x*\(_{0}\) ≠ *x*\(_{1}\), this equation may be rewritten as

\[
y=(x-x_0)\,\frac{y_1-y_0}{x_1-x_0}+y_0
\]

or

\[
y=x\,\frac{y_1-y_0}{x_1-x_0}+\frac{x_1y_0-x_0y_1}{x_1-x_0}\,.
\]
In two dimensions, the equation for non-vertical lines is often given in the *slope–intercept form*:


\[
y = mx + b
\]

where:
* *m* is the slope or gradient of the line.
* *b* is the y-intercept of the line.
* *x* is the independent variable of the function *y* = *f*(*x*).

The slope of the line through points \(A(x_a, y_a)\) and \(B(x_b, y_b)\), when \(x_a \neq x_b\), is given by \(m = (y_b - y_a)/(x_b - x_a)\) and the equation of this line can be written \(y = m (x - x_a) + y_a\).

As a note, lines in three dimensions may also be described as the simultaneous solutions of two linear equations

\[
a_1 x + b_1 y + c_1 z - d_1 = 0
\]


\[
a_2 x + b_2 y + c_2 z - d_2 = 0
\]

such that \((a_1,b_1,c_1)\) and \((a_2,b_2,c_2)\) are not proportional (the relations \(a_1 = t a_2, b_1 = t b_2, c_1 = t c_2\) imply \(t = 0\)). This follows since in three dimensions a single linear equation typically describes a plane and a line is what is common to two distinct intersecting planes.

### Parametric equation

Parametric equations are also used to specify lines, particularly in those in three dimensions or more because in more than two dimensions lines *cannot* be described by a single linear equation.

In three dimensions lines are frequently described by parametric equations:

\[
\begin{align}
x &= x_0 + at \\
y &= y_0 + bt \\
z &= z_0 + ct
\end{align}
\]

where:
* *x*, *y*, and *z* are all functions of the independent variable *t* which ranges over the real numbers.
* (*x*\(_{0}\), *y*\(_{0}\), *z*\(_{0}\)) is any point on the line.
* *a*, *b*, and *c* are related to the slope of the line, such that the direction vector (*a*, *b*, *c*) is parallel to the line.

Parametric equations for lines in higher dimensions are similar in that they are based on the specification of one point on the line and a direction vector.

### Hesse normal form


The *normal form* (also called the *Hesse normal form*, after the German mathematician Ludwig Otto Hesse), is based on the *normal segment* for a given line, which is defined to be the line segment drawn from the origin perpendicular to the line. This segment joins the origin with the closest point on the line to the origin. The normal form of the equation of a straight line on the plane is given by:

\[
x \cos \varphi + y \sin \varphi - p = 0 ,
\]

where \(\varphi\) is the angle of inclination of the normal segment (the oriented angle from the unit vector of the *x*-axis to this segment), and *p* is the (positive) length of the normal segment. The normal form can be derived from the standard form \(ax + by = c\) by dividing all of the coefficients by

\[
\sqrt{a^2 + b^2}.
\]

and also multiplying through by \(-1\) if \(c <0.\)

Unlike the slope-intercept and intercept forms, this form can represent any line but also requires only two finite parameters, \(\varphi\) and *p*, to be specified. If *p* > 0, then \(\varphi\) is uniquely defined modulo 2*π*. On the other hand, if the line is through the origin (*c* = *p* = 0), one drops the *c*/{abs|*c*} term to compute \(\sin\varphi\) and \(\cos\varphi\), and it follows that \(\varphi\) is only defined modulo π.

## Other representations
### Vectors
The vector equation of the line through points A and B is given by \(\mathbf{r} = \mathbf{OA} + \lambda\, \mathbf{AB}\) (where λ is a scalar).

If **a** is vector **OA** and **b** is vector **OB**, then the equation of the line can be written: \(\mathbf{r} = \mathbf{a} + \lambda (\mathbf{b} - \mathbf{a})\).

A ray starting at point *A* is described by limiting λ. One ray is obtained if λ ≥ 0, and the opposite ray comes from λ ≤ 0.

### Polar coordinates

In a Cartesian plane, polar coordinates (*r*, *θ*) are related to Cartesian coordinates by the parametric equations:
\[
x=r\cos\theta, \quad y=r\sin\theta.
\]


In polar coordinates, the equation of a line not passing through the origin—the point with coordinates (0, 0)—can be written

\[
r = \frac p {\cos (\theta-\varphi)},
\]

with *r* > 0 and \(\varphi-\pi/2 < \theta < \varphi + \pi/2.\) Here, \(p\) is the (positive) length of the line segment perpendicular to the line and delimited by the origin and the line, and \(\varphi\) is the (oriented) angle from the \(x\)-axis to this segment.

It may be useful to express the equation in terms of the angle \(\alpha=\varphi+\pi/2\) between the \(x\)-axis and the line. In this case, the equation becomes

\[
r=\frac p {\sin (\theta-\alpha)},
\]

with *r* > 0 and \(0 < \theta < \alpha + \pi.\)

These equations can be derived from the normal form of the line equation by setting \(x = r \cos\theta,\) and \(y = r \sin\theta,\) and then applying the angle difference identity for sine or cosine.

These equations can also be proven geometrically by applying right triangle definitions of sine and cosine to the right triangle that has a point of the line and the origin as vertices, and the line and its perpendicular through the origin as sides.

The previous forms do not apply for a line passing through the origin, but a simpler formula can be written: the polar coordinates \((r, \theta)\) of the points of a line passing through the origin and making an angle of \(\alpha\) with the \(x\)-axis, are the pairs \((r, \theta)\) such that

\[
r\ge 0,\qquad \text{and} \quad \theta=\alpha \quad\text{or}\quad \theta=\alpha +\pi.
\]


## Generalizations
Geometric space

In modern mathematics, given the multitude of geometries, the concept of a line is closely tied to the way the geometry is described. For instance, in analytic geometry, a line in the plane is often defined as the set of points whose coordinates satisfy a given linear equation, but in a more abstract setting, such as incidence geometry, a line may be an independent object, distinct from the set of points which lie on it.

When a geometry is described by a set of axioms, the notion of a line is usually left undefined (a so-called primitive object). The properties of lines are then determined by the axioms which refer to them. One advantage to this approach is the flexibility it gives to users of the geometry. Thus in differential geometry, a line may be interpreted as a geodesic (shortest path between points), while in some projective geometries, a line is a 2-dimensional vector space (all linear combinations of two independent vectors). This flexibility also extends beyond mathematics and, for example, permits physicists to think of the path of a light ray as being a line.

### Projective geometry


In many models of projective geometry, the representation of a line rarely conforms to the notion of the "straight curve" as it is visualised in Euclidean geometry. In elliptic geometry we see a typical example of this. In the spherical representation of elliptic geometry, lines are represented by great circles of a sphere with diametrically opposite points identified. In a different model of elliptic geometry, lines are represented by Euclidean planes passing through the origin. Even though these representations are visually distinct, they satisfy all the properties (such as, two points determining a unique line) that make them suitable representations for lines in this geometry.

The "shortness" and "straightness" of a line, interpreted as the property that the distance along the line between any two of its points is minimized (see triangle inequality), can be generalized and leads to the concept of geodesics in metric spaces.

## Related concepts
=== Ray ===
Orthant
Given a line and any point *A* on it, we may consider *A* as decomposing this line into two parts.
Each such part is called a **ray** (or **half-line**) and the point *A* is called its *initial point*. It is sometimes called a **half-axis** if it plays a distinct role, e.g., the positive and negative half-axes of a coordinate axis. It is a one-dimensional half-space. The point A is considered to be a member of the ray.On occasion we may consider a ray without its initial point. Such rays are called *open* rays, in contrast to the typical ray which would be said to be *closed*. Intuitively, a ray consists of those points on a line passing through *A* and proceeding indefinitely, starting at *A*, in one direction only along the line. However, in order to use this concept of a ray in proofs a more precise definition is required.

Given distinct points *A* and *B*, they determine a unique ray with initial point *A*. As two points define a unique line, this ray consists of all the points between *A* and *B* (including *A* and *B*) and all the points *C* on the line through *A* and *B* such that *B* is between *A* and *C*. This is, at times, also expressed as the set of all points *C* on the line determined by *A* and *B* such that *A* is not between *B* and *C*. A point *D*, on the line determined by *A* and *B* but not in the ray with initial point *A* determined by *B*, will determine another ray with initial point *A*. With respect to the *AB* ray, the *AD* ray is called the *opposite ray*.

Thus, we would say that two different points, *A* and *B*, define a line and a decomposition of this line into the disjoint union of an open segment *A*, *B* and two rays, *BC* and *AD* (the point *D* is not drawn in the diagram, but is to the left of *A* on the line *AB*). These are not opposite rays since they have different initial points.

In Euclidean geometry two rays with a common endpoint form an angle.

The definition of a ray depends upon the notion of betweenness for points on a line. It follows that rays exist only for geometries for which this notion exists, typically Euclidean geometry or affine geometry over an ordered field. On the other hand, rays do not exist in projective geometry nor in a geometry over a non-ordered field, like the complex numbers or any finite field.

### Line segment


A line segment is a part of a line that is bounded by two distinct end points and contains every point on the line between its end points. Depending on how the line segment is defined, either of the two end points may or may not be part of the line segment. Two or more line segments may have some of the same relationships as lines, such as being parallel, intersecting, or skew, but unlike lines they may be none of these, if they are coplanar and either do not intersect or are collinear.

### Number line


A point on number line corresponds to a real number and vice versa. Usually, integers are evenly spaced on the line, with positive numbers are on the right, negative numbers on the left. As an extension to the concept, an imaginary line representing imaginary numbers can be drawn perpendicular to the number line at zero. The two lines forms the complex plane, a geometrical representation of the set of complex numbers.

