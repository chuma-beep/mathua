# Importance sampling

Not every sample carries the same information. Importance sampling draws samples from a chosen distribution and reweights them so the estimator stays unbiased.

## The weighted estimator

If a sample is drawn with probability q instead of the uniform p, its contribution is weighted by p over q. The weighted average is unbiased for the original expectation and can have much lower variance.

## Choosing the distribution

A good proposal puts more mass where the integrand is large. The ideal proposal is proportional to the absolute integrand, which is rarely available, so practical schemes use norms or gradient magnitudes as proxies.

## The variance trade

Mismatched weights inflate variance and can make a single sample dominate. Effective sample size measures the damage, and clipping or mixing with the uniform distribution keeps the estimator stable.

To check your understanding, weight by p over q, put mass on large terms and watch effective size.

Importance sampling is the general form of the reweighting idea. Variance reduction methods for gradients are importance sampling over the sample index with an implicit proposal.
