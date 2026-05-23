# Separable Differential Equations

A first-order differential equation is **separable** if it can be written as:

$$\frac{dy}{dx} = f(x)g(y)$$

where the right side factors into a product of a function of $x$ and a function of $y$.

## Solution Method

1. Separate variables: $\frac{1}{g(y)} dy = f(x) dx$
2. Integrate both sides: $\int \frac{1}{g(y)} dy = \int f(x) dx$
3. Solve for $y$ if possible.

**Example 1:** Solve $\frac{dy}{dx} = 2xy$.

Separate: $\frac{1}{y} dy = 2x dx$
Integrate: $\int \frac{1}{y} dy = \int 2x dx$
$$\ln|y| = x^2 + C$$
$$y = \pm e^{x^2 + C} = Ae^{x^2}$$ where $A = \pm e^C$.

**Example 2:** Solve the initial value problem $\frac{dy}{dx} = y^2$, $y(0) = 1$.

Separate: $\frac{1}{y^2} dy = dx$
Integrate: $\int y^{-2} dy = \int dx$
$$-y^{-1} = x + C$$
$$y = -\frac{1}{x + C}$$
Use $y(0) = 1$: $1 = -\frac{1}{C}$, so $C = -1$.
$$y = -\frac{1}{x - 1} = \frac{1}{1 - x}$$

**Example 3:** Solve $\frac{dy}{dx} = 3x^2 e^{-y}$.

Separate: $e^y dy = 3x^2 dx$
Integrate: $e^y = x^3 + C$
$$y = \ln(x^3 + C)$$
