# Rank-one sampling

Some objectives are sums of rank-one terms, as in matrix factorisation. Sampling a row or column updates one factor block and leaves the other fixed, which makes each step cheap.

## The term structure

Each observation contributes an outer product of a row feature and a column feature. A sampled term touches one row and one column, so the gradient of the step is sparse in the factor matrices.

## The step

The sampled gradient is applied to the two touched blocks. Unsampled rows keep their values, which is what makes the update cheap and also what introduces the sampling noise.

## The variance

Variance depends on how unevenly the observations cover the matrix. Uniform sampling is unbiased but slow on skewed data, and weighted sampling by row or column norms is the usual remedy.

To check your understanding, form an outer product, apply the sampled gradient and expect skew from coverage.

Rank-one sampling shows how structure in the objective, not only in the sampler, controls the cost and the noise of a stochastic method.
