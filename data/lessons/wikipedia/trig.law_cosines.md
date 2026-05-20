> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Law_of_cosines) — CC BY-SA 4.0

# Law of cosines

In trigonometry, the **law of cosines** (also known as the **cosine formula** or **cosine rule**) relates the lengths of the sides of a triangle to the cosine of one of its angles. For a triangle with sides , , and , opposite respective angles , , and  (see Fig. 1), the law of cosines states:

\(\begin{align}
c^2 &= a^2 + b^2 - 2ab\cos\gamma, \\[3mu]
a^2 &= b^2+c^2-2bc\cos\alpha, \\[3mu]
b^2 &= a^2+c^2-2ac\cos\beta.
\end{align}\)

The law of cosines generalizes the Pythagorean theorem, which holds only for right triangles: if  is a right angle then , and the law of cosines reduces to .

The law of cosines is useful for solving a triangle when all three sides or two sides and their included angle are given.

## Use in solving triangles


The theorem is used in solution of triangles, i.e., to find (see Figure 3):
*the third side of a triangle if two sides and the angle between them is known:
\[
c = \sqrt{a^2+b^2-2ab\cos\gamma}\,;
\]

*the angles of a triangle if the three sides are known:
\[
\gamma = \arccos\left(\frac{a^2+b^2-c^2}{2ab}\right)\,;
\]

*the third side of a triangle if two sides and an angle opposite to one of them is known (this side can also be found by two applications of the law of sines):,  and angle ,  can be found using the law of sines, leaving up to two possibilities for angle . Either choice determines  because the three interior angles sum to a straight angle. Finally  can be found from , , and  by another application of the law of sines.}}
\[
a=b\cos\gamma \pm \sqrt{c^2 -b^2\sin^2\gamma}\,.
\]


These formulas produce high round-off errors in floating point calculations if the triangle is very acute, i.e., if *c* is small relative to *a* and *b* or *γ* is small compared to 1. It is even possible to obtain a result slightly greater than one for the cosine of an angle.

The third formula shown is the result of solving for *a* in the quadratic equation 1=*a*<sup>2</sup> − 2*ab* cos *γ* + *b*<sup>2</sup> − *c*<sup>2</sup> = 0. This equation can have 2, 1, or 0 positive solutions corresponding to the number of possible triangles given the data. It will have two positive solutions if *b* sin *γ* < *c* < *b*, only one positive solution if 1=*c* = *b* sin *γ*, and no solution if *c* < *b* sin *γ*. These different cases are also explained by the side-side-angle congruence ambiguity.

## History
Book II of Euclid's *Elements*, compiled c. 300 BC from material up to a century or two older, contains a geometric theorem corresponding to the law of cosines but expressed in the contemporary language of rectangle areas; Hellenistic trigonometry developed later, and sine and cosine per se first appeared centuries afterward in India.

The cases of obtuse triangles and acute triangles (corresponding to the two cases of negative or positive cosine) are treated separately, in Propositions II.12 and II.13:

 }}

Proposition 13 contains an analogous statement for acute triangles. In his (now-lost and only preserved through fragmentary quotations) commentary, Heron of Alexandria provided proofs of the converses of both II.12 and II.13.


Using notation as in Fig. 2, Euclid's statement of proposition II.12 can be represented more concisely (though anachronistically) by the formula

\(AB^2 = CA^2 + CB^2 + 2 (CA)(CH).\)

To transform this into the familiar expression for the law of cosines, substitute , , , and \(CH = a \cos(\pi - \gamma)\\){{nobr|\(\!{} = -a \cos \gamma\).}}

Proposition II.13 was not used in Euclid's time for the solution of triangles, but later it was used that way in the course of solving astronomical problems by al-Bīrūnī (11th century) and Johannes de Muris (14th century). Something equivalent to the spherical law of cosines was used (but not stated in general) by al-Khwārizmī (9th century), al-Battānī (9th century), and Nīlakaṇṭha (15th century).

The 13th century Persian mathematician Naṣīr al-Dīn al-Ṭūsī, in his  (*Book on the Complete Quadrilateral*, c. 1250), systematically described how to solve triangles from various combinations of given data. Given two sides and their included angle in a scalene triangle, he proposed finding the third side by dropping a perpendicular from the vertex of one of the unknown angles to the opposite base, reducing the problem to finding the legs of one right triangle from a known angle and hypotenuse using the law of sines and then finding the hypotenuse of another right triangle from two known sides by the Pythagorean theorem.

