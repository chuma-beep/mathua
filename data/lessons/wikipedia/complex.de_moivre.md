> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/De_Moivre%27s_formula) — CC BY-SA 4.0

# De Moivre's theorem

In mathematics, **de Moivre's formula** (also known as **de Moivre's theorem** and **de Moivre's identity**) states that for any real number  and integer ,

\[
\big(\cos x + i \sin x\big)^n = \cos nx + i \sin nx,
\]

where  is the imaginary unit (*i*<sup>2</sup> ). The formula is named after Abraham de Moivre, although he never stated it in his works. The expression cos *x* + *i* sin *x* is sometimes abbreviated to cis *x*.

The formula is important because it connects complex numbers and trigonometry. By expanding the left hand side and then comparing the real and imaginary parts under the assumption that  is real, it is possible to derive useful expressions for cos *nx* and sin *nx* in terms of cos *x* and sin *x*.

As written, the formula is not valid for non-integer powers . However, there are generalizations of this formula valid for other exponents. These can be used to give explicit expressions for the th roots of unity, that is, complex numbers  such that *z<sup>n</sup>* .

Using the standard extensions of the sine and cosine functions to complex numbers, the formula is valid even when  is an arbitrary complex number.

## Example
For \(x = \frac{\pi}{6}\) and \(n = 2\), de Moivre's formula asserts that
\(\left(\cos\bigg(\frac{\pi}{6}\bigg) + i \sin\bigg(\frac{\pi}{6}\bigg)\right)^2 = \cos\bigg(2 \cdot \frac{\pi}{6}\bigg) + i \sin \bigg(2 \cdot \frac{\pi}{6}\bigg),\)
or equivalently that
\(\left(\frac{\sqrt{3}}{2} + \frac{i}{2}\right)^2 = \frac{1}{2} + \frac{i\sqrt{3}}{2}.\)
In this example, it is easy to check the validity of the equation by multiplying out the left side.

## Relation to Euler's formula
De Moivre's formula is a precursor to Euler's formula
\(e^{ix} = \cos x + i\sin x,\)
with  expressed in radians rather than degrees, which establishes the fundamental relationship between the trigonometric functions and the complex exponential function.

One can derive de Moivre's formula using Euler's formula and the exponential law for integer powers

\(\left( e^{ix} \right)^n = e^{inx},\)

since Euler's formula implies that the left side is equal to \(\left(\cos x + i\sin x\right)^n\) while the right side is equal to \(\cos nx + i\sin nx.\)

## Proof by induction
The truth of de Moivre's theorem can be established by using mathematical induction for natural numbers, and extended to all integers from there. For an integer , call the following statement S(*n*):

\((\cos x + i \sin x)^n = \cos nx + i \sin nx.\)

For *n* > 0, we proceed by mathematical induction. S(1) is clearly true. For our hypothesis, we assume S(*k*) is true for some natural . That is, we assume

\(\left(\cos x + i \sin x\right)^k = \cos kx + i \sin kx.\)

Now, considering S(*k* + 1):

\(\begin{alignat}{2}
 \left(\cos x+i\sin x\right)^{k+1} & = \left(\cos x+i\sin x\right)^{k} \left(\cos x+i\sin x\right)\\
 & = \left(\cos kx + i\sin kx \right) \left(\cos x+i\sin x\right) &&\qquad \text{via induction hypothesis}\\
 & = \cos kx \cos x - \sin kx \sin x + i \left(\cos kx \sin x + \sin kx \cos x\right)\\
 & = \cos ((k+1)x) + i\sin ((k+1)x) &&\qquad \text{via trigonometric identities}
\end{alignat}\)

See angle sum and difference identities.

We deduce that S(*k*) implies S(*k* + 1). By the principle of mathematical induction it follows that the result is true for all natural numbers. Now, S(0) is clearly true since cos(0*x*) + *i* sin(0*x*)  1}}. Finally, for the negative integer cases, we consider an exponent of −*n* for natural .

\(\begin{align}
 \left(\cos x + i\sin x\right)^{-n} & = \big( \left(\cos x + i\sin x\right)^n \big)^{-1} \\
 & = \left(\cos nx + i\sin nx\right)^{-1} \\
 & =  \cos nx - i\sin nx \qquad\qquad(*)\\
 & = \cos(-nx) + i\sin (-nx).\\
\end{align}\)
The equation (*) is a result of the identity
\(z^{-1} = \frac{\bar z}{|z|^2},\)
for *z* . Hence, S(*n*) holds for all integers .

## Formulae for cosine and sine individually

For an equality of complex numbers, one necessarily has equality both of the real parts and of the imaginary parts of both members of the equation. If , and therefore also cos *x* and sin *x*, are real numbers, then the identity of these parts can be written using binomial coefficients. This formula was given by 16th century French mathematician François Viète:

\(\begin{align}
\sin nx &= \sum_{k=0}^n \binom{n}{k} (\cos x)^k\,(\sin x)^{n-k}\,\sin\frac{(n-k)\pi}{2} \\
\cos nx &= \sum_{k=0}^n \binom{n}{k} (\cos x)^k\,(\sin x)^{n-k}\,\cos\frac{(n-k)\pi}{2}.
\end{align}\)

In each of these two equations, the final trigonometric function equals one or minus one or zero, thus removing half the entries in each of the sums. These equations are in fact valid even for complex values of , because both sides are entire (that is, holomorphic on the whole complex plane) functions of , and two such functions that coincide on the real axis necessarily coincide everywhere. Here are the concrete instances of these equations for *n*  and *n* :

