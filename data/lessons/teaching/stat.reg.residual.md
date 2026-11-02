# Residuals and prediction error

A residual measures how far an observed value lies from the value predicted by a fitted model. Residuals show what the line missed.

## Observed and predicted values

For an observed response $y_i$ and a fitted value $\hat{y}_i$, the residual is

$$e_i = y_i - \hat{y}_i.$$

A positive residual means the observation lies above the line. A negative residual means it lies below. The fitted value comes from substituting the explanatory value into the regression equation.

## The residual sum of squares

Squaring removes the sign and penalises large misses. The residual sum of squares is

$$SSE = \sum_{i=1}^{n} e_i^2.$$

The least-squares line is the line with the smallest possible SSE for the given data. A smaller SSE means the line passes closer to the observations.

## Reading a residual plot

Plot the residuals against the explanatory variable. A random cloud around zero supports a straight-line fit. A curved pattern suggests a nonlinear relationship.

A funnel shape means the spread changes with the explanatory variable. A single point far from the others may be influential and worth checking.
