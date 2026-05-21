> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Absolute_value) — CC BY-SA 4.0

# Absolute value

In mathematics, the **absolute value** or **modulus** of a real number \(x\), denoted \(|x|\), is the (non-negative) magnitude of \(x\) measured without regard to its sign. Namely, \(|x|=x\) if \(x\) is a positive number, and \(|x|=-x\) if \(x\) is negative (in which case negating \(x\) makes \(-x\) positive), and \(|0|=0\). For example, the absolute value of 3 is 3, and the absolute value of −3 is also 3. The absolute value of a number may be thought of as its distance from zero.

Generalisations of the absolute value for real numbers occur in a wide variety of mathematical settings. For example, an absolute value is also defined for the complex numbers, the quaternions, ordered rings, fields and vector spaces. The absolute value is closely related to the notions of magnitude, distance, and norm in various mathematical and physical contexts.

## Terminology and notation
In 1806, Jean-Robert Argand introduced the term *module*, meaning *unit of measure* in French, specifically for the *complex* absolute value, and it was borrowed into English in 1866 as the Latin equivalent *modulus*. The term *absolute value* has been used in this sense from at least 1806 in French and 1857 in English.The notation |, with a vertical bar on each side, was introduced by Karl Weierstrass in 1841. Other names for *absolute value* include *numerical value* and *magnitude*. The absolute value of \(x\) has also been denoted \(\operatorname{abs} x\) in some mathematical publications, and in spreadsheets, programming languages, and computational software packages, the absolute value of \(x\) is generally represented by abs(*x*), or a similar expression, as it has been since the earliest days of high-level programming languages.

The vertical bar notation also appears in a number of other mathematical contexts: for example, when applied to a set, it denotes its cardinality; when applied to a matrix, it denotes its determinant. Vertical bars denote the absolute value only for algebraic objects for which the notion of an absolute value is defined, notably an element of a normed division algebra, for example, a real number, a complex number, or a quaternion. A closely related but distinct notation is the use of vertical bars for either the Euclidean norm or sup norm of a vector in \(\R^n\), although double vertical bars with subscripts (\(\|\cdot\|_2\) and \(\|\cdot\|_\infty\), respectively) are a more common and less ambiguous notation.

## Definition and properties
### Real numbers
For any real number \(x\), the absolute value or modulus of \(x\) is denoted by \(|x|\), with a vertical bar on each side of the quantity, and is defined as
\(|x| =
 \begin{cases}
 x, & \text{if } x \geq 0 \\
 -x, & \text{if } x < 0.
 \end{cases}\)

The absolute value of \(x\) is thus always either a positive number or zero, but never negative. When \(x\) itself is negative (\(x < 0\)), then its absolute value is necessarily positive (\(|x|=-x>0\)).

From an analytic geometry point of view, the absolute value of a real number is that number's distance from zero along the real number line, and more generally, the absolute value of the difference of two real numbers (their absolute difference) is the distance between them. The notion of an abstract distance function in mathematics can be seen to be a generalisation of the absolute value of the difference. See below.

Since the square root symbol represents the unique *positive* square root, when applied to a positive number, it follows that
\(|x| = \sqrt{x^2}.\)
This is equivalent to the definition above, and may be used as an alternative definition of the absolute value of real numbers.

The absolute value has the following four fundamental properties (\(a\), \(b\) are real numbers), that are used for generalization of this notion to other domains:

<dl>
  <dt>|\(|a| \ge 0\)</dt>
  <dd>Non-negativity</dd>
  <dt>\(|a| = 0 \iff a = 0\)</dt>
  <dd>Positive-definiteness</dd>
  <dt>\(|ab| = \left|a\right| \left|b\right|\)</dt>
  <dd>Multiplicativity</dd>
  <dt>\(|a+b| \le |a| + |b|\)</dt>
  <dd>Subadditivity, specifically the triangle inequality</dd>
</dl>

Non-negativity, positive definiteness, and multiplicativity are readily apparent from the definition. To see that subadditivity holds, first note that \(|a+b|=s(a+b)\) where \(s=\pm 1\), with its sign chosen to make the result positive. Now; since \(-1 \cdot x \le |x|\) and \(+1 \cdot x \le |x|\), it follows that, whichever of \(\pm1\) is the value of \(s\), one has \(s \cdot x\leq |x|\) for all real \(x\). Consequently, \(|a+b|=s \cdot (a+b) = s \cdot a + s \cdot b \leq |a| + |b|\), as desired.

Some additional useful properties are given below. These are either immediate consequences of the definition or implied by the four fundamental properties above.

<dl>
  <dt>|\(\bigl| \left|a\right| \bigr| = |a|\)</dt>
  <dd>Idempotence (the absolute value of the absolute value is the absolute value)</dd>
  <dt>|\(\left|-a\right| = |a|\)</dt>
  <dd>Evenness (reflection symmetry of the graph)</dd>
  <dt>\(|a - b| = 0 \iff a = b\)</dt>
  <dd>Identity of indiscernibles (equivalent to positive-definiteness)</dd>
  <dt>\(|a - b| \le |a - c| + |c - b|\)</dt>
  <dd>Triangle inequality (equivalent to subadditivity)</dd>
  <dt>\(\left|\frac{a}{b}\right| = \frac{|a|}{|b|}\) (if \(b \ne 0\))</dt>
  <dd>Preservation of division – equivalent to multiplicativity</dd>
  <dt>\(|a-b| \geq \bigl| \left|a\right| - \left|b\right| \bigr|\)</dt>
  <dd>Reverse triangle inequality – equivalent to subadditivity</dd>
</dl>

Two other useful properties concerning inequalities are:
| \(|a| \le b \iff -b \le a \le b\) |
| --- |
| \(|a| \ge b \iff a \le -b\) or \(a \ge b\) |

These relations may be used to solve inequalities involving absolute values. For example:

| \(|x-3| \le 9\) | \(\iff -9 \le x-3 \le 9\) |
| --- | --- |
| \(\iff -6 \le x \le 12\) |  |

The absolute value, as "distance from zero", is used to define the absolute difference between arbitrary real numbers, the standard metric on the real numbers.

### Complex numbers

Since the complex numbers are not ordered, the definition given at the top for the real absolute value cannot be directly applied to complex numbers. However, the geometric interpretation of the absolute value of a real number as its distance from 0 can be generalised. The absolute value of a complex number is defined by the Euclidean distance of its corresponding point in the complex plane from the origin. This can be computed using the Pythagorean theorem: for any complex number
\(z = x + iy,\)
where \(x\) and \(y\) are real numbers, the absolute value or modulus of \(z\) is denoted \(|z|\) and is defined by
\(|z| = \sqrt{\operatorname{Re}(z)^2 + \operatorname{Im}(z)^2}=\sqrt{x^2 + y^2},\)
the Pythagorean addition of \(x\) and \(y\), where \(\operatorname{Re}(z)=x\) and \(\operatorname{Im}(z)=y\) denote the real and imaginary parts of \(z\), respectively. When the imaginary part \(y\) is zero, this coincides with the definition of the absolute value of the real number \(x\).

When a complex number \(z\) is expressed in its polar form *x*
