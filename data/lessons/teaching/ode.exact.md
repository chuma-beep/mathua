# Exact Differential Equations

A differential equation $M(x,y) dx + N(x,y) dy = 0$ is **exact** if there exists a function $\psi(x,y)$ such that:

$$\frac{\partial \psi}{\partial x} = M, \quad \frac{\partial \psi}{\partial y} = N$$

**Exactness condition:** The equation is exact if and only if:

$$\frac{\partial M}{\partial y} = \frac{\partial N}{\partial x}$$

## Solution Method

1. Check exactness: verify $\frac{\partial M}{\partial y} = \frac{\partial N}{\partial x}$.
2. Find $\psi$ by integrating:
   $$\psi = \int M dx + h(y)$$
3. Determine $h(y)$ using $\frac{\partial \psi}{\partial y} = N$.
4. The general solution is $\psi(x, y) = C$.

**Example 1:** Solve $(2xy + 1) dx + (x^2 + 3y^2) dy = 0$.

$M = 2xy + 1$, $N = x^2 + 3y^2$
$\frac{\partial M}{\partial y} = 2x$, $\frac{\partial N}{\partial x} = 2x$ ✓ Exact

$\psi = \int (2xy + 1) dx = x^2y + x + h(y)$
$\frac{\partial \psi}{\partial y} = x^2 + h'(y) = x^2 + 3y^2$, so $h'(y) = 3y^2$, $h(y) = y^3$.

Solution: $x^2y + x + y^3 = C$

**Example 2:** Solve $(e^y + \cos x) dx + (xe^y + 2y) dy = 0$.

$M = e^y + \cos x$, $N = xe^y + 2y$
$\frac{\partial M}{\partial y} = e^y$, $\frac{\partial N}{\partial x} = e^y$ ✓ Exact

$\psi = \int (e^y + \cos x) dx = xe^y + \sin x + h(y)$
$\frac{\partial \psi}{\partial y} = xe^y + h'(y) = xe^y + 2y$, so $h'(y) = 2y$, $h(y) = y^2$.

Solution: $xe^y + \sin x + y^2 = C$

## Integrating Factors (Nonexact to Exact)

If $\frac{\partial M}{\partial y} \neq \frac{\partial N}{\partial x}$, we may find an integrating factor $\mu(x)$ or $\mu(y)$:

- If $\frac{M_y - N_x}{N}$ depends only on $x$: $\mu(x) = e^{\int \frac{M_y - N_x}{N} dx}$
- If $\frac{N_x - M_y}{M}$ depends only on $y$: $\mu(y) = e^{\int \frac{N_x - M_y}{M} dy}$
