> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Area_of_a_triangle) — CC BY-SA 4.0

# Area of a triangle

In geometry, calculating the area of a triangle is an elementary and widely encountered problem. The best known and simplest formula for the area of an arbitrary triangle is \(\text{area} = \tfrac12 \times \text{base} \times \text{height},\) where "base" means the length of any side of the triangle, and "height" means the shortest distance between the line of the base and the corresponding *apex* (opposite corner of the triangle).

Although simple, this formula is only useful if the height can be readily found, which is not always the case. For example, the land surveyor of a triangular field might find it relatively easy to measure the length of each side, but relatively difficult to construct a 'height'. Various methods may be used in practice, depending on what is known about the triangle. Other frequently used formulas for the area of a triangle use trigonometry, side lengths (Heron's formula), vectors, coordinates, line integrals, Pick's theorem, or other properties.

## Basic formula

The *area* of a plane figure (two-dimensional shape) is a quantity which describes its size or extent. Conventionally, area is measured in units based on a square whose side is one unit long (for instance, a square meter). When using such units, the area of a rectangle can be found by multiplying together the lengths of two adjacent sides, sometimes called the *base*, and *height*,. For example, a rectangle which is 2 meters long and 5 meters tall has area 10 square meters.

A right triangle, one with two sides (called the *legs*) that are perpendicular with a right angle between them, is half of a rectangle with the two perpendicular sides as its base and height, bisected along a diagonal. Thus the area of a right triangle can be found by halving the area of the corresponding rectangle; that is, by multiplying the lengths of the two perpendicular sides and dividing the result by two: \(T = \tfrac12 b h.\) Taking either leg of the right triangle to be the *base* of the triangle, the corresponding altitude of the triangle is the other leg.

Any other triangle, choosing an arbitrary side to be its base, is likewise half of a parallelogram, again bisected along a diagonal, which can then be cut into pieces and rearranged into a rectangle with the same base length and the same height; the triangle's area can therefore be calculated using the same formula, \(T = \tfrac12 b h.\) For a general triangle, the *altitude* does not coincide with any of the sides, but is a new line segment perpendicular to the base with one of its endpoints at the *apex* (opposite corner) and the other endpoint on the straight line containing the base. The length of this line segment is called the *height* or *altitude* of the triangle, and is the shortest distance between the apex and the line containing the base.

One consequence of the triangle area formula is the analog in the plane of Lexell's theorem about spherical triangles: all of the triangles with a fixed base side, a fixed area, and the apex on the same side of the base line have their apex on a specific line parallel to the base.

## Calculations from angles

The area *\(T\)* of a triangle can be found through the application of trigonometry.

### Knowing SAS (side-angle-side)
Using the labels in the image on the right, the *height* or *altitude* is \(h = a \sin \gamma\). Substituting this in the area formula \(T=\tfrac12 bh\) derived above, the area of the triangle can be expressed as:
\(T = \tfrac12 ab\sin \gamma = \tfrac12 bc\sin \alpha = \tfrac12 ca\sin \beta\)
Where \(a\) is the line ***BC***, \(b\) is the line ***AC***, \(c\) is the line ***AB***, \(\alpha\) is the interior angle at ***A***, \(\beta\) is the interior angle at ***B***, \(\gamma\) is the interior angle at ***C***.

Furthermore; since \(\sin \alpha = \sin (\pi - \alpha) = \sin(\beta + \gamma)\), and similarly for the other two angles:
\(T = \tfrac12 ab\sin (\alpha+\beta) = \tfrac12 bc\sin (\beta+\gamma) = \tfrac12 ca\sin (\gamma+\alpha)\)

### Knowing AAS (angle-angle-side)
Since \(\cot \gamma = -\cot(\pi - \gamma) = -\cot(\alpha+\beta)\), and similarly for the other two angles:
\(T =\frac{b^{2{2(\cot\alpha-\cot(\alpha+\beta))}= \frac {b^{2}\sin \alpha\sin (\alpha + \beta)}{2\sin \beta},\)

and analogously if the known side is \(a\) or \(c\).

### Knowing ASA (angle-side-angle)
\(T = \frac{a^{2{2(\cot \beta + \cot \gamma)} = \frac{a^{2}\sin \beta\sin \gamma}{2\sin(\beta + \gamma)},\)

and analogously if the known side is \(b\) or \(c\).

## Using side lengths (Heron's formula)

A triangle's shape is uniquely determined by the lengths of the sides, so its metrical properties, including area, can be described in terms of those lengths. By Heron's formula,
\(T = \sqrt{s(s-a)(s-b)(s-c)}\)

where \(s= \tfrac12(a+b+c)\) is the semiperimeter, or half of the triangle's perimeter.

Three other equivalent ways of writing Heron's formula are
\(\begin{align}
T &= \tfrac14 \sqrt{(a^2+b^2+c^2)^2-2(a^4+b^4+c^4)} \\[5mu]
&= \tfrac14 \sqrt{2(a^2b^2+a^2c^2+b^2c^2)-(a^4+b^4+c^4)} \\[5mu]
&= \tfrac14 \sqrt{(a+b+c)(-a+b+c)(a-b+c)(a+b-c)}.
\end{align}\)

### Formulas resembling Heron's formula
Three formulas have the same structure as Heron's formula but are expressed in terms of different variables. First, denoting the medians from sides \(a\), \(b\), and \(c\) respectively as \(m_a\), \(m_b\), and \(m_c\) and their semi-sum \(\sigma = (m_a + m_b + m_c)/2\), we have
\(T = \tfrac43 \sqrt{\sigma (\sigma - m_a)(\sigma - m_b)(\sigma - m_c)}.\)

Next, denoting the altitudes from sides \(a\), \(b\), and \(c\) respectively as \(h_a\), \(h_b\), and \(h_c\), and denoting the semi-sum of the reciprocals of the altitudes as \(H = (h_a^{-1} + h_b^{-1} + h_c^{-1})/2\) we have
\(T^{-1} = 4 \sqrt{H(H-h_a^{-1})(H-h_b^{-1})(H-h_c^{-1})}.\)

And denoting the semi-sum of the angles' sines as \(S = (\sin \alpha + \sin \beta + \sin \gamma)/2\), we have
\(T = D^{2} \sqrt{S(S-\sin \alpha)(S-\sin \beta)(S-\sin \gamma)}\)

where \(D\) is the diameter of the circumcircle: \(D=\tfrac{a}{\sin \alpha} = \tfrac{b}{\sin \beta} = \tfrac{c}{\sin \gamma}.\)

## Using vectors
The area of triangle ABC is half of the area of a parallelogram:
\(T = \tfrac12\bigl\|(\mathbf b - \mathbf a) \wedge (\mathbf c - \mathbf a)\bigr\|
= \tfrac12\bigl\|\mathbf a \wedge \mathbf b + \mathbf b \wedge \mathbf c + \mathbf c \wedge \mathbf a \bigr\|,\)

where , and are vectors to the triangle's vertices from any arbitrary origin point, so that and are the translation vectors from vertex to each of the others, and is the wedge product. If vertex is taken to be the origin, this simplifies to. The oriented relative area of a parallelogram in any affine space, a type of bivector, is defined as where and are translation vectors from one vertex of the parallelogram to each of the two adjacent vertices. In Euclidean space, the magnitude of this bivector is a well-defined scalar number representing the area of the parallelogram. (For vectors in three-dimensional space, the bivector-valued wedge product has the same magnitude as the vector-valued cross product, but unlike the cross product, which is only defined in three-dimensional Euclidean space, the wedge product is well-defined in an affine space of any dimension.)

The area of triangle *ABC* can also be expressed in terms of dot products. Taking vertex to be the origin and calling translation vectors to the other vertices and ,

\(T = \tfrac12 \sqrt{\mathbf{b}^2 \mathbf{c}^2 - (\mathbf{b} \cdot \mathbf{c})^2},\)

where for any Euclidean vector. This area formula can be derived from the previous one using the elementary vector identity

In two-dimensional Euclidean space, for a vector with coordinates and vector with coordinates , the magnitude of the wedge product is
\(\| \mathbf b \wedge \mathbf c \| = |x_B y_C - x_C y_B|.\)

(See the following section.)

## Using coordinates
If vertex *A* is located at the origin (0, 0) of a Cartesian coordinate system and the coordinates of the other two vertices are given by *B* and *C* , then the area can be computed as 1/2 times the absolute value of the determinant
\(T = \tfrac12\left|\det\begin{pmatrix}x_B & x_C \\ y_B & y_C \end{pmatrix}\right| = \tfrac12 |x_B y_C - x_C y_B|.\)

For three general vertices, the equation is:
\(T = \tfrac12 \left| \det\begin{pmatrix}x_A & x_B & x_C \\ y_A & y_B & y_C \\ 1 & 1 & 1\end{pmatrix} \right| = \tfrac12 \big| x_A y_B - x_A y_C + x_B y_C - x_B y_A + x_C y_A - x_C y_B \big|,\)

which can be written as
\(T = \tfrac12 \big| (x_A - x_C) (y_B - y_A) - (x_A - x_B) (y_C - y_A) \big|.\)

If the points are labeled sequentially in the counterclockwise direction, the above determinant expressions are positive and the absolute value signs can be omitted. The above formula is known as the shoelace formula or the surveyor's formula.

If we locate the vertices in the complex plane and denote them in counterclockwise sequence as *a* , *b* , and *c* , and denote their complex conjugates as \(\bar a\), \(\bar b\), and \(\bar c\), then the formula
\(T=\frac{i}{4}\begin{vmatrix}a & \bar a & 1 \\ b & \bar b & 1 \\ c & \bar c & 1 \end{vmatrix}\)

is equivalent to the shoelace formula.

In three dimensions, the area of a general triangle *A* , *B* and *C* ) is the Pythagorean sum of the areas of the respective projections on the three principal planes (i.e. *x* = 0, *y* = 0 and *z* = 0):
\(T = \tfrac12 \sqrt{\begin{vmatrix} x_A & x_B & x_C \\ y_A & y_B & y_C \\ 1 & 1 & 1 \end{vmatrix}^2 +
\begin{vmatrix} y_A & y_B & y_C \\ z_A & z_B & z_C \\ 1 & 1 & 1 \end{vmatrix}^2 +
\begin{vmatrix} z_A & z_B & z_C \\ x_A & x_B & x_C \\ 1 & 1 & 1 \end{vmatrix}^2 }.\)

## Using line integrals
The area within any closed curve, such as a triangle, is given by the line integral around the curve of the algebraic or signed distance of a point on the curve from an arbitrary oriented straight line *L*. Points to the right of *L* as oriented are taken to be at negative distance from *L*, while the weight for the integral is taken to be the component of arc length parallel to *L* rather than arc length itself.

This method is well suited to computation of the area of an arbitrary polygon. Taking *L* to be the *x*-axis, the line integral between consecutive vertices (*x_{i}*,*y_{i}*) and (*x*_{*i*+1},*y*_{*i*+1}) is given by the base times the mean height, namely (*x*_{*i*+1} − *x_{i}*)(*y_{i}* + *y*_{*i*+1})/2. The sign of the area is an overall indicator of the direction of traversal, with negative area indicating counterclockwise traversal. The area of a triangle then falls out as the case of a polygon with three sides.

While the line integral method has in common with other coordinate-based methods the arbitrary choice of a coordinate system, unlike the others it makes no arbitrary choice of vertex of the triangle as origin or of side as base. Furthermore, the choice of coordinate system defined by *L* commits to only two degrees of freedom rather than the usual three; since the weight is a local distance (e.g. *x*_{*i*+1} − *x_{i}* in the above) whence the method does not require choosing an axis normal to *L*.

When working in polar coordinates it is not necessary to convert to Cartesian coordinates to use line integration; since the line integral between consecutive vertices (*r_{i}*,θ_{*i*}) and (*r*_{*i*+1},θ_{*i*+1}) of a polygon is given directly by *r_{i}r*_{*i*+1}sin(θ_{*i*+1} − θ_{*i*})/2. This is valid for all values of θ, with some decrease in numerical accuracy when |θ| is many orders of magnitude greater than π. With this formulation negative area indicates clockwise traversal, which should be kept in mind when mixing polar and cartesian coordinates. Just as the choice of *y*-axis (*x* ) is immaterial for line integration in cartesian coordinates, so is the choice of zero heading (θ ) immaterial here.

## Using Pick's theorem
See Pick's theorem for a technique for finding the area of any arbitrary lattice polygon (one drawn on a grid with vertically and horizontally adjacent lattice points at equal distances, and with vertices on lattice points).

The theorem states:
\(T = I + \tfrac12 B - 1\)

where *\(I\)* is the number of internal lattice points and *B* is the number of lattice points lying on the border of the polygon.

## Other area formulas

Numerous other area formulas exist, such as
\(T = rs,\)

where *r* is the inradius, and *s* is the semiperimeter (in fact, this formula holds for *all* tangential polygons), and

\(T=r_a(s-a)=r_b(s-b)=r_c(s-c)\)

where \(r_a, \, r_b,\, r_c\) are the radii of the excircles tangent to sides *a, b, c* respectively.

We also have

\(T = \tfrac12 D^{2}(\sin \alpha)(\sin \beta)(\sin \gamma)\)

and
\(T = \frac{abc}{2D} = \frac{abc}{4R}\)

for circumdiameter *D*; and
\(T = \tfrac14(\tan \alpha)(b^{2}+c^{2}-a^{2})\)

for angle α ≠ 90°.

The area can also be expressed as
\(T = \sqrt{rr_ar_br_c}.\)

In 1885, Baker gave a collection of over a hundred distinct area formulas for the triangle. These include:
\(T = \tfrac12\sqrt[3]{abch_ah_bh_c},\)
\(T = \tfrac12 \sqrt{abh_ah_b},\)
\(T = \frac{a+b}{2(h_a^{-1} + h_b^{-1})},\)
\(T = \frac{Rh_bh_c}{a}\)

for circumradius (radius of the circumcircle) *R*, and
\(T = \frac{h_ah_b}{2 \sin \gamma}.\)

## Upper bound on the area
The area *T* of any triangle with perimeter *p* satisfies

\(T\le \tfrac{p^2}{12\sqrt{3,\)

with equality holding if and only if the triangle is equilateral.

Other upper bounds on the area *T* are given by

\(4\sqrt{3}T \leq a^2+b^2+c^2\)

and

\(4\sqrt{3}T \leq \frac{9abc}{a+b+c},\)

both with equality holding if and only if the triangle is equilateral.

## History
In 300 BCE Greek mathematician Euclid proved that the area of a triangle is half that of a parallelogram with the same base and height in his book *Elements of Geometry*. Heron of Alexandria found what is known as Heron's formula for the area of a triangle in terms of its sides, and a proof can be found in his book, *Metrica*, written around 60 CE. It has been suggested that Archimedes knew the formula over two centuries earlier, and since *Metrica* is a collection of the mathematical knowledge available in the ancient world, it is possible that the formula predates the reference given in that work.

In 499 CE Aryabhata, a mathematician-astronomer from the classical age of Indian mathematics and Indian astronomy, expressed the area of a triangle as one-half the base times the height in the *Aryabhatiya*.

A formula equivalent to Heron's was discovered by the Chinese independently of the Greeks. It was published in 1247 in *Shushu Jiuzhang* ("Mathematical Treatise in Nine Sections"), written by Qin Jiushao.

## Bisecting the area
There are infinitely many lines that bisect the area of a triangle. Three of them are the medians, which are the only area bisectors that go through the centroid. Three other area bisectors are parallel to the triangle's sides.

Any line through a triangle that splits both the triangle's area and its perimeter in half goes through the triangle's incenter. There can be one, two, or three of these for any given triangle.
