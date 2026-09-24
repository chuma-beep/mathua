# Regularization (L2 Weight Decay)

**Regularization:** Train on $L = L_{\text{data}} + \lambda\lVert w\rVert^2$: data fit plus a penalty $\lambda$ times squared weight size. The penalty term pulls weights toward zero every step, fighting memorization of noise.

## Worked: the penalty gradient

Take weight $w = 5.2$ with strength $\lambda = 0.1$:
1. Differentiate the penalty: $d/dw\, \lambda w^2 = 2\lambda w = 2\cdot 0.1\cdot 5.2 = 1.04$.
2. The update subtracts this pull: even with zero data gradient, $w$ shrinks by $0.1\cdot 1.04 \approx 0.10$ at rate 0.1.
3. Large weights feel a large pull, tiny weights feel almost none — shrinkage is proportional to size.

So L2 is weight decay with decay rate set by $\lambda$: the bigger the weight, the harder the inward pull.

## Worked: tuning the strength

Take a linear fit $y = wx$ on noisy points where the unregularized fit reaches $w = 5.2$:
1. With $\lambda = 0$: $w$ stays 5.2, fitting noise — low training error, high test error.
2. With $\lambda = 0.1$: the pull settles $w$ near 3.1, smoothing the fit and lowering test error.
3. With $\lambda = 10$: $w$ collapses near 0, underfitting both sets.

So $\lambda$ dials capacity: too small overfits, too large underfits, and the validation set picks the middle.

## When penalties hurt

With abundant clean data the penalty biases a good fit: shrinking already-correct weights adds error for no variance gain. Regularize where data is scarce or noisy, not where the training set already pins the function down.
