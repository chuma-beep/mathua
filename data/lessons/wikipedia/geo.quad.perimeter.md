> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Perimeter) — CC BY-SA 4.0

# Perimeter of polygons

A **perimeter** is the length of a closed boundary that encompasses, surrounds, or outlines either a two-dimensional shape or a one-dimensional line. The perimeter of a circle or an ellipse is called its circumference.

Calculating the perimeter has several practical applications. A calculated perimeter is the length of fence required to surround a yard or garden. The perimeter of a wheel/circle (its circumference) describes how far it will roll in one revolution. Similarly, the amount of string wound around a spool is related to the spool's perimeter; if the length of the string was exact, it would equal the perimeter.

## Formulas
{| class="wikitable sortable mw-collapsible"
|+
! shape !! formula || variables
|-
| circle || \(2 \pi r = \pi d\) || where \(r\) is the radius of the circle and \(d\) is the diameter.
|-
| semicircle || \((\pi+2)r\)|| where \(r\) is the radius of the semicircle.
|-
| triangle || \(a + b + c\,\) || where \(a\), \(b\) and \(c\) are the lengths of the sides of the triangle.
|-
| square/rhombus || \(4a\) || where \(a\) is the side length.
|-
| rectangle || \(2(l+w)\) || where \(l\) is the length and \(w\) is the width.
|-
| equilateral polygon || \(n \times a\,\) || where \(n\) is the number of sides and \(a\) is the length of one of the sides.
|-
| regular polygon || \(2nb \sin\left(\frac{\pi}{n}\right)\) || where \(n\) is the number of sides and \(b\) is the distance between center of the polygon and one of the vertices of the polygon.
|-
| general polygon || \(a_1 + a_2 + a_3 + \cdots + a_n = \sum_{i=1}^n a_i\) || where \(a_{i}\) is the length of the \(i\)-th (1st, 2nd, 3rd ... *n*th) side of an *n*-sided polygon.
|}

The perimeter is the distance around a shape. Perimeters for more general shapes can be calculated, as any path, with \(\int_0^L \mathrm{d}s\), where \(L\) is the length of the path and \(ds\) is an infinitesimal line element. Both of these must be replaced by algebraic forms in order to be practically calculated. If the perimeter is given as a closed piecewise smooth plane curve \(\gamma: [a,b] \to \mathbb{R}^2\) with
\(\gamma(t)=\begin{pmatrix}x(t)\\y(t)\end{pmatrix}\)
then its length \(L\) can be computed as follows:
\(L = \int_a^b \sqrt{x'(t)^2+y'(t)^2}\,\mathrm dt\)

A generalized notion of perimeter, which includes hypersurfaces bounding volumes in \(n\)-dimensional Euclidean spaces, is described by the theory of Caccioppoli sets.

## Polygons

Polygons are fundamental to determining perimeters, not only because they are the simplest shapes but also because the perimeters of many shapes are calculated by approximating them with sequences of polygons tending to these shapes. The first mathematician known to have used this kind of reasoning is Archimedes, who approximated the perimeter of a circle by surrounding it with regular polygons.

The perimeter of a polygon equals the sum of the lengths of its sides (edges). In particular, the perimeter of a rectangle of width \(w\) and length \(\ell\) equals \(2w + 2\ell.\)

An equilateral polygon is a polygon which has all sides of the same length (for example, a rhombus is a 4-sided equilateral polygon). To calculate the perimeter of an equilateral polygon, one must multiply the common length of the sides by the number of sides.

A regular polygon may be characterized by the number of its sides and by its circumradius, that is to say, the constant distance between its centre and each of its vertices. The length of its sides can be calculated using trigonometry. If *R* is a regular polygon's radius and *n* is the number of its sides, then its perimeter is
\(2nR \sin\left(\frac{180^{\circ}}{n}\right).\)

A splitter of a triangle is a cevian (a segment from a vertex to the opposite side) that divides the perimeter into two equal lengths, this common length being called the semiperimeter of the triangle. The three splitters of a triangle all intersect each other at the Nagel point of the triangle.

A cleaver of a triangle is a segment from the midpoint of a side of a triangle to the opposite side such that the perimeter is divided into two equal lengths. The three cleavers of a triangle all intersect each other at the triangle's Spieker center.

## Circumference of a circle


The perimeter of a circle, often called the circumference, is proportional to its diameter and its radius. That is to say, there exists a constant number pi,  (the Greek *p* for perimeter), such that if *P* is the circle's perimeter and *D* its diameter then,
\(P = \pi\cdot{D}.\!\)

In terms of the radius *r* of the circle, this formula becomes,

\(P=2\pi\cdot r.\)

To calculate a circle's perimeter, knowledge of its radius or diameter and the number  suffices. The problem is that  is not rational (it cannot be expressed as the quotient of two integers), nor is it algebraic (it is not a root of a polynomial equation with rational coefficients). So, obtaining an accurate approximation of  is important in the calculation. The computation of the digits of  is relevant to many fields, such as mathematical analysis, algorithmics and computer science.

## Perception of perimeter


The perimeter and the area are two main measures of geometric figures. Confusing them is a common error, as well as believing that the greater one of them is, the greater the other must be. Indeed, a commonplace observation is that an enlargement (or a reduction) of a shape make its area grow (or decrease) as well as its perimeter. For example, if a field is drawn on a 1/ scale map, the actual field perimeter can be calculated multiplying the drawing perimeter by . The real area is  times the area of the shape on the map. Nevertheless, there is no relation between the area and the perimeter of an ordinary shape. For example, the perimeter of a rectangle of width 0.001 and length 1000 is slightly above 2000, while the perimeter of a rectangle of width 0.5 and length 2 is 5. Both areas are equal to 1.

Proclus (5th century) reported that Greek peasants "fairly" parted fields relying on their perimeters. However, a field's production is proportional to its area, not to its perimeter, so many naive peasants may have gotten fields with long perimeters but small areas (thus, few crops).

If one removes a piece from a figure, its area decreases but its perimeter may not. The convex hull of a figure may be visualized as the shape formed by a rubber band stretched around it. In the animated picture on the left, all the figures have the same convex hull; the big, first hexagon.

## Isoperimetry


The isoperimetric problem is to determine a figure with the largest area, amongst those having a given perimeter. The solution is intuitive; it is the circle. In particular, this can be used to explain why drops of fat on a broth surface are circular.

This problem may seem simple, but its mathematical proof requires some sophisticated theorems. The isoperimetric problem is sometimes simplified by restricting the type of figures to be used. In particular, to find the quadrilateral, or the triangle, or another particular figure, with the largest area amongst those with the same shape having a given perimeter. The solution to the quadrilateral isoperimetric problem is the square, and the solution to the triangle problem is the equilateral triangle. In general, the polygon with *n* sides having the largest area and a given perimeter is the regular polygon, which is closer to being a circle than is any irregular polygon with the same number of sides.

## Etymology
The word comes from the Greek περίμετρος *perimetros*, from περί *peri* "around" and μέτρον *metron* "measure".

## See also

* Coastline paradox
* Girth (geometry)
* Wetted perimeter


## References


}}

## External links


*

