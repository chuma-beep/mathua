# Euler's Method

**Euler's method:** for $y' = f(x, y)$ with $y(x_0) = y_0$ and step $h$, advance by $y_{n+1} = y_n + h f(x_n, y_n)$ with $x_{n+1} = x_n + h$. It walks the tangent line at each step to trace the solution curve.

## Worked: one step for y' = y

Estimate $y(0.1)$ for $y' = y$, $y(0) = 1$, with $h = 0.1$:
1. Slope at the start: $f(0, 1) = 1$.
2. Step: $y_1 = 1 + 0.1 \cdot 1 = 1.1$.
3. Compare: exact $y(0.1) = e^{0.1} \approx 1.1052$, so the error is about 0.0052.

So one Euler step replaces the true curve by its tangent over a short run.

## Worked: a second step and its error

Continue the same problem to $x = 0.2$:
1. Slope at $(0.1, 1.1)$: $f = 1.1$.
2. Step: $y_2 = 1.1 + 0.1 \cdot 1.1 = 1.21$.
3. Compare: exact $e^{0.2} \approx 1.2214$, so the error has grown to about 0.0114.

So errors accumulate: each step inherits the last step's error and adds fresh tangent error.

## Error shrinks linearly with h

Each step's local error is $O(h^2)$, but a run of length 1 needs $1/h$ steps, leaving global error $O(h)$. Halving $h$ halves the error at double the cost — the price of a first-order method.
