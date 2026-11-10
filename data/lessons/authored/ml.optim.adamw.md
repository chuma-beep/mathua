# AdamW optimization

## Decoupled decay

AdamW separates weight decay from the adaptive gradient update. The decay term shrinks the weights directly instead of adding a penalty to the gradient.

Decoupling changes the effective regularization. The decay acts on every weight in the same way, independent of the gradient history.

## The adaptive step

The adaptive step uses the first and second moment estimates. In the simplest case the step is the learning rate times the normalized gradient. The normalization keeps the step size under control.

The adaptive step and the decay term are computed separately and then subtracted from the weight. Keeping them separate is the defining feature of AdamW.

## The combined update

The new weight is the old weight minus the adaptive step minus the decay term. A larger decay rate or a larger weight produces a larger decay term.

Check the signs. Both terms reduce the weight in the usual formulation, so the new weight is smaller than the old weight when the gradient is positive.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
