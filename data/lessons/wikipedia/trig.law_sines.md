> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Law_of_sines) — CC BY-SA 4.0

# Law of sines

In trigonometry, the **law of sines** (sometimes called the **sine formula** or **sine rule**) is a mathematical equation relating the lengths of the sides of any triangle to the sines of its angles. According to the law,

\[
\frac{a}{\sin{\alpha}} \,=\, \frac{b}{\sin{\beta}} \,=\, \frac{c}{\sin{\gamma}} \,=\, 2R,
\]

where *a*, *b*, and *c* are the lengths of the sides of a triangle, and *α*, *β*, and *γ* are the opposite angles (see figure 2), while *R* is the radius of the triangle's circumcircle. When the last part of the equation is not used, the law is sometimes stated using the reciprocals;

\[
\frac{\sin{\alpha}}{a} \,=\, \frac{\sin{\beta}}{b} \,=\, \frac{\sin{\gamma}}{c}.
\]

The law of sines can be used to compute the remaining sides of a triangle when two angles and a side are known—a technique known as triangulation. It can also be used when two sides and one of the non-enclosed angles are known. In some such cases, the triangle is not uniquely determined by this data (called the *ambiguous case*) and the technique gives two possible values for the enclosed angle.

The law of sines is one of two trigonometric equations commonly applied to find lengths and angles in scalene triangles, with the other being the law of cosines.

The law of sines can be generalized to higher dimensions on surfaces with constant curvature.

## Proof
With the side of length  as the base, the triangle's altitude can be computed as *b* sin *γ* or as *c* sin *β*.  Equating these two expressions gives
\(\frac{b}{\sin \beta} = \frac{c}{\sin \gamma}\,,\)
and similar equations arise by choosing the side of length  or the side of length  as the base of the triangle. For a proof that these expressions are equal to \(2R\), see Relation to the circumcircle.

## Ambiguous case of triangle solution
When using the law of sines to find a side of a triangle, an ambiguous case occurs when two separate triangles can be constructed from the data provided (i.e., there are two different possible solutions to the triangle). In the case shown below they are triangles *ABC* and *ABC′*.


Given a general triangle, the following conditions would need to be fulfilled for the case to be ambiguous:
* The only information known about the triangle is the angle *α* and the sides *a* and *c*.
* The angle *α* is acute (i.e., *α* < 90°).
* The side *a* is shorter than the side *c* (i.e., *a* < *c*).
* The side *a* is longer than the altitude *h* from angle *β*, where 1=*h* = *c* sin *α* (i.e., *a* > *h*).

If all the above conditions are true, then each of angles *β* and *β′* produces a valid triangle, meaning that both of the following are true:

\[
{\gamma}' = \arcsin\frac{c \sin{\alpha}}{a} \quad \text{or} \quad {\gamma} = \pi - \arcsin\frac{c \sin{\alpha}}{a}.
\]


From there we can find the corresponding *β* and *b* or *β′* and *b′* if required, where *b* is the side bounded by vertices *A* and *C* and *b′* is bounded by *A* and *C′*.

## Examples
The following are examples of how to solve a problem using the law of sines.

### Example 1

Given: side 1=*a* = 20, side 1=*c* = 24, and angle 1=*γ* = 40°. Angle *α* is desired.

Using the law of sines, we conclude that

\[
\frac{\sin \alpha}{20} = \frac{\sin (40^\circ)}{24}.
\]


\[
\alpha = \arcsin\left( \frac{20\sin (40^\circ)}{24} \right) \approx 32.39^\circ.
\]


Note that the potential solution 1=*α* = 147.61° is excluded because that would necessarily give 1=*α* + *β* + *γ* > 180°.

### Example 2

If the lengths of two sides of the triangle *a* and *b* are equal to *x*, the third side has length *c*, and the angles opposite the sides of lengths *a*, *b*, and *c* are *α*, *β*, and *γ* respectively then

\[
\begin{align}
& \alpha = \beta = \frac{180^\circ-\gamma}{2}= 90^\circ-\frac{\gamma}{2} \\[6pt]
& \sin \alpha = \sin \beta = \sin \left(90^\circ-\frac{\gamma}{2}\right) = \cos \left(\frac{\gamma}{2}\right) \\[6pt]
& \frac{c}{\sin \gamma}=\frac{a}{\sin \alpha}=\frac{x}{\cos \left(\frac{\gamma}{2}\right)} \\[6pt]
& \frac{c \cos \left(\frac{\gamma}{2}\right)}{\sin \gamma} = x
\end{align}
\]


