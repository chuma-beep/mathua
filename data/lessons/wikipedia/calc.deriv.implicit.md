> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Implicit_differentiation) — CC BY-SA 4.0

# Implicit differentiation

In calculus, **implicit differentiation** is a method of finding the derivative of an implicit function using the chain rule.

To differentiate an implicit function \(y(x)\), defined by an equation \(R(x, y) = 0\), it is not generally possible to solve it explicitly for \(y\) and then differentiate it. Instead, one can differentiate \(R(x, y) = 0\) with respect to \(x\) and \(y\) and then solve the resulting equation for \(\frac{dy}{dx}\).

## Formulation

If \(R(x, y) = 0\), the derivative of the implicit function \(y(x)\) is given by

\[
\frac{dy}{dx} = -\frac{\frac{\partial R}{\partial x}}{\frac{\partial R}{\partial y}} = -\frac{R_x}{R_y},
\]

where \(R_x\) and \(R_y\) indicate the partial derivatives of \(R\) with respect to \(x\) and \(y\).

The above formula comes from using the generalized chain rule to obtain the total derivative — with respect to \(x\) — of both sides of \(R(x, y) = 0\):

\[
\frac{\partial R}{\partial x}\frac{dx}{dx} + \frac{\partial R}{\partial y}\frac{dy}{dx} = 0,
\]

hence

\[
\frac{\partial R}{\partial x} + \frac{\partial R}{\partial y}\frac{dy}{dx} = 0,
\]

which, when solved for \(\frac{dy}{dx}\), gives the expression above.

## Examples

### Example 1

Consider

\[
y + x + 5 = 0.
\]

This equation is easy to solve for \(y\), giving

\[
y = -x - 5,
\]

where the right side is the explicit form of the function \(y(x)\). Differentiation then gives \(\frac{dy}{dx} = -1\).

Alternatively, one can differentiate the original equation:

\[
\begin{aligned}
\frac{dy}{dx} + \frac{dx}{dx} + \frac{d}{dx}(5) &= 0; \\
\frac{dy}{dx} + 1 + 0 &= 0.
\end{aligned}
\]

Solving for \(\frac{dy}{dx}\) gives

\[
\frac{dy}{dx} = -1,
\]

the same answer as obtained previously.

### Example 2

An example of an implicit function for which implicit differentiation is easier than using explicit differentiation is the function \(y(x)\) defined by the equation

\[
x^4 + 2y^2 = 8.
\]

To differentiate this explicitly with respect to \(x\), one has first to get

\[
y(x) = \pm \sqrt{\frac{8 - x^4}{2}},
\]

and then differentiate this function. This creates two derivatives: one for \(y \ge 0\) and another for \(y < 0\).

It is substantially easier to implicitly differentiate the original equation:

\[
4x^3 + 4y\frac{dy}{dx} = 0,
\]

giving

\[
\frac{dy}{dx} = \frac{-4x^3}{4y} = -\frac{x^3}{y}.
\]

### Example 3

Often, it is difficult or impossible to solve explicitly for \(y\), and implicit differentiation is the only feasible method of differentiation. An example is the equation

\[
y^5 - y = x.
\]

It is impossible to algebraically express \(y\) explicitly as a function of \(x\), and therefore one cannot find \(\frac{dy}{dx}\) by explicit differentiation. Using the implicit method, \(\frac{dy}{dx}\) can be obtained by differentiating the equation to obtain

\[
5y^4\frac{dy}{dx} - \frac{dy}{dx} = \frac{dx}{dx},
\]

where \(\frac{dx}{dx} = 1\). Factoring out \(\frac{dy}{dx}\) shows that

\[
\left(5y^4 - 1\right)\frac{dy}{dx} = 1,
\]

which yields the result

\[
\frac{dy}{dx} = \frac{1}{5y^4 - 1},
\]

which is defined for \(y \ne \pm \frac{1}{\sqrt[4]{5}}\).
