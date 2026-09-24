# Lyapunov Functions

**Lyapunov function:** $V$ with $V(x) > 0$ for $x \neq 0$ and $V(0) = 0$, whose derivative along trajectories satisfies $dV/dt < 0$. Then the equilibrium is asymptotically stable — proved without solving the ODE (the direct method).

## Worked: V = x^2 for x' = -3x

Take $V(x) = x^2$:
1. $V > 0$ away from zero and $V(0) = 0$: valid candidate.
2. Chain rule: $dV/dt = 2x \cdot (-3x) = -6x^2$.
3. At $x = 2$ this is $-24 < 0$; in fact $dV/dt < 0$ for every nonzero $x$.

So $V$ drains along every trajectory and 0 is asymptotically stable — no formula for $x(t)$ needed.

## Worked: stable but not asymptotic

Take $x' = 0$ with the same $V(x) = x^2$:
1. $V > 0$ away from zero, $V(0) = 0$: still a valid candidate.
2. But $dV/dt = 2x \cdot 0 = 0$ everywhere — never strictly negative.
3. Trajectories freeze where they start: nearby stays nearby (stable), nothing converges (not asymptotically stable).

So non-strict $dV/dt \le 0$ buys stability only; strict negativity buys convergence.

## Stability without solving

LaSalle extends the method to $dV/dt \le 0$ via invariant sets, and the converse holds too: every asymptotically stable equilibrium owns some Lyapunov function. Energy-like decay certifies stability.
