> Content sourced from [Algebrica](https://algebrica.org/secant-function/) — CC BY-NC 4.0

## Secant function

The [secant](<../secant-and-cosecant/>) function \\(f(x) = \sec(x)\\) is defined as the reciprocal of the [cosine function](<../cosine-function/>). For any real angle \\(x\\) (measured in radians), the secant takes the value:

\\[\sec(x) = \frac{1}{\cos(x)} \\]

as long as the \\(\cos(x) \neq 0)\\). This reciprocal relationship means that the behaviour of the secant function is entirely determined by the properties of the cosine function.

###### This section focuses on the analytical properties of the secant function. For a geometric interpretation based on the [unit circle](<../unit-circle/>), including how the secant arises from the extension of the radius and the corresponding right–triangle construction, see the dedicated entry.


Its graph is a periodic curve with period \\(2\pi\\). Because cosine reaches the value zero at isolated and regularly spaced points, the secant function exhibits vertical [asymptotes](<../asymptotes/>) at:

\\[x = \frac{\pi}{2} + k\pi \qquad k \in \mathbb{Z} \\]

where the reciprocal \\(1/\cos(x)\\) becomes undefined.

![Secant graph with asymptotic behaviour.](/diagrams/algebrica/secant-1-1.png)

These asymptotes separate the graph into distinct branches in which the function grows without bound as the angle approaches any of these points. The [domain](<../determining-the-domain-of-a-function/>) of \\(\sec(x)\\) is therefore the set of all real numbers except the points where \\(\cos(x)=0\\). Its range consists of the unbounded intervals \\((-\infty, -1] \cup [1, \infty)\\) reflecting the fact that the cosine function never takes values whose absolute value exceeds \\(1\\), making its reciprocal always greater than or equal to 1 in magnitude.

## Key properties

  * Domain: \\( { x \in \mathbb{R} : \cos(x) \neq 0 } = { x \in \mathbb{R} : x \neq \pi/2 + k\pi \text{ for all } k \in \mathbb{Z} } .\\)
  * Range: \\( y \in (-\infty, -1] \cup [1, \infty) .\\)
  * Periodicity: periodic in \\( x \\) with period \\( 2\pi .\\)
  * Parity: [even](<../even-and-odd-functions/>), \\( \sec(-x) = \sec(x) .\\)
  * The graph has vertical asymptotes at \\( x = \frac{\pi}{2} + k\pi .\\)


## Additional identity

There is a simple but meaningful relation that ties the secant and the [tangent](<../tangent-function/>) together. Starting from the [pythagorean identity](<../pythagorean-identity/>) for sine and cosine and rewriting everything in terms of cosine, we arrive at:

\\[\sec^{2}(x) = 1 + \tan^{2}(x) \\]

This identity shows how closely the two functions are linked: when the tangent becomes large, the secant grows as well, and both share the same vertical asymptotes. It is a handy relation that often appears in calculus, especially when dealing with derivatives, integrals, or trigonometric equations involving reciprocal functions.

## Limits, derivatives, and integrals of the secant function

Several limits help illustrate how the secant function behaves near key points of its domain. When the angle approaches values where the cosine is close to one, the secant remains bounded and approaches a finite value. As the angle nears those points at which the cosine tends to zero, the reciprocal grows without bound, giving rise to the vertical asymptotes characteristic of the function. These behaviours can be summarised through the following limits:

\\[1. \quad \lim_{x \to 0} \sec(x) = 1\\] \\[2. \quad \lim_{x \to \frac{\pi}{2}^-} \sec(x) = +\infty\\] \\[3. \quad \lim_{x \to \frac{\pi}{2}^+} \sec(x) = -\infty\\]


The secant function is [continuous](<../continuous-functions/>) and differentiable at every point where it is defined, that is, on the entire real line except at the angles where the cosine function vanishes. Within this domain it varies smoothly, and its rate of change follows from differentiating the reciprocal of the cosine. Using standard differentiation rules gives the derivative:

\\[4\. \quad \frac{d}{dx}\sec(x) = \sec(x)\tan(x) \\]

which expresses how the secant function grows or decreases depending on the combined behaviour of \\(\sec(x)\\) and \\(\tan(x)\\) at each point of its domain.


The antiderivative of the secant function can be obtained by a classical substitution that rewrites the integrand in a form suitable for logarithmic integration. This procedure leads to a compact expression involving both the secant and the tangent functions. The result is the following [indefinite integral](<../indefinite-integral/>):

\\[5\. \int \sec(x)\, dx = \ln\\!\left|\, \sec(x) + \tan(x) \,\right| + c \\]

##### A comprehensive overview of trigonometric integrals, together with the most useful transformation and substitution techniques for handling more complex cases, is available in the page on [trigonometric function integrals](<../integral-of-trigonometric-functions/>).


An alternative expression for the function \\(\sec(x)\\) can be obtained by rewriting the cosine in exponential form through [Euler’s](<../euler-number-limit-sequence/>) identity. This approach highlights the connection between trigonometric and complex exponential functions, and it often proves useful in contexts such as Fourier analysis or complex integration. Using the identity:

\\[6\. \quad \cos(x) = \frac{e^{ix} + e^{-ix}}{2} \\]

the secant function can be expressed as the reciprocal of this quantity, which yields:

\\[7\. \quad \sec(x) = \frac{2}{\,e^{ix} + e^{-ix}\,} \\]

This formulation emphasises the analytic structure of \\(\sec(x)\\) and provides a bridge between its trigonometric definition and its complex exponential representation.

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

1.7k

[Tangent Function](https://algebrica.org/tangent-function/)

2k

[Cotangent Function](https://algebrica.org/cotangent-function/)

1.4k

[Cosecant Function](https://algebrica.org/cosecant-function/)

1.2k

[Dirichlet Function](https://algebrica.org/dirichlet-function/)

1.4k

[Sigmoid Function](https://algebrica.org/sigmoid-function/)