## Relation to the circumcircle
In the identity

\[
\frac{a}{\sin{\alpha}} = \frac{b}{\sin{\beta}} = \frac{c}{\sin{\gamma}},
\]

the common value of the three fractions is actually the diameter of the triangle's circumcircle. This result dates back to Ptolemy.

### Proof


As shown in the figure, let there be a circle with inscribed \(\triangle ABC\) and another inscribed \(\triangle ADB\) that passes through the circle's center **O**. The \(\angle AOD\) has a central angle of \(180^\circ\) and thus \(\angle ABD = 90^\circ\), by Thales's theorem. Since \(\triangle ABD\) is a right triangle,

\[
\sin{\delta}= \frac{\text{opposite}}{\text{hypotenuse}}= \frac{c}{2R},
\]

where \(R= \frac{d}{2}\) is the radius of the circumscribing circle of the triangle. Angles \({\gamma}\) and \({\delta}\) lie on the same circle and subtend the same chord *c*; thus, by the inscribed angle theorem, {{nowrap|\({\gamma} = {\delta}\).}} Therefore,

\[
\sin{\delta} = \sin{\gamma} = \frac{c}{2R}.
\]


Rearranging yields

\[
2R = \frac{c}{\sin{\gamma}}.
\]


Repeating the process of creating \(\triangle ADB\) with other points gives

{{equation box 1|equation=\(\frac{a}{\sin{\alpha}} = \frac{b}{\sin{\beta}} = \frac{c}{\sin{\gamma}}=2R.\)}}

### Relationship to the area of the triangle
The area of a triangle is given by {{nowrap|\(T = \frac{1}{2}ab \sin \theta\),}} where \(\theta\) is the angle enclosed by the sides of lengths *a* and *b*. Substituting the sine law into this equation gives

\[
T=\frac{1}{2}ab \cdot \frac {c}{2R}.
\]


Taking \(R\) as the circumscribing radius,

{{equation box 1|equation=\(T=\frac{abc}{4R}.\)}}

It can also be shown that this equality implies

\[
\begin{align}
\frac{abc} {2T}
& = \frac{abc} {2\sqrt{s(s-a)(s-b)(s-c)}} \\[6pt]
& = \frac {2abc} {\sqrt{{(a^2+b^2+c^2)}^2-2(a^4+b^4+c^4) }},
\end{align}
\]

where *T* is the area of the triangle and *s* is the semiperimeter {{nowrap|\(s = \frac{1}{2}\left(a+b+c\right).\)}}

The second equality above readily simplifies to Heron's formula for the area.

The sine rule can also be used in deriving the following formula for the triangle's area: denoting the semi-sum of the angles' sines as {{nowrap|\(S =\frac{1}{2}\left(\sin A + \sin B + \sin C\right)\),}} we have

{{equation box 1|equation=\(T = 4R^{2} \sqrt{S \left(S - \sin A\right) \left(S - \sin B\right) \left(S - \sin C\right)}\)}}

where \(R\) is the radius of the circumcircle: {{nowrap|\(2R = \frac{a}{\sin A} = \frac{b}{\sin B} = \frac{c}{\sin C}\).}}

## Spherical law of sines
The spherical law of sines deals with triangles on a sphere, whose sides are arcs of great circles.

Suppose the radius of the sphere is 1. Let *a*, *b*, and *c* be the lengths of the great-arcs that are the sides of the triangle. Because it is a unit sphere, *a*, *b*, and *c* are the angles at the center of the sphere subtended by those arcs, in radians. Let *A*, *B*, and *C* be the angles opposite those respective sides. These are dihedral angles between the planes of the three great circles.

Then the spherical law of sines says:

\[
\frac{\sin A}{\sin a} = \frac{\sin B}{\sin b} = \frac{\sin C}{\sin c}.
\]


