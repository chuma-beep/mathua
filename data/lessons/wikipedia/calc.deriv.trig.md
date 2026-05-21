> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Differentiation_of_trigonometric_functions) — CC BY-SA 4.0

# Derivatives of sine and cosine

| Function | Derivative |
| --- | --- |
| \(\sin(x)\) | \(\cos(x)\) |
| \(\cos(x)\) | \(-\sin(x)\) |
| \(\tan(x)\) | \(\sec^2(x)\) |
| \(\cot(x)\) | \(-\csc^2(x)\) |
| \(\sec(x)\) | \(\sec(x)\tan(x)\) |
| \(\csc(x)\) | \(-\csc(x)\cot(x)\) |
| \(\arcsin(x)\) | \(\frac{1}{\sqrt{1-x^2}}\) |
| \(\arccos(x)\) | \(-\frac{1}{\sqrt{1-x^2}}\) |
| \(\arctan(x)\) | \(\frac{1}{x^2+1}\) |
| \(\arccot(x)\) | \(-\frac{1}{x^2+1}\) |
| \(\arcsec(x)\) | \(\frac{1}{|x|\sqrt{x^2-1}}\) |
| \(\arccsc(x)\) | \(-\frac{1}{|x|\sqrt{x^2-1}}\) |

The **differentiation of trigonometric functions** is the mathematical process of finding the derivative of a trigonometric function, or its rate of change with respect to a variable. For example, the derivative of the sine function is written (*a*) = cos(*a*), meaning that the rate of change of sin(*x*) at a particular angle *x = a* is given by the cosine of that angle.

All derivatives of circular trigonometric functions can be found from those of sin(*x*) and cos(*x*) by means of the quotient rule applied to functions such as tan(*x*) = sin(*x*)/cos(*x*). Knowing these derivatives, the derivatives of the inverse trigonometric functions are found using implicit differentiation.

## Proofs of derivatives of trigonometric functions
### Limit of sin(θ)/θ as θ tends to 0

The diagram at right shows a circle with centre *O* and radius *r =* 1. Let two radii *OA* and *OB* make an arc of θ radians. Since we are considering the limit as *θ* tends to zero, we may assume *θ* is a small positive number, say 0 < θ < π in the first quadrant.

In the diagram, let *R*\(_{1}\) be the triangle *OAB*, *R*\(_{2}\) the circular sector *OAB*, and *R*\(_{3}\) the triangle *OAC*.

The area of triangle *OAB* is:

\(\mathrm{Area}(R_1
)
=\tfrac{1}{2} \ |OA| \ |OB| \sin\theta = \tfrac{1}{2}\sin\theta \,. \)

The area of the circular sector *OAB* is:

\(\mathrm{Area}(R_2)
=\tfrac{1}{2}\theta \,. \)

The area of the triangle *OAC* is given by:

\(\mathrm{Area}(R_3
)
=\tfrac{1}{2} \ |OA| \ |AC| = \tfrac{1}{2} \tan\theta \,. \)

Since each region is contained in the next, one has:

\(\text{Area}(R_1) < \text{Area}(R_2) < \text{Area}(R_3) \implies
\tfrac{1}{2}\sin\theta < \tfrac{1}{2}\theta < \tfrac{1}{2}\tan\theta \,. \)

Moreover; since sin *θ* > 0 in the first quadrant, we may divide through by sin *θ*, giving:
\(1 < \frac{\theta}{\sin\theta} < \frac{1}{\cos\theta} \implies 1 > \frac{\sin\theta}{\theta} > \cos\theta \,. \)

In the last step we took the reciprocals of the three positive terms, reversing the inequities.

