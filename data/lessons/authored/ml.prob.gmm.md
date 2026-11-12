# Gaussian mixture models

A Gaussian mixture writes the density as a weighted sum of Gaussian components. Each component owns a mean and a covariance, and the weights form a probability vector.

## Parameter count

A mixture of \(K\) full-covariance components in \(d\) dimensions has \(K(2d+1)\) parameters before the weights. The weights add \(K-1\) free values.

## Responsibilities

The responsibility of a component is its posterior probability given a point. Responsibilities over all components sum to one for each point.

## The EM update

The M step sets each mean to the responsibility-weighted average. The covariance uses the same weights around that mean.

To check your understanding, count means and covariances, compute a responsibility and weight the mean.

A mixture models a multimodal density with local Gaussian pieces.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
