# Multiple regression

Multiple regression predicts one response from several explanatory variables at once. Each predictor has its own coefficient, and the coefficient describes the change in the predicted response when that predictor rises by one unit while the others are held fixed. The model is still linear in the coefficients, even though it can use several variables.

## The fitted value

The fitted model is

\[
\hat{y} = b_0 + b_1 x_1 + b_2 x_2 + \cdots + b_k x_k.
\]

The intercept \(b_0\) is the predicted response when every predictor is zero. That value is sometimes meaningful and sometimes outside the data range, so it should be read with care. Each slope \(b_j\) is the partial effect of \(x_j\), holding the other predictors constant.

Substitution is the basic operation. For a model \(\hat{y} = 12 + 3x_1 - 2x_2\), the prediction at \(x_1 = 4\) and \(x_2 = 5\) is \(12 + 12 - 10 = 14\). The arithmetic is ordinary, but the interpretation depends on the predictors being measured on their own scales.

The sign of a coefficient is not the same as the sign of a simple correlation. A predictor can have a positive slope in the multiple model and a negative raw correlation, because the other predictors absorb part of its variation. That reversal is one reason multiple regression is not just a stack of separate simple regressions.

## Residuals

The residual for an observation is

\[
e = y - \hat{y}.
\]

It is the part of the response the model did not explain. A positive residual means the observed value is above the prediction, and a negative residual means it is below. The residuals are the raw material for every regression diagnostic, from the residual plot to the standard error of the estimate.

Residuals should be plotted against the fitted values and against each predictor. A pattern in the plot suggests that the linear form is missing something, such as curvature or an interaction. A fan shape suggests that the variance changes with the fitted value, which violates the usual constant-variance assumption.

The sum of the residuals is zero when the model includes an intercept and is fitted by least squares. That identity is a consequence of the normal equations, not an extra assumption. It also means the residuals are not independent, because one is determined by the others.

## Adjusted R-squared

Adding predictors to a regression can only increase \(R^2\), even when the new predictors are noise. The adjusted version penalizes the addition:

\[
R^2_{\text{adj}} = 1 - (1-R^2)\frac{n-1}{n-k-1},
\]

where \(n\) is the number of observations and \(k\) is the number of predictors. The fraction is larger than one when \(k > 0\), so the adjusted value is lower than the ordinary \(R^2\).

The penalty grows as \(k\) approaches \(n-1\), and the formula is undefined when \(n-k-1 = 0\). A model with as many predictors as observations can fit the data perfectly and says nothing about a new observation, which is why the adjusted measure refuses that case.

Adjusted \(R^2\) is useful for comparing models with different numbers of predictors on the same data. It is not a measure of causal importance, and it does not replace a residual analysis. A high adjusted value with a patterned residual plot is still a model that is missing structure.
