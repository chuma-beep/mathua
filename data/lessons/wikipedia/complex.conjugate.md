> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Complex_conjugate) — CC BY-SA 4.0

# Complex conjugate

In mathematics, the **complex conjugate** of a complex number is the number with an equal real part and an imaginary part equal in magnitude but opposite in sign. That is, if \(a\) and \(b\) are real numbers, then the complex conjugate of \(a + bi\) is .   The complex conjugate of \(z\) is often denoted as \(\overline{z}\) or .

In polar form, if \(r\) and \(\varphi\) are real numbers then the conjugate of \(r e^{i \varphi}\) is -i \varphi}}. This can be shown using Euler's formula.

The product of a complex number and its conjugate is a real number: \(a^2 + b^2\) (or \(r^2\) in polar coordinates).

If a root of a univariate polynomial with real coefficients is complex, then its complex conjugate is also a root.

## Notation
The complex conjugate of a complex number \(z\) is written as \(\overline z\) or . The first notation, a vinculum, avoids confusion with the notation for the conjugate transpose of a matrix, which can be thought of as a generalization of the complex conjugate. The second is preferred in physics, where dagger (†) is used for the conjugate transpose, as well as electrical engineering and computer engineering, where bar notation can be confused for the logical negation ("NOT") Boolean algebra symbol, while the bar notation is more common in pure mathematics.

If a complex number is represented as a \(2 \times 2\) matrix, the notations are identical, and the complex conjugate corresponds to the matrix transpose, which is a flip along the diagonal.

## Properties
The following properties apply for all complex numbers \(z\) and , unless stated otherwise, and can be proved by writing \(z\) and \(w\) in the form .

For any two complex numbers, conjugation is distributive over addition, subtraction, multiplication and division:

\[
\begin{align}
                     \overline{z + w} &= \overline{z} + \overline{w}, \\
                     \overline{z - w} &= \overline{z} - \overline{w}, \\
                        \overline{zw} &= \overline{z} \; \overline{w}, \quad \text{and} \\
  \overline{\left(\frac{z}{w}\right)} &= \frac{\overline{z}}{\overline{w}},\quad \text{if } w \neq 0.
\end{align}
\]


A complex number is equal to its complex conjugate if its imaginary part is zero, that is, if the number is real.  In other words, real numbers are the only fixed points of conjugation.

Conjugation does not change the modulus of a complex number: \(\left| \overline{z} \right| = |z|.\)

Conjugation is an involution, that is, the conjugate of the conjugate of a complex number \(z\) is .  In symbols, {{tmath|1=\overline{ \overline{z} } = z}}.

The product of a complex number with its conjugate is equal to the square of the number's modulus:
\[
z\overline{z} = {\left| z \right|}^2.
\]
 This allows easy computation of the multiplicative inverse of a complex number given in rectangular coordinates:
\[
z^{-1} = \frac{\overline{z}}{{\left| z \right|}^2},\quad \text{ for all } z \neq 0.
\]


Conjugation is commutative under composition with exponentiation to integer powers, with the exponential function, and with the natural logarithm for nonzero arguments:

\[
\overline{z^n} = \left(\overline{z}\right)^n,\quad \text{ for all } n \in \Z
\]


\[
\exp\left(\overline{z}\right) = \overline{\exp(z)}
\]


\[
\ln\left(\overline{z}\right) = \overline{\ln(z)} \text{ if } z \text{ is not zero or a negative real number }
\]


If \(p\) is a polynomial with real coefficients and \(p(z) = 0,\) then \(p\left(\overline{z}\right) = 0\) as well. Thus, non-real roots of real polynomials occur in complex conjugate pairs (*see* Complex conjugate root theorem).

In general, if \(\varphi\) is a holomorphic function whose restriction to the real numbers is real-valued, and \(\varphi(z)\) and \(\varphi(\overline{z})\) are defined, then

\[
\varphi\left(\overline{z}\right) = \overline{\varphi(z)}.\,\!
\]


The map \(\sigma(z) = \overline{z}\) from \(\Complex\) to \(\Complex\) is a homeomorphism (where the topology on \(\Complex\) is taken to be the standard topology) and antilinear, if one considers \(\Complex\) as a complex vector space over itself. Even though it appears to be a well-behaved function, it is not holomorphic; it reverses orientation whereas holomorphic functions locally preserve orientation. It is bijective and compatible with the arithmetical operations, and hence is a field automorphism. As it keeps the real numbers fixed, it is an element of the Galois group of the field extension . This Galois group has only two elements: \(\sigma\) and the identity on . Thus the only two field automorphisms of \(\Complex\) that leave the real numbers fixed are the identity map and complex conjugation.

