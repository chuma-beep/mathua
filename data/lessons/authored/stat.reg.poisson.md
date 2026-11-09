# Poisson regression

Poisson regression models a count with a log link. The mean count is the exponential of the linear predictor.

## Apply the log link

The model says the log of the mean count is a linear function of the predictors. The mean is therefore the exponential of that linear function. The link keeps the mean positive.

## Compute a rate

Exponentiate the linear predictor at the chosen covariate values. The result is the predicted rate. For an intercept-only model it is the exponential of the intercept.

## Interpret a coefficient

A one-unit increase in a predictor multiplies the mean by the exponential of its coefficient. A coefficient of zero means no multiplicative effect. The interpretation is multiplicative, not additive.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
