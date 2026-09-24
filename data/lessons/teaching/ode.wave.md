# Wave Equation and d'Alembert Solution

**Wave equation:** $u_{tt} = c^2 u_{xx}$ — hyperbolic, finite propagation speed $c$. **d'Alembert:** $u(x,t) = f(x-ct) + g(x+ct)$, a right-traveller plus a left-traveller fixed by the initial data.

## Worked: d'Alembert with a parabola

Take $c = 2$, $u(x, 0) = x^2$, $u_t(x, 0) = 0$:
1. Zero initial velocity splits the data evenly: $u(x,t) = ((x-2t)^2 + (x+2t)^2)/2$.
2. Evaluate at $(3, 1)$: $(1^2 + 5^2)/2 = 13$.
3. Expand to check the shape: $u = x^2 + 4t^2$ — the parabola lifts rigidly, no spreading.

So disturbances keep their shape and just translate — unlike heat, nothing smooths.

## Worked: domain of dependence

For the same wave ($c = 2$), what fixes $u(3, 1)$?
1. Characteristics through $(3, 1)$ are $x \pm 2t = \text{const}$: $3 - 2 = 1$ and $3 + 2 = 5$.
2. So $u(3, 1)$ reads initial data on $[1, 5]$ only.
3. Change $f$ outside $[1, 5]$ and $u(3, 1)$ never notices.

So each point hears exactly the interval $[x - ct, x + ct]$ — finite speed, sharp fronts.

## Energy rides along

The energy $\int (u_t^2 + c^2 u_x^2) dx$ is conserved: waves carry it along characteristics without loss. Hyperbolic means signals, parabolic means spreading, elliptic means equilibrium — the wave equation is the signal case.
