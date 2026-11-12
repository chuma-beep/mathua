# Monte Carlo standard error

The standard error of a Monte Carlo average falls as the square root of the sample count. Four times the samples halve the error, whatever the dimension.

## The formula

For independent samples with variance \(\sigma^{2}\), the average has standard error \(\sigma/\sqrt{n}\). The estimate is unbiased.

## Rate

To gain one decimal digit of accuracy requires one hundred times as many samples. The rate does not degrade with dimension, which is the method's advantage.

## Estimation

The sample variance replaces the unknown \(\sigma^{2}\). The result is itself noisy, so confidence intervals use a t quantile for small n.

To check your understanding, divide by the root count, quadruple the samples and plug in the sample variance.

The square-root rate is both the guarantee and the limitation of Monte Carlo.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