We conclude that for 0 < θ < π, the quantity sin(*θ*)/*θ* is *always* less than 1 and *always* greater than cos(θ). Thus, as *θ* gets closer to 0, sin(*θ*)/*θ* is "squeezed" between a ceiling at height 1 and a floor at height cos *θ*, which rises towards 1; hence sin(*θ*)/*θ* must tend to 1 as *θ* tends to 0 from the positive side:\(\lim_{\theta \to 0^+} \frac{\sin\theta}{\theta} = 1 \,. \)For the case where *θ* is a small negative number – π < θ < 0, we use the fact that sine is an odd function:

\(\lim_{\theta \to 0^-}\! \frac{\sin\theta}{\theta}
\ =\
\lim_{\theta\to 0^+}\!\frac{\sin(-\theta)}{-\theta}
\ =\
\lim_{\theta \to 0^+}\!\frac{-\sin\theta}{-\theta}
\ =\
\lim_{\theta\to 0^+}\!\frac{\sin\theta}{\theta} \ =\
 1 \,. \)

### Limit of (cos(θ)-1)/θ as θ tends to 0
The last section enables us to calculate this new limit relatively easily. This is done by employing a simple trick. In this calculation, the sign of *θ* is unimportant.
\(\lim_{\theta \to 0}\, \frac{\cos\theta - 1}{\theta}
\ =\
\lim_{\theta \to 0} \left( \frac{\cos\theta - 1}{\theta} \right)\!\! \left( \frac{\cos\theta + 1}{\cos\theta + 1} \right)
\ =\
\lim_{\theta \to 0}\, \frac{\cos^2\!\theta - 1}{\theta\,(\cos\theta + 1)}. \)

Using cos\(^{2}\)*θ* – 1 = –sin\(^{2}\)*θ*,
the fact that the limit of a product is the product of limits, and the limit result from the previous section, we find that:

\(\lim_{\theta \to 0}\,\frac{\cos\theta - 1}{\theta}
\ =\
\lim_{\theta \to 0}\, \frac{-\sin^2\theta}{\theta(\cos\theta+1)}
\ =\
\left( -\lim_{\theta \to 0} \frac{\sin\theta}{\theta}\right)\! \left( \lim_{\theta \to 0}\,\frac{\sin\theta}{\cos\theta + 1} \right)
\ =\
(-1)\left(\frac{0}{2}\right) = 0 \,. \)

### Limit of tan(θ)/θ as θ tends to 0
Using the limit for the sine function, the fact that the tangent function is odd, and the fact that the limit of a product is the product of limits, we find:
\(\lim_{\theta\to 0} \frac{\tan\theta}{\theta}
 \ =\
 \left(\lim_{\theta\to 0} \frac{\sin\theta}{\theta}\right)\!
 \left( \lim_{\theta\to 0} \frac{1}{\cos\theta}\right)
 \ =\
 (1)(1)
 \ =\
 1 \,. \)

### Derivative of the sine function
We calculate the derivative of the sine function from the limit definition:
\(\frac{d}{d\theta}\,\sin\theta = \lim_{\delta \to 0} \frac{\sin(\theta + \delta) - \sin \theta}{\delta}. \)

Using the angle addition formula sin α cos β + sin β cos α, we have:
\(\frac{d}{d\theta}\,\sin\theta
 =
\lim_{\delta \to 0} \frac{\sin\theta\cos\delta + \sin\delta\cos\theta-\sin\theta}{\delta}
 =
\lim_{\delta \to 0} \left( \frac{\sin\delta}{\delta} \cos\theta
+ \frac{\cos\delta -1}{\delta}\sin\theta \right). \)
Using the limits for the sine and for the cosine functions:
\(\frac{d}{d\theta}\,\sin\theta
 =
(1)\cos\theta + (0)\sin\theta
 =
\cos\theta \,. \)

### Derivative of the cosine function
#### From the definition of derivative
We again calculate the derivative of the cosine function from the limit definition:

\(\frac{d}{d\theta}\,\cos\theta
 =
\lim_{\delta \to 0} \frac{\cos(\theta+\delta)-\cos\theta}{\delta}. \)

Using the angle addition formula cos α cos β – sin α sin β, we have:
\(\frac{d}{d\theta}\,\cos\theta
 =
\lim_{\delta \to 0} \frac{\cos\theta\cos\delta - \sin\theta\sin\delta-\cos\theta}{\delta}
 =
\lim_{\delta \to 0} \left(\frac{\cos\delta -1}{\delta}\cos\theta \,-\, \frac{\sin\delta}{\delta} \sin\theta \right). \)
Using the limits for the sine and cosine functions:
\(\frac{d}{d\theta}\,\cos\theta
 = (0) \cos\theta - (1) \sin\theta = -\sin\theta \,. \)

#### From the chain rule
To compute the derivative of the cosine function from the chain rule, first observe the following three facts:
\(\cos\theta = \sin\left(\tfrac{\pi}{2}-\theta\right)\)
\(\sin\theta = \cos\left(\tfrac{\pi}{2}-\theta\right)\)
\(\frac{\operatorname{d}}{\operatorname{d}\!\theta}\sin\theta = \cos\theta\)
The first and the second are trigonometric identities, and the third is proven above. Using these three facts, we can write the following,
\(\frac{\operatorname{d}}{\operatorname{d}\!\theta}\cos\theta = \frac{\operatorname{d}}{\operatorname{d}\!\theta}\sin\left(\frac{\pi}{2}-\theta\right)\)
We can differentiate this using the chain rule. Letting \(f(x) = \sin x,\ \ g(\theta) =\tfrac{\pi}{2}-\theta\), we have:

\(\frac{\operatorname{d}}{\operatorname{d}\!\theta} f\!\left(g\!\left(\theta\right)\right) = f^\prime\!\left(g\!\left(\theta\right)\right) \cdot g^\prime\!\left(\theta\right) = \cos\left(\frac{\pi}{2}-\theta\right) \cdot (0-1) = -\sin\theta\).

Therefore, we have proven that
\(\frac{\operatorname{d}}{\operatorname{d}\!\theta}\cos\theta = -\sin\theta\).

### Derivative of the tangent function
#### From the definition of derivative
To calculate the derivative of the tangent function tan *θ*, we use first principles. By definition:
\(\frac{d}{d\theta}\,\tan\theta
 = \lim_{\delta \to 0} \left( \frac{\tan(\theta+\delta)-\tan\theta}{\delta} \right). \)
Using the well-known angle formula (tan α + tan β) / (1 - tan α tan β), we have:
\(\frac{d}{d\theta}\,\tan\theta
 = \lim_{\delta \to 0} \left[ \frac{\frac{\tan\theta + \tan\delta}{1 - \tan\theta\tan\delta} - \tan\theta}{\delta} \right]
 = \lim_{\delta \to 0} \left[ \frac{\tan\theta + \tan\delta - \tan\theta + \tan^2\theta\tan\delta}{\delta \left( 1 - \tan\theta\tan\delta \right)} \right]. \)
Using the fact that the limit of a product is the product of the limits:
\(\frac{d}{d\theta}\,\tan\theta
 = \lim_{\delta \to 0} \frac{\tan\delta}{\delta} \times \lim_{\delta \to 0} \left( \frac{1 + \tan^2\theta}{1 - \tan\theta\tan\delta} \right). \)
Using the limit for the tangent function, and the fact that tan *δ* tends to 0 as δ tends to 0:
\(\frac{d}{d\theta}\,\tan\theta
 = 1 \times \frac{1 + \tan^2\theta}{1 - 0} = 1 + \tan^2\theta. \)
We see immediately that:
\(\frac{d}{d\theta}\,\tan\theta
 = 1 + \frac{\sin^2\theta}{\cos^2\theta}
 = \frac{\cos^2\theta + \sin^2\theta}{\cos^2\theta}
 = \frac{1}{\cos^2\theta}
 = \sec^2\theta \,. \)

#### From the quotient rule
One can also compute the derivative of the tangent function using the quotient rule.
\(\frac{d}{d\theta} \tan\theta
 = \frac{d}{d\theta} \frac{\sin\theta}{\cos\theta}
 = \frac{\left(\sin\theta\right)^\prime \cdot \cos\theta - \sin\theta \cdot \left(\cos\theta\right)^\prime}{ \cos^2 \theta }
 = \frac{\cos^2 \theta + \sin^2 \theta}{\cos^2 \theta}\)
The numerator can be simplified to 1 by the Pythagorean identity, giving us,
\(\frac{1}{\cos^2 \theta} = \sec^2 \theta\)
Therefore,
\(\frac{d}{d\theta} \tan\theta = \sec^2 \theta\)

## Proofs of derivatives of inverse trigonometric functions
The following derivatives are found by setting a variable *y* equal to the inverse trigonometric function that we wish to take the derivative of. Using implicit differentiation and then solving for *dy*/*dx*, the derivative of the inverse function is found in terms of *y*. To convert *dy*/*dx* back into being in terms of *x*, we can draw a reference triangle on the unit circle, letting *θ* be y. Using the Pythagorean theorem and the definition of the regular trigonometric functions, we can finally express *dy*/*dx* in terms of *x*.

