# Control variate estimators

A control variate is a correlated quantity whose mean is known. Subtracting its centred value from the estimator removes part of the noise.

## The estimator

\(\hat\theta=X-\beta(Y-\mu_Y)\), where \(\mu_Y\) is known. The correction has mean zero, so the estimator stays unbiased for any beta.

## Optimal coefficient

The variance is minimised at \(\beta^{*}=\mathrm{Cov}(X,Y)/\mathrm{Var}(Y)\). The remaining variance is \(\sigma_X^{2}(1-\rho^{2})\).

## Many controls

Several control variates can be combined by a small least-squares fit. The fit costs little compared with the sampling.

To check your understanding, subtract the centred control, use the covariance ratio and fit several coefficients.

A control variate spends a known mean to buy a smaller variance.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
