# Regularization (L2 Weight Decay)

**Regularization:** Add a penalty on model complexity to the loss,

$$L_{\text{reg}} = L_{\text{data}} + \lambda\|w\|^{2}$$

L2 (weight decay) shrinks weights toward zero, discouraging fitting noise.

## Why Regularize

### Controlling Capacity
Without penalty, large weights can create sharp decision boundaries that memorize training points. L2 adds $2\lambda w$ to the gradient, pulling weights inward each step.

### The Strength $\lambda$
Small $\lambda$: fits training data closely. Large $\lambda$: heavily shrinks weights, underfits. Tuning $\lambda$ on a validation set balances bias and variance.

## Example

Linear model $y=w x$ on noisy points. Unregularized $w=5.2$ fits noise; with $\lambda=0.1$, gradient $2\lambda w$ reduces $w$ to $3.1$, smoothing the fit and lowering test error.
