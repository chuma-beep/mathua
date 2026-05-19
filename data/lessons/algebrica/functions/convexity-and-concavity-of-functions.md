> Content sourced from [Algebrica](https://algebrica.org/convexity-and-concavity-of-functions/) — CC BY-NC 4.0

## Introduction

The study of a [function’s behaviour](<../analyzing-the-graphs-of-functions/>) involves not only determining where it [increases or decreases](<../increasing-and-decreasing-functions/>), but also understanding how its graph bends within an interval. A [function](<../functions/>) may be entirely increasing or decreasing on a region while still exhibiting distinct geometric shapes: the curve may arch upward or downward, affecting the overall appearance of the graph and the nature of its variation.

These qualitative features are described by the notions of **convexity** and **concavity**. From a geometric perspective, they characterise whether the graph lies above or below its tangents; analytically, they are determined through the sign of the second [derivative ](<../derivatives/>), which quantifies the curvature of the function. Establishing this connection provides a precise criterion for identifying how a function bends across its domain.

## Convexity

A function \\(f\\) is said to be convex on an interval when, for any two points \\(a\\) and \\(b\\) within that interval, the straight line that connects \\((a, f(a))\\) and \\((b, f(b))\\) remains above the graph of \\(f\\) at every point between \\(a\\) and \\(b\\). In other words, the chord linking these two points does not dip below the curve.

![](https://algebrica.org/wp-content/uploads/resources/images/convexity-1-1.png)

The equation of the secant line passing through the points \\((a, f(a))\\) and \\((b, f(b))\\) is

\\[h(x) = \frac{f(b)-f(a)}{b-a}(x-a) + f(a) \\]

For the function \\(f\\) to be convex, this line must lie above the graph of \\(f\\) for every \\(x\\) in the interval between \\(a\\) and \\(b\\). This requirement can be expressed by the inequality \\(h(x) > f(x)\\) which explicitly becomes:

\\[\frac{f(b)-f(a)}{b-a}(x-a) + f(a) > f(x) \\]

Rearranging the terms yields:

\\[\frac{f(b)-f(a)}{b-a}(x-a) > f(x) - f(a) \\]

and consequently:

\\[\frac{f(b)-f(a)}{b-a} > \frac{f(x)-f(a)}{x-a} \\]

##### The expression highlights how convexity forces the secant slope over the whole interval \\([a,b]\\) to dominate the slopes over any shorter subinterval starting at \\(a\\) (for example the line that connects \\((f, f(a)))\\) and \\((k, f(k))\\).

## Concavity

A function \\(f\\) is concave on an interval when, for any two points \\(a\\) and \\(b\\) in that interval, the line segment joining \\((a, f(a))\\) and \\((b, f(b))\\) lies below the graph of \\(f\\) for every point between \\(a\\) and \\(b\\). Equivalently, the chord connecting the two points never rises above the curve.

![The graph illustrates how the concavity of the function ensures that every secant line lies below the curve.](https://algebrica.org/wp-content/uploads/resources/images/convexity-1-2.png)

In an analogous way to the convex case, the equation of the secant line passing through the points ((a, f(a))) and ((b, f(b))) is

\\[h(x) = \frac{f(b)-f(a)}{b-a}(x-a) + f(a) \\]

For the function \\(f\\) to be concave, this line must lie _below_ the graph of \\(f\\) for every \\(x\\) in the interval between \\(a\\) and \\(b\\). This condition can be expressed by the inequality (h(x) < f(x)), which explicitly becomes:

\\[\frac{f(b)-f(a)}{b-a}(x-a) + f(a) < f(x) \\]

Rearranging the terms yields:

\\[\frac{f(b)-f(a)}{b-a}(x-a) < f(x) - f(a) \\]

and consequently:

\\[\frac{f(b)-f(a)}{b-a} < \frac{f(x)-f(a)}{x-a} \\]

##### The expression illustrates how concavity forces the secant slope over the entire interval \\([a,b]\\) to be smaller than the slopes over any shorter subinterval starting at \\(a\\) (such as the line connecting \\((a, f(a))\\) and \\((k, f(k))\\).

## Second derivative criteria for convexity and concavity

A direct way to determine the curvature of a function, in terms of concavity and convexity, is to examine its second derivative. If a function \\(f(x)\\) is differentiable, the behaviour of its second derivative \\(f’'(x)\\) identifies the intervals where the function is concave or convex. In other words, we have:

  * \\(f’'(x) > 0\\), the function \\(f(x)\\) is convex.
  * \\(f’'(x) < 0\\), the function \\(f(x)\\) is concave.
  * Points where \\(f’'(x) = 0\\) are candidates for changes in curvature, marking where the function may switch between concavity and convexity.


Let us consider, for example, a simple [polynomial function](<../polynomial-function/>):

\\[f(x) = x^{3} - 3x^{2} + 2x\\]

We examine the concavity and convexity of the function by analysing the behaviour of its second derivative. As a first step, we compute the first derivative of \\(f(x)\\), which, being a polynomial, can be obtained immediately:

\\[f’(x)=3x^2−6x+2\\]

By solving the associated [quadratic inequality](<../quadratic-inequalities/>), we identify the values for which the first derivative is positive or negative, and therefore the intervals on which \\(f(x)\\) is increasing or decreasing. For:

\\[f’(x)=3x^2−6x+2 \gt 0\\]

we obtain:

\\[x < \frac{3 - \sqrt{3}}{3} \quad \lor \quad x > \frac{3 + \sqrt{3}}{3}\\]

By plotting the solutions on the real line, we can visualize the intervals where the function \\(f(x)\\) is increasing or decreasing:

| | \\[\frac{3 - \sqrt{3}}{3} \\]| \\[\frac{3 + \sqrt{3}}{3} \\]  
---|---|---|---  
\\(f’(x) \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\(f(x) \\)| \\( \nearrow \\)| \\( \searrow \\)| \\( \nearrow \\)  
  
We now analyse the behaviour of \\(f(x)\\) with respect to its convexity and concavity. To do so, we compute the second derivative of \\(f(x)\\), which is given by:

\\[f’'(x) = 6x - 6 = 0 \\]

The equation is satisfied at \\(x = 1\\). For \\(x < 1\\), the second derivative is negative, which means that the function is concave on this interval. For \\(x > 1\\), the second derivative becomes positive, so the function is convex.

| | \\[1 \\]  
---|---|---  
\\(f’'(x) \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\(f(x) \\)| \\( \bigcap \\)| \\( \bigcup \\)  
Concavity| Downward| Upward  
  
At \\(x = 1\\), the second derivative becomes zero, which shows that the function has an [inflection point](<../maximum-minimum-and-inflection-points/>) at \\((1, 0)\\). In fact, we have:

\\[f(1)=13−3⋅12+2⋅1=0\\]

##### This confirms the effectiveness of the second derivative in describing changes in curvature and detecting inflection points.

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
