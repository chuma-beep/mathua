> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Absolute_value) — CC BY-SA 4.0

# Absolute value

In mathematics, the **absolute value** or **modulus** of a real number \(x\), denoted \(|x|\), is the (non-negative) magnitude of \(x\) measured without regard to its sign. Namely, \(|x|=x\) if \(x\) is a positive number, and \(|x|=-x\) if \(x\) is negative (in which case negating \(x\) makes \(-x\) positive), and \(|0|=0\). For example, the absolute value of 3 is 3, and the absolute value of −3 is also 3. The absolute value of a number may be thought of as its distance from zero.

Generalisations of the absolute value for real numbers occur in a wide variety of mathematical settings. For example, an absolute value is also defined for the complex numbers, the quaternions, ordered rings, fields and vector spaces. The absolute value is closely related to the notions of magnitude, distance, and norm in various mathematical and physical contexts.

## Terminology and notation
In 1806, Jean-Robert Argand introduced the term *module*, meaning *unit of measure* in French, specifically for the *complex* absolute value, and it was borrowed into English in 1866 as the Latin equivalent *modulus*. The term *absolute value* has been used in this sense from at least 1806 in French and 1857 in English.}} The notation |, with a vertical bar on each side, was introduced by Karl Weierstrass in 1841. Other names for *absolute value* include *numerical value* and *magnitude*. The absolute value of \(x\) has also been denoted \(\operatorname{abs} x\) in some mathematical publications, and in spreadsheets, programming languages, and computational software packages, the absolute value of \(x\) is generally represented by abs(*x*), or a similar expression, as it has been since the earliest days of high-level programming languages.

The vertical bar notation also appears in a number of other mathematical contexts: for example, when applied to a set, it denotes its cardinality; when applied to a matrix, it denotes its determinant. Vertical bars denote the absolute value only for algebraic objects for which the notion of an absolute value is defined, notably an element of a normed division algebra, for example, a real number, a complex number, or a quaternion. A closely related but distinct notation is the use of vertical bars for either the Euclidean norm or sup norm of a vector in \(\R^n\), although double vertical bars with subscripts (\(\|\cdot\|_2\) and \(\|\cdot\|_\infty\), respectively) are a more common and less ambiguous notation.

## Definition and properties
### Real numbers
For any real number \(x\), the absolute value or modulus of \(x\) is denoted by \(|x|\), with a vertical bar on each side of the quantity, and is defined as
\(|x| =
   \begin{cases}
     x, & \text{if }  x \geq 0 \\
     -x, & \text{if } x < 0.
   \end{cases}\)

The absolute value of \(x\) is thus always either a positive number or zero, but never negative. When \(x\) itself is negative (\(x < 0\)), then its absolute value is necessarily positive (\(|x|=-x>0\)).

From an analytic geometry point of view, the absolute value of a real number is that number's distance from zero along the real number line, and more generally, the absolute value of the difference of two real numbers (their absolute difference) is the distance between them. The notion of an abstract distance function in mathematics can be seen to be a generalisation of the absolute value of the difference. See  below.

Since the square root symbol represents the unique *positive* square root, when applied to a positive number, it follows that
\(|x| = \sqrt{x^2}.\)
This is equivalent to the definition above, and may be used as an alternative definition of the absolute value of real numbers.

The absolute value has the following four fundamental properties (\(a\), \(b\) are real numbers), that are used for generalization of this notion to other domains:

{| style="margin-left:1.6em"
|-
| style="width: 250px" |\(|a| \ge 0\)
| Non-negativity
|-
|\(|a| = 0 \iff a = 0\)
|Positive-definiteness
|-
|\(|ab| = \left|a\right| \left|b\right|\)
|Multiplicativity
|-
|\(|a+b| \le |a| + |b|\)
| Subadditivity, specifically the triangle inequality
|}