About two centuries later, another Persian mathematician, Jamshīd al-Kāshī, who computed the most accurate trigonometric tables of his era, also described the solution of triangles from various combinations of given data in his  (*Key of Arithmetic*, 1427), and repeated essentially al-Ṭūsī's method, now consolidated into one formula and including more explicit details, as follows:


,translation by Nuh Aydin, Lakhdar Hammoudi, and Ghada Bakbouk }}

Using modern algebraic notation and conventions this might be written

\(c = \sqrt{(b - a \cos \gamma)^2 + (a \sin \gamma)^2}\)

when  is acute or

\(c = \sqrt{\left(b + a \left|\cos \gamma\right| \right)^2 + \left(a \sin \gamma\right)^2}\)

when  is obtuse. (When  is obtuse, the modern convention is that  is negative and \(\cos(\pi-\gamma) = -\cos\gamma\) is positive; historically sines and cosines were considered to be line segments with non-negative lengths.) By squaring both sides, expanding the squared binomial, and then applying the Pythagorean trigonometric identity , we obtain the familiar law of cosines:

\(\begin{align}
c^2
&= b^2 - 2ba\cos \gamma + a^2\cos^2 \gamma + a^2\sin^2\gamma \\[5mu]
&= a^2 + b^2 - 2ab\cos \gamma.
\end{align}\)

In France, the law of cosines is sometimes referred to as the *théorème d'Al-Kashi*.

The same method used by al-Ṭūsī appeared in Europe as early as the 15th century, in Regiomontanus's *De triangulis omnimodis* (*On Triangles of All Kinds*, 1464), a comprehensive survey of plane and spherical trigonometry known at the time.

The theorem was first written using algebraic notation by François Viète in the 16th century. At the beginning of the 19th century, modern algebraic notation allowed the law of cosines to be written in its current symbolic form.

## Proofs
### Using the Pythagorean theorem


#### Case of an obtuse angle
Euclid proved this theorem by applying the Pythagorean theorem to each of the two right triangles in Fig. 2 (*AHB* and *CHB*). Using *a* to denote the line segment *CB*, *b* to denote the line segment *AC*, *c* to denote the line segment *AB*, *d* to denote the line segment *CH* and *h* for the height *BH*, triangle *AHB* gives us
\(c^2 = (b+d)^2 + h^2,\)

and triangle *CHB* gives
\(d^2 + h^2 = a^2.\)

Expanding the first equation gives
\(c^2 = b^2 + 2bd + d^2 +h^2.\)

Substituting the second equation into this, the following can be obtained:
\(c^2 = a^2 + b^2 + 2bd.\)

This is Euclid's Proposition 12 from Book 2 of the *Elements*. To transform it into the modern form of the law of cosines, note that
\(d = a\cos(\pi-\gamma)= -a\cos\gamma.\)

#### Case of an acute angle
Euclid's proof of his Proposition 13 proceeds along the same lines as his proof of Proposition 12: he applies the Pythagorean theorem to both right triangles formed by dropping the perpendicular onto one of the sides enclosing the angle *γ* and uses the square of a difference to simplify.

#### Another proof in the acute case


Using more trigonometry, the law of cosines can be deduced by using the Pythagorean theorem only once. In fact, by using the right triangle on the left hand side of Fig. 6 it can be shown that:

\(\begin{align}
 c^2 &= (b-a\cos\gamma)^2 + (a\sin\gamma)^2 \\
     &= b^2 - 2ab\cos\gamma + a^2\cos^2\gamma + a^2\sin^2\gamma \\
     &= b^2 + a^2 - 2ab\cos\gamma,
\end{align}\)

using the trigonometric identity .

This proof needs a slight modification if *b* < *a* cos(*γ*). In this case, the right triangle to which the Pythagorean theorem is applied moves *outside* the triangle *ABC*. The only effect this has on the calculation is that the quantity *b* − *a* cos(*γ*) is replaced by *a* cos(*γ*) − *b*. As this quantity enters the calculation only through its square, the rest of the proof is unaffected. However, this problem only occurs when *β* is obtuse, and may be avoided by reflecting the triangle about the bisector of *γ*.

