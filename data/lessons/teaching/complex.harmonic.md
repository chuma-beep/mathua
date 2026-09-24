# Harmonic Functions and Mean Value

**Harmonic:** $u(x,y)$ is harmonic if $\Delta u=u_{xx}+u_{yy}=0$. The real and imaginary parts of any analytic $f=u+iv$ are harmonic by CR; on a simply connected domain every harmonic $u$ has a conjugate $v$ making $u+iv$ analytic.

## Worked: x squared minus y squared is harmonic

Test $u=x^2-y^2$:
1. Differentiate twice in $x$: $u_x=2x$, then $u_{xx}=2$.
2. Differentiate twice in $y$: $u_y=-2y$, then $u_{yy}=-2$.
3. Add: $\Delta u=2+(-2)=0$.

So $u$ is harmonic; it is $\operatorname{Re} z^2$, and indeed the real part of every analytic function passes the Laplacian test.

## Worked: recovering the conjugate of x squared minus y squared

Find $v$ with $u+iv$ analytic for $u=x^2-y^2$:
1. Use $v_x=-u_y$: here $u_y=-2y$, so $v_x=2y$, giving $v=2xy+h(y)$.
2. Use $v_y=u_x$: here $u_x=2x$, and $v_y=2x+h'(y)$, so $h'(y)=0$.
3. Take $h=0$: the conjugate is $v=2xy$, and $u+iv=z^2$.

So harmonic conjugates come from integrating CR: the two equations determine $v$ up to an additive constant.

## Worked: mean value on x squared minus y squared

Average $u=x^2-y^2$ over the circle of radius $r$ around the origin:
1. Parametrize: $x=r\cos t$, $y=r\sin t$, so $u=r^2(\cos^2 t-\sin^2 t)=r^2\cos 2t$.
2. Average $t$ from 0 to $2\pi$: $\cos 2t$ averages to 0.
3. Compare with the center: $u(0,0)=0$, matching the average.

So harmonic functions satisfy the mean value property: $u(a)$ equals the average of $u$ over any surrounding circle, which also forces non-constant harmonic functions to take maxima only on the boundary.