Non-negativity, positive definiteness, and multiplicativity are readily apparent from the definition. To see that subadditivity holds, first note that \(|a+b|=s(a+b)\) where \(s=\pm 1\), with its sign chosen to make the result positive. Now, since \(-1 \cdot x \le |x|\) and \(+1 \cdot x \le |x|\), it follows that, whichever of \(\pm1\) is the value of \(s\), one has \(s \cdot x\leq |x|\) for all real \(x\). Consequently, \(|a+b|=s \cdot (a+b) = s \cdot a + s \cdot b \leq |a| + |b|\), as desired.

Some additional useful properties are given below. These are either immediate consequences of the definition or implied by the four fundamental properties above.

{| style="margin-left:1.6em"
|-
| style="width:250px" |\(\bigl| \left|a\right| \bigr| = |a|\)
|Idempotence (the absolute value of the absolute value is the absolute value)
|-
| style="width:250px" |\(\left|-a\right| = |a|\)
|Evenness (reflection symmetry of the graph)
|-
|\(|a - b| = 0 \iff a = b\)
|Identity of indiscernibles (equivalent to positive-definiteness)
|-
|\(|a - b|  \le |a - c| + |c - b|\)
|Triangle inequality (equivalent to subadditivity)
|-
|\(\left|\frac{a}{b}\right| = \frac{|a|}{|b|}\\) (if \(b \ne 0\))
|Preservation of division – equivalent to multiplicativity
|-
|\(|a-b| \geq \bigl| \left|a\right| - \left|b\right| \bigr|\)
|Reverse triangle inequality – equivalent to subadditivity
|}

Two other useful properties concerning inequalities are:
{| style="margin-left:1.6em"
|-
|\(|a| \le b \iff -b \le a \le b\)
|-
|\(|a| \ge b \iff a \le -b\\) or \(a \ge b\)
|}

These relations may be used to solve inequalities involving absolute values. For example:

{| style="margin-left:1.6em"
|-
|\(|x-3| \le 9\)
|\(\iff -9 \le x-3 \le 9\)
|-
|
|\(\iff -6 \le x \le 12\)
|}

The absolute value, as "distance from zero", is used to define the absolute difference between arbitrary real numbers, the standard metric on the real numbers.

### Complex numbers


Since the complex numbers are not ordered, the definition given at the top for the real absolute value cannot be directly applied to complex numbers. However, the geometric interpretation of the absolute value of a real number as its distance from 0 can be generalised. The absolute value of a complex number is defined by the Euclidean distance of its corresponding point in the complex plane from the origin.  This can be computed using the Pythagorean theorem: for any complex number
\(z = x + iy,\)
where \(x\) and \(y\) are real numbers, the absolute value or modulus of \(z\) is denoted \(|z|\) and is defined by
\(|z| = \sqrt{\operatorname{Re}(z)^2 + \operatorname{Im}(z)^2}=\sqrt{x^2 + y^2},\)
the Pythagorean addition of \(x\) and \(y\), where \(\operatorname{Re}(z)=x\) and \(\operatorname{Im}(z)=y\) denote the real and imaginary parts of \(z\), respectively. When the imaginary part \(y\) is zero, this coincides with the definition of the absolute value of the real number \(x\).

When a complex number \(z\) is expressed in its polar form {{nowrap|as \(z = r e^{i \theta},\)}} its absolute value is \(|z| = r.\)

Since the product of any complex number \(z\) and its complex conjugate \(\bar z = x - iy\), with the same absolute value, is always the non-negative real number \(\left(x^2 + y^2\right)\), the absolute value of a complex number \(z\) is the square root {{nowrap|of  \(z \cdot \overline{z},\)}} which is therefore called the absolute square or *squared modulus* of \(z\):
\(|z| = \sqrt{z \cdot \overline{z}}.\)
This generalizes the alternative definition for reals: {{nowrap|\(|x| = \sqrt{x\cdot x}\).}}

The complex absolute value shares the four fundamental properties given above for the real absolute value. The identity \(|z|^n = |z^n|\) is a special case of multiplicativity that is often useful by itself.

