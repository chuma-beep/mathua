# The Wishart distribution

The Wishart distribution is the law of a sum of outer products of independent Gaussian vectors. It generalises the chi-squared distribution to matrices.

## Definition

For \(n\) independent centred Gaussian vectors with covariance \(\Sigma\), the sum of their outer products is Wishart with \(n\) degrees of freedom.

## Mean

The expectation is \(n\Sigma\), so the normalised matrix is an unbiased estimator of the covariance.

## Special case

In one dimension the Wishart law reduces to a scaled chi-squared law.

To check your understanding, draw Gaussian vectors, multiply the covariance and take one dimension.

The Wishart is the sample covariance matrix's law.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