### Differentiating the inverse sine function
We let

\(y=\arcsin x\,\!\)

Where

\(-\frac{\pi}{2}\le y \le \frac{\pi}{2}\)

Then

\(\sin y=x\,\!\)

Taking the derivative with respect to \(x\) on both sides and solving for dy/dx:

\({d \over dx}\sin y={d \over dx}x\)

\(\cos y \cdot {dy \over dx} = 1\,\!\)

Substituting \(\cos y = \sqrt{1-\sin^2 y}\) in from above,

\(\sqrt{1-\sin^2 y} \cdot {dy \over dx} =1\)

Substituting \(x=\sin y\) in from above,

\(\sqrt{1-x^2} \cdot {dy \over dx} =1\)

\(\frac{dy}{dx}=\frac{1}{\sqrt{1-x^2}}\)

### Differentiating the inverse cosine function
We let

\(y=\arccos x\,\!\)

Where

\(0 \le y \le \pi\)

Then

\(\cos y=x\,\!}\)

Taking the derivative with respect to \(x}\) on both sides and solving for dy/dx:

\({d \over dx}\cos y={d \over dx}x\)

\(-\sin y \cdot {dy \over dx} =1\)

Substituting \(\sin y = \sqrt{1-\cos^2 y}\,\!\) in from above, we get

