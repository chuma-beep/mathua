> Content sourced from [Algebrica](https://algebrica.org/derivative-of-composite-power-functions/) — CC BY-NC 4.0

## Composite Power Functions and Derivatives

We have previously introduced how to calculate the [derivative](<../derivatives>) of a function at a point using the definition of the [difference quotient](<../difference-quotient>). We also studied how to differentiate simple functions and [composite functions](<the-derivative-of-a-composite-function/>). Now, let’s see how to differentiate power functions of the form:

\\[D[f(x)]^{g(x)} \\]

To calculate the derivative of such a function, a combination of the [logarithmic](<../logarithms>) rule and the derivative of [exponential functions](<../exponential-function>) is used. The general formula for the derivative of \\(f(x)^g(x)\\), with \\(f\\) and \\(g\\) differentiable, is as follows:

\\[D[f(x)]^{g(x)} = f(x)^{g(x)} \left[ g’(x) \ln f(x) + g(x) \frac{f’(x)}{f(x)} \right ] \\]

Where:

  * \\( f(x)^{g(x)} \\) is the original function.
  * \\( f’(x) \\) is the derivative of \\( f(x) \\).
  * \\( \ln f(x) \\) is the natural logarithm of \\( f(x) \\).
  * \\( g’(x) \\) is the derivative of \\( g(x) \\).


## Example

Let’s consider the function \\( y = x^{2x} \\) as an example, and calculate its derivative.


First, let’s rewrite the function by applying the logarithm to both sides:

\\[\ln y = \ln(x^{2x})\\]

For the properties of logarithms \\(\log_a(b^c) = c \cdot \log_a(b)\\)

The equality can be rewritten as:

\\[\ln y = 2x \cdot \ln(x)\\]


Since \\(\ln y\\) is a composite function, its derivative is

\\[\frac{1}{y} \cdot y’\\]

Let’s compute the derivative for the element on the right-hand side of the equality \\(2x \cdot \ln(x)\\):

\\[2 \cdot \ln(x) + 2x \cdot \frac{1}{x}\\]

We obtain:

\\[\frac{1}{y} \cdot y’ = 2 \cdot \ln(x) + 2x \cdot \frac{1}{x}\\]


The equality can be rewritten as:

\\[y’ = y \cdot (2 \cdot \ln(x) + 2)\\]

Since \\(y = x^{2x}\\), we have:

\\[y’ = x^{2x} \cdot (2 \cdot \ln(x) + 2)\\]

Therefore, the derivative of \\( y = x^2 \\) is equal to:

\\[x^{2x} \cdot (2 \cdot \ln(x) + 2)\\]

## Test yourself

  * \\[\text{1. } \quad y = x^{2\cos(x)}\\] [solution](<../derivative-a-1>)

  * \\[\text{2. } \quad y = x^{\ln(x)}\\] [solution](<../derivative-a-2>)


##### The proposed functions are designed to help you consolidate your understanding of composite function derivatives. Try solving them independently before checking the solutions provided.

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

7.4k

[Maximum, Minimum, and Inflection Points](https://algebrica.org/maximum-minimum-and-inflection-points/)

1.4k

[Partial Derivatives](https://algebrica.org/partial-derivatives/)