## Use as a variable
Once a complex number \(z = x + yi\) or \(z = re^{i\theta}\) is given, its conjugate is sufficient to reproduce the parts of the \(z\)-variable:
* Real part: \(x = \operatorname{Re}(z) = \dfrac{z + \overline{z}}{2}\)
* Imaginary part: \(y = \operatorname{Im}(z) = \dfrac{z - \overline{z}}{2i}\)
* Modulus (or absolute value): \(r= \left| z \right| = \sqrt{z\overline{z}}\)
* Argument: \(e^{i\theta} = e^{i\arg z} = \sqrt{\dfrac{z}{\overline z}},\) so \(\theta = \arg z = \dfrac{1}{i} \ln\sqrt{\frac{z}{\overline{z}}} = \dfrac{\ln z - \ln \overline{z}}{2i}\)

Furthermore, \(\overline{z}\) can be used to specify lines in the plane: the set

\[
\left\{z : z \overline{r} + \overline{z} r = 0 \right\}
\]

is a line through the origin and perpendicular to \({r},\) since the real part of \(z\cdot\overline{r}\) is zero only when the cosine of the angle between \(z\) and \({r}\) is zero. Similarly, for a fixed complex unit \(u = e^{i b},\) the equation

\[
\frac{z - z_0}{\overline{z} - \overline{z_0}} = u^2
\]

determines the line through \(z_0\) parallel to the line through 0 and .

These uses of the conjugate of \(z\) as a variable are illustrated in Frank Morley's book *Inversive Geometry* (1933), written with his son Frank Vigor Morley.

## Generalizations
The other planar real unital algebras, dual numbers, and split-complex numbers are also analyzed using complex conjugation.

For matrices of complex numbers, \(\overline{\mathbf{AB}} = \left(\overline{\mathbf{A}}\right) \left(\overline{\mathbf{B}}\right),\) where \(\overline{\mathbf{A}}\) represents the element-by-element conjugation of \(\mathbf{A}.\) Contrast this to the property \(\left(\mathbf{AB}\right)^*=\mathbf{B}^* \mathbf{A}^*,\) where \(\mathbf{A}^*\) represents the conjugate transpose of \(\mathbf{A}.\)

Taking the conjugate transpose (or adjoint) of complex matrices generalizes complex conjugation. Even more general is the concept of adjoint operator for operators on (possibly infinite-dimensional) complex Hilbert spaces. All this is subsumed by the *-operations of C*-algebras.

One may also define a conjugation for quaternions and split-quaternions: the conjugate of \(a + bi + cj + dk\) is .

All these generalizations are multiplicative only if the factors are reversed:

\[
{\left(zw\right)}^* = w^* z^*.
\]


Since the multiplication of planar real algebras is commutative, this reversal is not needed there.

There is also an abstract notion of conjugation for vector spaces \(V\) over the complex numbers. In this context, any antilinear map \(\varphi: V \to V\) that satisfies

# \(\varphi^2 = \operatorname{id}_V\,,\) where \(\varphi^2 = \varphi \circ \varphi\) and \(\operatorname{id}_V\) is the identity map on ,
# \(\varphi(zv) = \overline{z} \varphi(v)\) for all \(v \in V, z \in \Complex,\) and
# \(\varphi\left(v_1 + v_2\right) = \varphi\left(v_1\right) + \varphi\left(v_2\right)\,\) for all \(v_1, v_2 \in V,\)

is called a , or a real structure. As the involution \(\varphi\) is antilinear, it cannot be the identity map on .

Of course, \(\varphi\) is a \(\R\)-linear transformation of \(V,\) if one notes that every complex space \(V\) has a real form obtained by taking the same vectors as in the original space and restricting the scalars to be real. The above properties actually define a real structure on the complex vector space .

One example of this notion is the conjugate transpose operation of complex matrices defined above. However, on generic complex vector spaces, there is no  notion of complex conjugation.

## See also
*
*
*
*
*
*
*
*

## References


## Footnotes


## Bibliography
* Budinich, P. and Trautman, A. *The Spinorial Chessboard*. Springer-Verlag, 1988. . (antilinear maps are discussed in section 3.3).

