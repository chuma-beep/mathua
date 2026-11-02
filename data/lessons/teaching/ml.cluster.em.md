# Expectation-maximization

## Latent assignments

Expectation-maximization fits models with unobserved variables, such as the component that generated each point in a mixture. The expectation step computes the posterior probability of each assignment given the current parameters.

For two equally weighted components with densities p and q, the responsibility of the first component is p divided by p plus q. The responsibilities of all components for one point sum to one.

## The maximization step

The maximization step updates each component using the responsibilities as soft weights. Means become weighted averages, and covariances become weighted outer products of deviations.

Each full iteration cannot decrease the observed-data likelihood. The guarantee is monotone improvement rather than a fixed number of iterations.

## Counting parameters

A Gaussian mixture with k components in d dimensions and full covariances stores k means, k covariance matrices and k minus one free mixture weights. Each full covariance contributes d times d plus one divided by two entries.

The covariance part alone therefore has k times d times d plus one divided by two parameters, which grows quadratically in the dimension.
