# Huber loss

## Below the threshold

For a residual within the threshold delta, Huber loss is one half of the squared residual. This is the same as squared error near zero, so small errors are smoothed.

The threshold controls the switch. A larger delta keeps more residuals in the quadratic region.

## Above the threshold

For a residual beyond the threshold, Huber loss becomes linear. The loss is delta times the absolute residual minus half of delta. The linear growth reduces the influence of large outliers.

The two pieces meet with matching values and slopes at the threshold. This makes the loss differentiable there.

## Comparing losses

Squared error grows quadratically and is sensitive to outliers. Absolute error grows linearly and is robust but less smooth near zero. Huber loss combines the quadratic behaviour near zero with the linear behaviour far away.

The difference between squared error and Huber loss is largest for large residuals. For small residuals the two losses are identical.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
