# Laplace Transforms

The **Laplace transform** converts a differential equation into an algebraic equation. It is especially useful for initial value problems and discontinuous forcing functions.

## Definition

$$\mathcal{L}\{f(t)\} = F(s) = \int_0^{\infty} e^{-st} f(t) dt$$

## Key Transforms

| $f(t)$ | $F(s)$ |
|--------|--------|
| $1$ | $\frac{1}{s}$, $s > 0$ |
| $t^n$ | $\frac{n!}{s^{n+1}}$, $s > 0$ |
| $e^{at}$ | $\frac{1}{s-a}$, $s > a$ |
| $\sin(at)$ | $\frac{a}{s^2 + a^2}$, $s > 0$ |
| $\cos(at)$ | $\frac{s}{s^2 + a^2}$, $s > 0$ |
| $t^n e^{at}$ | $\frac{n!}{(s-a)^{n+1}}$, $s > a$ |
| $e^{at}\sin(bt)$ | $\frac{b}{(s-a)^2 + b^2}$ |
| $e^{at}\cos(bt)$ | $\frac{s-a}{(s-a)^2 + b^2}$ |

## Transform of Derivatives

$$\mathcal{L}\{y'\} = sY(s) - y(0)$$
$$\mathcal{L}\{y''\} = s^2Y(s) - sy(0) - y'(0)$$

## Solution Method

1. Take Laplace transform of both sides of the ODE.
2. Solve for $Y(s)$.
3. Take the inverse Laplace transform to find $y(t)$.

**Example 1:** Solve $y'' + y = 1$, $y(0) = 0$, $y'(0) = 0$.

Transform: $s^2Y - 0 - 0 + Y = \frac{1}{s}$
$Y(s^2 + 1) = \frac{1}{s}$
$Y = \frac{1}{s(s^2 + 1)} = \frac{1}{s} - \frac{s}{s^2 + 1}$ (partial fractions)

Inverse: $y(t) = 1 - \cos t$

**Example 2:** Solve $y'' - 3y' + 2y = e^{3t}$, $y(0) = 1$, $y'(0) = 0$.

Transform: $s^2Y - s - 3(sY - 1) + 2Y = \frac{1}{s-3}$
$(s^2 - 3s + 2)Y - s + 3 = \frac{1}{s-3}$
$Y = \frac{s - 3 + \frac{1}{s-3}}{(s-1)(s-2)}$

Inverse: $y(t) = \frac{1}{2}e^{3t} + e^{t} - \frac{1}{2}e^{2t}$ (after partial fractions)
