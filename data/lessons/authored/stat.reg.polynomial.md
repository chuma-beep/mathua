# Polynomial regression models

A polynomial regression fits a response as a polynomial in one predictor, adding powers of that predictor as separate terms so the fitted curve can bend.

## Curvature in the model

A straight line assumes the change in the response per unit of the predictor is constant. A squared term lets that change itself change. Fitting \(y = \beta_0 + \beta_1 x + \beta_2 x^2\) to points whose rise slows with x gives a negative \(\beta_2\), so the curve flattens. The model remains linear in the coefficients, so least squares applies unchanged.

## Coefficients and predictors

The powers \(x\) and \(x^2\) are separate predictors whose coefficient standard errors are correlated. A t test on \(\beta_2\) therefore answers whether curvature is needed rather than whether x matters. A quadratic fitted to 20 points with \(\beta_2 = -0.5\) and standard error 0.15 gives a t statistic of about minus 3.33. Centring x by subtracting its mean usually reduces that correlation.

## Interactions

A polynomial in one variable and an interaction between two variables are different devices that are often confused. The term \(x_1 x_2\) says the effect of one predictor depends on the level of another, while \(x^2\) says the effect of one predictor depends on its own level. A model with 3 predictors and all pairwise interactions has 3 interaction terms, whereas a cubic in one predictor has 2 power terms.

The choice between the two is a claim about mechanism, so it should be made from the subject matter before the fit is inspected.
