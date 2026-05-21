> Content sourced from [Algebrica](https://algebrica.org/cotangent-function/) — CC BY-NC 4.0

## Cotangent function

The cotangent function \\( f(x) = \cot(x) \\) assigns to each angle \\( x \\), expressed in radians, its corresponding [cotangent](<../tangent-and-cotangent>) value. Its graph is a periodic curve with a period of \\( \pi \\) and features vertical [asymptotes](<../asymptotes>) where the sine of \\( x \\) equals zero, specifically at \\( x = k\pi \\) for \\( k \in \mathbb{Z} \\). The function \\( f(x) = \cot(x) \\) has a [domain](<../determining-the-domain-of-a-function/>) of all real numbers except these points, and its range is all real numbers.

![](/diagrams/algebrica/cotangent-chart-1.png)

  * Domain: \\( { x \in \mathbb{R} : x \neq k\pi \text{ for all } k \in \mathbb{Z} } \\)
  * Range: \\( y \in \mathbb{R} \\)
  * Periodicity: periodic in \\( x \\) with period \\( \pi \\)
  * Parity: odd, \\( \cot(-x) = -\cot(x) \\)


  * The cotangent of \\( x \\) is defined as the ratio between the [cosine and sine](<../sine-and-cosine>) of the angle \\( x \\). \\[\cot(x) = \frac{\cos(x)}{\sin(x)} \\]


  * Roots: \\( x = \frac{\pi}{2} + \pi n, \quad n \in \mathbb{Z} \\)
  * Fundamental root: \\( x = \frac{\pi}{2} \\)


  * Notable limits:

\\[\lim\limits_{x \to 0} x \cot(x) = 1 \\]

\\[\lim_{x\to0^+} \cot(x) = +\infty \quad \text{and} \quad \lim_{x\to0^-} \cot(x) = -\infty \\]


  * The function is [continuous](<../continuous-functions/>) and differentiable on its domain.
  * [Derivative](<../derivatives>): \\[\frac{d}{dx} \cot(x) = -\csc^2(x) \\]


  * [Indefinite integral](<../indefinite-integrals/>): \\[\int \cot(x) dx = \ln |\sin(x)| + c \\]


##### A comprehensive overview of trigonometric integrals, together with the most useful transformation and substitution techniques for handling more complex cases, is available in the page on [trigonometric function integrals](<../integral-of-trigonometric-functions/>).


  * An alternative form of the function \\( \cot(x) \\) using imaginary numbers is given by Euler’s formula. Here, \\( e^{ix} \\) is the [exponential function](<../exponential-function>) with base \\( e \\) and \\( i \\) is the [imaginary](<../complex-numbers>) unit. By expressing sine and cosine as \\[\sin(x) = \frac{e^{ix} - e^{-ix}}{2i} \quad \text{and} \quad \cos(x) = \frac{e^{ix} + e^{-ix}}{2} \\] we obtain the cotangent function as \\[\cot(x) = \frac{\cos(x)}{\sin(x)} = i \frac{e^{ix} + e^{-ix}}{e^{ix} - e^{-ix}}. \\]
