# Stochastic averaging

Averaging the iterates of a stochastic method reduces the variance of the final estimate without changing the steps.

## The average

The returned point is the mean of the iterates rather than the last iterate. The mean cancels part of the noise that the last point carries.

## Variance

Averaging many roughly independent iterates divides their variance by the count. The benefit is largest when the steps are small and the iterates are close.

## Weighting

A uniform average treats every iterate equally. A weighted average can favour later iterates when the early ones are far from the solution.

To check your understanding, average two iterates, divide the variance and choose weights.

Averaging is the cheapest variance reduction available once the iterates exist.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
