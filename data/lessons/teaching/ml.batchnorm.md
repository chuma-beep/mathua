# Batch Normalization

**Batch norm:** For a mini-batch with mean $\mu$ and variance $\sigma^2$, replace each activation $x$ by $\hat{x} = (x-\mu)/\sqrt{\sigma^2+\epsilon}$, then apply learnable scale and shift $y = \gamma\hat{x} + \beta$. Normalize first, restore capacity after.

## Worked: normalizing one batch

Take a layer input batch $[1, 3, 5]$ with $\gamma = 1$, $\beta = 0$:
1. Mean: $\mu = (1+3+5)/3 = 3$.
2. Variance: $((1-3)^2 + 0 + (5-3)^2)/3 = 8/3 \approx 2.67$.
3. Standardize: $\hat{x} = ([1, 3, 5]-3)/\sqrt{2.67} \approx [-1.22, 0, 1.22]$ — mean 0, variance 1.

So each layer sees stable inputs regardless of how earlier weights drift: mean 0 and variance 1, by construction, every batch.

## Worked: test time without batches

Take inference on a single example, where batch statistics do not exist:
1. During training, track running averages of $\mu$ and $\sigma^2$ across batches.
2. At test time, normalize with those frozen running averages — never the test batch.
3. Keep the learned $\gamma, \beta$: the layer still scales and shifts, just with population statistics.

So batch norm behaves differently in train and test modes: batch statistics while training, running averages at inference — using test-batch statistics leaks information across examples.

## When normalization stabilizes

Standardized inputs tolerate higher learning rates and tame internal covariate shift: gradients neither explode on huge activations nor vanish on tiny ones. The $\gamma, \beta$ pair guarantees no loss of expressiveness — with $\gamma = \sigma$ and $\beta = \mu$ the layer recovers the identity.