Referring to Fig. 6 it is worth noting that if the angle opposite side *a* is *α* then:
\(\tan\alpha = \frac{a\sin\gamma}{b-a\cos\gamma}.\)

This is useful for direct calculation of a second angle when two sides and an included angle are given.

### From three altitudes

The altitude through vertex  is a segment perpendicular to side . The distance from the foot of the altitude to vertex  plus the distance from the foot of the altitude to vertex  is equal to the length of side  (see Fig. 5). Each of these distances can be written as one of the other sides multiplied by the cosine of the adjacent angle,
\(c=a\cos\beta+b\cos\alpha.\)

(This is still true if *α* or *β* is obtuse, in which case the perpendicular falls outside the triangle.)  Multiplying both sides by *c* yields
\(c^2 = ac\cos\beta + bc\cos\alpha.\)

The same steps work just as well when treating either of the other sides as the base of the triangle:
\(\begin{align}
a^2 &= ac\cos\beta + ab\cos\gamma, \\[3mu]
b^2 &= bc\cos\alpha + ab\cos\gamma.
\end{align}\)

Taking the equation for  and subtracting the equations for  and ,
\(\begin{align}

c^2 - a^2 - b^2 &= {\color{BlueGreen}\cancel{\color{Black}ac\cos\beta}} + {\color{Peach}\cancel{\color{Black}bc\cos\alpha}} - {\color{BlueGreen}\cancel{\color{Black}ac\cos\beta}} - {\color{Peach}\cancel{\color{Black}bc\cos\alpha}} - 2ab\cos\gamma \\
c^2 &= a^2 + b^2- 2ab\cos\gamma.
\end{align}\)

This proof is independent of the Pythagorean theorem, insofar as it is based only on the right-triangle definition of cosine and obtains squared side lengths algebraically. Other proofs typically invoke the Pythagorean theorem explicitly, and are more geometric, treating *a* cos *γ* as a label for the length of a certain line segment.

Unlike many proofs, this one handles the cases of obtuse and acute angles *γ* in a unified fashion.

### Cartesian coordinates


Consider a triangle with sides of length *a*, *b*, *c*, where *θ* is the measurement of the angle opposite the side of length *c*. This triangle can be placed on the Cartesian coordinate system with side *a* aligned along the *x* axis and angle *θ* placed at the origin, by plotting the components of the 3 points of the triangle as shown in Fig. 4:
\(A = (b \cos\theta, b \sin\theta), B = (a, 0), \text{ and } C = (0, 0).\)

By the distance formula,

\(c = \sqrt{(a - b \cos\theta)^2 + (0 - b \sin\theta)^2}.\)

Squaring both sides and simplifying
\(\begin{align}
c^2 &=  (a - b \cos\theta)^2 + (- b \sin\theta)^2 \\
    &=  a^2 - 2 a b \cos\theta+ b^2 \cos^2 \theta+ b^2 \sin^2 \theta\\
    &=  a^2 + b^2 (\sin^2 \theta+ \cos^2 \theta) - 2 a b \cos\theta\\
    &=  a^2 + b^2  - 2 a b \cos\theta.
\end{align}\)

An advantage of this proof is that it does not require the consideration of separate cases depending on whether the angle  is acute, right, or obtuse. However, the cases treated separately in *Elements* II.12–13 and later by al-Ṭūsī, al-Kāshī, and others could themselves be combined by using concepts of signed lengths and areas and a concept of signed cosine, without needing a full Cartesian coordinate system.

### Using Ptolemy's theorem

Referring to the diagram, triangle *ABC* with sides *AB* = *c*, *BC* = *a* and *AC* = *b* is drawn inside its circumcircle as shown. Triangle *ABD* is constructed congruent to triangle *ABC* with *AD* = *BC* and *BD* = *AC*. Perpendiculars from *D* and *C* meet base *AB* at *E* and *F* respectively. Then:
\(\begin{align}
& BF=AE=BC\cos\hat{B}=a\cos\hat{B} \\
\Rightarrow \ & DC=EF=AB-2BF=c-2a\cos\hat{B}.
\end{align}\)

Now the law of cosines is rendered by a straightforward application of Ptolemy's theorem to cyclic quadrilateral *ABCD*:
\(\begin{align}
& AD \times BC + AB \times DC = AC \times BD \\
\Rightarrow \ & a^2 + c(c-2a\cos\hat{B})=b^2 \\
\Rightarrow \ & a^2+c^2-2ac \cos\hat{B}=b^2.
\end{align}\)

