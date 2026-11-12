# Sample complexity

Sample complexity counts the samples needed to reach a target accuracy. It separates the statistical cost from the computational one.

## The convex rate

For a strongly convex objective with bounded variance the sample cost to reach an accuracy is proportional to one over the accuracy. The constant holds the condition number and the variance.

## The batch trade

A larger batch reduces variance but costs more per step. The total sample count is steps times batch size, and the best split depends on the regime.

## The nonconvex rate

Without convexity the guarantee is stationarity and the sample cost scales as one over the accuracy squared. That gap is the price of the weaker target.

To check your understanding, read the inverse accuracy, multiply steps by batch and square the accuracy.

Sample complexity is the currency of stochastic optimisation, and the target decides the exchange rate.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
