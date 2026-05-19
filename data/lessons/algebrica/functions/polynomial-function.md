> Content sourced from [Algebrica](https://algebrica.org/polynomial-function/) — CC BY-NC 4.0

Concept

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Easy

2

Requires

0

Enables

The following concepts, [Functions](https://algebrica.org/functions/), [Polynomials](https://algebrica.org/polynomials/), are required as prerequisites for this entry.

## Introduction

A polynomial function is a [function](<../functions/>) that consists of [polynomials](<../polynomials/>) expressed in the following form:

\\[f(x) = a_n x^n + a_{n-1} x^{n-1} + \dotsb + a_1 x + a_0 \\]

  * \\( n \\) is a non-negative [integer](<../integers/>).
  * The coefficients \\( a_0, a_1, \ldots, a_n \\) are [real numbers](<../properties-of-real-numbers/>), with \\( a_n \neq 0 \\).
  * The integer \\( n \\) is the degree of the polynomial function
  * \\( a_n \\) is the leading coefficient. The term \\( a_0 \\) is the constant term and coincides with the value of the function at the origin, since \\( f(0) = a_0 \\).


###### Polynomial functions constitute the simplest and most structurally regular class of real functions. They are defined for every real number, they possess derivatives of all orders, and their graphs are smooth curves without corners, cusps, or discontinuities of any kind.

## Properties

The [domain](<../determining-the-domain-of-a-function/>) of any polynomial function is the entire real line \\( \mathbb{R} \\), since the expression \\( a_n x^n + \dotsb + a_0 \\) involves only addition and multiplication of real numbers, neither of which imposes any restriction on the input. The range depends on the degree and the sign of the leading coefficient \\(a_n\\), and it may be all of \\( \mathbb{R} \\) or an [interval](<../intervals/>) of the form \\( [m, +\infty) \\) or \\( (-\infty, M] \\).

The following properties hold for every polynomial function, regardless of its degree.

  * Domain: \\( \mathbb{R} \\).
  * A polynomial function is [continuous](<../continuous-functions/>) over all of \\( \mathbb{R} \\).
  * A polynomial function is differentiable over all of \\( \mathbb{R} \\), with [derivatives](<../derivatives/>) of every order.
  * A polynomial function has no [asymptotes](<../asymptotes/>) of any kind, since it is defined and finite for every finite value of \\( x \\), and diverges to infinity as \\( |x| \to +\infty \\).
  * The graph of a polynomial function has no corners, cusps, or [discontinuities](<../discontinuities-of-real-functions/>).


## Polynomial function explorer

This interactive graph illustrates how the shape of a polynomial function depends on its degree and the sign of its leading coefficient. By selecting different configurations, one can observe the corresponding changes in end behavior, symmetry, and the number of turning points, providing a direct visual interpretation of the underlying algebraic structure.


Degree

1 2 3 4 5

Leading coefficient

a > 0 a < 0

## Degree 1: linear functions

A polynomial function of degree 1 takes the form:

\\[f(x) = mx + q \\]

where \\( m \neq 0 \\) is the slope and \\( q \\) is the y-intercept. Its graph is a straight line. The function is [strictly increasing](<../increasing-and-decreasing-functions/>) when \\( m > 0 \\) and strictly decreasing when \\( m < 0 \\).

  * Domain: \\( \mathbb{R} \\).
  * Range: \\( \mathbb{R} \\).
  * Monotonicity: strictly monotone over \\( \mathbb{R} \\).
  * The function is bijective from \\( \mathbb{R} \\) to \\( \mathbb{R} \\).
  * It has no maximum or minimum points.


## Degree 2: quadratic functions

A polynomial function of degree 2 takes the form:

\\[f(x) = ax^2 + bx + c \\]

where \\( a \neq 0 \\). Its graph is a [parabola](<../parabola/>) with vertical axis of symmetry. The vertex of the parabola is located at:

\\[x_v = -\frac{b}{2a}\\] \\[\qquad y_v = f(x_v) = c-\frac{b^2}{4a}\\]

  * Domain: \\( \mathbb{R} \\).
  * Range: \\( \left[ y_v, +\infty \right) \\) if \\( a > 0 \\); \\( \left( -\infty, y_v \right] \\) if \\( a < 0 \\).
  * When \\( a > 0 \\), the parabola opens upward and the vertex is a global minimum.
  * When \\( a < 0 \\), the parabola opens downward and the vertex is a global maximum.
  * The function is not monotone over all of \\( \mathbb{R} \\), but it is strictly monotone on each of the two half-lines separated by the vertex.


## Degree 3: cubic functions

A polynomial function of degree 3 takes the form:

\\[f(x) = ax^3 + bx^2 + cx + d \\]

where \\( a \neq 0 \\). Unlike the quadratic case, a cubic function has no [global maximum or minimum](<../maximum-minimum-and-inflection-points/>).

  * Domain: \\( \mathbb{R} \\).
  * Range: \\( \mathbb{R} \\).
  * The function is bijective from \\( \mathbb{R} \\) to \\( \mathbb{R} \\) if and only if it has no local extrema, that is, if its derivative has no real roots.
  * It may have one or two local extrema and exactly one [inflection point](https://algebrica.org/maximum-minimum-and-inflection-points/).
  * Limits at infinity: \\[\begin{align} \lim_{x \to -\infty} f(x) &= -\infty \quad \text{if } a > 0 \\\\[6pt] \lim_{x \to +\infty} f(x) &= +\infty \quad \text{if } a > 0 \end{align} \\]


## End behavior

The end behavior of a polynomial function is determined exclusively by its leading term \\( a_n x^n \\). As \\( |x| \\) grows without bound, the contributions of all lower-degree terms become negligible in comparison, and the behavior of \\( f(x) \\) approaches that of the power function \\( a_n x^n \\). The result depends on two parameters: the parity of \\( n \\) and the sign of \\( a_n \\).

When the degree is even, \\( x^n \\) is non-negative for all \\( x \\), so both ends of the graph point in the same vertical direction. When the degree is odd, \\( x^n \\) changes sign with \\( x \\), and the two ends of the graph point in opposite directions. More precisely:

Degree \\( n \\)| \\( a_n \\)| \\( \lim_{x \to -\infty} f(x) \\)| \\( \lim_{x \to +\infty} f(x) \\)| \\(x \to -\infty\\)| \\(x \to +\infty\\)  
---|---|---|---|---|---  
even| \\( >0 \\)| \\( +\infty \\)| \\( +\infty \\)| \\(\nwarrow\\)| \\(\nearrow\\)  
even| \\( <0 \\)| \\( -\infty \\)| \\( -\infty \\)| \\(\swarrow\\)| \\(\searrow\\)  
odd| \\( >0 \\)| \\( -\infty \\)| \\( +\infty \\)| \\(\swarrow\\)| \\(\nearrow\\)  
odd| \\( <0 \\)| \\( +\infty \\)| \\( -\infty \\)| \\(\nwarrow\\)| \\(\searrow\\)  
  
###### This relationship between the leading term and the global shape of the graph is particularly useful when studying the behavior of functions, as it allows one to predict the qualitative appearance of the curve without performing a complete analysis.

## Symmetry

Polynomial functions may exhibit symmetry with respect to the y-axis or to the origin, depending on the parity of their terms.

  * A polynomial function is [even](<../even-and-odd-functions/>) if all of its terms have even degree, meaning that \\( f(-x) = f(x) \\) for all \\( x \\). In this case, the graph is symmetric with respect to the y-axis.
  * A polynomial function is odd if all of its terms have odd degree, meaning that \\( f(-x) = -f(x) \\) for all \\( x \\), and its graph is symmetric with respect to the origin.


Most polynomial functions are neither even nor odd, since they contain terms of mixed parity. For instance, \\( f(x) = x^3 + x^2 \\) is neither even nor odd.

## Roots and intersections with the axes

A root of a polynomial function is a value \\( \alpha \in \mathbb{R} \\) such that \\( f(\alpha) = 0 \\). Geometrically, the roots correspond to the points where the graph crosses or touches the x-axis. [The Fundamental Theorem of Algebra](<../roots-of-a-polynomial/>) guarantees that every non-constant polynomial has exactly \\( n \\) roots in \\( \mathbb{C} \\), counted with multiplicity. Over \\( \mathbb{R} \\), fewer real roots may be present, as some roots may be complex conjugate pairs.

A root \\( \alpha \\) has multiplicity \\( k \\) if \\( (x-\alpha)^k \\) divides \\( P(x) \\) but \\( (x-\alpha)^{k+1} \\) does not. The multiplicity determines the local behavior of the graph near the root.

  * If \\( k \\) is odd, the graph crosses the x-axis at \\( x = \alpha \\).
  * If \\( k \\) is even, the graph is tangent to the x-axis at \\( x = \alpha \\) and does not cross it.


The y-intercept is always \\( (0, a_0) \\), since \\( f(0) = a_0 \\).

## Derivative and integral of a polynomial function

The [derivative](<../derivatives/>) of a polynomial function is computed term by term by applying the power rule. For the general monomial \\( a_k x^k \\), the rule gives:

\\[\frac{d}{dx}\left( a_k x^k \right) = k\,a_k\,x^{k-1} \\]

Applying this to the full polynomial yields:

\\[f’(x) = n\,a_n x^{n-1} + (n-1)\,a_{n-1} x^{n-2} + \dotsb + a_1 \\]

The derivative is a polynomial of degree \\( n-1 \\), and the constant term disappears since the derivative of a constant is zero. A polynomial function is infinitely differentiable over \\( \mathbb{R} \\). The derivative of a polynomial function of degree \\( n \\) is itself a polynomial function of degree \\( n-1 \\), and this operation can be repeated until the zero polynomial is reached.


The [indefinite integral](<../indefinite-integrals/>) of a polynomial function is computed by applying the power rule for integration to each term:

\\[\int x^k \,dx = \frac{x^{k+1}}{k+1} + c \\]

Integrating the full polynomial gives:

\\[\int f(x) \,dx = \frac{a_n}{n+1}\,x^{n+1} + \frac{a_{n-1}}{n}\,x^n + \dotsb + a_1 x + a_0 x + c \\]

where \\( c \in \mathbb{R} \\) is the constant of integration. The result is a polynomial of degree \\( n+1 \\).

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

2.6k

[Rational Functions](https://algebrica.org/rational-functions/)

2k

[Logarithmic Function](https://algebrica.org/logarithmic-function/)

3.1k

[Exponential Function](https://algebrica.org/exponential-function/)

1.5k

[Absolute Value Function](https://algebrica.org/absolute-value-function/)

847

[Sign Function](https://algebrica.org/sign-function/)

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
