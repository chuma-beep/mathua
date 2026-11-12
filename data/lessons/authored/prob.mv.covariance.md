# Covariance matrices

A covariance matrix collects the variances on its diagonal and the covariances off it. It is symmetric and positive semidefinite.

## Properties

For any vector \(a\), the variance of \(a^TX\) is \(a^T\Sigma a\), which is non-negative, so \(\Sigma\) is positive semidefinite.

## Decomposition

An eigendecomposition writes the matrix as a rotation, a diagonal scaling and the inverse rotation.

## Rank

The rank equals the dimension of the support of the centred random vector, so a singular matrix signals a linear constraint.

To check your understanding, take a linear combination, diagonalise and count positive eigenvalues.

Covariance matrices encode variances and linear dependence together.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
