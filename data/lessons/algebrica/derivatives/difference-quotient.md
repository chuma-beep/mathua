> Content sourced from [Algebrica](https://algebrica.org/difference-quotient/) — CC BY-NC 4.0

## What is the difference quotient

Consider a function \\(y = f(x)\\) defined on the interval \\([a, b]\\), and two real numbers \\(c\\) and \\(c + h\\) with \\(h \neq 0\\), both lying within the interval \\([a, b]\\). The difference quotient of \\(f\\) at the point \\(c\\) is defined as the ratio:

\\[\frac{\Delta y}{\Delta x} = \frac{f(c+h)-f \left(c \right)}{h} \\]

The condition \\( h \neq 0 \\) is necessary: for \\( h = 0 \\) the points \\( A \\) and \\( B \\) coincide, the secant line is not defined, and the ratio reduces to \\( \frac{0}{0} \\). Consider the points \\(A\\) and \\(B\\) with:

  * \\(A(c, f \left(c \right))\\)
  * \\(B(c+h, f(c+h))\\)


the **difference quotient** of \\(f\\) at the point \\(c\\) is the slope of the line passing through \\(A\\) and \\(B\\).

![](https://algebrica.org/wp-content/uploads/resources/images/difference-quotient-2.png)

The difference quotient is fundamental to the definition of the [derivative](<../derivatives>). The derivative of a function at a point is the [limit](<../limits/>) of the difference quotient as \\(h\\) approaches zero.

\\[f’ \left(c \right) = \lim_{h \to 0} \frac{f(c+h)-f \left(c \right)}{h} \\]

This process, known as the limit of the difference quotient, provides the instantaneous rate of change of the function at that point, or equivalently, the slope of the tangent line to the graph of the function.


In general, the difference quotient measures how a function changes over a finite displacement. As the interval shrinks, it transitions from a global measure of variation to a local one.

  * The difference quotient provides an approximation of the rate of change. It is calculated over a finite interval \\( [x, x + \Delta x] \\) and represents the average rate of change.

  * The derivative provides the exact rate of change. It is calculated by taking the limit as \\( \Delta x \to 0 \\) and represents the instantaneous rate of change.


It is worth noting that the difference quotient and the [derivative](<../derivatives>) measure the same geometric quantity at different scales.

  * The difference quotient measures the slope of a [secant](<../secant-and-cosecant/>) line over a finite interval.
  * The derivative measures the slope of the [tangent](<../tangent-and-cotangent/>) line at a point.


## Alternative forms

  * \\[\text{1.} \quad \frac{f(c+h) - f(c)}{h} \\]

  * \\[\text{2.} \quad \frac{f(x_1) - f(x_0)}{x_1 - x_0} \\]

  * \\[\text{3.} \quad \frac{f(x + \Delta x) - f(x)}{\Delta x} \\]

  * \\[\text{4.} \quad \frac{f(x + dx) - f(x)}{dx} \\]


###### The expressions above represent the same quantity and differ only in how the two points are labeled, depending on the notation in use.

## Example 1

Let us calculate the difference quotient of the function \\(y = f(x) = 3x^2-x\\) at the point \\(c = 1\\) for a generic \\(h\\).

Determine \\(f(c+h) = f(1+h)\\):

\\[\begin{align*} f(1+h) &= 3(1+h)^2-(1+h) \\\\[0.4em] &= 3(1 + h^2 + 2h)-1-h \\\\[0.4em] &= 3 + 3h^2 + 6h-1-h \\\\[0.4em] &= 3h^2 + 5h + 2 \end{align*} \\]


Determine \\(f \left(c \right) = f(1)\\):  
\\[f(1) = 3(1)^2 - 1 = 3 - 1 = 2 \\]


Calculate the difference quotient: \\[\begin{align*} \frac{f(1+h)-f(1)}{h} &= \frac{(3h^2 + 5h + 2)-2}{h}\\\\[0.4em] &= \frac{3h^2 + 5h + 2-2}{h} \\\\[0.4em] &= \frac{3h^2 + 5h}{h}\\\\[0.4em] &= 3h + 5 \end{align*} \\]

The expression \\(3h + 5\\) represents, as \\(h\\) varies, the slope of a secant line passing through point \\(A\\) on the graph with an abscissa of 1.

## Example 2

Let us now consider a function that is not polynomial: \\( f(x) = \sqrt{x} \\),calculated at the point \\( c = 4 \\). The procedure is the same as before, but the simplification of the difference quotient requires an additional step. Determine \\( f(4+h) \\):

\\[f(4+h) = \sqrt{4+h} \\]

Determine \\( f(4) \\):

\\[f(4) = \sqrt{4} = 2 \\]

The difference quotient takes the form:

\\[\frac{f(4+h) - f(4)}{h} = \frac{\sqrt{4+h} - 2}{h} \\]

As it stands, this expression cannot be simplified directly since numerator and denominator share no obvious common factor. The standard approach is to rationalize the numerator by multiplying both numerator and denominator by the conjugate expression \\( \sqrt{4+h} + 2 \\):

\\[\begin{align*} \frac{\sqrt{4+h} - 2}{h} \cdot \frac{\sqrt{4+h} + 2}{\sqrt{4+h} + 2} &= \frac{(4+h) - 4}{h\left(\sqrt{4+h} + 2\right)} \\\\[0.4em] &= \frac{h}{h\left(\sqrt{4+h} + 2\right)} \\\\[0.4em] &= \frac{1}{\sqrt{4+h} + 2} \end{align*} \\]

The factor \\( h \\) cancels and the result is well-defined for every \\( h \neq 0 \\).

The expression \\( \dfrac{1}{\sqrt{4+h}+2} \\) represents the slope of the secant through \\( A = (4,\ 2) \\) and \\( B = (4+h,\ \sqrt{4+h}) \\) as \\( h \\) varies.

As \\( h \to 0 \\), this slope approaches \\( \dfrac{1}{4} \\) which is precisely the derivative of \\( \sqrt{x} \\) at \\( x = 4 \\).

## Selected references

  * **Purdue University, N. Egbert**. [The Difference Quotient](https://www.math.purdue.edu/~egbertn/fa2016/notes/lesson8.pdf)

  * **Dartmouth College**. [The Difference Quotient and the Derivative](https://math.dartmouth.edu/opencalc2/cole/lecture21.pdf)

  * **Princeton University**. [Differential and Integral Calculus](https://imai.fas.harvard.edu/teaching/files/calculus.pdf)

  * **Harvard University, O. Knill**. [Introduction to Calculus](https://people.math.harvard.edu/~knill/teaching/math1a_2012/handouts/math1a_2012.pdf)

  * **California State University San Marcos**. [Difference Quotient](https://www.csusm.edu/mathlab/documents/differencequotient-r6.pdf)


Derivatives

The derivative describes a function’s rate of change.

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

7.4k

[Maximum, Minimum, and Inflection Points](https://algebrica.org/maximum-minimum-and-inflection-points/)

1.4k

[Partial Derivatives](https://algebrica.org/partial-derivatives/)
