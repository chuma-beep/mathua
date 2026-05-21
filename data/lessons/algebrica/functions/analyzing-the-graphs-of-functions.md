> Content sourced from [Algebrica](https://algebrica.org/analyzing-the-graphs-of-functions/) — CC BY-NC 4.0

## Introduction

Analyzing the graph of a [function](<../functions>) \\(y=f(x)\\) allows us to analyze its behavior and key characteristics, providing valuable insights into its mathematical properties. It is a structured process that can be efficiently carried out by following a precise analytical framework, which consists of the following steps.

  * Determine the [domain](<../determining-the-domain-of-a-function/>) by identifying the set of real numbers where the function is defined.
  * Check for symmetry to see if the function is [even or odd](<../even-and-odd-functions/>).
  * Find intercepts by locating points where the graph crosses the \\(x\\)-axis and \\(y\\)-axis.
  * Analyze the sign of the function to determine where it is positive or negative.
  * Identify [asymptotes](<../asymptotes>), including horizontal, vertical, or oblique boundaries.
  * Use the first derivative to find increasing and decreasing intervals.
  * Apply the second derivative to study concavity and detect [inflection points](<../https://algebrica.org/maximum-minimum-and-inflection-points/>).
  * Provide a qualitative representation of the function’s graph.


## Domain

The first step in analyzing a function is to determine its [domain](<../determining-the-domain-of-a-function/>), the set of all real numbers for which the function is well-defined. The domain depends on the type of function being analyzed. Identifying the domain is crucial because it establishes the input values for which the function can be evaluated and graphed. To find the domain, we examine the function’s expression and identify any mathematical constraints that may restrict certain values of \\( x \\).

Let us consider the function \\( y = x^3 - 2x \\) as an example. This is a polynomial function of the form: \\[y = a_n x^n + a_{n-1} x^{n-1} + \dots + a_1 x + a_0\\] where \\( a_n, a_{n-1}, \dots, a_0 \\) are real coefficients, and \\( n \\) is the degree of the [polynomial](<../polynomials>).

Since polynomial functions are defined for all real numbers, the domain of this function is \\( \mathbb{R} \\).

## Symmetry

The following step is determining whether the function exhibits [symmetry](<../even-and-odd-functions/>) with respect to the \\( y \\)-axis or the origin.

  * A function \\(y=f(x)\\) is classified as [even](<../even-and-odd-functions/>) if it satisfies the condition \\(f(-x) = f(x) \quad \forall x \in D.\\) This implies that the graph is symmetric with respect to the \\( y \\)-axis.

  * A function \\(y=f(x)\\) is classified as [odd](<../even-and-odd-functions/>) if it satisfies the condition \\(f(-x) = -f(x) \quad \forall x \in D.\\) This implies that the graph is symmetric with respect to the origin.


To determine whether the function \\(f(x) = x^3 - 2x\\) is even or odd, we compute \\( f(-x) \\). Substituting ( -x ) in place of ( x ) we have: \\[\begin{align} f(-x) &= (-x)^3 - 2(-x) \\\\[0.5em] &= -x^3 + 2x \end{align} \\]


Calculating \\( -f(x) \\) we obtain:  
\\[\begin{align} -f(x) &= -(x^3 - 2x) \\\\[0.5em] &= -x^3 + 2x \end{align} \\]

Since \\( f(-x) = -f(x) \\), the function is odd, meaning it is symmetric with respect to the origin.

## Intersections with the cartesian axes

Once the symmetry has been studied, the next step is to determine the points where the function intersects the Cartesian axes. These include:

  * \\( x \\)-intercepts, obtained by solving \\( f(x) = 0 \\).
  * \\( y \\)-intercept, given by \\( f(0) \\) when the function is defined at \\( x = 0 \\).


To determine the intersection points with the \\( y \\)-axis, we set \\( x = 0 \\) in the function \\( f(x) \\). We have: \\[f(0) = 0^3 - 2(0) = 0\\]

Thus, the function intersects the \\( y \\)-axis at the point: \\[(0, 0)\\]


Next, to find the intersection points with the \\( x \\)-axis, we set \\( f(x) = 0 \\):

\\[x^3 - 2x = 0\\] \\[x(x^2 - 2) = 0\\]

Solving for \\( x \\) we have: \\[x = 0 \quad \text{or} \quad x^2 - 2 = 0\\]

From \\( x^2 - 2 = 0 \\), we obtain:  
\\[x = \pm \sqrt{2}\\]

Thus, the function intersects the \\( x \\)-axis at the points:

\\[(0,0), \quad (\sqrt{2},0), \quad (-\sqrt{2},0)\\]

Therefore, the intersections with the Cartesian axes are:

\\( y \\)-axis: \\( (0,0)\\)  
\\( x \\)-axis: \\( (0,0) \\), \\( (\sqrt{2},0) \\), \\( (-\sqrt{2},0) \\)

## Sign analysis

Next, we analyze the sign of the function, identifying the intervals where it is positive and negative. This is done by solving the inequality \\(f(x) > 0\\) which determines where the function takes positive values. The complementary intervals where \\( f(x) < 0 \\) indicate where the function is negative.

##### In this context, since we will be dealing with inequalities, it is useful to recall how to perform [sign analysis](<../sign-analysis-in-inequalities/>) in inequalities.

To analyze the sign of the function, we determine where \\( f(x) \\) is positive, negative, or zero by solving the inequality:  
\\[x^3 - 2x \gt 0\\]


Factoring the expression, we have: \\[x(x^2 - 2) \gt 0\\]

From the first factor we obtain: \\[x \gt 0\\]

From the second factor we obtain: \\[x^2-2 \gt 0 \Rightarrow x \gt \sqrt{2} \quad x \lt -\sqrt{2}\\]


By multiplying the signs of the first and second factors, we obtain (in black) the intervals where the function is positive.

| \\[-\sqrt{2}\\]| \\[0\\]| \\[+\sqrt{2}\\]|   
---|---|---|---|---  
| | | |   
| | | |   
| | | |   
| | | |   
  
Therefore, the function \\( x(x^2 - 2) \\) is positive for:

\\[x \in (-\sqrt{2},0) \cup (+\sqrt{2},+\infty)\\]

##### For the sake of completeness, we recall that the sign analysis of a function, as in the given example, requires examining the signs of its individual factors and determining the overall sign for each interval by computing the product of these signs.


We then represent on the Cartesian plane the intervals where the function must be located, excluding those in gray.

![](/diagrams/algebrica/graph-functions-1-2.png)

## Asymptotes

Another important step in function analysis is examining its behavior at the boundaries of the domain. This involves computing [limits](<../limits/>) to determine how the function behaves as \\( x \\) approaches the endpoints of its domain or extends toward infinity. By evaluating these limits, we can identify the presence of asymptotes.

  * A function \\( f(x) \\) has a horizontal asymptote if: \\[\lim\limits_{x \to \pm\infty} f(x) = L\\] where \\( L \\) is a finite real number. In this case, the line \\( y = L \\) represents the asymptote, describing the function’s end behavior.

  * A function \\( f(x) \\) has a vertical asymptote at \\( x = x_0 \\) if  
\\[\lim\limits_{x \to x_0^\pm} f(x) = \pm\infty\\] In this case, the line \\( x = a \\) represents the asymptote, indicating that the function grows unbounded near \\( x = a \\).

  * A function \\( f(x) \\) has an oblique asymptote of the form \\( y = mx + q \\) if the following limits exist and are finite:  
\\[m = \lim\limits_{x \to \pm\infty} \frac{f(x)}{x}\\] \\[q = \lim\limits_{x \to \pm\infty} \left[ f(x) - mx \right]\\]


To determine whether horizontal asymptotes exist, we verify whether the following limit exists and is finite: \\[\lim\limits_{x \to \pm\infty} f(x)\\]

We have:  
\\[\lim\limits_{x \to \pm\infty} (x^3 - 2x) = \pm\infty\\]

The function tends to infinity in both directions thus, there are no horizontal asymptotes.


Vertical asymptotes occur at points where a function is undefined and its values grow unbounded. The given function, \\(y = x^3 - 2x\\) is a polynomial, which is defined for all \\( x \in \mathbb{R} \\). Since polynomials do not have points of discontinuity or infinite limits at finite values of \\( x \\), we conclude that no vertical asymptotes exist.


To determine the existence of oblique asymptotes, we compute the slope \\(m\\) using the following limit:

\\[m = \lim\limits_{x \to \pm\infty} \frac{x^3 - 2x}{x} = +\infty\\]

Since the limit is not finite, oblique asymptotes do not exist.

Therefore, there are no horizontal, vertical, or oblique asymptotes.

## Maximum and minimum points

Now, by analyzing the first derivative, we first identify its domain and zeros, determining where \\( f^{\prime}(x) = 0 \\). By studying its sign, we establish the intervals where the function is increasing \\( f^{\prime}(x) > 0 \\) and consequently those where it is decreasing \\( f^{\prime}(x) < 0 \\).

Next, we identify possible [local maxima and minima](<../maximum-minimum-and-inflection-points/>) by evaluating the critical points of \\( f(x) \\). Additionally, we examine the function for points of inflection, where the concavity changes, and for points where \\( f(x) \\) is [not differentiable](<../points-of-non-differentiability/>).

We calculate the first derivative of \\( f(x) \\) and analyze its sign. \\[f^{\prime}(x) = 3x^2 - 2\\]

For \\[3x^2 - 2 > 0 \Rightarrow x < -\sqrt{\frac{2}{3}} \quad x > \sqrt{\frac{2}{3}} \\]


From this, it follows that \\(f(x)\\) is increasing for:

\\[x < -\sqrt{\frac{2}{3}} \quad \text{or} \quad x > \sqrt{\frac{2}{3}}\\]

and decreasing for:

\\[-\sqrt{\frac{2}{3}} < x < \sqrt{\frac{2}{3}}\\]


From the sign analysis, it follows that there is a local minimum at \\( \sqrt{\frac{2}{3}} \\). We now compute the function value at this point:

\\[\begin{align} f\left(\sqrt{\frac{2}{3}}\right) &= \left(\sqrt{\frac{2}{3}}\right)^3 - 2\left(\sqrt{\frac{2}{3}}\right) \\\\[0.5em] &= \frac{2\sqrt{2}}{3\sqrt{3}} - 2\sqrt{\frac{2}{3}} \\\\[0.5em] &= \frac{2\sqrt{6}}{9} - \frac{6\sqrt{6}}{9} \\\\[0.5em] &= \frac{-4\sqrt{6}}{9} \end{align} \\]


From the sign analysis, it also follows that there is a local maximum at ( -\sqrt{\frac{2}{3}} ). We now compute the function value at this point:

\\[\begin{align} f\left(-\sqrt{\frac{2}{3}}\right) &= \left(-\sqrt{\frac{2}{3}}\right)^3 - 2\left(-\sqrt{\frac{2}{3}}\right) \\\\[0.5em] &= -\frac{2\sqrt{2}}{3\sqrt{3}} + 2\sqrt{\frac{2}{3}} \\\\[0.5em] &= -\frac{2\sqrt{6}}{9} + \frac{6\sqrt{6}}{9} \\\\[0.5em] &= \frac{4\sqrt{6}}{9} \end{align} \\]

The function has:

  * A local maximum at:  
\\[x = -\sqrt{\frac{2}{3}}, \quad y = \frac{4\sqrt{6}}{9}\\]

  * A local minimum at:  
\\[x = \sqrt{\frac{2}{3}}, \quad y = \frac{-4\sqrt{6}}{9} \\]


## Inflection points

Finally, by analyzing the second derivative \\( f^{\prime\prime}(x) \\), we determine the intervals where the graph is concave up \\( f^{\prime\prime}(x) > 0 \\) or concave down \\( f^{\prime\prime}(x) < 0 \\).

Now, we identify the [inflection points](<../maximum-minimum-and-inflection-points/>) by analyzing the sign of the second derivative. The second derivative of f(x) is:

\\[f^{\prime\prime}(x) = 6x\\]

Setting \\( f’'(x) = 0 \\), we find the inflection point at \\(x = 0\\). Evaluating \\( f(0) \\), we obtain \\( y = 0 \\).

Thus, the inflection point is: \\[(0,0)\\]

## The final graph

At this point, we have all the necessary information to construct a qualitative-quantitative graph of our function, considering its behavior, possible asymptotes, local maxima and minima, and inflection points.

In the given example, we obtain the following graph:

![](/diagrams/algebrica/graph-functions-2.png)

In conclusion, studying the graph of a function requires a structured approach that involves identifying key properties such as domain, symmetry, intercepts, sign analysis, asymptotes, monotonicity, concavity, and critical points. Following these steps systematically ensures a precise and thorough understanding of the function’s behavior.

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
