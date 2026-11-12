# The multivariate normal

The multivariate normal is determined by a mean vector and a positive semidefinite covariance matrix. Every linear combination of its coordinates is normal.

## Density

When the covariance is invertible, the density is proportional to the exponential of minus half the squared Mahalanobis distance.

## Linear maps

An affine map sends a multivariate normal to another multivariate normal with transformed mean and covariance.

## Independence

Jointly normal coordinates are independent exactly when their covariance is zero, a property that fails for general distributions.

To check your understanding, invert the covariance, transform the mean and check the covariance.

The multivariate normal is closed under affine maps and conditioning.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
