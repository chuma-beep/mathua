# Boundary Value Problems

**BVP:** an ODE $y'' = f(x, y, y')$ with conditions pinned at two points, $y(a) = \alpha$ and $y(b) = \beta$ — unlike an IVP, which prescribes $y$ and $y'$ at one point. A BVP may have no solution, one, or infinitely many.

## Worked: pinning down y'' = 0

Solve $y'' = 0$ with $y(0) = 0$ and $y(1) = 3$:
1. Integrate twice: $y' = C$, then $y = Cx + D$.
2. Apply $y(0) = 0$: $D = 0$, so $y = Cx$.
3. Apply $y(1) = 3$: $C = 3$, so $y = 3x$.

So this BVP has exactly one solution, $y = 3x$ — the straight line through the two pinned points.

## Worked: same ODE, different pins

Keep $y'' = 0$ but change the pins to $y(0) = 1$ and $y(1) = 1$:
1. Again $y = Cx + D$ from two integrations.
2. Now $y(0) = 1$ gives $D = 1$, so $y = Cx + 1$.
3. Then $y(1) = 1$ gives $C + 1 = 1$, so $C = 0$ and $y = 1$.

So the constant function $y = 1$ is the unique solution — the pins alone select it from the whole family $Cx + D$.

## IVP always answers, BVP sometimes

Picard guarantees an IVP has exactly one local solution, but a BVP can have zero, one, or infinitely many — resonance ($y'' + \pi^2 y = 0$ with $y(0) = y(1) = 0$) admits a whole line of solutions. The shooting method turns a BVP into repeated IVPs that home in on the far pin, while Green's functions and the Fredholm alternative state exactly when a solution exists.
