# The lasso

The lasso uses an absolute-value penalty instead of a quadratic one. The penalty drives some coefficients exactly to zero.

## The objective

The objective is the squared residual plus a positive multiple of the sum of absolute coefficients. The absolute value is not differentiable at zero, which is where sparsity comes from.

## Sparsity

The solution sets coefficients to zero when the data does not justify them. The number of nonzero coefficients is not fixed in advance, so the method selects variables.

## A closed form

In one dimension the solution is a soft threshold of the unpenalised estimate. Values inside the threshold become zero.

To check your understanding, evaluate the penalty, count nonzeros and apply a soft threshold.

The lasso is least squares with a penalty that can afford to discard variables.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
