> Content sourced from [Algebrica](https://algebrica.org/increasing-and-decreasing-functions/) — CC BY-NC 4.0

## Introduction

Understanding the behavior of [functions](<../functions>) is fundamental in mathematics. Depending on how their output values change with respect to the input, functions can be classified as:

  * increasing
  * decreasing
  * monotonic


##### These characteristics are closely tied to the geometry of their graphs in the Cartesian plane, revealing whether a function rises, falls, or maintains a consistent directional trend.


Let \\( y = f(x)\\) be a function defined on a [domain](<../determining-the-domain-of-a-function/>) \\( X \subseteq \mathbb{R} \\). We say that \\( f \\) is **strictly increasing** on an interval \\(I \subseteq X \\) if, for any two values \\( x_1, x_2 \in I \\) such that \\( x_1 < x_2 \\), the following condition holds: \\[f(x_1) < f(x_2) \\]

![](/diagrams/algebrica/increasing-decreasing-functions-1-3.png)

This means that as the input \\( x \\) increases within the interval \\( I \\), the output \\( f(x) \\) also strictly increases, without any flat or decreasing segments.


Let \\( y = f(x) \\) be a function defined on a domain \\( X \subseteq \mathbb{R} \\). We say that \\( f \\) is **strictly decreasing** on an interval \\( I \subseteq X \\) if, for any two values \\( x_1, x_2 \in I \\) such that \\( x_1 < x_2 \\), the following condition holds:

\\[f(x_1) > f(x_2) \\]

![](/diagrams/algebrica/increasing-decreasing-functions-2-1.png)

This means that as the input \\( x \\) increases within the interval \\( I \\), the output \\( f(x) \\) strictly decreases, with no flat or increasing sections.


A function with domain \\( X \subseteq \mathbb{R} \\) is said to be **strictly monotonic** on an interval \\( I \subseteq X \\) if it is either strictly increasing or strictly decreasing throughout the entire interval \\( I \\), with no change in direction or flat segments. In other words, the function maintains a consistent trend, either upward or downward, across \\( I \\).

To summarize: let \\( X \subseteq \mathbb{R} \\), and let \\( x_1, x_2 \in X \\) with \\( x_1 < x_2 \\). Then the function \\( f : X \to \mathbb{R} \\) is said to be:

  * Increasing: if \\( f(x_1) \leq f(x_2) \\).
  * Strictly increasing: if \\( f(x_1) < f(x_2) \\).
  * Decreasing: if \\( f(x_1) \geq f(x_2) \\).
  * Strictly decreasing: if \\( f(x_1) > f(x_2) \\).
  * (Strictly) monotonic: if the function is either (strictly) increasing or (strictly) decreasing.


## Derivatives and monotonic behavior

We know that [derivatives](<../derivatives>) are used to describe the shape and graph of functions. In particular, the first derivative of a function, \\( f’(x) \\), can indicate the intervals where the original function \\( f(x) \\) is increasing and where it is decreasing.

In general, given a function \\( y = f(x) \\) that is [continuous](<../continuous-functions/>) on an interval \\( I \\) and differentiable at the interior points of \\( I \\):

  * If \\( f’(x) > 0 \\) for every \\( x \\) in the interior of \\( I \\), then \\( f(x) \\) is increasing on \\( I \\).
  * If \\( f’(x) < 0 \\) for every \\( x \\) in the interior of \\( I \\), then \\( f(x) \\) is decreasing on \\( I \\).
  * If \\( f’(x) = 0 \\) for every \\( x \\) in the interior of \\( I \\), then \\( f(x) \\) is constant on \\( I \\).


To demonstrate these properties, we use the [Lagrange’s Theorem](<../lagrange-theorem>). Let’s imagine having two points \\( a \\) and \\( b \\) \\(\in I \\) with \\( a < b \\). Next, let us consider a point \\( c \\) belonging to the interval \\( ]a, b[ \\).

![](/diagrams/algebrica/increasing-decreasing-function.png)

By the Lagrange’s Theorem, we have:

\\[f^{\prime}\left (c \right ) = \frac{f(b) - f(a)}{b - a} \\]

Since we have \\( b - a > 0 \\) and \\( f’\left (c \right) > 0 \\), it follows that \\( f(b) - f(a) > 0 \\), which implies \\( f(b) > f(a) \\). Since \\( a \\) and \\( b \\) are arbitrary points in \\( I \\), the function is increasing on \\( I \\).

Similarly, considering the opposite case, since we have \\( b - a > 0 \\) and \\( f’\left( c \right) < 0 \\), it follows that \\( f(b) - f(a) < 0 \\), which implies \\( f(b) < f(a) \\). Since \\( a \\) and \\( b \\) are arbitrary points in \\( I \\), the function is decreasing on \\( I \\).

## Example 1

Let us consider the function: \\[f(x) = \frac{x^4}{4} - \frac{x^2}{2}\\]


Let us compute its derivative: \\[f’(x) = x(x^2 - 1)\\]

Let us find the intervals where the derivative is greater than zero. We have:

\\[x > 0\\] \\[x^2 - 1 > 0 \implies x < -1 \text{ or } x > 1\\]

By multiplying the signs of the first and second factors, we obtain the intervals where the derivative is positive.

| | \\[-1 \\]| \\[0 \\]| \\[1 \\]  
---|---|---|---|---  
\\( x > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\( x^2 - 1 > 0 \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\[f’(x)\\]| \\( \boldsymbol{-}\\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\(\boldsymbol{+} \\)  
  
Therefore, the derivative \\( x(x^2 - 1) \\) is positive for:

\\[x \in (-1,0) \cup (1,+\infty)\\]

##### For the sake of completeness, we recall that the sign analysis of a function, as in the given example, requires examining the signs of its individual factors and determining the overall sign for each interval by computing the product of these signs.


Graphically, its behavior is as follows:

![](/diagrams/algebrica/increasing-decreasing-functions-2.png)

Therefore, the function is increasing in the interval \\((-1,0) \cup (1,+\infty)\\) and decreasing in the interval \\((-\infty, -1) \cup (0,1)\\).

Functions

A function maps each input to a unique output.

17.6k

[Functions](https://algebrica.org/functions/)

3 comments

10.2k

[Determining the Domain of a Function](https://algebrica.org/determining-the-domain-of-a-function/)

1.6k

[Even and Odd Functions](https://algebrica.org/even-and-odd-functions/)

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
