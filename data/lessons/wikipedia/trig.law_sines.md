> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Law_of_sines) — CC BY-SA 4.0

# Law of sines

In trigonometry, the **law of sines** (sometimes called the **sine formula** or **sine rule**) is a mathematical equation relating the lengths of the sides of any triangle to the sines of its angles. According to the law,

\[
\frac{a}{\sin \alpha} = \frac{b}{\sin \beta} = \frac{c}{\sin \gamma} = 2R,
\]

where *a*, *b*, and *c* are the lengths of the sides of a triangle, and *α*, *β*, and *γ* are the opposite angles (see figure 2), while *R* is the radius of the triangle's circumcircle. When the last part of the equation is not used, the law is sometimes stated using the reciprocals;

\[
\frac{\sin \alpha}{a} = \frac{\sin \beta}{b} = \frac{\sin \gamma}{c}.
\]

The law of sines can be used to compute the remaining sides of a triangle when two angles and a side are known—a technique known as triangulation. It can also be used when two sides and one of the non-enclosed angles are known. In some such cases, the triangle is not uniquely determined by this data (called the *ambiguous case*) and the technique gives two possible values for the enclosed angle.

The law of sines is one of two trigonometric equations commonly applied to find lengths and angles in scalene triangles, with the other being the law of cosines.

The law of sines can be generalized to higher dimensions on surfaces with constant curvature.

## Proof
With the side of length as the base, the triangle's altitude can be computed as *b* sin *γ* or as *c* sin *β*. Equating these two expressions gives
\(\frac{b}{\sin \beta} = \frac{c}{\sin \gamma}\, \)
and similar equations arise by choosing the side of length or the side of length as the base of the triangle. For a proof that these expressions are equal to \(2R\), see Relation to the circumcircle.

## Ambiguous case of triangle solution
When using the law of sines to find a side of a triangle, an ambiguous case occurs when two separate triangles can be constructed from the data provided (i.e., there are two different possible solutions to the triangle). In the case shown below they are triangles *ABC* and *ABC′*.

Given a general triangle, the following conditions would need to be fulfilled for the case to be ambiguous:
* The only information known about the triangle is the angle *α* and the sides *a* and *c*.
* The angle *α* is acute (i.e., *α* < 90°).
* The side *a* is shorter than the side *c* (i.e., *a* < *c*).
* The side *a* is longer than the altitude *h* from angle *β*, where *h* = *c* sin *α* (i.e., *a* > *h*).

If all the above conditions are true, then each of angles *β* and *β′* produces a valid triangle, meaning that both of the following are true:

\[
{\gamma}' = \arcsin\frac{c \sin \alpha}{a} \quad \text{or} \quad {\gamma} = \pi - \arcsin\frac{c \sin \alpha}{a}.
\]

From there we can find the corresponding *β* and *b* or *β′* and *b′* if required, where *b* is the side bounded by vertices *A* and *C* and *b′* is bounded by *A* and *C′*.

## Examples
The following are examples of how to solve a problem using the law of sines.

### Example 1

Given: side *a* = 20, side *c* = 24, and angle *γ* = 40°. Angle *α* is desired.

Using the law of sines, we conclude that

\[
\frac{\sin \alpha}{20} = \frac{\sin (40^\circ)}{24}.
\]

\[
\alpha = \arcsin\left( \frac{20\sin (40^\circ)}{24} \right) \approx 32.39^\circ.
\]

Note that the potential solution *α* = 147.61° is excluded because that would necessarily give *α* + *β* + *γ* > 180°.

### Example 2

If the lengths of two sides of the triangle *a* and *b* are equal to *x*, the third side has length *c*, and the angles opposite the sides of lengths *a*, *b*, and *c* are *α*, *β*, and *γ* respectively then

\[
\begin{aligned}
& \alpha = \beta = \frac{180^\circ-\gamma}{2}= 90^\circ-\frac{\gamma}{2} \\[6pt]
& \sin \alpha = \sin \beta = \sin \left(90^\circ-\frac{\gamma}{2}\right) = \cos \left(\frac{\gamma}{2}\right) \\[6pt]
& \frac{c}{\sin \gamma}=\frac{a}{\sin \alpha}=\frac{x}{\cos \left(\frac{\gamma}{2}\right)} \\[6pt]
& \frac{c \cos \left(\frac{\gamma}{2}\right)}{\sin \gamma} = x
\end{aligned}
\]

## Relation to the circumcircle
In the identity

\[
\frac{a}{\sin \alpha} = \frac{b}{\sin \beta} = \frac{c}{\sin \gamma},
\]

the common value of the three fractions is actually the diameter of the triangle's circumcircle. This result dates back to Ptolemy.

### Proof

As shown in the figure, let there be a circle with inscribed \(\triangle ABC\) and another inscribed \(\triangle ADB\) that passes through the circle's center **O**. The \(\angle AOD\) has a central angle of \(180^\circ\) and thus \(\angle ABD = 90^\circ\), by Thales's theorem. Since \(\triangle ABD\) is a right triangle,

\[
\sin{\delta}= \frac{\text{opposite{\text{hypotenuse= \frac{c}{2R},
\]

where \(R= \frac{d}{2}\) is the radius of the circumscribing circle of the triangle. Angles \({\gamma}\) and \({\delta}\) lie on the same circle and subtend the same chord *c*; thus, by the inscribed angle theorem, \det{V}
