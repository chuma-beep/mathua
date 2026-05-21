> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Surface_area) — CC BY-SA 4.0

# Surface area of rectangular prism

The **surface area** (symbol ***A***) of a solid object is a measure of the total area that the surface of the object occupies. The mathematical definition of surface area in the presence of curved surfaces is considerably more involved than the definition of arc length of one-dimensional curves, or of the surface area for polyhedra (i.e., objects with flat polygonal faces), for which the surface area is the sum of the areas of its faces. Smooth surfaces, such as a sphere, are assigned surface area using their representation as parametric surfaces. This definition of surface area is based on methods of infinitesimal calculus and involves partial derivatives and double integration.

A general definition of surface area was sought by Henri Lebesgue and Hermann Minkowski at the turn of the twentieth century. Their work led to the development of geometric measure theory, which studies various notions of surface area for irregular objects of any dimension. An important example is the Minkowski content of a surface.

## Definition
While the areas of many simple surfaces have been known since antiquity, a rigorous mathematical *definition* of area requires a great deal of care.
This should provide a function

\(S \mapsto A(S)\)

which assigns a positive real number to a certain class of surfaces that satisfies several natural requirements. The most fundamental property of the surface area is its **additivity**: *the area of the whole is the sum of the areas of the parts*. More rigorously, if a surface *S* is a union of finitely many pieces *S*\(_{1}\), …, *S*\(_{*r*}\) which do not overlap except at their boundaries, then
\(A(S) = A(S_1) + \cdots + A(S_r).\)

Surface areas of flat polygonal shapes must agree with their geometrically defined area. Since surface area is a geometric notion, areas of congruent surfaces must be the same and the area must depend only on the shape of the surface, but not on its position and orientation in space. This means that surface area is invariant under the group of Euclidean motions. These properties uniquely characterize surface area for a wide class of geometric surfaces called *piecewise smooth*. Such surfaces consist of finitely many pieces that can be represented in the parametric form

\(S_D: \vec{r}=\vec{r}(u,v), \quad (u,v)\in D\)

with a continuously differentiable function \(\vec{r}.\) The area of an individual piece is defined by the formula

\(A(S_D) = \iint_D\left |\vec{r}_u\times\vec{r}_v\right | \, du \, dv.\)

Thus the area of *S*\(_{*D*}\) is obtained by integrating the length of the normal vector \(\vec{r}_u\times\vec{r}_v\) to the surface over the appropriate region *D* in the parametric *uv* plane. The area of the whole surface is then obtained by adding together the areas of the pieces, using additivity of surface area. The main formula can be specialized to different classes of surfaces, giving, in particular, formulas for areas of graphs *z* = *f*(*x*,*y*) and surfaces of revolution.

One of the subtleties of surface area, as compared to arc length of curves, is that surface area cannot be defined simply as the limit of areas of polyhedral shapes approximating a given smooth surface. It was demonstrated by Hermann Schwarz that already for the cylinder, different choices of approximating flat surfaces can lead to different limiting values of the area; this example is known as the Schwarz lantern.

Various approaches to a general definition of surface area were developed in the late nineteenth and the early twentieth century by Henri Lebesgue and Hermann Minkowski. While for piecewise smooth surfaces there is a unique natural notion of surface area, if a surface is very irregular, or rough, then it may not be possible to assign an area to it at all. A typical example is given by a surface with spikes spread throughout in a dense fashion. Many surfaces of this type occur in the study of fractals. Extensions of the notion of area which partially fulfill its function and may be defined even for very badly irregular surfaces are studied in geometric measure theory. A specific example of such an extension is the Minkowski content of the surface.

## Common formulas

