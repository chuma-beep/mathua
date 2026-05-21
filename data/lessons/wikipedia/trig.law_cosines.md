> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Law_of_cosines) — CC BY-SA 4.0

# Law of cosines

In trigonometry, the **law of cosines** (also known as the **cosine formula** or **cosine rule**) relates the lengths of the sides of a triangle to the cosine of one of its angles. For a triangle with sides , and , opposite respective angles , and (see Fig. 1), the law of cosines states:

\(\begin{align}
c^2 &= a^2 + b^2 - 2ab\cos\gamma, \\[3mu]
a^2 &= b^2+c^2-2bc\cos\alpha, \\[3mu]
b^2 &= a^2+c^2-2ac\cos\beta.
\end{align}\)

The law of cosines generalizes the Pythagorean theorem, which holds only for right triangles: if is a right angle then , and the law of cosines reduces to. The law of cosines is useful for solving a triangle when all three sides or two sides and their included angle are given.

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

*the third side of a triangle if two sides and an angle opposite to one of them is known (this side can also be found by two applications of the law of sines):, and angle , can be found using the law of sines, leaving up to two possibilities for angle. Either choice determines because the three interior angles sum to a straight angle. Finally can be found from , and by another application of the law of sines.
\[
a=b\cos\gamma \pm \sqrt{c^2 -b^2\sin^2\gamma}\,.
\]

These formulas produce high round-off errors in floating point calculations if the triangle is very acute, i.e., if *c* is small relative to *a* and *b* or *γ* is small compared to 1. It is even possible to obtain a result slightly greater than one for the cosine of an angle.

The third formula shown is the result of solving for *a* in the quadratic equation 1=*a*\(^{2}\) − 2*ab* cos *γ* + *b*\(^{2}\) − *c*\(^{2}\) = 0. This equation can have 2, 1, or 0 positive solutions corresponding to the number of possible triangles given the data. It will have two positive solutions if *b* sin *γ* < *c* < *b*, only one positive solution if 1=*c* = *b* sin *γ*, and no solution if *c* < *b* sin *γ*. These different cases are also explained by the side-side-angle congruence ambiguity.

## History
Book II of Euclid's *Elements*, compiled c. 300 BC from material up to a century or two older, contains a geometric theorem corresponding to the law of cosines but expressed in the contemporary language of rectangle areas; Hellenistic trigonometry developed later, and sine and cosine per se first appeared centuries afterward in India.

The cases of obtuse triangles and acute triangles (corresponding to the two cases of negative or positive cosine) are treated separately, in Propositions II.12 and II.13:

 

Proposition 13 contains an analogous statement for acute triangles. In his (now-lost and only preserved through fragmentary quotations) commentary, Heron of Alexandria provided proofs of the converses of both II.12 and II.13.

Using notation as in Fig. 2, Euclid's statement of proposition II.12 can be represented more concisely (though anachronistically) by the formula

\(AB^2 = CA^2 + CB^2 + 2 (CA)(CH).\)

To transform this into the familiar expression for the law of cosines, substitute , and \(CH = a \cos(\pi - \gamma)\\)\cos \gamma\right
