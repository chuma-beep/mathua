# Meta-gradients

A meta-gradient is the derivative of a training outcome with respect to a parameter of the training procedure. It answers how the final model would change if an initialisation or a step size changed.

## The two levels

The inner level runs a training algorithm for a fixed number of steps. The outer level measures the resulting model on a validation set and asks for the derivative with respect to the procedure's parameters.

## The backward pass

The outer derivative is computed by differentiating through the inner steps, either by unrolling or by the implicit function theorem at the inner solution. Unrolling is exact for a truncated run and costly in memory.

## The use

Meta-gradients tune learning rates, initialisations and loss weights, and they underlie few-shot adaptation. Their practical limit is variance, because a derivative through many stochastic steps is noisy.

To check your understanding, run inner training, unroll or use implicit form and tune learning rates.

A meta-gradient is a gradient through a learning process, and its cost is the memory or the approximation used to cross the inner loop.