{| class="wikitable"
|+ Surface areas of common solids
|-
!Shape
!Formula/Equation
!Variables
|-
|Cube
|\(6a^2\)
|*a* = side length
|-
|Cuboid
|\(2\left(lb+lh+bh\right)\)
|
**l* = length
**b* = breadth
**h* = height

|-
|Triangular prism
|\(bh+l\left(p+q+r\right)\)
|
**b* = base length of triangle,
**h* = height of triangle,
**l* = distance between triangular bases,
**p*, *q*, *r* = sides of triangle

|-
|All prisms
|\(2B+Ph\)
|
**B* = the area of one base
**P* = the perimeter of one base
**h* = height

|-
|Sphere
|\(4\pi r^2=\pi d^2\)
|
**r* = radius of sphere
**d* = diameter

|-
|Hemisphere
|\(3\pi r^2\)
|*r* = radius of the hemisphere
|-
|Hemispherical shell
|\(\pi \left(3R^2+r^2\right)\)
|
**R* = external radius of hemisphere
**r* = internal radius of hemisphere

|-
|Spherical lune
|\(2r^2\theta\)
|
**r* = radius of sphere
**θ* = dihedral angle

|-
|Torus
|\(\left(2\pi r\right)\left(2\pi R\right)=4\pi^2Rr\)
|
**r* = minor radius (radius of the tube)
**R* = major radius (distance from center of tube to center of torus)

|-
|Closed cylinder
|\(2\pi r^2+2\pi rh=2\pi r\left(r+h\right)\)
|
**r* = radius of the circular base
**h* = height of the cylinder

|-
|Cylindrical annulus
|\(\begin{align}
& 2\pi Rh + 2\pi rh + 2\left(\pi R^2 - \pi r^2\right) \\
&= 2\pi \left(R+r\right)\left(R-r+h\right)
\end{align}\)
|
**R* = External radius
**r* = Internal radius
**h* = height

|-
|Capsule
|\(2\pi r(2r+h)\)
|
**r* = radius of the hemispheres and cylinder
**h* = height of the cylinder

|-
|Curved surface area of a cone
|\(\pi r\sqrt{r^2+h^2}=\pi rs\)
|
*\(s=\sqrt{r^2+h^2}\)
**s* = slant height of the cone
**r* = radius of the circular base
**h* = height of the cone

|-
|Full surface area of a cone
|\(\pi r\left(r+\sqrt{r^2+h^2}\right)=\pi r\left(r +s\right)\)
|
**s* = slant height of the cone
**r* = radius of the circular base
**h* = height of the cone

|-
|Regular Pyramid
|\(B+\frac{Ps}{2}\)
|
**B* = area of base
**P* = perimeter of base
**s* = slant height

|-
|Square pyramid
|\(b^2 + 2bs = b^2+ 2b\sqrt{\left(\frac{b}{2}\right)^2+h^2}\)
|
**b* = base length
**s* = slant height
**h* = vertical height

|-
|Rectangular pyramid
|\(lb+l\sqrt{\left(\frac{b}{2}\right)^2+h^2}+ b\sqrt{\left(\frac{l}{2}\right)^2+h^2}\)
|
**l* = length
**b* = breadth
**h* = height

|-
|Tetrahedron
|\(\sqrt{3}a^2\)
|*a* = side length
|-
|Surface of revolution
|\(2\pi \int_a^b {f(x) \sqrt{1+(f'(x))^2} dx}\)
|
|-
|Parametric surface
|\(\iint_D \left \vert \vec{r}_u \times \vec{r}_v \right \vert dA\)
|
*\(\vec{r}\) = parametric vector equation of surface
*\(\vec{r}_u\) = partial derivative of \(\vec{r}\) with respect to \(u\)
*\(\vec{r}_v\) = partial derivative of \(\vec{r}\) with respect to \(v\)
*\(D\) = shadow region

|}

### Ratio of surface areas of a sphere and cylinder of the same radius and height

The below given formulas can be used to show that the surface area of a sphere and cylinder of the same radius and height are in the ratio **2 : 3**, as follows.

Let the radius be *r* and the height be *h* (which is 2*r* for the sphere).

\[
\begin{array}{rlll}
\text{Sphere surface area} & = 4 \pi r^2 & & = (2 \pi r^2) \times 2 \\
\text{Cylinder surface area} & = 2 \pi r (h + r) & = 2 \pi r (2r + r) & = (2 \pi r^2) \times 3
\end{array}
\]

The discovery of this ratio is credited to Archimedes.

## In chemistry

Surface area is important in chemical kinetics. Increasing the surface area of a substance generally increases the rate of a chemical reaction. For example, iron in a fine powder will combust, while in solid blocks it is stable enough to use in structures. For different applications a minimal or maximal surface area may be desired.

## In biology

The surface area of an organism is important in several considerations, such as regulation of body temperature and digestion. Animals use their teeth to grind food down into smaller particles, increasing the surface area available for digestion. The epithelial tissue lining the digestive tract contains microvilli, greatly increasing the area available for absorption. Elephants have large ears, allowing them to regulate their own body temperature. In other instances, animals will need to minimize surface area; for example, people will fold their arms over their chest when cold to minimize heat loss.

The surface area to volume ratio (SA:V) of a cell imposes upper limits on size, as the volume increases much faster than does the surface area, thus limiting the rate at which substances diffuse from the interior across the cell membrane to interstitial spaces or to other cells. Indeed, representing a cell as an idealized sphere of radius , the volume and surface area are, respectively, 1=*V* = (4/3)*πr*\(^{3}\) and 1=*SA* = 4*πr*\(^{2}\). The resulting surface area to volume ratio is therefore 3/*r*. Thus, if a cell has a radius of 1 μm, the SA:V ratio is 3; whereas if the radius of the cell is instead 10 μm, then the SA:V ratio becomes 0.3. With a cell radius of 100, SA:V ratio is 0.03. Thus, the surface area falls off steeply with increasing volume.
