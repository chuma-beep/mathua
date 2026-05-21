> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Area_of_a_circle) — CC BY-SA 4.0

# Area of a circle

In geometry, the area enclosed by a circle of radius is π*r*\(^{2}\). Here, the Greek letter represents the constant ratio of the circumference of any circle to its diameter, approximately equal to 3.14159.

One method of deriving this formula, which originated with Archimedes, involves viewing the circle as the limit of a sequence of regular polygons with an increasing number of sides. The area of a regular polygon is half its perimeter multiplied by the distance from its center to its sides, and because the sequence tends to a circle, the corresponding formula–that the area is half the circumference times the radius–namely, *A* = , holds for a circle.

## Terminology
Although often referred to as the area of a circle in informal contexts, strictly speaking, the term *disk* refers to the interior region of the circle, while circle is reserved for the boundary only, which is a curve and covers no area itself. Therefore, the *area of a disk* is the more precise phrase for the area enclosed by a circle.

## History
Modern mathematics can obtain the area using the methods of integral calculus or its more sophisticated offspring, real analysis. However, the area of a disk was studied by the Ancient Greeks. Eudoxus of Cnidus in the fifth century B.C. had found that the area of a disk is proportional to its radius squared. Archimedes used the tools of Euclidean geometry to show that the area inside a circle is equal to that of a right triangle whose base has the length of the circle's circumference and whose height equals the circle's radius in his book *Measurement of a Circle*. The circumference is 2*πr*, and the area of a triangle is half the base times the height, yielding the area *πr*\(^{2}\) for the disk. Prior to Archimedes, Hippocrates of Chios was the first to show that the area of a disk is proportional to the square of its diameter, as part of his quadrature of the lune of Hippocrates, but did not identify the constant of proportionality.

## Historical arguments
A variety of arguments have been advanced historically to establish the equation \(A=\pi r^2\) to varying degrees of mathematical rigor. The most famous of these is Archimedes' method of exhaustion, one of the earliest uses of the mathematical concept of a limit, as well as the origin of Archimedes' axiom which remains part of the standard analytical treatment of the real number system. The original proof of Archimedes is not rigorous by modern standards, because it assumes that we can compare the length of arc of a circle to the length of a secant and a tangent line, and similar statements about the area, as geometrically evident.

### Using polygons
The area of a regular polygon is half its perimeter times the apothem. As the number of sides of the regular polygon increases, the polygon tends to a circle, and the apothem tends to the radius. This suggests that the area of a disk is half the circumference of its bounding circle times the radius.

### Archimedes' proof
Following Archimedes' argument in *The Measurement of a Circle* (c. 260 BCE), compare the area enclosed by a circle to a right triangle whose base has the length of the circle's circumference and whose height equals the circle's radius. If the area of the circle is not equal to that of the triangle, then it must be either greater or less. We eliminate each of these by contradiction, leaving equality as the only possibility. We use regular polygons in the same way.

#### Not greater

Suppose that the area *C* enclosed by the circle is greater than the area *T* = *cr*/2 of the triangle. Let *E* denote the excess amount. Inscribe a square in the circle, so that its four corners lie on the circle. Between the square and the circle are four segments. If the total area of those gaps, *G*\(_{4}\), is greater than *E*, split each arc in half. This makes the inscribed square into an inscribed octagon, and produces eight segments with a smaller total gap, *G*\(_{8}\). Continue splitting until the total gap area, *G\(_{n}\)*, is less than *E*. Now the area of the inscribed polygon, *P\(_{n}\)* = *C* − *G\(_{n}\)*, must be greater than that of the triangle.
\(\begin{align}
 E &{}= C - T \\
 &{}> G_n \\
 P_n &{}= C - G_n \\
 &{}> C - E \\
 P_n &{}> T
\end{align}\)
But this forces a contradiction, as follows. Draw a perpendicular from the center to the midpoint of a side of the polygon; its length, *h*, is less than the circle radius. Also, let each side of the polygon have length *s*; then the sum of the sides is *ns*, which is less than the circle circumference. The polygon area consists of *n* equal triangles with height *h* and base *s*, thus equals *nhs*/2. But since *h* < *r* and *ns* < *c*, the polygon area must be less than the triangle area, *cr*/2, a contradiction. Therefore, our supposition that *C* might be greater than *T* must be wrong.

