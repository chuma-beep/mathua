# First-Order Linear Differential Equations

A first-order linear ODE has the standard form:

$$\frac{dy}{dx} + P(x)y = Q(x)$$

## Integrating Factor Method

1. Find the integrating factor: $\mu(x) = e^{\int P(x) dx}$
2. Multiply both sides by $\mu(x)$:
   $$\mu(x)\frac{dy}{dx} + \mu(x)P(x)y = \mu(x)Q(x)$$
   The left side becomes $\frac{d}{dx}[\mu(x)y]$.
3. Integrate:
   $$\mu(x)y = \int \mu(x)Q(x) dx + C$$
4. Solve for $y$.

**Example 1:** Solve $\frac{dy}{dx} + 2y = 6e^{x}$.

$P(x) = 2$, so $\mu = e^{\int 2 dx} = e^{2x}$.
Multiply: $e^{2x}\frac{dy}{dx} + 2e^{2x}y = 6e^{3x}$
The left is $\frac{d}{dx}[e^{2x}y] = 6e^{3x}$.
Integrate: $e^{2x}y = \int 6e^{3x} dx = 2e^{3x} + C$
$$y = 2e^{x} + Ce^{-2x}$$

**Example 2:** Solve $\frac{dy}{dx} + \frac{1}{x}y = x^2$, $x > 0$.

$P(x) = \frac{1}{x}$, so $\mu = e^{\int \frac{1}{x} dx} = e^{\ln x} = x$.
Multiply: $x\frac{dy}{dx} + y = x^3$
The left is $\frac{d}{dx}[xy] = x^3$.
Integrate: $xy = \int x^3 dx = \frac{x^4}{4} + C$
$$y = \frac{x^3}{4} + \frac{C}{x}$$

**Example 3:** Solve the IVP $\frac{dy}{dx} + 3y = \sin x$, $y(0) = 1$.

$\mu = e^{3x}$.
$\frac{d}{dx}[e^{3x}y] = e^{3x}\sin x$
$e^{3x}y = \int e^{3x}\sin x dx = \frac{e^{3x}}{10}(3\sin x - \cos x) + C$
$$y = \frac{1}{10}(3\sin x - \cos x) + Ce^{-3x}$$
Using $y(0) = 1$: $1 = \frac{1}{10}(-1) + C$, so $C = \frac{11}{10}$.
$$y = \frac{1}{10}(3\sin x - \cos x) + \frac{11}{10}e^{-3x}$$
