> Content sourced from [Algebrica](https://algebrica.org/tangent-function/) — CC BY-NC 4.0

## Tangent function

The tangent function \\(f(x) = \tan(x)\\) assigns to each angle \\(x\\), expressed in radians, its corresponding [tangent](<../tangent-and-cotangent>) value. Its graph is a periodic curve with a period of \\(\pi\\) and features vertical [asymptotes](<../asymptotes/>) where the cosine of \\(x\\) equals zero, specifically at \\(x = \pi/2 + k\pi\\) for \\(k \in \mathbb{Z}\\). The function \\(f(x) = \tan(x)\\) has a [domain](<../determining-the-domain-of-a-function/>) of all real numbers except these points, and its range is all real numbers.

![](https://algebrica.org/wp-content/uploads/resources/images/tangent-function.png)

A useful way to read this graph is to keep in mind that the tangent is defined as  
\\[\tan(x) = \frac{\sin(x)}{\cos(x)} \\] Thinking of it as a ratio helps make sense of the curve: the tangent varies gently where the underlying [sine and cosine](<../sine-and-cosine/>) change smoothly, while it rises or falls sharply as the cosine approaches zero, shaping the overall appearance of the graph.

###### The graph also shows that near the origin the tangent function behaves almost like a straight line: for small values of \\(x\\), \\( \tan(x) \\) increases smoothly before its growth becomes more pronounced as it approaches the discontinuities.

## Properties

  * [Domain](<../determining-the-domain-of-a-function/>): \\( { x \in \mathbb{R} : x \neq \frac{\pi}{2} + k\pi \text{ for all } k \in \mathbb{Z} } \\)
  * Range: \\( y \in \mathbb{R} \\)
  * Periodicity: periodic in \\( x \\) with period \\( \pi \\)
  * Parity: [odd](<../even-and-odd-functions/>), \\( \tan(-x) = -\tan(x) \\)
  * Roots: \\(x = \pi n, \qquad n \in \mathbb{Z}\\)
  * [Integer](<../integers/>) root: \\(x = 0\\)


## Limits, derivatives, and integrals of the tangent function

The tangent of \\( x \\) is defined as the ratio between the [sine and cosine](<../sine-and-cosine>) of the angle \\( x \\). \\[\tan(x) = \frac{\sin(x)}{\cos(x)} \\]


A useful limit to remember is: \\[\lim_{x \to 0} \frac{\tan(x)}{x} = 1\\] which shows that, near the origin, the tangent behaves almost like the function \\(x\\). The behaviour of the tangent near its first vertical asymptote is also well described by limits. As \\(x\\) approaches \\(\pi/2\\) from the left, the function grows without bound: \\[\lim_{x \to \frac{\pi}{2}^-} \tan(x) = +\infty \\] Coming from the right, the values instead diverge negatively: \\[\lim_{x \to \frac{\pi}{2}^+} \tan(x) = -\infty \\]


The function is [continuous](<../continuous-functions/>) and differentiable on its domain. The [derivative](<../derivatives>) is: \\[\frac{d}{dx} \tan(x) = \sec^2(x) \\]


The [indefinite integral](<../indefinite-integrals/>) is: \\[\int \tan(x) dx = -\ln |\cos(x)| + c \\]

##### A comprehensive overview of trigonometric integrals, together with the most useful transformation and substitution techniques for handling more complex cases, is available in the page on [trigonometric function integrals](<../integral-of-trigonometric-functions/>).


An alternative form of the function \\( \tan(x) \\) using imaginary numbers is given by Euler’s formula. Here, \\( e^{ix} \\) is the [exponential function](<../exponential-function>) with base \\( e \\) and \\( i \\) is the [imaginary](<../complex-numbers>) unit: \\[\tan(x) = \frac{e^{ix} - e^{-ix}}{i \left(e^{ix} + e^{-ix}\right)} \\]

Functions

A function maps each input to a unique output.

17.6k

[Functions](https://algebrica.org/functions/)

3 comments

10.2k

[Determining the Domain of a Function](https://algebrica.org/determining-the-domain-of-a-function/)

1.6k

[Even and Odd Functions](https://algebrica.org/even-and-odd-functions/)

2.7k

[Increasing, Decreasing and Monotonic Functions](https://algebrica.org/increasing-and-decreasing-functions/)

1.9k

[Convexity and Concavity of Functions](https://algebrica.org/convexity-and-concavity-of-functions/)

1.3k

[Composite Functions](https://algebrica.org/composite-functions/)

1.2k

[Inverse Function](https://algebrica.org/inverse-function/)

1.9k

[Continuous Functions](https://algebrica.org/continuous-functions/)

2 comments

1.1k

[Uniform Continuity](https://algebrica.org/uniform-continuity/)

1.3k

[Discontinuities of Real Functions](https://algebrica.org/discontinuities-of-real-functions/)

1.9k

[Analyzing the Graphs of Functions](https://algebrica.org/analyzing-the-graphs-of-functions/)

1.1k

[Polynomial Function](https://algebrica.org/polynomial-function/)

2.6k

[Rational Functions](https://algebrica.org/rational-functions/)

2k

[Logarithmic Function](https://algebrica.org/logarithmic-function/)

3.1k

[Exponential Function](https://algebrica.org/exponential-function/)

1.5k

[Absolute Value Function](https://algebrica.org/absolute-value-function/)

848

[Sign Function](https://algebrica.org/sign-function/)

2.2k

[Sine Function](https://algebrica.org/sine-function/)

2.1k

[Cosine Function](https://algebrica.org/cosine-function/)

2k

[Cotangent Function](https://algebrica.org/cotangent-function/)

1.8k

[Secant Function](https://algebrica.org/secant-function/)

1.4k

[Cosecant Function](https://algebrica.org/cosecant-function/)

1.2k

[Dirichlet Function](https://algebrica.org/dirichlet-function/)

1.4k

[Sigmoid Function](https://algebrica.org/sigmoid-function/)
