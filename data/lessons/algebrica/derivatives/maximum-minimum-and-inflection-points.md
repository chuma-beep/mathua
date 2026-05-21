> Content sourced from [Algebrica](https://algebrica.org/maximum-minimum-and-inflection-points/) — CC BY-NC 4.0

## Global maximum and minimum points

The maximum and minimum of a function \\( f(x) \\) represent, respectively, the highest and lowest values that the function can attain within its [domain](<../determining-the-domain-of-a-function/>). In other words, they indicate the extreme points of the function, showing where \\( f(x) \\) reaches its greatest possible value or its smallest possible value for all permissible values of \\( x \\) in the given domain.

Given a function \\( y = f(x) \\) with domain \\( D \\), a point \\( x_0 \in D \\) is a **global maximum** if \\( f(x_0) \geq f(x) \\) for every \\( x \in D \\). The value \\( f(x_0) = M \\) is the global maximum of the function.

![Graph of a function f\(x\) showing a maximum point, where the curve reaches its highest value at a smooth peak.](/diagrams/algebrica/max-min-1.png)


Given a function \\( y = f(x) \\) with domain \\( D \\), a point \\( x_0 \in D \\) is a **global minimum** if \\( f(x_0) \leq f(x) \\) for every \\( x \in D \\). The value \\( f(x_0) = m \\) is the global minimum of the function.

![Graph of a function f\(x\) showing a minimum point, where the curve reaches its lowest value at a smooth valley.](/diagrams/algebrica/max-min-2.png)

If the global maximum and global minimum of a function exist, they are unique. By the [Weierstrass’s Theorem](<../weierstrass-theorem/>), if a function is [continuous](<../continuous-functions/>) on a closed and bounded interval \\([a, b]\\), then it attains both a global maximum and a global minimum on that interval.

## Local maximum and minimum points

In some cases, a function can display more than one peak or valley within a particular interval. Such points are known as local maxima and local minima. They correspond to positions where the function reaches a relatively highest or lowest value compared to its immediate surroundings, without necessarily being the absolute extremes over the entire domain.

Given a function \\( y = f(x) \\) defined on an interval \\([a, b]\\), the point \\( x_0 \in [a, b] \\) is a **local maximum** if there exists a neighborhood \\( I \\) of the point \\( x_0 \\) such that \\( f(x_0) \geq f(x) \\) for every \\( x \\) in the interval \\( I \\).

![Graph of a function f\(x\) showing a local maximum point, where the curve reaches a temporary highest value compared to nearby points.](/diagrams/algebrica/max-min-3-1.png)

In more formal terms, given a function \\( y = f(x) \\) that is defined and continuous in a neighborhood of the point \\( x_0 \\), and differentiable in the same neighborhood for every \\( x \neq x_0 \\), if for every \\( x \\) in the neighborhood the following conditions hold:

\\[\begin{align} f’(x) &> 0 \quad \text{for} \quad x < x_0 \\\ f’(x) &< 0 \quad \text{for} \quad x > x_0 \end{align} \\]

then, \\( x_0 \\) is a point of **local maximum** for the function \\( f(x) \\). We have:

| | \\[x_0 \\]  
---|---|---  
\\( f’(x) \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)  
\\( f(x) \\)| \\( \boldsymbol{\nearrow} \\)| \\( \boldsymbol{\searrow} \\)  
  

Given a function \\( y = f(x) \\) defined on an interval \\([a, b]\\), the point \\( x_0 \in [a, b] \\) is a **local minimum** if there exists a neighborhood \\( I \\) of the point \\( x_0 \\) such that \\( f(x_0) \leq f(x) \\) for every \\( x \\) in the interval \\( I \\).

![Graph of a function f\(x\) showing a local minimum point, where the curve reaches a temporary lowest value compared to nearby points.](/diagrams/algebrica/max-min-4.png)

If the following conditions hold:

\\[\begin{align} f’(x) &< 0 \quad \text{for} \quad x < x_0 \\\ f’(x) &> 0 \quad \text{for} \quad x > x_0 \end{align} \\]

then, \\( x_0 \\) is a point of **local minimum** for the function \\( f(x) \\). We have:

| | \\[x_0 \\]  
---|---|---  
\\( f’(x) \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\( f(x) \\)| \\( \boldsymbol{\searrow} \\)| \\( \boldsymbol{\nearrow} \\)  
  

A function can have at most one global maximum and at most one global minimum, but it can have multiple local maxima and local minima within its domain. By [Fermat’s theorem](<../fermat-theorem/>), the relative maximum and minimum points of a differentiable function, located within the domain of the function, are **stationary points**. This implies that the tangent line at a point of a relative maximum or minimum is parallel to the x-axis. In this case, the [derivative](<../derivatives/>) of the function at \\(x_0\\) is zero, and we have \\(f{\prime}(x_0) = 0\\).

## Upward and downward concavity

We say that the function \\( f(x) \\) is **concave upward** at \\( x_0 \\) if there exists a neighborhood \\( I \\) of \\( x_0 \\) such that, for every \\( x \in I \\) with \\( x \neq x_0 \\), the function \\( f(x) \\) takes values greater than those of the line \\( y = t(x) \\), which is the tangent line to the graph of \\( f(x) \\) at \\( x_0 \\).

\\[f(x) > t(x) \quad \forall x \in I - \lbrace x_0 \rbrace \\]

![](/diagrams/algebrica/max-min-5.png)


Similarly, we say that the function \\( f(x) \\) is **concave downward** at \\( x_0 \\) if there exists a neighborhood \\( I \\) of \\( x_0 \\) such that, for every \\( x \in I \\) with \\( x \neq x_0 \\), the function \\( f(x) \\) takes values less than those of the line \\( y = t(x) \\).

\\[f(x) < t(x) \quad \forall x \in I - \lbrace x_0 \rbrace \\]

![](/diagrams/algebrica/maximum-minimum-6-1.png)

The concepts of concavity and convexity are discussed in detail and in their analytical formulation in the entry [Convexity and Concavity of Functions](<../convexity-and-concavity-of-functions/>)

## Inflection points and change in concavity

An **inflection point** is a point where the concavity of a function changes.

Let us consider the case where a function \\( y = f(x) \\) is defined on an interval \\( (a, b) \\), and let \\( x_0 \in (a, b) \\) be either a point where \\( f(x) \\) is differentiable, or a point where:

\\[\lim_{x \to x_0} f’(x) = +\infty \quad \text{or} \quad \lim_{x \to x_0} f’(x) = -\infty \\]

The point \\( x_0 \\) is defined as an inflection point if the function changes concavity at \\( x_0 \\).

![An _inflection point_ is a point where the concavity of a function changes.](/diagrams/algebrica/maximum-minimum-7.png)

An inflection point is called horizontal if the tangent at the inflection point is parallel to the x-axis. When the tangent is parallel to the y-axis, the inflection point is called vertical. In all other cases, as in the case shown in the figure, it is called oblique.

![](/diagrams/algebrica/maximum-minimum-8.png)

\\( x_0 \\) is an **horizontal inflection point** for a function \\( f(x) \\) if \\( f’(x) = 0 \\) and the sign of \\( f’(x) \\) is the same\\(^1\\) for every \\( x \neq x_0 \\) in the neighborhood \\( I \\).

| | \\[x_0 \\]  
---|---|---  
\\( f’(x) \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\( f(x) \\)| \\( \boldsymbol{\nearrow} \\)| \\( \boldsymbol{\nearrow} \\)  
  
###### The signs in the neighborhood of \\(x_0\\) can be both positive (as in the scheme above) or both negative.

## How to calculate the points of local maximum and minimum

Given a continuous function, to find the local maximum and minimum points, we analyze the sign of the first derivative. The procedure involves the following steps:

  * Compute the derivative \\( f’(x) \\) and determine its domain to identify points where the function is not differentiable (e.g., cusps, corners).

  * Study the sign of the derivative by analyzing where \\( f’(x) \\) is positive, negative, or zero.

  * Identify local maxima and minima: a point \\( x_0 \\) is a local maximum if \\( f’(x) \\) changes from positive to negative around \\( x_0 \\). A point \\( x_0 \\) is a local minimum if \\( f’(x) \\) changes from negative to positive around \\( x_0 \\).


## Example 1

Let us calculate the local maximum and minimum points of the following function:

\\[y = f(x) = x^3 - \frac{1}{2}x^2 \\]


Being a [polynomial function](<../polynomial-function/>), it is continuous and differentiable for all \\( x \in \mathbb{R} \\). Therefore, it does not have any points of discontinuity within its domain. Let us now calculate the first derivative of the function. We obtain:

\\[f’(x) = 3x^2 - x\\]


Now, we study the sign of the derivative by imposing: \\[3x^2 - x > 0 \\]

Passing to the associated equation, we obtain:

\\[3x^2 - x = 0 \implies x(3x -1) = 0\\]

The equation is satisfied for \\( x = 0 \\) and \\( x = \dfrac{1}{3} \\).


Returning to the inequality, we obtain that \\( f’(x) > 0 \\) for \\( x < 0 \\) and \\( x > \dfrac{1}{3} \\).


Let us now represent the sign chart and observe that the function is increasing for \\( x < 0 \\), decreasing for \\( 0 < x < \dfrac{1}{3} \\), and increasing again for \\( x > \dfrac{1}{3} \\).

| | \\[0 \\]| \\[+\dfrac{1}{3} \\]  
---|---|---|---  
\\( f’(x) \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\( f(x) \\)| \\( \boldsymbol{\nearrow} \\)| \\( \boldsymbol{\searrow} \\)| \\( \boldsymbol{\nearrow} \\)  
  

For ( x = 0 ) the function takes the value \\( f(0) = 0^3 - \dfrac{1}{2}0^2 = 0 \\). The point \\( (0,0) \\) is therefore a local maximum.


For \\( x = \dfrac{1}{3} \\), the function takes the value: \\[\begin{align} f\left( \dfrac{1}{3} \right) &= \left( \dfrac{1}{3} \right)^3 - \dfrac{1}{2} \left( \dfrac{1}{3} \right)^2 \\\\[0.5em] &= \dfrac{1}{27} - \dfrac{1}{2} \times \dfrac{1}{9} \\\\[0.5em] &= -\dfrac{1}{54}\\\ \end{align} \\]

The point \\( \left( \dfrac{1}{3}, -\dfrac{1}{54} \right) \\) is therefore a local minimum. In this way, we have found the local maximum and minimum points of the function \\(f(x)\\).

## How to determine the concavity of a function

Let \\( y = f(x) \\) be a continuous function defined in a neighborhood of the point \\( x_0 \\), along with its first and second derivatives.

If at \\( x_0 \\) we have \\( f’'(x_0) \neq 0 \\), then:

  * The function is concave upward if \\( f’'(x_0) > 0 \\).
  * The function is concave downward if \\( f’'(x_0) < 0 \\).


## Example 2

Let us consider the function from Example 1 and determine its [convexity and concavity](<../convexity-and-concavity-of-functions/>):

\\[y = f(x) = x^3 - \frac{1}{2}x^2 \\]


The second derivative of the function is:

\\[f’'(x) = 6x - 1 \\]

Let us now study the sign by imposing:

\\[6x - 1 > 0 \implies x > \frac{1}{6} \\]


Let’s represent the sign chart, obtaining the intervals in which the function is concave upward or concave downward.

| | \\[0 \\]  
---|---|---  
\\( f’'(x) \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)  
\\( f(x) \\)| \\( \boldsymbol{\bigcup} \\)| \\( \boldsymbol{\bigcap} \\)  
Concavity| Upward| Downward  
  
In this way, we have obtained the intervals of concavity of the function.

## Identifying inflection points

An inflection point occurs when the concavity of a function changes sign. This change indicates a transition from a concave upward shape to a concave downward shape, or vice versa. To determine if a point is truly an inflection point, we need to verify if the second derivative \\(f{\prime}{\prime}(x)\\) changes sign as we pass through that point.

  * A point \\( x_0 \\) is a horizontal inflection point if: \\[f’(x_0) = 0, \quad f’'(x_0) = 0 \\] but the concavity changes sign in the neighborhood of \\( x_0 \\). In this case, the tangent line at \\( x_0 \\) is horizontal.

  * A point \\( x_0 \\) is a vertical inflection point if the function is not differentiable at \\( x_0 \\) and the concavity changes sign around \\( x_0 \\). This type of inflection point often occurs at points with sharp corners or cusps where the function is continuous but not smooth.

  * A point \\( x_0 \\) is an oblique inflection point if: \\[f’(x_0) \neq 0, \quad f’'(x_0) = 0 \\] and the concavity changes sign around \\( x_0 \\). In this case, the tangent line is neither horizontal nor vertical but has a non-zero slope.


## Exercises to find maxima, minima, and inflection points of functions

  * \\[\text{1. } \quad f(x) = x^3 - 6x^2 + 9x \\] [solution](<#>)

  * \\[\text{2. } \quad f(x) = \dfrac{x^2}{x^2 + 1} \\] [solution](<#>)

  * \\[\text{3. } \quad f(x) = \ln(x^2 + 1) \\] [solution](<#>)

  * \\[\text{4. } \quad f(x) = x e^{-x} \\] [solution](<#>)

  * \\[\text{5. } \quad f(x) = \sin(x) + \cos(x) \\] [solution](<#>)

  * \\[\text{5. } \quad f(x) = x^2 \ln(x) \\] [solution](<#>)


##### The proposed functions are carefully designed to help you consolidate your understanding of local maxima, minima, and inflection points. Each function requires you to compute the first and second derivatives, identify critical points, and analyze concavity changes. Some are more direct, while others involve algebraic manipulation or mixed expressions ([polynomial](<../polynomials>), [exponential](<../exponential-function/>), [logarithmic](<../logarithms/>), or [trigonometric](<../unit-circle/>)). Try to determine and classify all relevant points independently before checking the solutions.

Derivatives

The derivative describes a function’s rate of change.

9.1k

[Difference Quotient](https://algebrica.org/difference-quotient/)

7.8k

[Derivatives](https://algebrica.org/derivatives/)

1.6k

[Derivative of a Composite Function](https://algebrica.org/the-derivative-of-a-composite-function/)

3.3k

[Non-Differentiable Points](https://algebrica.org/points-of-non-differentiability/)

1k

[Differential of a Function](https://algebrica.org/differential-of-a-function/)

1.4k

[Derivative of Composite Power Functions](https://algebrica.org/derivative-of-composite-power-functions/)

1.4k

[Partial Derivatives](https://algebrica.org/partial-derivatives/)
