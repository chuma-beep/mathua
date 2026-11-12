# Gaussian conditioning

Conditioning a multivariate normal on some coordinates leaves the others normal, with a mean that is linear in the observed values.

## Formula

The conditional mean is the unconditional mean plus the cross-covariance times the inverse covariance of the observed block times the centred observation.

## Variance

The conditional covariance is the Schur complement of the observed block, which does not depend on the observed values.

## Use

The formulas are the basis of the Kalman filter and of Gaussian process regression.

To check your understanding, split the covariance, form the Schur complement and update a state estimate.

Gaussian conditioning is linear updating with a constant variance.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
