
# De Moivre’s Theorem

Source: algebrica.org — CC BY-NC 4.0
https://algebrica.org/de-moivre-theorem/

## Motivation for De Moivre’s Theorem

Suppose we want to compute the power of a complex number \\( z \in \mathbb{C} \\). The most straightforward approach is to start from its [algebraic form](../complex-numbers-introduction/) and expand the expression directly. For example, given \\(z = a + ib\\) we may want to calculate its square. We have:

\\[
\begin{align}
z^2 &= (a + ib)^2 \\\\[0.5em]
&= a^2 + i2ab - b^2
\end{align}
\\]

While this method is valid, as the exponent increases beyond three, the calculations become increasingly tedious and impractical. Expanding higher powers algebraically yields lengthy expressions and more terms, reducing the practicality of this approach. In these situations, [De Moivre's Theorem](../de-moivre-theorem/) provides a more efficient and elegant solution.

- - -
## De Moivre's theorem and exponential notation for complex numbers

De Moivre's theorem provides a method for computing powers and roots of complex numbers, whether written in [trigonometric](../complex-numbers-trigonometric-form/) or [exponential form](../complex-numbers-exponential-form/). Consider a complex number \\( z \\) raised to an integer power \\( n \in \mathbb{Z} \\). That is,

\\[
z^n \quad n \in \mathbb{Z}
\\]

Rewrite the number \\(z\\) in trigonometric form:

\\[
z = r(\cos\theta + i\sin\theta)
\\]

For any [integer](../integers/) \\( n \\), the power \\( z^n \\) can be computed by raising the modulus to the \\( n \\)-th power and multiplying the angle by \\( n \\). The result is a new complex number in polar form. We have:

\\[
z^n = r^n \left(\cos(n\theta) + i\sin(n\theta)\right)
\\]

This identity holds for all integers \\( n \\), including negative ones. When \\( n \\) is a rational number \\( n = p/q \\), the formula still applies but yields one of the \\( q \\) distinct roots; the full set of roots requires considering all values of the argument of the form \\( \theta + 2k\pi \\) for \\( k = 0, 1, \dots, q-1 \\).

---

Now rewrite the complex number \\( z \\) using Euler's identity, in exponential form:

\\[
e^{i\theta} = \cos\theta + i\sin\theta
\\]

We obtain:

\\[
z = re^{i\theta}
\\]

This formulation allows us to interpret complex numbers in a way that's deeply aligned with the structure of exponentiation. It also turns De Moivre's Theorem into something almost automatic. When we raise \\( z \\) to an integer power, we simply apply the usual exponent laws:

\\[
z^n = (re^{i\theta})^n = r^n e^{in\theta}
\\]

There's no algebra to expand, no trigonometric identities to manipulate. The modulus is raised to the power \\( n \\), and the argument is multiplied by \\( n \\).

- - -
## A proof by induction

De Moivre's Theorem states that for any integer \\( n \\) and any complex number \\( z = r(\cos\theta + i\sin\theta) \\):

\\[
z^n = r^n\bigl(\cos(n\theta) + i\sin(n\theta)\bigr)
\\]

The formula can be established by [induction](../principle-of-mathematical-induction/) on \\( n \\). The argument has two parts:

+ Verifying the base case.
+ Showing that validity at step \\( n \\) forces validity at step \\( n + 1 \\).

- - -

For the base case, setting \\( n = 1 \\) reduces the formula to \\( z = r(\cos\theta + i\sin\theta) \\), which is the trigonometric form of \\( z \\) by definition. For the inductive step, suppose the formula holds for some integer \\( n \geq 1 \\):

\\[
z^n = r^n\bigl(\cos(n\theta) + i\sin(n\theta)\bigr)
\\]

Multiplying both sides by \\( z = r(\cos\theta + i\sin\theta) \\) and expanding the product we obtain:

\\[
z^{n+1} = r^{n+1}\Bigl[\bigl(\cos(n\theta)\cos\theta - \sin(n\theta)\sin\theta\bigr) + i\bigl(\sin(n\theta)\cos\theta + \cos(n\theta)\sin\theta\bigr)\Bigr]
\\]

The two expressions in brackets are the addition formulas for [cosine and sine](../sine-and-cosine/) respectively and applying them yields:

\\[
z^{n+1} = r^{n+1}\bigl(\cos\bigl((n+1)\theta\bigr) + i\sin\bigl((n+1)\theta\bigr)\bigr)
\\]

The identity holds at step \\( n + 1 \\), which completes the induction.

- - -
## Example 1

For example, squaring the complex number \\( z = re^{i\theta} \\) gives:

\\[
z^2 = (re^{i\theta})^2 = r^2 e^{i2\theta}
\\]

The result is a new complex number whose modulus is \\( r^2 \\) and whose argument is \\( 2\theta. \\) In geometric terms, this means the [vector](../vectors/) is stretched by a factor of \\( r^2 \\) and rotated to double its original angle.

- - -

## Example 2

Let’s try to compute \\( z^4 \\) for the complex number \\( z = 2 + 2i \\).
First, we determine the modulus of \\( z \\):

\\[
|z| = \sqrt{2^2 + 2^2} = \sqrt{8} = 2\sqrt{2}
\\]

> The modulus of a complex number represents its distance from the origin in the complex plane. It is calculated using the [Pythagorean theorem](../pythagorean-theorem/).

- - -

Next, we determine the argument of \\( z \\):

\\[
\theta = \arg(z) = \arctan\left(\frac{2}{2}\right) = \frac{\pi}{4}
\\]

> The argument of a complex number is the angle it makes with the positive real axis, measured counterclockwise. In this case, since both the real and imaginary parts are equal, the angle is exactly \\( 45^\circ \\), or \\( \frac{\pi}{4} \\) radians.

- - -

We can now express \\( z \\) in exponential form:

\\[
z = 2\sqrt{2} \cdot e^{i\frac{\pi}{4}}
\\]

This compact form makes it much easier to raise \\( z \\) to a power, since we can apply the rules of exponents directly. Applying De Moivre’s Theorem, we compute:

\\[
z^4 = (2\sqrt{2})^4 \cdot e^{i \cdot 4 \cdot \frac{\pi}{4}} = (2\sqrt{2})^4 \cdot e^{i\pi}
\\]

Let’s simplify:

\\[
(2\sqrt{2})^4 = (2^1 \cdot 2^{1/2})^4 = 2^6 = 64
\\]

- - -

Since \\(e^{i\pi} = -1\\) we find:

\\[
z^4 = 64 \cdot (-1) = -64
\\]

> The same result can also be obtained by expanding the expression algebraically.

So the fourth power of \\( z = 2 + 2i \\) is the real number \\( -64 \\).

- - -
## Deriving trigonometric identities

One of the most practical applications of De Moivre's Theorem is the derivation of explicit formulas for [sine and cosine](../sine-and-cosine/), in particular for \\( \cos(n\theta) \\) and \\( \sin(n\theta) \\) in terms of powers of \\( \cos\theta \\) and \\( \sin\theta \\). The idea is straightforward: expand the left-hand side of the theorem using the [binomial formula](../binomial-coefficient/), then separate real and imaginary parts.

For \\( n = 3 \\), the theorem gives:

\\[
(\cos\theta + i\sin\theta)^3 = \cos(3\theta) + i\sin(3\theta)
\\]

Expanding the left-hand side with the binomial formula:

\\[
(\cos\theta + i\sin\theta)^3 = \cos^3\theta + 3i\cos^2\theta\sin\theta + 3i^2\cos\theta\sin^2\theta + i^3\sin^3\theta
\\]

Using \\( i^2 = -1 \\) and \\( i^3 = -i \\):

\\[
= \bigl(\cos^3\theta - 3\cos\theta\sin^2\theta\bigr) + i\bigl(3\cos^2\theta\sin\theta - \sin^3\theta\bigr)
\\]

Equating real and imaginary parts with the right-hand side:

\\[
\cos(3\theta) = \cos^3\theta - 3\cos\theta\sin^2\theta
\\]

\\[
\sin(3\theta) = 3\cos^2\theta\sin\theta - \sin^3\theta
\\]

These are the triple angle formulas for cosine and sine. Both follow directly from a single application of the binomial expansion, with no need for repeated use of [addition formulas](../reduction-formulas-and-reference-angles/) or any other intermediate result. The same procedure extends to any integer \\( n \\): the binomial expansion of \\( (\cos\theta + i\sin\theta)^n \\) always yields \\( \cos(n\theta) \\) as its real part and \\( \sin(n\theta) \\) as its imaginary part.

- - -
## Finding complex roots with De Moivre's theorem

De Moivre's Theorem isn't just useful for powers. It also gives us a clean and elegant way to find the roots of a complex number. Suppose we want to solve:

\\[
z^n = w
\\]

where \\( w \in \mathbb{C} \\). This means we're looking for all the complex numbers \\( z \\) such that raising them to the \\( n \\)-th power gives \\( w \\). First, we write \\( w \\) in exponential form. Since the argument of a complex number is defined up to multiples of \\( 2\pi \\), we write:

\\[
w = r e^{i(\theta + 2k\pi)}, \quad k \in \mathbb{Z}
\\]

Applying De Moivre's Theorem to \\( z^n = w \\) and taking the \\( n \\)-th root of both sides, we obtain the general formula for the \\( n \\)-th roots:

\\[
z_k = \sqrt[n]{r} \cdot e^{i\left(\frac{\theta + 2k\pi}{n}\right)}, \quad \text{for } k = 0, 1, \dots, n - 1
\\]

This gives all the \\( n \\) distinct complex roots. They lie on a circle of radius \\( \sqrt[n]{r} \\), equally spaced by an angle of \\( \dfrac{2\pi}{n} \\). This means the roots are arranged like the vertices of a regular polygon with \\( n \\) sides inscribed in a circle of radius \\( \sqrt[n]{r} \\). In the case of cube roots, we get three points on a circle, each separated by an angle of \\( \dfrac{2\pi}{3} \\), forming an equilateral triangle in the complex plane.

- - -
## Example 3

Let's find all the complex solutions to the equation:

\\[
z^3 = 1
\\]

At first glance, it seems obvious that \\( z = 1 \\) is a solution. But since we're working in the complex plane, we know there are three cube roots in total, equally spaced around the [unit circle](../unit-circle).

Since the argument of a complex number is defined up to multiples of \\( 2\pi \\), we write \\( 1 \\) in exponential form as:

\\[
1 = e^{i \cdot 2k\pi}, \quad k \in \mathbb{Z}
\\]

Applying the general root formula with \\( r = 1 \\) and \\( \theta = 0 \\), we obtain:

\\[
z_k = \sqrt[3]{1} \cdot e^{i\left(\frac{0 + 2k\pi}{3}\right)} = e^{i \cdot \frac{2k\pi}{3}}, \quad \text{for } k = 0, 1, 2
\\]

Let's now evaluate the three roots explicitly.

For \\( k = 0 \\):

\\[
z_0 = e^{i \cdot 0} = \cos(0) + i\sin(0) = 1
\\]

For \\( k = 1 \\):

\\[
z_1 = e^{i \cdot \frac{2\pi}{3}} = \cos\left(\frac{2\pi}{3}\right) + i\sin\left(\frac{2\pi}{3}\right) = -\frac{1}{2} + \frac{\sqrt{3}}{2}i
\\]

For \\( k = 2 \\):

\\[
z_2 = e^{i \cdot \frac{4\pi}{3}} = \cos\left(\frac{4\pi}{3}\right) + i\sin\left(\frac{4\pi}{3}\right) = -\frac{1}{2} - \frac{\sqrt{3}}{2}i
\\]

These are the three cube roots of 1, arranged in the complex plane like the vertices of an equilateral triangle. Together, they form what are known as the [cube roots of unity](../roots-of-unity/).