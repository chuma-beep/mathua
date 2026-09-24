# Variation of Parameters

**Variation of parameters:** for $y'' + p(x)y' + q(x)y = g(x)$ with homogeneous basis $y_1, y_2$ and Wronskian $W = y_1 y_2' - y_1' y_2$, a particular solution is $y_p = u_1 y_1 + u_2 y_2$ where $u_1' = -y_2 g/W$ and $u_2' = y_1 g/W$. It works for any $g$ — even where undetermined coefficients fails.

## Worked: the Wronskian of cos and sin

Compute $W$ for $y_1 = \cos x$, $y_2 = \sin x$:
1. Derivatives: $y_1' = -\sin x$, $y_2' = \cos x$.
2. $W = \cos x \cdot \cos x - (-\sin x) \cdot \sin x = \cos^2 x + \sin^2 x$.
3. So $W = 1$: constant, never zero — the pair is independent everywhere.

So this basis is as clean as it gets: unit Wronskian, no vanishing points.

## Worked: a constant forcing term

Solve $y'' + y = 1$ with the same basis:
1. Rates: $u_1' = -\sin x \cdot 1/1 = -\sin x$; $u_2' = \cos x \cdot 1/1 = \cos x$.
2. Integrate: $u_1 = \cos x$, $u_2 = \sin x$ (drop constants — they rebuild the homogeneous part).
3. Assemble: $y_p = \cos^2 x + \sin^2 x = 1$. Check: $y_p'' + y_p = 0 + 1 = 1$.

So the full solution is $y = C_1 \cos x + C_2 \sin x + 1$.

## It needs the homogeneous basis first

No $y_1, y_2$, no $W$, no $u_1, u_2$ — the method builds on the homogeneous solution. Where both methods apply (constant coefficients, friendly $g$) they agree; variation of parameters alone survives hostile $g$ like $\sec x$.