### Vector proof
Consider a unit sphere with three unit vectors **OA**, **OB** and **OC** drawn from the origin to the vertices of the triangle. Thus the angles *α*, *β*, and *γ* are the angles *a*, *b*, and *c*, respectively. The arc BC subtends an angle of magnitude *a* at the centre. Introduce a Cartesian basis with **OA** along the *z*-axis and **OB** in the *xz*-plane making an angle *c* with the *z*-axis. The vector **OC** projects to ON in the *xy*-plane and the angle between ON and the *x*-axis is *A*. Therefore, the three vectors have components:

\[
\mathbf{OA} = \begin{pmatrix}0 \\ 0 \\ 1\end{pmatrix}, \quad
\mathbf{OB} = \begin{pmatrix}\sin c \\ 0 \\ \cos c\end{pmatrix}, \quad
\mathbf{OC} = \begin{pmatrix}\sin b\cos A \\ \sin b\sin A \\ \cos b\end{pmatrix}.
\]


The scalar triple product, **OA** ⋅ (**OB** × **OC**) is the volume of the parallelepiped formed by the position vectors of the vertices of the spherical triangle **OA**, **OB** and **OC**. This volume is invariant to the specific coordinate system used to represent **OA**, **OB** and **OC**. The value of the scalar triple product **OA** ⋅ (**OB** × **OC**) is the 3 × 3 determinant with **OA**, **OB** and **OC** as its rows. With the *z*-axis along **OA** the square of this determinant is

\[
\begin{align}
\bigl(\mathbf{OA} \cdot (\mathbf{OB} \times \mathbf{OC})\bigr)^2
& = \left(\det \begin{pmatrix}\mathbf{OA} & \mathbf{OB} & \mathbf{OC}\end{pmatrix}\right)^2 \\[4pt]
& =  \begin{vmatrix}
0 & 0 & 1 \\
\sin c & 0 & \cos c \\
\sin b \cos A & \sin b \sin A & \cos b
\end{vmatrix} ^2
= \left(\sin b \sin c \sin A\right)^2.
\end{align}
\]

Repeating this calculation with the *z*-axis along **OB** gives (sin *c* sin *a* sin *B*)<sup>2</sup>, while with the *z*-axis along **OC** it is (sin *a* sin *b* sin *C*)<sup>2</sup>. Equating these expressions and dividing throughout by (sin *a* sin *b* sin *c*)<sup>2</sup> gives

\[
\frac{\sin^2 A}{\sin^2 a}
= \frac{\sin^2 B}{\sin^2 b}
= \frac{\sin^2 C}{\sin^2 c}
= \frac{V^2}{\sin^2 (a) \sin^2 (b) \sin^2 (c)},
\]

where  is the volume of the parallelepiped formed by the position vector of the vertices of the spherical triangle. Consequently, the result follows.

It is easy to see how for small spherical triangles, when the radius of the sphere is much greater than the sides of the triangle, this formula becomes the planar formula at the limit, since

\[
\lim_{a \to 0} \frac{\sin a}{a} = 1
\]

and the same for sin *b* and sin *c*.


### Geometric proof
Consider a unit sphere with:

\[
OA = OB = OC = 1
\]


Construct point \(D\) and point \(E\) such that \(\angle ADO = \angle AEO = 90^\circ\)

