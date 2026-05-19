> Content sourced from [Algebrica](https://algebrica.org/sign-function/) — CC BY-NC 4.0

## Introduction

The sign function assigns to each real number its sign, disregarding its magnitude. The function is defined as follows:

\\[\operatorname{sgn}(x) = \begin{cases} -1 & \text{if } x < 0 \\\\[6pt] 0 & \text{if } x = 0 \\\\[6pt] 1 & \text{if } x > 0 \end{cases} \quad \forall \, x \in \mathbb{R} \\]

Specifically, \\(\operatorname{sgn}(x)\\) returns \\(-1\\) for negative values, \\(0\\) when \\(x = 0\\), and \\(1\\) for positive values. The function does not quantify the magnitude of \\(x\\) but solely indicates the position of \\(x\\) relative to zero. For example, applying the definition we have:

\\[\operatorname{sgn}(-7) = -1 \qquad \operatorname{sgn}(0) = 0 \qquad \operatorname{sgn}(4) = 1 \\]


The graph of \\(y = \operatorname{sgn}(x)\\) consists of two horizontal rays and one isolated point. The ray on \\(y = -1\\) extends over all \\(x < 0\\), the ray on \\(y = 1\\) extends over all \\(x > 0\\), and the isolated point at the origin \\((0, 0)\\) lies on \\(y = 0\\). The two rays approach but do not intersect the y-axis.

![Sign function.](https://algebrica.org/wp-content/uploads/resources/images/sign-function-1.png)

The sign function is classified as an [odd function](<../even-and-odd-functions/>) because it satisfies the identity:

\\[\operatorname{sgn}(-x) = -\operatorname{sgn}(x) \quad \forall \, x \in \mathbb{R} \\]

## Properties

  * [Domain](<../determining-the-domain-of-a-function/>): \\(\mathbb{R}\\).
  * Range: \\({-1,\, 0,\, 1}\\).
  * The function is [odd](<../even-and-odd-functions/>), since \\(\operatorname{sgn}(-x) = -\operatorname{sgn}(x)\\).
  * The function has exactly one root at \\(x = 0\\), since \\(\operatorname{sgn}(x) = 0\\) only when \\(x = 0\\).
  * The function is constant on \\((-\infty, 0)\\) and on \\((0, +\infty)\\), hence [increasing](<../increasing-and-decreasing-functions/>) in the sense that it is non-decreasing over \\(\mathbb{R}\\).
  * The function has a [jump discontinuity](<../discontinuities-of-real-functions/>) at \\(x = 0\\); it is [continuous](<../continuous-functions/>) everywhere else.
  * The function is not differentiable at \\(x = 0\\). It is differentiable, with zero [derivative](<../derivatives/>), at every other point.
  * Limits approaching \\(x = 0\\) from either side: \\[\begin{align} \lim_{x \to 0^-} \operatorname{sgn}(x) &= -1 \\\\[0.5em] \lim_{x \to 0^+} \operatorname{sgn}(x) &= 1 \end{align} \\]


Since the two one-sided limits differ, the two-sided the following limit does not exist: \\[\lim_{x \to 0} \operatorname{sgn}(x)\\]

## Limits at infinity

As \\(x\\) moves away from the origin in either direction, the sign function stabilizes at a constant value:

\\[\begin{align} \lim_{x \to -\infty} \operatorname{sgn}(x) &= -1 \\\\[6pt] \lim_{x \to +\infty} \operatorname{sgn}(x) &= 1 \end{align} \\]

These are simply a consequence of the fact that \\(\operatorname{sgn}(x) = -1\\) for all \\(x < 0\\) and \\(\operatorname{sgn}(x) = 1\\) for all \\(x > 0\\), so the function value does not change as \\(x\\) moves further from zero.

## Derivative and integral

On each of the two open half-lines where the sign function is constant, its [derivative](<../derivatives/>) is zero:

\\[\frac{d}{dx} \operatorname{sgn}(x) = 0 \quad \text{for } x \neq 0 \\]

At \\(x = 0\\), the derivative does not exist because the function is discontinuous there. In the sense of distributions, however, the derivative of the sign function is:

\\[\frac{d}{dx} \operatorname{sgn}(x) = 2\delta(x) \\]

\\(\delta(x)\\) is the Dirac delta, a generalised function that is zero everywhere except at the origin and integrates to one over the entire real line. This distributional identity demonstrates that the sign function exhibits a discontinuity of amplitude \\(2\\) at the origin, as:

\\[\lim_{x \to 0^-} \operatorname{sgn}(x) = -1 \\] \\[\lim_{x \to 0^+} \operatorname{sgn}(x) = 1 \\]

This amplitude explains the presence of the factor \\(2\\) preceding the Dirac delta function.


The [indefinite integral](<../indefinite-integrals/>) of the sign function, computed away from the origin, gives back the [absolute value](<../absolute-value-function/>):

\\[\int \operatorname{sgn}(x) \, dx = |x| + c \\]

This is consistent with the fact that the derivative of \\(|x|\\) equals \\(\operatorname{sgn}(x)\\) wherever the former is differentiable.

## Relationship with the absolute value function

A direct algebraic relationship exists between the sign function and the [absolute value](<../absolute-value-function/>). For any \\(x \neq 0\\), the following identity holds:

\\[\operatorname{sgn}(x) = \frac{x}{|x|} \\]

This result is consistent with the definition: when \\(x > 0\\), the ratio \\(x / |x| = x/x = 1\\). When \\(x < 0\\), the ratio \\(x / |x| = x/(-x) = -1\\). The formula is undefined at \\(x = 0\\), so the value \\(\operatorname{sgn}(0) = 0\\) is assigned separately by convention.

Conversely, the absolute value can be expressed in terms of the sign function using the following identity:

\\[|x| = x \cdot \operatorname{sgn}(x) \\]

This identity holds for all \\(x \in \mathbb{R}\\), including \\(x = 0\\), where both sides are zero. Together, these two identities demonstrate that \\(|x|\\) and \\(\operatorname{sgn}(x)\\) are complementary: the absolute value preserves magnitude and omits sign, whereas the sign function preserves sign and omits magnitude.


Two additional identities arise from the relationship between the sign function and the absolute value. The first identity is as follows:

\\[x = |x| \cdot \operatorname{sgn}(x) \\]

This identity expresses any real number as the product of its magnitude and its sign. The second identity utilises the equality \\(|x| = \sqrt{x^2}\\), which is valid for all \\(x \in \mathbb{R}\\), and provides an alternative representation of the sign function for \\(x \neq 0\\):

\\[\operatorname{sgn}(x) = \frac{x}{\sqrt{x^2}} \\]

## Relationship with the Heaviside step function

The [Heaviside step function](<../heaviside-function/>) \\(H(x)\\) is defined as:

\\[H(x) = \begin{cases} 0 & \text{if } x < 0 \\\\[6pt] \dfrac{1}{2} & \text{if } x = 0 \\\\[8pt] 1 & \text{if } x > 0 \end{cases} \\]

![Heaviside step function.](https://algebrica.org/wp-content/uploads/resources/images/Heaviside-function-1.png)

The sign function and the Heaviside step function are related by a simple linear transformation. Specifically:

\\[\operatorname{sgn}(x) = 2H(x) - 1 \\]

Equivalently, we have:

\\[H(x) = \frac{1 + \operatorname{sgn}(x)}{2} \\]

This relationship is frequently utilised as converting between these representations can simplify calculations. The Heaviside function maps \\((-\infty, 0)\\) to \\(0\\) and \\((0, +\infty)\\) to \\(1\\), and may therefore be interpreted as a shifted and rescaled form of the sign function.

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

2.2k

[Sine Function](https://algebrica.org/sine-function/)

2k

[Cosine Function](https://algebrica.org/cosine-function/)

1.7k

[Tangent Function](https://algebrica.org/tangent-function/)

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
