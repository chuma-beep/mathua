# Adam Optimizer Details

**Adam:** Track a first moment $m$ (mean gradient, momentum) and second moment $v$ (mean squared gradient, RMSProp-style scale), bias-correct both, and update $w \leftarrow w - \eta\hat{m}/(\sqrt{\hat{v}}+\epsilon)$. Defaults $\beta_1 = 0.9$, $\beta_2 = 0.999$, $\epsilon = 10^{-8}$: adaptive step per weight.

## Worked: the first update

Take first step $t = 1$ with gradient $g = 2$, $m = 0$, $v = 0$, $\eta = 0.1$:
1. Update moments: $m = 0.9\cdot 0 + 0.1\cdot 2 = 0.2$; $v = 0.999\cdot 0 + 0.001\cdot 4 = 0.004$.
2. Bias-correct (early averages start at zero): $\hat{m} = 0.2/(1-0.9) = 2$; $\hat{v} = 0.004/(1-0.999) = 4$.
3. Step: $0.1\cdot 2/(\sqrt{4}+\epsilon) = 0.1$ — the correction recovered the true gradient scale exactly.

So bias correction matters early: without it the first step would use $m = 0.2$ raw and crawl at a tenth of the intended size.

## Worked: why moments adapt

Take gradients alternating $+5, -5$ across valley walls with steady $+0.3$ downhill:
1. The first moment averages the walls to near zero while the 0.3 accumulates — momentum damps oscillation, keeps descent.
2. The second moment records wall variance $\approx 25$, so wall steps scale by $1/5$ while the downhill direction scales larger.
3. Each weight thus gets its own effective rate $\eta/\sqrt{v}$: steep noisy directions throttled, quiet steady ones promoted.

So Adam is momentum plus per-weight normalization: direction from $m$, step size from $v$.

## When Adam is not the answer

The adaptive scaling can generalize worse than plain SGD with momentum on some vision tasks, and the $\epsilon$ placement matters at low precision. Default to Adam for fast prototyping and sparse gradients; switch to tuned SGD where final generalization dominates.
