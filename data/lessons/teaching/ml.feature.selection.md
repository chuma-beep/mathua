# Feature selection

## Why select features

Feature selection removes inputs that add noise, cost or redundancy. A smaller feature set can improve generalization, reduce measurement cost and make a model easier to explain.

Filter methods rank features by a statistic computed without a model. Wrapper methods repeatedly train a model on candidate subsets. Embedded methods select features as part of fitting, as the lasso does through its penalty.

## Searching subsets

The number of subsets of d features is 2 to the power d, so exhaustive search becomes impractical quickly. Forward selection adds one feature at a time, and backward elimination removes one at a time.

The number of distinct pairs among d features is d times d minus one divided by two, which is the cost of building a full pairwise correlation ranking.

## Collinearity

The variance inflation factor measures how much the variance of an estimated coefficient grows because a feature is linearly related to the others. It equals one divided by one minus the squared multiple correlation of that feature against the rest.

Values near one indicate little collinearity, while large values signal that the coefficient is poorly determined and that the feature may be redundant.
