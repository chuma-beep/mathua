# Logistic regression for binary outcomes

Logistic regression for binary outcomes is part of the statistics curriculum. This lesson develops the idea in three sections, from the defining picture through the main computation to the interpretation.

## From Probability to Log-Odds

A binary outcome takes one of two values, so a straight regression line is not suitable: it can predict probabilities below zero or above one. Logistic regression instead models the log-odds of success as a linear function of the predictors. If p is the probability of success, the odds are p/(1-p) and the log-odds are ln(p/(1-p)). The model sets this log-odds equal to b0 + b1 x1 + and further terms for additional predictors.

The inverse transformation turns a linear score z into a probability through the sigmoid function 1/(1+e^(-z)). A score of zero gives probability one half, large positive scores approach one and large negative scores approach zero. This keeps every fitted probability inside the valid interval while preserving a linear structure on the log-odds scale.

The odds ratio is the exponential of a coefficient. If b1 equals 0.7, then increasing x1 by one unit multiplies the odds by about two, holding the other predictors fixed. That interpretation is the main reason logistic regression remains widely used: the coefficients describe changes in odds rather than changes in the probability itself, and the latter depend on the starting probability.

## Fitting the Model

Logistic regression is usually fitted by maximum likelihood rather than least squares. The likelihood of the observed outcomes is the product of the predicted probabilities for successes and one minus those probabilities for failures. Maximising the log-likelihood gives a concave problem with a unique solution when the data are not separated, and numerical methods such as Newton or gradient descent find the coefficients.

The fitted model produces a probability for every case, and a threshold converts that probability into a predicted class. The default threshold of one half is convenient but not mandatory. In a medical screening problem the cost of a missed case may justify a lower threshold, while a spam filter may prefer a higher one. The threshold is a decision rule, not a parameter estimated by the regression.

Model assessment uses the same ideas as other classifiers: confusion matrices, sensitivity, specificity and the area under the ROC curve. Calibration is also important, because a model can rank cases well while its probabilities are systematically too high or too low. A reliability plot compares predicted probabilities with observed frequencies and reveals that kind of error.

## Interpreting Coefficients

Each coefficient describes the change in log-odds for a one-unit increase in its predictor, holding the others constant. Exponentiation converts that change into an odds ratio. A coefficient of zero means no association on the log-odds scale, a positive coefficient means higher odds of success and a negative coefficient means lower odds. The sign is usually more robust to interpret than the exact size.

Because the model is nonlinear in the probability, a one-unit change does not produce a constant change in probability. The effect is largest near a probability of one half and smaller near the extremes. Reporting an average marginal effect, or the effect at a representative value, is often clearer than quoting a single probability change for all cases.

Correlated predictors complicate interpretation in the same way they do in ordinary regression. Coefficients can change substantially when another predictor is added or removed, and a large standard error may signal collinearity rather than a weak relationship. Standardising predictors or examining the correlation structure before fitting helps avoid misleading conclusions about which variables matter.
