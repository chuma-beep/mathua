# The SAGA estimator

SVRG needs a full pass to refresh. SAGA stores one gradient per sample and refreshes a single entry at a time, so it never needs a full gradient after the start.

## The gradient table

The method keeps a table with one stored gradient per sample. A step draws a sample, reads its stored gradient, and replaces that entry with the freshly computed gradient. The estimator combines the new gradient, the old entry and the table average.

## The update

The corrected direction is the new sample gradient minus the old stored entry plus the running table average. The average is updated incrementally, so the cost of a step does not grow with the sample count.

## Memory and bias

The table costs one gradient vector per sample, which is the price of avoiding full passes. The estimator is unbiased at every step, and the stored entries converge to the local gradients as the iterate settles.

To check your understanding, store one gradient per sample, subtract the old entry and cost one vector per sample.

SAGA trades memory proportional to the sample count for a step cost that stays constant, which suits problems where a full gradient is the expensive operation.
