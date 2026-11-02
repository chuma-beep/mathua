# Monte Carlo Integration

## Random Samples

Monte Carlo integration estimates an average value by sampling points at random. The integral is approximated by the sample mean multiplied by the volume of the region.

The method does not require a grid or a smooth integrand. It only needs the ability to evaluate the function at arbitrary points and to draw uniform samples.

The estimate is unbiased, so its expected value equals the true integral for any number of samples. The error is measured by the standard deviation of the estimator.

## Variance and Error

The standard error of the estimate is the standard deviation of the integrand divided by the square root of the sample count. Quadrupling the samples halves the error.

This rate is slower than the high-order convergence of Gaussian quadrature in one dimension. In high dimensions, however, the cost of a grid grows exponentially while Monte Carlo stays polynomial.

The error estimate is itself random. A confidence interval is formed from the sample variance, which gives a practical way to report the accuracy of the result.

## Variance Reduction

Importance sampling draws points from a distribution that concentrates where the integrand is large. The samples are reweighted so the estimate remains unbiased.

Antithetic sampling pairs each draw with its mirror image. If the integrand is monotone, the pair averages away part of the fluctuation.

Control variates subtract a function with a known integral and add its exact value back. A well-chosen control can reduce the variance by a large factor at almost no extra cost.
