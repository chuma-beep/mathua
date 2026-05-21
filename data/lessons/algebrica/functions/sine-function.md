> Content sourced from [Algebrica](https://algebrica.org/sine-function/) — CC BY-NC 4.0

## Sine function

The sine function \\(f(x) = \sin(x) \\) assigns to each angle \\(x\\), expressed in radians, its corresponding [sine](<../sine-and-cosine>) value. Its graph is a periodic wave with a period of \\(2 \pi \\) and an amplitude of 1, oscillating between \\(-1\\) and \\(1\\). The function \\( f(x) = \cos x \\) has all real numbers in its [domain](<../determining-the-domain-of-a-function/>), but its range is \\( -1 \leq \cos(x) \leq 1 \\).

![](/diagrams/algebrica/sine-cosine-3-1.png)

##### Together with the [cosine function](<../cosine-function>), it represents one of the fundamental models of periodic waves, and is widely used to describe cyclic phenomena in physics, engineering, and mathematics. For example, in simple [harmonic motion](<../simple-harmonic-motion/>) in physics, the sine function typically appears in the equations for both displacement and [velocity](<../velocity>), describing the oscillatory behavior of systems like springs and pendulums.

## Properties

  * [Domain](<../determining-the-domain-of-a-function/>): \\(x \in \mathbb{R} \\)
  * Range: \\(y \in \mathbb{R} : -1 \leq y \leq\ 1 \\)
  * Periodicity: periodic in \\(x\\) with period \\( 2 \pi \\)
  * Parity: [odd](<../even-and-odd-functions/>), \\( \sin(-x) = -\sin(x)\\)
  * Roots: \\(x = \pi n, \quad n \in \mathbb{Z} \\)
  * [Integer](<../integers/>) root: \\(x=0\\)
  * [Maximum and minimum points](<../maximum-minimum-and-inflection-points/>): \\( \sin(x) \\) reaches its maximum \\(1\\) at \\( x = \dfrac{\pi}{2} + 2k \pi \\) with \\( k \in \mathbb{Z} \\) and its minimum \\(-1\\) at \\( x = \dfrac{3\pi}{2} + 2k \pi \\) with \\( k \in \mathbb{Z} \\).


## Limits, derivatives, and integrals of the cosine function

A fundamental [limit](<../limits/>) involving the sine function captures how \\(\sin(x)\\) behaves in a neighbourhood of the origin and plays a central role in differential calculus. As \\(x\\) approaches zero, the value of \\(\sin(x)\\) becomes increasingly close to \\(x\\) itself when both are measured in radians. This reflects the fact that, near the origin, the sine curve is almost indistinguishable from the line \\(y = x\\). This relationship is formalised through the following limit:

\\[\text{1.} \quad \lim_{x \to 0} \frac{\sin(x)}{x} = 1 \\]


The function \\(\sin(x)\\) is [continuous](<../continuous-functions/>) and differentiable for every real value of \\(x\\). Its behaviour is smooth and regular across the entire real line, with no points of discontinuity or non-differentiability. Moreover, since \\(\sin(x)\\) is differentiable everywhere, its [derivative](<../derivatives/>) is defined at all real numbers and is given by:

\\[2\. \quad \frac{d}{dx}\,\sin(x) = \cos(x) \\]


Since the derivative of \\( -\cos(x) \\) is \\( \sin(x) \\), the [indefinite integral](<../indefinite-integrals/>) of the sine function can be written as: \\[3. \quad \int \sin(x) dx = -\cos(x) + c\\]

##### A comprehensive overview of trigonometric integrals, together with the most useful transformation and substitution techniques for handling more complex cases, is available in the page on [trigonometric function integrals](<../integral-of-trigonometric-functions/>).


An alternative form of the function \\( \sin(x) \\) using imaginary numbers is given by Euler’s formula, where \\( e^{ix} \\) is the [exponential function](<../exponential-function>) with base \\( e \\) and \\( I \\) is the [imaginary](<../complex-numbers>) unit: \\[4. \quad \sin(x) = \frac{e^{ix} - e^{-ix}}{2i} \\]
