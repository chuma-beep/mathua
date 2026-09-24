# Laplace Equation

**Laplace equation:** $\Delta u = u_{xx} + u_{yy} = 0$. Solutions are **harmonic**: they satisfy the mean value property, take their maxima on the boundary, and solve the Dirichlet problem (prescribed boundary values, harmonic inside).

## Worked: checking u = xy is harmonic

Test $u(x, y) = xy$:
1. Differentiate in $x$ twice: $u_x = y$, so $u_{xx} = 0$.
2. Differentiate in $y$ twice: $u_y = x$, so $u_{yy} = 0$.
3. Add: $\Delta u = 0 + 0 = 0$.

So $xy$ is harmonic — while $x^2 + y^2$ is not, since $\Delta(x^2+y^2) = 2 + 2 = 4$.

## Worked: the center equals the average

Take the harmonic $u(x, y) = x^2 - y^2$ on the disk of radius 2:
1. Average $u$ over the circle: opposite points $(\pm a, \pm b)$ cancel, so the average is 0.
2. The mean value property says $u$ at the center equals that average: $u(0, 0) = 0$.
3. Check directly: $0^2 - 0^2 = 0$.

So interior values are forced by surrounding values — harmonic functions cannot freelance.

## Maxima live on the boundary

A non-constant harmonic function never attains its maximum (or minimum) inside the domain — both sit on the boundary. That is the maximum principle, and it makes the Dirichlet problem well-posed: boundary data $f$ extends inward to exactly one harmonic $u$.
