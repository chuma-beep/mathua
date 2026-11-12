# Learning-rate warmup

Adaptive methods build their statistics from the first few gradients, which can be noisy. Warmup starts with a small rate and increases it so those statistics settle before full steps are taken.

## Linear warmup

The rate rises linearly from near zero to its target over a fixed number of steps. The ramp is short relative to training and prevents a large early step from a poorly scaled denominator.

## Why it helps

Early second moments are estimated from few samples and can be far too small, which produces large steps. Warmup keeps the product of rate and step moderate until the estimates improve.

## The cost

Warmup spends steps at a reduced rate. Too long a ramp wastes budget, and too short a ramp does not protect the statistics, so the length is usually a small fraction of the run.

To check your understanding, ramp from zero, settle early statistics and spend reduced steps.

Warmup is a safety measure for adaptive scaling rather than a convergence device, and it is most valuable when the initial gradients are large or sparse.
