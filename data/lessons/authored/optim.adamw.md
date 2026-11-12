# AdamW and decoupled weight decay

Adding weight decay to the gradient couples it to the adaptive denominator, so the effective decay depends on the gradient history. AdamW applies the decay directly to the parameters instead.

## Coupled decay

In the coupled form the decay term enters the gradient before the adaptive scaling. A coordinate with a large gradient history then receives a decay that is scaled down along with its gradient.

## Decoupled decay

The decoupled form shrinks each parameter by a factor before the adaptive step is added. The decay is the same for every coordinate and independent of the gradient statistics, which is the intended regularisation.

## The effect

The two forms agree only when the adaptive denominator is constant. Elsewhere they regularise differently, and the decoupled form usually gives the stronger and more predictable shrinkage.

To check your understanding, add decay to the gradient, shrink the parameters and agree only when constant.

The distinction matters because regularisation is a modelling choice, and coupling it to the optimiser's gradient history makes the model depend on an implementation detail.
