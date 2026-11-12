# Antithetic variates

Antithetic variates pair each sample with a negatively correlated partner. The average of the pair has lower variance when the response is monotone.

## Pairing

For a uniform input U the partner is 1-U. The pair average is used as one observation, so n pairs cost 2n evaluations.

## Variance

If the two responses have correlation \(\rho\), the pair average has variance \(\tfrac{\sigma^{2}}{2}(1+\rho)\). Negative correlation therefore helps.

## Limits

For a non-monotone response the correlation can be positive and the variance can rise. The method is cheap but not universally beneficial.

To check your understanding, reflect the uniform, use the correlation and check monotonicity.

Antithetic sampling is a free variance reduction only when the response is monotone.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
