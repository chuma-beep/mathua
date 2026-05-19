> Content sourced from [Algebrica](https://algebrica.org/derivatives/) — CC BY-NC 4.0

## Introduction to derivatives

Consider a function \\( y = f(x) \\) defined on an interval \\( [a,b] \\). The derivative of \\( f \\) at a point \\( c \in (a,b) \\), denoted by \\( f’(c) \\), is defined, if the limit exists and is finite, as the limit of the [difference quotient](<../difference-quotient/>) as \\( h \to 0 \\):

\\[f’(c) = \lim_{h \to 0} \frac{f(c+h) - f(c)}{h} \\]

If this limit exists for every \\( x \\) in an interval, then the derivative defines a new function \\( f’(x) \\), called the derivative of \\( f \\). The value \\( f’(x) \\) represents the slope of the tangent line to the graph of \\( f \\) at the point \\( x \\).


When \\( h \to 0 \\), the point \\( B \\) approaches the point \\( A \\), and the [line](<../lines/>) \\( AB \\) becomes the tangent line to the curve at point \\( A \\). The slope of the tangent line at \\( A \\) is called the derivative of the function at point \\( c \\).

![](https://algebrica.org/wp-content/uploads/resources/images/derivatives-1.png)

With reference to the tangent at point \\(A\\) of the function \\( y = mx + q \\), the derivative \\( f’(x) \\) represents the value of the slope coefficient \\( m \\).


A function is differentiable at a point \\(c\\) if the derivative \\(f’ \left(c \right)\\) exists. If a function is differentiable:

  * The function is defined in a neighborhood of the point \\(c\\).
  * The limit of the [difference quotient](<../difference-quotient>), with respect to \\(c\\), exists and is finite as \\( h \to 0 \\).
  * The right-hand and left-hand limits of the difference quotient exist and are equal.


The inverse operation of differentiation is [integration](<../indefinite-integrals/>). This deep connection between derivatives and integrals is formalized by the [Fundamental Theorem of Calculus](<../fundamental-theorem-of-calculus/>), which establishes that differentiation and integration are inverse processes.

##### Derivatives play a fundamental role in physics as well. A classic example is [velocity](<../velocity>), which is defined as the derivative of position with respect to time. This simple concept forms the basis for understanding motion and change in the physical world.

If a function \\( f(x) \\) is differentiable at a point \\( c \\), then the function is also [continuous](<../continuous-functions/>) at that point. However, not all functions that are continuous at a point \\( c \\) are differentiable. In other words, differentiable functions form a subset of continuous functions.

## Example 1

Let’s calculate the derivative of the function \\( f(x) = 2x^2-3x \\) at \\( c = 2 \\) using the definition of the derivative.

\\[f’(2) = \lim_{h \to 0} \frac{f(2+h) - f(2)}{h} \\]


Let’s calculate the value of \\( f(2+h) \\):

\\[\begin{aligned} f(2+h) &= 2(2+h)^2 - 3(2+h) \\\\[0.5em] &= 2(4 + 4h + h^2)-3(2 + h) \\\\[0.5em] &= 8 + 8h + 2h^2-6-3h \\\\[0.5em] &= 2 + 5h + 2h^2 \end{aligned} \\]


Now, calculate the value of \\( f(2) \\):

\\[f(2) = 2(2^2)-3(2) = 8-6 = 2 \\]


Next, calculate the difference \\( f(2+h)-f(2) \\):

\\[f(2+h)-f(2) = (2 + 5h + 2h^2)-2 = 5h + 2h^2 \\]

Now, calculate the difference quotient:

\\[\frac{f(2+h)-f(2)}{h} = \frac{5h + 2h^2}{h} = 5 + 2h \\]


Finally, calculate the limit:

\\[\lim_{h \to 0} (5 + 2h) = 5\\]

Thus, the derivative of \\( f(x) = 2x^2 - 3x \\) at \\( c = 2 \\) is: \\[f’(2) = 5 \\]

## Right-hand and left-hand derivatives

Since the derivative is the limit of the difference quotient, as in the case of limits, it is possible to define the **right-hand** and **left-hand derivatives** of a function \\(y=f(x)\\).

The right-hand derivative is: \\[f_{+}^{\prime} \left(c \right) = \lim_{h \to 0^+} \frac{f(c+h)-f \left(c \right)}{h} \\]

The left-hand derivative is: \\[f_{-}^{\prime} \left(c \right) = \lim_{h \to 0^-} \frac{f(c+h)-f \left(c \right)}{h} \\]

A function is differentiable at a point \\(c\\) if the right-hand derivative and the left-hand derivative at the point exist and are equal to each other. More generally, a function \\( y = f(x) \\) is differentiable on an interval \\([A, B]\\) if it is differentiable at all interior points of the interval and if the right-hand derivative at point \\(A\\) and the left-hand derivative at point \\(B\\) exist and are finite.

## Fundamental derivatives

  * \\[f(x) = c \quad f’(x) = 0 \\]

  * \\[f(x) = x \quad f’(x) = 1 \\]

  * \\[f(x) = x^a \quad a \in \mathbb{R}, x > 0 \quad f’(x) = ax^{a-1} \\]

  * \\[f(x) = \sqrt{x} \quad x > 0 \quad f’(x) = \frac{1}{2\sqrt{x}} \\]

  * \\[f(x) = a^x \quad f’(x) = a^x \ln(a) \\]

  * \\[f(x) = \log_a(x) \quad f’(x) = \frac{1}{x \log(a)} \\]

  * \\[f(x) = \ln(x) \quad f’(x) = \frac{1}{x} \\]

  * \\[f(x) = e^x \quad f’(x) = e^x \\]

  * \\[f(x) = \sin(x) \quad f’(x) = \cos(x) \\]

  * \\[f(x) = \cos(x) \quad f’(x) = -\sin(x) \\]

  * \\[f(x) = \tan(x) \quad f’(x) = 1 + \tan^2(x) \\]

  * \\[f(x) = \cot(x) \quad f’(x) = -(1 + \cot^2(x)) \\]

  * \\[f(x) = \arcsin(x) \quad f’(x) = \frac{1}{\sqrt{1 - x^2}} \\]

  * \\[f(x) = \arccos(x) \quad f’(x) = \frac{-1}{\sqrt{1 - x^2}} \\]

  * \\[f(x) = \arctan(x) \quad f’(x) = \frac{1}{1 + x^2} \\]

  * \\[f(x) = \text{arccot}(x) \quad f’(x) = \frac{-1}{1 + x^2} \\]


##### \\( f(x) = c \\); Derivative: \\( f’(x) = 0 \\), since a line \\( y = c \\) is parallel to the x-axis, its slope \\( m \\) is equal to 0.

##### \\( f(x) = x \\); Derivative: \\( f’(x) = 1 \\). The function \\( y = x \\) is the bisector of the first and third quadrants, and its slope \\( m \\) is equal to 1.

##### \\( f(x) = x^a \\), \\( a \in \mathbb{R}, x > 0 \\); Derivative: \\( f’(x) = ax^{a-1} \\)

## Operations with derivatives

The derivative of the product of a constant \\( c \\) and a differentiable function \\( f(x) \\) is equal to the product of the constant and the derivative of the function. This is expressed as:

\\[D[c \cdot f(x)] = c \cdot f’(x) \\]

For example, if \\( c = 3 \\) and \\( f(x) = x^2 \\), then:

\\[D[3 \cdot x^2] = 3 \cdot f’[x^2] = 3 \cdot 2x = 6x \\]


The derivative of the sum of two functions \\( f(x) \\) and \\( g(x) \\) is equal to the sum of their derivatives. This is expressed as:

\\[D[f(x) + g(x)] = f’(x) + g’(x) \\]

For example, if \\( f(x) = x^2 \\) and \\( g(x) = 3x \\), then:

\\[D[x^2 + 3x] = f’(x^2) + g’(3x) = 2x + 3 \\]


The derivative of the product of two functions \\( f(x) \\) and \\( g(x) \\) is given by the product rule. This is expressed as:

\\[D[f(x) \cdot g(x)] = f’(x) \cdot g(x) + f(x) \cdot g’(x) \\]

For example, if \\( f(x) = x^2 \\) and \\( g(x) = 3x \\), then:

\\[\begin{aligned} D[x^2 \cdot 3x] &= f’(x^2) \cdot g(3x) + f(x^2) \cdot g’(3x)\\\\[0.5em] & = 2x \cdot 3x + x^2 \cdot 3 \\\\[0.5em] & = 6x^2 + 3x^2 = 9x^2\\\ \end{aligned} \\]


The derivative of the quotient of two functions \\( f(x) \\) and \\( g(x) \\), where \\( g(x) \neq 0 \\), is given by the quotient rule. This is expressed as:

\\[D\left[\frac{f(x)}{g(x)}\right] = \frac{f’(x) \cdot g(x)-f(x) \cdot g’(x)}{g^2(x)} \\]

For example, if \\( f(x) = x^2 \\) and \\( g(x) = 3x + 1 \\), then:

\\[\begin{aligned} D\left[\frac{x^2}{3x + 1}\right] &= \frac{2x \cdot (3x + 1)-x^2 \cdot 3}{(3x + 1)^2}\\\\[0.5em] &= \frac{6x^2 + 2x-3x^2}{(3x + 1)^2} \\\\[0.5em] &= \frac{3x^2 + 2x}{(3x + 1)^2}\\\\[0.5em] \end{aligned} \\]


The derivative of the reciprocal of a function \\( f(x) \\), where \\( f(x) \neq 0 \\), is given by:

\\[D\left[\frac{1}{f(x)}\right] = -\frac{f’(x)}{f^2(x)} \\]

For example, if \\( f(x) = 3x + 1 \\), then:

\\[D\left[\frac{1}{3x + 1}\right] = -\frac{3}{(3x + 1)^2} \\]


When differentiating a composition of two functions, these rules are not sufficient. In that case, it is necessary to apply the [derivative of a composite function](<../the-derivative-of-a-composite-function/>), also known as the chain rule.

## Higher-order derivatives

In general, the derivatives we have discussed so far are the first derivatives of a function \\( y = f(x) \\). The differentiation process can be iterated to compute the higher-order derivatives of a first derivative, such as the second derivative and the third derivative.

For example, let the function \\( y = f(x) = 3x^3 - 2x^2 + 1 \\).

  * The first derivative of the function is: \\[f’(x) = 9x^2 - 4x \\]

  * The second derivative is: \\[f’'(x) = 18x - 4 \\]

  * The third derivative \\[f^{\prime \prime \prime}(x) = 18 \\]


First and second derivatives play a fundamental role in analyzing the local behavior of functions, particularly in identifying [minimum and maximum points](<../maximum-minimum-and-inflection-points/>), as well as inflection points.

## Key theorems in differential calculus

Derivatives are also at the foundation of some of the most important theorems in differential calculus.

  * [Weierstrass’s Theorem](<../weierstrass-theorem/>) guarantees that a continuous function on a closed and bounded interval attains both its maximum and minimum values.
  * [Fermat’s Theorem](<../fermats-theorem/>) establishes a necessary condition for local extrema.
  * [Rolle’s Theorem](<../rolles-theorem/>) and [Lagrange’s Theorem](<../lagranges-theorem/>), also known as the Mean Value Theorem, describe fundamental properties of differentiable functions on a closed interval.
  * Further generalizations are provided by [Cauchy’s Theorem](<../cauchy-theorem//>) and [L’Hôpital’s Rule](<../hopital-rule//>), which extends the use of derivatives to the evaluation of [indeterminate forms](<../indeterminate-forms/>) of limits.


## Equation of the tangent line

The slope of the tangent line to the graph of a function \\( f(x) \\) at a point \\( x_0 \\) is given by the derivative \\( f’(x_0) \\). This value represents the instantaneous rate of change of the function at that point and coincides with the limit of the slopes of the secant lines approaching \\( x_0 \\).

More precisely, if we consider a second point \\( x_0 + h \\), the slope of the secant line through the points \\( (x_0, f(x_0)) \\) and \\( (x_0 + h, f(x_0 + h)) \\) is:

\\[\frac{f(x_0 + h) - f(x_0)}{h} \\]

If the limit of this expression exists as \\( h \to 0 \\), the function is differentiable at \\( x_0 \\), and this limit equals \\( f’(x_0) \\). The tangent line is therefore understood as the limiting position of the secant lines.

If the derivative exists and is finite, the tangent line is not vertical. In that case, its equation can be written in point–slope form. Since the line passes through the point \\( (x_0, f(x_0)) \\) and has slope \\( f’(x_0) \\), its equation is:

\\[y - f(x_0) = f’(x_0)(x - x_0) \\]

This linear function provides the best linear approximation of \\( f \\) near \\( x_0 \\). In fact, for values of \\( x \\) close to \\( x_0 \\), the increment of the function satisfies:

\\[f(x) \approx f(x_0) + f’(x_0)(x - x_0) \\]

which expresses the idea that, at sufficiently small scales, a differentiable function behaves approximately like its tangent line.

## Example 2

Let us consider the [parabola](<../parabola>) defined by the equation \\( y = 2x^2 + 3x \\), and determine the tangent line at the point \\( P(1, 5) \\).


First, we compute the derivative \\( f’(x) \\), and we obtain:

\\[2x+3\\]

We compute the slope of the tangent line at \\( x = 1 \\):

\\[f’(1) = 2(1) + 3 = 5 \\]

Therefore, the slope of the tangent line is \\( m = 5 \\). We have:

\\[y - f(1) = f{\prime}(1)(x - 1) \rightarrow y - 5 = 5(x - 1)\\]

Completing the calculations, we obtain the equation of the tangent line:

\\[y = 5x \\]

## Partial derivatives

Derivatives are fundamental in multivariable calculus. When a function depends on multiple variables, the rate of change with respect to a single variable, while keeping all other variables constant, is described by the partial derivative:

\\[\frac{\partial f}{\partial x_i}(x_0) \;=\; \lim_{h \to 0} \frac{f(x_1^0, \ldots, x_i^0 + h, \ldots, x_n^0) - f(x_0)}{h} \\]

The [vector](<../vectors/>) of all partial derivatives forms the gradient \\( \nabla f \\), which indicates the direction of the steepest increase of \\( f \\). For a comprehensive discussion, including higher-order derivatives, Schwarz’s theorem, the Jacobian matrix, and the chain rule for functions of several variables, refer to the entry on [partial derivatives](<../partial-derivatives>).

## Selected references

  * **MIT OpenCourseWare**. [Single Variable Calculus – Lecture Notes (Derivatives)](https://ocw.mit.edu/courses/18-01-single-variable-calculus-fall-2006/pages/lecture-notes/)

  * **University of British Columbia**. [Differentiation – Lecture Notes](https://www.math.ubc.ca/~feldman/m101/clp/clp_notes_100.pdf)

  * **University College London (UCL)**. [On Differentiation I](https://www.ucl.ac.uk/~uczlcfe/my%20website%20notes/_maths%20-%20diff%201%20-%2000%20-%20complete.pdf)

  * **Simon Fraser University**. [Calculus III – Partial Derivatives](https://www.sfu.ca/~vjungic/Calculus%203/Calculus3.pdf)

  * **Portland State University**. [Calculus Problems and Exercises](https://web.pdx.edu/~erdman/CALCULUS/CALCULUS_pdf.pdf)

  * **Northwestern University**. [Real Analysis – Lecture Notes](https://sites.math.northwestern.edu/scg479/courses/notes/lecture-notes-320-3.pdf)

  * **Penn State University**. [Calculus – Derivatives and Applications](https://www.math.psu.edu/lamy/misc/Calc_Notes.pdf)


Derivatives

The derivative describes a function’s rate of change.

9.1k

[Difference Quotient](https://algebrica.org/difference-quotient/)

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
