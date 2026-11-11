# The Armijo condition

The Armijo condition accepts a step only when the objective decreases by at least a fixed fraction of what the first-order model predicts. It prevents a line search from accepting a step that is technically downhill but barely so, and it is the sufficient-decrease half of the Wolfe conditions.

## Testing sufficient decrease

Compute the predicted decrease as the step size times the inner product of the gradient with the direction. The condition compares the actual decrease with a fraction, usually small, of that prediction. If the actual decrease is at least the fraction, the step passes. The test needs two evaluations per trial: the current point and the trial point.

## Choosing a backtracking step

Start with a trial size and shrink it by a fixed factor until the condition holds. Each shrink costs one evaluation. The number of shrinks is usually small near a good direction and larger when the direction is poor. The final accepted size is what the method uses for the update.

## Comparing two trial steps

Two trial sizes can be compared by their predicted and actual decreases. The larger size that still satisfies the condition is generally preferred, because it makes more progress. A size that fails the condition is rejected regardless of how large its predicted decrease was.

Sufficient decrease is the cheapest safeguard against a useless step. It is one inequality, it needs no second derivative, and it is the part of a line search that guarantees the objective sequence decreases.
