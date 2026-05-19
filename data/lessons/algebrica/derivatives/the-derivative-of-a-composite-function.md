> Content sourced from [Algebrica](https://algebrica.org/the-derivative-of-a-composite-function/) — CC BY-NC 4.0

## The chain rule

Let \\( g \\) be differentiable at \\( x \\), and let \\( f \\) be differentiable at \\( z = g(x) \\). Then the composite function \\( y = f(g(x)) \\) is differentiable at \\( x \\), and its [derivative](<../derivative/>) is the product of the derivative of \\( f \\) evaluated at \\( g(x) \\) and the derivative of \\( g \\) at \\( x \\):

\\[D[f(g(x))] = f’(g(x)) \cdot g’(x) \\]

This result is known as the chain rule. It states that to differentiate a composite function, one multiplies the derivative of the outer function, evaluated at the inner function, by the derivative of the inner function.

In Leibniz notation, if \\( y = f(u) \\) and \\( u = g(x) \\), the chain rule takes the form:

\\[\frac{dy}{dx} = \frac{dy}{du} \cdot \frac{du}{dx} \\]

## Proof

To prove that \\(D[f(g(x))] = f’(g(x)) \cdot g’(x)\\) we calculate the following [limit](<../limits/>):

\\[D[f(g(x))] = \lim_{h \to 0} \frac{f(g(x + h))-f(g(x))}{h} \\]


Let \\( z = g(x) \\), then \\( g(x+h)-g(x) = \Delta z. \\) This implies that \\( g(x+h) = g(x) + \Delta z. \\) The limit becomes:

\\[D[f(g(x))] = \lim_{h \to 0} \frac{f(z+\Delta z)-f(z)}{h} \\]


Multiplying both the numerator and the denominator by \\( \Delta z \\), we get:

\\[\begin{aligned} D[f(g(x))] &= \lim_{h \to 0} \frac{f(z + \Delta z)-f(z)}{\Delta z} \cdot \frac{\Delta z}{h} \\\\[0.5em] &= \lim_{h \to 0} \frac{f(z + \Delta z)-f(z)}{\Delta z} \cdot \frac{g(x + h)-g(x)}{h} \\\\[0.8em] &= f’(z) \cdot g’(x)\\\\[1em] &= f’(g(x)) \cdot g’(x) \end{aligned} \\]

This argument assumes \\( \Delta z \neq 0 \\) for \\( h \\) sufficiently small. A complete proof handles the case \\( \Delta z = 0 \\) separately via an auxiliary function; the conclusion is the same.


In the case of powers of a function, the rule generalizes as follows:

\\[D[f(x)^a] = a[f(x)]^{a-1}f’(x) \\]

## Example 1

Let’s compute the derivative of the following composite function:

\\[y = f(g(x)) = \sin(3x^2 + 2x) \\]

In this case, we have:

  * The inner function \\( g(x) = 3x^2 + 2x \\)
  * The outer function \\( f(t) = \sin(t) \\), where \\( t = g(x) = 3x^2 + 2x \\)


The outer function is \\( f(t) = \sin(t) \\). Its derivative is:

\\[f’(t) = \cos(t) \\]

Substituting \\( t = g(x) \\):

\\[f’(g(x)) = \cos(3x^2 + 2x) \\]


The inner function is \\( g(x) = 3x^2 + 2x \\). Its derivative is:

\\[g’(x) = 6x + 2 \\]

Applying the chain rule we obtain:

\\[D[f(g(x))] = f’(g(x)) \cdot g’(x) \\]

The result is:

\\[(6x + 2)\cos(3x^2 + 2x)\\]

###### Explore the case of [composite power functions](<derivative-of-composite-power-functions/>), specifically the calculation of the derivative of functions of the type: \\[D[f(x)^{g(x)}] \\]

## Extension to multiple compositions

The chain rule can be extended to compositions involving three or more functions. For example, given \\( y = f(g(h(x))) \\), the derivative is:

\\[D[f(g(h(x)))] = f’(g(h(x))) \cdot g’(h(x)) \cdot h’(x) \\]

Each factor represents the derivative of a function in the composition, evaluated at the composition of all subsequent functions. This pattern generalises to any finite number of nested functions. For \\( y = f_1(f_2(\cdots f_n(x)\cdots)) \\), the derivative is given by the product:

\\[f_1’(f_2(\cdots f_n(x)\cdots)) \cdot f_2’(f_3(\cdots f_n(x)\cdots)) \cdots f_{n-1}'(f_n(x)) \cdot f_n’(x) \\]

In practical applications, differentiation proceeds from the outermost function inward, with each derivative computed in sequence and the results multiplied together.


As an example, consider \\( y = \sin(e^{3x}) \\). The composition involves three functions:

\\[\begin{align} h(x) &= 3x \\\\[6pt] g(t) &= e^t \\\\[6pt] f(s) &= \sin(s) \end{align} \\]

Applying the chain rule from the outside inward we obtain:

\\[\begin{align} D[\sin(e^{3x})] &= \cos(e^{3x}) \cdot e^{3x} \cdot 3 \\\\[6pt] &= 3e^{3x}\cos(e^{3x}) \end{align} \\]

## Example 2

Consider the following function:

\\[y = \ln\left(e^{x^2} + 1\right) \\]

The composition involves three functions:

\\[\begin{align} h(x) &= x^2 \\\\[6pt] g(t) &= e^t + 1 \\\\[6pt] f(s) &= \ln(s) \end{align} \\]


The derivative of the outer function \\( f(s) = \ln(s) \\) is \\( f’(s) = \frac{1}{s} \\), evaluated at \\( s = g(h(x)) = e^{x^2} + 1 \\):

\\[f’(g(h(x))) = \frac{1}{e^{x^2} + 1} \\]

The derivative of the middle function \\( g(t) = e^t + 1 \\) is \\( g’(t) = e^t \\), evaluated at \\( t = h(x) = x^2 \\):

\\[g’(h(x)) = e^{x^2} \\]

The derivative of the inner function \\( h(x) = x^2 \\) is:

\\[h’(x) = 2x \\]


Applying the chain rule from the outside inward:

\\[\begin{align} D\left[\ln\left(e^{x^2} + 1\right)\right] &= \frac{1}{e^{x^2} + 1} \cdot e^{x^2} \cdot 2x \\\\[6pt] &= \frac{2x e^{x^2}}{e^{x^2} + 1} \end{align} \\]

The result is:

\\[\frac{2x e^{x^2}}{e^{x^2} + 1}\\]

## Selected references

  * **Harvard University, O. Knill**. [Chain Rule](https://people.math.harvard.edu/~knill/teaching/math1a2020/handouts/lecture10.pdf)
  * **MIT OpenCourseWare, G. Strang**. [Derivatives by the Chain Rule](https://ocw.mit.edu/courses/res-18-001-calculus-fall-2023/mitres_18_001_f17_ch04.pdf)
  * **University of Toronto, J. Campesato**. [Differentiability and the Chain Rule](https://www.math.toronto.edu/campesat/ens/1920/1114.pdf)


Derivatives

The derivative describes a function’s rate of change.

9.1k

[Difference Quotient](https://algebrica.org/difference-quotient/)

7.8k

[Derivatives](https://algebrica.org/derivatives/)

3.3k

[Non-Differentiable Points](https://algebrica.org/points-of-non-differentiability/)

1k

[Differential of a Function](https://algebrica.org/differential-of-a-function/)

1.4k

[Derivative of Composite Power Functions](https://algebrica.org/derivative-of-composite-power-functions/)

7.4k

[Maximum, Minimum, and Inflection Points](https://algebrica.org/maximum-minimum-and-inflection-points/)

1.4k

[Partial Derivatives](https://algebrica.org/partial-derivatives/)
