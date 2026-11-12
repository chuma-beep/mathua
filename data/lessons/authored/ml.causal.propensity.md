# Propensity scores

The propensity score is the probability of treatment given covariates. Conditioning on it balances the covariates between treated and control groups.

## The score

A logistic model is the usual estimate. The score is a one-dimensional summary of many covariates.

## Inverse weighting

Treated units are weighted by \(1/e(x)\) and controls by \(1/(1-e(x))\). The weights reconstruct a population in which treatment is independent of the covariates.

## Matching

Nearest-neighbour matching pairs units with close scores. The balance check compares covariate means after matching.

To check your understanding, fit the model, weight the treated and find a neighbour.

The propensity score compresses confounding into one dimension.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
