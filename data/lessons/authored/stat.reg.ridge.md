# Ridge regression

Ridge regression fits a linear model with a penalty on the size of the coefficients, trading a little bias for a large reduction in variance.

## The ridge penalty

Ordinary least squares minimises the residual sum of squares. Ridge minimises that sum plus \(\lambda \sum \beta_j^2\), so large coefficients cost something and the fit prefers modest ones. The intercept is normally left unpenalised. With \(\lambda = 0\) the fit is least squares, while a large \(\lambda\) drives every slope toward zero and the model approaches the mean of the response.

## The closed form

The solution is \(\hat{\beta} = (X^\top X + \lambda I)^{-1} X^\top y\), with I the identity on the penalised coefficients. The added diagonal makes the matrix invertible even when the predictors outnumber the observations, which is what made ridge useful for wide data. If the ordinary estimate is 2.5 and the penalty shrinks it to 1.9, the bias added is 0.6 while the variance saved is often far larger.

## Choosing lambda

The penalty is chosen by cross-validation, usually over a grid such as 0.01, 0.1, 1 and 10; the value with the lowest held-out error wins. The predictors are normally standardised first, because the penalty treats a coefficient of 0.2 on a variable measured in thousands and a coefficient of 2 on a variable measured in units as unequal purely because of scale.

Ridge never sets a coefficient exactly to zero, so it is a smoother rather than a selector. That distinction decides when it is the right tool.