\(\begin{alignat}{2}
 \cos 2x &= \left(\cos x\right)^2 +\left(\left(\cos x\right)^2-1\right)        &{}={}& 2\left(\cos x\right)^2-1       \\
 \sin 2x &= 2\left(\sin x\right)\left(\cos x\right)                            &     &                                \\
 \cos 3x &= \left(\cos x\right)^3 +3\cos x\left(\left(\cos x\right)^2-1\right) &{}={}& 4\left(\cos x\right)^3-3\cos x \\
 \sin 3x &= 3\left(\cos x\right)^2\left(\sin x\right)-\left(\sin x\right)^3    &{}={}& 3\sin x-4\left(\sin x\right)^3.
\end{alignat}\)

The right-hand side of the formula for cos *nx* is in fact the value *T*<sub>*n*</sub>(cos *x*) of the Chebyshev polynomial *T*<sub>*n*</sub> at cos *x*.

## Failure for non-integer powers, and generalization
De Moivre's formula does not hold for non-integer powers. The derivation of de Moivre's formula above involves a complex number raised to the integer power . If a complex number is raised to a non-integer power, the result is multiple-valued (see failure of power and logarithm identities).

### Roots of complex numbers
A modest extension of the version of de Moivre's formula given in this article can be used to find the -th roots of a complex number for a non-zero integer .  If  is a complex number, written in polar form as

\[
z=r\left(\cos x+i\sin x\right),
\]

then the -th roots of  are given by

\[
r^\frac1n \left( \cos \frac{x+2\pi k}{n} + i\sin \frac{x+2\pi k}{n} \right)
\]

where  varies over the integer values from 0 to |*n*| − 1.
This formula is also sometimes known as de Moivre's formula.

### Complex numbers raised to an arbitrary power
Generally, if \(z=r\left(\cos x+i\sin x\right)\) (in polar form) and  are arbitrary complex numbers, then the set of possible values is
\(z^w = r^w \left(\cos x + i\sin x\right)^w = \lbrace r^w \cos(xw + 2\pi kw) + i r^w \sin(xw + 2\pi kw) | k \in \mathbb{Z}\rbrace\,.\)
(Note that if  is a rational number that equals *p* / *q* in lowest terms then this set will have exactly  distinct values rather than infinitely many.  In particular, if  is an integer then the set will have exactly one value, as previously discussed.)  In contrast, de Moivre's formula gives
\(r^w (\cos xw + i\sin xw)\,,\)
which is just the single value from this set corresponding to 1=*k* = 0.

## Analogues in other settings
### Hyperbolic trigonometry
Since cosh *x* + sinh *x* , an analog to de Moivre's formula also applies to the hyperbolic trigonometry. For all integers ,

**\((\cosh x + \sinh x)^n = \cosh nx + \sinh nx.\)**

If  is a rational number (but not necessarily an integer), then cosh *nx* + sinh *nx* will be one of the values of (cosh *x* + sinh *x*)<sup>*n*</sup>.

### Extension to complex numbers
For any integer , the formula holds for any complex number \(z=x+iy\)
\(( \cos z + i \sin z)^n = \cos {nz} + i \sin {nz}.\)
where
\(\begin{align} \cos z = \cos(x + iy) &= \cos x \cosh y - i \sin x \sinh y\, , \\
\sin z = \sin(x + iy) &= \sin x \cosh y + i \cos x \sinh y\, . \end{align}\)

### Quaternions
To find the roots of a quaternion there is an analogous form of de Moivre's formula. A quaternion in the form
\(q = d + a\mathbf{\hat i} + b\mathbf{\hat j} + c\mathbf{\hat k}\)
can be represented in the form
\(q = k(\cos \theta + \varepsilon \sin \theta) \qquad \mbox{for } 0 \leq \theta < 2 \pi.\)
In this representation,
\(k = \sqrt{d^2 + a^2 + b^2 + c^2},\)
and the trigonometric functions are defined as
\(\cos \theta = \frac{d}{k} \quad \mbox{and} \quad \sin \theta = \pm \frac{\sqrt{a^2 + b^2 + c^2}}{k}.\)
In the case that *a*<sup>2</sup> + *b*<sup>2</sup> + *c*<sup>2</sup> ≠ 0,
\(\varepsilon = \pm \frac{a\mathbf{\hat i} + b\mathbf{\hat j} + c\mathbf{\hat k}}{\sqrt{a^2 + b^2 + c^2}},\)
that is, the unit vector. This leads to the variation of De Moivre's formula:

\(q^n = k^n(\cos n \theta + \varepsilon \sin n \theta).\)

#### Example
To find the cube roots of
\(Q = 1 + \mathbf{\hat i} + \mathbf{\hat j}+ \mathbf{\hat k},\)
write the quaternion in the form
\(Q = 2\left(\cos \frac{\pi}{3} + \varepsilon \sin \frac{\pi}{3}\right) \qquad \mbox{where } \varepsilon = \frac{\mathbf{\hat i} + \mathbf{\hat j}+ \mathbf{\hat k}}{\sqrt 3}.\)
Then the cube roots are given by:
\(\sqrt[3]{Q} = \sqrt[3]{2}(\cos \theta + \varepsilon \sin \theta) \qquad \mbox{for } \theta = \frac{\pi}{9}, \frac{7\pi}{9}, \frac{13\pi}{9}.\)

### 2 × 2 matrices
With matrices, \(\begin{pmatrix}\cos\phi & -\sin\phi \\ \sin\phi & \cos\phi \end{pmatrix}^n=\begin{pmatrix}\cos n\phi & -\sin n\phi \\ \sin n\phi & \cos n\phi \end{pmatrix}\) when  is an integer. This is a direct consequence of the isomorphism between the matrices of type \(\begin{pmatrix}a & -b \\ b & a \end{pmatrix}\) and the complex plane.

## References
*.


## External links
* De Moivre's Theorem for Trig Identities by Michael Croucher, Wolfram Demonstrations Project.