## Absolute value function


The real absolute value function is continuous everywhere. It is differentiable everywhere except for 1=*x* = 0.  It is monotonically decreasing on the interval  and monotonically increasing on the interval . Since a real number and its opposite have the same absolute value, it is an even function, and is hence not invertible. The real absolute value function is a piecewise linear, convex function.

For both real and complex numbers, the absolute value function is idempotent (meaning that the absolute value of any absolute value is itself).

### Relationship to the sign function
The absolute value function of a real number returns its value irrespective of its sign, whereas the sign (or signum) function returns a number's sign irrespective of its value. The following equations show the relationship between these two functions:
\(|x| = x \sgn(x),\)
or
\(|x| \sgn(x) = x,\)
and for *x* ≠ 0,
\(\sgn(x) = \frac{|x|}{x} = \frac{x}{|x|}.\)

### Relationship to the max and min functions
Let \(s,t\in\R\), then the following relationship to the minimum and maximum functions hold:
\(|t-s|= -2 \min(s,t)+s+t\)
and
\(|t-s|=2 \max(s,t)-s-t.\)

The formulas can be derived by considering each case \(s>t\) and \(t>s\) separately.

From the last formula one can derive also \(|t|= \max(t,-t)\).

### Derivative
The real absolute value function has a derivative for every *x* ≠ 0, given by a step function equal to the sign function except at 1=*x* = 0 where the absolute value function is not differentiable:

\[
\begin{align}
\frac{d\left|x\right|}{dx}
&= \frac{x}{|x|} = \begin{cases} -1 & x<0 \\  1 & x>0 \end{cases} \\[7mu]
&= \sgn x\quad \text{for } x \ne 0.
\end{align}
\]


The real absolute value function is an example of a continuous function that achieves a global minimum where the derivative does not exist.

The subdifferential of | at 1=*x* = 0 is the interval .

The complex absolute value function is continuous everywhere but complex differentiable *nowhere* because it violates the Cauchy–Riemann equations.

The second derivative of | with respect to  is zero everywhere except zero, where it does not exist. As a generalised function, the second derivative may be taken as two times the Dirac delta function.

### Antiderivative
The antiderivative (indefinite integral) of the real absolute value function is

\(\int \left|x\right| dx = \frac{x\left|x\right|}{2} + C,\)

where  is an arbitrary constant of integration. This is not a complex antiderivative because complex antiderivatives can only exist for complex-differentiable (holomorphic) functions, which the complex absolute value function is not.

### Derivatives and antiderivatives of compositions
The following three formulae are special cases of the chain rule:

\({\text{d}^n \over \text{d}x^n} f(|x|)= (\sgn x)^n f^{(n)}(|x|)\quad \text{for } x \ne 0\,,\)

if the absolute value is inside a function, and

\({\text{d}^n \over \text{d}x^n} |f(x)|=\sgn(f(x)) f^{(n)}(x)\quad \text{for } f(x) \ne 0\,,\)

if another function is inside the absolute value. Combining both, the result is:

\({\text{d}^n \over \text{d}x^n} |f(|x|)|=(\sgn x)^n \sgn(f(|x|)) f^{(n)}(|x|)\quad \text{for } x \ne 0, f(|x|) \ne 0\,.\)

From these formulae and using integration by parts, antiderivatives can also be obtained:

\(\int {f(|x|)\text{d}x}= \sgn (x) F(|x|)\quad \text{for } x \ne 0\,,\)

\(\int {|f(x)| \text{d}x} = \int {{|f(x)| \over f(x)} f(x) \text{d}x} = \sgn(f(x)) F(x)\quad \text{for } f(x) \ne 0\,,\)

\(\int {|f(|x|)| \text{d}x} = \int {{|f(|x|)| \over f(|x|)} f(|x|) \text{d}x} = \sgn (x) \sgn(f(|x|)) F(|x|)\quad \text{for } x \ne 0, f(|x|) \ne 0\,,\)

