# Dual averaging

Dual averaging accumulates gradients in a dual space and returns a primal point from the running sum. It is the lazy counterpart of mirror descent, which applies each gradient immediately.

## The running average

The method adds each gradient to a dual accumulator and then maps the accumulator back through the divergence. The primal iterate is therefore a function of the whole history rather than of the last step.

## Lazy versus greedy

Mirror descent is greedy: it moves as soon as a gradient arrives. Dual averaging is lazy: it stays near the current average and moves only when the accumulated evidence shifts, which gives smoother iterates and simpler regularisation.

## The regulariser view

The accumulator can be read as a penalty that pulls the iterate toward the origin of the dual space. Changing the penalty changes the method, and the proximal form of the penalty recovers several known algorithms.

To check your understanding, accumulate gradients in the dual, contrast with mirror descent and read the accumulator as a penalty.

Dual averaging trades immediate responsiveness for a stable average, which is why it suits problems where the gradient is noisy and the regulariser is part of the objective.