#### Not less

Suppose that the area enclosed by the circle is less than the area *T* of the triangle. Let *D* denote the deficit amount. Circumscribe a square, so that the midpoint of each edge lies on the circle. If the total area gap between the square and the circle, *G*\(_{4}\), is greater than *D*, slice off the corners with circle tangents to make a circumscribed octagon, and continue slicing until the gap area is less than *D*. The area of the polygon, *P\(_{n}\)*, must be less than *T*.
\(\begin{align}
 D &{}= T - C \\
 &{}> G_n \\
 P_n &{}= C + G_n \\
 &{}< C + D \\
 P_n &{}< T
\end{align}\)
This, too, forces a contradiction. For, a perpendicular to the midpoint of each polygon side is a radius, of length *r*. And since the total side length is greater than the circumference, the polygon consists of *n* identical triangles with total area greater than *T*. Again we have a contradiction, so our supposition that *C* might be less than *T* must be wrong as well.

Therefore, it must be the case that the area enclosed by the circle is precisely the same as the area of the triangle. This concludes the proof.

### Rearrangement proof

Following Satō Moshun , Nicholas of Cusa and Leonardo da Vinci , we can use inscribed regular polygons in a different way. Suppose we inscribe a hexagon. Cut the hexagon into six triangles by splitting it from the center. Two opposite triangles both touch two common diameters; slide them along one so the radial edges are adjacent. They now form a parallelogram, with the hexagon sides making two opposite edges, one of which is the base, *s*. Two radial edges form slanted sides, and the height, *h* is equal to its apothem (as in the Archimedes proof). In fact, we can also assemble all the triangles into one big parallelogram by putting successive pairs next to each other. The same is true if we increase it to eight sides and so on. For a polygon with 2*n* sides, the parallelogram will have a base of length *ns*, and a height *h*. As the number of sides increases, the length of the parallelogram base approaches half the circle circumference, and its height approaches the circle radius. In the limit, the parallelogram becomes a rectangle with width *r* and height *r*.

{| class="wikitable" frame="vsides" style="text-align:center" cellspacing="0" cellpadding="3"
|+ ***Unit disk area by rearranging n polygons.***
|-
! colspan="2" | polygon
| rowspan="11" style="padding:1px;"|
! colspan="3" | parallelogram
|-
! *n* !! side !! base !! height !! area
|-
| align="right" | 4 || 1.4142136 || 2.8284271 || 0.7071068 || 2.0000000
|-
| align="right" | 6 || 1.0000000 || 3.0000000 || 0.8660254 || 2.5980762
|-
| align="right" | 8 || 0.7653669 || 3.0614675 || 0.9238795 || 2.8284271
|-
| align="right" | 10 || 0.6180340 || 3.0901699 || 0.9510565 || 2.9389263
|-
| align="right" | 12 || 0.5176381 || 3.1058285 || 0.9659258 || 3.0000000
|-
| align="right" | 14 || 0.4450419 || 3.1152931 || 0.9749279 || 3.0371862
|-
| align="right" | 16 || 0.3901806 || 3.1214452 || 0.9807853 || 3.0614675
|-
| align="right" | 96 || 0.0654382 || 3.1410320 || 0.9994646 || 3.1393502
|-
| ∞ || 1/∞ || || 1 ||
|}

## Modern proofs
There are various equivalent definitions of the constant π. The conventional definition in pre-calculus geometry is the ratio of the circumference of a circle to its diameter:
\(\pi=\frac{C}{D}.\)
However, because the circumference of a circle is not a primitive analytical concept, this definition is not suitable in modern rigorous treatments. A standard modern definition is that is equal to twice the least positive root of the cosine function or, equivalently, the half-period of the sine (or cosine) function. The cosine function can be defined either as a power series, or as the solution of a certain differential equation. This avoids any reference to circles in the definition of , so that statements about the relation of to the circumference and area of circles are actually theorems, rather than definitions, that follow from the analytical definitions of concepts like "area" and "circumference".