Construct point \(A'\) such that \(\angle A'DO = \angle A'EO = 90^\circ\)

It can therefore be seen that \(\angle ADA' = B\) and \(\angle AEA' = C\)

Notice that \(A'\) is the projection of \(A\) on plane \(OBC\). Therefore \(\angle AA'D = \angle AA'E = 90^\circ\)

By basic trigonometry, we have:

\[
\begin{align}
  AD &= \sin c \\
  AE &= \sin b
\end{align}
\]


But \(AA' = AD \sin B = AE \sin C\)

Combining them we have:

\[
\begin{align}
\sin c \sin B &= \sin b \sin C \\
\Rightarrow \frac{\sin B}{\sin b} &=\frac{\sin C}{\sin c}
\end{align}
\]


By applying similar reasoning, we obtain the spherical law of sines:

\[
\frac{\sin A}{\sin a} =\frac{\sin B}{\sin b} =\frac{\sin C}{\sin c}
\]


### Other proofs
A purely algebraic proof can be constructed from the spherical law of cosines. From the identity \(\sin^2 A = 1 - \cos^2 A\) and the explicit expression for \(\cos A\) from the spherical law of cosines

\[
\begin{align}
   \sin^2\!A &= 1-\left(\frac{\cos a  - \cos b\, \cos c}{\sin b \,\sin c}\right)^2\\
   &=\frac{\left(1-\cos^2\!b\right) \left(1-\cos^2\!c\right)-\left(\cos a  - \cos b\, \cos c\right)^2}
          {\sin^2\!b \,\sin^2\!c}\\[8pt]
 \frac{\sin A}{\sin a}
 &= \frac{\left[1-\cos^2\!a-\cos^2\!b-\cos^2\!c + 2\cos a\cos b\cos c\right]^{1/2}}{\sin a\sin b\sin c}.
\end{align}
\]

Since the right hand side is invariant under a cyclic permutation of \(a,\;b,\;c\) the spherical sine rule follows immediately.

The figure used in the Geometric proof above is used by and also provided in Banerjee (see Figure 3 in this paper) to derive the sine law using elementary linear algebra and projection matrices.

## Hyperbolic case
In hyperbolic geometry when the curvature is −1, the law of sines becomes

\[
\frac{\sin A}{\sinh a} = \frac{\sin B}{\sinh b} = \frac{\sin C}{\sinh c} \,.
\]


In the special case when *B* is a right angle, one gets

\[
\sin C = \frac{\sinh c}{\sinh b}
\]


which is the analog of the formula in Euclidean geometry expressing the sine of an angle as the opposite side divided by the hypotenuse.


## The case of surfaces of constant curvature
Define a generalized sine function, depending also on a real parameter \(\kappa\):

\[
\sin_\kappa(x) = x - \frac{\kappa}{3!}x^3 + \frac{\kappa^2}{5!}x^5 - \frac{\kappa^3}{7!}x^7 + \cdots = \sum_{n=0}^\infty \frac{(-1)^n \kappa^n}{(2n+1)!}x^{2n+1}.
\]


The law of sines in constant curvature \(\kappa\) reads as

\[
\frac{\sin A}{\sin_\kappa a} = \frac{\sin B}{\sin_\kappa b} = \frac{\sin C}{\sin_\kappa c} \,.
\]


By substituting \(\kappa=0\), \(\kappa=1\), and \(\kappa=-1\), one obtains respectively \(\sin_{0}(x) = x\), \(\sin_{1}(x) = \sin x\), and \(\sin_{-1}(x) = \sinh x\), that is, the Euclidean, spherical, and hyperbolic cases of the law of sines described above.

Let \(p_\kappa(r)\) indicate the circumference of a circle of radius \(r\) in a space of constant curvature \(\kappa\). Then \(p_\kappa(r)=2\pi\sin_\kappa(r)\). Therefore, the law of sines can also be expressed as:

\[
\frac{\sin A}{p_\kappa(a)} = \frac{\sin B}{p_\kappa(b)} = \frac{\sin C}{p_\kappa(c)} \,.
\]


This formulation was discovered by János Bolyai.

## Higher dimensions
A tetrahedron has four triangular facets. The absolute value of the polar sine (psin) of the normal vectors to the three facets that share a vertex of the tetrahedron, divided by the area of the fourth facet will not depend upon the choice of the vertex:


\[
\begin{align}
& \frac{\left|\operatorname{psin}(\mathbf{b}, \mathbf{c}, \mathbf{d})\right|}{\mathrm{Area}_a} =
  \frac{\left|\operatorname{psin}(\mathbf{a}, \mathbf{c}, \mathbf{d})\right|}{\mathrm{Area}_b} =
  \frac{\left|\operatorname{psin}(\mathbf{a}, \mathbf{b}, \mathbf{d})\right|}{\mathrm{Area}_c} =
  \frac{\left|\operatorname{psin}(\mathbf{a}, \mathbf{b}, \mathbf{c})\right|}{\mathrm{Area}_d} \\[4pt]
= {} & \frac{(3~\mathrm{Volume}_\mathrm{tetrahedron})^2}{2~\mathrm{Area}_a \mathrm{Area}_b \mathrm{Area}_c \mathrm{Area}_d}\,.
\end{align}
\]


More generally, for an *n*-dimensional simplex (i.e., triangle (1=*n* = 2), tetrahedron (1=*n* = 3), pentatope (1=*n* = 4), etc.) in *n*-dimensional Euclidean space, the absolute value of the polar sine of the normal vectors of the facets that meet at a vertex, divided by the hyperarea of the facet opposite the vertex is independent of the choice of the vertex. Writing *V* for the hypervolume of the *n*-dimensional simplex and *P* for the product of the hyperareas of its (*n* − 1)-dimensional facets, the common ratio is

\[
\frac{\left|\operatorname{psin}(\mathbf{b}, \ldots, \mathbf{z})\right|}{\mathrm{Area}_a} = \cdots = \frac{\left|\operatorname{psin}(\mathbf{a}, \ldots, \mathbf{y})\right|}{\mathrm{Area}_z} = \frac{(nV)^{n-1}}{(n-1)! P}.
\]


Note that when the vectors **v**}}, from a selected vertex to each of the other vertices, are the columns of a matrix  then the columns of the matrix \(N = -V (V^TV)^{-1} \sqrt{\det{V^TV}} / (n-1)!\) are outward-facing normal vectors of those facets that meet at the selected vertex.  This formula also works when the vectors are in a -dimensional space having *m* > *n*.  In the 1=*m* = *n* case that  is square, the formula simplifies to \(N = -(V^T)^{-1} |\det{V}| / (n-1)!\,.\)

## History
An equivalent of the law of sines, that the sides of a triangle are proportional to the chords of double the opposite angles, was known to the 2nd century Hellenistic astronomer Ptolemy and used occasionally in his *Almagest*.

Statements related to the law of sines appear in the astronomical and trigonometric work of 7th century Indian mathematician Brahmagupta. In his *Brāhmasphuṭasiddhānta*, Brahmagupta expresses the circumradius of a triangle as the product of two sides divided by twice the altitude; the law of sines can be derived by alternately expressing the altitude as the sine of one or the other base angle times its opposite side, then equating the two resulting variants. An equation even closer to the modern law of sines appears in Brahmagupta's *Khaṇḍakhādyaka*, in a method for finding the distance between the Earth and a planet following an epicycle; however, Brahmagupta never treated the law of sines as an independent subject or used it systematically for solving triangles.

The spherical law of sines is sometimes credited to 10th century scholars Abu-Mahmud Khujandi or Abū al-Wafāʾ (it appears in his *Almagest*), but it is given prominence in Abū Naṣr Manṣūr's *Treatise on the Determination of Spherical Arcs*, and was credited to Abū Naṣr Manṣūr by his student al-Bīrūnī in his *Keys to Astronomy*. Ibn Muʿādh al-Jayyānī's 11th-century *Book of Unknown Arcs of a Sphere* also contains the spherical law of sines.

The 13th-century Persian mathematician Naṣīr al-Dīn al-Ṭūsī stated and proved the planar law of sines: In any plane triangle, the ratio of the sides is equal to the ratio of the sines of the angles opposite to those sides. That is, in triangle ABC, we have AB : AC = Sin(∠ACB) : Sin(∠ABC) By employing the law of sines, al-Tusi could solve triangles where either two angles and a side were known or two sides and an angle opposite one of them were given. For triangles with two sides and the included angle, he divided them into right triangles that he could then solve. When three sides were given, he dropped a perpendicular line and then used Proposition II-13 of Euclid's *Elements* (a geometric version of the law of cosines). Al-Tusi established the important result that if the sum or difference of two arcs is provided along with the ratio of their sines, then the arcs can be calculated.

According to Glen Van Brummelen, "The Law of Sines is really Regiomontanus's foundation for his solutions of right-angled triangles in Book IV, and these solutions are in turn the bases for his solutions of general triangles." Regiomontanus was a 15th-century German mathematician.

## See also
*
* Half-side formula for solving spherical triangles
* Law of cosines
* Law of tangents
* Law of cotangents
* Mollweide's formula for checking solutions of triangles
* Solution of triangles
* Surveying

## References


## External links

*
* The Law of Sines at cut-the-knot
* Degree of Curvature
* Finding the Sine of 1 Degree
* Generalized law of sines to higher dimensions

