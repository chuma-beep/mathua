# Dropout Regularization

**Dropout:** During training, zero each activation independently with probability $p$, and scale survivors by $1/(1-p)$ (inverted dropout) to keep expectations fixed. At test time dropout is off and all units stay on.

## Worked: one dropout mask

Take a layer of 100 units with $p = 0.5$:
1. Each unit survives with probability $1 - p = 0.5$, so about 50 stay active per forward pass.
2. Scale survivors by $1/(1-0.5) = 2$: a unit outputting 3 contributes 6 when kept, 0 when dropped — expectation stays 3 either way.
3. At test time all 100 units run unscaled, averaging the subnetworks seen in training.

So inverted dropout keeps train and test expectations matched: the network never sees a scale shift at deployment.

## Worked: counting the ensemble

Take $n$ units with independent masks:
1. Each unit is on or off: 2 states per unit, $2^n$ possible masks total.
2. With $n = 100$ that is $2^{100}$ subnetworks sharing weights — far more than could ever train separately.
3. Each training step samples one mask and updates only its active weights, so every subnetwork gets some training.

So dropout trains an enormous weight-shared ensemble cheaply, and the test-time full network approximates its average.

## When dropout hurts

Dropout noise helps large networks but starves small ones: a 10-unit layer at $p = 0.5$ keeps 5 units and may never form stable features. Use it where capacity is ample; skip it (or lower $p$) on tiny layers and at test time, always.