The analytical definitions are seen to be equivalent, if it is agreed that the circumference of the circle is measured as a rectifiable curve by means of the integral
\(C = 2\int_{-R}^R \frac{R\,dx}{\sqrt{R^2-x^2= 2R\int_{-1}^1\frac{dx}{\sqrt{1-x^2.\)
The integral appearing on the right is an abelian integral whose value is a half-period of the sine function, equal to. Thus \(C=2\pi R=\pi D\) is seen to be true as a theorem.

Several of the arguments that follow use only concepts from elementary calculus to reproduce the formula \(A=\pi r^2\), but in many cases to regard these as actual proofs, they rely implicitly on the fact that one can develop trigonometric functions and the fundamental constant in a way that is totally independent of their relation to geometry. We have indicated where appropriate how each of these proofs can be made totally independent of all trigonometry, but in some cases that requires more sophisticated mathematical ideas than those afforded by elementary calculus.

### Onion proof

Using calculus, we can sum the area incrementally, partitioning the disk into thin concentric rings like the layers of an onion. This is the method of shell integration in two dimensions. For an infinitesimally thin ring of the "onion" of radius *t*, the accumulated area is 2*t dt*, the circumferential length of the ring times its infinitesimal width (one can approximate this ring by a rectangle with width=2*t* and height=*dt*). This gives an elementary integral for a disk of radius *r*.
\(\begin{align}
 \mathrm{Area}(r) &{}= \int_0^{r} 2 \pi t \, dt \\
 &{}= 2\pi \left[\frac{t^2}{2} \right]_{0}^{r}\\
 &{}= \pi r^2.
\end{align}\)

It is rigorously justified by the multivariate substitution rule in polar coordinates. Namely, the area is given by a double integral of the constant function 1 over the disk itself. If *D* denotes the disk, then the double integral can be computed in polar coordinates as follows:
\(\begin{align}
 \mathrm{Area}(r) &{}= \iint_D 1\ d(x, y)\\
 &{} = \iint_D t\ dt\ d\theta\\
 &{} = \int_0^r \int_0^{2\pi} t\ d\theta\ dt\\
 &{} = \int_0^r \left[ t\theta \right]_{0}^{2\pi} dt\\
 &{}= \int_0^{r} 2 \pi t \, dt \\
 \end{align}\)
which is the same result as obtained above.

An equivalent rigorous justification, without relying on the special coordinates of trigonometry, uses the coarea formula. Define a function \(\rho:\mathbb R^2\to\mathbb R\) by \(\rho(x,y)=\sqrt{x^2+y^2}\). Note ρ is a Lipschitz function whose gradient is a unit vector \(|\nabla\rho|=1\) (almost everywhere). Let *D* be the disc \(\rho<1\) in \(\mathbb R^2\). We will show that \(\mathcal L^2(D)=\pi\), where \(\mathcal L^2\) is the two-dimensional Lebesgue measure in \(\mathbb R^2\). We shall assume that the one-dimensional Hausdorff measure of the circle \(\rho=r\) is \(2\pi r\), the circumference of the circle of radius *r*. (This can be taken as the definition of circumference.) Then, by the coarea formula,
\(\begin{align}
\mathcal L^2(D) &= \iint_D |\nabla \rho|\,d\mathcal{L}^2\\
&= \int_{\mathbb R} \mathcal H^1(\rho^{-1}(r)\cap D)\,dr\\
&= \int_0^1\mathcal H^1(\rho^{-1}(r))\,dr \\
&= \int_0^1 2\pi r\, dr= \pi.
\end{align}\)

### Triangle proof

Similar to the onion proof outlined above, we could exploit calculus in a different way in order to arrive at the formula for the area of a disk. Consider unwrapping the concentric circles to straight strips. This will form a right angled triangle with r as its height and 2r (being the outer slice of onion) as its base.

Finding the area of this triangle will give the area of the disk
\(\begin{align}
\text{Area} &{}= \frac{1}{2} \cdot \text{base} \cdot \text{height} \\[6pt]
 &{}= \frac{1}{2} \cdot 2 \pi r \cdot r \\[6pt]
 &{}= \pi r^2
\end{align}\)

The opposite and adjacent angles for this triangle are respectively in degrees 9.0430611..., 80.956939... and in radians 0.1578311... , 1.4129651....

Explicitly, we imagine dividing up a circle into triangles, each with a height equal to the circle's radius and a base that is infinitesimally small. The area of each of these triangles is equal to \(1/2\cdot r \cdot du\). By summing up (integrating) all of the areas of these triangles, we arrive at the formula for the circle's area:
\(\begin{align}
 \mathrm{Area}(r) &{}= \int_0^{2\pi r} \frac{1}{2} r \, du \\[6pt]
 &{}= \left[ \frac{1}{2} r u \right]_{0}^{2 \pi r}\\[6pt]
 &{}= \pi r^2.
\end{align}\)

It too can be justified by a double integral of the constant function 1 over the disk by reversing the order of integration and using a change of variables in the above iterated integral:
\(\begin{align}
 \mathrm{Area}(r) &{}= \iint_D 1\ d(x, y)\\
 &{} = \iint_D t\ dt\ d\theta\\
 &{} = \int_0^{2\pi} \int_0^r t\ dt\ d\theta\\
 &{} = \int_0^{2\pi} \frac{1}{2}r^2\ d\theta\\
 \end{align}\)
Making the substitution \(u = r\theta,\ du = r\ d\theta\) converts the integral to
\(\int_0^{2\pi r} \frac{1}{2} \frac{r^2}{r} du = \int_0^{2\pi r} \frac{1}{2} r\ du\)
which is the same as the above result.

The triangle proof can be reformulated as an application of Green's theorem in flux-divergence form (i.e. a two-dimensional version of the divergence theorem), in a way that avoids all mention of trigonometry and the constant. Consider the vector field \(\mathbf r = x\mathbf i + y\mathbf j\) in the plane. So the divergence of **r** is equal to two, and hence the area of a disc *D* is equal to
\(A = \frac 1 2 \iint_D \operatorname{div}\mathbf r\, dA.\)
By Green's theorem, this is the same as the outward flux of **r** across the circle bounding *D*:
\(A = \frac 1 2 \oint_{\partial D} \mathbf r\cdot\mathbf n\, ds\)
where **n** is the unit normal and *ds* is the arc length measure. For a circle of radius *R* centered at the origin, we have \(|\mathbf r|=R\) and \(\mathbf n=\mathbf r/R\), so the above equality is
\(A = \frac 1 2\oint_{\partial D} \mathbf r\cdot\frac{\mathbf r}{R}\, ds = \frac{R}{2}\oint_{\partial D} \,ds.\)

The integral of *ds* over the whole circle \(\partial D\) is just the arc length, which is its circumference, so this shows that the area *A* enclosed by the circle is equal to \(R/2\) times the circumference of the circle.

Another proof that uses triangles considers the area enclosed by a circle to be made up of an infinite number of triangles (i.e. the triangles each have an angle of *dθ* at the center of the circle), each with an area of ½·*r*·*r*·sin(*dθ*) = ½·*r*^{2}·*dθ* (derived from the expression for the area of a triangle: ½·*ab*·sin(*C*)). Note that sin(*dθ*) ≈ *dθ* due to small angle approximation. Through summing the areas of the triangles, the expression for the area of the circle can therefore be found:

\[
\begin{align}
 \mathrm{Area}&{}= \int_0^{2\pi} \frac{1}{2} r^2 \, d \theta \\
 &{}= \left[ \frac{1}{2} r^2 \theta \right]_{0}^{2 \pi}\\
 &{}= \pi r^2.
\end{align}
\]

### Semicircle proof
Note that the area of a semicircle of radius *r* can be computed by the integral \(\int_{-r}^r \sqrt{r^2 - x^2}\,dx\).

By trigonometric substitution, we substitute \(x=r \sin\theta\), hence \(dx=r\cos \theta\, d\theta.\)

\[
\begin{align}\int_{-r}^r \sqrt{r^2 - x^2}\,dx
&=\int_{-\frac{\pi}{2^{\frac{\pi}{2\sqrt{r^2\left(1-\sin ^2 \theta\right)} \cdot r \cos \theta\, d \theta \\[5pt]
&=2r^2\int_{0}^{\frac{\pi}{2\cos ^2 \theta\, d \theta \\[5pt]
&=\frac{\pi r^2}{2}.
\end{align}
\]

The last step follows since the trigonometric identity \(\cos(\theta)=\sin(\pi/2-\theta)\) implies that \(\cos^2\theta\) and \(\sin^2\theta\) have equal integrals over the interval \([0,\pi/2]\), using integration by substitution. But on the other hand; since \(\cos^2\theta+\sin^2\theta=1\), the sum of the two integrals is the length of that interval, which is \(\pi/2\). Consequently, the integral of \(\cos^2 \theta\) is equal to half the length of that interval, which is \(\pi/4\).

Therefore, the area of a circle of radius *r*, which is twice the area of the semi-circle, is equal to \(2 \cdot \frac{\pi r^2}{2} = \pi r^2\).

This particular proof may appear to beg the question, if the sine and cosine functions involved in the trigonometric substitution are regarded as being defined in relation to circles. However, as noted earlier, it is possible to define sine, cosine, and in a way that is totally independent of trigonometry, in which case the proof is valid by the change of variables formula and Fubini's theorem, assuming the basic properties of sine and cosine (which can also be proved without assuming anything about their relation to circles).

## Isoperimetric inequality
The circle is the closed curve of least perimeter that encloses the maximum area. This is known as the isoperimetric inequality, which states that if a rectifiable Jordan curve in the Euclidean plane has perimeter *C* and encloses an area *A* (by the Jordan curve theorem) then
\(4\pi A\le C^2.\)
Moreover, equality holds in this inequality if and only if the curve is a circle, in which case \(A=\pi r^2\) and \(C=2\pi r\).

## Fast approximation
The calculations Archimedes used to approximate the area numerically were laborious, and he stopped with a polygon of 96 sides. A faster method uses ideas of Willebrord Snell (*Cyclometricus*, 1621), further developed by Christiaan Huygens (*De Circuli Magnitudine Inventa*, 1654), described in. ### Archimedes' doubling method
Given a circle, let *u\(_{n}\)* be the perimeter of an inscribed regular *n-*gon, and let *U\(_{n}\)* be the perimeter of a circumscribed regular *n-*gon. Then *u\(_{n}\)* and *U\(_{n}\)* are lower and upper bounds for the circumference of the circle that become sharper and sharper as *n* increases, and their average (*u\(_{n}\)* + *U\(_{n}\)*)/2 is an especially good approximation to the circumference. To compute *u\(_{n}\)* and *U\(_{n}\)* for large *n*, Archimedes derived the following doubling formulae:

\(u_{2n} = \sqrt{U_{2n} u_n}\)   (geometric mean), and

\(U_{2n} = \frac{2 U_n u_n}{ U_n + u_n}\)    (harmonic mean).

Starting from a hexagon, Archimedes doubled *n* four times to get a 96-gon, which gave him a good approximation to the circumference of the circle.

In modern notation, we can reproduce his computation (and go further) as follows.
For a unit circle, an inscribed hexagon has *u*\(_{6}\) = 6, and a circumscribed hexagon has *U*\(_{6}\) = 4.
Doubling seven times yields

{| class="wikitable" frame="vsides" style="text-align:center" cellspacing="0" cellpadding="3"
|+ **Archimedes doubling seven times; *n* = 6 × 2\(^{*k*}\).**
|- style="background-color:#eeeeee"
! *k* !! *n* !! *u\(_{n}\)* !! *U\(_{n}\)* !!
|-
| 0 || 6 || 6.0000000 || 6.9282032 || 3.2320508
|-
| 1 || 12 || 6.2116571 || 6.4307806 || 3.1606094
|-
| 2 || 24 || 6.2652572 || 6.3193199 || 3.1461443
|-
| 3 || 48 || 6.2787004 || 6.2921724 || 3.1427182
|-
| 4 || 96 || 6.2820639 || 6.2854292 || 3.1418733
|-
| 5 || 192 || 6.2829049 || 6.2837461 || 3.1416628
|-
| 6 || 384 || 6.2831152 || 6.2833255 || 3.1416102
|-
| 7 || 768 || 6.2831678 || 6.2832204 || 3.1415970
|}

(Here approximates the circumference of the unit circle, which is 2, so approximates. )

The last entry of the table has \(^{355}\)⁄\(_{113}\) as one of its best rational approximations;
i.e., there is no better approximation among rational numbers with denominator up to 113.
The number \(^{355}\)⁄\(_{113}\) is also an excellent approximation to , attributed to Chinese mathematician Zu Chongzhi, who named it Milü. This approximation is better than any other rational number with denominator less than 16,604.

### The Snell–Huygens refinement
Snell proposed (and Huygens proved) a tighter bound than Archimedes':
\(n \frac{3 \sin \frac{\pi}{n{2+\cos\frac{\pi}{n< \pi < n \left(2 \sin \frac{\pi}{3 n} + \tan \frac{\pi}{3 n}\right).\)
This for *n* = 48 gives a better approximation (about 3.14159292) than Archimedes' method for *n* = 768.

### Derivation of Archimedes' doubling formulae

Let one side of an inscribed regular *n-*gon have length *s\(_{n}\)* and touch the circle at points A and B. Let A′ be the point opposite A on the circle, so that A′A is a diameter, and A′AB is an inscribed triangle on a diameter. By Thales' theorem, this is a right triangle with right angle at B. Let the length of A′B be *c\(_{n}\)*, which we call the complement of *s\(_{n}\)*; thus *c\(_{n}\)*\(^{2}\)+*s\(_{n}\)*\(^{2}\) = (2*r*)\(^{2}\). Let C bisect the arc from A to B, and let C′ be the point opposite C on the circle. Thus the length of CA is *s*\(_{2*n*}\), the length of C′A is *c*\(_{2*n*}\), and C′CA is itself a right triangle on diameter C′C. Because C bisects the arc from A to B, C′C perpendicularly bisects the chord from A to B, say at P. Triangle C′AP is thus a right triangle, and is similar to C′CA since they share the angle at C′. Thus all three corresponding sides are in the same proportion; in particular, we have C′A : C′C = C′P : C′A and AP : C′A = CA : C′C. The center of the circle, O, bisects A′A, so we also have triangle OAP similar to A′AB, with OP half the length of A′B. In terms of side lengths, this gives us
\(\begin{align}
 c_{2n}^2 &{}= \left( r + \frac{1}{2} c_n \right) 2r \\
 c_{2n} &{}= \frac{s_n}{s_{2n.
\end{align}\)
In the first equation C′P is C′O+OP, length *r* + *c\(_{n}\)*, and C′C is the diameter, 2*r*. For a unit circle we have the famous doubling equation of Ludolph van Ceulen,
\(c_{2n} = \sqrt{2+c_n}. \)
If we now circumscribe a regular *n-*gon, with side A″B″ parallel to AB, then OAB and OA″B″ are similar triangles, with A″B″ : AB = OC : OP. Call the circumscribed side *S\(_{n}\)*; then this is *S\(_{n}\)* : *s\(_{n}\)* = 1 : \(^{1}\)⁄\(_{2}\)*c\(_{n}\)*. (We have again used that OP is half the length of A′B.) Thus we obtain
\(c_n = 2\frac{s_n}{S_n}. \)
Call the inscribed perimeter *u\(_{n}\)* = *ns\(_{n}\)*, and the circumscribed perimeter *U\(_{n}\)* = *nS\(_{n}\)*. Then combining equations, we have
\(c_{2n} = \frac{s_n}{s_{2n= 2 \frac{s_{2n{S_{2n,\)
so that
\(u_{2n}^2 = u_n U_{2n}. \)
This gives a geometric mean equation.

We can also deduce
\(2 \frac{s_{2n{S_{2n\frac{s_n}{s_{2n= 2 + 2 \frac{s_n}{S_n} ,\)
or
\(\frac{2}{U_{2n= \frac{1}{u_n} + \frac{1}{U_n}. \)
This gives a harmonic mean equation.

## Dart approximation

When more efficient methods of finding areas are not available, we can resort to "throwing darts". This Monte Carlo method uses the fact that if random samples are taken uniformly scattered across the surface of a square in which a disk resides, the proportion of samples that hit the disk approximates the ratio of the area of the disk to the area of the square. This should be considered a method of last resort for computing the area of a disk (or any shape), as it requires an enormous number of samples to get useful accuracy; an estimate good to 10\(^{−*n*}\) requires about 100\(^{*n*}\) random samples. ## Finite rearrangement
We have seen that by partitioning the disk into an infinite number of pieces we can reassemble the pieces into a rectangle. A remarkable fact discovered relatively recently is that we can dissect the disk into a large but *finite* number of pieces and then reassemble the pieces into a square of equal area. This is called Tarski's circle-squaring problem. The nature of Laczkovich's proof is such that it proves the existence of such a partition (in fact, of many such partitions) but does not exhibit any particular partition.

## Non-Euclidean circles
Circles can be defined in non-Euclidean geometry, and in particular in the hyperbolic and elliptic planes.

For example, the unit sphere \(S^2(1)\) is a model for the two-dimensional elliptic plane. It carries an intrinsic metric that arises by measuring geodesic length. The geodesic circles are the parallels in a geodesic coordinate system.

More precisely, fix a point \(\mathbf z\in S^2(1)\) that we place at the zenith. Associated to that zenith is a geodesic polar coordinate system \((\varphi,\theta)\), \(0\le\varphi\le\pi\), \(0\le\theta< 2\pi\), where **z** is the point \(\varphi=0\). In these coordinates, the geodesic distance from **z** to any other point \(\mathbf x\in S^2(1)\) having coordinates \((\varphi,\theta)\) is the value of \(\varphi\) at **x**. A spherical circle is the set of points a geodesic distance *R* from the zenith point **z**. Equivalently, with a fixed embedding into \(\mathbb R^3\), the spherical circle of radius \(R\le\pi\) centered at **z** is the set of **x** in \(S^2(1)\) such that \(\mathbf x\cdot\mathbf z = \cos R\).

We can also measure the area of the spherical disk enclosed within a spherical circle, using the intrinsic surface area measure on the sphere. The area of the disk of radius *R* is then given by
\(A = \int_0^{2\pi} \int_0^R \sin(\varphi) \, d\varphi\,d\theta = 2\pi(1-\cos R).\)

More generally, if a sphere \(S^2(\rho)\) has radius of curvature \(\rho\), then the area of the disk of radius *R* is given by
\(A = 2\pi\rho^2(1-\cos(R/\rho)).\)
Observe that, as an application of L'Hôpital's rule, this tends to the Euclidean area \(\pi R^2\) in the flat limit \(\rho\to\infty\).

The hyperbolic case is similar, with the area of a disk of intrinsic radius *R* in the (constant curvature \(-1\)) hyperbolic plane given by
\(A = 2\pi(1-\cosh R)\)
where cosh is the hyperbolic cosine. More generally, for the constant curvature \(-k\) hyperbolic plane, the answer is
\(A = 2\pi k^{-2}(1-\cosh(kR)).\)

These identities are important for comparison inequalities in geometry. For example, the area enclosed by a circle of radius *R* in a flat space is always greater than the area of a spherical circle and smaller than a hyperbolic circle, provided all three circles have the same (intrinsic) radius. That is,
\(2\pi(1-\cos R) < \pi R^2 < 2\pi(1-\cosh R)\)
for all \(R>0\). Intuitively, this is because the sphere tends to curve back on itself, yielding circles of smaller area than those in the plane, whilst the hyperbolic plane, when immersed into space, develops fringes that produce additional area. It is more generally true that the area of the circle of a fixed radius *R* is a strictly decreasing function of the curvature.

In all cases, if \(k\) is the curvature (constant, positive or negative), then the isoperimetric inequality for a domain with area *A* and perimeter *L* is
\(L^2\ge 4\pi A - kA^2\)
where equality is achieved precisely for the circle.

## Generalizations
We can stretch a disk to form an ellipse. Because this stretch is a linear transformation of the plane, it has a distortion factor which will change the area but preserve *ratios* of areas. This observation can be used to compute the area of an arbitrary ellipse from the area of a unit circle.

Consider the unit circle circumscribed by a square of side length 2. The transformation sends the circle to an ellipse by stretching or shrinking the horizontal and vertical diameters to the major and minor axes of the ellipse. The square gets sent to a rectangle circumscribing the ellipse. The ratio of the area of the circle to the square is /4, which means the ratio of the ellipse to the rectangle is also /4. Suppose *a* and *b* are the lengths of the major and minor axes of the ellipse. Since the area of the rectangle is *ab*, the area of the ellipse is *ab*/4.

We can also consider analogous measurements in higher dimensions. For example, we may wish to find the volume inside a sphere. When we have a formula for the surface area, we can use the same kind of "onion" approach we used for the disk.