supposing the derivative of the sign function is 0.

### Power rule for expressions with absolute values
Using chain and product rules, the power rule for expressions of the type \(x^n |x|^m\) can be written as:

\({\text{d} \over \text{d}x} x^n |x|^m = (n+m)x^{n-1}|x|^m\,.\)

This holds true even for \(n=0\):

\({\text{d} \over \text{d}x} |x|^m = mx^{-1}|x|^m\,.\)

## Distance

The absolute value is closely related to the idea of distance. As noted above, the absolute value of a real or complex number is the distance from that number to the origin, along the real number line, for real numbers, or in the complex plane, for complex numbers, and more generally, the absolute value of the difference of two real or complex numbers is the distance between them.

The standard Euclidean distance between two points

\[
a = (a_1, a_2, \dots , a_n)
\]

and

\[
b = (b_1, b_2, \dots , b_n)
\]

in Euclidean -space is defined as:

\[
\sqrt{\textstyle\sum_{i=1}^n(a_i-b_i)^2}.
\]


This can be seen as a generalisation, since for \(a_1\) and \(b_1\) real, i.e. in a 1-space, according to the alternative definition of the absolute value,

\(|a_1 - b_1| = \sqrt{(a_1 - b_1)^2} = \sqrt{\textstyle\sum_{i=1}^1(a_i-b_i)^2},\)

and for \(a = a_1 + i a_2\) and \(b = b_1 + i b_2\) complex numbers, i.e. in a 2-space,

{|
|-
|\(|a - b|\)
|\(= |(a_1 + i a_2) - (b_1 + i b_2)|\)
|-
|
|\(= |(a_1 - b_1) + i(a_2 - b_2)|\)
|-
|
|\(= \sqrt{(a_1 - b_1)^2 + (a_2 - b_2)^2} = \sqrt{\textstyle\sum_{i=1}^2(a_i-b_i)^2}.\)
|}

The above shows that the "absolute value"-distance, for real and complex numbers, agrees with the standard Euclidean distance, which they inherit as a result of considering them as one and two-dimensional Euclidean spaces, respectively.

The properties of the absolute value of the difference of two real or complex numbers: non-negativity, identity of indiscernibles, symmetry and the triangle inequality given above, can be seen to motivate the more general notion of a distance function as follows:

A real valued function  on a set *X* × *X* is called a metric (or a *distance function*) on , if it satisfies the following four axioms:
{|
|-
|style="width:250px" | \(d(a, b) \ge 0\)
|Non-negativity
|-
|\(d(a, b) = 0 \iff a = b\)
|Identity of indiscernibles
|-
|\(d(a, b) = d(b, a)\)
|Symmetry
|-
|\(d(a, b)  \le d(a, c) + d(c, b)\)
|Triangle inequality
|}

## Generalizations
### Ordered rings
The definition of absolute value given for real numbers above can be extended to any ordered ring. That is, if  is an element of an ordered ring *R*, then the **absolute value** of , denoted by |*a*|, is defined to be:

