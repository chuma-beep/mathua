# Gaussian process regression

Gaussian process regression is Bayesian linear regression in a feature space fixed by a kernel. The posterior mean is a kernel-weighted combination of the observed targets.

## The posterior mean

The mean at a test point is \(k_*^T(K+\sigma^2 I)^{-1}y\). The vector \(k_*\) holds the kernel values between the test point and the training inputs.

## The posterior variance

The variance is the prior variance minus a quadratic form in the same kernel vector. It falls to the noise level at an observed input.

## The marginal likelihood

The log marginal likelihood has a fit term, a complexity term and a constant. Maximising it tunes the kernel hyperparameters.

To check your understanding, build the kernel vector, subtract the quadratic form and name the three terms.

Gaussian process regression is kernel ridge regression with a calibrated variance.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
