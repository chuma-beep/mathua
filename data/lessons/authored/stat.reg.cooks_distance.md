# Cook's distance

Cook's distance measures how much the fitted values of a regression change when one observation is removed, so it ranks observations by influence.

## Influence of one point

An observation is influential when deleting it moves the fitted line noticeably, which is not the same as being an outlier. A point with a large residual but little leverage moves nothing, while a point with a modest residual at the edge of the predictor range can move everything. Cook's distance combines the two, weighting the residual by the leverage \(h_i\) through \(D_i = r_i^2 h_i/(p(1-h_i)^2)\), with p the number of parameters.

## The 4/n cutoff

A common rule flags observations with \(D_i\) above \(4/n\). For 50 observations the threshold is \(4/50 = 0.08\), a small number that is easy to exceed and is meant to prompt inspection rather than deletion. The threshold is a screening device rather than a test, so it says nothing about whether the observation is wrong.

## A worked check

With 40 observations the cutoff is \(4/40 = 0.1\). Suppose one point has leverage 0.2, a standardised residual of 2 and the model has 2 parameters. Then \(D = 4 \times 0.2/(2 \times 0.64) = 0.8/1.28 = 0.625\), more than six times the cutoff, so the point is worth examining. Dropping it and refitting shows the change in the slope, which is the number a reader needs.

A large distance is a question about the data rather than a verdict, so the answer belongs in the report beside the fitted model.
