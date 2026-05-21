> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Pythagorean_trigonometric_identity) — CC BY-SA 4.0

# Pythagorean identity sin²θ + cos²θ = 1

The **Pythagorean trigonometric identity**, also called simply the **Pythagorean identity**, is an identity expressing the Pythagorean theorem in terms of trigonometric functions. Along with the sum-of-angles formulae, it is one of the basic relations between the sine and cosine functions.

The identity is
\(\sin^2 \theta + \cos^2 \theta = 1\),

where \(\sin^2 \theta\) means \((\sin\theta)^2\) and \(\cos^2 \theta\) means \((\cos\theta)^2\).

## Proofs and their relationships to the Pythagorean theorem

### Proof based on right-angle triangles
Any similar triangles have the property that if we select the same angle in all of them, the ratio of the two sides defining the angle is the same regardless of which similar triangle is selected, regardless of its actual size: the ratios depend upon the three angles, not the lengths of the sides. Thus for either of the similar right triangles in the figure, the ratio of its horizontal side to its hypotenuse is the same, namely cos *θ*.

The elementary definitions of the sine and cosine functions in terms of the sides of a right triangle are:

\[
\begin{alignat}{3}
\sin \theta &= \frac{\mathrm{opposite{\mathrm{hypotenuse= \frac{b}{c} \\
\cos \theta &= \frac{\mathrm{adjacent{\mathrm{hypotenuse= \frac{a}{c}
\end{alignat}
\]

The Pythagorean identity follows by squaring both definitions above, and adding; the left-hand side of the identity then becomes

\[
\frac{\mathrm{opposite}^2 + \mathrm{adjacent}^2}{\mathrm{hypotenuse}^2}
\]

which by the Pythagorean theorem is equal to 1. This definition is valid for all angles, due to the definition of defining *x* = cos *θ* and *y* sin *θ* for the unit circle and thus *x* = *c* and *y* = *c* for a circle of radius and reflecting our triangle in the and setting *a* = *x* and *b* = *y*.

Alternatively, the identities found at Trigonometric symmetry, shifts, and periodicity may be employed. By the periodicity identities we can say if the formula is true for −*π* < *θ* ≤ *π* then it is true for all real. Next we prove the identity in the range. To do this we let *t* , will now be in the range 0 < *t* ≤. We can then make use of squared versions of some basic shift identities (squaring conveniently removes the minus signs):

\[
\sin^2\theta + \cos^2\theta = \sin^2\left(t + \tfrac{\pi}{2}\right) + \cos^2\left(t + \tfrac{\pi}{2}\right) = \cos^2 t + \sin^2 t = 1.
\]

Finally, it remains is to prove the formula for −*π* < *θ* < 0; this can be done by squaring the symmetry identities to get

\[
\sin^2\theta = \sin^2(-\theta)\text{ and }\cos^2\theta = \cos^2(-\theta).
\]

#### Related identities

The two identities
\[
\begin{align}
1 + \tan^2 \theta &= \sec^2 \theta \\
1 + \cot^2 \theta &= \csc^2 \theta
\end{align}
\]
 are also called Pythagorean trigonometric identities. If one leg of a right triangle has length 1, then the tangent of the angle adjacent to that leg is the length of the other leg, and the secant of the angle is the length of the hypotenuse.

\[
\begin{align}
\tan \theta &= \frac{b}{a}\, \\
\sec \theta &= \frac{c}{a}\,.
\end{align}
\]

In this way, this trigonometric identity involving the tangent and the secant follows from the Pythagorean theorem. The angle opposite the leg of length 1 (this angle can be labeled *φ* = ) has cotangent equal to the length of the other leg, and cosecant equal to the length of the hypotenuse. In that way, this trigonometric identity involving the cotangent and the cosecant also follows from the Pythagorean theorem.

The following table gives the identities with the factor or divisor that relates them to the main identity.

| scope="col" | Divisor | scope="col" | Divisor Equation | scope="col" | Derived Identity | scope="col" | Derived Identity (Alternate) |
| --- | --- | --- | --- |
| scope="row" | | \(\frac{\sin^2 \theta}{\cos^2 \theta} + \frac{\cos^2 \theta}{\cos^2 \theta} = \frac{1}{\cos^2 \theta}\) | | \(\tan^2 \theta + 1 = \sec^2 \theta\) | \(\begin{align} \sec^2\theta - \tan^2\theta = 1\\ (\sec\theta - \tan\theta)(\sec\theta + \tan\theta) = 1\\ \end{align}\) |
| scope="row" | | \(\frac{\sin^2 \theta}{\sin^2 \theta} + \frac{\cos^2 \theta}{\sin^2 \theta} = \frac{1}{\sin^2 \theta}\) | | \(1 + \cot^2 \theta = \csc^2 \theta\) | \(\begin{align} \csc^2\theta - \cot^2\theta = 1\\ (\csc\theta - \cot\theta)(\csc\theta + \cot\theta) = 1\\ \end{align}\) |

### Proof using the unit circle

The unit circle centered at the origin in the Euclidean plane is defined by the equation:

\(x^2 + y^2 = 1.\)

Given an angle *θ*, there is a unique point *P* on the unit circle at an anticlockwise angle of *θ* from the *x*-axis, and the *x*- and *y*-coordinates of *P* are:

\[
x = \cos\theta \ \text{ and }\ y = \sin\theta.
\]

Consequently, from the equation for the unit circle,
\[
\cos^2 \theta + \sin^2 \theta = 1,
\]

the Pythagorean identity.

In the figure, the point has a -coordinate, and is appropriately given by *x* = cos *θ*, which is a negative number: cos *θ* = −cos(*π* − *θ*). Point has a positive -coordinate, and sin *θ* = sin(*π* − *θ*) > 0. As increases from zero to the full circle *θ* = 2*π*, the sine and cosine change signs in the various quadrants to keep and with the correct signs. The figure shows how the sign of the sine function varies as the angle changes quadrant.

Because the - and -axes are perpendicular, this Pythagorean identity is equivalent to the Pythagorean theorem for triangles with hypotenuse of length 1 (which is in turn equivalent to the full Pythagorean theorem by applying a similar-triangles argument). See Unit circle for a short explanation.

### Proof using power series
The trigonometric functions may also be defined using power series, namely for (an angle measured in radians):

\[
\begin{align}
 \sin x &= \sum_{n = 0}^\infty \frac{(-1)^n}{(2n + 1)!} x^{2n + 1},\\
 \cos x &= \sum_{n = 0}^\infty \frac{(-1)^n}{(2n)!} x^{2n}.
\end{align}
\]

Using the multiplication formula for power series at Multiplication and division of power series (suitably modified to account for the form of the series here) we obtain

\[
\begin{align}
\sin^2 x & = \sum_{i = 0}^\infty \sum_{j = 0}^\infty \frac{(-1)^i}{(2i + 1)!} \frac{(-1)^j}{(2j + 1)!} x^{(2i + 1) + (2j + 1)} \\
& = \sum_{n = 1}^\infty \left(\sum_{i = 0}^{n - 1} \frac{(-1)^{n - 1{(2i + 1)!(2(n - i - 1) + 1)!}\right) x^{2n} \\
& = \sum_{n = 1}^\infty \left( \sum_{i = 0}^{n - 1} {2n \choose 2i + 1} \right) \frac{(-1)^{n - 1{(2n)!} x^{2n},\\
\cos^2 x & = \sum_{i = 0}^\infty \sum_{j = 0}^\infty \frac{(-1)^i}{(2i)!} \frac{(-1)^j}{(2j)!} x^{(2i) + (2j)} \\
& = \sum_{n = 0}^\infty \left(\sum_{i = 0}^n \frac{(-1)^n}{(2i)!(2(n - i))!}\right) x^{2n} \\
& = \sum_{n = 0}^\infty \left( \sum_{i = 0}^n {2n \choose 2i} \right) \frac{(-1)^n}{(2n)!} x^{2n}.
\end{align}
\]

In the expression for sin\(^{2}\), must be at least 1, while in the expression for cos\(^{2}\), the constant term is equal to 1. The remaining terms of their sum are (with common factors removed)

\[
\begin{align}
\sum_{i = 0}^n {2n \choose 2i} - \sum_{i = 0}^{n - 1} {2n \choose 2i + 1} &= \sum_{j = 0}^{2n} (-1)^j {2n \choose j} \\
&= (1 - 1)^{2n} = 0
\end{align}
\]

by the binomial theorem. Consequently,

\[
\sin^2 x + \cos^2 x = 1,
\]

which is the Pythagorean trigonometric identity.

When the trigonometric functions are defined in this way, the identity in combination with the Pythagorean theorem shows that these power series parameterize the unit circle, which we used in the previous section. This definition constructs the sine and cosine functions in a rigorous fashion and proves that they are differentiable, so that in fact it subsumes the previous two.

### Proof using the differential equation
Sine and cosine can be defined as the two solutions to the differential equation:

\[
y'' + y = 0
\]

satisfying respectively *y*(0) = 0, and *y*(0) = 1,. It follows from the theory of ordinary differential equations that the first solution, sine, has the second, cosine, as its derivative, and it follows from this that the derivative of cosine is the negative of the sine. The identity is equivalent to the assertion that the function

\[
z = \sin^2 x + \cos^2 x
\]

is constant and equal to 1. Differentiating using the chain rule gives:

\[
\frac{d}{dx} z = 2 \sin x \cos x + 2 \cos x(-\sin x) = 0,
\]

so is constant. A calculation confirms that *z*(0) = 1, and is a constant so *z* = 1 for all , so the Pythagorean identity is established.

A similar proof can be completed using power series as above to establish that the sine has as its derivative the cosine, and the cosine has as its derivative the negative sine. In fact, the definitions by ordinary differential equation and by power series lead to similar derivations of most identities.

This proof of the identity has no direct connection with Euclid's demonstration of the Pythagorean theorem.

### Proof using Euler's formula
Factoring \(\cos^2 \theta + \sin^2 \theta\) as the complex difference of two squares and using Euler's formula \(e^{i\theta} = \cos\theta + i\sin\theta\),

\[
\begin{align}
\cos^2 \theta + \sin^2 \theta
&= \cos^2 \theta - i^2 \sin^2 \theta \\[3mu]
&= (\cos\theta + i\sin\theta)(\cos\theta - i\sin\theta) \\[3mu]
&= e^{i\theta}e^{-i\theta} \\
&= e^{i\theta-i\theta} = e^0 = 1
\end{align}
\]
