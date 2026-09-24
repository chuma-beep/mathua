# Gradient Descent

**Gradient descent:** For parameters $w$ and loss $L$, repeat $w \leftarrow w - \eta \nabla L(w)$: step opposite the gradient, scaled by the learning rate $\eta$. The negative gradient is the locally steepest downhill direction.

## Worked: one gradient step

Take weight $w = 1$, gradient $\nabla L = 2$, learning rate 0.1:
1. Scale the gradient: $0.1 \cdot 2 = 0.2$.
2. Step against it: $w \leftarrow 1 - 0.2 = 0.8$.
3. The loss falls: on $L = w^2$ it drops from 1 to 0.64.

So each step trades gradient information for loss decrease — small, repeated, downhill.

## Worked: stochastic versus batch cost

Take 1000 training points and mini-batches of 32:
1. Batch gradient descent computes all 1000 gradients, then takes 1 step per epoch.
2. Mini-batch SGD takes $1000/32 \approx 31$ steps per epoch on noisy estimates.
3. SGD with one sample per update takes 1000 noisy steps per epoch — fastest per step, noisiest direction.

So batch size trades gradient accuracy against update count: more noise per step, many more steps per pass.

## When descent diverges

The same rule with $\eta$ too large overshoots: on $L = w^2$ from $w = 1$ with $\eta = 1.5$, the update jumps to $1 - 1.5\cdot 2 = -2$, and the loss grows from 1 to 4. Too large diverges, too small crawls — the rate must be tuned between the two.
