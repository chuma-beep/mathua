# Girsanov's theorem

Girsanov's theorem changes the drift of a Brownian motion by reweighting the probability measure with an exponential martingale.

## The density

The Radon-Nikodym density is the exponential of a stochastic integral minus half its quadratic variation.

## Effect on drift

Under the new measure, \(B_t-\int\theta_s ds\) is a Brownian motion, so the drift is removed.

## Condition

The density is a true martingale when Novikov's condition holds, which bounds the exponential moment of the squared drift.

To check your understanding, build the exponential martingale, shift the process and state Novikov.

Girsanov moves drift into the measure.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