\(-\sqrt{1-\cos^2 y} \cdot {dy \over dx} =1\)

Substituting \(x=\cos y\,\!\) in from above, we get

\(-\sqrt{1-x^2} \cdot {dy \over dx} =1\)

\(\frac{dy}{dx} = -\frac{1}{\sqrt{1-x^2}}\)

Alternatively, once the derivative of \(\arcsin x\) is established, the derivative of \(\arccos x\) follows immediately by differentiating the identity \(\arcsin x+\arccos x=\pi/2\) so that \((\arccos x)'=-(\arcsin x)'}\).

### Differentiating the inverse tangent function
We let

\(y=\arctan x\,\!}\)

Where

\(-\frac{\pi}{2} < y < \frac{\pi}{2}\)

Then

\(\tan y=x\,\!\)

Taking the derivative with respect to \(x\) on both sides and solving for dy/dx:

\({d \over dx}\tan y={d \over dx}x\)

Left side:

\({d \over dx}\tan y
 = \sec^2 y \cdot {dy \over dx}
 = (1 + \tan^2 y) {dy \over dx}\) using the Pythagorean identity

Right side:

\({d \over dx}x = 1\)

Therefore,

\((1+\tan^2 y){dy \over dx}=1\)

Substituting \(x=\tan y\,\!\) in from above, we get

\((1+x^2){dy \over dx}=1\)

\({dy \over dx}=\frac{1}{1+x^2}\)

### Differentiating the inverse cotangent function
We let

\(y=\arccot x\)

where \(0<y<\pi\). Then

\(\cot y=x\)

Taking the derivative with respect to \(x\) on both sides and solving for dy/dx:

\(\frac{d}{dx}\cot y=\frac{d}{dx}x\)

Left side:

\({d \over dx}\cot y
 = -\csc^2 y \cdot {dy \over dx}
 = -(1 + \cot^2 y) {dy \over dx}\) using the Pythagorean identity

Right side:

\({d \over dx}x = 1\)

Therefore,

\(-(1+\cot^2y)\frac{dy}{dx}=1\)

Substituting \(x=\cot y\),

\(-(1+x^2)\frac{dy}{dx}=1\)