\[
|a| = \left\{
   \begin{array}{rl}
     a, & \text{if }  a \geq 0 \\
     -a, & \text{if } a < 0.
   \end{array}\right.
\]

where −*a* is the additive inverse of , 0 is the additive identity, and < and ≥ have the usual meaning with respect to the ordering in the ring.

### Fields

The four fundamental properties of the absolute value for real numbers can be used to generalise the notion of absolute value to an arbitrary field, as follows.

A real-valued function  on a field  is called an *absolute value* (also a *modulus*, *magnitude*, *value*, or *valuation*) if it satisfies the following four axioms:

{| cellpadding=10
|-
|\(v(a) \ge 0\)
|Non-negativity
|-
|\(v(a) = 0 \iff a = \mathbf{0}\)
|Positive-definiteness
|-
|\(v(ab) = v(a) v(b)\)
|Multiplicativity
|-
|\(v(a+b)  \le v(a) + v(b)\)
|Subadditivity or the triangle inequality
|}

Where **0** denotes the additive identity of . It follows from positive-definiteness and multiplicativity that 1=*v*(**1**) = 1, where **1** denotes the multiplicative identity of . The real and complex absolute values defined above are examples of absolute values for an arbitrary field.

If  is an absolute value on , then the function  on *F* × *F*, defined by 1=*d*(*a*, *b*) = *v*(*a* − *b*), is a metric and the following are equivalent:

*  satisfies the ultrametric inequality \(d(x, y) \leq \max(d(x,z),d(y,z))\) for all , ,  in .
* \(\left\{ v\left( \sum_{k=1}^n \mathbf{1}\right) : n \in \N \right\}\) is bounded in **R**.
* \(v\left({\textstyle \sum_{k=1}^n } \mathbf{1}\right) \le 1\\) for every \(n \in \N\).
* \(v(a) \le 1 \Rightarrow v(1+a) \le 1\\) for all \(a \in F\).
* \(v(a + b) \le \max \{v(a), v(b)\}\\) for all \(a, b \in F\).

An absolute value which satisfies any (hence all) of the above conditions is said to be **non-Archimedean**, otherwise it is said to be Archimedean.

### Vector spaces

Again the fundamental properties of the absolute value for real numbers can be used, with a slight modification, to generalise the notion to an arbitrary vector space.

A real-valued function on a vector space  over a field , represented as , is called an **absolute value**, but more usually a **norm**, if it satisfies the following axioms:

For all  in , and **v**, **u** in ,

{| cellpadding=10
|-
|\(\|\mathbf{v}\|  \ge 0\)
|Non-negativity
|-
|\(\|\mathbf{v}\| = 0 \iff \mathbf{v} = 0\)
|Positive-definiteness
|-
|\(\|a \mathbf{v}\| = \left|a\right| \left\|\mathbf{v}\right\|\)
|Absolute homogeneity or positive scalability
|-
|\(\|\mathbf{v} + \mathbf{u}\| \le \|\mathbf{v}\| + \|\mathbf{u}\|\)
|Subadditivity or the triangle inequality
|}

The norm of a vector is also called its *length* or *magnitude*.

In the case of Euclidean space \(\mathbb{R}^n\), the function defined by

\(\|(x_1, x_2, \dots , x_n) \| = \sqrt{\textstyle\sum_{i=1}^{n} x_i^2}\)

is a norm called the Euclidean norm. When the real numbers \(\mathbb{R}\) are considered as the one-dimensional vector space \(\mathbb{R}^1\), the absolute value is a norm, and is the -norm (see L<sup>p</sup> space) for any . In fact the absolute value is the "only" norm on \(\mathbb{R}^1\), in the sense that, for every norm  on \(\mathbb{R}^1\), 1= ⋅ |*x*|}}.

The complex absolute value is a special case of the norm in an inner product space, which is identical to the Euclidean norm when the complex plane is identified as the Euclidean plane \(\mathbb{R}^2\).

### Composition algebras

Every composition algebra *A* has an involution *x* → *x** called its **conjugation**. The product in *A* of an element *x* and its conjugate *x** is written *N*(*x*) = *x x** and called the **norm of x**.

The real numbers \(\mathbb{R}\), complex numbers \(\mathbb{C}\), and quaternions \(\mathbb{H}\) are all composition algebras with norms given by definite quadratic forms. The absolute value in these division algebras is given by the square root of the composition algebra norm.

In general, the norm of a composition algebra may be a quadratic form that is not definite and has null vectors. However, as in the case of division algebras, when an element *x* has a non-zero norm, then *x* has a multiplicative inverse given by *x**/*N*(*x*).

## See also
*Least absolute values

## Notes


## Footnotes


## References
*
 | title = Calculus Problems
 | doi = 10.1007/978-3-319-15428-2
 | isbn = 978-3-319-15428-2
}}
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*

## External links
*
*
*

