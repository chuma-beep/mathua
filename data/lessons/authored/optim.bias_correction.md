# Bias correction in adaptive methods

Adaptive averages start at zero and are therefore biased towards zero early in training. Bias correction divides by one minus the decay raised to the step count.

## The size of the bias

After k steps a decayed average has weight one minus the decay to the k on the current sample. The missing weight is the bias, and it is largest in the first few steps.

## The correction

Dividing by one minus the decay to the k restores the scale of an average over the samples seen. Both the first and second moments are corrected in the same way.

## The effect on steps

Without correction the first steps are too small and the denominator is too small. Correcting both moments makes the early steps well scaled, which matters most when training is short.

To check your understanding, compute the missing weight, divide by one minus decay to k and avoid tiny early steps.

Bias correction is a small factor with a large effect on the first hundred steps, which is exactly where a badly scaled update can waste a run.
