# Batch Normalization

**Batch norm:** For a mini-batch, normalize activations to $\hat{x}=(x-μ)/\sqrt{σ^{2}+ε}$, then scale/shift $y=γ\hat{x}+β$ with learnable $γ,β$.

## Stabilizing Training

### Mean Zero, Variance One
Per-batch mean $μ$ and variance $σ^{2}$ standardize inputs to each layer, allowing higher learning rates and reducing internal covariate shift.

### Test Time
At inference, use running averages of $μ,σ^{2}$ estimated during training, not batch statistics.

## Example

Layer input batch $[1,3,5]$: $μ=3$, $σ^{2}=8/3$, $\hat{x}=[-1.22,0,1.22]$, then $y=γ\hat{x}+β$ restores capacity while keeping distribution stable.