Plainly if angle *B* is right, then *ABCD* is a rectangle and application of Ptolemy's theorem yields the Pythagorean theorem:
\(a^2+c^2=b^2.\)

### By comparing areas


One can also prove the law of cosines by calculating areas. The change of sign as the angle *γ* becomes obtuse makes a case distinction necessary.

Recall that
**a*<sup>2</sup>, *b*<sup>2</sup>, and *c*<sup>2</sup> are the areas of the squares with sides *a*, *b*, and *c*, respectively;
*if *γ* is acute, then *ab* cos *γ* is the area of the parallelogram with sides *a* and *b* forming an angle of *γ′*  − *γ*}};
*if *γ* is obtuse, and so cos *γ* is negative, then −*ab* cos *γ* is the area of the parallelogram with sides *a* and *b* forming an angle of *γ′* }}.

**Acute case.** Figure 7a shows a heptagon cut into smaller pieces (in two different ways) to yield a proof of the law of cosines. The various pieces are
*in pink, the areas *a*<sup>2</sup>, *b*<sup>2</sup> on the left and the areas 2*ab* cos *γ* and *c*<sup>2</sup> on the right;
*in blue, the triangle *ABC*, on the left and on the right;
*in grey, auxiliary triangles, all congruent to *ABC*, an equal number (namely 2) both on the left and on the right.

The equality of areas on the left and on the right gives
\(a^2 + b^2 = c^2 + 2ab\cos\gamma.\)

**Obtuse case.** Figure 7b cuts a hexagon in two different ways into smaller pieces, yielding a proof of the law of cosines in the case that the angle *γ* is obtuse. We have
*in pink, the areas *a*<sup>2</sup>, *b*<sup>2</sup>, and −2*ab* cos *γ* on the left and *c*<sup>2</sup> on the right;
*in blue, the triangle *ABC* twice, on the left, as well as on the right.

The equality of areas on the left and on the right gives
\(a^2 + b^2 - 2ab\cos(\gamma) = c^2.\)

The rigorous proof will have to include proofs that various shapes are congruent and therefore have equal area. This will use the theory of congruent triangles.


### Using circle geometry


Using the geometry of the circle, it is possible to give a more geometric proof than using the Pythagorean theorem alone. Algebraic manipulations (in particular the binomial theorem) are avoided.

