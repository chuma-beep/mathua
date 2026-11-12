# Control variates

A control variate is a random quantity with a known mean that is correlated with the estimator. Subtracting it and adding its mean reduces variance without changing the expectation.

## The estimator

For a sample gradient G and a control variate H with known mean, the corrected estimator is G - H + mean(H). The correction is unbiased for any H, and it is effective when H tracks the part of G that fluctuates.

## Choosing the coefficient

A scalar coefficient scales the control variate. The variance-minimising coefficient is the covariance of G and H divided by the variance of H. In practice the coefficient is estimated or fixed from a pilot run.

## The variance reduction

The remaining variance is the original variance times one minus the squared correlation. A control variate that correlates strongly with the gradient can remove most of the noise at almost no cost.

To check your understanding, subtract the variate, use the covariance ratio and square the correlation.

Control variates explain why snapshot gradients work. The snapshot is the control variate, its mean is the full gradient, and the correlation comes from evaluating the same sample at two nearby points.
