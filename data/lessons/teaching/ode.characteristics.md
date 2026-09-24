# First-Order PDEs and Characteristics

**First-order PDE:** $a(x,y) u_x + b(x,y) u_y = c$, with **characteristics** the curves $x(t), y(t)$ satisfying $dx/dt = a$ and $dy/dt = b$. Along a characteristic the PDE collapses to the ODE $du/dt = c$.

## Worked: riding the transport equation

Solve $u_t + 2 u_x = 0$ with $u(x, 0) = x^2$:
1. Characteristics satisfy $dx/dt = 2$, so $x - 2t$ is constant along each one.
2. With $c = 0$, $du/dt = 0$ along characteristics: $u$ is constant on each line $x - 2t = \text{const}$.
3. Trace back to $t = 0$: $u(x, t) = (x - 2t)^2$. At the point $(5, 1)$ this gives $u = 9$.

So the initial parabola slides right at speed 2 unchanged — transport just shifts data along characteristics.

## Worked: reading char ODEs

Take $3 u_x + u_y = 0$:
1. Write the char ODEs $dx/3 = dy/1$: along characteristics $x - 3y$ is constant.
2. With zero right side, $u$ is constant on each characteristic, so $u = f(x - 3y)$ for some $f$.
3. Check: $u_x = f'$, $u_y = -3f'$, and $3f' - 3f' = 0$.

So $dx/a = dy/b$ finds the curves, and constancy of $u$ along them writes the solution.

## When characteristics cross

If characteristics intersect, they carry two different data values to one point and the solution turns multi-valued — a shock forms. Quasi-linear equations, where $u$ feeds back into $a$ and $b$, bend their own characteristics and are the usual shock producers.
