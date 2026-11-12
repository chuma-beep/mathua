# Snapshots and epochs

SVRG organises its work in epochs. One epoch is one pass over the data, and its boundaries decide when the snapshot is refreshed.

## Epoch length

The inner loop runs for a fixed number of steps, often equal to the sample count. A shorter epoch refreshes more often and lowers bias, while a longer epoch amortises the full gradient over more cheap steps.

## The snapshot point

The snapshot is taken at the iterate where the epoch begins. If the iterate moves far during the epoch the correction weakens, which is why the step size is kept small relative to the strong convexity.

## Cost per epoch

Each epoch costs one full gradient plus one sample gradient per inner step. The full pass is a small fraction of the epoch when the sample count is large, so the method stays close to the sample cost.

To check your understanding, fix the inner steps, take it at the epoch start and pay one full gradient.

Treating the snapshot as an epoch boundary makes the memory and the cost accounting explicit, which is how SVRG is implemented and how its learning rate is tuned.
