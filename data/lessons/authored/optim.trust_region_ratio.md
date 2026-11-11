# The trust-region ratio

A trust-region method builds a quadratic model, restricts the step to a ball and then compares the predicted decrease with the actual decrease. That ratio decides whether the step is accepted and whether the radius grows, shrinks or stays. The ratio is the feedback signal that keeps the model honest.

## Computing the ratio

The numerator is the actual decrease in the objective from the current point to the trial point. The denominator is the decrease predicted by the quadratic model over the same step. The ratio is the first divided by the second. A ratio near one means the model is accurate, a ratio near zero or negative means it is not.

## Accepting or rejecting a step

A common rule accepts the step when the ratio exceeds a small threshold and rejects it otherwise. A rejected step leaves the point unchanged and reduces the radius. Acceptance does not require a ratio of one, only enough agreement that progress is real.

## Updating the radius

A high ratio expands the radius, a moderate ratio keeps it and a low ratio contracts it. The thresholds and factors are parameters, and the logic is the same in every implementation. The sequence of radii records how much the model has earned the right to be trusted.

The ratio turns a model into a controlled step. It is why trust-region methods are robust when curvature changes quickly and why their parameters have a clear meaning.
