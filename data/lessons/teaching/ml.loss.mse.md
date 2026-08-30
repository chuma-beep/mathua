## Definition of mean squared error

Mean squared error (MSE) measures the average squared distance between predictions and targets: $\frac{1}{n}\sum (y_i - t_i)^2$.

### Example

Prediction $2$, target $3$: MSE $= (2-3)^2 = 1$.

## Properties

MSE is non-negative and zero only when predictions equal targets. It penalizes large errors more than small ones.

### Example

Predictions $[0,1]$ vs targets $[1,1]$: MSE $= ((0-1)^2 + 0)/2 = 0.5$.

## Derivative and use in regression

The derivative of $\frac12 (y-t)^2$ with respect to $y$ is $y-t$. MSE is standard for regression tasks.

### Example

At $y=2, t=3$, the gradient is $-1$, pointing toward the target.