**Case of acute angle *γ*, where *a* > 2*b* cos *γ*.** Drop the perpendicular from *A* onto *a* = *BC*, creating a line segment of length *b* cos *γ*. Duplicate the *right triangle* to form the isosceles triangle *ACP*. Construct the circle with center *A* and radius *b*, and its tangent *h*  through *B*. The tangent *h* forms a right angle with the radius *b* (Euclid's *Elements*: Book 3, Proposition 18; or see here), so the yellow triangle in Figure 8 is right. Apply  the Pythagorean theorem to obtain
\(c^2 = b^2 + h^2.\)

Then use the *tangent secant theorem* (Euclid's *Elements*: Book 3, Proposition 36), which says that the square on the tangent through a point *B* outside the circle is equal to the product of the two lines segments (from *B*) created by any secant of the circle through *B*. In the present case: *BH*<sup>2</sup> , or
\(h^2 = a(a - 2b\cos\gamma).\)

Substituting into the previous equation gives the law of cosines:
\(c^2 = b^2 + a(a - 2b\cos\gamma).\)

Note that *h*<sup>2</sup> is the power of the point *B* with respect to the circle. The use of the Pythagorean theorem and the tangent secant theorem can be replaced by a single application of the power of a point theorem.

**Case of acute angle *γ*, where *a* < 2*b* cos *γ*.** Drop the perpendicular from *A* onto *a* = *BC*, creating a line segment of length *b* cos *γ*. Duplicate the right triangle to form the isosceles triangle *ACP*. Construct the circle with center *A* and radius *b*, and a chord through *B* perpendicular to *c*  half of which is *h*  Apply  the Pythagorean theorem to obtain
\(b^2 = c^2 + h^2.\)

Now use the *chord theorem* (Euclid's *Elements*: Book 3, Proposition 35), which says that if two chords intersect, the product of the two line segments obtained on one chord is equal to the product of the two line segments obtained on the other chord. In the present case: *BH*<sup>2</sup>  or
\(h^2 = a(2b\cos\gamma - a).\)

Substituting into the previous equation gives the law of cosines:
\(b^2 = c^2 + a(2b\cos\gamma - a).\)

Note that the power of the point *B* with respect to the circle has the negative value −*h*<sup>2</sup>.

**Case of obtuse angle *γ*.** This proof uses the power of a point theorem directly, without the auxiliary triangles obtained by constructing a tangent or a chord. Construct a circle with center *B* and radius *a* (see Figure 9), which intersects the secant through *A* and *C* in *C* and *K*. The power of the point *A* with respect to the circle is equal to both *AB*<sup>2</sup> − *BC*<sup>2</sup> and *AC*·*AK*. Therefore,
\(\begin{align}
c^2 - a^2 & {} = b(b + 2a\cos(\pi - \gamma)) \\
& {} = b(b - 2a\cos\gamma),
\end{align}\)

which is the law of cosines.

Using algebraic measures for line segments (allowing negative numbers as lengths of segments) the case of obtuse angle (*CK* > 0) and acute angle (*CK* < 0) can be treated simultaneously.

### Using the law of sines
The law of cosines can be proven algebraically from the law of sines and a few standard trigonometric identities. To start, three angles of a triangle sum to a straight angle (\(\alpha + \beta + \gamma = \pi\) radians). Thus by the angle sum identities for sine and cosine,

\(\begin{alignat}{3}
\sin \gamma &= \phantom{-}\sin(\pi - \gamma)
&&= \phantom{-}\sin(\alpha + \beta)
&&= \sin\alpha\,\cos\beta + \cos\alpha\,\sin\beta, \\[5mu]
\cos \gamma &= -\cos(\pi - \gamma)
&&= -\cos(\alpha + \beta)
&&= \sin\alpha\,\sin\beta - \cos\alpha\,\cos\beta.
\end{alignat}\)

Squaring the first of these identities, then substituting \(\cos\alpha\,\cos\beta = {}\)\(\sin\alpha\,\sin\beta - \cos \gamma\) from the second, and finally replacing \(\cos^2 \alpha + \sin^2 \alpha = {}\)\(\cos^2 \beta + \sin^2 \beta = 1,\) the Pythagorean trigonometric identity, we have:

\(\begin{align}
\sin^2 \gamma
&= (\sin\alpha\,\cos\beta + \cos\alpha\,\sin\beta)^2 \\[3mu]
&= \sin^2\alpha\,\cos^2\beta + 2\sin\alpha\,\sin\beta\,\cos\alpha\,\cos\beta + \cos^2\alpha\,\sin^2\beta \\[3mu]
&= \sin^2\alpha\,\cos^2\beta + 2\sin\alpha\,\sin\beta(\sin\alpha\,\sin\beta - \cos \gamma) + \cos^2\alpha\,\sin^2\beta \\[3mu]
&= \sin^2\alpha(\cos^2\beta + \sin^2\beta) + \sin^2\beta(\cos^2\alpha + \sin^2\alpha) - 2\sin\alpha\,\sin\beta\,\cos \gamma \\[3mu]
&= \sin^2\alpha + \sin^2\beta - 2\sin\alpha\,\sin\beta\,\cos \gamma.
\end{align}\)

The law of sines holds that
\(\frac{a}{\sin \alpha\vphantom{\beta}} = \frac{b}{\sin \beta} = \frac{c}{\sin \gamma\vphantom{\beta}} = k,\)

so to prove the law of cosines, we multiply both sides of our previous identity by :

\(\begin{align}
\sin^2 \gamma \frac{c^2}{\sin^2 \gamma}
&= \sin^2 \alpha \frac{a^2}{\sin^2 \alpha} + \sin^2 \beta \frac{b^2}{\sin^2 \beta} - 2\sin\alpha\,\sin\beta\,\cos\gamma\frac{ab}{\sin\alpha\,\sin\beta\vphantom{\sin^2}} \\[10mu]
c^2 &= a^2 + b^2 - 2ab\cos\gamma.
\end{align}\)

This concludes the proof.

### Using vectors


Denote

\(\overrightarrow{CB}=\vec{a},\ \overrightarrow{CA}=\vec{b},\ \overrightarrow{AB}=\vec{c}\)

Therefore,
\(\vec{c} = \vec{a}-\vec{b}\)

Taking the dot product of each side with itself:
\(\begin{align}
\vec{c}\cdot\vec{c} &= (\vec{a}-\vec{b})\cdot(\vec{a}-\vec{b}) \\
\Vert\vec{c}\Vert^2 &= \Vert\vec{a}\Vert^2 + \Vert\vec{b}\Vert^2 - 2\,\vec{a}\cdot\vec{b}
\end{align}\)

Using the identity

\(\vec{u}\cdot\vec{v} = \Vert\vec{u}\Vert\,\Vert\vec{v}\Vert \cos\angle(\vec{u},\ \vec{v})\)

leads to

\(\Vert\vec{c}\Vert^2 = \Vert\vec{a}\Vert^2 + {\Vert\vec{b}\Vert}^2 - 2\,\Vert\vec{a}\Vert\!\;\Vert\vec{b}\Vert \cos\angle(\vec{a}, \ \vec{b})\)

The result follows.

## Isosceles case
When *a* , i.e., when the triangle is isosceles with the two sides incident to the angle *γ* equal, the law of cosines simplifies significantly. Namely, because *a*<sup>2</sup> + *b*<sup>2</sup>  2*ab*}}, the law of cosines becomes
\(\cos\gamma = 1 - \frac{c^2}{2a^2}\)

or
\(c^2 = 2a^2 (1 - \cos\gamma).\)

## Analogue for tetrahedra
Given an arbitrary tetrahedron whose four faces have areas , , , and , with dihedral angle {{tmath|\varphi_{ab} }} between faces  and , etc., a higher-dimensional analogue of the law of cosines is:
\(A^2 = B^2 + C^2 + D^2 - 2\left(BC \cos \varphi_{bc} + CD \cos\varphi_{cd} + DB\cos \varphi_{db}\right).\)

## Version suited to small angles
When the angle, *γ*, is small and the adjacent sides, *a* and *b*, are of similar length, the right hand side of the standard form of the law of cosines is subject to catastrophic cancellation in numerical approximations. In situations where this is an important concern, a mathematically equivalent version of the law of cosines, similar to the haversine formula, can prove useful:
\(\begin{align}
c^2 &= (a - b)^2 + 4ab\sin^2\left(\frac{\gamma}{2}\right) \\
& = (a - b)^2 + 4ab\operatorname{haversin}(\gamma).
\end{align}\)

In the limit of an infinitesimal angle, the law of cosines degenerates into the circular arc length formula, *c* .

## In non-Euclidean geometry


As in Euclidean geometry, one can use the law of cosines to determine the angles *A*, *B*, *C* from the knowledge of the sides *a*, *b*, *c*. In contrast to Euclidean geometry, the reverse is also possible in both non-Euclidean models: the angles *A*, *B*, *C* determine the sides *a*, *b*, *c*.

A triangle is defined by three points **u**, **v**, and **w** on the unit sphere, and the arcs of great circles connecting those points. If these great circles make angles *A*, *B*, and *C* with opposite sides *a*, *b*, *c* then the spherical law of cosines asserts that all of the following relationships hold:

\(\begin{align}
\cos a &= \cos b\cos c + \sin b\sin c\cos A\\
\cos A &= -\cos B\cos C + \sin B\sin C\cos a\\
\cos a &= \frac{\cos A + \cos B\cos C} {\sin B \sin C}.
\end{align}\)


In hyperbolic geometry, a pair of equations are collectively known as the hyperbolic law of cosines. The first is
\(\cosh a = \cosh b\cosh c - \sinh b \sinh c \cos A\)

where sinh and cosh are the hyperbolic sine and cosine, and the second is
\(\cos A = -\cos B \cos C + \sin B\sin C\cosh a.\)

The length of the sides can be computed by:

\(\cosh a = \frac{\cos A + \cos B\cos C} {\sin B \sin C}.\)

## Polyhedra
The law of cosines can be generalized to all polyhedra by considering any polyhedron with vector sides and invoking the divergence theorem.

## See also
*Half-side formula
*Law of sines
*Law of tangents
*Law of cotangents
*List of trigonometric identities
*Mollweide's formula


## References


## External links

*
* Several derivations of the Cosine Law, including Euclid's at cut-the-knot
* Interactive applet of Law of Cosines

