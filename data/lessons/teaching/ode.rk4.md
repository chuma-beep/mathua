# Runge-Kutta 4th Order Method

**RK4:** for $y' = f(x, y)$ with step $h$, sample four slopes $k_1 = f(x_n, y_n)$, $k_2 = f(x_n + h/2, y_n + hk_1/2)$, $k_3 = f(x_n + h/2, y_n + hk_2/2)$, $k_4 = f(x_n + h, y_n + hk_3)$, then $y_{n+1} = y_n + h(k_1 + 2k_2 + 2k_3 + k_4)/6$. A weighted average of slopes instead of Euler's single tangent.

## Worked: one RK4 step for y' = y

Take $y(0) = 1$, $h = 0.1$:
1. $k_1 = 1$; $k_2 = 1 + 0.05 \cdot 1 = 1.05$; $k_3 = 1 + 0.05 \cdot 1.05 = 1.0525$; $k_4 = 1 + 0.1 \cdot 1.0525 = 1.10525$.
2. Average: $(1 + 2.1 + 2.105 + 1.10525)/6 \approx 1.0517$.
3. Step: $y_1 = 1 + 0.1 \cdot 1.0517 \approx 1.10517$.

So the four probes feel the curvature that Euler ignores.

## Worked: RK4 against Euler

Compare at $x = 0.1$ (exact $e^{0.1} \approx 1.10517$):
1. RK4 gives 1.10517: error around $10^{-7}$.
2. Euler gave 1.1: error around 0.0052.
3. Same step size, five orders of magnitude apart — the extra probes pay off.

So RK4 dominates Euler whenever function evaluations are affordable.

## Order costs evaluations

RK4 needs four evaluations per step for local error $O(h^5)$ and global error $O(h^4)$ — halving $h$ cuts global error sixteenfold, versus Euler's twofold. Higher order per step beats smaller steps per order.
