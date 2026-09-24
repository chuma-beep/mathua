# Mean Squared Error Loss

**Mean squared error:** For predictions $y$ and targets $t$ over $n$ samples, $\text{MSE} = (1/n)\sum (y_i - t_i)^2$: the average squared distance. It is non-negative, zero only at perfect prediction, and punishes large errors more than small ones.

## Worked: MSE of two predictions

Take predictions $[0, 1]$ against targets $[1, 1]$:
1. Square each error: $(0-1)^2 = 1$ and $(1-1)^2 = 0$.
2. Add them: $1 + 0 = 1$.
3. Divide by $n=2$: MSE $= 0.5$.

So doubling one error from 1 to 2 would quadruple its contribution from 1 to 4 — large misses dominate the loss.

## Worked: one gradient step

Minimize $(1/2)(y-t)^2$ with prediction $y=2$, target $t=3$, learning rate 0.1:
1. Differentiate: $d/dy\, (1/2)(y-t)^2 = y - t$, so the gradient is $2 - 3 = -1$.
2. Step against the gradient: $y \leftarrow 2 - 0.1\cdot(-1) = 2.1$.
3. Recheck the loss: $(2.1-3)^2 = 0.81$, down from 1 — the negative gradient pointed at the target.

So the gradient $y - t$ always points from target back to prediction, and descending it pulls predictions toward targets.

## When MSE is the wrong loss

MSE measures Euclidean distance for continuous targets, so on classification it is flat when confidently wrong: prediction 0.1 for label 1 gives MSE 0.81, barely worse than predicting 0.5. Cross-entropy instead grows without bound there, which is why classification uses it.
