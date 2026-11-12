# Consistency

An estimator is consistent when it converges in probability to the parameter as the sample grows.

## Definition

Consistency requires that for every positive tolerance, the probability of an error larger than that tolerance tends to zero.

## Sufficient conditions

Vanishing bias and vanishing variance together imply consistency by Chebyshev's inequality.

## Failure

An estimator can be unbiased at every sample size and still fail to converge if its variance does not shrink.

To check your understanding, fix a tolerance, send the bias to zero and keep the estimator unbiased.

Consistency is a large-sample statement about concentration.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
