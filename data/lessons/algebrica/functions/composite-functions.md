> Content sourced from [Algebrica](https://algebrica.org/composite-functions/) — CC BY-NC 4.0

## What are composite functions

When we talk about composite functions, we refer to the process of applying one [function](<../functions>) to the result of another. In other words, given two functions \\( f(x) \\) and \\( g(x) \\), a composite function is formed by evaluating \\( g \\) at the output of \\( f \\). This is denoted by:

\\[g \circ f= g(f(x)) \\]

This means that we first apply \\( f \\) to the input \\( x \\), and then apply \\( g \\) to the result.

![](https://algebrica.org/wp-content/uploads/resources/images/composite-functions-1-1.png)

##### This diagram illustrates the concept of a composite function: the input \\( x \\) from set \\( A \\) is first mapped to \\( f(x) \\) in set \\( B \\), and then \\( f(x) \\) is mapped to \\( g(f(x)) \\) in set \\( C \\), resulting in the composition \\( g \circ f \\).


More formally, let two functions \\( f(x) \\) and \\( g(x) \\) be given such that:

  * \\(f \colon A \rightarrow B \\)
  * \\(g \colon B \rightarrow C \\)
  * \\(f(A) \subseteq B\\)


The composite function is defined as follows:

\\[g \circ f \colon x \in A \rightarrow g(f(x)) \in C \\]

This means that the function \\( g \circ f \\) maps each element \\( x \\) in the [domain](<../determining-the-domain-of-a-function/>) \\( A \\) to the value \\( g(f(x)) \\), provided that the image of \\( f \\) is contained in the domain of \\( g \\).

## Example

Consider the following functions:

  * \\( f(x) = 2x + 3 \\)
  * \\( g(x) = x^2 \\)


We wanto to define the composite function \\(g \circ f = g(f(x)) \\). Let’s start by evaluating \\( f(x) \\):

\\[f(x) = 2x + 3 \\]

Now plug that into \\( g(x) \\):

\\[g(f(x)) = g(2x + 3) = (2x + 3)^2 \\]

Therefore, the composite function is: \\[g \circ f = (2x + 3)^2 \\]

## Composition with the inverse function

If a function \\( f \\) is composed with its [inverse](<../inverse-function>) \\( f^{-1} \\), the result is the identity function, which maps each element of a set to itself:

\\[f(f^{-1}(x)) = f^{-1}(f(x)) = x \\]

##### This operation is valid only if the function \\( f \\) is invertible, meaning that it is both one-to-one (injective) and onto (surjective) over its domain.

When the composition between two functions is well-defined, that is, when the output of the first function lies within the domain of the second, we can write:

\\[g(f(x)) \equiv g \circ f \quad \text{and} \quad f(g(x)) \equiv f \circ g \\]

This notation highlights that function composition is **not commutative** : in general, the order in which functions are composed affects the outcome, and the following holds:

\\[g \circ f \neq f \circ g \\]

## Example

Let’s demonstrate with a simple example that function composition is not a commutative operation, that is, in general, \\( g \circ f \neq f \circ g \\). Consider the two functions:

  * \\( f(x) = e^x \\)
  * \\( g(x) = x + 1 \\)


Compute \\( f \circ g \\): \\[f \circ g = f(g(x)) = f(x + 1) = e^{x + 1} \\]

Compute \\( g \circ f \\): \\[g \circ f = g(f(x)) = g(e^x) = e^x + 1 \\]

We have:

  * \\( (f \circ g)(x) = e^{x + 1} = e \cdot e^x \\)
  * \\( (g \circ f)(x) = e^x + 1 \\)


These expressions are not equal and this proves that function composition is not commutative.

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
