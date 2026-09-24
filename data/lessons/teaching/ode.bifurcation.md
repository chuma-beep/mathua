# Bifurcations in Nonlinear ODEs

**Bifurcation:** a qualitative change in the phase portrait as a parameter $\mu$ crosses a critical value — the number or stability of equilibria changes. The picture above the threshold is not a continuation of the picture below it.

## Worked: saddle-node creation

Take $x' = \mu - x^2$ with $\mu = 4$:
1. Set $4 - x^2 = 0$, so equilibria sit at $x = 2$ and $x = -2$.
2. Between the roots $x' > 0$ (flow moves right); outside them $x' < 0$ (flow moves left). So $x = 2$ attracts from the left only, and $x = -2$ repels to the right only.
3. Drop to $\mu = -1$: then $x' = -1 - x^2 < 0$ everywhere, so there are no equilibria at all.

So the saddle-node at $\mu = 0$ creates two equilibria out of none as $\mu$ turns positive.

## Worked: pitchfork branches

Take $x' = \mu x - x^3$ with $\mu = 1$:
1. Factor $x(1 - x^2) = 0$: equilibria at $x = 0$, $x = 1$, and $x = -1$.
2. For large positive $x$ the cubic term wins and $x' < 0$; checking the sign on each interval gives $x = 0$ unstable, while $x = 1$ and $x = -1$ are stable.
3. Drop to $\mu = -1$: then $x' = -x - x^3$ has the sign of $-x$ everywhere, so only $x = 0$ remains, and it is stable.

So the pitchfork at $\mu = 0$ splits one stable equilibrium into an unstable one plus two stable branches at $\pm\sqrt{\mu}$ — the normal form of symmetry breaking.

## When the portrait changes

Saddle-node creates or destroys pairs, transcritical ($x' = \mu x - x^2$) swaps stability between two equilibria that cross, pitchfork splits one into three, and Hopf births a limit cycle when a complex pair crosses the imaginary axis. In every case the test is the same: count equilibria and read flow signs on each side of the threshold.
