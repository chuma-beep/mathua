# Heat Equation and Separation of Variables

**Heat equation:** $u_t = k u_{xx}$ on a rod $[0, L]$. **Separation:** assume $u = X(x)T(t)$; then $X''/X = T'/(kT) = -\lambda$, one ODE per factor.

## Worked: splitting into two ODEs

Insert $u = X(x)T(t)$ into $u_t = k u_{xx}$:
1. Get $X T' = k X'' T$, then divide by $kXT$: $T'/(kT) = X''/X$.
2. Left side depends only on $t$, right only on $x$ — both equal one constant $-\lambda$.
3. So $T' = -k\lambda T$ gives $T = C e^{-k\lambda t}$, and $X'' + \lambda X = 0$ carries the spatial shape.

So separation trades one PDE for a decay ODE in time and an oscillation ODE in space.

## Worked: modes on a unit rod

Take $L = \pi$ with ends pinned at zero:
1. $X'' + \lambda X = 0$ with $X(0) = X(\pi) = 0$ forces $\lambda_n = n^2$ and $X_n = \sin(nx)$.
2. Each mode decays as $e^{-k n^2 t}$: higher $n$ dies faster.
3. Initial $u(x, 0) = \sin x$ is already mode one, so $u(x, t) = e^{-kt} \sin x$.

So a pure sine just fades in place, and general data is a sum of such fading modes.

## Smoothing and the maximum principle

Heat peaks cannot sit inside the rod: the maximum over space-time lives on the boundary (initial data or rod ends). Even jagged initial data is instantly smooth for $t > 0$ — diffusion irons every corner.
