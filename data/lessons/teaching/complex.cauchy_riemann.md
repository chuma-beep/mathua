# Cauchy-Riemann Equations

**Cauchy-Riemann equations:** Write $f(z)=u(x,y)+iv(x,y)$ with $z=x+iy$. If $f$ is complex-differentiable at a point then $u_x=v_y$ and $u_y=-v_x$ there; with continuous partials these equations are also sufficient.

## Worked: z squared satisfies CR everywhere

Test $f(z)=z^2$ with $u=x^2-y^2$ and $v=2xy$:
1. Differentiate in $x$: $u_x=2x$ and $v_x=2y$.
2. Differentiate in $y$: $u_y=-2y$ and $v_y=2x$.
3. Compare: $u_x=2x=v_y$ and $u_y=-2y=-v_x$ at every point.

So $z^2$ satisfies CR everywhere and is entire; polynomials pass the test because their real and imaginary parts are matched in degree.

## Worked: the conjugate fails CR everywhere

Test $f(z)=\bar{z}$ with $u=x$ and $v=-y$:
1. Differentiate in $x$: $u_x=1$ and $v_x=0$.
2. Differentiate in $y$: $u_y=0$ and $v_y=-1$.
3. Compare: $u_x=1 \ne -1=v_y$, so the first equation already fails.

So $\bar{z}$ is analytic nowhere: reflection across the real axis reverses orientation, which no complex derivative can do.

## Worked: e to the z satisfies CR everywhere

Test $f(z)=e^z$ with $u=e^x\cos y$ and $v=e^x\sin y$:
1. Differentiate in $x$: $u_x=e^x\cos y$ and $v_x=e^x\sin y$.
2. Differentiate in $y$: $u_y=-e^x\sin y$ and $v_y=e^x\cos y$.
3. Compare: $u_x=v_y$ and $u_y=-v_x$ at every point.

So $e^z$ is entire: the exponential differentiates itself in $x$ while the $y$-derivative rotates through sine and cosine to match.
