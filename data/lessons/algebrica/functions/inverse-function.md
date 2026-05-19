> Content sourced from [Algebrica](https://algebrica.org/inverse-function/) — CC BY-NC 4.0

## What is an inverse function

In the introduction to [functions](<../functions>), we saw that a function \\( f: X \to Y \\) is called bijective if it is both injective and surjective, that is, for every \\( y \in Y \\), there exists a unique \\( x \in X \\) such that \\( f(x) = y \\).

  * \\( X \\) is the domain.
  * \\( Y \\) is the codomain.
  * A function is called injective if for every \\( x_1, x_2 \in X \\), with \\( x_1 \ne x_2 \\), we have \\( f(x_1) \ne f(x_2) \\). In other words, for every \\( y \in Y \\), there exists at most one \\( x \in X \\) such that \\( f(x) = y \\).
  * A function is called surjective if for every \\( y \in Y \\), there exists at least one \\( x \in X \\) such that \\( f(x) = y \\).


A function \\( f : X \to Y \\) is bijective if and only if there exists a function \\( g : Y \to X \\) such that:

  * \\( (g \circ f)(x) = g(f(x)) = x \\) for every \\( x \in X \\)
  * \\( (f \circ g)(y) = f(g(y)) = y \\) for every \\( y \in Y \\)


In this case, the function \\( g \\) is unique and is called the **inverse function** of \\( f \\), denoted by: \\[f^{-1} = g \\]

##### \\( (g \circ f)(x) = g(f(x)) \\) is called the composite function, which means applying \\( f \\) to \\( x \\) first, and then applying \\( g \\) to the result.

## Making a function invertible by restricting its domain

Consider the function \\( f(x) = x^2 \\), defined on \\( \mathbb{R} \\). This is a quadratic function, represented by a [parabola](<../parabola>) with its vertex at the origin of the Cartesian plane. On its full [domain](<../determining-the-domain-of-a-function/>) \\( \mathbb{R} \\), the function is not invertible, since it is not injective: distinct inputs can produce the same output, for example \\( f(-2) = f(2) \\).

However, if we restrict the domain to \\( [0, +\infty) \\), the function becomes bijective and therefore invertible. In this case, the inverse function is:

\\[f(x) = x^2 \rightarrow f^{-1}(x) = \sqrt{x} \quad \text{for} \; x \geq 0 \\]

![](https://algebrica.org/wp-content/uploads/resources/images/inverse-function-1.png)

The graph of a function and that of its inverse are symmetric with respect to the line \\(y = x\\), which is the diagonal bisecting the first and third quadrants of the Cartesian plane.

If a function \\( f \\) is [composed](<../composite-functions/>) with its inverse \\( f^{-1} \\), the result is the **identity function** , which maps each element of a set to itself:

\\[f(f^{-1}(x)) = f^{-1}(f(x)) = x \\]

## How to find the inverse of a general function

  * Check whether the function is bijective, or determine a restriction of its domain that makes it bijective.

  * Replace \\( f(x) \\) with \\( y \\), so that you work with the equation \\( y = f(x) \\).

  * Swap the variables \\( x \\) and \\( y \\): write \\( x = f(y) \\). This reflects the idea of inverting input and output.

  * Solve the equation for \\( y \\), isolating it explicitly.

  * Rewrite the result as \\( f^{-1}(x) = \ldots \\), using \\( x \\) as the input variable for the inverse.


## Example

We want to find its inverse of the function \\( f(x) = \dfrac{2x - 1}{x + 3} \\).

##### The function \\( f \\) is bijective on its domain \\( \mathbb{R} \setminus {-3} \\), because it is strictly increasing: its derivative is always positive. This ensures that \\( f \\) is injective, and since the image of \\( f \\) covers all real numbers except a single point, it is also surjective onto its codomain.


Write the function as an equation: \\[y = \dfrac{2x - 1}{x + 3} \\]

Swap \\( x \\) and \\( y \\)  
\\[x = \dfrac{2y - 1}{y + 3} \\]


Solve for \\( y \\). Multiply both sides by \\( y + 3 \\):  
\\[x(y + 3) = 2y - 1 \\]

Distribute the left-hand side:  
\\[xy + 3x = 2y - 1 \\]


Bring all terms to one side and factor \\( y \\) on the left-hand side:  
\\[\begin{align} &xy - 2y = -1 - 3x\\\\[0.5em] &y(x - 2) = -1 - 3x \end{align} \\]

Solve for \\( y \\):  
\\[y = \dfrac{-1 - 3x}{x - 2} \\]

The inverse function is: \\[f^{-1}(x) = \dfrac{-1 - 3x}{x - 2} \\]

## Inverse function theorem

A useful result from basic analysis is the one–dimensional version of the inverse function theorem. The idea is quite intuitive: if a function behaves regularly on an interval, then it can be inverted without difficulty. More precisely, suppose a function \\(f\\) is continuous and differentiable on an interval \\(I\\), and its [derivative](<../derivatives>) never vanishes:

\\[f’(x) \neq 0 \quad \forall \, x \in I \\]

Under these conditions, the function is strictly monotonic on \\(I\\), which guarantees that it is invertible on that interval. As a consequence, an inverse function \\(f^{-1}\\) exists on \\(f(I)\\). This inverse is not only continuous but also differentiable, and its derivative is given by the relation:

\\[\bigl(f^{-1}\bigr)'(y) = \frac{1}{f’\\!\bigl(f^{-1}(y)\bigr)} \\]

This result shows how a local condition (the derivative never becomes zero) ensures a global property such as invertibility.

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