\(\frac{dy}{dx}=-\frac{1}{1+x^2}\)
Alternatively, as the derivative of \(\arctan x\) is derived as shown above, then using the identity \(\arctan x+\arccot x=\dfrac{\pi}{2}\) follows immediately that
\[
\begin{align}
\dfrac{d}{dx}\arccot x
&=\dfrac{d}{dx}\left(\dfrac{\pi}{2}-\arctan x\right)\\
&=-\dfrac{1}{1+x^2}
\end{align}
\]

### Differentiating the inverse secant function
#### Using implicit differentiation
Let

\(y = \arcsec x\ \mid |x| \geq 1\)

Then

\(x = \sec y \mid \ y \in \left [0,\frac{\pi}{2} \right )\cup \left (\frac{\pi}{2},\pi \right]\)

\(\frac{dx}{dy} = \sec y \tan y = |x|\sqrt{x^2-1}\)
(The absolute value in the expression is necessary as the product of secant and tangent in the interval of y is always nonnegative, while the radical \(\sqrt{x^2-1}\) is always nonnegative by definition of the principal square root, so the remaining factor must also be nonnegative, which is achieved by using the absolute value of x.)

\(\frac{dy}{dx} = \frac{1}{|x|\sqrt{x^2-1}}\)

#### Using the chain rule
Alternatively, the derivative of arcsecant may be derived from the derivative of arccosine using the chain rule.

Let

\(y = \arcsec x = \arccos \left(\frac{1}{x}\right)\)

Where

\(|x| \geq 1\) and \(y \in \left[0, \frac{\pi}{2}\right) \cup \left(\frac{\pi}{2}, \pi\right]\)

Then, applying the chain rule to \(\arccos \left(\frac{1}{x}\right)\):

\(\frac{dy}{dx} = -\frac{1}{\sqrt{1-\left(\frac{1}{x}\right)^2}} \cdot \left(-\frac{1}{x^2}\right)
 = \frac{1}{x^2\sqrt{1-\frac{1}{x^2}}}
 = \frac{1}{x^2} \cdot \frac{\sqrt{x^2-1}}{\sqrt{x^2}}
 = \frac{1}{\sqrt{x^2}\sqrt{x^2-1}}
 = \frac{1}{|x|\sqrt{x^2-1}}\)

### Differentiating the inverse cosecant function
#### Using implicit differentiation
Let

\(y = \arccsc x\ \mid |x| \geq 1\)

Then

\(x = \csc y\ \mid \ y \in \left [-\frac{\pi}{2},0 \right )\cup \left (0,\frac{\pi}{2} \right]\)

\(\frac{dx}{dy} = -\csc y \cot y = -|x|\sqrt{x^2-1}\)
(The absolute value in the expression is necessary as the product of cosecant and cotangent in the interval of y is always nonnegative, while the radical \(\sqrt{x^2-1}\) is always nonnegative by definition of the principal square root, so the remaining factor must also be nonnegative, which is achieved by using the absolute value of x.)
\(\frac{dy}{dx} = \frac{-1}{|x|\sqrt{x^2-1}}\))

#### Using the chain rule
Alternatively, the derivative of arccosecant may be derived from the derivative of arcsine using the chain rule.

Let

\(y = \arccsc x = \arcsin \left(\frac{1}{x}\right)\)

Where

\(|x| \geq 1\) and \(y \in \left[-\frac{\pi}{2}, 0\right) \cup \left(0, \frac{\pi}{2}\right]\)

Then, applying the chain rule to \(\arcsin \left(\frac{1}{x}\right)\):

\(\frac{dy}{dx} = \frac{1}{\sqrt{1-(\frac{1}{x})^2}} \cdot \left(-\frac{1}{x^2}\right)
 = -\frac{1}{x^2\sqrt{1-\frac{1}{x^2}}}
 = -\frac{1}{x^2} \cdot \frac{\sqrt{x^2-1}}{\sqrt{x^2}}
 = -\frac{1}{\sqrt{x^2}\sqrt{x^2-1}}
 = -\frac{1}{|x|\sqrt{x^2-1}}\)
