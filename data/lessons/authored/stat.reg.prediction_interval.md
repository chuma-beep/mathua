# Prediction intervals

A prediction interval estimates where a single new observation will fall, while a confidence interval estimates where the mean response lies. The prediction interval is always wider, because a new observation carries its own random error on top of the uncertainty in the fitted mean. Confusing the two is one of the most common mistakes in applied regression.

## A new observation is harder to predict

The fitted line gives an estimate of the mean response at a chosen predictor value. A new observation at that value differs from the fitted mean for two reasons: the mean itself is estimated and the observation has its own noise around the mean. The variance of the prediction error therefore has two parts:

\[
\operatorname{Var}(y_{\text{new}} - \hat{y}) = \sigma^2 + \operatorname{Var}(\hat{y}).
\]

The first term is the noise of the new observation. The second is the uncertainty of the fitted mean. A confidence interval for the mean uses only the second term, so it is narrower.

The extra term does not shrink to zero as the sample grows. Even with a perfect estimate of the mean, a new observation still varies around that mean. That is why a prediction interval has a floor, while a confidence interval can become arbitrarily narrow with enough data.

## The half-width

The prediction interval has the form

\[
\hat{y} \pm t^* s_{\text{pred}},
\]

where \(s_{\text{pred}}\) is the standard error of prediction and \(t^*\) is the critical value with the regression degrees of freedom. The half-width is the product of the two, and it is the quantity to report when the interval is described.

For a prediction with \(t^* = 2.306\) and \(s_{\text{pred}} = 3.4\), the half-width is \(2.306 \times 3.4 \approx 7.84\). If the fitted value is 50, the interval runs from about 42.16 to 57.84. The upper and lower endpoints are symmetric around the fitted value.

The critical value is larger for small samples and for higher confidence levels. A 95 percent interval uses a smaller multiplier than a 99 percent interval, and a regression with few degrees of freedom uses a larger t value than one with many.

## Combining the two standard errors

The standard error of prediction is

\[
s_{\text{pred}} = \sqrt{s_e^2 + s_{\text{mean}}^2},
\]

where \(s_e\) is the residual standard error and \(s_{\text{mean}}\) is the standard error of the fitted mean at the chosen predictor value. The two contributions add in quadrature, not directly, because they are independent sources of variation.

If the residual standard error is 3 and the standard error of the mean is 4, the prediction standard error is \(\sqrt{9+16} = 5\). The Pythagorean numbers make the combination exact and show that the prediction error is always at least as large as either component.

The standard error of the mean is smallest near the centre of the predictor data and grows toward the edges. A prediction interval is therefore narrowest near the mean of the predictors, and it flares at the extremes. That flare is a warning that predictions far outside the observed range are less trustworthy.
