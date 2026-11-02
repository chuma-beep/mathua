# The bias-variance tradeoff

## Decomposing the error

The expected squared prediction error of a learned model splits into three parts: the squared bias of the model family, the variance of the fitted model around its own average, and the irreducible noise in the target.

Bias measures how far the average prediction lies from the truth. Variance measures how much the prediction moves when the training sample changes.

## Model complexity

Simple models tend to have high bias and low variance, while flexible models tend to have low bias and high variance. Moving along that spectrum trades one component against the other.

The best expected error occurs where the sum of squared bias and variance is smallest, which is not necessarily at the most flexible model.

## Averaging reduces variance

If n independent predictions each have variance v, their average has variance v divided by n. Ensembles exploit this effect by combining many models whose errors are not perfectly correlated.

Regularisation also reduces variance by constraining the fitted parameters, at the cost of introducing a small amount of bias.
