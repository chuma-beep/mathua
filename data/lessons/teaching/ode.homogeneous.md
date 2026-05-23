# Homogeneous Differential Equations

A first-order ODE is **homogeneous** if it can be written as:

$$\frac{dy}{dx} = F\left(\frac{y}{x}\right)$$

## Substitution Method

Use the substitution $v = \frac{y}{x}$, so $y = vx$ and $\frac{dy}{dx} = v + x\frac{dv}{dx}$.

The equation becomes:
$$v + x\frac{dv}{dx} = F(v)$$
$$x\frac{dv}{dx} = F(v) - v$$

This is separable in $v$ and $x$.

**Example 1:** Solve $\frac{dy}{dx} = \frac{x^2 + y^2}{xy}$.

Rewrite: $\frac{dy}{dx} = \frac{x}{y} + \frac{y}{x} = \frac{1}{v} + v$ where $v = y/x$.

Substitute $y = vx$:
$v + x\frac{dv}{dx} = \frac{1}{v} + v$
$x\frac{dv}{dx} = \frac{1}{v}$
$v dv = \frac{1}{x} dx$

Integrate: $\frac{v^2}{2} = \ln|x| + C$
$v^2 = 2\ln|x| + C$
$\frac{y^2}{x^2} = \ln(x^2) + C$
$$y^2 = x^2\ln(Cx^2)$$

**Example 2:** Solve $\frac{dy}{dx} = \frac{y + \sqrt{x^2 + y^2}}{x}$.

Divide numerator and denominator by $x$:
$\frac{dy}{dx} = \frac{y}{x} + \sqrt{1 + (\frac{y}{x})^2} = v + \sqrt{1 + v^2}$

Substitute $y = vx$:
$v + x\frac{dv}{dx} = v + \sqrt{1 + v^2}$
$x\frac{dv}{dx} = \sqrt{1 + v^2}$
$\frac{dv}{\sqrt{1 + v^2}} = \frac{dx}{x}$

Integrate: $\sinh^{-1}(v) = \ln|x| + C$
$v = \sinh(\ln|x| + C)$
$$y = x \sinh(\ln|x| + C)$$

**Note:** A homogeneous ODE is different from a linear homogeneous ODE (where $Q(x) = 0$). The term "homogeneous" has two different meanings in ODE theory.
