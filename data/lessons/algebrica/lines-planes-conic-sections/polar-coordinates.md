> Content sourced from [Algebrica](https://algebrica.org/polar-coordinates/) — CC BY-NC 4.0

## Radial–angular description of a point

The Cartesian coordinate system describes a point in the plane by projecting it onto two perpendicular axes. This representation privileges horizontal and vertical directions. There are situations, however, in which distance from a fixed point and direction relative to a fixed ray are more natural geometric descriptors. This leads to the polar coordinate system. Fix in the plane:

  * A point \\( O \\), called the pole.
  * A reference half-line starting from \\( O \\), called the polar axis.
  * A counterclockwise orientation.


Every point \\( Q \neq O \\) determines a distance from the pole \\(\rho = |OQ|\\) and an oriented angle \\( \theta \\) between the polar axis and the ray \\( OQ \\). The ordered pair \\((\rho, \theta)\\)is called a system of polar coordinates for \\( Q \\). The first component \\( \rho \\) is called the radius [vector](<../vectors/>), and the second component \\( \theta \\) is called the anomaly.

![Polar coordinates.](/diagrams/algebrica/polar-coordinates-1-1.png)

Let \\( (x,y) \\) be the Cartesian coordinates of \\( Q \\), and \\( (\rho,\theta) \\) its polar coordinates. Consider the [right triangle](<../right-triangle-trigonometry/>) formed by the origin \\( O \\), the point \\( Q \\), and its projection onto the \\( x \\)-axis. The hypotenuse is the segment \\( OQ \\) of length \\( \rho \\), and the angle at the origin is \\( \theta \\). Resolving this segment into its horizontal and vertical components yields:

\\[\begin{align} x &= \rho \cos\theta \\\\[6pt] y &= \rho \sin\theta \end{align} \\]

The point is thus reconstructed by projecting the segment of length \\( \rho \\), inclined at angle \\( \theta \\) with respect to the polar axis, onto each coordinate axis in turn. Conversely, starting from the Cartesian description of \\( Q \\), the radial coordinate follows directly from the [Pythagorean theorem](<../pythagorean-theorem/>) applied to the same right triangle:

\\[\rho = \sqrt{x^2 + y^2} \\]

The quantity \\( \rho \\) therefore measures the Euclidean distance of the point from the origin, regardless of the direction from which it is approached.

## Determining the angle \\(\theta\\)

Given the Cartesian coordinates of a point, recovering the radial coordinate \\( \rho \\) is straightforward. Recovering the angular coordinate \\( \theta \\) requires more care, because the trigonometric functions involved are not injective on the full circle and the naive approach through the tangent leads to an ambiguity that must be resolved geometrically. If \\( x \neq 0 \\), we can divide the second transformation formula by the first:

\\[\frac{y}{x} = \frac{\rho \sin\theta}{\rho \cos\theta} \\]

Since \\( \rho > 0 \\), the factor \\( \rho \\) cancels:

\\[\frac{y}{x} = \frac{\sin\theta}{\cos\theta} \\]

By definition of the [tangent function](<../tangent-function/>), we obtain:

\\[\tan\theta = \frac{y}{x} \\]

However, this equation alone does not determine \\( \theta \\) uniquely, since:

\\[\tan\theta = \tan(\theta + \pi) \\]

The correct angle must be chosen according to the quadrant in which the point lies. In practice, \\( \theta \\) is the angle satisfying:

\\[\cos\theta = \frac{x}{\rho} \qquad \sin\theta = \frac{y}{\rho} \\]

These two conditions together pin down \\( \theta \\) unambiguously, since the signs of \\( \cos\theta \\) and \\( \sin\theta \\) identify the quadrant, resolving the ambiguity left open by the tangent alone.

If \\( Q = O \\), then \\(\rho = 0\\) In this case, the angular coordinate becomes irrelevant. Every ray from the pole passes through the origin, so the angle carries no geometric meaning. Thus, the origin is represented by

\\[O = (0, \theta) \quad \forall \, \theta \\]

This reflects a structural feature of the system: angular information collapses when radial distance vanishes.

## Non-uniqueness of polar representation

Polar coordinates are inherently non-unique. For every integer \\( k \\) we have:

\\[(\rho,\theta) = (\rho,\theta + 2\pi k) \\]

Moreover, the same geometric point can also be represented by changing the sign of the radial coordinate while simultaneously adding \\( \pi \\) to the angular coordinate. Indeed, reversing the direction of the radial segment and rotating it by half a turn leaves the endpoint unchanged. Thus, we obtain:

\\[(\rho,\theta) = (-\rho, \theta + \pi) \\]

Thus, infinitely many pairs represent the same geometric point. To obtain a canonical representation, one typically imposes:

\\[\rho \ge 0, \qquad \theta \in [0,2\pi) \\]

This non-uniqueness reflects the rotational symmetry of the plane: angles are periodic, and directions can be traversed in opposite orientations.

## Example 1

To consolidate the procedure described above, we work through a complete conversion from Cartesian to polar coordinates. The point chosen lies in the second quadrant, which is precisely where the ambiguity of the tangent becomes relevant and must be resolved explicitly. Consider the point:

\\[(x,y) = (-3,\ \sqrt{3}) \\]

We begin by computing the radial coordinate. Applying the Pythagorean relation, we have:

\\[\rho = \sqrt{(-3)^2 + (\sqrt{3})^2} = \sqrt{9 + 3} = 2\sqrt{3} \\]

We then turn to the angular coordinate. The tangent of \\( \theta \\) is given by:

\\[\tan\theta = \frac{\sqrt{3}}{-3} = -\frac{\sqrt{3}}{3} \\]

As noted in the previous section, this [equation](<../equations/>) admits two solutions in \\( [0, 2\pi) \\), differing by \\( \pi \\). To select the correct one, we observe that the point has \\( x < 0 \\) and \\( y > 0 \\), placing it in the second quadrant. The angle consistent with this quadrant is:

\\[\theta = \frac{5\pi}{6} \\]

One can verify this directly:

\\[\cos\frac{5\pi}{6} = -\frac{\sqrt{3}}{2} < 0 \\] \\[\sin\frac{5\pi}{6} = \frac{1}{2} > 0 \\]

which matches the signs of \\( x \\) and \\( y \\) respectively.

Therefore, one polar representation of the point is:

\\[\left(2\sqrt{3},\ \frac{5\pi}{6}\right) \\]

## Example 2

The conversion in the opposite direction, from polar to Cartesian coordinates, is more direct, since it requires no quadrant analysis. We choose a point whose polar form is clean but whose Cartesian form is not obvious at first glance, so that the calculation is worth carrying out in full. Consider, for example, the point given in polar coordinates by:

\\[\left(\sqrt{6},\ \frac{7\pi}{4}\right)\\]

The angle \\( \frac{7\pi}{4} \\) lies in the fourth quadrant, just short of a full rotation. We apply the transformation formulas directly:

\\[x = \rho\cos\theta = \sqrt{6}\cdot\cos\frac{7\pi}{4} \\]

Recalling that \\( \cos\frac{7\pi}{4} = \frac{\sqrt{2}}{2} \\), we obtain:

\\[x = \sqrt{6}\cdot\frac{\sqrt{2}}{2} = \frac{\sqrt{12}}{2} = \frac{2\sqrt{3}}{2} = \sqrt{3} \\]

For the vertical component, we have:

\\[y = \rho\sin\theta = \sqrt{6}\cdot\sin\frac{7\pi}{4} \\]

Since \\( \sin\dfrac{7\pi}{4} = -\dfrac{\sqrt{2}}{2} \\), the same calculation gives:

\\[y = \sqrt{6}\cdot\left(-\frac{\sqrt{2}}{2}\right) = -\sqrt{3} \\]

###### The angle \\( \frac{7\pi}{4} \\) lies in the fourth quadrant and can be written as \\( 2\pi - \frac{\pi}{4} \\). Since the sine function is negative in the fourth quadrant and has the same absolute value as the reference angle \\( \frac{\pi}{4} \\), we obtain \\(\sin\left(\frac{7\pi}{4}\right) = -\frac{\sqrt{2}}{2}.\\)

The Cartesian representation of the point is therefore:

\\[(x, y) = \left(\sqrt{3},\ -\sqrt{3}\right). \\]

## Canonical representation and bijectivity

The non-uniqueness of polar coordinates suggests a natural question: although many pairs \\( (\rho, \theta) \\) may represent the same point, can we choose a single, preferred representative? This is indeed possible, once we restrict the range of the coordinates in a suitable way.

Let us consider the correspondence \\( f \\) that assigns to every point \\( Q \neq O \\) the pair \\( (\rho, \theta) \\) satisfying:

\\[\rho > 0, \qquad \theta \in [0, 2\pi) \\]

With these restrictions in place, each point of the punctured plane \\( \mathbb{R}^2 \setminus {O} \\) determines exactly one such pair, and conversely each admissible pair determines exactly one point. In other words, \\( f \\) establishes a bijection between the punctured plane and the half-open strip \\( (0,+\infty) \times [0,2\pi) \\).

To see why injectivity holds, take two distinct points \\( Q \\) and \\( Q’ \\), neither equal to \\( O \\). Each determines a unique ray from the pole.

  * If these rays are different, then they form different angles with the polar axis, so the corresponding angular coordinates must differ.
  * If instead the two points lie on the same ray, then they must sit at different distances from the pole, and therefore their radial coordinates are different. In either situation, the associated pairs cannot coincide.


The reverse implication is obtained by simply reversing the construction. Given any pair \\( (\rho, \theta) \\) with \\( \rho > 0 \\) and \\( \theta \in [0,2\pi) \\), we consider the ray forming angle \\( \theta \\) with the polar axis and mark on it the point at distance \\( \rho \\) from the pole. This produces a unique point \\( Q \neq O \\), and by construction it is mapped back to the original pair. No further ambiguity appears once \\( \rho \\) is required to be strictly positive.


The origin must be treated separately. When \\( \rho = 0 \\) every ray from the pole passes through the same point. Including the origin in the domain would therefore break injectivity. For this reason, it is customary to write its coordinates as \\( (0,\theta) \\) for an arbitrary \\( \theta \\), while understanding that the angular component carries no geometric information in this degenerate case.

## Polar coordinates in space

The idea of polar coordinates extends naturally to three dimensions by introducing a radial distance and two angular parameters. Let \\( O \\) be the origin of a Cartesian reference frame in space. For a point \\( Q \\) with Cartesian coordinates \\( (x, y, z) \\), we denote by \\( \rho \\) the Euclidean distance from the origin:

\\[\rho = \sqrt{x^2 + y^2 + z^2} \\]

To describe the direction of \\( Q \\), we proceed in two steps.

  * First, we consider the orthogonal projection of \\( Q \\) onto the plane \\( XY \\), and denote by \\( \theta \\) the polar angle of this projection with respect to the positive \\( x \\)-axis.
  * Second, we introduce an angle \\( \psi \\), measured from the positive \\( z \\)-axis to the segment \\( OQ \\).
  * The triple \\( (\rho, \theta, \psi) \\) provides a radial–angular description of the point in space.


![Polar coordinates in space.](/diagrams/algebrica/polar-coordinates-2-1.png)

Here \\( \rho \\) denotes the radial distance, \\( \theta \\) the azimuthal angle in the \\( XY \\)-plane, and \\( \psi \\) the zenith measured from the positive \\( z \\)-axis. From the [right-triangle relations](<../right-triangle-trigonometry/>) in the associated geometry, one obtains:

\\[\begin{align} x &= \rho \sin\psi \cos\theta \\\\[6pt] y &= \rho \sin\psi \sin\theta \\\\[6pt] z &= \rho \cos\psi \end{align} \\]

These formulas show that the Cartesian coordinates are recovered by decomposing the segment \\( OQ \\) into a vertical component of length \\( \rho \cos\psi \\) and a horizontal component of length \\( \rho \sin\psi \\), the latter being further resolved in the plane through the usual planar polar relations. This coordinate system is particularly suited to problems exhibiting spherical symmetry, where distance from the origin plays a more fundamental role than alignment with coordinate axes.

## Integration in Polar Coordinates

One of the main reasons for introducing polar coordinates appears in [integration](<../indefinite-integrals/>) over planar regions possessing radial symmetry. When a region is more naturally described in terms of distance from the origin and angular spread, Cartesian coordinates may obscure the underlying structure. In polar coordinates, the elementary area element is not simply the product of two independent differentials. A small region determined by variations ( d\rho ) and ( d\theta ) has area:

\\[dA = \rho \, d\rho \, d\theta \\]

The additional factor \\( \rho \\) reflects the fact that circular arcs increase in length proportionally to the distance from the origin. Accordingly, a double integral over a region \\( D \\) can be rewritten as:

\\[\iint_D f(x,y)\, dx\,dy = \iint_{D’} f(\rho\cos\theta,\rho\sin\theta)\, \rho \, d\rho\, d\theta \\]

where \\( D’ \\) denotes the corresponding region in the \\( (\rho,\theta) \\)-plane. This transformation often simplifies computations when the geometry of the problem is naturally expressed in radial terms.

## Summary

Property| Plane| Space  
---|---|---  
Coordinates| \\( (\rho,\ \theta) \\)| \\( (\rho,\ \theta,\ \psi) \\)  
Radial distance| \\( \rho = \sqrt{x^2 + y^2} \\)| \\( \rho = \sqrt{x^2 + y^2 + z^2} \\)  
From polar to Cartesian| \\( x = \rho\cos\theta \\) \\( y = \rho\sin\theta \\)| \\( x = \rho\sin\psi\cos\theta \\) \\( y = \rho\sin\psi\sin\theta \\) \\( z = \rho\cos\psi \\)  
Angular constraint| \\( \theta \in [0, 2\pi) \\)| \\( \theta \in [0, 2\pi),\ \psi \in [0, \pi] \\)  
Canonical form| \\( \rho \ge 0,\ \theta \in [0, 2\pi) \\)| \\( \rho \ge 0,\ \theta \in [0, 2\pi),\ \psi \in [0,\pi] \\)  
  
## Selected references

  * **MIT, S. Widnall, J. Peraire**. [Other Coordinate Systems: Polar, Cylindrical and Spherical](https://ocw.mit.edu/courses/16-07-dynamics-fall-2009/57081b546fff23e6b88dbac0ab859c7d_MIT16_07F09_Lec05.pdf)

  * **University of Maryland**. [Polar Coordinates](https://math.umd.edu/~tjp/141%2010.1%20lecture%20notes.pdf)

  * **University of South Carolina**. [Polar Coordinates](https://people.math.sc.edu/josephcf/Teaching/142/Files/Lecture%20Notes/Chapter11/11.3.pdf)

  * **Purdue University**. [Polar Coordinates and Integration](https://www.math.purdue.edu/~neptamin/324Au17/Notes/15.4/15.4.pdf)
